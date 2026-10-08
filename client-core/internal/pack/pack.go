// Package pack 负责文件的流式打包与解包。
//
// 全程流式：边遍历边压缩边写出，不在内存或磁盘上留完整副本。
// 复制一个 2GB 的文件夹时，这决定了是立刻开始传输还是先卡住几分钟。
package pack

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// Stats 描述一次打包的规模，用于在传输前告知对端与界面。
type Stats struct {
	Files int
	Dirs  int
	Bytes int64
}

// Scan 预扫描待传内容的规模。
//
// 用于在发送前给出总量，让进度条有分母。大目录扫描本身要耗时，
// 因此限制条目数上限，超出即停止并返回已统计的部分。
func Scan(paths []string, maxEntries int) (Stats, error) {
	var s Stats
	for _, root := range paths {
		err := filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // 个别文件读不到不该中断整次扫描
			}
			if info.IsDir() {
				s.Dirs++
			} else {
				s.Files++
				s.Bytes += info.Size()
			}
			if maxEntries > 0 && s.Files+s.Dirs >= maxEntries {
				return errTooMany
			}
			return nil
		})
		if err == errTooMany {
			return s, nil // 只是规模超限，不是错误
		}
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

var errTooMany = fmt.Errorf("条目数超过上限")

// Pack 把若干路径打包成 tar+zstd 流写入 w。
//
// 目录会被递归打包；多个路径各自作为顶层条目，
// 因此对端解包后能还原出与复制时相同的结构。
func Pack(w io.Writer, paths []string) error {
	enc, err := zstd.NewWriter(w,
		// 同步文件传输更看重速度而非极致压缩比——
		// 多数待传内容（图片、压缩包、二进制）本来就压不动。
		zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		return fmt.Errorf("创建压缩流: %w", err)
	}
	defer enc.Close()

	tw := tar.NewWriter(enc)
	defer tw.Close()

	for _, p := range paths {
		if err := addPath(tw, p); err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("关闭 tar 流: %w", err)
	}
	return enc.Close()
}

func addPath(tw *tar.Writer, root string) error {
	info, err := os.Lstat(root)
	if err != nil {
		return fmt.Errorf("读取 %s: %w", root, err)
	}
	// 以被复制项自身为顶层名，不带上它在源机器上的完整路径
	base := filepath.Base(root)

	if !info.IsDir() {
		return addFile(tw, root, base, info)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过读不了的条目，尽量把能传的传过去
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if skipWhenPacking(info.Name()) {
			return nil
		}
		name := base
		if rel != "." {
			name = filepath.Join(base, rel)
		}
		return addFile(tw, path, name, info)
	})
}

func addFile(tw *tar.Writer, path, name string, info os.FileInfo) error {
	// 符号链接按链接本身打包，不跟随——跟随可能把整个外部目录卷进来
	var link string
	if info.Mode()&os.ModeSymlink != 0 {
		var err error
		if link, err = os.Readlink(path); err != nil {
			return nil
		}
	}

	// 普通文件先打开再写头：头里已经声明了大小，写完头才发现打不开的话，
	// 整个 tar 流就坏了，对端一个文件也收不到。Windows 上被别的程序锁住的文件很常见
	var f *os.File
	if info.Mode().IsRegular() {
		var err error
		if f, err = os.Open(path); err != nil {
			return nil // 打开失败就跳过，不因单个文件中断整次传输
		}
		defer f.Close()
	}

	hdr, err := tar.FileInfoHeader(info, link)
	if err != nil {
		return fmt.Errorf("构造 tar 头 %s: %w", path, err)
	}
	// tar 内统一用正斜杠，保证 Windows 与 macOS 之间可互解
	hdr.Name = filepath.ToSlash(name)
	if info.IsDir() {
		hdr.Name += "/"
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("写 tar 头 %s: %w", name, err)
	}
	if f == nil {
		return nil // 目录与符号链接没有内容体
	}
	if _, err := io.Copy(tw, f); err != nil {
		return fmt.Errorf("写入 %s 内容: %w", name, err)
	}
	return nil
}

// Unpack 把 tar+zstd 流解到 destDir，返回顶层条目的绝对路径。
//
// 返回顶层条目是因为写入剪贴板时需要的正是这些路径——
// 用户复制的是"这几个文件"，粘贴出来也该是这几个。
func Unpack(r io.Reader, destDir string) ([]string, error) {
	if err := os.MkdirAll(destDir, 0o700); err != nil {
		return nil, err
	}

	dec, err := zstd.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("创建解压流: %w", err)
	}
	defer dec.Close()

	tr := tar.NewReader(dec)
	seen := make(map[string]struct{})
	var tops []string
	names := newDestNames()

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return tops, fmt.Errorf("读取 tar 条目: %w", err)
		}

		// Windows 上先把名字改成合法的（见 names.go），其他平台原样使用
		name := hdr.Name
		if windowsNames {
			if skipOnWindows(path.Base(strings.TrimSuffix(name, "/"))) {
				continue
			}
			if name, err = names.rel(hdr.Name, hdr.Typeflag == tar.TypeDir); err != nil {
				return tops, err
			}
		}

		target, err := safeJoin(destDir, name)
		if err != nil {
			return tops, err
		}

		// 记录顶层条目，供写入剪贴板使用
		if top := topLevel(name); top != "" {
			if _, ok := seen[top]; !ok {
				seen[top] = struct{}{}
				tops = append(tops, filepath.Join(destDir, top))
			}
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o700); err != nil {
				return tops, err
			}
		case tar.TypeSymlink:
			os.Remove(target)
			if err := os.Symlink(hdr.Linkname, target); err != nil {
				if windowsNames {
					// 普通用户在 Windows 上没有创建符号链接的权限（除非开了开发者模式）。
					// 跳过这一个，别让整次解包失败
					continue
				}
				return tops, fmt.Errorf("创建符号链接 %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
				return tops, err
			}
			f, err := os.OpenFile(target,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode).Perm())
			if err != nil {
				return tops, fmt.Errorf("创建 %s: %w", target, err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return tops, fmt.Errorf("写入 %s: %w", target, err)
			}
			f.Close()
		}
	}
	return tops, nil
}

// safeJoin 阻止 tar 条目用 ../ 逃出目标目录（Zip Slip）。
//
// 归档来自对端设备，即便已配对也不该无条件信任其中的路径。
func safeJoin(base, name string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(name))
	// VolumeName 拦住 Windows 上的「C:foo」「\\server\share」：它们不算绝对路径，Join 之后却能跑出去
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") || filepath.VolumeName(clean) != "" ||
		strings.HasPrefix(clean, string(filepath.Separator)) { // Windows 上「\etc」不算绝对路径
		return "", fmt.Errorf("归档中含非法路径: %q", name)
	}
	target := filepath.Join(base, clean)
	if !strings.HasPrefix(target, filepath.Clean(base)+string(os.PathSeparator)) &&
		target != filepath.Clean(base) {
		return "", fmt.Errorf("归档条目逃出目标目录: %q", name)
	}
	return target, nil
}

func topLevel(name string) string {
	clean := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(filepath.FromSlash(name))), "./")
	if clean == "." || clean == "" {
		return ""
	}
	if i := strings.Index(clean, "/"); i > 0 {
		return clean[:i]
	}
	return clean
}

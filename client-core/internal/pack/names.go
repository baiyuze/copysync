package pack

// 解包到 Windows 时的文件名处理。
//
// Mac 上合法的文件名在 Windows 上未必合法：轻则解包失败，重则出安全问题——
// 名字里的冒号会被 NTFS 当作「备用数据流」，`a.txt:x` 写进的是 a.txt 的隐藏流。
// 这里把每一级名字改成 Windows 能接受、又尽量保持原样的形式。只在 Windows 上启用，
// 其他平台解包时名字原样落盘。

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// windowsNames 决定解包时是否按 Windows 的规则处理文件名。是变量而不是常量，
// 好让测试在任何平台上都能验证这套规则。
var windowsNames = runtime.GOOS == "windows"

// Windows 不允许出现在文件名里的字符，换成对应的全角字符：名字仍然可读，
// 也不会与已有的文件撞名。
var fullwidth = map[rune]rune{
	':': '：', '*': '＊', '?': '？', '"': '＂',
	'<': '＜', '>': '＞', '|': '｜', '\\': '＼',
}

// 设备保留名：不管带不带扩展名（nul.txt 也不行），都打不开、删不掉
var reservedNames = map[string]bool{
	"CON": true, "PRN": true, "AUX": true, "NUL": true, "CONIN$": true, "CONOUT$": true,
	"COM1": true, "COM2": true, "COM3": true, "COM4": true, "COM5": true,
	"COM6": true, "COM7": true, "COM8": true, "COM9": true,
	"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true, "LPT5": true,
	"LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
}

// windowsSafeName 把一级名字改成 Windows 上合法的形式。
func windowsSafeName(name string) string {
	// Mac 上的文件名常是分解形式（「が」存成「か」加浊点两个码位），
	// 资源管理器会显示成两个字符、也搜不到。统一成组合形式
	name = norm.NFC.String(name)

	var b strings.Builder
	for _, r := range name {
		switch {
		case fullwidth[r] != 0:
			b.WriteRune(fullwidth[r])
		case r < 0x20:
			b.WriteByte('_') // 控制字符
		default:
			b.WriteRune(r)
		}
	}
	s := b.String()

	// 结尾的点和空格会被 Windows 悄悄去掉，「notes.」就成了「notes」
	if strings.TrimRight(s, ". ") != s {
		s += "_"
	}
	stem, _, _ := strings.Cut(s, ".")
	if reservedNames[strings.ToUpper(strings.TrimRight(stem, " "))] {
		s = "_" + s
	}
	if s == "" {
		s = "_"
	}
	return s
}

// skipOnWindows 是不该带到 Windows 上的 Mac 系统文件：
// Finder 的目录设置，以及在非 APFS 磁盘上保存扩展属性的 AppleDouble 文件。
func skipOnWindows(name string) bool {
	return name == ".DS_Store" || strings.HasPrefix(name, "._")
}

// skipWhenPacking 是在 Windows 上打包时跳过的系统文件：缩略图缓存与文件夹设置。
func skipWhenPacking(name string) bool {
	if !windowsNames {
		return false
	}
	switch strings.ToLower(name) {
	case "thumbs.db", "desktop.ini":
		return true
	}
	return false
}

// destNames 给解包出来的每个条目决定落盘的相对路径（正斜杠分隔）。
//
// 名字经过改写后可能与别的条目重名，Windows 又不区分大小写（Readme.md 与 README.md
// 是同一个文件），所以要记下已经用掉的名字：同名的目录合并，同名的文件改名为
// 「README (2).md」，不能互相覆盖。
type destNames struct {
	used map[string]usedName // 小写的落盘路径 → 实际占用的名字
	dirs map[string]string   // 归档里的目录 → 落盘路径
}

type usedName struct {
	path  string
	isDir bool
}

func newDestNames() *destNames {
	return &destNames{used: map[string]usedName{}, dirs: map[string]string{}}
}

// rel 返回归档条目 name 的落盘相对路径。含 .. 或绝对路径的条目直接拒绝。
func (d *destNames) rel(name string, isDir bool) (string, error) {
	clean := path.Clean(strings.TrimSuffix(name, "/"))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") {
		return "", fmt.Errorf("归档中含非法路径: %q", name)
	}
	parent := d.dir(path.Dir(clean))
	leaf := d.claim(parent, windowsSafeName(path.Base(clean)), isDir)
	if isDir {
		d.dirs[clean] = leaf
	}
	return leaf, nil
}

// dir 返回归档里某个目录的落盘路径；归档里没有单独的目录条目时，在这里补上。
func (d *destNames) dir(orig string) string {
	if orig == "." || orig == "" {
		return ""
	}
	if mapped, ok := d.dirs[orig]; ok {
		return mapped
	}
	mapped := d.claim(d.dir(path.Dir(orig)), windowsSafeName(path.Base(orig)), true)
	d.dirs[orig] = mapped
	return mapped
}

// claim 在 parent 下占用一个不重名的名字。
func (d *destNames) claim(parent, name string, isDir bool) string {
	candidate := path.Join(parent, name)
	ext := path.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	if isDir {
		stem, ext = name, ""
	}
	for n := 2; ; n++ {
		prev, taken := d.used[strings.ToLower(candidate)]
		if !taken {
			d.used[strings.ToLower(candidate)] = usedName{candidate, isDir}
			return candidate
		}
		if isDir && prev.isDir {
			// 两个只差大小写的目录在 Windows 上就是同一个：合并，沿用先占用的名字，
			// 免得复制列表里出现两个顶层条目
			return prev.path
		}
		candidate = path.Join(parent, fmt.Sprintf("%s (%d)%s", stem, n, ext))
	}
}

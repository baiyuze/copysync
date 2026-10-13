package diag

import (
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// 每个文件只读末尾这么多：要的最多一千来行，读整个文件没必要。Mac 上的日志由 launchd 写，不轮转，可能很大
const tailBytes = 512 << 10

// 单行太长的截断：一次网络出错的报错可能带出很长的内容，几行就能占满一份报告
const maxLineRunes = 400

// tailLines 取日志末尾的最多 n 行。files 按时间先后排列，轮转出去的旧文件在前：
// 最新的文件不够 n 行时，再从前一个文件的末尾补。读不到的文件跳过。
func tailLines(files []string, n int) []string {
	var out []string
	for i := len(files) - 1; i >= 0 && len(out) < n; i-- {
		out = append(tailFile(files[i], n-len(out)), out...)
	}
	return out
}

func tailFile(path string, n int) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil
	}
	start := max(0, info.Size()-tailBytes)
	if _, err := f.Seek(start, io.SeekStart); err != nil {
		return nil
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil
	}
	text := strings.TrimRight(strings.ToValidUTF8(string(data), "�"), "\r\n")
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	if start > 0 {
		lines = lines[1:] // 从文件中间开始读，第一行只读到了后半截
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	for i, l := range lines {
		l = strings.TrimRight(l, "\r")
		if utf8.RuneCountInString(l) > maxLineRunes {
			l = string([]rune(l)[:maxLineRunes]) + "…"
		}
		lines[i] = l
	}
	return lines
}

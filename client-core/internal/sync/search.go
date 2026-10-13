package sync

import (
	"slices"
	"strings"
	"unicode"

	"github.com/baiyuze/copysync/client-core/internal/store"
)

// SearchHistory 搜索复制记录：文字内容或文件名包含 query 的，不分大小写，按时间倒序，最多 limit 条。
//
// 数据库只能粗筛（见 store.SearchClips），这里按用户看得到的文字再核对一遍：
// HTML 记录看纯文本，不看标签。匹配处不在预览片段里的文字记录，预览换成匹配处前后的一段，
// 列表里一眼能看出为什么搜到了它。
func (e *Engine) SearchHistory(query string, limit int) ([]store.Clip, error) {
	q := []rune(strings.Map(unicode.ToLower, strings.TrimSpace(query)))
	if len(q) == 0 {
		return nil, nil
	}
	candidates, err := e.store.SearchClips(string(q), limit)
	if err != nil {
		return nil, err
	}
	out := candidates[:0]
	for _, c := range candidates {
		if c.Kind != store.KindImage && slices.ContainsFunc(c.Items, func(it store.Item) bool {
			return indexFold(it.Name, q) >= 0
		}) {
			out = append(out, c)
			continue
		}
		if c.Kind != store.KindText && c.Kind != store.KindHTML {
			continue
		}
		text := plainText(c)
		at := indexFold(text, q)
		if at < 0 {
			continue
		}
		if indexFold(c.TextPreview, q) < 0 {
			c.TextPreview = snippet(text, at)
		}
		out = append(out, c)
	}
	return out, nil
}

// indexFold 返回 q（已转小写）在 s 里第一次出现的位置，按 rune 计，不分大小写；找不到返回 -1。
func indexFold(s string, q []rune) int {
	r := []rune(strings.Map(unicode.ToLower, s)) // 逐个 rune 转换，位置与原文一一对应
	for i := 0; i+len(q) <= len(r); i++ {
		if slices.Equal(r[i:i+len(q)], q) {
			return i
		}
	}
	return -1
}

// snippet 从匹配处往前留一小段上下文开始截取预览，前面被截掉的部分用省略号表示。
func snippet(text string, at int) string {
	const lead = 16
	r := []rune(text)
	if at <= lead {
		return preview(text)
	}
	return "…" + preview(string(r[at-lead:]))
}

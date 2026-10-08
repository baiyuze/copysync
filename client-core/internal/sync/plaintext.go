package sync

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/baiyuze/copysync/client-core/internal/store"
)

// plainText 返回记录的纯文本形式：写入剪贴板的纯文本类型、生成预览都用它。
//
// HTML 记录必须单独给出纯文本——只认纯文本的输入框会把 HTML 原样粘出来，
// 满屏 <meta charset='utf-8'> 之类的标签。
func plainText(c store.Clip) string {
	if c.Kind != store.KindHTML {
		return c.TextContent
	}
	if c.TextPlain != "" {
		return c.TextPlain
	}
	// 升级前存下的记录、旧版本对端发来的内容都没有纯文本，只能从 HTML 里提取
	return htmlToText(c.TextContent)
}

// htmlToText 把 HTML 粗略转成纯文本：去掉标签与不可见内容，块级元素换行。
// 只是兜底，源应用自己放进剪贴板的纯文本永远比这准确。
func htmlToText(s string) string {
	z := html.NewTokenizer(strings.NewReader(s))
	var w textWriter
	hidden, pre := 0, 0 // 处于 script/style 内、<pre> 内的嵌套深度

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return strings.TrimSpace(w.String())

		case html.TextToken:
			if hidden > 0 {
				continue
			}
			if pre > 0 {
				w.raw(string(z.Text()))
			} else {
				w.collapsed(string(z.Text()))
			}

		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken:
			name, _ := z.TagName()
			a := atom.Lookup(name)
			start := tt != html.EndTagToken
			switch {
			case a == atom.Script || a == atom.Style || a == atom.Title || a == atom.Template:
				if tt == html.StartTagToken {
					hidden++
				} else if tt == html.EndTagToken && hidden > 0 {
					hidden--
				}
			case a == atom.Br:
				w.raw("\n")
			case a == atom.Pre:
				if tt == html.StartTagToken {
					pre++
				} else if tt == html.EndTagToken && pre > 0 {
					pre--
				}
				w.lineBreak()
			case a == atom.Td || a == atom.Th:
				if start {
					w.separator()
				}
			case blockElements[a]:
				w.lineBreak()
			}
		}
	}
}

var blockElements = map[atom.Atom]bool{
	atom.P: true, atom.Div: true, atom.Li: true, atom.Ul: true, atom.Ol: true,
	atom.Tr: true, atom.Table: true, atom.Blockquote: true, atom.Hr: true,
	atom.H1: true, atom.H2: true, atom.H3: true, atom.H4: true, atom.H5: true, atom.H6: true,
	atom.Section: true, atom.Article: true, atom.Header: true, atom.Footer: true,
	atom.Dt: true, atom.Dd: true,
}

// textWriter 负责空白处理：HTML 源码里的换行与缩进不算内容，
// 换行只由块级元素与 <br> 产生。
type textWriter struct {
	b    strings.Builder
	last rune
}

func (w *textWriter) String() string { return w.b.String() }

func (w *textWriter) atLineStart() bool { return w.b.Len() == 0 || w.last == '\n' }

func (w *textWriter) raw(s string) {
	if s == "" {
		return
	}
	w.b.WriteString(s)
	w.last, _ = utf8.DecodeLastRuneInString(s)
}

// collapsed 按浏览器的规则把连续空白折叠成一个空格，行首空白直接丢弃。
func (w *textWriter) collapsed(s string) {
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !w.atLineStart() && w.last != ' ' {
				w.b.WriteByte(' ')
				w.last = ' '
			}
			continue
		}
		w.b.WriteRune(r)
		w.last = r
	}
}

func (w *textWriter) lineBreak() {
	if !w.atLineStart() {
		w.trimTrailingSpace()
		w.b.WriteByte('\n')
		w.last = '\n'
	}
}

func (w *textWriter) separator() {
	if !w.atLineStart() && w.last != ' ' {
		w.b.WriteByte(' ')
		w.last = ' '
	}
}

func (w *textWriter) trimTrailingSpace() {
	if w.last != ' ' {
		return
	}
	s := strings.TrimRight(w.b.String(), " ")
	w.b.Reset()
	w.b.WriteString(s)
	w.last, _ = utf8.DecodeLastRuneInString(s) // 空串时为 RuneError，atLineStart 看的是长度
}

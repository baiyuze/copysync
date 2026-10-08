package sync

import (
	"testing"

	"github.com/baiyuze/copysync/client-core/internal/store"
)

func TestHTMLToText(t *testing.T) {
	cases := []struct{ name, in, want string }{
		// 浏览器、Electron 应用复制文字时放进剪贴板的 HTML 就长这样
		{"浏览器的 meta 前缀", "<meta charset='utf-8'>用新签名覆盖安装，因为我快过期了签名",
			"用新签名覆盖安装，因为我快过期了签名"},
		{"实体解码", "<span>a &amp; b&nbsp;&lt;c&gt;</span>", "a & b <c>"},
		{"块级元素换行", "<div>第一行</div><div>第二行</div>", "第一行\n第二行"},
		{"br 换行", "一<br>二<br/>三", "一\n二\n三"},
		{"源码缩进不算内容", "<ul>\n  <li>甲</li>\n  <li>乙</li>\n</ul>", "甲\n乙"},
		{"文本内空白折叠", "<p>\n    hello\n    world\n</p>", "hello world"},
		{"不可见内容丢弃", "<style>p{color:red}</style><p>正文</p><script>alert(1)</script>", "正文"},
		{"pre 保留格式", "<pre>a\n  b</pre>", "a\n  b"},
		{"表格", "<table><tr><td>A</td><td>B</td></tr><tr><td>C</td><td>D</td></tr></table>",
			"A B\nC D"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := htmlToText(c.in); got != c.want {
				t.Errorf("htmlToText(%q) = %q，期望 %q", c.in, got, c.want)
			}
		})
	}
}

func TestPlainText(t *testing.T) {
	// 源应用给出的纯文本优先，它比从 HTML 提取的准确
	html := store.Clip{Kind: store.KindHTML, TextContent: "<b>粗体</b>", TextPlain: "**粗体**"}
	if got := plainText(html); got != "**粗体**" {
		t.Errorf("有纯文本时应原样使用，得到 %q", got)
	}

	// 升级前的记录、旧版本对端发来的内容没有纯文本
	legacy := store.Clip{Kind: store.KindHTML, TextContent: "<meta charset='utf-8'>你好"}
	if got := plainText(legacy); got != "你好" {
		t.Errorf("缺纯文本时应从 HTML 提取，得到 %q", got)
	}

	// 纯文本记录里的尖括号是用户复制的内容本身（比如一段代码），绝不能当标签剥掉
	text := store.Clip{Kind: store.KindText, TextContent: "<div>不是标签</div>"}
	if got := plainText(text); got != "<div>不是标签</div>" {
		t.Errorf("纯文本记录应原样返回，得到 %q", got)
	}
}

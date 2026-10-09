package main

import "testing"

func TestTrayLabels(t *testing.T) {
	cases := map[string]string{"zh-Hans": "打开 CopySync", "en": "Open CopySync", "ja": "CopySync を開く"}
	for lang, want := range cases {
		if open, _ := trayLabels(lang); open != want {
			t.Errorf("trayLabels(%q) = %q，应为 %q", lang, open, want)
		}
	}
	// 跟随系统时总有一个结果，不会是空的
	if open, quit := trayLabels(""); open == "" || quit == "" {
		t.Error("跟随系统时菜单文字为空")
	}
}

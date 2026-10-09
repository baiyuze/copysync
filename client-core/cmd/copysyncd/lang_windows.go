//go:build windows

package main

import "golang.org/x/sys/windows"

var procGetUserDefaultUILanguage = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetUserDefaultUILanguage")

// systemLanguage 按 Windows 的界面语言返回 "zh"、"ja" 或 "en"。
func systemLanguage() string {
	id, _, _ := procGetUserDefaultUILanguage.Call()
	switch id & 0x3ff { // 低 10 位是主语言
	case 0x04: // LANG_CHINESE
		return "zh"
	case 0x11: // LANG_JAPANESE
		return "ja"
	default:
		return "en"
	}
}

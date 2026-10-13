//go:build windows

package main

import "golang.org/x/sys/windows"

var procGetUserDefaultUILanguage = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetUserDefaultUILanguage")

// systemLanguage 按 Windows 的界面语言返回 "zh"、"zh-Hant-TW"、"zh-Hant-HK"、"ja" 或 "en"。
func systemLanguage() string {
	id, _, _ := procGetUserDefaultUILanguage.Call()
	switch id & 0x3ff { // 低 10 位是主语言，高 6 位是地区
	case 0x04: // LANG_CHINESE
		switch id {
		case 0x0c04, 0x1404: // 香港、澳门
			return "zh-Hant-HK"
		case 0x0404, 0x7c04: // 台湾，以及没写地区的繁体
			return "zh-Hant-TW"
		}
		return "zh"
	case 0x11: // LANG_JAPANESE
		return "ja"
	default:
		return "en"
	}
}

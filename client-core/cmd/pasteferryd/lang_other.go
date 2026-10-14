//go:build !windows

package main

// 只有 Windows 的托盘用得到；其他平台没有托盘。
func systemLanguage() string { return "en" }

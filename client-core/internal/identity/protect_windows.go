//go:build windows

package identity

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// keyProtected 表示私钥落盘前要加密。Windows 上用 DPAPI：密钥由系统按当前用户派生，
// 换个用户、或把文件拷到别的电脑上都解不开。
const keyProtected = true

const cryptprotectUIForbidden = 0x1

func protectKey(plain []byte) ([]byte, error) {
	in := windows.DataBlob{Size: uint32(len(plain)), Data: &plain[0]}
	var out windows.DataBlob
	if err := windows.CryptProtectData(&in, nil, nil, 0, nil, cryptprotectUIForbidden, &out); err != nil {
		return nil, fmt.Errorf("加密私钥: %w", err)
	}
	return takeBlob(out), nil
}

func unprotectKey(sealed []byte) ([]byte, error) {
	if len(sealed) == 0 {
		return nil, fmt.Errorf("加密的私钥为空")
	}
	in := windows.DataBlob{Size: uint32(len(sealed)), Data: &sealed[0]}
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, cryptprotectUIForbidden, &out); err != nil {
		return nil, fmt.Errorf("解密私钥（身份文件可能来自别的用户或别的电脑）: %w", err)
	}
	return takeBlob(out), nil
}

// takeBlob 复制出系统分配的结果并释放它。
func takeBlob(b windows.DataBlob) []byte {
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(b.Data)))
	return append([]byte(nil), unsafe.Slice(b.Data, b.Size)...)
}

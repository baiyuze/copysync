//go:build !windows

package identity

// 其他平台私钥以 0600 权限明文存放，与之前的版本一致。
const keyProtected = false

func protectKey(plain []byte) ([]byte, error)    { return plain, nil }
func unprotectKey(sealed []byte) ([]byte, error) { return sealed, nil }

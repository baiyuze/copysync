// Package deviceid 定义设备 ID 与指纹的派生规则。
//
// 这属于协议契约的一部分：客户端与服务器必须用完全相同的规则，
// 否则服务器无法校验「这个 device_id 确实属于这把公钥」。
// 因此放在 proto 模块里与 .proto 一同维护，而不是各自实现一份。
package deviceid

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base32"
	"strings"
)

// 不带填充的 Base32：字符集无易混字符，适合展示与口头转述。
var enc = base32.StdEncoding.WithPadding(base32.NoPadding)

// For 由公钥确定性地派生设备 ID。
//
// 做成确定性派生，任何一方都能独立校验 ID 与公钥的对应关系——
// 服务器因此无需维护用户表也能拒绝伪造的 device_id。
func For(pub ed25519.PublicKey) string {
	sum := sha256.Sum256(pub)
	return strings.ToLower(enc.EncodeToString(sum[:16]))
}

// Fingerprint 返回供用户肉眼核对的短指纹，形如 "K7X2-M9PL-4RT8"。
//
// 配对时双方各自显示，用户确认一致后才建立互信。这一步是
// 防止信令服务器在配对阶段做中间人的关键——服务器可以替换公钥，
// 但无法让两端显示出相同的指纹。
func Fingerprint(pub ed25519.PublicKey) string {
	sum := sha256.Sum256(pub)
	s := enc.EncodeToString(sum[:8])
	parts := make([]string, 0, 3)
	for i := 0; i+4 <= len(s) && len(parts) < 3; i += 4 {
		parts = append(parts, s[i:i+4])
	}
	return strings.Join(parts, "-")
}

// Valid 校验 id 是否确实由 pub 派生而来。
func Valid(id string, pub ed25519.PublicKey) bool {
	if len(pub) != ed25519.PublicKeySize {
		return false
	}
	return id == For(pub)
}

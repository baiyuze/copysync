// Package pairingcode 定义配对码的格式。
//
// 配对码既可能由信令服务器签发，也可能由设备在局域网里自己生成；
// 两边必须是同一种格式，用户在界面里输入时才不必区分来源。
// 因此与 deviceid 一样放在 proto 模块里，由服务器和客户端共用。
package pairingcode

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// Alphabet 刻意剔除易混字符（0/O、1/I/L），因为用户可能要口头转述。
const Alphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

const (
	Length = 6
	TTL    = 5 * time.Minute
)

// Generate 生成一个新的配对码。
func Generate() (string, error) {
	out := make([]byte, Length)
	max := big.NewInt(int64(len(Alphabet)))
	for i := range out {
		// 用密码学随机源：配对码在有效期内是唯一的准入凭证
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("生成配对码: %w", err)
		}
		out[i] = Alphabet[n.Int64()]
	}
	return string(out), nil
}

// Normalize 把用户输入的配对码整理成标准形式：去掉空白与连字符，转成大写。
func Normalize(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	return strings.NewReplacer(" ", "", "-", "").Replace(code)
}

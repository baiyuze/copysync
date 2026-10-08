package rpc

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Endpoint 是 daemon 写给 UI 的握手信息。
//
// daemon 监听 127.0.0.1 的随机端口，把端口和一次性 token 写进这个文件
// （权限 0600）。UI 启动时读取它来连接。用 token 是因为本机任何进程都能
// 连上 loopback 端口，仅靠"只监听 127.0.0.1"不足以防止其他程序乱连。
type Endpoint struct {
	Port    int    `json:"port"`
	Token   string `json:"token"`
	PID     int    `json:"pid"`
	Version string `json:"version"`
}

func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("生成 token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func WriteEndpoint(path string, ep Endpoint) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ep, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("写入 endpoint 文件: %w", err)
	}
	return os.Rename(tmp, path)
}

func ReadEndpoint(path string) (Endpoint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Endpoint{}, err
	}
	var ep Endpoint
	if err := json.Unmarshal(data, &ep); err != nil {
		return Endpoint{}, fmt.Errorf("解析 endpoint 文件: %w", err)
	}
	return ep, nil
}

// RemoveEndpoint 在 daemon 退出时清理，避免 UI 连到一个已死的端口。
func RemoveEndpoint(path string) { _ = os.Remove(path) }

//go:build windows

package identity

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestKeyIsProtectedOnDisk(t *testing.T) {
	path := filepath.Join(t.TempDir(), "identity.json")
	id, err := LoadOrCreate(path)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	var kf keyFile
	if err := json.Unmarshal(raw, &kf); err != nil {
		t.Fatal(err)
	}
	if len(kf.PrivateKey) != 0 || len(kf.ProtectedKey) == 0 {
		t.Fatalf("磁盘上应只有加密后的私钥：%s", raw)
	}
	again, err := LoadOrCreate(path)
	if err != nil {
		t.Fatal(err)
	}
	if again.DeviceID != id.DeviceID || !bytes.Equal(again.privateKey, id.privateKey) {
		t.Error("重新读取后身份不一致")
	}
}

// 旧格式（明文私钥）读取后自动改存为加密格式，设备身份不变。
func TestPlaintextKeyIsMigrated(t *testing.T) {
	path := filepath.Join(t.TempDir(), "identity.json")
	id, err := LoadOrCreate(path)
	if err != nil {
		t.Fatal(err)
	}
	plain, _ := json.Marshal(keyFile{DeviceID: id.DeviceID, PublicKey: id.PublicKey, PrivateKey: id.privateKey})
	os.WriteFile(path, plain, 0o600)

	if _, err := LoadOrCreate(path); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	var kf keyFile
	json.Unmarshal(raw, &kf)
	if len(kf.PrivateKey) != 0 || len(kf.ProtectedKey) == 0 {
		t.Errorf("明文私钥没有被加密存回：%s", raw)
	}
}

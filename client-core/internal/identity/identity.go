// Package identity 管理本机的 Ed25519 设备身份。
//
// 设备身份是整套安全模型的根：配对时双方交换公钥并由用户肉眼核对指纹，
// 之后所有经信令服务器转发的消息都带签名。服务器因此无法伪造 peer，
// 也就无法对 P2P 握手做中间人攻击。
package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/baiyuze/copysync/proto/deviceid"
)

type Identity struct {
	DeviceID   string
	PublicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

// 磁盘格式。私钥以 0600 权限存放在用户目录下——与 Syncthing 等
// 同类工具一致。后续可换成 Keychain / DPAPI，接口不变。
type keyFile struct {
	DeviceID   string `json:"device_id"`
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}

// LoadOrCreate 读取已有身份；不存在则生成一份并落盘。
func LoadOrCreate(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		var kf keyFile
		if err := json.Unmarshal(data, &kf); err != nil {
			return nil, fmt.Errorf("解析身份文件 %s: %w", path, err)
		}
		if len(kf.PrivateKey) != ed25519.PrivateKeySize ||
			len(kf.PublicKey) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("身份文件 %s 已损坏（密钥长度不对）", path)
		}
		id := &Identity{
			DeviceID:   kf.DeviceID,
			PublicKey:  kf.PublicKey,
			privateKey: kf.PrivateKey,
		}
		// 防御损坏或被篡改的文件：ID 必须能由公钥推出
		if want := DeviceIDFor(id.PublicKey); want != id.DeviceID {
			return nil, fmt.Errorf("身份文件 %s 的 device_id 与公钥不匹配", path)
		}
		return id, nil

	case os.IsNotExist(err):
		return create(path)

	default:
		return nil, fmt.Errorf("读取身份文件: %w", err)
	}
}

func create(path string) (*Identity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成密钥: %w", err)
	}
	id := &Identity{
		DeviceID:   DeviceIDFor(pub),
		PublicKey:  pub,
		privateKey: priv,
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(keyFile{
		DeviceID:   id.DeviceID,
		PublicKey:  pub,
		PrivateKey: priv,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	// 0600：私钥只有当前用户可读
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return nil, fmt.Errorf("写入身份文件: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return nil, err
	}
	return id, nil
}

func (i *Identity) Sign(data []byte) []byte {
	return ed25519.Sign(i.privateKey, data)
}

func Verify(pub ed25519.PublicKey, data, sig []byte) bool {
	if len(pub) != ed25519.PublicKeySize {
		return false
	}
	return ed25519.Verify(pub, data, sig)
}

// Fingerprint 返回供用户肉眼核对的短指纹，形如 "K7X2-M9PL-4RT8"。
func (i *Identity) Fingerprint() string { return deviceid.Fingerprint(i.PublicKey) }

// 派生规则属于协议契约，与服务器共用同一份实现（proto/deviceid）。
var (
	FingerprintFor = deviceid.Fingerprint
	DeviceIDFor    = deviceid.For
)

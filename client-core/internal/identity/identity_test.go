package identity

import (
	"crypto/ed25519"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/baiyuze/copysync/proto/deviceid"
)

func TestCreateThenLoadIsStable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "identity.json")

	first, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("创建身份: %v", err)
	}
	second, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("重新加载: %v", err)
	}

	// 身份必须跨重启稳定，否则每次启动都会变成"新设备"，已配对关系全失效
	if first.DeviceID != second.DeviceID {
		t.Errorf("device_id 不稳定: %q vs %q", first.DeviceID, second.DeviceID)
	}
	if !first.PublicKey.Equal(second.PublicKey) {
		t.Error("公钥不稳定")
	}
	// 私钥也要能正常加载：签名要能被原公钥验证
	sig := second.Sign([]byte("payload"))
	if !Verify(first.PublicKey, []byte("payload"), sig) {
		t.Error("重新加载后的私钥与原公钥不配对")
	}
}

func TestPrivateKeyFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows 不用权限位，私钥另用 DPAPI 加密，见 protect_windows_test.go")
	}
	path := filepath.Join(t.TempDir(), "identity.json")
	if _, err := LoadOrCreate(path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	// 文件里是设备私钥，同机其他用户绝不能读到
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("身份文件权限 = %o，期望 0600", perm)
	}
}

func TestSignAndVerify(t *testing.T) {
	id, err := LoadOrCreate(filepath.Join(t.TempDir(), "identity.json"))
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("signal payload")
	sig := id.Sign(msg)

	if !Verify(id.PublicKey, msg, sig) {
		t.Error("自己的签名验证失败")
	}
	if Verify(id.PublicKey, []byte("tampered"), sig) {
		t.Error("篡改后的内容仍通过验签")
	}

	other, _ := LoadOrCreate(filepath.Join(t.TempDir(), "other.json"))
	if Verify(other.PublicKey, msg, sig) {
		t.Error("用他人公钥验签竟然通过了")
	}
	// 长度不对的公钥不能 panic，信令里的公钥是对端传来的
	if Verify(ed25519.PublicKey("short"), msg, sig) {
		t.Error("非法长度的公钥应当直接判定失败")
	}
}

// device_id 由公钥确定性派生，服务器据此校验"这个 ID 确实属于这把公钥"。
// 客户端与服务器必须用同一套规则，否则设备根本连不上。
func TestDeviceIDDerivationMatchesSharedRule(t *testing.T) {
	id, err := LoadOrCreate(filepath.Join(t.TempDir(), "identity.json"))
	if err != nil {
		t.Fatal(err)
	}
	if id.DeviceID != deviceid.For(id.PublicKey) {
		t.Errorf("device_id 与共享派生规则不一致: %q vs %q",
			id.DeviceID, deviceid.For(id.PublicKey))
	}
	if !deviceid.Valid(id.DeviceID, id.PublicKey) {
		t.Error("自身的 device_id 未通过校验")
	}
	if id.Fingerprint() != deviceid.Fingerprint(id.PublicKey) {
		t.Error("指纹与共享规则不一致")
	}
}

// 指纹要给用户逐字比对，必须短、分组、无易混字符。
func TestFingerprintIsHumanComparable(t *testing.T) {
	id, _ := LoadOrCreate(filepath.Join(t.TempDir(), "identity.json"))
	fp := id.Fingerprint()

	parts := strings.Split(fp, "-")
	if len(parts) != 3 {
		t.Errorf("指纹 = %q，期望三段分组便于朗读", fp)
	}
	for _, p := range parts {
		if len(p) != 4 {
			t.Errorf("分段 %q 长度应为 4", p)
		}
	}
	// Base32 字符集不含 0/1/8/9，避免与 O/I/B/g 混淆
	for _, bad := range []rune{'0', '1', '8', '9'} {
		if strings.ContainsRune(fp, bad) {
			t.Errorf("指纹含易混字符 %q: %s", bad, fp)
		}
	}
}

// 身份文件被篡改（换了公钥但留着旧 ID）时必须拒绝加载，
// 否则会用一个与公钥不匹配的 ID 去连服务器，且永远连不上。
func TestTamperedIdentityRejected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "identity.json")
	if _, err := LoadOrCreate(path); err != nil {
		t.Fatal(err)
	}

	raw, _ := os.ReadFile(path)
	var kf keyFile
	if err := json.Unmarshal(raw, &kf); err != nil {
		t.Fatal(err)
	}
	kf.DeviceID = "someoneelsesdeviceid"
	patched, _ := json.Marshal(kf)
	if err := os.WriteFile(path, patched, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadOrCreate(path); err == nil {
		t.Error("device_id 与公钥不匹配的身份文件应被拒绝")
	}
}

func TestCorruptIdentityRejected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "identity.json")
	if err := os.WriteFile(path, []byte(`{"device_id":"x","public_key":"AAAA"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadOrCreate(path); err == nil {
		t.Error("密钥长度不合法的文件应被拒绝")
	}
}

func TestDistinctDevicesGetDistinctIDs(t *testing.T) {
	a, _ := LoadOrCreate(filepath.Join(t.TempDir(), "a.json"))
	b, _ := LoadOrCreate(filepath.Join(t.TempDir(), "b.json"))

	if a.DeviceID == b.DeviceID {
		t.Error("两台设备生成了相同的 device_id")
	}
	if a.Fingerprint() == b.Fingerprint() {
		t.Error("两台设备生成了相同的指纹，用户将无法通过比对发现异常")
	}
}

package pairing

import (
	"strings"
	"testing"
	"time"
)

func peerNamed(id string) Peer {
	return Peer{DeviceID: id, DeviceName: id, PublicKey: []byte("key-" + id)}
}

func TestCodeFormat(t *testing.T) {
	b := NewBroker()
	code, _, _, err := b.Create(peerNamed("alice"))
	if err != nil {
		t.Fatalf("生成配对码: %v", err)
	}
	if len(code) != CodeLength {
		t.Errorf("长度 = %d，期望 %d", len(code), CodeLength)
	}
	// 字符集刻意剔除了易混字符，用户可能要口头转述配对码
	for _, c := range code {
		if !strings.ContainsRune(codeAlphabet, c) {
			t.Errorf("配对码含字符集外的字符 %q（可能与 0/O、1/I 混淆）", c)
		}
	}
	for _, bad := range []string{"0", "O", "1", "I", "L"} {
		if strings.Contains(codeAlphabet, bad) {
			t.Errorf("字符集不该包含易混字符 %q", bad)
		}
	}
}

// 配对码是一次性凭证：兑换后必须立即失效，否则第三方可以重复使用它接入。
func TestCodeIsSingleUse(t *testing.T) {
	b := NewBroker()
	code, _, redeemed, _ := b.Create(peerNamed("alice"))

	initiator, err := b.Redeem(code, peerNamed("bob"))
	if err != nil {
		t.Fatalf("首次兑换应成功: %v", err)
	}
	if initiator.DeviceID != "alice" {
		t.Errorf("拿到的发起方 = %q，期望 alice", initiator.DeviceID)
	}

	if _, err := b.Redeem(code, peerNamed("mallory")); err == nil {
		t.Error("配对码被重复兑换了，应当一次性失效")
	}

	// 发起方必须收到兑换方信息，否则只能建立单向互信
	select {
	case got := <-redeemed:
		if got.DeviceID != "bob" {
			t.Errorf("回传的兑换方 = %q，期望 bob", got.DeviceID)
		}
	case <-time.After(time.Second):
		t.Error("发起方未收到兑换通知")
	}
}

func TestRedeemUnknownCode(t *testing.T) {
	b := NewBroker()
	if _, err := b.Redeem("ZZZZZZ", peerNamed("bob")); err != ErrCodeNotFound {
		t.Errorf("err = %v，期望 ErrCodeNotFound", err)
	}
}

// 与自己配对没有意义，且会让设备列表里出现自身。
func TestSelfPairingRejected(t *testing.T) {
	b := NewBroker()
	code, _, _, _ := b.Create(peerNamed("alice"))
	if _, err := b.Redeem(code, peerNamed("alice")); err != ErrSelfPairing {
		t.Errorf("err = %v，期望 ErrSelfPairing", err)
	}
}

func TestExpiredCodeRejected(t *testing.T) {
	b := NewBroker()
	now := time.Now()
	b.now = func() time.Time { return now }

	code, expires, redeemed, _ := b.Create(peerNamed("alice"))
	if !expires.Equal(now.Add(CodeTTL)) {
		t.Errorf("过期时间 = %v，期望 %v", expires, now.Add(CodeTTL))
	}

	// 拨快到过期之后
	b.now = func() time.Time { return now.Add(CodeTTL + time.Second) }

	if _, err := b.Redeem(code, peerNamed("bob")); err != ErrCodeNotFound {
		t.Errorf("过期配对码仍可兑换，err = %v", err)
	}
	// 过期时应关闭通道，避免发起方永久挂起
	select {
	case _, open := <-redeemed:
		if open {
			t.Error("过期后不该再回传兑换方信息")
		}
	case <-time.After(time.Second):
		t.Error("过期后通道未关闭，发起方会一直等待")
	}
	if b.Len() != 0 {
		t.Errorf("过期条目未被清理，剩余 %d 条", b.Len())
	}
}

func TestCancelReleasesWaiter(t *testing.T) {
	b := NewBroker()
	code, _, redeemed, _ := b.Create(peerNamed("alice"))

	b.Cancel(code)

	select {
	case _, open := <-redeemed:
		if open {
			t.Error("取消后不该再收到兑换方信息")
		}
	case <-time.After(time.Second):
		t.Error("取消后通道未关闭")
	}
	if _, err := b.Redeem(code, peerNamed("bob")); err != ErrCodeNotFound {
		t.Error("已取消的配对码仍可兑换")
	}
}

func TestConcurrentCodesAreUnique(t *testing.T) {
	b := NewBroker()
	seen := make(map[string]bool)
	for i := 0; i < 200; i++ {
		code, _, _, err := b.Create(peerNamed("alice"))
		if err != nil {
			t.Fatalf("第 %d 次生成失败: %v", i, err)
		}
		if seen[code] {
			t.Fatalf("生成了重复的配对码 %q", code)
		}
		seen[code] = true
	}
	if b.Len() != 200 {
		t.Errorf("在存配对码数 = %d，期望 200", b.Len())
	}
}

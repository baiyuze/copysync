package lan

import (
	"context"
	"errors"
	"testing"
	"time"
)

// waitNearby 等到 n 的「附近的设备」里正好有 want 台。
func waitNearby(t *testing.T, n *node, want int) []NearbyInstance {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for {
		got := n.t.Nearby()
		if len(got) == want {
			return got
		}
		if time.Now().After(deadline) {
			t.Fatalf("%s 附近应有 %d 台设备，得到 %+v", n.id.DeviceID[:6], want, got)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// 设备名只在窗口打开期间公开；两台都打开才互相看得到，关掉就消失。
func TestNearbyOnlyWhileShown(t *testing.T) {
	net := NewMemoryNetwork()
	a, b := newNode(t, net), newNode(t, net)

	if got := a.t.Nearby(); len(got) != 0 {
		t.Fatalf("没打开窗口时不该有附近的设备: %+v", got)
	}
	a.t.mu.Lock()
	for _, si := range a.t.instances {
		if _, ok := nearbyName(si.inst.Text); ok {
			t.Errorf("没打开窗口的设备不该广播设备名: %v", si.inst.Text)
		}
	}
	a.t.mu.Unlock()

	hideA := a.t.ShowNearby()
	defer hideA()
	waitNearby(t, a, 0) // b 还没打开窗口

	hideB := b.t.ShowNearby()
	got := waitNearby(t, a, 1)
	if got[0].Name != "node-"+b.id.DeviceID[:4] || got[0].Platform != "test" || got[0].PeerID != "" {
		t.Fatalf("附近的设备信息不对: %+v", got[0])
	}
	waitNearby(t, b, 1)

	hideB()
	hideB() // 多调一次不该出错
	waitNearby(t, a, 0)
}

// 已配对的设备也显示，按令牌认出来是谁。
func TestNearbyRecognizesPairedDevice(t *testing.T) {
	net := NewMemoryNetwork()
	a, b := newNode(t, net), newNode(t, net)
	pair(a, b)
	defer a.t.ShowNearby()()
	defer b.t.ShowNearby()()

	if got := waitNearby(t, a, 1); got[0].PeerID != b.id.DeviceID {
		t.Fatalf("应当认出已配对的设备: %+v", got[0])
	}
}

// 点了附近的设备：不用配对码就交换身份，对方的上层收到请求、由用户决定。
func TestNearbyPairing(t *testing.T) {
	net := NewMemoryNetwork()
	a, b := newNode(t, net), newNode(t, net)
	defer a.t.ShowNearby()()
	hideB := b.t.ShowNearby()
	defer hideB()

	target := waitNearby(t, a, 1)[0]
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	peer, err := a.t.PairNearby(ctx, target.ID)
	if err != nil {
		t.Fatalf("配对失败: %v", err)
	}
	if peer.ID != b.id.DeviceID || string(peer.PublicKey) != string(b.id.PublicKey) {
		t.Fatalf("拿到的不是对方的身份: %+v", peer)
	}
	b.mu.Lock()
	got := b.nearbyReq
	b.mu.Unlock()
	if len(got) != 1 || got[0].ID != a.id.DeviceID {
		t.Fatalf("对方应当收到一次配对请求，得到 %+v", got)
	}
	if len(b.accepted) != 0 {
		t.Fatal("不该走配对码那条路")
	}

	// 对方关掉了窗口：拒收，很快返回「找不到」，而不是等到超时
	hideB()
	start := time.Now()
	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()
	if _, err := a.t.PairNearby(ctx2, target.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("对方关掉窗口后应当找不到，得到 %v", err)
	}
	if time.Since(start) > 2*time.Second {
		t.Errorf("被拒收后应当很快返回，用了 %v", time.Since(start))
	}
	if _, err := a.t.PairNearby(ctx2, "copysync-nobody"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("不认识的实例应当找不到，得到 %v", err)
	}
}

// mDNS 不报下线：对方不再发 ping，几秒后就从列表里消失；一直发着就一直在。
func TestNearbyNeedsPings(t *testing.T) {
	net := NewMemoryNetwork()
	a, b := newNode(t, net), newNode(t, net)
	defer a.t.ShowNearby()()
	defer b.t.ShowNearby()()
	target := waitNearby(t, a, 1)[0]

	age := func() {
		a.t.mu.Lock()
		si := a.t.instances[target.ID]
		si.seen = time.Now().Add(-time.Minute)
		a.t.instances[target.ID] = si
		a.t.mu.Unlock()
	}
	age()
	if got := a.t.Nearby(); len(got) != 0 {
		t.Fatalf("太久没有消息的设备应当消失: %+v", got)
	}
	age()
	b.t.pingNearby()
	waitNearby(t, a, 1)
}

func TestNearbyNameEncoding(t *testing.T) {
	text := map[string]string{}
	tr := &Transport{visible: 1, visibleName: "光光的 MacBook Pro（公司）"}
	tr.nearbyAdvertLocked(text)
	if name, ok := nearbyName(text); !ok || name != "光光的 MacBook Pro（公司）" {
		t.Fatalf("设备名往返后不一致: %q %v", name, ok)
	}

	long := ""
	for i := 0; i < 40; i++ {
		long += "设"
	}
	got := truncateUTF8(long, maxNearbyName)
	if len(got) > maxNearbyName || len(got)%3 != 0 {
		t.Fatalf("截断要落在字符边界上: %d 字节", len(got))
	}
	if _, ok := nearbyName(map[string]string{"n": "!!"}); ok {
		t.Fatal("编码不对的设备名不该接受")
	}
}

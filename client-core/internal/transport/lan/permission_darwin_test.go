package lan

import (
	"net/netip"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/identity"
)

// 探测本地网络权限的目标在默认路由所在的网段里，又不是本机。没有网络的机器上跳过。
func TestProbeTargetIsOnLink(t *testing.T) {
	ip, ok := probeTarget()
	if !ok {
		t.Skip("这台机器没有默认路由所在的私有网段")
	}
	if !ip.Is4() || !onLink(ip) {
		t.Fatalf("探测目标应当是直连网段里别人的地址，得到 %v", ip)
	}
	t.Logf("探测目标 %v", ip)
}

// 本机自己的地址不算直连网段里的别人：发往它的包不受本地网络权限限制。
func TestOnLinkExcludesOwnAddress(t *testing.T) {
	for _, ip := range LocalAddrs() {
		if onLink(ip) {
			t.Errorf("本机地址 %v 不该算作直连网段里的地址", ip)
		}
	}
}

func idleTransport(t *testing.T) *Transport {
	t.Helper()
	id, err := identity.LoadOrCreate(filepath.Join(t.TempDir(), "identity.json"))
	if err != nil {
		t.Fatal(err)
	}
	return New(Options{Identity: id, Port: -1, Discovery: NewMemoryNetwork().Discovery()})
}

// 新程序刚启动时被短暂拦下：马上再探测，探测成功就重启 mDNS，不报告被拦。
// 被拦期间发给先联系过本机的设备照样成功（算回包），不能当作拦截解除。
func TestBlockedStateFollowsProbe(t *testing.T) {
	tr := idleTransport(t)
	tr.probeAddr = netip.MustParseAddr("192.168.1.253")
	now := time.Now()
	tr.nextProbe = now.Add(probeEvery)

	tr.noteSend(syscall.EHOSTUNREACH, true, now)
	if tr.nextProbe.After(now) {
		t.Fatal("刚被拦时应当下一轮就再探测")
	}
	tr.noteSend(nil, false, now.Add(time.Second))
	if tr.rediscover || tr.deniedSince.IsZero() {
		t.Fatal("发给先联系过本机的设备成功了，不能当作拦截解除")
	}
	if !tr.blockedLocked(now.Add(blockGrace)) {
		t.Fatal("持续被拦过了宽限期应当报告")
	}
	tr.noteSend(nil, true, now.Add(2*time.Second))
	if !tr.rediscover {
		t.Fatal("探测成功后应当重启 mDNS")
	}
	if tr.blockedLocked(now.Add(blockGrace)) {
		t.Fatal("探测成功后不该再报告被拦")
	}
}

// 被拦超过宽限期后放慢探测；找不到探测地址时，任何一次发送成功都算拦截解除。
func TestBlockedWithoutProbeTarget(t *testing.T) {
	tr := idleTransport(t)
	now := time.Now()
	tr.nextProbe = now.Add(probeEvery)
	tr.noteSend(syscall.EHOSTUNREACH, false, now)
	tr.nextProbe = now.Add(blockGrace + probeEvery) // 宽限期末又探测了一次，probe 先排到 10 分钟后
	tr.noteSend(syscall.EHOSTUNREACH, true, now.Add(blockGrace))
	if want := now.Add(blockGrace + probeEveryBlocked); !tr.nextProbe.Equal(want) {
		t.Fatalf("宽限期过后应当每 %v 探测一次，下次探测在 %v，应为 %v", probeEveryBlocked, tr.nextProbe, want)
	}
	tr.noteSend(nil, false, now.Add(blockGrace+time.Second))
	if !tr.rediscover || !tr.deniedSince.IsZero() {
		t.Fatal("没有探测地址时，发送成功应当算拦截解除")
	}
}

package lan

import "testing"

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

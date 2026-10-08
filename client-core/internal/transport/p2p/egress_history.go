package p2p

import (
	"encoding/json"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// egressHistory 记住每个网络里见过的出口 IP，补上某一轮探测漏掉的线路。
//
// 探测服务器的 IP 恰好都没落到某条线路上时，这一轮就看不到它。若这条线路不改端口，
// 它在本端口上的公网地址就是「出口 IP + 本地端口」，可以直接推算出来发给对端。
// 会改端口的线路无法推算，不记。
type egressHistory struct {
	path string
	now  func() time.Time

	mu   sync.Mutex
	data historyFile
}

type historyFile struct {
	// 键是网络标识（见 localNetworkKey）：换了网络，出口自然也不一样了
	Networks map[string][]historyEntry `json:"networks"`
}

type historyEntry struct {
	IP            string    `json:"ip"`
	PortPreserved bool      `json:"port_preserved"`
	LastSeen      time.Time `json:"last_seen"`
}

// loadEgressHistory 读取历史；path 为空时只保存在内存里。文件损坏就当没有历史。
func loadEgressHistory(path string) *egressHistory {
	h := &egressHistory{path: path, now: time.Now, data: historyFile{Networks: map[string][]historyEntry{}}}
	if path == "" {
		return h
	}
	if raw, err := os.ReadFile(path); err == nil {
		var f historyFile
		if json.Unmarshal(raw, &f) == nil && f.Networks != nil {
			h.data = f
		}
	}
	return h
}

// record 把一轮实测结果并入历史，并清掉过期的条目。
func (h *egressHistory) record(network string, r EgressReport) {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := h.now()
	byIP := map[string]historyEntry{}
	for _, e := range h.data.Networks[network] {
		if now.Sub(e.LastSeen) < historyTTL {
			byIP[e.IP] = e
		}
	}
	for _, e := range r.Egresses {
		if e.Guessed() {
			continue // 推算出来的不能反过来当作见过的证据
		}
		ip := e.Addr.Addr().String()
		byIP[ip] = historyEntry{IP: ip, PortPreserved: e.PortPreserved, LastSeen: now}
	}

	entries := make([]historyEntry, 0, len(byIP))
	for _, e := range byIP {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].IP < entries[j].IP })
	h.data.Networks[network] = entries
	h.saveLocked()
}

// augment 给本轮结果补上历史里见过、这一轮没探测到、并且不改端口的出口。
func (h *egressHistory) augment(network string, r EgressReport) EgressReport {
	h.mu.Lock()
	defer h.mu.Unlock()

	present := map[netip.Addr]bool{}
	for _, e := range r.Egresses {
		present[e.Addr.Addr()] = true
	}
	now := h.now()
	for _, e := range h.data.Networks[network] {
		ip, err := netip.ParseAddr(e.IP)
		if err != nil || !e.PortPreserved || present[ip] || now.Sub(e.LastSeen) >= historyTTL {
			continue
		}
		r.Egresses = append(r.Egresses, Egress{
			Addr:          netip.AddrPortFrom(ip, uint16(r.LocalPort)),
			PortPreserved: true,
		})
	}
	return r
}

func (h *egressHistory) saveLocked() {
	if h.path == "" {
		return
	}
	raw, err := json.MarshalIndent(h.data, "", "  ")
	if err != nil {
		return
	}
	// 先写临时文件再改名，避免进程中途退出留下半个文件
	tmp := h.path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(h.path), 0o700); err != nil {
		return
	}
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return
	}
	_ = os.Rename(tmp, h.path)
}

// localNetworkKey 用本机各网卡的 IPv4 地址标识「当前所在的网络」。
// 换了 Wi-Fi、插拔网线，地址随之变化，历史与探测结果也就跟着换一套。
func localNetworkKey() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var addrs []string
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 || unstableAdapter(ifc.Name) {
			continue
		}
		list, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, a := range list {
			if ipn, ok := a.(*net.IPNet); ok && ipn.IP.To4() != nil {
				addrs = append(addrs, ipn.IP.String())
			}
		}
	}
	sort.Strings(addrs)
	return strings.Join(addrs, ",")
}

// unstableAdapter 判断一块网卡是否该排除在网络标识之外。
//
// Windows 上 WSL、Hyper-V、Docker 的虚拟网卡（名字以 vEthernet 开头）每次开机地址都会变，
// 算进来的话，同一个网络每次开机都像是新网络，出口历史就用不上了。其他平台不排除，
// 与之前的行为一致。
func unstableAdapter(name string) bool {
	return runtime.GOOS == "windows" && strings.HasPrefix(name, "vEthernet")
}

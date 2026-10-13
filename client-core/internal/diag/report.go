// Package diag 组装「反馈问题」用的诊断报告：版本与系统、设置、设备与网络的状态、最近的日志。
//
// 报告由用户看过之后自己贴进 GitHub 的 issue，那里是公开的，所以在这里先脱敏（见 redactor）。
// 复制的内容本来就不写进日志，这里不用再处理。报告的栏目名用英文，日志本身是什么语言就是什么语言。
package diag

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/config"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// Input 是生成报告要用的全部信息，由后台服务收集好传进来。
type Input struct {
	Now       time.Time
	Version   string // 后台服务的版本
	UIVersion string // 界面的版本，界面请求报告时带上
	Config    config.Config
	Status    *pb.Status
	Self      *pb.Device
	Peers     []*pb.Device
	Network   *pb.NetworkInfo // 还没探测过时为 nil

	// 详细日志打开到什么时候；零值表示没打开
	VerboseUntil time.Time

	LogFiles []string // 按时间先后，轮转出去的旧文件在前
	CrashLog string   // 未处理的 panic 写在这里；Mac 上没有，panic 直接进日志
	LogLines int

	Home      string   // 用户目录，换成 ~
	FileNames []string // 复制过的文件名，换成 <file>
}

// Build 生成脱敏后的报告文本。
func Build(in Input) string {
	// 整份报告最后统一脱敏；对齐成列的设备名先换掉，不然换完长短变了，列就歪了
	r := newRedactor(in)
	var b strings.Builder
	line := func(format string, a ...any) { fmt.Fprintf(&b, format+"\n", a...) }
	row := func(label string, value any) { line("  %-22s %v", label, value) }

	line("CopySync diagnostics report")
	line("Names, server address, public IP addresses, the user folder and file names are replaced.")
	line("")
	row("Generated", in.Now.Format("2006-01-02 15:04:05 -0700"))
	row("Background service", in.Version)
	if in.UIVersion != "" {
		row("App", in.UIVersion)
	}
	row("System", fmt.Sprintf("%s (%s/%s)", osVersion(), runtime.GOOS, runtime.GOARCH))

	c := in.Config
	line("")
	line("Settings")
	row("Server", or(c.SignalingURL, "not set"))
	row("Find devices on LAN", onOff(!c.LANDisabled))
	row("Public STUN servers", onOff(!c.OnlyOwnSTUN))
	row("Auto-sync limit", fmt.Sprintf("%d MB", c.AutoSyncThresholdBytes>>20))
	row("History kept", duration(time.Duration(c.HistoryTTL)))
	row("Files kept", duration(time.Duration(c.CacheTTL)))
	row("Sync", syncKinds(c))
	row("Language", or(c.Language, "follow system"))

	if st := in.Status; st != nil {
		line("")
		line("Status")
		conn := "not connected"
		if st.GetSignalingConnected() {
			conn = "connected"
		}
		row("Server connection", conn)
		row("LAN", lanState(st))
		row("Clipboard access", strings.ToLower(strings.TrimPrefix(st.GetClipboardPermission().String(), "CLIPBOARD_PERMISSION_")))
		row("Cache", fmt.Sprintf("%d MB", st.GetCacheBytesUsed()>>20))
	}
	if in.VerboseUntil.IsZero() {
		row("Detailed logging", "off")
	} else {
		row("Detailed logging", "on until "+in.VerboseUntil.Format("15:04"))
	}

	line("")
	line("Devices")
	if in.Self != nil {
		line("  %-22s %s (this device)", r.apply(in.Self.GetName()), in.Self.GetPlatform())
	}
	for _, p := range in.Peers {
		line("  %-22s %-8s %s", r.apply(p.GetName()), p.GetPlatform(),
			strings.ToLower(strings.TrimPrefix(p.GetConnection().String(), "CONNECTION_KIND_")))
	}

	if n := in.Network; n != nil {
		line("")
		line("Network")
		row("Hole-punching port", fmt.Sprintf("UDP %d", n.GetLocalPort()))
		probed := "not yet"
		if n.GetProbedAtUnix() > 0 {
			probed = fmt.Sprintf("%s, %d of %d servers answered",
				time.Unix(n.GetProbedAtUnix(), 0).Format("15:04:05"), n.GetAnswered(), n.GetProbed())
		}
		row("Last probe", probed)
		for _, e := range n.GetEgresses() {
			mode := "changes port"
			if e.GetPortPreserved() {
				mode = "keeps port"
			}
			via := strings.Join(e.GetVia(), ", ")
			if e.GetGuessed() {
				via = "guessed from history"
			}
			row("Egress", fmt.Sprintf("%s  %s  %s", e.GetAddress(), mode, via))
		}
		if lan := n.GetLan(); lan != nil && lan.GetEnabled() {
			row("LAN port", fmt.Sprintf("UDP %d, addresses %s", lan.GetPort(), strings.Join(lan.GetLocalAddresses(), " ")))
			row("Found on LAN", fmt.Sprintf("%d CopySync instances", lan.GetDiscovered()))
			for _, p := range lan.GetPeers() {
				state := "not found"
				if p.GetReachable() {
					state = "reachable at " + p.GetAddress()
				}
				if known := p.GetKnownAddresses(); len(known) > 0 {
					state += ", remembered " + strings.Join(known, " ")
				}
				line("  %-22s %s", r.apply(p.GetDeviceName()), state)
			}
		}
	}

	logs := tailLines(in.LogFiles, in.LogLines)
	line("")
	line("Log (%s)", lastLines(len(logs)))
	for _, l := range logs {
		line("%s", l)
	}
	if crash := tailLines([]string{in.CrashLog}, 100); in.CrashLog != "" && len(crash) > 0 {
		line("")
		line("Crash log (%s)", lastLines(len(crash)))
		for _, l := range crash {
			line("%s", l)
		}
	}
	return r.apply(b.String())
}

func lastLines(n int) string {
	if n == 1 {
		return "last line"
	}
	return fmt.Sprintf("last %d lines", n)
}

func or(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func onOff(on bool) string {
	if on {
		return "on"
	}
	return "off"
}

func lanState(st *pb.Status) string {
	switch {
	case !st.GetLanEnabled():
		return "off"
	case st.GetLanBlocked():
		return "blocked by the system"
	default:
		return "on"
	}
}

func syncKinds(c config.Config) string {
	var kinds []string
	for _, k := range []struct {
		on   bool
		name string
	}{{c.SyncText, "text"}, {c.SyncHTML, "rich text"}, {c.SyncImage, "images"}, {c.SyncFile, "files"}} {
		if k.on {
			kinds = append(kinds, k.name)
		}
	}
	if len(kinds) == 0 {
		return "nothing"
	}
	return strings.Join(kinds, ", ")
}

// duration 写成 3 days、12 hours 这样，设置里的时长都是整天或整小时
func duration(d time.Duration) string {
	plural := func(n time.Duration, unit string) string {
		if n == 1 {
			return "1 " + unit
		}
		return fmt.Sprintf("%d %ss", n, unit)
	}
	switch {
	case d >= 24*time.Hour && d%(24*time.Hour) == 0:
		return plural(d/(24*time.Hour), "day")
	case d >= time.Hour && d%time.Hour == 0:
		return plural(d/time.Hour, "hour")
	default:
		return d.String()
	}
}

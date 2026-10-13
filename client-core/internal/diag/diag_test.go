package diag

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/config"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

func sampleInput(t *testing.T, log string) Input {
	dir := t.TempDir()
	path := filepath.Join(dir, "daemon.log")
	if err := os.WriteFile(path, []byte(log), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default()
	cfg.DeviceName = "张三的 MacBook Pro"
	cfg.SignalingURL = "ws://sync.example.com:8787/signal"
	return Input{
		Now:       time.Date(2026, 10, 13, 20, 50, 0, 0, time.UTC),
		Version:   "1.7.0",
		Config:    cfg,
		Status:    &pb.Status{SignalingConnected: true, LanEnabled: true},
		Self:      &pb.Device{Id: "self", Name: "张三的 MacBook Pro", Platform: "darwin"},
		Peers:     []*pb.Device{{Id: "b", Name: "DESKTOP-LISI", Platform: "windows", Connection: pb.ConnectionKind_CONNECTION_KIND_DIRECT}},
		LogFiles:  []string{path},
		LogLines:  50,
		Home:      "/Users/zhangsan",
		FileNames: []string{"离职申请.docx", "a"},
	}
}

func TestReportRedacts(t *testing.T) {
	log := strings.Join([]string{
		`time=2026-10-13T20:41:02.123+08:00 level=INFO msg=配对成功 peer="DESKTOP-LISI" device=b`,
		`time=2026-10-13T20:41:03.000+08:00 level=INFO msg=信令已连接 url=ws://sync.example.com:8787/signal`,
		`time=2026-10-13T20:41:04.000+08:00 level=WARN msg=解包失败 err="open /Users/zhangsan/Downloads/离职申请.docx: permission denied"`,
		`time=2026-10-13T20:41:05.000+08:00 level=DEBUG msg=候选 addr=203.0.113.7:51234 lan=192.168.1.20:8788 ts=100.88.20.5`,
		`time=2026-10-13T20:41:06.000+08:00 level=DEBUG msg=候选 v6=2001:db8::1 local=fe80::1 mac=aa:bb:cc:dd:ee:ff`,
	}, "\n") + "\n"
	report := Build(sampleInput(t, log))

	for _, secret := range []string{"张三", "LISI", "sync.example.com", "zhangsan", "离职申请", "203.0.113.7", "2001:db8::1"} {
		if strings.Contains(report, secret) {
			t.Errorf("报告里还有 %q：\n%s", secret, report)
		}
	}
	for _, kept := range []string{
		"<this-device>", "<device-1>", "ws://<server>:8787/signal", "~/Downloads/<file>", "<public-ip>:51234",
		"192.168.1.20:8788", "100.88.20.5", "fe80::1", "aa:bb:cc:dd:ee:ff", "20:41:06.000+08:00",
		"Log (last 5 lines)", "1.7.0",
	} {
		if !strings.Contains(report, kept) {
			t.Errorf("报告里应当有 %q：\n%s", kept, report)
		}
	}
	// 太短的文件名不替换，免得把日志里的普通字词也换掉
	if !strings.Contains(report, `err="open ~/Downloads/<file>: permission denied"`) {
		t.Errorf("路径没有按预期脱敏：\n%s", report)
	}
	if strings.Count(report, "<file>") != 1 {
		t.Errorf("只有那个文件名该被替换：\n%s", report)
	}
}

func TestTailLines(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "daemon.log.1")
	cur := filepath.Join(dir, "daemon.log")
	_ = os.WriteFile(old, []byte("old 1\nold 2\nold 3\n"), 0o600)
	_ = os.WriteFile(cur, []byte("new 1\r\nnew 2\n"), 0o600)

	// 最新的文件不够时，从轮转出去的旧文件末尾补，顺序不乱
	got := strings.Join(tailLines([]string{old, cur, filepath.Join(dir, "missing")}, 4), "|")
	if got != "old 2|old 3|new 1|new 2" {
		t.Errorf("tailLines = %q", got)
	}

	// 超长的一行截断
	long := filepath.Join(dir, "long.log")
	_ = os.WriteFile(long, []byte(strings.Repeat("长", 1000)+"\n"), 0o600)
	if l := tailLines([]string{long}, 1)[0]; len([]rune(l)) != maxLineRunes+1 || !strings.HasSuffix(l, "…") {
		t.Errorf("长行没有截断：%d 个字", len([]rune(l)))
	}

	// 只读末尾一段时，读到一半的第一行丢掉
	big := filepath.Join(dir, "big.log")
	_ = os.WriteFile(big, []byte(strings.Repeat("x", tailBytes)+"\nlast\n"), 0o600)
	if got := tailLines([]string{big}, 5); len(got) != 1 || got[0] != "last" {
		t.Errorf("大文件的末尾：%q", got)
	}
}

func TestDuration(t *testing.T) {
	for d, want := range map[time.Duration]string{
		72 * time.Hour: "3 days", 24 * time.Hour: "1 day", 12 * time.Hour: "12 hours", 90 * time.Minute: "1h30m0s",
	} {
		if got := duration(d); got != want {
			t.Errorf("duration(%v) = %q，应为 %q", d, got, want)
		}
	}
}

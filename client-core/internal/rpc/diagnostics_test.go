package rpc

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/config"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

func TestExportDiagnostics(t *testing.T) {
	log := filepath.Join(t.TempDir(), "daemon.log")
	var lines []string
	for i := range 300 {
		lines = append(lines, "level=INFO msg=第"+string(rune('0'+i%10))+"行 peer=\"李四的 ThinkPad\"")
	}
	if err := os.WriteFile(log, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default()
	cfg.DeviceName = "张三的 MacBook Pro"
	s := NewServer(Deps{
		Version:    "1.7.0",
		LoadConfig: func() config.Config { return cfg },
		Devices: func() (*pb.Device, []*pb.Device) {
			return &pb.Device{Id: "a", Name: cfg.DeviceName}, []*pb.Device{{Id: "b", Name: "李四的 ThinkPad"}}
		},
		LogFiles: func() ([]string, string) { return []string{log}, "" },
	}, NewHub(), "test-token")

	got, err := s.ExportDiagnostics(context.Background(), &pb.ExportDiagnosticsRequest{LogLines: 20, UiVersion: "1.7.0"})
	if err != nil {
		t.Fatal(err)
	}
	r := got.GetReport()
	if strings.Contains(r, "张三") || strings.Contains(r, "李四") {
		t.Errorf("设备名没有换掉：\n%s", r)
	}
	if !strings.Contains(r, "Log (last 20 lines)") || strings.Count(r, "<device-1>") < 20 {
		t.Errorf("日志行数或代号不对：\n%s", r)
	}
}

func TestSetVerboseLog(t *testing.T) {
	var asked time.Duration
	until := time.Now().Add(30 * time.Minute)
	s := NewServer(Deps{
		LoadConfig:   config.Default,
		VerboseLog:   func(d time.Duration) time.Time { asked = d; return until },
		VerboseUntil: func() time.Time { return until },
	}, NewHub(), "test-token")

	st, err := s.SetVerboseLog(context.Background(), &pb.SetVerboseLogRequest{Minutes: 600})
	if err != nil {
		t.Fatal(err)
	}
	if asked != 2*time.Hour {
		t.Errorf("没有限制在两小时内：%v", asked)
	}
	if st.GetVerboseLogUntilUnix() != until.Unix() {
		t.Errorf("状态里没有详细日志的截止时间：%d", st.GetVerboseLogUntilUnix())
	}
}

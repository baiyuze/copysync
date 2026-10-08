package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadCreatesDefaultsWhenMissing(t *testing.T) {
	p := PathsUnder(t.TempDir())
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("加载: %v", err)
	}
	if cfg.AutoSyncThresholdBytes != 50<<20 {
		t.Errorf("默认阈值 = %d，期望 50 MiB", cfg.AutoSyncThresholdBytes)
	}
	// 首次加载应把默认配置写盘，用户才能看到有哪些可调项
	if _, err := os.Stat(p.Config); err != nil {
		t.Errorf("默认配置未落盘: %v", err)
	}
}

// 新增配置项时，老配置文件里没有该字段，必须继承默认值而不是零值。
func TestLoadInheritsDefaultsForMissingFields(t *testing.T) {
	p := PathsUnder(t.TempDir())
	if err := p.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p.Config,
		[]byte(`{"device_name":"我的电脑"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("加载: %v", err)
	}
	if cfg.DeviceName != "我的电脑" {
		t.Errorf("设备名 = %q", cfg.DeviceName)
	}
	if cfg.AutoSyncThresholdBytes != 50<<20 {
		t.Errorf("缺失字段未继承默认值: %d", cfg.AutoSyncThresholdBytes)
	}
	if !cfg.SyncText {
		t.Error("缺失的布尔项未继承默认值 true")
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	p := PathsUnder(t.TempDir())
	want := Default()
	want.DeviceName = "测试设备"
	want.SignalingURL = "wss://example.com/signal"
	want.HistoryTTL = Duration(6 * time.Hour)

	if err := Save(p, want); err != nil {
		t.Fatalf("保存: %v", err)
	}
	got, err := Load(p)
	if err != nil {
		t.Fatalf("加载: %v", err)
	}
	if got.DeviceName != want.DeviceName || got.SignalingURL != want.SignalingURL {
		t.Errorf("往返后不一致: %+v", got)
	}
	if got.HistoryTTL.Std() != 6*time.Hour {
		t.Errorf("时长往返失败: %v", got.HistoryTTL.Std())
	}
}

// 时长写成 "72h" 而非纳秒整数——配置文件是给人读和改的。
func TestDurationIsHumanReadableInJSON(t *testing.T) {
	data, err := json.Marshal(Duration(72 * time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"72h0m0s"` {
		t.Errorf("序列化结果 = %s，期望可读字符串", data)
	}

	p := PathsUnder(t.TempDir())
	if err := Save(p, Default()); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(p.Config)
	if strings.Contains(string(raw), "259200000000000") {
		t.Error("配置文件里出现了纳秒整数，应为可读时长")
	}
}

func TestDurationAcceptsSeconds(t *testing.T) {
	// 手写配置时可能直接写数字，按秒解释比报错友好
	var d Duration
	if err := json.Unmarshal([]byte(`90`), &d); err != nil {
		t.Fatalf("解析数字失败: %v", err)
	}
	if d.Std() != 90*time.Second {
		t.Errorf("= %v，期望 90s", d.Std())
	}

	if err := json.Unmarshal([]byte(`"bad"`), &d); err == nil {
		t.Error("非法时长字符串应当报错")
	}
}

func TestValidateClampsBadValues(t *testing.T) {
	c := Config{
		DeviceName:             "",
		SignalingURL:           "",
		AutoSyncThresholdBytes: -1,
		HistoryTTL:             Duration(time.Second), // 过短
		CacheTTL:               Duration(0),
	}
	c.Validate()

	if c.DeviceName == "" || c.SignalingURL == "" {
		t.Error("空值未被补为默认值")
	}
	if c.AutoSyncThresholdBytes <= 0 {
		t.Errorf("负阈值未被修正: %d", c.AutoSyncThresholdBytes)
	}
	// TTL 过短会让 GC 立刻删掉刚写入的记录
	if c.HistoryTTL.Std() < time.Minute {
		t.Errorf("过短的 history_ttl 未被修正: %v", c.HistoryTTL.Std())
	}
}

// 缓存活得比记录久没有意义：记录没了就再也找不到那些文件。
func TestValidateCapsCacheTTLByHistoryTTL(t *testing.T) {
	c := Default()
	c.HistoryTTL = Duration(time.Hour)
	c.CacheTTL = Duration(24 * time.Hour)
	c.Validate()

	if c.CacheTTL > c.HistoryTTL {
		t.Errorf("cache_ttl (%v) 不应超过 history_ttl (%v)",
			c.CacheTTL.Std(), c.HistoryTTL.Std())
	}
}

func TestLoadRejectsCorruptFile(t *testing.T) {
	p := PathsUnder(t.TempDir())
	if err := p.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p.Config, []byte("{ 这不是 json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(p); err == nil {
		t.Error("损坏的配置文件应当报错，而不是静默用默认值覆盖用户设置")
	}
}

func TestPathsLayout(t *testing.T) {
	root := t.TempDir()
	p := PathsUnder(root)

	if p.Cache != filepath.Join(root, "cache") {
		t.Errorf("缓存目录 = %q", p.Cache)
	}
	if filepath.Base(p.Endpoint) != "daemon.json" {
		t.Errorf("endpoint 文件名 = %q，UI 依赖这个约定", p.Endpoint)
	}
	if err := p.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(p.Cache)
	if err != nil {
		t.Fatal(err)
	}
	// 缓存里是用户复制过的文件，不该让同机其他用户读到
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Errorf("缓存目录权限 = %o，期望 0700", perm)
	}
}

func TestSaveWritesPrivatePermissions(t *testing.T) {
	p := PathsUnder(t.TempDir())
	if err := Save(p, Default()); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(p.Config)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("配置文件权限 = %o，期望 0600", perm)
	}
}

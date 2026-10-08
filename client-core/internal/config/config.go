// Package config 管理 daemon 的持久化配置与数据目录布局。
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Config 是用户可调的全部设置。字段与 proto 的 copysync.v1.Config 一一对应。
type Config struct {
	DeviceID     string `json:"device_id"`
	DeviceName   string `json:"device_name"`
	SignalingURL string `json:"signaling_url"`

	// 小于此值的内容在复制时自动推送给对端；大于则只同步元数据，
	// 等用户在 UI 里手动拉取。这是 M-1 否决纯懒模式后的折中点。
	AutoSyncThresholdBytes int64 `json:"auto_sync_threshold_bytes"`

	HistoryTTL Duration `json:"history_ttl"` // 历史记录元数据保留期
	CacheTTL   Duration `json:"cache_ttl"`   // 已落盘内容的保留期

	SyncText  bool `json:"sync_text"`
	SyncHTML  bool `json:"sync_html"`
	SyncImage bool `json:"sync_image"`
	SyncFile  bool `json:"sync_file"`

	AutoApplyToClipboard bool `json:"auto_apply_to_clipboard"`
	LaunchAtLogin        bool `json:"launch_at_login"`

	// OnlyOwnSTUN 只用自己的服务器探测网络出口。公共 STUN 服务器能看到本机的公网 IP
	// （与任何 WebRTC 应用一样，看不到内容），介意的用户可以关掉；代价是多出口网络里
	// 可能探测不全，直连机会变小。
	OnlyOwnSTUN bool `json:"only_own_stun"`
}

func Default() Config {
	host, _ := os.Hostname()
	if host == "" {
		host = "未命名设备"
	}
	return Config{
		DeviceName:             host,
		SignalingURL:           "ws://127.0.0.1:8787/signal",
		AutoSyncThresholdBytes: 50 << 20, // 50 MiB
		HistoryTTL:             Duration(3 * 24 * time.Hour),
		CacheTTL:               Duration(3 * 24 * time.Hour),
		SyncText:               true,
		SyncHTML:               true,
		SyncImage:              true,
		SyncFile:               true,
		AutoApplyToClipboard:   true,
		LaunchAtLogin:          false,
	}
}

// Validate 修正越界值，避免用户或损坏的配置文件把 daemon 带进坏状态。
func (c *Config) Validate() {
	d := Default()
	if c.DeviceName == "" {
		c.DeviceName = d.DeviceName
	}
	if c.SignalingURL == "" {
		c.SignalingURL = d.SignalingURL
	}
	if c.AutoSyncThresholdBytes <= 0 {
		c.AutoSyncThresholdBytes = d.AutoSyncThresholdBytes
	}
	// TTL 至少一分钟，否则 GC 会把刚写入的记录立刻删掉
	if time.Duration(c.HistoryTTL) < time.Minute {
		c.HistoryTTL = d.HistoryTTL
	}
	if time.Duration(c.CacheTTL) < time.Minute {
		c.CacheTTL = d.CacheTTL
	}
	// 缓存活得比历史记录久没有意义：记录没了就找不到那些文件了
	if c.CacheTTL > c.HistoryTTL {
		c.CacheTTL = c.HistoryTTL
	}
}

// ─────────────────────────── 目录布局 ───────────────────────────

// Paths 是 daemon 用到的全部路径。集中在一处，便于测试时重定向。
type Paths struct {
	Root   string // 配置与数据库
	Cache  string // 收到的文件内容
	Config string
	DB     string
	// daemon 启动后把 gRPC 端口和一次性 token 写在这里，供 UI 发现
	Endpoint string
}

func DefaultPaths() (Paths, error) {
	var root string
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return Paths{}, err
		}
		root = filepath.Join(home, "Library", "Application Support", "CopySync")
	case "windows":
		base := os.Getenv("LOCALAPPDATA")
		if base == "" {
			var err error
			if base, err = os.UserConfigDir(); err != nil {
				return Paths{}, err
			}
		}
		root = filepath.Join(base, "CopySync")
	default:
		base, err := os.UserConfigDir()
		if err != nil {
			return Paths{}, err
		}
		root = filepath.Join(base, "copysync")
	}
	return pathsUnder(root), nil
}

func pathsUnder(root string) Paths {
	return Paths{
		Root:     root,
		Cache:    filepath.Join(root, "cache"),
		Config:   filepath.Join(root, "config.json"),
		DB:       filepath.Join(root, "copysync.db"),
		Endpoint: filepath.Join(root, "daemon.json"),
	}
}

// PathsUnder 让测试把整套目录指到临时位置。
func PathsUnder(root string) Paths { return pathsUnder(root) }

func (p Paths) EnsureDirs() error {
	for _, d := range []string{p.Root, p.Cache} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return fmt.Errorf("创建目录 %s: %w", d, err)
		}
	}
	return nil
}

// ─────────────────────────── 读写 ───────────────────────────

// Load 读取配置；文件不存在时写入并返回默认配置。
func Load(p Paths) (Config, error) {
	data, err := os.ReadFile(p.Config)
	if os.IsNotExist(err) {
		c := Default()
		if err := Save(p, c); err != nil {
			return c, err
		}
		return c, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("读取配置: %w", err)
	}

	c := Default() // 以默认值打底，缺失字段自动继承
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("解析配置 %s: %w", p.Config, err)
	}
	c.Validate()
	return c, nil
}

func Save(p Paths, c Config) error {
	if err := p.EnsureDirs(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	// 先写临时文件再改名，避免写一半断电留下损坏的配置
	tmp := p.Config + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p.Config)
}

// ─────────────────────────── Duration ───────────────────────────

// Duration 让 time.Duration 在 JSON 里表现为 "72h" 这样的可读字符串，
// 而不是一长串纳秒整数——配置文件是给人看的。
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		parsed, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("无法解析时长 %q: %w", value, err)
		}
		*d = Duration(parsed)
	case float64: // 兼容手写成秒数的情况
		*d = Duration(time.Duration(value) * time.Second)
	default:
		return fmt.Errorf("时长字段类型不支持: %T", v)
	}
	return nil
}

func (d Duration) Std() time.Duration { return time.Duration(d) }

// copysyncd 是 CopySync 的常驻后台进程，承载全部核心逻辑：
// 剪贴板监听、设备配对、P2P 传输、本地缓存与回收。
//
// Flutter UI 只是控制台——关掉 UI 不影响后台同步。UI 通过
// 本地 gRPC（127.0.0.1 随机端口 + 一次性 token）连接本进程。
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/cache"
	"github.com/baiyuze/copysync/client-core/internal/clipboard"
	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/gc"
	"github.com/baiyuze/copysync/client-core/internal/identity"
	"github.com/baiyuze/copysync/client-core/internal/peers"
	"github.com/baiyuze/copysync/client-core/internal/rpc"
	"github.com/baiyuze/copysync/client-core/internal/store"
	syncengine "github.com/baiyuze/copysync/client-core/internal/sync"
	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

var version = "0.1.0-dev"

func main() {
	// 必须在任何其他 goroutine 启动前锁定主线程。
	//
	// macOS 的剪贴板 API 在隐私机制开启时会经 ViewBridge 与系统 UI 服务通信，
	// 其内部 NSCFRunLoopSemaphore 依赖 run loop；在非主线程调用会直接
	// SIGSEGV（M-1 实测，见 spikes/results.md）。所以主线程被专门留出来
	// 跑 MainLoop，业务逻辑一律在其他 goroutine 里，通过 Post/Call 投递。
	runtime.LockOSThread()

	var (
		dataDir = flag.String("data-dir", "", "数据目录（默认为系统标准位置）")
		verbose = flag.Bool("v", false, "输出调试日志")
		probe   = flag.Bool("probe-clipboard", false,
			"诊断用：在主线程直接读一次剪贴板并打印结果后退出")
	)
	flag.Parse()

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	if *probe {
		probeClipboard()
		return
	}

	if err := run(*dataDir); err != nil {
		slog.Error("daemon 退出", "err", err)
		os.Exit(1)
	}
}

// probeClipboard 在主线程直接走一遍探测与读取，用来把剪贴板问题
// 与 daemon 的其余部分隔离开。
func probeClipboard() {
	clipboard.InitPlatform()
	cb := clipboard.New()

	fmt.Println("授权状态:", cb.Permission())

	snap, err := cb.Peek()
	if err != nil {
		fmt.Println("Peek 失败:", err)
		return
	}
	fmt.Printf("changeCount: %d\n分类: %v\ncontentType: %q\ntypes: %v\n",
		snap.ChangeCount, snap.Kind, snap.ContentType, snap.Types)

	content, err := cb.Read()
	if err != nil {
		fmt.Println("Read 失败:", err)
		return
	}
	fmt.Printf("读取结果: kind=%v 文本=%d字符 HTML=%d字符 图片=%d字节 文件=%v\n",
		content.Kind, len([]rune(content.Text)), len(content.HTML),
		len(content.Image), content.Files)
	if content.Text != "" {
		fmt.Printf("文本内容: %q\n", content.Text)
	}
}

func run(dataDir string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ── 路径与配置 ──
	paths, err := resolvePaths(dataDir)
	if err != nil {
		return err
	}
	if err := paths.EnsureDirs(); err != nil {
		return err
	}
	cfg, err := config.Load(paths)
	if err != nil {
		return err
	}

	// ── 设备身份 ──
	id, err := identity.LoadOrCreate(filepath.Join(paths.Root, "identity.json"))
	if err != nil {
		return err
	}
	if cfg.DeviceID != id.DeviceID {
		cfg.DeviceID = id.DeviceID
		if err := config.Save(paths, cfg); err != nil {
			return err
		}
	}
	slog.Info("设备身份就绪",
		"device_id", id.DeviceID, "name", cfg.DeviceName, "fingerprint", id.Fingerprint())

	// ── 存储与缓存 ──
	db, err := store.Open(paths.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	blobs, err := cache.New(paths.Cache)
	if err != nil {
		return err
	}
	// 上次异常退出可能留下半截的临时文件
	if err := blobs.CleanTemp(); err != nil {
		slog.Debug("清理临时文件失败", "err", err)
	}

	// ── 配置的并发安全读写 ──
	var cfgMu sync.RWMutex
	loadCfg := func() config.Config {
		cfgMu.RLock()
		defer cfgMu.RUnlock()
		return cfg
	}
	saveCfg := func(c config.Config) error {
		cfgMu.Lock()
		defer cfgMu.Unlock()
		if err := config.Save(paths, c); err != nil {
			return err
		}
		cfg = c
		return nil
	}

	hub := rpc.NewHub()

	// ── 剪贴板（必须在主线程使用）──
	clipboard.InitPlatform()
	loop := clipboard.NewMainLoop()
	watcher := clipboard.NewWatcher(clipboard.WatcherOptions{
		Clipboard: clipboard.New(),
		MainLoop:  loop,
		OnPermission: func(p clipboard.Permission) {
			hub.Publish(rpc.EventPermissionChanged(permissionToProto(p)))
		},
	})

	// ── 设备配对与信令（引擎稍后注入）──
	var engine *syncengine.Engine

	// 状态快照要由 gRPC 服务组装，而它晚于 peerMgr 创建，这里先记下"状态变了"，
	// 等服务就绪后统一推送。容量 1 让连续的变化合并成一次刷新。
	statusDirty := make(chan struct{}, 1)

	peerMgr, err := peers.NewManager(peers.Options{
		Identity:   id,
		Store:      db,
		LoadConfig: loadCfg,
		OnDeviceChanged: func(d *pb.Device) {
			hub.Publish(rpc.EventDeviceChanged(d))
		},
		OnStatusChanged: func() {
			select {
			case statusDirty <- struct{}{}:
			default:
			}
		},
		OnPeerMessage: func(deviceID string, msg *pb.PeerMessage) {
			if engine != nil {
				engine.HandlePeerMessage(ctx, deviceID, msg)
			}
		},
		OnPeerStream: func(deviceID, label string, stream *p2p.Stream) {
			if engine != nil {
				engine.HandlePeerStream(deviceID, label, stream)
			}
		},
	})
	if err != nil {
		return err
	}

	// ── 同步引擎 ──
	engine = syncengine.New(syncengine.Options{
		Store:      db,
		Cache:      blobs,
		Watcher:    watcher,
		DeviceID:   id.DeviceID,
		DeviceName: func() string { return loadCfg().DeviceName },
		LoadConfig: loadCfg,
		Send:       peerMgr.P2P().Send,
		Broadcast:  peerMgr.P2P().Broadcast,
		OpenStream: peerMgr.P2P().OpenStream,
		OnlinePeers: func() []string {
			_, list := peerMgr.Devices()
			var online []string
			for _, d := range list {
				if d.GetOnline() {
					online = append(online, d.GetId())
				}
			}
			return online
		},
		OnRecord: func(r *pb.ClipRecord) {
			hub.Publish(rpc.EventClipAdded(r))
		},
		OnProgress: func(p *pb.TransferProgress) {
			hub.Publish(rpc.EventProgress(p))
		},
	})

	// 剪贴板变化在独立 goroutine 里处理：读内容可能耗时，
	// 不能拖住轮询循环，更不能占着主线程。
	watcher.SetOnChange(func(snap clipboard.Snapshot) {
		go engine.HandleClipboardChange(ctx, snap)
	})

	go peerMgr.Run(ctx)
	go watcher.Run(ctx)

	// ── 过期回收 ──
	collector := gc.New(gc.Options{
		Store:      db,
		Cache:      blobs,
		LoadConfig: loadCfg,
		OnRemoved: func(ids []string) {
			for _, cid := range ids {
				hub.Publish(rpc.EventClipRemoved(cid))
			}
		},
	})
	go collector.Run(ctx)

	// 改了服务器地址要立即重连，否则用户得重启 daemon 才生效
	saveCfgAndApply := func(c config.Config) error {
		prev := loadCfg()
		if err := saveCfg(c); err != nil {
			return err
		}
		if c.SignalingURL != prev.SignalingURL {
			slog.Info("信令服务器地址已变更", "from", prev.SignalingURL, "to", c.SignalingURL)
			peerMgr.SetSignalingURL(c.SignalingURL)
		}
		return nil
	}

	// ── gRPC 服务 ──
	token, err := rpc.NewToken()
	if err != nil {
		return err
	}
	srv := rpc.NewServer(rpc.Deps{
		Store:              db,
		Paths:              paths,
		Version:            version,
		DeviceID:           id.DeviceID,
		LoadConfig:         loadCfg,
		SaveConfig:         saveCfgAndApply,
		Devices:            peerMgr.Devices,
		CreatePairingCode:  peerMgr.CreatePairingCode,
		RedeemPairingCode:  peerMgr.RedeemPairingCode,
		ConfirmPairing:     peerMgr.ConfirmPairing,
		Unpair:             peerMgr.Unpair,
		SignalingConnected: peerMgr.Connected,
		Fetch:              engine.Fetch,
		ApplyToClipboard:   engine.ApplyToClipboard,
		CacheBytesUsed:     blobs.Size,
		Permission: func() clipboard.Permission {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			return watcher.Permission(ctx)
		},
		OpenPermissionSettings: func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return watcher.OpenPermissionSettings(ctx)
		},
	}, hub, token)

	// 信令连接与在线对端数的变化推给 UI。Subscribe 只在建立订阅时推一次状态，
	// 没有这里的话，UI 会一直停留在它打开那一刻的连接状态。
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-statusDirty:
				if st, err := srv.GetStatus(ctx, &pb.Empty{}); err == nil {
					hub.Publish(rpc.EventStatusChanged(st))
				}
			}
		}
	}()

	port, grpcServer, err := srv.Serve(ctx)
	if err != nil {
		return err
	}
	defer grpcServer.GracefulStop()

	// UI 靠这个文件发现 daemon
	if err := rpc.WriteEndpoint(paths.Endpoint, rpc.Endpoint{
		Port: port, Token: token, PID: os.Getpid(), Version: version,
	}); err != nil {
		return err
	}
	defer rpc.RemoveEndpoint(paths.Endpoint)

	slog.Info("daemon 已就绪",
		"port", port, "endpoint", paths.Endpoint, "cache", paths.Cache)

	// ── 主线程交给 MainLoop ──
	// 剪贴板相关的一切都在这条线程上执行，同时驱动 macOS run loop。
	loop.Run(ctx)

	slog.Info("正在退出")
	return nil
}

func resolvePaths(dataDir string) (config.Paths, error) {
	if dataDir != "" {
		abs, err := filepath.Abs(dataDir)
		if err != nil {
			return config.Paths{}, err
		}
		return config.PathsUnder(abs), nil
	}
	return config.DefaultPaths()
}

func permissionToProto(p clipboard.Permission) pb.ClipboardPermission {
	switch p {
	case clipboard.PermissionNotApplicable:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_NOT_APPLICABLE
	case clipboard.PermissionDefault:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_DEFAULT
	case clipboard.PermissionAsk:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_ASK
	case clipboard.PermissionAlwaysAllow:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_ALWAYS_ALLOW
	case clipboard.PermissionAlwaysDeny:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_ALWAYS_DENY
	default:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_UNSPECIFIED
	}
}

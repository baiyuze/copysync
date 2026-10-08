// copysync-server 是 CopySync 的信令服务器。
//
// 它只做三件事：校验连上来的设备确实持有其声称公钥的私钥、
// 转达配对码、在设备间转发带签名的信令。它不存用户，
// 也读不懂设备间传输的内容——P2P 打通后数据根本不经过它，
// 打不通时走 TURN 中转，而 TURN 只转发 UDP 包，看不到 DTLS 内层。
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	ossignal "os/signal" // 与内部的 signal 包重名，加别名区分
	"strings"
	"syscall"
	"time"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
	"github.com/baiyuze/copysync/server/relay"
	"github.com/baiyuze/copysync/server/signal"
)

var version = "0.1.0-dev"

func main() {
	var (
		addr     = flag.String("addr", ":8787", "HTTP 监听地址")
		stunURLs = flag.String("stun", "stun:stun.l.google.com:19302",
			"STUN 服务器，逗号分隔")
		turnIP = flag.String("turn-ip", "",
			"TURN 对外公布的公网地址。留空则不启用中转，P2P 打不通时将无法传输")
		turnPort   = flag.Int("turn-port", 3478, "TURN 监听端口")
		turnSecret = flag.String("turn-secret", "",
			"TURN 凭证派生密钥。留空自动生成（重启后已签发的凭证失效）")
		verbose     = flag.Bool("v", false, "输出调试日志")
		showVersion = flag.Bool("version", false, "打印版本号后退出")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println("copysync-server", version)
		return
	}

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(log)

	cfg := serverConfig{
		addr:       *addr,
		stunURLs:   *stunURLs,
		turnIP:     *turnIP,
		turnPort:   *turnPort,
		turnSecret: *turnSecret,
	}
	if err := run(cfg, log); err != nil {
		log.Error("服务器退出", "err", err)
		os.Exit(1)
	}
}

type serverConfig struct {
	addr       string
	stunURLs   string
	turnIP     string
	turnPort   int
	turnSecret string
}

func run(cfg serverConfig, log *slog.Logger) error {
	ctx, stop := ossignal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	stun := splitNonEmpty(cfg.stunURLs)

	// TURN 中转：P2P 打洞失败时兜底。它只转发 UDP 包，看不到 DTLS 内层明文。
	var turnSrv *relay.Server
	if cfg.turnIP != "" {
		secret := cfg.turnSecret
		if secret == "" {
			buf := make([]byte, 32)
			if _, err := rand.Read(buf); err != nil {
				return err
			}
			secret = hex.EncodeToString(buf)
			log.Warn("未指定 -turn-secret，已自动生成；重启后先前签发的凭证会失效")
		}
		var err error
		turnSrv, err = relay.New(relay.Options{
			PublicIP: cfg.turnIP,
			Port:     cfg.turnPort,
			Secret:   secret,
			Logger:   log,
		})
		if err != nil {
			return err
		}
		defer turnSrv.Close()
	} else {
		log.Warn("未配置 -turn-ip，仅提供 STUN；对称 NAT 环境下设备可能无法互通")
	}

	sig := signal.NewServer(signal.Options{
		Logger: log,
		// 每台设备拿到的 TURN 凭证都是按其 device_id 单独签发、带有效期的
		ICEServers: func(deviceID string) []*pb.IceServer {
			var servers []*pb.IceServer
			if len(stun) > 0 {
				servers = append(servers, &pb.IceServer{Urls: stun})
			}
			if turnSrv != nil {
				c := turnSrv.Issue(deviceID)
				servers = append(servers, &pb.IceServer{
					Urls:          c.URLs,
					Username:      c.Username,
					Credential:    c.Password,
					ExpiresAtUnix: c.ExpiresAt.Unix(),
				})
			}
			return servers
		},
	})

	mux := http.NewServeMux()
	mux.Handle("/signal", sig)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","version":"` + version + `"}`))
	})

	srv := &http.Server{
		Addr:    cfg.addr,
		Handler: mux,
		// WebSocket 是长连接，不能设 WriteTimeout；
		// 空闲与读头超时仍然保留，防止慢速连接耗尽资源。
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Info("信令服务器已启动", "addr", cfg.addr, "stun", stun,
		"turn", cfg.turnIP != "", "version", version)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	log.Info("已停止")
	return nil
}

func splitNonEmpty(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

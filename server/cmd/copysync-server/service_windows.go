package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
)

// 与 server/deploy/install.ps1 注册的服务名一致
const serviceName = "CopySyncServer"

// isService 判断是否由 Windows 服务管理器启动。
func isService() bool {
	ok, err := svc.IsWindowsService()
	return err == nil && ok
}

func defaultServiceLog() string {
	return filepath.Join(os.Getenv("ProgramData"), "CopySync Server", "server.log")
}

// runService 在服务管理器下运行，收到停止或关机通知时退出。
func runService(cfg serverConfig, log *slog.Logger) error {
	return svc.Run(serviceName, &service{cfg: cfg, log: log})
}

type service struct {
	cfg serverConfig
	log *slog.Logger
}

func (s *service) Execute(_ []string, req <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- run(ctx, s.cfg, s.log) }()
	status <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	for {
		select {
		case err := <-done:
			if err != nil {
				// 端口被占用之类：以非零退出码结束，服务管理器按安装时设置的恢复策略重启
				s.log.Error("服务器退出", "err", err)
				return true, 1
			}
			return false, 0
		case c := <-req:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				status <- svc.Status{State: svc.StopPending}
				cancel()
				<-done
				return false, 0
			}
		}
	}
}

// copysync-cli 是开发调试用的命令行客户端。
//
// 它走与 Flutter UI 完全相同的 gRPC 接口，因此既能在 UI 就绪前验证 daemon，
// 也能在之后用来区分「是 daemon 的问题还是 UI 的问题」。
//
//	copysync-cli status      查看 daemon 状态
//	copysync-cli watch       订阅事件流
//	copysync-cli history     列出历史记录
//	copysync-cli config      查看配置
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/rpc"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

func main() {
	dataDir := flag.String("data-dir", "", "数据目录（默认为系统标准位置）")
	flag.Parse()

	cmd := flag.Arg(0)
	if cmd == "" {
		fmt.Fprintln(os.Stderr,
			"用法: copysync-cli [-data-dir DIR] <status|watch|history|config|devices|pair|join CODE>")
		os.Exit(2)
	}

	if err := run(cmd, *dataDir); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}

func run(cmd, dataDir string) error {
	paths, err := resolvePaths(dataDir)
	if err != nil {
		return err
	}

	ep, err := rpc.ReadEndpoint(paths.Endpoint)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("daemon 未运行（找不到 %s）", paths.Endpoint)
		}
		return err
	}

	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", ep.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("连接 daemon: %w", err)
	}
	defer conn.Close()

	client := pb.NewDaemonServiceClient(conn)
	// token 随每次调用发送，daemon 侧用常数时间比较校验
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-copysync-token", ep.Token)

	switch cmd {
	case "status":
		return showStatus(ctx, client, ep)
	case "watch":
		return watch(ctx, client)
	case "history":
		return showHistory(ctx, client)
	case "config":
		return showConfig(ctx, client)
	case "devices":
		return showDevices(ctx, client)
	case "pair":
		return createPairing(ctx, client)
	case "join":
		return joinPairing(ctx, client, flag.Arg(1))
	default:
		return fmt.Errorf("未知命令 %q", cmd)
	}
}

func showDevices(ctx context.Context, c pb.DaemonServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.ListDevices(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	self := resp.GetSelf()
	fmt.Printf("本机  %s  %s\n", self.GetName(), self.GetPublicKeyFingerprint())
	fmt.Printf("      %s\n\n", self.GetId())

	if len(resp.GetPeers()) == 0 {
		fmt.Println("（尚未配对任何设备，用 pair / join 建立配对）")
		return nil
	}
	for _, p := range resp.GetPeers() {
		fmt.Printf("%-16s %-10s %-8s %s\n",
			p.GetName(), p.GetPlatform(),
			connectionLabel(p), p.GetPublicKeyFingerprint())
	}
	return nil
}

// connectionLabel 区分直连与中转，便于排查网络问题：
// 一直显示「中转」说明 NAT 打洞没成功，走的是服务器带宽。
func connectionLabel(d *pb.Device) string {
	if !d.GetOnline() {
		return "离线"
	}
	switch d.GetConnection() {
	case pb.ConnectionKind_CONNECTION_KIND_DIRECT:
		return "直连"
	case pb.ConnectionKind_CONNECTION_KIND_RELAY:
		return "中转"
	default:
		return "在线"
	}
}

// createPairing 生成配对码并等待对方兑换。
func createPairing(ctx context.Context, c pb.DaemonServiceClient) error {
	// 先订阅事件流，才能在对方兑换后收到设备变化通知
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := c.Subscribe(streamCtx, &pb.SubscribeRequest{})
	if err != nil {
		return err
	}

	createCtx, createCancel := context.WithTimeout(ctx, 15*time.Second)
	resp, err := c.CreatePairingCode(createCtx, &pb.Empty{})
	createCancel()
	if err != nil {
		return err
	}

	fmt.Printf("\n  配对码：%s\n", resp.GetCode())
	fmt.Printf("  有效期至 %s\n\n",
		time.Unix(resp.GetExpiresAtUnix(), 0).Format("15:04:05"))
	fmt.Println("  在另一台设备上运行：copysync-cli join " + resp.GetCode())
	fmt.Println("  等待对方兑换…（Ctrl+C 取消）")

	for {
		ev, err := stream.Recv()
		if err != nil {
			return err
		}
		d := ev.GetDeviceChanged()
		if d == nil || d.GetPairingSession() == "" {
			continue // 非待确认的配对事件（普通的在线状态变化）
		}
		fmt.Printf("\n  对方已兑换：%s（%s）\n", d.GetName(), d.GetPlatform())
		if err := printFingerprints(ctx, c, d); err != nil {
			return err
		}
		fmt.Print("  确认配对？[y/N] ")

		var answer string
		fmt.Scanln(&answer)
		accept := answer == "y" || answer == "Y"

		confirmCtx, confirmCancel := context.WithTimeout(ctx, 10*time.Second)
		_, confirmErr := c.ConfirmPairing(confirmCtx, &pb.ConfirmPairingRequest{
			PairingSession: d.GetPairingSession(),
			Accept:         accept,
		})
		confirmCancel()
		if confirmErr != nil {
			return confirmErr
		}
		if accept {
			fmt.Println("\n  ✓ 配对成功")
		} else {
			fmt.Println("\n  已取消")
		}
		return nil
	}
}

// printFingerprints 打印配对时要核对的两行：本机与对方的公钥指纹，按指纹排序。
// 两台设备打印出的内容完全相同，用户比对两块屏幕即可；规则与界面一致
// （见 ui/lib/widgets/pairing_dialog.dart 的 pairingFingerprintRows）。
func printFingerprints(ctx context.Context, c pb.DaemonServiceClient, peer *pb.Device) error {
	resp, err := c.ListDevices(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	rows := []*pb.Device{resp.GetSelf(), peer}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].GetPublicKeyFingerprint() < rows[j].GetPublicKeyFingerprint()
	})
	fmt.Println("  安全指纹：")
	for _, d := range rows {
		fmt.Printf("    %s  %s\n", d.GetPublicKeyFingerprint(), d.GetName())
	}
	fmt.Println("\n  两台设备上显示的这两行应当完全相同，一致才继续。")
	return nil
}

func joinPairing(ctx context.Context, c pb.DaemonServiceClient, code string) error {
	if code == "" {
		return errors.New("用法: copysync-cli join <配对码>")
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	resp, err := c.RedeemPairingCode(ctx, &pb.RedeemPairingCodeRequest{Code: code})
	if err != nil {
		return err
	}
	peer := resp.GetPeer()
	fmt.Printf("\n  找到设备：%s（%s）\n", peer.GetName(), peer.GetPlatform())
	if err := printFingerprints(ctx, c, peer); err != nil {
		return err
	}
	fmt.Print("  确认配对？[y/N] ")

	var answer string
	fmt.Scanln(&answer)
	accept := answer == "y" || answer == "Y"

	if _, err := c.ConfirmPairing(ctx, &pb.ConfirmPairingRequest{
		PairingSession: resp.GetPairingSession(),
		Accept:         accept,
	}); err != nil {
		return err
	}
	if accept {
		fmt.Println("\n  ✓ 配对成功")
	} else {
		fmt.Println("\n  已取消")
	}
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

func showStatus(ctx context.Context, c pb.DaemonServiceClient, ep rpc.Endpoint) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	st, err := c.GetStatus(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	fmt.Printf("daemon      PID %d  端口 %d  版本 %s\n", ep.PID, ep.Port, st.GetVersion())
	fmt.Printf("设备        %s\n", st.GetDeviceName())
	fmt.Printf("设备 ID     %s\n", st.GetDeviceId())
	fmt.Printf("信令        %s\n", boolText(st.GetSignalingConnected(), "已连接", "未连接"))
	fmt.Printf("在线对端    %d\n", st.GetPeersOnline())
	fmt.Printf("剪贴板授权  %s\n", st.GetClipboardPermission())
	fmt.Printf("缓存占用    %s\n", humanBytes(st.GetCacheBytesUsed()))
	return nil
}

func watch(ctx context.Context, c pb.DaemonServiceClient) error {
	stream, err := c.Subscribe(ctx, &pb.SubscribeRequest{})
	if err != nil {
		return err
	}
	fmt.Println("已订阅事件流，Ctrl+C 退出")
	for {
		ev, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("流已结束")
			return nil
		}
		if err != nil {
			return err
		}
		printEvent(ev)
	}
}

func printEvent(ev *pb.Event) {
	ts := time.Now().Format("15:04:05")
	switch p := ev.GetPayload().(type) {
	case *pb.Event_ClipAdded:
		r := p.ClipAdded
		fmt.Printf("[%s] + 新记录  %s  %s  %s  %q\n", ts,
			r.GetKind(), r.GetStatus(), humanBytes(r.GetTotalSize()), r.GetTextPreview())
	case *pb.Event_ClipUpdated:
		r := p.ClipUpdated
		fmt.Printf("[%s] ~ 更新    %s  %s\n", ts, r.GetId(), r.GetStatus())
	case *pb.Event_ClipRemoved:
		fmt.Printf("[%s] - 删除    %s\n", ts, p.ClipRemoved)
	case *pb.Event_Progress:
		pr := p.Progress
		fmt.Printf("[%s] ↓ 进度    %s  %s/%s\n", ts, pr.GetClipId(),
			humanBytes(pr.GetTransferred()), humanBytes(pr.GetTotal()))
	case *pb.Event_DeviceChanged:
		d := p.DeviceChanged
		fmt.Printf("[%s] ⇄ 设备    %s  %s  %s\n", ts,
			d.GetName(), boolText(d.GetOnline(), "在线", "离线"), d.GetConnection())
	case *pb.Event_PermissionChanged:
		fmt.Printf("[%s] ! 授权    %s\n", ts, p.PermissionChanged)
	case *pb.Event_StatusChanged:
		s := p.StatusChanged
		fmt.Printf("[%s] i 状态    信令=%v 对端=%d 授权=%s\n", ts,
			s.GetSignalingConnected(), s.GetPeersOnline(), s.GetClipboardPermission())
	default:
		fmt.Printf("[%s] ? 未知事件 %T\n", ts, p)
	}
}

func showHistory(ctx context.Context, c pb.DaemonServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.ListHistory(ctx, &pb.ListHistoryRequest{Limit: 20})
	if err != nil {
		return err
	}
	if len(resp.GetRecords()) == 0 {
		fmt.Println("（暂无记录）")
		return nil
	}
	for _, r := range resp.GetRecords() {
		dir := "←收到"
		if r.GetOutgoing() {
			dir = "→本机"
		}
		fmt.Printf("%s  %s  %-6s %-14s %8s  %s\n",
			time.Unix(r.GetCreatedAtUnix(), 0).Format("01-02 15:04:05"),
			dir, r.GetKind(), r.GetStatus(),
			humanBytes(r.GetTotalSize()), preview(r))
	}
	return nil
}

func preview(r *pb.ClipRecord) string {
	if t := r.GetTextPreview(); t != "" {
		return fmt.Sprintf("%q", t)
	}
	switch n := len(r.GetItems()); n {
	case 0:
		return ""
	case 1:
		return r.GetItems()[0].GetName()
	default:
		return fmt.Sprintf("%s 等 %d 项", r.GetItems()[0].GetName(), n)
	}
}

func showConfig(ctx context.Context, c pb.DaemonServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cfg, err := c.GetConfig(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	fmt.Printf("设备名称      %s\n", cfg.GetDeviceName())
	fmt.Printf("信令地址      %s\n", cfg.GetSignalingUrl())
	fmt.Printf("自动同步阈值  %s\n", humanBytes(cfg.GetAutoSyncThresholdBytes()))
	fmt.Printf("历史保留      %s\n", time.Duration(cfg.GetHistoryTtlSeconds())*time.Second)
	fmt.Printf("缓存保留      %s\n", time.Duration(cfg.GetCacheTtlSeconds())*time.Second)
	fmt.Printf("同步类型      文本=%v HTML=%v 图片=%v 文件=%v\n",
		cfg.GetSyncText(), cfg.GetSyncHtml(), cfg.GetSyncImage(), cfg.GetSyncFile())
	fmt.Printf("自动写剪贴板  %v\n", cfg.GetAutoApplyToClipboard())
	return nil
}

func boolText(b bool, yes, no string) string {
	if b {
		return yes
	}
	return no
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}

// Package rpc 是 daemon 暴露给 Flutter UI 的本地 gRPC 服务。
package rpc

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/pasteferry/client-core/internal/clipboard"
	"github.com/baiyuze/pasteferry/client-core/internal/config"
	"github.com/baiyuze/pasteferry/client-core/internal/diag"
	"github.com/baiyuze/pasteferry/client-core/internal/peers"
	"github.com/baiyuze/pasteferry/client-core/internal/store"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/lan"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

const tokenHeader = "x-pasteferry-token"

// Deps 是 Server 依赖的 daemon 内部能力。用接口而非具体类型，
// 让 rpc 层不必反向依赖尚未实现的 sync / transport 包。
type Deps struct {
	Store    *store.Store
	Paths    config.Paths
	Version  string
	DeviceID string

	// LoadConfig / SaveConfig 由 daemon 提供，保证配置变更能即时生效
	LoadConfig func() config.Config
	SaveConfig func(config.Config) error

	// 以下在对应里程碑接上，未实现时为 nil，方法返回 Unimplemented
	Fetch                  func(ctx context.Context, clipID string) error
	FetchForPreview        func(ctx context.Context, clipID string) error
	ApplyToClipboard       func(ctx context.Context, clipID string) error
	ImagePreview           func(ctx context.Context, clipID string, index int) (path, name string, count int, err error)
	Permission             func() clipboard.Permission
	OpenPermissionSettings func() error
	// SearchHistory 按文字内容与文件名搜索复制记录，见 sync.Engine.SearchHistory
	SearchHistory func(query string, limit int) ([]store.Clip, error)
	// 诊断报告要读的日志文件（轮转出去的旧文件在前）与崩溃日志
	LogFiles func() (logs []string, crash string)
	// VerboseLog 临时打开详细日志，d 不大于 0 时关闭；VerboseUntil 是打开到什么时候
	VerboseLog   func(d time.Duration) time.Time
	VerboseUntil func() time.Time
	// Network 返回最近一轮网络出口探测的结果
	Network           func() *pb.NetworkInfo
	Devices           func() (self *pb.Device, peers []*pb.Device)
	CreatePairingCode func(ctx context.Context) (code string, expiresAt time.Time, lanAddrs []string, err error)
	RedeemPairingCode func(ctx context.Context, code, address string) (*pb.Device, string, error)
	ConfirmPairing    func(ctx context.Context, session string, accept bool) error
	Unpair            func(ctx context.Context, deviceID string) error
	SetDeviceAddress  func(ctx context.Context, deviceID, address string) error
	// ShowNearby / Nearby / PairNearby 是「搜索附近设备」，见 lan/nearby.go
	ShowNearby         func() (hide func())
	Nearby             func() []*pb.NearbyDevice
	PairNearby         func(ctx context.Context, id string) (*pb.Device, string, error)
	SignalingConnected func() bool
	ServerConfigured   func() bool
	LANEnabled         func() bool
	LANBlocked         func() bool
	// OpenLocalNetworkSettings 让系统放行局域网访问（macOS 打开设置面板，Windows 放行防火墙）
	OpenLocalNetworkSettings func() error
	CacheBytesUsed           func() int64
	OfficialUsage            func() *pb.OfficialUsage
	NotificationStatus       func(context.Context) (*pb.NotificationStatus, error)
	RequestNotifications     func(context.Context) (*pb.NotificationStatus, error)
	ActivateNotification     func(context.Context, string, string) error
}

func (s *Server) GetNotificationStatus(ctx context.Context, _ *pb.Empty) (*pb.NotificationStatus, error) {
	if s.deps.NotificationStatus == nil {
		return &pb.NotificationStatus{Permission: "unavailable"}, nil
	}
	return s.deps.NotificationStatus(ctx)
}

func (s *Server) RequestNotificationPermission(ctx context.Context, _ *pb.Empty) (*pb.NotificationStatus, error) {
	if s.deps.RequestNotifications == nil {
		return &pb.NotificationStatus{Permission: "unavailable"}, nil
	}
	return s.deps.RequestNotifications(ctx)
}

func (s *Server) ActivateNotification(ctx context.Context, req *pb.NotificationAction) (*pb.Empty, error) {
	if s.deps.ActivateNotification == nil {
		return nil, status.Error(codes.Unimplemented, "Notifications unavailable")
	}
	if err := s.deps.ActivateNotification(ctx, req.GetClipId(), req.GetAction()); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &pb.Empty{}, nil
}

type Server struct {
	pb.UnimplementedDaemonServiceServer
	deps  Deps
	hub   *Hub
	token string
}

func NewServer(deps Deps, hub *Hub, token string) *Server {
	return &Server{deps: deps, hub: hub, token: token}
}

// Serve 在 127.0.0.1 的随机端口上启动 gRPC，返回实际端口。
// 只监听回环地址：daemon 不对外提供服务。
func (s *Server) Serve(ctx context.Context) (int, *grpc.Server, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, nil, fmt.Errorf("监听本地端口: %w", err)
	}

	gs := grpc.NewServer(
		grpc.UnaryInterceptor(s.authUnary),
		grpc.StreamInterceptor(s.authStream),
	)
	pb.RegisterDaemonServiceServer(gs, s)

	go func() {
		<-ctx.Done()
		gs.GracefulStop()
	}()
	go func() {
		if err := gs.Serve(lis); err != nil {
			// GracefulStop 会让 Serve 正常返回，这里只记录异常退出
			_ = err
		}
	}()

	return lis.Addr().(*net.TCPAddr).Port, gs, nil
}

// ─────────────────────────── 认证 ───────────────────────────

// 本机回环端口对同机其他进程同样可达，因此用一次性 token 把关。
func (s *Server) checkToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "缺少认证信息")
	}
	vals := md.Get(tokenHeader)
	if len(vals) == 0 {
		return status.Error(codes.Unauthenticated, "缺少 token")
	}
	// 常数时间比较，避免时序侧信道
	if subtle.ConstantTimeCompare([]byte(vals[0]), []byte(s.token)) != 1 {
		return status.Error(codes.Unauthenticated, "token 无效")
	}
	return nil
}

func (s *Server) authUnary(ctx context.Context, req any,
	_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if err := s.checkToken(ctx); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func (s *Server) authStream(srv any, ss grpc.ServerStream,
	_ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := s.checkToken(ss.Context()); err != nil {
		return err
	}
	return handler(srv, ss)
}

// ─────────────────────────── 事件流 ───────────────────────────

func (s *Server) Subscribe(req *pb.SubscribeRequest,
	stream grpc.ServerStreamingServer[pb.Event]) error {
	ch, cancel := s.hub.Subscribe()
	defer cancel()

	// 先补发断线期间错过的记录，再转入实时推送
	if req.GetSinceUnix() > 0 {
		since := time.Unix(req.GetSinceUnix(), 0)
		clips, err := s.deps.Store.ListClips(200, time.Time{}, store.KindUnspecified)
		if err != nil {
			return status.Errorf(codes.Internal, "读取历史: %v", err)
		}
		// ListClips 是倒序，补发时正序推送更符合 UI 预期
		for i := len(clips) - 1; i >= 0; i-- {
			if clips[i].CreatedAt.After(since) {
				if err := stream.Send(EventClipAdded(clipToProto(clips[i]))); err != nil {
					return err
				}
			}
		}
	}

	// 立即推一次状态，让 UI 不必等到下次变化才知道当前情况
	if st, err := s.GetStatus(stream.Context(), &pb.Empty{}); err == nil {
		if err := stream.Send(EventStatusChanged(st)); err != nil {
			return err
		}
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case e, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(e); err != nil {
				return err
			}
		}
	}
}

// ─────────────────────────── 历史记录 ───────────────────────────

func (s *Server) ListHistory(_ context.Context,
	req *pb.ListHistoryRequest) (*pb.ListHistoryResponse, error) {
	if q := strings.TrimSpace(req.GetQuery()); q != "" {
		return s.searchHistory(q, req)
	}
	var before time.Time
	if req.GetBeforeUnix() > 0 {
		before = time.Unix(req.GetBeforeUnix(), 0)
	}
	clips, err := s.deps.Store.ListClips(int(req.GetLimit()), before, store.Kind(req.GetKindFilter()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "读取历史: %v", err)
	}
	out := make([]*pb.ClipRecord, 0, len(clips))
	for _, c := range clips {
		out = append(out, clipToProto(c))
	}
	return &pb.ListHistoryResponse{Records: out}, nil
}

func (s *Server) searchHistory(query string, req *pb.ListHistoryRequest) (*pb.ListHistoryResponse, error) {
	if s.deps.SearchHistory == nil {
		return nil, status.Error(codes.Unimplemented, "搜索未接入")
	}
	clips, err := s.deps.SearchHistory(query, int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "搜索历史: %v", err)
	}
	out := make([]*pb.ClipRecord, 0, len(clips))
	for _, c := range clips {
		if k := req.GetKindFilter(); k != pb.ClipKind_CLIP_KIND_UNSPECIFIED && pb.ClipKind(c.Kind) != k {
			continue
		}
		out = append(out, clipToProto(c))
	}
	return &pb.ListHistoryResponse{Records: out, Searched: true}, nil
}

func (s *Server) DeleteHistory(_ context.Context, req *pb.DeleteHistoryRequest) (*pb.Empty, error) {
	var err error
	if req.GetAll() {
		err = s.deps.Store.DeleteAllClips()
	} else {
		err = s.deps.Store.DeleteClips(req.GetIds())
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "删除历史: %v", err)
	}
	for _, id := range req.GetIds() {
		s.hub.Publish(EventClipRemoved(id))
	}
	return &pb.Empty{}, nil
}

func (s *Server) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.Empty, error) {
	fetch := s.deps.Fetch
	if req.GetPreserveClipboard() {
		fetch = s.deps.FetchForPreview
	}
	if fetch == nil {
		return nil, status.Error(codes.Unimplemented, "传输层尚未接入（M4）")
	}
	if err := fetch(ctx, req.GetClipId()); err != nil {
		return nil, status.Errorf(codes.Internal, "拉取失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Server) ApplyToClipboard(ctx context.Context,
	req *pb.ApplyToClipboardRequest) (*pb.Empty, error) {
	if s.deps.ApplyToClipboard == nil {
		return nil, status.Error(codes.Unimplemented, "剪贴板写入尚未接入（M3）")
	}
	if err := s.deps.ApplyToClipboard(ctx, req.GetClipId()); err != nil {
		return nil, status.Errorf(codes.Internal, "写入剪贴板失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Server) GetImagePreview(ctx context.Context,
	req *pb.GetImagePreviewRequest) (*pb.GetImagePreviewResponse, error) {
	if req.GetClipId() == "" {
		return nil, status.Error(codes.InvalidArgument, "缺少图片记录 ID")
	}
	if s.deps.ImagePreview == nil {
		return nil, status.Error(codes.Unimplemented, "请更新后台服务后再预览图片")
	}
	path, name, count, err := s.deps.ImagePreview(ctx, req.GetClipId(), int(req.GetIndex()))
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &pb.GetImagePreviewResponse{Path: path, Name: name, Count: int32(count)}, nil
}

// ─────────────────────────── 配置与状态 ───────────────────────────

func (s *Server) GetConfig(context.Context, *pb.Empty) (*pb.Config, error) {
	return configToProto(s.deps.LoadConfig()), nil
}

func (s *Server) UpdateConfig(_ context.Context, in *pb.Config) (*pb.Config, error) {
	c := configFromProto(in, s.deps.LoadConfig())
	c.Validate()
	if err := s.deps.SaveConfig(c); err != nil {
		return nil, status.Errorf(codes.Internal, "保存配置: %v", err)
	}
	return configToProto(c), nil
}

func (s *Server) GetStatus(context.Context, *pb.Empty) (*pb.Status, error) {
	cfg := s.deps.LoadConfig()
	st := &pb.Status{
		DeviceId:   s.deps.DeviceID,
		DeviceName: cfg.DeviceName,
		Version:    s.deps.Version,
	}
	if s.deps.OfficialUsage != nil {
		st.OfficialUsage = s.deps.OfficialUsage()
	}
	if s.deps.Permission != nil {
		st.ClipboardPermission = permissionToProto(s.deps.Permission())
	}
	if s.deps.SignalingConnected != nil {
		st.SignalingConnected = s.deps.SignalingConnected()
	}
	if s.deps.ServerConfigured != nil {
		st.ServerConfigured = s.deps.ServerConfigured()
	}
	if s.deps.LANEnabled != nil {
		st.LanEnabled = s.deps.LANEnabled()
	}
	if s.deps.LANBlocked != nil {
		st.LanBlocked = s.deps.LANBlocked()
	}
	if s.deps.VerboseUntil != nil {
		if until := s.deps.VerboseUntil(); !until.IsZero() {
			st.VerboseLogUntilUnix = until.Unix()
		}
	}
	if s.deps.CacheBytesUsed != nil {
		st.CacheBytesUsed = s.deps.CacheBytesUsed()
	}
	if s.deps.Devices != nil {
		_, peers := s.deps.Devices()
		for _, p := range peers {
			if p.GetOnline() {
				st.PeersOnline++
			}
		}
	}
	return st, nil
}

func (s *Server) GetNetwork(context.Context, *pb.Empty) (*pb.NetworkInfo, error) {
	if s.deps.Network == nil {
		return nil, status.Error(codes.Unimplemented, "当前版本不支持网络诊断")
	}
	return s.deps.Network(), nil
}

// ExportDiagnostics 生成反馈问题用的诊断报告，脱敏在 diag 包里做。
func (s *Server) ExportDiagnostics(ctx context.Context, req *pb.ExportDiagnosticsRequest) (*pb.Diagnostics, error) {
	lines := int(req.GetLogLines())
	if lines <= 0 || lines > 2000 {
		lines = 200
	}
	st, err := s.GetStatus(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	in := diag.Input{
		Now:       time.Now(),
		Version:   s.deps.Version,
		UIVersion: req.GetUiVersion(),
		Config:    s.deps.LoadConfig(),
		Status:    st,
		LogLines:  lines,
	}
	if s.deps.Devices != nil {
		in.Self, in.Peers = s.deps.Devices()
	}
	if s.deps.Network != nil {
		in.Network = s.deps.Network()
	}
	if s.deps.VerboseUntil != nil {
		in.VerboseUntil = s.deps.VerboseUntil()
	}
	if s.deps.LogFiles != nil {
		in.LogFiles, in.CrashLog = s.deps.LogFiles()
	}
	in.Home, _ = os.UserHomeDir()
	if s.deps.Store != nil {
		// 名字取不全只是少换掉几个，报告照样能用
		in.FileNames, _ = s.deps.Store.FileNames(5000)
	}
	return &pb.Diagnostics{Report: diag.Build(in)}, nil
}

func (s *Server) SetVerboseLog(ctx context.Context, req *pb.SetVerboseLogRequest) (*pb.Status, error) {
	if s.deps.VerboseLog == nil {
		return nil, status.Error(codes.Unimplemented, "当前版本不支持临时打开详细日志")
	}
	// 最长两小时：详细日志量大，忘了关也不至于一直写下去
	minutes := min(max(req.GetMinutes(), 0), 120)
	s.deps.VerboseLog(time.Duration(minutes) * time.Minute)
	return s.GetStatus(ctx, &pb.Empty{})
}

func (s *Server) RequestClipboardPermission(context.Context, *pb.Empty) (*pb.Empty, error) {
	if s.deps.OpenPermissionSettings == nil {
		return nil, status.Error(codes.Unimplemented, "当前平台无需授权")
	}
	if err := s.deps.OpenPermissionSettings(); err != nil {
		return nil, status.Errorf(codes.Internal, "打开设置失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Server) RequestLocalNetworkPermission(context.Context, *pb.Empty) (*pb.Empty, error) {
	if s.deps.OpenLocalNetworkSettings == nil {
		return nil, status.Error(codes.Unimplemented, "当前平台无需授权")
	}
	if err := s.deps.OpenLocalNetworkSettings(); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return nil, status.Error(codes.Unimplemented, "当前平台无需授权")
		}
		return nil, status.Errorf(codes.Internal, "放行局域网访问失败: %v", err)
	}
	return &pb.Empty{}, nil
}

// ─────────────────────────── 设备与配对 ───────────────────────────

// 配对出错时，消息以「原因: 说明」开头，界面按原因显示本地化的提示，
// 不认识的原因再直接显示说明：
//
//	not_found     局域网里没找到正在配对的设备，服务器也没连（界面会请用户输入对方地址）
//	wrong_code    配对码不对（或服务器那边码不存在、已过期，服务器不区分）
//	code_expired  配对码已失效，或输错次数太多
//	bad_address   输入的地址无法识别
//	no_channel    局域网发现已关闭，也没连服务器
func pairingStatus(err error) error {
	reason, code := "", codes.Internal
	switch {
	case errors.Is(err, lan.ErrNotFound):
		reason, code = "not_found", codes.NotFound
	case errors.Is(err, lan.ErrWrongCode), errors.Is(err, peers.ErrCodeRejected):
		reason, code = "wrong_code", codes.InvalidArgument
	case errors.Is(err, lan.ErrCodeClosed), errors.Is(err, lan.ErrTooManyAttempts):
		reason, code = "code_expired", codes.FailedPrecondition
	case errors.Is(err, lan.ErrBadAddress):
		reason, code = "bad_address", codes.InvalidArgument
	case errors.Is(err, peers.ErrNoChannel), errors.Is(err, lan.ErrDisabled):
		reason, code = "no_channel", codes.FailedPrecondition
	}
	if reason == "" {
		return status.Error(code, err.Error())
	}
	return status.Error(code, reason+": "+err.Error())
}

func (s *Server) ListDevices(context.Context, *pb.Empty) (*pb.ListDevicesResponse, error) {
	if s.deps.Devices == nil {
		return nil, status.Error(codes.Unimplemented, "配对尚未接入（M1）")
	}
	self, peers := s.deps.Devices()
	return &pb.ListDevicesResponse{Self: self, Peers: peers}, nil
}

func (s *Server) CreatePairingCode(ctx context.Context,
	_ *pb.Empty) (*pb.CreatePairingCodeResponse, error) {
	if s.deps.CreatePairingCode == nil {
		return nil, status.Error(codes.Unimplemented, "配对尚未接入（M1）")
	}
	code, exp, addrs, err := s.deps.CreatePairingCode(ctx)
	if err != nil {
		return nil, pairingStatus(err)
	}
	return &pb.CreatePairingCodeResponse{Code: code, ExpiresAtUnix: exp.Unix(), LanAddresses: addrs}, nil
}

func (s *Server) RedeemPairingCode(ctx context.Context,
	req *pb.RedeemPairingCodeRequest) (*pb.RedeemPairingCodeResponse, error) {
	if s.deps.RedeemPairingCode == nil {
		return nil, status.Error(codes.Unimplemented, "配对尚未接入（M1）")
	}
	peer, session, err := s.deps.RedeemPairingCode(ctx, req.GetCode(), req.GetAddress())
	if err != nil {
		return nil, pairingStatus(err)
	}
	return &pb.RedeemPairingCodeResponse{Peer: peer, PairingSession: session}, nil
}

func (s *Server) ConfirmPairing(ctx context.Context, req *pb.ConfirmPairingRequest) (*pb.Empty, error) {
	if s.deps.ConfirmPairing == nil {
		return nil, status.Error(codes.Unimplemented, "配对尚未接入（M1）")
	}
	if err := s.deps.ConfirmPairing(ctx, req.GetPairingSession(), req.GetAccept()); err != nil {
		return nil, status.Errorf(codes.Internal, "确认配对: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Server) Unpair(ctx context.Context, req *pb.UnpairRequest) (*pb.Empty, error) {
	if s.deps.Unpair == nil {
		return nil, status.Error(codes.Unimplemented, "配对尚未接入（M1）")
	}
	if err := s.deps.Unpair(ctx, req.GetDeviceId()); err != nil {
		return nil, status.Errorf(codes.Internal, "解除配对: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Server) SetDeviceAddress(ctx context.Context, req *pb.SetDeviceAddressRequest) (*pb.Empty, error) {
	if s.deps.SetDeviceAddress == nil {
		return nil, status.Error(codes.Unimplemented, "当前版本不支持手动输入地址")
	}
	if err := s.deps.SetDeviceAddress(ctx, req.GetDeviceId(), req.GetAddress()); err != nil {
		if errors.Is(err, peers.ErrNotPaired) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, pairingStatus(err)
	}
	return &pb.Empty{}, nil
}

// WatchNearby 在调用期间公开本机（界面上「搜索附近设备」的窗口开着），
// 附近的设备有变化时推送最新的列表。窗口关掉、界面退出时调用随之结束，本机也随之撤下。
func (s *Server) WatchNearby(_ *pb.Empty, stream grpc.ServerStreamingServer[pb.NearbyDevices]) error {
	if s.deps.ShowNearby == nil || s.deps.Nearby == nil {
		return status.Error(codes.Unimplemented, "当前版本不支持搜索附近设备")
	}
	hide := s.deps.ShowNearby()
	defer hide()

	tk := time.NewTicker(time.Second)
	defer tk.Stop()
	var last *pb.NearbyDevices
	for {
		cur := &pb.NearbyDevices{Devices: s.deps.Nearby()}
		if last == nil || !proto.Equal(cur, last) {
			if err := stream.Send(cur); err != nil {
				return err
			}
			last = cur
		}
		select {
		case <-stream.Context().Done():
			return nil
		case <-tk.C:
		}
	}
}

func (s *Server) PairNearby(ctx context.Context, req *pb.PairNearbyRequest) (*pb.RedeemPairingCodeResponse, error) {
	if s.deps.PairNearby == nil {
		return nil, status.Error(codes.Unimplemented, "当前版本不支持搜索附近设备")
	}
	peer, session, err := s.deps.PairNearby(ctx, req.GetId())
	if err != nil {
		return nil, pairingStatus(err)
	}
	return &pb.RedeemPairingCodeResponse{Peer: peer, PairingSession: session}, nil
}

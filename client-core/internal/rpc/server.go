// Package rpc 是 daemon 暴露给 Flutter UI 的本地 gRPC 服务。
package rpc

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/baiyuze/copysync/client-core/internal/clipboard"
	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/store"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

const tokenHeader = "x-copysync-token"

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
	ImagePreviewPath       func(ctx context.Context, clipID string) (string, error)
	Permission             func() clipboard.Permission
	OpenPermissionSettings func() error
	// Network 返回最近一轮网络出口探测的结果
	Network            func() *pb.NetworkInfo
	Devices            func() (self *pb.Device, peers []*pb.Device)
	CreatePairingCode  func(ctx context.Context) (code string, expiresAt time.Time, err error)
	RedeemPairingCode  func(ctx context.Context, code string) (*pb.Device, string, error)
	ConfirmPairing     func(ctx context.Context, session string, accept bool) error
	Unpair             func(ctx context.Context, deviceID string) error
	SignalingConnected func() bool
	CacheBytesUsed     func() int64
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
	if s.deps.ImagePreviewPath == nil {
		return nil, status.Error(codes.Unimplemented, "请更新后台服务后再预览图片")
	}
	path, err := s.deps.ImagePreviewPath(ctx, req.GetClipId())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &pb.GetImagePreviewResponse{Path: path}, nil
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
	if s.deps.Permission != nil {
		st.ClipboardPermission = permissionToProto(s.deps.Permission())
	}
	if s.deps.SignalingConnected != nil {
		st.SignalingConnected = s.deps.SignalingConnected()
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

func (s *Server) RequestClipboardPermission(context.Context, *pb.Empty) (*pb.Empty, error) {
	if s.deps.OpenPermissionSettings == nil {
		return nil, status.Error(codes.Unimplemented, "当前平台无需授权")
	}
	if err := s.deps.OpenPermissionSettings(); err != nil {
		return nil, status.Errorf(codes.Internal, "打开设置失败: %v", err)
	}
	return &pb.Empty{}, nil
}

// ─────────────────────────── 设备与配对 ───────────────────────────

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
	code, exp, err := s.deps.CreatePairingCode(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "生成配对码: %v", err)
	}
	return &pb.CreatePairingCodeResponse{Code: code, ExpiresAtUnix: exp.Unix()}, nil
}

func (s *Server) RedeemPairingCode(ctx context.Context,
	req *pb.RedeemPairingCodeRequest) (*pb.RedeemPairingCodeResponse, error) {
	if s.deps.RedeemPairingCode == nil {
		return nil, status.Error(codes.Unimplemented, "配对尚未接入（M1）")
	}
	peer, session, err := s.deps.RedeemPairingCode(ctx, req.GetCode())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "兑换配对码: %v", err)
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

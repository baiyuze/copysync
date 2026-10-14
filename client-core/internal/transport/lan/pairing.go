package lan

import (
	"context"
	cryptorand "crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/netip"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/signed"
	"github.com/baiyuze/pasteferry/proto/deviceid"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"github.com/baiyuze/pasteferry/proto/pairingcode"
)

// 在局域网里配对：显示配对码的一方把码「开放」出来（mDNS 广播里带 p=1），
// 输入配对码的一方把码发给所有正在配对的实例，以及用户输入的地址。
//
// 配对码在局域网里明文传输，与经服务器配对时服务器能看到配对码是一回事：
// 拿到配对码最多能冒充一方发起配对，但两端显示的指纹会对不上——
// 防中间人靠的是用户核对指纹，不是配对码。
// 显示码的一方限制错误次数，挡住在 5 分钟内暴力猜码。

const maxPairAttempts = 10

var (
	ErrNotFound        = errors.New("没在局域网里找到正在配对的设备")
	ErrWrongCode       = errors.New("配对码不对")
	ErrCodeClosed      = errors.New("配对码已失效")
	ErrTooManyAttempts = errors.New("配对码输错次数太多，已失效")
)

type pairingState struct {
	code     string
	expires  time.Time
	attempts int
}

// cachedResponse 是发出去的配对回复：对方没收到而重发请求时原样再发一次。
type cachedResponse struct {
	pkt []byte
	at  time.Time
}

type pairWaiter struct {
	resp       chan pairResult
	rejects    map[string]bool
	rejectedAt time.Time
}

type pairResult struct {
	info *pb.LanPairInfo
	src  netip.AddrPort
}

// OpenPairing 在局域网里开放一个配对码，有效期到 expires。
func (t *Transport) OpenPairing(code string, expires time.Time) {
	t.mu.Lock()
	t.pairing = &pairingState{code: pairingcode.Normalize(code), expires: expires}
	t.mu.Unlock()
	t.updateAdvert()
}

// ClosePairing 撤回正在开放的配对码。
func (t *Transport) ClosePairing() {
	t.mu.Lock()
	was := t.pairing != nil
	t.pairing = nil
	t.mu.Unlock()
	if was {
		t.updateAdvert()
	}
}

// Pair 用配对码在局域网里找到对方并交换公钥，直到 ctx 结束。extra 是用户输入的地址，
// 不在同一网段（mDNS 过不去）时靠它。
func (t *Transport) Pair(ctx context.Context, code string, extra []netip.AddrPort) (PeerInfo, error) {
	if !t.Running() {
		return PeerInfo{}, ErrDisabled
	}
	req := t.selfPairInfo()
	req.Code = pairingcode.Normalize(code)
	return t.exchange(ctx, req, pairRequest, func() []netip.AddrPort { return t.pairTargets(extra) }, true)
}

// exchange 每 500 毫秒把配对请求发往 targets 给出的地址，直到收到回复、被明确拒绝或 ctx 结束。
// seek 为真时广播里带上 s=1，显示码的一方会先往本机发包，见 found。
func (t *Transport) exchange(ctx context.Context, req *pb.LanPairInfo, kind pairKind,
	targets func() []netip.AddrPort, seek bool) (PeerInfo, error) {
	nonce := make([]byte, 16)
	if _, err := cryptorand.Read(nonce); err != nil {
		return PeerInfo{}, err
	}
	req.Nonce = nonce
	pkt, err := t.encodePairInfo(req, kind)
	if err != nil {
		return PeerInfo{}, err
	}

	key := hex.EncodeToString(nonce)
	w := &pairWaiter{resp: make(chan pairResult, 1), rejects: make(map[string]bool)}
	t.mu.Lock()
	t.waiters[key] = w
	if seek {
		t.seeking++
	}
	t.mu.Unlock()
	if seek {
		t.updateAdvert()
	}
	defer func() {
		t.mu.Lock()
		delete(t.waiters, key)
		if seek {
			t.seeking--
		}
		t.mu.Unlock()
		if seek {
			t.updateAdvert()
		}
	}()

	tk := time.NewTicker(500 * time.Millisecond)
	defer tk.Stop()
	for round := 0; ; round++ {
		dsts := targets()
		if round%4 == 0 {
			t.log.Debug("发送配对请求", "目标", dsts)
		}
		for _, dst := range dsts {
			_ = t.write(dst, pkt)
		}
		select {
		case r := <-w.resp:
			// 回复带着对方的签名：它此刻就在这个地址，确认配对后的第一次连接不必再等打招呼
			t.sawPeer(r.info.GetDeviceId(), r.src, nil)
			return peerFromInfo(r.info, r.src), nil
		case <-ctx.Done():
			if err := t.rejectError(w, true); err != nil {
				return PeerInfo{}, err
			}
			return PeerInfo{}, ErrNotFound
		case <-tk.C:
			// 有人明确拒绝了：再等一轮，万一另一台设备才是对的那台；之后就不必等到超时
			if err := t.rejectError(w, false); err != nil {
				return PeerInfo{}, err
			}
		}
	}
}

// rejectError 把收到的拒绝换成错误。final 为假时，只在拒绝到达半秒以后才算数。
func (t *Transport) rejectError(w *pairWaiter, final bool) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if w.rejectedAt.IsZero() || (!final && time.Since(w.rejectedAt) < 500*time.Millisecond) {
		return nil
	}
	switch {
	case w.rejects["wrong_code"]:
		return ErrWrongCode
	case w.rejects["too_many_attempts"]:
		return ErrTooManyAttempts
	}
	return ErrCodeClosed
}

// pairTargets 是配对请求要发往的地址：用户输入的，加上正在广播配对的实例（最多 16 个）。
func (t *Transport) pairTargets(extra []netip.AddrPort) []netip.AddrPort {
	out := append([]netip.AddrPort(nil), extra...)
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, si := range t.instances {
		if len(out) >= 16+len(extra) {
			break
		}
		if si.inst.Text["p"] == "1" && now.Sub(si.seen) < time.Minute {
			out = append(out, si.inst.Addr)
		}
	}
	return out
}

func (t *Transport) handlePairRequest(ps *pb.LanPairSigned, src netip.AddrPort) {
	info, ok := t.openPairInfo(ps)
	if !ok || len(info.GetNonce()) < 8 || t.resendResponse(info.GetNonce(), src) {
		return
	}
	now := time.Now()

	t.mu.Lock()
	p := t.pairing
	reason := ""
	closeAdvert := false
	switch {
	case p == nil || now.After(p.expires):
		reason = "closed"
	case subtle.ConstantTimeCompare([]byte(pairingcode.Normalize(info.GetCode())), []byte(p.code)) != 1:
		p.attempts++
		reason = "wrong_code"
		if p.attempts >= maxPairAttempts {
			t.log.Warn("配对码输错次数太多，已作废", "src", src)
			t.pairing, reason, closeAdvert = nil, "too_many_attempts", true
		}
	default:
		// 码对了：先撤下，同一个码不会再被第二个人用到
		t.pairing, closeAdvert = nil, true
	}
	t.mu.Unlock()
	if closeAdvert {
		t.updateAdvert()
	}
	if reason != "" {
		t.reject(src, info.GetNonce(), reason)
		return
	}

	peer := peerFromInfo(info, src)
	if t.opts.AcceptPair != nil && !t.opts.AcceptPair(peer) {
		t.reject(src, info.GetNonce(), "closed")
		return
	}
	t.log.Info("局域网配对请求已接受，等待用户核对指纹", "peer", peer.Name, "src", src)
	t.respondPair(info, peer, src)
}

// resendResponse：对方没收到回复而重发了同一个请求时，把那份回复原样再发一次。
func (t *Transport) resendResponse(nonce []byte, src netip.AddrPort) bool {
	t.mu.Lock()
	cached, ok := t.responses[hex.EncodeToString(nonce)]
	t.mu.Unlock()
	if ok {
		_ = t.write(src, cached.pkt)
	}
	return ok
}

// respondPair 回复一个已接受的配对请求：本机的信息，回显请求的 nonce。
func (t *Transport) respondPair(info *pb.LanPairInfo, peer PeerInfo, src netip.AddrPort) {
	resp := t.selfPairInfo()
	resp.Nonce = make([]byte, 16)
	_, _ = cryptorand.Read(resp.Nonce)
	resp.RequestNonce = info.GetNonce()
	pkt, err := t.encodePairInfo(resp, pairResponse)
	if err != nil {
		return
	}
	t.mu.Lock()
	t.responses[hex.EncodeToString(info.GetNonce())] = cachedResponse{pkt: pkt, at: time.Now()}
	t.mu.Unlock()
	t.sawPeer(peer.ID, src, nil)
	_ = t.write(src, pkt)
	time.AfterFunc(resendDelay, func() { _ = t.write(src, pkt) })
}

func peerFromInfo(info *pb.LanPairInfo, src netip.AddrPort) PeerInfo {
	return PeerInfo{
		ID:        info.GetDeviceId(),
		Name:      info.GetDeviceName(),
		Platform:  info.GetPlatform(),
		PublicKey: info.GetPublicKey(),
		Addrs:     mergeAddrs(info.GetAddresses(), src),
	}
}

func (t *Transport) handlePairResponse(ps *pb.LanPairSigned, src netip.AddrPort) {
	info, ok := t.openPairInfo(ps)
	if !ok {
		return
	}
	t.mu.Lock()
	w := t.waiters[hex.EncodeToString(info.GetRequestNonce())]
	t.mu.Unlock()
	// 按回显的 nonce 找到发出请求的那次 Pair：回复只能对应本机真正发过的请求
	if w == nil {
		return
	}
	select {
	case w.resp <- pairResult{info: info, src: src}:
	default:
	}
}

func (t *Transport) handlePairReject(r *pb.LanPairReject) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if w := t.waiters[hex.EncodeToString(r.GetRequestNonce())]; w != nil {
		w.rejects[r.GetReason()] = true
		if w.rejectedAt.IsZero() {
			w.rejectedAt = time.Now()
		}
	}
}

func (t *Transport) reject(dst netip.AddrPort, nonce []byte, reason string) {
	pkt, err := encode(&pb.LanPacket{Payload: &pb.LanPacket_PairReject{PairReject: &pb.LanPairReject{
		RequestNonce: nonce,
		Reason:       reason,
	}}})
	if err == nil {
		_ = t.write(dst, pkt)
	}
}

func (t *Transport) selfPairInfo() *pb.LanPairInfo {
	name, platform := "", ""
	if t.opts.DeviceInfo != nil {
		name, platform = t.opts.DeviceInfo()
	}
	t.mu.Lock()
	addrs := append([]string(nil), t.selfAddrs...)
	t.mu.Unlock()
	return &pb.LanPairInfo{
		DeviceId:      t.id.DeviceID,
		DeviceName:    name,
		Platform:      platform,
		PublicKey:     t.id.PublicKey,
		Addresses:     addrs,
		TimestampUnix: time.Now().Unix(),
	}
}

// pairKind 是签名的配对消息放在 LanPacket 的哪个位置。
type pairKind int

const (
	pairRequest pairKind = iota
	pairResponse
	nearbyRequest // 在「搜索附近设备」里点了对方，不带配对码
)

func (t *Transport) encodePairInfo(info *pb.LanPairInfo, kind pairKind) ([]byte, error) {
	payload, err := proto.Marshal(info)
	if err != nil {
		return nil, err
	}
	ps := &pb.LanPairSigned{Payload: payload, Signature: t.id.Sign(payload)}
	switch kind {
	case pairRequest:
		return encode(&pb.LanPacket{Payload: &pb.LanPacket_PairRequest{PairRequest: ps}})
	case nearbyRequest:
		return encode(&pb.LanPacket{Payload: &pb.LanPacket_NearbyRequest{NearbyRequest: ps}})
	}
	return encode(&pb.LanPacket{Payload: &pb.LanPacket_PairResponse{PairResponse: ps}})
}

// openPairInfo 解出并校验配对消息：签名对得上自带的公钥、ID 由公钥派生、时间戳不离谱、不是自己。
func (t *Transport) openPairInfo(ps *pb.LanPairSigned) (*pb.LanPairInfo, bool) {
	var info pb.LanPairInfo
	if err := proto.Unmarshal(ps.GetPayload(), &info); err != nil {
		return nil, false
	}
	if !identity.Verify(info.GetPublicKey(), ps.GetPayload(), ps.GetSignature()) ||
		!deviceid.Valid(info.GetDeviceId(), info.GetPublicKey()) ||
		info.GetDeviceId() == t.id.DeviceID {
		return nil, false
	}
	skew := time.Since(time.Unix(info.GetTimestampUnix(), 0))
	if skew > signed.MaxClockSkew || skew < -signed.MaxClockSkew {
		return nil, false
	}
	return &info, true
}

// mergeAddrs 把对方自报的地址与实际来源合在一起，来源排第一。
func mergeAddrs(reported []string, src netip.AddrPort) []string {
	out := []string{src.String()}
	for _, a := range reported {
		if ap, err := netip.ParseAddrPort(a); err == nil && ap.Addr().Is4() && ap != src {
			out = append(out, ap.String())
		}
	}
	return out
}

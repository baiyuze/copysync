// Package signed 负责信令的签名与校验，供服务器和局域网两条通道共用。
//
// 信任模型与通道无关：无论信令经服务器转发还是在局域网里直接收发，
// 发送方都用设备私钥签名，接收方用配对时存下的公钥验签。转发的一方——
// 服务器，或同一网络里的任何人——改不了内容，也冒充不了任何一方。
//
// 两条通道共用同一个 Verifier，也就共用同一份 nonce 记录：同一条信令
// 从两条路都到达时，只有先到的那份会被交给上层。
package signed

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

// MaxClockSkew 是防重放的时间窗：时间戳偏离超过它的消息一律拒收。
const MaxClockSkew = 5 * time.Minute

var (
	ErrUnknownSender  = errors.New("发送方未配对")
	ErrBadSignature   = errors.New("签名无效")
	ErrMalformed      = errors.New("消息无法解析")
	ErrStale          = errors.New("时间戳超出容忍范围")
	ErrReplay         = errors.New("重复的消息")
	ErrWrongRecipient = errors.New("消息不是发给本机的")
)

// Sign 给消息填上收件人、时间戳与 nonce，签名后返回可直接发送的 Signed。
// 会改写 msg。
func Sign(id *identity.Identity, to string, msg *pb.SignalMessage) (*pb.Signed, error) {
	msg.ToDeviceId = to
	msg.TimestampUnix = time.Now().Unix()
	if msg.Nonce == nil {
		nonce := make([]byte, 16)
		if _, err := cryptorand.Read(nonce); err != nil {
			return nil, err
		}
		msg.Nonce = nonce
	}
	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &pb.Signed{
		Payload:      payload,
		Signature:    id.Sign(payload),
		FromDeviceId: id.DeviceID,
	}, nil
}

// Verifier 校验收到的 Signed。
type Verifier struct {
	selfID string
	// publicKeyOf 返回已配对（或正在配对）设备的公钥；nil 表示不认识
	publicKeyOf func(deviceID string) ed25519.PublicKey
	now         func() time.Time

	mu   sync.Mutex
	seen map[string]time.Time
}

func NewVerifier(selfID string, publicKeyOf func(deviceID string) ed25519.PublicKey) *Verifier {
	return &Verifier{
		selfID:      selfID,
		publicKeyOf: publicKeyOf,
		now:         time.Now,
		seen:        make(map[string]time.Time),
	}
}

// Verify 验签并防重放。通过时返回发送方与消息。
//
// 签名不覆盖 from_device_id（服务器会改写它），所以必须按 from 去查公钥：
// 冒用别人的 from 只会让验签失败。
func (v *Verifier) Verify(s *pb.Signed) (string, *pb.SignalMessage, error) {
	from := s.GetFromDeviceId()
	pub := v.publicKeyOf(from)
	if pub == nil {
		return from, nil, ErrUnknownSender
	}
	if !identity.Verify(pub, s.GetPayload(), s.GetSignature()) {
		return from, nil, ErrBadSignature
	}

	var msg pb.SignalMessage
	if err := proto.Unmarshal(s.GetPayload(), &msg); err != nil {
		return from, nil, ErrMalformed
	}
	if msg.GetToDeviceId() != v.selfID {
		return from, nil, ErrWrongRecipient
	}
	ts := time.Unix(msg.GetTimestampUnix(), 0)
	if skew := v.now().Sub(ts); skew > MaxClockSkew || skew < -MaxClockSkew {
		return from, nil, ErrStale
	}
	if v.seenNonce(msg.GetNonce()) {
		return from, nil, ErrReplay
	}
	return from, &msg, nil
}

// seenNonce 记录并判断 nonce 是否重复，同时清理过期条目。
func (v *Verifier) seenNonce(nonce []byte) bool {
	if len(nonce) == 0 {
		return false
	}
	key := string(nonce)

	v.mu.Lock()
	defer v.mu.Unlock()

	now := v.now()
	if _, dup := v.seen[key]; dup {
		return true
	}
	// 超出时钟容忍窗口的 nonce 不可能再被接受，可以安全丢弃
	for k, t := range v.seen {
		if now.Sub(t) > MaxClockSkew*2 {
			delete(v.seen, k)
		}
	}
	v.seen[key] = now
	return false
}

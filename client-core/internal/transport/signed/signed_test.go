package signed

import (
	"crypto/ed25519"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

func newIdentity(t *testing.T) *identity.Identity {
	t.Helper()
	id, err := identity.LoadOrCreate(filepath.Join(t.TempDir(), "identity.json"))
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func keysOf(ids ...*identity.Identity) func(string) ed25519.PublicKey {
	return func(deviceID string) ed25519.PublicKey {
		for _, id := range ids {
			if id.DeviceID == deviceID {
				return id.PublicKey
			}
		}
		return nil
	}
}

func ping() *pb.SignalMessage {
	return &pb.SignalMessage{Payload: &pb.SignalMessage_Ping{Ping: &pb.Ping{SentAtUnix: 1}}}
}

func TestSignAndVerify(t *testing.T) {
	alice, bob := newIdentity(t), newIdentity(t)
	v := NewVerifier(bob.DeviceID, keysOf(alice))

	s, err := Sign(alice, bob.DeviceID, ping())
	if err != nil {
		t.Fatal(err)
	}
	from, msg, err := v.Verify(s)
	if err != nil {
		t.Fatalf("验签失败: %v", err)
	}
	if from != alice.DeviceID || msg.GetPing() == nil {
		t.Fatalf("from=%s msg=%v", from, msg)
	}
}

// 同一条信令从服务器和局域网两条路都到达时，只交给上层一次。
func TestReplayAcrossChannelsIsDropped(t *testing.T) {
	alice, bob := newIdentity(t), newIdentity(t)
	v := NewVerifier(bob.DeviceID, keysOf(alice))

	s, _ := Sign(alice, bob.DeviceID, ping())
	copied := proto.Clone(s).(*pb.Signed)
	if _, _, err := v.Verify(s); err != nil {
		t.Fatal(err)
	}
	if _, _, err := v.Verify(copied); !errors.Is(err, ErrReplay) {
		t.Fatalf("第二份应当被当作重复消息，得到 %v", err)
	}
}

// 局域网里谁都能把截到的消息转发给第三台设备：收件人不是本机的一律拒收。
func TestWrongRecipientIsRejected(t *testing.T) {
	alice, bob, carol := newIdentity(t), newIdentity(t), newIdentity(t)
	carolVerifier := NewVerifier(carol.DeviceID, keysOf(alice))

	s, _ := Sign(alice, bob.DeviceID, ping())
	if _, _, err := carolVerifier.Verify(s); !errors.Is(err, ErrWrongRecipient) {
		t.Fatalf("发给 bob 的消息不该被 carol 接受，得到 %v", err)
	}
}

func TestRejectsUnknownForgedAndStale(t *testing.T) {
	alice, bob, mallory := newIdentity(t), newIdentity(t), newIdentity(t)
	v := NewVerifier(bob.DeviceID, keysOf(alice))

	s, _ := Sign(mallory, bob.DeviceID, ping())
	if _, _, err := v.Verify(s); !errors.Is(err, ErrUnknownSender) {
		t.Errorf("未配对设备: %v", err)
	}

	// mallory 冒用 alice 的 ID：签名对不上 alice 的公钥
	s.FromDeviceId = alice.DeviceID
	if _, _, err := v.Verify(s); !errors.Is(err, ErrBadSignature) {
		t.Errorf("冒用 ID: %v", err)
	}

	old, _ := Sign(alice, bob.DeviceID, ping())
	v.now = func() time.Time { return time.Now().Add(MaxClockSkew + time.Minute) }
	if _, _, err := v.Verify(old); !errors.Is(err, ErrStale) {
		t.Errorf("过期消息: %v", err)
	}
}

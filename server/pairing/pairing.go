// Package pairing 管理配对码的签发与兑换。
//
// 配对码只在内存中存活 5 分钟，服务器不持久化任何设备信息——
// 双方一旦交换到对方公钥，此后的互信就完全不依赖服务器了。
package pairing

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// 配对码字符集刻意剔除易混字符（0/O、1/I/L），因为用户可能要口头转述。
const codeAlphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

const (
	CodeLength = 6
	CodeTTL    = 5 * time.Minute
)

var (
	ErrCodeNotFound = errors.New("配对码不存在或已过期")
	ErrSelfPairing  = errors.New("不能与自己配对")
)

// Peer 是配对双方互相看到的信息。
type Peer struct {
	DeviceID   string
	DeviceName string
	Platform   string
	PublicKey  []byte
}

type pending struct {
	code      string
	initiator Peer
	expiresAt time.Time
	// 兑换方的信息经此通道回传给发起方，让发起方也能拿到对端公钥
	redeemed chan Peer
}

type Broker struct {
	mu    sync.Mutex
	codes map[string]*pending
	// 便于测试时注入假时钟
	now func() time.Time
}

func NewBroker() *Broker {
	return &Broker{
		codes: make(map[string]*pending),
		now:   time.Now,
	}
}

// Create 为发起方签发一个配对码。
// 返回的通道在有人兑换时收到对方信息；超时或 Cancel 时关闭。
func (b *Broker) Create(initiator Peer) (code string, expiresAt time.Time, ch <-chan Peer, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.gcLocked()

	// 极小概率撞码，重试几次即可
	for attempt := 0; attempt < 8; attempt++ {
		c, genErr := generateCode()
		if genErr != nil {
			return "", time.Time{}, nil, genErr
		}
		if _, taken := b.codes[c]; taken {
			continue
		}
		p := &pending{
			code:      c,
			initiator: initiator,
			expiresAt: b.now().Add(CodeTTL),
			redeemed:  make(chan Peer, 1),
		}
		b.codes[c] = p
		return c, p.expiresAt, p.redeemed, nil
	}
	return "", time.Time{}, nil, errors.New("生成配对码失败：连续撞码")
}

// Redeem 用配对码换取发起方信息，同时把兑换方信息回传给发起方。
// 配对码一次性使用，兑换后立即失效。
func (b *Broker) Redeem(code string, redeemer Peer) (Peer, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.gcLocked()

	p, ok := b.codes[code]
	if !ok {
		return Peer{}, ErrCodeNotFound
	}
	if p.initiator.DeviceID == redeemer.DeviceID {
		return Peer{}, ErrSelfPairing
	}

	delete(b.codes, code)
	// 缓冲为 1，发起方即使暂时没在读也不会阻塞这里
	p.redeemed <- redeemer
	close(p.redeemed)

	return p.initiator, nil
}

// Cancel 撤销尚未被兑换的配对码（发起方关闭配对界面或断线时）。
func (b *Broker) Cancel(code string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if p, ok := b.codes[code]; ok {
		delete(b.codes, code)
		close(p.redeemed)
	}
}

// gcLocked 清理过期配对码。调用方必须已持有锁。
func (b *Broker) gcLocked() {
	now := b.now()
	for code, p := range b.codes {
		if now.After(p.expiresAt) {
			delete(b.codes, code)
			close(p.redeemed)
		}
	}
}

func (b *Broker) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.codes)
}

func generateCode() (string, error) {
	out := make([]byte, CodeLength)
	max := big.NewInt(int64(len(codeAlphabet)))
	for i := range out {
		// 用密码学随机源：配对码在有效期内是唯一的准入凭证
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("生成配对码: %w", err)
		}
		out[i] = codeAlphabet[n.Int64()]
	}
	return string(out), nil
}

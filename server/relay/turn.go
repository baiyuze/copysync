// Package relay 是自建的 TURN 中转服务，供 P2P 打洞失败时兜底。
//
// TURN 只转发 UDP 数据包，看不到 DTLS 内层——即便走中转，
// 设备间的内容对服务器仍然是密文。
package relay

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pion/turn/v4"
)

// Server 封装 pion/turn，并按标准的 TURN REST API 方式签发短期凭证。
type Server struct {
	turn     *turn.Server
	secret   []byte
	publicIP string
	port     int
	ttl      time.Duration
	log      *slog.Logger
}

type Options struct {
	// PublicIP 是客户端能访问到的地址。NAT 后部署时必须显式指定，
	// 否则 TURN 会把内网地址写进 relay candidate，外部设备连不上。
	PublicIP string
	// ListenIP 是本机监听的地址，默认 0.0.0.0。服务器有多个 IP 时指定它，
	// 回包的源地址才会与客户端发往的地址一致——否则客户端的 NAT 会把回包当作陌生来源丢掉。
	ListenIP string
	Port     int
	Realm    string
	// Secret 用于派生短期凭证。与签发方共享，不下发给客户端。
	Secret string
	// CredentialTTL 决定签发的凭证多久过期
	CredentialTTL time.Duration
	Logger        *slog.Logger
}

func New(opts Options) (*Server, error) {
	if opts.PublicIP == "" {
		return nil, fmt.Errorf("必须指定 TURN 的公网地址")
	}
	if opts.Port == 0 {
		opts.Port = 3478
	}
	if opts.ListenIP == "" {
		opts.ListenIP = "0.0.0.0"
	}
	if opts.Realm == "" {
		opts.Realm = "copysync"
	}
	if opts.CredentialTTL == 0 {
		opts.CredentialTTL = 12 * time.Hour
	}
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}

	udpConn, err := net.ListenPacket("udp4", net.JoinHostPort(opts.ListenIP, strconv.Itoa(opts.Port)))
	if err != nil {
		return nil, fmt.Errorf("监听 TURN 端口 %d: %w", opts.Port, err)
	}

	secret := []byte(opts.Secret)
	s := &Server{
		secret:   secret,
		publicIP: opts.PublicIP,
		port:     opts.Port,
		ttl:      opts.CredentialTTL,
		log:      log,
	}

	ts, err := turn.NewServer(turn.ServerConfig{
		Realm: opts.Realm,
		// 凭证是 HMAC 派生的，无需维护用户表：
		// username = "<过期时间戳>:<device_id>"，password = HMAC-SHA1(username, secret)
		AuthHandler: func(username, realm string, srcAddr net.Addr) ([]byte, bool) {
			key, ok := s.verify(username, realm)
			if !ok {
				log.Debug("TURN 认证失败", "username", username, "from", srcAddr)
			}
			return key, ok
		},
		PacketConnConfigs: []turn.PacketConnConfig{{
			PacketConn: udpConn,
			RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
				RelayAddress: net.ParseIP(opts.PublicIP),
				Address:      opts.ListenIP,
			},
		}},
	})
	if err != nil {
		udpConn.Close()
		return nil, fmt.Errorf("启动 TURN 服务: %w", err)
	}
	s.turn = ts

	log.Info("TURN 中转已启动", "addr", opts.PublicIP, "port", opts.Port, "realm", opts.Realm)
	return s, nil
}

func (s *Server) Close() error {
	if s.turn == nil {
		return nil
	}
	return s.turn.Close()
}

// Credential 是签发给客户端的一次性 TURN 凭证。
type Credential struct {
	Username  string
	Password  string
	ExpiresAt time.Time
	URLs      []string
}

// Issue 为指定设备签发短期 TURN 凭证。
//
// 采用标准的 TURN REST API 约定：用户名里编码过期时间，
// 密码由服务端密钥 HMAC 派生。服务器因此无需存储任何凭证。
func (s *Server) Issue(deviceID string) Credential {
	expires := time.Now().Add(s.ttl)
	username := fmt.Sprintf("%d:%s", expires.Unix(), deviceID)

	mac := hmac.New(sha1.New, s.secret)
	mac.Write([]byte(username))
	password := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return Credential{
		Username:  username,
		Password:  password,
		ExpiresAt: expires,
		URLs: []string{
			fmt.Sprintf("turn:%s:%d?transport=udp", s.publicIP, s.port),
		},
	}
}

// verify 校验凭证并返回 TURN 所需的鉴权 key。
func (s *Server) verify(username, realm string) ([]byte, bool) {
	parts := strings.SplitN(username, ":", 2)
	if len(parts) != 2 {
		return nil, false
	}
	exp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || time.Now().After(time.Unix(exp, 0)) {
		return nil, false // 格式错误或已过期
	}

	mac := hmac.New(sha1.New, s.secret)
	mac.Write([]byte(username))
	password := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return turn.GenerateAuthKey(username, realm, password), true
}

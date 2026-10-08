module github.com/baiyuze/copysync/client-core

go 1.26.3

require modernc.org/sqlite v1.59.0

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pion/datachannel v1.6.2 // indirect
	github.com/pion/dtls/v3 v3.1.8 // indirect
	github.com/pion/interceptor v0.1.48 // indirect
	github.com/pion/logging v0.2.4 // indirect
	github.com/pion/mdns/v2 v2.2.0 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.17 // indirect
	github.com/pion/rtp v1.10.5 // indirect
	github.com/pion/sctp v1.11.1 // indirect
	github.com/pion/sdp/v3 v3.0.19 // indirect
	github.com/pion/srtp/v3 v3.0.13 // indirect
	github.com/pion/stun/v3 v3.0.1 // indirect
	github.com/pion/transport/v3 v3.0.8 // indirect
	github.com/pion/transport/v4 v4.1.0 // indirect
	github.com/pion/turn/v4 v4.1.4 // indirect
	github.com/pion/turn/v5 v5.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/wlynxg/anet v0.0.5 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	modernc.org/libc v1.75.7 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
)

require (
	github.com/baiyuze/copysync/proto v0.0.0
	github.com/coder/websocket v1.8.15
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

replace github.com/baiyuze/copysync/proto => ../proto

require (
	github.com/baiyuze/copysync/server v0.0.0
	github.com/klauspost/compress v1.20.0
	github.com/pion/ice/v4 v4.4.2
	github.com/pion/stun/v4 v4.0.0
	github.com/pion/webrtc/v4 v4.2.20
	golang.org/x/net v0.57.0
)

replace github.com/baiyuze/copysync/server => ../server

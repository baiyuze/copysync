module github.com/baiyuze/copysync/server

go 1.26.3

replace github.com/baiyuze/copysync/proto => ../proto

require (
	github.com/baiyuze/copysync/proto v0.0.0-00010101000000-000000000000
	github.com/coder/websocket v1.8.15
	github.com/pion/turn/v4 v4.1.4
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/pion/dtls/v3 v3.0.7 // indirect
	github.com/pion/logging v0.2.4 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/stun/v3 v3.0.1 // indirect
	github.com/pion/transport/v3 v3.0.8 // indirect
	github.com/pion/transport/v4 v4.0.1 // indirect
	github.com/wlynxg/anet v0.0.5 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/grpc v1.84.0 // indirect
)

# Building from source

[简体中文](building.zh-CN.md) · [Back to README](../README.md)

The Mac app needs Go 1.26, Flutter 3.44 and Xcode:

```bash
./scripts/build.sh       # build everything into dist/
./scripts/package.sh     # build the DMG and server packages into release/
```

The Windows app needs Go, Flutter, Visual Studio ("Desktop development with C++") and Inno Setup 6. On Windows:

```powershell
powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1   # installer and portable zip into dist\windows\
```

The background service is pure Go and also cross-compiles from a Mac: `GOOS=windows go build ./cmd/copysyncd`. CI builds the Windows installer on every commit.

The server is pure Go too. `package.sh` cross-compiles it for Linux (amd64, arm64), macOS and Windows, and bundles each with its install script from `server/deploy/`.

During development:

```bash
cd client-core && go test ./...          # service tests
./tools/natlab/docker.sh                 # NAT lab: direct vs relay across 11 network topologies
cd ui && flutter test                    # app tests
./scripts/install-macos.sh               # install the service from dist/ as a login item
./dist/copysync-cli status               # service status
./dist/copysync-cli watch                # live event stream
```

After changing a `.proto` file, run `./scripts/gen-proto.sh` (uses buf; install with `go install`, no protoc needed). After changing the UI, run `./tools/screenshots/render.sh` to regenerate the screenshots from demo data.

```
client-core/    background service (Go)
  internal/clipboard    clipboard (Mac: cgo + Objective-C; Windows: Win32; main thread only)
  internal/transport    signaling (WebSocket) and P2P (WebRTC)
  internal/sync         sync engine: routing, packing, transfer, storage
ui/             app (Flutter)
server/         signaling + TURN relay server (Go); deployment files in server/deploy/
proto/          message definitions and device ID derivation, shared by all three
docs/           website (GitHub Pages)
spikes/         technical validation done before development
design/         design documents
tools/natlab/   NAT traversal lab
tools/installer/ Windows installer (Inno Setup)
```

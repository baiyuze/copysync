# Building from source

[简体中文](building.zh-CN.md) · [Back to README](../README.md)

The Mac app needs Go 1.26, Flutter 3.44 and Xcode:

```bash
./scripts/build.sh       # build everything into dist/
./scripts/package.sh     # build the DMGs (Apple silicon, Intel, universal) and server packages into release/
```

The Windows app needs Go, Flutter, Visual Studio ("Desktop development with C++") and Inno Setup 6. On Windows:

```powershell
powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1   # installer and portable zip into dist\windows\
```

The background service is pure Go and also cross-compiles from a Mac: `GOOS=windows go build ./cmd/pasteferryd`. CI builds the Windows installer on every commit.

The Linux app needs Go, Flutter and its Linux toolchain (`sudo apt install clang cmake ninja-build pkg-config libgtk-3-dev`). Flutter can't cross-compile Linux desktop apps, so build on the architecture you want (CI builds both x64 and ARM64 on Ubuntu 22.04):

```bash
./scripts/package-linux.sh   # .deb and tar.gz for this machine's architecture into dist/linux/
```

The Linux clipboard is pure Go (X11 protocol) and also cross-compiles from a Mac: `GOOS=linux go build ./cmd/pasteferryd`. Its tests need an X display: `xvfb-run -a go test ./internal/clipboard/` with `xclip` installed.

The server is pure Go too. `package.sh` cross-compiles it for Linux (amd64, arm64), macOS and Windows, and bundles each with its install script from `server/deploy/`.

During development:

```bash
cd client-core && go test ./...          # service tests
./tools/natlab/docker.sh                 # NAT lab: direct vs relay across 11 network topologies
cd ui && flutter test                    # app tests
./scripts/install-macos.sh               # install the service from dist/ as a login item
./dist/pasteferry-cli status               # service status
./dist/pasteferry-cli watch                # live event stream
```

After changing a `.proto` file, run `./scripts/gen-proto.sh` (uses buf; install with `go install`, no protoc needed). After changing the UI, run `./tools/screenshots/render.sh` to regenerate the screenshots from demo data.

```
client-core/    background service (Go)
  internal/clipboard    clipboard (Mac: cgo + Objective-C; Windows: Win32; Linux: X11; main thread only)
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
tools/linux/    Linux menu entry, icons and install script
```

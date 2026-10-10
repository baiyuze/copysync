# 从源码构建

[English](building.md) · [回到 README](../README.zh-CN.md)

Mac 版需要 Go 1.26、Flutter 3.44 和 Xcode：

```bash
./scripts/build.sh       # 构建全部到 dist/
./scripts/package.sh     # 打出 DMG 与服务器发布包到 release/
```

Windows 版需要 Go、Flutter、Visual Studio（「使用 C++ 的桌面开发」）与 Inno Setup 6，在 Windows 上运行：

```powershell
powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1   # 安装程序与便携版到 dist\windows\
```

后台服务是纯 Go，也可以在 Mac 上交叉编译：`GOOS=windows go build ./cmd/copysyncd`。CI 在每次提交时构建 Windows 安装包。

服务器也是纯 Go。`package.sh` 把它交叉编译成 Linux（amd64、arm64）、macOS 与 Windows 版本，各自带上 `server/deploy/` 里对应的安装脚本。

开发时可以单独运行各部分：

```bash
cd client-core && go test ./...          # 后台服务的测试
./tools/natlab/docker.sh                 # NAT 打洞实验室：11 种网络拓扑下的直连与中转
cd ui && flutter test                    # 界面的测试
./scripts/install-macos.sh               # 把 dist/ 里的后台服务装成登录项
./dist/copysync-cli status               # 查看后台服务状态
./dist/copysync-cli watch                # 实时事件流
```

改了 `.proto` 之后运行 `./scripts/gen-proto.sh` 重新生成代码（用 buf，`go install` 即可安装，无需 protoc）。改了界面之后运行 `./tools/screenshots/render.sh` 重新生成截图，截图用演示数据渲染，不会带出本机的真实内容。

```
client-core/    后台服务（Go）
  internal/clipboard    剪贴板（Mac：cgo + Objective-C；Windows：Win32；都必须跑在主线程）
  internal/transport    信令（WebSocket）与 P2P（WebRTC）
  internal/sync         同步引擎：分流、打包、传输、落盘
ui/             界面（Flutter）
server/         信令 + TURN 中转服务器（Go），部署文件在 server/deploy/
proto/          消息定义与设备 ID 派生规则，三方共用
docs/           项目主页（GitHub Pages）
spikes/         开工前的技术验证
design/         技术方案文档
tools/natlab/   NAT 打洞实验室
tools/installer/ Windows 安装程序（Inno Setup）
```

# CopySync 界面

CopySync 的 macOS 界面（Flutter）：复制记录、设备配对、设置，以及首次启动时启用后台同步。

界面本身不碰剪贴板，也不联网：所有操作都经本机 gRPC 交给后台服务 `copysyncd`
（见 `../client-core`）。连接信息由后台服务写在
`~/Library/Application Support/CopySync/daemon.json`。

```bash
flutter run -d macos     # 开发运行，需要后台服务已在运行（../scripts/install-macos.sh）
flutter test             # 界面测试
```

| 目录 | 内容 |
|---|---|
| `lib/app_state.dart` | 唯一的状态源：连接、断线重连、事件合并 |
| `lib/background_service.dart` | 把内置的后台服务注册为登录项、升级后重启 |
| `lib/pages/` | 复制记录、设备、设置、首次启动引导 |
| `lib/theme.dart` | 颜色、字号、间距 |
| `lib/gen/` | 由 `../proto` 生成的代码，不要手改 |
| `tool/screenshots/` | 用演示数据渲染 README 与网站的截图，见 `../tools/screenshots/render.sh` |

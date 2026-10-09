# 桌面入口、配对弹窗与图片预览验证

2026-10-09。本轮功能代码已包含在工作区提交 `a78913e`。测试环境为 Windows 11 build 26200（Parallels）及 macOS，Flutter 3.44.9。

## 结果

| 功能 | 验证结果 |
|---|---|
| 桌面快捷方式 | 升级安装自动生成 CopySync.lnk，目标指向当前用户安装目录；实际双击入口可启动界面 |
| 托盘图标 | 后台服务持有图标；关闭界面后保留，点击可重开界面 |
| 托盘菜单 | 实际显示“打开 CopySync”“退出 CopySync”；退出后界面、daemon、守护进程均结束，等待 12 秒不重启，图标清理正常；恢复启动后原配对正常 |
| Explorer 重建通知 | 发送 TaskbarCreated 后图标仍可用；未强制重启用户的 Explorer |
| 配对弹窗 | Windows 实际生成配对码，临时 Mac 端兑换，Windows 核对指纹并确认后返回设备页；显示配对成功且没有残留弹窗 |
| 输入配对码的竞态 | Mac 和 Windows 上的回归测试覆盖 RPC/事件先后顺序、重复事件、接受/拒绝、过期和处理旧会话时新请求到达，均只显示一个确认框 |
| 图片预览 | Windows 新版界面实际点击历史图片后显示图片；支持缩放、拖动、适应窗口和关闭 |
| 图片异常与下载 | 两个平台的测试覆盖远端下载后显示、下载失败重试、读取失败重试、记录删除/过期、文件缺失/损坏；全部通过 |
| 剪贴板隔离 | 预览路径只读；“下载并预览”专用请求在完成接收时不写剪贴板，后台测试实际调用接收完成路径验证 |
| 图片接口 | Windows 实际运行 gRPC 和存储测试：无 token 拒绝、空 ID 拒绝、缺失文件提示、中文路径返回正确 |

## 自动检查

- macOS：`flutter analyze` 无问题，`flutter test` 全部 **26 项通过**；`flutter build macos --release` 成功。
- Windows 虚拟机：`flutter test` 全部 **26 项通过**，图片 RPC/存储/下载测试全部通过。
- 后台：`go vet ./...`、`go test ./...` 通过；Windows 目标 `go vet ./...` 和 daemon 交叉编译通过。
- `git diff --check` 通过。

真实鼠标点击验证了 Windows 的生成配对码与确认流程。输入端的完整流程由两平台的界面回归测试覆盖；未将命令行兑换冒充为界面输入测试。

## 产物与恢复

Windows 本地测试包已安装：新的 Dart AOT/资源和同版本 Flutter 引擎，配合原 CI 的未修改原生 runner、新编译 daemon 及当前 Inno Setup 安装脚本。版本仍为 1.1.2，属于本地验证产物，发布应使用正式发布流程重建。

- 安装包：`~/Library/Caches/copysync-tray-20261009/artifacts/CopySync-Setup.exe`
- SHA-256：`0b8fb27c7449b252df392bb21b914c16bcd0d4d03b42bf834655474c80a6f68a`
- 便携包 SHA-256：`206a38abdc37e94418c8b3ebd63cf9a801816c3ec63bfc0be30f06808c42b722`
- 日志与截图：`~/Library/Caches/copysync-tray-20261009/` 和 `~/Library/Caches/copysync-validation-20261009/`。

临时测试设备已从 Windows 解除配对，临时 Mac daemon 已停止，临时编译器已卸载；Windows 与用户原 Mac 恢复直连。Mac 正在使用的生产 App/daemon 未替换，Mac 新界面修复需安装新构建后生效。Windows 桌面由 Parallels 共享，快捷方式位于其映射的桌面目录。托盘图标可能由 Windows 收纳在右下角“^”内。

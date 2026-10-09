# Windows 客户端技术方案

> 状态：已实现，随 1.2.0 发布（2026-10-09）。实机验证记录见 [windows-validation-2026-10-09.md](windows-validation-2026-10-09.md)
> 目标版本：1.2.0（Mac 与 Windows 同时发布）
> 相关：[多出口网络的 NAT 打洞](nat-traversal.md)

## 一、目标

让 Windows 电脑加入 CopySync：与 Mac、与其他 Windows 电脑互相同步文本、带格式文本、图片和文件，体验与 Mac 版一致。

- 在一台上 `Ctrl+C`，另一台上 `Ctrl+V`（或 `⌘V`）。粘贴文件时得到的是真实文件，可以直接粘贴进资源管理器、微信、Office。
- 开机自启，关掉窗口照常同步，崩溃自动拉起。Windows 安装时创建桌面快捷方式；后台服务持有托盘图标，点击可打开界面，右键可打开界面或退出 CopySync。
- 网络行为与 Mac 版完全相同：直连优先、多出口打洞、自动中转、端到端加密。
- 与 Mac 版互通，Mac 端不升级也能和 Windows 配对、同步（协议不变）。

**一期不做：**

| 不做 | 原因 | 何时 |
|---|---|---|
| ARM64 原生包 | Flutter 3.47 才正式支持 Windows ARM64，且只能在 ARM 机器上构建；x64 版在 ARM 上可以模拟运行 | 二期，用 GitHub 的 ARM Windows 构建机 |
| 虚拟文件（Outlook 附件、压缩包里直接复制的文件） | 剪贴板里没有真实路径，只有流，需要实现 COM 接口 | 视需求 |
| 自动更新、微软商店 | 先跑通分发与签名 | 二期 |

支持的系统：Windows 10（21H2 及以上）与 Windows 11，x64。这是 Flutter 桌面版支持的范围。

## 二、现状：哪些能直接用

后台服务除剪贴板外全部是纯 Go：数据库用的是纯 Go 的 `modernc.org/sqlite`，压缩是纯 Go 的 zstd，cgo 只出现在 `internal/clipboard/darwin.go`。实测在 Mac 上已经可以直接交叉编译出 Windows 版：

```
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./cmd/...    # copysyncd.exe、copysync-cli.exe
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build ./cmd/copysyncd
```

| 模块 | Windows 上 | 说明 |
|---|---|---|
| 信令、配对、签名（`transport/signaling`、`peers`、`identity`） | 直接复用 | |
| P2P、出口探测、打洞、中转（`transport/p2p`） | 直接复用 | 防火墙与虚拟网卡见第七节 |
| 同步引擎、补发、存储、缓存、回收（`sync`、`store`、`cache`、`gc`） | 直接复用 | |
| 本机接口（`rpc`，127.0.0.1 + 一次性 token） | 直接复用 | |
| 数据目录（`config.DefaultPaths`） | 已有 Windows 分支 | `%LOCALAPPDATA%\CopySync` |
| 打包解包（`pack`） | 要改 | 文件名、符号链接，见第五节 |
| 剪贴板（`clipboard`） | **要新写** | 目前是返回「不支持」的占位实现，见第四节 |
| 进程生命周期（`cmd/copysyncd`） | 要改 | 自启、守护、日志、单实例、注销，见第六节 |
| 界面（`ui/`） | 要改 | 已有 Flutter 的 `windows/` 工程骨架；后台服务管理、窗口外观、文案，见第八节 |
| 构建与发布（`scripts/`、CI） | 要新写 | Flutter 的 Windows 版只能在 Windows 上构建，见第十节 |

## 三、总体架构

与 Mac 版同构，三个可执行文件：

```
CopySync.exe（Flutter 界面）
   │  gRPC，127.0.0.1 随机端口 + token（%LOCALAPPDATA%\CopySync\daemon.json）
   ▼
copysyncd.exe（Go 后台服务，用户会话内运行，无控制台窗口）
   ├─ 剪贴板：隐藏的消息窗口 + Win32 剪贴板 API
   ├─ 同步引擎 / 存储 / 缓存
   └─ 信令（WebSocket）+ P2P（WebRTC）──► 信令服务器 / 对端设备

copysync-cli.exe（命令行，排障用）
```

后台服务必须作为**登录用户的普通进程**运行，不能做成 Windows 服务。Windows 服务运行在会话 0，与用户桌面隔离，访问不到用户的剪贴板。这与 Mac 上用 LaunchAgent、不用 LaunchDaemon 是同一个原因。

## 四、剪贴板

这是一期的主要工作量。用纯 Go 实现（`golang.org/x/sys/windows` 加按需加载的 user32、kernel32、shell32 函数），不用 cgo，这样仍能在 Mac 上交叉编译。

### 1. 线程与消息循环

现有的 `MainLoop` 设计可以原样沿用，只是「事件循环」换成 Windows 的消息泵：

- 主线程（已 `LockOSThread`）创建一个**隐藏的顶层窗口**作为剪贴板的拥有者。不用「仅消息窗口」（HWND_MESSAGE），因为它收不到会话结束、电源变化等广播消息。
- `pumpRunLoop` 改为 `PeekMessage`/`DispatchMessage` 循环，空闲时 `MsgWaitForMultipleObjects` 等待，不空转。
- 用 `AddClipboardFormatListener` 订阅 `WM_CLIPBOARDUPDATE`，变化时立即触发一次探测，比 Mac 的 300 毫秒轮询更跟手；轮询保留作兜底。
- `ChangeCount` 用 `GetClipboardSequenceNumber()`。它与 Mac 的 changeCount 语义一致，`Watcher` 里「跳过本机写入」（selfWrites）的逻辑不用改。
- 权限：Windows 没有剪贴板授权机制，`Permission()` 恒为 `PermissionNotApplicable`，界面上的授权提示自然不显示。

### 2. 格式对照

| 类型 | 读取（优先级从高到低） | 写入 |
|---|---|---|
| 文件 | `CF_HDROP`（`DragQueryFileW` 取路径） | `CF_HDROP`（`DROPFILES` 结构 + 双 NUL 结尾的 UTF-16 路径表），加上 `Preferred DropEffect = DROPEFFECT_COPY`，资源管理器粘贴时执行复制而不是移动 |
| 图片 | 注册格式 `PNG`（Chrome、Office、截图工具都提供）→ `CF_DIBV5` → `CF_DIB`，后两者转成 PNG | `PNG` + `CF_DIBV5`（带透明通道）；`CF_DIB`、`CF_BITMAP` 由系统自动合成 |
| 带格式文本 | 注册格式 `HTML Format`（CF_HTML），解析头部偏移量取出 HTML | `HTML Format`（生成带偏移量的头部）+ `CF_UNICODETEXT` |
| 纯文本 | `CF_UNICODETEXT` | `CF_UNICODETEXT`（`CF_TEXT` 由系统合成） |

判型顺序：文件最优先（资源管理器复制文件时剪贴板里同时有文件名文本）；只有图片、没有文本时才算图片（截图、画图、浏览器里「复制图片」）；Excel、Word 复制的内容同时带着一份渲染出来的图片，但有文本，按带格式文本处理。

**文件粘贴不需要 COM。**占位实现的注释里写着「文件粘贴需要在 Go 里实现 COM IDataObject」，那是按延迟渲染设想的。Mac 版已经验证过延迟渲染行不通，改成了「先落盘、再写剪贴板」；Windows 也一样，收到的文件已在本地缓存里，直接给 `CF_HDROP` 真实路径即可。

### 3. 不同步敏感内容

密码管理器复制密码时，会在剪贴板里放一个标记，告诉剪贴板工具「别记录、别同步」。Windows 上的约定（微软文档「Cloud Clipboard and Clipboard History Formats」，以及老的通行做法）：

| 标记 | 规则 | 谁在用 |
|---|---|---|
| `ExcludeClipboardContentFromMonitorProcessing` | 存在即跳过 | KeePass、KeePassXC |
| `Clipboard Viewer Ignore` | 存在即跳过 | KeePass 等较老的工具 |
| `CanIncludeInClipboardHistory` | 值为 DWORD 0 时跳过（只看存在与否不够） | Bitwarden 只写这一个 |
| `CanUploadToCloudClipboard` | 值为 0 时跳过：它的含义正是「不要同步到别的设备」 | |

命中就既不同步，也不写进复制记录。**Mac 版目前没有做这件事**（应检查 `org.nspasteboard.ConcealedType` 和 `org.nspasteboard.TransientType`），1Password、Bitwarden 复制的密码会被同步过去并留在记录里。这一条应先在 Mac 版单独修掉，不等 Windows。

### 4. 其他细节

- **剪贴板被占用**：`OpenClipboard` 在别的程序正打开剪贴板时会失败，按 20 毫秒间隔重试 10 次，仍失败就放弃这一次，等下一次变化。
- **写入时的标记**：把对端内容写进本机剪贴板时，加上 `CanUploadToCloudClipboard = 0`，免得 Windows 自带的云剪贴板再把它同步一遍，在用户的其他 Windows 电脑上出现重复。写入本机剪贴板历史（`Win+V`）不受影响。
- **不支持的内容**：只有虚拟文件（`FileGroupDescriptorW`/`FileContents`）而没有 `CF_HDROP` 时判为未知类型，跳过并记日志。

## 五、跨平台内容兼容

网络上传输的格式保持现状（与 Mac 版一致），差异都在 Windows 一侧的剪贴板层与解包时消化：

| 项目 | 网络上的约定 | Windows 侧的处理 |
|---|---|---|
| 文本换行 | `LF` | 读取时把 `CRLF` 统一成 `LF`；写入时转成 `CRLF`（记事本等程序依赖它） |
| 带格式文本 | 完整的 HTML 文本 | 读取时去掉 CF_HTML 头部；写入时生成头部，用 `<!--StartFragment-->` 标出片段 |
| 图片 | PNG | DIB（倒序扫描行、`BI_BITFIELDS`、32 位透明通道）与 PNG 互转 |
| 条目类型（`ClipItem.content_type`） | 沿用 Mac 的 UTType 字符串，如 `public.png` | Windows 侧按扩展名映射成同样的字符串，不另起一套 |
| 文件与文件夹 | tar + zstd | 解包时按下表处理文件名 |

**文件名**是最容易出事的地方。Mac 上合法的文件名在 Windows 上未必合法，处理不好轻则解包失败，重则有安全问题：

| 情况 | 例子 | 处理 |
|---|---|---|
| Windows 不允许的字符 `: * ? " < > \| \` | `会议:纪要.txt` | 换成对应的全角字符（`：＊？＂＜＞｜＼`），名字仍然可读。冒号还关系到安全：`a.txt:stream` 会写进 NTFS 的备用数据流 |
| 结尾是点或空格 | `notes.` | 末尾加 `_` |
| 设备保留名 | `CON`、`nul.txt`、`COM1` | 前面加 `_` |
| 只差大小写 | `Readme.md` 与 `README.md` | 第二个改名为 `README (2).md`，不能互相覆盖 |
| 分解形式的 Unicode | Mac 上的 `が` 可能是 `か` + 濁点两个码位 | 统一成组合形式（NFC），否则资源管理器里显示成两个字符、搜索也搜不到。中文不受影响 |
| 路径超过 260 个字符 | 深层文件夹 | Go 的 `os` 包会自动加 `\\?\` 前缀，读写没问题；只是资源管理器对超长路径支持不好，记一条已知限制 |
| 符号链接 | 文件夹里的软链接 | 普通用户没有创建符号链接的权限（除非开启了开发者模式）。改为跳过并记录，不能让整次解包失败 |
| 系统垃圾文件 | `.DS_Store`、`Thumbs.db`、`desktop.ini` | 打包时跳过（Mac 版同步修改） |

`pack.safeJoin` 的防目录穿越检查保留，在**清洗文件名之后**做；Windows 上额外拒绝驱动器号开头（`C:foo`）和 `\\` 开头的路径。

## 六、后台服务的生命周期

| 需求 | Mac 版的做法 | Windows 版的做法 |
|---|---|---|
| 开机自启 | LaunchAgent | `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` 下写一项。只影响当前用户，不需要管理员权限 |
| 崩溃自动拉起 | launchd `KeepAlive`，10 秒节流 | `copysyncd.exe` 自带守护模式：父进程拉起子进程，子进程异常退出就 10 秒后重启 |
| 只运行一个 | launchd 保证 | 命名互斥量 `Local\CopySync.Daemon`，第二个实例直接退出 |
| 不弹黑框 | 本来就没有 | 以 GUI 子系统编译（`-ldflags -H=windowsgui`） |
| 日志 | launchd 把标准输出重定向到文件 | 自己写 `%LOCALAPPDATA%\CopySync\logs\daemon.log`，按大小轮转（Mac 版的日志目前不轮转，一并加上） |
| 注销、关机 | SIGTERM | 隐藏窗口处理 `WM_QUERYENDSESSION`/`WM_ENDSESSION`，正常收尾 |
| 睡眠唤醒 | 靠重连巡检 | 处理 `WM_POWERBROADCAST`：唤醒后立即重新探测出口、重连，不等巡检 |
| 升级 | 打开新版 App 时校正登录项并重启服务 | 安装程序借助「重启管理器」关掉旧进程，装完再启动；界面同样校验版本、必要时重启服务 |

也评估过「计划任务」：它自带失败重启，但创建「登录时运行」的任务通常要管理员权限，与「不需要管理员即可安装」的目标冲突，所以不用。

## 七、网络

打洞、多出口探测、中转、补发的代码原样复用。Windows 上要注意三件事：

1. **防火墙弹窗**。后台服务在 0.0.0.0 上收发 UDP，Windows 防火墙第一次会弹窗询问是否允许。允许最好；即使用户拒绝，防火墙对 UDP 也是有状态的，本端先往对端发过包，对端的回包就能进来。打洞时双方本来就会互相发包，所以预计仍能直连，**需要实测确认**。安装程序不加防火墙规则，因为那需要管理员权限。
2. **虚拟网卡**。装了 WSL、Hyper-V、Docker 或 VPN 的电脑有很多虚拟网卡，其中 WSL 的网卡地址每次开机都会变。出口历史是按「本机各网卡的地址」区分网络的（`localNetworkKey`），会因此失效。Windows 上计算时排除名字以 `vEthernet` 开头的网卡和已知的 VPN 虚拟网卡。
3. **「公用网络」**。连上咖啡馆等网络时，Windows 默认把它设为「公用网络」，防火墙更严。行为同第 1 条，实测时覆盖。

## 八、界面

继续用 Flutter，一套代码。Flutter 的 `windows/` 工程骨架已经在仓库里，`daemon_client.dart` 也已有 Windows 的数据目录分支。

### 1. 外观：像 Windows 11 自带的应用

Mac 版刻意做成系统自带应用的样子（石墨 + 系统蓝、透明标题栏、侧边栏延伸到红绿灯下）。Windows 版沿用同一个原则，参照 Windows 11 的「设置」应用：

| | Mac | Windows |
|---|---|---|
| 窗口 | 透明标题栏，红绿灯在左上 | 自绘标题栏，最小化、最大化、关闭在右上；侧边栏顶部放应用名，与「设置」应用一致 |
| 背景 | 系统材质 | Mica 材质（在 `windows/runner` 里调用 `DwmSetWindowAttribute`，不引入插件） |
| 字体 | SF Pro、苹方（跟随系统） | Segoe UI Variable、微软雅黑 UI（跟随系统） |
| 图标 | Cupertino | Fluent 图标 |
| 强调色 | 系统蓝 | 跟随系统强调色，取不到时用同一个蓝 |

实现上在 `theme.dart` 加一层平台差异（字体、图标、标题栏高度、圆角），页面结构与交互不变：复制记录、设备、设置三页。

### 2. 后台服务管理

`background_service.dart` 目前全是 launchd 的逻辑，抽成接口，分 Mac 与 Windows 两个实现。Windows 的状态判定：

| 状态 | 判定 |
|---|---|
| 未安装 | `Run` 下没有 CopySync 项 |
| 指向旧位置 | `Run` 项的路径不是当前这份 `copysyncd.exe` |
| 已停止 | `daemon.json` 不存在，或其中的进程号已不在运行 |
| 运行中 | 能连上，且版本与界面一致（不一致就重启服务） |

### 3. 文案与小处

- 「这台 Mac」「两台 Mac」等十余处文案改为按平台显示（「这台电脑」「两台设备」）。Mac 版同步修改，因为对端可能是 Windows。
- 设备列表的平台图标已有 `windows` 分支，换成 Fluent 图标即可。
- 打开网页用 `open` 命令的地方改为跨平台的写法。
- 截图脚本（`tools/screenshots`）加 Windows 外观的版本，用于 README 与网站。

## 九、安全

| 项目 | 做法 |
|---|---|
| 设备私钥（`identity.json`） | 用 DPAPI（`CryptProtectData`，当前用户范围）加密后落盘。别的用户、拷走文件到别的电脑都解不开。Mac 版目前是权限 0600 的明文文件，后续可迁到钥匙串 |
| 数据目录 | `%LOCALAPPDATA%` 本身只有当前用户可访问，不另设权限 |
| 本机接口 | 与 Mac 相同：只监听 127.0.0.1，必须带 token |
| 密码等敏感内容 | 见第四节第 3 条：命中标记就不同步、不记录 |
| 收到的文件 | 一期不加「来自网络」标记（MOTW）：内容来自自己的设备，加了会让 Office 以受保护视图打开、可执行文件弹警告。二期考虑「源文件带标记才带过去」 |

## 十、安装、签名与发布

### 1. 安装包

用 Inno Setup 做**按用户安装**的安装包（类似 VS Code 的「用户安装版」）：

- 装到 `%LOCALAPPDATA%\Programs\CopySync`，不需要管理员权限，不弹 UAC。
- 创建开始菜单项；装完直接打开 CopySync，由界面引导「启用后台同步」，流程与 Mac 一致。
- 卸载时停掉后台服务、删掉自启项，询问是否保留复制记录与设备身份。
- 另出一个 zip 便携版，给不想运行安装程序的人。

也评估过 MSIX：好处是安装干净、能自动更新；但要求必须有受信任的签名证书，且应用写入 `%LOCALAPPDATA%` 的数据会被重定向到包的私有目录，命令行工具与后台服务对不上路径。一期不用。

### 2. 代码签名

不签名的安装包在下载、运行时会被 SmartScreen 拦一下（「Windows 已保护你的电脑」→「更多信息」→「仍要运行」），与 Mac 版没经过公证时「仍要打开」是一回事。可选的路子：

| 方案 | 费用 | 效果 |
|---|---|---|
| 不签名（一期） | 0 | 每个新版本都被 SmartScreen 拦，README 里写清楚怎么放行 |
| 普通代码签名证书（OV） | 每年数百美元起 | 有发布者名称；SmartScreen 信誉要随下载量慢慢积累 |
| 微软的云签名服务 | 按月计费 | 价格低，但对个人开发者有地区限制，需要确认是否可用 |

### 3. 构建与发布

Flutter 的 Windows 版只能在 Windows 上构建，现有的 `scripts/package.sh` 跑在 Mac 上，所以 Windows 包由 CI 构建：

- 新增 GitHub Actions 任务（`windows-latest`）：编译 Go（`GOOS=windows`）→ `flutter build windows` → Inno Setup 打包 → 上传到同一个 Release。
- 发布时在 Mac 上照常运行 `package.sh` 出 DMG，推送标签后 CI 补上 Windows 包。
- 发布资源沿用「不带版本号」的命名，下载直链不会失效：`CopySync-Setup.exe`、`CopySync-windows-x64.zip`。
- CI 的常规检查增加 Windows：`go vet`、`go test`、`flutter analyze`、`flutter test`。

## 十一、测试与验证

1. **纯函数单测**（任何系统上都能跑）：CF_HTML 头部的生成与解析、DIB 与 PNG 互转（含倒序扫描行、透明通道）、`DROPFILES` 编解码、换行转换、文件名清洗（第五节每一行一个用例）、敏感标记判定。
2. **Windows 上的集成测试**：真实剪贴板的读写往返（文本、HTML、图片、文件）、跳过本机写入、密码管理器标记。GitHub 的 Windows 构建机能否访问剪贴板需要先验证；不行的话这部分在本地 Windows 机器上跑。
3. **Mac ↔ Windows 联调清单**：
   - 内容：中英文文本；从网页、Word、Excel 复制的带格式文本；截图、带透明背景的 PNG；单个文件、多个文件、深层文件夹、中文与日文文件名、含 `:` 的文件名；1 GB 大文件按需拉取。
   - 粘贴目标：资源管理器、微信、Office、记事本。
   - 网络：同一局域网、公司多出口网络、家用宽带、防火墙弹窗点「取消」、睡眠唤醒、断网恢复。
   - 生命周期：注销与重启、升级安装、卸载后重装。
4. **NAT 打洞实验室**：网络代码与平台无关，现有 11 个场景继续覆盖。

## 十二、里程碑

按一个人全职估算，约 6 周：

| 阶段 | 内容 | 产出 | 估时 |
|---|---|---|---|
| M1 | 剪贴板实现（第四节）与后台服务生命周期（第六节） | `copysyncd.exe` 在 Windows 上跑通：用 `copysync-cli` 能与 Mac 配对、互相同步文本与图片 | 2 周 |
| M2 | 文件与内容兼容（第五节）、敏感内容过滤 | 联调清单里的内容类全部通过 | 1 周 |
| M3 | 界面适配（第八节） | Windows 外观的界面，能启用、管理后台服务 | 1.5 周 |
| M4 | 安装包、CI、文档与网站 | Release 里出现 Windows 安装包，README 与网站补上 Windows | 1 周 |
| M5 | 实机联调与修复，发布 1.2.0 | | 0.5 周起 |

M1 结束时就能验证最大的不确定点：剪贴板方案、防火墙行为、Windows 上的打洞。之后的阶段风险较小。

## 十三、需要拍板的事

1. **代码签名**：一期不签名（推荐，先验证需求），还是一开始就买证书？
2. **托盘图标**：已按用户需求加入一期，由后台服务持有，随登录自启、关闭主窗口后保留；退出时停止当前后台同步，但保留下次登录的自启设置。
3. **ARM64**：一期只出 x64（推荐，ARM 机器可以模拟运行），还是同时出 ARM64？
4. **Mac 版的敏感内容漏洞**（第四节第 3 条）：建议马上单独修，发 1.1.3。

## 附：要改动的代码

| 位置 | 改动 |
|---|---|
| `client-core/internal/clipboard/windows.go`（新） | Win32 剪贴板实现、消息窗口、格式转换 |
| `client-core/internal/clipboard/other.go` | 构建标签改为排除 darwin 与 windows |
| `client-core/internal/clipboard/sensitive.go`（新） | 敏感内容判定（两个平台共用接口） |
| `client-core/internal/pack/pack.go` | 文件名清洗、符号链接、垃圾文件过滤 |
| `client-core/internal/transport/p2p/egress_history.go` | 网络标识排除虚拟网卡 |
| `client-core/internal/identity/identity.go` | Windows 上用 DPAPI 保护私钥 |
| `client-core/cmd/copysyncd/main.go` | 守护模式、单实例、日志文件、会话结束与电源事件 |
| `ui/lib/background_service.dart` | 拆成接口与 Mac、Windows 两个实现 |
| `ui/lib/theme.dart`、`ui/lib/main.dart` | 平台外观、自绘标题栏 |
| `ui/windows/runner/` | Mica 背景、窗口尺寸记忆、应用图标 |
| `scripts/`、`.github/workflows/` | Windows 构建、安装包、发布 |
| `tools/installer/`（新） | Inno Setup 脚本 |

# Linux 客户端

> 状态：已实现，随 1.6.0 发布
> 代码：`client-core/internal/clipboard/linux.go`（剪贴板）、`client-core/cmd/pasteferryd/platform_linux.go` 与 `supervisor.go`（后台服务）、`ui/linux/`（GTK 外壳）、`ui/lib/background_service_linux.dart`（开机自启）、`scripts/package-linux.sh` 与 `tools/linux/`（打包）
> 验证：`internal/clipboard/linux_test.go`（Xvfb + xclip、Wayland + wl-clipboard）、CI 的 Linux 任务（x64 与 ARM64）、Ubuntu 24.04 GNOME（Wayland）虚拟机实测
> 相关：[Windows 客户端](windows-client.md)、[局域网直连](lan-mode.md)

## 一、目标

让 Linux 电脑加入隔空贴，与 Mac、Windows 互相同步文本、带格式文本、图片和文件，体验与另外两个平台一致：开机自启、关掉窗口照常同步、崩溃自动拉起，网络行为完全相同。

支持：x64 与 ARM64；X11 桌面，或者带 XWayland 的 Wayland 桌面（GNOME、KDE 等主流桌面默认都带）；GTK 3。

**不做：**

| 不做 | 原因 |
|---|---|
| 托盘图标 | GNOME 默认不显示托盘图标；与 Mac 一样，从应用菜单打开窗口 |
| 没有 XWayland 的纯 Wayland 桌面 | Wayland 不允许后台程序读剪贴板；wlr-data-control 只有 wlroots 系与 KDE 支持，GNOME 不支持。见第二节 |
| AppImage、Flatpak、Snap | AppImage 每次运行挂载到临时目录，开机自启要另行处理；Flatpak 的沙盒拿不到剪贴板。先提供 .deb 与 tar.gz |
| HEIC、TIFF 预览 | Mac 上靠系统的 ImageIO 转换，Linux 上没有现成的，按需再做 |

## 二、剪贴板

### 1. 走 X11，Wayland 上靠 XWayland 桥接

X11 的剪贴板是「谁复制谁持有」：没有一块共享的内容，读的时候请持有者把内容转换成某种格式放到自己窗口的属性上。实现是纯 Go 的 X11 协议（`github.com/jezek/xgb`），不链接 libX11，也不调用 xclip 之类的外部程序。后台服务因此是静态链接的纯 Go 程序，不依赖发行版的库版本。

Wayland 的设计不允许后台程序读剪贴板（只有拿到键盘焦点的窗口能读）。但 GNOME、KDE 都带着 XWayland，混成器在 Wayland 剪贴板与 X11 的 CLIPBOARD 选区之间双向同步：Wayland 程序复制，混成器就成为 X11 选区的持有者，X11 程序能收到变化通知、能读；X11 程序成为持有者，Wayland 程序粘贴时混成器替它去取。所以同一套 X11 实现在 X11 和 Wayland 会话里都能用。

在 Ubuntu 24.04 的 GNOME（Wayland）上实测（`TestLinuxWaylandBridge`，用 wl-copy、wl-paste 扮演原生 Wayland 程序）：两个方向的文本、图片都通，变化通知即时到达。

### 2. 变化、读取与写入

| 做什么 | 怎么做 |
|---|---|
| 知道剪贴板变了 | XFixes 订阅 CLIPBOARD 的持有者变化（设置、窗口销毁、客户端断开），每次计数加一并通知 Watcher 立即探测，不用轮询内容 |
| 判断类型 | 问持有者要 TARGETS（能提供的格式列表），按规则分类；同一次变化只问一次 |
| 读 | ConvertSelection 请持有者把内容放到本机窗口的属性上；超过请求上限的走 INCR 分段 |
| 写 | 本机成为持有者，在 SelectionRequest 里按请求的格式回应；大块数据同样走 INCR |

**写入时的计数：** Watcher 写入后马上 Peek，用新的计数记下「这是自己写的」。所以写入成功时本地先把计数加一，XFixes 随后报来的、持有者是本机窗口的那次变化不再计数。

**INCR 的两个坑：**

- X 协议不开 BIG-REQUESTS 时单个请求最多 262140 字节，还要扣掉请求头。一开始按 256 KB 整块写属性，超长的请求把整条连接写坏，之后所有请求都卡住。分段改为 128 KB。
- 持有者写入 INCR 标记时会产生一次「属性有新值」的事件，等本机读完标记（顺带删掉）再去等下一段，先拿到的就是这个旧事件，此时属性已不存在，读出来是空的，被当成了结束。属性不存在（类型为 None）的事件跳过即可：真正的结束段有类型、长度为零。

### 3. 格式对照

| 内容 | 读 | 写 |
|---|---|---|
| 文本 | UTF8_STRING、text/plain;charset=utf-8、text/plain、TEXT、STRING（Latin-1，但不少程序放的是 UTF-8，按 UTF-8 能解就按 UTF-8） | UTF8_STRING、text/plain;charset=utf-8、text/plain、TEXT；纯 ASCII 时再加 STRING |
| 带格式文本 | text/html（Firefox 给的是带 BOM 的 UTF-16） | text/html（UTF-8）加上文本 |
| 图片 | image/png；只有 JPEG、GIF、BMP 时转成 PNG（BMP 去掉文件头后就是 Windows 的 DIB，复用那边的解码） | image/png |
| 文件 | x-special/gnome-copied-files、text/uri-list，只取还存在的本地文件 | x-special/gnome-copied-files（第一行 copy：文件管理器粘贴时复制而不是移动）、text/uri-list、application/x-kde-cutselection=0，再加一份路径文本 |

分类规则与 Windows 一致：文件优先；只有图片没有文本时才算图片（办公软件复制文字时会附带一张渲染出来的图）；再是 HTML、文本。浏览器复制网页内容时也可能带 text/uri-list，所以只有 uri-list、同时没有 HTML 和图片时才当作文件；读出来不是本地文件（比如网址）就按文本处理。

KeePassXC 等密码管理器复制密码时带着 KDE 约定的 `x-kde-passwordManagerHint`，看到就标为敏感：不同步，也不进复制记录（与 Windows 上的做法相同）。

### 4. 会话结束

与 X 服务器的连接断开几乎总是因为用户注销。这时按会话结束处理（与 Windows 收到 WM_ENDSESSION 一样），后台服务正常退出，守护进程见到退出码 0 也随之结束；下次登录由开机自启重新拉起，而不是留下一个连不上新会话的后台服务。

## 三、后台服务

| | Mac | Windows | Linux |
|---|---|---|---|
| 开机自启 | LaunchAgent | HKCU\...\Run | XDG 自启动项（`~/.config/autostart/pasteferry-daemon.desktop`） |
| 崩溃拉起 | launchd 的 KeepAlive | `pasteferryd -supervise` 守护进程 | 同 Windows |
| 只运行一个 | launchd | 命名互斥量 | `flock`（数据目录下的 daemon.lock、supervisor.lock） |
| 日志 | launchd 重定向 | `logs/daemon.log`，10 MB 轮转 | 同 Windows |
| 数据目录 | `~/Library/Application Support/PasteFerry` | `%LOCALAPPDATA%\PasteFerry` | `$XDG_DATA_HOME/pasteferry`，默认 `~/.local/share/pasteferry` |

用 XDG 自启动项而不是 systemd 的用户服务：自启动项在登录后的图形会话里启动，带着 DISPLAY、XAUTHORITY，后台服务才连得上剪贴板；systemd 用户服务不一定拿得到这些环境变量，各发行版的做法也不一致。GNOME、KDE、Xfce 都认 XDG 自启动项。

守护进程与 Windows 共用（`supervisor.go`）。Linux 上子进程带 `Pdeathsig`：守护进程被结束时后台服务随之结束，界面重启服务时只要结束守护进程，不会留下没人看着的后台服务。

## 四、界面

- **GTK 外壳**（`ui/linux/`）：应用 ID `com.pasteferry.app`（与 Mac 的包名一致，桌面据此把窗口与应用菜单项对应起来，Dock 里显示正确的图标）；默认 960×640、最小 760×520；只运行一个实例，再次打开时激活已有的窗口；窗口图标随程序发布。为了在 Ubuntu 22.04（GLib 2.72）上编译，不用 2.74 才有的 `G_APPLICATION_DEFAULT_FLAGS`。
- **开机自启**（`background_service_linux.dart`）：写自启动项，路径按 Desktop Entry 规范加引号转义（文件夹名可能带空格）；状态按自启动项里的路径与 `pgrep` 判断；重启时结束当前用户的 pasteferryd 再脱离界面启动守护进程，环境变量随之传下去。
- **字体**：主字体跟随系统（fontconfig 默认的无衬线字体）；汉字指定 Noto Sans CJK SC 等（日文界面用 JP），理由与 Windows 相同：系统逐字回退时，中日共用的汉字可能落到另一种语言的字形上。等宽字体交给 fontconfig 的 `monospace` 别名。
- **叫法与图标**：与 Windows 一样说「这台电脑」，用 Fluent 图标（线条风格与 GNOME、KDE 的图标接近）；打开网址用 `xdg-open`。

## 五、打包与发布

`scripts/package-linux.sh` 在 Linux 上构建后台服务、命令行与界面，产出：

- `PasteFerry-linux-<arch>.deb`：装在 `/opt/PasteFerry`，应用菜单项与图标装进系统目录，`/usr/bin` 下有 `pasteferry`、`pasteferry-cli`。依赖只写 `libgtk-3-0`（Ubuntu 24.04 起改名为 libgtk-3-0t64，新包也提供旧名字，两边都能装），推荐 `fonts-noto-cjk`。卸载时结束正在运行的后台服务。
- `PasteFerry-linux-<arch>.tar.gz`：解压即用，`install.sh` 把应用菜单项和图标装进用户目录（不需要 root），`uninstall.sh` 撤销并停用后台服务。

Flutter 不能交叉编译 Linux 桌面程序，CI 在 GitHub 的 x64 与 ARM64 两种 Ubuntu 22.04 机器上各跑一次：先在 Xvfb 里跑后台服务的全部测试（含 X11 剪贴板），再跑界面测试，最后打包上传。在 22.04 上构建，产物依赖的 glibc 与 GTK 版本较老，新的发行版也能运行。

## 六、验证

**单元测试与集成测试**（`linux_test.go`，CI 在 Xvfb 里用 xclip 扮演别的程序）：读写文本（含中文、换行）；本机写入后 Peek 看到新的计数、别的程序复制后本机不再持有；约 1.4 MB 的图片两个方向都走 INCR；文件的两种格式与中文文件名、网址不当作文件；HTML；密码管理器的敏感标记；分类规则、UTF-16 的 HTML、Latin-1、BMP 转 PNG、URI 解析。

**Ubuntu 24.04 GNOME（Wayland）虚拟机实测**，与 Mac 上的测试设备（fakeclip）配对：

- 装 .deb：程序、应用菜单项、图标、命令行都就位；从应用菜单项启动后 Dock 显示正确的图标。
- 点「启用后台同步」：写入自启动项，守护进程与后台服务启动，界面连上。
- 在雷达上点 Mac 的测试设备：两边指纹一致，确认后直连。
- 用 wl-copy、wl-paste（原生 Wayland 程序）测同步：文本、图片、带中文名与子目录的文件夹、单个文件，两个方向都通，SHA-256 一致；从 Mac 来的文件在剪贴板里带着 x-special/gnome-copied-files，文件管理器能直接粘贴。
- 重启虚拟机：自动登录后，自启动项拉起后台服务，带着会话的 DISPLAY 与 XAUTHORITY，1 秒内与 Mac 恢复直连，剪贴板照常同步。
- Ubuntu 自带的 avahi 与后台服务共用 UDP 5353，互不影响，局域网发现两个方向都正常。

**没测：** KDE、Xfce 等其他桌面；X11 会话（只在 Xvfb 里测过）；ARM64 真机（CI 构建并跑了测试）；Fedora 等非 Debian 系发行版上的 tar.gz。

## 七、局限

- 需要 XWayland。GNOME 46 起 XWayland 按需启动，后台服务连上 DISPLAY 时它就会启动，不需要额外设置。
- 后台服务写进剪贴板的内容由它自己持有，后台服务退出后就粘贴不出来了。这与 X11 下所有程序的行为一致；GNOME 没有剪贴板管理器替程序保管内容。
- 没有托盘图标，也没有代码签名。

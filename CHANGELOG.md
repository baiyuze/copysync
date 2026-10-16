# 更新记录 / Changelog

## 2.1.0 — 2026-10-16

Sync across networks without hosting a server. Enable optional official access,
try direct connections first and use encrypted relay only when needed.
The official website is now https://pasteferry.cc/.

### Changes

- **Official service**: enable one switch in Settings on each computer, with no
  server address to configure. Both peers try direct NAT traversal before
  separately authorizing an encrypted WSS relay tunnel to standard TURN.
- **Relay allowance**: view upload, download, remaining allowance and reset time
  in a dedicated sidebar page, with MB shown to two decimal places. The current
  default is 100 MB per machine quota group per day, upload plus download,
  resetting at midnight Beijing time. Direct traffic does not use this allowance.
- **Clear settings**: save custom servers with Save and connect or Enter;
  unsaved edits do not change the active connection. Shorter settings copy keeps
  full privacy and service-disable information in the action confirmations.
- **Large-file feedback**: receive files explicitly from History, follow their
  transfer state and receive desktop notifications when supported and enabled.
- **Privacy**: official access is off on fresh installs. Public handshakes use
  generic device labels, and real names travel through the encrypted peer
  channel. Machine quota identification happens only when relay is requested.
- **Website and guides**: publish official-service onboarding in five languages
  at pasteferry.cc, with direct-first routing, allowance and privacy information.
  Existing GitHub Pages addresses redirect to the new official website.
- **Dependencies**: update the Go toolchain and networking dependencies.

### Fixes

- Resume relay automatically after allowance is increased and retire closed
  tunnels immediately so interrupted clipboard updates enter the retry queue.
- Release authenticated signaling sockets after abrupt or half-closed TCP
  disconnects so reconnects do not accumulate server connections.
- Preserve received history when writing the clipboard fails, and report the
  failure so a completed transfer is not confused with a failed clipboard action.
- Isolate pairing checks from public STUN and use loopback relay in server tests.

### Upgrade notes

- **Mac, Windows and Linux apps**: update both peers to 2.1.0 for official relay.
  Existing identities, pairings and history are retained. Official access must
  be enabled on each computer; current limits and availability appear in the app.
- **Server**: updating is recommended for dependency and signaling fixes.
  Self-hosted servers remain independent of the private official gateway and
  require no private authorization API, callback or SDK.
- The device-to-device protocol adds optional signed negotiation and encrypted
  display metadata. Existing LAN and self-hosted transfers remain compatible.
  Older apps retain official signaling and direct support, but cannot use or
  bypass the new official relay gateway.

## 2.0.0 — 2026-10-14

CopySync is now PasteFerry (隔空贴 in Chinese). Everything carries the new name, and your pairings, history and settings come along.

### Changes

- **New name**: the app, the installers, the background service (`pasteferryd`), the command-line tool (`pasteferry-cli`), the server (`pasteferry-server`) and the data folders are now called PasteFerry, and 隔空贴 in Chinese. The project moved to github.com/baiyuze/pasteferry and the website to baiyuze.github.io/pasteferry; old links redirect.
- **Moving over from CopySync**: the first time PasteFerry opens, it stops and removes the old CopySync background service, and turns background sync on for PasteFerry if it was on before. PasteFerry's background service then moves the old data folder into place, so you don't need to pair again. The Windows installer removes the old installation, and the Linux package replaces the `copysync` package. On a Mac, drag the old CopySync app to the Trash.

### Upgrade notes

- **Mac, Windows and Linux apps**: upgrade every computer. On the same network, PasteFerry 2.0.0 and CopySync 1.x no longer find each other.
- **Server**: optional. The server protocol is unchanged, so a CopySync 1.x server keeps working. To rename it too, run the new install script: it removes the old `copysync-server` service and keeps its settings. With Docker, remove the old container first (`docker rm -f copysync-server`).
- The device-to-device protocol changed on the local network: discovery and LAN packets use new identifiers, so devices on 1.x and 2.0.0 don't see each other there. Connections through a server are unchanged.

## 1.7.0 — 2026-10-13

Search your history, use CopySync in Traditional Chinese, and report problems with a diagnostics report you can read before posting.

### Changes

- **Search the history**: a search field on the History page finds earlier copies by their text or file name, across everything you keep, not only the newest items. Matches are highlighted, and for long text the list shows the part around the match. Press ⌘F (Ctrl+F on Windows and Linux) to jump to it and Esc to clear it. It works together with the All, Text, Images and Files filter.
- **Traditional Chinese**: the app is now in Traditional Chinese for Taiwan and for Hong Kong, each in its own wording. It follows the system language (Hong Kong and Macau use the Hong Kong version), or you can pick it in Settings. The website has Traditional Chinese pages too, and a language menu in the top bar.
- **Report a problem**: Settings → Report a problem shows a diagnostics report with the app version, your system, settings, connection state and recent logs. Device names, the server address, public IP addresses, your user folder and file names are replaced before you see it, and what you copy is never logged. Copy it or save it as a file, then open a GitHub issue with the version filled in; nothing is uploaded for you. For connection problems, turn on detailed logging for 30 minutes, make the problem happen again, and regenerate the report.

### Upgrade notes

- **Mac, Windows and Linux apps**: upgrading is recommended on every computer; search, Traditional Chinese and the problem report come with the new app.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged.

## 1.6.4 — 2026-10-13

A fix for the Mac app: the background service no longer uses more and more memory the longer it runs.

### Fixes

- On a Mac, the background service (CopySync Daemon) leaked a few kilobytes every time it checked the clipboard, which it does several times a second. Its memory grew by about 1 GB a day, and Activity Monitor could show it at 1 GB or more. Copying images and files leaked their data too. Its memory now stays flat. Your clipboard history and copied files were never kept in memory: the history is stored in a database on disk and files in a cache folder, so nothing was lost.

### Upgrade notes

- **Mac, Windows and Linux apps**: upgrading the Mac app is recommended. Until you upgrade, quitting and reopening CopySync frees the memory. Nothing changed for Windows and Linux; upgrading them is optional.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged.

## 1.6.3 — 2026-10-12

The Windows app now brings the Visual C++ runtime it needs, so it also starts on computers that don't have it installed.

### Fixes

- On Windows, the app needs three Visual C++ runtime files (`msvcp140.dll`, `vcruntime140.dll`, `vcruntime140_1.dll`) that Windows doesn't always include. On a computer without the Visual C++ Redistributable, CopySync didn't open and Windows reported that `MSVCP140.dll` was not found. The installer and the portable zip now include these files next to `CopySync.exe`; you still don't need administrator rights or a separate download.

### Upgrade notes

- **Mac, Windows and Linux apps**: upgrading the Windows app is recommended, and needed if it didn't start before. Nothing changed for Mac and Linux; upgrading them is optional.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged.

## 1.6.2 — 2026-10-12

The Windows installer now runs on every version of Windows 10.

### Fixes

- On Windows 10 versions before 21H2, the installer stopped with "This program does not support the version of Windows your computer is running", although the app itself works there. It now installs on any version of Windows 10, and on Windows 11 as before.

### Upgrade notes

- **Mac, Windows and Linux apps**: the apps themselves are unchanged; upgrading is optional. On Windows 10 before 21H2, you can now use the installer instead of the portable zip.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged.

## 1.6.1 — 2026-10-12

A fix for the Mac app: discovery on the local network no longer restarts every minute or two while a paired computer is offline.

### Fixes

- On a Mac, while a paired computer on the same network was offline, the background service mistook the failed sends to it for macOS blocking local network access. It then restarted discovery on the local network every one to two minutes, which could make the Search nearby list flicker. Now it tells the two cases apart, and still shows the local network notice when macOS really blocks access.

### Upgrade notes

- **Mac, Windows and Linux apps**: upgrading the Mac app is recommended. Nothing changed for Windows and Linux; upgrading them is optional.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged.

## 1.6.0 — 2026-10-12

CopySync now runs on Linux, and the Mac app comes in separate downloads for Apple silicon and Intel.

### Changes

- **Linux client**: syncs text, rich text, images, files and folders with Mac and Windows, on the same network or through a server. It runs on x64 and ARM64 with X11, or with Wayland through XWayland (GNOME and KDE include it). Get the `.deb` for Debian, Ubuntu and derivatives, or the `tar.gz` for other distributions. The background service starts at login and keeps syncing after you close the window.
- **Mac downloads per chip**: `CopySync-mac-apple-silicon.dmg` and `CopySync-mac-intel.dmg` are about half the size of the universal `CopySync.dmg`, which is still available.
- **Download page**: the website has one download button that opens a page listing every package. It suggests the build for your system, including your Mac's chip.

### Fixes

- On macOS 13 and 14, the background service and the command line tool could not start: they were built to require macOS 15. They now run on macOS 13 and later, like the app.

### Upgrade notes

- **Mac, Windows and Linux apps**: install the Linux app on each Linux computer and pair it as usual. Mac and Windows don't need to upgrade to sync with it. Upgrade Macs on macOS 13 or 14; elsewhere upgrading is optional. On a Mac, you can switch to the build for your chip; your settings and pairings stay.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged.

## 1.5.0 — 2026-10-12

Pair two computers on the same network without a code: open Search nearby on both and click the other one. This release also fixes mixed fonts on Windows.

### Changes

- **Search nearby**: a new button next to **Add device** opens a radar showing the nearby computers that have the same window open. Click one to pair; both screens then show the usual fingerprints to compare. **Add device** and pairing codes work as before and are still the way to pair across subnets.
- **Privacy on the network**: a computer broadcasts its name only while Search nearby is open. The rest of the time, others on the network still see only an anonymous entry that changes daily.
- **Command line**: `copysync-cli nearby [NAME]` does the same from a terminal.
- **Website**: the same network guide explains pairing without a code.

### Fixes

- On Windows, Chinese text mixed glyphs from different fonts within one line. The app now uses Segoe UI with Microsoft YaHei UI for Chinese (Yu Gothic UI in the Japanese interface), and Consolas for codes, fingerprints and addresses.

### Upgrade notes

- **Mac and Windows apps**: upgrade both computers to use Search nearby. A computer on an older version doesn't appear on the radar but still pairs with a code and keeps syncing.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged. Search nearby adds two messages on the local network and two mDNS TXT fields; older versions ignore them.

## 1.4.0 — 2026-10-11

Computers on the same network now pair and sync without a server. You only need a server to connect computers on different networks.

### Changes

- **No server on the same network**: computers on one network find each other and pair with the usual six-character code. There is nothing to deploy.
- **Different subnets**: the pairing dialog shows this computer's address under the code. If the two computers are on different subnets, type that address on the other one. Addresses are remembered afterwards; if both change, use **Enter its address…** in the device's menu.
- **Tailscale without a server**: pair by typing the other computer's Tailscale (100.x) address.
- **The server is optional**: leave its address empty in Settings to sync on your network only. Devices paired on the network also sync through a server once you add one. With both, computers on the same network connect directly and keep syncing when the server is down.
- **Settings → Local network**: turn local discovery on or off. On Windows, **Allow** adds a firewall rule for CopySync with one administrator prompt. On a Mac, CopySync tells you when macOS blocks local network access and links to the setting.
- **Privacy on the network**: other people on the network see only an anonymous entry that changes daily, never the device name. Messages between devices are signed as before, and pairing still asks you to compare fingerprints.
- **Website**: a new guide for computers on the same network. The Tailscale guide no longer needs a server, and the LAN server guide stays for networks that keep devices apart, such as guest Wi-Fi.

### Fixes

- Unpairing a device now closes its connection right away.
- A device now accepts a signed message only when the message is addressed to it.

### Upgrade notes

- **Mac and Windows apps**: upgrade both computers to sync without a server; with a server, 1.4.0 still works with older versions. If you never changed the server address (the built-in `127.0.0.1` placeholder), it is cleared and the app runs without a server. On macOS 15, allow CopySync Daemon when macOS asks about finding devices on your local network; on Windows, allow CopySync when the firewall asks.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged. Computers on the same network also exchange messages over UDP port 8788 and find each other over mDNS (`_copysync._udp`); older versions ignore both.

## 1.3.4 — 2026-10-11

New step-by-step guides on the website for getting started and for choosing where to run the server. The apps and the server are unchanged apart from the version number.

### Changes

- **Quick start**: a new [quick start](https://baiyuze.github.io/copysync/en/start/) walks through installing the apps, choosing where to run the server, entering its address and pairing.
- **Three ways to run the server**, each with its own guide: a [public server](https://baiyuze.github.io/copysync/en/cloud/) reachable from any network, your [LAN](https://baiyuze.github.io/copysync/en/server/), or [LAN + Tailscale](https://baiyuze.github.io/copysync/en/tailscale/) to sync across networks without a public server.
- The website, README and server deployment guide are reorganized around these three options, in English, Chinese and Japanese.

### Upgrade notes

- **Mac and Windows apps**: no functional changes; upgrading is optional.
- **Server**: no changes; upgrading is optional.
- The device-to-device protocol is unchanged.

## 1.3.3 — 2026-10-11

- **服务器可以装在 Windows 上了**：新增 `copysync-server-windows-amd64.zip`，解压后双击 `install.cmd`，注册为开机自启的 Windows 服务（低权限的 LOCAL SERVICE 账户），并在 Windows 防火墙里只为它放行；双击 `uninstall.cmd` 卸载。
  **The server now runs on Windows.** Extract `copysync-server-windows-amd64.zip` and double-click `install.cmd`: it registers a Windows service under LOCAL SERVICE that starts with Windows, and adds a firewall rule for the server only. `uninstall.cmd` removes it.
- **Mac 上一条命令装好服务器**：`copysync-server-macos.tar.gz` 里新增 `install.sh`，注册为 LaunchAgent，登录后自动运行，不需要管理员权限。
  **One-command server install on a Mac**: the macOS package now includes `install.sh`, which registers a LaunchAgent; no administrator rights needed.
- 三个平台的安装脚本默认用本机的局域网地址启用 TURN 中转，也可以指定地址；重复运行即为升级，设置保留；都支持 `uninstall`。Linux 的 `sudo ./install.sh` 不带参数时也改为使用局域网地址（此前不启用中转）。
  All install scripts default to the machine's LAN address for the TURN relay, accept an explicit address, upgrade in place and can uninstall. On Linux, `sudo ./install.sh` without an address now uses the LAN address instead of disabling the relay.
- 服务器新增 `-log <文件>` 参数，日志写入文件，超过 10 MB 自动换新。
  New `-log <file>` server option; the log rotates at 10 MB.
- 官网新增[在局域网里部署服务器](https://baiyuze.github.io/copysync/server/)图文教程（中文、英文、日文），顶部导航可以直达。README 精简，工作原理、排障、从源码构建移到 `guide/`。
  New [LAN server guide](https://baiyuze.github.io/copysync/en/server/) on the website, linked from the top bar. The README is shorter; how it works, troubleshooting and building moved to `guide/`.
- 客户端没有变化，不需要升级；设备间协议不变。
  No client changes and no protocol changes; clients don't need to upgrade.

## 1.3.2 — 2026-10-10

- 修复与 UU 远程同时使用时，`.uuremote_aeawv…` 临时占位文件不断进入复制记录并被转发的问题。发送前过滤占位文件，混合选择保留真实文件；接收端忽略旧版发来的纯占位记录与后续传输头。
  Fixed repeated UU Remote placeholder file records and forwarding. Filter placeholders before sending; preserve ordinary files in mixed selections and ignore placeholder-only offers from older clients.
- 启动时清理历史中仅包含已确认占位文件的记录，不删除 UU 源文件，不清空正常记录。普通空文件、隐藏文件与大文件保持原行为。
  Remove placeholder-only history on startup without deleting UU source files or ordinary history. Normal empty, hidden and large files retain their behavior.
- 50 MB 自动传输阈值和 UDP P2P 设置无需改动。建议两端都升级，服务器不需要升级，设备间协议不变。
  No change to the 50 MB auto-transfer threshold, UDP P2P settings or wire protocol. Upgrade both clients; no server upgrade required.

## 1.3.1 — 2026-10-09

- **点一下记录就能预览图片，去掉了单独的「预览」按钮**（它和悬停才出现的按钮挤在一起，位置看起来是错的）。
  **Click a history item to preview it**; the separate Preview button is gone.
- **复制的图片文件也能预览**：PNG、JPEG、GIF、WebP、BMP，Mac 上还有 iPhone 照片常见的 HEIC 和 TIFF
  （由系统转换，按拍摄方向转正）。一次复制了几张图片的，可以用左右箭头或方向键翻看。
  **Copied image files can be previewed too**: PNG, JPEG, GIF, WebP, BMP, plus HEIC and TIFF on Mac. Step through
  several with the arrows or arrow keys.
- 单张图片文件在列表里显示图片图标。 Single image files show an image icon in the list.

## 1.3.0 — 2026-10-09

- **界面支持简体中文、英文、日文**：默认跟随系统语言，也可以在「设置 → 这台 Mac / 这台电脑 → 语言」里切换。
  Windows 托盘菜单也跟着变。图片和文件记录的标题按当前语言显示。
  **The interface is now in English, Simplified Chinese and Japanese.** It follows the system language and can be
  switched in Settings; the Windows tray menu follows too.
- **许可改为 GNU AGPL-3.0，另提供商业授权。**个人和公司都可以免费使用和修改；分发、或通过网络提供修改过的版本时，
  需要以同样的许可公开源码。闭源产品、不公开的修改、托管服务等情况可以购买商业授权，见 `LICENSING.md`。
  1.2.0 及更早的版本仍适用 MIT。
  **License changed to the GNU AGPL-3.0, with commercial licenses available.** Versions 1.2.0 and earlier remain MIT.
  See `LICENSING.md`.
- 复制记录：鼠标移到一行上时，「预览」按钮不再被挤开换位置，不容易点错。
  History: the Preview button no longer shifts when you hover a row.
- 图片预览按应用窗口的比例定大小（窗口的 80%），小图按原始大小显示、不放大。
  Image preview is sized relative to the app window, and small images are no longer upscaled.
- 配对对话框在文字较长（英文、日文、系统字号调大）时改为滚动，不再溢出。
  The pairing dialog scrolls instead of overflowing with longer text.
- 网站改版：整页深色，新增日文页面，首次打开按浏览器语言显示；技术方案放到网站上，不用再跳到 GitHub。
  New website design with a Japanese page and language detection; design notes are now hosted on the site.
- 服务器不需要升级；设备之间的协议没有变化。 No server upgrade needed; the device-to-device protocol is unchanged.

## 1.2.0 — 2026-10-09

- **新增 Windows 客户端**：Windows 10（21H2 起）与 Windows 11，x64。与 Mac 互通，Mac 与 Mac、
  Mac 与 Windows、Windows 与 Windows 都可以配对同步。文本、带格式文本、图片、文件、文件夹都支持，
  粘贴文件时得到的是真实文件，能直接粘进资源管理器、微信、Office。
  **New: Windows client** (Windows 10 21H2+ and 11, x64). Works with Macs and with other Windows PCs; pasted
  files are real files.
  - 安装程序按用户安装，不需要管理员权限；另有便携版 zip。开机自启、崩溃后自动重启，
    关掉窗口后在右下角托盘里继续运行，桌面上有快捷方式。
    Per-user installer with no admin rights, plus a portable zip. Starts at login, restarts after a crash, and
    stays in the notification area; desktop shortcut included.
  - Mac 上合法、Windows 上不合法的文件名（含冒号、问号等）换成对应的全角字符；只差大小写的
    文件不会互相覆盖；跳过 `.DS_Store` 等 Mac 系统文件。
    File names Windows doesn't allow get full-width characters; names differing only in case never overwrite
    each other.
  - 密码管理器标记为「不要记录」的内容不同步、不进复制记录（Mac 版尚未支持）。
    Content that password managers mark as sensitive is skipped (not yet on Mac).
  - 设备私钥用 Windows 的 DPAPI 加密存放。 The device key is encrypted with DPAPI.
- **复制记录里的图片可以点开预览**：缩放、拖动、适应窗口；还没下载的先拉取，预览用的下载
  不会改写当前剪贴板。 **Image preview** in the history, with zoom and pan; downloading for a preview doesn't
  touch your clipboard.
- 修复：配对成功后多出一个确认框，点了也没用。 Fixed: an extra, useless confirmation dialog after pairing.
- 修复：打包文件夹时其中某个文件打不开（比如正被别的程序锁住），整个传输都会失败。现在跳过那一个。
  Fixed: one unreadable file inside a folder broke the whole transfer; it's now skipped.
- 网站与 README 增加 Windows 的安装说明；技术方案见 `design/windows-client.md`，实机验证记录见
  `design/windows-validation-2026-10-09.md`。
- 服务器不需要升级。新增的 `GetImagePreview` 只是界面与本机后台服务之间的接口，设备之间的协议没有变化。
  No server upgrade needed; the device-to-device protocol is unchanged.

## 1.1.2 — 2026-10-09

- **修复：家里那台一重启就只能中转，重试也换不回直连。**谁的包先到，谁的路由器就先收到陌生来源的包；
  实测公司路由器会为它占住端口，本机随后往外发只好换端口，打洞失败。对方重启、或 ICE 重启时，总是对方
  先拿到本机的地址、先发包。现在由规则决定谁先发包：一方先公布出口地址，另一方收到后先往那边发包再公布
  自己的；初次握手由发起方等，换直连的重试里两种顺序轮流用。
  **Fixed: restarting the home Mac always fell back to the relay, and retries never got back to direct.** Some
  routers lock a port when the other side's packet arrives first. A rule now decides who sends first, and
  retries alternate the order.
- NAT 打洞实验室新增两种「路由器会为外来的包占住端口」的场景，并给信令加上延迟，共 11 种。
  The NAT lab adds two port-locking router scenarios with signaling delay (11 in total).
- 协议：`SdpOffer` 新增 `answerer_waits` 字段。服务器不需要升级。 Protocol: `SdpOffer.answerer_waits`. No server upgrade needed.

## 1.1.1 — 2026-10-09

- **修复：对方断线后，复制的内容发不过去。**连接悄悄断了（对方睡眠、断网、重启），本端却以为还连着，
  不重连，复制的内容全部发进一条死连接，日志还记成「已发送」。现在连接失效会立即关闭并重连；对方重建了
  连接，本端也随之换新；没真正送达的不再算成功。
  **Fixed: items copied after the other Mac dropped off never arrived.** A connection that died silently was
  never closed or reconnected, and everything copied afterwards went into it. Dead connections are now closed
  and rebuilt, and a transfer only counts once the data has actually left.
- **补发**：复制时对方连不上，就记下最近 20 条，连接恢复后补发。补发的只进复制记录，不写进对方的剪贴板。
  **Catch-up:** items copied while the other Mac is unreachable (latest 20) are sent when it reconnects. They go
  into its history, not its clipboard.
- **先中转，再换直连**：对端刚启动时出口地址可能晚到，第一次握手落到中转后就一直中转。现在中转时由发起方
  在空闲时用 ICE 重启再打一次洞，打通就换成直连（10 秒后第一次，之后间隔逐步拉长到 15 分钟）。
  **Relay first, direct later:** a connection that lands on the relay now retries the direct path with an ICE
  restart while idle, and switches when it works.
- 出口探测并发解析域名，整轮从约 1.6 秒降到 0.2 秒以内；中转等待从 3 秒放宽到 5 秒。
  Uplink probing resolves hostnames in parallel (about 1.6 s down to under 0.2 s); the relay wait is now 5 s.
- NAT 打洞实验室新增「一开始打不通，先中转后换直连」场景，共 9 种。 The NAT lab adds a relay-then-direct scenario (9 in total).
- 协议：`ClipOffer` 新增 `backlog` 字段。服务器不需要升级。 Protocol: `ClipOffer.backlog`. No server upgrade needed.

## 1.1.0 — 2026-10-08

- **多出口网络也能直连**：公司双线、运营商 NAT 地址池这类网络有多个出口，并按目标地址选出口。
  以前只向自己的服务器探测一次，只知道其中一个出口的地址，发往对端的包从别的出口出去就对不上，
  只能走中转。现在所有连接共用一个本地端口，在这个端口上向多台服务器探测，拿到每个出口的地址
  一并发给对端；会改端口的出口也能拿到真实映射，探测漏掉的不改端口出口按历史推算。
  **Direct connections on multi-uplink networks.** Office networks with several uplinks pick an uplink per
  destination; CopySync now probes every uplink from one shared port and sends all of them to the peer.
- 中转路径要等 3 秒才允许被选中，给直连留出打通的时间。 Relay paths wait 3 s before they can be selected.
- 设置里新增「用公共服务器探测网络出口」开关，关闭后只用自己的服务器。 New setting to probe with your own server only.
- 新增 `copysync-cli nat`：查看本机有几个出口、各设备走直连还是中转。 New `copysync-cli nat` diagnostics.
- 新增 NAT 打洞实验室（`tools/natlab`），在 CI 里用真实的网络拓扑验证 8 种场景。
  New NAT lab that checks 8 network topologies in CI.
- 服务器：TURN 可以绑定指定 IP（`relay.Options.ListenIP`）。 Server: TURN can bind to a specific IP.

## 1.0.2 — 2026-10-08

- 修复：一端经中转时，另一端误显示为「直连」，两台设备显示不一致。只有发送方能确定自己的数据
  是否走 TURN，现在连接建立后两端互相告知，任一端经中转即两端都显示「中转」。
  Fixed: when only one side was relayed, the other side showed "direct". Only the sender knows whether its
  traffic goes through TURN, so the two sides now tell each other after connecting; if either is relayed,
  both show "relay".
- 修复：收到图片或文件后，复制记录里同一条出现两行，其中一行进度条一直在动。
  Fixed: a received image or file appeared twice in the history, one row with a progress bar that never stopped.
- 修复：本机复制的图片在发送后一直显示进度条。 Fixed: items copied on this Mac kept showing a progress bar after sending.
- 传输进度按真实比例显示，并做了节流，大文件传输时不再挤掉「传输完成」的事件。
  Transfer progress now shows the actual percentage and is throttled, so large transfers no longer crowd out the completion event.

## 1.0.1 — 2026-10-08

- 修复：配对时两台 Mac 显示的安全指纹不一样，无从核对。之前每台只显示「对方」的指纹；
  现在两台显示同样的两行（双方各一行，按固定顺序），两块屏幕逐行一致即可确认。命令行的
  `pair` / `join` 同步修改。
  Fixed: the two Macs showed different fingerprints during pairing, so they could never be compared.
  Each Mac used to show only the other's fingerprint; both now show the same two lines, one per device,
  in a fixed order. The CLI's `pair` / `join` follow the same rule.
- 对话框与菜单加上投影。 Dialogs and menus now have a drop shadow.
- 发布资源改用不带版本号的文件名（如 `CopySync.dmg`），下载直链不再随版本失效。
  Release assets now use versionless names (e.g. `CopySync.dmg`) so direct download links keep working.

## 1.0.0 — 2026-10-08

第一个正式版本。 First public release.

**Mac 客户端 / Mac app**

- 同步文本、带格式的文本、图片、文件与文件夹；50 MB 以内复制即推送，更大的按需拉取。
  Syncs text, rich text, images, files and folders; items under 50 MB are pushed on copy, larger ones are pulled on demand.
- 复制记录：查看最近 3 天所有设备上复制过的内容，随时放回剪贴板。
  History of everything copied on any of your Macs in the last 3 days.
- 6 位配对码 + 安全指纹核对。 Pairing with a 6-character code and fingerprint check.
- DMG 拖拽安装；首次打开时由 App 自己启用后台同步，不需要终端。
  Drag-to-install DMG; the app enables background sync itself on first launch, no Terminal needed.
- 覆盖安装新版、或移动 App 位置后，自动校正后台服务。
  After an in-place update or moving the app, the background service is re-pointed and restarted automatically.
- 带格式的文本同时写入纯文本，粘贴到只认纯文本的输入框不再出现 HTML 标签。
  Rich text is also written as plain text, so no HTML markup leaks into plain-text fields.
- 界面在信令连接与对端上下线时实时刷新。 Connection status updates live.
- Intel 与 Apple 芯片通用；支持深色外观。 Universal binary; dark appearance.

**服务器 / Server**

- 信令 + 自带 STUN + TURN 中转，单个静态二进制。 Signaling, built-in STUN and TURN relay in one static binary.
- 提供 systemd 一键安装脚本与 Docker Compose 配置。 One-step systemd installer and Docker Compose setup.
- 新增 `-version` 参数。 Added the `-version` flag.

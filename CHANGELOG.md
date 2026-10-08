# 更新记录 / Changelog

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

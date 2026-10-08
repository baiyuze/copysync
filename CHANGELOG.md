# 更新记录 / Changelog

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

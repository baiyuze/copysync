"""界面文案的唯一来源：每条是 (键, 英文, 简体中文, 日文)。

    python3 ui/tool/l10n/strings.py

生成 ui/lib/l10n/app_en.arb、app_zh.arb、app_ja.arb，再由 flutter gen-l10n 生成 Dart 代码。
三种语言放在同一行里，改一条时不会漏掉另外两种。

占位符写成 {name}；英文里的复数用 ICU 写法，例如 {count, plural, =1{…} other{…}}。
"""

import json
import re
from pathlib import Path

S = [
    # ── 侧边栏与连接状态 ──
    ("navHistory", "History", "复制记录", "履歴"),
    ("navDevices", "Devices", "设备", "デバイス"),
    ("navSettings", "Settings", "设置", "設定"),
    ("linkConnecting", "Connecting…", "正在连接…", "接続中…"),
    ("linkServiceNotRunning", "Service not running", "后台服务未运行", "サービス停止中"),
    ("linkNoServer", "No server connection", "未连接服务器", "サーバーに未接続"),
    ("linkOnlinePeers", "{count, plural, =1{1 device online} other{{count} devices online}}", "{count} 台设备在线", "{count} 台がオンライン"),
    ("linkServerConnected", "Connected to the server", "已连接服务器", "サーバーに接続済み"),
    ("linkConnectionLost", "Connection lost: {detail}", "连接中断：{detail}", "接続が切れました：{detail}"),
    ("linkServiceDisconnected", "Service disconnected", "后台服务已断开", "サービスとの接続が切れました"),
    ("linkServiceDisabled", "Background sync is off", "后台服务已停用", "同期はオフです"),
    ("errServiceNotConnected", "Not connected to the background service", "后台服务未连接", "バックグラウンドサービスに接続していません"),
    ("errUnimplemented", "This isn't available yet", "该功能尚未接入", "この機能はまだ使えません"),
    ("errUnavailable", "Can't reach the background service", "无法连接后台服务", "バックグラウンドサービスに接続できません"),
    ("errServiceStart", "Couldn't start the background service: {detail}", "启动后台服务失败：{detail}", "バックグラウンドサービスを起動できませんでした：{detail}"),
    ("errServiceRestart", "Couldn't restart the background service: {detail}", "重启后台服务失败：{detail}", "バックグラウンドサービスを再起動できませんでした：{detail}"),
    ("errAutostart", "Couldn't set CopySync to start at login: {detail}", "登记开机自启失败：{detail}", "ログイン時の自動起動を設定できませんでした：{detail}"),

    # ── 对设备的叫法：Mac 上说「这台 Mac」，Windows 上说「这台电脑」。
    #    InText 用在句子中间，中日文里「Mac」后面紧跟汉字或假名，所以带一个空格 ──
    ("thisMacTitle", "This Mac", "这台 Mac", "この Mac"),
    ("thisPcTitle", "This computer", "这台电脑", "このコンピューター"),
    ("thisMacInText", "this Mac", "这台 Mac ", "この Mac "),
    ("thisPcInText", "this computer", "这台电脑", "このコンピューター"),
    ("bothMacs", "both Macs", "两台 Mac ", "両方の Mac "),
    ("bothDevices", "both devices", "两台设备", "両方のデバイス"),
    ("thisDeviceShort", "This device", "本机", "このデバイス"),
    ("fromDevice", "From {name}", "来自 {name}", "{name} から"),

    # ── 通用动作 ──
    ("cancel", "Cancel", "取消", "キャンセル"),
    ("close", "Close", "关闭", "閉じる"),
    ("save", "Save", "保存", "保存"),
    ("retry", "Try again", "重试", "再試行"),
    ("continueAction", "Continue", "继续", "続ける"),

    # ── 时间与时长 ──
    ("justNow", "Just now", "刚刚", "たった今"),
    ("secondsAgo", "{n}s ago", "{n} 秒前", "{n} 秒前"),
    ("minutesAgo", "{n, plural, =1{1 minute ago} other{{n} minutes ago}}", "{n} 分钟前", "{n} 分前"),
    ("hoursAgo", "{n, plural, =1{1 hour ago} other{{n} hours ago}}", "{n} 小时前", "{n} 時間前"),
    ("daysAgo", "{n, plural, =1{Yesterday} other{{n} days ago}}", "{n} 天前", "{n} 日前"),
    ("monthDay", "{month}/{day}", "{month}月{day}日", "{month}月{day}日"),
    ("durationHours", "{n, plural, =1{1 hour} other{{n} hours}}", "{n} 小时", "{n} 時間"),
    ("durationDays", "{n, plural, =1{1 day} other{{n} days}}", "{n} 天", "{n} 日"),

    # ── 复制记录 ──
    ("historyEmptySubtitle", "Copy on any of your devices and it shows up here", "在任意一台设备上复制，内容会出现在这里", "どのデバイスでコピーしても、ここに表示されます"),
    ("historyCount", "{count, plural, =1{1 item} other{{count} items}}", "共 {count} 条", "{count} 件"),
    ("historyCountKept", "{count, plural, =1{1 item} other{{count} items}}, kept for {ttl}", "共 {count} 条，保留 {ttl}", "{count} 件、{ttl}保存"),
    ("filterAll", "All", "全部", "すべて"),
    ("filterText", "Text", "文本", "テキスト"),
    ("filterImage", "Images", "图片", "画像"),
    ("filterFile", "Files", "文件", "ファイル"),
    ("historyClearTooltip", "Clear history", "清空记录", "履歴を消去"),
    ("historyEmptyTitle", "Nothing here yet", "还没有记录", "まだ履歴がありません"),
    ("historyEmptyFilteredTitle", "Nothing of this kind", "没有这一类的记录", "この種類の履歴はありません"),
    ("historyEmptyBody", "Copy text, an image or a file on any paired device and it appears here.", "在任意一台已配对的设备上复制文本、图片或文件，记录会出现在这里。", "ペアリングしたどのデバイスでテキスト、画像、ファイルをコピーしても、ここに表示されます。"),
    ("clearConfirmTitle", "Clear the whole history?", "清空全部记录？", "履歴をすべて消去しますか？"),
    ("clearConfirmBody", "Everything synced from your devices, and the cached files, will be removed from {device}. The original files aren't affected.", "所有设备上同步来的记录和已缓存的文件都会从{device}上删除。原始文件不受影响。", "同期した履歴とキャッシュしたファイルが{device}から削除されます。元のファイルには影響しません。"),
    ("clear", "Clear", "清空", "消去"),
    ("historyCleared", "History cleared", "记录已清空", "履歴を消去しました"),
    ("permDeniedTitle", "CopySync can't read the clipboard", "已禁止 CopySync 读取剪贴板", "CopySync はクリップボードを読み取れません"),
    ("permAskTitle", "Allow CopySync to read the clipboard", "允许 CopySync 读取剪贴板", "CopySync にクリップボードの読み取りを許可してください"),
    ("permDeniedBody", "What you copy on this Mac won't be synced. In System Settings, set CopySync Daemon to Allow.", "在这台 Mac 上复制的内容不会同步出去。到系统设置里把 CopySync Daemon 改为「允许」。", "この Mac でコピーした内容は同期されません。システム設定で CopySync Daemon を「許可」にしてください。"),
    ("permAskBody", "Otherwise macOS asks every time. In System Settings, set CopySync Daemon to Allow.", "否则 macOS 每次都会弹窗询问。到系统设置里把 CopySync Daemon 改为「允许」。", "許可しないと macOS が毎回確認します。システム設定で CopySync Daemon を「許可」にしてください。"),
    ("openSystemSettings", "Open System Settings", "打开系统设置", "システム設定を開く"),
    ("fetch", "Pull to this device", "拉取到本机", "このデバイスに取り込む"),
    ("fetchStarted", "Pulling…", "开始拉取", "取り込みを開始しました"),
    ("transferFailed", "Transfer failed", "传输失败", "転送に失敗しました"),
    ("expired", "Expired", "已过期", "期限切れ"),
    ("preview", "Preview", "预览", "プレビュー"),
    ("putOnClipboard", "Copy to clipboard", "放入剪贴板", "クリップボードに入れる"),
    ("putOnClipboardDone", "On the clipboard", "已放入剪贴板", "クリップボードに入れました"),
    ("deleteRecord", "Delete this item", "删除这条记录", "この履歴を削除"),
    ("deleted", "Deleted", "已删除", "削除しました"),
    ("emptyRecord", "(empty)", "（空）", "（空）"),
    ("imageRecord", "Image", "图片", "画像"),
    ("itemsSummary", "{name} and {others, plural, =1{1 more} other{{others} more}}", "{name} 等 {count} 项", "{name} ほか {others} 件"),

    # ── 图片预览 ──
    ("imagePreviewTitle", "Image preview", "图片预览", "画像プレビュー"),
    ("closePreview", "Close preview", "关闭预览", "プレビューを閉じる"),
    ("previewHint", "Scroll or pinch to zoom, drag to move around", "滚轮或触控板缩放，拖动查看", "ホイールやトラックパッドで拡大縮小、ドラッグで移動"),
    ("fitWindow", "Fit to window", "适应窗口", "ウィンドウに合わせる"),
    ("previewDeleted", "This image has been deleted from the history", "这条图片记录已被删除", "この画像は履歴から削除されました"),
    ("previewExpired", "This image has expired and can't be previewed", "图片已过期，无法预览", "この画像は期限切れのためプレビューできません"),
    ("previewDownloading", "Downloading the image…", "正在下载图片…", "画像をダウンロード中…"),
    ("previewDownloadFailed", "The download failed. Try again.", "图片下载失败，请重试", "ダウンロードに失敗しました。もう一度お試しください。"),
    ("previewNotDownloaded", "This image isn't on this computer yet", "图片尚未下载到这台电脑", "この画像はまだこのコンピューターにありません"),
    ("downloadAndPreview", "Download and preview", "下载并预览", "ダウンロードしてプレビュー"),
    ("previewUnavailable", "This image can't be previewed right now", "图片暂时无法预览", "現在この画像はプレビューできません"),
    ("imageSemantic", "Image from the history", "复制记录中的图片", "履歴の画像"),
    ("previewUnreadable", "Can't read the image. The file may have been cleaned up or damaged.", "无法读取图片，文件可能已被清理或损坏", "画像を読み込めません。ファイルが削除されたか破損している可能性があります。"),

    # ── 设备 ──
    ("devicesSubtitleEmpty", "Once paired, your devices share one clipboard", "配对之后，设备之间就会同步剪贴板", "ペアリングすると、デバイス間でクリップボードが同期されます"),
    ("devicesSubtitle", "{paired} paired, {online} online", "已配对 {paired} 台，{online} 台在线", "ペアリング済み {paired} 台、オンライン {online} 台"),
    ("addDevice", "Add device", "添加设备", "デバイスを追加"),
    ("thisDeviceFootnote", "When pairing, {both} show the same two lines of security fingerprints, one for each device. Confirm only if they match character for character.", "配对时，{both}会显示同样的两行安全指纹：这台的和对方的。逐字核对一致才能确认。", "ペアリング時、{both}に同じ 2 行のセキュリティ指紋（このデバイスと相手のもの）が表示されます。1 文字ずつ一致を確かめてから確定してください。"),
    ("noDevicesTitle", "No paired devices yet", "还没有配对的设备", "ペアリングしたデバイスはまだありません"),
    ("noDevicesBody", "Open CopySync on your other computer and link the two with a 6-character pairing code.", "在另一台电脑上打开 CopySync，用 6 位配对码把两台连起来。", "もう一台のコンピューターで CopySync を開き、6 桁のペアリングコードでつなげます。"),
    ("noDevicesNeedServer", "Connect to a signaling server in Settings before pairing.", "先在「设置」里连接信令服务器，才能配对设备。", "ペアリングの前に、「設定」でシグナリングサーバーに接続してください。"),
    ("pairedDevices", "Paired devices", "已配对的设备", "ペアリング済みのデバイス"),
    ("pairedFootnote", "Direct: data goes straight between the two devices. Relay: when the network blocks a direct path, data goes through your server and stays end-to-end encrypted.", "直连：数据在两台设备之间直接传输。中转：网络环境打不通直连时，经你的服务器转发，内容依然是端到端加密的。", "直接接続：データは 2 台のデバイス間で直接やり取りされます。中継：直接つながらないネットワークでは、あなたのサーバーを経由します。内容はエンドツーエンドで暗号化されたままです。"),
    ("offline", "Offline", "离线", "オフライン"),
    ("direct", "Direct", "直连", "直接接続"),
    ("relay", "Relay", "中转", "中継"),
    ("online", "Online", "在线", "オンライン"),
    ("moreActions", "More", "更多操作", "その他の操作"),
    ("copyFingerprint", "Copy fingerprint", "拷贝安全指纹", "指紋をコピー"),
    ("unpairEllipsis", "Unpair…", "解除配对…", "ペアリングを解除…"),
    ("fingerprintCopied", "Fingerprint copied", "安全指纹已拷贝", "指紋をコピーしました"),
    ("unpairTitle", "Unpair “{name}”?", "解除与「{name}」的配对？", "「{name}」とのペアリングを解除しますか？"),
    ("unpairBody", "The two devices will stop syncing. To sync again, you'll need to pair them and check the fingerprints again.", "两台设备将不再同步。以后想恢复，需要重新配对并核对指纹。", "2 台のデバイスは同期しなくなります。再開するには、もう一度ペアリングして指紋を確かめる必要があります。"),
    ("unpair", "Unpair", "解除配对", "ペアリングを解除"),
    ("unpaired", "Unpaired", "已解除配对", "ペアリングを解除しました"),

    # ── 配对 ──
    ("generateCode", "Show a code", "生成配对码", "コードを表示"),
    ("enterCode", "Enter a code", "输入配对码", "コードを入力"),
    ("codeFailedTitle", "Couldn't create a pairing code", "没能生成配对码", "ペアリングコードを作成できませんでした"),
    ("enterOnOther", "Enter this code on your other device", "在另一台设备上输入这个配对码", "もう一台のデバイスでこのコードを入力してください"),
    ("codeCopied", "Code copied", "配对码已拷贝", "コードをコピーしました"),
    ("codeHint", "Click to copy. Valid for 5 minutes.", "点击可拷贝，5 分钟内有效", "クリックでコピー。有効期間は 5 分です。"),
    ("afterEnterFootnote", "Once it's entered, {both} show the same two lines of security fingerprints. Confirm only if they match.", "对方输入后，{both}会显示同样的两行安全指纹，核对一致再确认。", "入力されると、{both}に同じ 2 行のセキュリティ指紋が表示されます。一致を確かめてから確定してください。"),
    ("codeInvalid", "A pairing code is 6 letters or digits", "配对码是 6 位字母或数字", "ペアリングコードは 6 桁の英数字です"),
    ("enterCodeShown", "Enter the code shown on your other device", "输入另一台设备上显示的配对码", "もう一台のデバイスに表示されたコードを入力してください"),
    ("pairingDone", "Paired", "配对完成", "ペアリングしました"),
    ("pairingRejected", "Pairing declined", "已拒绝配对", "ペアリングを拒否しました"),
    ("verifyTitle", "Check the security fingerprints", "核对安全指纹", "セキュリティ指紋を確認"),
    ("pairingWith", "Pairing with “{name}”", "正在与「{name}」配对", "「{name}」とペアリング中"),
    ("verifyBody", "The two lines on {both} must be identical. If even one character differs, someone may be tampering with the connection, so decline.", "{both}上显示的这两行应当完全相同。有任何一个字符不同，说明连接可能被第三方篡改，请拒绝。", "{both}に表示された 2 行は完全に同じはずです。1 文字でも違う場合は、第三者が通信を改ざんしている可能性があるため、拒否してください。"),
    ("rejectMismatch", "Mismatch, decline", "不一致，拒绝", "不一致、拒否"),
    ("confirmMatch", "They match, pair", "一致，确认配对", "一致、ペアリング"),

    # ── 设置 ──
    ("server", "Server", "服务器", "サーバー"),
    ("serverFootnote", "The server only helps devices find each other and set up a direct connection. It can't see what they send.", "服务器只负责让设备找到彼此、协商直连。设备之间传输的内容它看不到。", "サーバーはデバイス同士を見つけ、直接接続を取り持つだけです。やり取りする内容は見えません。"),
    ("signalingUrl", "Signaling server", "信令服务器地址", "シグナリングサーバー"),
    ("signalingUrlDesc", "Press Return to save; it reconnects right away", "按回车保存，保存后立即重新连接", "Enter で保存し、すぐに再接続します"),
    ("signalingUrlPlaceholder", "ws://server-address:8787/signal", "ws://服务器地址:8787/signal", "ws://サーバーアドレス:8787/signal"),
    ("publicStun", "Probe with public servers", "用公共服务器探测网络出口", "公開サーバーで出口を調べる"),
    ("publicStunDesc", "On networks with several uplinks, such as offices with two ISPs, this finds more direct paths. Public servers only see your public IP.", "公司双线这类多出口网络里，能找到更多可以直连的路径。公共服务器只会看到你的公网 IP", "複数の回線を持つ会社などのネットワークで、直接つながる経路を多く見つけます。公開サーバーに見えるのはグローバル IP だけです。"),
    ("connectionStatus", "Connection", "连接状态", "接続状態"),
    ("connectedDesc", "Connected to the signaling server", "已连接到信令服务器", "シグナリングサーバーに接続しています"),
    ("disconnectedDesc", "Can't reach the server. Check the address and that the server is running.", "连不上服务器。检查地址是否正确、服务器是否在运行。", "サーバーに接続できません。アドレスとサーバーの稼働を確認してください。"),
    ("connected", "Connected", "已连接", "接続済み"),
    ("notConnected", "Not connected", "未连接", "未接続"),
    ("deviceName", "Device name", "设备名称", "デバイス名"),
    ("deviceNameDesc", "Shown in the list on your other devices", "显示在其他设备的列表里", "他のデバイスの一覧に表示されます"),
    ("language", "Language", "语言", "言語"),
    ("languageSystem", "System", "跟随系统", "システムに合わせる"),
    ("sync", "Sync", "同步", "同期"),
    ("syncFootnote", "Items under the limit are pushed to your other devices as you copy. Larger ones sync as a record; pull them from History when you need them, so big files don't waste bandwidth or disk space.", "小于阈值的内容在复制时直接推送到其他设备；超过阈值的只同步一条记录，需要时在「复制记录」里点「拉取到本机」，避免大文件无谓地占用带宽和磁盘。", "上限より小さい内容はコピーした時点で他のデバイスに送られます。大きいものは記録だけを同期し、必要なときに「履歴」から取り込むので、帯域やディスクを無駄に使いません。"),
    ("autoSyncLimit", "Auto-sync limit", "自动同步上限", "自動同期の上限"),
    ("autoSyncLimitDesc", "Larger files are pulled by hand", "超过这个大小的文件改为手动拉取", "これより大きいファイルは手動で取り込みます"),
    ("autoApply", "Put received items on the clipboard", "收到后直接写入剪贴板", "受け取った内容をクリップボードに入れる"),
    ("autoApplyDesc", "When off, put them on the clipboard from History yourself", "关闭后，需要在「复制记录」里手动放入剪贴板", "オフの場合は「履歴」から手動でクリップボードに入れます"),
    ("whatToSync", "What to sync", "同步哪些内容", "同期する内容"),
    ("syncFiles", "Files and folders", "文件与文件夹", "ファイルとフォルダ"),
    ("syncText", "Plain text", "纯文本", "プレーンテキスト"),
    ("syncRich", "Rich text", "带格式的文本", "書式付きテキスト"),
    ("syncImages", "Images and screenshots", "图片与截图", "画像とスクリーンショット"),
    ("storage", "Storage", "存储", "ストレージ"),
    ("storageFootnote", "Expired items and cached files are cleaned up automatically.", "到期的记录与缓存文件会被自动清理。", "期限切れの履歴とキャッシュは自動で削除されます。"),
    ("keepHistory", "Keep history for", "记录保留", "履歴の保存期間"),
    ("keepCache", "Keep cached files for", "缓存文件保留", "キャッシュの保存期間"),
    ("keepCacheDesc", "Never longer than the history", "不会超过记录的保留时长", "履歴の保存期間を超えません"),
    ("cacheUsage", "Cache size", "缓存占用", "キャッシュ容量"),
    ("backgroundSync", "Background sync", "后台同步", "バックグラウンド同期"),
    ("backgroundFootnote", "When off, nothing syncs and CopySync doesn't start at login. History and pairings are kept; open the app to turn it back on.", "停用后不再同步，也不会开机启动。历史记录与配对关系会保留，重新打开 App 即可再次启用。", "オフにすると同期せず、ログイン時にも起動しません。履歴とペアリングは残り、アプリを開けば再び有効にできます。"),
    ("backgroundService", "Background service", "后台服务", "バックグラウンドサービス"),
    ("backgroundServiceDesc", "Starts at login and keeps syncing with the window closed", "随系统登录自动启动，关掉窗口也照常同步", "ログイン時に起動し、ウィンドウを閉じても同期を続けます"),
    ("restart", "Restart", "重新启动", "再起動"),
    ("restarted", "Background service restarted", "后台服务已重新启动", "バックグラウンドサービスを再起動しました"),
    ("disableEllipsis", "Turn off…", "停用…", "オフにする…"),
    ("about", "About", "关于", "このアプリについて"),
    ("aboutDesc", "Open source under the AGPL-3.0", "开源软件，AGPL-3.0 许可", "オープンソース（AGPL-3.0）"),
    ("homepage", "Website", "项目主页", "ウェブサイト"),
    ("feedback", "Report a problem", "反馈问题", "問題を報告"),
    ("disableTitle", "Turn off background sync?", "停用后台同步？", "バックグラウンド同期をオフにしますか？"),
    ("disableBody", "CopySync on {device} will stop syncing and won't start at login. History and pairings are kept.", "{device}将停止同步，也不再开机启动。历史记录和配对关系会保留。", "{device}は同期を停止し、ログイン時にも起動しなくなります。履歴とペアリングは残ります。"),
    ("disable", "Turn off", "停用", "オフにする"),
    ("disabled", "Background sync is off", "后台同步已停用", "バックグラウンド同期をオフにしました"),
    ("urlMustBeWs", "The address must start with ws:// or wss://", "地址需以 ws:// 或 wss:// 开头", "アドレスは ws:// または wss:// で始めてください"),
    ("serverUpdated", "Server updated. Reconnecting…", "服务器地址已更新，正在重连", "サーバーを更新しました。再接続しています…"),
    ("deviceNameUpdated", "Device name updated", "设备名称已更新", "デバイス名を更新しました"),

    # ── 启用后台服务的引导页 ──
    ("svcExtractTitle", "Install or extract CopySync first", "先解压或安装 CopySync", "先に CopySync をインストールまたは展開してください"),
    ("svcExtractBody", "CopySync is running straight from the zip, and its temporary folder disappears when you close it. Run the installer, or extract the whole folder and open it from there, to turn on background sync.", "CopySync 现在是直接从压缩包里打开的，关掉之后它所在的临时文件夹就会被删除。运行安装程序，或者把整个文件夹解压出来再打开，就可以启用后台同步。", "CopySync は zip から直接開かれていて、閉じると一時フォルダごと消えます。インストーラーを実行するか、フォルダ全体を展開してから開くと、バックグラウンド同期を有効にできます。"),
    ("svcMoveTitle", "Move CopySync to Applications first", "先把 CopySync 移到「应用程序」", "先に CopySync を「アプリケーション」に移動してください"),
    ("svcMoveBody", "CopySync is running straight from the disk image. Drag it into Applications, eject the disk image, and open it from Applications to turn on background sync.", "CopySync 现在是直接从安装盘里打开的。把它拖进「应用程序」文件夹，推出安装盘，再从「应用程序」里打开，就可以启用后台同步。", "CopySync はディスクイメージから直接開かれています。「アプリケーション」にドラッグしてディスクイメージを取り出し、「アプリケーション」から開くと、バックグラウンド同期を有効にできます。"),
    ("svcStartTitle", "Get started with CopySync", "开始使用 CopySync", "CopySync を始める"),
    ("svcStartBody", "CopySync runs a small background process that notices clipboard changes and stays connected to your other devices. Once on, it starts at login and keeps syncing even with this window closed.", "CopySync 会在后台运行一个小进程，用来感知剪贴板的变化，并与你的其他设备保持连接。启用后它随系统登录自动启动，关掉这个窗口也会照常同步。", "CopySync は小さなバックグラウンドプロセスを動かし、クリップボードの変化を検知して他のデバイスとつながり続けます。有効にするとログイン時に起動し、このウィンドウを閉じても同期を続けます。"),
    ("svcEnable", "Turn on background sync", "启用后台同步", "バックグラウンド同期を有効にする"),
    ("svcStaleTitle", "Turn background sync on again", "重新启用后台同步", "バックグラウンド同期をもう一度有効にする"),
    ("svcStaleBodyWin", "CopySync has moved, but starting at login still points to the old place. Turn it on again so the background process uses this copy. History and pairings are kept.", "CopySync 换了位置，开机自启还指向原来的地方。重新启用一次，让后台进程指向当前这份 CopySync。历史记录与配对关系都会保留。", "CopySync の場所が変わりましたが、自動起動は以前の場所を指しています。もう一度有効にすると、このコピーを使うようになります。履歴とペアリングは残ります。"),
    ("svcStaleBodyMac", "CopySync has moved, or was installed with an old script. Turn it on again so the background process uses this app. History and pairings are kept.", "CopySync 换了位置，或者之前用旧版安装脚本装过。重新启用一次，让后台进程指向当前这个 App。历史记录与配对关系都会保留。", "CopySync の場所が変わったか、古いスクリプトでインストールされています。もう一度有効にすると、このアプリを使うようになります。履歴とペアリングは残ります。"),
    ("svcReenable", "Turn on again", "重新启用", "もう一度有効にする"),
    ("svcStoppedTitle", "The background service isn't running", "后台服务没有运行", "バックグラウンドサービスが動いていません"),
    ("svcStoppedBody", "It may have just quit; the system usually restarts it within seconds. If you stay on this page, restart it with the button below.", "它可能刚刚退出，系统通常会在几秒内把它重新拉起。如果一直停在这一页，点下面的按钮重新启动。", "終了した直後かもしれません。通常は数秒で再起動されます。この画面のままなら、下のボタンで再起動してください。"),
    ("svcDevBodyWin", "This is a development build with no copysyncd.exe next to it. Run scripts\\build-windows.ps1 in the source folder to build a complete version.", "这是开发构建，旁边没有 copysyncd.exe。在源码目录运行 scripts\\build-windows.ps1 构建完整的版本。", "開発ビルドのため copysyncd.exe がありません。ソースフォルダで scripts\\build-windows.ps1 を実行して完全版をビルドしてください。"),
    ("svcDevBodyMac", "This is a development build without the background service. Run ./scripts/install-macos.sh in the source folder to install it.", "这是开发构建，App 里没有内置后台服务。在源码目录运行 ./scripts/install-macos.sh 安装它。", "開発ビルドのためバックグラウンドサービスが含まれていません。ソースフォルダで ./scripts/install-macos.sh を実行してインストールしてください。"),
    ("svcLogPath", "Log: {path}", "运行日志：{path}", "ログ：{path}"),
    ("svcCanDisable", "You can turn background sync off anytime in Settings.", "随时可以在「设置」中停用后台同步。", "バックグラウンド同期は「設定」でいつでもオフにできます。"),
]

# 数字类占位符的类型；其余都是字符串
INTS = {"count", "n", "paired", "online", "others", "month", "day"}


def placeholders(text: str) -> list[str]:
    # 只取最外层的 {name} 与 {name, plural, …}；plural 分支里的 =1{Yesterday} 不是占位符
    names, depth = [], 0
    for i, ch in enumerate(text):
        if ch == "{":
            if depth == 0:
                m = re.match(r"\{(\w+)\s*[,}]", text[i:])
                if m and m.group(1) not in names:
                    names.append(m.group(1))
            depth += 1
        elif ch == "}":
            depth -= 1
    return names


def main() -> None:
    out = Path(__file__).resolve().parents[2] / "lib" / "l10n"
    out.mkdir(parents=True, exist_ok=True)
    keys = [k for k, *_ in S]
    dupes = {k for k in keys if keys.count(k) > 1}
    assert not dupes, f"重复的键：{dupes}"

    en = {"@@locale": "en"}
    for key, text, zh, ja in S:
        names = placeholders(text) or placeholders(zh)
        for other in (zh, ja):
            names += [n for n in placeholders(other) if n not in names]
        en[key] = text
        if names:
            en["@" + key] = {
                "placeholders": {n: {"type": "int" if n in INTS else "String"} for n in names}
            }
    for locale, index in (("zh", 2), ("ja", 3)):
        data = {"@@locale": locale}
        for row in S:
            data[row[0]] = row[index]
        (out / f"app_{locale}.arb").write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    (out / "app_en.arb").write_text(json.dumps(en, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"✓ {len(S)} 条，写入 {out}")


if __name__ == "__main__":
    main()

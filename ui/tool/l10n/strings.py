"""界面文案的唯一来源：每条是 (键, 英文, 简体中文, 日文, 繁体中文（台湾）, 繁体中文（香港）)。

    python3 ui/tool/l10n/strings.py

生成 ui/lib/l10n/app_en.arb、app_zh.arb、app_ja.arb、app_zh_Hant_TW.arb、app_zh_Hant_HK.arb，
再由 flutter gen-l10n 生成 Dart 代码。各种语言放在同一行里，改一条时不会漏掉别的。
两种繁体用词不同：台湾写位址、網路、軟體，香港写地址、網絡、軟件，按当地的写法各写一份。

占位符写成 {name}；英文里的复数用 ICU 写法，例如 {count, plural, =1{…} other{…}}。
"""

import json
import re
from pathlib import Path

S = [
    # ── 侧边栏与连接状态 ──
    ("navHistory", "History", "复制记录", "履歴", "複製記錄", "複製記錄"),
    ("navDevices", "Devices", "设备", "デバイス", "裝置", "裝置"),
    ("navSettings", "Settings", "设置", "設定", "設定", "設定"),
    ("linkConnecting", "Connecting…", "正在连接…", "接続中…", "正在連線…", "正在連線…"),
    ("linkServiceNotRunning", "Service not running", "后台服务未运行", "サービス停止中", "背景服務未執行", "背景服務未在執行"),
    ("linkLanMode", "Looking for devices on this network", "在局域网中查找设备", "このネットワークでデバイスを探しています", "在區域網路中尋找裝置", "在區域網絡中尋找裝置"),
    ("linkServerUnreachable", "Can't reach the server; syncing on this network", "连不上服务器，局域网内照常同步", "サーバーに接続できません。このネットワーク内では同期します", "連不上伺服器，區域網路內照常同步", "連不上伺服器，區域網絡內照常同步"),
    ("linkNothingOn", "Local network discovery is off and no server is set", "局域网发现已关闭，也没有设置服务器", "ローカルネットワークの検出がオフで、サーバーも未設定です", "區域網路探索已關閉，也沒有設定伺服器", "區域網絡發現已關閉，也沒有設定伺服器"),
    ("linkLanBlocked", "macOS is blocking the local network", "macOS 阻止了局域网访问", "macOS がローカルネットワークをブロックしています", "macOS 阻擋了區域網路存取", "macOS 阻擋了區域網絡存取"),
    ("linkOnlinePeers", "{count, plural, =1{1 device online} other{{count} devices online}}", "{count} 台设备在线", "{count} 台がオンライン", "{count} 台裝置在線上", "{count} 部裝置在線"),
    ("linkServerConnected", "Connected to the server", "已连接服务器", "サーバーに接続済み", "已連線到伺服器", "已連線到伺服器"),
    ("linkConnectionLost", "Connection lost: {detail}", "连接中断：{detail}", "接続が切れました：{detail}", "連線中斷：{detail}", "連線中斷：{detail}"),
    ("linkServiceDisconnected", "Service disconnected", "后台服务已断开", "サービスとの接続が切れました", "背景服務已中斷連線", "背景服務已中斷連線"),
    ("linkServiceDisabled", "Background sync is off", "后台服务已停用", "同期はオフです", "背景服務已停用", "背景服務已停用"),
    ("errServiceNotConnected", "Not connected to the background service", "后台服务未连接", "バックグラウンドサービスに接続していません", "背景服務未連線", "未連線到背景服務"),
    ("errUnimplemented", "This isn't available yet", "该功能尚未接入", "この機能はまだ使えません", "這項功能尚未提供", "這項功能仍未提供"),
    ("errUnavailable", "Can't reach the background service", "无法连接后台服务", "バックグラウンドサービスに接続できません", "無法連線到背景服務", "無法連線到背景服務"),
    ("errServiceStart", "Couldn't start the background service: {detail}", "启动后台服务失败：{detail}", "バックグラウンドサービスを起動できませんでした：{detail}", "啟動背景服務失敗：{detail}", "無法啟動背景服務：{detail}"),
    ("errServiceRestart", "Couldn't restart the background service: {detail}", "重启后台服务失败：{detail}", "バックグラウンドサービスを再起動できませんでした：{detail}", "重新啟動背景服務失敗：{detail}", "無法重新啟動背景服務：{detail}"),
    ("errAutostart", "Couldn't set CopySync to start at login: {detail}", "登记开机自启失败：{detail}", "ログイン時の自動起動を設定できませんでした：{detail}", "設定開機自動啟動失敗：{detail}", "無法設定開機自動啟動：{detail}"),

    # ── 对设备的叫法：Mac 上说「这台 Mac」，Windows 上说「这台电脑」。
    #    InText 用在句子中间，中日文里「Mac」后面紧跟汉字或假名，所以带一个空格 ──
    ("thisMacTitle", "This Mac", "这台 Mac", "この Mac", "這台 Mac", "這部 Mac"),
    ("thisPcTitle", "This computer", "这台电脑", "このコンピューター", "這台電腦", "這部電腦"),
    ("thisMacInText", "this Mac", "这台 Mac ", "この Mac ", "這台 Mac ", "這部 Mac "),
    ("thisPcInText", "this computer", "这台电脑", "このコンピューター", "這台電腦", "這部電腦"),
    ("bothMacs", "both Macs", "两台 Mac ", "両方の Mac ", "兩台 Mac ", "兩部 Mac "),
    ("bothDevices", "both devices", "两台设备", "両方のデバイス", "兩台裝置", "兩部裝置"),
    ("thisDeviceShort", "This device", "本机", "このデバイス", "本機", "本機"),
    ("fromDevice", "From {name}", "来自 {name}", "{name} から", "來自 {name}", "來自 {name}"),

    # ── 通用动作 ──
    ("cancel", "Cancel", "取消", "キャンセル", "取消", "取消"),
    ("close", "Close", "关闭", "閉じる", "關閉", "關閉"),
    ("save", "Save", "保存", "保存", "儲存", "儲存"),
    ("retry", "Try again", "重试", "再試行", "重試", "重試"),
    ("continueAction", "Continue", "继续", "続ける", "繼續", "繼續"),

    # ── 时间与时长 ──
    ("justNow", "Just now", "刚刚", "たった今", "剛剛", "剛剛"),
    ("secondsAgo", "{n}s ago", "{n} 秒前", "{n} 秒前", "{n} 秒前", "{n} 秒前"),
    ("minutesAgo", "{n, plural, =1{1 minute ago} other{{n} minutes ago}}", "{n} 分钟前", "{n} 分前", "{n} 分鐘前", "{n} 分鐘前"),
    ("hoursAgo", "{n, plural, =1{1 hour ago} other{{n} hours ago}}", "{n} 小时前", "{n} 時間前", "{n} 小時前", "{n} 小時前"),
    ("daysAgo", "{n, plural, =1{Yesterday} other{{n} days ago}}", "{n} 天前", "{n} 日前", "{n} 天前", "{n} 天前"),
    ("monthDay", "{month}/{day}", "{month}月{day}日", "{month}月{day}日", "{month}月{day}日", "{month}月{day}日"),
    ("durationHours", "{n, plural, =1{1 hour} other{{n} hours}}", "{n} 小时", "{n} 時間", "{n} 小時", "{n} 小時"),
    ("durationDays", "{n, plural, =1{1 day} other{{n} days}}", "{n} 天", "{n} 日", "{n} 天", "{n} 天"),

    # ── 复制记录 ──
    ("historyEmptySubtitle", "Copy on any of your devices and it shows up here", "在任意一台设备上复制，内容会出现在这里", "どのデバイスでコピーしても、ここに表示されます", "在任何一台裝置上複製，內容會出現在這裡", "在任何一部裝置上複製，內容會出現在這裏"),
    ("historyCount", "{count, plural, =1{1 item} other{{count} items}}", "共 {count} 条", "{count} 件", "共 {count} 筆", "共 {count} 條"),
    ("historyCountKept", "{count, plural, =1{1 item} other{{count} items}}, kept for {ttl}", "共 {count} 条，保留 {ttl}", "{count} 件、{ttl}保存", "共 {count} 筆，保留 {ttl}", "共 {count} 條，保留 {ttl}"),
    ("filterAll", "All", "全部", "すべて", "全部", "全部"),
    ("filterText", "Text", "文本", "テキスト", "文字", "文字"),
    ("filterImage", "Images", "图片", "画像", "圖片", "圖片"),
    ("filterFile", "Files", "文件", "ファイル", "檔案", "檔案"),
    ("historyClearTooltip", "Clear history", "清空记录", "履歴を消去", "清空記錄", "清除記錄"),
    ("historyEmptyTitle", "Nothing here yet", "还没有记录", "まだ履歴がありません", "還沒有記錄", "仍未有記錄"),
    ("historyEmptyFilteredTitle", "Nothing of this kind", "没有这一类的记录", "この種類の履歴はありません", "沒有這一類的記錄", "沒有這一類的記錄"),
    ("historyEmptyBody", "Copy text, an image or a file on any paired device and it appears here.", "在任意一台已配对的设备上复制文本、图片或文件，记录会出现在这里。", "ペアリングしたどのデバイスでテキスト、画像、ファイルをコピーしても、ここに表示されます。", "在任何一台已配對的裝置上複製文字、圖片或檔案，記錄會出現在這裡。", "在任何一部已配對的裝置上複製文字、圖片或檔案，記錄會出現在這裏。"),
    ("clearConfirmTitle", "Clear the whole history?", "清空全部记录？", "履歴をすべて消去しますか？", "要清空全部記錄嗎？", "清除全部記錄？"),
    ("clearConfirmBody", "Everything synced from your devices, and the cached files, will be removed from {device}. The original files aren't affected.", "所有设备上同步来的记录和已缓存的文件都会从{device}上删除。原始文件不受影响。", "同期した履歴とキャッシュしたファイルが{device}から削除されます。元のファイルには影響しません。", "從所有裝置同步過來的記錄和已快取的檔案，都會從{device}上刪除。原始檔案不受影響。", "所有裝置同步過來的記錄和已快取的檔案都會從{device}上刪除。原始檔案不受影響。"),
    ("clear", "Clear", "清空", "消去", "清空", "清除"),
    ("historyCleared", "History cleared", "记录已清空", "履歴を消去しました", "記錄已清空", "記錄已清除"),
    ("permDeniedTitle", "CopySync can't read the clipboard", "已禁止 CopySync 读取剪贴板", "CopySync はクリップボードを読み取れません", "已禁止 CopySync 讀取剪貼簿", "已禁止 CopySync 讀取剪貼簿"),
    ("permAskTitle", "Allow CopySync to read the clipboard", "允许 CopySync 读取剪贴板", "CopySync にクリップボードの読み取りを許可してください", "允許 CopySync 讀取剪貼簿", "允許 CopySync 讀取剪貼簿"),
    ("permDeniedBody", "What you copy on this Mac won't be synced. In System Settings, set CopySync Daemon to Allow.", "在这台 Mac 上复制的内容不会同步出去。到系统设置里把 CopySync Daemon 改为「允许」。", "この Mac でコピーした内容は同期されません。システム設定で CopySync Daemon を「許可」にしてください。", "在這台 Mac 上複製的內容不會同步出去。請到系統設定裡把 CopySync Daemon 改為「允許」。", "在這部 Mac 上複製的內容不會同步出去。請到系統設定裏把 CopySync Daemon 改為「允許」。"),
    ("permAskBody", "Otherwise macOS asks every time. In System Settings, set CopySync Daemon to Allow.", "否则 macOS 每次都会弹窗询问。到系统设置里把 CopySync Daemon 改为「允许」。", "許可しないと macOS が毎回確認します。システム設定で CopySync Daemon を「許可」にしてください。", "否則 macOS 每次都會跳出視窗詢問。請到系統設定裡把 CopySync Daemon 改為「允許」。", "否則 macOS 每次都會彈出視窗詢問。請到系統設定裏把 CopySync Daemon 改為「允許」。"),
    ("openSystemSettings", "Open System Settings", "打开系统设置", "システム設定を開く", "打開系統設定", "開啟系統設定"),
    ("fetch", "Pull to this device", "拉取到本机", "このデバイスに取り込む", "拉取到本機", "拉取到本機"),
    ("fetchStarted", "Pulling…", "开始拉取", "取り込みを開始しました", "開始拉取", "開始拉取"),
    ("transferFailed", "Transfer failed", "传输失败", "転送に失敗しました", "傳輸失敗", "傳輸失敗"),
    ("expired", "Expired", "已过期", "期限切れ", "已過期", "已過期"),
    ("putOnClipboard", "Copy to clipboard", "放入剪贴板", "クリップボードに入れる", "放入剪貼簿", "放入剪貼簿"),
    ("putOnClipboardDone", "On the clipboard", "已放入剪贴板", "クリップボードに入れました", "已放入剪貼簿", "已放入剪貼簿"),
    ("deleteRecord", "Delete this item", "删除这条记录", "この履歴を削除", "刪除這筆記錄", "刪除這條記錄"),
    ("deleted", "Deleted", "已删除", "削除しました", "已刪除", "已刪除"),
    ("emptyRecord", "(empty)", "（空）", "（空）", "（空）", "（空）"),
    ("imageRecord", "Image", "图片", "画像", "圖片", "圖片"),
    ("itemsSummary", "{name} and {others, plural, =1{1 more} other{{others} more}}", "{name} 等 {count} 项", "{name} ほか {others} 件", "{name} 等 {count} 項", "{name} 等 {count} 項"),

    # ── 图片预览 ──
    ("imagePreviewTitle", "Image preview", "图片预览", "画像プレビュー", "圖片預覽", "圖片預覽"),
    ("closePreview", "Close preview", "关闭预览", "プレビューを閉じる", "關閉預覽", "關閉預覽"),
    ("previousImage", "Previous image", "上一张", "前の画像", "上一張", "上一張"),
    ("nextImage", "Next image", "下一张", "次の画像", "下一張", "下一張"),
    ("previewHint", "Scroll or pinch to zoom, drag to move around", "滚轮或触控板缩放，拖动查看", "ホイールやトラックパッドで拡大縮小、ドラッグで移動", "用滾輪或觸控板縮放，拖曳查看", "用滾輪或觸控板縮放，拖移查看"),
    ("fitWindow", "Fit to window", "适应窗口", "ウィンドウに合わせる", "符合視窗大小", "符合視窗大小"),
    ("previewDeleted", "This image has been deleted from the history", "这条图片记录已被删除", "この画像は履歴から削除されました", "這筆圖片記錄已被刪除", "這條圖片記錄已被刪除"),
    ("previewExpired", "This image has expired and can't be previewed", "图片已过期，无法预览", "この画像は期限切れのためプレビューできません", "圖片已過期，無法預覽", "圖片已過期，無法預覽"),
    ("previewDownloading", "Downloading the image…", "正在下载图片…", "画像をダウンロード中…", "正在下載圖片…", "正在下載圖片…"),
    ("previewDownloadFailed", "The download failed. Try again.", "图片下载失败，请重试", "ダウンロードに失敗しました。もう一度お試しください。", "圖片下載失敗，請重試", "圖片下載失敗，請重試"),
    ("previewNotDownloaded", "This image isn't on this computer yet", "图片尚未下载到这台电脑", "この画像はまだこのコンピューターにありません", "圖片尚未下載到這台電腦", "圖片仍未下載到這部電腦"),
    ("downloadAndPreview", "Download and preview", "下载并预览", "ダウンロードしてプレビュー", "下載並預覽", "下載並預覽"),
    ("previewUnavailable", "This image can't be previewed right now", "图片暂时无法预览", "現在この画像はプレビューできません", "圖片暫時無法預覽", "圖片暫時無法預覽"),
    ("imageSemantic", "Image from the history", "复制记录中的图片", "履歴の画像", "複製記錄中的圖片", "複製記錄中的圖片"),
    ("previewUnreadable", "Can't read the image. The file may have been cleaned up or damaged.", "无法读取图片，文件可能已被清理或损坏", "画像を読み込めません。ファイルが削除されたか破損している可能性があります。", "無法讀取圖片，檔案可能已被清除或損毀", "無法讀取圖片，檔案可能已被清理或損壞"),

    # ── 设备 ──
    ("devicesSubtitleEmpty", "Once paired, your devices share one clipboard", "配对之后，设备之间就会同步剪贴板", "ペアリングすると、デバイス間でクリップボードが同期されます", "配對之後，裝置之間就會同步剪貼簿", "配對之後，裝置之間就會同步剪貼簿"),
    ("devicesSubtitle", "{paired} paired, {online} online", "已配对 {paired} 台，{online} 台在线", "ペアリング済み {paired} 台、オンライン {online} 台", "已配對 {paired} 台，{online} 台在線上", "已配對 {paired} 部，{online} 部在線"),
    ("addDevice", "Add device", "添加设备", "デバイスを追加", "新增裝置", "新增裝置"),
    ("thisDeviceFootnote", "When pairing, {both} show the same two lines of security fingerprints, one for each device. Confirm only if they match character for character.", "配对时，{both}会显示同样的两行安全指纹：这台的和对方的。逐字核对一致才能确认。", "ペアリング時、{both}に同じ 2 行のセキュリティ指紋（このデバイスと相手のもの）が表示されます。1 文字ずつ一致を確かめてから確定してください。", "配對時，{both}會顯示同樣的兩行安全指紋：這台的和對方的。逐字核對一致才能確認。", "配對時，{both}會顯示同樣的兩行安全指紋：這部的和對方的。逐字核對一致才可確認。"),
    ("noDevicesTitle", "No paired devices yet", "还没有配对的设备", "ペアリングしたデバイスはまだありません", "還沒有配對的裝置", "仍未有已配對的裝置"),
    ("noDevicesBody", "Open CopySync on your other computer and link the two with a 6-character pairing code. On the same network, you can also open Search nearby on both and click the other one.", "在另一台电脑上打开 CopySync，用 6 位配对码把两台连起来。两台在同一个网络里时，也可以在两台上都打开「搜索附近设备」，点一下对方。", "もう一台のコンピューターで CopySync を開き、6 桁のペアリングコードでつなげます。同じネットワーク内なら、両方で「近くを検索」を開いて相手をクリックすることもできます。", "在另一台電腦上打開 CopySync，用 6 碼配對碼把兩台連起來。兩台在同一個網路裡時，也可以在兩台上都打開「搜尋附近裝置」，點一下對方。", "在另一部電腦上開啟 CopySync，用 6 位配對碼把兩部連起來。兩部在同一個網絡裏時，也可以在兩部上都開啟「搜尋附近裝置」，按一下對方。"),
    ("pairedDevices", "Paired devices", "已配对的设备", "ペアリング済みのデバイス", "已配對的裝置", "已配對的裝置"),
    ("pairedFootnote", "Direct: data goes straight between the two devices. Relay: when the network blocks a direct path, data goes through your server and stays end-to-end encrypted.", "直连：数据在两台设备之间直接传输。中转：网络环境打不通直连时，经你的服务器转发，内容依然是端到端加密的。", "直接接続：データは 2 台のデバイス間で直接やり取りされます。中継：直接つながらないネットワークでは、あなたのサーバーを経由します。内容はエンドツーエンドで暗号化されたままです。", "直連：資料在兩台裝置之間直接傳輸。中繼：網路環境打不通直連時，經你的伺服器轉送，內容依然是端對端加密的。", "直連：資料在兩部裝置之間直接傳輸。中轉：網絡環境打不通直連時，經你的伺服器轉發，內容依然是端對端加密的。"),
    ("offline", "Offline", "离线", "オフライン", "離線", "離線"),
    ("direct", "Direct", "直连", "直接接続", "直連", "直連"),
    ("relay", "Relay", "中转", "中継", "中繼", "中轉"),
    ("online", "Online", "在线", "オンライン", "線上", "在線"),
    ("moreActions", "More", "更多操作", "その他の操作", "更多操作", "更多操作"),
    ("copyFingerprint", "Copy fingerprint", "拷贝安全指纹", "指紋をコピー", "複製安全指紋", "拷貝安全指紋"),
    ("unpairEllipsis", "Unpair…", "解除配对…", "ペアリングを解除…", "解除配對…", "解除配對…"),
    ("fingerprintCopied", "Fingerprint copied", "安全指纹已拷贝", "指紋をコピーしました", "安全指紋已複製", "安全指紋已拷貝"),
    ("unpairTitle", "Unpair “{name}”?", "解除与「{name}」的配对？", "「{name}」とのペアリングを解除しますか？", "要解除與「{name}」的配對嗎？", "解除與「{name}」的配對？"),
    ("unpairBody", "The two devices will stop syncing. To sync again, you'll need to pair them and check the fingerprints again.", "两台设备将不再同步。以后想恢复，需要重新配对并核对指纹。", "2 台のデバイスは同期しなくなります。再開するには、もう一度ペアリングして指紋を確かめる必要があります。", "兩台裝置將不再同步。以後想恢復，需要重新配對並核對指紋。", "兩部裝置將不再同步。日後想恢復，需要重新配對並核對指紋。"),
    ("unpair", "Unpair", "解除配对", "ペアリングを解除", "解除配對", "解除配對"),
    ("unpaired", "Unpaired", "已解除配对", "ペアリングを解除しました", "已解除配對", "已解除配對"),
    ("enterAddressEllipsis", "Enter its address…", "输入它的地址…", "アドレスを入力…", "輸入它的位址…", "輸入它的地址…"),
    ("setAddressTitle", "Address of “{name}”", "「{name}」的地址", "「{name}」のアドレス", "「{name}」的位址", "「{name}」的地址"),
    ("setAddressBody", "Use this when the two devices are on different subnets, or both changed addresses. Enter the address shown on “{name}” under Add device → Show a code.", "两台不在同一网段、或者都换了地址时用。输入「{name}」上「添加设备 → 生成配对码」下面显示的地址。", "2 台が別のサブネットにある場合や、どちらもアドレスが変わった場合に使います。「{name}」の「デバイスを追加 → コードを表示」に表示されるアドレスを入力してください。", "兩台不在同一網段、或都換了位址時使用。輸入「{name}」上「新增裝置 → 產生配對碼」下方顯示的位址。", "兩部不在同一網段，或者都換了地址時使用。輸入「{name}」上「新增裝置 → 產生配對碼」下面顯示的地址。"),
    ("addressSaved", "Looking for it at that address", "已记下，正在那个地址上找它", "そのアドレスで探しています", "已記下，正在那個位址上尋找它", "已記下，正在那個地址上尋找它"),
    ("lanBlockedTitle", "macOS is blocking the local network", "macOS 阻止了局域网访问", "macOS がローカルネットワークへのアクセスをブロックしています", "macOS 阻擋了區域網路存取", "macOS 阻擋了區域網絡存取"),
    ("lanBlockedBody", "Turn on CopySync Daemon in System Settings → Privacy & Security → Local Network. Until then, {device} can't connect to devices on this network.", "在「系统设置 → 隐私与安全性 → 本地网络」里打开 CopySync Daemon。打开之前，{device}连不上同一网络里的设备。", "「システム設定 → プライバシーとセキュリティ → ローカルネットワーク」で CopySync Daemon をオンにしてください。オンにするまで、{device}は同じネットワークのデバイスに接続できません。", "在「系統設定 → 隱私權與安全性 → 本機網路」裡開啟 CopySync Daemon。開啟之前，{device}連不上同一個網路裡的裝置。", "在「系統設定 → 私隱與保安 → 本地網絡」裏開啟 CopySync Daemon。開啟之前，{device}連不上同一網絡裏的裝置。"),

    # ── 搜索附近设备 ──
    ("searchNearby", "Search nearby", "搜索附近设备", "近くを検索", "搜尋附近裝置", "搜尋附近裝置"),
    ("nearbyTitle", "Nearby devices", "附近的设备", "近くのデバイス", "附近的裝置", "附近的裝置"),
    ("nearbySearching", "Looking for nearby devices…", "正在查找附近的设备…", "近くのデバイスを探しています…", "正在尋找附近的裝置…", "正在尋找附近的裝置…"),
    ("nearbyHint", "Open Search nearby on the other computer too, and it shows up here. Click it to pair.", "在另一台电脑上也打开「搜索附近设备」，它就会出现在这里，点一下即可配对。", "もう一台でも「近くを検索」を開くと、ここに表示されます。クリックしてペアリングします。", "在另一台電腦上也打開「搜尋附近裝置」，它就會出現在這裡，點一下即可配對。", "在另一部電腦上也開啟「搜尋附近裝置」，它就會出現在這裏，按一下即可配對。"),
    ("nearbyVisible", "While this window is open, devices on this network can see the name of {device}.", "窗口打开期间，同一网络里的设备能看到{device}的名字。", "このウィンドウを開いている間は、同じネットワークのデバイスから{device}の名前が見えます。", "這個視窗開著的期間，同一個網路裡的裝置能看到{device}的名稱。", "視窗開啟期間，同一網絡裏的裝置能看到{device}的名稱。"),
    ("nearbyPaired", "Paired", "已配对", "ペアリング済み", "已配對", "已配對"),
    ("nearbyRequesting", "Contacting “{name}”…", "正在联系「{name}」…", "「{name}」に接続しています…", "正在聯絡「{name}」…", "正在聯絡「{name}」…"),
    ("nearbyNoReply", "“{name}” didn't answer. Make sure Search nearby is still open on it.", "「{name}」没有回应。确认它还开着「搜索附近设备」。", "「{name}」から応答がありません。相手で「近くを検索」が開いたままか確認してください。", "「{name}」沒有回應。請確認它還開著「搜尋附近裝置」。", "「{name}」沒有回應。請確認它仍開着「搜尋附近裝置」。"),
    ("nearbyLanOffTitle", "Local network discovery is off", "局域网发现已关闭", "ローカルネットワークの検出がオフです", "區域網路探索已關閉", "區域網絡發現已關閉"),
    ("nearbyLanOffBody", "Turn on “Find devices on this network” in Settings to see nearby devices, or pair with a code instead.", "在「设置」里打开「在局域网中查找设备」才能看到附近的设备，或者改用配对码。", "近くのデバイスを表示するには「設定」で「このネットワークでデバイスを探す」をオンにするか、コードでペアリングしてください。", "在「設定」裡開啟「在區域網路中尋找裝置」才能看到附近的裝置，或者改用配對碼。", "在「設定」裏開啟「在區域網絡中尋找裝置」才能看到附近的裝置，或者改用配對碼。"),

    # ── 配对 ──
    ("generateCode", "Show a code", "生成配对码", "コードを表示", "產生配對碼", "產生配對碼"),
    ("enterCode", "Enter a code", "输入配对码", "コードを入力", "輸入配對碼", "輸入配對碼"),
    ("codeFailedTitle", "Couldn't create a pairing code", "没能生成配对码", "ペアリングコードを作成できませんでした", "無法產生配對碼", "未能產生配對碼"),
    ("enterOnOther", "Enter this code on your other device", "在另一台设备上输入这个配对码", "もう一台のデバイスでこのコードを入力してください", "在另一台裝置上輸入這組配對碼", "在另一部裝置上輸入這個配對碼"),
    ("codeCopied", "Code copied", "配对码已拷贝", "コードをコピーしました", "配對碼已複製", "配對碼已拷貝"),
    ("codeHint", "Click to copy. Valid for 5 minutes.", "点击可拷贝，5 分钟内有效", "クリックでコピー。有効期間は 5 分です。", "點一下即可複製，5 分鐘內有效", "按一下即可拷貝，5 分鐘內有效"),
    ("afterEnterFootnote", "Once it's entered, {both} show the same two lines of security fingerprints. Confirm only if they match.", "对方输入后，{both}会显示同样的两行安全指纹，核对一致再确认。", "入力されると、{both}に同じ 2 行のセキュリティ指紋が表示されます。一致を確かめてから確定してください。", "對方輸入後，{both}會顯示同樣的兩行安全指紋，核對一致再確認。", "對方輸入後，{both}會顯示同樣的兩行安全指紋，核對一致再確認。"),
    ("codeInvalid", "A pairing code is 6 letters or digits", "配对码是 6 位字母或数字", "ペアリングコードは 6 桁の英数字です", "配對碼是 6 個字母或數字", "配對碼是 6 位字母或數字"),
    ("localAddress", "Address", "本机地址", "アドレス", "本機位址", "本機地址"),
    ("addressHint", "On a different subnet? Enter this address on the other device as well.", "不在同一网段时，另一台还要输入这个地址", "別のサブネットにある場合は、もう一台でこのアドレスも入力します", "不在同一網段時，另一台還要輸入這個位址", "不在同一網段時，另一部還要輸入這個地址"),
    ("addressCopied", "Address copied", "地址已拷贝", "アドレスをコピーしました", "位址已複製", "地址已拷貝"),
    ("differentSubnet", "On a different subnet?", "不在同一网段？", "別のサブネットにありますか？", "不在同一網段？", "不在同一網段？"),
    ("peerAddress", "Address shown under the code", "对方配对码下面显示的地址", "コードの下に表示されたアドレス", "對方配對碼下方顯示的位址", "對方配對碼下面顯示的地址"),
    ("pairNotFound", "Couldn't find that device. If it's on a different subnet, enter the address shown under its code.", "没找到那台设备。不在同一网段时，请输入它配对码下面显示的地址。", "そのデバイスが見つかりません。別のサブネットにある場合は、コードの下に表示されたアドレスを入力してください。", "找不到那台裝置。不在同一網段時，請輸入它配對碼下方顯示的位址。", "找不到那部裝置。不在同一網段時，請輸入它配對碼下面顯示的地址。"),
    ("pairNotFoundAddress", "Nothing answered at that address. Check the address and that the other device is still showing its code. A firewall on that device may also be blocking it.", "那个地址上没有回应。检查地址是否正确、对方是否还停在配对码页面；对方的防火墙也可能拦住了。", "そのアドレスから応答がありません。アドレスと、相手がまだコードを表示しているかを確認してください。相手のファイアウォールがブロックしている可能性もあります。", "那個位址沒有回應。請檢查位址是否正確、對方是否還停在配對碼畫面；也可能是對方的防火牆擋住了。", "那個地址沒有回應。請檢查地址是否正確、對方是否仍停留在配對碼頁面；對方的防火牆也可能阻擋了連線。"),
    ("codeWrong", "That code isn't right, or it has expired", "配对码不对，或者已经失效", "コードが違うか、有効期限が切れています", "配對碼不正確，或已經失效", "配對碼不正確，或者已經失效"),
    ("codeExpired", "This code has expired. Show a new one on the other device.", "配对码已失效，请在另一台上重新生成", "このコードは無効になりました。もう一台で新しいコードを表示してください。", "配對碼已失效，請在另一台上重新產生", "配對碼已失效，請在另一部上重新產生"),
    ("addressInvalid", "That address isn't valid. Enter something like 192.168.1.20.", "地址格式不对，应当像 192.168.1.20 这样", "アドレスの形式が正しくありません。192.168.1.20 のように入力してください。", "位址格式不正確，應該像 192.168.1.20 這樣", "地址格式不正確，應該像 192.168.1.20 這樣"),
    ("pairNoChannel", "Turn on local network discovery or set a server in Settings first.", "请先在「设置」里开启局域网发现，或者填写服务器地址。", "先に「設定」でローカルネットワークの検出をオンにするか、サーバーを設定してください。", "請先在「設定」裡開啟區域網路探索，或者填寫伺服器位址。", "請先在「設定」裏開啟區域網絡發現，或者填寫伺服器地址。"),
    ("enterCodeShown", "Enter the code shown on your other device", "输入另一台设备上显示的配对码", "もう一台のデバイスに表示されたコードを入力してください", "輸入另一台裝置上顯示的配對碼", "輸入另一部裝置上顯示的配對碼"),
    ("pairingDone", "Paired", "配对完成", "ペアリングしました", "配對完成", "配對完成"),
    ("pairingRejected", "Pairing declined", "已拒绝配对", "ペアリングを拒否しました", "已拒絕配對", "已拒絕配對"),
    ("verifyTitle", "Check the security fingerprints", "核对安全指纹", "セキュリティ指紋を確認", "核對安全指紋", "核對安全指紋"),
    ("pairingWith", "Pairing with “{name}”", "正在与「{name}」配对", "「{name}」とペアリング中", "正在與「{name}」配對", "正在與「{name}」配對"),
    ("verifyBody", "The two lines on {both} must be identical. If even one character differs, someone may be tampering with the connection, so decline.", "{both}上显示的这两行应当完全相同。有任何一个字符不同，说明连接可能被第三方篡改，请拒绝。", "{both}に表示された 2 行は完全に同じはずです。1 文字でも違う場合は、第三者が通信を改ざんしている可能性があるため、拒否してください。", "{both}上顯示的這兩行應該完全相同。只要有任何一個字元不同，就表示連線可能遭到第三方竄改，請拒絕。", "{both}上顯示的這兩行應該完全相同。只要有一個字元不同，就代表連線可能被第三方篡改，請拒絕。"),
    ("rejectMismatch", "Mismatch, decline", "不一致，拒绝", "不一致、拒否", "不一致，拒絕", "不一致，拒絕"),
    ("confirmMatch", "They match, pair", "一致，确认配对", "一致、ペアリング", "一致，確認配對", "一致，確認配對"),

    # ── 设置 ──
    ("server", "Server (optional)", "服务器（可选）", "サーバー（任意）", "伺服器（選用）", "伺服器（可選）"),
    ("serverFootnote", "Only needed to sync between different networks, such as home and office. Devices on the same network connect without it. The server only helps devices find each other and can't see what they send.", "只在不同网络之间同步时需要，比如家里和公司；同一网络里的设备不用服务器也能直接连上。服务器只负责让设备找到彼此，看不到传输的内容。", "家と会社のように、別々のネットワーク間で同期するときだけ必要です。同じネットワーク内のデバイスはサーバーなしで直接つながります。サーバーはデバイス同士を引き合わせるだけで、やり取りする内容は見えません。", "只在不同網路之間同步時需要，例如家裡和公司；同一個網路裡的裝置不用伺服器也能直接連上。伺服器只負責讓裝置找到彼此，看不到傳輸的內容。", "只在不同網絡之間同步時需要，例如家裏和公司；同一網絡裏的裝置不用伺服器也能直接連上。伺服器只負責讓裝置找到彼此，看不到傳輸的內容。"),
    ("signalingUrl", "Signaling server", "信令服务器地址", "シグナリングサーバー", "信令伺服器位址", "信令伺服器地址"),
    ("signalingUrlDesc", "Press Return to save. Leave it empty to sync on this network only.", "按回车保存。留空表示只在局域网里同步", "Enter で保存。空欄にするとこのネットワーク内だけで同期します", "按 Enter 儲存。留空表示只在區域網路裡同步", "按 Enter 儲存。留空表示只在區域網絡裏同步"),
    ("signalingUrlPlaceholder", "ws://server-address:8787/signal", "ws://服务器地址:8787/signal", "ws://サーバーアドレス:8787/signal", "ws://伺服器位址:8787/signal", "ws://伺服器地址:8787/signal"),
    ("publicStun", "Probe with public servers", "用公共服务器探测网络出口", "公開サーバーで出口を調べる", "用公共伺服器偵測網路出口", "用公共伺服器探測網絡出口"),
    ("publicStunDesc", "On networks with several uplinks, such as offices with two ISPs, this finds more direct paths. Public servers only see your public IP.", "公司双线这类多出口网络里，能找到更多可以直连的路径。公共服务器只会看到你的公网 IP", "複数の回線を持つ会社などのネットワークで、直接つながる経路を多く見つけます。公開サーバーに見えるのはグローバル IP だけです。", "在公司雙線路這類多出口網路裡，能找到更多可以直連的路徑。公共伺服器只會看到你的公用 IP", "在公司雙線這類多出口網絡裏，能找到更多可以直連的路徑。公共伺服器只會看到你的公網 IP"),
    ("connectionStatus", "Connection", "连接状态", "接続状態", "連線狀態", "連線狀態"),
    ("connectedDesc", "Connected to the signaling server", "已连接到信令服务器", "シグナリングサーバーに接続しています", "已連線到信令伺服器", "已連線到信令伺服器"),
    ("disconnectedDesc", "Can't reach the server. Check the address and that the server is running.", "连不上服务器。检查地址是否正确、服务器是否在运行。", "サーバーに接続できません。アドレスとサーバーの稼働を確認してください。", "連不上伺服器。請檢查位址是否正確、伺服器是否在執行。", "連不上伺服器。請檢查地址是否正確、伺服器是否正在執行。"),
    ("connected", "Connected", "已连接", "接続済み", "已連線", "已連線"),
    ("notConnected", "Not connected", "未连接", "未接続", "未連線", "未連線"),
    ("notSet", "Not set", "未设置", "未設定", "未設定", "未設定"),
    ("notSetDesc", "Syncing with devices on this network only", "只与同一网络里的设备同步", "同じネットワークのデバイスとだけ同期します", "只與同一個網路裡的裝置同步", "只與同一網絡裏的裝置同步"),
    ("serverCleared", "Server removed. Syncing on this network only.", "已清除服务器地址，只在局域网里同步", "サーバーを削除しました。このネットワーク内だけで同期します", "已清除伺服器位址，只在區域網路裡同步", "已清除伺服器地址，只在區域網絡裏同步"),
    ("localNetwork", "Local network", "局域网", "ローカルネットワーク", "區域網路", "區域網絡"),
    ("lanFootnote", "Uses UDP port 8788. If the devices are on different subnets, enter the other device's address when pairing.", "使用 UDP 端口 8788。两台不在同一网段时，配对时输入对方的地址。", "UDP ポート 8788 を使います。2 台が別のサブネットにある場合は、ペアリング時に相手のアドレスを入力します。", "使用 UDP 連接埠 8788。兩台不在同一網段時，配對時輸入對方的位址。", "使用 UDP 端口 8788。兩部不在同一網段時，配對時要輸入對方的地址。"),
    ("lanDiscovery", "Find devices on this network", "在局域网中查找设备", "このネットワークでデバイスを探す", "在區域網路中尋找裝置", "在區域網絡中尋找裝置"),
    ("lanDiscoveryDesc", "Pair and sync without a server. Others on the network only see an anonymous entry that changes daily, and the device name only while Search nearby is open.", "不用服务器就能配对和同步。网络里的其他人只能看到一个每天都变的匿名标识；打开「搜索附近设备」期间才看得到设备名", "サーバーなしでペアリングと同期ができます。ネットワーク上の他の人には、毎日変わる匿名の識別子しか見えません。デバイス名が見えるのは「近くを検索」を開いている間だけです", "不用伺服器就能配對和同步。網路裡的其他人只能看到一個每天都會變的匿名識別碼；打開「搜尋附近裝置」期間才看得到裝置名稱", "不用伺服器就能配對和同步。網絡裏的其他人只能看到一個每天都會改變的匿名標識；開啟「搜尋附近裝置」期間才看得到裝置名稱"),
    ("macLocalNetwork", "Local Network access", "本地网络权限", "ローカルネットワークへのアクセス", "本機網路權限", "本地網絡權限"),
    ("macLocalNetworkDesc", "macOS asks once whether CopySync may connect to devices on your network. If you declined, turn on CopySync Daemon there.", "macOS 会问一次是否允许 CopySync 连接网络里的设备。拒绝过的话，到那里打开 CopySync Daemon", "CopySync がネットワーク上のデバイスに接続してよいか、macOS が一度だけ確認します。拒否した場合は、そこで CopySync Daemon をオンにしてください", "macOS 會詢問一次是否允許 CopySync 連線到網路裡的裝置。拒絕過的話，到那裡開啟 CopySync Daemon", "macOS 會詢問一次是否允許 CopySync 連接網絡裏的裝置。如果拒絕過，請到那裏開啟 CopySync Daemon"),
    ("allowLan", "Windows Firewall", "Windows 防火墙", "Windows ファイアウォール", "Windows 防火牆", "Windows 防火牆"),
    ("allowLanDesc", "Allow this when other devices can't reach this one, for example when pairing from another subnet. Asks for administrator approval once.", "其他设备连不上这台时放行，比如从别的网段配对。会请求一次管理员权限", "他のデバイスからこのデバイスに接続できないとき（別のサブネットからのペアリングなど）に許可します。管理者の承認を 1 回求めます", "其他裝置連不上這台時允許連入，例如從別的網段配對。會要求一次系統管理員權限", "其他裝置連不上這部電腦時允許連線，例如從別的網段配對。會要求一次管理員權限"),
    ("allow", "Allow", "放行", "許可", "允許", "允許"),
    ("allowLanDone", "Allowed. Other devices can now reach this one.", "已放行，其他设备现在可以连到这台了", "許可しました。他のデバイスからこのデバイスに接続できます。", "已允許，其他裝置現在可以連到這台了", "已允許，其他裝置現在可以連到這部電腦了"),
    ("deviceName", "Device name", "设备名称", "デバイス名", "裝置名稱", "裝置名稱"),
    ("deviceNameDesc", "Shown in the list on your other devices", "显示在其他设备的列表里", "他のデバイスの一覧に表示されます", "顯示在其他裝置的清單裡", "顯示在其他裝置的列表裏"),
    ("language", "Language", "语言", "言語", "語言", "語言"),
    ("languageSystem", "System", "跟随系统", "システムに合わせる", "跟隨系統", "跟隨系統"),
    ("sync", "Sync", "同步", "同期", "同步", "同步"),
    ("syncFootnote", "Items under the limit are pushed to your other devices as you copy. Larger ones sync as a record; pull them from History when you need them, so big files don't waste bandwidth or disk space.", "小于阈值的内容在复制时直接推送到其他设备；超过阈值的只同步一条记录，需要时在「复制记录」里点「拉取到本机」，避免大文件无谓地占用带宽和磁盘。", "上限より小さい内容はコピーした時点で他のデバイスに送られます。大きいものは記録だけを同期し、必要なときに「履歴」から取り込むので、帯域やディスクを無駄に使いません。", "小於上限的內容在複製時直接推送到其他裝置；超過上限的只同步一筆記錄，需要時在「複製記錄」裡點「拉取到本機」，避免大檔案白白佔用頻寬和磁碟空間。", "小於上限的內容在複製時直接推送到其他裝置；超過上限的只同步一條記錄，需要時在「複製記錄」裏按「拉取到本機」，避免大檔案無謂地佔用頻寬和磁碟空間。"),
    ("autoSyncLimit", "Auto-sync limit", "自动同步上限", "自動同期の上限", "自動同步上限", "自動同步上限"),
    ("autoSyncLimitDesc", "Larger files are pulled by hand", "超过这个大小的文件改为手动拉取", "これより大きいファイルは手動で取り込みます", "超過這個大小的檔案改為手動拉取", "超過這個大小的檔案改為手動拉取"),
    ("autoApply", "Put received items on the clipboard", "收到后直接写入剪贴板", "受け取った内容をクリップボードに入れる", "收到後直接寫入剪貼簿", "收到後直接寫入剪貼簿"),
    ("autoApplyDesc", "When off, put them on the clipboard from History yourself", "关闭后，需要在「复制记录」里手动放入剪贴板", "オフの場合は「履歴」から手動でクリップボードに入れます", "關閉後，需要在「複製記錄」裡手動放入剪貼簿", "關閉後，需要在「複製記錄」裏手動放入剪貼簿"),
    ("whatToSync", "What to sync", "同步哪些内容", "同期する内容", "同步哪些內容", "同步哪些內容"),
    ("syncFiles", "Files and folders", "文件与文件夹", "ファイルとフォルダ", "檔案與資料夾", "檔案與資料夾"),
    ("syncText", "Plain text", "纯文本", "プレーンテキスト", "純文字", "純文字"),
    ("syncRich", "Rich text", "带格式的文本", "書式付きテキスト", "帶格式的文字", "帶格式的文字"),
    ("syncImages", "Images and screenshots", "图片与截图", "画像とスクリーンショット", "圖片與螢幕截圖", "圖片與截圖"),
    ("storage", "Storage", "存储", "ストレージ", "儲存空間", "儲存空間"),
    ("storageFootnote", "Expired items and cached files are cleaned up automatically.", "到期的记录与缓存文件会被自动清理。", "期限切れの履歴とキャッシュは自動で削除されます。", "到期的記錄與快取檔案會自動清除。", "到期的記錄與快取檔案會自動清理。"),
    ("keepHistory", "Keep history for", "记录保留", "履歴の保存期間", "記錄保留", "記錄保留"),
    ("keepCache", "Keep cached files for", "缓存文件保留", "キャッシュの保存期間", "快取檔案保留", "快取檔案保留"),
    ("keepCacheDesc", "Never longer than the history", "不会超过记录的保留时长", "履歴の保存期間を超えません", "不會超過記錄的保留時間", "不會超過記錄的保留時間"),
    ("cacheUsage", "Cache size", "缓存占用", "キャッシュ容量", "快取佔用", "快取佔用空間"),
    ("backgroundSync", "Background sync", "后台同步", "バックグラウンド同期", "背景同步", "背景同步"),
    ("backgroundFootnote", "When off, nothing syncs and CopySync doesn't start at login. History and pairings are kept; open the app to turn it back on.", "停用后不再同步，也不会开机启动。历史记录与配对关系会保留，重新打开 App 即可再次启用。", "オフにすると同期せず、ログイン時にも起動しません。履歴とペアリングは残り、アプリを開けば再び有効にできます。", "停用後不再同步，也不會開機自動啟動。歷史記錄與配對關係會保留，重新打開 App 即可再次啟用。", "停用後不再同步，也不會開機啟動。歷史記錄與配對關係會保留，重新開啟 App 即可再次啟用。"),
    ("backgroundService", "Background service", "后台服务", "バックグラウンドサービス", "背景服務", "背景服務"),
    ("backgroundServiceDesc", "Starts at login and keeps syncing with the window closed", "随系统登录自动启动，关掉窗口也照常同步", "ログイン時に起動し、ウィンドウを閉じても同期を続けます", "隨系統登入自動啟動，關掉視窗也照常同步", "隨系統登入自動啟動，關閉視窗也照常同步"),
    ("restart", "Restart", "重新启动", "再起動", "重新啟動", "重新啟動"),
    ("restarted", "Background service restarted", "后台服务已重新启动", "バックグラウンドサービスを再起動しました", "背景服務已重新啟動", "背景服務已重新啟動"),
    ("disableEllipsis", "Turn off…", "停用…", "オフにする…", "停用…", "停用…"),
    ("about", "About", "关于", "このアプリについて", "關於", "關於"),
    ("aboutDesc", "Open source under the AGPL-3.0", "开源软件，AGPL-3.0 许可", "オープンソース（AGPL-3.0）", "開源軟體，採用 AGPL-3.0 授權", "開源軟件，以 AGPL-3.0 授權"),
    ("homepage", "Website", "项目主页", "ウェブサイト", "專案首頁", "項目主頁"),
    ("feedback", "Report a problem", "反馈问题", "問題を報告", "回報問題", "回報問題"),
    ("disableTitle", "Turn off background sync?", "停用后台同步？", "バックグラウンド同期をオフにしますか？", "要停用背景同步嗎？", "停用背景同步？"),
    ("disableBody", "CopySync on {device} will stop syncing and won't start at login. History and pairings are kept.", "{device}将停止同步，也不再开机启动。历史记录和配对关系会保留。", "{device}は同期を停止し、ログイン時にも起動しなくなります。履歴とペアリングは残ります。", "{device}將停止同步，也不會再開機自動啟動。歷史記錄和配對關係會保留。", "{device}將停止同步，也不再開機啟動。歷史記錄和配對關係會保留。"),
    ("disable", "Turn off", "停用", "オフにする", "停用", "停用"),
    ("disabled", "Background sync is off", "后台同步已停用", "バックグラウンド同期をオフにしました", "背景同步已停用", "背景同步已停用"),
    ("urlMustBeWs", "The address must start with ws:// or wss://", "地址需以 ws:// 或 wss:// 开头", "アドレスは ws:// または wss:// で始めてください", "位址必須以 ws:// 或 wss:// 開頭", "地址必須以 ws:// 或 wss:// 開頭"),
    ("serverUpdated", "Server updated. Reconnecting…", "服务器地址已更新，正在重连", "サーバーを更新しました。再接続しています…", "伺服器位址已更新，正在重新連線", "伺服器地址已更新，正在重新連線"),
    ("deviceNameUpdated", "Device name updated", "设备名称已更新", "デバイス名を更新しました", "裝置名稱已更新", "裝置名稱已更新"),

    # ── 启用后台服务的引导页 ──
    ("svcExtractTitle", "Install or extract CopySync first", "先解压或安装 CopySync", "先に CopySync をインストールまたは展開してください", "先解壓縮或安裝 CopySync", "先解壓縮或安裝 CopySync"),
    ("svcExtractBody", "CopySync is running straight from the zip, and its temporary folder disappears when you close it. Run the installer, or extract the whole folder and open it from there, to turn on background sync.", "CopySync 现在是直接从压缩包里打开的，关掉之后它所在的临时文件夹就会被删除。运行安装程序，或者把整个文件夹解压出来再打开，就可以启用后台同步。", "CopySync は zip から直接開かれていて、閉じると一時フォルダごと消えます。インストーラーを実行するか、フォルダ全体を展開してから開くと、バックグラウンド同期を有効にできます。", "CopySync 現在是直接從壓縮檔裡打開的，關掉之後它所在的暫存資料夾就會被刪除。執行安裝程式，或者把整個資料夾解壓縮出來再打開，就可以啟用背景同步。", "CopySync 現在是直接從壓縮檔裏開啟的，關閉之後它所在的臨時資料夾就會被刪除。執行安裝程式，或者把整個資料夾解壓縮出來再開啟，就可以啟用背景同步。"),
    ("svcMoveTitle", "Move CopySync to Applications first", "先把 CopySync 移到「应用程序」", "先に CopySync を「アプリケーション」に移動してください", "先把 CopySync 移到「應用程式」", "先把 CopySync 移到「應用程式」"),
    ("svcMoveBody", "CopySync is running straight from the disk image. Drag it into Applications, eject the disk image, and open it from Applications to turn on background sync.", "CopySync 现在是直接从安装盘里打开的。把它拖进「应用程序」文件夹，推出安装盘，再从「应用程序」里打开，就可以启用后台同步。", "CopySync はディスクイメージから直接開かれています。「アプリケーション」にドラッグしてディスクイメージを取り出し、「アプリケーション」から開くと、バックグラウンド同期を有効にできます。", "CopySync 現在是直接從磁碟映像檔裡打開的。把它拖進「應用程式」資料夾，退出磁碟映像檔，再從「應用程式」裡打開，就可以啟用背景同步。", "CopySync 現在是直接從磁碟映像裏開啟的。把它拖進「應用程式」資料夾，推出磁碟映像，再從「應用程式」裏開啟，就可以啟用背景同步。"),
    ("svcStartTitle", "Get started with CopySync", "开始使用 CopySync", "CopySync を始める", "開始使用 CopySync", "開始使用 CopySync"),
    ("svcStartBody", "CopySync runs a small background process that notices clipboard changes and stays connected to your other devices. Once on, it starts at login and keeps syncing even with this window closed.", "CopySync 会在后台运行一个小进程，用来感知剪贴板的变化，并与你的其他设备保持连接。启用后它随系统登录自动启动，关掉这个窗口也会照常同步。", "CopySync は小さなバックグラウンドプロセスを動かし、クリップボードの変化を検知して他のデバイスとつながり続けます。有効にするとログイン時に起動し、このウィンドウを閉じても同期を続けます。", "CopySync 會在背景執行一個小程式，用來偵測剪貼簿的變化，並與你的其他裝置保持連線。啟用後它會隨系統登入自動啟動，關掉這個視窗也會照常同步。", "CopySync 會在背景執行一個小型程序，用來偵測剪貼簿的變化，並與你的其他裝置保持連線。啟用後它會隨系統登入自動啟動，關閉這個視窗也會照常同步。"),
    ("svcEnable", "Turn on background sync", "启用后台同步", "バックグラウンド同期を有効にする", "啟用背景同步", "啟用背景同步"),
    ("svcStaleTitle", "Turn background sync on again", "重新启用后台同步", "バックグラウンド同期をもう一度有効にする", "重新啟用背景同步", "重新啟用背景同步"),
    ("svcStaleBodyWin", "CopySync has moved, but starting at login still points to the old place. Turn it on again so the background process uses this copy. History and pairings are kept.", "CopySync 换了位置，开机自启还指向原来的地方。重新启用一次，让后台进程指向当前这份 CopySync。历史记录与配对关系都会保留。", "CopySync の場所が変わりましたが、自動起動は以前の場所を指しています。もう一度有効にすると、このコピーを使うようになります。履歴とペアリングは残ります。", "CopySync 換了位置，開機自動啟動還指向原來的地方。重新啟用一次，讓背景程序指向目前這份 CopySync。歷史記錄與配對關係都會保留。", "CopySync 換了位置，開機自動啟動仍指向原來的位置。重新啟用一次，讓背景程序指向目前這份 CopySync。歷史記錄與配對關係都會保留。"),
    ("svcStaleBodyMac", "CopySync has moved, or was installed with an old script. Turn it on again so the background process uses this app. History and pairings are kept.", "CopySync 换了位置，或者之前用旧版安装脚本装过。重新启用一次，让后台进程指向当前这个 App。历史记录与配对关系都会保留。", "CopySync の場所が変わったか、古いスクリプトでインストールされています。もう一度有効にすると、このアプリを使うようになります。履歴とペアリングは残ります。", "CopySync 換了位置，或者之前用舊版安裝腳本裝過。重新啟用一次，讓背景程序指向目前這個 App。歷史記錄與配對關係都會保留。", "CopySync 換了位置，或者之前用舊版安裝腳本安裝過。重新啟用一次，讓背景程序指向目前這個 App。歷史記錄與配對關係都會保留。"),
    ("svcReenable", "Turn on again", "重新启用", "もう一度有効にする", "重新啟用", "重新啟用"),
    ("svcStoppedTitle", "The background service isn't running", "后台服务没有运行", "バックグラウンドサービスが動いていません", "背景服務沒有在執行", "背景服務沒有執行"),
    ("svcStoppedBody", "It may have just quit; the system usually restarts it within seconds. If you stay on this page, restart it with the button below.", "它可能刚刚退出，系统通常会在几秒内把它重新拉起。如果一直停在这一页，点下面的按钮重新启动。", "終了した直後かもしれません。通常は数秒で再起動されます。この画面のままなら、下のボタンで再起動してください。", "它可能剛剛結束，系統通常會在幾秒內把它重新啟動。如果一直停在這一頁，請點下面的按鈕重新啟動。", "它可能剛剛結束，系統通常會在幾秒內重新啟動它。如果一直停留在這一頁，請按下面的按鈕重新啟動。"),
    ("svcDevBodyWin", "This is a development build with no copysyncd.exe next to it. Run scripts\\build-windows.ps1 in the source folder to build a complete version.", "这是开发构建，旁边没有 copysyncd.exe。在源码目录运行 scripts\\build-windows.ps1 构建完整的版本。", "開発ビルドのため copysyncd.exe がありません。ソースフォルダで scripts\\build-windows.ps1 を実行して完全版をビルドしてください。", "這是開發版本，旁邊沒有 copysyncd.exe。在原始碼資料夾執行 scripts\\build-windows.ps1 建置完整的版本。", "這是開發版本，旁邊沒有 copysyncd.exe。在原始碼目錄執行 scripts\\build-windows.ps1 建置完整的版本。"),
    ("svcDevBodyLinux", "This is a development build with no copysyncd next to it. Run ./scripts/package-linux.sh in the source folder to build a complete version.", "这是开发构建，旁边没有 copysyncd。在源码目录运行 ./scripts/package-linux.sh 构建完整的版本。", "開発ビルドのため copysyncd がありません。ソースフォルダで ./scripts/package-linux.sh を実行して完全版をビルドしてください。", "這是開發版本，旁邊沒有 copysyncd。在原始碼資料夾執行 ./scripts/package-linux.sh 建置完整的版本。", "這是開發版本，旁邊沒有 copysyncd。在原始碼目錄執行 ./scripts/package-linux.sh 建置完整的版本。"),
    ("svcDevBodyMac", "This is a development build without the background service. Run ./scripts/install-macos.sh in the source folder to install it.", "这是开发构建，App 里没有内置后台服务。在源码目录运行 ./scripts/install-macos.sh 安装它。", "開発ビルドのためバックグラウンドサービスが含まれていません。ソースフォルダで ./scripts/install-macos.sh を実行してインストールしてください。", "這是開發版本，App 裡沒有內建背景服務。在原始碼資料夾執行 ./scripts/install-macos.sh 安裝它。", "這是開發版本，App 裏沒有內建背景服務。在原始碼目錄執行 ./scripts/install-macos.sh 安裝它。"),
    ("svcLogPath", "Log: {path}", "运行日志：{path}", "ログ：{path}", "記錄檔：{path}", "執行日誌：{path}"),
    ("svcCanDisable", "You can turn background sync off anytime in Settings.", "随时可以在「设置」中停用后台同步。", "バックグラウンド同期は「設定」でいつでもオフにできます。", "隨時可以在「設定」中停用背景同步。", "隨時可以在「設定」中停用背景同步。"),
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
    short = [k for k, *rest in S if len(rest) != 5]
    assert not short, f"这些条目没写全五种语言：{short}"

    en = {"@@locale": "en"}
    for key, text, zh, ja, tw, hk in S:
        names = placeholders(text) or placeholders(zh)
        for other in (zh, ja, tw, hk):
            names += [n for n in placeholders(other) if n not in names]
        en[key] = text
        if names:
            en["@" + key] = {
                "placeholders": {n: {"type": "int" if n in INTS else "String"} for n in names}
            }
    for locale, index in (("zh", 2), ("ja", 3), ("zh_Hant_TW", 4), ("zh_Hant_HK", 5)):
        data = {"@@locale": locale}
        for row in S:
            data[row[0]] = row[index]
        (out / f"app_{locale}.arb").write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    (out / "app_en.arb").write_text(json.dumps(en, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"✓ {len(S)} 条，写入 {out}")


if __name__ == "__main__":
    main()

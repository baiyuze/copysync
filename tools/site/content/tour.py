"""Tour copy. Every row is zh, zh-tw, zh-hk, en, ja; no locale fallback.

Keep film text separate from the player UI. The template uses semantic keys,
so translations never replace substrings or accidentally alter JavaScript.
"""

LANGUAGES = ("zh", "zh-tw", "zh-hk", "en", "ja")

COPY = {
    "opening_1": ("两台电脑之间复制粘贴，", "兩台電腦之間複製貼上，", "兩部電腦之間複製貼上，", "Copying between two computers", "2 台のパソコンでコピペするだけ。"),
    "opening_2": ("为什么还这么麻烦？", "為什麼還這麼麻煩？", "為甚麼還這麼麻煩？", "shouldn’t be this much work.", "もっと簡単にならない？"),
    "pain_chat": ("发给自己", "傳給自己", "傳給自己", "Sending it to yourself", "自分宛てに送信"),
    "pain_chat_detail": ("发过去，点开，下载，再找文件在哪", "傳過去，點開，下載，再找檔案在哪", "傳過去，開啟，下載，再找檔案在哪", "Send, open, download, then find the file", "送って、開いて、保存して、探して…"),
    "airdrop": ("隔空投送", "AirDrop", "AirDrop", "AirDrop", "AirDrop"),
    "pain_airdrop_detail": ("每次都要点接收，也进不了剪贴板", "每次都要接受，也不會進剪貼簿", "每次都要接受，也不會進剪貼簿", "Accept every time. Still no clipboard sync.", "毎回受信を許可。クリップボードには入らない"),
    "universal": ("通用剪贴板", "通用剪貼板", "通用剪貼板", "Universal Clipboard", "ユニバーサルクリップボード"),
    "pain_universal_detail": ("同一个 Apple ID，还得离得近", "同一個 Apple ID，還得靠得近", "同一個 Apple ID，還要靠得近", "Same Apple ID, and the devices must be close", "同じ Apple ID と、近くにある端末が必要"),
    "pain_cross_platform": ("Mac 和 Windows 之间", "Mac 和 Windows 之間", "Mac 和 Windows 之間", "Between Mac and Windows", "Mac と Windows の間では"),
    "pain_cross_platform_detail": ("系统自带的办法，一个都用不上", "系統內建的方法，一個都用不上", "系統內置的方法，一個都用不上", "None of the built-in options work", "標準機能だけでは、うまくいかない"),
    "simpler": ("其实，可以简单得多。", "其實，可以簡單得多。", "其實，可以簡單得多。", "It can be much simpler.", "もっとシンプルに。"),
    "tagline": ("在这台电脑复制，到那台电脑粘贴。", "在這台電腦複製，到那台電腦貼上。", "在這部電腦複製，到那部電腦貼上。", "Copy on this computer. Paste on that one.", "こっちでコピーして、あっちで貼り付け。"),
    "notes": ("备忘录", "備忘錄", "備忘錄", "Notes", "メモ"),
    "decision_1": ("会议结论：周二先发布 10%，", "會議結論：週二先發布 10%，", "會議結論：星期二先發佈 10%，", "Roll out to 10% on Tuesday.", "火曜日に 10% へリリース。"),
    "decision_2": ("观察两天没问题，周四全量。", "觀察兩天沒問題，週四全面發布。", "觀察兩天沒問題，星期四全面發佈。", "All users on Thursday if all goes well.", "問題なければ木曜日に全体へ。"),
    "screenshot": ("屏幕截图.png", "螢幕截圖.png", "螢幕截圖.png", "Screenshot.png", "スクリーンショット.png"),
    "notepad": ("会议记录.txt - 记事本", "會議記錄.txt - 記事本", "會議記錄.txt - 記事本", "Meeting notes.txt - Notepad", "会議メモ.txt - メモ帳"),
    "word": ("周报.docx - Word", "週報.docx - Word", "週報.docx - Word", "Weekly report.docx - Word", "週報.docx - Word"),
    "demo_file": ("Mac 上复制文件，Windows 上直接粘贴", "Mac 上複製檔案，Windows 上直接貼上", "Mac 上複製檔案，Windows 上直接貼上", "Copy a file on Mac. Paste it on Windows.", "Mac でファイルをコピー。Windows に貼り付け。"),
    "demo_text": ("反过来也一样：Windows 复制，Mac 粘贴", "反過來也一樣：Windows 複製，Mac 貼上", "反過來也一樣：Windows 複製，Mac 貼上", "And back: copy on Windows, paste on Mac.", "逆も同じ。Windows でコピー、Mac に貼り付け。"),
    "demo_image": ("截图、图片也是，粘贴出来就是图片", "截圖、圖片也一樣，貼上就是圖片", "截圖、圖片也一樣，貼上就是圖片", "Screenshots too. They paste as images.", "スクリーンショットも、そのまま画像として貼り付け。"),
    "types": ("复制什么，粘贴出来就是什么", "複製什麼，貼上就是什麼", "複製甚麼，貼上就是甚麼", "What you copy is what you paste", "コピーしたものを、そのまま貼り付け"),
    "text": ("文字", "文字", "文字", "Text", "テキスト"),
    "rich_text": ("带格式的文字", "格式化文字", "格式化文字", "Rich text", "書式付きテキスト"),
    "images": ("图片", "圖片", "圖片", "Images", "画像"),
    "files": ("文件", "檔案", "檔案", "Files", "ファイル"),
    "folders": ("整个文件夹", "整個資料夾", "整個資料夾", "Whole folders", "フォルダごと"),
    "small_files": ("小于 50 MB 的内容，复制时就开始传输；", "小於 50 MB 的內容，複製時就開始傳輸；", "小於 50 MB 的內容，複製時就開始傳輸；", "Under 50 MB? Transfer starts as you copy.", "50 MB 未満なら、コピーと同時に転送開始。"),
    "large_files": ("更大的文件先同步记录，需要时再拉取。", "更大的檔案先同步記錄，需要時再下載。", "較大的檔案先同步記錄，需要時再下載。", "Larger files arrive on demand, with one click.", "大きなファイルは、必要なときにワンクリックで取得。"),
    "history_1": ("所有设备复制过的，", "所有裝置複製過的，", "所有裝置複製過的，", "Everything you copied,", "どの端末のコピーも、"),
    "history_2": ("都在这里", "都在這裡", "都在這裏", "all in one place", "履歴はここに。"),
    "retention": ("默认保留最近 3 天，可以调", "預設保留最近 3 天，可調整", "預設保留最近 3 天，可調整", "3 days of history by default", "標準で 3 日分を保存。変更も可能"),
    "restore": ("点一下，放回剪贴板", "點一下，放回剪貼簿", "按一下，放回剪貼簿", "Click to copy it again", "クリックで、もう一度コピー"),
    "preview": ("图片点开就能预览", "點開圖片就能預覽", "開啟圖片就能預覽", "Preview images in place", "画像はその場でプレビュー"),
    "pull": ("大文件，需要时再拉到本机", "大檔案，需要時再下載到本機", "大檔案，需要時再下載到本機", "Pull large files when you need them", "大きなファイルは必要なときに取得"),
    "pairing": ("第一次用，配对一次就好", "第一次用，配對一次就好", "首次使用，配對一次就好", "Pair once, and you’re set", "一度ペアリングすれば、準備完了"),
    "nearby": ("同一网络搜索附近设备，点一下对方就能配对", "同一網路搜尋附近裝置，點一下就能配對", "同一網絡搜尋附近裝置，按一下就能配對", "On the same network, find a nearby computer and pair", "同じネットワークの端末を検索して、ペアリング"),
    "fingerprints": ("核对两边的安全指纹，一致再确认", "核對兩邊的安全指紋，一致才確認", "核對兩邊的安全指紋，一致才確認", "Check the fingerprints on both screens before confirming", "両方の画面でフィンガープリントが一致することを確認"),
    "networks": ("同一网络不用服务器，跨网络也能用", "同一網路不用伺服器，跨網路也能用", "同一網絡不用伺服器，跨網絡也能用", "One network or many. Your choice.", "同じネットワークでも、離れた場所でも"),
    "home": ("家里或办公室，", "家裡或辦公室，", "家裏或辦公室，", "At home or at the office,", "自宅やオフィスで、"),
    "lan_detail": ("装好就能用，不需要服务器", "裝好就能用，不需要伺服器", "安裝好就能用，不需要伺服器", "no server to set up", "サーバーの準備は不要"),
    "cloud_detail": ("自己部署，一台小云主机就够，", "自行部署，一台小雲端主機就夠，", "自行部署，一部小雲端主機就夠，", "Host a small cloud VM,", "小さなクラウド VM を用意して、"),
    "anywhere": ("在哪儿都能同步", "在哪裡都能同步", "在哪裏都能同步", "sync from anywhere", "どこからでも同期"),
    "across": ("跨网络，", "跨網路，", "跨網絡，", "Across networks,", "ネットワークをまたいでも、"),
    "tailscale_detail": ("也不需要服务器", "也不需要伺服器", "也不需要伺服器", "without your own server", "自前のサーバーは不要"),
    "encrypted": ("设备之间优先直连，端到端加密。", "裝置之間優先直連，端對端加密。", "裝置之間優先直連，端對端加密。", "A direct connection when possible. End-to-end encrypted.", "可能な限り端末同士で直接接続。エンドツーエンド暗号化。"),
    "relay": ("自建服务器模式下，直连失败时由你的服务器中转，内容仍然加密。", "自架伺服器模式下，直連失敗時由你的伺服器轉送，內容仍然加密。", "自建伺服器模式下，直連失敗時由你的伺服器轉送，內容仍然加密。", "In server mode, your server relays encrypted data if a direct link fails.", "サーバーモードでは、直接接続できない場合も暗号化したまま中継。"),
    "comparison": ("和常用的办法比一比", "和常用的方法比一比", "和常用的方法比一比", "How it compares", "いつもの方法と比べてみると"),
    "free": ("免费开源", "免費開源", "免費開源", "Free & open source", "無料・オープンソース"),
    "download": ("官网下载", "官網下載", "官網下載", "Download", "ダウンロード"),
    "source": ("源代码", "原始碼", "原始碼", "Source", "ソースコード"),
    "closing": ("少一步传来传去，多一点顺手。", "少一點傳來傳去，多一點順手。", "少一點傳來傳去，多一點順手。", "Less sending things around. More getting things done.", "送る手間を減らして、作業をもっとスムーズに。"),
}

UI = {
    "watch": ("看完整演示", "看完整示範", "觀看完整示範", "Watch the full tour", "紹介ムービーを見る"),
    "title": ("一分钟多，了解隔空贴", "一分多鐘，認識隔空貼", "一分多鐘，認識隔空貼", "Meet PasteFerry in 85 seconds", "85 秒でわかる PasteFerry"),
    "play": ("播放", "播放", "播放", "Play", "再生"),
    "pause": ("暂停", "暫停", "暫停", "Pause", "一時停止"),
    "replay": ("重新播放", "重新播放", "重新播放", "Replay", "もう一度再生"),
    "progress": ("播放进度", "播放進度", "播放進度", "Playback position", "再生位置"),
    "chapters": ("选择章节", "選擇章節", "選擇章節", "Choose a chapter", "チャプターを選択"),
    "fullscreen": ("全屏", "全螢幕", "全螢幕", "Full screen", "全画面"),
    "exit_fullscreen": ("退出全屏", "離開全螢幕", "離開全螢幕", "Exit full screen", "全画面を終了"),
    "close": ("关闭宣传片", "關閉宣傳片", "關閉宣傳片", "Close tour", "ムービーを閉じる"),
    "back": ("返回官网", "返回官網", "返回官網", "Back to the site", "公式サイトへ"),
    "summary": ("影片要点", "影片重點", "影片重點", "Tour highlights", "ムービーのポイント"),
    "silent": ("无声演示", "無聲示範", "無聲示範", "Silent tour", "音声なし"),
}

CHAPTERS = (
    ("复制的麻烦", "複製的麻煩", "複製的麻煩", "Too many steps", "コピペの悩み"),
    ("认识隔空贴", "認識隔空貼", "認識隔空貼", "Meet PasteFerry", "PasteFerry とは"),
    ("跨电脑粘贴", "跨電腦貼上", "跨電腦貼上", "Copy and paste", "端末間でコピペ"),
    ("支持的内容", "支援的內容", "支援的內容", "What you can sync", "同期できるもの"),
    ("复制历史", "複製歷史", "複製歷史", "Clipboard history", "クリップボード履歴"),
    ("安全配对", "安全配對", "安全配對", "Secure pairing", "安全にペアリング"),
    ("连接方式", "連線方式", "連線方式", "Ways to connect", "接続方法"),
    ("功能比较", "功能比較", "功能比較", "Compare options", "機能を比較"),
    ("开始使用", "開始使用", "開始使用", "Get started", "使ってみよう"),
)


def localized(code, site):
    column = LANGUAGES.index(code)
    copy = {key: row[column] for key, row in COPY.items()}
    ui = {key: row[column] for key, row in UI.items()}
    copy.update(brand=site["brand"], file=site["stage_file"],
                windows=site["stage_win"], direct=site["stage_arc"],
                lan=site["plans"]["lan"]["name"], cloud=site["plans"]["cloud"]["name"],
                chat=site["compare_cols"][3])
    # Reuse the carefully localized website comparison, including ✓/— labels.
    keys = ("apple_ids", "cross_platform", "different_networks", "instant",
            "files_folders", "history", "privacy")
    for key, (label, _) in zip(keys, site["compare_rows"]):
        copy[key] = label
    for key, row, col in (("accept", 0, 2), ("unreliable", 4, 1),
                          ("zipped", 4, 3), ("scroll", 5, 3)):
        copy[key] = site["compare_rows"][row][1][col][1]
    return copy, ui, [row[column] for row in CHAPTERS]

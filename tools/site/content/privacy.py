"""Localized public-service privacy notice. No installation analytics."""

COPY = {
    "zh": ("隐私与必要数据", [
        ("你决定是否使用官方服务", "新安装默认只使用局域网。开启官方公网后才连接官方信令并记录必要接入统计；自建服务器模式不向官方上报。先尝试直连，失败后才申请中转授权和读取机器额度摘要。"),
        ("仅保留服务所需数据", "官方网关保留脱敏设备编号、首次及最近接入时间、每日中转用量、必要额度绑定和累计汇总。每日设备明细保留 90 天。统计表示官方服务设备数，不能代表真实人数或全部安装量。"),
        ("内容与机器标识", "剪贴板、文件内容与真实设备名通过设备间加密通道传输。网关不收集文件名、路径、用户名、MAC、UA 或行为轨迹；原始系统标识不上传，机器摘要仅用于中转额度归组，不能保证不可伪造。"),
        ("网络可见性和额度", "联网端点必然能看到临时网络地址和路由信息，但官方统计库及普通代理日志不保存这些信息。默认每日 100 MB 中转额度，上下行合计，北京时间零点重置；直连不消耗额度。关闭官方服务会停止其请求、轮询与中转通道。"),
        ("网站与可审计边界", "本站使用既有的 Cloudflare Web Analytics，网站访问不等于软件用户统计。客户端、信令和 TURN 保持开源，公开所有官方网关协议与字段；私有生产后台的实际存储行为仍需要部署审计。"),
    ]),
    "zh-tw": ("隱私與必要資料", [
        ("由你決定使用官方服務", "新安裝預設僅使用區域網路。啟用官方公網後才連線官方信令並記錄必要連線統計；自訂伺服器不向官方回報。先嘗試直連，失敗才申請中繼並讀取機器額度摘要。"),
        ("最少必要資料", "網關保留去識別化裝置編號、首次與最近連線時間、每日中繼用量、必要額度綁定和累計彙總。每日裝置明細保留 90 天。這是官方服務裝置數，不是人數或全部安裝量。"),
        ("內容與機器識別碼", "內容與真實裝置名稱透過裝置間加密通道傳送。不收集檔名、路徑、使用者名稱、MAC、UA 或行為軌跡。原始系統識別碼不上傳；摘要只用於額度分組，並非不可偽造。"),
        ("網路資訊與額度", "網路端點暫時可見路由和網路位址，官方統計庫及一般代理日誌不保存這些資訊。每日 100 MB 中繼額度，含上下行，北京時間午夜重設；直連不計入。關閉官方服務會停止其請求、輪詢及通道。"),
        ("網站與稽核", "本站使用既有的 Cloudflare Web Analytics，瀏覽量不代表軟體使用者數。用戶端、信令及 TURN 持續開源，網關協定和欄位公開；私有正式環境的儲存行為仍需部署稽核。"),
    ]),
    "zh-hk": ("私隱與必要資料", [
        ("由你決定使用官方服務", "新安裝預設只使用區域網絡。啟用官方公網後才連接官方信令並記錄必要連線統計；自訂伺服器不向官方回報。先嘗試直連，失敗才申請中轉並讀取機器額度摘要。"),
        ("最少必要資料", "網關保留去識別化裝置編號、首次與最近連線時間、每日中轉用量、必要額度綁定和累計彙總。每日裝置明細保留 90 天。這是官方服務裝置數，不是人數或全部安裝量。"),
        ("內容與機器識別碼", "內容與真實裝置名稱透過裝置間加密通道傳送。不收集檔案名稱、路徑、用戶名稱、MAC、UA 或行為軌跡。原始系統識別碼不上傳；摘要只用於額度分組，並非不可偽造。"),
        ("網絡資料與額度", "網絡端點暫時可見路由和網絡地址，官方統計庫及一般代理日誌不保存這些資料。每日 100 MB 中轉額度，含上下行，北京時間午夜重設；直連不計入。關閉官方服務會停止其請求、輪詢及通道。"),
        ("網站與審計", "本站使用既有的 Cloudflare Web Analytics，瀏覽量不代表軟件用戶數。用戶端、信令及 TURN 持續開源，網關協定和欄位公開；私有正式環境的儲存行為仍需部署審計。"),
    ]),
    "en": ("Privacy and essential data", [
        ("Official service is optional", "Fresh installs use LAN mode. Official access and essential connection statistics begin only when enabled. Custom servers do not report to the official service. Direct connections are attempted first; machine quota identification and relay authorization happen only after direct failure."),
        ("Minimal service records", "The gateway retains pseudonymous device numbers, first/last access times, daily relay usage, essential quota bindings and cumulative totals. Daily device detail is retained for 90 days. Official device counts are not people or total installations."),
        ("Content and identifiers", "Content and real device names travel through the encrypted device channel. The gateway does not collect filenames, paths, usernames, MAC addresses, UA or behavior trails. Raw system identifiers are not uploaded; the application-specific digest groups relay quota and can be forged by a modified client."),
        ("Network visibility and allowance", "Network endpoints necessarily see transient routing/address metadata; official statistics and ordinary proxy logs do not retain it. The default allowance is 100 MB/day, upload plus download, reset at midnight Beijing time. Direct traffic is not charged. Disabling official access stops its requests, polling and tunnels."),
        ("Website and audit boundary", "This website uses its existing Cloudflare Web Analytics beacon. Visits are not software-user counts. The client and signaling/TURN server remain open source, with gateway fields and protocol published. Actual private production retention still requires deployment audit."),
    ]),
    "ja": ("プライバシーと必要なデータ", [
        ("公式サービスは任意です", "新規インストールではLANのみを使用します。有効にした場合のみ公式シグナリングに接続し、必要な接続統計を記録します。独自サーバーは公式へ報告しません。まず直接接続を試し、失敗した場合のみ中継の認証と端末の残量識別を行います。"),
        ("最小限の記録", "仮名化した端末番号、初回・最終接続時刻、日別中継使用量、必要な上限の紐付けと累計を保持します。端末の日別明細は90日間保持します。公式端末数は人数や全インストール数ではありません。"),
        ("内容と識別子", "内容と実際の端末名は端末間の暗号化通信で送信します。ファイル名、パス、ユーザー名、MAC、UA、行動履歴を収集しません。元のシステム識別子は送信せず、専用ダイジェストを中継上限の分類にのみ使います。改変クライアントによる偽装は防ぎきれません。"),
        ("ネットワークと残量", "通信先は一時的な経路・アドレスを認識しますが、公式統計と通常のプロキシログには保存しません。中継上限は送受信合計100 MB/日で、北京時間午前0時にリセットします。直接接続は消費しません。無効化すると公式の要求、照会とトンネルを停止します。"),
        ("ウェブサイトと監査", "当サイトは既存のCloudflare Web Analyticsを使います。閲覧数はソフトの利用者数ではありません。クライアント、シグナリングとTURNはオープンソースを維持し、ゲートウェイの仕様と項目を公開します。非公開の本番環境の保持実態には配備監査が必要です。"),
    ]),
}

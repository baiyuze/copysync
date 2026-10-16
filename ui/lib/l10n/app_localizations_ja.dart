// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Japanese (`ja`).
class AppLocalizationsJa extends AppLocalizations {
  AppLocalizationsJa([String locale = 'ja']) : super(locale);

  @override
  String get appName => 'PasteFerry';

  @override
  String get navHistory => '履歴';

  @override
  String get navDevices => 'デバイス';

  @override
  String get navSettings => '設定';

  @override
  String get linkConnecting => '接続中…';

  @override
  String get linkServiceNotRunning => 'サービス停止中';

  @override
  String get linkLanMode => 'このネットワークでデバイスを探しています';

  @override
  String get linkServerUnreachable => 'サーバーに接続できません。このネットワーク内では同期します';

  @override
  String get linkNothingOn => 'ローカルネットワークの検出がオフで、サーバーも未設定です';

  @override
  String get linkLanBlocked => 'macOS がローカルネットワークをブロックしています';

  @override
  String linkOnlinePeers(int count) {
    return '$count 台がオンライン';
  }

  @override
  String get linkServerConnected => 'サーバーに接続済み';

  @override
  String linkConnectionLost(String detail) {
    return '接続が切れました：$detail';
  }

  @override
  String get linkServiceDisconnected => 'サービスとの接続が切れました';

  @override
  String get linkServiceDisabled => '同期はオフです';

  @override
  String get errServiceNotConnected => 'バックグラウンドサービスに接続していません';

  @override
  String get errUnimplemented => 'この機能はまだ使えません';

  @override
  String get errUnavailable => 'バックグラウンドサービスに接続できません';

  @override
  String errServiceStart(String detail) {
    return 'バックグラウンドサービスを起動できませんでした：$detail';
  }

  @override
  String errServiceRestart(String detail) {
    return 'バックグラウンドサービスを再起動できませんでした：$detail';
  }

  @override
  String errAutostart(String detail) {
    return 'ログイン時の自動起動を設定できませんでした：$detail';
  }

  @override
  String get thisMacTitle => 'この Mac';

  @override
  String get thisPcTitle => 'このコンピューター';

  @override
  String get thisMacInText => 'この Mac ';

  @override
  String get thisPcInText => 'このコンピューター';

  @override
  String get bothMacs => '両方の Mac ';

  @override
  String get bothDevices => '両方のデバイス';

  @override
  String get thisDeviceShort => 'このデバイス';

  @override
  String fromDevice(String name) {
    return '$name から';
  }

  @override
  String get cancel => 'キャンセル';

  @override
  String get close => '閉じる';

  @override
  String get save => '保存';

  @override
  String get retry => '再試行';

  @override
  String get continueAction => '続ける';

  @override
  String get justNow => 'たった今';

  @override
  String secondsAgo(int n) {
    return '$n 秒前';
  }

  @override
  String minutesAgo(int n) {
    return '$n 分前';
  }

  @override
  String hoursAgo(int n) {
    return '$n 時間前';
  }

  @override
  String daysAgo(int n) {
    return '$n 日前';
  }

  @override
  String monthDay(int month, int day) {
    return '$month月$day日';
  }

  @override
  String durationHours(int n) {
    return '$n 時間';
  }

  @override
  String durationDays(int n) {
    return '$n 日';
  }

  @override
  String get historyEmptySubtitle => 'どのデバイスでコピーしても、ここに表示されます';

  @override
  String historyCount(int count) {
    return '$count 件';
  }

  @override
  String historyCountKept(int count, String ttl) {
    return '$count 件、$ttl保存';
  }

  @override
  String get filterAll => 'すべて';

  @override
  String get filterText => 'テキスト';

  @override
  String get filterImage => '画像';

  @override
  String get filterFile => 'ファイル';

  @override
  String get historyClearTooltip => '履歴を消去';

  @override
  String get historyEmptyTitle => 'まだ履歴がありません';

  @override
  String get historyEmptyFilteredTitle => 'この種類の履歴はありません';

  @override
  String get historyEmptyBody => 'ペアリングしたどのデバイスでテキスト、画像、ファイルをコピーしても、ここに表示されます。';

  @override
  String get clearConfirmTitle => '履歴をすべて消去しますか？';

  @override
  String clearConfirmBody(String device) {
    return '同期した履歴とキャッシュしたファイルが$deviceから削除されます。元のファイルには影響しません。';
  }

  @override
  String get clear => '消去';

  @override
  String get historyCleared => '履歴を消去しました';

  @override
  String get permDeniedTitle => 'PasteFerry はクリップボードを読み取れません';

  @override
  String get permAskTitle => 'PasteFerry にクリップボードの読み取りを許可してください';

  @override
  String get permDeniedBody =>
      'この Mac でコピーした内容は同期されません。システム設定で PasteFerry Daemon を「許可」にしてください。';

  @override
  String get permAskBody =>
      '許可しないと macOS が毎回確認します。システム設定で PasteFerry Daemon を「許可」にしてください。';

  @override
  String get openSystemSettings => 'システム設定を開く';

  @override
  String get fetch => 'このデバイスに取り込む';

  @override
  String get fetchStarted => '取り込みを開始しました';

  @override
  String get transferFailed => '転送に失敗しました';

  @override
  String get expired => '期限切れ';

  @override
  String get putOnClipboard => 'クリップボードに入れる';

  @override
  String get putOnClipboardDone => 'クリップボードに入れました';

  @override
  String get deleteRecord => 'この履歴を削除';

  @override
  String get deleted => '削除しました';

  @override
  String get emptyRecord => '（空）';

  @override
  String get imageRecord => '画像';

  @override
  String itemsSummary(String name, int others, int count) {
    return '$name ほか $others 件';
  }

  @override
  String get searchHint => 'テキストとファイル名を検索';

  @override
  String get searchClear => '検索をクリア';

  @override
  String searchResultCount(int count) {
    return '$count 件見つかりました';
  }

  @override
  String searchEmptyTitle(String query) {
    return '「$query」は見つかりません';
  }

  @override
  String get searchEmptyBody => 'コピーしたテキストとファイル名から検索します。';

  @override
  String get imagePreviewTitle => '画像プレビュー';

  @override
  String get closePreview => 'プレビューを閉じる';

  @override
  String get previousImage => '前の画像';

  @override
  String get nextImage => '次の画像';

  @override
  String get previewHint => 'ホイールやトラックパッドで拡大縮小、ドラッグで移動';

  @override
  String get fitWindow => 'ウィンドウに合わせる';

  @override
  String get previewDeleted => 'この画像は履歴から削除されました';

  @override
  String get previewExpired => 'この画像は期限切れのためプレビューできません';

  @override
  String get previewDownloading => '画像をダウンロード中…';

  @override
  String get previewDownloadFailed => 'ダウンロードに失敗しました。もう一度お試しください。';

  @override
  String get previewNotDownloaded => 'この画像はまだこのコンピューターにありません';

  @override
  String get downloadAndPreview => 'ダウンロードしてプレビュー';

  @override
  String get previewUnavailable => '現在この画像はプレビューできません';

  @override
  String get imageSemantic => '履歴の画像';

  @override
  String get previewUnreadable => '画像を読み込めません。ファイルが削除されたか破損している可能性があります。';

  @override
  String get devicesSubtitleEmpty => 'ペアリングすると、デバイス間でクリップボードが同期されます';

  @override
  String devicesSubtitle(int paired, int online) {
    return 'ペアリング済み $paired 台、オンライン $online 台';
  }

  @override
  String get addDevice => 'デバイスを追加';

  @override
  String thisDeviceFootnote(String both) {
    return 'ペアリング時、$bothに同じ 2 行のセキュリティ指紋（このデバイスと相手のもの）が表示されます。1 文字ずつ一致を確かめてから確定してください。';
  }

  @override
  String get noDevicesTitle => 'ペアリングしたデバイスはまだありません';

  @override
  String get noDevicesBody =>
      'もう一台のコンピューターで PasteFerry を開き、6 桁のペアリングコードでつなげます。同じネットワーク内なら、両方で「近くを検索」を開いて相手をクリックすることもできます。';

  @override
  String get pairedDevices => 'ペアリング済みのデバイス';

  @override
  String get pairedFootnote =>
      '直接接続：データは 2 台のデバイス間で直接やり取りされます。中継：直接つながらないネットワークでは、あなたのサーバーを経由します。内容はエンドツーエンドで暗号化されたままです。';

  @override
  String get offline => 'オフライン';

  @override
  String get direct => '直接接続';

  @override
  String get relay => '中継';

  @override
  String get online => 'オンライン';

  @override
  String get moreActions => 'その他の操作';

  @override
  String get copyFingerprint => '指紋をコピー';

  @override
  String get unpairEllipsis => 'ペアリングを解除…';

  @override
  String get fingerprintCopied => '指紋をコピーしました';

  @override
  String unpairTitle(String name) {
    return '「$name」とのペアリングを解除しますか？';
  }

  @override
  String get unpairBody =>
      '2 台のデバイスは同期しなくなります。再開するには、もう一度ペアリングして指紋を確かめる必要があります。';

  @override
  String get unpair => 'ペアリングを解除';

  @override
  String get unpaired => 'ペアリングを解除しました';

  @override
  String get enterAddressEllipsis => 'アドレスを入力…';

  @override
  String setAddressTitle(String name) {
    return '「$name」のアドレス';
  }

  @override
  String setAddressBody(String name) {
    return '2 台が別のサブネットにある場合や、どちらもアドレスが変わった場合に使います。「$name」の「デバイスを追加 → コードを表示」に表示されるアドレスを入力してください。';
  }

  @override
  String get addressSaved => 'そのアドレスで探しています';

  @override
  String get lanBlockedTitle => 'macOS がローカルネットワークへのアクセスをブロックしています';

  @override
  String lanBlockedBody(String device) {
    return '「システム設定 → プライバシーとセキュリティ → ローカルネットワーク」で PasteFerry Daemon をオンにしてください。オンにするまで、$deviceは同じネットワークのデバイスに接続できません。';
  }

  @override
  String get searchNearby => '近くを検索';

  @override
  String get nearbyTitle => '近くのデバイス';

  @override
  String get nearbySearching => '近くのデバイスを探しています…';

  @override
  String get nearbyHint => 'もう一台でも「近くを検索」を開くと、ここに表示されます。クリックしてペアリングします。';

  @override
  String nearbyVisible(String device) {
    return 'このウィンドウを開いている間は、同じネットワークのデバイスから$deviceの名前が見えます。';
  }

  @override
  String get nearbyPaired => 'ペアリング済み';

  @override
  String nearbyRequesting(String name) {
    return '「$name」に接続しています…';
  }

  @override
  String nearbyNoReply(String name) {
    return '「$name」から応答がありません。相手で「近くを検索」が開いたままか確認してください。';
  }

  @override
  String get nearbyLanOffTitle => 'ローカルネットワークの検出がオフです';

  @override
  String get nearbyLanOffBody =>
      '近くのデバイスを表示するには「設定」で「このネットワークでデバイスを探す」をオンにするか、コードでペアリングしてください。';

  @override
  String get generateCode => 'コードを表示';

  @override
  String get enterCode => 'コードを入力';

  @override
  String get codeFailedTitle => 'ペアリングコードを作成できませんでした';

  @override
  String get enterOnOther => 'もう一台のデバイスでこのコードを入力してください';

  @override
  String get codeCopied => 'コードをコピーしました';

  @override
  String get codeHint => 'クリックでコピー。有効期間は 5 分です。';

  @override
  String afterEnterFootnote(String both) {
    return '入力されると、$bothに同じ 2 行のセキュリティ指紋が表示されます。一致を確かめてから確定してください。';
  }

  @override
  String get codeInvalid => 'ペアリングコードは 6 桁の英数字です';

  @override
  String get localAddress => 'アドレス';

  @override
  String get addressHint => '別のサブネットにある場合は、もう一台でこのアドレスも入力します';

  @override
  String get addressCopied => 'アドレスをコピーしました';

  @override
  String get differentSubnet => '別のサブネットにありますか？';

  @override
  String get peerAddress => 'コードの下に表示されたアドレス';

  @override
  String get pairNotFound =>
      'そのデバイスが見つかりません。別のサブネットにある場合は、コードの下に表示されたアドレスを入力してください。';

  @override
  String get pairNotFoundAddress =>
      'そのアドレスから応答がありません。アドレスと、相手がまだコードを表示しているかを確認してください。相手のファイアウォールがブロックしている可能性もあります。';

  @override
  String get codeWrong => 'コードが違うか、有効期限が切れています';

  @override
  String get codeExpired => 'このコードは無効になりました。もう一台で新しいコードを表示してください。';

  @override
  String get addressInvalid => 'アドレスの形式が正しくありません。192.168.1.20 のように入力してください。';

  @override
  String get pairNoChannel => '先に「設定」でローカルネットワークの検出をオンにするか、サーバーを設定してください。';

  @override
  String get enterCodeShown => 'もう一台のデバイスに表示されたコードを入力してください';

  @override
  String get pairingDone => 'ペアリングしました';

  @override
  String get pairingRejected => 'ペアリングを拒否しました';

  @override
  String get verifyTitle => 'セキュリティ指紋を確認';

  @override
  String pairingWith(String name) {
    return '「$name」とペアリング中';
  }

  @override
  String verifyBody(String both) {
    return '$bothに表示された 2 行は完全に同じはずです。1 文字でも違う場合は、第三者が通信を改ざんしている可能性があるため、拒否してください。';
  }

  @override
  String get rejectMismatch => '不一致、拒否';

  @override
  String get confirmMatch => '一致、ペアリング';

  @override
  String get server => 'サーバー（任意）';

  @override
  String get serverFootnote =>
      '家と会社のように、別々のネットワーク間で同期するときだけ必要です。同じネットワーク内のデバイスはサーバーなしで直接つながります。サーバーはデバイス同士を引き合わせるだけで、やり取りする内容は見えません。';

  @override
  String get signalingUrl => 'シグナリングサーバー';

  @override
  String get signalingUrlDesc =>
      '「保存して接続」をクリックするか Enter を押します。空欄にするとこのネットワーク内だけで同期します';

  @override
  String get signalingUrlPlaceholder => 'ws://サーバーアドレス:8787/signal';

  @override
  String get publicStun => '公開サーバーで出口を調べる';

  @override
  String get publicStunDesc =>
      '複数の回線を持つ会社などのネットワークで、直接つながる経路を多く見つけます。公開サーバーに見えるのはグローバル IP だけです。';

  @override
  String get connectionStatus => '接続状態';

  @override
  String get connectedDesc => 'シグナリングサーバーに接続しています';

  @override
  String get disconnectedDesc => 'サーバーに接続できません。アドレスとサーバーの稼働を確認してください。';

  @override
  String get connected => '接続済み';

  @override
  String get notConnected => '未接続';

  @override
  String get notSet => '未設定';

  @override
  String get notSetDesc => '同じネットワークのデバイスとだけ同期します';

  @override
  String get serverCleared => 'サーバーを削除しました。このネットワーク内だけで同期します';

  @override
  String get localNetwork => 'ローカルネットワーク';

  @override
  String get lanFootnote =>
      'UDP ポート 8788 を使います。2 台が別のサブネットにある場合は、ペアリング時に相手のアドレスを入力します。';

  @override
  String get lanDiscovery => 'このネットワークでデバイスを探す';

  @override
  String get lanDiscoveryDesc =>
      'サーバーなしでペアリングと同期ができます。ネットワーク上の他の人には、毎日変わる匿名の識別子しか見えません。デバイス名が見えるのは「近くを検索」を開いている間だけです';

  @override
  String get macLocalNetwork => 'ローカルネットワークへのアクセス';

  @override
  String get macLocalNetworkDesc =>
      'PasteFerry がネットワーク上のデバイスに接続してよいか、macOS が一度だけ確認します。拒否した場合は、そこで PasteFerry Daemon をオンにしてください';

  @override
  String get allowLan => 'Windows ファイアウォール';

  @override
  String get allowLanDesc =>
      '他のデバイスからこのデバイスに接続できないとき（別のサブネットからのペアリングなど）に許可します。管理者の承認を 1 回求めます';

  @override
  String get allow => '許可';

  @override
  String get allowLanDone => '許可しました。他のデバイスからこのデバイスに接続できます。';

  @override
  String get deviceName => 'デバイス名';

  @override
  String get deviceNameDesc => '他のデバイスの一覧に表示されます';

  @override
  String get language => '言語';

  @override
  String get languageSystem => 'システムに合わせる';

  @override
  String get sync => '同期';

  @override
  String get syncFootnote =>
      '上限より小さい内容はコピーした時点で他のデバイスに送られます。大きいものは記録だけを同期し、必要なときに「履歴」から取り込むので、帯域やディスクを無駄に使いません。';

  @override
  String get autoSyncLimit => '自動同期の上限';

  @override
  String get autoSyncLimitDesc => 'これより大きいファイルは手動で取り込みます';

  @override
  String get autoApply => '受け取った内容をクリップボードに入れる';

  @override
  String get autoApplyDesc => 'オフの場合は「履歴」から手動でクリップボードに入れます';

  @override
  String get whatToSync => '同期する内容';

  @override
  String get syncFiles => 'ファイルとフォルダ';

  @override
  String get syncText => 'プレーンテキスト';

  @override
  String get syncRich => '書式付きテキスト';

  @override
  String get syncImages => '画像とスクリーンショット';

  @override
  String get storage => 'ストレージ';

  @override
  String get storageFootnote => '期限切れの履歴とキャッシュは自動で削除されます。';

  @override
  String get keepHistory => '履歴の保存期間';

  @override
  String get keepCache => 'キャッシュの保存期間';

  @override
  String get keepCacheDesc => '履歴の保存期間を超えません';

  @override
  String get cacheUsage => 'キャッシュ容量';

  @override
  String get backgroundSync => 'バックグラウンド同期';

  @override
  String get backgroundFootnote =>
      'オフにすると同期せず、ログイン時にも起動しません。履歴とペアリングは残り、アプリを開けば再び有効にできます。';

  @override
  String get backgroundService => 'バックグラウンドサービス';

  @override
  String get backgroundServiceDesc => 'ログイン時に起動し、ウィンドウを閉じても同期を続けます';

  @override
  String get restart => '再起動';

  @override
  String get restarted => 'バックグラウンドサービスを再起動しました';

  @override
  String get disableEllipsis => 'オフにする…';

  @override
  String get about => 'このアプリについて';

  @override
  String get aboutDesc => 'オープンソース（AGPL-3.0）';

  @override
  String get homepage => 'ウェブサイト';

  @override
  String get feedback => '問題を報告';

  @override
  String get disableTitle => 'バックグラウンド同期をオフにしますか？';

  @override
  String disableBody(String device) {
    return '$deviceは同期を停止し、ログイン時にも起動しなくなります。履歴とペアリングは残ります。';
  }

  @override
  String get disable => 'オフにする';

  @override
  String get disabled => 'バックグラウンド同期をオフにしました';

  @override
  String get urlMustBeWs => 'アドレスは ws:// または wss:// で始めてください';

  @override
  String get serverUpdated => 'サーバーを更新しました。再接続しています…';

  @override
  String get deviceNameUpdated => 'デバイス名を更新しました';

  @override
  String get diagIntro =>
      'このレポートには、バージョン、システム、設定、デバイスの状態、最近のログが含まれます。デバイス名、サーバーのアドレス、グローバル IP アドレス、ユーザーフォルダ、ファイル名は置き換え済みで、コピーした内容はログに書かれません。GitHub の Issue は公開されるので、投稿する前に目を通してください。';

  @override
  String get diagVerboseTitle => '詳細ログ';

  @override
  String get diagVerboseDesc =>
      '接続の問題を調べるときに使います。オンにして問題をもう一度起こしてから、ここでレポートを作り直してください。30 分後に自動でオフになります。';

  @override
  String get diagVerboseOn => '30 分間オンにする';

  @override
  String diagVerboseUntil(String time) {
    return '$time までオン';
  }

  @override
  String get diagVerboseOff => 'オフにする';

  @override
  String get diagRefresh => '作り直す';

  @override
  String get diagSave => 'ファイルに保存';

  @override
  String diagSaved(String name) {
    return '「ダウンロード」フォルダに保存しました：$name';
  }

  @override
  String get diagCopy => 'レポートをコピー';

  @override
  String get diagCopied => 'レポートをコピーしました';

  @override
  String get diagNewIssue => 'GitHub で Issue を作成';

  @override
  String get diagIssueHint =>
      'レポートをコピーしました。Issue の「Diagnostics report」欄に貼り付けてください';

  @override
  String diagFailed(String detail) {
    return 'レポートを作成できませんでした：$detail';
  }

  @override
  String get svcExtractTitle => '先に PasteFerry をインストールまたは展開してください';

  @override
  String get svcExtractBody =>
      'PasteFerry は zip から直接開かれていて、閉じると一時フォルダごと消えます。インストーラーを実行するか、フォルダ全体を展開してから開くと、バックグラウンド同期を有効にできます。';

  @override
  String get svcMoveTitle => '先に PasteFerry を「アプリケーション」に移動してください';

  @override
  String get svcMoveBody =>
      'PasteFerry はディスクイメージから直接開かれています。「アプリケーション」にドラッグしてディスクイメージを取り出し、「アプリケーション」から開くと、バックグラウンド同期を有効にできます。';

  @override
  String get svcStartTitle => 'PasteFerry を始める';

  @override
  String get svcStartBody =>
      'PasteFerry は小さなバックグラウンドプロセスを動かし、クリップボードの変化を検知して他のデバイスとつながり続けます。有効にするとログイン時に起動し、このウィンドウを閉じても同期を続けます。';

  @override
  String get svcEnable => 'バックグラウンド同期を有効にする';

  @override
  String get svcStaleTitle => 'バックグラウンド同期をもう一度有効にする';

  @override
  String get svcStaleBodyWin =>
      'PasteFerry の場所が変わりましたが、自動起動は以前の場所を指しています。もう一度有効にすると、このコピーを使うようになります。履歴とペアリングは残ります。';

  @override
  String get svcStaleBodyMac =>
      'PasteFerry の場所が変わったか、古いスクリプトでインストールされています。もう一度有効にすると、このアプリを使うようになります。履歴とペアリングは残ります。';

  @override
  String get svcReenable => 'もう一度有効にする';

  @override
  String get svcStoppedTitle => 'バックグラウンドサービスが動いていません';

  @override
  String get svcStoppedBody =>
      '終了した直後かもしれません。通常は数秒で再起動されます。この画面のままなら、下のボタンで再起動してください。';

  @override
  String get svcDevBodyWin =>
      '開発ビルドのため pasteferryd.exe がありません。ソースフォルダで scripts\\build-windows.ps1 を実行して完全版をビルドしてください。';

  @override
  String get svcDevBodyLinux =>
      '開発ビルドのため pasteferryd がありません。ソースフォルダで ./scripts/package-linux.sh を実行して完全版をビルドしてください。';

  @override
  String get svcDevBodyMac =>
      '開発ビルドのためバックグラウンドサービスが含まれていません。ソースフォルダで ./scripts/install-macos.sh を実行してインストールしてください。';

  @override
  String svcLogPath(String path) {
    return 'ログ：$path';
  }

  @override
  String get svcCanDisable => 'バックグラウンド同期は「設定」でいつでもオフにできます。';

  @override
  String get officialService => '公式サービス';

  @override
  String get useOfficialService => '公式サービスを使う';

  @override
  String get officialPrivacy =>
      '有効にした場合のみ公式サービスに接続し、必要最小限の接続情報を記録します。まず直接接続を試み、失敗した場合のみ中継を申請します。内容、実際の端末名、元のシステム識別子は送信しません。';

  @override
  String get officialEnableTitle => '公式インターネットサービスを有効にしますか？';

  @override
  String get officialSwitchConfirm => '現在のカスタムサーバーから切り替えます。保存済みのアドレスは保持されます。';

  @override
  String get customServer => 'カスタムサーバー';

  @override
  String get saveAndConnect => '保存して接続';

  @override
  String get officialUnused => '中継は未申請です。残量は未確認です';

  @override
  String get officialUnavailable => '使用量を取得できません。表示データは古い可能性があります';

  @override
  String get officialQuotaTitle => '本日の中継上限に達しました';

  @override
  String get officialQuotaBody =>
      '公式中継は北京時間の午前0時にリセットされるまで停止します。LANと直接接続は引き続き利用できます。失敗した転送は後で再試行できます。';

  @override
  String get officialPeerQuota => '相手の中継残量がありません。直接接続は引き続き利用できます';

  @override
  String get officialPaused => '公式中継は現在利用できません。直接接続は引き続き利用できます';

  @override
  String get officialViewUsage => '使用量を表示';

  @override
  String get officialAcknowledge => '確認しました';

  @override
  String get officialToday => '本日の中継使用量';

  @override
  String get officialUsed => '使用済み';

  @override
  String get officialRemaining => '残量';

  @override
  String get officialDailyLimit => '1日の上限';

  @override
  String get officialUpload => 'アップロード';

  @override
  String get officialDownload => 'ダウンロード';

  @override
  String get officialReset => 'リセット時刻（北京時間）';

  @override
  String get officialUpdated => '最終更新';

  @override
  String get officialDeviceID => '端末 ID';

  @override
  String get officialWarning => '本日の中継上限の90%を使用しました。直接接続ではこの残量は消費しません。';

  @override
  String get officialAccounting =>
      '暗号化の通信量を含む送受信の合計です。公式中継のみ計測し、直接接続の内容や通信量は記録しません。';
}

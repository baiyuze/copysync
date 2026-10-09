// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Japanese (`ja`).
class AppLocalizationsJa extends AppLocalizations {
  AppLocalizationsJa([String locale = 'ja']) : super(locale);

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
  String get linkNoServer => 'サーバーに未接続';

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
  String get permDeniedTitle => 'CopySync はクリップボードを読み取れません';

  @override
  String get permAskTitle => 'CopySync にクリップボードの読み取りを許可してください';

  @override
  String get permDeniedBody =>
      'この Mac でコピーした内容は同期されません。システム設定で CopySync Daemon を「許可」にしてください。';

  @override
  String get permAskBody =>
      '許可しないと macOS が毎回確認します。システム設定で CopySync Daemon を「許可」にしてください。';

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
  String get noDevicesBody => 'もう一台のコンピューターで CopySync を開き、6 桁のペアリングコードでつなげます。';

  @override
  String get noDevicesNeedServer => 'ペアリングの前に、「設定」でシグナリングサーバーに接続してください。';

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
  String get server => 'サーバー';

  @override
  String get serverFootnote => 'サーバーはデバイス同士を見つけ、直接接続を取り持つだけです。やり取りする内容は見えません。';

  @override
  String get signalingUrl => 'シグナリングサーバー';

  @override
  String get signalingUrlDesc => 'Enter で保存し、すぐに再接続します';

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
  String get svcExtractTitle => '先に CopySync をインストールまたは展開してください';

  @override
  String get svcExtractBody =>
      'CopySync は zip から直接開かれていて、閉じると一時フォルダごと消えます。インストーラーを実行するか、フォルダ全体を展開してから開くと、バックグラウンド同期を有効にできます。';

  @override
  String get svcMoveTitle => '先に CopySync を「アプリケーション」に移動してください';

  @override
  String get svcMoveBody =>
      'CopySync はディスクイメージから直接開かれています。「アプリケーション」にドラッグしてディスクイメージを取り出し、「アプリケーション」から開くと、バックグラウンド同期を有効にできます。';

  @override
  String get svcStartTitle => 'CopySync を始める';

  @override
  String get svcStartBody =>
      'CopySync は小さなバックグラウンドプロセスを動かし、クリップボードの変化を検知して他のデバイスとつながり続けます。有効にするとログイン時に起動し、このウィンドウを閉じても同期を続けます。';

  @override
  String get svcEnable => 'バックグラウンド同期を有効にする';

  @override
  String get svcStaleTitle => 'バックグラウンド同期をもう一度有効にする';

  @override
  String get svcStaleBodyWin =>
      'CopySync の場所が変わりましたが、自動起動は以前の場所を指しています。もう一度有効にすると、このコピーを使うようになります。履歴とペアリングは残ります。';

  @override
  String get svcStaleBodyMac =>
      'CopySync の場所が変わったか、古いスクリプトでインストールされています。もう一度有効にすると、このアプリを使うようになります。履歴とペアリングは残ります。';

  @override
  String get svcReenable => 'もう一度有効にする';

  @override
  String get svcStoppedTitle => 'バックグラウンドサービスが動いていません';

  @override
  String get svcStoppedBody =>
      '終了した直後かもしれません。通常は数秒で再起動されます。この画面のままなら、下のボタンで再起動してください。';

  @override
  String get svcDevBodyWin =>
      '開発ビルドのため copysyncd.exe がありません。ソースフォルダで scripts\\build-windows.ps1 を実行して完全版をビルドしてください。';

  @override
  String get svcDevBodyMac =>
      '開発ビルドのためバックグラウンドサービスが含まれていません。ソースフォルダで ./scripts/install-macos.sh を実行してインストールしてください。';

  @override
  String svcLogPath(String path) {
    return 'ログ：$path';
  }

  @override
  String get svcCanDisable => 'バックグラウンド同期は「設定」でいつでもオフにできます。';
}

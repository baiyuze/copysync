// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Chinese (`zh`).
class AppLocalizationsZh extends AppLocalizations {
  AppLocalizationsZh([String locale = 'zh']) : super(locale);

  @override
  String get appName => '隔空贴';

  @override
  String get navHistory => '复制记录';

  @override
  String get navDevices => '设备';

  @override
  String get navSettings => '设置';

  @override
  String get linkConnecting => '正在连接…';

  @override
  String get linkServiceNotRunning => '后台服务未运行';

  @override
  String get linkLanMode => '在局域网中查找设备';

  @override
  String get linkServerUnreachable => '连不上服务器，局域网内照常同步';

  @override
  String get linkNothingOn => '局域网发现已关闭，也没有设置服务器';

  @override
  String get linkLanBlocked => 'macOS 阻止了局域网访问';

  @override
  String linkOnlinePeers(int count) {
    return '$count 台设备在线';
  }

  @override
  String get linkServerConnected => '已连接服务器';

  @override
  String linkConnectionLost(String detail) {
    return '连接中断：$detail';
  }

  @override
  String get linkServiceDisconnected => '后台服务已断开';

  @override
  String get linkServiceDisabled => '后台服务已停用';

  @override
  String get errServiceNotConnected => '后台服务未连接';

  @override
  String get errUnimplemented => '该功能尚未接入';

  @override
  String get errUnavailable => '无法连接后台服务';

  @override
  String errServiceStart(String detail) {
    return '启动后台服务失败：$detail';
  }

  @override
  String errServiceRestart(String detail) {
    return '重启后台服务失败：$detail';
  }

  @override
  String errAutostart(String detail) {
    return '登记开机自启失败：$detail';
  }

  @override
  String get thisMacTitle => '这台 Mac';

  @override
  String get thisPcTitle => '这台电脑';

  @override
  String get thisMacInText => '这台 Mac ';

  @override
  String get thisPcInText => '这台电脑';

  @override
  String get bothMacs => '两台 Mac ';

  @override
  String get bothDevices => '两台设备';

  @override
  String get thisDeviceShort => '本机';

  @override
  String fromDevice(String name) {
    return '来自 $name';
  }

  @override
  String get cancel => '取消';

  @override
  String get close => '关闭';

  @override
  String get save => '保存';

  @override
  String get retry => '重试';

  @override
  String get continueAction => '继续';

  @override
  String get justNow => '刚刚';

  @override
  String secondsAgo(int n) {
    return '$n 秒前';
  }

  @override
  String minutesAgo(int n) {
    return '$n 分钟前';
  }

  @override
  String hoursAgo(int n) {
    return '$n 小时前';
  }

  @override
  String daysAgo(int n) {
    return '$n 天前';
  }

  @override
  String monthDay(int month, int day) {
    return '$month月$day日';
  }

  @override
  String durationHours(int n) {
    return '$n 小时';
  }

  @override
  String durationDays(int n) {
    return '$n 天';
  }

  @override
  String get historyEmptySubtitle => '在任意一台设备上复制，内容会出现在这里';

  @override
  String historyCount(int count) {
    return '共 $count 条';
  }

  @override
  String historyCountKept(int count, String ttl) {
    return '共 $count 条，保留 $ttl';
  }

  @override
  String get filterAll => '全部';

  @override
  String get filterText => '文本';

  @override
  String get filterImage => '图片';

  @override
  String get filterFile => '文件';

  @override
  String get historyClearTooltip => '清空记录';

  @override
  String get historyEmptyTitle => '还没有记录';

  @override
  String get historyEmptyFilteredTitle => '没有这一类的记录';

  @override
  String get historyEmptyBody => '在任意一台已配对的设备上复制文本、图片或文件，记录会出现在这里。';

  @override
  String get clearConfirmTitle => '清空全部记录？';

  @override
  String clearConfirmBody(String device) {
    return '所有设备上同步来的记录和已缓存的文件都会从$device上删除。原始文件不受影响。';
  }

  @override
  String get clear => '清空';

  @override
  String get historyCleared => '记录已清空';

  @override
  String get permDeniedTitle => '已禁止隔空贴读取剪贴板';

  @override
  String get permAskTitle => '允许隔空贴读取剪贴板';

  @override
  String get permDeniedBody =>
      '在这台 Mac 上复制的内容不会同步出去。到系统设置里把 PasteFerry Daemon 改为「允许」。';

  @override
  String get permAskBody =>
      '否则 macOS 每次都会弹窗询问。到系统设置里把 PasteFerry Daemon 改为「允许」。';

  @override
  String get openSystemSettings => '打开系统设置';

  @override
  String get fetch => '接收';

  @override
  String get fetchStarted => '正在接收，可在复制记录查看进度';

  @override
  String get transferFailed => '传输失败';

  @override
  String get expired => '已过期';

  @override
  String get putOnClipboard => '放入剪贴板';

  @override
  String get putOnClipboardDone => '已放入剪贴板';

  @override
  String get deleteRecord => '删除这条记录';

  @override
  String get deleted => '已删除';

  @override
  String get emptyRecord => '（空）';

  @override
  String get imageRecord => '图片';

  @override
  String itemsSummary(String name, int others, int count) {
    return '$name 等 $count 项';
  }

  @override
  String get searchHint => '搜索文字和文件名';

  @override
  String get searchClear => '清除搜索';

  @override
  String searchResultCount(int count) {
    return '找到 $count 条';
  }

  @override
  String searchEmptyTitle(String query) {
    return '没有找到「$query」';
  }

  @override
  String get searchEmptyBody => '搜索的是复制过的文字和文件名。';

  @override
  String get imagePreviewTitle => '图片预览';

  @override
  String get closePreview => '关闭预览';

  @override
  String get previousImage => '上一张';

  @override
  String get nextImage => '下一张';

  @override
  String get previewHint => '滚轮或触控板缩放，拖动查看';

  @override
  String get fitWindow => '适应窗口';

  @override
  String get previewDeleted => '这条图片记录已被删除';

  @override
  String get previewExpired => '图片已过期，无法预览';

  @override
  String get previewDownloading => '正在下载图片…';

  @override
  String get previewDownloadFailed => '图片下载失败，请重试';

  @override
  String get previewNotDownloaded => '图片尚未下载到这台电脑';

  @override
  String get downloadAndPreview => '下载并预览';

  @override
  String get previewUnavailable => '图片暂时无法预览';

  @override
  String get imageSemantic => '复制记录中的图片';

  @override
  String get previewUnreadable => '无法读取图片，文件可能已被清理或损坏';

  @override
  String get devicesSubtitleEmpty => '配对之后，设备之间就会同步剪贴板';

  @override
  String devicesSubtitle(int paired, int online) {
    return '已配对 $paired 台，$online 台在线';
  }

  @override
  String get addDevice => '添加设备';

  @override
  String thisDeviceFootnote(String both) {
    return '配对时，$both会显示同样的两行安全指纹：这台的和对方的。逐字核对一致才能确认。';
  }

  @override
  String get noDevicesTitle => '还没有配对的设备';

  @override
  String get noDevicesBody =>
      '在另一台电脑上打开隔空贴，用 6 位配对码把两台连起来。两台在同一个网络里时，也可以在两台上都打开「搜索附近设备」，点一下对方。';

  @override
  String get pairedDevices => '已配对的设备';

  @override
  String get pairedFootnote =>
      '直连：数据在两台设备之间直接传输。中转：网络环境打不通直连时，经你的服务器转发，内容依然是端到端加密的。';

  @override
  String get offline => '离线';

  @override
  String get direct => '直连';

  @override
  String get relay => '中转';

  @override
  String get online => '在线';

  @override
  String get moreActions => '更多操作';

  @override
  String get copyFingerprint => '拷贝安全指纹';

  @override
  String get unpairEllipsis => '解除配对…';

  @override
  String get fingerprintCopied => '安全指纹已拷贝';

  @override
  String unpairTitle(String name) {
    return '解除与「$name」的配对？';
  }

  @override
  String get unpairBody => '两台设备将不再同步。以后想恢复，需要重新配对并核对指纹。';

  @override
  String get unpair => '解除配对';

  @override
  String get unpaired => '已解除配对';

  @override
  String get enterAddressEllipsis => '输入它的地址…';

  @override
  String setAddressTitle(String name) {
    return '「$name」的地址';
  }

  @override
  String setAddressBody(String name) {
    return '两台不在同一网段、或者都换了地址时用。输入「$name」上「添加设备 → 生成配对码」下面显示的地址。';
  }

  @override
  String get addressSaved => '已记下，正在那个地址上找它';

  @override
  String get lanBlockedTitle => 'macOS 阻止了局域网访问';

  @override
  String lanBlockedBody(String device) {
    return '在「系统设置 → 隐私与安全性 → 本地网络」里打开 PasteFerry Daemon。打开之前，$device连不上同一网络里的设备。';
  }

  @override
  String get searchNearby => '搜索附近设备';

  @override
  String get nearbyTitle => '附近的设备';

  @override
  String get nearbySearching => '正在查找附近的设备…';

  @override
  String get nearbyHint => '在另一台电脑上也打开「搜索附近设备」，它就会出现在这里，点一下即可配对。';

  @override
  String nearbyVisible(String device) {
    return '窗口打开期间，同一网络里的设备能看到$device的名字。';
  }

  @override
  String get nearbyPaired => '已配对';

  @override
  String nearbyRequesting(String name) {
    return '正在联系「$name」…';
  }

  @override
  String nearbyNoReply(String name) {
    return '「$name」没有回应。确认它还开着「搜索附近设备」。';
  }

  @override
  String get nearbyLanOffTitle => '局域网发现已关闭';

  @override
  String get nearbyLanOffBody => '在「设置」里打开「在局域网中查找设备」才能看到附近的设备，或者改用配对码。';

  @override
  String get generateCode => '生成配对码';

  @override
  String get enterCode => '输入配对码';

  @override
  String get codeFailedTitle => '没能生成配对码';

  @override
  String get enterOnOther => '在另一台设备上输入这个配对码';

  @override
  String get codeCopied => '配对码已拷贝';

  @override
  String get codeHint => '点击可拷贝，5 分钟内有效';

  @override
  String afterEnterFootnote(String both) {
    return '对方输入后，$both会显示同样的两行安全指纹，核对一致再确认。';
  }

  @override
  String get codeInvalid => '配对码是 6 位字母或数字';

  @override
  String get localAddress => '本机地址';

  @override
  String get addressHint => '不在同一网段时，另一台还要输入这个地址';

  @override
  String get addressCopied => '地址已拷贝';

  @override
  String get differentSubnet => '不在同一网段？';

  @override
  String get peerAddress => '对方配对码下面显示的地址';

  @override
  String get pairNotFound => '没找到那台设备。不在同一网段时，请输入它配对码下面显示的地址。';

  @override
  String get pairNotFoundAddress =>
      '那个地址上没有回应。检查地址是否正确、对方是否还停在配对码页面；对方的防火墙也可能拦住了。';

  @override
  String get codeWrong => '配对码不对，或者已经失效';

  @override
  String get codeExpired => '配对码已失效，请在另一台上重新生成';

  @override
  String get addressInvalid => '地址格式不对，应当像 192.168.1.20 这样';

  @override
  String get pairNoChannel => '请先在「设置」里开启局域网发现，或者填写服务器地址。';

  @override
  String get enterCodeShown => '输入另一台设备上显示的配对码';

  @override
  String get pairingDone => '配对完成';

  @override
  String get pairingRejected => '已拒绝配对';

  @override
  String get verifyTitle => '核对安全指纹';

  @override
  String pairingWith(String name) {
    return '正在与「$name」配对';
  }

  @override
  String verifyBody(String both) {
    return '$both上显示的这两行应当完全相同。有任何一个字符不同，说明连接可能被第三方篡改，请拒绝。';
  }

  @override
  String get rejectMismatch => '不一致，拒绝';

  @override
  String get confirmMatch => '一致，确认配对';

  @override
  String get server => '服务器（可选）';

  @override
  String get signalingUrl => '信令服务器地址';

  @override
  String get signalingUrlDesc => '留空仅在局域网同步';

  @override
  String get signalingUrlPlaceholder => 'ws://服务器地址:8787/signal';

  @override
  String get publicStun => '增强公网直连';

  @override
  String get publicStunDesc => '公共服务器可见本机公网 IP';

  @override
  String get connectionStatus => '连接状态';

  @override
  String get disconnectedDesc => '请检查地址与服务器状态';

  @override
  String get connected => '已连接';

  @override
  String get notConnected => '未连接';

  @override
  String get notSet => '未设置';

  @override
  String get serverCleared => '已清除服务器地址，只在局域网里同步';

  @override
  String get localNetwork => '局域网';

  @override
  String get lanDiscovery => '在局域网中查找设备';

  @override
  String get lanDiscoveryDesc => '匿名发现，仅搜索附近设备时显示名称';

  @override
  String get macLocalNetwork => '本地网络权限';

  @override
  String get macLocalNetworkDesc => '在系统设置中允许 PasteFerry Daemon';

  @override
  String get allowLan => 'Windows 防火墙';

  @override
  String get allowLanDesc => '连接受阻时放行，需管理员权限';

  @override
  String get allow => '放行';

  @override
  String get allowLanDone => '已放行，其他设备现在可以连到这台了';

  @override
  String get deviceName => '设备名称';

  @override
  String get language => '语言';

  @override
  String get languageSystem => '跟随系统';

  @override
  String get sync => '同步';

  @override
  String get autoSyncLimit => '自动同步上限';

  @override
  String get autoSyncLimitDesc => '超出上限后，在复制记录中手动接收';

  @override
  String get autoApply => '收到后直接写入剪贴板';

  @override
  String get whatToSync => '同步哪些内容';

  @override
  String get syncFiles => '文件与文件夹';

  @override
  String get syncText => '纯文本';

  @override
  String get syncRich => '带格式的文本';

  @override
  String get syncImages => '图片与截图';

  @override
  String get storage => '存储';

  @override
  String get storageFootnote => '到期自动清理记录与缓存。';

  @override
  String get keepHistory => '记录保留';

  @override
  String get keepCache => '缓存文件保留';

  @override
  String get keepCacheDesc => '不超过记录保留时长';

  @override
  String get cacheUsage => '缓存占用';

  @override
  String get backgroundSync => '后台同步';

  @override
  String get backgroundService => '后台服务';

  @override
  String get backgroundServiceDesc => '登录自启，关闭窗口后继续同步';

  @override
  String get restart => '重新启动';

  @override
  String get restarted => '后台服务已重新启动';

  @override
  String get disableEllipsis => '停用…';

  @override
  String get about => '关于';

  @override
  String get aboutDesc => '开源软件，AGPL-3.0 许可';

  @override
  String get homepage => '项目主页';

  @override
  String get feedback => '反馈问题';

  @override
  String get disableTitle => '停用后台同步？';

  @override
  String disableBody(String device) {
    return '$device将停止同步，也不再开机启动。历史记录和配对关系会保留。';
  }

  @override
  String get disable => '停用';

  @override
  String get disabled => '后台同步已停用';

  @override
  String get urlMustBeWs => '地址需以 ws:// 或 wss:// 开头';

  @override
  String get serverUpdated => '服务器地址已更新，正在重连';

  @override
  String get deviceNameUpdated => '设备名称已更新';

  @override
  String get diagIntro =>
      '报告里有版本、系统、设置、设备状态和最近的日志。设备名、服务器地址、公网 IP、用户目录和文件名都已替换，复制的内容从不写进日志。GitHub 上的问题是公开的，贴出去之前请看一遍。';

  @override
  String get diagVerboseTitle => '详细日志';

  @override
  String get diagVerboseDesc => '连接方面的问题常要看它：打开后重现一次问题，再回来生成报告。30 分钟后自动关闭。';

  @override
  String get diagVerboseOn => '打开 30 分钟';

  @override
  String diagVerboseUntil(String time) {
    return '已打开，$time 自动关闭';
  }

  @override
  String get diagVerboseOff => '关闭';

  @override
  String get diagRefresh => '重新生成';

  @override
  String get diagSave => '保存为文件';

  @override
  String diagSaved(String name) {
    return '已保存到「下载」文件夹：$name';
  }

  @override
  String get diagCopy => '拷贝报告';

  @override
  String get diagCopied => '报告已拷贝';

  @override
  String get diagNewIssue => '在 GitHub 上提问题';

  @override
  String get diagIssueHint => '报告已拷贝，粘贴到问题里的「诊断报告」一栏';

  @override
  String diagFailed(String detail) {
    return '生成报告失败：$detail';
  }

  @override
  String get svcExtractTitle => '先解压或安装隔空贴';

  @override
  String get svcExtractBody =>
      '隔空贴现在是直接从压缩包里打开的，关掉之后它所在的临时文件夹就会被删除。运行安装程序，或者把整个文件夹解压出来再打开，就可以启用后台同步。';

  @override
  String get svcMoveTitle => '先把隔空贴移到「应用程序」';

  @override
  String get svcMoveBody =>
      '隔空贴现在是直接从安装盘里打开的。把它拖进「应用程序」文件夹，推出安装盘，再从「应用程序」里打开，就可以启用后台同步。';

  @override
  String get svcStartTitle => '开始使用隔空贴';

  @override
  String get svcStartBody =>
      '隔空贴会在后台运行一个小进程，用来感知剪贴板的变化，并与你的其他设备保持连接。启用后它随系统登录自动启动，关掉这个窗口也会照常同步。';

  @override
  String get svcEnable => '启用后台同步';

  @override
  String get svcStaleTitle => '重新启用后台同步';

  @override
  String get svcStaleBodyWin =>
      '隔空贴换了位置，开机自启还指向原来的地方。重新启用一次，让后台进程指向当前这份隔空贴。历史记录与配对关系都会保留。';

  @override
  String get svcStaleBodyMac =>
      '隔空贴换了位置，或者之前用旧版安装脚本装过。重新启用一次，让后台进程指向当前这个 App。历史记录与配对关系都会保留。';

  @override
  String get svcReenable => '重新启用';

  @override
  String get svcStoppedTitle => '后台服务没有运行';

  @override
  String get svcStoppedBody => '它可能刚刚退出，系统通常会在几秒内把它重新拉起。如果一直停在这一页，点下面的按钮重新启动。';

  @override
  String get svcDevBodyWin =>
      '这是开发构建，旁边没有 pasteferryd.exe。在源码目录运行 scripts\\build-windows.ps1 构建完整的版本。';

  @override
  String get svcDevBodyLinux =>
      '这是开发构建，旁边没有 pasteferryd。在源码目录运行 ./scripts/package-linux.sh 构建完整的版本。';

  @override
  String get svcDevBodyMac =>
      '这是开发构建，App 里没有内置后台服务。在源码目录运行 ./scripts/install-macos.sh 安装它。';

  @override
  String svcLogPath(String path) {
    return '运行日志：$path';
  }

  @override
  String get svcCanDisable => '随时可以在「设置」中停用后台同步。';

  @override
  String get officialService => '官方服务';

  @override
  String get useOfficialService => '使用官方服务';

  @override
  String get officialServiceDesc => '跨网络同步，优先直连';

  @override
  String get officialPrivacy =>
      '仅开启后连接官方服务并统计必要接入信息。先尝试直连，失败才申请中转；不上传内容、真实设备名或原始系统标识。';

  @override
  String get officialEnableTitle => '启用官方公网服务？';

  @override
  String get officialSwitchConfirm => '这会切换当前自定义服务器，已保存的地址仍会保留。';

  @override
  String get customServer => '自定义服务器';

  @override
  String get saveAndConnect => '保存并连接';

  @override
  String get officialUnused => '尚未申请中转，余额待核验';

  @override
  String get officialUnavailable => '用量暂不可用，已有数据可能过期';

  @override
  String get officialQuotaTitle => '今日中转额度已用完';

  @override
  String get officialQuotaBody => '官方中转已暂停，北京时间零点重置。局域网和公网直连仍可使用；传输失败时可稍后重试。';

  @override
  String get officialPeerQuota => '对端中转额度不足，直连仍可使用';

  @override
  String get officialPaused => '官方中转暂不可用，直连仍可使用';

  @override
  String get officialViewUsage => '查看用量';

  @override
  String get officialAcknowledge => '知道了';

  @override
  String get officialToday => '今日中转用量';

  @override
  String get officialUsed => '已使用';

  @override
  String get officialRemaining => '剩余额度';

  @override
  String get officialDailyLimit => '每日额度';

  @override
  String get officialUpload => '上传';

  @override
  String get officialDownload => '下载';

  @override
  String get officialReset => '重置时间（北京时间）';

  @override
  String get officialUpdated => '最后更新';

  @override
  String get officialDeviceID => '设备 ID';

  @override
  String get officialWarning => '今日中转额度已使用 90%，直连不消耗此额度。';

  @override
  String get officialAccounting => '上下行合计，包含加密传输开销；只统计官方中转，不统计直连内容或流量。';

  @override
  String buildDetails(String uiVersion, String serviceVersion, String build) {
    return '界面 $uiVersion · 后台 $serviceVersion · 构建 $build';
  }

  @override
  String get notReceived => '尚未接收，暂不可粘贴';

  @override
  String get receivedPreserved => '已接收，点击放入剪贴板后可粘贴';

  @override
  String get clipboardWriteFailed => '已接收，但未能写入剪贴板';

  @override
  String get receivingStatus => '正在接收…';

  @override
  String get largeFileNotifications => '大文件接收提醒';

  @override
  String get notificationsAllowed => '已允许系统通知，接收进度在复制记录中查看。';

  @override
  String get notificationsNotRequested => '开启系统通知，避免大文件等待接收时没有提示；也可随时从复制记录接收。';

  @override
  String get notificationsDenied => '系统通知已关闭，请在系统设置中开启；待接收内容仍可从复制记录接收。';

  @override
  String get notificationsUnknown => '通知显示和勿扰由系统控制；未看到提醒时，可在复制记录中接收。';

  @override
  String get notificationsUnavailable => '系统通知暂不可用，请从复制记录接收大文件。';

  @override
  String get enableNotifications => '开启提醒';

  @override
  String get notificationSettings => '通知设置';

  @override
  String get notifyPendingTitle => '接收后即可粘贴';

  @override
  String notifyPendingBody(String device, String name, String size) {
    return '$device · $name · $size\n尚未接收，暂不可粘贴。也可从复制记录接收。';
  }

  @override
  String get notifyReceivingTitle => '正在接收…';

  @override
  String notifyReceivingBody(String name) {
    return '$name\n可在复制记录查看接收进度。';
  }

  @override
  String get notifyReadyTitle => '已就绪，可以粘贴';

  @override
  String notifyReadyBody(String name) {
    return '$name\n已接收并放入剪贴板。';
  }

  @override
  String get notifyReceivedTitle => '已接收';

  @override
  String notifyReceivedBody(String name) {
    return '$name\n已保留当前剪贴板，点击放入剪贴板后可粘贴，也可从复制记录操作。';
  }

  @override
  String get notifyFailedTitle => '接收失败';

  @override
  String notifyFailedBody(String name, String error) {
    return '$name\n$error\n可从此通知或复制记录重试。';
  }

  @override
  String get notifyClipboardFailedTitle => '已接收，但未能写入剪贴板';

  @override
  String notifyClipboardFailedBody(String name, String error) {
    return '$name\n$error\n可从此通知或复制记录重新放入剪贴板。';
  }
}

/// The translations for Chinese, as used in Hong Kong, using the Han script (`zh_Hant_HK`).
class AppLocalizationsZhHantHk extends AppLocalizationsZh {
  AppLocalizationsZhHantHk() : super('zh_Hant_HK');

  @override
  String get appName => '隔空貼';

  @override
  String get navHistory => '複製記錄';

  @override
  String get navDevices => '裝置';

  @override
  String get navSettings => '設定';

  @override
  String get linkConnecting => '正在連線…';

  @override
  String get linkServiceNotRunning => '背景服務未在執行';

  @override
  String get linkLanMode => '在區域網絡中尋找裝置';

  @override
  String get linkServerUnreachable => '連不上伺服器，區域網絡內照常同步';

  @override
  String get linkNothingOn => '區域網絡發現已關閉，也沒有設定伺服器';

  @override
  String get linkLanBlocked => 'macOS 阻擋了區域網絡存取';

  @override
  String linkOnlinePeers(int count) {
    return '$count 部裝置在線';
  }

  @override
  String get linkServerConnected => '已連線到伺服器';

  @override
  String linkConnectionLost(String detail) {
    return '連線中斷：$detail';
  }

  @override
  String get linkServiceDisconnected => '背景服務已中斷連線';

  @override
  String get linkServiceDisabled => '背景服務已停用';

  @override
  String get errServiceNotConnected => '未連線到背景服務';

  @override
  String get errUnimplemented => '這項功能仍未提供';

  @override
  String get errUnavailable => '無法連線到背景服務';

  @override
  String errServiceStart(String detail) {
    return '無法啟動背景服務：$detail';
  }

  @override
  String errServiceRestart(String detail) {
    return '無法重新啟動背景服務：$detail';
  }

  @override
  String errAutostart(String detail) {
    return '無法設定開機自動啟動：$detail';
  }

  @override
  String get thisMacTitle => '這部 Mac';

  @override
  String get thisPcTitle => '這部電腦';

  @override
  String get thisMacInText => '這部 Mac ';

  @override
  String get thisPcInText => '這部電腦';

  @override
  String get bothMacs => '兩部 Mac ';

  @override
  String get bothDevices => '兩部裝置';

  @override
  String get thisDeviceShort => '本機';

  @override
  String fromDevice(String name) {
    return '來自 $name';
  }

  @override
  String get cancel => '取消';

  @override
  String get close => '關閉';

  @override
  String get save => '儲存';

  @override
  String get retry => '重試';

  @override
  String get continueAction => '繼續';

  @override
  String get justNow => '剛剛';

  @override
  String secondsAgo(int n) {
    return '$n 秒前';
  }

  @override
  String minutesAgo(int n) {
    return '$n 分鐘前';
  }

  @override
  String hoursAgo(int n) {
    return '$n 小時前';
  }

  @override
  String daysAgo(int n) {
    return '$n 天前';
  }

  @override
  String monthDay(int month, int day) {
    return '$month月$day日';
  }

  @override
  String durationHours(int n) {
    return '$n 小時';
  }

  @override
  String durationDays(int n) {
    return '$n 天';
  }

  @override
  String get historyEmptySubtitle => '在任何一部裝置上複製，內容會出現在這裏';

  @override
  String historyCount(int count) {
    return '共 $count 條';
  }

  @override
  String historyCountKept(int count, String ttl) {
    return '共 $count 條，保留 $ttl';
  }

  @override
  String get filterAll => '全部';

  @override
  String get filterText => '文字';

  @override
  String get filterImage => '圖片';

  @override
  String get filterFile => '檔案';

  @override
  String get historyClearTooltip => '清除記錄';

  @override
  String get historyEmptyTitle => '仍未有記錄';

  @override
  String get historyEmptyFilteredTitle => '沒有這一類的記錄';

  @override
  String get historyEmptyBody => '在任何一部已配對的裝置上複製文字、圖片或檔案，記錄會出現在這裏。';

  @override
  String get clearConfirmTitle => '清除全部記錄？';

  @override
  String clearConfirmBody(String device) {
    return '所有裝置同步過來的記錄和已快取的檔案都會從$device上刪除。原始檔案不受影響。';
  }

  @override
  String get clear => '清除';

  @override
  String get historyCleared => '記錄已清除';

  @override
  String get permDeniedTitle => '已禁止隔空貼讀取剪貼簿';

  @override
  String get permAskTitle => '允許隔空貼讀取剪貼簿';

  @override
  String get permDeniedBody =>
      '在這部 Mac 上複製的內容不會同步出去。請到系統設定裏把 PasteFerry Daemon 改為「允許」。';

  @override
  String get permAskBody =>
      '否則 macOS 每次都會彈出視窗詢問。請到系統設定裏把 PasteFerry Daemon 改為「允許」。';

  @override
  String get openSystemSettings => '開啟系統設定';

  @override
  String get fetch => '接收';

  @override
  String get fetchStarted => '正在接收，可在複製記錄查看進度';

  @override
  String get transferFailed => '傳輸失敗';

  @override
  String get expired => '已過期';

  @override
  String get putOnClipboard => '放入剪貼簿';

  @override
  String get putOnClipboardDone => '已放入剪貼簿';

  @override
  String get deleteRecord => '刪除這條記錄';

  @override
  String get deleted => '已刪除';

  @override
  String get emptyRecord => '（空）';

  @override
  String get imageRecord => '圖片';

  @override
  String itemsSummary(String name, int others, int count) {
    return '$name 等 $count 項';
  }

  @override
  String get searchHint => '搜尋文字和檔案名稱';

  @override
  String get searchClear => '清除搜尋';

  @override
  String searchResultCount(int count) {
    return '找到 $count 條';
  }

  @override
  String searchEmptyTitle(String query) {
    return '找不到「$query」';
  }

  @override
  String get searchEmptyBody => '搜尋範圍是複製過的文字和檔案名稱。';

  @override
  String get imagePreviewTitle => '圖片預覽';

  @override
  String get closePreview => '關閉預覽';

  @override
  String get previousImage => '上一張';

  @override
  String get nextImage => '下一張';

  @override
  String get previewHint => '用滾輪或觸控板縮放，拖移查看';

  @override
  String get fitWindow => '符合視窗大小';

  @override
  String get previewDeleted => '這條圖片記錄已被刪除';

  @override
  String get previewExpired => '圖片已過期，無法預覽';

  @override
  String get previewDownloading => '正在下載圖片…';

  @override
  String get previewDownloadFailed => '圖片下載失敗，請重試';

  @override
  String get previewNotDownloaded => '圖片仍未下載到這部電腦';

  @override
  String get downloadAndPreview => '下載並預覽';

  @override
  String get previewUnavailable => '圖片暫時無法預覽';

  @override
  String get imageSemantic => '複製記錄中的圖片';

  @override
  String get previewUnreadable => '無法讀取圖片，檔案可能已被清理或損壞';

  @override
  String get devicesSubtitleEmpty => '配對之後，裝置之間就會同步剪貼簿';

  @override
  String devicesSubtitle(int paired, int online) {
    return '已配對 $paired 部，$online 部在線';
  }

  @override
  String get addDevice => '新增裝置';

  @override
  String thisDeviceFootnote(String both) {
    return '配對時，$both會顯示同樣的兩行安全指紋：這部的和對方的。逐字核對一致才可確認。';
  }

  @override
  String get noDevicesTitle => '仍未有已配對的裝置';

  @override
  String get noDevicesBody =>
      '在另一部電腦上開啟隔空貼，用 6 位配對碼把兩部連起來。兩部在同一個網絡裏時，也可以在兩部上都開啟「搜尋附近裝置」，按一下對方。';

  @override
  String get pairedDevices => '已配對的裝置';

  @override
  String get pairedFootnote =>
      '直連：資料在兩部裝置之間直接傳輸。中轉：網絡環境打不通直連時，經你的伺服器轉發，內容依然是端對端加密的。';

  @override
  String get offline => '離線';

  @override
  String get direct => '直連';

  @override
  String get relay => '中轉';

  @override
  String get online => '在線';

  @override
  String get moreActions => '更多操作';

  @override
  String get copyFingerprint => '拷貝安全指紋';

  @override
  String get unpairEllipsis => '解除配對…';

  @override
  String get fingerprintCopied => '安全指紋已拷貝';

  @override
  String unpairTitle(String name) {
    return '解除與「$name」的配對？';
  }

  @override
  String get unpairBody => '兩部裝置將不再同步。日後想恢復，需要重新配對並核對指紋。';

  @override
  String get unpair => '解除配對';

  @override
  String get unpaired => '已解除配對';

  @override
  String get enterAddressEllipsis => '輸入它的地址…';

  @override
  String setAddressTitle(String name) {
    return '「$name」的地址';
  }

  @override
  String setAddressBody(String name) {
    return '兩部不在同一網段，或者都換了地址時使用。輸入「$name」上「新增裝置 → 產生配對碼」下面顯示的地址。';
  }

  @override
  String get addressSaved => '已記下，正在那個地址上尋找它';

  @override
  String get lanBlockedTitle => 'macOS 阻擋了區域網絡存取';

  @override
  String lanBlockedBody(String device) {
    return '在「系統設定 → 私隱與保安 → 本地網絡」裏開啟 PasteFerry Daemon。開啟之前，$device連不上同一網絡裏的裝置。';
  }

  @override
  String get searchNearby => '搜尋附近裝置';

  @override
  String get nearbyTitle => '附近的裝置';

  @override
  String get nearbySearching => '正在尋找附近的裝置…';

  @override
  String get nearbyHint => '在另一部電腦上也開啟「搜尋附近裝置」，它就會出現在這裏，按一下即可配對。';

  @override
  String nearbyVisible(String device) {
    return '視窗開啟期間，同一網絡裏的裝置能看到$device的名稱。';
  }

  @override
  String get nearbyPaired => '已配對';

  @override
  String nearbyRequesting(String name) {
    return '正在聯絡「$name」…';
  }

  @override
  String nearbyNoReply(String name) {
    return '「$name」沒有回應。請確認它仍開着「搜尋附近裝置」。';
  }

  @override
  String get nearbyLanOffTitle => '區域網絡發現已關閉';

  @override
  String get nearbyLanOffBody => '在「設定」裏開啟「在區域網絡中尋找裝置」才能看到附近的裝置，或者改用配對碼。';

  @override
  String get generateCode => '產生配對碼';

  @override
  String get enterCode => '輸入配對碼';

  @override
  String get codeFailedTitle => '未能產生配對碼';

  @override
  String get enterOnOther => '在另一部裝置上輸入這個配對碼';

  @override
  String get codeCopied => '配對碼已拷貝';

  @override
  String get codeHint => '按一下即可拷貝，5 分鐘內有效';

  @override
  String afterEnterFootnote(String both) {
    return '對方輸入後，$both會顯示同樣的兩行安全指紋，核對一致再確認。';
  }

  @override
  String get codeInvalid => '配對碼是 6 位字母或數字';

  @override
  String get localAddress => '本機地址';

  @override
  String get addressHint => '不在同一網段時，另一部還要輸入這個地址';

  @override
  String get addressCopied => '地址已拷貝';

  @override
  String get differentSubnet => '不在同一網段？';

  @override
  String get peerAddress => '對方配對碼下面顯示的地址';

  @override
  String get pairNotFound => '找不到那部裝置。不在同一網段時，請輸入它配對碼下面顯示的地址。';

  @override
  String get pairNotFoundAddress =>
      '那個地址沒有回應。請檢查地址是否正確、對方是否仍停留在配對碼頁面；對方的防火牆也可能阻擋了連線。';

  @override
  String get codeWrong => '配對碼不正確，或者已經失效';

  @override
  String get codeExpired => '配對碼已失效，請在另一部上重新產生';

  @override
  String get addressInvalid => '地址格式不正確，應該像 192.168.1.20 這樣';

  @override
  String get pairNoChannel => '請先在「設定」裏開啟區域網絡發現，或者填寫伺服器地址。';

  @override
  String get enterCodeShown => '輸入另一部裝置上顯示的配對碼';

  @override
  String get pairingDone => '配對完成';

  @override
  String get pairingRejected => '已拒絕配對';

  @override
  String get verifyTitle => '核對安全指紋';

  @override
  String pairingWith(String name) {
    return '正在與「$name」配對';
  }

  @override
  String verifyBody(String both) {
    return '$both上顯示的這兩行應該完全相同。只要有一個字元不同，就代表連線可能被第三方篡改，請拒絕。';
  }

  @override
  String get rejectMismatch => '不一致，拒絕';

  @override
  String get confirmMatch => '一致，確認配對';

  @override
  String get server => '伺服器（可選）';

  @override
  String get signalingUrl => '信令伺服器地址';

  @override
  String get signalingUrlDesc => '留空只在區域網絡同步';

  @override
  String get signalingUrlPlaceholder => 'ws://伺服器地址:8787/signal';

  @override
  String get publicStun => '提升公網直連成功率';

  @override
  String get publicStunDesc => '公共伺服器可見本機公網 IP';

  @override
  String get connectionStatus => '連線狀態';

  @override
  String get disconnectedDesc => '請檢查地址與伺服器狀態';

  @override
  String get connected => '已連線';

  @override
  String get notConnected => '未連線';

  @override
  String get notSet => '未設定';

  @override
  String get serverCleared => '已清除伺服器地址，只在區域網絡裏同步';

  @override
  String get localNetwork => '區域網絡';

  @override
  String get lanDiscovery => '在區域網絡中尋找裝置';

  @override
  String get lanDiscoveryDesc => '匿名發現，只在搜尋附近裝置時顯示名稱';

  @override
  String get macLocalNetwork => '本地網絡權限';

  @override
  String get macLocalNetworkDesc => '在系統設定中允許 PasteFerry Daemon';

  @override
  String get allowLan => 'Windows 防火牆';

  @override
  String get allowLanDesc => '連線受阻時允許連線，需管理員權限';

  @override
  String get allow => '允許';

  @override
  String get allowLanDone => '已允許，其他裝置現在可以連到這部電腦了';

  @override
  String get deviceName => '裝置名稱';

  @override
  String get language => '語言';

  @override
  String get languageSystem => '跟隨系統';

  @override
  String get sync => '同步';

  @override
  String get autoSyncLimit => '自動同步上限';

  @override
  String get autoSyncLimitDesc => '超出上限後，在複製記錄中手動接收';

  @override
  String get autoApply => '收到後直接寫入剪貼簿';

  @override
  String get whatToSync => '同步哪些內容';

  @override
  String get syncFiles => '檔案與資料夾';

  @override
  String get syncText => '純文字';

  @override
  String get syncRich => '帶格式的文字';

  @override
  String get syncImages => '圖片與截圖';

  @override
  String get storage => '儲存空間';

  @override
  String get storageFootnote => '到期自動清理記錄與快取。';

  @override
  String get keepHistory => '記錄保留';

  @override
  String get keepCache => '快取檔案保留';

  @override
  String get keepCacheDesc => '不超過記錄保留時間';

  @override
  String get cacheUsage => '快取佔用空間';

  @override
  String get backgroundSync => '背景同步';

  @override
  String get backgroundService => '背景服務';

  @override
  String get backgroundServiceDesc => '登入自動啟動，關閉視窗後繼續同步';

  @override
  String get restart => '重新啟動';

  @override
  String get restarted => '背景服務已重新啟動';

  @override
  String get disableEllipsis => '停用…';

  @override
  String get about => '關於';

  @override
  String get aboutDesc => '開源軟件，以 AGPL-3.0 授權';

  @override
  String get homepage => '項目主頁';

  @override
  String get feedback => '回報問題';

  @override
  String get disableTitle => '停用背景同步？';

  @override
  String disableBody(String device) {
    return '$device將停止同步，也不再開機啟動。歷史記錄和配對關係會保留。';
  }

  @override
  String get disable => '停用';

  @override
  String get disabled => '背景同步已停用';

  @override
  String get urlMustBeWs => '地址必須以 ws:// 或 wss:// 開頭';

  @override
  String get serverUpdated => '伺服器地址已更新，正在重新連線';

  @override
  String get deviceNameUpdated => '裝置名稱已更新';

  @override
  String get diagIntro =>
      '報告裏有版本、系統、設定、裝置狀態和最近的日誌。裝置名稱、伺服器地址、公網 IP、用戶資料夾和檔案名稱都已替換，複製的內容從不寫進日誌。GitHub 上的問題是公開的，貼出去之前請看一遍。';

  @override
  String get diagVerboseTitle => '詳細日誌';

  @override
  String get diagVerboseDesc => '連線方面的問題常要看它：開啟後重現一次問題，再回來產生報告。30 分鐘後自動關閉。';

  @override
  String get diagVerboseOn => '開啟 30 分鐘';

  @override
  String diagVerboseUntil(String time) {
    return '已開啟，$time 自動關閉';
  }

  @override
  String get diagVerboseOff => '關閉';

  @override
  String get diagRefresh => '重新產生';

  @override
  String get diagSave => '儲存為檔案';

  @override
  String diagSaved(String name) {
    return '已儲存到「下載」資料夾：$name';
  }

  @override
  String get diagCopy => '拷貝報告';

  @override
  String get diagCopied => '報告已拷貝';

  @override
  String get diagNewIssue => '在 GitHub 上回報';

  @override
  String get diagIssueHint => '報告已拷貝，貼到問題裏的「診斷報告」一欄';

  @override
  String diagFailed(String detail) {
    return '產生報告失敗：$detail';
  }

  @override
  String get svcExtractTitle => '先解壓縮或安裝隔空貼';

  @override
  String get svcExtractBody =>
      '隔空貼現在是直接從壓縮檔裏開啟的，關閉之後它所在的臨時資料夾就會被刪除。執行安裝程式，或者把整個資料夾解壓縮出來再開啟，就可以啟用背景同步。';

  @override
  String get svcMoveTitle => '先把隔空貼移到「應用程式」';

  @override
  String get svcMoveBody =>
      '隔空貼現在是直接從磁碟映像裏開啟的。把它拖進「應用程式」資料夾，推出磁碟映像，再從「應用程式」裏開啟，就可以啟用背景同步。';

  @override
  String get svcStartTitle => '開始使用隔空貼';

  @override
  String get svcStartBody =>
      '隔空貼會在背景執行一個小型程序，用來偵測剪貼簿的變化，並與你的其他裝置保持連線。啟用後它會隨系統登入自動啟動，關閉這個視窗也會照常同步。';

  @override
  String get svcEnable => '啟用背景同步';

  @override
  String get svcStaleTitle => '重新啟用背景同步';

  @override
  String get svcStaleBodyWin =>
      '隔空貼換了位置，開機自動啟動仍指向原來的位置。重新啟用一次，讓背景程序指向目前這份隔空貼。歷史記錄與配對關係都會保留。';

  @override
  String get svcStaleBodyMac =>
      '隔空貼換了位置，或者之前用舊版安裝腳本安裝過。重新啟用一次，讓背景程序指向目前這個 App。歷史記錄與配對關係都會保留。';

  @override
  String get svcReenable => '重新啟用';

  @override
  String get svcStoppedTitle => '背景服務沒有執行';

  @override
  String get svcStoppedBody => '它可能剛剛結束，系統通常會在幾秒內重新啟動它。如果一直停留在這一頁，請按下面的按鈕重新啟動。';

  @override
  String get svcDevBodyWin =>
      '這是開發版本，旁邊沒有 pasteferryd.exe。在原始碼目錄執行 scripts\\build-windows.ps1 建置完整的版本。';

  @override
  String get svcDevBodyLinux =>
      '這是開發版本，旁邊沒有 pasteferryd。在原始碼目錄執行 ./scripts/package-linux.sh 建置完整的版本。';

  @override
  String get svcDevBodyMac =>
      '這是開發版本，App 裏沒有內建背景服務。在原始碼目錄執行 ./scripts/install-macos.sh 安裝它。';

  @override
  String svcLogPath(String path) {
    return '執行日誌：$path';
  }

  @override
  String get svcCanDisable => '隨時可以在「設定」中停用背景同步。';

  @override
  String get officialService => '官方服務';

  @override
  String get useOfficialService => '使用官方服務';

  @override
  String get officialServiceDesc => '跨網絡同步，優先直連';

  @override
  String get officialPrivacy =>
      '只在啟用後連接官方服務並統計必要連線資料。先嘗試直連，失敗才申請中轉；不傳送內容、真實裝置名稱或原始系統識別碼。';

  @override
  String get officialEnableTitle => '啟用官方公網服務？';

  @override
  String get officialSwitchConfirm => '這會切換目前的自訂伺服器，已儲存的地址仍會保留。';

  @override
  String get customServer => '自訂伺服器';

  @override
  String get saveAndConnect => '儲存並連接';

  @override
  String get officialUnused => '尚未申請中轉，餘額待確認';

  @override
  String get officialUnavailable => '用量暫時無法取得，現有資料可能過期';

  @override
  String get officialQuotaTitle => '今日中轉額度已用完';

  @override
  String get officialQuotaBody => '官方中轉已暫停，北京時間午夜重設。區域網絡和公網直連仍可使用；傳輸失敗時可稍後重試。';

  @override
  String get officialPeerQuota => '對方中轉額度不足，直連仍可使用';

  @override
  String get officialPaused => '官方中轉暫時無法使用，直連仍可使用';

  @override
  String get officialViewUsage => '查看用量';

  @override
  String get officialAcknowledge => '知道了';

  @override
  String get officialToday => '今日中轉用量';

  @override
  String get officialUsed => '已使用';

  @override
  String get officialRemaining => '剩餘額度';

  @override
  String get officialDailyLimit => '每日額度';

  @override
  String get officialUpload => '上傳';

  @override
  String get officialDownload => '下載';

  @override
  String get officialReset => '重設時間（北京時間）';

  @override
  String get officialUpdated => '上次更新';

  @override
  String get officialDeviceID => '裝置 ID';

  @override
  String get officialWarning => '今日中轉額度已使用 90%，直連不消耗此額度。';

  @override
  String get officialAccounting => '上下行合計，包含加密傳輸開銷；只統計官方中轉，不統計直連內容或流量。';

  @override
  String buildDetails(String uiVersion, String serviceVersion, String build) {
    return '介面 $uiVersion · 背景服務 $serviceVersion · 版本識別 $build';
  }

  @override
  String get notReceived => '尚未接收，暫不可貼上';

  @override
  String get receivedPreserved => '已接收，按放入剪貼簿後即可貼上';

  @override
  String get clipboardWriteFailed => '已接收，但無法寫入剪貼簿';

  @override
  String get receivingStatus => '正在接收…';

  @override
  String get largeFileNotifications => '大檔案接收提醒';

  @override
  String get notificationsAllowed => '已允許系統通知，接收進度可在複製記錄中查看。';

  @override
  String get notificationsNotRequested => '開啟系統通知，避免大檔案等待接收時沒有提示；亦可隨時從複製記錄接收。';

  @override
  String get notificationsDenied => '系統通知已關閉，請在系統設定中開啟；待接收內容仍可從複製記錄接收。';

  @override
  String get notificationsUnknown => '通知顯示和勿擾由系統控制；未看到提示時，可在複製記錄中接收。';

  @override
  String get notificationsUnavailable => '系統通知暫時無法使用，請從複製記錄接收大檔案。';

  @override
  String get enableNotifications => '開啟提醒';

  @override
  String get notificationSettings => '通知設定';

  @override
  String get notifyPendingTitle => '接收後即可貼上';

  @override
  String notifyPendingBody(String device, String name, String size) {
    return '$device · $name · $size\n尚未接收，暫不可貼上。亦可從複製記錄接收。';
  }

  @override
  String get notifyReceivingTitle => '正在接收…';

  @override
  String notifyReceivingBody(String name) {
    return '$name\n可在複製記錄查看接收進度。';
  }

  @override
  String get notifyReadyTitle => '已就緒，可以貼上';

  @override
  String notifyReadyBody(String name) {
    return '$name\n已接收並放入剪貼簿。';
  }

  @override
  String get notifyReceivedTitle => '已接收';

  @override
  String notifyReceivedBody(String name) {
    return '$name\n已保留目前剪貼簿，按放入剪貼簿後即可貼上，亦可從複製記錄操作。';
  }

  @override
  String get notifyFailedTitle => '接收失敗';

  @override
  String notifyFailedBody(String name, String error) {
    return '$name\n$error\n可從此通知或複製記錄重試。';
  }

  @override
  String get notifyClipboardFailedTitle => '已接收，但無法寫入剪貼簿';

  @override
  String notifyClipboardFailedBody(String name, String error) {
    return '$name\n$error\n可從此通知或複製記錄重新放入剪貼簿。';
  }
}

/// The translations for Chinese, as used in Taiwan, using the Han script (`zh_Hant_TW`).
class AppLocalizationsZhHantTw extends AppLocalizationsZh {
  AppLocalizationsZhHantTw() : super('zh_Hant_TW');

  @override
  String get appName => '隔空貼';

  @override
  String get navHistory => '複製記錄';

  @override
  String get navDevices => '裝置';

  @override
  String get navSettings => '設定';

  @override
  String get linkConnecting => '正在連線…';

  @override
  String get linkServiceNotRunning => '背景服務未執行';

  @override
  String get linkLanMode => '在區域網路中尋找裝置';

  @override
  String get linkServerUnreachable => '連不上伺服器，區域網路內照常同步';

  @override
  String get linkNothingOn => '區域網路探索已關閉，也沒有設定伺服器';

  @override
  String get linkLanBlocked => 'macOS 阻擋了區域網路存取';

  @override
  String linkOnlinePeers(int count) {
    return '$count 台裝置在線上';
  }

  @override
  String get linkServerConnected => '已連線到伺服器';

  @override
  String linkConnectionLost(String detail) {
    return '連線中斷：$detail';
  }

  @override
  String get linkServiceDisconnected => '背景服務已中斷連線';

  @override
  String get linkServiceDisabled => '背景服務已停用';

  @override
  String get errServiceNotConnected => '背景服務未連線';

  @override
  String get errUnimplemented => '這項功能尚未提供';

  @override
  String get errUnavailable => '無法連線到背景服務';

  @override
  String errServiceStart(String detail) {
    return '啟動背景服務失敗：$detail';
  }

  @override
  String errServiceRestart(String detail) {
    return '重新啟動背景服務失敗：$detail';
  }

  @override
  String errAutostart(String detail) {
    return '設定開機自動啟動失敗：$detail';
  }

  @override
  String get thisMacTitle => '這台 Mac';

  @override
  String get thisPcTitle => '這台電腦';

  @override
  String get thisMacInText => '這台 Mac ';

  @override
  String get thisPcInText => '這台電腦';

  @override
  String get bothMacs => '兩台 Mac ';

  @override
  String get bothDevices => '兩台裝置';

  @override
  String get thisDeviceShort => '本機';

  @override
  String fromDevice(String name) {
    return '來自 $name';
  }

  @override
  String get cancel => '取消';

  @override
  String get close => '關閉';

  @override
  String get save => '儲存';

  @override
  String get retry => '重試';

  @override
  String get continueAction => '繼續';

  @override
  String get justNow => '剛剛';

  @override
  String secondsAgo(int n) {
    return '$n 秒前';
  }

  @override
  String minutesAgo(int n) {
    return '$n 分鐘前';
  }

  @override
  String hoursAgo(int n) {
    return '$n 小時前';
  }

  @override
  String daysAgo(int n) {
    return '$n 天前';
  }

  @override
  String monthDay(int month, int day) {
    return '$month月$day日';
  }

  @override
  String durationHours(int n) {
    return '$n 小時';
  }

  @override
  String durationDays(int n) {
    return '$n 天';
  }

  @override
  String get historyEmptySubtitle => '在任何一台裝置上複製，內容會出現在這裡';

  @override
  String historyCount(int count) {
    return '共 $count 筆';
  }

  @override
  String historyCountKept(int count, String ttl) {
    return '共 $count 筆，保留 $ttl';
  }

  @override
  String get filterAll => '全部';

  @override
  String get filterText => '文字';

  @override
  String get filterImage => '圖片';

  @override
  String get filterFile => '檔案';

  @override
  String get historyClearTooltip => '清空記錄';

  @override
  String get historyEmptyTitle => '還沒有記錄';

  @override
  String get historyEmptyFilteredTitle => '沒有這一類的記錄';

  @override
  String get historyEmptyBody => '在任何一台已配對的裝置上複製文字、圖片或檔案，記錄會出現在這裡。';

  @override
  String get clearConfirmTitle => '要清空全部記錄嗎？';

  @override
  String clearConfirmBody(String device) {
    return '從所有裝置同步過來的記錄和已快取的檔案，都會從$device上刪除。原始檔案不受影響。';
  }

  @override
  String get clear => '清空';

  @override
  String get historyCleared => '記錄已清空';

  @override
  String get permDeniedTitle => '已禁止隔空貼讀取剪貼簿';

  @override
  String get permAskTitle => '允許隔空貼讀取剪貼簿';

  @override
  String get permDeniedBody =>
      '在這台 Mac 上複製的內容不會同步出去。請到系統設定裡把 PasteFerry Daemon 改為「允許」。';

  @override
  String get permAskBody =>
      '否則 macOS 每次都會跳出視窗詢問。請到系統設定裡把 PasteFerry Daemon 改為「允許」。';

  @override
  String get openSystemSettings => '打開系統設定';

  @override
  String get fetch => '接收';

  @override
  String get fetchStarted => '正在接收，可在複製記錄查看進度';

  @override
  String get transferFailed => '傳輸失敗';

  @override
  String get expired => '已過期';

  @override
  String get putOnClipboard => '放入剪貼簿';

  @override
  String get putOnClipboardDone => '已放入剪貼簿';

  @override
  String get deleteRecord => '刪除這筆記錄';

  @override
  String get deleted => '已刪除';

  @override
  String get emptyRecord => '（空）';

  @override
  String get imageRecord => '圖片';

  @override
  String itemsSummary(String name, int others, int count) {
    return '$name 等 $count 項';
  }

  @override
  String get searchHint => '搜尋文字和檔名';

  @override
  String get searchClear => '清除搜尋';

  @override
  String searchResultCount(int count) {
    return '找到 $count 筆';
  }

  @override
  String searchEmptyTitle(String query) {
    return '找不到「$query」';
  }

  @override
  String get searchEmptyBody => '搜尋範圍是複製過的文字和檔名。';

  @override
  String get imagePreviewTitle => '圖片預覽';

  @override
  String get closePreview => '關閉預覽';

  @override
  String get previousImage => '上一張';

  @override
  String get nextImage => '下一張';

  @override
  String get previewHint => '用滾輪或觸控板縮放，拖曳查看';

  @override
  String get fitWindow => '符合視窗大小';

  @override
  String get previewDeleted => '這筆圖片記錄已被刪除';

  @override
  String get previewExpired => '圖片已過期，無法預覽';

  @override
  String get previewDownloading => '正在下載圖片…';

  @override
  String get previewDownloadFailed => '圖片下載失敗，請重試';

  @override
  String get previewNotDownloaded => '圖片尚未下載到這台電腦';

  @override
  String get downloadAndPreview => '下載並預覽';

  @override
  String get previewUnavailable => '圖片暫時無法預覽';

  @override
  String get imageSemantic => '複製記錄中的圖片';

  @override
  String get previewUnreadable => '無法讀取圖片，檔案可能已被清除或損毀';

  @override
  String get devicesSubtitleEmpty => '配對之後，裝置之間就會同步剪貼簿';

  @override
  String devicesSubtitle(int paired, int online) {
    return '已配對 $paired 台，$online 台在線上';
  }

  @override
  String get addDevice => '新增裝置';

  @override
  String thisDeviceFootnote(String both) {
    return '配對時，$both會顯示同樣的兩行安全指紋：這台的和對方的。逐字核對一致才能確認。';
  }

  @override
  String get noDevicesTitle => '還沒有配對的裝置';

  @override
  String get noDevicesBody =>
      '在另一台電腦上打開隔空貼，用 6 碼配對碼把兩台連起來。兩台在同一個網路裡時，也可以在兩台上都打開「搜尋附近裝置」，點一下對方。';

  @override
  String get pairedDevices => '已配對的裝置';

  @override
  String get pairedFootnote =>
      '直連：資料在兩台裝置之間直接傳輸。中繼：網路環境打不通直連時，經你的伺服器轉送，內容依然是端對端加密的。';

  @override
  String get offline => '離線';

  @override
  String get direct => '直連';

  @override
  String get relay => '中繼';

  @override
  String get online => '線上';

  @override
  String get moreActions => '更多操作';

  @override
  String get copyFingerprint => '複製安全指紋';

  @override
  String get unpairEllipsis => '解除配對…';

  @override
  String get fingerprintCopied => '安全指紋已複製';

  @override
  String unpairTitle(String name) {
    return '要解除與「$name」的配對嗎？';
  }

  @override
  String get unpairBody => '兩台裝置將不再同步。以後想恢復，需要重新配對並核對指紋。';

  @override
  String get unpair => '解除配對';

  @override
  String get unpaired => '已解除配對';

  @override
  String get enterAddressEllipsis => '輸入它的位址…';

  @override
  String setAddressTitle(String name) {
    return '「$name」的位址';
  }

  @override
  String setAddressBody(String name) {
    return '兩台不在同一網段、或都換了位址時使用。輸入「$name」上「新增裝置 → 產生配對碼」下方顯示的位址。';
  }

  @override
  String get addressSaved => '已記下，正在那個位址上尋找它';

  @override
  String get lanBlockedTitle => 'macOS 阻擋了區域網路存取';

  @override
  String lanBlockedBody(String device) {
    return '在「系統設定 → 隱私權與安全性 → 本機網路」裡開啟 PasteFerry Daemon。開啟之前，$device連不上同一個網路裡的裝置。';
  }

  @override
  String get searchNearby => '搜尋附近裝置';

  @override
  String get nearbyTitle => '附近的裝置';

  @override
  String get nearbySearching => '正在尋找附近的裝置…';

  @override
  String get nearbyHint => '在另一台電腦上也打開「搜尋附近裝置」，它就會出現在這裡，點一下即可配對。';

  @override
  String nearbyVisible(String device) {
    return '這個視窗開著的期間，同一個網路裡的裝置能看到$device的名稱。';
  }

  @override
  String get nearbyPaired => '已配對';

  @override
  String nearbyRequesting(String name) {
    return '正在聯絡「$name」…';
  }

  @override
  String nearbyNoReply(String name) {
    return '「$name」沒有回應。請確認它還開著「搜尋附近裝置」。';
  }

  @override
  String get nearbyLanOffTitle => '區域網路探索已關閉';

  @override
  String get nearbyLanOffBody => '在「設定」裡開啟「在區域網路中尋找裝置」才能看到附近的裝置，或者改用配對碼。';

  @override
  String get generateCode => '產生配對碼';

  @override
  String get enterCode => '輸入配對碼';

  @override
  String get codeFailedTitle => '無法產生配對碼';

  @override
  String get enterOnOther => '在另一台裝置上輸入這組配對碼';

  @override
  String get codeCopied => '配對碼已複製';

  @override
  String get codeHint => '點一下即可複製，5 分鐘內有效';

  @override
  String afterEnterFootnote(String both) {
    return '對方輸入後，$both會顯示同樣的兩行安全指紋，核對一致再確認。';
  }

  @override
  String get codeInvalid => '配對碼是 6 個字母或數字';

  @override
  String get localAddress => '本機位址';

  @override
  String get addressHint => '不在同一網段時，另一台還要輸入這個位址';

  @override
  String get addressCopied => '位址已複製';

  @override
  String get differentSubnet => '不在同一網段？';

  @override
  String get peerAddress => '對方配對碼下方顯示的位址';

  @override
  String get pairNotFound => '找不到那台裝置。不在同一網段時，請輸入它配對碼下方顯示的位址。';

  @override
  String get pairNotFoundAddress =>
      '那個位址沒有回應。請檢查位址是否正確、對方是否還停在配對碼畫面；也可能是對方的防火牆擋住了。';

  @override
  String get codeWrong => '配對碼不正確，或已經失效';

  @override
  String get codeExpired => '配對碼已失效，請在另一台上重新產生';

  @override
  String get addressInvalid => '位址格式不正確，應該像 192.168.1.20 這樣';

  @override
  String get pairNoChannel => '請先在「設定」裡開啟區域網路探索，或者填寫伺服器位址。';

  @override
  String get enterCodeShown => '輸入另一台裝置上顯示的配對碼';

  @override
  String get pairingDone => '配對完成';

  @override
  String get pairingRejected => '已拒絕配對';

  @override
  String get verifyTitle => '核對安全指紋';

  @override
  String pairingWith(String name) {
    return '正在與「$name」配對';
  }

  @override
  String verifyBody(String both) {
    return '$both上顯示的這兩行應該完全相同。只要有任何一個字元不同，就表示連線可能遭到第三方竄改，請拒絕。';
  }

  @override
  String get rejectMismatch => '不一致，拒絕';

  @override
  String get confirmMatch => '一致，確認配對';

  @override
  String get server => '伺服器（選用）';

  @override
  String get signalingUrl => '信令伺服器位址';

  @override
  String get signalingUrlDesc => '留空僅在區域網路同步';

  @override
  String get signalingUrlPlaceholder => 'ws://伺服器位址:8787/signal';

  @override
  String get publicStun => '改善公網直接連線';

  @override
  String get publicStunDesc => '公共伺服器可見本機的公用 IP';

  @override
  String get connectionStatus => '連線狀態';

  @override
  String get disconnectedDesc => '請檢查位址與伺服器狀態';

  @override
  String get connected => '已連線';

  @override
  String get notConnected => '未連線';

  @override
  String get notSet => '未設定';

  @override
  String get serverCleared => '已清除伺服器位址，只在區域網路裡同步';

  @override
  String get localNetwork => '區域網路';

  @override
  String get lanDiscovery => '在區域網路中尋找裝置';

  @override
  String get lanDiscoveryDesc => '匿名探索，僅搜尋附近裝置時顯示名稱';

  @override
  String get macLocalNetwork => '本機網路權限';

  @override
  String get macLocalNetworkDesc => '在系統設定中允許 PasteFerry Daemon';

  @override
  String get allowLan => 'Windows 防火牆';

  @override
  String get allowLanDesc => '連線受阻時允許連入，需系統管理員權限';

  @override
  String get allow => '允許';

  @override
  String get allowLanDone => '已允許，其他裝置現在可以連到這台了';

  @override
  String get deviceName => '裝置名稱';

  @override
  String get language => '語言';

  @override
  String get languageSystem => '跟隨系統';

  @override
  String get sync => '同步';

  @override
  String get autoSyncLimit => '自動同步上限';

  @override
  String get autoSyncLimitDesc => '超過上限後，在複製記錄中手動接收';

  @override
  String get autoApply => '收到後直接寫入剪貼簿';

  @override
  String get whatToSync => '同步哪些內容';

  @override
  String get syncFiles => '檔案與資料夾';

  @override
  String get syncText => '純文字';

  @override
  String get syncRich => '帶格式的文字';

  @override
  String get syncImages => '圖片與螢幕截圖';

  @override
  String get storage => '儲存空間';

  @override
  String get storageFootnote => '到期自動清除記錄與快取。';

  @override
  String get keepHistory => '記錄保留';

  @override
  String get keepCache => '快取檔案保留';

  @override
  String get keepCacheDesc => '不超過記錄保留時間';

  @override
  String get cacheUsage => '快取佔用';

  @override
  String get backgroundSync => '背景同步';

  @override
  String get backgroundService => '背景服務';

  @override
  String get backgroundServiceDesc => '登入自動啟動，關閉視窗後繼續同步';

  @override
  String get restart => '重新啟動';

  @override
  String get restarted => '背景服務已重新啟動';

  @override
  String get disableEllipsis => '停用…';

  @override
  String get about => '關於';

  @override
  String get aboutDesc => '開源軟體，採用 AGPL-3.0 授權';

  @override
  String get homepage => '專案首頁';

  @override
  String get feedback => '回報問題';

  @override
  String get disableTitle => '要停用背景同步嗎？';

  @override
  String disableBody(String device) {
    return '$device將停止同步，也不會再開機自動啟動。歷史記錄和配對關係會保留。';
  }

  @override
  String get disable => '停用';

  @override
  String get disabled => '背景同步已停用';

  @override
  String get urlMustBeWs => '位址必須以 ws:// 或 wss:// 開頭';

  @override
  String get serverUpdated => '伺服器位址已更新，正在重新連線';

  @override
  String get deviceNameUpdated => '裝置名稱已更新';

  @override
  String get diagIntro =>
      '報告裡有版本、系統、設定、裝置狀態和最近的記錄檔。裝置名稱、伺服器位址、公用 IP、使用者資料夾和檔名都已替換，複製的內容從不寫進記錄檔。GitHub 上的問題是公開的，貼出去之前請看一遍。';

  @override
  String get diagVerboseTitle => '詳細記錄檔';

  @override
  String get diagVerboseDesc => '連線方面的問題常要看它：開啟後重現一次問題，再回來產生報告。30 分鐘後自動關閉。';

  @override
  String get diagVerboseOn => '開啟 30 分鐘';

  @override
  String diagVerboseUntil(String time) {
    return '已開啟，$time 自動關閉';
  }

  @override
  String get diagVerboseOff => '關閉';

  @override
  String get diagRefresh => '重新產生';

  @override
  String get diagSave => '儲存為檔案';

  @override
  String diagSaved(String name) {
    return '已儲存到「下載」資料夾：$name';
  }

  @override
  String get diagCopy => '複製報告';

  @override
  String get diagCopied => '報告已複製';

  @override
  String get diagNewIssue => '在 GitHub 上回報';

  @override
  String get diagIssueHint => '報告已複製，貼到問題裡的「診斷報告」一欄';

  @override
  String diagFailed(String detail) {
    return '產生報告失敗：$detail';
  }

  @override
  String get svcExtractTitle => '先解壓縮或安裝隔空貼';

  @override
  String get svcExtractBody =>
      '隔空貼現在是直接從壓縮檔裡打開的，關掉之後它所在的暫存資料夾就會被刪除。執行安裝程式，或者把整個資料夾解壓縮出來再打開，就可以啟用背景同步。';

  @override
  String get svcMoveTitle => '先把隔空貼移到「應用程式」';

  @override
  String get svcMoveBody =>
      '隔空貼現在是直接從磁碟映像檔裡打開的。把它拖進「應用程式」資料夾，退出磁碟映像檔，再從「應用程式」裡打開，就可以啟用背景同步。';

  @override
  String get svcStartTitle => '開始使用隔空貼';

  @override
  String get svcStartBody =>
      '隔空貼會在背景執行一個小程式，用來偵測剪貼簿的變化，並與你的其他裝置保持連線。啟用後它會隨系統登入自動啟動，關掉這個視窗也會照常同步。';

  @override
  String get svcEnable => '啟用背景同步';

  @override
  String get svcStaleTitle => '重新啟用背景同步';

  @override
  String get svcStaleBodyWin =>
      '隔空貼換了位置，開機自動啟動還指向原來的地方。重新啟用一次，讓背景程序指向目前這份隔空貼。歷史記錄與配對關係都會保留。';

  @override
  String get svcStaleBodyMac =>
      '隔空貼換了位置，或者之前用舊版安裝腳本裝過。重新啟用一次，讓背景程序指向目前這個 App。歷史記錄與配對關係都會保留。';

  @override
  String get svcReenable => '重新啟用';

  @override
  String get svcStoppedTitle => '背景服務沒有在執行';

  @override
  String get svcStoppedBody => '它可能剛剛結束，系統通常會在幾秒內把它重新啟動。如果一直停在這一頁，請點下面的按鈕重新啟動。';

  @override
  String get svcDevBodyWin =>
      '這是開發版本，旁邊沒有 pasteferryd.exe。在原始碼資料夾執行 scripts\\build-windows.ps1 建置完整的版本。';

  @override
  String get svcDevBodyLinux =>
      '這是開發版本，旁邊沒有 pasteferryd。在原始碼資料夾執行 ./scripts/package-linux.sh 建置完整的版本。';

  @override
  String get svcDevBodyMac =>
      '這是開發版本，App 裡沒有內建背景服務。在原始碼資料夾執行 ./scripts/install-macos.sh 安裝它。';

  @override
  String svcLogPath(String path) {
    return '記錄檔：$path';
  }

  @override
  String get svcCanDisable => '隨時可以在「設定」中停用背景同步。';

  @override
  String get officialService => '官方服務';

  @override
  String get useOfficialService => '使用官方服務';

  @override
  String get officialServiceDesc => '跨網路同步，優先直連';

  @override
  String get officialPrivacy =>
      '僅開啟後連線官方服務並統計必要連線資訊。先嘗試直連，失敗才申請中繼；不傳送內容、真實裝置名稱或原始系統識別碼。';

  @override
  String get officialEnableTitle => '啟用官方公網服務？';

  @override
  String get officialSwitchConfirm => '這會切換目前的自訂伺服器，已儲存的位址仍會保留。';

  @override
  String get customServer => '自訂伺服器';

  @override
  String get saveAndConnect => '儲存並連線';

  @override
  String get officialUnused => '尚未申請中繼，餘額待確認';

  @override
  String get officialUnavailable => '用量暫時無法取得，現有資料可能過期';

  @override
  String get officialQuotaTitle => '今日中繼額度已用完';

  @override
  String get officialQuotaBody => '官方中繼已暫停，北京時間午夜重設。區域網路和公網直連仍可使用；傳輸失敗時可稍後重試。';

  @override
  String get officialPeerQuota => '對方中繼額度不足，直連仍可使用';

  @override
  String get officialPaused => '官方中繼暫時無法使用，直連仍可使用';

  @override
  String get officialViewUsage => '查看用量';

  @override
  String get officialAcknowledge => '知道了';

  @override
  String get officialToday => '今日中繼用量';

  @override
  String get officialUsed => '已使用';

  @override
  String get officialRemaining => '剩餘額度';

  @override
  String get officialDailyLimit => '每日額度';

  @override
  String get officialUpload => '上傳';

  @override
  String get officialDownload => '下載';

  @override
  String get officialReset => '重設時間（北京時間）';

  @override
  String get officialUpdated => '上次更新';

  @override
  String get officialDeviceID => '裝置 ID';

  @override
  String get officialWarning => '今日中繼額度已使用 90%，直連不消耗此額度。';

  @override
  String get officialAccounting => '上下行合計，包含加密傳輸開銷；僅統計官方中繼，不統計直連內容或流量。';

  @override
  String buildDetails(String uiVersion, String serviceVersion, String build) {
    return '介面 $uiVersion · 背景服務 $serviceVersion · 建置 $build';
  }

  @override
  String get notReceived => '尚未接收，暫不可貼上';

  @override
  String get receivedPreserved => '已接收，點選放入剪貼簿後即可貼上';

  @override
  String get clipboardWriteFailed => '已接收，但無法寫入剪貼簿';

  @override
  String get receivingStatus => '正在接收…';

  @override
  String get largeFileNotifications => '大檔案接收提醒';

  @override
  String get notificationsAllowed => '已允許系統通知，接收進度可在複製記錄中查看。';

  @override
  String get notificationsNotRequested => '開啟系統通知，避免大檔案等待接收時沒有提醒；也可隨時從複製記錄接收。';

  @override
  String get notificationsDenied => '系統通知已關閉，請在系統設定中開啟；待接收內容仍可從複製記錄接收。';

  @override
  String get notificationsUnknown => '通知顯示和勿擾由系統控制；未看到提醒時，可在複製記錄中接收。';

  @override
  String get notificationsUnavailable => '系統通知暫時無法使用，請從複製記錄接收大檔案。';

  @override
  String get enableNotifications => '開啟提醒';

  @override
  String get notificationSettings => '通知設定';

  @override
  String get notifyPendingTitle => '接收後即可貼上';

  @override
  String notifyPendingBody(String device, String name, String size) {
    return '$device · $name · $size\n尚未接收，暫不可貼上。也可從複製記錄接收。';
  }

  @override
  String get notifyReceivingTitle => '正在接收…';

  @override
  String notifyReceivingBody(String name) {
    return '$name\n可在複製記錄查看接收進度。';
  }

  @override
  String get notifyReadyTitle => '已就緒，可以貼上';

  @override
  String notifyReadyBody(String name) {
    return '$name\n已接收並放入剪貼簿。';
  }

  @override
  String get notifyReceivedTitle => '已接收';

  @override
  String notifyReceivedBody(String name) {
    return '$name\n已保留目前剪貼簿，點選放入剪貼簿後即可貼上，也可從複製記錄操作。';
  }

  @override
  String get notifyFailedTitle => '接收失敗';

  @override
  String notifyFailedBody(String name, String error) {
    return '$name\n$error\n可從此通知或複製記錄重試。';
  }

  @override
  String get notifyClipboardFailedTitle => '已接收，但無法寫入剪貼簿';

  @override
  String notifyClipboardFailedBody(String name, String error) {
    return '$name\n$error\n可從此通知或複製記錄重新放入剪貼簿。';
  }
}

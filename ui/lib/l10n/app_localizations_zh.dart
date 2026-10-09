// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Chinese (`zh`).
class AppLocalizationsZh extends AppLocalizations {
  AppLocalizationsZh([String locale = 'zh']) : super(locale);

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
  String get linkNoServer => '未连接服务器';

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
  String get permDeniedTitle => '已禁止 CopySync 读取剪贴板';

  @override
  String get permAskTitle => '允许 CopySync 读取剪贴板';

  @override
  String get permDeniedBody =>
      '在这台 Mac 上复制的内容不会同步出去。到系统设置里把 CopySync Daemon 改为「允许」。';

  @override
  String get permAskBody => '否则 macOS 每次都会弹窗询问。到系统设置里把 CopySync Daemon 改为「允许」。';

  @override
  String get openSystemSettings => '打开系统设置';

  @override
  String get fetch => '拉取到本机';

  @override
  String get fetchStarted => '开始拉取';

  @override
  String get transferFailed => '传输失败';

  @override
  String get expired => '已过期';

  @override
  String get preview => '预览';

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
  String get imagePreviewTitle => '图片预览';

  @override
  String get closePreview => '关闭预览';

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
  String get noDevicesBody => '在另一台电脑上打开 CopySync，用 6 位配对码把两台连起来。';

  @override
  String get noDevicesNeedServer => '先在「设置」里连接信令服务器，才能配对设备。';

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
  String get server => '服务器';

  @override
  String get serverFootnote => '服务器只负责让设备找到彼此、协商直连。设备之间传输的内容它看不到。';

  @override
  String get signalingUrl => '信令服务器地址';

  @override
  String get signalingUrlDesc => '按回车保存，保存后立即重新连接';

  @override
  String get signalingUrlPlaceholder => 'ws://服务器地址:8787/signal';

  @override
  String get publicStun => '用公共服务器探测网络出口';

  @override
  String get publicStunDesc => '公司双线这类多出口网络里，能找到更多可以直连的路径。公共服务器只会看到你的公网 IP';

  @override
  String get connectionStatus => '连接状态';

  @override
  String get connectedDesc => '已连接到信令服务器';

  @override
  String get disconnectedDesc => '连不上服务器。检查地址是否正确、服务器是否在运行。';

  @override
  String get connected => '已连接';

  @override
  String get notConnected => '未连接';

  @override
  String get deviceName => '设备名称';

  @override
  String get deviceNameDesc => '显示在其他设备的列表里';

  @override
  String get language => '语言';

  @override
  String get languageSystem => '跟随系统';

  @override
  String get sync => '同步';

  @override
  String get syncFootnote =>
      '小于阈值的内容在复制时直接推送到其他设备；超过阈值的只同步一条记录，需要时在「复制记录」里点「拉取到本机」，避免大文件无谓地占用带宽和磁盘。';

  @override
  String get autoSyncLimit => '自动同步上限';

  @override
  String get autoSyncLimitDesc => '超过这个大小的文件改为手动拉取';

  @override
  String get autoApply => '收到后直接写入剪贴板';

  @override
  String get autoApplyDesc => '关闭后，需要在「复制记录」里手动放入剪贴板';

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
  String get storageFootnote => '到期的记录与缓存文件会被自动清理。';

  @override
  String get keepHistory => '记录保留';

  @override
  String get keepCache => '缓存文件保留';

  @override
  String get keepCacheDesc => '不会超过记录的保留时长';

  @override
  String get cacheUsage => '缓存占用';

  @override
  String get backgroundSync => '后台同步';

  @override
  String get backgroundFootnote =>
      '停用后不再同步，也不会开机启动。历史记录与配对关系会保留，重新打开 App 即可再次启用。';

  @override
  String get backgroundService => '后台服务';

  @override
  String get backgroundServiceDesc => '随系统登录自动启动，关掉窗口也照常同步';

  @override
  String get restart => '重新启动';

  @override
  String get restarted => '后台服务已重新启动';

  @override
  String get disableEllipsis => '停用…';

  @override
  String get about => '关于';

  @override
  String get aboutDesc => '开源软件，MIT 许可';

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
  String get svcExtractTitle => '先解压或安装 CopySync';

  @override
  String get svcExtractBody =>
      'CopySync 现在是直接从压缩包里打开的，关掉之后它所在的临时文件夹就会被删除。运行安装程序，或者把整个文件夹解压出来再打开，就可以启用后台同步。';

  @override
  String get svcMoveTitle => '先把 CopySync 移到「应用程序」';

  @override
  String get svcMoveBody =>
      'CopySync 现在是直接从安装盘里打开的。把它拖进「应用程序」文件夹，推出安装盘，再从「应用程序」里打开，就可以启用后台同步。';

  @override
  String get svcStartTitle => '开始使用 CopySync';

  @override
  String get svcStartBody =>
      'CopySync 会在后台运行一个小进程，用来感知剪贴板的变化，并与你的其他设备保持连接。启用后它随系统登录自动启动，关掉这个窗口也会照常同步。';

  @override
  String get svcEnable => '启用后台同步';

  @override
  String get svcStaleTitle => '重新启用后台同步';

  @override
  String get svcStaleBodyWin =>
      'CopySync 换了位置，开机自启还指向原来的地方。重新启用一次，让后台进程指向当前这份 CopySync。历史记录与配对关系都会保留。';

  @override
  String get svcStaleBodyMac =>
      'CopySync 换了位置，或者之前用旧版安装脚本装过。重新启用一次，让后台进程指向当前这个 App。历史记录与配对关系都会保留。';

  @override
  String get svcReenable => '重新启用';

  @override
  String get svcStoppedTitle => '后台服务没有运行';

  @override
  String get svcStoppedBody => '它可能刚刚退出，系统通常会在几秒内把它重新拉起。如果一直停在这一页，点下面的按钮重新启动。';

  @override
  String get svcDevBodyWin =>
      '这是开发构建，旁边没有 copysyncd.exe。在源码目录运行 scripts\\build-windows.ps1 构建完整的版本。';

  @override
  String get svcDevBodyMac =>
      '这是开发构建，App 里没有内置后台服务。在源码目录运行 ./scripts/install-macos.sh 安装它。';

  @override
  String svcLogPath(String path) {
    return '运行日志：$path';
  }

  @override
  String get svcCanDisable => '随时可以在「设置」中停用后台同步。';
}

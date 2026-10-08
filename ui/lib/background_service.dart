import 'dart:io';

/// 构建时注入的版本号（build.sh 传 --dart-define=APP_VERSION）。
/// 用来判断后台服务是不是跟界面同一个版本。
const appVersion = String.fromEnvironment('APP_VERSION', defaultValue: 'dev');

/// 后台服务的安装状态。
enum ServiceState {
  /// 开发构建：App 包里没有内置后台服务，只能用源码目录的脚本安装
  unavailable,

  /// 直接从 DMG 里打开、或被系统随机化路径（App Translocation）运行：
  /// 这时注册的登录项会指向一个随时消失的路径，必须先把 App 移到「应用程序」
  mustMove,

  /// 尚未注册为登录项
  notInstalled,

  /// 已注册，但指向的不是当前这个 App（被移动过，或旧版脚本装的）
  stale,

  /// 已注册，launchd 中没有在运行
  stopped,

  running,
}

/// 管理随 App 一起分发的后台服务（剪贴板监听与同步进程）。
///
/// 后台服务内置在 `CopySync.app/Contents/Helpers/` 下，由界面在首次启动时
/// 注册成 launchd 的 LaunchAgent：开机自启、崩溃自动拉起，关掉窗口也照常同步。
/// 用 LaunchAgent 而非 LaunchDaemon：剪贴板属于登录用户的会话，系统级进程访问不到。
class BackgroundService {
  BackgroundService({String? executable, String? home})
      : _executable = executable ?? Platform.resolvedExecutable,
        _home = home ?? Platform.environment['HOME'] ?? '';

  static const label = 'com.copysync.daemon';

  final String _executable;
  final String _home;

  /// `…/CopySync.app`，由可执行文件 `…/CopySync.app/Contents/MacOS/CopySync` 上溯三级。
  String get _bundle => File(_executable).parent.parent.parent.path;

  String get helperApp => '$_bundle/Contents/Helpers/CopySyncDaemon.app';
  String get helperBinary => '$helperApp/Contents/MacOS/copysyncd';
  String get plistPath => '$_home/Library/LaunchAgents/$label.plist';
  String get logPath => '$_home/Library/Logs/CopySync/daemon.log';

  /// 旧版 install.sh 装在这里，迁移时清理掉
  String get _legacyApp => '$_home/Applications/CopySyncDaemon.app';

  bool get available => Platform.isMacOS && File(helperBinary).existsSync();

  bool get _runningFromTemporaryLocation =>
      _bundle.startsWith('/Volumes/') || _bundle.contains('/AppTranslocation/');

  Future<ServiceState> state() async {
    if (!available) return ServiceState.unavailable;
    if (_runningFromTemporaryLocation) return ServiceState.mustMove;
    final plist = File(plistPath);
    if (!await plist.exists()) return ServiceState.notInstalled;
    if (!(await plist.readAsString()).contains(helperBinary)) return ServiceState.stale;
    final r = await Process.run('launchctl', ['print', await _target()]);
    if (r.exitCode != 0) return ServiceState.stopped;
    return '${r.stdout}'.contains('state = running')
        ? ServiceState.running
        : ServiceState.stopped;
  }

  /// 注册并启动。重复调用即为修复：先注销旧的再重新注册。
  Future<void> install() async {
    await _bootout();

    final legacy = Directory(_legacyApp);
    if (await legacy.exists()) await legacy.delete(recursive: true);

    // 用户打开 App 时已经放行过整个包，内置的后台服务不该再被隔离属性拦下：
    // launchd 拉起它时没有任何弹窗可以让用户确认
    await Process.run('xattr', ['-dr', 'com.apple.quarantine', helperApp]);

    await Directory(File(plistPath).parent.path).create(recursive: true);
    await Directory(File(logPath).parent.path).create(recursive: true);
    await File(plistPath).writeAsString(_plist());

    final r = await Process.run('launchctl', ['bootstrap', await _domain(), plistPath]);
    if (r.exitCode != 0) {
      throw ServiceException('启动后台服务失败：${'${r.stderr}'.trim()}');
    }
  }

  /// 让 launchd 结束旧进程并按当前的 App 重新拉起。用于 App 更新后。
  Future<void> restart() async {
    final r = await Process.run('launchctl', ['kickstart', '-k', await _target()]);
    if (r.exitCode != 0) {
      throw ServiceException('重启后台服务失败：${'${r.stderr}'.trim()}');
    }
  }

  /// 停用并移除登录项。数据目录保留，重新启用后历史记录与配对关系都还在。
  Future<void> uninstall() async {
    await _bootout();
    final plist = File(plistPath);
    if (await plist.exists()) await plist.delete();
  }

  Future<void> _bootout() async {
    await Process.run('launchctl', ['bootout', await _target()]);
    // bootout 是异步的：旧服务没注销干净就 bootstrap 会报 "Bootstrap failed: 5"
    for (var i = 0; i < 20; i++) {
      final r = await Process.run('launchctl', ['print', await _target()]);
      if (r.exitCode != 0) return;
      await Future<void>.delayed(const Duration(milliseconds: 250));
    }
  }

  String? _uid;
  Future<String> _domain() async {
    _uid ??= '${(await Process.run('id', ['-u'])).stdout}'.trim();
    return 'gui/$_uid';
  }

  Future<String> _target() async => '${await _domain()}/$label';

  String _plist() => '''
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>              <string>$label</string>
    <key>ProgramArguments</key>
    <array>
        <string>${_xml(helperBinary)}</string>
    </array>
    <key>RunAtLoad</key>          <true/>
    <key>KeepAlive</key>          <true/>
    <key>ThrottleInterval</key>   <integer>10</integer>
    <key>StandardOutPath</key>    <string>${_xml(logPath)}</string>
    <key>StandardErrorPath</key>  <string>${_xml(logPath)}</string>
    <key>ProcessType</key>        <string>Interactive</string>
</dict>
</plist>
''';

  static String _xml(String s) =>
      s.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;');
}

class ServiceException implements Exception {
  ServiceException(this.message);
  final String message;
  @override
  String toString() => message;
}

import 'dart:io';

import 'background_service.dart';

/// Windows：后台服务 `copysyncd.exe` 与 `CopySync.exe` 放在同一个目录。
///
/// 开机自启登记在 `HKCU\...\Run`，只影响当前用户，不需要管理员权限。登记的命令是
/// `copysyncd.exe -supervise`：守护进程拉起真正干活的子进程，子进程崩溃就隔 10 秒
/// 再拉起，相当于 Mac 上 launchd 的 KeepAlive。
///
/// 不做成 Windows 服务：服务运行在会话 0，与用户桌面隔离，访问不到用户的剪贴板。
class WindowsBackgroundService implements BackgroundService {
  WindowsBackgroundService({String? executable, Map<String, String>? environment})
      : _executable = executable ?? Platform.resolvedExecutable,
        _env = environment ?? Platform.environment;

  static const _runKey = r'HKCU\Software\Microsoft\Windows\CurrentVersion\Run';
  static const _valueName = 'CopySync';
  static const _image = 'copysyncd.exe';

  final String _executable;
  final Map<String, String> _env;

  String get daemonPath => '${File(_executable).parent.path}\\$_image';

  /// 登记到自启项里的命令行
  String get _command => '"$daemonPath" -supervise';

  String get _localAppData => _env['LOCALAPPDATA'] ?? '${_env['USERPROFILE']}\\AppData\\Local';

  String get logPath => '$_localAppData\\CopySync\\logs\\daemon.log';

  @override
  String get displayLogPath => logPath.replaceFirst(_localAppData, '%LOCALAPPDATA%');

  @override
  bool get available => Platform.isWindows && File(daemonPath).existsSync();

  /// 在资源管理器里直接双击压缩包里的程序时，它其实运行在临时目录里，
  /// 关掉之后整个目录就被删了
  bool get _runningFromTemporaryLocation {
    final temp = (_env['TEMP'] ?? '').toLowerCase();
    final exe = _executable.toLowerCase();
    return temp.isNotEmpty && exe.startsWith(temp);
  }

  @override
  Future<ServiceState> state() async {
    if (!available) return ServiceState.unavailable;
    if (_runningFromTemporaryLocation) return ServiceState.mustMove;
    final registered = await _registeredCommand();
    if (registered == null) return ServiceState.notInstalled;
    if (!registered.toLowerCase().contains(daemonPath.toLowerCase())) return ServiceState.stale;
    return await _running() ? ServiceState.running : ServiceState.stopped;
  }

  @override
  Future<void> install() async {
    final r = await Process.run('reg', [
      'add', _runKey, '/v', _valueName, '/t', 'REG_SZ', '/d', _command, '/f', //
    ]);
    if (r.exitCode != 0) {
      throw ServiceException('登记开机自启失败：${_output(r)}');
    }
    await restart();
  }

  @override
  Future<void> restart() async {
    await _stop();
    // 脱离界面进程运行：关掉窗口后照常同步
    await Process.start(daemonPath, ['-supervise'], mode: ProcessStartMode.detached);
  }

  @override
  Future<void> uninstall() async {
    await Process.run('reg', ['delete', _runKey, '/v', _valueName, '/f']);
    await _stop();
  }

  /// 守护进程与后台服务一起结束：只结束后台服务的话，守护进程会把它重新拉起来
  Future<void> _stop() async {
    await Process.run('taskkill', ['/F', '/T', '/IM', _image]);
    for (var i = 0; i < 20 && await _running(); i++) {
      await Future<void>.delayed(const Duration(milliseconds: 250));
    }
  }

  Future<bool> _running() async {
    final r = await Process.run('tasklist', ['/FI', 'IMAGENAME eq $_image', '/NH']);
    return '${r.stdout}'.toLowerCase().contains(_image);
  }

  /// 自启项里登记的命令；没有登记时为 null。
  ///
  /// `reg query` 的输出形如：
  ///
  ///     HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
  ///         CopySync    REG_SZ    "C:\...\copysyncd.exe" -supervise
  Future<String?> _registeredCommand() async {
    final r = await Process.run('reg', ['query', _runKey, '/v', _valueName]);
    if (r.exitCode != 0) return null;
    for (final line in '${r.stdout}'.split('\n')) {
      final i = line.indexOf('REG_SZ');
      if (line.trimLeft().startsWith(_valueName) && i >= 0) {
        return line.substring(i + 'REG_SZ'.length).trim();
      }
    }
    return null;
  }

  static String _output(ProcessResult r) => '${r.stderr}${r.stdout}'.trim();
}

import 'dart:io';

import 'background_service.dart';
import 'i18n.dart';

/// Linux：后台服务 `pasteferryd` 与界面程序放在同一个目录（.deb 装在 /opt/PasteFerry）。
///
/// 开机自启用 XDG 自启动项（`~/.config/autostart/` 下的 .desktop 文件），GNOME、KDE、Xfce 等桌面都认。
/// 它在登录后的图形会话里启动，带着 DISPLAY、XAUTHORITY 等环境变量，后台服务才连得上剪贴板。
/// 登记的命令是 `pasteferryd -supervise`：守护进程拉起真正干活的子进程，崩溃就隔 10 秒再拉起，
/// 与 Windows 的做法相同。
///
/// 不用 systemd 的用户服务：它不一定拿得到图形会话的环境变量，各发行版的做法也不一致。
class LinuxBackgroundService implements BackgroundService {
  LinuxBackgroundService({String? executable, Map<String, String>? environment})
      : _executable = executable ?? Platform.resolvedExecutable,
        _env = environment ?? Platform.environment;

  static const _image = 'pasteferryd';

  final String _executable;
  final Map<String, String> _env;

  String get daemonPath => '${File(_executable).parent.path}/$_image';

  String get _home => _env['HOME'] ?? '';

  /// XDG 目录：环境变量要是绝对路径才算数，否则用默认位置
  String _xdg(String name, String fallback) {
    final v = _env[name] ?? '';
    return v.startsWith('/') ? v : '$_home/$fallback';
  }

  String get autostartPath => '${_xdg('XDG_CONFIG_HOME', '.config')}/autostart/pasteferry-daemon.desktop';

  String get logPath => '${_xdg('XDG_DATA_HOME', '.local/share')}/pasteferry/logs/daemon.log';

  @override
  String get displayLogPath => _home.isEmpty ? logPath : logPath.replaceFirst(_home, '~');

  @override
  bool get available => Platform.isLinux && File(daemonPath).existsSync();

  /// 自启动项里的命令行
  String get _command => '${quoteExecArg(daemonPath)} -supervise';

  @override
  Future<ServiceState> state() async {
    if (!available) return ServiceState.unavailable;
    final file = File(autostartPath);
    if (!file.existsSync()) return ServiceState.notInstalled;
    final exec = execLine(await file.readAsString());
    if (exec == null || !exec.contains(quoteExecArg(daemonPath))) return ServiceState.stale;
    return await _running() ? ServiceState.running : ServiceState.stopped;
  }

  @override
  Future<void> install() async {
    try {
      final file = File(autostartPath);
      await file.parent.create(recursive: true);
      await file.writeAsString(desktopEntry(_command));
    } on FileSystemException catch (e) {
      throw ServiceException(appL10n.errAutostart(e.message));
    }
    await restart();
  }

  @override
  Future<void> restart() async {
    await _stop();
    // 脱离界面进程运行：关掉窗口后照常同步。环境变量随之传下去，后台服务才连得上剪贴板
    await Process.start(daemonPath, ['-supervise'], mode: ProcessStartMode.detached);
  }

  @override
  Future<void> uninstall() async {
    final file = File(autostartPath);
    if (file.existsSync()) await file.delete();
    await _stop();
  }

  @override
  Future<bool> removeLegacy() async {
    const legacyImage = 'copysyncd';
    final old = File('${_xdg('XDG_CONFIG_HOME', '.config')}/autostart/copysync-daemon.desktop');
    final enabled = old.existsSync();
    if (enabled) await old.delete();
    if (await _running(legacyImage)) await _stop(legacyImage);
    return enabled;
  }

  /// 守护进程与后台服务一起结束：只结束后台服务的话，守护进程会把它重新拉起来
  Future<void> _stop([String image = _image]) async {
    await Process.run('pkill', [..._ownedBy, '-x', image]);
    for (var i = 0; i < 20 && await _running(image); i++) {
      await Future<void>.delayed(const Duration(milliseconds: 250));
    }
  }

  Future<bool> _running([String image = _image]) async {
    final r = await Process.run('pgrep', [..._ownedBy, '-x', image]);
    return r.exitCode == 0;
  }

  /// 只管当前用户的进程：同一台电脑上别的用户也可能在运行 PasteFerry
  List<String> get _ownedBy {
    final user = _env['USER'] ?? '';
    return user.isEmpty ? const [] : ['-u', user];
  }
}

/// XDG 自启动项。NoDisplay：不出现在应用菜单里；它只是后台服务的开机启动入口。
String desktopEntry(String exec) => '''[Desktop Entry]
Type=Application
Name=PasteFerry
Comment=Keeps the clipboard in sync with your other devices
Exec=$exec
Icon=com.pasteferry.app
Terminal=false
NoDisplay=true
X-GNOME-Autostart-enabled=true
''';

/// 按 Desktop Entry 规范给 Exec 里的一个参数加引号：路径里可能有空格。
/// 引号里的 `"`、`` ` ``、`$`、`\` 要加反斜杠转义，而 .desktop 文件本身又把每个 `\` 转义一次：
/// 所以前三个写成两个反斜杠加原字符，反斜杠本身写成四个。
String quoteExecArg(String arg) {
  final escaped = arg.replaceAllMapped(
      RegExp(r'["`$\\]'), (m) => m[0] == r'\' ? r'\\\\' : r'\\' + m[0]!);
  return '"$escaped"';
}

/// 取出 .desktop 文件里 [Desktop Entry] 一节的 Exec 一行。
String? execLine(String content) {
  for (final line in content.split('\n')) {
    if (line.startsWith('Exec=')) return line.substring('Exec='.length).trim();
  }
  return null;
}

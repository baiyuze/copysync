import 'dart:io';

import 'i18n.dart';

/// 随平台变化的叫法：Mac 上说「这台 Mac」，Windows 与 Linux 上说「这台电脑」。
abstract final class Wording {
  static final bool _win = !Platform.isMacOS;

  /// 说到本机时的叫法，用作标题
  static String get thisDevice => _win ? appL10n.thisPcTitle : appL10n.thisMacTitle;

  /// 用在句子中间。中日文里「Mac」后面紧跟汉字，文案里已带好空格
  static String get thisDeviceInText => _win ? appL10n.thisPcInText : appL10n.thisMacInText;

  /// 说到配对双方时的叫法，用在句子中间：Windows 与 Linux 上对方可能是 Mac，也可能是另一台电脑
  static String get bothDevicesInText => _win ? appL10n.bothDevices : appL10n.bothMacs;
}

/// 用默认浏览器打开网址。
void openUrl(String url) {
  if (Platform.isWindows) {
    // explorer 把网址交给默认浏览器；不经 cmd /c start，免得网址里的 & 被当成命令分隔符
    Process.run('explorer.exe', [url]);
  } else if (Platform.isLinux) {
    Process.run('xdg-open', [url]);
  } else {
    Process.run('open', [url]);
  }
}

/// 用户的「下载」文件夹。Linux 上这个文件夹按系统语言命名（比如「下载」），以 XDG 的设置为准。
String downloadsDir() {
  final env = Platform.environment;
  if (Platform.isWindows) return '${env['USERPROFILE']}\\Downloads';
  final home = env['HOME'] ?? '';
  if (Platform.isLinux) {
    final config = env['XDG_CONFIG_HOME'] ?? '$home/.config';
    try {
      for (final line in File('$config/user-dirs.dirs').readAsLinesSync()) {
        final m = RegExp(r'^XDG_DOWNLOAD_DIR="(.+)"$').firstMatch(line.trim());
        if (m != null) return m.group(1)!.replaceFirst(r'$HOME', home);
      }
    } on FileSystemException {
      // 没有这个文件就用默认的
    }
  }
  return '$home/Downloads';
}

/// 在访达、资源管理器或文件管理器里显示这个文件。Linux 上的文件管理器各不相同，打开它所在的文件夹。
void revealFile(String path) {
  if (Platform.isWindows) {
    Process.run('explorer.exe', ['/select,', path]);
  } else if (Platform.isLinux) {
    Process.run('xdg-open', [File(path).parent.path]);
  } else {
    Process.run('open', ['-R', path]);
  }
}

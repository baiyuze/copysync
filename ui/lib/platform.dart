import 'dart:io';

/// 随平台变化的叫法。Mac 上保持原来的文案。
abstract final class Wording {
  static final bool _win = Platform.isWindows;

  /// 说到本机时的叫法，用作标题
  static String get thisDevice => _win ? '这台电脑' : '这台 Mac';

  /// 用在句子中间：后面紧跟汉字，Mac 上要在「Mac」后留一个空格
  static String get thisDeviceInText => _win ? '这台电脑' : '这台 Mac ';

  /// 说到配对双方时的叫法，用在句子中间：Windows 上对方可能是 Mac，也可能是另一台电脑
  static String get bothDevicesInText => _win ? '两台设备' : '两台 Mac ';
}

/// 用默认浏览器打开网址。
void openUrl(String url) {
  if (Platform.isWindows) {
    // explorer 把网址交给默认浏览器；不经 cmd /c start，免得网址里的 & 被当成命令分隔符
    Process.run('explorer.exe', [url]);
  } else {
    Process.run('open', [url]);
  }
}

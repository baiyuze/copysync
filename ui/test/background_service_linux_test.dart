import 'package:pasteferry_ui/background_service_linux.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('自启动项里的路径按 Desktop Entry 规范加引号', () {
    expect(quoteExecArg('/opt/PasteFerry/pasteferryd'), '"/opt/PasteFerry/pasteferryd"');
    // 路径里有空格时整体放在引号里
    expect(quoteExecArg('/home/a/下载/PasteFerry 1/pasteferryd'), '"/home/a/下载/PasteFerry 1/pasteferryd"');
    // 引号里的特殊字符要转义，.desktop 文件本身又把反斜杠转义一次
    expect(quoteExecArg(r'/x/$a"b`c\d'), r'"/x/\\$a\\"b\\`c\\\\d"');
  });

  test('能从写出去的自启动项里读回同一条命令', () {
    final exec = '${quoteExecArg('/opt/PasteFerry/pasteferryd')} -supervise';
    final entry = desktopEntry(exec);
    expect(entry, startsWith('[Desktop Entry]\n'));
    expect(entry, contains('NoDisplay=true'));
    expect(execLine(entry), exec);
    expect(execLine('[Desktop Entry]\nName=x\n'), isNull);
  });
}

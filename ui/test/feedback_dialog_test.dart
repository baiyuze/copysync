import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:pasteferry_ui/app_state.dart';
import 'package:pasteferry_ui/gen/pasteferry/v1/daemon.pb.dart';
import 'package:pasteferry_ui/main.dart';
import 'package:pasteferry_ui/widgets/feedback_dialog.dart';

/// 不连后台服务：报告与详细日志的开关都由测试给出。
/// 「在 GitHub 上提问题」与「保存为文件」不在测试里按：一个会打开浏览器，一个会写进真的「下载」文件夹。
class _FeedbackState extends AppState {
  int exports = 0;
  final verboseCalls = <int>[];
  Object? failWith;

  @override
  Future<String> exportDiagnostics({int logLines = 200}) async {
    exports++;
    if (failWith != null) throw failWith!;
    return 'PasteFerry diagnostics report\n  Background service     1.7.0\n#$exports';
  }

  @override
  Future<void> setVerboseLog(int minutes) async {
    verboseCalls.add(minutes);
    final until = minutes == 0 ? 0 : DateTime(2030, 1, 1, 21, 20).millisecondsSinceEpoch ~/ 1000;
    debugSeed(status: status!..verboseLogUntilUnix = Int64(until), config: config!, self: self!);
  }
}

_FeedbackState _seed() => _FeedbackState()
  ..debugSeed(
    status: Status(signalingConnected: true, version: '1.7.0'),
    config: Config(language: 'en'),
    self: Device(id: 'self', name: 'MacBook Pro'),
  );

/// 常见的窗口大小。最小窗口下有没有挤坏，由 i18n_test 里逐个语言检查
Future<void> _open(WidgetTester tester, AppState state) async {
  tester.view.physicalSize = const Size(1100, 800);
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);
  await tester.pumpWidget(PasteFerryApp(state: state));
  showFeedbackDialog(tester.element(find.byType(Shell)));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('打开时生成报告给用户看，拷贝的就是看到的这份', (tester) async {
    String? copied;
    tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(SystemChannels.platform, (call) async {
      if (call.method == 'Clipboard.setData') copied = (call.arguments as Map)['text'] as String;
      return null;
    });
    addTearDown(() => tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(SystemChannels.platform, null));

    final state = _seed();
    await _open(tester, state);
    expect(find.textContaining('PasteFerry diagnostics report'), findsOneWidget);

    await tester.tap(find.text('Copy report'));
    await tester.pumpAndSettle();
    expect(copied, contains('Background service     1.7.0'));
    expect(find.text('Report copied'), findsOneWidget);
  });

  testWidgets('打开详细日志后显示到几点，并重新生成报告；可以提前关掉', (tester) async {
    final state = _seed();
    await _open(tester, state);
    expect(state.exports, 1);

    await tester.tap(find.text('Turn on for 30 minutes'));
    await tester.pumpAndSettle();
    expect(state.verboseCalls, [30]);
    expect(find.text('On until 21:20'), findsOneWidget);
    expect(state.exports, 2);
    expect(find.textContaining('#2'), findsOneWidget);

    await tester.tap(find.text('Turn off'));
    await tester.pumpAndSettle();
    expect(state.verboseCalls, [30, 0]);
    expect(find.text('Turn on for 30 minutes'), findsOneWidget);
  });

  testWidgets('生成失败时说明原因，可以重试；拷贝与保存先不能按', (tester) async {
    final state = _seed()..failWith = Exception('daemon gone');
    await _open(tester, state);
    expect(find.textContaining("Couldn't create the report"), findsOneWidget);
    expect(tester.widget<OutlinedButton>(find.widgetWithText(OutlinedButton, 'Copy report')).onPressed, isNull);
    expect(tester.widget<OutlinedButton>(find.widgetWithText(OutlinedButton, 'Save as file')).onPressed, isNull);

    state.failWith = null;
    await tester.tap(find.text('Try again'));
    await tester.pumpAndSettle();
    expect(find.textContaining('PasteFerry diagnostics report'), findsOneWidget);
  });
}

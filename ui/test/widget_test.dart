import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:copysync_ui/daemon_client.dart';
import 'package:copysync_ui/main.dart';
import 'package:copysync_ui/widgets/common.dart';
import 'package:copysync_ui/widgets/pairing_dialog.dart';

void main() {
  testWidgets('应用启动后显示三个导航区', (tester) async {
    await tester.pumpWidget(const CopySyncApp());
    await tester.pump();

    expect(find.text('CopySync'), findsOneWidget);
    // 「复制记录」同时出现在侧边栏与页面标题中
    expect(find.text('复制记录'), findsWidgets);
    expect(find.text('设备'), findsWidgets);
    expect(find.text('设置'), findsWidgets);
  });

  testWidgets('daemon 未运行时给出下一步指引而非空白页', (tester) async {
    await tester.pumpWidget(const CopySyncApp());
    await tester.pump(const Duration(milliseconds: 300));

    // 空状态要说明该做什么，不能只写"暂无数据"
    expect(find.text('还没有记录'), findsOneWidget);
    expect(find.byType(EmptyState), findsOneWidget);
  });

  testWidgets('可以切换到设备页', (tester) async {
    await tester.pumpWidget(const CopySyncApp());
    await tester.pump();

    // 点侧边栏的「设备」（第一个匹配项是导航条目）
    await tester.tap(find.text('设备').first);
    await tester.pump();

    expect(find.text('还没有配对的设备'), findsOneWidget);
    expect(find.text('添加设备'), findsWidgets);
  });

  testWidgets('配对对话框里能取到 AppState', (tester) async {
    await tester.pumpWidget(const CopySyncApp());
    await tester.pump();

    // 测试环境没有 daemon，「添加设备」按钮是禁用的，直接按页面里的方式弹出
    showDialog<void>(
      context: tester.element(find.byType(Shell)),
      builder: (_) => const PairingDialog(),
    );
    await tester.pumpAndSettle();

    // 对话框挂在 Navigator 上而非 home 之下，AppScope 若放在 home 里，
    // 这里会报 "Null check operator used on a null value"
    expect(find.text('生成配对码'), findsOneWidget);
    expect(find.textContaining('Null check'), findsNothing);
  });

  test('数据目录路径与 Go 侧 config.DefaultPaths 保持一致', () {
    final dir = DaemonEndpoint.defaultDataDir();
    // 路径不一致会导致 UI 根本找不到 daemon，这里锁住约定
    expect(dir, contains('CopySync'));
    expect(DaemonEndpoint.endpointFile(dir).path, endsWith('daemon.json'));
  });

  group('格式化', () {
    test('humanBytes 在各量级下可读', () {
      expect(humanBytes(512), '512 B');
      expect(humanBytes(1536), '1.5 KB');
      expect(humanBytes(50 * 1024 * 1024), '50.0 MB');
      // 三位数时省略小数，避免 "150.0 MB" 这样冗余的显示
      expect(humanBytes(150 * 1024 * 1024), '150 MB');
    });

    test('relativeTime 对新近事件用口语化表述', () {
      final now = DateTime.now();
      expect(relativeTime(now), '刚刚');
      expect(relativeTime(now.subtract(const Duration(minutes: 5))), '5 分钟前');
      expect(relativeTime(now.subtract(const Duration(hours: 3))), '3 小时前');
    });
  });
}

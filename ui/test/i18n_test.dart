import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:copysync_ui/app_state.dart';
import 'package:copysync_ui/background_service.dart';
import 'package:copysync_ui/gen/copysync/v1/daemon.pb.dart';
import 'package:copysync_ui/i18n.dart';
import 'package:copysync_ui/main.dart';
import 'package:copysync_ui/widgets/pairing_dialog.dart';

class _DemoState extends AppState {
  @override
  Future<String> createPairingCode() async => 'K7MP2X';
}

AppState _demo(String language, {LinkState link = LinkState.online, ServiceState? service}) =>
    _DemoState()
      ..debugSeed(
        link: link,
        serviceState: service,
        status: Status(signalingConnected: true, version: '1.3.0'),
        config: Config(
          language: language,
          signalingUrl: 'wss://sync.example.com/signal',
          deviceName: 'MacBook Pro',
          autoSyncThresholdBytes: Int64(50 << 20),
          historyTtlSeconds: Int64(3 * 86400),
          cacheTtlSeconds: Int64(86400),
        ),
        self: Device(id: 'self', name: 'MacBook Pro', publicKeyFingerprint: 'H3QD-7KXA-M2PV'),
        peers: [
          Device(
            id: 'air',
            name: 'MacBook Air',
            online: true,
            connection: ConnectionKind.CONNECTION_KIND_RELAY,
            publicKeyFingerprint: 'R8NF-2WTC-QL5J',
          ),
        ],
        records: [
          ClipRecord(
            id: '1',
            kind: ClipKind.CLIP_KIND_IMAGE,
            originDeviceName: 'MacBook Air',
            textPreview: '图片 2.4 MB',
            totalSize: Int64(2 << 20),
          ),
          ClipRecord(
            id: '2',
            kind: ClipKind.CLIP_KIND_FILE,
            status: ClipStatus.CLIP_STATUS_REMOTE_ONLY,
            originDeviceName: 'MacBook Air',
            textPreview: 'report.key 等 2 项',
            totalSize: Int64(900 << 20),
            items: [ClipItem(name: 'report.key'), ClipItem(name: 'notes.txt')],
          ),
        ],
      );

void main() {
  test('按系统的语言偏好取第一个支持的，都不支持时用英文', () {
    expect(resolveLocale(const [Locale('ja', 'JP')]), const Locale('ja'));
    expect(
      resolveLocale(const [Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hant', countryCode: 'TW')]),
      const Locale('zh'),
    );
    expect(resolveLocale(const [Locale('fr'), Locale('ja')]), const Locale('ja'));
    expect(resolveLocale(const [Locale('fr')]), const Locale('en'));
    expect(resolveLocale(null), const Locale('en'));
  });

  test('设置里的空字符串表示跟随系统', () {
    expect(localeFromSetting(''), isNull);
    expect(localeFromSetting('zh-Hans'), const Locale('zh'));
    expect(localeFromSetting('en'), const Locale('en'));
    expect(localeFromSetting('ja'), const Locale('ja'));
  });

  testWidgets('设置里选的语言优先于系统语言', (tester) async {
    await tester.pumpWidget(CopySyncApp(state: _demo('en')));
    expect(find.text('History'), findsWidgets);
    expect(find.text('复制记录'), findsNothing);
  });

  testWidgets('跟随系统时用系统语言', (tester) async {
    await tester.pumpWidget(CopySyncApp(state: _demo('')));
    expect(find.text('复制记录'), findsWidgets);
  });

  testWidgets('配置里的语言一变，界面立即换过来', (tester) async {
    final state = _demo('');
    await tester.pumpWidget(CopySyncApp(state: state));
    expect(find.text('复制记录'), findsWidgets);

    state.config!.language = 'ja';
    state.notifyListeners();
    await tester.pump();
    expect(find.text('履歴'), findsWidgets);
    expect(find.text('复制记录'), findsNothing);
  });

  testWidgets('图片与文件的标题按界面语言显示，不用后台服务存下的中文摘要', (tester) async {
    await tester.pumpWidget(CopySyncApp(state: _demo('en')));
    expect(find.text('Image'), findsOneWidget);
    expect(find.text('report.key and 1 more'), findsOneWidget);
    expect(find.textContaining('图片'), findsNothing);
  });

  // 英文、日文通常比中文长：在最小窗口尺寸下把每一页、每个对话框都打开一遍，确认没有文字溢出
  for (final language in ['en', 'ja', 'zh-Hans']) {
    group('最小窗口下的 $language 界面', () {
      setUp(() {
        final view = TestWidgetsFlutterBinding.instance.platformDispatcher.views.first;
        view.physicalSize = const Size(760, 520);
        view.devicePixelRatio = 1;
      });
      tearDown(() => TestWidgetsFlutterBinding.instance.platformDispatcher.views.first.reset());

      Future<AppLocalizations> open(WidgetTester tester, AppState state) async {
        await tester.pumpWidget(CopySyncApp(state: state));
        await tester.pump();
        return AppLocalizations.of(tester.element(find.byType(Shell)));
      }

      testWidgets('三个页面', (tester) async {
        final l = await open(tester, _demo(language));
        await tester.tap(find.text(l.navDevices).first);
        await tester.pump();
        expect(find.text(l.pairedDevices), findsOneWidget);
        await tester.tap(find.text(l.navSettings).first);
        await tester.pump();
        await tester.dragUntilVisible(find.text(l.language), find.byType(ListView), const Offset(0, -100));
        await tester.dragUntilVisible(find.text(l.about), find.byType(ListView), const Offset(0, -200));
      });

      testWidgets('配对对话框', (tester) async {
        final l = await open(tester, _demo(language));
        showPairingDialog(tester.element(find.byType(Shell)));
        await tester.pumpAndSettle();
        expect(find.text('K7MP2X'), findsOneWidget);
        await tester.tap(find.text(l.enterCode));
        await tester.pumpAndSettle();
        expect(find.text(l.enterCodeShown), findsOneWidget);
      });

      testWidgets('核对指纹', (tester) async {
        final state = _demo(language);
        state.debugSeed(
          status: state.status!,
          config: state.config!,
          self: state.self!,
          pendingPairing: Device(name: 'MacBook Air', publicKeyFingerprint: 'R8NF-2WTC-QL5J', pairingSession: 's'),
        );
        final l = await open(tester, state);
        await tester.pumpAndSettle();
        expect(find.text(l.confirmMatch), findsOneWidget);
      });

      testWidgets('启用后台服务的引导页', (tester) async {
        final l = await open(tester, _demo(language, link: LinkState.offline, service: ServiceState.notInstalled));
        expect(find.text(l.svcEnable), findsOneWidget);
      });
    });
  }
}

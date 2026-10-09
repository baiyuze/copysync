// 生成 README 与网站用的界面截图：演示数据，不连接后台服务，不读真实剪贴板。
//
// 不要直接运行，用仓库根目录的 tools/screenshots/render.sh——
// 它会先从本机系统字体生成渲染用的字体文件，渲染完再加上 macOS 窗口外框。
//
// 放在 tool/ 而非 test/：它依赖本机字体，不该随 flutter test 在 CI 上跑。
// ignore_for_file: invalid_use_of_visible_for_testing_member
import 'dart:io';
import 'dart:ui' show ImageByteFormat;

import 'package:copysync_ui/app_state.dart';
import 'package:copysync_ui/background_service.dart';
import 'package:copysync_ui/gen/copysync/v1/daemon.pb.dart';
import 'package:copysync_ui/i18n.dart';
import 'package:copysync_ui/main.dart';
import 'package:copysync_ui/theme.dart';
import 'package:copysync_ui/widgets/pairing_dialog.dart';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

/// 截图的界面语言：zh 或 en，由 render.sh 设置。演示数据里的文件名、文字也跟着换。
final _lang = Platform.environment['SCREENSHOT_LANG'] ?? 'zh';
final _locale = Locale(_lang);
String _t(String zh, String en) => _lang == 'en' ? en : zh;

/// 生成配对码不走后台服务，直接给一个固定值。
class _DemoState extends AppState {
  @override
  Future<String> createPairingCode() async => 'K7MP2X';
}

Future<void> _loadFonts() async {
  final dir = Platform.environment['SCREENSHOT_FONTS'];
  if (dir == null) throw StateError('用 tools/screenshots/render.sh 运行，它会设置 SCREENSHOT_FONTS');

  Future<void> family(String name, List<String> files) async {
    final loader = FontLoader(name);
    for (final f in files) {
      loader.addFont(Future.value(ByteData.sublistView(File('$dir/$f').readAsBytesSync())));
    }
    await loader.load();
  }

  await family('SF', ['SF-Regular.ttf', 'SF-Medium.ttf', 'SF-Semibold.ttf']);
  await family('PingFang', ['PingFang-Regular.ttf', 'PingFang-Medium.ttf', 'PingFang-Semibold.ttf']);
  await family(monoFamily, ['Menlo-Regular.ttf', 'Menlo-Bold.ttf']);
  // 图标字体随依赖包打进测试资源，但测试环境不会自动加载
  final icons = FontLoader('packages/cupertino_icons/CupertinoIcons')
    ..addFont(rootBundle.load('packages/cupertino_icons/assets/CupertinoIcons.ttf'));
  await icons.load();
}

final _now = DateTime.now();
Int64 _ago(Duration d) => Int64(_now.subtract(d).millisecondsSinceEpoch ~/ 1000);

ClipRecord _clip(
  String id,
  ClipKind kind, {
  String preview = '',
  bool outgoing = false,
  String from = 'MacBook Air',
  int size = 0,
  required Duration ago,
  ClipStatus status = ClipStatus.CLIP_STATUS_READY,
  List<ClipItem> items = const [],
}) =>
    ClipRecord(
      id: id,
      kind: kind,
      status: status,
      outgoing: outgoing,
      originDeviceName: outgoing ? 'MacBook Pro' : from,
      textPreview: preview,
      totalSize: Int64(size),
      createdAtUnix: _ago(ago),
      items: items,
    );

const _mb = 1024 * 1024;
final _quarterly = _t('季度汇报-终版.key', 'Q3 review - final.key');

AppState _demo({LinkState link = LinkState.online,
    ServiceState? service, Device? pending, bool empty = false}) {
  final state = _DemoState();
  state.debugSeed(
    link: link,
    serviceState: service,
    pendingPairing: pending,
    status: Status(
      version: '1.3.0',
      deviceName: 'MacBook Pro',
      signalingConnected: link == LinkState.online,
      clipboardPermission: ClipboardPermission.CLIPBOARD_PERMISSION_ALWAYS_ALLOW,
      cacheBytesUsed: Int64(412 * _mb),
      peersOnline: 2,
    ),
    config: Config(
      signalingUrl: 'wss://sync.example.com/signal',
      deviceName: 'MacBook Pro',
      autoSyncThresholdBytes: Int64(50 * _mb),
      historyTtlSeconds: Int64(3 * 86400),
      cacheTtlSeconds: Int64(86400),
      syncText: true,
      syncHtml: true,
      syncImage: true,
      syncFile: true,
      autoApplyToClipboard: true,
    ),
    self: Device(id: 'self', name: 'MacBook Pro', platform: 'darwin', publicKeyFingerprint: 'H3QD-7KXA-M2PV'),
    peers: empty ? const [] : [
      Device(
        id: 'air',
        name: 'MacBook Air',
        platform: 'darwin',
        online: true,
        connection: ConnectionKind.CONNECTION_KIND_DIRECT,
        publicKeyFingerprint: 'R8NF-2WTC-QL5J',
      ),
      Device(
        id: 'mini',
        name: 'Mac mini',
        platform: 'darwin',
        online: true,
        connection: ConnectionKind.CONNECTION_KIND_RELAY,
        publicKeyFingerprint: 'B6YE-X4MH-9DKS',
      ),
      Device(id: 'imac', name: 'iMac', platform: 'darwin', publicKeyFingerprint: 'T2VA-LP7Q-8ZRN'),
    ],
    records: empty ? const [] : [
      _clip('1', ClipKind.CLIP_KIND_TEXT,
          preview: 'git push origin release/1.0 && gh release create v1.0.0',
          outgoing: true, size: 54, ago: const Duration(seconds: 4)),
      _clip('2', ClipKind.CLIP_KIND_IMAGE,
          preview: _t('图片 2.4 MB', 'Image 2.4 MB'), size: (2.4 * _mb).round(), ago: const Duration(minutes: 2),
          items: [ClipItem(name: 'image.png', size: Int64((2.4 * _mb).round()))]),
      _clip('3', ClipKind.CLIP_KIND_FILE,
          size: (38.2 * _mb).round(), ago: const Duration(minutes: 9),
          items: [ClipItem(name: _quarterly, size: Int64((38.2 * _mb).round()))]),
      _clip('4', ClipKind.CLIP_KIND_FILE,
          size: (1.3 * 1024 * _mb).round(), ago: const Duration(minutes: 14),
          from: 'Mac mini', status: ClipStatus.CLIP_STATUS_REMOTE_ONLY,
          items: [ClipItem(name: _t('发布演示录屏.mov', 'launch-demo.mov'), size: Int64((1.3 * 1024 * _mb).round()))]),
      _clip('5', ClipKind.CLIP_KIND_HTML,
          preview: _t('会议结论：周二灰度 10%，观察两天没问题周四全量；回滚预案见文档第 3 节',
              'Decision: ship to 10% on Tuesday, everyone on Thursday if two days look clean. Rollback plan in section 3'),
          from: 'Mac mini', size: 2310, ago: const Duration(minutes: 26)),
      _clip('6', ClipKind.CLIP_KIND_FILE,
          outgoing: true, size: 126 * _mb, ago: const Duration(minutes: 41),
          items: [ClipItem(name: 'design-assets', size: Int64(126 * _mb), isDir: true)]),
      _clip('7', ClipKind.CLIP_KIND_FILE,
          size: (18.6 * _mb).round(), ago: const Duration(hours: 1),
          status: ClipStatus.CLIP_STATUS_FETCHING,
          items: [ClipItem(name: _t('合同扫描件.pdf', 'signed-contract.pdf'), size: Int64((18.6 * _mb).round()))]),
      _clip('8', ClipKind.CLIP_KIND_TEXT,
          preview: 'https://github.com/baiyuze/copysync/releases', outgoing: true, size: 44,
          ago: const Duration(hours: 3)),
      _clip('9', ClipKind.CLIP_KIND_FILE,
          size: (4.1 * _mb).round(), ago: const Duration(hours: 5),
          items: [
            ClipItem(name: 'IMG_2041.HEIC', size: Int64((1.4 * _mb).round())),
            ClipItem(name: 'IMG_2042.HEIC', size: Int64((1.3 * _mb).round())),
            ClipItem(name: 'IMG_2043.HEIC', size: Int64((1.4 * _mb).round())),
          ]),
      _clip('10', ClipKind.CLIP_KIND_IMAGE,
          preview: _t('图片 864 KB', 'Image 864 KB'), outgoing: true, size: 864 * 1024, ago: const Duration(days: 1),
          items: [ClipItem(name: 'image.png', size: Int64(864 * 1024))]),
    ],
  );
  return state;
}

final _boundary = GlobalKey();

Future<void> _shoot(
  WidgetTester tester,
  String name, {
  required AppState state,
  Section section = Section.history,
  Brightness brightness = Brightness.light,
  Future<void> Function()? arrange,
  bool settle = true,
}) async {
  // 960×640 的窗口，按 Retina 2 倍渲染
  tester.view.physicalSize = const Size(1920, 1280);
  tester.view.devicePixelRatio = 2;
  addTearDown(tester.view.reset);
  // 平台覆盖必须在测试体内复原：框架在执行 tearDown 之前就会检查它
  debugDefaultTargetPlatformOverride = TargetPlatform.macOS;
  // 测试框架默认把阴影画成实心描边（便于图片比对），截图要真实的模糊阴影
  debugDisableShadows = false;
  try {
    await _render(tester, name, state: state, section: section, brightness: brightness,
        arrange: arrange, settle: settle);
  } finally {
    debugDefaultTargetPlatformOverride = null;
    debugDisableShadows = true;
  }
}

Future<void> _render(
  WidgetTester tester,
  String name, {
  required AppState state,
  required Section section,
  required Brightness brightness,
  required Future<void> Function()? arrange,
  required bool settle,
}) async {

  final theme = buildTheme(brightness, fontFamily: 'SF', fontFamilyFallback: const ['PingFang']);
  await tester.pumpWidget(RepaintBoundary(
    key: _boundary,
    child: AppScope(
      state: state,
      child: MaterialApp(
        debugShowCheckedModeBanner: false,
        theme: theme,
        locale: _locale,
        supportedLocales: supportedLocales,
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        home: Shell(initialSection: section),
      ),
    ),
  ));
  await tester.pump();
  // 图片解码是真实的异步操作，在测试的假时钟里不会自己完成
  await tester.runAsync(() async {
    for (final el in find.byType(Image).evaluate()) {
      await precacheImage((el.widget as Image).image, el);
    }
  });
  await tester.pump();
  if (arrange != null) await arrange();
  // 「拉取中」的进度条是无限动画，有它在的页面不能等动画结束
  settle ? await tester.pumpAndSettle() : await tester.pump(const Duration(milliseconds: 400));
  await _capture(tester, name);
}

/// 自己按 2 倍截取：matchesGoldenFile 只按 1 倍输出，放到 Retina 屏上会糊。
Future<void> _capture(WidgetTester tester, String name) async {
  final boundary = tester.renderObject<RenderRepaintBoundary>(find.byKey(_boundary));
  await tester.runAsync(() async {
    final image = await boundary.toImage(pixelRatio: 2);
    final png = await image.toByteData(format: ImageByteFormat.png);
    File('tool/screenshots/out/$name.png')
      ..createSync(recursive: true)
      ..writeAsBytesSync(png!.buffer.asUint8List());
  });
}

void main() {
  setUpAll(() async {
    TestWidgetsFlutterBinding.ensureInitialized();
    appL10n = lookupAppLocalizations(_locale);
    await _loadFonts();
  });

  testWidgets('history', (tester) async {
    await _shoot(tester, 'history', state: _demo(), settle: false, arrange: () async {
      // 悬停在一条记录上，展示「放入剪贴板」
      final mouse = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await mouse.addPointer(location: Offset.zero);
      addTearDown(mouse.removePointer);
      await mouse.moveTo(tester.getCenter(find.text(_quarterly)));
      await tester.pump();
    });
  });

  testWidgets('history-dark', (tester) async {
    await _shoot(tester, 'history-dark',
        state: _demo(), brightness: Brightness.dark, settle: false);
  });

  testWidgets('devices', (tester) async {
    await _shoot(tester, 'devices', state: _demo(), section: Section.devices);
  });

  testWidgets('pairing', (tester) async {
    await _shoot(tester, 'pairing', state: _demo(), section: Section.devices, arrange: () async {
      showDialog<void>(
        context: tester.element(find.byType(Shell)),
        builder: (_) => const PairingDialog(),
      );
    });
  });

  testWidgets('verify', (tester) async {
    await _shoot(tester, 'verify',
        section: Section.devices,
        state: _demo(
          pending: Device(
            id: 'new',
            name: 'MacBook Air',
            platform: 'darwin',
            publicKeyFingerprint: 'R8NF-2WTC-QL5J',
            pairingSession: 's1',
          ),
        ));
  });

  testWidgets('settings', (tester) async {
    await _shoot(tester, 'settings', state: _demo(), section: Section.settings);
  });

  // 网站上与文字并排展示的对话框：单独渲染，透明背景，保留对话框自己的圆角与阴影
  for (final dark in [false, true]) {
    testWidgets(dark ? 'card-verify-dark' : 'card-verify', (tester) async {
      tester.view.physicalSize = const Size(1120, 800);
      tester.view.devicePixelRatio = 2;
      addTearDown(tester.view.reset);
      debugDefaultTargetPlatformOverride = TargetPlatform.macOS;
      debugDisableShadows = false;
      try {
        final brightness = dark ? Brightness.dark : Brightness.light;
        await tester.pumpWidget(RepaintBoundary(
          key: _boundary,
          // 对话框要同时显示本机的指纹，从 AppScope 里取
          child: AppScope(
            state: _demo(),
            child: MaterialApp(
              debugShowCheckedModeBanner: false,
              theme: buildTheme(brightness, fontFamily: 'SF', fontFamilyFallback: const ['PingFang']),
              locale: _locale,
              supportedLocales: supportedLocales,
              localizationsDelegates: AppLocalizations.localizationsDelegates,
              home: Scaffold(
                backgroundColor: Colors.transparent,
                body: Center(
                  child: PairingConfirmDialog(
                    peer: Device(name: 'MacBook Air', publicKeyFingerprint: 'R8NF-2WTC-QL5J'),
                  ),
                ),
              ),
            ),
          ),
        ));
        await tester.pumpAndSettle();
        await _capture(tester, dark ? 'card-verify-dark' : 'card-verify');
      } finally {
        debugDefaultTargetPlatformOverride = null;
        debugDisableShadows = true;
      }
    });
  }

  testWidgets('welcome', (tester) async {
    // 首次打开时还没有任何记录和设备
    final state = _demo(link: LinkState.offline, service: ServiceState.notInstalled, empty: true);
    await _shoot(tester, 'welcome', state: state);
  });
}

// 生成 README 与网站用的界面截图：演示数据，不连接后台服务，不读真实剪贴板。
//
// 不要直接运行，用仓库根目录的 tools/screenshots/render.sh——
// 它会先从本机系统字体生成渲染用的字体文件，渲染完再加上 macOS 窗口外框。
//
// 放在 tool/ 而非 test/：它依赖本机字体，不该随 flutter test 在 CI 上跑。
// ignore_for_file: invalid_use_of_visible_for_testing_member
import 'dart:async';
import 'dart:io';
import 'dart:ui' show ImageByteFormat;

import 'package:copysync_ui/app_state.dart';
import 'package:copysync_ui/background_service.dart';
import 'package:copysync_ui/gen/copysync/v1/daemon.pb.dart';
import 'package:copysync_ui/i18n.dart';
import 'package:copysync_ui/main.dart';
import 'package:copysync_ui/theme.dart';
import 'package:copysync_ui/widgets/feedback_dialog.dart';
import 'package:copysync_ui/widgets/nearby_dialog.dart';
import 'package:copysync_ui/widgets/pairing_dialog.dart';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

/// 截图的界面语言：zh、zh-tw、zh-hk、en 或 ja，由 render.sh 设置。演示数据里的文件名、文字也跟着换。
final _lang = Platform.environment['SCREENSHOT_LANG'] ?? 'zh';
final _locale = switch (_lang) { 'zh-tw' => zhHantTW, 'zh-hk' => zhHantHK, _ => Locale(_lang) };
String _t(String zh, String en, String ja, {required String tw, required String hk}) =>
    switch (_lang) { 'en' => en, 'ja' => ja, 'zh-tw' => tw, 'zh-hk' => hk, _ => zh };

/// 汉字的字形跟着语言走：日文界面用冬青黑体，繁体用苹方的台湾、香港字形，否则用简体的苹方
final _cjk = switch (_lang) {
  'ja' => const ['Hiragino', 'PingFang'],
  'zh-tw' => const ['PingFangTC'],
  'zh-hk' => const ['PingFangHK'],
  _ => const ['PingFang'],
};

/// 生成配对码、搜索附近设备、搜索复制记录都不走后台服务，直接给固定的演示数据。
class _DemoState extends AppState {
  @override
  Future<PairingCode> createPairingCode() async => (code: 'K7MP2X', addresses: ['192.168.1.20']);

  @override
  Future<ListHistoryResponse> requestSearch(String query) async =>
      ListHistoryResponse(records: records.where((r) => matchesSearch(r, query)), searched: true);

  @override
  Future<String> exportDiagnostics({int logLines = 200}) async => _sampleReport;

  // 不关闭的流：后台服务那边的订阅一直开着，结束就表示后台服务断开了
  @override
  Stream<List<NearbyDevice>> watchNearby() => (StreamController<List<NearbyDevice>>()
        ..add([
          NearbyDevice(id: 'n1', name: _t('设计部的 PC', 'Design PC', 'デザイン部の PC', tw: '設計部的 PC', hk: '設計部的 PC'), platform: 'windows'),
          NearbyDevice(id: 'n2', name: 'MacBook Air', platform: 'darwin', paired: true),
          NearbyDevice(id: 'n3', name: _t('会议室的 iMac', 'Meeting room iMac', '会議室の iMac', tw: '會議室的 iMac', hk: '會議室的 iMac'), platform: 'darwin'),
        ]))
      .stream;
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
  await family('PingFangTC', ['PingFangTC-Regular.ttf', 'PingFangTC-Medium.ttf', 'PingFangTC-Semibold.ttf']);
  await family('PingFangHK', ['PingFangHK-Regular.ttf', 'PingFangHK-Medium.ttf', 'PingFangHK-Semibold.ttf']);
  await family('Hiragino', ['Hiragino-Regular.ttf', 'Hiragino-Medium.ttf', 'Hiragino-Semibold.ttf']);
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
final _quarterly = _t('季度汇报-终版.key', 'Q3 review - final.key', '四半期報告_最終版.key',
    tw: '季度報告-定稿.key', hk: '季度報告-最終版.key');

AppState _demo({LinkState link = LinkState.online,
    ServiceState? service, Device? pending, bool empty = false,
    String server = 'wss://sync.example.com/signal'}) {
  final state = _DemoState();
  state.debugSeed(
    link: link,
    serviceState: service,
    pendingPairing: pending,
    status: Status(
      version: '1.5.0',
      deviceName: 'MacBook Pro',
      signalingConnected: link == LinkState.online && server.isNotEmpty,
      serverConfigured: server.isNotEmpty,
      lanEnabled: true,
      clipboardPermission: ClipboardPermission.CLIPBOARD_PERMISSION_ALWAYS_ALLOW,
      cacheBytesUsed: Int64(412 * _mb),
      peersOnline: 2,
    ),
    config: Config(
      signalingUrl: server,
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
          preview: _t('图片 2.4 MB', 'Image 2.4 MB', '画像 2.4 MB', tw: '圖片 2.4 MB', hk: '圖片 2.4 MB'), size: (2.4 * _mb).round(), ago: const Duration(minutes: 2),
          items: [ClipItem(name: 'image.png', size: Int64((2.4 * _mb).round()))]),
      _clip('3', ClipKind.CLIP_KIND_FILE,
          size: (38.2 * _mb).round(), ago: const Duration(minutes: 9),
          items: [ClipItem(name: _quarterly, size: Int64((38.2 * _mb).round()))]),
      _clip('4', ClipKind.CLIP_KIND_FILE,
          size: (1.3 * 1024 * _mb).round(), ago: const Duration(minutes: 14),
          from: 'Mac mini', status: ClipStatus.CLIP_STATUS_REMOTE_ONLY,
          items: [ClipItem(name: _t('发布演示录屏.mov', 'launch-demo.mov', '発表デモ.mov', tw: '發表會示範錄影.mov', hk: '發佈會示範錄影.mov'), size: Int64((1.3 * 1024 * _mb).round()))]),
      _clip('5', ClipKind.CLIP_KIND_HTML,
          preview: _t('会议结论：周二灰度 10%，观察两天没问题周四全量；回滚预案见文档第 3 节',
              'Decision: ship to 10% on Tuesday, everyone on Thursday if two days look clean. Rollback plan in section 3',
              '決定事項：火曜に 10% へ先行公開、2 日問題なければ木曜に全体公開。切り戻し手順は資料の第 3 節',
              tw: '會議結論：週二先開放 10%，觀察兩天沒問題就在週四全面上線；回滾方案見文件第 3 節',
              hk: '會議結論：星期二先開放 10%，觀察兩日沒問題便於星期四全面推出；回滾方案見文件第 3 節'),
          from: 'Mac mini', size: 2310, ago: const Duration(minutes: 26)),
      _clip('6', ClipKind.CLIP_KIND_FILE,
          outgoing: true, size: 126 * _mb, ago: const Duration(minutes: 41),
          items: [ClipItem(name: 'design-assets', size: Int64(126 * _mb), isDir: true)]),
      _clip('7', ClipKind.CLIP_KIND_FILE,
          size: (18.6 * _mb).round(), ago: const Duration(hours: 1),
          status: ClipStatus.CLIP_STATUS_FETCHING,
          items: [ClipItem(name: _t('合同扫描件.pdf', 'signed-contract.pdf', '契約書_署名済み.pdf', tw: '合約掃描檔.pdf', hk: '合約掃描本.pdf'), size: Int64((18.6 * _mb).round()))]),
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
          preview: _t('图片 864 KB', 'Image 864 KB', '画像 864 KB', tw: '圖片 864 KB', hk: '圖片 864 KB'), outgoing: true, size: 864 * 1024, ago: const Duration(days: 1),
          items: [ClipItem(name: 'image.png', size: Int64(864 * 1024))]),
    ],
  );
  return state;
}

/// 诊断报告的样子：后台服务生成，名字、服务器、公网 IP 都已换成占位符
const _sampleReport = '''CopySync diagnostics report
Names, server address, public IP addresses, the user folder and file names are replaced.

  Generated              2026-10-13 20:50:11 +0800
  Background service     1.7.0
  App                    1.7.0
  System                 macOS 15.6 (darwin/arm64)

Settings
  Server                 ws://<server>:8787/signal
  Find devices on LAN    on
  Public STUN servers    on
  Auto-sync limit        50 MB
  History kept           3 days

Devices
  <this-device>          darwin (this device)
  <device-1>             windows  direct
  <device-2>             darwin   relay

Log (last 200 lines)
time=2026-10-13T20:41:02.118+08:00 level=INFO msg=信令已连接 url=ws://<server>:8787/signal
time=2026-10-13T20:41:02.530+08:00 level=INFO msg="P2P 连接状态" device=9f2c41 state=connected''';

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

  final theme = buildTheme(brightness, fontFamily: 'SF', fontFamilyFallback: _cjk);
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

  // 搜索 release：命令行与链接两条都含有它，看得到多条结果和标出的匹配处
  testWidgets('history-search', (tester) async {
    await _shoot(tester, 'history-search', state: _demo(), settle: false, arrange: () async {
      await tester.enterText(find.byType(TextField), 'release');
      await tester.pump(const Duration(milliseconds: 300));
      await tester.pump();
    });
  });

  testWidgets('feedback', (tester) async {
    await _shoot(tester, 'feedback', state: _demo(), section: Section.settings, arrange: () async {
      showFeedbackDialog(tester.element(find.byType(Shell)));
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

  for (final dark in [false, true]) {
    testWidgets(dark ? 'nearby-dark' : 'nearby', (tester) async {
      await _shoot(tester, dark ? 'nearby-dark' : 'nearby',
          state: _demo(),
          section: Section.devices,
          brightness: dark ? Brightness.dark : Brightness.light,
          settle: false, // 雷达的波纹一直在动
          arrange: () async {
        showDialog<void>(
          context: tester.element(find.byType(Shell)),
          builder: (_) => const NearbyDialog(),
        );
        await tester.pump();
        await tester.pump();
        await tester.pump(const Duration(milliseconds: 300));
        await tester.pump(const Duration(milliseconds: 900));
      });
    });
  }

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

  // 网站「在局域网里部署服务器」：服务器地址填局域网 IP，已连接
  testWidgets('settings-lan-dark', (tester) async {
    await _shoot(tester, 'settings-lan-dark',
        state: _demo(server: 'ws://192.168.1.20:8787/signal'),
        section: Section.settings, brightness: Brightness.dark);
  });

  // 网站「部署到公网服务器」：服务器地址填公网 IP（文档专用的示例地址段）
  testWidgets('settings-cloud-dark', (tester) async {
    await _shoot(tester, 'settings-cloud-dark',
        state: _demo(server: 'ws://203.0.113.10:8787/signal'),
        section: Section.settings, brightness: Brightness.dark);
  });

  // 网站「快速开始」：第一次打开、生成配对码，深色界面配深色网页
  testWidgets('welcome-dark', (tester) async {
    final state = _demo(link: LinkState.offline, service: ServiceState.notInstalled, empty: true);
    await _shoot(tester, 'welcome-dark', state: state, brightness: Brightness.dark);
  });

  testWidgets('pairing-dark', (tester) async {
    await _shoot(tester, 'pairing-dark', state: _demo(), section: Section.devices,
        brightness: Brightness.dark, arrange: () async {
      showDialog<void>(
        context: tester.element(find.byType(Shell)),
        builder: (_) => const PairingDialog(),
      );
    });
  });

  // 网站「同一网络」的跨网段一节与「Tailscale」：输入配对码，展开地址框，填上对方显示的地址
  for (final (name, address) in [('pairing-join-dark', '192.168.1.20'), ('pairing-join-ts-dark', '100.92.41.17')]) {
    testWidgets(name, (tester) async {
      await _shoot(tester, name, state: _demo(), section: Section.devices,
          brightness: Brightness.dark, arrange: () async {
        showDialog<void>(
          context: tester.element(find.byType(Shell)),
          builder: (_) => const PairingDialog(),
        );
        await tester.pumpAndSettle();
        await tester.tap(find.text(appL10n.enterCode));
        await tester.pumpAndSettle();
        await tester.enterText(find.byType(TextField).last, 'K7MP2X');
        await tester.tap(find.text(appL10n.differentSubnet));
        await tester.pumpAndSettle();
        await tester.enterText(find.byType(TextField).last, address);
      });
    });
  }

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
              theme: buildTheme(brightness, fontFamily: 'SF', fontFamilyFallback: _cjk),
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

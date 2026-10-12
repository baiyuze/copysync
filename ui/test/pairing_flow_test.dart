import 'dart:async';

import 'package:copysync_ui/app_state.dart';
import 'package:copysync_ui/gen/copysync/v1/daemon.pb.dart';
import 'package:copysync_ui/main.dart';
import 'package:copysync_ui/widgets/nearby_dialog.dart';
import 'package:copysync_ui/widgets/pairing_dialog.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:grpc/grpc.dart';

Device pendingPeer([String session = 'session-1']) => Device(
  id: 'peer',
  name: '另一台电脑',
  publicKeyFingerprint: 'BBBB-CCCC-DDDD',
  pairingSession: session,
);

class PairingState extends AppState {
  final redemption = Completer<Device>();
  final confirmations = <(String, bool)>[];
  final redeemedWith = <(String, String)>[];
  Object? nextRedeemFailure;
  bool failConfirmation = false;

  bool nearbyCancelled = false;
  late final nearby = StreamController<List<NearbyDevice>>(onCancel: () => nearbyCancelled = true);
  final pairedNearby = <String>[];
  Object? nextNearbyFailure;

  // 没有服务器、只开着局域网：这是新装的默认状态，照样能配对
  PairingState() {
    debugSeed(
      status: Status(lanEnabled: true),
      config: Config(),
      self: Device(
        id: 'self',
        name: '本机',
        publicKeyFingerprint: 'AAAA-BBBB-CCCC',
      ),
    );
  }

  void receive(Device peer) => debugEvent(Event(deviceChanged: peer));

  @override
  Future<PairingCode> createPairingCode() async => (code: 'ABC123', addresses: ['192.168.1.20']);

  @override
  Future<Device> redeemPairingCode(String code, {String address = ''}) async {
    redeemedWith.add((code, address));
    final failure = nextRedeemFailure;
    if (failure != null) {
      nextRedeemFailure = null;
      throw failure;
    }
    final peer = await redemption.future;
    receive(peer);
    return peer;
  }

  @override
  Stream<List<NearbyDevice>> watchNearby() => nearby.stream;

  @override
  Future<Device> pairNearby(String id) async {
    pairedNearby.add(id);
    final failure = nextNearbyFailure;
    if (failure != null) {
      nextNearbyFailure = null;
      throw failure;
    }
    final peer = pendingPeer();
    receive(peer);
    return peer;
  }

  @override
  Future<void> confirmPairing(String session, bool accept) async {
    confirmations.add((session, accept));
    if (failConfirmation) throw StateError('配对会话已过期');
    dismissPendingPairing(session);
    if (accept) receive(pendingPeer()..clearPairingSession());
  }
}

Future<void> openSetup(WidgetTester tester, PairingState state) async {
  await tester.pumpWidget(CopySyncApp(state: state));
  await tester.tap(find.text('设备').first);
  await tester.pumpAndSettle();
  await tester.tap(find.text('添加设备').first);
  await tester.pumpAndSettle();
}

/// 打开「搜索附近设备」。雷达的波纹一直在动，不能用 pumpAndSettle 等它停下。
Future<void> openNearby(WidgetTester tester, PairingState state) async {
  await tester.pumpWidget(CopySyncApp(state: state));
  await tester.tap(find.text('设备').first);
  await tester.pumpAndSettle();
  await tester.tap(find.text('搜索附近设备').first);
  await tester.pump();
  await tester.pump(const Duration(milliseconds: 300));
}

final nearbyDevices = [
  NearbyDevice(id: 'copysync-1', name: '光光的电脑', platform: 'windows'),
  NearbyDevice(id: 'copysync-2', name: '旧 Mac', platform: 'darwin', paired: true),
];

void expectSingleConfirmation() {
  // 包括被第二个对话框盖住的路由，避免只检查最上层而漏掉叠加弹窗。
  expect(
    find.byType(PairingConfirmDialog, skipOffstage: false),
    findsOneWidget,
  );
  expect(find.byType(PairingDialog, skipOffstage: false), findsNothing);
}

void main() {
  for (final eventFirst in [false, true]) {
    testWidgets('输入配对码只确认一次，${eventFirst ? "事件先到、RPC 迟到" : "RPC 先到、事件重复"}', (
      tester,
    ) async {
      final state = PairingState();
      addTearDown(state.dispose);
      await openSetup(tester, state);
      await tester.tap(find.text('输入配对码'));
      await tester.pumpAndSettle();
      await tester.enterText(find.byType(TextField), 'ABC123');
      await tester.tap(find.text('继续'));
      if (eventFirst) {
        state.receive(pendingPeer());
      } else {
        state.redemption.complete(pendingPeer());
      }
      await tester.pumpAndSettle();
      expectSingleConfirmation();

      state.receive(pendingPeer());
      await tester.pumpAndSettle();
      expectSingleConfirmation();
      await tester.tap(find.text('一致，确认配对'));
      await tester.pumpAndSettle();

      // 用户可能在兑换 RPC 返回前就通过事件弹出的确认框完成了配对。
      if (eventFirst) state.redemption.complete(pendingPeer());
      state.receive(pendingPeer());
      await tester.pumpAndSettle();
      expect(state.confirmations, [('session-1', true)]);
      expect(state.pendingPairing, isNull);
      expect(
        find.byType(PairingConfirmDialog, skipOffstage: false),
        findsNothing,
      );
      expect(find.byType(PairingDialog, skipOffstage: false), findsNothing);
      expect(tester.takeException(), isNull);
    });
  }

  for (final accept in [false, true]) {
    testWidgets('生成配对码后${accept ? "接受" : "拒绝"}请求，移除旧页面且不再次弹窗', (tester) async {
      final state = PairingState();
      addTearDown(state.dispose);
      await openSetup(tester, state);
      state.receive(pendingPeer());
      await tester.pumpAndSettle();
      expectSingleConfirmation();
      await tester.tap(find.text(accept ? '一致，确认配对' : '不一致，拒绝'));
      await tester.pumpAndSettle();
      state.receive(pendingPeer());
      await tester.pumpAndSettle();
      expect(state.confirmations, [('session-1', accept)]);
      expect(
        find.byType(PairingConfirmDialog, skipOffstage: false),
        findsNothing,
      );
      expect(find.byType(PairingDialog, skipOffstage: false), findsNothing);
      expect(tester.takeException(), isNull);
    });
  }

  testWidgets('未打开添加设备时收到请求仍主动确认，过期后不会反复弹出', (tester) async {
    final state = PairingState()..failConfirmation = true;
    addTearDown(state.dispose);
    await tester.pumpWidget(CopySyncApp(state: state));
    state.receive(pendingPeer());
    await tester.pumpAndSettle();
    expectSingleConfirmation();
    await tester.tap(find.text('一致，确认配对'));
    await tester.pumpAndSettle();
    state.receive(pendingPeer());
    await tester.pumpAndSettle();
    expect(
      find.byType(PairingConfirmDialog, skipOffstage: false),
      findsNothing,
    );
    expect(state.pendingPairing, isNull);
    expect(tester.takeException(), isNull);
  });

  test('处理旧会话不清掉新的配对请求', () {
    final state = PairingState();
    addTearDown(state.dispose);
    state.receive(pendingPeer());
    state.receive(pendingPeer('session-2'));
    state.dismissPendingPairing('session-1');
    state.receive(pendingPeer());
    expect(state.pendingPairing?.pairingSession, 'session-2');
  });

  testWidgets('生成配对码时显示本机地址，跨网段时对方要输入它', (tester) async {
    final state = PairingState();
    addTearDown(state.dispose);
    await openSetup(tester, state);
    expect(find.text('ABC123'), findsOneWidget);
    expect(find.textContaining('192.168.1.20', findRichText: true), findsOneWidget);
    expect(find.text('不在同一网段时，另一台还要输入这个地址'), findsOneWidget);
  });

  testWidgets('局域网里找不到对方时展开地址框，带上地址重试', (tester) async {
    final state = PairingState()
      ..nextRedeemFailure = GrpcError.notFound('not_found: 没在局域网里找到正在配对的设备');
    addTearDown(state.dispose);
    await openSetup(tester, state);
    await tester.tap(find.text('输入配对码'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextField), 'abc123');
    await tester.tap(find.text('继续'));
    await tester.pumpAndSettle();

    expect(find.text('没找到那台设备。不在同一网段时，请输入它配对码下面显示的地址。'), findsOneWidget);
    final fields = find.byType(TextField);
    expect(fields, findsNWidgets(2));
    await tester.enterText(fields.last, '192.168.1.20');
    await tester.tap(find.text('继续'));
    state.redemption.complete(pendingPeer());
    await tester.pumpAndSettle();

    expect(state.redeemedWith, [('ABC123', ''), ('ABC123', '192.168.1.20')]);
    expectSingleConfirmation();
  });

  testWidgets('配对码不对时给出本地化的原因', (tester) async {
    final state = PairingState()
      ..nextRedeemFailure = GrpcError.invalidArgument('wrong_code: 配对码不对');
    addTearDown(state.dispose);
    await openSetup(tester, state);
    await tester.tap(find.text('输入配对码'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextField), 'ABC123');
    await tester.tap(find.text('继续'));
    await tester.pumpAndSettle();
    expect(find.text('配对码不对，或者已经失效'), findsOneWidget);
    // 不是「找不到」，地址框保持收起
    expect(find.byType(TextField), findsOneWidget);
  });

  group('搜索附近设备', () {
    testWidgets('原来的添加设备不变，雷达是单独的入口', (tester) async {
      final state = PairingState();
      addTearDown(state.dispose);
      await openSetup(tester, state);
      expect(find.byType(PairingDialog), findsOneWidget);
      expect(find.byType(NearbyDialog), findsNothing);
      expect(find.text('ABC123'), findsOneWidget);
      expect(state.nearby.hasListener, isFalse, reason: '没打开雷达时不该公开本机');
    });

    testWidgets('显示附近的设备，点一下弹出核对指纹并关掉雷达', (tester) async {
      final state = PairingState();
      addTearDown(state.dispose);
      await openNearby(tester, state);
      expect(find.byType(NearbyDialog), findsOneWidget);
      expect(find.text('正在查找附近的设备…'), findsOneWidget);
      expect(state.nearby.hasListener, isTrue);

      state.nearby.add(nearbyDevices);
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 300));
      expect(find.text('光光的电脑'), findsOneWidget);
      expect(find.text('已配对'), findsOneWidget);
      expect(find.text('正在查找附近的设备…'), findsNothing);

      // 已配对的点了没反应
      await tester.tap(find.text('旧 Mac'));
      await tester.pump();
      expect(state.pairedNearby, isEmpty);

      await tester.tap(find.text('光光的电脑'));
      await tester.pumpAndSettle();
      expect(state.pairedNearby, ['copysync-1']);
      expectSingleConfirmation();
      expect(find.byType(NearbyDialog, skipOffstage: false), findsNothing);
      expect(state.nearbyCancelled, isTrue, reason: '雷达关掉后要撤下本机');
      expect(tester.takeException(), isNull);
    });

    testWidgets('对方没有回应时给出原因，关掉窗口就撤下本机', (tester) async {
      final state = PairingState()
        ..nextNearbyFailure = GrpcError.notFound('not_found: 没在局域网里找到正在配对的设备');
      addTearDown(state.dispose);
      await openNearby(tester, state);
      state.nearby.add(nearbyDevices);
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 300));

      await tester.tap(find.text('光光的电脑'));
      await tester.pump();
      await tester.pump();
      expect(find.text('「光光的电脑」没有回应。确认它还开着「搜索附近设备」。'), findsOneWidget);
      expect(find.byType(NearbyDialog), findsOneWidget);

      await tester.tap(find.byTooltip('关闭'));
      await tester.pumpAndSettle();
      expect(find.byType(NearbyDialog), findsNothing);
      expect(state.nearbyCancelled, isTrue);
    });

    testWidgets('局域网发现关闭时说明原因', (tester) async {
      final state = PairingState();
      addTearDown(state.dispose);
      state.debugSeed(status: Status(), config: Config(), self: Device(id: 'self', name: '本机'));
      await openNearby(tester, state);
      expect(find.text('局域网发现已关闭'), findsOneWidget);
      expect(find.text('光光的电脑'), findsNothing);
    });
  });

  group('侧边栏的连接状态', () {
    Future<void> pumpWith(WidgetTester tester, Status status) async {
      final state = PairingState();
      addTearDown(state.dispose);
      state.debugSeed(status: status, config: Config(), self: Device(id: 'self', name: '本机'));
      await tester.pumpWidget(CopySyncApp(state: state));
      await tester.pumpAndSettle();
    }

    testWidgets('没设置服务器时是局域网模式，不算异常', (tester) async {
      await pumpWith(tester, Status(lanEnabled: true));
      expect(find.text('在局域网中查找设备'), findsOneWidget);
    });

    testWidgets('设置了服务器但连不上时提示仍在局域网里同步', (tester) async {
      await pumpWith(tester, Status(lanEnabled: true, serverConfigured: true));
      expect(find.text('连不上服务器，局域网内照常同步'), findsOneWidget);
    });

    testWidgets('macOS 拦下局域网时在设备页给出处理办法', (tester) async {
      await pumpWith(tester, Status(lanEnabled: true, lanBlocked: true));
      await tester.tap(find.text('设备').first);
      await tester.pumpAndSettle();
      expect(find.text('macOS 阻止了局域网访问'), findsWidgets);
      expect(find.text('打开系统设置'), findsOneWidget);
    });
  });
}

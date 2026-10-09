import 'dart:async';

import 'package:copysync_ui/app_state.dart';
import 'package:copysync_ui/gen/copysync/v1/daemon.pb.dart';
import 'package:copysync_ui/main.dart';
import 'package:copysync_ui/widgets/pairing_dialog.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

Device pendingPeer([String session = 'session-1']) => Device(
  id: 'peer',
  name: '另一台电脑',
  publicKeyFingerprint: 'BBBB-CCCC-DDDD',
  pairingSession: session,
);

class PairingState extends AppState {
  final redemption = Completer<Device>();
  final confirmations = <(String, bool)>[];
  bool failConfirmation = false;

  PairingState() {
    debugSeed(
      status: Status(signalingConnected: true),
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
  Future<String> createPairingCode() async => 'ABC123';

  @override
  Future<Device> redeemPairingCode(String code) async {
    final peer = await redemption.future;
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
}

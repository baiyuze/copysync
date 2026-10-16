import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:pasteferry_ui/app_state.dart';
import 'package:pasteferry_ui/gen/pasteferry/v1/daemon.pb.dart';
import 'package:pasteferry_ui/main.dart';
import 'package:pasteferry_ui/pages/official_service_page.dart';
import 'package:pasteferry_ui/pages/settings_page.dart';

class OfficialTestState extends AppState {
  bool failSave = false;
  int saves = 0;
  OfficialTestState({
    bool official = false,
    String language = 'zh-Hans',
    String ack = '',
  }) {
    debugSeed(
      status: Status(deviceId: 'abcdefghijklmnopqrstuvwxyz', lanEnabled: true),
      config: Config(
        serviceMode: official
            ? ServiceMode.SERVICE_MODE_OFFICIAL
            : ServiceMode.SERVICE_MODE_LAN,
        language: language,
        officialNoticeAck: ack,
      ),
      self: Device(id: 'abcdefghijklmnopqrstuvwxyz'),
    );
  }
  @override
  Future<void> updateConfig(Config next) async {
    saves++;
    if (failSave) throw StateError('save failed');
    debugSeed(status: status!, config: next, self: self!);
  }

  @override
  Future<void> acknowledgeOfficialNotice(String key) async =>
      updateConfig(config!.deepCopy()..officialNoticeAck = key);
  void usage(OfficialUsage u) =>
      debugEvent(Event(statusChanged: status!.deepCopy()..officialUsage = u));
}

void main() {
  testWidgets(
    'official navigation is opt-in and does not replace startup history',
    (tester) async {
      final state = OfficialTestState();
      await tester.pumpWidget(PasteFerryApp(state: state));
      expect(find.text('官方服务'), findsNothing);
      state.debugSeed(
        status: state.status!,
        config: state.config!.deepCopy()
          ..serviceMode = ServiceMode.SERVICE_MODE_OFFICIAL,
        self: state.self!,
      );
      await tester.pumpAndSettle();
      expect(find.text('还没有记录'), findsOneWidget);
      await tester.tap(find.text('官方服务').first);
      await tester.pumpAndSettle();
      expect(find.byType(OfficialServicePage), findsOneWidget);
      expect(find.text('尚未申请中转，余额待核验'), findsWidgets);
      expect(find.text('0.00 MB'), findsNothing);
      state.debugSeed(
        status: state.status!,
        config: state.config!.deepCopy()
          ..serviceMode = ServiceMode.SERVICE_MODE_LAN,
        self: state.self!,
      );
      await tester.pumpAndSettle();
      expect(find.byType(OfficialServicePage), findsNothing);
      expect(find.byType(SettingsPage), findsOneWidget);
    },
  );
  testWidgets(
    'custom address is saved explicitly and a failure preserves the draft',
    (tester) async {
      final state = OfficialTestState()..failSave = true;
      await tester.pumpWidget(PasteFerryApp(state: state));
      await tester.tap(find.text('设置').first);
      await tester.pumpAndSettle();
      final input = find.byType(TextField).first;
      await tester.enterText(input, 'wss://my.example/signal');
      expect(state.saves, 0);
      await tester.tap(find.text('保存并连接'));
      await tester.pumpAndSettle();
      expect(state.config!.signalingUrl, isEmpty);
      expect(
        tester.widget<TextField>(input).controller!.text,
        'wss://my.example/signal',
      );
      state.failSave = false;
      await tester.tap(find.text('保存并连接'));
      await tester.pumpAndSettle();
      expect(state.config!.customSignalingUrl, 'wss://my.example/signal');
      expect(state.config!.serviceMode, ServiceMode.SERVICE_MODE_CUSTOM);
    },
  );
  testWidgets(
    'quota notice waits for foreground and acknowledgement survives restart',
    (tester) async {
      final state = OfficialTestState(official: true);
      await tester.pumpWidget(PasteFerryApp(state: state));
      tester.binding.handleAppLifecycleStateChanged(AppLifecycleState.inactive);
      final u = OfficialUsage(
        available: true,
        day: '2026-10-16',
        quotaGroup: 'bucket',
        reason: 'quota_exhausted',
        limitBytes: Int64(100000000),
      );
      state.usage(u);
      await tester.pumpAndSettle();
      expect(find.byType(AlertDialog), findsNothing);
      tester.binding.handleAppLifecycleStateChanged(AppLifecycleState.resumed);
      await tester.pumpAndSettle();
      expect(find.text('今日中转额度已用完'), findsOneWidget);
      await tester.tap(find.text('知道了'));
      await tester.pumpAndSettle();
      expect(
        state.config!.officialNoticeAck,
        contains('bucket:2026-10-16:quota_exhausted'),
      );
      state.usage(u);
      await tester.pumpAndSettle();
      expect(find.byType(AlertDialog), findsNothing);
      final restarted = OfficialTestState(
        official: true,
        ack: state.config!.officialNoticeAck,
      )..usage(u);
      await tester.pumpWidget(
        PasteFerryApp(key: const ValueKey('restart'), state: restarted),
      );
      await tester.pumpAndSettle();
      expect(find.byType(AlertDialog), findsNothing);
    },
  );
  for (final language in ['en', 'zh-Hans', 'zh-Hant-TW', 'zh-Hant-HK', 'ja']) {
    testWidgets('official usage fits minimum window in $language', (
      tester,
    ) async {
      tester.view.physicalSize = const Size(760, 560);
      tester.view.devicePixelRatio = 1;
      addTearDown(tester.view.resetPhysicalSize);
      addTearDown(tester.view.resetDevicePixelRatio);
      final state = OfficialTestState(official: true, language: language);
      final now = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      state.usage(
        OfficialUsage(
          available: true,
          day: '2026-10-16',
          quotaGroup: 'bucket',
          uploadBytes: Int64(90001234),
          downloadBytes: Int64(456789),
          limitBytes: Int64(100000000),
          remainingBytes: Int64(9541977),
          resetAtUnix: Int64(now + 3600),
          updatedAtUnix: Int64(now),
        ),
      );
      await tester.pumpWidget(PasteFerryApp(state: state));
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.cloud_outlined).first);
      await tester.pumpAndSettle();
      expect(find.byType(OfficialServicePage), findsOneWidget);
      expect(find.text('90.46 MB'), findsOneWidget);
      expect(find.text('100.00 MB'), findsOneWidget);
      expect(find.text('9.54 MB'), findsOneWidget);
      expect(find.text('90.00 MB'), findsOneWidget);
      expect(find.text('0.46 MB'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  }
}

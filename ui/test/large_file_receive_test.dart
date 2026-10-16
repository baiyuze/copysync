import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:pasteferry_ui/app_state.dart';
import 'package:pasteferry_ui/gen/pasteferry/v1/daemon.pb.dart';
import 'package:pasteferry_ui/main.dart';

class _ReceiveState extends AppState {
  final received = <String>[];
  @override
  Future<void> fetch(String id) async { received.add(id); }
}

_ReceiveState _seed(ClipStatus status, {String language = 'en'}) => _ReceiveState()
  ..debugSeed(status: Status(), config: Config(language: language), self: Device(id: 'self'), records: [
    ClipRecord(id: 'large', kind: ClipKind.CLIP_KIND_FILE, status: status,
      items: [ClipItem(name: 'archive.zip')], totalSize: Int64(80 << 20), originDeviceName: 'Office PC',
      error: status == ClipStatus.CLIP_STATUS_FAILED ? 'Source is offline' : ''),
  ]);

void main() {
  testWidgets('pending item explains why paste is unavailable and offers receive', (tester) async {
    final state = _seed(ClipStatus.CLIP_STATUS_REMOTE_ONLY);
    await tester.pumpWidget(PasteFerryApp(state: state));
    expect(find.text('Not received — receive before pasting'), findsOneWidget);
    await tester.tap(find.widgetWithText(OutlinedButton, 'Receive'));
    await tester.pump();
    expect(state.received, ['large']);
  });

  testWidgets('failed receive keeps the reason and a retry action', (tester) async {
    final state = _seed(ClipStatus.CLIP_STATUS_FAILED);
    await tester.pumpWidget(PasteFerryApp(state: state));
    expect(find.text('Source is offline'), findsOneWidget);
    await tester.tap(find.widgetWithText(OutlinedButton, 'Try again'));
    await tester.pump();
    expect(state.received, ['large']);
  });

  testWidgets('downloaded item preserves clipboard result and exposes apply without hover', (tester) async {
    final state = _seed(ClipStatus.CLIP_STATUS_READY);
    await tester.pumpWidget(PasteFerryApp(state: state));
    state.debugEvent(Event(clipboardWritten: ClipboardWriteResult(clipId: 'large',
      outcome: ClipboardWriteOutcome.CLIPBOARD_WRITE_OUTCOME_PRESERVED)));
    await tester.pump();
    expect(find.text('Received; use Copy to clipboard when ready'), findsOneWidget);
    expect(find.widgetWithText(TextButton, 'Copy to clipboard').hitTestable(), findsOneWidget);
    state.debugEvent(Event(clipboardWritten: ClipboardWriteResult(clipId: 'large',
      outcome: ClipboardWriteOutcome.CLIPBOARD_WRITE_OUTCOME_FAILED, error: 'Clipboard is busy')));
    await tester.pump();
    expect(find.text('Received, but the clipboard could not be updated'), findsOneWidget);
    expect(state.records.single.status, ClipStatus.CLIP_STATUS_READY);
  });
}

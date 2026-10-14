import 'dart:io';

import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:pasteferry_ui/app_state.dart';
import 'package:pasteferry_ui/gen/pasteferry/v1/daemon.pb.dart';
import 'package:pasteferry_ui/main.dart';
import 'package:pasteferry_ui/pages/history_page.dart';

/// 不连后台服务：搜索请求记下来，按测试给的规则回答。
class _SearchState extends AppState {
  _SearchState(this.respond);

  final ListHistoryResponse Function(String query) respond;
  final queries = <String>[];

  @override
  Future<ListHistoryResponse> requestSearch(String query) async {
    queries.add(query);
    return respond(query);
  }
}

final _report = ClipRecord(
  id: 'report',
  kind: ClipKind.CLIP_KIND_FILE,
  originDeviceName: 'MacBook Air',
  items: [ClipItem(name: 'notes.txt'), ClipItem(name: 'Quarterly report.key')],
  totalSize: Int64(2048),
);
final _note = ClipRecord(id: 'note', kind: ClipKind.CLIP_KIND_TEXT, textPreview: 'Lunch at noon', outgoing: true);

_SearchState _seed(ListHistoryResponse Function(String) respond) => _SearchState(respond)
  ..debugSeed(
    status: Status(signalingConnected: true, version: '1.7.0'),
    config: Config(language: 'en', historyTtlSeconds: Int64(3 * 86400)),
    self: Device(id: 'self', name: 'MacBook Pro'),
    records: [_note, _report],
  );

Future<void> _type(WidgetTester tester, String text) async {
  await tester.enterText(find.byType(TextField), text);
  await tester.pump(const Duration(milliseconds: 250)); // 过了停顿时间才查询
  await tester.pump();
}

void main() {
  testWidgets('停顿后才向后台服务查询，显示它返回的结果并标出匹配处', (tester) async {
    final state = _seed((q) => ListHistoryResponse(records: [_report], searched: true));
    await tester.pumpWidget(PasteFerryApp(state: state));

    await tester.enterText(find.byType(TextField), 'quar');
    await tester.pump(const Duration(milliseconds: 100));
    await tester.enterText(find.byType(TextField), 'quarterly');
    await tester.pump(const Duration(milliseconds: 100));
    expect(state.queries, isEmpty);
    await tester.pump(const Duration(milliseconds: 150));
    await tester.pump();
    expect(state.queries, ['quarterly']);

    // 搜到的不是第一个文件时，标题用搜到的那个打头
    expect(find.text('Quarterly report.key and 1 more', findRichText: true), findsOneWidget);
    expect(find.text('Lunch at noon', findRichText: true), findsNothing);
    expect(find.text('1 match'), findsOneWidget);
    final title = tester.widget<RichText>(find.text('Quarterly report.key and 1 more', findRichText: true));
    final marked = <String>[];
    title.text.visitChildren((span) {
      if (span is TextSpan && span.style?.backgroundColor != null) marked.add(span.text!);
      return true;
    });
    expect(marked, ['Quarterly']);
    // 搜索时不放「清空记录」，免得以为只清空搜到的这些
    expect(find.byTooltip('Clear history'), findsNothing);
  });

  testWidgets('没有结果时说明搜的是什么，清空搜索词后回到完整列表', (tester) async {
    final state = _seed((q) => ListHistoryResponse(searched: true));
    await tester.pumpWidget(PasteFerryApp(state: state));

    await _type(tester, 'zebra');
    expect(find.text('No matches for “zebra”'), findsOneWidget);
    expect(find.text('Search looks through copied text and the names of copied files.'), findsOneWidget);

    await tester.tap(find.byTooltip('Clear search'));
    await tester.pump();
    expect(find.text('Lunch at noon', findRichText: true), findsOneWidget);
    expect(find.byTooltip('Clear history'), findsOneWidget);
  });

  testWidgets('旧版后台服务不认识搜索词时，退回到在已加载的记录里找', (tester) async {
    // 旧版照常返回全部记录，没有 searched
    final state = _seed((q) => ListHistoryResponse(records: [_note, _report]));
    await tester.pumpWidget(PasteFerryApp(state: state));

    await _type(tester, 'LUNCH');
    expect(find.text('Lunch at noon', findRichText: true), findsOneWidget);
    expect(find.text('notes.txt and 1 more', findRichText: true), findsNothing);
  });

  testWidgets('快捷键跳到搜索框，Esc 清空搜索词', (tester) async {
    final state = _seed((q) => ListHistoryResponse(records: [_note], searched: true));
    await tester.pumpWidget(PasteFerryApp(state: state));
    await tester.pump();

    final modifier = Platform.isMacOS ? LogicalKeyboardKey.metaLeft : LogicalKeyboardKey.controlLeft;
    await tester.sendKeyDownEvent(modifier);
    await tester.sendKeyEvent(LogicalKeyboardKey.keyF);
    await tester.sendKeyUpEvent(modifier);
    await tester.pump();
    final field = tester.widget<TextField>(find.byType(TextField));
    expect(field.focusNode!.hasFocus, isTrue);

    await _type(tester, 'lunch');
    expect(state.searching, isTrue);
    await tester.sendKeyEvent(LogicalKeyboardKey.escape);
    await tester.pump();
    expect(state.searching, isFalse);
    expect(field.controller!.text, isEmpty);
  });

  test('标出所有匹配处，不分大小写', () {
    const mark = TextStyle(backgroundColor: Color(0xFF000000));
    final spans = highlightMatches('Copy, copy, COPY', 'copy', mark);
    expect(spans.where((s) => s.style == mark).map((s) => s.text), ['Copy', 'copy', 'COPY']);
    expect(spans.map((s) => s.text).join(), 'Copy, copy, COPY');
    expect(highlightMatches('abc', '', mark).single.text, 'abc');
  });

  test('界面里的匹配：预览片段与文件名，剪贴板图片的文件名不算', () {
    expect(matchesSearch(_note, 'NOON'), isTrue);
    expect(matchesSearch(_report, 'quarterly'), isTrue);
    expect(matchesSearch(ClipRecord(kind: ClipKind.CLIP_KIND_IMAGE, items: [ClipItem(name: 'image.png')]), 'image'),
        isFalse);
  });
}

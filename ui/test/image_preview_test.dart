import 'dart:io';

import 'package:copysync_ui/app_state.dart';
import 'package:copysync_ui/gen/copysync/v1/daemon.pb.dart';
import 'package:copysync_ui/main.dart';
import 'package:copysync_ui/widgets/image_preview_dialog.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

class PreviewState extends AppState {
  PreviewState(this.path, {ClipStatus status = ClipStatus.CLIP_STATUS_READY}) {
    debugSeed(
      status: Status(signalingConnected: true),
      config: Config(),
      self: Device(id: 'self'),
      records: [
        ClipRecord(
          id: 'image',
          kind: ClipKind.CLIP_KIND_IMAGE,
          status: status,
          textPreview: '测试图片',
          originDeviceName: '测试电脑',
        ),
      ],
    );
  }
  String path;
  int previewCalls = 0;
  int fetchCalls = 0;
  int clipboardCalls = 0;
  bool failPreview = false;
  bool failFetch = false;

  @override
  Future<String> imagePreviewPath(String clipId) async {
    previewCalls++;
    if (failPreview) throw StateError('图片已不在本机，可能已被清理');
    return path;
  }

  @override
  Future<void> applyToClipboard(String clipId) async {
    clipboardCalls++;
  }

  @override
  Future<void> fetchImagePreview(String clipId) async {
    fetchCalls++;
    if (failFetch) throw StateError('对方设备暂时离线');
    setStatus(ClipStatus.CLIP_STATUS_FETCHING);
  }

  void setStatus(ClipStatus status) => debugEvent(
    Event(clipUpdated: records.single.deepCopy()..status = status),
  );
}

Future<void> openPreview(
  WidgetTester tester,
  PreviewState state, {
  bool button = false,
}) async {
  await tester.pumpWidget(CopySyncApp(state: state));
  await tester.pump();
  await tester.tap(button ? find.text('预览') : find.textContaining('来自 测试电脑'));
  await tester.pump(const Duration(milliseconds: 250));
  expect(find.byType(ImagePreviewDialog), findsOneWidget);
}

Future<RawImage> decodedImage(WidgetTester tester) async {
  final finder = find.descendant(
    of: find.byKey(const ValueKey('image-preview-content')),
    matching: find.byType(RawImage),
  );
  for (var i = 0; i < 60; i++) {
    await tester.runAsync(
      () => Future<void>.delayed(const Duration(milliseconds: 50)),
    );
    await tester.pump();
    if (finder.evaluate().isNotEmpty) {
      final image = tester.widget<RawImage>(finder);
      if (image.image != null) return image;
    }
  }
  throw StateError('图片未成功解码');
}

void main() {
  testWidgets('点击图片记录可预览、缩放和关闭，不修改剪贴板', (tester) async {
    final state = PreviewState(
      File('test/fixtures/image-preview.png').absolute.path,
    );
    addTearDown(state.dispose);
    await openPreview(tester, state);
    final image = await decodedImage(tester);
    expect(image.image!.width, 16);
    expect(image.image!.height, 16);
    expect(state.previewCalls, 1);
    expect(state.clipboardCalls, 0);
    final viewer = tester.widget<InteractiveViewer>(
      find.byType(InteractiveViewer),
    );
    viewer.transformationController!.value = Matrix4.diagonal3Values(2, 2, 1);
    await tester.tap(find.text('适应窗口'));
    await tester.pump();
    expect(viewer.transformationController!.value, Matrix4.identity());
    await tester.tap(find.byTooltip('关闭预览'));
    await tester.pumpAndSettle();
    expect(find.byType(ImagePreviewDialog), findsNothing);
    expect(state.clipboardCalls, 0);
    expect(tester.takeException(), isNull);
  });

  testWidgets('预览按钮同样打开图片', (tester) async {
    final state = PreviewState(
      File('test/fixtures/image-preview.png').absolute.path,
    );
    addTearDown(state.dispose);
    await openPreview(tester, state, button: true);
    await decodedImage(tester);
    expect(state.clipboardCalls, 0);
  });

  testWidgets('远端图片明确下载，进度结束后自动显示预览', (tester) async {
    final state = PreviewState(
      File('test/fixtures/image-preview.png').absolute.path,
      status: ClipStatus.CLIP_STATUS_REMOTE_ONLY,
    );
    addTearDown(state.dispose);
    await openPreview(tester, state);
    expect(find.text('图片尚未下载到这台电脑'), findsOneWidget);
    expect(state.previewCalls, 0);
    await tester.tap(find.text('下载并预览'));
    await tester.pump();
    expect(find.text('正在下载图片…'), findsOneWidget);
    expect(state.fetchCalls, 1);
    state.setStatus(ClipStatus.CLIP_STATUS_READY);
    await tester.pump();
    await decodedImage(tester);
    expect(state.previewCalls, 1);
    expect(state.clipboardCalls, 0);
  });

  testWidgets('下载失败显示原因并允许重试', (tester) async {
    final state = PreviewState('', status: ClipStatus.CLIP_STATUS_REMOTE_ONLY)
      ..failFetch = true;
    addTearDown(state.dispose);
    await openPreview(tester, state);
    await tester.tap(find.text('下载并预览'));
    await tester.pumpAndSettle();
    expect(find.text('对方设备暂时离线'), findsOneWidget);
    expect(find.text('下载并预览'), findsOneWidget);
    expect(state.previewCalls, 0);
  });

  testWidgets('图片过期或记录被删除时不再读取本地文件', (tester) async {
    final state = PreviewState('', status: ClipStatus.CLIP_STATUS_EXPIRED);
    addTearDown(state.dispose);
    await openPreview(tester, state);
    expect(find.text('图片已过期，无法预览'), findsOneWidget);
    state.debugEvent(Event(clipRemoved: 'image'));
    await tester.pumpAndSettle();
    expect(find.text('这条图片记录已被删除'), findsOneWidget);
    expect(state.previewCalls, 0);
  });

  testWidgets('读取接口失败可重试，预览期间记录过期立即隐藏图片', (tester) async {
    final state = PreviewState(
      File('test/fixtures/image-preview.png').absolute.path,
    )..failPreview = true;
    addTearDown(state.dispose);
    await openPreview(tester, state);
    await tester.pumpAndSettle();
    expect(find.text('图片已不在本机，可能已被清理'), findsOneWidget);
    state.failPreview = false;
    await tester.tap(find.text('重试'));
    await tester.pump();
    await decodedImage(tester);
    expect(state.previewCalls, 2);
    state.setStatus(ClipStatus.CLIP_STATUS_EXPIRED);
    await tester.pumpAndSettle();
    expect(find.text('图片已过期，无法预览'), findsOneWidget);
    expect(find.byKey(const ValueKey('image-preview-content')), findsNothing);
    expect(state.clipboardCalls, 0);
  });

  for (final corrupt in [false, true]) {
    testWidgets('文件${corrupt ? "损坏" : "被清理"}时显示提示', (tester) async {
      final dir = Directory.systemTemp.createTempSync('copysync-preview-');
      addTearDown(() => dir.deleteSync(recursive: true));
      final file = File('${dir.path}/invalid.png');
      if (corrupt) file.writeAsStringSync('not an image');
      final state = PreviewState(file.path);
      addTearDown(state.dispose);
      await openPreview(tester, state);
      for (
        var i = 0;
        i < 60 && find.text('无法读取图片，文件可能已被清理或损坏').evaluate().isEmpty;
        i++
      ) {
        await tester.runAsync(
          () => Future<void>.delayed(const Duration(milliseconds: 50)),
        );
        await tester.pump();
      }
      expect(find.text('无法读取图片，文件可能已被清理或损坏'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  }

  testWidgets('预览框按窗口比例定大小，小图按原始大小显示、不放大', (tester) async {
    tester.view.physicalSize = const Size(1000, 700);
    tester.view.devicePixelRatio = 1;
    addTearDown(tester.view.reset);
    final state = PreviewState(
      File('test/fixtures/image-preview.png').absolute.path,
    );
    addTearDown(state.dispose);
    await openPreview(tester, state);
    final box = tester.getSize(
      find.descendant(of: find.byType(Dialog), matching: find.byType(SizedBox)).first,
    );
    expect(box.width, closeTo(800, 0.5));
    expect(box.height, closeTo(560, 0.5));
    final image = await decodedImage(tester);
    expect(image.fit, BoxFit.scaleDown);
  });

  testWidgets('鼠标移到记录上时，「预览」按钮不挪位置', (tester) async {
    tester.view.physicalSize = const Size(1000, 700);
    tester.view.devicePixelRatio = 1;
    addTearDown(tester.view.reset);
    final state = PreviewState(
      File('test/fixtures/image-preview.png').absolute.path,
    );
    addTearDown(state.dispose);
    await tester.pumpWidget(CopySyncApp(state: state));
    await tester.pump();
    final before = tester.getCenter(find.text('预览'));
    expect(find.text('放入剪贴板').hitTestable(), findsNothing);

    final mouse = await tester.createGesture(kind: PointerDeviceKind.mouse);
    addTearDown(mouse.removePointer);
    await mouse.addPointer(location: Offset.zero);
    await mouse.moveTo(tester.getCenter(find.textContaining('来自 测试电脑')));
    await tester.pump();

    expect(find.text('放入剪贴板').hitTestable(), findsOneWidget);
    expect(tester.getCenter(find.text('预览')), before);
  });
}

import 'dart:io';

import 'package:flutter/material.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../icons.dart';
import '../main.dart';
import '../theme.dart';

Future<void> showImagePreview(BuildContext context, String clipId) =>
    showDialog<void>(
      context: context,
      builder: (_) => ImagePreviewDialog(clipId: clipId),
    );

class ImagePreviewDialog extends StatefulWidget {
  const ImagePreviewDialog({super.key, required this.clipId});

  final String clipId;

  @override
  State<ImagePreviewDialog> createState() => _ImagePreviewDialogState();
}

class _ImagePreviewDialogState extends State<ImagePreviewDialog> {
  final _transform = TransformationController();
  Future<String>? _preview;
  bool _requesting = false;
  String? _downloadError;

  @override
  void dispose() {
    _transform.dispose();
    super.dispose();
  }

  Future<void> _download() async {
    if (_requesting) return;
    setState(() {
      _requesting = true;
      _downloadError = null;
      _preview = null;
    });
    try {
      await AppScope.read(context).fetchImagePreview(widget.clipId);
    } catch (e) {
      if (mounted) setState(() => _downloadError = AppState.describeError(e));
    } finally {
      if (mounted) setState(() => _requesting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final record = state.records
        .where((r) => r.id == widget.clipId)
        .firstOrNull;
    // 预览框按应用窗口的比例定大小：窗口小，预览也小，不会一打开就把整个窗口盖满
    final window = MediaQuery.sizeOf(context);
    return Dialog(
      insetPadding: EdgeInsets.zero,
      child: SizedBox(
        width: (window.width * 0.8).clamp(280.0, 1400.0),
        height: (window.height * 0.8).clamp(240.0, 1000.0),
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(
                Insets.xl,
                Insets.md,
                Insets.md,
                Insets.md,
              ),
              child: Row(
                children: [
                  Expanded(
                    child: Text('图片预览', style: context.text.titleMedium),
                  ),
                  IconButton(
                    tooltip: '关闭预览',
                    onPressed: () => Navigator.pop(context),
                    icon: Icon(AppIcons.xmark),
                  ),
                ],
              ),
            ),
            Expanded(
              child: Container(
                width: double.infinity,
                color: context.palette.fill,
                child: _content(state, record),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(
                horizontal: Insets.xl,
                vertical: Insets.sm,
              ),
              child: Row(
                children: [
                  Expanded(
                    child: Text('滚轮或触控板缩放，拖动查看', style: context.text.bodySmall),
                  ),
                  TextButton(
                    onPressed: () => _transform.value = Matrix4.identity(),
                    child: const Text('适应窗口'),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _content(AppState state, ClipRecord? record) {
    if (record == null) return _message('这条图片记录已被删除');
    if (record.status == ClipStatus.CLIP_STATUS_EXPIRED) {
      return _message('图片已过期，无法预览');
    }
    if (_requesting || record.status == ClipStatus.CLIP_STATUS_FETCHING) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            CircularProgressIndicator(
              value: state.progressOf(record.id),
              strokeWidth: 2,
            ),
            const SizedBox(height: Insets.lg),
            const Text('正在下载图片…'),
          ],
        ),
      );
    }
    if (record.status == ClipStatus.CLIP_STATUS_REMOTE_ONLY ||
        record.status == ClipStatus.CLIP_STATUS_FAILED) {
      return _message(
        _downloadError ??
            (record.status == ClipStatus.CLIP_STATUS_FAILED
                ? '图片下载失败，请重试'
                : '图片尚未下载到这台电脑'),
        action: FilledButton(onPressed: _download, child: const Text('下载并预览')),
      );
    }
    if (record.status != ClipStatus.CLIP_STATUS_READY) {
      return _message('图片暂时无法预览');
    }

    _preview ??= state.imagePreviewPath(record.id);
    return FutureBuilder<String>(
      future: _preview,
      builder: (context, snapshot) {
        if (snapshot.hasError) {
          return _message(
            AppState.describeError(snapshot.error!),
            action: OutlinedButton(
              onPressed: () => setState(() => _preview = null),
              child: const Text('重试'),
            ),
          );
        }
        if (!snapshot.hasData) {
          return const Center(child: CircularProgressIndicator(strokeWidth: 2));
        }
        return InteractiveViewer(
          transformationController: _transform,
          minScale: 0.5,
          maxScale: 8,
          trackpadScrollCausesScale: true,
          child: SizedBox.expand(
            child: Image(
              key: const ValueKey('image-preview-content'),
              // 预览限制解码尺寸并保持比例，超大截图不占用整张原图的内存。
              image: ResizeImage(
                FileImage(File(snapshot.data!)),
                width: 2048,
                height: 2048,
                policy: ResizeImagePolicy.fit,
              ),
              // 比预览框小的图片按原始大小显示，大的才缩小到放得下：小图标不会被放大成一整屏
              fit: BoxFit.scaleDown,
              semanticLabel: '复制记录中的图片',
              frameBuilder: (_, child, frame, synchronous) =>
                  frame != null || synchronous
                  ? child
                  : const Center(
                      child: CircularProgressIndicator(strokeWidth: 2),
                    ),
              errorBuilder: (_, _, _) => _message('无法读取图片，文件可能已被清理或损坏'),
            ),
          ),
        );
      },
    );
  }

  Widget _message(String text, {Widget? action}) => Center(
    child: Padding(
      padding: const EdgeInsets.all(Insets.xl),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            text,
            textAlign: TextAlign.center,
            style: context.text.bodyMedium,
          ),
          if (action != null) ...[const SizedBox(height: Insets.lg), action],
        ],
      ),
    ),
  );
}

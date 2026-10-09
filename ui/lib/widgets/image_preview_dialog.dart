import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../i18n.dart';
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
  Future<GetImagePreviewResponse>? _preview;
  bool _requesting = false;
  String? _downloadError;

  // 一条记录里有多张图片时（复制了几张照片），左右翻看
  int _index = 0;
  int _count = 1;
  String _name = '';

  Future<GetImagePreviewResponse> _load(AppState state, String clipId) =>
      state.imagePreview(clipId, index: _index).then((response) {
        if (mounted) {
          setState(() {
            _count = response.count < 1 ? 1 : response.count;
            _name = response.name;
          });
        }
        return response;
      });

  void _go(int delta) {
    if (_count < 2) return;
    setState(() {
      _index = (_index + delta) % _count;
      if (_index < 0) _index += _count;
      _preview = null;
      _transform.value = Matrix4.identity();
    });
  }

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
    final l = context.l10n;
    final dialog = Dialog(
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
                    child: Text(
                      _name.isEmpty ? l.imagePreviewTitle : _name,
                      style: context.text.titleMedium,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  if (_count > 1) ...[
                    IconButton(
                      tooltip: l.previousImage,
                      onPressed: () => _go(-1),
                      icon: Icon(AppIcons.chevronLeft, size: 16),
                    ),
                    Text('${_index + 1} / $_count', style: context.text.bodySmall),
                    IconButton(
                      tooltip: l.nextImage,
                      onPressed: () => _go(1),
                      icon: Icon(AppIcons.chevronRight, size: 16),
                    ),
                    const SizedBox(width: Insets.sm),
                  ],
                  IconButton(
                    tooltip: l.closePreview,
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
                    child: Text(context.l10n.previewHint, style: context.text.bodySmall),
                  ),
                  TextButton(
                    onPressed: () => _transform.value = Matrix4.identity(),
                    child: Text(context.l10n.fitWindow),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
    // 左右方向键翻看上一张、下一张
    return CallbackShortcuts(
      bindings: {
        const SingleActivator(LogicalKeyboardKey.arrowLeft): () => _go(-1),
        const SingleActivator(LogicalKeyboardKey.arrowRight): () => _go(1),
      },
      child: Focus(autofocus: true, child: dialog),
    );
  }

  Widget _content(AppState state, ClipRecord? record) {
    if (record == null) return _message(context.l10n.previewDeleted);
    if (record.status == ClipStatus.CLIP_STATUS_EXPIRED) {
      return _message(context.l10n.previewExpired);
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
            Text(context.l10n.previewDownloading),
          ],
        ),
      );
    }
    if (record.status == ClipStatus.CLIP_STATUS_REMOTE_ONLY ||
        record.status == ClipStatus.CLIP_STATUS_FAILED) {
      return _message(
        _downloadError ??
            (record.status == ClipStatus.CLIP_STATUS_FAILED
                ? context.l10n.previewDownloadFailed
                : context.l10n.previewNotDownloaded),
        action: FilledButton(onPressed: _download, child: Text(context.l10n.downloadAndPreview)),
      );
    }
    if (record.status != ClipStatus.CLIP_STATUS_READY) {
      return _message(context.l10n.previewUnavailable);
    }

    _preview ??= _load(state, record.id);
    return FutureBuilder<GetImagePreviewResponse>(
      future: _preview,
      builder: (context, snapshot) {
        if (snapshot.hasError) {
          return _message(
            AppState.describeError(snapshot.error!),
            action: OutlinedButton(
              onPressed: () => setState(() => _preview = null),
              child: Text(context.l10n.retry),
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
                FileImage(File(snapshot.data!.path)),
                width: 2048,
                height: 2048,
                policy: ResizeImagePolicy.fit,
              ),
              // 比预览框小的图片按原始大小显示，大的才缩小到放得下：小图标不会被放大成一整屏
              fit: BoxFit.scaleDown,
              semanticLabel: context.l10n.imageSemantic,
              frameBuilder: (_, child, frame, synchronous) =>
                  frame != null || synchronous
                  ? child
                  : const Center(
                      child: CircularProgressIndicator(strokeWidth: 2),
                    ),
              errorBuilder: (_, _, _) => _message(context.l10n.previewUnreadable),
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

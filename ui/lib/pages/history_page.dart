import 'package:flutter/material.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../i18n.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';
import '../widgets/image_preview_dialog.dart';
import '../icons.dart';
import '../platform.dart';

enum _Filter { all, text, image, file }

class HistoryPage extends StatefulWidget {
  const HistoryPage({super.key});

  @override
  State<HistoryPage> createState() => _HistoryPageState();
}

class _HistoryPageState extends State<HistoryPage> {
  _Filter _filter = _Filter.all;

  bool _matches(ClipRecord r) => switch (_filter) {
        _Filter.all => true,
        // 富文本也是文字，用户不关心它在剪贴板里是不是 HTML
        _Filter.text =>
          r.kind == ClipKind.CLIP_KIND_TEXT || r.kind == ClipKind.CLIP_KIND_HTML,
        _Filter.image => r.kind == ClipKind.CLIP_KIND_IMAGE,
        _Filter.file => r.kind == ClipKind.CLIP_KIND_FILE,
      };

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final records = state.records.where(_matches).toList();
    final ttl = state.config?.historyTtlSeconds.toInt();
    final l = context.l10n;
    final count = state.records.length;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PageHeader(
          title: l.navHistory,
          subtitle: state.records.isEmpty
              ? l.historyEmptySubtitle
              : ttl == null
                  ? l.historyCount(count)
                  : l.historyCountKept(count, durationLabel(l, ttl)),
          actions: [
            if (state.records.isNotEmpty) ...[
              Segmented<_Filter>(
                value: _filter,
                segments: [
                  (_Filter.all, l.filterAll),
                  (_Filter.text, l.filterText),
                  (_Filter.image, l.filterImage),
                  (_Filter.file, l.filterFile),
                ],
                onChanged: (f) => setState(() => _filter = f),
              ),
              IconAction(
                icon: AppIcons.trash,
                tooltip: l.historyClearTooltip,
                onTap: () => _confirmClear(context, state),
              ),
            ],
          ],
        ),
        if (_needsPermission(state.status?.clipboardPermission))
          Padding(
            padding: const EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.md),
            child: _PermissionNotice(permission: state.status!.clipboardPermission),
          ),
        Expanded(
          child: records.isEmpty
              ? EmptyState(
                  icon: _filter == _Filter.all
                      ? AppIcons.docOnClipboard
                      : AppIcons.filter,
                  title: _filter == _Filter.all ? l.historyEmptyTitle : l.historyEmptyFilteredTitle,
                  description: _filter == _Filter.all
                      ? l.historyEmptyBody
                      : null,
                )
              : CustomScrollView(
                  slivers: [
                    SliverPadding(
                      padding: const EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.xl),
                      sliver: DecoratedSliver(
                        decoration: BoxDecoration(
                          color: context.palette.surface,
                          borderRadius: BorderRadius.circular(Radii.md),
                          border: Border.all(color: context.palette.border),
                        ),
                        sliver: SliverPadding(
                          padding: const EdgeInsets.all(Insets.xs),
                          sliver: SliverList.separated(
                            itemCount: records.length,
                            separatorBuilder: (_, _) => const Divider(indent: 44, endIndent: 8),
                            itemBuilder: (_, i) => _RecordRow(record: records[i]),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
        ),
      ],
    );
  }

  static bool _needsPermission(ClipboardPermission? p) =>
      p == ClipboardPermission.CLIPBOARD_PERMISSION_DEFAULT ||
      p == ClipboardPermission.CLIPBOARD_PERMISSION_ASK ||
      p == ClipboardPermission.CLIPBOARD_PERMISSION_ALWAYS_DENY;

  Future<void> _confirmClear(BuildContext context, AppState state) async {
    final l = context.l10n;
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(l.clearConfirmTitle),
        content: Text(l.clearConfirmBody(Wording.thisDeviceInText)),
        actions: [
          OutlinedButton(onPressed: () => Navigator.pop(ctx, false), child: Text(l.cancel)),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: FilledButton.styleFrom(backgroundColor: ctx.palette.danger),
            child: Text(l.clear),
          ),
        ],
      ),
    );
    if (ok != true || !context.mounted) return;
    try {
      await state.deleteRecords(all: true);
      if (context.mounted) showToast(context, l.historyCleared);
    } catch (e) {
      if (context.mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }
}

/// 读剪贴板需要用户授权。没授权时整个产品等于不工作，所以放在最显眼的位置。
class _PermissionNotice extends StatelessWidget {
  const _PermissionNotice({required this.permission});

  final ClipboardPermission permission;

  @override
  Widget build(BuildContext context) {
    final denied = permission == ClipboardPermission.CLIPBOARD_PERMISSION_ALWAYS_DENY;
    final l = context.l10n;
    return Notice(
      title: denied ? l.permDeniedTitle : l.permAskTitle,
      message: denied ? l.permDeniedBody : l.permAskBody,
      action: OutlinedButton(
        onPressed: () async {
          try {
            await AppScope.read(context).requestClipboardPermission();
          } catch (e) {
            if (context.mounted) showToast(context, AppState.describeError(e), error: true);
          }
        },
        child: Text(l.openSystemSettings),
      ),
    );
  }
}

class _RecordRow extends StatefulWidget {
  const _RecordRow({required this.record});

  final ClipRecord record;

  @override
  State<_RecordRow> createState() => _RecordRowState();
}

class _RecordRowState extends State<_RecordRow> {
  bool _hovered = false;
  bool _busy = false;

  @override
  Widget build(BuildContext context) {
    final r = widget.record;
    final p = context.palette;
    // 剪贴板里的图片、以及复制的图片文件，点一下这一行就能预览，不另放按钮
    final previewable = r.imageCount > 0;

    return MouseRegion(
      cursor: previewable ? SystemMouseCursors.click : MouseCursor.defer,
      onEnter: (_) => setState(() => _hovered = true),
      onExit: (_) => setState(() => _hovered = false),
      child: InkWell(
        onTap: previewable ? () => showImagePreview(context, r.id) : null,
        borderRadius: BorderRadius.circular(Radii.sm),
        child: Container(
          constraints: const BoxConstraints(minHeight: 52),
          padding: const EdgeInsets.symmetric(horizontal: Insets.md, vertical: Insets.sm),
          decoration: BoxDecoration(
            color: _hovered ? p.hover : Colors.transparent,
            borderRadius: BorderRadius.circular(Radii.sm),
          ),
          child: Row(
            children: [
              Icon(kindIcon(r), size: 18, color: p.textDim),
              const SizedBox(width: Insets.md),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      recordTitle(r),
                      style: context.text.labelLarge,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      [
                        r.outgoing ? context.l10n.thisDeviceShort : context.l10n.fromDevice(r.originDeviceName),
                        if (r.totalSize > 0) humanBytes(r.totalSize.toInt()),
                        relativeTime(DateTime.fromMillisecondsSinceEpoch(r.createdAtUnix.toInt() * 1000)),
                      ].join(' · '),
                      style: context.text.bodySmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              const SizedBox(width: Insets.md),
              _trailing(context, r),
            ],
          ),
        ),
      ),
    );
  }

  Widget _trailing(BuildContext context, ClipRecord r) {
    final p = context.palette;
    final l = context.l10n;
    if (_busy) {
      return const SizedBox(
        width: 16,
        height: 16,
        child: CircularProgressIndicator(strokeWidth: 2),
      );
    }
    switch (r.status) {
      case ClipStatus.CLIP_STATUS_REMOTE_ONLY:
        // 超过自动同步阈值的大文件：只有元数据，需要时再拉
        return OutlinedButton(
          onPressed: () => _run(() => AppScope.read(context).fetch(r.id), l.fetchStarted),
          child: Text(l.fetch),
        );
      case ClipStatus.CLIP_STATUS_FETCHING:
        // 知道总量时显示真实进度，否则显示不确定进度
        return SizedBox(
          width: 72,
          child: LinearProgressIndicator(
            minHeight: 3,
            value: AppScope.of(context).progressOf(r.id),
          ),
        );
      case ClipStatus.CLIP_STATUS_FAILED:
        return Tooltip(
          message: r.error.isEmpty ? l.transferFailed : r.error,
          child: Text(l.transferFailed, style: context.text.bodySmall?.copyWith(color: p.danger)),
        );
      case ClipStatus.CLIP_STATUS_EXPIRED:
        return Text(l.expired, style: context.text.bodySmall?.copyWith(color: p.textFaint));
      default:
        // 悬停时才显示「放入剪贴板」与删除，静态时列表保持干净。
        // 不悬停时它们只是看不见，位置照样留着：鼠标移上来时标题不会被挤短一截
        Widget hoverOnly(Widget child) => Visibility(
              visible: _hovered,
              maintainSize: true,
              maintainAnimation: true,
              maintainState: true,
              child: child,
            );
        return Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            hoverOnly(TextButton(
              onPressed: () =>
                  _run(() => AppScope.read(context).applyToClipboard(r.id), l.putOnClipboardDone),
              child: Text(l.putOnClipboard),
            )),
            hoverOnly(IconAction(
              icon: AppIcons.trash,
              tooltip: l.deleteRecord,
              onTap: () => _run(() => AppScope.read(context).deleteRecords(ids: [r.id]), l.deleted),
            )),
          ],
        );
    }
  }

  Future<void> _run(Future<void> Function() action, String success) async {
    setState(() => _busy = true);
    try {
      await action();
      if (mounted) showToast(context, success);
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }
}

/// 按内容类型区分的图标：只靠形状区分，不靠颜色。
IconData kindIcon(ClipRecord r) => switch (r.kind) {
      ClipKind.CLIP_KIND_TEXT => AppIcons.textAlignLeft,
      ClipKind.CLIP_KIND_HTML => AppIcons.docRichtext,
      ClipKind.CLIP_KIND_IMAGE => AppIcons.photo,
      // 单张图片文件也用图片图标，一眼看出点开能预览
      ClipKind.CLIP_KIND_FILE when r.items.length == 1 && r.imageCount == 1 => AppIcons.photo,
      ClipKind.CLIP_KIND_FILE when r.items.length > 1 => AppIcons.docOnDoc,
      ClipKind.CLIP_KIND_FILE when r.items.isNotEmpty && r.items.first.isDir =>
        AppIcons.folder,
      ClipKind.CLIP_KIND_FILE => AppIcons.doc,
      _ => AppIcons.questionCircle,
    };

String recordTitle(ClipRecord r) {
  // 图片和文件的标题按当前界面语言拼：后台服务存下的摘要是写入时的中文，不会跟着换语言
  if (r.kind == ClipKind.CLIP_KIND_IMAGE) return appL10n.imageRecord;
  if (r.kind != ClipKind.CLIP_KIND_FILE && r.textPreview.isNotEmpty) {
    return r.textPreview.replaceAll('\n', ' ').trim();
  }
  if (r.items.isEmpty) return appL10n.emptyRecord;
  if (r.items.length == 1) return r.items.first.name;
  return appL10n.itemsSummary(r.items.first.name, r.items.length - 1, r.items.length);
}

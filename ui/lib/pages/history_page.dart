import 'package:flutter/material.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';
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

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PageHeader(
          title: '复制记录',
          subtitle: state.records.isEmpty
              ? '在任意一台设备上复制，内容会出现在这里'
              : '共 ${state.records.length} 条${ttl == null ? '' : '，保留 ${_ttlLabel(ttl)}'}',
          actions: [
            if (state.records.isNotEmpty) ...[
              Segmented<_Filter>(
                value: _filter,
                segments: const [
                  (_Filter.all, '全部'),
                  (_Filter.text, '文本'),
                  (_Filter.image, '图片'),
                  (_Filter.file, '文件'),
                ],
                onChanged: (f) => setState(() => _filter = f),
              ),
              IconAction(
                icon: AppIcons.trash,
                tooltip: '清空记录',
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
                  title: _filter == _Filter.all ? '还没有记录' : '没有这一类的记录',
                  description: _filter == _Filter.all
                      ? '在任意一台已配对的设备上复制文本、图片或文件，记录会出现在这里。'
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

  static String _ttlLabel(int seconds) =>
      seconds % 86400 == 0 ? '${seconds ~/ 86400} 天' : '${(seconds / 3600).round()} 小时';

  Future<void> _confirmClear(BuildContext context, AppState state) async {
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('清空全部记录？'),
        content: Text('所有设备上同步来的记录和已缓存的文件都会从${Wording.thisDeviceInText}上删除。原始文件不受影响。'),
        actions: [
          OutlinedButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: FilledButton.styleFrom(backgroundColor: ctx.palette.danger),
            child: const Text('清空'),
          ),
        ],
      ),
    );
    if (ok != true || !context.mounted) return;
    try {
      await state.deleteRecords(all: true);
      if (context.mounted) showToast(context, '记录已清空');
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
    return Notice(
      title: denied ? '已禁止 CopySync 读取剪贴板' : '允许 CopySync 读取剪贴板',
      message: denied
          ? '在这台 Mac 上复制的内容不会同步出去。到系统设置里把 CopySync Daemon 改为「允许」。'
          : '否则 macOS 每次都会弹窗询问。到系统设置里把 CopySync Daemon 改为「允许」。',
      action: OutlinedButton(
        onPressed: () async {
          try {
            await AppScope.read(context).requestClipboardPermission();
          } catch (e) {
            if (context.mounted) showToast(context, AppState.describeError(e), error: true);
          }
        },
        child: const Text('打开系统设置'),
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

    return MouseRegion(
      onEnter: (_) => setState(() => _hovered = true),
      onExit: (_) => setState(() => _hovered = false),
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
                      r.outgoing ? '本机' : '来自 ${r.originDeviceName}',
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
    );
  }

  Widget _trailing(BuildContext context, ClipRecord r) {
    final p = context.palette;
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
          onPressed: () => _run(() => AppScope.read(context).fetch(r.id), '开始拉取'),
          child: const Text('拉取到本机'),
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
          message: r.error.isEmpty ? '传输失败' : r.error,
          child: Text('传输失败', style: context.text.bodySmall?.copyWith(color: p.danger)),
        );
      case ClipStatus.CLIP_STATUS_EXPIRED:
        return Text('已过期', style: context.text.bodySmall?.copyWith(color: p.textFaint));
      default:
        // 悬停时才显示操作，静态时列表保持干净
        if (!_hovered) return const SizedBox(height: 28);
        return Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextButton(
              onPressed: () =>
                  _run(() => AppScope.read(context).applyToClipboard(r.id), '已放入剪贴板'),
              child: const Text('放入剪贴板'),
            ),
            IconAction(
              icon: AppIcons.trash,
              tooltip: '删除这条记录',
              onTap: () => _run(() => AppScope.read(context).deleteRecords(ids: [r.id]), '已删除'),
            ),
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
      ClipKind.CLIP_KIND_FILE when r.items.length > 1 => AppIcons.docOnDoc,
      ClipKind.CLIP_KIND_FILE when r.items.isNotEmpty && r.items.first.isDir =>
        AppIcons.folder,
      ClipKind.CLIP_KIND_FILE => AppIcons.doc,
      _ => AppIcons.questionCircle,
    };

String recordTitle(ClipRecord r) {
  if (r.textPreview.isNotEmpty) return r.textPreview.replaceAll('\n', ' ').trim();
  if (r.items.isEmpty) return '（空）';
  if (r.items.length == 1) return r.items.first.name;
  return '${r.items.first.name} 等 ${r.items.length} 项';
}

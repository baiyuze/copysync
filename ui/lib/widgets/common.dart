import 'package:flutter/material.dart';

import '../theme.dart';
import '../i18n.dart';
import '../icons.dart';

/// 页面统一的标题区。保证各页视觉节奏一致。
class PageHeader extends StatelessWidget {
  const PageHeader({
    super.key,
    required this.title,
    this.subtitle,
    this.actions = const [],
  });

  final String title;
  final String? subtitle;
  final List<Widget> actions;

  @override
  Widget build(BuildContext context) {
    return Padding(
      // 顶部留出 macOS 红绿灯所在的标题栏高度
      padding: const EdgeInsets.fromLTRB(Insets.xl, 44, Insets.xl, Insets.lg),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: context.text.headlineSmall),
                if (subtitle != null) ...[
                  const SizedBox(height: 2),
                  Text(subtitle!, style: context.text.bodySmall),
                ],
              ],
            ),
          ),
          for (final a in actions) ...[const SizedBox(width: Insets.sm), a],
        ],
      ),
    );
  }
}

/// 分组容器：白底、细边框，内部各行之间用缩进的细线分隔。
/// 与「系统设置」的分组列表同一种做法——用分隔线而不是一张张独立卡片。
class GroupBox extends StatelessWidget {
  const GroupBox({super.key, required this.children, this.dividerIndent = Insets.lg});

  final List<Widget> children;
  final double dividerIndent;

  @override
  Widget build(BuildContext context) {
    return Container(
      clipBehavior: Clip.antiAlias,
      decoration: BoxDecoration(
        color: context.palette.surface,
        borderRadius: BorderRadius.circular(Radii.md),
        border: Border.all(color: context.palette.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          for (var i = 0; i < children.length; i++) ...[
            if (i > 0) Divider(indent: dividerIndent),
            children[i],
          ],
        ],
      ),
    );
  }
}

/// 分组的标题与脚注。
class GroupSection extends StatelessWidget {
  const GroupSection({super.key, required this.title, required this.child, this.footnote});

  final String title;
  final String? footnote;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: Insets.xl),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Padding(
            padding: const EdgeInsets.only(left: 2, bottom: Insets.sm),
            child: Text(title, style: context.text.titleSmall),
          ),
          child,
          if (footnote != null)
            Padding(
              padding: const EdgeInsets.only(left: 2, top: Insets.sm, right: Insets.xl),
              child: Text(footnote!, style: context.text.bodySmall),
            ),
        ],
      ),
    );
  }
}

/// 设置项行：左侧标题与说明，右侧控件。
class SettingRow extends StatelessWidget {
  const SettingRow({super.key, required this.title, required this.child, this.description});

  final String title;
  final String? description;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
      constraints: const BoxConstraints(minHeight: 48),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: Insets.lg, vertical: 10),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(title, style: context.text.labelLarge),
                  if (description != null) ...[
                    const SizedBox(height: 2),
                    Text(description!, style: context.text.bodySmall),
                  ],
                ],
              ),
            ),
            const SizedBox(width: Insets.lg),
            child,
          ],
        ),
      ),
    );
  }
}

/// macOS 尺寸的开关：Material 默认的开关对桌面来说太大。
class SmallSwitch extends StatelessWidget {
  const SmallSwitch({super.key, required this.value, required this.onChanged});

  final bool value;
  final ValueChanged<bool>? onChanged;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 24,
      child: FittedBox(
        child: Switch(value: value, onChanged: onChanged),
      ),
    );
  }
}

/// 分段控件，对应 macOS 的 NSSegmentedControl：浅灰轨道，选中段浮起为白色。
class Segmented<T> extends StatelessWidget {
  const Segmented({
    super.key,
    required this.value,
    required this.segments,
    required this.onChanged,
  });

  final T value;
  final List<(T, String)> segments;
  final ValueChanged<T> onChanged;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    return Container(
      padding: const EdgeInsets.all(2),
      decoration: BoxDecoration(
        color: p.fill,
        borderRadius: BorderRadius.circular(Radii.sm + 1),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          for (final (v, label) in segments)
            MouseRegion(
              cursor: SystemMouseCursors.click,
              child: GestureDetector(
                onTap: () => onChanged(v),
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 120),
                  padding: const EdgeInsets.symmetric(horizontal: Insets.md, vertical: 4),
                  decoration: BoxDecoration(
                    color: v == value ? p.surface : Colors.transparent,
                    borderRadius: BorderRadius.circular(Radii.sm - 1),
                    boxShadow: v == value
                        ? const [
                            BoxShadow(
                                color: Color(0x1A000000), blurRadius: 2, offset: Offset(0, 1)),
                          ]
                        : null,
                  ),
                  child: Text(
                    label,
                    style: context.text.labelMedium?.copyWith(
                      color: v == value ? p.text : p.textDim,
                    ),
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}

/// 空状态。给出下一步该做什么，而不是只说"暂无数据"。
class EmptyState extends StatelessWidget {
  const EmptyState({
    super.key,
    required this.icon,
    required this.title,
    this.description,
    this.action,
  });

  final IconData icon;
  final String title;
  final String? description;
  final Widget? action;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 340),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 36, color: context.palette.textFaint),
            const SizedBox(height: Insets.lg),
            Text(title, style: context.text.titleMedium, textAlign: TextAlign.center),
            if (description != null) ...[
              const SizedBox(height: 6),
              Text(description!, style: context.text.bodySmall, textAlign: TextAlign.center),
            ],
            if (action != null) ...[
              const SizedBox(height: Insets.xl),
              action!,
            ],
          ],
        ),
      ),
    );
  }
}

/// 行内提示条：说明问题，并给出解决它的按钮。
class Notice extends StatelessWidget {
  const Notice({
    super.key,
    required this.title,
    required this.message,
    this.action,
    this.icon,
    this.iconColor,
  });

  final String title;
  final String message;
  final Widget? action;
  final IconData? icon;
  final Color? iconColor;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    final icon = this.icon ?? AppIcons.exclamationmarkCircle;
    return Container(
      padding: const EdgeInsets.fromLTRB(Insets.lg, Insets.md, Insets.md, Insets.md),
      decoration: BoxDecoration(
        color: p.surface,
        borderRadius: BorderRadius.circular(Radii.md),
        border: Border.all(color: p.border),
      ),
      child: Row(
        children: [
          Icon(icon, size: 18, color: iconColor ?? p.warning),
          const SizedBox(width: Insets.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: context.text.labelLarge),
                const SizedBox(height: 2),
                Text(message, style: context.text.bodySmall),
              ],
            ),
          ),
          if (action != null) ...[const SizedBox(width: Insets.md), action!],
        ],
      ),
    );
  }
}

/// 状态圆点 + 文字。设备在线状态、连接状态都用它，保持一致。
class StatusDot extends StatelessWidget {
  const StatusDot({super.key, required this.color, required this.label, this.labelColor});

  final Color color;
  final String label;
  final Color? labelColor;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 7,
          height: 7,
          decoration: BoxDecoration(color: color, shape: BoxShape.circle),
        ),
        const SizedBox(width: 6),
        Text(label, style: context.text.bodySmall?.copyWith(color: labelColor)),
      ],
    );
  }
}

/// App 图标。侧边栏、欢迎页、关于里都用同一张图，而不是另画一个符号。
class AppMark extends StatelessWidget {
  const AppMark({super.key, this.size = 20});

  final double size;

  @override
  Widget build(BuildContext context) {
    return Image.asset(
      'assets/app_icon.png',
      width: size,
      height: size,
      filterQuality: FilterQuality.medium,
    );
  }
}

/// 可点击的图标按钮，悬停时出现浅色底。
class IconAction extends StatefulWidget {
  const IconAction({super.key, required this.icon, required this.tooltip, this.onTap});

  final IconData icon;
  final String tooltip;
  final VoidCallback? onTap;

  @override
  State<IconAction> createState() => _IconActionState();
}

class _IconActionState extends State<IconAction> {
  bool _hovered = false;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    return Tooltip(
      message: widget.tooltip,
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        onEnter: (_) => setState(() => _hovered = true),
        onExit: (_) => setState(() => _hovered = false),
        child: GestureDetector(
          onTap: widget.onTap,
          child: Container(
            width: 28,
            height: 28,
            decoration: BoxDecoration(
              color: _hovered ? p.hover : Colors.transparent,
              borderRadius: BorderRadius.circular(Radii.sm),
            ),
            child: Icon(widget.icon, size: 16, color: _hovered ? p.text : p.textDim),
          ),
        ),
      ),
    );
  }
}

/// 字节数的可读形式。
String humanBytes(int n) {
  if (n < 1024) return '$n B';
  const units = ['KB', 'MB', 'GB', 'TB'];
  var v = n / 1024;
  var i = 0;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i++;
  }
  return '${v.toStringAsFixed(v >= 100 ? 0 : 1)} ${units[i]}';
}

/// 时长：整天的按天说，否则按小时。
String durationLabel(AppLocalizations l, int seconds) => seconds % 86400 == 0
    ? l.durationDays(seconds ~/ 86400)
    : l.durationHours((seconds / 3600).round());

/// 相对时间。刚刚发生的事用"刚刚"比"14:03:22"更符合直觉。
String relativeTime(DateTime t, {DateTime? now}) {
  final d = (now ?? DateTime.now()).difference(t);
  final l = appL10n;
  if (d.inSeconds < 10) return l.justNow;
  if (d.inMinutes < 1) return l.secondsAgo(d.inSeconds);
  if (d.inHours < 1) return l.minutesAgo(d.inMinutes);
  if (d.inDays < 1) return l.hoursAgo(d.inHours);
  if (d.inDays < 7) return l.daysAgo(d.inDays);
  return l.monthDay(t.month, t.day);
}

void showToast(BuildContext context, String message, {bool error = false}) {
  final p = context.palette;
  ScaffoldMessenger.of(context)
    ..hideCurrentSnackBar()
    ..showSnackBar(SnackBar(
      content: Row(
        children: [
          Icon(
            error ? AppIcons.exclamationmarkCircle : AppIcons.checkmarkCircle,
            size: 16,
            color: error ? PaletteColors.onDark(p.danger) : PaletteColors.onDark(p.online),
          ),
          const SizedBox(width: Insets.sm),
          Expanded(child: Text(message)),
        ],
      ),
      duration: const Duration(seconds: 3),
      width: 360,
    ));
}

/// 提示条是深色底，浅色模式下的深绿、深红在上面看不清，换成亮色版本。
abstract final class PaletteColors {
  static Color onDark(Color c) {
    if (c == AppPalette.light.online) return AppPalette.dark.online;
    if (c == AppPalette.light.danger) return AppPalette.dark.danger;
    return c;
  }
}

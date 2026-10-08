import 'dart:io';

import 'package:flutter/material.dart';

import 'app_state.dart';
import 'background_service.dart';
import 'pages/devices_page.dart';
import 'pages/history_page.dart';
import 'pages/service_page.dart';
import 'pages/settings_page.dart';
import 'theme.dart';
import 'widgets/common.dart';
import 'widgets/pairing_dialog.dart';
import 'icons.dart';

void main() => runApp(const CopySyncApp());

class CopySyncApp extends StatefulWidget {
  const CopySyncApp({super.key, this.state});

  /// 测试与截图用：注入预先填好数据的状态，不去连真实的后台服务。
  final AppState? state;

  @override
  State<CopySyncApp> createState() => _CopySyncAppState();
}

class _CopySyncAppState extends State<CopySyncApp> {
  late final AppState _state;

  @override
  void initState() {
    super.initState();
    _state = widget.state ?? (AppState()..connect());
  }

  @override
  void dispose() {
    if (widget.state == null) _state.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // AppScope 必须包在 MaterialApp 外面：showDialog 推出的路由挂在 Navigator 上，
    // 与 home 是兄弟而非后代，放在 home 里的话对话框里就取不到 AppState。
    return AppScope(
      state: _state,
      child: MaterialApp(
        title: 'CopySync',
        debugShowCheckedModeBanner: false,
        theme: buildTheme(Brightness.light),
        darkTheme: buildTheme(Brightness.dark),
        home: const Shell(),
      ),
    );
  }
}

/// 把 AppState 传递给整棵子树。用 InheritedNotifier 而非引入
/// provider 之类的依赖——状态只有一个，没必要为此多一个包。
class AppScope extends InheritedNotifier<AppState> {
  const AppScope({super.key, required AppState state, required super.child})
      : super(notifier: state);

  static AppState of(BuildContext context) {
    final scope = context.dependOnInheritedWidgetOfExactType<AppScope>();
    assert(scope != null, 'AppScope 未在上层提供');
    return scope!.notifier!;
  }

  /// 只读取一次、不订阅变化。用于事件回调里拿实例。
  static AppState read(BuildContext context) =>
      context.getInheritedWidgetOfExactType<AppScope>()!.notifier!;
}

enum Section { history, devices, settings }

class Shell extends StatefulWidget {
  const Shell({super.key, this.initialSection = Section.history});

  final Section initialSection;

  @override
  State<Shell> createState() => _ShellState();
}

class _ShellState extends State<Shell> {
  late Section _section = widget.initialSection;
  bool _pairingDialogOpen = false;

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);

    // 对端发起配对时主动弹出确认——指纹核对不能等用户自己去翻页面
    final pending = state.pendingPairing;
    if (pending != null && !_pairingDialogOpen) {
      _pairingDialogOpen = true;
      WidgetsBinding.instance.addPostFrameCallback((_) async {
        await showPairingConfirmDialog(context, pending);
        if (mounted) setState(() => _pairingDialogOpen = false);
      });
    }

    // 后台服务没装或没在跑时，其他页面都无从谈起，先引导把它启用
    final service = state.serviceState;
    final needsService = state.link == LinkState.offline &&
        service != null &&
        service != ServiceState.running;

    return Scaffold(
      body: Row(
        children: [
          _Sidebar(current: _section, onSelect: (s) => setState(() => _section = s)),
          Expanded(
            child: needsService
                ? const ServicePage()
                : switch (_section) {
                    Section.history => const HistoryPage(),
                    Section.devices => const DevicesPage(),
                    Section.settings => const SettingsPage(),
                  },
          ),
        ],
      ),
    );
  }
}

class _Sidebar extends StatelessWidget {
  const _Sidebar({required this.current, required this.onSelect});

  final Section current;
  final ValueChanged<Section> onSelect;

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final p = context.palette;

    return Container(
      width: 200,
      decoration: BoxDecoration(
        color: p.sidebar,
        border: Border(right: BorderSide(color: p.border)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Mac 上顶部留白是红绿灯按钮所在的标题栏（内容铺满了整个窗口）；
          // Windows 用系统标题栏，内容从标题栏下面开始
          SizedBox(height: Platform.isMacOS ? 52 : Insets.lg),
          Padding(
            padding: const EdgeInsets.fromLTRB(Insets.lg, 0, Insets.lg, Insets.lg),
            child: Row(
              children: [
                const AppMark(size: 22),
                const SizedBox(width: Insets.sm),
                Text('CopySync', style: context.text.titleSmall),
              ],
            ),
          ),
          _NavItem(
            icon: AppIcons.clock,
            label: '复制记录',
            selected: current == Section.history,
            trailing: state.records.isEmpty ? null : '${state.records.length}',
            onTap: () => onSelect(Section.history),
          ),
          _NavItem(
            icon: AppIcons.deviceLaptop,
            label: '设备',
            selected: current == Section.devices,
            onTap: () => onSelect(Section.devices),
          ),
          _NavItem(
            icon: AppIcons.gear,
            label: '设置',
            selected: current == Section.settings,
            onTap: () => onSelect(Section.settings),
          ),
          const Spacer(),
          const _ConnectionFooter(),
        ],
      ),
    );
  }
}

class _NavItem extends StatefulWidget {
  const _NavItem({
    required this.icon,
    required this.label,
    required this.selected,
    required this.onTap,
    this.trailing,
  });

  final IconData icon;
  final String label;
  final bool selected;
  final VoidCallback onTap;
  final String? trailing;

  @override
  State<_NavItem> createState() => _NavItemState();
}

class _NavItemState extends State<_NavItem> {
  bool _hovered = false;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    final selected = widget.selected;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: Insets.sm, vertical: 1),
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        onEnter: (_) => setState(() => _hovered = true),
        onExit: (_) => setState(() => _hovered = false),
        child: GestureDetector(
          onTap: widget.onTap,
          child: Container(
            height: 30,
            padding: const EdgeInsets.symmetric(horizontal: 10),
            decoration: BoxDecoration(
              color: selected ? p.selection : (_hovered ? p.hover : Colors.transparent),
              borderRadius: BorderRadius.circular(Radii.sm),
            ),
            child: Row(
              children: [
                Icon(widget.icon, size: 16, color: selected ? p.accent : p.textDim),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(
                    widget.label,
                    style: context.text.labelLarge?.copyWith(
                      fontWeight: selected ? FontWeight.w600 : FontWeight.w500,
                    ),
                  ),
                ),
                if (widget.trailing != null)
                  Text(widget.trailing!, style: context.text.labelSmall),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

/// 侧边栏底部的连接状态。常驻可见，用户不必进设置页就知道能不能用。
class _ConnectionFooter extends StatelessWidget {
  const _ConnectionFooter();

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final p = context.palette;
    final signaling = state.status?.signalingConnected ?? false;
    final online = state.onlinePeerCount;

    final (color, label) = switch (state.link) {
      LinkState.connecting => (p.textFaint, '正在连接…'),
      LinkState.offline => (p.danger, state.error ?? '后台服务未运行'),
      LinkState.online when !signaling => (p.warning, '未连接服务器'),
      LinkState.online when online > 0 => (p.online, '$online 台设备在线'),
      LinkState.online => (p.online, '已连接服务器'),
    };

    return Container(
      padding: const EdgeInsets.fromLTRB(Insets.lg, Insets.md, Insets.lg, Insets.md),
      decoration: BoxDecoration(border: Border(top: BorderSide(color: p.border))),
      child: Row(
        children: [
          Container(
            width: 7,
            height: 7,
            decoration: BoxDecoration(color: color, shape: BoxShape.circle),
          ),
          const SizedBox(width: Insets.sm),
          Expanded(
            child: Text(
              label,
              style: context.text.bodySmall,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ),
        ],
      ),
    );
  }
}

import 'dart:async';
import 'dart:io';
import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../i18n.dart';
import '../icons.dart';
import '../main.dart';
import '../pages/devices_page.dart' show LanBlockedNotice, platformIcon;
import '../platform.dart';
import '../theme.dart';
import 'common.dart';
import 'pairing_dialog.dart';

Future<void> showNearbyDialog(BuildContext context) =>
    pushPairingSetupRoute(context, (_) => const NearbyDialog());

/// 搜索附近设备：本机在中间，附近同样打开了这个窗口的设备出现在周围，点一下发起配对。
///
/// 窗口开着时订阅后台服务的 WatchNearby：订阅期间本机在局域网里公开设备名，
/// 窗口一关订阅就断，本机随之撤下。点了对方之后与输入配对码一样，双方核对指纹再确认。
class NearbyDialog extends StatefulWidget {
  const NearbyDialog({super.key});

  @override
  State<NearbyDialog> createState() => _NearbyDialogState();
}

/// 设备在雷达上的位置（角度，0 是正右方，顺时针）。先占左右两侧，再补正上、正下，
/// 相邻两台的名字不会叠在一起；最多显示这么多台。
const _slotAngles = [-30.0, 150.0, 30.0, 210.0, -90.0, 90.0];

class _NearbyDialogState extends State<NearbyDialog> with SingleTickerProviderStateMixin {
  late final AnimationController _pulse =
      AnimationController(vsync: this, duration: const Duration(milliseconds: 2800));
  StreamSubscription<List<NearbyDevice>>? _sub;
  List<NearbyDevice> _devices = const [];

  /// 每台设备占一个固定的位置：别的设备出现或离开时，已有的不会挪动
  final Map<String, int> _slots = {};
  String? _requesting; // 正在联系的设备
  String? _message; // 出错时的提示
  bool _lost = false; // 与后台服务的订阅断了

  @override
  void initState() {
    super.initState();
    _watch();
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    // 系统设置里关掉了动画时，雷达静止
    if (MediaQuery.disableAnimationsOf(context)) {
      _pulse.stop();
    } else if (!_pulse.isAnimating) {
      _pulse.repeat();
    }
  }

  @override
  void dispose() {
    _sub?.cancel();
    _pulse.dispose();
    super.dispose();
  }

  void _watch() {
    _sub?.cancel();
    _lost = false;
    _message = null;
    try {
      _sub = AppScope.read(context).watchNearby().listen(
            _onDevices,
            onError: (Object e) => _onLost(AppState.describeError(e)),
            onDone: () => _onLost(appL10n.linkServiceDisconnected),
          );
    } catch (e) {
      _onLost(AppState.describeError(e));
    }
  }

  void _onLost(String message) {
    if (!mounted) return;
    setState(() {
      _lost = true;
      _message = message;
      _devices = const [];
    });
  }

  void _onDevices(List<NearbyDevice> list) {
    if (!mounted) return;
    final ids = {for (final d in list) d.id};
    _slots.removeWhere((id, _) => !ids.contains(id));
    final shown = <NearbyDevice>[];
    for (final d in list) {
      var slot = _slots[d.id];
      if (slot == null) {
        final used = _slots.values.toSet();
        slot = List.generate(_slotAngles.length, (i) => i).where((i) => !used.contains(i)).firstOrNull;
        if (slot == null) continue; // 位置占满了
        _slots[d.id] = slot;
      }
      shown.add(d);
    }
    setState(() => _devices = shown);
  }

  Future<void> _pair(NearbyDevice d) async {
    if (_requesting != null || d.paired) return;
    setState(() {
      _requesting = d.id;
      _message = null;
    });
    try {
      // 成功后 Shell 弹出核对指纹，并移除这个窗口
      await AppScope.read(context).pairNearby(d.id);
    } catch (e) {
      if (!mounted) return;
      final l = context.l10n;
      setState(() => _message = switch (AppState.pairingFailure(e)) {
            PairingFailure.notFound => l.nearbyNoReply(d.name),
            PairingFailure.noChannel => l.pairNoChannel,
            _ => AppState.describeError(e),
          });
    } finally {
      if (mounted) setState(() => _requesting = null);
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final p = context.palette;
    final l = context.l10n;
    final requesting = _devices.where((d) => d.id == _requesting).firstOrNull;
    final status = _message ??
        (requesting != null
            ? l.nearbyRequesting(requesting.name)
            : _devices.isEmpty
                ? l.nearbySearching
                : null);

    return Dialog(
      child: SizedBox(
        width: 460,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(Insets.xl, Insets.lg, Insets.md, 0),
              child: Row(
                children: [
                  Expanded(child: Text(l.nearbyTitle, style: context.text.titleMedium)),
                  IconAction(
                    icon: AppIcons.xmark,
                    tooltip: l.close,
                    onTap: () => Navigator.pop(context),
                  ),
                ],
              ),
            ),
            if (state.lanBlocked)
              const Padding(
                padding: EdgeInsets.fromLTRB(Insets.xl, Insets.md, Insets.xl, 0),
                child: LanBlockedNotice(),
              ),
            SizedBox(
              height: 300,
              child: state.lanEnabled
                  ? _radar(context, state)
                  : Padding(
                      padding: const EdgeInsets.all(Insets.xl),
                      child: EmptyState(
                        icon: AppIcons.radar,
                        title: l.nearbyLanOffTitle,
                        description: l.nearbyLanOffBody,
                      ),
                    ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.xl),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  if (state.lanEnabled) ...[
                    SizedBox(
                      height: 20,
                      child: status == null
                          ? null
                          : Text(
                              status,
                              textAlign: TextAlign.center,
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: context.text.bodyMedium?.copyWith(
                                color: _message != null ? p.danger : p.textDim,
                              ),
                            ),
                    ),
                    if (_lost)
                      Center(child: TextButton(onPressed: () => setState(_watch), child: Text(l.retry))),
                    const SizedBox(height: 2),
                    Text(l.nearbyHint, textAlign: TextAlign.center, style: context.text.bodySmall),
                  ],
                  const SizedBox(height: Insets.lg),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Icon(AppIcons.lockShield, size: 15, color: p.textDim),
                      const SizedBox(width: Insets.sm),
                      Expanded(
                        child: Text(
                          l.nearbyVisible(Wording.thisDeviceInText),
                          style: context.text.bodySmall,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _radar(BuildContext context, AppState state) {
    final p = context.palette;
    final selfPlatform = state.self?.platform ?? (Platform.isWindows ? 'windows' : 'darwin');
    return LayoutBuilder(
      builder: (context, box) {
        final center = box.biggest.center(Offset.zero);
        // 设备排在一个椭圆上：横向空间多，名字不容易挤在一起
        final rx = math.min(150.0, box.maxWidth / 2 - _NearbyNode.width / 2);
        final ry = math.min(100.0, box.maxHeight / 2 - 50);
        return Stack(
          clipBehavior: Clip.none,
          children: [
            Positioned.fill(
              child: RepaintBoundary(
                child: CustomPaint(
                  painter: _RadarPainter(pulse: _pulse, ring: p.separator, wave: p.accent),
                ),
              ),
            ),
            Positioned(
              left: center.dx - _NearbyNode.width / 2,
              top: center.dy - _NearbyNode.avatar / 2,
              width: _NearbyNode.width,
              child: _SelfNode(platform: selfPlatform),
            ),
            for (final d in _devices)
              Builder(builder: (context) {
                final a = _slotAngles[_slots[d.id]!] * math.pi / 180;
                final at = center + Offset(rx * math.cos(a), ry * math.sin(a));
                return Positioned(
                  key: ValueKey(d.id),
                  left: at.dx - _NearbyNode.width / 2,
                  top: at.dy - _NearbyNode.avatar / 2,
                  width: _NearbyNode.width,
                  child: _NearbyNode(
                    device: d,
                    busy: d.id == _requesting,
                    onTap: _requesting == null ? () => _pair(d) : null,
                  ),
                );
              }),
          ],
        );
      },
    );
  }
}

/// 雷达的底：几圈很淡的同心圆，两道波纹错开半个周期从本机向外扩散、逐渐变淡。
class _RadarPainter extends CustomPainter {
  _RadarPainter({required this.pulse, required this.ring, required this.wave}) : super(repaint: pulse);

  final Animation<double> pulse;
  final Color ring;
  final Color wave;

  @override
  void paint(Canvas canvas, Size size) {
    final c = size.center(Offset.zero);
    final maxR = size.shortestSide / 2;
    final rings = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1
      ..color = ring;
    for (final f in const [0.36, 0.68, 1.0]) {
      canvas.drawCircle(c, maxR * f, rings);
    }
    for (final offset in const [0.0, 0.5]) {
      final t = (pulse.value + offset) % 1.0;
      final r = 28 + (maxR - 28) * Curves.easeOut.transform(t);
      canvas.drawCircle(
        c,
        r,
        Paint()
          ..style = PaintingStyle.stroke
          ..strokeWidth = 1.5
          ..color = wave.withValues(alpha: 0.3 * (1 - t)),
      );
    }
  }

  @override
  bool shouldRepaint(_RadarPainter old) => old.ring != ring || old.wave != wave;
}

class _SelfNode extends StatelessWidget {
  const _SelfNode({required this.platform});

  final String platform;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: _NearbyNode.avatar,
          height: _NearbyNode.avatar,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: p.surface,
            border: Border.all(color: p.accent, width: 1.5),
          ),
          child: Icon(platformIcon(platform), size: 22, color: p.accent),
        ),
        const SizedBox(height: 6),
        Text(
          Wording.thisDevice,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          textAlign: TextAlign.center,
          style: context.text.labelMedium,
        ),
      ],
    );
  }
}

/// 雷达上的一台附近设备：圆形图标加名字。已配对的变淡、不能点。
class _NearbyNode extends StatefulWidget {
  const _NearbyNode({required this.device, required this.busy, this.onTap});

  static const width = 112.0;
  static const avatar = 46.0;

  final NearbyDevice device;
  final bool busy;
  final VoidCallback? onTap;

  @override
  State<_NearbyNode> createState() => _NearbyNodeState();
}

class _NearbyNodeState extends State<_NearbyNode> {
  bool _hover = false;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    final d = widget.device;
    final enabled = !d.paired && widget.onTap != null;
    final highlight = enabled && _hover;

    final node = Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        SizedBox.square(
          dimension: _NearbyNode.avatar,
          child: Stack(
            alignment: Alignment.center,
            children: [
              AnimatedContainer(
                duration: const Duration(milliseconds: 120),
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: p.surface,
                  border: Border.all(color: highlight ? p.accent : p.border, width: highlight ? 1.5 : 1),
                  boxShadow: const [BoxShadow(color: Color(0x14000000), blurRadius: 6, offset: Offset(0, 2))],
                ),
              ),
              Icon(platformIcon(d.platform), size: 22, color: highlight ? p.accent : p.textDim),
              if (widget.busy) const SizedBox.expand(child: CircularProgressIndicator(strokeWidth: 1.5)),
            ],
          ),
        ),
        const SizedBox(height: 6),
        Text(
          d.name,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          textAlign: TextAlign.center,
          style: context.text.labelMedium?.copyWith(color: p.text),
        ),
        if (d.paired) Text(context.l10n.nearbyPaired, style: context.text.labelSmall),
      ],
    );

    // 新出现的设备从小到大淡入
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0, end: 1),
      duration: const Duration(milliseconds: 260),
      curve: Curves.easeOut,
      builder: (context, v, child) => Opacity(
        opacity: v,
        child: Transform.scale(scale: 0.7 + 0.3 * v, child: child),
      ),
      child: Semantics(
        button: enabled,
        label: d.name,
        child: MouseRegion(
          cursor: enabled ? SystemMouseCursors.click : SystemMouseCursors.basic,
          onEnter: (_) => setState(() => _hover = true),
          onExit: (_) => setState(() => _hover = false),
          child: GestureDetector(
            behavior: HitTestBehavior.opaque,
            onTap: enabled ? widget.onTap : null,
            child: Opacity(opacity: d.paired ? 0.55 : 1, child: node),
          ),
        ),
      ),
    );
  }
}

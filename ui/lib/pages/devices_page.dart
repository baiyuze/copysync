import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';
import '../widgets/pairing_dialog.dart';

class DevicesPage extends StatelessWidget {
  const DevicesPage({super.key});

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final canPair = state.isOnline && (state.status?.signalingConnected ?? false);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PageHeader(
          title: '设备',
          subtitle: state.peers.isEmpty
              ? '配对之后，设备之间就会同步剪贴板'
              : '已配对 ${state.peers.length} 台，${state.onlinePeerCount} 台在线',
          actions: [
            FilledButton.icon(
              onPressed: canPair ? () => _startPairing(context) : null,
              icon: const Icon(CupertinoIcons.plus, size: 14),
              label: const Text('添加设备'),
            ),
          ],
        ),
        Expanded(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.xl),
            children: [
              if (state.self != null)
                GroupSection(
                  title: '这台 Mac',
                  footnote: '配对时，两台 Mac 会显示同样的两行安全指纹：这台的和对方的。逐字核对一致才能确认。',
                  child: GroupBox(children: [_DeviceRow(device: state.self!, isSelf: true)]),
                ),
              if (state.peers.isEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: Insets.xl),
                  child: EmptyState(
                    icon: CupertinoIcons.device_laptop,
                    title: '还没有配对的设备',
                    description: canPair
                        ? '在另一台电脑上打开 CopySync，用 6 位配对码把两台连起来。'
                        : '先在「设置」里连接信令服务器，才能配对设备。',
                    action: canPair
                        ? FilledButton(
                            onPressed: () => _startPairing(context),
                            child: const Text('添加设备'),
                          )
                        : null,
                  ),
                )
              else
                GroupSection(
                  title: '已配对的设备',
                  footnote: '直连：数据在两台设备之间直接传输。'
                      '中转：网络环境打不通直连时，经你的服务器转发，内容依然是端到端加密的。',
                  child: GroupBox(
                    dividerIndent: 52,
                    children: [for (final d in state.peers) _DeviceRow(device: d)],
                  ),
                ),
            ],
          ),
        ),
      ],
    );
  }

  void _startPairing(BuildContext context) =>
      showDialog(context: context, builder: (_) => const PairingDialog());
}

class _DeviceRow extends StatelessWidget {
  const _DeviceRow({required this.device, this.isSelf = false});

  final Device device;
  final bool isSelf;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    return Padding(
      padding: const EdgeInsets.fromLTRB(Insets.lg, Insets.md, Insets.sm, Insets.md),
      child: Row(
        children: [
          Icon(platformIcon(device.platform), size: 22, color: p.textDim),
          const SizedBox(width: Insets.md + 2),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(device.name, style: context.text.labelLarge),
                const SizedBox(height: 2),
                Text(
                  device.publicKeyFingerprint,
                  // 等宽：用户要逐字符比对，等宽能显著降低看错的概率
                  style: context.text.bodySmall?.copyWith(fontFamily: monoFamily, letterSpacing: 0.4),
                ),
              ],
            ),
          ),
          if (isSelf)
            Padding(
              padding: const EdgeInsets.only(right: Insets.sm),
              child: Text('本机', style: context.text.bodySmall),
            )
          else ...[
            StatusDot(
              color: device.online ? p.online : p.textFaint,
              label: device.online ? connectionLabel(device) : '离线',
            ),
            const SizedBox(width: Insets.sm),
            _DeviceMenu(device: device),
          ],
        ],
      ),
    );
  }
}

String connectionLabel(Device d) => switch (d.connection) {
      ConnectionKind.CONNECTION_KIND_DIRECT => '直连',
      ConnectionKind.CONNECTION_KIND_RELAY => '中转',
      _ => '在线',
    };

IconData platformIcon(String platform) => switch (platform) {
      'windows' => CupertinoIcons.desktopcomputer,
      _ => CupertinoIcons.device_laptop,
    };

class _DeviceMenu extends StatelessWidget {
  const _DeviceMenu({required this.device});

  final Device device;

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<String>(
      tooltip: '更多操作',
      padding: EdgeInsets.zero,
      position: PopupMenuPosition.under,
      itemBuilder: (ctx) => [
        const PopupMenuItem(value: 'copy', height: 34, child: Text('拷贝安全指纹')),
        PopupMenuItem(
          value: 'unpair',
          height: 34,
          child: Text('解除配对…', style: TextStyle(color: ctx.palette.danger)),
        ),
      ],
      onSelected: (v) => v == 'copy' ? _copy(context) : _unpair(context),
      child: const IgnorePointer(
        child: IconAction(icon: CupertinoIcons.ellipsis, tooltip: ''),
      ),
    );
  }

  Future<void> _copy(BuildContext context) async {
    await Clipboard.setData(ClipboardData(text: device.publicKeyFingerprint));
    if (context.mounted) showToast(context, '安全指纹已拷贝');
  }

  Future<void> _unpair(BuildContext context) async {
    final state = AppScope.read(context);
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('解除与「${device.name}」的配对？'),
        content: const Text('两台设备将不再同步。以后想恢复，需要重新配对并核对指纹。'),
        actions: [
          OutlinedButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: FilledButton.styleFrom(backgroundColor: ctx.palette.danger),
            child: const Text('解除配对'),
          ),
        ],
      ),
    );
    if (ok != true) return;
    try {
      await state.unpair(device.id);
      if (context.mounted) showToast(context, '已解除配对');
    } catch (e) {
      if (context.mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }
}

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../i18n.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';
import '../widgets/pairing_dialog.dart';
import '../icons.dart';
import '../platform.dart';

class DevicesPage extends StatelessWidget {
  const DevicesPage({super.key});

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final canPair = state.isOnline && (state.status?.signalingConnected ?? false);
    final l = context.l10n;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PageHeader(
          title: l.navDevices,
          subtitle: state.peers.isEmpty
              ? l.devicesSubtitleEmpty
              : l.devicesSubtitle(state.peers.length, state.onlinePeerCount),
          actions: [
            FilledButton.icon(
              onPressed: canPair ? () => _startPairing(context) : null,
              icon: Icon(AppIcons.plus, size: 14),
              label: Text(l.addDevice),
            ),
          ],
        ),
        Expanded(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.xl),
            children: [
              if (state.self != null)
                GroupSection(
                  title: Wording.thisDevice,
                  footnote: l.thisDeviceFootnote(Wording.bothDevicesInText),
                  child: GroupBox(children: [_DeviceRow(device: state.self!, isSelf: true)]),
                ),
              if (state.peers.isEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: Insets.xl),
                  child: EmptyState(
                    icon: AppIcons.deviceLaptop,
                    title: l.noDevicesTitle,
                    description: canPair ? l.noDevicesBody : l.noDevicesNeedServer,
                    action: canPair
                        ? FilledButton(
                            onPressed: () => _startPairing(context),
                            child: Text(l.addDevice),
                          )
                        : null,
                  ),
                )
              else
                GroupSection(
                  title: l.pairedDevices,
                  footnote: l.pairedFootnote,
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

  void _startPairing(BuildContext context) => showPairingDialog(context);
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
              child: Text(context.l10n.thisDeviceShort, style: context.text.bodySmall),
            )
          else ...[
            StatusDot(
              color: device.online ? p.online : p.textFaint,
              label: device.online ? connectionLabel(device) : context.l10n.offline,
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
      ConnectionKind.CONNECTION_KIND_DIRECT => appL10n.direct,
      ConnectionKind.CONNECTION_KIND_RELAY => appL10n.relay,
      _ => appL10n.online,
    };

IconData platformIcon(String platform) => switch (platform) {
      'windows' => AppIcons.desktopComputer,
      _ => AppIcons.deviceLaptop,
    };

class _DeviceMenu extends StatelessWidget {
  const _DeviceMenu({required this.device});

  final Device device;

  @override
  Widget build(BuildContext context) {
    final l = context.l10n;
    return PopupMenuButton<String>(
      tooltip: l.moreActions,
      padding: EdgeInsets.zero,
      position: PopupMenuPosition.under,
      itemBuilder: (ctx) => [
        PopupMenuItem(value: 'copy', height: 34, child: Text(l.copyFingerprint)),
        PopupMenuItem(
          value: 'unpair',
          height: 34,
          child: Text(l.unpairEllipsis, style: TextStyle(color: ctx.palette.danger)),
        ),
      ],
      onSelected: (v) => v == 'copy' ? _copy(context) : _unpair(context),
      child: IgnorePointer(
        child: IconAction(icon: AppIcons.ellipsis, tooltip: ''),
      ),
    );
  }

  Future<void> _copy(BuildContext context) async {
    await Clipboard.setData(ClipboardData(text: device.publicKeyFingerprint));
    if (context.mounted) showToast(context, context.l10n.fingerprintCopied);
  }

  Future<void> _unpair(BuildContext context) async {
    final state = AppScope.read(context);
    final l = context.l10n;
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(l.unpairTitle(device.name)),
        content: Text(l.unpairBody),
        actions: [
          OutlinedButton(onPressed: () => Navigator.pop(ctx, false), child: Text(l.cancel)),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: FilledButton.styleFrom(backgroundColor: ctx.palette.danger),
            child: Text(l.unpair),
          ),
        ],
      ),
    );
    if (ok != true) return;
    try {
      await state.unpair(device.id);
      if (context.mounted) showToast(context, l.unpaired);
    } catch (e) {
      if (context.mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }
}

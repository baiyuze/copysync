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
    // 同一网络里不需要服务器，后台服务在运行就能配对；连不上对方时由配对框给出原因
    final canPair = state.isOnline;
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
        if (state.lanBlocked)
          const Padding(
            padding: EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.md),
            child: LanBlockedNotice(),
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
                    description: l.noDevicesBody,
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
        // 局域网里找不到对方时（跨网段、两边都换了地址），手动告诉后台服务去哪里找
        if (!device.online)
          PopupMenuItem(value: 'address', height: 34, child: Text(l.enterAddressEllipsis)),
        PopupMenuItem(
          value: 'unpair',
          height: 34,
          child: Text(l.unpairEllipsis, style: TextStyle(color: ctx.palette.danger)),
        ),
      ],
      onSelected: (v) => switch (v) {
        'copy' => _copy(context),
        'address' => _setAddress(context),
        _ => _unpair(context),
      },
      child: IgnorePointer(
        child: IconAction(icon: AppIcons.ellipsis, tooltip: ''),
      ),
    );
  }

  Future<void> _copy(BuildContext context) async {
    await Clipboard.setData(ClipboardData(text: device.publicKeyFingerprint));
    if (context.mounted) showToast(context, context.l10n.fingerprintCopied);
  }

  Future<void> _setAddress(BuildContext context) async {
    final state = AppScope.read(context);
    final l = context.l10n;
    final controller = TextEditingController();
    final address = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(l.setAddressTitle(device.name)),
        content: SizedBox(
          width: 360,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(l.setAddressBody(device.name)),
              const SizedBox(height: Insets.lg),
              TextField(
                controller: controller,
                autofocus: true,
                style: const TextStyle(fontFamily: monoFamily, fontSize: 14),
                decoration: const InputDecoration(hintText: '192.168.1.20'),
                onSubmitted: (v) => Navigator.pop(ctx, v),
              ),
            ],
          ),
        ),
        actions: [
          OutlinedButton(onPressed: () => Navigator.pop(ctx), child: Text(l.cancel)),
          FilledButton(onPressed: () => Navigator.pop(ctx, controller.text), child: Text(l.save)),
        ],
      ),
    );
    controller.dispose();
    if (address == null || address.trim().isEmpty) return;
    try {
      await state.setDeviceAddress(device.id, address);
      if (context.mounted) showToast(context, l.addressSaved);
    } catch (e) {
      if (!context.mounted) return;
      final failure = AppState.pairingFailure(e);
      showToast(
        context,
        switch (failure) {
          PairingFailure.badAddress => l.addressInvalid,
          PairingFailure.noChannel => l.pairNoChannel,
          _ => AppState.describeError(e),
        },
        error: true,
      );
    }
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

/// macOS 拦下了发往局域网的包。不放行的话，同一网络里的设备连不上这台。
class LanBlockedNotice extends StatelessWidget {
  const LanBlockedNotice({super.key});

  @override
  Widget build(BuildContext context) {
    final l = context.l10n;
    return Notice(
      title: l.lanBlockedTitle,
      message: l.lanBlockedBody(Wording.thisDeviceInText),
      action: OutlinedButton(
        onPressed: () async {
          try {
            await AppScope.read(context).requestLocalNetworkPermission();
          } catch (e) {
            if (context.mounted) showToast(context, AppState.describeError(e), error: true);
          }
        },
        child: Text(l.openSystemSettings),
      ),
    );
  }
}

import 'dart:io';

import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';

import '../app_state.dart';
import '../background_service.dart';
import '../gen/pasteferry/v1/daemon.pb.dart';
import '../i18n.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';
import '../widgets/feedback_dialog.dart';
import '../icons.dart';
import '../platform.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  final _serverController = TextEditingController();
  final _nameController = TextEditingController();
  bool _serverDirty = false;
  bool _serverSaving = false;
  bool _nameDirty = false;
  String? _lastLoadedServer;
  String? _lastLoadedName;

  @override
  void dispose() {
    _serverController.dispose();
    _nameController.dispose();
    super.dispose();
  }

  /// 用户正在编辑时不要用服务端下发的值覆盖输入框，
  /// 否则打字打到一半会被"吃掉"。
  void _syncControllers(Config cfg) {
    final custom = cfg.customSignalingUrl.isNotEmpty
        ? cfg.customSignalingUrl
        : cfg.signalingUrl;
    if (!_serverDirty && custom != _lastLoadedServer) {
      _serverController.text = custom;
      _lastLoadedServer = custom;
    }
    if (!_nameDirty && cfg.deviceName != _lastLoadedName) {
      _nameController.text = cfg.deviceName;
      _lastLoadedName = cfg.deviceName;
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final cfg = state.config;

    if (cfg == null) {
      return const Center(child: CircularProgressIndicator(strokeWidth: 2.5));
    }
    _syncControllers(cfg);

    final p = context.palette;
    final l = context.l10n;
    final connected = state.status?.signalingConnected == true;
    // 服务器是可选的：没填地址不是故障，用中性色
    final (
      Color statusColor,
      String statusLabel,
      String? statusDesc,
    ) = switch ((state.serverConfigured, connected)) {
      (false, _) => (p.textFaint, l.notSet, null),
      (true, true) => (p.online, l.connected, null),
      (true, false) => (p.danger, l.notConnected, l.disconnectedDesc),
    };

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PageHeader(title: l.navSettings),
        Expanded(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(
              Insets.xl,
              0,
              Insets.xl,
              Insets.xxl,
            ),
            children: [
              GroupSection(
                title: l.server,
                child: GroupBox(
                  children: [
                    SettingRow(
                      title: l.useOfficialService,
                      description: l.officialServiceDesc,
                      child: SmallSwitch(
                        value: state.officialEnabled,
                        onChanged: _serverSaving
                            ? null
                            : (v) => _setOfficial(state, cfg, v),
                      ),
                    ),
                    SettingRow(
                      title: l.customServer,
                      description: l.signalingUrlDesc,
                      child: SizedBox(
                        width: 360,
                        child: Row(
                          children: [
                            Expanded(
                              child: TextField(
                                controller: _serverController,
                                style: context.text.bodyMedium?.copyWith(
                                  fontFamily: monoFamily,
                                  fontSize: 12,
                                ),
                                decoration: InputDecoration(
                                  hintText: l.signalingUrlPlaceholder,
                                ),
                                onChanged: (_) =>
                                    setState(() => _serverDirty = true),
                                onSubmitted: (_) => _saveServer(state, cfg),
                              ),
                            ),
                            const SizedBox(width: 8),
                            FilledButton(
                              onPressed: _serverSaving
                                  ? null
                                  : () => _saveServer(state, cfg),
                              child: Text(l.saveAndConnect),
                            ),
                          ],
                        ),
                      ),
                    ),
                    SettingRow(
                      title: l.publicStun,
                      description: l.publicStunDesc,
                      child: SmallSwitch(
                        value: !cfg.onlyOwnStun,
                        onChanged: (v) => _update(state, cfg..onlyOwnStun = !v),
                      ),
                    ),
                    SettingRow(
                      title: l.connectionStatus,
                      description: statusDesc,
                      child: StatusDot(color: statusColor, label: statusLabel),
                    ),
                  ],
                ),
              ),
              GroupSection(
                title: l.localNetwork,
                child: GroupBox(
                  children: [
                    SettingRow(
                      title: l.lanDiscovery,
                      description: l.lanDiscoveryDesc,
                      child: SmallSwitch(
                        value: !cfg.lanDisabled,
                        onChanged: (v) => _update(state, cfg..lanDisabled = !v),
                      ),
                    ),
                    if (Platform.isMacOS)
                      SettingRow(
                        title: l.macLocalNetwork,
                        description: l.macLocalNetworkDesc,
                        child: OutlinedButton(
                          onPressed: () => _allowLocalNetwork(state, null),
                          child: Text(l.openSystemSettings),
                        ),
                      ),
                    if (Platform.isWindows)
                      SettingRow(
                        title: l.allowLan,
                        description: l.allowLanDesc,
                        child: OutlinedButton(
                          onPressed: () =>
                              _allowLocalNetwork(state, l.allowLanDone),
                          child: Text(l.allow),
                        ),
                      ),
                  ],
                ),
              ),
              GroupSection(
                title: Wording.thisDevice,
                child: GroupBox(
                  children: [
                    SettingRow(
                      title: l.deviceName,
                      child: SizedBox(
                        width: 220,
                        child: TextField(
                          controller: _nameController,
                          decoration: InputDecoration(
                            suffixIcon: _nameDirty
                                ? IconButton(
                                    tooltip: l.save,
                                    icon: Icon(
                                      AppIcons.checkmark,
                                      size: 14,
                                      color: p.accent,
                                    ),
                                    onPressed: () => _saveName(state, cfg),
                                  )
                                : null,
                          ),
                          onChanged: (_) => setState(() => _nameDirty = true),
                          onSubmitted: (_) => _saveName(state, cfg),
                        ),
                      ),
                    ),
                    SettingRow(
                      title: l.language,
                      child: _Dropdown<String>(
                        value: cfg.language,
                        // 语言名用各自的写法，选错了语言的人也认得出自己的
                        items: [
                          MapEntry('', l.languageSystem),
                          const MapEntry('zh-Hans', '简体中文'),
                          const MapEntry('zh-Hant-TW', '繁體中文（台灣）'),
                          const MapEntry('zh-Hant-HK', '繁體中文（香港）'),
                          const MapEntry('en', 'English'),
                          const MapEntry('ja', '日本語'),
                        ],
                        onChanged: (v) => _update(state, cfg..language = v),
                      ),
                    ),
                  ],
                ),
              ),
              GroupSection(
                title: l.sync,
                child: GroupBox(
                  children: [
                    SettingRow(
                      title: l.autoSyncLimit,
                      description: l.autoSyncLimitDesc,
                      child: _ThresholdSelector(
                        value: cfg.autoSyncThresholdBytes.toInt(),
                        onChanged: (v) => _update(
                          state,
                          cfg..autoSyncThresholdBytes = Int64(v),
                        ),
                      ),
                    ),
                    SettingRow(
                      title: l.largeFileNotifications,
                      description:
                          switch (state.notificationStatus?.permission) {
                            'allowed' => l.notificationsAllowed,
                            'not_requested' => l.notificationsNotRequested,
                            'denied' => l.notificationsDenied,
                            'unknown' => l.notificationsUnknown,
                            _ => l.notificationsUnavailable,
                          },
                      child:
                          state.notificationStatus?.supported == true &&
                              !(Platform.isLinux &&
                                  state.notificationStatus?.permission ==
                                      'unknown')
                          ? TextButton(
                              onPressed: () async {
                                try {
                                  await state.requestNotifications();
                                } catch (e) {
                                  if (context.mounted) {
                                    showToast(
                                      context,
                                      AppState.describeError(e),
                                      error: true,
                                    );
                                  }
                                }
                              },
                              child: Text(
                                state.notificationStatus?.permission ==
                                        'not_requested'
                                    ? l.enableNotifications
                                    : l.notificationSettings,
                              ),
                            )
                          : TextButton(
                              onPressed: state.refreshNotificationStatus,
                              child: Text(l.retry),
                            ),
                    ),
                    SettingRow(
                      title: l.autoApply,
                      child: SmallSwitch(
                        value: cfg.autoApplyToClipboard,
                        onChanged: (v) =>
                            _update(state, cfg..autoApplyToClipboard = v),
                      ),
                    ),
                  ],
                ),
              ),
              GroupSection(
                title: l.whatToSync,
                child: GroupBox(
                  children: [
                    _ToggleRow(
                      title: l.syncFiles,
                      value: cfg.syncFile,
                      onChanged: (v) => _update(state, cfg..syncFile = v),
                    ),
                    _ToggleRow(
                      title: l.syncText,
                      value: cfg.syncText,
                      onChanged: (v) => _update(state, cfg..syncText = v),
                    ),
                    _ToggleRow(
                      title: l.syncRich,
                      value: cfg.syncHtml,
                      onChanged: (v) => _update(state, cfg..syncHtml = v),
                    ),
                    _ToggleRow(
                      title: l.syncImages,
                      value: cfg.syncImage,
                      onChanged: (v) => _update(state, cfg..syncImage = v),
                    ),
                  ],
                ),
              ),
              GroupSection(
                title: l.storage,
                footnote: l.storageFootnote,
                child: GroupBox(
                  children: [
                    SettingRow(
                      title: l.keepHistory,
                      child: _DurationSelector(
                        seconds: cfg.historyTtlSeconds.toInt(),
                        onChanged: (v) =>
                            _update(state, cfg..historyTtlSeconds = Int64(v)),
                      ),
                    ),
                    SettingRow(
                      title: l.keepCache,
                      description: l.keepCacheDesc,
                      child: _DurationSelector(
                        seconds: cfg.cacheTtlSeconds.toInt(),
                        onChanged: (v) =>
                            _update(state, cfg..cacheTtlSeconds = Int64(v)),
                      ),
                    ),
                    SettingRow(
                      title: l.cacheUsage,
                      child: Text(
                        humanBytes(state.status?.cacheBytesUsed.toInt() ?? 0),
                        style: context.text.bodyMedium?.copyWith(
                          color: p.textDim,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              if (state.service.available)
                GroupSection(
                  title: l.backgroundSync,
                  child: GroupBox(
                    children: [
                      SettingRow(
                        title: l.backgroundService,
                        description: l.backgroundServiceDesc,
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            OutlinedButton(
                              onPressed: () => _serviceAction(
                                state.restartService,
                                l.restarted,
                              ),
                              child: Text(l.restart),
                            ),
                            const SizedBox(width: Insets.sm),
                            OutlinedButton(
                              onPressed: () => _confirmDisable(state),
                              child: Text(l.disableEllipsis),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              GroupSection(
                title: l.about,
                child: GroupBox(
                  children: [
                    SettingRow(
                      title: '${l.appName} ${state.status?.version ?? ''}',
                      description: appBuild.isEmpty
                          ? l.aboutDesc
                          : '${l.aboutDesc}\n${l.buildDetails(appVersion, state.status?.version ?? '…', appBuild)}',
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          TextButton(
                            onPressed: () => _open('https://pasteferry.cc/'),
                            child: Text(l.homepage),
                          ),
                          TextButton(
                            onPressed: () => showFeedbackDialog(context),
                            child: Text(l.feedback),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  static void _open(String url) => openUrl(url);

  Future<void> _serviceAction(
    Future<void> Function() action,
    String message,
  ) async {
    try {
      await action();
      if (mounted) showToast(context, message);
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }

  Future<void> _confirmDisable(AppState state) async {
    final l = context.l10n;
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(l.disableTitle),
        content: Text(l.disableBody(Wording.thisDeviceInText)),
        actions: [
          OutlinedButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: Text(l.cancel),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: Text(l.disable),
          ),
        ],
      ),
    );
    if (ok == true) await _serviceAction(state.disableService, l.disabled);
  }

  Future<void> _saveServer(AppState state, Config cfg) async {
    if (_serverSaving) return;
    final l = context.l10n;
    final url = _serverController.text.trim();
    // 留空就是不用服务器，只在局域网里同步
    if (url.isNotEmpty &&
        !url.startsWith('ws://') &&
        !url.startsWith('wss://')) {
      showToast(context, l.urlMustBeWs, error: true);
      return;
    }
    final next = cfg.deepCopy()
      ..signalingUrl = url
      ..customSignalingUrl = url
      ..serviceMode = url.isEmpty
          ? ServiceMode.SERVICE_MODE_LAN
          : ServiceMode.SERVICE_MODE_CUSTOM;
    setState(() => _serverSaving = true);
    try {
      await state.updateConfig(next);
      if (!mounted) return;
      setState(() {
        _serverDirty = false;
        _lastLoadedServer = url;
      });
      showToast(context, url.isEmpty ? l.serverCleared : l.serverUpdated);
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    } finally {
      if (mounted) setState(() => _serverSaving = false);
    }
  }

  Future<void> _setOfficial(AppState state, Config cfg, bool value) async {
    if (value) {
      final ok = await showDialog<bool>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: Text(context.l10n.officialEnableTitle),
          content: Text(
            '${context.l10n.officialPrivacy}${cfg.serviceMode == ServiceMode.SERVICE_MODE_CUSTOM ? '\n\n${context.l10n.officialSwitchConfirm}' : ''}',
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(ctx, false),
              child: Text(context.l10n.cancel),
            ),
            FilledButton(
              onPressed: () => Navigator.pop(ctx, true),
              child: Text(context.l10n.useOfficialService),
            ),
          ],
        ),
      );
      if (ok != true || !mounted) return;
    }
    setState(() => _serverSaving = true);
    try {
      final next = cfg.deepCopy()
        ..serviceMode = (value
            ? ServiceMode.SERVICE_MODE_OFFICIAL
            : ServiceMode.SERVICE_MODE_LAN)
        ..signalingUrl = '';
      await state.updateConfig(next);
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    } finally {
      if (mounted) setState(() => _serverSaving = false);
    }
  }

  /// macOS 打开「本地网络」设置面板；Windows 弹出 UAC，在防火墙里放行后台服务。
  Future<void> _allowLocalNetwork(AppState state, String? doneMessage) async {
    try {
      await state.requestLocalNetworkPermission();
      if (mounted && doneMessage != null) showToast(context, doneMessage);
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }

  Future<void> _saveName(AppState state, Config cfg) async {
    final name = _nameController.text.trim();
    if (name.isEmpty) return;
    setState(() => _nameDirty = false);
    _lastLoadedName = name;
    await _update(
      state,
      cfg..deviceName = name,
      message: context.l10n.deviceNameUpdated,
    );
  }

  Future<void> _update(AppState state, Config next, {String? message}) async {
    try {
      await state.updateConfig(next);
      if (mounted && message != null) showToast(context, message);
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }
}

class _ToggleRow extends StatelessWidget {
  const _ToggleRow({
    required this.title,
    required this.value,
    required this.onChanged,
  });

  final String title;
  final bool value;
  final ValueChanged<bool> onChanged;

  @override
  Widget build(BuildContext context) {
    return SettingRow(
      title: title,
      child: SmallSwitch(value: value, onChanged: onChanged),
    );
  }
}

class _ThresholdSelector extends StatelessWidget {
  const _ThresholdSelector({required this.value, required this.onChanged});

  final int value;
  final ValueChanged<int> onChanged;

  static const _options = <int, String>{
    1 << 20: '1 MB',
    10 << 20: '10 MB',
    50 << 20: '50 MB',
    200 << 20: '200 MB',
    1024 << 20: '1 GB',
  };

  @override
  Widget build(BuildContext context) {
    // 当前值不在预设里时补一项，避免 DropdownButton 因找不到值而报错
    final options = Map<int, String>.from(_options);
    options.putIfAbsent(value, () => humanBytes(value));

    return _Dropdown<int>(
      value: value,
      items: options.entries.toList()..sort((a, b) => a.key.compareTo(b.key)),
      onChanged: onChanged,
    );
  }
}

class _DurationSelector extends StatelessWidget {
  const _DurationSelector({required this.seconds, required this.onChanged});

  final int seconds;
  final ValueChanged<int> onChanged;

  static const _options = [3600, 86400, 3 * 86400, 7 * 86400, 30 * 86400];

  @override
  Widget build(BuildContext context) {
    final l = context.l10n;
    final options = {for (final s in _options) s: durationLabel(l, s)};
    options.putIfAbsent(seconds, () => durationLabel(l, seconds));

    return _Dropdown<int>(
      value: seconds,
      items: options.entries.toList()..sort((a, b) => a.key.compareTo(b.key)),
      onChanged: onChanged,
    );
  }
}

/// macOS 弹出按钮：当前值 + 上下箭头，点开是菜单。
class _Dropdown<T> extends StatelessWidget {
  const _Dropdown({
    required this.value,
    required this.items,
    required this.onChanged,
  });

  final T value;
  final List<MapEntry<T, String>> items;
  final ValueChanged<T> onChanged;

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    final current = items.firstWhere(
      (e) => e.key == value,
      orElse: () => items.first,
    );
    return PopupMenuButton<T>(
      tooltip: '',
      initialValue: value,
      position: PopupMenuPosition.under,
      onSelected: onChanged,
      itemBuilder: (_) => [
        for (final e in items)
          PopupMenuItem(value: e.key, height: 32, child: Text(e.value)),
      ],
      child: Container(
        padding: const EdgeInsets.fromLTRB(10, 5, 8, 5),
        decoration: BoxDecoration(
          color: p.surface,
          borderRadius: BorderRadius.circular(Radii.sm),
          border: Border.all(color: p.border),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(current.value, style: context.text.labelLarge),
            const SizedBox(width: 6),
            Icon(AppIcons.chevronUpDown, size: 12, color: p.textDim),
          ],
        ),
      ),
    );
  }
}

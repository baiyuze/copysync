
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../i18n.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';
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
    if (!_serverDirty && cfg.signalingUrl != _lastLoadedServer) {
      _serverController.text = cfg.signalingUrl;
      _lastLoadedServer = cfg.signalingUrl;
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

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PageHeader(title: l.navSettings),
        Expanded(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.xxl),
            children: [
              GroupSection(
                title: l.server,
                footnote: l.serverFootnote,
                child: GroupBox(children: [
                  SettingRow(
                    title: l.signalingUrl,
                    description: l.signalingUrlDesc,
                    child: SizedBox(
                      width: 290,
                      child: TextField(
                        controller: _serverController,
                        style: context.text.bodyMedium?.copyWith(fontFamily: monoFamily, fontSize: 12),
                        decoration: InputDecoration(
                          hintText: l.signalingUrlPlaceholder,
                          suffixIcon: _serverDirty
                              ? IconButton(
                                  tooltip: l.save,
                                  icon: Icon(AppIcons.checkmark, size: 14, color: p.accent),
                                  onPressed: () => _saveServer(state, cfg),
                                )
                              : null,
                        ),
                        onChanged: (_) => setState(() => _serverDirty = true),
                        onSubmitted: (_) => _saveServer(state, cfg),
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
                    description: connected ? l.connectedDesc : l.disconnectedDesc,
                    child: StatusDot(
                      color: connected ? p.online : p.danger,
                      label: connected ? l.connected : l.notConnected,
                    ),
                  ),
                ]),
              ),
              GroupSection(
                title: Wording.thisDevice,
                child: GroupBox(children: [
                  SettingRow(
                    title: l.deviceName,
                    description: l.deviceNameDesc,
                    child: SizedBox(
                      width: 220,
                      child: TextField(
                        controller: _nameController,
                        decoration: InputDecoration(
                          suffixIcon: _nameDirty
                              ? IconButton(
                                  tooltip: l.save,
                                  icon: Icon(AppIcons.checkmark, size: 14, color: p.accent),
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
                        const MapEntry('en', 'English'),
                        const MapEntry('ja', '日本語'),
                      ],
                      onChanged: (v) => _update(state, cfg..language = v),
                    ),
                  ),
                ]),
              ),
              GroupSection(
                title: l.sync,
                footnote: l.syncFootnote,
                child: GroupBox(children: [
                  SettingRow(
                    title: l.autoSyncLimit,
                    description: l.autoSyncLimitDesc,
                    child: _ThresholdSelector(
                      value: cfg.autoSyncThresholdBytes.toInt(),
                      onChanged: (v) => _update(state, cfg..autoSyncThresholdBytes = Int64(v)),
                    ),
                  ),
                  SettingRow(
                    title: l.autoApply,
                    description: l.autoApplyDesc,
                    child: SmallSwitch(
                      value: cfg.autoApplyToClipboard,
                      onChanged: (v) => _update(state, cfg..autoApplyToClipboard = v),
                    ),
                  ),
                ]),
              ),
              GroupSection(
                title: l.whatToSync,
                child: GroupBox(children: [
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
                ]),
              ),
              GroupSection(
                title: l.storage,
                footnote: l.storageFootnote,
                child: GroupBox(children: [
                  SettingRow(
                    title: l.keepHistory,
                    child: _DurationSelector(
                      seconds: cfg.historyTtlSeconds.toInt(),
                      onChanged: (v) => _update(state, cfg..historyTtlSeconds = Int64(v)),
                    ),
                  ),
                  SettingRow(
                    title: l.keepCache,
                    description: l.keepCacheDesc,
                    child: _DurationSelector(
                      seconds: cfg.cacheTtlSeconds.toInt(),
                      onChanged: (v) => _update(state, cfg..cacheTtlSeconds = Int64(v)),
                    ),
                  ),
                  SettingRow(
                    title: l.cacheUsage,
                    child: Text(
                      humanBytes(state.status?.cacheBytesUsed.toInt() ?? 0),
                      style: context.text.bodyMedium?.copyWith(color: p.textDim),
                    ),
                  ),
                ]),
              ),
              if (state.service.available)
                GroupSection(
                  title: l.backgroundSync,
                  footnote: l.backgroundFootnote,
                  child: GroupBox(children: [
                    SettingRow(
                      title: l.backgroundService,
                      description: l.backgroundServiceDesc,
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          OutlinedButton(
                            onPressed: () => _serviceAction(state.restartService, l.restarted),
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
                  ]),
                ),
              GroupSection(
                title: l.about,
                child: GroupBox(children: [
                  SettingRow(
                    title: 'CopySync ${state.status?.version ?? ''}',
                    description: l.aboutDesc,
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        TextButton(
                          onPressed: () => _open('https://baiyuze.github.io/copysync/'),
                          child: Text(l.homepage),
                        ),
                        TextButton(
                          onPressed: () => _open('https://github.com/baiyuze/copysync/issues'),
                          child: Text(l.feedback),
                        ),
                      ],
                    ),
                  ),
                ]),
              ),
            ],
          ),
        ),
      ],
    );
  }

  static void _open(String url) => openUrl(url);

  Future<void> _serviceAction(Future<void> Function() action, String message) async {
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
          OutlinedButton(onPressed: () => Navigator.pop(ctx, false), child: Text(l.cancel)),
          FilledButton(onPressed: () => Navigator.pop(ctx, true), child: Text(l.disable)),
        ],
      ),
    );
    if (ok == true) await _serviceAction(state.disableService, l.disabled);
  }

  Future<void> _saveServer(AppState state, Config cfg) async {
    final url = _serverController.text.trim();
    if (!url.startsWith('ws://') && !url.startsWith('wss://')) {
      showToast(context, context.l10n.urlMustBeWs, error: true);
      return;
    }
    setState(() => _serverDirty = false);
    _lastLoadedServer = url;
    await _update(state, cfg..signalingUrl = url, message: context.l10n.serverUpdated);
  }

  Future<void> _saveName(AppState state, Config cfg) async {
    final name = _nameController.text.trim();
    if (name.isEmpty) return;
    setState(() => _nameDirty = false);
    _lastLoadedName = name;
    await _update(state, cfg..deviceName = name, message: context.l10n.deviceNameUpdated);
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
    final current = items.firstWhere((e) => e.key == value, orElse: () => items.first);
    return PopupMenuButton<T>(
      tooltip: '',
      initialValue: value,
      position: PopupMenuPosition.under,
      onSelected: onChanged,
      itemBuilder: (_) => [
        for (final e in items) PopupMenuItem(value: e.key, height: 32, child: Text(e.value)),
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

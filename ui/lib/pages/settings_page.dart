
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
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
    final connected = state.status?.signalingConnected == true;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const PageHeader(title: '设置'),
        Expanded(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(Insets.xl, 0, Insets.xl, Insets.xxl),
            children: [
              GroupSection(
                title: '服务器',
                footnote: '服务器只负责让设备找到彼此、协商直连。设备之间传输的内容它看不到。',
                child: GroupBox(children: [
                  SettingRow(
                    title: '信令服务器地址',
                    description: '按回车保存，保存后立即重新连接',
                    child: SizedBox(
                      width: 290,
                      child: TextField(
                        controller: _serverController,
                        style: context.text.bodyMedium?.copyWith(fontFamily: monoFamily, fontSize: 12),
                        decoration: InputDecoration(
                          hintText: 'ws://服务器地址:8787/signal',
                          suffixIcon: _serverDirty
                              ? IconButton(
                                  tooltip: '保存',
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
                    title: '用公共服务器探测网络出口',
                    description: '公司双线这类多出口网络里，能找到更多可以直连的路径。公共服务器只会看到你的公网 IP',
                    child: SmallSwitch(
                      value: !cfg.onlyOwnStun,
                      onChanged: (v) => _update(state, cfg..onlyOwnStun = !v),
                    ),
                  ),
                  SettingRow(
                    title: '连接状态',
                    description: connected ? '已连接到信令服务器' : '连不上服务器。检查地址是否正确、服务器是否在运行。',
                    child: StatusDot(
                      color: connected ? p.online : p.danger,
                      label: connected ? '已连接' : '未连接',
                    ),
                  ),
                ]),
              ),
              GroupSection(
                title: Wording.thisDevice,
                child: GroupBox(children: [
                  SettingRow(
                    title: '设备名称',
                    description: '显示在其他设备的列表里',
                    child: SizedBox(
                      width: 220,
                      child: TextField(
                        controller: _nameController,
                        decoration: InputDecoration(
                          suffixIcon: _nameDirty
                              ? IconButton(
                                  tooltip: '保存',
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
                ]),
              ),
              GroupSection(
                title: '同步',
                footnote: '小于阈值的内容在复制时直接推送到其他设备；超过阈值的只同步一条记录，'
                    '需要时在「复制记录」里点「拉取到本机」，避免大文件无谓地占用带宽和磁盘。',
                child: GroupBox(children: [
                  SettingRow(
                    title: '自动同步上限',
                    description: '超过这个大小的文件改为手动拉取',
                    child: _ThresholdSelector(
                      value: cfg.autoSyncThresholdBytes.toInt(),
                      onChanged: (v) => _update(state, cfg..autoSyncThresholdBytes = Int64(v)),
                    ),
                  ),
                  SettingRow(
                    title: '收到后直接写入剪贴板',
                    description: '关闭后，需要在「复制记录」里手动放入剪贴板',
                    child: SmallSwitch(
                      value: cfg.autoApplyToClipboard,
                      onChanged: (v) => _update(state, cfg..autoApplyToClipboard = v),
                    ),
                  ),
                ]),
              ),
              GroupSection(
                title: '同步哪些内容',
                child: GroupBox(children: [
                  _ToggleRow(
                    title: '文件与文件夹',
                    value: cfg.syncFile,
                    onChanged: (v) => _update(state, cfg..syncFile = v),
                  ),
                  _ToggleRow(
                    title: '纯文本',
                    value: cfg.syncText,
                    onChanged: (v) => _update(state, cfg..syncText = v),
                  ),
                  _ToggleRow(
                    title: '带格式的文本',
                    value: cfg.syncHtml,
                    onChanged: (v) => _update(state, cfg..syncHtml = v),
                  ),
                  _ToggleRow(
                    title: '图片与截图',
                    value: cfg.syncImage,
                    onChanged: (v) => _update(state, cfg..syncImage = v),
                  ),
                ]),
              ),
              GroupSection(
                title: '存储',
                footnote: '到期的记录与缓存文件会被自动清理。',
                child: GroupBox(children: [
                  SettingRow(
                    title: '记录保留',
                    child: _DurationSelector(
                      seconds: cfg.historyTtlSeconds.toInt(),
                      onChanged: (v) => _update(state, cfg..historyTtlSeconds = Int64(v)),
                    ),
                  ),
                  SettingRow(
                    title: '缓存文件保留',
                    description: '不会超过记录的保留时长',
                    child: _DurationSelector(
                      seconds: cfg.cacheTtlSeconds.toInt(),
                      onChanged: (v) => _update(state, cfg..cacheTtlSeconds = Int64(v)),
                    ),
                  ),
                  SettingRow(
                    title: '缓存占用',
                    child: Text(
                      humanBytes(state.status?.cacheBytesUsed.toInt() ?? 0),
                      style: context.text.bodyMedium?.copyWith(color: p.textDim),
                    ),
                  ),
                ]),
              ),
              if (state.service.available)
                GroupSection(
                  title: '后台同步',
                  footnote: '停用后不再同步，也不会开机启动。历史记录与配对关系会保留，重新打开 App 即可再次启用。',
                  child: GroupBox(children: [
                    SettingRow(
                      title: '后台服务',
                      description: '随系统登录自动启动，关掉窗口也照常同步',
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          OutlinedButton(
                            onPressed: () => _serviceAction(state.restartService, '后台服务已重新启动'),
                            child: const Text('重新启动'),
                          ),
                          const SizedBox(width: Insets.sm),
                          OutlinedButton(
                            onPressed: () => _confirmDisable(state),
                            child: const Text('停用…'),
                          ),
                        ],
                      ),
                    ),
                  ]),
                ),
              GroupSection(
                title: '关于',
                child: GroupBox(children: [
                  SettingRow(
                    title: 'CopySync ${state.status?.version ?? ''}',
                    description: '开源软件，MIT 许可',
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        TextButton(
                          onPressed: () => _open('https://baiyuze.github.io/copysync/'),
                          child: const Text('项目主页'),
                        ),
                        TextButton(
                          onPressed: () => _open('https://github.com/baiyuze/copysync/issues'),
                          child: const Text('反馈问题'),
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
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('停用后台同步？'),
        content: Text('${Wording.thisDeviceInText}将停止同步，也不再开机启动。历史记录和配对关系会保留。'),
        actions: [
          OutlinedButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('停用')),
        ],
      ),
    );
    if (ok == true) await _serviceAction(state.disableService, '后台同步已停用');
  }

  Future<void> _saveServer(AppState state, Config cfg) async {
    final url = _serverController.text.trim();
    if (!url.startsWith('ws://') && !url.startsWith('wss://')) {
      showToast(context, '地址需以 ws:// 或 wss:// 开头', error: true);
      return;
    }
    setState(() => _serverDirty = false);
    _lastLoadedServer = url;
    await _update(state, cfg..signalingUrl = url, message: '服务器地址已更新，正在重连');
  }

  Future<void> _saveName(AppState state, Config cfg) async {
    final name = _nameController.text.trim();
    if (name.isEmpty) return;
    setState(() => _nameDirty = false);
    _lastLoadedName = name;
    await _update(state, cfg..deviceName = name, message: '设备名称已更新');
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

  static const _options = <int, String>{
    3600: '1 小时',
    86400: '1 天',
    3 * 86400: '3 天',
    7 * 86400: '7 天',
    30 * 86400: '30 天',
  };

  @override
  Widget build(BuildContext context) {
    final options = Map<int, String>.from(_options);
    options.putIfAbsent(seconds, () => '${(seconds / 3600).round()} 小时');

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

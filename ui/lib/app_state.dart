import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:grpc/grpc.dart';

import 'background_service.dart';
import 'daemon_client.dart';
import 'gen/copysync/v1/daemon.pb.dart';
import 'gen/copysync/v1/daemon.pbgrpc.dart' show DaemonServiceClient;

/// daemon 连接的生命周期状态。
enum LinkState { connecting, online, offline }

/// 应用的单一状态源。
///
/// 所有与 daemon 的交互都经过这里：页面只读状态、调方法，不直接碰 gRPC。
/// 这样断线重连、事件合并这些琐事集中在一处，页面代码保持干净。
class AppState extends ChangeNotifier {
  AppState({BackgroundService? service}) : service = service ?? BackgroundService();

  final BackgroundService service;

  DaemonConnection? _conn;
  StreamSubscription<Event>? _events;
  Timer? _retry;

  LinkState _link = LinkState.connecting;
  String? _error;
  Status? _status;
  Config? _config;
  Device? _self;
  List<Device> _peers = const [];
  final List<ClipRecord> _records = [];

  /// 待用户核对指纹的配对请求。非空时界面弹出确认。
  Device? _pendingPairing;

  /// 连不上后台服务时，它的安装状态。界面据此决定显示欢迎页还是重启提示。
  ServiceState? _serviceState;
  bool _serviceChecked = false;

  LinkState get link => _link;
  String? get error => _error;
  Status? get status => _status;
  Config? get config => _config;
  Device? get self => _self;
  List<Device> get peers => List.unmodifiable(_peers);
  List<ClipRecord> get records => List.unmodifiable(_records);
  Device? get pendingPairing => _pendingPairing;
  ServiceState? get serviceState => _serviceState;

  bool get isOnline => _link == LinkState.online;
  int get onlinePeerCount => _peers.where((p) => p.online).length;

  @override
  void dispose() {
    _retry?.cancel();
    _events?.cancel();
    _conn?.close();
    super.dispose();
  }

  // ─────────────────────────── 连接 ───────────────────────────

  Future<void> connect() async {
    _retry?.cancel();
    _setLink(LinkState.connecting, null);

    try {
      final conn = await DaemonConnection.connect();
      _conn = conn;

      // 首屏需要的数据一次取齐，避免界面分几次闪烁
      final results = await Future.wait([
        conn.client.listHistory(ListHistoryRequest(limit: 200)),
        conn.client.getConfig(Empty()),
        conn.client.getStatus(Empty()),
      ]);

      _records
        ..clear()
        ..addAll((results[0] as ListHistoryResponse).records);
      _config = results[1] as Config;
      _status = results[2] as Status;

      await _loadDevices();

      _events = conn.client.subscribe(SubscribeRequest()).listen(
            _onEvent,
            onError: (Object e) => _scheduleRetry('连接中断：${_describe(e)}'),
            onDone: () => _scheduleRetry('后台服务已断开'),
          );

      _serviceState = null;
      _setLink(LinkState.online, null);
      _ensureServiceCurrent();
    } on DaemonNotRunning {
      await _refreshServiceState();
      _scheduleRetry('后台服务未运行');
    } catch (e) {
      await _refreshServiceState();
      _scheduleRetry(_describe(e));
    }
  }

  Future<void> _refreshServiceState() async {
    try {
      _serviceState = await service.state();
    } catch (_) {
      _serviceState = null;
    }
  }

  /// 连上的后台服务不一定是这个 App 内置的那份，每次启动核对一次：
  /// - 登录项指向别处（旧版脚本装的、或 App 换了位置）：改为指向当前 App
  /// - 覆盖安装了新版 App：launchd 里跑的还是旧进程，按新的可执行文件重启
  /// 两种情况都会让当前连接断开，随后的自动重连会连上新进程。
  Future<void> _ensureServiceCurrent() async {
    if (!service.available || _serviceChecked) return;
    _serviceChecked = true;
    try {
      final s = await service.state();
      if (s == ServiceState.stale) {
        await service.install();
        return;
      }
      final running = _status?.version ?? '';
      if (s == ServiceState.running &&
          appVersion != 'dev' &&
          running.isNotEmpty &&
          running != appVersion) {
        await service.restart();
      }
    } catch (_) {
      // 核对失败不影响使用：旧进程仍在正常工作
    }
  }

  /// 截图与界面测试用：直接填入数据，不连接后台服务。
  @visibleForTesting
  void debugSeed({
    required Status status,
    required Config config,
    required Device self,
    List<Device> peers = const [],
    List<ClipRecord> records = const [],
    Device? pendingPairing,
    LinkState link = LinkState.online,
    ServiceState? serviceState,
  }) {
    _status = status;
    _config = config;
    _self = self;
    _peers = peers;
    _records
      ..clear()
      ..addAll(records);
    _pendingPairing = pendingPairing;
    _serviceState = serviceState;
    _link = link;
    notifyListeners();
  }

  // ─────────────────────────── 后台服务 ───────────────────────────

  /// 注册后台服务并连接。用于首次启用，也用于修复被移动过的 App。
  Future<void> enableService() async {
    await service.install();
    await _refreshServiceState();
    await connect();
  }

  Future<void> restartService() async {
    await service.restart();
    await connect();
  }

  /// 停用后台服务：关闭同步并取消开机自启。
  Future<void> disableService() async {
    _retry?.cancel();
    await service.uninstall();
    await _refreshServiceState();
    _scheduleRetry('后台服务已停用');
  }

  Future<void> _loadDevices() async {
    final conn = _conn;
    if (conn == null) return;
    try {
      final resp = await conn.client.listDevices(Empty());
      _self = resp.self;
      _peers = resp.peers;
      notifyListeners();
    } on GrpcError catch (e) {
      // 配对能力尚未接入时返回 Unimplemented，属预期情况
      if (e.code != StatusCode.unimplemented) rethrow;
    }
  }

  void _scheduleRetry(String message) {
    _events?.cancel();
    _events = null;
    _conn?.close();
    _conn = null;
    _setLink(LinkState.offline, message);
    // daemon 可能正在重启，自动重连省去用户手动操作
    _retry = Timer(const Duration(seconds: 3), connect);
  }

  void _setLink(LinkState state, String? error) {
    _link = state;
    _error = error;
    notifyListeners();
  }

  void _onEvent(Event ev) {
    final added = ev.whichPayload();
    switch (added) {
      case Event_Payload.clipAdded:
        _records.insert(0, ev.clipAdded);
        if (_records.length > 500) _records.removeLast();

      case Event_Payload.clipUpdated:
        final i = _records.indexWhere((r) => r.id == ev.clipUpdated.id);
        if (i >= 0) {
          _records[i] = ev.clipUpdated;
        } else {
          _records.insert(0, ev.clipUpdated);
        }

      case Event_Payload.clipRemoved:
        _records.removeWhere((r) => r.id == ev.clipRemoved);

      case Event_Payload.deviceChanged:
        _onDeviceChanged(ev.deviceChanged);

      case Event_Payload.statusChanged:
        _status = ev.statusChanged;

      case Event_Payload.progress:
        final i = _records.indexWhere((r) => r.id == ev.progress.clipId);
        if (i >= 0) {
          _records[i] = _records[i]..status = ClipStatus.CLIP_STATUS_FETCHING;
        }

      default:
        break;
    }
    notifyListeners();
  }

  void _onDeviceChanged(Device d) {
    // 带 pairingSession 的不是已配对设备，而是一条待确认的配对请求
    if (d.pairingSession.isNotEmpty) {
      _pendingPairing = d;
      return;
    }
    final i = _peers.indexWhere((p) => p.id == d.id);
    final next = List<Device>.from(_peers);
    if (i >= 0) {
      next[i] = d;
    } else {
      next.add(d);
    }
    _peers = next;
  }

  // ─────────────────────────── 操作 ───────────────────────────

  DaemonServiceClient get _client {
    final c = _conn;
    if (c == null) throw StateError('后台服务未连接');
    return c.client;
  }

  Future<String> createPairingCode() async {
    final resp = await _client.createPairingCode(Empty());
    return resp.code;
  }

  /// 兑换配对码，返回待核对指纹的对端信息。
  Future<Device> redeemPairingCode(String code) async {
    final resp = await _client.redeemPairingCode(
      RedeemPairingCodeRequest(code: code.trim().toUpperCase()),
    );
    _pendingPairing = resp.peer..pairingSession = resp.pairingSession;
    notifyListeners();
    return _pendingPairing!;
  }

  Future<void> confirmPairing(String session, bool accept) async {
    await _client.confirmPairing(
      ConfirmPairingRequest(pairingSession: session, accept: accept),
    );
    _pendingPairing = null;
    await _loadDevices();
  }

  void dismissPendingPairing() {
    _pendingPairing = null;
    notifyListeners();
  }

  Future<void> unpair(String deviceId) async {
    await _client.unpair(UnpairRequest(deviceId: deviceId));
    _peers = _peers.where((p) => p.id != deviceId).toList();
    notifyListeners();
  }

  Future<void> updateConfig(Config next) async {
    _config = await _client.updateConfig(next);
    notifyListeners();
  }

  Future<void> deleteRecords({List<String>? ids, bool all = false}) async {
    await _client.deleteHistory(
      DeleteHistoryRequest(ids: ids ?? const [], all: all),
    );
    if (all) {
      _records.clear();
    } else if (ids != null) {
      _records.removeWhere((r) => ids.contains(r.id));
    }
    notifyListeners();
  }

  Future<void> fetch(String clipId) =>
      _client.fetch(FetchRequest(clipId: clipId));

  Future<void> applyToClipboard(String clipId) =>
      _client.applyToClipboard(ApplyToClipboardRequest(clipId: clipId));

  Future<void> requestClipboardPermission() =>
      _client.requestClipboardPermission(Empty());

  /// 把 gRPC 异常翻译成用户能看懂的话，而不是抛一串堆栈。
  static String _describe(Object e) {
    if (e is GrpcError) {
      if (e.code == StatusCode.unimplemented) return '该功能尚未接入';
      if (e.code == StatusCode.unavailable) return '无法连接后台服务';
      return e.message ?? e.codeName;
    }
    if (e is DaemonNotRunning) return '后台服务未运行';
    if (e is StateError) return e.message;
    return e.toString();
  }

  static String describeError(Object e) => _describe(e);
}

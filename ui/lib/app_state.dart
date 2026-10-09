import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:grpc/grpc.dart';

import 'background_service.dart';
import 'daemon_client.dart';
import 'i18n.dart';
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

  /// 正在传输的记录的进度（0–1），只在界面上展示，传完即移除。
  final Map<String, double> _progress = {};

  /// 待用户核对指纹的配对请求。非空时界面弹出确认。
  Device? _pendingPairing;
  final Set<String> _handledPairingSessions = {};

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

  /// 传输进度；总量未知或尚未开始时为 null。
  double? progressOf(String clipId) => _progress[clipId];
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
            onError: (Object e) => _scheduleRetry(appL10n.linkConnectionLost(_describe(e))),
            onDone: () => _scheduleRetry(appL10n.linkServiceDisconnected),
          );

      _serviceState = null;
      _setLink(LinkState.online, null);
      _ensureServiceCurrent();
    } on DaemonNotRunning {
      await _refreshServiceState();
      _scheduleRetry(appL10n.linkServiceNotRunning);
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
    _scheduleRetry(appL10n.linkServiceDisabled);
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
      // 同一条记录会被推送多次（开始接收时一次、接收完成时一次），一律按 id 合并，
      // 否则界面上会出现两行，其中「接收中」那行的进度条永远停不下来
      case Event_Payload.clipAdded:
        _upsert(ev.clipAdded);

      case Event_Payload.clipUpdated:
        _upsert(ev.clipUpdated);

      case Event_Payload.clipRemoved:
        _records.removeWhere((r) => r.id == ev.clipRemoved);

      case Event_Payload.deviceChanged:
        _onDeviceChanged(ev.deviceChanged);

      case Event_Payload.statusChanged:
        _status = ev.statusChanged;

      case Event_Payload.progress:
        _onProgress(ev.progress);

      default:
        break;
    }
    notifyListeners();
  }

  void _upsert(ClipRecord r) {
    final i = _records.indexWhere((x) => x.id == r.id);
    if (i >= 0) {
      _records[i] = r;
    } else {
      _records.insert(0, r);
      if (_records.length > 500) _records.removeLast();
    }
    // 传完（或失败）了就不再显示进度
    if (r.status != ClipStatus.CLIP_STATUS_FETCHING) _progress.remove(r.id);
  }

  /// 进度事件只用于「从对端接收」的记录：本机复制的记录在发送时也会产生进度，
  /// 若把它们也标成「接收中」，就再也没有事件能把状态改回来。
  void _onProgress(TransferProgress p) {
    final i = _records.indexWhere((r) => r.id == p.clipId);
    if (i < 0 || _records[i].outgoing) return;
    final r = _records[i];
    if (r.status == ClipStatus.CLIP_STATUS_REMOTE_ONLY) {
      // 用户点了「拉取到本机」，开始传了
      _records[i] = r.deepCopy()..status = ClipStatus.CLIP_STATUS_FETCHING;
    } else if (r.status != ClipStatus.CLIP_STATUS_FETCHING) {
      return; // 已经传完的记录，迟到的进度事件不能把它拉回「接收中」
    }
    if (p.total > 0) {
      _progress[p.clipId] = (p.transferred.toInt() / p.total.toInt()).clamp(0.0, 1.0);
    }
  }

  /// 测试用：直接注入一条后台服务的事件。
  @visibleForTesting
  void debugEvent(Event ev) => _onEvent(ev);

  void _onDeviceChanged(Device d) {
    // 带 pairingSession 的不是已配对设备，而是一条待确认的配对请求
    if (d.pairingSession.isNotEmpty) {
      _setPendingPairing(d);
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
    if (c == null) throw StateError(appL10n.errServiceNotConnected);
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
    final peer = resp.peer..pairingSession = resp.pairingSession;
    _setPendingPairing(peer);
    notifyListeners();
    return peer;
  }

  Future<void> confirmPairing(String session, bool accept) async {
    await _client.confirmPairing(
      ConfirmPairingRequest(pairingSession: session, accept: accept),
    );
    dismissPendingPairing(session);
    await _loadDevices();
  }

  void _setPendingPairing(Device peer) {
    // 事件流和兑换 RPC 返回同一个会话，先后顺序不固定。
    if (!_handledPairingSessions.contains(peer.pairingSession)) {
      _pendingPairing = peer;
    }
  }

  void dismissPendingPairing([String? session]) {
    final handled = session ?? _pendingPairing?.pairingSession;
    if (handled != null) _handledPairingSessions.add(handled);
    // 确认旧会话期间可能来了另一条请求，不能把新请求一并清掉。
    if (_pendingPairing?.pairingSession == handled) _pendingPairing = null;
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

  Future<void> fetchImagePreview(String clipId) =>
      _client.fetch(FetchRequest(clipId: clipId, preserveClipboard: true));

  Future<void> applyToClipboard(String clipId) =>
      _client.applyToClipboard(ApplyToClipboardRequest(clipId: clipId));

  Future<String> imagePreviewPath(String clipId) async {
    final response = await _client.getImagePreview(GetImagePreviewRequest(clipId: clipId),
        options: CallOptions(timeout: const Duration(seconds: 10)));
    return response.path;
  }

  Future<void> requestClipboardPermission() =>
      _client.requestClipboardPermission(Empty());

  /// 把 gRPC 异常翻译成用户能看懂的话，而不是抛一串堆栈。
  static String _describe(Object e) {
    if (e is GrpcError) {
      if (e.code == StatusCode.unimplemented) return appL10n.errUnimplemented;
      if (e.code == StatusCode.unavailable) return appL10n.errUnavailable;
      return e.message ?? e.codeName;
    }
    if (e is DaemonNotRunning) return appL10n.linkServiceNotRunning;
    if (e is StateError) return e.message;
    return e.toString();
  }

  static String describeError(Object e) => _describe(e);
}

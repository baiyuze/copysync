import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:grpc/grpc.dart';

import 'gen/copysync/v1/daemon.pbgrpc.dart';

/// daemon 写在数据目录下的握手信息。
///
/// daemon 监听 127.0.0.1 的随机端口，把端口与一次性 token 写进 daemon.json。
/// 用 token 是因为本机任意进程都能连上回环端口，仅"只监听 127.0.0.1"不够。
class DaemonEndpoint {
  const DaemonEndpoint({
    required this.port,
    required this.token,
    required this.pid,
    required this.version,
  });

  final int port;
  final String token;
  final int pid;
  final String version;

  factory DaemonEndpoint.fromJson(Map<String, dynamic> json) => DaemonEndpoint(
        port: json['port'] as int,
        token: json['token'] as String,
        pid: json['pid'] as int? ?? 0,
        version: json['version'] as String? ?? '',
      );

  /// 数据目录。必须与 Go 侧 config.DefaultPaths() 完全一致，
  /// 否则 UI 找不到 daemon。
  ///
  /// 刻意不用 path_provider：它在 macOS 上返回的是
  /// `~/Library/Application Support/<bundle_id>`，与 daemon 用的
  /// `~/Library/Application Support/CopySync` 并不是同一个位置。
  static String defaultDataDir() {
    if (Platform.isMacOS) {
      final home = Platform.environment['HOME'] ?? '';
      return '$home/Library/Application Support/CopySync';
    }
    if (Platform.isWindows) {
      final base = Platform.environment['LOCALAPPDATA'] ??
          Platform.environment['APPDATA'] ??
          '';
      return '$base\\CopySync';
    }
    final home = Platform.environment['HOME'] ?? '';
    return '$home/.config/copysync';
  }

  static File endpointFile([String? dataDir]) =>
      File('${dataDir ?? defaultDataDir()}${Platform.pathSeparator}daemon.json');

  static Future<DaemonEndpoint> read([String? dataDir]) async {
    final file = endpointFile(dataDir);
    if (!await file.exists()) {
      throw DaemonNotRunning('未找到 ${file.path}');
    }
    final raw = jsonDecode(await file.readAsString()) as Map<String, dynamic>;
    return DaemonEndpoint.fromJson(raw);
  }
}

class DaemonNotRunning implements Exception {
  DaemonNotRunning(this.detail);
  final String detail;
  @override
  String toString() => 'daemon 未运行（$detail）';
}

/// 与 daemon 的一条 gRPC 连接。
class DaemonConnection {
  DaemonConnection._(this._channel, this.client, this.endpoint);

  final ClientChannel _channel;
  final DaemonServiceClient client;
  final DaemonEndpoint endpoint;

  static Future<DaemonConnection> connect({String? dataDir}) async {
    final ep = await DaemonEndpoint.read(dataDir);
    final channel = ClientChannel(
      '127.0.0.1',
      port: ep.port,
      options: const ChannelOptions(
        credentials: ChannelCredentials.insecure(),
        // daemon 在本机，握手应当很快；连不上就尽快报错而不是干等
        connectionTimeout: Duration(seconds: 5),
      ),
    );
    final client = DaemonServiceClient(
      channel,
      options: CallOptions(metadata: {'x-copysync-token': ep.token}),
    );
    return DaemonConnection._(channel, client, ep);
  }

  Future<void> close() => _channel.shutdown();
}

// This is a generated file - do not edit.
//
// Generated from copysync/v1/daemon.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'daemon.pb.dart' as $0;

export 'daemon.pb.dart';

@$pb.GrpcServiceName('copysync.v1.DaemonService')
class DaemonServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DaemonServiceClient(super.channel, {super.options, super.interceptors});

  /// 订阅事件流。UI 启动后第一件事，断线需重连。
  $grpc.ResponseStream<$0.Event> subscribe(
    $0.SubscribeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$subscribe, $async.Stream.fromIterable([request]),
        options: options);
  }

  $grpc.ResponseFuture<$0.ListHistoryResponse> listHistory(
    $0.ListHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listHistory, request, options: options);
  }

  $grpc.ResponseFuture<$0.Empty> deleteHistory(
    $0.DeleteHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteHistory, request, options: options);
  }

  /// 拉取仅有元数据的大文件记录。立即返回，进度走事件流。
  $grpc.ResponseFuture<$0.Empty> fetch(
    $0.FetchRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$fetch, request, options: options);
  }

  /// 把某条已就绪的记录重新放进本机剪贴板
  $grpc.ResponseFuture<$0.Empty> applyToClipboard(
    $0.ApplyToClipboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$applyToClipboard, request, options: options);
  }

  /// 只返回已就绪图片在本机的路径，不改动剪贴板，也不把路径发给对端。
  $grpc.ResponseFuture<$0.GetImagePreviewResponse> getImagePreview(
    $0.GetImagePreviewRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getImagePreview, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListDevicesResponse> listDevices(
    $0.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDevices, request, options: options);
  }

  $grpc.ResponseFuture<$0.CreatePairingCodeResponse> createPairingCode(
    $0.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPairingCode, request, options: options);
  }

  $grpc.ResponseFuture<$0.RedeemPairingCodeResponse> redeemPairingCode(
    $0.RedeemPairingCodeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$redeemPairingCode, request, options: options);
  }

  /// 配对双方各自确认指纹后才真正互信
  $grpc.ResponseFuture<$0.Empty> confirmPairing(
    $0.ConfirmPairingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$confirmPairing, request, options: options);
  }

  $grpc.ResponseFuture<$0.Empty> unpair(
    $0.UnpairRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$unpair, request, options: options);
  }

  $grpc.ResponseFuture<$0.Config> getConfig(
    $0.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConfig, request, options: options);
  }

  $grpc.ResponseFuture<$0.Config> updateConfig(
    $0.Config request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateConfig, request, options: options);
  }

  $grpc.ResponseFuture<$0.Status> getStatus(
    $0.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getStatus, request, options: options);
  }

  /// 引导用户授予剪贴板权限（打开系统设置面板）
  $grpc.ResponseFuture<$0.Empty> requestClipboardPermission(
    $0.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$requestClipboardPermission, request,
        options: options);
  }

  /// 网络出口与打洞诊断：本机有几个出口、各自的公网地址，见 copysync-cli nat
  $grpc.ResponseFuture<$0.NetworkInfo> getNetwork(
    $0.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getNetwork, request, options: options);
  }

  // method descriptors

  static final _$subscribe = $grpc.ClientMethod<$0.SubscribeRequest, $0.Event>(
      '/copysync.v1.DaemonService/Subscribe',
      ($0.SubscribeRequest value) => value.writeToBuffer(),
      $0.Event.fromBuffer);
  static final _$listHistory =
      $grpc.ClientMethod<$0.ListHistoryRequest, $0.ListHistoryResponse>(
          '/copysync.v1.DaemonService/ListHistory',
          ($0.ListHistoryRequest value) => value.writeToBuffer(),
          $0.ListHistoryResponse.fromBuffer);
  static final _$deleteHistory =
      $grpc.ClientMethod<$0.DeleteHistoryRequest, $0.Empty>(
          '/copysync.v1.DaemonService/DeleteHistory',
          ($0.DeleteHistoryRequest value) => value.writeToBuffer(),
          $0.Empty.fromBuffer);
  static final _$fetch = $grpc.ClientMethod<$0.FetchRequest, $0.Empty>(
      '/copysync.v1.DaemonService/Fetch',
      ($0.FetchRequest value) => value.writeToBuffer(),
      $0.Empty.fromBuffer);
  static final _$applyToClipboard =
      $grpc.ClientMethod<$0.ApplyToClipboardRequest, $0.Empty>(
          '/copysync.v1.DaemonService/ApplyToClipboard',
          ($0.ApplyToClipboardRequest value) => value.writeToBuffer(),
          $0.Empty.fromBuffer);
  static final _$getImagePreview =
      $grpc.ClientMethod<$0.GetImagePreviewRequest, $0.GetImagePreviewResponse>(
          '/copysync.v1.DaemonService/GetImagePreview',
          ($0.GetImagePreviewRequest value) => value.writeToBuffer(),
          $0.GetImagePreviewResponse.fromBuffer);
  static final _$listDevices =
      $grpc.ClientMethod<$0.Empty, $0.ListDevicesResponse>(
          '/copysync.v1.DaemonService/ListDevices',
          ($0.Empty value) => value.writeToBuffer(),
          $0.ListDevicesResponse.fromBuffer);
  static final _$createPairingCode =
      $grpc.ClientMethod<$0.Empty, $0.CreatePairingCodeResponse>(
          '/copysync.v1.DaemonService/CreatePairingCode',
          ($0.Empty value) => value.writeToBuffer(),
          $0.CreatePairingCodeResponse.fromBuffer);
  static final _$redeemPairingCode = $grpc.ClientMethod<
          $0.RedeemPairingCodeRequest, $0.RedeemPairingCodeResponse>(
      '/copysync.v1.DaemonService/RedeemPairingCode',
      ($0.RedeemPairingCodeRequest value) => value.writeToBuffer(),
      $0.RedeemPairingCodeResponse.fromBuffer);
  static final _$confirmPairing =
      $grpc.ClientMethod<$0.ConfirmPairingRequest, $0.Empty>(
          '/copysync.v1.DaemonService/ConfirmPairing',
          ($0.ConfirmPairingRequest value) => value.writeToBuffer(),
          $0.Empty.fromBuffer);
  static final _$unpair = $grpc.ClientMethod<$0.UnpairRequest, $0.Empty>(
      '/copysync.v1.DaemonService/Unpair',
      ($0.UnpairRequest value) => value.writeToBuffer(),
      $0.Empty.fromBuffer);
  static final _$getConfig = $grpc.ClientMethod<$0.Empty, $0.Config>(
      '/copysync.v1.DaemonService/GetConfig',
      ($0.Empty value) => value.writeToBuffer(),
      $0.Config.fromBuffer);
  static final _$updateConfig = $grpc.ClientMethod<$0.Config, $0.Config>(
      '/copysync.v1.DaemonService/UpdateConfig',
      ($0.Config value) => value.writeToBuffer(),
      $0.Config.fromBuffer);
  static final _$getStatus = $grpc.ClientMethod<$0.Empty, $0.Status>(
      '/copysync.v1.DaemonService/GetStatus',
      ($0.Empty value) => value.writeToBuffer(),
      $0.Status.fromBuffer);
  static final _$requestClipboardPermission =
      $grpc.ClientMethod<$0.Empty, $0.Empty>(
          '/copysync.v1.DaemonService/RequestClipboardPermission',
          ($0.Empty value) => value.writeToBuffer(),
          $0.Empty.fromBuffer);
  static final _$getNetwork = $grpc.ClientMethod<$0.Empty, $0.NetworkInfo>(
      '/copysync.v1.DaemonService/GetNetwork',
      ($0.Empty value) => value.writeToBuffer(),
      $0.NetworkInfo.fromBuffer);
}

@$pb.GrpcServiceName('copysync.v1.DaemonService')
abstract class DaemonServiceBase extends $grpc.Service {
  $core.String get $name => 'copysync.v1.DaemonService';

  DaemonServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SubscribeRequest, $0.Event>(
        'Subscribe',
        subscribe_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.SubscribeRequest.fromBuffer(value),
        ($0.Event value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListHistoryRequest, $0.ListHistoryResponse>(
            'ListHistory',
            listHistory_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListHistoryRequest.fromBuffer(value),
            ($0.ListHistoryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteHistoryRequest, $0.Empty>(
        'DeleteHistory',
        deleteHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteHistoryRequest.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.FetchRequest, $0.Empty>(
        'Fetch',
        fetch_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.FetchRequest.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ApplyToClipboardRequest, $0.Empty>(
        'ApplyToClipboard',
        applyToClipboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ApplyToClipboardRequest.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetImagePreviewRequest,
            $0.GetImagePreviewResponse>(
        'GetImagePreview',
        getImagePreview_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetImagePreviewRequest.fromBuffer(value),
        ($0.GetImagePreviewResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.Empty, $0.ListDevicesResponse>(
        'ListDevices',
        listDevices_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Empty.fromBuffer(value),
        ($0.ListDevicesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.Empty, $0.CreatePairingCodeResponse>(
        'CreatePairingCode',
        createPairingCode_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Empty.fromBuffer(value),
        ($0.CreatePairingCodeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RedeemPairingCodeRequest,
            $0.RedeemPairingCodeResponse>(
        'RedeemPairingCode',
        redeemPairingCode_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RedeemPairingCodeRequest.fromBuffer(value),
        ($0.RedeemPairingCodeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConfirmPairingRequest, $0.Empty>(
        'ConfirmPairing',
        confirmPairing_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ConfirmPairingRequest.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UnpairRequest, $0.Empty>(
        'Unpair',
        unpair_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UnpairRequest.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.Empty, $0.Config>(
        'GetConfig',
        getConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Empty.fromBuffer(value),
        ($0.Config value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.Config, $0.Config>(
        'UpdateConfig',
        updateConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Config.fromBuffer(value),
        ($0.Config value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.Empty, $0.Status>(
        'GetStatus',
        getStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Empty.fromBuffer(value),
        ($0.Status value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.Empty, $0.Empty>(
        'RequestClipboardPermission',
        requestClipboardPermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Empty.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.Empty, $0.NetworkInfo>(
        'GetNetwork',
        getNetwork_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Empty.fromBuffer(value),
        ($0.NetworkInfo value) => value.writeToBuffer()));
  }

  $async.Stream<$0.Event> subscribe_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SubscribeRequest> $request) async* {
    yield* subscribe($call, await $request);
  }

  $async.Stream<$0.Event> subscribe(
      $grpc.ServiceCall call, $0.SubscribeRequest request);

  $async.Future<$0.ListHistoryResponse> listHistory_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListHistoryRequest> $request) async {
    return listHistory($call, await $request);
  }

  $async.Future<$0.ListHistoryResponse> listHistory(
      $grpc.ServiceCall call, $0.ListHistoryRequest request);

  $async.Future<$0.Empty> deleteHistory_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteHistoryRequest> $request) async {
    return deleteHistory($call, await $request);
  }

  $async.Future<$0.Empty> deleteHistory(
      $grpc.ServiceCall call, $0.DeleteHistoryRequest request);

  $async.Future<$0.Empty> fetch_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.FetchRequest> $request) async {
    return fetch($call, await $request);
  }

  $async.Future<$0.Empty> fetch(
      $grpc.ServiceCall call, $0.FetchRequest request);

  $async.Future<$0.Empty> applyToClipboard_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ApplyToClipboardRequest> $request) async {
    return applyToClipboard($call, await $request);
  }

  $async.Future<$0.Empty> applyToClipboard(
      $grpc.ServiceCall call, $0.ApplyToClipboardRequest request);

  $async.Future<$0.GetImagePreviewResponse> getImagePreview_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetImagePreviewRequest> $request) async {
    return getImagePreview($call, await $request);
  }

  $async.Future<$0.GetImagePreviewResponse> getImagePreview(
      $grpc.ServiceCall call, $0.GetImagePreviewRequest request);

  $async.Future<$0.ListDevicesResponse> listDevices_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Empty> $request) async {
    return listDevices($call, await $request);
  }

  $async.Future<$0.ListDevicesResponse> listDevices(
      $grpc.ServiceCall call, $0.Empty request);

  $async.Future<$0.CreatePairingCodeResponse> createPairingCode_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Empty> $request) async {
    return createPairingCode($call, await $request);
  }

  $async.Future<$0.CreatePairingCodeResponse> createPairingCode(
      $grpc.ServiceCall call, $0.Empty request);

  $async.Future<$0.RedeemPairingCodeResponse> redeemPairingCode_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RedeemPairingCodeRequest> $request) async {
    return redeemPairingCode($call, await $request);
  }

  $async.Future<$0.RedeemPairingCodeResponse> redeemPairingCode(
      $grpc.ServiceCall call, $0.RedeemPairingCodeRequest request);

  $async.Future<$0.Empty> confirmPairing_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ConfirmPairingRequest> $request) async {
    return confirmPairing($call, await $request);
  }

  $async.Future<$0.Empty> confirmPairing(
      $grpc.ServiceCall call, $0.ConfirmPairingRequest request);

  $async.Future<$0.Empty> unpair_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.UnpairRequest> $request) async {
    return unpair($call, await $request);
  }

  $async.Future<$0.Empty> unpair(
      $grpc.ServiceCall call, $0.UnpairRequest request);

  $async.Future<$0.Config> getConfig_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Empty> $request) async {
    return getConfig($call, await $request);
  }

  $async.Future<$0.Config> getConfig($grpc.ServiceCall call, $0.Empty request);

  $async.Future<$0.Config> updateConfig_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Config> $request) async {
    return updateConfig($call, await $request);
  }

  $async.Future<$0.Config> updateConfig(
      $grpc.ServiceCall call, $0.Config request);

  $async.Future<$0.Status> getStatus_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Empty> $request) async {
    return getStatus($call, await $request);
  }

  $async.Future<$0.Status> getStatus($grpc.ServiceCall call, $0.Empty request);

  $async.Future<$0.Empty> requestClipboardPermission_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Empty> $request) async {
    return requestClipboardPermission($call, await $request);
  }

  $async.Future<$0.Empty> requestClipboardPermission(
      $grpc.ServiceCall call, $0.Empty request);

  $async.Future<$0.NetworkInfo> getNetwork_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Empty> $request) async {
    return getNetwork($call, await $request);
  }

  $async.Future<$0.NetworkInfo> getNetwork(
      $grpc.ServiceCall call, $0.Empty request);
}

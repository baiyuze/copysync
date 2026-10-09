// This is a generated file - do not edit.
//
// Generated from copysync/v1/daemon.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'daemon.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'daemon.pbenum.dart';

class Empty extends $pb.GeneratedMessage {
  factory Empty() => Empty._();

  Empty._();

  factory Empty.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Empty()..mergeFromBuffer(data, registry);
  factory Empty.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Empty()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Empty',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: Empty.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Empty clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Empty copyWith(void Function(Empty) updates) =>
      super.copyWith((message) => updates(message as Empty)) as Empty;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Empty() / Empty.new instead')
  static Empty create() => Empty._();
  static $pb.GeneratedMessage $_createMessage() => Empty._();
  @$core.override
  Empty createEmptyInstance() => Empty._();
  @$core.pragma('dart2js:noInline')
  static Empty getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Empty>(Empty.$_createMessage);
  static Empty? _defaultInstance;
}

class ClipItem extends $pb.GeneratedMessage {
  factory ClipItem({
    $core.String? name,
    $fixnum.Int64? size,
    $core.bool? isDir,
    $core.String? contentType,
  }) {
    final result = ClipItem._();
    if (name != null) result.name = name;
    if (size != null) result.size = size;
    if (isDir != null) result.isDir = isDir;
    if (contentType != null) result.contentType = contentType;
    return result;
  }

  ClipItem._();

  factory ClipItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ClipItem()..mergeFromBuffer(data, registry);
  factory ClipItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ClipItem()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ClipItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: ClipItem.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aInt64(2, _omitFieldNames ? '' : 'size')
    ..aOB(3, _omitFieldNames ? '' : 'isDir')
    ..aOS(4, _omitFieldNames ? '' : 'contentType')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClipItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClipItem copyWith(void Function(ClipItem) updates) =>
      super.copyWith((message) => updates(message as ClipItem)) as ClipItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ClipItem() / ClipItem.new instead')
  static ClipItem create() => ClipItem._();
  static $pb.GeneratedMessage $_createMessage() => ClipItem._();
  @$core.override
  ClipItem createEmptyInstance() => ClipItem._();
  @$core.pragma('dart2js:noInline')
  static ClipItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ClipItem>(ClipItem.$_createMessage);
  static ClipItem? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get size => $_getI64(1);
  @$pb.TagNumber(2)
  set size($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSize() => $_has(1);
  @$pb.TagNumber(2)
  void clearSize() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get isDir => $_getBF(2);
  @$pb.TagNumber(3)
  set isDir($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasIsDir() => $_has(2);
  @$pb.TagNumber(3)
  void clearIsDir() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get contentType => $_getSZ(3);
  @$pb.TagNumber(4)
  set contentType($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasContentType() => $_has(3);
  @$pb.TagNumber(4)
  void clearContentType() => $_clearField(4);
}

class ClipRecord extends $pb.GeneratedMessage {
  factory ClipRecord({
    $core.String? id,
    ClipKind? kind,
    ClipStatus? status,
    $core.String? originDeviceId,
    $core.String? originDeviceName,
    $core.bool? outgoing,
    $core.Iterable<ClipItem>? items,
    $fixnum.Int64? totalSize,
    $core.String? textPreview,
    $fixnum.Int64? createdAtUnix,
    $fixnum.Int64? expiresAtUnix,
    $core.String? error,
  }) {
    final result = ClipRecord._();
    if (id != null) result.id = id;
    if (kind != null) result.kind = kind;
    if (status != null) result.status = status;
    if (originDeviceId != null) result.originDeviceId = originDeviceId;
    if (originDeviceName != null) result.originDeviceName = originDeviceName;
    if (outgoing != null) result.outgoing = outgoing;
    if (items != null) result.items.addAll(items);
    if (totalSize != null) result.totalSize = totalSize;
    if (textPreview != null) result.textPreview = textPreview;
    if (createdAtUnix != null) result.createdAtUnix = createdAtUnix;
    if (expiresAtUnix != null) result.expiresAtUnix = expiresAtUnix;
    if (error != null) result.error = error;
    return result;
  }

  ClipRecord._();

  factory ClipRecord.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ClipRecord()..mergeFromBuffer(data, registry);
  factory ClipRecord.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ClipRecord()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ClipRecord',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: ClipRecord.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aE<ClipKind>(2, _omitFieldNames ? '' : 'kind',
        enumValues: ClipKind.values)
    ..aE<ClipStatus>(3, _omitFieldNames ? '' : 'status',
        enumValues: ClipStatus.values)
    ..aOS(4, _omitFieldNames ? '' : 'originDeviceId')
    ..aOS(5, _omitFieldNames ? '' : 'originDeviceName')
    ..aOB(6, _omitFieldNames ? '' : 'outgoing')
    ..pPM<ClipItem>(7, _omitFieldNames ? '' : 'items',
        subBuilder: ClipItem.$_createMessage)
    ..aInt64(8, _omitFieldNames ? '' : 'totalSize')
    ..aOS(9, _omitFieldNames ? '' : 'textPreview')
    ..aInt64(10, _omitFieldNames ? '' : 'createdAtUnix')
    ..aInt64(11, _omitFieldNames ? '' : 'expiresAtUnix')
    ..aOS(12, _omitFieldNames ? '' : 'error')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClipRecord clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClipRecord copyWith(void Function(ClipRecord) updates) =>
      super.copyWith((message) => updates(message as ClipRecord)) as ClipRecord;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ClipRecord() / ClipRecord.new instead')
  static ClipRecord create() => ClipRecord._();
  static $pb.GeneratedMessage $_createMessage() => ClipRecord._();
  @$core.override
  ClipRecord createEmptyInstance() => ClipRecord._();
  @$core.pragma('dart2js:noInline')
  static ClipRecord getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ClipRecord>(ClipRecord.$_createMessage);
  static ClipRecord? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  ClipKind get kind => $_getN(1);
  @$pb.TagNumber(2)
  set kind(ClipKind value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasKind() => $_has(1);
  @$pb.TagNumber(2)
  void clearKind() => $_clearField(2);

  @$pb.TagNumber(3)
  ClipStatus get status => $_getN(2);
  @$pb.TagNumber(3)
  set status(ClipStatus value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasStatus() => $_has(2);
  @$pb.TagNumber(3)
  void clearStatus() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get originDeviceId => $_getSZ(3);
  @$pb.TagNumber(4)
  set originDeviceId($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasOriginDeviceId() => $_has(3);
  @$pb.TagNumber(4)
  void clearOriginDeviceId() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get originDeviceName => $_getSZ(4);
  @$pb.TagNumber(5)
  set originDeviceName($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasOriginDeviceName() => $_has(4);
  @$pb.TagNumber(5)
  void clearOriginDeviceName() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.bool get outgoing => $_getBF(5);
  @$pb.TagNumber(6)
  set outgoing($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasOutgoing() => $_has(5);
  @$pb.TagNumber(6)
  void clearOutgoing() => $_clearField(6);

  @$pb.TagNumber(7)
  $pb.PbList<ClipItem> get items => $_getList(6);

  @$pb.TagNumber(8)
  $fixnum.Int64 get totalSize => $_getI64(7);
  @$pb.TagNumber(8)
  set totalSize($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(8)
  $core.bool hasTotalSize() => $_has(7);
  @$pb.TagNumber(8)
  void clearTotalSize() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get textPreview => $_getSZ(8);
  @$pb.TagNumber(9)
  set textPreview($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasTextPreview() => $_has(8);
  @$pb.TagNumber(9)
  void clearTextPreview() => $_clearField(9);

  @$pb.TagNumber(10)
  $fixnum.Int64 get createdAtUnix => $_getI64(9);
  @$pb.TagNumber(10)
  set createdAtUnix($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(10)
  $core.bool hasCreatedAtUnix() => $_has(9);
  @$pb.TagNumber(10)
  void clearCreatedAtUnix() => $_clearField(10);

  @$pb.TagNumber(11)
  $fixnum.Int64 get expiresAtUnix => $_getI64(10);
  @$pb.TagNumber(11)
  set expiresAtUnix($fixnum.Int64 value) => $_setInt64(10, value);
  @$pb.TagNumber(11)
  $core.bool hasExpiresAtUnix() => $_has(10);
  @$pb.TagNumber(11)
  void clearExpiresAtUnix() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get error => $_getSZ(11);
  @$pb.TagNumber(12)
  set error($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasError() => $_has(11);
  @$pb.TagNumber(12)
  void clearError() => $_clearField(12);
}

class SubscribeRequest extends $pb.GeneratedMessage {
  factory SubscribeRequest({
    $fixnum.Int64? sinceUnix,
  }) {
    final result = SubscribeRequest._();
    if (sinceUnix != null) result.sinceUnix = sinceUnix;
    return result;
  }

  SubscribeRequest._();

  factory SubscribeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SubscribeRequest()..mergeFromBuffer(data, registry);
  factory SubscribeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SubscribeRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: SubscribeRequest.$_createMessage)
    ..aInt64(1, _omitFieldNames ? '' : 'sinceUnix')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeRequest copyWith(void Function(SubscribeRequest) updates) =>
      super.copyWith((message) => updates(message as SubscribeRequest))
          as SubscribeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use SubscribeRequest() / SubscribeRequest.new instead')
  static SubscribeRequest create() => SubscribeRequest._();
  static $pb.GeneratedMessage $_createMessage() => SubscribeRequest._();
  @$core.override
  SubscribeRequest createEmptyInstance() => SubscribeRequest._();
  @$core.pragma('dart2js:noInline')
  static SubscribeRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SubscribeRequest>(
          SubscribeRequest.$_createMessage);
  static SubscribeRequest? _defaultInstance;

  /// 断线重连时带上最后收到的记录时间，daemon 补发这之后的记录
  @$pb.TagNumber(1)
  $fixnum.Int64 get sinceUnix => $_getI64(0);
  @$pb.TagNumber(1)
  set sinceUnix($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSinceUnix() => $_has(0);
  @$pb.TagNumber(1)
  void clearSinceUnix() => $_clearField(1);
}

class TransferProgress extends $pb.GeneratedMessage {
  factory TransferProgress({
    $core.String? clipId,
    $fixnum.Int64? transferred,
    $fixnum.Int64? total,
    $core.double? bytesPerSecond,
  }) {
    final result = TransferProgress._();
    if (clipId != null) result.clipId = clipId;
    if (transferred != null) result.transferred = transferred;
    if (total != null) result.total = total;
    if (bytesPerSecond != null) result.bytesPerSecond = bytesPerSecond;
    return result;
  }

  TransferProgress._();

  factory TransferProgress.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TransferProgress()..mergeFromBuffer(data, registry);
  factory TransferProgress.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TransferProgress()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransferProgress',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: TransferProgress.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'clipId')
    ..aInt64(2, _omitFieldNames ? '' : 'transferred')
    ..aInt64(3, _omitFieldNames ? '' : 'total')
    ..aD(4, _omitFieldNames ? '' : 'bytesPerSecond')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransferProgress clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransferProgress copyWith(void Function(TransferProgress) updates) =>
      super.copyWith((message) => updates(message as TransferProgress))
          as TransferProgress;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TransferProgress() / TransferProgress.new instead')
  static TransferProgress create() => TransferProgress._();
  static $pb.GeneratedMessage $_createMessage() => TransferProgress._();
  @$core.override
  TransferProgress createEmptyInstance() => TransferProgress._();
  @$core.pragma('dart2js:noInline')
  static TransferProgress getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TransferProgress>(
          TransferProgress.$_createMessage);
  static TransferProgress? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get clipId => $_getSZ(0);
  @$pb.TagNumber(1)
  set clipId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasClipId() => $_has(0);
  @$pb.TagNumber(1)
  void clearClipId() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get transferred => $_getI64(1);
  @$pb.TagNumber(2)
  set transferred($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTransferred() => $_has(1);
  @$pb.TagNumber(2)
  void clearTransferred() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get total => $_getI64(2);
  @$pb.TagNumber(3)
  set total($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTotal() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotal() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.double get bytesPerSecond => $_getN(3);
  @$pb.TagNumber(4)
  set bytesPerSecond($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(4)
  $core.bool hasBytesPerSecond() => $_has(3);
  @$pb.TagNumber(4)
  void clearBytesPerSecond() => $_clearField(4);
}

class Device extends $pb.GeneratedMessage {
  factory Device({
    $core.String? id,
    $core.String? name,
    $core.String? platform,
    $core.bool? online,
    ConnectionKind? connection,
    $core.String? publicKeyFingerprint,
    $fixnum.Int64? pairedAtUnix,
    $core.String? pairingSession,
  }) {
    final result = Device._();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (platform != null) result.platform = platform;
    if (online != null) result.online = online;
    if (connection != null) result.connection = connection;
    if (publicKeyFingerprint != null)
      result.publicKeyFingerprint = publicKeyFingerprint;
    if (pairedAtUnix != null) result.pairedAtUnix = pairedAtUnix;
    if (pairingSession != null) result.pairingSession = pairingSession;
    return result;
  }

  Device._();

  factory Device.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Device()..mergeFromBuffer(data, registry);
  factory Device.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Device()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Device',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: Device.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'platform')
    ..aOB(4, _omitFieldNames ? '' : 'online')
    ..aE<ConnectionKind>(5, _omitFieldNames ? '' : 'connection',
        enumValues: ConnectionKind.values)
    ..aOS(6, _omitFieldNames ? '' : 'publicKeyFingerprint')
    ..aInt64(7, _omitFieldNames ? '' : 'pairedAtUnix')
    ..aOS(8, _omitFieldNames ? '' : 'pairingSession')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Device clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Device copyWith(void Function(Device) updates) =>
      super.copyWith((message) => updates(message as Device)) as Device;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Device() / Device.new instead')
  static Device create() => Device._();
  static $pb.GeneratedMessage $_createMessage() => Device._();
  @$core.override
  Device createEmptyInstance() => Device._();
  @$core.pragma('dart2js:noInline')
  static Device getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Device>(Device.$_createMessage);
  static Device? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get platform => $_getSZ(2);
  @$pb.TagNumber(3)
  set platform($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPlatform() => $_has(2);
  @$pb.TagNumber(3)
  void clearPlatform() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get online => $_getBF(3);
  @$pb.TagNumber(4)
  set online($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasOnline() => $_has(3);
  @$pb.TagNumber(4)
  void clearOnline() => $_clearField(4);

  @$pb.TagNumber(5)
  ConnectionKind get connection => $_getN(4);
  @$pb.TagNumber(5)
  set connection(ConnectionKind value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasConnection() => $_has(4);
  @$pb.TagNumber(5)
  void clearConnection() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get publicKeyFingerprint => $_getSZ(5);
  @$pb.TagNumber(6)
  set publicKeyFingerprint($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasPublicKeyFingerprint() => $_has(5);
  @$pb.TagNumber(6)
  void clearPublicKeyFingerprint() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get pairedAtUnix => $_getI64(6);
  @$pb.TagNumber(7)
  set pairedAtUnix($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasPairedAtUnix() => $_has(6);
  @$pb.TagNumber(7)
  void clearPairedAtUnix() => $_clearField(7);

  /// 非空表示这是一条**待用户核对指纹**的配对请求，而非已配对设备。
  /// 用此值调 ConfirmPairing 完成或拒绝配对。
  /// 配对码的发起方是被动收到兑换通知的，靠这个字段把确认入口推给界面。
  @$pb.TagNumber(8)
  $core.String get pairingSession => $_getSZ(7);
  @$pb.TagNumber(8)
  set pairingSession($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasPairingSession() => $_has(7);
  @$pb.TagNumber(8)
  void clearPairingSession() => $_clearField(8);
}

enum Event_Payload {
  clipAdded,
  clipUpdated,
  clipRemoved,
  progress,
  deviceChanged,
  permissionChanged,
  statusChanged,
  notSet
}

class Event extends $pb.GeneratedMessage {
  factory Event({
    ClipRecord? clipAdded,
    ClipRecord? clipUpdated,
    $core.String? clipRemoved,
    TransferProgress? progress,
    Device? deviceChanged,
    ClipboardPermission? permissionChanged,
    Status? statusChanged,
  }) {
    final result = Event._();
    if (clipAdded != null) result.clipAdded = clipAdded;
    if (clipUpdated != null) result.clipUpdated = clipUpdated;
    if (clipRemoved != null) result.clipRemoved = clipRemoved;
    if (progress != null) result.progress = progress;
    if (deviceChanged != null) result.deviceChanged = deviceChanged;
    if (permissionChanged != null) result.permissionChanged = permissionChanged;
    if (statusChanged != null) result.statusChanged = statusChanged;
    return result;
  }

  Event._();

  factory Event.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Event()..mergeFromBuffer(data, registry);
  factory Event.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Event()..mergeFromJson(json, registry);

  static const $core.Map<$core.int, Event_Payload> _Event_PayloadByTag = {
    1: Event_Payload.clipAdded,
    2: Event_Payload.clipUpdated,
    3: Event_Payload.clipRemoved,
    4: Event_Payload.progress,
    5: Event_Payload.deviceChanged,
    6: Event_Payload.permissionChanged,
    7: Event_Payload.statusChanged,
    0: Event_Payload.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Event',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: Event.$_createMessage)
    ..oo(0, [1, 2, 3, 4, 5, 6, 7])
    ..aOM<ClipRecord>(1, _omitFieldNames ? '' : 'clipAdded',
        subBuilder: ClipRecord.$_createMessage)
    ..aOM<ClipRecord>(2, _omitFieldNames ? '' : 'clipUpdated',
        subBuilder: ClipRecord.$_createMessage)
    ..aOS(3, _omitFieldNames ? '' : 'clipRemoved')
    ..aOM<TransferProgress>(4, _omitFieldNames ? '' : 'progress',
        subBuilder: TransferProgress.$_createMessage)
    ..aOM<Device>(5, _omitFieldNames ? '' : 'deviceChanged',
        subBuilder: Device.$_createMessage)
    ..aE<ClipboardPermission>(6, _omitFieldNames ? '' : 'permissionChanged',
        enumValues: ClipboardPermission.values)
    ..aOM<Status>(7, _omitFieldNames ? '' : 'statusChanged',
        subBuilder: Status.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Event clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Event copyWith(void Function(Event) updates) =>
      super.copyWith((message) => updates(message as Event)) as Event;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Event() / Event.new instead')
  static Event create() => Event._();
  static $pb.GeneratedMessage $_createMessage() => Event._();
  @$core.override
  Event createEmptyInstance() => Event._();
  @$core.pragma('dart2js:noInline')
  static Event getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Event>(Event.$_createMessage);
  static Event? _defaultInstance;

  @$pb.TagNumber(1)
  @$pb.TagNumber(2)
  @$pb.TagNumber(3)
  @$pb.TagNumber(4)
  @$pb.TagNumber(5)
  @$pb.TagNumber(6)
  @$pb.TagNumber(7)
  Event_Payload whichPayload() => _Event_PayloadByTag[$_whichOneof(0)]!;
  @$pb.TagNumber(1)
  @$pb.TagNumber(2)
  @$pb.TagNumber(3)
  @$pb.TagNumber(4)
  @$pb.TagNumber(5)
  @$pb.TagNumber(6)
  @$pb.TagNumber(7)
  void clearPayload() => $_clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  ClipRecord get clipAdded => $_getN(0);
  @$pb.TagNumber(1)
  set clipAdded(ClipRecord value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasClipAdded() => $_has(0);
  @$pb.TagNumber(1)
  void clearClipAdded() => $_clearField(1);
  @$pb.TagNumber(1)
  ClipRecord ensureClipAdded() => $_ensure(0);

  @$pb.TagNumber(2)
  ClipRecord get clipUpdated => $_getN(1);
  @$pb.TagNumber(2)
  set clipUpdated(ClipRecord value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasClipUpdated() => $_has(1);
  @$pb.TagNumber(2)
  void clearClipUpdated() => $_clearField(2);
  @$pb.TagNumber(2)
  ClipRecord ensureClipUpdated() => $_ensure(1);

  @$pb.TagNumber(3)
  $core.String get clipRemoved => $_getSZ(2);
  @$pb.TagNumber(3)
  set clipRemoved($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasClipRemoved() => $_has(2);
  @$pb.TagNumber(3)
  void clearClipRemoved() => $_clearField(3);

  @$pb.TagNumber(4)
  TransferProgress get progress => $_getN(3);
  @$pb.TagNumber(4)
  set progress(TransferProgress value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasProgress() => $_has(3);
  @$pb.TagNumber(4)
  void clearProgress() => $_clearField(4);
  @$pb.TagNumber(4)
  TransferProgress ensureProgress() => $_ensure(3);

  @$pb.TagNumber(5)
  Device get deviceChanged => $_getN(4);
  @$pb.TagNumber(5)
  set deviceChanged(Device value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasDeviceChanged() => $_has(4);
  @$pb.TagNumber(5)
  void clearDeviceChanged() => $_clearField(5);
  @$pb.TagNumber(5)
  Device ensureDeviceChanged() => $_ensure(4);

  @$pb.TagNumber(6)
  ClipboardPermission get permissionChanged => $_getN(5);
  @$pb.TagNumber(6)
  set permissionChanged(ClipboardPermission value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasPermissionChanged() => $_has(5);
  @$pb.TagNumber(6)
  void clearPermissionChanged() => $_clearField(6);

  @$pb.TagNumber(7)
  Status get statusChanged => $_getN(6);
  @$pb.TagNumber(7)
  set statusChanged(Status value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasStatusChanged() => $_has(6);
  @$pb.TagNumber(7)
  void clearStatusChanged() => $_clearField(7);
  @$pb.TagNumber(7)
  Status ensureStatusChanged() => $_ensure(6);
}

class ListHistoryRequest extends $pb.GeneratedMessage {
  factory ListHistoryRequest({
    $core.int? limit,
    $fixnum.Int64? beforeUnix,
    ClipKind? kindFilter,
  }) {
    final result = ListHistoryRequest._();
    if (limit != null) result.limit = limit;
    if (beforeUnix != null) result.beforeUnix = beforeUnix;
    if (kindFilter != null) result.kindFilter = kindFilter;
    return result;
  }

  ListHistoryRequest._();

  factory ListHistoryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListHistoryRequest()..mergeFromBuffer(data, registry);
  factory ListHistoryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListHistoryRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHistoryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: ListHistoryRequest.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'limit')
    ..aInt64(2, _omitFieldNames ? '' : 'beforeUnix')
    ..aE<ClipKind>(3, _omitFieldNames ? '' : 'kindFilter',
        enumValues: ClipKind.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHistoryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHistoryRequest copyWith(void Function(ListHistoryRequest) updates) =>
      super.copyWith((message) => updates(message as ListHistoryRequest))
          as ListHistoryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ListHistoryRequest() / ListHistoryRequest.new instead')
  static ListHistoryRequest create() => ListHistoryRequest._();
  static $pb.GeneratedMessage $_createMessage() => ListHistoryRequest._();
  @$core.override
  ListHistoryRequest createEmptyInstance() => ListHistoryRequest._();
  @$core.pragma('dart2js:noInline')
  static ListHistoryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHistoryRequest>(
          ListHistoryRequest.$_createMessage);
  static ListHistoryRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(1)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(1)
  void clearLimit() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get beforeUnix => $_getI64(1);
  @$pb.TagNumber(2)
  set beforeUnix($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBeforeUnix() => $_has(1);
  @$pb.TagNumber(2)
  void clearBeforeUnix() => $_clearField(2);

  @$pb.TagNumber(3)
  ClipKind get kindFilter => $_getN(2);
  @$pb.TagNumber(3)
  set kindFilter(ClipKind value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasKindFilter() => $_has(2);
  @$pb.TagNumber(3)
  void clearKindFilter() => $_clearField(3);
}

class ListHistoryResponse extends $pb.GeneratedMessage {
  factory ListHistoryResponse({
    $core.Iterable<ClipRecord>? records,
  }) {
    final result = ListHistoryResponse._();
    if (records != null) result.records.addAll(records);
    return result;
  }

  ListHistoryResponse._();

  factory ListHistoryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListHistoryResponse()..mergeFromBuffer(data, registry);
  factory ListHistoryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListHistoryResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHistoryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: ListHistoryResponse.$_createMessage)
    ..pPM<ClipRecord>(1, _omitFieldNames ? '' : 'records',
        subBuilder: ClipRecord.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHistoryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHistoryResponse copyWith(void Function(ListHistoryResponse) updates) =>
      super.copyWith((message) => updates(message as ListHistoryResponse))
          as ListHistoryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use ListHistoryResponse() / ListHistoryResponse.new instead')
  static ListHistoryResponse create() => ListHistoryResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListHistoryResponse._();
  @$core.override
  ListHistoryResponse createEmptyInstance() => ListHistoryResponse._();
  @$core.pragma('dart2js:noInline')
  static ListHistoryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHistoryResponse>(
          ListHistoryResponse.$_createMessage);
  static ListHistoryResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<ClipRecord> get records => $_getList(0);
}

class DeleteHistoryRequest extends $pb.GeneratedMessage {
  factory DeleteHistoryRequest({
    $core.Iterable<$core.String>? ids,
    $core.bool? all,
  }) {
    final result = DeleteHistoryRequest._();
    if (ids != null) result.ids.addAll(ids);
    if (all != null) result.all = all;
    return result;
  }

  DeleteHistoryRequest._();

  factory DeleteHistoryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteHistoryRequest()..mergeFromBuffer(data, registry);
  factory DeleteHistoryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteHistoryRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteHistoryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: DeleteHistoryRequest.$_createMessage)
    ..pPS(1, _omitFieldNames ? '' : 'ids')
    ..aOB(2, _omitFieldNames ? '' : 'all')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHistoryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHistoryRequest copyWith(void Function(DeleteHistoryRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteHistoryRequest))
          as DeleteHistoryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DeleteHistoryRequest() / DeleteHistoryRequest.new instead')
  static DeleteHistoryRequest create() => DeleteHistoryRequest._();
  static $pb.GeneratedMessage $_createMessage() => DeleteHistoryRequest._();
  @$core.override
  DeleteHistoryRequest createEmptyInstance() => DeleteHistoryRequest._();
  @$core.pragma('dart2js:noInline')
  static DeleteHistoryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteHistoryRequest>(
          DeleteHistoryRequest.$_createMessage);
  static DeleteHistoryRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get ids => $_getList(0);

  @$pb.TagNumber(2)
  $core.bool get all => $_getBF(1);
  @$pb.TagNumber(2)
  set all($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAll() => $_has(1);
  @$pb.TagNumber(2)
  void clearAll() => $_clearField(2);
}

class FetchRequest extends $pb.GeneratedMessage {
  factory FetchRequest({
    $core.String? clipId,
    $core.bool? preserveClipboard,
  }) {
    final result = FetchRequest._();
    if (clipId != null) result.clipId = clipId;
    if (preserveClipboard != null) result.preserveClipboard = preserveClipboard;
    return result;
  }

  FetchRequest._();

  factory FetchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      FetchRequest()..mergeFromBuffer(data, registry);
  factory FetchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      FetchRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FetchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: FetchRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'clipId')
    ..aOB(2, _omitFieldNames ? '' : 'preserveClipboard')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FetchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FetchRequest copyWith(void Function(FetchRequest) updates) =>
      super.copyWith((message) => updates(message as FetchRequest))
          as FetchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use FetchRequest() / FetchRequest.new instead')
  static FetchRequest create() => FetchRequest._();
  static $pb.GeneratedMessage $_createMessage() => FetchRequest._();
  @$core.override
  FetchRequest createEmptyInstance() => FetchRequest._();
  @$core.pragma('dart2js:noInline')
  static FetchRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FetchRequest>(
          FetchRequest.$_createMessage);
  static FetchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get clipId => $_getSZ(0);
  @$pb.TagNumber(1)
  set clipId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasClipId() => $_has(0);
  @$pb.TagNumber(1)
  void clearClipId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get preserveClipboard => $_getBF(1);
  @$pb.TagNumber(2)
  set preserveClipboard($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPreserveClipboard() => $_has(1);
  @$pb.TagNumber(2)
  void clearPreserveClipboard() => $_clearField(2);
}

class ApplyToClipboardRequest extends $pb.GeneratedMessage {
  factory ApplyToClipboardRequest({
    $core.String? clipId,
  }) {
    final result = ApplyToClipboardRequest._();
    if (clipId != null) result.clipId = clipId;
    return result;
  }

  ApplyToClipboardRequest._();

  factory ApplyToClipboardRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ApplyToClipboardRequest()..mergeFromBuffer(data, registry);
  factory ApplyToClipboardRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ApplyToClipboardRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApplyToClipboardRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: ApplyToClipboardRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'clipId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplyToClipboardRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplyToClipboardRequest copyWith(
          void Function(ApplyToClipboardRequest) updates) =>
      super.copyWith((message) => updates(message as ApplyToClipboardRequest))
          as ApplyToClipboardRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ApplyToClipboardRequest() / ApplyToClipboardRequest.new instead')
  static ApplyToClipboardRequest create() => ApplyToClipboardRequest._();
  static $pb.GeneratedMessage $_createMessage() => ApplyToClipboardRequest._();
  @$core.override
  ApplyToClipboardRequest createEmptyInstance() => ApplyToClipboardRequest._();
  @$core.pragma('dart2js:noInline')
  static ApplyToClipboardRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ApplyToClipboardRequest>(
          ApplyToClipboardRequest.$_createMessage);
  static ApplyToClipboardRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get clipId => $_getSZ(0);
  @$pb.TagNumber(1)
  set clipId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasClipId() => $_has(0);
  @$pb.TagNumber(1)
  void clearClipId() => $_clearField(1);
}

class ListDevicesResponse extends $pb.GeneratedMessage {
  factory ListDevicesResponse({
    Device? self,
    $core.Iterable<Device>? peers,
  }) {
    final result = ListDevicesResponse._();
    if (self != null) result.self = self;
    if (peers != null) result.peers.addAll(peers);
    return result;
  }

  ListDevicesResponse._();

  factory ListDevicesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDevicesResponse()..mergeFromBuffer(data, registry);
  factory ListDevicesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDevicesResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDevicesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: ListDevicesResponse.$_createMessage)
    ..aOM<Device>(1, _omitFieldNames ? '' : 'self',
        subBuilder: Device.$_createMessage)
    ..pPM<Device>(2, _omitFieldNames ? '' : 'peers',
        subBuilder: Device.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDevicesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDevicesResponse copyWith(void Function(ListDevicesResponse) updates) =>
      super.copyWith((message) => updates(message as ListDevicesResponse))
          as ListDevicesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use ListDevicesResponse() / ListDevicesResponse.new instead')
  static ListDevicesResponse create() => ListDevicesResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListDevicesResponse._();
  @$core.override
  ListDevicesResponse createEmptyInstance() => ListDevicesResponse._();
  @$core.pragma('dart2js:noInline')
  static ListDevicesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDevicesResponse>(
          ListDevicesResponse.$_createMessage);
  static ListDevicesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Device get self => $_getN(0);
  @$pb.TagNumber(1)
  set self(Device value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasSelf() => $_has(0);
  @$pb.TagNumber(1)
  void clearSelf() => $_clearField(1);
  @$pb.TagNumber(1)
  Device ensureSelf() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Device> get peers => $_getList(1);
}

class CreatePairingCodeResponse extends $pb.GeneratedMessage {
  factory CreatePairingCodeResponse({
    $core.String? code,
    $fixnum.Int64? expiresAtUnix,
  }) {
    final result = CreatePairingCodeResponse._();
    if (code != null) result.code = code;
    if (expiresAtUnix != null) result.expiresAtUnix = expiresAtUnix;
    return result;
  }

  CreatePairingCodeResponse._();

  factory CreatePairingCodeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreatePairingCodeResponse()..mergeFromBuffer(data, registry);
  factory CreatePairingCodeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreatePairingCodeResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePairingCodeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: CreatePairingCodeResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'code')
    ..aInt64(2, _omitFieldNames ? '' : 'expiresAtUnix')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePairingCodeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePairingCodeResponse copyWith(
          void Function(CreatePairingCodeResponse) updates) =>
      super.copyWith((message) => updates(message as CreatePairingCodeResponse))
          as CreatePairingCodeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CreatePairingCodeResponse() / CreatePairingCodeResponse.new instead')
  static CreatePairingCodeResponse create() => CreatePairingCodeResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      CreatePairingCodeResponse._();
  @$core.override
  CreatePairingCodeResponse createEmptyInstance() =>
      CreatePairingCodeResponse._();
  @$core.pragma('dart2js:noInline')
  static CreatePairingCodeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePairingCodeResponse>(
          CreatePairingCodeResponse.$_createMessage);
  static CreatePairingCodeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get code => $_getSZ(0);
  @$pb.TagNumber(1)
  set code($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCode() => $_has(0);
  @$pb.TagNumber(1)
  void clearCode() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get expiresAtUnix => $_getI64(1);
  @$pb.TagNumber(2)
  set expiresAtUnix($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasExpiresAtUnix() => $_has(1);
  @$pb.TagNumber(2)
  void clearExpiresAtUnix() => $_clearField(2);
}

class RedeemPairingCodeRequest extends $pb.GeneratedMessage {
  factory RedeemPairingCodeRequest({
    $core.String? code,
  }) {
    final result = RedeemPairingCodeRequest._();
    if (code != null) result.code = code;
    return result;
  }

  RedeemPairingCodeRequest._();

  factory RedeemPairingCodeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RedeemPairingCodeRequest()..mergeFromBuffer(data, registry);
  factory RedeemPairingCodeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RedeemPairingCodeRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RedeemPairingCodeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: RedeemPairingCodeRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'code')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedeemPairingCodeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedeemPairingCodeRequest copyWith(
          void Function(RedeemPairingCodeRequest) updates) =>
      super.copyWith((message) => updates(message as RedeemPairingCodeRequest))
          as RedeemPairingCodeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RedeemPairingCodeRequest() / RedeemPairingCodeRequest.new instead')
  static RedeemPairingCodeRequest create() => RedeemPairingCodeRequest._();
  static $pb.GeneratedMessage $_createMessage() => RedeemPairingCodeRequest._();
  @$core.override
  RedeemPairingCodeRequest createEmptyInstance() =>
      RedeemPairingCodeRequest._();
  @$core.pragma('dart2js:noInline')
  static RedeemPairingCodeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RedeemPairingCodeRequest>(
          RedeemPairingCodeRequest.$_createMessage);
  static RedeemPairingCodeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get code => $_getSZ(0);
  @$pb.TagNumber(1)
  set code($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCode() => $_has(0);
  @$pb.TagNumber(1)
  void clearCode() => $_clearField(1);
}

class RedeemPairingCodeResponse extends $pb.GeneratedMessage {
  factory RedeemPairingCodeResponse({
    Device? peer,
    $core.String? pairingSession,
  }) {
    final result = RedeemPairingCodeResponse._();
    if (peer != null) result.peer = peer;
    if (pairingSession != null) result.pairingSession = pairingSession;
    return result;
  }

  RedeemPairingCodeResponse._();

  factory RedeemPairingCodeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RedeemPairingCodeResponse()..mergeFromBuffer(data, registry);
  factory RedeemPairingCodeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RedeemPairingCodeResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RedeemPairingCodeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: RedeemPairingCodeResponse.$_createMessage)
    ..aOM<Device>(1, _omitFieldNames ? '' : 'peer',
        subBuilder: Device.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'pairingSession')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedeemPairingCodeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedeemPairingCodeResponse copyWith(
          void Function(RedeemPairingCodeResponse) updates) =>
      super.copyWith((message) => updates(message as RedeemPairingCodeResponse))
          as RedeemPairingCodeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RedeemPairingCodeResponse() / RedeemPairingCodeResponse.new instead')
  static RedeemPairingCodeResponse create() => RedeemPairingCodeResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      RedeemPairingCodeResponse._();
  @$core.override
  RedeemPairingCodeResponse createEmptyInstance() =>
      RedeemPairingCodeResponse._();
  @$core.pragma('dart2js:noInline')
  static RedeemPairingCodeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RedeemPairingCodeResponse>(
          RedeemPairingCodeResponse.$_createMessage);
  static RedeemPairingCodeResponse? _defaultInstance;

  /// 返回对端信息供用户肉眼核对指纹，确认后再调 ConfirmPairing
  @$pb.TagNumber(1)
  Device get peer => $_getN(0);
  @$pb.TagNumber(1)
  set peer(Device value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPeer() => $_has(0);
  @$pb.TagNumber(1)
  void clearPeer() => $_clearField(1);
  @$pb.TagNumber(1)
  Device ensurePeer() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.String get pairingSession => $_getSZ(1);
  @$pb.TagNumber(2)
  set pairingSession($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPairingSession() => $_has(1);
  @$pb.TagNumber(2)
  void clearPairingSession() => $_clearField(2);
}

class ConfirmPairingRequest extends $pb.GeneratedMessage {
  factory ConfirmPairingRequest({
    $core.String? pairingSession,
    $core.bool? accept,
  }) {
    final result = ConfirmPairingRequest._();
    if (pairingSession != null) result.pairingSession = pairingSession;
    if (accept != null) result.accept = accept;
    return result;
  }

  ConfirmPairingRequest._();

  factory ConfirmPairingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ConfirmPairingRequest()..mergeFromBuffer(data, registry);
  factory ConfirmPairingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ConfirmPairingRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConfirmPairingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: ConfirmPairingRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'pairingSession')
    ..aOB(2, _omitFieldNames ? '' : 'accept')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfirmPairingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfirmPairingRequest copyWith(
          void Function(ConfirmPairingRequest) updates) =>
      super.copyWith((message) => updates(message as ConfirmPairingRequest))
          as ConfirmPairingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ConfirmPairingRequest() / ConfirmPairingRequest.new instead')
  static ConfirmPairingRequest create() => ConfirmPairingRequest._();
  static $pb.GeneratedMessage $_createMessage() => ConfirmPairingRequest._();
  @$core.override
  ConfirmPairingRequest createEmptyInstance() => ConfirmPairingRequest._();
  @$core.pragma('dart2js:noInline')
  static ConfirmPairingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConfirmPairingRequest>(
          ConfirmPairingRequest.$_createMessage);
  static ConfirmPairingRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get pairingSession => $_getSZ(0);
  @$pb.TagNumber(1)
  set pairingSession($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPairingSession() => $_has(0);
  @$pb.TagNumber(1)
  void clearPairingSession() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get accept => $_getBF(1);
  @$pb.TagNumber(2)
  set accept($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAccept() => $_has(1);
  @$pb.TagNumber(2)
  void clearAccept() => $_clearField(2);
}

class UnpairRequest extends $pb.GeneratedMessage {
  factory UnpairRequest({
    $core.String? deviceId,
  }) {
    final result = UnpairRequest._();
    if (deviceId != null) result.deviceId = deviceId;
    return result;
  }

  UnpairRequest._();

  factory UnpairRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UnpairRequest()..mergeFromBuffer(data, registry);
  factory UnpairRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UnpairRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnpairRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: UnpairRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'deviceId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnpairRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnpairRequest copyWith(void Function(UnpairRequest) updates) =>
      super.copyWith((message) => updates(message as UnpairRequest))
          as UnpairRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use UnpairRequest() / UnpairRequest.new instead')
  static UnpairRequest create() => UnpairRequest._();
  static $pb.GeneratedMessage $_createMessage() => UnpairRequest._();
  @$core.override
  UnpairRequest createEmptyInstance() => UnpairRequest._();
  @$core.pragma('dart2js:noInline')
  static UnpairRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnpairRequest>(
          UnpairRequest.$_createMessage);
  static UnpairRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deviceId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deviceId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeviceId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeviceId() => $_clearField(1);
}

class Config extends $pb.GeneratedMessage {
  factory Config({
    $fixnum.Int64? autoSyncThresholdBytes,
    $fixnum.Int64? historyTtlSeconds,
    $fixnum.Int64? cacheTtlSeconds,
    $core.bool? syncText,
    $core.bool? syncHtml,
    $core.bool? syncImage,
    $core.bool? syncFile,
    $core.bool? autoApplyToClipboard,
    $core.bool? launchAtLogin,
    $core.String? deviceName,
    $core.String? signalingUrl,
    $core.bool? onlyOwnStun,
  }) {
    final result = Config._();
    if (autoSyncThresholdBytes != null)
      result.autoSyncThresholdBytes = autoSyncThresholdBytes;
    if (historyTtlSeconds != null) result.historyTtlSeconds = historyTtlSeconds;
    if (cacheTtlSeconds != null) result.cacheTtlSeconds = cacheTtlSeconds;
    if (syncText != null) result.syncText = syncText;
    if (syncHtml != null) result.syncHtml = syncHtml;
    if (syncImage != null) result.syncImage = syncImage;
    if (syncFile != null) result.syncFile = syncFile;
    if (autoApplyToClipboard != null)
      result.autoApplyToClipboard = autoApplyToClipboard;
    if (launchAtLogin != null) result.launchAtLogin = launchAtLogin;
    if (deviceName != null) result.deviceName = deviceName;
    if (signalingUrl != null) result.signalingUrl = signalingUrl;
    if (onlyOwnStun != null) result.onlyOwnStun = onlyOwnStun;
    return result;
  }

  Config._();

  factory Config.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Config()..mergeFromBuffer(data, registry);
  factory Config.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Config()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Config',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: Config.$_createMessage)
    ..aInt64(1, _omitFieldNames ? '' : 'autoSyncThresholdBytes')
    ..aInt64(2, _omitFieldNames ? '' : 'historyTtlSeconds')
    ..aInt64(3, _omitFieldNames ? '' : 'cacheTtlSeconds')
    ..aOB(4, _omitFieldNames ? '' : 'syncText')
    ..aOB(5, _omitFieldNames ? '' : 'syncHtml')
    ..aOB(6, _omitFieldNames ? '' : 'syncImage')
    ..aOB(7, _omitFieldNames ? '' : 'syncFile')
    ..aOB(8, _omitFieldNames ? '' : 'autoApplyToClipboard')
    ..aOB(9, _omitFieldNames ? '' : 'launchAtLogin')
    ..aOS(10, _omitFieldNames ? '' : 'deviceName')
    ..aOS(11, _omitFieldNames ? '' : 'signalingUrl')
    ..aOB(12, _omitFieldNames ? '' : 'onlyOwnStun')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Config clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Config copyWith(void Function(Config) updates) =>
      super.copyWith((message) => updates(message as Config)) as Config;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Config() / Config.new instead')
  static Config create() => Config._();
  static $pb.GeneratedMessage $_createMessage() => Config._();
  @$core.override
  Config createEmptyInstance() => Config._();
  @$core.pragma('dart2js:noInline')
  static Config getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Config>(Config.$_createMessage);
  static Config? _defaultInstance;

  /// 小于此值自动预取，大于此值只同步元数据等用户手动拉取
  @$pb.TagNumber(1)
  $fixnum.Int64 get autoSyncThresholdBytes => $_getI64(0);
  @$pb.TagNumber(1)
  set autoSyncThresholdBytes($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAutoSyncThresholdBytes() => $_has(0);
  @$pb.TagNumber(1)
  void clearAutoSyncThresholdBytes() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get historyTtlSeconds => $_getI64(1);
  @$pb.TagNumber(2)
  set historyTtlSeconds($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasHistoryTtlSeconds() => $_has(1);
  @$pb.TagNumber(2)
  void clearHistoryTtlSeconds() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get cacheTtlSeconds => $_getI64(2);
  @$pb.TagNumber(3)
  set cacheTtlSeconds($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCacheTtlSeconds() => $_has(2);
  @$pb.TagNumber(3)
  void clearCacheTtlSeconds() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get syncText => $_getBF(3);
  @$pb.TagNumber(4)
  set syncText($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSyncText() => $_has(3);
  @$pb.TagNumber(4)
  void clearSyncText() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.bool get syncHtml => $_getBF(4);
  @$pb.TagNumber(5)
  set syncHtml($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSyncHtml() => $_has(4);
  @$pb.TagNumber(5)
  void clearSyncHtml() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.bool get syncImage => $_getBF(5);
  @$pb.TagNumber(6)
  set syncImage($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasSyncImage() => $_has(5);
  @$pb.TagNumber(6)
  void clearSyncImage() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.bool get syncFile => $_getBF(6);
  @$pb.TagNumber(7)
  set syncFile($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(7)
  $core.bool hasSyncFile() => $_has(6);
  @$pb.TagNumber(7)
  void clearSyncFile() => $_clearField(7);

  /// 收到对端内容后是否自动覆盖本机剪贴板
  @$pb.TagNumber(8)
  $core.bool get autoApplyToClipboard => $_getBF(7);
  @$pb.TagNumber(8)
  set autoApplyToClipboard($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(8)
  $core.bool hasAutoApplyToClipboard() => $_has(7);
  @$pb.TagNumber(8)
  void clearAutoApplyToClipboard() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.bool get launchAtLogin => $_getBF(8);
  @$pb.TagNumber(9)
  set launchAtLogin($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(9)
  $core.bool hasLaunchAtLogin() => $_has(8);
  @$pb.TagNumber(9)
  void clearLaunchAtLogin() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get deviceName => $_getSZ(9);
  @$pb.TagNumber(10)
  set deviceName($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasDeviceName() => $_has(9);
  @$pb.TagNumber(10)
  void clearDeviceName() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.String get signalingUrl => $_getSZ(10);
  @$pb.TagNumber(11)
  set signalingUrl($core.String value) => $_setString(10, value);
  @$pb.TagNumber(11)
  $core.bool hasSignalingUrl() => $_has(10);
  @$pb.TagNumber(11)
  void clearSignalingUrl() => $_clearField(11);

  /// 只用自己的服务器探测网络出口，不用公共 STUN 服务器。
  /// 默认 false（即允许公共服务器）：旧版界面不发这个字段时，不会被误关掉。
  @$pb.TagNumber(12)
  $core.bool get onlyOwnStun => $_getBF(11);
  @$pb.TagNumber(12)
  set onlyOwnStun($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(12)
  $core.bool hasOnlyOwnStun() => $_has(11);
  @$pb.TagNumber(12)
  void clearOnlyOwnStun() => $_clearField(12);
}

/// 最近一轮网络出口探测的结果。
class NetworkInfo extends $pb.GeneratedMessage {
  factory NetworkInfo({
    $core.int? localPort,
    $core.Iterable<NetworkEgress>? egresses,
    $core.int? probed,
    $core.int? answered,
    $fixnum.Int64? probedAtUnix,
  }) {
    final result = NetworkInfo._();
    if (localPort != null) result.localPort = localPort;
    if (egresses != null) result.egresses.addAll(egresses);
    if (probed != null) result.probed = probed;
    if (answered != null) result.answered = answered;
    if (probedAtUnix != null) result.probedAtUnix = probedAtUnix;
    return result;
  }

  NetworkInfo._();

  factory NetworkInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NetworkInfo()..mergeFromBuffer(data, registry);
  factory NetworkInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NetworkInfo()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NetworkInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: NetworkInfo.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'localPort')
    ..pPM<NetworkEgress>(2, _omitFieldNames ? '' : 'egresses',
        subBuilder: NetworkEgress.$_createMessage)
    ..aI(3, _omitFieldNames ? '' : 'probed')
    ..aI(4, _omitFieldNames ? '' : 'answered')
    ..aInt64(5, _omitFieldNames ? '' : 'probedAtUnix')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkInfo copyWith(void Function(NetworkInfo) updates) =>
      super.copyWith((message) => updates(message as NetworkInfo))
          as NetworkInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use NetworkInfo() / NetworkInfo.new instead')
  static NetworkInfo create() => NetworkInfo._();
  static $pb.GeneratedMessage $_createMessage() => NetworkInfo._();
  @$core.override
  NetworkInfo createEmptyInstance() => NetworkInfo._();
  @$core.pragma('dart2js:noInline')
  static NetworkInfo getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<NetworkInfo>(
          NetworkInfo.$_createMessage);
  static NetworkInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get localPort => $_getIZ(0);
  @$pb.TagNumber(1)
  set localPort($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLocalPort() => $_has(0);
  @$pb.TagNumber(1)
  void clearLocalPort() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<NetworkEgress> get egresses => $_getList(1);

  @$pb.TagNumber(3)
  $core.int get probed => $_getIZ(2);
  @$pb.TagNumber(3)
  set probed($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasProbed() => $_has(2);
  @$pb.TagNumber(3)
  void clearProbed() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.int get answered => $_getIZ(3);
  @$pb.TagNumber(4)
  set answered($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAnswered() => $_has(3);
  @$pb.TagNumber(4)
  void clearAnswered() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get probedAtUnix => $_getI64(4);
  @$pb.TagNumber(5)
  set probedAtUnix($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasProbedAtUnix() => $_has(4);
  @$pb.TagNumber(5)
  void clearProbedAtUnix() => $_clearField(5);
}

class NetworkEgress extends $pb.GeneratedMessage {
  factory NetworkEgress({
    $core.String? address,
    $core.bool? portPreserved,
    $core.bool? guessed,
    $core.Iterable<$core.String>? via,
  }) {
    final result = NetworkEgress._();
    if (address != null) result.address = address;
    if (portPreserved != null) result.portPreserved = portPreserved;
    if (guessed != null) result.guessed = guessed;
    if (via != null) result.via.addAll(via);
    return result;
  }

  NetworkEgress._();

  factory NetworkEgress.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NetworkEgress()..mergeFromBuffer(data, registry);
  factory NetworkEgress.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NetworkEgress()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NetworkEgress',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: NetworkEgress.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOB(2, _omitFieldNames ? '' : 'portPreserved')
    ..aOB(3, _omitFieldNames ? '' : 'guessed')
    ..pPS(4, _omitFieldNames ? '' : 'via')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkEgress clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkEgress copyWith(void Function(NetworkEgress) updates) =>
      super.copyWith((message) => updates(message as NetworkEgress))
          as NetworkEgress;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use NetworkEgress() / NetworkEgress.new instead')
  static NetworkEgress create() => NetworkEgress._();
  static $pb.GeneratedMessage $_createMessage() => NetworkEgress._();
  @$core.override
  NetworkEgress createEmptyInstance() => NetworkEgress._();
  @$core.pragma('dart2js:noInline')
  static NetworkEgress getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<NetworkEgress>(
          NetworkEgress.$_createMessage);
  static NetworkEgress? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get portPreserved => $_getBF(1);
  @$pb.TagNumber(2)
  set portPreserved($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPortPreserved() => $_has(1);
  @$pb.TagNumber(2)
  void clearPortPreserved() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get guessed => $_getBF(2);
  @$pb.TagNumber(3)
  set guessed($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasGuessed() => $_has(2);
  @$pb.TagNumber(3)
  void clearGuessed() => $_clearField(3);

  @$pb.TagNumber(4)
  $pb.PbList<$core.String> get via => $_getList(3);
}

class Status extends $pb.GeneratedMessage {
  factory Status({
    $core.String? deviceId,
    $core.String? deviceName,
    $core.bool? signalingConnected,
    ClipboardPermission? clipboardPermission,
    $fixnum.Int64? cacheBytesUsed,
    $core.int? peersOnline,
    $core.String? version,
  }) {
    final result = Status._();
    if (deviceId != null) result.deviceId = deviceId;
    if (deviceName != null) result.deviceName = deviceName;
    if (signalingConnected != null)
      result.signalingConnected = signalingConnected;
    if (clipboardPermission != null)
      result.clipboardPermission = clipboardPermission;
    if (cacheBytesUsed != null) result.cacheBytesUsed = cacheBytesUsed;
    if (peersOnline != null) result.peersOnline = peersOnline;
    if (version != null) result.version = version;
    return result;
  }

  Status._();

  factory Status.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Status()..mergeFromBuffer(data, registry);
  factory Status.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Status()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Status',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: Status.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'deviceId')
    ..aOS(2, _omitFieldNames ? '' : 'deviceName')
    ..aOB(3, _omitFieldNames ? '' : 'signalingConnected')
    ..aE<ClipboardPermission>(4, _omitFieldNames ? '' : 'clipboardPermission',
        enumValues: ClipboardPermission.values)
    ..aInt64(5, _omitFieldNames ? '' : 'cacheBytesUsed')
    ..aI(6, _omitFieldNames ? '' : 'peersOnline')
    ..aOS(7, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Status clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Status copyWith(void Function(Status) updates) =>
      super.copyWith((message) => updates(message as Status)) as Status;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Status() / Status.new instead')
  static Status create() => Status._();
  static $pb.GeneratedMessage $_createMessage() => Status._();
  @$core.override
  Status createEmptyInstance() => Status._();
  @$core.pragma('dart2js:noInline')
  static Status getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Status>(Status.$_createMessage);
  static Status? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deviceId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deviceId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeviceId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeviceId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get deviceName => $_getSZ(1);
  @$pb.TagNumber(2)
  set deviceName($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDeviceName() => $_has(1);
  @$pb.TagNumber(2)
  void clearDeviceName() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get signalingConnected => $_getBF(2);
  @$pb.TagNumber(3)
  set signalingConnected($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSignalingConnected() => $_has(2);
  @$pb.TagNumber(3)
  void clearSignalingConnected() => $_clearField(3);

  @$pb.TagNumber(4)
  ClipboardPermission get clipboardPermission => $_getN(3);
  @$pb.TagNumber(4)
  set clipboardPermission(ClipboardPermission value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasClipboardPermission() => $_has(3);
  @$pb.TagNumber(4)
  void clearClipboardPermission() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get cacheBytesUsed => $_getI64(4);
  @$pb.TagNumber(5)
  set cacheBytesUsed($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasCacheBytesUsed() => $_has(4);
  @$pb.TagNumber(5)
  void clearCacheBytesUsed() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get peersOnline => $_getIZ(5);
  @$pb.TagNumber(6)
  set peersOnline($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasPeersOnline() => $_has(5);
  @$pb.TagNumber(6)
  void clearPeersOnline() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get version => $_getSZ(6);
  @$pb.TagNumber(7)
  set version($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasVersion() => $_has(6);
  @$pb.TagNumber(7)
  void clearVersion() => $_clearField(7);
}

class GetImagePreviewRequest extends $pb.GeneratedMessage {
  factory GetImagePreviewRequest({
    $core.String? clipId,
  }) {
    final result = GetImagePreviewRequest._();
    if (clipId != null) result.clipId = clipId;
    return result;
  }

  GetImagePreviewRequest._();

  factory GetImagePreviewRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetImagePreviewRequest()..mergeFromBuffer(data, registry);
  factory GetImagePreviewRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetImagePreviewRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetImagePreviewRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: GetImagePreviewRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'clipId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImagePreviewRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImagePreviewRequest copyWith(
          void Function(GetImagePreviewRequest) updates) =>
      super.copyWith((message) => updates(message as GetImagePreviewRequest))
          as GetImagePreviewRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetImagePreviewRequest() / GetImagePreviewRequest.new instead')
  static GetImagePreviewRequest create() => GetImagePreviewRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetImagePreviewRequest._();
  @$core.override
  GetImagePreviewRequest createEmptyInstance() => GetImagePreviewRequest._();
  @$core.pragma('dart2js:noInline')
  static GetImagePreviewRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetImagePreviewRequest>(
          GetImagePreviewRequest.$_createMessage);
  static GetImagePreviewRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get clipId => $_getSZ(0);
  @$pb.TagNumber(1)
  set clipId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasClipId() => $_has(0);
  @$pb.TagNumber(1)
  void clearClipId() => $_clearField(1);
}

class GetImagePreviewResponse extends $pb.GeneratedMessage {
  factory GetImagePreviewResponse({
    $core.String? path,
  }) {
    final result = GetImagePreviewResponse._();
    if (path != null) result.path = path;
    return result;
  }

  GetImagePreviewResponse._();

  factory GetImagePreviewResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetImagePreviewResponse()..mergeFromBuffer(data, registry);
  factory GetImagePreviewResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetImagePreviewResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetImagePreviewResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'copysync.v1'),
      createEmptyInstance: GetImagePreviewResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'path')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImagePreviewResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImagePreviewResponse copyWith(
          void Function(GetImagePreviewResponse) updates) =>
      super.copyWith((message) => updates(message as GetImagePreviewResponse))
          as GetImagePreviewResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetImagePreviewResponse() / GetImagePreviewResponse.new instead')
  static GetImagePreviewResponse create() => GetImagePreviewResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetImagePreviewResponse._();
  @$core.override
  GetImagePreviewResponse createEmptyInstance() => GetImagePreviewResponse._();
  @$core.pragma('dart2js:noInline')
  static GetImagePreviewResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetImagePreviewResponse>(
          GetImagePreviewResponse.$_createMessage);
  static GetImagePreviewResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get path => $_getSZ(0);
  @$pb.TagNumber(1)
  set path($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPath() => $_has(0);
  @$pb.TagNumber(1)
  void clearPath() => $_clearField(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');

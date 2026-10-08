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

import 'package:protobuf/protobuf.dart' as $pb;

class ClipKind extends $pb.ProtobufEnum {
  static const ClipKind CLIP_KIND_UNSPECIFIED =
      ClipKind._(0, _omitEnumNames ? '' : 'CLIP_KIND_UNSPECIFIED');
  static const ClipKind CLIP_KIND_TEXT =
      ClipKind._(1, _omitEnumNames ? '' : 'CLIP_KIND_TEXT');
  static const ClipKind CLIP_KIND_HTML =
      ClipKind._(2, _omitEnumNames ? '' : 'CLIP_KIND_HTML');
  static const ClipKind CLIP_KIND_IMAGE =
      ClipKind._(3, _omitEnumNames ? '' : 'CLIP_KIND_IMAGE');
  static const ClipKind CLIP_KIND_FILE =
      ClipKind._(4, _omitEnumNames ? '' : 'CLIP_KIND_FILE');

  static const $core.List<ClipKind> values = <ClipKind>[
    CLIP_KIND_UNSPECIFIED,
    CLIP_KIND_TEXT,
    CLIP_KIND_HTML,
    CLIP_KIND_IMAGE,
    CLIP_KIND_FILE,
  ];

  static final $core.List<ClipKind?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ClipKind? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ClipKind._(super.value, super.name);
}

class ClipStatus extends $pb.ProtobufEnum {
  static const ClipStatus CLIP_STATUS_UNSPECIFIED =
      ClipStatus._(0, _omitEnumNames ? '' : 'CLIP_STATUS_UNSPECIFIED');
  static const ClipStatus CLIP_STATUS_READY =
      ClipStatus._(1, _omitEnumNames ? '' : 'CLIP_STATUS_READY');
  static const ClipStatus CLIP_STATUS_REMOTE_ONLY =
      ClipStatus._(2, _omitEnumNames ? '' : 'CLIP_STATUS_REMOTE_ONLY');
  static const ClipStatus CLIP_STATUS_FETCHING =
      ClipStatus._(3, _omitEnumNames ? '' : 'CLIP_STATUS_FETCHING');
  static const ClipStatus CLIP_STATUS_FAILED =
      ClipStatus._(4, _omitEnumNames ? '' : 'CLIP_STATUS_FAILED');
  static const ClipStatus CLIP_STATUS_EXPIRED =
      ClipStatus._(5, _omitEnumNames ? '' : 'CLIP_STATUS_EXPIRED');

  static const $core.List<ClipStatus> values = <ClipStatus>[
    CLIP_STATUS_UNSPECIFIED,
    CLIP_STATUS_READY,
    CLIP_STATUS_REMOTE_ONLY,
    CLIP_STATUS_FETCHING,
    CLIP_STATUS_FAILED,
    CLIP_STATUS_EXPIRED,
  ];

  static final $core.List<ClipStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ClipStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ClipStatus._(super.value, super.name);
}

class ConnectionKind extends $pb.ProtobufEnum {
  static const ConnectionKind CONNECTION_KIND_UNSPECIFIED =
      ConnectionKind._(0, _omitEnumNames ? '' : 'CONNECTION_KIND_UNSPECIFIED');
  static const ConnectionKind CONNECTION_KIND_OFFLINE =
      ConnectionKind._(1, _omitEnumNames ? '' : 'CONNECTION_KIND_OFFLINE');
  static const ConnectionKind CONNECTION_KIND_DIRECT =
      ConnectionKind._(2, _omitEnumNames ? '' : 'CONNECTION_KIND_DIRECT');
  static const ConnectionKind CONNECTION_KIND_RELAY =
      ConnectionKind._(3, _omitEnumNames ? '' : 'CONNECTION_KIND_RELAY');

  static const $core.List<ConnectionKind> values = <ConnectionKind>[
    CONNECTION_KIND_UNSPECIFIED,
    CONNECTION_KIND_OFFLINE,
    CONNECTION_KIND_DIRECT,
    CONNECTION_KIND_RELAY,
  ];

  static final $core.List<ConnectionKind?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ConnectionKind? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionKind._(super.value, super.name);
}

/// macOS 剪贴板授权状态，对应 NSPasteboardAccessBehavior。
/// 非 ALWAYS_ALLOW 时 daemon 不会去读剪贴板内容，避免卡在授权弹窗上。
class ClipboardPermission extends $pb.ProtobufEnum {
  static const ClipboardPermission CLIPBOARD_PERMISSION_UNSPECIFIED =
      ClipboardPermission._(
          0, _omitEnumNames ? '' : 'CLIPBOARD_PERMISSION_UNSPECIFIED');
  static const ClipboardPermission CLIPBOARD_PERMISSION_NOT_APPLICABLE =
      ClipboardPermission._(
          1, _omitEnumNames ? '' : 'CLIPBOARD_PERMISSION_NOT_APPLICABLE');
  static const ClipboardPermission CLIPBOARD_PERMISSION_DEFAULT =
      ClipboardPermission._(
          2, _omitEnumNames ? '' : 'CLIPBOARD_PERMISSION_DEFAULT');
  static const ClipboardPermission CLIPBOARD_PERMISSION_ASK =
      ClipboardPermission._(
          3, _omitEnumNames ? '' : 'CLIPBOARD_PERMISSION_ASK');
  static const ClipboardPermission CLIPBOARD_PERMISSION_ALWAYS_ALLOW =
      ClipboardPermission._(
          4, _omitEnumNames ? '' : 'CLIPBOARD_PERMISSION_ALWAYS_ALLOW');
  static const ClipboardPermission CLIPBOARD_PERMISSION_ALWAYS_DENY =
      ClipboardPermission._(
          5, _omitEnumNames ? '' : 'CLIPBOARD_PERMISSION_ALWAYS_DENY');

  static const $core.List<ClipboardPermission> values = <ClipboardPermission>[
    CLIPBOARD_PERMISSION_UNSPECIFIED,
    CLIPBOARD_PERMISSION_NOT_APPLICABLE,
    CLIPBOARD_PERMISSION_DEFAULT,
    CLIPBOARD_PERMISSION_ASK,
    CLIPBOARD_PERMISSION_ALWAYS_ALLOW,
    CLIPBOARD_PERMISSION_ALWAYS_DENY,
  ];

  static final $core.List<ClipboardPermission?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ClipboardPermission? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ClipboardPermission._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');

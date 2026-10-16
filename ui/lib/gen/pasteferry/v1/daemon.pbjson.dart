// This is a generated file - do not edit.
//
// Generated from pasteferry/v1/daemon.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use clipKindDescriptor instead')
const ClipKind$json = {
  '1': 'ClipKind',
  '2': [
    {'1': 'CLIP_KIND_UNSPECIFIED', '2': 0},
    {'1': 'CLIP_KIND_TEXT', '2': 1},
    {'1': 'CLIP_KIND_HTML', '2': 2},
    {'1': 'CLIP_KIND_IMAGE', '2': 3},
    {'1': 'CLIP_KIND_FILE', '2': 4},
  ],
};

/// Descriptor for `ClipKind`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List clipKindDescriptor = $convert.base64Decode(
    'CghDbGlwS2luZBIZChVDTElQX0tJTkRfVU5TUEVDSUZJRUQQABISCg5DTElQX0tJTkRfVEVYVB'
    'ABEhIKDkNMSVBfS0lORF9IVE1MEAISEwoPQ0xJUF9LSU5EX0lNQUdFEAMSEgoOQ0xJUF9LSU5E'
    'X0ZJTEUQBA==');

@$core.Deprecated('Use clipStatusDescriptor instead')
const ClipStatus$json = {
  '1': 'ClipStatus',
  '2': [
    {'1': 'CLIP_STATUS_UNSPECIFIED', '2': 0},
    {'1': 'CLIP_STATUS_READY', '2': 1},
    {'1': 'CLIP_STATUS_REMOTE_ONLY', '2': 2},
    {'1': 'CLIP_STATUS_FETCHING', '2': 3},
    {'1': 'CLIP_STATUS_FAILED', '2': 4},
    {'1': 'CLIP_STATUS_EXPIRED', '2': 5},
  ],
};

/// Descriptor for `ClipStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List clipStatusDescriptor = $convert.base64Decode(
    'CgpDbGlwU3RhdHVzEhsKF0NMSVBfU1RBVFVTX1VOU1BFQ0lGSUVEEAASFQoRQ0xJUF9TVEFUVV'
    'NfUkVBRFkQARIbChdDTElQX1NUQVRVU19SRU1PVEVfT05MWRACEhgKFENMSVBfU1RBVFVTX0ZF'
    'VENISU5HEAMSFgoSQ0xJUF9TVEFUVVNfRkFJTEVEEAQSFwoTQ0xJUF9TVEFUVVNfRVhQSVJFRB'
    'AF');

@$core.Deprecated('Use connectionKindDescriptor instead')
const ConnectionKind$json = {
  '1': 'ConnectionKind',
  '2': [
    {'1': 'CONNECTION_KIND_UNSPECIFIED', '2': 0},
    {'1': 'CONNECTION_KIND_OFFLINE', '2': 1},
    {'1': 'CONNECTION_KIND_DIRECT', '2': 2},
    {'1': 'CONNECTION_KIND_RELAY', '2': 3},
  ],
};

/// Descriptor for `ConnectionKind`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionKindDescriptor = $convert.base64Decode(
    'Cg5Db25uZWN0aW9uS2luZBIfChtDT05ORUNUSU9OX0tJTkRfVU5TUEVDSUZJRUQQABIbChdDT0'
    '5ORUNUSU9OX0tJTkRfT0ZGTElORRABEhoKFkNPTk5FQ1RJT05fS0lORF9ESVJFQ1QQAhIZChVD'
    'T05ORUNUSU9OX0tJTkRfUkVMQVkQAw==');

@$core.Deprecated('Use clipboardPermissionDescriptor instead')
const ClipboardPermission$json = {
  '1': 'ClipboardPermission',
  '2': [
    {'1': 'CLIPBOARD_PERMISSION_UNSPECIFIED', '2': 0},
    {'1': 'CLIPBOARD_PERMISSION_NOT_APPLICABLE', '2': 1},
    {'1': 'CLIPBOARD_PERMISSION_DEFAULT', '2': 2},
    {'1': 'CLIPBOARD_PERMISSION_ASK', '2': 3},
    {'1': 'CLIPBOARD_PERMISSION_ALWAYS_ALLOW', '2': 4},
    {'1': 'CLIPBOARD_PERMISSION_ALWAYS_DENY', '2': 5},
  ],
};

/// Descriptor for `ClipboardPermission`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List clipboardPermissionDescriptor = $convert.base64Decode(
    'ChNDbGlwYm9hcmRQZXJtaXNzaW9uEiQKIENMSVBCT0FSRF9QRVJNSVNTSU9OX1VOU1BFQ0lGSU'
    'VEEAASJwojQ0xJUEJPQVJEX1BFUk1JU1NJT05fTk9UX0FQUExJQ0FCTEUQARIgChxDTElQQk9B'
    'UkRfUEVSTUlTU0lPTl9ERUZBVUxUEAISHAoYQ0xJUEJPQVJEX1BFUk1JU1NJT05fQVNLEAMSJQ'
    'ohQ0xJUEJPQVJEX1BFUk1JU1NJT05fQUxXQVlTX0FMTE9XEAQSJAogQ0xJUEJPQVJEX1BFUk1J'
    'U1NJT05fQUxXQVlTX0RFTlkQBQ==');

@$core.Deprecated('Use clipboardWriteOutcomeDescriptor instead')
const ClipboardWriteOutcome$json = {
  '1': 'ClipboardWriteOutcome',
  '2': [
    {'1': 'CLIPBOARD_WRITE_OUTCOME_UNSPECIFIED', '2': 0},
    {'1': 'CLIPBOARD_WRITE_OUTCOME_APPLIED', '2': 1},
    {'1': 'CLIPBOARD_WRITE_OUTCOME_PRESERVED', '2': 2},
    {'1': 'CLIPBOARD_WRITE_OUTCOME_FAILED', '2': 3},
  ],
};

/// Descriptor for `ClipboardWriteOutcome`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List clipboardWriteOutcomeDescriptor = $convert.base64Decode(
    'ChVDbGlwYm9hcmRXcml0ZU91dGNvbWUSJwojQ0xJUEJPQVJEX1dSSVRFX09VVENPTUVfVU5TUE'
    'VDSUZJRUQQABIjCh9DTElQQk9BUkRfV1JJVEVfT1VUQ09NRV9BUFBMSUVEEAESJQohQ0xJUEJP'
    'QVJEX1dSSVRFX09VVENPTUVfUFJFU0VSVkVEEAISIgoeQ0xJUEJPQVJEX1dSSVRFX09VVENPTU'
    'VfRkFJTEVEEAM=');

@$core.Deprecated('Use serviceModeDescriptor instead')
const ServiceMode$json = {
  '1': 'ServiceMode',
  '2': [
    {'1': 'SERVICE_MODE_UNSPECIFIED', '2': 0},
    {'1': 'SERVICE_MODE_LAN', '2': 1},
    {'1': 'SERVICE_MODE_OFFICIAL', '2': 2},
    {'1': 'SERVICE_MODE_CUSTOM', '2': 3},
  ],
};

/// Descriptor for `ServiceMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List serviceModeDescriptor = $convert.base64Decode(
    'CgtTZXJ2aWNlTW9kZRIcChhTRVJWSUNFX01PREVfVU5TUEVDSUZJRUQQABIUChBTRVJWSUNFX0'
    '1PREVfTEFOEAESGQoVU0VSVklDRV9NT0RFX09GRklDSUFMEAISFwoTU0VSVklDRV9NT0RFX0NV'
    'U1RPTRAD');

@$core.Deprecated('Use emptyDescriptor instead')
const Empty$json = {
  '1': 'Empty',
};

/// Descriptor for `Empty`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emptyDescriptor =
    $convert.base64Decode('CgVFbXB0eQ==');

@$core.Deprecated('Use clipItemDescriptor instead')
const ClipItem$json = {
  '1': 'ClipItem',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'size', '3': 2, '4': 1, '5': 3, '10': 'size'},
    {'1': 'is_dir', '3': 3, '4': 1, '5': 8, '10': 'isDir'},
    {'1': 'content_type', '3': 4, '4': 1, '5': 9, '10': 'contentType'},
  ],
};

/// Descriptor for `ClipItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clipItemDescriptor = $convert.base64Decode(
    'CghDbGlwSXRlbRISCgRuYW1lGAEgASgJUgRuYW1lEhIKBHNpemUYAiABKANSBHNpemUSFQoGaX'
    'NfZGlyGAMgASgIUgVpc0RpchIhCgxjb250ZW50X3R5cGUYBCABKAlSC2NvbnRlbnRUeXBl');

@$core.Deprecated('Use clipRecordDescriptor instead')
const ClipRecord$json = {
  '1': 'ClipRecord',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'kind',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ClipKind',
      '10': 'kind'
    },
    {
      '1': 'status',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ClipStatus',
      '10': 'status'
    },
    {'1': 'origin_device_id', '3': 4, '4': 1, '5': 9, '10': 'originDeviceId'},
    {
      '1': 'origin_device_name',
      '3': 5,
      '4': 1,
      '5': 9,
      '10': 'originDeviceName'
    },
    {'1': 'outgoing', '3': 6, '4': 1, '5': 8, '10': 'outgoing'},
    {
      '1': 'items',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.pasteferry.v1.ClipItem',
      '10': 'items'
    },
    {'1': 'total_size', '3': 8, '4': 1, '5': 3, '10': 'totalSize'},
    {'1': 'text_preview', '3': 9, '4': 1, '5': 9, '10': 'textPreview'},
    {'1': 'created_at_unix', '3': 10, '4': 1, '5': 3, '10': 'createdAtUnix'},
    {'1': 'expires_at_unix', '3': 11, '4': 1, '5': 3, '10': 'expiresAtUnix'},
    {'1': 'error', '3': 12, '4': 1, '5': 9, '10': 'error'},
    {'1': 'image_count', '3': 13, '4': 1, '5': 5, '10': 'imageCount'},
  ],
};

/// Descriptor for `ClipRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clipRecordDescriptor = $convert.base64Decode(
    'CgpDbGlwUmVjb3JkEg4KAmlkGAEgASgJUgJpZBIrCgRraW5kGAIgASgOMhcucGFzdGVmZXJyeS'
    '52MS5DbGlwS2luZFIEa2luZBIxCgZzdGF0dXMYAyABKA4yGS5wYXN0ZWZlcnJ5LnYxLkNsaXBT'
    'dGF0dXNSBnN0YXR1cxIoChBvcmlnaW5fZGV2aWNlX2lkGAQgASgJUg5vcmlnaW5EZXZpY2VJZB'
    'IsChJvcmlnaW5fZGV2aWNlX25hbWUYBSABKAlSEG9yaWdpbkRldmljZU5hbWUSGgoIb3V0Z29p'
    'bmcYBiABKAhSCG91dGdvaW5nEi0KBWl0ZW1zGAcgAygLMhcucGFzdGVmZXJyeS52MS5DbGlwSX'
    'RlbVIFaXRlbXMSHQoKdG90YWxfc2l6ZRgIIAEoA1IJdG90YWxTaXplEiEKDHRleHRfcHJldmll'
    'dxgJIAEoCVILdGV4dFByZXZpZXcSJgoPY3JlYXRlZF9hdF91bml4GAogASgDUg1jcmVhdGVkQX'
    'RVbml4EiYKD2V4cGlyZXNfYXRfdW5peBgLIAEoA1INZXhwaXJlc0F0VW5peBIUCgVlcnJvchgM'
    'IAEoCVIFZXJyb3ISHwoLaW1hZ2VfY291bnQYDSABKAVSCmltYWdlQ291bnQ=');

@$core.Deprecated('Use subscribeRequestDescriptor instead')
const SubscribeRequest$json = {
  '1': 'SubscribeRequest',
  '2': [
    {'1': 'since_unix', '3': 1, '4': 1, '5': 3, '10': 'sinceUnix'},
  ],
};

/// Descriptor for `SubscribeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeRequestDescriptor = $convert.base64Decode(
    'ChBTdWJzY3JpYmVSZXF1ZXN0Eh0KCnNpbmNlX3VuaXgYASABKANSCXNpbmNlVW5peA==');

@$core.Deprecated('Use transferProgressDescriptor instead')
const TransferProgress$json = {
  '1': 'TransferProgress',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {'1': 'transferred', '3': 2, '4': 1, '5': 3, '10': 'transferred'},
    {'1': 'total', '3': 3, '4': 1, '5': 3, '10': 'total'},
    {'1': 'bytes_per_second', '3': 4, '4': 1, '5': 1, '10': 'bytesPerSecond'},
  ],
};

/// Descriptor for `TransferProgress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transferProgressDescriptor = $convert.base64Decode(
    'ChBUcmFuc2ZlclByb2dyZXNzEhcKB2NsaXBfaWQYASABKAlSBmNsaXBJZBIgCgt0cmFuc2Zlcn'
    'JlZBgCIAEoA1ILdHJhbnNmZXJyZWQSFAoFdG90YWwYAyABKANSBXRvdGFsEigKEGJ5dGVzX3Bl'
    'cl9zZWNvbmQYBCABKAFSDmJ5dGVzUGVyU2Vjb25k');

@$core.Deprecated('Use deviceDescriptor instead')
const Device$json = {
  '1': 'Device',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'platform', '3': 3, '4': 1, '5': 9, '10': 'platform'},
    {'1': 'online', '3': 4, '4': 1, '5': 8, '10': 'online'},
    {
      '1': 'connection',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ConnectionKind',
      '10': 'connection'
    },
    {
      '1': 'public_key_fingerprint',
      '3': 6,
      '4': 1,
      '5': 9,
      '10': 'publicKeyFingerprint'
    },
    {'1': 'paired_at_unix', '3': 7, '4': 1, '5': 3, '10': 'pairedAtUnix'},
    {'1': 'pairing_session', '3': 8, '4': 1, '5': 9, '10': 'pairingSession'},
  ],
};

/// Descriptor for `Device`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceDescriptor = $convert.base64Decode(
    'CgZEZXZpY2USDgoCaWQYASABKAlSAmlkEhIKBG5hbWUYAiABKAlSBG5hbWUSGgoIcGxhdGZvcm'
    '0YAyABKAlSCHBsYXRmb3JtEhYKBm9ubGluZRgEIAEoCFIGb25saW5lEj0KCmNvbm5lY3Rpb24Y'
    'BSABKA4yHS5wYXN0ZWZlcnJ5LnYxLkNvbm5lY3Rpb25LaW5kUgpjb25uZWN0aW9uEjQKFnB1Ym'
    'xpY19rZXlfZmluZ2VycHJpbnQYBiABKAlSFHB1YmxpY0tleUZpbmdlcnByaW50EiQKDnBhaXJl'
    'ZF9hdF91bml4GAcgASgDUgxwYWlyZWRBdFVuaXgSJwoPcGFpcmluZ19zZXNzaW9uGAggASgJUg'
    '5wYWlyaW5nU2Vzc2lvbg==');

@$core.Deprecated('Use eventDescriptor instead')
const Event$json = {
  '1': 'Event',
  '2': [
    {
      '1': 'clip_added',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.ClipRecord',
      '9': 0,
      '10': 'clipAdded'
    },
    {
      '1': 'clip_updated',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.ClipRecord',
      '9': 0,
      '10': 'clipUpdated'
    },
    {'1': 'clip_removed', '3': 3, '4': 1, '5': 9, '9': 0, '10': 'clipRemoved'},
    {
      '1': 'progress',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.TransferProgress',
      '9': 0,
      '10': 'progress'
    },
    {
      '1': 'device_changed',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.Device',
      '9': 0,
      '10': 'deviceChanged'
    },
    {
      '1': 'permission_changed',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ClipboardPermission',
      '9': 0,
      '10': 'permissionChanged'
    },
    {
      '1': 'status_changed',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.Status',
      '9': 0,
      '10': 'statusChanged'
    },
    {
      '1': 'clipboard_written',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.ClipboardWriteResult',
      '9': 0,
      '10': 'clipboardWritten'
    },
  ],
  '8': [
    {'1': 'payload'},
  ],
};

/// Descriptor for `Event`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDescriptor = $convert.base64Decode(
    'CgVFdmVudBI6CgpjbGlwX2FkZGVkGAEgASgLMhkucGFzdGVmZXJyeS52MS5DbGlwUmVjb3JkSA'
    'BSCWNsaXBBZGRlZBI+CgxjbGlwX3VwZGF0ZWQYAiABKAsyGS5wYXN0ZWZlcnJ5LnYxLkNsaXBS'
    'ZWNvcmRIAFILY2xpcFVwZGF0ZWQSIwoMY2xpcF9yZW1vdmVkGAMgASgJSABSC2NsaXBSZW1vdm'
    'VkEj0KCHByb2dyZXNzGAQgASgLMh8ucGFzdGVmZXJyeS52MS5UcmFuc2ZlclByb2dyZXNzSABS'
    'CHByb2dyZXNzEj4KDmRldmljZV9jaGFuZ2VkGAUgASgLMhUucGFzdGVmZXJyeS52MS5EZXZpY2'
    'VIAFINZGV2aWNlQ2hhbmdlZBJTChJwZXJtaXNzaW9uX2NoYW5nZWQYBiABKA4yIi5wYXN0ZWZl'
    'cnJ5LnYxLkNsaXBib2FyZFBlcm1pc3Npb25IAFIRcGVybWlzc2lvbkNoYW5nZWQSPgoOc3RhdH'
    'VzX2NoYW5nZWQYByABKAsyFS5wYXN0ZWZlcnJ5LnYxLlN0YXR1c0gAUg1zdGF0dXNDaGFuZ2Vk'
    'ElIKEWNsaXBib2FyZF93cml0dGVuGAggASgLMiMucGFzdGVmZXJyeS52MS5DbGlwYm9hcmRXcm'
    'l0ZVJlc3VsdEgAUhBjbGlwYm9hcmRXcml0dGVuQgkKB3BheWxvYWQ=');

@$core.Deprecated('Use notificationStatusDescriptor instead')
const NotificationStatus$json = {
  '1': 'NotificationStatus',
  '2': [
    {'1': 'supported', '3': 1, '4': 1, '5': 8, '10': 'supported'},
    {
      '1': 'actions_supported',
      '3': 2,
      '4': 1,
      '5': 8,
      '10': 'actionsSupported'
    },
    {'1': 'permission', '3': 3, '4': 1, '5': 9, '10': 'permission'},
  ],
};

/// Descriptor for `NotificationStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationStatusDescriptor = $convert.base64Decode(
    'ChJOb3RpZmljYXRpb25TdGF0dXMSHAoJc3VwcG9ydGVkGAEgASgIUglzdXBwb3J0ZWQSKwoRYW'
    'N0aW9uc19zdXBwb3J0ZWQYAiABKAhSEGFjdGlvbnNTdXBwb3J0ZWQSHgoKcGVybWlzc2lvbhgD'
    'IAEoCVIKcGVybWlzc2lvbg==');

@$core.Deprecated('Use notificationActionDescriptor instead')
const NotificationAction$json = {
  '1': 'NotificationAction',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {'1': 'action', '3': 2, '4': 1, '5': 9, '10': 'action'},
  ],
};

/// Descriptor for `NotificationAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationActionDescriptor = $convert.base64Decode(
    'ChJOb3RpZmljYXRpb25BY3Rpb24SFwoHY2xpcF9pZBgBIAEoCVIGY2xpcElkEhYKBmFjdGlvbh'
    'gCIAEoCVIGYWN0aW9u');

@$core.Deprecated('Use clipboardWriteResultDescriptor instead')
const ClipboardWriteResult$json = {
  '1': 'ClipboardWriteResult',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {
      '1': 'outcome',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ClipboardWriteOutcome',
      '10': 'outcome'
    },
    {'1': 'error', '3': 3, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ClipboardWriteResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clipboardWriteResultDescriptor = $convert.base64Decode(
    'ChRDbGlwYm9hcmRXcml0ZVJlc3VsdBIXCgdjbGlwX2lkGAEgASgJUgZjbGlwSWQSPgoHb3V0Y2'
    '9tZRgCIAEoDjIkLnBhc3RlZmVycnkudjEuQ2xpcGJvYXJkV3JpdGVPdXRjb21lUgdvdXRjb21l'
    'EhQKBWVycm9yGAMgASgJUgVlcnJvcg==');

@$core.Deprecated('Use listHistoryRequestDescriptor instead')
const ListHistoryRequest$json = {
  '1': 'ListHistoryRequest',
  '2': [
    {'1': 'limit', '3': 1, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'before_unix', '3': 2, '4': 1, '5': 3, '10': 'beforeUnix'},
    {
      '1': 'kind_filter',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ClipKind',
      '10': 'kindFilter'
    },
    {'1': 'query', '3': 4, '4': 1, '5': 9, '10': 'query'},
  ],
};

/// Descriptor for `ListHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHistoryRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0SGlzdG9yeVJlcXVlc3QSFAoFbGltaXQYASABKAVSBWxpbWl0Eh8KC2JlZm9yZV91bm'
    'l4GAIgASgDUgpiZWZvcmVVbml4EjgKC2tpbmRfZmlsdGVyGAMgASgOMhcucGFzdGVmZXJyeS52'
    'MS5DbGlwS2luZFIKa2luZEZpbHRlchIUCgVxdWVyeRgEIAEoCVIFcXVlcnk=');

@$core.Deprecated('Use listHistoryResponseDescriptor instead')
const ListHistoryResponse$json = {
  '1': 'ListHistoryResponse',
  '2': [
    {
      '1': 'records',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.pasteferry.v1.ClipRecord',
      '10': 'records'
    },
    {'1': 'searched', '3': 2, '4': 1, '5': 8, '10': 'searched'},
  ],
};

/// Descriptor for `ListHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHistoryResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0SGlzdG9yeVJlc3BvbnNlEjMKB3JlY29yZHMYASADKAsyGS5wYXN0ZWZlcnJ5LnYxLk'
    'NsaXBSZWNvcmRSB3JlY29yZHMSGgoIc2VhcmNoZWQYAiABKAhSCHNlYXJjaGVk');

@$core.Deprecated('Use deleteHistoryRequestDescriptor instead')
const DeleteHistoryRequest$json = {
  '1': 'DeleteHistoryRequest',
  '2': [
    {'1': 'ids', '3': 1, '4': 3, '5': 9, '10': 'ids'},
    {'1': 'all', '3': 2, '4': 1, '5': 8, '10': 'all'},
  ],
};

/// Descriptor for `DeleteHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteHistoryRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVIaXN0b3J5UmVxdWVzdBIQCgNpZHMYASADKAlSA2lkcxIQCgNhbGwYAiABKAhSA2'
    'FsbA==');

@$core.Deprecated('Use fetchRequestDescriptor instead')
const FetchRequest$json = {
  '1': 'FetchRequest',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {
      '1': 'preserve_clipboard',
      '3': 2,
      '4': 1,
      '5': 8,
      '10': 'preserveClipboard'
    },
  ],
};

/// Descriptor for `FetchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fetchRequestDescriptor = $convert.base64Decode(
    'CgxGZXRjaFJlcXVlc3QSFwoHY2xpcF9pZBgBIAEoCVIGY2xpcElkEi0KEnByZXNlcnZlX2NsaX'
    'Bib2FyZBgCIAEoCFIRcHJlc2VydmVDbGlwYm9hcmQ=');

@$core.Deprecated('Use applyToClipboardRequestDescriptor instead')
const ApplyToClipboardRequest$json = {
  '1': 'ApplyToClipboardRequest',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
  ],
};

/// Descriptor for `ApplyToClipboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List applyToClipboardRequestDescriptor =
    $convert.base64Decode(
        'ChdBcHBseVRvQ2xpcGJvYXJkUmVxdWVzdBIXCgdjbGlwX2lkGAEgASgJUgZjbGlwSWQ=');

@$core.Deprecated('Use listDevicesResponseDescriptor instead')
const ListDevicesResponse$json = {
  '1': 'ListDevicesResponse',
  '2': [
    {
      '1': 'self',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.Device',
      '10': 'self'
    },
    {
      '1': 'peers',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.pasteferry.v1.Device',
      '10': 'peers'
    },
  ],
};

/// Descriptor for `ListDevicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDevicesResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0RGV2aWNlc1Jlc3BvbnNlEikKBHNlbGYYASABKAsyFS5wYXN0ZWZlcnJ5LnYxLkRldm'
    'ljZVIEc2VsZhIrCgVwZWVycxgCIAMoCzIVLnBhc3RlZmVycnkudjEuRGV2aWNlUgVwZWVycw==');

@$core.Deprecated('Use createPairingCodeResponseDescriptor instead')
const CreatePairingCodeResponse$json = {
  '1': 'CreatePairingCodeResponse',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'expires_at_unix', '3': 2, '4': 1, '5': 3, '10': 'expiresAtUnix'},
    {'1': 'lan_addresses', '3': 3, '4': 3, '5': 9, '10': 'lanAddresses'},
  ],
};

/// Descriptor for `CreatePairingCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPairingCodeResponseDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVQYWlyaW5nQ29kZVJlc3BvbnNlEhIKBGNvZGUYASABKAlSBGNvZGUSJgoPZXhwaX'
    'Jlc19hdF91bml4GAIgASgDUg1leHBpcmVzQXRVbml4EiMKDWxhbl9hZGRyZXNzZXMYAyADKAlS'
    'DGxhbkFkZHJlc3Nlcw==');

@$core.Deprecated('Use redeemPairingCodeRequestDescriptor instead')
const RedeemPairingCodeRequest$json = {
  '1': 'RedeemPairingCodeRequest',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `RedeemPairingCodeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redeemPairingCodeRequestDescriptor =
    $convert.base64Decode(
        'ChhSZWRlZW1QYWlyaW5nQ29kZVJlcXVlc3QSEgoEY29kZRgBIAEoCVIEY29kZRIYCgdhZGRyZX'
        'NzGAIgASgJUgdhZGRyZXNz');

@$core.Deprecated('Use redeemPairingCodeResponseDescriptor instead')
const RedeemPairingCodeResponse$json = {
  '1': 'RedeemPairingCodeResponse',
  '2': [
    {
      '1': 'peer',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.Device',
      '10': 'peer'
    },
    {'1': 'pairing_session', '3': 2, '4': 1, '5': 9, '10': 'pairingSession'},
  ],
};

/// Descriptor for `RedeemPairingCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redeemPairingCodeResponseDescriptor = $convert.base64Decode(
    'ChlSZWRlZW1QYWlyaW5nQ29kZVJlc3BvbnNlEikKBHBlZXIYASABKAsyFS5wYXN0ZWZlcnJ5Ln'
    'YxLkRldmljZVIEcGVlchInCg9wYWlyaW5nX3Nlc3Npb24YAiABKAlSDnBhaXJpbmdTZXNzaW9u');

@$core.Deprecated('Use confirmPairingRequestDescriptor instead')
const ConfirmPairingRequest$json = {
  '1': 'ConfirmPairingRequest',
  '2': [
    {'1': 'pairing_session', '3': 1, '4': 1, '5': 9, '10': 'pairingSession'},
    {'1': 'accept', '3': 2, '4': 1, '5': 8, '10': 'accept'},
  ],
};

/// Descriptor for `ConfirmPairingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmPairingRequestDescriptor = $convert.base64Decode(
    'ChVDb25maXJtUGFpcmluZ1JlcXVlc3QSJwoPcGFpcmluZ19zZXNzaW9uGAEgASgJUg5wYWlyaW'
    '5nU2Vzc2lvbhIWCgZhY2NlcHQYAiABKAhSBmFjY2VwdA==');

@$core.Deprecated('Use nearbyDeviceDescriptor instead')
const NearbyDevice$json = {
  '1': 'NearbyDevice',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'platform', '3': 3, '4': 1, '5': 9, '10': 'platform'},
    {'1': 'paired', '3': 4, '4': 1, '5': 8, '10': 'paired'},
  ],
};

/// Descriptor for `NearbyDevice`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nearbyDeviceDescriptor = $convert.base64Decode(
    'CgxOZWFyYnlEZXZpY2USDgoCaWQYASABKAlSAmlkEhIKBG5hbWUYAiABKAlSBG5hbWUSGgoIcG'
    'xhdGZvcm0YAyABKAlSCHBsYXRmb3JtEhYKBnBhaXJlZBgEIAEoCFIGcGFpcmVk');

@$core.Deprecated('Use nearbyDevicesDescriptor instead')
const NearbyDevices$json = {
  '1': 'NearbyDevices',
  '2': [
    {
      '1': 'devices',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.pasteferry.v1.NearbyDevice',
      '10': 'devices'
    },
  ],
};

/// Descriptor for `NearbyDevices`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nearbyDevicesDescriptor = $convert.base64Decode(
    'Cg1OZWFyYnlEZXZpY2VzEjUKB2RldmljZXMYASADKAsyGy5wYXN0ZWZlcnJ5LnYxLk5lYXJieU'
    'RldmljZVIHZGV2aWNlcw==');

@$core.Deprecated('Use pairNearbyRequestDescriptor instead')
const PairNearbyRequest$json = {
  '1': 'PairNearbyRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `PairNearbyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pairNearbyRequestDescriptor =
    $convert.base64Decode('ChFQYWlyTmVhcmJ5UmVxdWVzdBIOCgJpZBgBIAEoCVICaWQ=');

@$core.Deprecated('Use unpairRequestDescriptor instead')
const UnpairRequest$json = {
  '1': 'UnpairRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
  ],
};

/// Descriptor for `UnpairRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unpairRequestDescriptor = $convert.base64Decode(
    'Cg1VbnBhaXJSZXF1ZXN0EhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQ=');

@$core.Deprecated('Use setDeviceAddressRequestDescriptor instead')
const SetDeviceAddressRequest$json = {
  '1': 'SetDeviceAddressRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `SetDeviceAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setDeviceAddressRequestDescriptor =
    $convert.base64Decode(
        'ChdTZXREZXZpY2VBZGRyZXNzUmVxdWVzdBIbCglkZXZpY2VfaWQYASABKAlSCGRldmljZUlkEh'
        'gKB2FkZHJlc3MYAiABKAlSB2FkZHJlc3M=');

@$core.Deprecated('Use configDescriptor instead')
const Config$json = {
  '1': 'Config',
  '2': [
    {
      '1': 'auto_sync_threshold_bytes',
      '3': 1,
      '4': 1,
      '5': 3,
      '10': 'autoSyncThresholdBytes'
    },
    {
      '1': 'history_ttl_seconds',
      '3': 2,
      '4': 1,
      '5': 3,
      '10': 'historyTtlSeconds'
    },
    {'1': 'cache_ttl_seconds', '3': 3, '4': 1, '5': 3, '10': 'cacheTtlSeconds'},
    {'1': 'sync_text', '3': 4, '4': 1, '5': 8, '10': 'syncText'},
    {'1': 'sync_html', '3': 5, '4': 1, '5': 8, '10': 'syncHtml'},
    {'1': 'sync_image', '3': 6, '4': 1, '5': 8, '10': 'syncImage'},
    {'1': 'sync_file', '3': 7, '4': 1, '5': 8, '10': 'syncFile'},
    {
      '1': 'auto_apply_to_clipboard',
      '3': 8,
      '4': 1,
      '5': 8,
      '10': 'autoApplyToClipboard'
    },
    {'1': 'launch_at_login', '3': 9, '4': 1, '5': 8, '10': 'launchAtLogin'},
    {'1': 'device_name', '3': 10, '4': 1, '5': 9, '10': 'deviceName'},
    {'1': 'signaling_url', '3': 11, '4': 1, '5': 9, '10': 'signalingUrl'},
    {'1': 'only_own_stun', '3': 12, '4': 1, '5': 8, '10': 'onlyOwnStun'},
    {'1': 'language', '3': 13, '4': 1, '5': 9, '10': 'language'},
    {'1': 'lan_disabled', '3': 14, '4': 1, '5': 8, '10': 'lanDisabled'},
    {
      '1': 'service_mode',
      '3': 15,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ServiceMode',
      '10': 'serviceMode'
    },
    {
      '1': 'custom_signaling_url',
      '3': 16,
      '4': 1,
      '5': 9,
      '10': 'customSignalingUrl'
    },
    {
      '1': 'official_notice_ack',
      '3': 17,
      '4': 1,
      '5': 9,
      '10': 'officialNoticeAck'
    },
  ],
};

/// Descriptor for `Config`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List configDescriptor = $convert.base64Decode(
    'CgZDb25maWcSOQoZYXV0b19zeW5jX3RocmVzaG9sZF9ieXRlcxgBIAEoA1IWYXV0b1N5bmNUaH'
    'Jlc2hvbGRCeXRlcxIuChNoaXN0b3J5X3R0bF9zZWNvbmRzGAIgASgDUhFoaXN0b3J5VHRsU2Vj'
    'b25kcxIqChFjYWNoZV90dGxfc2Vjb25kcxgDIAEoA1IPY2FjaGVUdGxTZWNvbmRzEhsKCXN5bm'
    'NfdGV4dBgEIAEoCFIIc3luY1RleHQSGwoJc3luY19odG1sGAUgASgIUghzeW5jSHRtbBIdCgpz'
    'eW5jX2ltYWdlGAYgASgIUglzeW5jSW1hZ2USGwoJc3luY19maWxlGAcgASgIUghzeW5jRmlsZR'
    'I1ChdhdXRvX2FwcGx5X3RvX2NsaXBib2FyZBgIIAEoCFIUYXV0b0FwcGx5VG9DbGlwYm9hcmQS'
    'JgoPbGF1bmNoX2F0X2xvZ2luGAkgASgIUg1sYXVuY2hBdExvZ2luEh8KC2RldmljZV9uYW1lGA'
    'ogASgJUgpkZXZpY2VOYW1lEiMKDXNpZ25hbGluZ191cmwYCyABKAlSDHNpZ25hbGluZ1VybBIi'
    'Cg1vbmx5X293bl9zdHVuGAwgASgIUgtvbmx5T3duU3R1bhIaCghsYW5ndWFnZRgNIAEoCVIIbG'
    'FuZ3VhZ2USIQoMbGFuX2Rpc2FibGVkGA4gASgIUgtsYW5EaXNhYmxlZBI9CgxzZXJ2aWNlX21v'
    'ZGUYDyABKA4yGi5wYXN0ZWZlcnJ5LnYxLlNlcnZpY2VNb2RlUgtzZXJ2aWNlTW9kZRIwChRjdX'
    'N0b21fc2lnbmFsaW5nX3VybBgQIAEoCVISY3VzdG9tU2lnbmFsaW5nVXJsEi4KE29mZmljaWFs'
    'X25vdGljZV9hY2sYESABKAlSEW9mZmljaWFsTm90aWNlQWNr');

@$core.Deprecated('Use officialUsageDescriptor instead')
const OfficialUsage$json = {
  '1': 'OfficialUsage',
  '2': [
    {'1': 'available', '3': 1, '4': 1, '5': 8, '10': 'available'},
    {'1': 'day', '3': 2, '4': 1, '5': 9, '10': 'day'},
    {'1': 'quota_group', '3': 3, '4': 1, '5': 9, '10': 'quotaGroup'},
    {'1': 'upload_bytes', '3': 4, '4': 1, '5': 3, '10': 'uploadBytes'},
    {'1': 'download_bytes', '3': 5, '4': 1, '5': 3, '10': 'downloadBytes'},
    {'1': 'limit_bytes', '3': 6, '4': 1, '5': 3, '10': 'limitBytes'},
    {'1': 'remaining_bytes', '3': 7, '4': 1, '5': 3, '10': 'remainingBytes'},
    {'1': 'reset_at_unix', '3': 8, '4': 1, '5': 3, '10': 'resetAtUnix'},
    {'1': 'updated_at_unix', '3': 9, '4': 1, '5': 3, '10': 'updatedAtUnix'},
    {'1': 'reason', '3': 10, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `OfficialUsage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List officialUsageDescriptor = $convert.base64Decode(
    'Cg1PZmZpY2lhbFVzYWdlEhwKCWF2YWlsYWJsZRgBIAEoCFIJYXZhaWxhYmxlEhAKA2RheRgCIA'
    'EoCVIDZGF5Eh8KC3F1b3RhX2dyb3VwGAMgASgJUgpxdW90YUdyb3VwEiEKDHVwbG9hZF9ieXRl'
    'cxgEIAEoA1ILdXBsb2FkQnl0ZXMSJQoOZG93bmxvYWRfYnl0ZXMYBSABKANSDWRvd25sb2FkQn'
    'l0ZXMSHwoLbGltaXRfYnl0ZXMYBiABKANSCmxpbWl0Qnl0ZXMSJwoPcmVtYWluaW5nX2J5dGVz'
    'GAcgASgDUg5yZW1haW5pbmdCeXRlcxIiCg1yZXNldF9hdF91bml4GAggASgDUgtyZXNldEF0VW'
    '5peBImCg91cGRhdGVkX2F0X3VuaXgYCSABKANSDXVwZGF0ZWRBdFVuaXgSFgoGcmVhc29uGAog'
    'ASgJUgZyZWFzb24=');

@$core.Deprecated('Use networkInfoDescriptor instead')
const NetworkInfo$json = {
  '1': 'NetworkInfo',
  '2': [
    {'1': 'local_port', '3': 1, '4': 1, '5': 5, '10': 'localPort'},
    {
      '1': 'egresses',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.pasteferry.v1.NetworkEgress',
      '10': 'egresses'
    },
    {'1': 'probed', '3': 3, '4': 1, '5': 5, '10': 'probed'},
    {'1': 'answered', '3': 4, '4': 1, '5': 5, '10': 'answered'},
    {'1': 'probed_at_unix', '3': 5, '4': 1, '5': 3, '10': 'probedAtUnix'},
    {
      '1': 'lan',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.LanInfo',
      '10': 'lan'
    },
  ],
};

/// Descriptor for `NetworkInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkInfoDescriptor = $convert.base64Decode(
    'CgtOZXR3b3JrSW5mbxIdCgpsb2NhbF9wb3J0GAEgASgFUglsb2NhbFBvcnQSOAoIZWdyZXNzZX'
    'MYAiADKAsyHC5wYXN0ZWZlcnJ5LnYxLk5ldHdvcmtFZ3Jlc3NSCGVncmVzc2VzEhYKBnByb2Jl'
    'ZBgDIAEoBVIGcHJvYmVkEhoKCGFuc3dlcmVkGAQgASgFUghhbnN3ZXJlZBIkCg5wcm9iZWRfYX'
    'RfdW5peBgFIAEoA1IMcHJvYmVkQXRVbml4EigKA2xhbhgGIAEoCzIWLnBhc3RlZmVycnkudjEu'
    'TGFuSW5mb1IDbGFu');

@$core.Deprecated('Use lanInfoDescriptor instead')
const LanInfo$json = {
  '1': 'LanInfo',
  '2': [
    {'1': 'enabled', '3': 1, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'port', '3': 2, '4': 1, '5': 5, '10': 'port'},
    {'1': 'local_addresses', '3': 3, '4': 3, '5': 9, '10': 'localAddresses'},
    {'1': 'discovered', '3': 4, '4': 1, '5': 5, '10': 'discovered'},
    {
      '1': 'peers',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.pasteferry.v1.LanPeer',
      '10': 'peers'
    },
    {'1': 'blocked', '3': 6, '4': 1, '5': 8, '10': 'blocked'},
  ],
};

/// Descriptor for `LanInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lanInfoDescriptor = $convert.base64Decode(
    'CgdMYW5JbmZvEhgKB2VuYWJsZWQYASABKAhSB2VuYWJsZWQSEgoEcG9ydBgCIAEoBVIEcG9ydB'
    'InCg9sb2NhbF9hZGRyZXNzZXMYAyADKAlSDmxvY2FsQWRkcmVzc2VzEh4KCmRpc2NvdmVyZWQY'
    'BCABKAVSCmRpc2NvdmVyZWQSLAoFcGVlcnMYBSADKAsyFi5wYXN0ZWZlcnJ5LnYxLkxhblBlZX'
    'JSBXBlZXJzEhgKB2Jsb2NrZWQYBiABKAhSB2Jsb2NrZWQ=');

@$core.Deprecated('Use lanPeerDescriptor instead')
const LanPeer$json = {
  '1': 'LanPeer',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'device_name', '3': 2, '4': 1, '5': 9, '10': 'deviceName'},
    {'1': 'reachable', '3': 3, '4': 1, '5': 8, '10': 'reachable'},
    {'1': 'address', '3': 4, '4': 1, '5': 9, '10': 'address'},
    {'1': 'known_addresses', '3': 5, '4': 3, '5': 9, '10': 'knownAddresses'},
    {'1': 'last_seen_unix', '3': 6, '4': 1, '5': 3, '10': 'lastSeenUnix'},
  ],
};

/// Descriptor for `LanPeer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lanPeerDescriptor = $convert.base64Decode(
    'CgdMYW5QZWVyEhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQSHwoLZGV2aWNlX25hbWUYAi'
    'ABKAlSCmRldmljZU5hbWUSHAoJcmVhY2hhYmxlGAMgASgIUglyZWFjaGFibGUSGAoHYWRkcmVz'
    'cxgEIAEoCVIHYWRkcmVzcxInCg9rbm93bl9hZGRyZXNzZXMYBSADKAlSDmtub3duQWRkcmVzc2'
    'VzEiQKDmxhc3Rfc2Vlbl91bml4GAYgASgDUgxsYXN0U2VlblVuaXg=');

@$core.Deprecated('Use networkEgressDescriptor instead')
const NetworkEgress$json = {
  '1': 'NetworkEgress',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    {'1': 'port_preserved', '3': 2, '4': 1, '5': 8, '10': 'portPreserved'},
    {'1': 'guessed', '3': 3, '4': 1, '5': 8, '10': 'guessed'},
    {'1': 'via', '3': 4, '4': 3, '5': 9, '10': 'via'},
  ],
};

/// Descriptor for `NetworkEgress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkEgressDescriptor = $convert.base64Decode(
    'Cg1OZXR3b3JrRWdyZXNzEhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3MSJQoOcG9ydF9wcmVzZX'
    'J2ZWQYAiABKAhSDXBvcnRQcmVzZXJ2ZWQSGAoHZ3Vlc3NlZBgDIAEoCFIHZ3Vlc3NlZBIQCgN2'
    'aWEYBCADKAlSA3ZpYQ==');

@$core.Deprecated('Use statusDescriptor instead')
const Status$json = {
  '1': 'Status',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'device_name', '3': 2, '4': 1, '5': 9, '10': 'deviceName'},
    {
      '1': 'signaling_connected',
      '3': 3,
      '4': 1,
      '5': 8,
      '10': 'signalingConnected'
    },
    {
      '1': 'clipboard_permission',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.pasteferry.v1.ClipboardPermission',
      '10': 'clipboardPermission'
    },
    {'1': 'cache_bytes_used', '3': 5, '4': 1, '5': 3, '10': 'cacheBytesUsed'},
    {'1': 'peers_online', '3': 6, '4': 1, '5': 5, '10': 'peersOnline'},
    {'1': 'version', '3': 7, '4': 1, '5': 9, '10': 'version'},
    {
      '1': 'server_configured',
      '3': 8,
      '4': 1,
      '5': 8,
      '10': 'serverConfigured'
    },
    {'1': 'lan_enabled', '3': 9, '4': 1, '5': 8, '10': 'lanEnabled'},
    {'1': 'lan_blocked', '3': 10, '4': 1, '5': 8, '10': 'lanBlocked'},
    {
      '1': 'verbose_log_until_unix',
      '3': 11,
      '4': 1,
      '5': 3,
      '10': 'verboseLogUntilUnix'
    },
    {
      '1': 'official_usage',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.pasteferry.v1.OfficialUsage',
      '10': 'officialUsage'
    },
  ],
};

/// Descriptor for `Status`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statusDescriptor = $convert.base64Decode(
    'CgZTdGF0dXMSGwoJZGV2aWNlX2lkGAEgASgJUghkZXZpY2VJZBIfCgtkZXZpY2VfbmFtZRgCIA'
    'EoCVIKZGV2aWNlTmFtZRIvChNzaWduYWxpbmdfY29ubmVjdGVkGAMgASgIUhJzaWduYWxpbmdD'
    'b25uZWN0ZWQSVQoUY2xpcGJvYXJkX3Blcm1pc3Npb24YBCABKA4yIi5wYXN0ZWZlcnJ5LnYxLk'
    'NsaXBib2FyZFBlcm1pc3Npb25SE2NsaXBib2FyZFBlcm1pc3Npb24SKAoQY2FjaGVfYnl0ZXNf'
    'dXNlZBgFIAEoA1IOY2FjaGVCeXRlc1VzZWQSIQoMcGVlcnNfb25saW5lGAYgASgFUgtwZWVyc0'
    '9ubGluZRIYCgd2ZXJzaW9uGAcgASgJUgd2ZXJzaW9uEisKEXNlcnZlcl9jb25maWd1cmVkGAgg'
    'ASgIUhBzZXJ2ZXJDb25maWd1cmVkEh8KC2xhbl9lbmFibGVkGAkgASgIUgpsYW5FbmFibGVkEh'
    '8KC2xhbl9ibG9ja2VkGAogASgIUgpsYW5CbG9ja2VkEjMKFnZlcmJvc2VfbG9nX3VudGlsX3Vu'
    'aXgYCyABKANSE3ZlcmJvc2VMb2dVbnRpbFVuaXgSQwoOb2ZmaWNpYWxfdXNhZ2UYDCABKAsyHC'
    '5wYXN0ZWZlcnJ5LnYxLk9mZmljaWFsVXNhZ2VSDW9mZmljaWFsVXNhZ2U=');

@$core.Deprecated('Use exportDiagnosticsRequestDescriptor instead')
const ExportDiagnosticsRequest$json = {
  '1': 'ExportDiagnosticsRequest',
  '2': [
    {'1': 'log_lines', '3': 1, '4': 1, '5': 5, '10': 'logLines'},
    {'1': 'ui_version', '3': 2, '4': 1, '5': 9, '10': 'uiVersion'},
  ],
};

/// Descriptor for `ExportDiagnosticsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportDiagnosticsRequestDescriptor =
    $convert.base64Decode(
        'ChhFeHBvcnREaWFnbm9zdGljc1JlcXVlc3QSGwoJbG9nX2xpbmVzGAEgASgFUghsb2dMaW5lcx'
        'IdCgp1aV92ZXJzaW9uGAIgASgJUgl1aVZlcnNpb24=');

@$core.Deprecated('Use diagnosticsDescriptor instead')
const Diagnostics$json = {
  '1': 'Diagnostics',
  '2': [
    {'1': 'report', '3': 1, '4': 1, '5': 9, '10': 'report'},
  ],
};

/// Descriptor for `Diagnostics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List diagnosticsDescriptor = $convert
    .base64Decode('CgtEaWFnbm9zdGljcxIWCgZyZXBvcnQYASABKAlSBnJlcG9ydA==');

@$core.Deprecated('Use setVerboseLogRequestDescriptor instead')
const SetVerboseLogRequest$json = {
  '1': 'SetVerboseLogRequest',
  '2': [
    {'1': 'minutes', '3': 1, '4': 1, '5': 5, '10': 'minutes'},
  ],
};

/// Descriptor for `SetVerboseLogRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setVerboseLogRequestDescriptor =
    $convert.base64Decode(
        'ChRTZXRWZXJib3NlTG9nUmVxdWVzdBIYCgdtaW51dGVzGAEgASgFUgdtaW51dGVz');

@$core.Deprecated('Use getImagePreviewRequestDescriptor instead')
const GetImagePreviewRequest$json = {
  '1': 'GetImagePreviewRequest',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {'1': 'index', '3': 2, '4': 1, '5': 5, '10': 'index'},
  ],
};

/// Descriptor for `GetImagePreviewRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImagePreviewRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRJbWFnZVByZXZpZXdSZXF1ZXN0EhcKB2NsaXBfaWQYASABKAlSBmNsaXBJZBIUCgVpbm'
        'RleBgCIAEoBVIFaW5kZXg=');

@$core.Deprecated('Use getImagePreviewResponseDescriptor instead')
const GetImagePreviewResponse$json = {
  '1': 'GetImagePreviewResponse',
  '2': [
    {'1': 'path', '3': 1, '4': 1, '5': 9, '10': 'path'},
    {'1': 'count', '3': 2, '4': 1, '5': 5, '10': 'count'},
    {'1': 'name', '3': 3, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetImagePreviewResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImagePreviewResponseDescriptor =
    $convert.base64Decode(
        'ChdHZXRJbWFnZVByZXZpZXdSZXNwb25zZRISCgRwYXRoGAEgASgJUgRwYXRoEhQKBWNvdW50GA'
        'IgASgFUgVjb3VudBISCgRuYW1lGAMgASgJUgRuYW1l');

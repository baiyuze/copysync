// This is a generated file - do not edit.
//
// Generated from copysync/v1/daemon.proto.

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
      '6': '.copysync.v1.ClipKind',
      '10': 'kind'
    },
    {
      '1': 'status',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.copysync.v1.ClipStatus',
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
      '6': '.copysync.v1.ClipItem',
      '10': 'items'
    },
    {'1': 'total_size', '3': 8, '4': 1, '5': 3, '10': 'totalSize'},
    {'1': 'text_preview', '3': 9, '4': 1, '5': 9, '10': 'textPreview'},
    {'1': 'created_at_unix', '3': 10, '4': 1, '5': 3, '10': 'createdAtUnix'},
    {'1': 'expires_at_unix', '3': 11, '4': 1, '5': 3, '10': 'expiresAtUnix'},
    {'1': 'error', '3': 12, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ClipRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clipRecordDescriptor = $convert.base64Decode(
    'CgpDbGlwUmVjb3JkEg4KAmlkGAEgASgJUgJpZBIpCgRraW5kGAIgASgOMhUuY29weXN5bmMudj'
    'EuQ2xpcEtpbmRSBGtpbmQSLwoGc3RhdHVzGAMgASgOMhcuY29weXN5bmMudjEuQ2xpcFN0YXR1'
    'c1IGc3RhdHVzEigKEG9yaWdpbl9kZXZpY2VfaWQYBCABKAlSDm9yaWdpbkRldmljZUlkEiwKEm'
    '9yaWdpbl9kZXZpY2VfbmFtZRgFIAEoCVIQb3JpZ2luRGV2aWNlTmFtZRIaCghvdXRnb2luZxgG'
    'IAEoCFIIb3V0Z29pbmcSKwoFaXRlbXMYByADKAsyFS5jb3B5c3luYy52MS5DbGlwSXRlbVIFaX'
    'RlbXMSHQoKdG90YWxfc2l6ZRgIIAEoA1IJdG90YWxTaXplEiEKDHRleHRfcHJldmlldxgJIAEo'
    'CVILdGV4dFByZXZpZXcSJgoPY3JlYXRlZF9hdF91bml4GAogASgDUg1jcmVhdGVkQXRVbml4Ei'
    'YKD2V4cGlyZXNfYXRfdW5peBgLIAEoA1INZXhwaXJlc0F0VW5peBIUCgVlcnJvchgMIAEoCVIF'
    'ZXJyb3I=');

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
      '6': '.copysync.v1.ConnectionKind',
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
    '0YAyABKAlSCHBsYXRmb3JtEhYKBm9ubGluZRgEIAEoCFIGb25saW5lEjsKCmNvbm5lY3Rpb24Y'
    'BSABKA4yGy5jb3B5c3luYy52MS5Db25uZWN0aW9uS2luZFIKY29ubmVjdGlvbhI0ChZwdWJsaW'
    'Nfa2V5X2ZpbmdlcnByaW50GAYgASgJUhRwdWJsaWNLZXlGaW5nZXJwcmludBIkCg5wYWlyZWRf'
    'YXRfdW5peBgHIAEoA1IMcGFpcmVkQXRVbml4EicKD3BhaXJpbmdfc2Vzc2lvbhgIIAEoCVIOcG'
    'FpcmluZ1Nlc3Npb24=');

@$core.Deprecated('Use eventDescriptor instead')
const Event$json = {
  '1': 'Event',
  '2': [
    {
      '1': 'clip_added',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.ClipRecord',
      '9': 0,
      '10': 'clipAdded'
    },
    {
      '1': 'clip_updated',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.ClipRecord',
      '9': 0,
      '10': 'clipUpdated'
    },
    {'1': 'clip_removed', '3': 3, '4': 1, '5': 9, '9': 0, '10': 'clipRemoved'},
    {
      '1': 'progress',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.TransferProgress',
      '9': 0,
      '10': 'progress'
    },
    {
      '1': 'device_changed',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.Device',
      '9': 0,
      '10': 'deviceChanged'
    },
    {
      '1': 'permission_changed',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.copysync.v1.ClipboardPermission',
      '9': 0,
      '10': 'permissionChanged'
    },
    {
      '1': 'status_changed',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.Status',
      '9': 0,
      '10': 'statusChanged'
    },
  ],
  '8': [
    {'1': 'payload'},
  ],
};

/// Descriptor for `Event`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDescriptor = $convert.base64Decode(
    'CgVFdmVudBI4CgpjbGlwX2FkZGVkGAEgASgLMhcuY29weXN5bmMudjEuQ2xpcFJlY29yZEgAUg'
    'ljbGlwQWRkZWQSPAoMY2xpcF91cGRhdGVkGAIgASgLMhcuY29weXN5bmMudjEuQ2xpcFJlY29y'
    'ZEgAUgtjbGlwVXBkYXRlZBIjCgxjbGlwX3JlbW92ZWQYAyABKAlIAFILY2xpcFJlbW92ZWQSOw'
    'oIcHJvZ3Jlc3MYBCABKAsyHS5jb3B5c3luYy52MS5UcmFuc2ZlclByb2dyZXNzSABSCHByb2dy'
    'ZXNzEjwKDmRldmljZV9jaGFuZ2VkGAUgASgLMhMuY29weXN5bmMudjEuRGV2aWNlSABSDWRldm'
    'ljZUNoYW5nZWQSUQoScGVybWlzc2lvbl9jaGFuZ2VkGAYgASgOMiAuY29weXN5bmMudjEuQ2xp'
    'cGJvYXJkUGVybWlzc2lvbkgAUhFwZXJtaXNzaW9uQ2hhbmdlZBI8Cg5zdGF0dXNfY2hhbmdlZB'
    'gHIAEoCzITLmNvcHlzeW5jLnYxLlN0YXR1c0gAUg1zdGF0dXNDaGFuZ2VkQgkKB3BheWxvYWQ=');

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
      '6': '.copysync.v1.ClipKind',
      '10': 'kindFilter'
    },
  ],
};

/// Descriptor for `ListHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHistoryRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0SGlzdG9yeVJlcXVlc3QSFAoFbGltaXQYASABKAVSBWxpbWl0Eh8KC2JlZm9yZV91bm'
    'l4GAIgASgDUgpiZWZvcmVVbml4EjYKC2tpbmRfZmlsdGVyGAMgASgOMhUuY29weXN5bmMudjEu'
    'Q2xpcEtpbmRSCmtpbmRGaWx0ZXI=');

@$core.Deprecated('Use listHistoryResponseDescriptor instead')
const ListHistoryResponse$json = {
  '1': 'ListHistoryResponse',
  '2': [
    {
      '1': 'records',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.copysync.v1.ClipRecord',
      '10': 'records'
    },
  ],
};

/// Descriptor for `ListHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHistoryResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0SGlzdG9yeVJlc3BvbnNlEjEKB3JlY29yZHMYASADKAsyFy5jb3B5c3luYy52MS5DbG'
    'lwUmVjb3JkUgdyZWNvcmRz');

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
      '6': '.copysync.v1.Device',
      '10': 'self'
    },
    {
      '1': 'peers',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.copysync.v1.Device',
      '10': 'peers'
    },
  ],
};

/// Descriptor for `ListDevicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDevicesResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0RGV2aWNlc1Jlc3BvbnNlEicKBHNlbGYYASABKAsyEy5jb3B5c3luYy52MS5EZXZpY2'
    'VSBHNlbGYSKQoFcGVlcnMYAiADKAsyEy5jb3B5c3luYy52MS5EZXZpY2VSBXBlZXJz');

@$core.Deprecated('Use createPairingCodeResponseDescriptor instead')
const CreatePairingCodeResponse$json = {
  '1': 'CreatePairingCodeResponse',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'expires_at_unix', '3': 2, '4': 1, '5': 3, '10': 'expiresAtUnix'},
  ],
};

/// Descriptor for `CreatePairingCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPairingCodeResponseDescriptor =
    $convert.base64Decode(
        'ChlDcmVhdGVQYWlyaW5nQ29kZVJlc3BvbnNlEhIKBGNvZGUYASABKAlSBGNvZGUSJgoPZXhwaX'
        'Jlc19hdF91bml4GAIgASgDUg1leHBpcmVzQXRVbml4');

@$core.Deprecated('Use redeemPairingCodeRequestDescriptor instead')
const RedeemPairingCodeRequest$json = {
  '1': 'RedeemPairingCodeRequest',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
  ],
};

/// Descriptor for `RedeemPairingCodeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redeemPairingCodeRequestDescriptor =
    $convert.base64Decode(
        'ChhSZWRlZW1QYWlyaW5nQ29kZVJlcXVlc3QSEgoEY29kZRgBIAEoCVIEY29kZQ==');

@$core.Deprecated('Use redeemPairingCodeResponseDescriptor instead')
const RedeemPairingCodeResponse$json = {
  '1': 'RedeemPairingCodeResponse',
  '2': [
    {
      '1': 'peer',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.Device',
      '10': 'peer'
    },
    {'1': 'pairing_session', '3': 2, '4': 1, '5': 9, '10': 'pairingSession'},
  ],
};

/// Descriptor for `RedeemPairingCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redeemPairingCodeResponseDescriptor = $convert.base64Decode(
    'ChlSZWRlZW1QYWlyaW5nQ29kZVJlc3BvbnNlEicKBHBlZXIYASABKAsyEy5jb3B5c3luYy52MS'
    '5EZXZpY2VSBHBlZXISJwoPcGFpcmluZ19zZXNzaW9uGAIgASgJUg5wYWlyaW5nU2Vzc2lvbg==');

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
    'FuZ3VhZ2U=');

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
      '6': '.copysync.v1.NetworkEgress',
      '10': 'egresses'
    },
    {'1': 'probed', '3': 3, '4': 1, '5': 5, '10': 'probed'},
    {'1': 'answered', '3': 4, '4': 1, '5': 5, '10': 'answered'},
    {'1': 'probed_at_unix', '3': 5, '4': 1, '5': 3, '10': 'probedAtUnix'},
  ],
};

/// Descriptor for `NetworkInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkInfoDescriptor = $convert.base64Decode(
    'CgtOZXR3b3JrSW5mbxIdCgpsb2NhbF9wb3J0GAEgASgFUglsb2NhbFBvcnQSNgoIZWdyZXNzZX'
    'MYAiADKAsyGi5jb3B5c3luYy52MS5OZXR3b3JrRWdyZXNzUghlZ3Jlc3NlcxIWCgZwcm9iZWQY'
    'AyABKAVSBnByb2JlZBIaCghhbnN3ZXJlZBgEIAEoBVIIYW5zd2VyZWQSJAoOcHJvYmVkX2F0X3'
    'VuaXgYBSABKANSDHByb2JlZEF0VW5peA==');

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
      '6': '.copysync.v1.ClipboardPermission',
      '10': 'clipboardPermission'
    },
    {'1': 'cache_bytes_used', '3': 5, '4': 1, '5': 3, '10': 'cacheBytesUsed'},
    {'1': 'peers_online', '3': 6, '4': 1, '5': 5, '10': 'peersOnline'},
    {'1': 'version', '3': 7, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `Status`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statusDescriptor = $convert.base64Decode(
    'CgZTdGF0dXMSGwoJZGV2aWNlX2lkGAEgASgJUghkZXZpY2VJZBIfCgtkZXZpY2VfbmFtZRgCIA'
    'EoCVIKZGV2aWNlTmFtZRIvChNzaWduYWxpbmdfY29ubmVjdGVkGAMgASgIUhJzaWduYWxpbmdD'
    'b25uZWN0ZWQSUwoUY2xpcGJvYXJkX3Blcm1pc3Npb24YBCABKA4yIC5jb3B5c3luYy52MS5DbG'
    'lwYm9hcmRQZXJtaXNzaW9uUhNjbGlwYm9hcmRQZXJtaXNzaW9uEigKEGNhY2hlX2J5dGVzX3Vz'
    'ZWQYBSABKANSDmNhY2hlQnl0ZXNVc2VkEiEKDHBlZXJzX29ubGluZRgGIAEoBVILcGVlcnNPbm'
    'xpbmUSGAoHdmVyc2lvbhgHIAEoCVIHdmVyc2lvbg==');

@$core.Deprecated('Use getImagePreviewRequestDescriptor instead')
const GetImagePreviewRequest$json = {
  '1': 'GetImagePreviewRequest',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
  ],
};

/// Descriptor for `GetImagePreviewRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImagePreviewRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRJbWFnZVByZXZpZXdSZXF1ZXN0EhcKB2NsaXBfaWQYASABKAlSBmNsaXBJZA==');

@$core.Deprecated('Use getImagePreviewResponseDescriptor instead')
const GetImagePreviewResponse$json = {
  '1': 'GetImagePreviewResponse',
  '2': [
    {'1': 'path', '3': 1, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `GetImagePreviewResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImagePreviewResponseDescriptor =
    $convert.base64Decode(
        'ChdHZXRJbWFnZVByZXZpZXdSZXNwb25zZRISCgRwYXRoGAEgASgJUgRwYXRo');

// This is a generated file - do not edit.
//
// Generated from copysync/v1/wire.proto.

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

@$core.Deprecated('Use compressionDescriptor instead')
const Compression$json = {
  '1': 'Compression',
  '2': [
    {'1': 'COMPRESSION_UNSPECIFIED', '2': 0},
    {'1': 'COMPRESSION_NONE', '2': 1},
    {'1': 'COMPRESSION_ZSTD', '2': 2},
  ],
};

/// Descriptor for `Compression`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List compressionDescriptor = $convert.base64Decode(
    'CgtDb21wcmVzc2lvbhIbChdDT01QUkVTU0lPTl9VTlNQRUNJRklFRBAAEhQKEENPTVBSRVNTSU'
    '9OX05PTkUQARIUChBDT01QUkVTU0lPTl9aU1REEAI=');

@$core.Deprecated('Use packagingDescriptor instead')
const Packaging$json = {
  '1': 'Packaging',
  '2': [
    {'1': 'PACKAGING_UNSPECIFIED', '2': 0},
    {'1': 'PACKAGING_RAW', '2': 1},
    {'1': 'PACKAGING_TAR', '2': 2},
  ],
};

/// Descriptor for `Packaging`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List packagingDescriptor = $convert.base64Decode(
    'CglQYWNrYWdpbmcSGQoVUEFDS0FHSU5HX1VOU1BFQ0lGSUVEEAASEQoNUEFDS0FHSU5HX1JBVx'
    'ABEhEKDVBBQ0tBR0lOR19UQVIQAg==');

@$core.Deprecated('Use signedDescriptor instead')
const Signed$json = {
  '1': 'Signed',
  '2': [
    {'1': 'payload', '3': 1, '4': 1, '5': 12, '10': 'payload'},
    {'1': 'signature', '3': 2, '4': 1, '5': 12, '10': 'signature'},
    {'1': 'from_device_id', '3': 3, '4': 1, '5': 9, '10': 'fromDeviceId'},
  ],
};

/// Descriptor for `Signed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signedDescriptor = $convert.base64Decode(
    'CgZTaWduZWQSGAoHcGF5bG9hZBgBIAEoDFIHcGF5bG9hZBIcCglzaWduYXR1cmUYAiABKAxSCX'
    'NpZ25hdHVyZRIkCg5mcm9tX2RldmljZV9pZBgDIAEoCVIMZnJvbURldmljZUlk');

@$core.Deprecated('Use signalMessageDescriptor instead')
const SignalMessage$json = {
  '1': 'SignalMessage',
  '2': [
    {'1': 'timestamp_unix', '3': 1, '4': 1, '5': 3, '10': 'timestampUnix'},
    {'1': 'nonce', '3': 2, '4': 1, '5': 12, '10': 'nonce'},
    {'1': 'to_device_id', '3': 3, '4': 1, '5': 9, '10': 'toDeviceId'},
    {
      '1': 'offer',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.SdpOffer',
      '9': 0,
      '10': 'offer'
    },
    {
      '1': 'answer',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.SdpAnswer',
      '9': 0,
      '10': 'answer'
    },
    {
      '1': 'candidate',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.IceCandidate',
      '9': 0,
      '10': 'candidate'
    },
    {
      '1': 'ping',
      '3': 13,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.Ping',
      '9': 0,
      '10': 'ping'
    },
  ],
  '8': [
    {'1': 'payload'},
  ],
};

/// Descriptor for `SignalMessage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signalMessageDescriptor = $convert.base64Decode(
    'Cg1TaWduYWxNZXNzYWdlEiUKDnRpbWVzdGFtcF91bml4GAEgASgDUg10aW1lc3RhbXBVbml4Eh'
    'QKBW5vbmNlGAIgASgMUgVub25jZRIgCgx0b19kZXZpY2VfaWQYAyABKAlSCnRvRGV2aWNlSWQS'
    'LQoFb2ZmZXIYCiABKAsyFS5jb3B5c3luYy52MS5TZHBPZmZlckgAUgVvZmZlchIwCgZhbnN3ZX'
    'IYCyABKAsyFi5jb3B5c3luYy52MS5TZHBBbnN3ZXJIAFIGYW5zd2VyEjkKCWNhbmRpZGF0ZRgM'
    'IAEoCzIZLmNvcHlzeW5jLnYxLkljZUNhbmRpZGF0ZUgAUgljYW5kaWRhdGUSJwoEcGluZxgNIA'
    'EoCzIRLmNvcHlzeW5jLnYxLlBpbmdIAFIEcGluZ0IJCgdwYXlsb2Fk');

@$core.Deprecated('Use sdpOfferDescriptor instead')
const SdpOffer$json = {
  '1': 'SdpOffer',
  '2': [
    {'1': 'sdp', '3': 1, '4': 1, '5': 9, '10': 'sdp'},
    {'1': 'dtls_fingerprint', '3': 2, '4': 1, '5': 9, '10': 'dtlsFingerprint'},
  ],
};

/// Descriptor for `SdpOffer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sdpOfferDescriptor = $convert.base64Decode(
    'CghTZHBPZmZlchIQCgNzZHAYASABKAlSA3NkcBIpChBkdGxzX2ZpbmdlcnByaW50GAIgASgJUg'
    '9kdGxzRmluZ2VycHJpbnQ=');

@$core.Deprecated('Use sdpAnswerDescriptor instead')
const SdpAnswer$json = {
  '1': 'SdpAnswer',
  '2': [
    {'1': 'sdp', '3': 1, '4': 1, '5': 9, '10': 'sdp'},
    {'1': 'dtls_fingerprint', '3': 2, '4': 1, '5': 9, '10': 'dtlsFingerprint'},
  ],
};

/// Descriptor for `SdpAnswer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sdpAnswerDescriptor = $convert.base64Decode(
    'CglTZHBBbnN3ZXISEAoDc2RwGAEgASgJUgNzZHASKQoQZHRsc19maW5nZXJwcmludBgCIAEoCV'
    'IPZHRsc0ZpbmdlcnByaW50');

@$core.Deprecated('Use iceCandidateDescriptor instead')
const IceCandidate$json = {
  '1': 'IceCandidate',
  '2': [
    {'1': 'candidate', '3': 1, '4': 1, '5': 9, '10': 'candidate'},
    {'1': 'sdp_mid', '3': 2, '4': 1, '5': 9, '10': 'sdpMid'},
    {'1': 'sdp_mline_index', '3': 3, '4': 1, '5': 13, '10': 'sdpMlineIndex'},
  ],
};

/// Descriptor for `IceCandidate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List iceCandidateDescriptor = $convert.base64Decode(
    'CgxJY2VDYW5kaWRhdGUSHAoJY2FuZGlkYXRlGAEgASgJUgljYW5kaWRhdGUSFwoHc2RwX21pZB'
    'gCIAEoCVIGc2RwTWlkEiYKD3NkcF9tbGluZV9pbmRleBgDIAEoDVINc2RwTWxpbmVJbmRleA==');

@$core.Deprecated('Use pingDescriptor instead')
const Ping$json = {
  '1': 'Ping',
  '2': [
    {'1': 'sent_at_unix', '3': 1, '4': 1, '5': 3, '10': 'sentAtUnix'},
  ],
};

/// Descriptor for `Ping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pingDescriptor = $convert
    .base64Decode('CgRQaW5nEiAKDHNlbnRfYXRfdW5peBgBIAEoA1IKc2VudEF0VW5peA==');

@$core.Deprecated('Use clientHelloDescriptor instead')
const ClientHello$json = {
  '1': 'ClientHello',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'device_name', '3': 2, '4': 1, '5': 9, '10': 'deviceName'},
    {'1': 'platform', '3': 3, '4': 1, '5': 9, '10': 'platform'},
    {'1': 'public_key', '3': 4, '4': 1, '5': 12, '10': 'publicKey'},
    {'1': 'signature', '3': 5, '4': 1, '5': 12, '10': 'signature'},
    {'1': 'paired_device_ids', '3': 6, '4': 3, '5': 9, '10': 'pairedDeviceIds'},
  ],
};

/// Descriptor for `ClientHello`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clientHelloDescriptor = $convert.base64Decode(
    'CgtDbGllbnRIZWxsbxIbCglkZXZpY2VfaWQYASABKAlSCGRldmljZUlkEh8KC2RldmljZV9uYW'
    '1lGAIgASgJUgpkZXZpY2VOYW1lEhoKCHBsYXRmb3JtGAMgASgJUghwbGF0Zm9ybRIdCgpwdWJs'
    'aWNfa2V5GAQgASgMUglwdWJsaWNLZXkSHAoJc2lnbmF0dXJlGAUgASgMUglzaWduYXR1cmUSKg'
    'oRcGFpcmVkX2RldmljZV9pZHMYBiADKAlSD3BhaXJlZERldmljZUlkcw==');

@$core.Deprecated('Use serverHelloDescriptor instead')
const ServerHello$json = {
  '1': 'ServerHello',
  '2': [
    {'1': 'challenge', '3': 1, '4': 1, '5': 12, '10': 'challenge'},
    {
      '1': 'ice_servers',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.copysync.v1.IceServer',
      '10': 'iceServers'
    },
    {'1': 'server_time_unix', '3': 3, '4': 1, '5': 3, '10': 'serverTimeUnix'},
  ],
};

/// Descriptor for `ServerHello`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverHelloDescriptor = $convert.base64Decode(
    'CgtTZXJ2ZXJIZWxsbxIcCgljaGFsbGVuZ2UYASABKAxSCWNoYWxsZW5nZRI3CgtpY2Vfc2Vydm'
    'VycxgCIAMoCzIWLmNvcHlzeW5jLnYxLkljZVNlcnZlclIKaWNlU2VydmVycxIoChBzZXJ2ZXJf'
    'dGltZV91bml4GAMgASgDUg5zZXJ2ZXJUaW1lVW5peA==');

@$core.Deprecated('Use iceServerDescriptor instead')
const IceServer$json = {
  '1': 'IceServer',
  '2': [
    {'1': 'urls', '3': 1, '4': 3, '5': 9, '10': 'urls'},
    {'1': 'username', '3': 2, '4': 1, '5': 9, '10': 'username'},
    {'1': 'credential', '3': 3, '4': 1, '5': 9, '10': 'credential'},
    {'1': 'expires_at_unix', '3': 4, '4': 1, '5': 3, '10': 'expiresAtUnix'},
  ],
};

/// Descriptor for `IceServer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List iceServerDescriptor = $convert.base64Decode(
    'CglJY2VTZXJ2ZXISEgoEdXJscxgBIAMoCVIEdXJscxIaCgh1c2VybmFtZRgCIAEoCVIIdXNlcm'
    '5hbWUSHgoKY3JlZGVudGlhbBgDIAEoCVIKY3JlZGVudGlhbBImCg9leHBpcmVzX2F0X3VuaXgY'
    'BCABKANSDWV4cGlyZXNBdFVuaXg=');

@$core.Deprecated('Use peerPresenceDescriptor instead')
const PeerPresence$json = {
  '1': 'PeerPresence',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'online', '3': 2, '4': 1, '5': 8, '10': 'online'},
  ],
};

/// Descriptor for `PeerPresence`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerPresenceDescriptor = $convert.base64Decode(
    'CgxQZWVyUHJlc2VuY2USGwoJZGV2aWNlX2lkGAEgASgJUghkZXZpY2VJZBIWCgZvbmxpbmUYAi'
    'ABKAhSBm9ubGluZQ==');

@$core.Deprecated('Use pairingCreateRequestDescriptor instead')
const PairingCreateRequest$json = {
  '1': 'PairingCreateRequest',
  '2': [
    {'1': 'device_name', '3': 1, '4': 1, '5': 9, '10': 'deviceName'},
    {'1': 'public_key', '3': 2, '4': 1, '5': 12, '10': 'publicKey'},
  ],
};

/// Descriptor for `PairingCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pairingCreateRequestDescriptor = $convert.base64Decode(
    'ChRQYWlyaW5nQ3JlYXRlUmVxdWVzdBIfCgtkZXZpY2VfbmFtZRgBIAEoCVIKZGV2aWNlTmFtZR'
    'IdCgpwdWJsaWNfa2V5GAIgASgMUglwdWJsaWNLZXk=');

@$core.Deprecated('Use pairingCreateResponseDescriptor instead')
const PairingCreateResponse$json = {
  '1': 'PairingCreateResponse',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'expires_at_unix', '3': 2, '4': 1, '5': 3, '10': 'expiresAtUnix'},
  ],
};

/// Descriptor for `PairingCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pairingCreateResponseDescriptor = $convert.base64Decode(
    'ChVQYWlyaW5nQ3JlYXRlUmVzcG9uc2USEgoEY29kZRgBIAEoCVIEY29kZRImCg9leHBpcmVzX2'
    'F0X3VuaXgYAiABKANSDWV4cGlyZXNBdFVuaXg=');

@$core.Deprecated('Use pairingRedeemRequestDescriptor instead')
const PairingRedeemRequest$json = {
  '1': 'PairingRedeemRequest',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'device_id', '3': 2, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'device_name', '3': 3, '4': 1, '5': 9, '10': 'deviceName'},
    {'1': 'public_key', '3': 4, '4': 1, '5': 12, '10': 'publicKey'},
  ],
};

/// Descriptor for `PairingRedeemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pairingRedeemRequestDescriptor = $convert.base64Decode(
    'ChRQYWlyaW5nUmVkZWVtUmVxdWVzdBISCgRjb2RlGAEgASgJUgRjb2RlEhsKCWRldmljZV9pZB'
    'gCIAEoCVIIZGV2aWNlSWQSHwoLZGV2aWNlX25hbWUYAyABKAlSCmRldmljZU5hbWUSHQoKcHVi'
    'bGljX2tleRgEIAEoDFIJcHVibGljS2V5');

@$core.Deprecated('Use pairingRedeemResponseDescriptor instead')
const PairingRedeemResponse$json = {
  '1': 'PairingRedeemResponse',
  '2': [
    {'1': 'peer_device_id', '3': 1, '4': 1, '5': 9, '10': 'peerDeviceId'},
    {'1': 'peer_device_name', '3': 2, '4': 1, '5': 9, '10': 'peerDeviceName'},
    {'1': 'peer_platform', '3': 3, '4': 1, '5': 9, '10': 'peerPlatform'},
    {'1': 'peer_public_key', '3': 4, '4': 1, '5': 12, '10': 'peerPublicKey'},
  ],
};

/// Descriptor for `PairingRedeemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pairingRedeemResponseDescriptor = $convert.base64Decode(
    'ChVQYWlyaW5nUmVkZWVtUmVzcG9uc2USJAoOcGVlcl9kZXZpY2VfaWQYASABKAlSDHBlZXJEZX'
    'ZpY2VJZBIoChBwZWVyX2RldmljZV9uYW1lGAIgASgJUg5wZWVyRGV2aWNlTmFtZRIjCg1wZWVy'
    'X3BsYXRmb3JtGAMgASgJUgxwZWVyUGxhdGZvcm0SJgoPcGVlcl9wdWJsaWNfa2V5GAQgASgMUg'
    '1wZWVyUHVibGljS2V5');

@$core.Deprecated('Use pairingRedeemedDescriptor instead')
const PairingRedeemed$json = {
  '1': 'PairingRedeemed',
  '2': [
    {'1': 'peer_device_id', '3': 1, '4': 1, '5': 9, '10': 'peerDeviceId'},
    {'1': 'peer_device_name', '3': 2, '4': 1, '5': 9, '10': 'peerDeviceName'},
    {'1': 'peer_platform', '3': 3, '4': 1, '5': 9, '10': 'peerPlatform'},
    {'1': 'peer_public_key', '3': 4, '4': 1, '5': 12, '10': 'peerPublicKey'},
  ],
};

/// Descriptor for `PairingRedeemed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pairingRedeemedDescriptor = $convert.base64Decode(
    'Cg9QYWlyaW5nUmVkZWVtZWQSJAoOcGVlcl9kZXZpY2VfaWQYASABKAlSDHBlZXJEZXZpY2VJZB'
    'IoChBwZWVyX2RldmljZV9uYW1lGAIgASgJUg5wZWVyRGV2aWNlTmFtZRIjCg1wZWVyX3BsYXRm'
    'b3JtGAMgASgJUgxwZWVyUGxhdGZvcm0SJgoPcGVlcl9wdWJsaWNfa2V5GAQgASgMUg1wZWVyUH'
    'VibGljS2V5');

@$core.Deprecated('Use iceConfigDescriptor instead')
const IceConfig$json = {
  '1': 'IceConfig',
  '2': [
    {
      '1': 'ice_servers',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.copysync.v1.IceServer',
      '10': 'iceServers'
    },
  ],
};

/// Descriptor for `IceConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List iceConfigDescriptor = $convert.base64Decode(
    'CglJY2VDb25maWcSNwoLaWNlX3NlcnZlcnMYASADKAsyFi5jb3B5c3luYy52MS5JY2VTZXJ2ZX'
    'JSCmljZVNlcnZlcnM=');

@$core.Deprecated('Use serverEnvelopeDescriptor instead')
const ServerEnvelope$json = {
  '1': 'ServerEnvelope',
  '2': [
    {
      '1': 'hello',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.ServerHello',
      '9': 0,
      '10': 'hello'
    },
    {
      '1': 'forwarded',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.Signed',
      '9': 0,
      '10': 'forwarded'
    },
    {
      '1': 'presence',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.PeerPresence',
      '9': 0,
      '10': 'presence'
    },
    {
      '1': 'pairing_created',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.PairingCreateResponse',
      '9': 0,
      '10': 'pairingCreated'
    },
    {
      '1': 'pairing_redeemed_result',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.PairingRedeemResponse',
      '9': 0,
      '10': 'pairingRedeemedResult'
    },
    {
      '1': 'pairing_redeemed',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.PairingRedeemed',
      '9': 0,
      '10': 'pairingRedeemed'
    },
    {
      '1': 'error',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.ServerError',
      '9': 0,
      '10': 'error'
    },
    {
      '1': 'ice_config',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.IceConfig',
      '9': 0,
      '10': 'iceConfig'
    },
  ],
  '8': [
    {'1': 'payload'},
  ],
};

/// Descriptor for `ServerEnvelope`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverEnvelopeDescriptor = $convert.base64Decode(
    'Cg5TZXJ2ZXJFbnZlbG9wZRIwCgVoZWxsbxgBIAEoCzIYLmNvcHlzeW5jLnYxLlNlcnZlckhlbG'
    'xvSABSBWhlbGxvEjMKCWZvcndhcmRlZBgCIAEoCzITLmNvcHlzeW5jLnYxLlNpZ25lZEgAUglm'
    'b3J3YXJkZWQSNwoIcHJlc2VuY2UYAyABKAsyGS5jb3B5c3luYy52MS5QZWVyUHJlc2VuY2VIAF'
    'IIcHJlc2VuY2USTQoPcGFpcmluZ19jcmVhdGVkGAQgASgLMiIuY29weXN5bmMudjEuUGFpcmlu'
    'Z0NyZWF0ZVJlc3BvbnNlSABSDnBhaXJpbmdDcmVhdGVkElwKF3BhaXJpbmdfcmVkZWVtZWRfcm'
    'VzdWx0GAUgASgLMiIuY29weXN5bmMudjEuUGFpcmluZ1JlZGVlbVJlc3BvbnNlSABSFXBhaXJp'
    'bmdSZWRlZW1lZFJlc3VsdBJJChBwYWlyaW5nX3JlZGVlbWVkGAYgASgLMhwuY29weXN5bmMudj'
    'EuUGFpcmluZ1JlZGVlbWVkSABSD3BhaXJpbmdSZWRlZW1lZBIwCgVlcnJvchgHIAEoCzIYLmNv'
    'cHlzeW5jLnYxLlNlcnZlckVycm9ySABSBWVycm9yEjcKCmljZV9jb25maWcYCCABKAsyFi5jb3'
    'B5c3luYy52MS5JY2VDb25maWdIAFIJaWNlQ29uZmlnQgkKB3BheWxvYWQ=');

@$core.Deprecated('Use clientEnvelopeDescriptor instead')
const ClientEnvelope$json = {
  '1': 'ClientEnvelope',
  '2': [
    {
      '1': 'hello',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.ClientHello',
      '9': 0,
      '10': 'hello'
    },
    {
      '1': 'forward',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.Signed',
      '9': 0,
      '10': 'forward'
    },
    {
      '1': 'pairing_create',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.PairingCreateRequest',
      '9': 0,
      '10': 'pairingCreate'
    },
    {
      '1': 'pairing_redeem',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.PairingRedeemRequest',
      '9': 0,
      '10': 'pairingRedeem'
    },
  ],
  '8': [
    {'1': 'payload'},
  ],
};

/// Descriptor for `ClientEnvelope`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clientEnvelopeDescriptor = $convert.base64Decode(
    'Cg5DbGllbnRFbnZlbG9wZRIwCgVoZWxsbxgBIAEoCzIYLmNvcHlzeW5jLnYxLkNsaWVudEhlbG'
    'xvSABSBWhlbGxvEi8KB2ZvcndhcmQYAiABKAsyEy5jb3B5c3luYy52MS5TaWduZWRIAFIHZm9y'
    'd2FyZBJKCg5wYWlyaW5nX2NyZWF0ZRgDIAEoCzIhLmNvcHlzeW5jLnYxLlBhaXJpbmdDcmVhdG'
    'VSZXF1ZXN0SABSDXBhaXJpbmdDcmVhdGUSSgoOcGFpcmluZ19yZWRlZW0YBCABKAsyIS5jb3B5'
    'c3luYy52MS5QYWlyaW5nUmVkZWVtUmVxdWVzdEgAUg1wYWlyaW5nUmVkZWVtQgkKB3BheWxvYW'
    'Q=');

@$core.Deprecated('Use serverErrorDescriptor instead')
const ServerError$json = {
  '1': 'ServerError',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServerError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverErrorDescriptor = $convert.base64Decode(
    'CgtTZXJ2ZXJFcnJvchISCgRjb2RlGAEgASgJUgRjb2RlEhgKB21lc3NhZ2UYAiABKAlSB21lc3'
    'NhZ2U=');

@$core.Deprecated('Use peerMessageDescriptor instead')
const PeerMessage$json = {
  '1': 'PeerMessage',
  '2': [
    {
      '1': 'offer',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.ClipOffer',
      '9': 0,
      '10': 'offer'
    },
    {
      '1': 'fetch',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.PeerFetch',
      '9': 0,
      '10': 'fetch'
    },
    {
      '1': 'reject',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.FetchReject',
      '9': 0,
      '10': 'reject'
    },
    {
      '1': 'header',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.TransferHeader',
      '9': 0,
      '10': 'header'
    },
    {
      '1': 'done',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.TransferDone',
      '9': 0,
      '10': 'done'
    },
    {
      '1': 'link',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.copysync.v1.LinkInfo',
      '9': 0,
      '10': 'link'
    },
  ],
  '8': [
    {'1': 'payload'},
  ],
};

/// Descriptor for `PeerMessage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerMessageDescriptor = $convert.base64Decode(
    'CgtQZWVyTWVzc2FnZRIuCgVvZmZlchgBIAEoCzIWLmNvcHlzeW5jLnYxLkNsaXBPZmZlckgAUg'
    'VvZmZlchIuCgVmZXRjaBgCIAEoCzIWLmNvcHlzeW5jLnYxLlBlZXJGZXRjaEgAUgVmZXRjaBIy'
    'CgZyZWplY3QYAyABKAsyGC5jb3B5c3luYy52MS5GZXRjaFJlamVjdEgAUgZyZWplY3QSNQoGaG'
    'VhZGVyGAQgASgLMhsuY29weXN5bmMudjEuVHJhbnNmZXJIZWFkZXJIAFIGaGVhZGVyEi8KBGRv'
    'bmUYBSABKAsyGS5jb3B5c3luYy52MS5UcmFuc2ZlckRvbmVIAFIEZG9uZRIrCgRsaW5rGAYgAS'
    'gLMhUuY29weXN5bmMudjEuTGlua0luZm9IAFIEbGlua0IJCgdwYXlsb2Fk');

@$core.Deprecated('Use linkInfoDescriptor instead')
const LinkInfo$json = {
  '1': 'LinkInfo',
  '2': [
    {'1': 'relayed', '3': 1, '4': 1, '5': 8, '10': 'relayed'},
  ],
};

/// Descriptor for `LinkInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List linkInfoDescriptor =
    $convert.base64Decode('CghMaW5rSW5mbxIYCgdyZWxheWVkGAEgASgIUgdyZWxheWVk');

@$core.Deprecated('Use clipOfferDescriptor instead')
const ClipOffer$json = {
  '1': 'ClipOffer',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {
      '1': 'kind',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.copysync.v1.ClipKind',
      '10': 'kind'
    },
    {
      '1': 'items',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.copysync.v1.ClipItem',
      '10': 'items'
    },
    {'1': 'total_size', '3': 4, '4': 1, '5': 3, '10': 'totalSize'},
    {'1': 'text_content', '3': 5, '4': 1, '5': 9, '10': 'textContent'},
    {'1': 'created_at_unix', '3': 6, '4': 1, '5': 3, '10': 'createdAtUnix'},
    {'1': 'will_push', '3': 7, '4': 1, '5': 8, '10': 'willPush'},
    {'1': 'plain_text', '3': 8, '4': 1, '5': 9, '10': 'plainText'},
    {'1': 'backlog', '3': 9, '4': 1, '5': 8, '10': 'backlog'},
  ],
};

/// Descriptor for `ClipOffer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clipOfferDescriptor = $convert.base64Decode(
    'CglDbGlwT2ZmZXISFwoHY2xpcF9pZBgBIAEoCVIGY2xpcElkEikKBGtpbmQYAiABKA4yFS5jb3'
    'B5c3luYy52MS5DbGlwS2luZFIEa2luZBIrCgVpdGVtcxgDIAMoCzIVLmNvcHlzeW5jLnYxLkNs'
    'aXBJdGVtUgVpdGVtcxIdCgp0b3RhbF9zaXplGAQgASgDUgl0b3RhbFNpemUSIQoMdGV4dF9jb2'
    '50ZW50GAUgASgJUgt0ZXh0Q29udGVudBImCg9jcmVhdGVkX2F0X3VuaXgYBiABKANSDWNyZWF0'
    'ZWRBdFVuaXgSGwoJd2lsbF9wdXNoGAcgASgIUgh3aWxsUHVzaBIdCgpwbGFpbl90ZXh0GAggAS'
    'gJUglwbGFpblRleHQSGAoHYmFja2xvZxgJIAEoCFIHYmFja2xvZw==');

@$core.Deprecated('Use peerFetchDescriptor instead')
const PeerFetch$json = {
  '1': 'PeerFetch',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
  ],
};

/// Descriptor for `PeerFetch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerFetchDescriptor =
    $convert.base64Decode('CglQZWVyRmV0Y2gSFwoHY2xpcF9pZBgBIAEoCVIGY2xpcElk');

@$core.Deprecated('Use fetchRejectDescriptor instead')
const FetchReject$json = {
  '1': 'FetchReject',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `FetchReject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fetchRejectDescriptor = $convert.base64Decode(
    'CgtGZXRjaFJlamVjdBIXCgdjbGlwX2lkGAEgASgJUgZjbGlwSWQSFgoGcmVhc29uGAIgASgJUg'
    'ZyZWFzb24=');

@$core.Deprecated('Use transferHeaderDescriptor instead')
const TransferHeader$json = {
  '1': 'TransferHeader',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {'1': 'stream_id', '3': 2, '4': 1, '5': 9, '10': 'streamId'},
    {'1': 'total_bytes', '3': 3, '4': 1, '5': 3, '10': 'totalBytes'},
    {
      '1': 'compression',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.copysync.v1.Compression',
      '10': 'compression'
    },
    {
      '1': 'packaging',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.copysync.v1.Packaging',
      '10': 'packaging'
    },
  ],
};

/// Descriptor for `TransferHeader`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transferHeaderDescriptor = $convert.base64Decode(
    'Cg5UcmFuc2ZlckhlYWRlchIXCgdjbGlwX2lkGAEgASgJUgZjbGlwSWQSGwoJc3RyZWFtX2lkGA'
    'IgASgJUghzdHJlYW1JZBIfCgt0b3RhbF9ieXRlcxgDIAEoA1IKdG90YWxCeXRlcxI6Cgtjb21w'
    'cmVzc2lvbhgEIAEoDjIYLmNvcHlzeW5jLnYxLkNvbXByZXNzaW9uUgtjb21wcmVzc2lvbhI0Cg'
    'lwYWNrYWdpbmcYBSABKA4yFi5jb3B5c3luYy52MS5QYWNrYWdpbmdSCXBhY2thZ2luZw==');

@$core.Deprecated('Use transferDoneDescriptor instead')
const TransferDone$json = {
  '1': 'TransferDone',
  '2': [
    {'1': 'clip_id', '3': 1, '4': 1, '5': 9, '10': 'clipId'},
    {'1': 'stream_id', '3': 2, '4': 1, '5': 9, '10': 'streamId'},
    {'1': 'ok', '3': 3, '4': 1, '5': 8, '10': 'ok'},
    {'1': 'error', '3': 4, '4': 1, '5': 9, '10': 'error'},
    {'1': 'sha256', '3': 5, '4': 1, '5': 12, '10': 'sha256'},
  ],
};

/// Descriptor for `TransferDone`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transferDoneDescriptor = $convert.base64Decode(
    'CgxUcmFuc2ZlckRvbmUSFwoHY2xpcF9pZBgBIAEoCVIGY2xpcElkEhsKCXN0cmVhbV9pZBgCIA'
    'EoCVIIc3RyZWFtSWQSDgoCb2sYAyABKAhSAm9rEhQKBWVycm9yGAQgASgJUgVlcnJvchIWCgZz'
    'aGEyNTYYBSABKAxSBnNoYTI1Ng==');

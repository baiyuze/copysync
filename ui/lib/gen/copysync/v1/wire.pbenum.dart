// This is a generated file - do not edit.
//
// Generated from copysync/v1/wire.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class Compression extends $pb.ProtobufEnum {
  static const Compression COMPRESSION_UNSPECIFIED =
      Compression._(0, _omitEnumNames ? '' : 'COMPRESSION_UNSPECIFIED');
  static const Compression COMPRESSION_NONE =
      Compression._(1, _omitEnumNames ? '' : 'COMPRESSION_NONE');
  static const Compression COMPRESSION_ZSTD =
      Compression._(2, _omitEnumNames ? '' : 'COMPRESSION_ZSTD');

  static const $core.List<Compression> values = <Compression>[
    COMPRESSION_UNSPECIFIED,
    COMPRESSION_NONE,
    COMPRESSION_ZSTD,
  ];

  static final $core.List<Compression?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static Compression? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Compression._(super.value, super.name);
}

class Packaging extends $pb.ProtobufEnum {
  static const Packaging PACKAGING_UNSPECIFIED =
      Packaging._(0, _omitEnumNames ? '' : 'PACKAGING_UNSPECIFIED');
  static const Packaging PACKAGING_RAW =
      Packaging._(1, _omitEnumNames ? '' : 'PACKAGING_RAW');
  static const Packaging PACKAGING_TAR =
      Packaging._(2, _omitEnumNames ? '' : 'PACKAGING_TAR');

  static const $core.List<Packaging> values = <Packaging>[
    PACKAGING_UNSPECIFIED,
    PACKAGING_RAW,
    PACKAGING_TAR,
  ];

  static final $core.List<Packaging?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static Packaging? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Packaging._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');

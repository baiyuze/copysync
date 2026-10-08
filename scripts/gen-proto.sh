#!/usr/bin/env bash
# 由 .proto 生成 Go 与 Dart 代码。
#
# 用 buf 而非 protoc：这台开发机访问不到 GitHub releases（下载预编译 protoc 超时），
# 而 Homebrew 在 Intel Mac 上没有 protobuf 预编译包、需从源码编译十余分钟。
# buf 是纯 Go 实现，可经 goproxy.cn 用 go install 装上，自带 protobuf 编译器。
#
# 首次准备环境：
#   go install github.com/bufbuild/buf/cmd/buf@latest
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#   dart pub global activate protoc_plugin

set -euo pipefail
cd "$(dirname "$0")/.."

export PATH="$PATH:$(go env GOPATH)/bin:$HOME/.pub-cache/bin"

for tool in buf protoc-gen-go protoc-gen-go-grpc protoc-gen-dart; do
    command -v "$tool" >/dev/null || {
        echo "缺少 $tool —— 见本脚本顶部的环境准备步骤" >&2
        exit 1
    }
done

echo "▶ lint"
buf lint proto

echo "▶ 生成"
mkdir -p ui/lib/gen
buf generate proto

echo "✓ 完成"
echo "  Go   : $(find proto/gen -name '*.go' | wc -l | tr -d ' ') 个文件"
echo "  Dart : $(find ui/lib/gen -name '*.dart' | wc -l | tr -d ' ') 个文件"

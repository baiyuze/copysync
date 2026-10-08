#!/usr/bin/env bash
# 在 macOS 上运行 NAT 打洞实验室：编译 Linux 版的 natlab，在特权容器里执行 run.sh。
#
#   ./tools/natlab/docker.sh
#
# 国内网络拉取 Alpine 软件包慢的话：APK_MIRROR=mirrors.aliyun.com ./tools/natlab/docker.sh
# 只跑部分场景：NATLAB_ONLY=晚到 ./tools/natlab/docker.sh

set -euo pipefail
cd "$(dirname "$0")/../.."

OUT=$(mktemp -d)
trap 'rm -rf "$OUT"' EXIT
ARCH=$(docker version --format '{{.Server.Arch}}')

echo "▶ 编译 natlab（linux/$ARCH）"
(cd client-core && CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH go build -o "$OUT/natlab" ./cmd/natlab)

echo "▶ 准备实验室镜像"
docker build -q -t copysync-natlab --build-arg APK_MIRROR="${APK_MIRROR:-}" tools/natlab >/dev/null

docker run --rm --privileged -e NATLAB_ONLY -v "$OUT:/lab:ro" -v "$PWD/tools/natlab:/scripts:ro" \
    copysync-natlab bash /scripts/run.sh /lab/natlab

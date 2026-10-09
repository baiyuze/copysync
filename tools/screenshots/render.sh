#!/usr/bin/env bash
# 重新生成 README 与网站用的界面截图：中文界面输出到 docs/assets/screenshots/，
# 英文、日文界面分别输出到 docs/assets/screenshots/en/ 与 ja/。
#
#   ./tools/screenshots/render.sh
#
# 截图用演示数据渲染，不连接后台服务，不会出现本机真实的剪贴板内容。
# 需要 macOS、Flutter，以及 fonttools（pip install fonttools）。

set -euo pipefail
cd "$(dirname "$0")/../.."

FONTS=$(mktemp -d)
trap 'rm -rf "$FONTS"' EXIT

echo "▶ 准备字体"
python3 tools/screenshots/prepare_fonts.py "$FONTS"

for lang in zh en ja; do
  out=docs/assets/screenshots
  [ "$lang" = zh ] || out="$out/$lang"

  echo "▶ 渲染界面（$lang）"
  rm -rf ui/tool/screenshots/out
  (cd ui && SCREENSHOT_FONTS="$FONTS" SCREENSHOT_LANG="$lang" flutter test tool/screenshots/screenshots_test.dart 2>&1 | tail -1)

  echo "▶ 加窗口外框（$lang）"
  mkdir -p "$out"
  swift tools/screenshots/frame.swift ui/tool/screenshots/out "$out"
done

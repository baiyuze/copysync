#!/usr/bin/env bash
# 把两个 spike 编译成 .app bundle。
#
# 必须是 bundle 而不是裸二进制：macOS 的 TCC 隐私机制以 bundle identifier
# 为授权主体，裸命令行工具的提示行为与真实 App 不一致，测出来的结论不可用。

set -euo pipefail
cd "$(dirname "$0")"

OUT=build
rm -rf "$OUT" && mkdir -p "$OUT"

build_app() {
    local name=$1 src=$2 bundle_id=$3
    local app="$OUT/$name.app"

    mkdir -p "$app/Contents/MacOS"

    cat > "$app/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>          <string>$name</string>
    <key>CFBundleIdentifier</key>          <string>$bundle_id</string>
    <key>CFBundleName</key>                <string>$name</string>
    <key>CFBundlePackageType</key>         <string>APPL</string>
    <key>CFBundleShortVersionString</key>  <string>0.1</string>
    <key>CFBundleVersion</key>             <string>1</string>
    <key>LSMinimumSystemVersion</key>      <string>13.0</string>
    <key>LSUIElement</key>                 <true/>
</dict>
</plist>
PLIST

    clang -fobjc-arc -O0 -g -Wall \
        -framework Cocoa \
        -framework UniformTypeIdentifiers \
        -o "$app/Contents/MacOS/$name" "$src"

    # ad-hoc 签名：TCC 需要一个可识别的签名身份才能记住授权决定
    codesign --force --sign - "$app"

    echo "  ✓ $app"
}

echo "编译中…"
build_app SpikeA privacy/main.m com.copysync.spike.privacy
build_app SpikeB promise/main.m  com.copysync.spike.promise
build_app SpikeC writer/main.m   com.copysync.spike.writer
build_app SpikeD receiver/main.m com.copysync.spike.receiver

cat <<EOF

完成。

  ./$OUT/SpikeA.app/Contents/MacOS/SpikeA [轮数]   # 隐私机制探测（读方）
  ./$OUT/SpikeB.app/Contents/MacOS/SpikeB          # promise 可用性
  ./$OUT/SpikeC.app/Contents/MacOS/SpikeC <mode>   # 复制方模拟（file|text|image）

完整的隐私对照实验：./test-privacy.sh
EOF

#!/usr/bin/env bash
# 构建全部产物到 dist/。
#
#   dist/CopySync.app         发给用户的 App：界面 + 内置的后台服务与命令行工具
#   dist/CopySyncDaemon.app   后台服务单独一份，供 install-macos.sh 在开发时安装
#   dist/copysync-server      信令 + TURN 中转服务器（macOS 版）
#   dist/copysync-cli         调试用命令行
#
# daemon 刻意打包成 .app 而非裸二进制：macOS 的剪贴板授权（TCC）以
# bundle identifier 为主体，下个大版本起剪贴板隐私机制将强制执行，
# 届时没有稳定 bundle id 的进程无法保留用户的授权选择。

set -euo pipefail
cd "$(dirname "$0")/.."

DIST=dist
VERSION=${VERSION:-1.3.3}
rm -rf "$DIST" && mkdir -p "$DIST"

# 构建 x86_64 + arm64 通用二进制：两台 Mac 的芯片未必相同，单架构的
# daemon 到 Apple Silicon 上要靠 Rosetta，而 Rosetta 默认并未安装。
# daemon 含 cgo，交叉编译时必须给 clang 显式指定 -arch。
build_universal() { # <模块目录> <包路径> <输出路径> <ldflags>
    local dir=$1 pkg=$2 out=$3 ldflags=$4 tmp
    tmp=$(mktemp -d)
    for arch in amd64 arm64; do
        local clang_arch=$arch
        [[ $arch == amd64 ]] && clang_arch=x86_64
        (cd "$dir" && CGO_ENABLED=1 GOOS=darwin GOARCH=$arch CC="clang -arch $clang_arch" \
            go build -trimpath -ldflags "$ldflags" -o "$tmp/$arch" "$pkg")
    done
    lipo -create -output "$out" "$tmp/amd64" "$tmp/arm64"
    rm -rf "$tmp"
}

echo "▶ 构建 daemon"
build_universal client-core ./cmd/copysyncd "$DIST/copysyncd" \
    "-s -w -X main.version=$VERSION"

echo "▶ 构建 CLI"
build_universal client-core ./cmd/copysync-cli "$DIST/copysync-cli" "-s -w"

echo "▶ 构建服务器"
build_universal server ./cmd/copysync-server "$DIST/copysync-server" \
    "-s -w -X main.version=$VERSION"

# ── 把 daemon 包成 .app ──
APP="$DIST/CopySyncDaemon.app"
mkdir -p "$APP/Contents/MacOS"
mv "$DIST/copysyncd" "$APP/Contents/MacOS/copysyncd"
cat > "$APP/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>         <string>copysyncd</string>
    <key>CFBundleIdentifier</key>         <string>com.copysync.daemon</string>
    <key>CFBundleName</key>               <string>CopySync Daemon</string>
    <key>CFBundleIconFile</key>           <string>AppIcon</string>
    <key>CFBundlePackageType</key>        <string>APPL</string>
    <key>CFBundleShortVersionString</key> <string>$VERSION</string>
    <key>CFBundleVersion</key>            <string>$VERSION</string>
    <key>LSMinimumSystemVersion</key>     <string>13.0</string>
    <!-- 后台常驻：不出现在 Dock 与任务切换器里 -->
    <key>LSUIElement</key>                <true/>
</dict>
</plist>
PLIST
# 系统设置的「隐私与安全性」列表里会显示后台服务，给它同一个图标
mkdir -p "$APP/Contents/Resources"
cp assets/brand/AppIcon.icns "$APP/Contents/Resources/AppIcon.icns"
# ad-hoc 签名即可让 TCC 记住授权决定；正式分发应换成 Developer ID
codesign --force --sign - "$APP"
echo "  ✓ $APP"

echo "▶ 构建图形界面"
(cd ui && flutter build macos --release \
    --build-name "$VERSION" --dart-define=APP_VERSION="$VERSION" 2>&1 | tail -1)
UI="$DIST/CopySync.app"
cp -R ui/build/macos/Build/Products/Release/CopySync.app "$UI"

# 后台服务与命令行工具内置进 App：用户只需把一个 App 拖进「应用程序」，
# 首次打开时由界面把后台服务注册为登录项（见 ui/lib/background_service.dart）
mkdir -p "$UI/Contents/Helpers"
cp -R "$APP" "$UI/Contents/Helpers/"
cp "$DIST/copysync-cli" "$UI/Contents/Helpers/copysync-cli"

# 签名由内向外：先签内置的组件，最后签外层 App。
# Flutter 增量构建在 Dart 代码变化后会重新生成 App.framework 却不重签，
# 留下失效签名，Apple Silicon 上会直接拒绝运行，这里一并重签。
codesign --force --sign - "$UI/Contents/Helpers/copysync-cli"
codesign --force --sign - "$UI/Contents/Helpers/CopySyncDaemon.app"
codesign --force --sign - "$UI/Contents/Frameworks/App.framework"
codesign --force --sign - --preserve-metadata=entitlements,requirements,flags "$UI"
codesign --verify --deep --strict "$UI" "$APP"
echo "  ✓ $UI"

echo
echo "完成。产物："
ls -1 "$DIST"
cat <<'EOF'

下一步：
  open dist/CopySync.app              首次打开会引导启用后台同步
  ./scripts/package.sh                打出 DMG 与服务器发布包
  ./dist/copysync-server -h           查看服务器参数
EOF

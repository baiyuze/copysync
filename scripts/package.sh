#!/usr/bin/env bash
# 构建并打出可分发的发布包到 release/。
#
#   CopySync.dmg                          macOS 客户端，Intel 与 Apple 芯片通用
#   copysync-server-linux-amd64.tar.gz    信令服务器 + systemd / Docker 部署文件
#   copysync-server-linux-arm64.tar.gz
#   copysync-server-macos.tar.gz          把服务器放在其中一台 Mac 上时用
#
# Windows 的安装程序与便携版（CopySync-Setup.exe、CopySync-windows-x64.zip）只能在 Windows 上
# 构建：由 CI 的 Windows 任务运行 scripts/build-windows.ps1 产出，发布时从构建产物里取。
#
# 文件名不带版本号：网站与 README 用 releases/latest/download/<文件名> 直链下载，
# 带版本号的话每发一版这些链接都会失效。版本号见 Release 标题与 copysync-server -version。
#   SHA256SUMS
#
# 服务器是纯 Go，直接交叉编译；客户端依赖 cgo 与 Flutter，只能在 Mac 上构建。

set -euo pipefail
cd "$(dirname "$0")/.."

export VERSION=${VERSION:-1.3.0}
OUT=release

./scripts/build.sh
echo

rm -rf "$OUT" && mkdir -p "$OUT"
STAGE=$(mktemp -d)
trap 'rm -rf "$STAGE"' EXIT

# ─────────────────────────── DMG ───────────────────────────

# Finder 窗口的版式与 tools/dmg/background.swift 画的背景一一对应。
# 需要「自动化」权限去控制 Finder；拿不到就跳过，DMG 照样能装，只是图标排布是默认的。
layout_dmg() {
    osascript >/dev/null <<EOF &
tell application "Finder"
    tell disk "$1"
        open
        set current view of container window to icon view
        set toolbar visible of container window to false
        set statusbar visible of container window to false
        set bounds of container window to {200, 120, 840, 548}
        set opts to icon view options of container window
        set arrangement of opts to not arranged
        set icon size of opts to 112
        set text size of opts to 13
        set background picture of opts to file ".background:background.tiff"
        set position of item "CopySync.app" of container window to {160, 190}
        set position of item "Applications" of container window to {480, 190}
        update without registering applications
        delay 1
        close
    end tell
end tell
EOF
    local pid=$!
    for _ in {1..60}; do
        kill -0 "$pid" 2>/dev/null || { wait "$pid"; return $?; }
        sleep 1
    done
    kill "$pid" 2>/dev/null
    return 1
}

echo "▶ 打包 DMG"
VOL=CopySync
SRC="$STAGE/dmg"
mkdir -p "$SRC/.background"
ditto dist/CopySync.app "$SRC/CopySync.app"
ln -s /Applications "$SRC/Applications"
swift tools/dmg/background.swift "$STAGE/bg" >/dev/null
tiffutil -cathidpicheck "$STAGE/bg/background.png" "$STAGE/bg/background@2x.png" \
    -out "$SRC/.background/background.tiff" 2>/dev/null

# 同名的卷还挂着的话，Finder 会把布局写到另一个窗口上
hdiutil detach "/Volumes/$VOL" -quiet 2>/dev/null || true
hdiutil create -volname "$VOL" -srcfolder "$SRC" -fs HFS+ -format UDRW -ov "$STAGE/rw.dmg" -quiet
MNT=$(hdiutil attach "$STAGE/rw.dmg" -readwrite -noverify -noautoopen | awk -F'\t' '/\/Volumes\//{print $NF}')
if layout_dmg "$VOL"; then
    echo "  ✓ 安装窗口布局"
else
    echo "  ! 没能设置安装窗口布局（Finder 自动化权限未授予），DMG 仍可正常安装"
fi
# 卷图标放在布局之后：-srcfolder 不会带上 .VolumeIcon.icns，
# 而 Finder 保存布局时又会清掉卷的自定义图标标记
cp assets/brand/AppIcon.icns "$MNT/.VolumeIcon.icns"
SetFile -a C "$MNT"
rm -rf "$MNT/.fseventsd"
sync
hdiutil detach "$MNT" -quiet
hdiutil convert "$STAGE/rw.dmg" -format UDZO -imagekey zlib-level=9 \
    -o "$OUT/CopySync.dmg" -quiet
echo "  ✓ CopySync.dmg（$VERSION）"

# ─────────────────────────── 服务器 ───────────────────────────

echo "▶ 打包服务器"
for arch in amd64 arm64; do
    name="copysync-server-linux-$arch"
    d="$STAGE/$name"
    mkdir -p "$d"
    (cd server && CGO_ENABLED=0 GOOS=linux GOARCH=$arch go build -trimpath \
        -ldflags "-s -w -X main.version=$VERSION" \
        -o "$d/copysync-server" ./cmd/copysync-server)
    cp server/deploy/{install.sh,copysync-server.service,Dockerfile,docker-compose.yml,env.example,README.md} "$d/"
    tar -C "$STAGE" -czf "$OUT/$name.tar.gz" "$name"
    echo "  ✓ $name.tar.gz"
done

name="copysync-server-macos"
mkdir -p "$STAGE/$name"
cp dist/copysync-server server/deploy/README.md "$STAGE/$name/"
tar -C "$STAGE" -czf "$OUT/$name.tar.gz" "$name"
echo "  ✓ $name.tar.gz"

(cd "$OUT" && shasum -a 256 -- * > SHA256SUMS)

echo
echo "完成。发布包："
ls -lh "$OUT" | tail -n +2

#!/usr/bin/env bash
# 在 Linux 上构建并打包 PasteFerry 客户端：界面（Flutter，GTK）、后台服务 pasteferryd、命令行 pasteferry-cli。
#
#   ./scripts/package-linux.sh
#
# 产出放在 dist/linux/：
#   PasteFerry-linux-<arch>.deb     Debian、Ubuntu 及其衍生版，装在 /opt/PasteFerry，自动出现在应用菜单里
#   PasteFerry-linux-<arch>.tar.gz  其他发行版：解压即用，运行里面的 install.sh 加进应用菜单
#
# <arch> 是 amd64 或 arm64，取决于运行脚本的机器：Flutter 不能交叉编译 Linux 桌面程序，
# 所以 CI 在两种架构的机器上各跑一次。
#
# 需要 Go、Flutter，以及 Flutter 的 Linux 构建依赖：
#   sudo apt install clang cmake ninja-build pkg-config libgtk-3-dev
# 文件名不带版本号，网站上的 releases/latest/download/<文件名> 链接才不会随发版失效。
set -euo pipefail
cd "$(dirname "$0")/.."

VERSION=${VERSION:-$(sed -n 's/^VERSION=\${VERSION:-\(.*\)}$/\1/p' scripts/build.sh)}
case "$(uname -m)" in
    x86_64) ARCH=amd64 FLUTTER_ARCH=x64 ;;
    aarch64 | arm64) ARCH=arm64 FLUTTER_ARCH=arm64 ;;
    *)
        echo "不支持的架构：$(uname -m)" >&2
        exit 1
        ;;
esac

OUT=dist/linux
APP=$OUT/PasteFerry
rm -rf "$OUT" && mkdir -p "$APP"

echo "▶ 后台服务与命令行（$VERSION，$ARCH）"
# 纯 Go，不依赖 glibc 版本：剪贴板经 X11 协议直接通信，不链接 libX11
(cd client-core && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$VERSION" \
    -o "../$APP/pasteferryd" ./cmd/pasteferryd)
(cd client-core && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o "../$APP/pasteferry-cli" ./cmd/pasteferry-cli)

echo "▶ 界面"
(cd ui && flutter build linux --release --build-name "$VERSION" --dart-define=APP_VERSION="$VERSION")
cp -R "ui/build/linux/$FLUTTER_ARCH/release/bundle/." "$APP/"

# ── tar.gz：解压即用 ──
echo "▶ PasteFerry-linux-$ARCH.tar.gz"
TAR=$OUT/tar/PasteFerry
mkdir -p "$TAR/share"
cp -R "$APP/." "$TAR/"
cp tools/linux/install.sh tools/linux/uninstall.sh "$TAR/"
cp tools/linux/com.pasteferry.app.desktop tools/linux/icons/*.png "$TAR/share/"
tar -C "$OUT/tar" -czf "$OUT/PasteFerry-linux-$ARCH.tar.gz" --owner=0 --group=0 PasteFerry

# ── deb：装在 /opt/PasteFerry ──
echo "▶ PasteFerry-linux-$ARCH.deb"
DEB=$OUT/deb
mkdir -p "$DEB/DEBIAN" "$DEB/opt" "$DEB/usr/bin" "$DEB/usr/share/applications"
cp -R "$APP" "$DEB/opt/PasteFerry"
sed 's#@EXEC@#/opt/PasteFerry/pasteferry#' tools/linux/com.pasteferry.app.desktop \
    > "$DEB/usr/share/applications/com.pasteferry.app.desktop"
for size in 256 512; do
    mkdir -p "$DEB/usr/share/icons/hicolor/${size}x$size/apps"
    cp "tools/linux/icons/com.pasteferry.app-$size.png" "$DEB/usr/share/icons/hicolor/${size}x$size/apps/com.pasteferry.app.png"
done
ln -s /opt/PasteFerry/pasteferry "$DEB/usr/bin/pasteferry"
ln -s /opt/PasteFerry/pasteferry-cli "$DEB/usr/bin/pasteferry-cli"

# libgtk-3-0：Ubuntu 24.04 起改名为 libgtk-3-0t64，新包也提供旧名字，写旧名字两边都能装
cat > "$DEB/DEBIAN/control" <<EOF
Package: pasteferry
Version: $VERSION
Architecture: $ARCH
Maintainer: baiyuze <27543832+baiyuze@users.noreply.github.com>
Installed-Size: $(du -sk "$DEB" | cut -f1)
Depends: libgtk-3-0 (>= 3.22)
Recommends: fonts-noto-cjk
Conflicts: copysync
Replaces: copysync
Section: utils
Priority: optional
Homepage: https://baiyuze.github.io/pasteferry/
Description: Sync what you copy between your computers
 PasteFerry keeps the clipboard in sync across Mac, Windows and Linux:
 text, images and files, directly over the local network or through
 your own server. Content is end-to-end encrypted.
EOF
# 2.0.0 之前的包叫 copysync：上面的 Conflicts/Replaces 让 apt 装这个包时自动卸掉旧的，
# 旧的后台服务与开机自启项由界面第一次启动时清理（ui/lib/background_service_linux.dart 的 removeLegacy）
#
# 卸载时结束正在运行的后台服务：程序文件已经没了。升级时不结束，界面打开时会按新版本重启它
cat > "$DEB/DEBIAN/prerm" <<'EOF'
#!/bin/sh
set -e
if [ "$1" = remove ]; then
    pkill -x pasteferryd || true
fi
EOF
chmod 755 "$DEB/DEBIAN/prerm"
dpkg-deb --build --root-owner-group "$DEB" "$OUT/PasteFerry-linux-$ARCH.deb" >/dev/null

rm -rf "$OUT/tar" "$DEB"
echo
echo "完成："
ls -lh "$OUT"/PasteFerry-linux-*

#!/usr/bin/env bash
# 把解压出来的 CopySync 加进当前用户的应用菜单（不需要 root）。
#
#   ./install.sh
#
# 程序留在原地运行：之后挪了这个文件夹，重新运行一次即可。卸载运行 ./uninstall.sh。
set -euo pipefail
here=$(cd "$(dirname "$0")" && pwd)
data=${XDG_DATA_HOME:-$HOME/.local/share}

mkdir -p "$data/applications" "$data/icons/hicolor/256x256/apps" "$data/icons/hicolor/512x512/apps"
# Exec 里的路径按 Desktop Entry 规范放在引号里，允许文件夹名带空格
sed "s#@EXEC@#\"$here/copysync\"#" "$here/share/com.copysync.app.desktop" > "$data/applications/com.copysync.app.desktop"
cp "$here/share/com.copysync.app-256.png" "$data/icons/hicolor/256x256/apps/com.copysync.app.png"
cp "$here/share/com.copysync.app-512.png" "$data/icons/hicolor/512x512/apps/com.copysync.app.png"
command -v update-desktop-database >/dev/null && update-desktop-database -q "$data/applications" || true
command -v gtk-update-icon-cache >/dev/null && gtk-update-icon-cache -q -t "$data/icons/hicolor" || true

echo "已把 CopySync 加进应用菜单。也可以直接运行：$here/copysync"

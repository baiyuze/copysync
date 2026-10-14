#!/usr/bin/env bash
# 撤销 install.sh：从应用菜单移除 PasteFerry，停用后台同步与开机自启。
# 历史记录与配对关系保存在 ~/.local/share/pasteferry，不会删除；不再需要时手动删掉这个文件夹。
set -euo pipefail
data=${XDG_DATA_HOME:-$HOME/.local/share}
config=${XDG_CONFIG_HOME:-$HOME/.config}

pkill -u "$(id -u)" -x pasteferryd 2>/dev/null || true
rm -f "$config/autostart/pasteferry-daemon.desktop" \
      "$data/applications/com.pasteferry.app.desktop" \
      "$data/icons/hicolor/256x256/apps/com.pasteferry.app.png" \
      "$data/icons/hicolor/512x512/apps/com.pasteferry.app.png"
command -v update-desktop-database >/dev/null && update-desktop-database -q "$data/applications" || true
echo "已移除。现在可以删掉 PasteFerry 所在的文件夹。"

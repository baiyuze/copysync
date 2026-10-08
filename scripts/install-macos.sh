#!/usr/bin/env bash
# 把 daemon 安装到用户目录并注册为开机自启（launchd LaunchAgent）。
#
# 用 LaunchAgent 而非 LaunchDaemon：剪贴板属于用户会话，
# 系统级守护进程根本访问不到登录用户的剪贴板。

set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"

# 发布包里 app 与本脚本同级；源码仓库里则在 dist/ 下
if [[ -d "$HERE/CopySyncDaemon.app" ]]; then
    APP_SRC="$HERE/CopySyncDaemon.app"
else
    APP_SRC="$HERE/../dist/CopySyncDaemon.app"
fi
INSTALL_DIR="$HOME/Applications"
APP_DST="$INSTALL_DIR/CopySyncDaemon.app"
BIN="$APP_DST/Contents/MacOS/copysyncd"
LABEL="com.copysync.daemon"
PLIST="$HOME/Library/LaunchAgents/$LABEL.plist"
LOG_DIR="$HOME/Library/Logs/CopySync"

[[ -d "$APP_SRC" ]] || { echo "找不到 $APP_SRC，请先运行 ./scripts/build.sh"; exit 1; }

echo "▶ 停止已在运行的实例"
launchctl bootout "gui/$(id -u)/$LABEL" 2>/dev/null || true
pkill -f copysyncd 2>/dev/null || true
# bootout 是异步的：旧服务没注销干净就 bootstrap 会报 "Bootstrap failed: 5"
for _ in {1..20}; do
    launchctl print "gui/$(id -u)/$LABEL" &>/dev/null || break
    sleep 0.25
done

echo "▶ 安装到 $APP_DST"
mkdir -p "$INSTALL_DIR" "$LOG_DIR"
rm -rf "$APP_DST"
cp -R "$APP_SRC" "$APP_DST"
# 经 AirDrop、浏览器、聊天软件传过来的包带隔离属性，ad-hoc 签名的 app
# 会被 Gatekeeper 拦下，而 launchd 拉起时没有弹窗让用户放行
xattr -dr com.apple.quarantine "$APP_DST" 2>/dev/null || true

echo "▶ 写入 LaunchAgent"
cat > "$PLIST" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>              <string>$LABEL</string>
    <key>ProgramArguments</key>
    <array>
        <string>$BIN</string>
    </array>
    <key>RunAtLoad</key>          <true/>
    <key>KeepAlive</key>          <true/>
    <!-- 崩溃后不要疯狂重启，给系统留出喘息 -->
    <key>ThrottleInterval</key>   <integer>10</integer>
    <key>StandardOutPath</key>    <string>$LOG_DIR/daemon.log</string>
    <key>StandardErrorPath</key>  <string>$LOG_DIR/daemon.log</string>
    <key>ProcessType</key>        <string>Interactive</string>
</dict>
</plist>
PLIST

echo "▶ 启动"
launchctl bootstrap "gui/$(id -u)" "$PLIST"
sleep 2

if pgrep -f copysyncd > /dev/null; then
    echo "  ✓ daemon 已启动"
else
    echo "  ✗ 启动失败，查看 $LOG_DIR/daemon.log"
    exit 1
fi

cat <<EOF

安装完成。

  日志      $LOG_DIR/daemon.log
  数据目录  ~/Library/Application Support/CopySync
  停止      launchctl bootout gui/\$(id -u)/$LABEL
  卸载      上面的命令 + rm -rf "$APP_DST" "$PLIST"

接下来：
  1. 打开 CopySync.app 配置信令服务器地址并配对设备
  2. 若系统弹出「想要从其他 App 粘贴」，选择允许；
     之后可在 系统设置 → 隐私与安全性 → 从其他 App 粘贴 中改为「始终允许」
EOF

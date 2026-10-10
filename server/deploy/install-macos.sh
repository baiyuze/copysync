#!/usr/bin/env bash
# 在 Mac 上安装 copysync-server，登录后自动在后台运行（launchd LaunchAgent），不需要管理员权限。
# 发布包 copysync-server-macos.tar.gz 里的 install.sh 就是这个文件。
#
#   ./install.sh                 用本机的局域网地址启用 TURN 中转
#   ./install.sh <IP>            指定 TURN 对外公布的地址
#   ./install.sh uninstall       卸载
#
# 重复执行即为升级：替换程序并重启，已有设置（含 TURN 密钥）保留；指定了地址时改用新地址。

set -euo pipefail
cd "$(dirname "$0")"

# 中文系统显示中文，其他显示英文
LANGS=$(defaults read -g AppleLanguages 2>/dev/null | tr -d ' \n' || true)
[[ ${LANGS#(\"} == zh* || ${LC_ALL:-${LANG:-}} == zh* ]] && ZH=1 || ZH=
say() { if [[ $ZH ]]; then echo "$1"; else echo "$2"; fi; }

LABEL=com.copysync.server
DIR="$HOME/Library/Application Support/CopySync Server"
BIN="$DIR/copysync-server"
CONF="$DIR/server.conf"
PLIST="$HOME/Library/LaunchAgents/$LABEL.plist"
LOG="$HOME/Library/Logs/CopySync/server.log"
DOMAIN="gui/$(id -u)"

stop() {
    launchctl bootout "$DOMAIN/$LABEL" 2>/dev/null || true
    # bootout 是异步的：旧服务没注销干净就 bootstrap 会报 "Bootstrap failed: 5"
    for _ in {1..20}; do
        launchctl print "$DOMAIN/$LABEL" &>/dev/null || break
        sleep 0.25
    done
}

if [[ ${1:-} == uninstall ]]; then
    stop
    rm -rf "$DIR" "$PLIST"
    say "✓ 已卸载（日志留在 $LOG）" "✓ Uninstalled (the log is kept at $LOG)"
    exit 0
fi

# 本机的局域网地址：默认路由所用网卡的地址。代理软件的 TUN 模式会把默认路由指到
# utun（198.18.0.0/15），这时改看 Wi‑Fi 与有线网卡。
lan_ip() {
    local ifc ip
    ifc=$(route -n get default 2>/dev/null | awk '/interface:/{print $2}')
    [[ -n $ifc ]] && ip=$(ipconfig getifaddr "$ifc" 2>/dev/null || true)
    if [[ -z ${ip:-} || $ip == 198.1[89].* ]]; then
        ip=
        for ifc in en0 en1 en2 en3 en4 en5; do
            ip=$(ipconfig getifaddr "$ifc" 2>/dev/null || true)
            [[ -n $ip ]] && break
        done
    fi
    echo "$ip"
}

TURN_IP= TURN_SECRET=
# shellcheck source=/dev/null
[[ -f $CONF ]] && source "$CONF"
if [[ -n ${1:-} ]]; then
    TURN_IP=$1
elif [[ -z $TURN_IP ]]; then
    TURN_IP=$(lan_ip)
fi
[[ -n $TURN_IP ]] || {
    say "找不到本机的局域网地址，请指定：$0 <IP>" "Couldn't find this Mac's LAN address. Specify it: $0 <IP>"
    exit 1
}
# 固定 secret：每次重启都换的话，客户端手里已签发的 TURN 凭证会失效
[[ -n $TURN_SECRET ]] || TURN_SECRET=$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')

say "▶ 停止旧版本" "▶ Stopping the running version"
stop

say "▶ 安装到 $DIR" "▶ Installing to $DIR"
mkdir -p "$DIR" "$(dirname "$LOG")" "$(dirname "$PLIST")"
cp copysync-server "$BIN"
chmod 755 "$BIN"
# 用浏览器下载的包带隔离属性，Gatekeeper 会拦下没有公证的程序，而 launchd 拉起时没有弹窗可以放行
xattr -d com.apple.quarantine "$BIN" 2>/dev/null || true
printf 'TURN_IP=%s\nTURN_SECRET=%s\n' "$TURN_IP" "$TURN_SECRET" > "$CONF"
chmod 600 "$CONF"

say "▶ 写入 LaunchAgent（中转地址 $TURN_IP）" "▶ Writing the LaunchAgent (relay address $TURN_IP)"
cat > "$PLIST" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>              <string>$LABEL</string>
    <key>ProgramArguments</key>
    <array>
        <string>$BIN</string>
        <string>-addr</string>        <string>:8787</string>
        <string>-turn-ip</string>     <string>$TURN_IP</string>
        <string>-turn-secret</string> <string>$TURN_SECRET</string>
        <string>-stun</string>        <string>stun:$TURN_IP:3478,stun:stun.l.google.com:19302</string>
        <string>-log</string>         <string>$LOG</string>
    </array>
    <key>RunAtLoad</key>          <true/>
    <key>KeepAlive</key>          <true/>
    <!-- 起不来时（比如端口被占用）不要疯狂重启 -->
    <key>ThrottleInterval</key>   <integer>10</integer>
    <!-- 参数错误之类在打开日志文件之前的输出 -->
    <key>StandardErrorPath</key>  <string>$LOG</string>
</dict>
</plist>
PLIST
chmod 600 "$PLIST"

say "▶ 启动" "▶ Starting"
launchctl bootstrap "$DOMAIN" "$PLIST"
ok=
for _ in {1..20}; do
    curl -fsS --noproxy '*' -m 1 http://127.0.0.1:8787/healthz &>/dev/null && { ok=1; break; }
    sleep 0.25
done
if [[ $ok ]]; then
    say "  ✓ 已启动" "  ✓ Running"
else
    say "  ✗ 启动失败，查看日志：tail -n 50 \"$LOG\"" "  ✗ Failed to start. See the log: tail -n 50 \"$LOG\""
    exit 1
fi

VERSION=$("$BIN" -version)
if [[ $ZH ]]; then
    cat <<EOF

安装完成：$VERSION，登录这台 Mac 后自动运行
在每台电脑的 CopySync「设置 → 信令服务器地址」填：

    ws://$TURN_IP:8787/signal

  日志   $LOG
  重启   launchctl kickstart -k $DOMAIN/$LABEL
  卸载   $0 uninstall
EOF
else
    cat <<EOF

Installed: $VERSION. It starts whenever you log in to this Mac.
In CopySync on each computer, set Settings → Signaling server to:

    ws://$TURN_IP:8787/signal

  Logs       $LOG
  Restart    launchctl kickstart -k $DOMAIN/$LABEL
  Uninstall  $0 uninstall
EOF
fi

# 系统防火墙开着时，没有放行的程序收不到其他电脑的连接
FW=/usr/libexec/ApplicationFirewall/socketfilterfw
if [[ -x $FW ]] && "$FW" --getglobalstate 2>/dev/null | grep -q "enabled"; then
    echo
    say "系统防火墙已开启，放行 copysync-server（需要输入密码）：" \
        "The macOS firewall is on. Allow copysync-server through it (asks for your password):"
    echo "  sudo $FW --add \"$BIN\" && sudo $FW --unblockapp \"$BIN\""
fi

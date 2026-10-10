#!/usr/bin/env bash
# 在 Linux 上安装 copysync-server 并注册为开机自启的 systemd 服务。
#
#   sudo ./install.sh                 用本机的局域网地址启用 TURN 中转（局域网部署）
#   sudo ./install.sh <公网IP>        用公网 IP 启用 TURN 中转（云服务器，跨公网使用）
#   sudo ./install.sh uninstall       卸载
#
# 重复执行即为升级：替换二进制并重启，已有的 /etc/copysync/server.env 不会被覆盖。

set -euo pipefail
cd "$(dirname "$0")"

# 中文系统显示中文，其他显示英文
[[ ${LC_ALL:-${LC_MESSAGES:-${LANG:-}}} == zh* ]] && ZH=1 || ZH=
say() { if [[ $ZH ]]; then echo "$1"; else echo "$2"; fi; }

[[ $EUID -eq 0 ]] || { say "请用 root 运行：sudo $0 $*" "Run as root: sudo $0 $*"; exit 1; }

ENV_FILE=/etc/copysync/server.env

if [[ ${1:-} == uninstall ]]; then
    systemctl disable --now copysync-server 2>/dev/null || true
    rm -f /usr/local/bin/copysync-server /etc/systemd/system/copysync-server.service
    rm -rf /etc/copysync
    systemctl daemon-reload
    say "✓ 已卸载" "✓ Uninstalled"
    exit 0
fi

# 本机的局域网地址：默认路由所用网卡的地址。代理软件的 TUN 模式会把默认路由指到
# 198.18.0.0/15，这时改用其他网卡上的地址；Docker 的网桥也跳过。
lan_ip() {
    local ip
    ip=$(ip -4 route get 1.1.1.1 2>/dev/null | sed -n 's/.* src \([0-9.]*\).*/\1/p' | head -n1)
    if [[ -z $ip || $ip == 198.1[89].* ]]; then
        ip=$(ip -4 -o addr show scope global 2>/dev/null | awk '{print $4}' | cut -d/ -f1 |
            grep -Ev '^(198\.1[89]\.|172\.17\.)' | head -n1)
    fi
    echo "$ip"
}

say "▶ 安装二进制到 /usr/local/bin" "▶ Installing the binary to /usr/local/bin"
install -m 755 copysync-server /usr/local/bin/copysync-server

say "▶ 写入 systemd 单元" "▶ Writing the systemd unit"
install -m 644 copysync-server.service /etc/systemd/system/copysync-server.service

if [[ -f $ENV_FILE ]]; then
    say "▶ 保留已有配置 $ENV_FILE" "▶ Keeping the existing config $ENV_FILE"
    TURN_IP=$(sed -n 's/.*-turn-ip \([^ "]*\).*/\1/p' "$ENV_FILE")
    if [[ -n ${1:-} && ${1:-} != "$TURN_IP" ]]; then
        say "  要改用 $1，编辑 $ENV_FILE 里的 -turn-ip 后执行 systemctl restart copysync-server" \
            "  To switch to $1, edit -turn-ip in $ENV_FILE and run systemctl restart copysync-server"
    fi
else
    TURN_IP=${1:-$(lan_ip)}
    [[ -n $TURN_IP ]] || {
        say "找不到本机的局域网地址，请指定：sudo $0 <IP>" "Couldn't find this machine's LAN address. Specify it: sudo $0 <IP>"
        exit 1
    }
    say "▶ 生成配置 $ENV_FILE（中转地址 $TURN_IP）" "▶ Writing $ENV_FILE (relay address $TURN_IP)"
    mkdir -p /etc/copysync
    # 固定 secret：自动生成的话每次重启都会让已签发的 TURN 凭证失效
    # STUN 优先用自己：国内网络访问公共 STUN（如 Google 的）时通时断
    OPTIONS="-addr :8787 -turn-ip $TURN_IP -turn-secret $(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')"
    OPTIONS="$OPTIONS -stun stun:$TURN_IP:3478,stun:stun.l.google.com:19302"
    printf 'OPTIONS="%s"\n' "$OPTIONS" > "$ENV_FILE"
    chmod 600 "$ENV_FILE"
fi

say "▶ 启动" "▶ Starting"
systemctl daemon-reload
systemctl enable copysync-server >/dev/null 2>&1
systemctl restart copysync-server
sleep 1

if systemctl is-active --quiet copysync-server; then
    say "  ✓ 已启动" "  ✓ Running"
else
    say "  ✗ 启动失败，查看日志：journalctl -u copysync-server -n 50" \
        "  ✗ Failed to start. See the log: journalctl -u copysync-server -n 50"
    exit 1
fi

ADDR=${TURN_IP:-$(lan_ip)}
PORTS=$(tr '\t' '-' < /proc/sys/net/ipv4/ip_local_port_range 2>/dev/null || echo 32768-60999)
VERSION=$(/usr/local/bin/copysync-server -version)
if [[ $ZH ]]; then
    cat <<EOF

安装完成：$VERSION
在每台电脑的 CopySync「设置 → 信令服务器地址」填：

    ws://$ADDR:8787/signal

  配置   $ENV_FILE
  日志   journalctl -u copysync-server -f
  重启   systemctl restart copysync-server
  卸载   sudo $0 uninstall

防火墙需放行：TCP 8787（信令）、UDP 3478（STUN 与 TURN）、UDP $PORTS（中转端口，系统随机分配）
EOF
else
    cat <<EOF

Installed: $VERSION
In CopySync on each computer, set Settings → Signaling server to:

    ws://$ADDR:8787/signal

  Config     $ENV_FILE
  Logs       journalctl -u copysync-server -f
  Restart    systemctl restart copysync-server
  Uninstall  sudo $0 uninstall

Open in your firewall: TCP 8787 (signaling), UDP 3478 (STUN and TURN), UDP $PORTS (relay ports, assigned by the OS)
EOF
fi

# 开着的防火墙直接给出放行命令
if command -v ufw >/dev/null && ufw status 2>/dev/null | grep -q "Status: active"; then
    say "ufw 已启用，放行命令：" "ufw is active. To open the ports:"
    echo "  sudo ufw allow 8787/tcp && sudo ufw allow 3478/udp && sudo ufw allow ${PORTS/-/:}/udp"
elif command -v firewall-cmd >/dev/null && firewall-cmd --state &>/dev/null; then
    say "firewalld 已启用，放行命令：" "firewalld is running. To open the ports:"
    echo "  sudo firewall-cmd --permanent --add-port=8787/tcp --add-port=3478/udp --add-port=$PORTS/udp && sudo firewall-cmd --reload"
fi

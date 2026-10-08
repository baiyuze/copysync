#!/usr/bin/env bash
# 在 Linux 服务器上安装 copysync-server 并注册为 systemd 服务。
#
#   sudo ./install.sh                 仅信令 + STUN（同一局域网够用）
#   sudo ./install.sh <公网IP>        同时启用 TURN 中转（跨公网、对称 NAT 时需要）
#
# 重复执行即为升级：替换二进制并重启，已有的 /etc/copysync/server.env 不会被覆盖。

set -euo pipefail
cd "$(dirname "$0")"

[[ $EUID -eq 0 ]] || { echo "请用 root 运行：sudo $0 $*"; exit 1; }

TURN_IP=${1:-}
ENV_FILE=/etc/copysync/server.env

echo "▶ 安装二进制到 /usr/local/bin"
install -m 755 copysync-server /usr/local/bin/copysync-server

echo "▶ 写入 systemd 单元"
install -m 644 copysync-server.service /etc/systemd/system/copysync-server.service

if [[ -f $ENV_FILE ]]; then
    echo "▶ 保留已有配置 $ENV_FILE"
else
    echo "▶ 生成配置 $ENV_FILE"
    mkdir -p /etc/copysync
    if [[ -n $TURN_IP ]]; then
        # 固定 secret：自动生成的话每次重启都会让已签发的 TURN 凭证失效
        # STUN 优先用自己：国内网络访问公共 STUN（如 Google 的）时通时断
        OPTIONS="-addr :8787 -turn-ip $TURN_IP -turn-secret $(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')"
        OPTIONS="$OPTIONS -stun stun:$TURN_IP:3478,stun:stun.l.google.com:19302"
    else
        OPTIONS="-addr :8787"
    fi
    printf 'OPTIONS="%s"\n' "$OPTIONS" > "$ENV_FILE"
    chmod 600 "$ENV_FILE"
fi

echo "▶ 启动"
systemctl daemon-reload
systemctl enable copysync-server >/dev/null
systemctl restart copysync-server
sleep 1

if systemctl is-active --quiet copysync-server; then
    echo "  ✓ 已启动"
else
    echo "  ✗ 启动失败，查看日志：journalctl -u copysync-server -n 50"
    exit 1
fi

cat <<EOF

安装完成：$(/usr/local/bin/copysync-server -version)
客户端「设置 → 信令服务器地址」填：ws://${TURN_IP:-<本机地址>}:8787/signal

  配置   $ENV_FILE
  日志   journalctl -u copysync-server -f
  重启   systemctl restart copysync-server
  卸载   systemctl disable --now copysync-server && rm /usr/local/bin/copysync-server /etc/systemd/system/copysync-server.service

防火墙 / 云安全组需放行：
  TCP 8787                信令（WebSocket）
  UDP 3478                TURN（仅启用中转时）
  UDP 32768-60999         TURN 中转分配的端口（仅启用中转时；由系统随机分配，
                          实际范围见 /proc/sys/net/ipv4/ip_local_port_range）
EOF

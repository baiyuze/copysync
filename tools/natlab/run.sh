#!/usr/bin/env bash
# NAT 打洞实验室：用 Linux 网络命名空间与 iptables 搭出真实的网络拓扑，
# 验证两台设备最终走直连还是中转。
#
#   sudo ./tools/natlab/run.sh <natlab 可执行文件>
#
# 需要 root、iproute2、iptables。macOS 上用 tools/natlab/docker.sh，它会在容器里跑本脚本。
#
# 拓扑（全部用内网地址扮演公网，natlab 以 ProbeAllowPrivate 运行）：
#
#   [a] 192.168.10.2 ── [rA] 出口 10.255.0.11–14 ──┐
#                                                  [inet] 信令 10.200.0.2  TURN 10.200.0.3
#   [b] 192.168.20.2 ── [rB] 出口 10.254.0.21–22 ──┘        STUN 10.200.1.1–4
#
# 10.200.1.1 扮演自己的服务器（信令服务器下发的 STUN），10.200.1.2–4 扮演公共 STUN。
# rA 按目标地址选出口，模拟公司双线、运营商地址池；入站只放行有连接记录的包，
# 与真实路由器一样「只认打过招呼的地址」。

set -euo pipefail

BIN=$(cd "$(dirname "${1:?用法: run.sh <natlab 可执行文件>}")" && pwd)/$(basename "$1")
WORK=$(mktemp -d)
PROBE="10.200.1.2:3478,10.200.1.3:3478,10.200.1.4:3478"

ns() { ip netns exec "$@"; }

INFRA_PID=""

# 按进程号结束 infra：用 pkill -f 按路径匹配会连本脚本自己一起杀掉（命令行里也有这个路径）
cleanup() {
    if [[ -n $INFRA_PID ]]; then
        kill "$INFRA_PID" 2>/dev/null || true
        wait "$INFRA_PID" 2>/dev/null || true
        INFRA_PID=""
    fi
    for n in inet rA a rB b; do ip netns del "$n" 2>/dev/null || true; done
}
trap 'cleanup; rm -rf "$WORK"' EXIT

veth() { # <命名空间1> <网卡1> <命名空间2> <网卡2>
    ip link add "$2" netns "$1" type veth peer name "$4" netns "$3"
}

# setup <A 侧出口模式> <B 侧模式>
#   A：single 单出口 | dual 两个出口，去对端走出口 2、探测服务器有一台走出口 2
#      hidden 两个出口，但没有任何探测服务器走出口 2 | quad 四个出口
#   B：cone 家用路由器 | symmetric 端口随机分配 | dual 两个出口
setup() {
    cleanup
    for n in inet rA a rB b; do
        ip netns add "$n"
        ns "$n" ip link set lo up
    done
    veth inet i-a rA wan
    veth inet i-b rB wan
    veth rA lan a eth0
    veth rB lan b eth0

    ns inet ip addr add 10.255.0.1/24 dev i-a
    ns inet ip addr add 10.254.0.1/24 dev i-b
    for ip in 10.200.0.2 10.200.0.3 10.200.1.1 10.200.1.2 10.200.1.3 10.200.1.4; do
        ns inet ip addr add "$ip/32" dev lo
    done
    ns inet ip link set i-a up
    ns inet ip link set i-b up
    ns inet sysctl -qw net.ipv4.ip_forward=1

    for i in 1 2 3 4; do ns rA ip addr add "10.255.0.1$i/24" dev wan; done
    ns rA ip addr add 192.168.10.1/24 dev lan
    ns rA ip link set wan up
    ns rA ip link set lan up
    ns rA ip route add default via 10.255.0.1
    ns rA sysctl -qw net.ipv4.ip_forward=1
    ns a ip addr add 192.168.10.2/24 dev eth0
    ns a ip link set eth0 up
    ns a ip route add default via 192.168.10.1

    for i in 1 2; do ns rB ip addr add "10.254.0.2$i/24" dev wan; done
    ns rB ip addr add 192.168.20.1/24 dev lan
    ns rB ip link set wan up
    ns rB ip link set lan up
    ns rB ip route add default via 10.254.0.1
    ns rB sysctl -qw net.ipv4.ip_forward=1
    ns b ip addr add 192.168.20.2/24 dev eth0
    ns b ip link set eth0 up
    ns b ip route add default via 192.168.20.1

    # 路由器的防火墙：丢弃外网主动发来的新连接，与真实的家用、公司路由器一致。
    # 不能省：被丢弃的包不会留下连接记录。若放行到路由器本机，Linux 会为它建一条
    # 「未应答」的连接记录，随后本端往外发时源端口与之冲突，SNAT 只好换一个端口，
    # 事先告诉对端的地址就失效了——真实路由器不会这样。
    for r in rA rB; do
        ns "$r" iptables -A INPUT -i wan -m conntrack --ctstate NEW -j DROP
        ns "$r" iptables -A FORWARD -i wan -m conntrack --ctstate NEW -j DROP
    done

    local snat="iptables -t nat -A POSTROUTING -o wan"
    case $1 in
    single) ns rA $snat -j SNAT --to-source 10.255.0.11 ;;
    dual)
        ns rA $snat -d 10.254.0.0/24 -j SNAT --to-source 10.255.0.12
        ns rA $snat -d 10.200.1.2 -j SNAT --to-source 10.255.0.12
        ns rA $snat -j SNAT --to-source 10.255.0.11
        ;;
    hidden)
        ns rA $snat -d 10.254.0.0/24 -j SNAT --to-source 10.255.0.12
        ns rA $snat -j SNAT --to-source 10.255.0.11
        ;;
    quad)
        ns rA $snat -d 10.254.0.0/24 -j SNAT --to-source 10.255.0.14
        for i in 2 3 4; do ns rA $snat -d "10.200.1.$i" -j SNAT --to-source "10.255.0.1$i"; done
        ns rA $snat -j SNAT --to-source 10.255.0.11
        ;;
    esac
    case $2 in
    cone) ns rB $snat -j SNAT --to-source 10.254.0.21 ;;
    symmetric) ns rB $snat -j SNAT --to-source 10.254.0.21 --random-fully ;;
    dual)
        ns rB $snat -d 10.255.0.0/24 -j SNAT --to-source 10.254.0.22
        ns rB $snat -d 10.200.1.3 -j SNAT --to-source 10.254.0.22
        ns rB $snat -j SNAT --to-source 10.254.0.21
        ;;
    esac

    ns inet "$BIN" infra -http 10.200.0.2:9000 -turn 10.200.0.3 \
        -stun 10.200.1.1,10.200.1.2,10.200.1.3,10.200.1.4 >"$WORK/infra.log" 2>&1 &
    INFRA_PID=$!
    for _ in $(seq 50); do
        ns a bash -c 'exec 3<>/dev/tcp/10.200.0.2/9000' 2>/dev/null && return
        sleep 0.1
    done
    echo "infra 没有起来：" && cat "$WORK/infra.log" && exit 1
}

PASS=0
FAIL=0

# check <场景名> <A 侧> <B 侧> <direct|relay> [legacy|history]
check() {
    # NATLAB_ONLY 只跑名字里含这段文字的场景，排查问题时用
    [[ -n ${NATLAB_ONLY:-} && $1 != *"$NATLAB_ONLY"* ]] && return
    local name=$1 expect=$4 flag=${5:-} aflags=() state="$WORK/a-state.json"
    rm -f "$state" # 每个场景从空白开始：上一个场景的出口历史会让这一个「意外」直连
    setup "$2" "$3"
    case $flag in
    legacy) aflags=(-legacy) ;;
    history)
        # 预置一份历史：这个网络里见过不改端口的出口 2
        printf '{"networks":{"192.168.10.2":[{"ip":"10.255.0.12","port_preserved":true,"last_seen":"%s"}]}}' \
            "$(date -u +%Y-%m-%dT%H:%M:%SZ)" >"$state"
        ;;
    esac
    local common=(-infra http://10.200.0.2:9000 -stun 10.200.1.1:3478 -timeout 30s)
    ns b "$BIN" peer -id b -peer a "${common[@]}" -probe "$PROBE" >"$WORK/b.json" 2>"$WORK/b.log" &
    local bpid=$!
    ns a "$BIN" peer -id a -peer b "${common[@]}" -probe "$PROBE" -state "$state" "${aflags[@]}" \
        >"$WORK/a.json" 2>"$WORK/a.log" || true
    wait "$bpid" || true

    local a b
    a=$(sed -n 's/.*"state":"\([a-z]*\)".*/\1/p' "$WORK/a.json")
    b=$(sed -n 's/.*"state":"\([a-z]*\)".*/\1/p' "$WORK/b.json")
    if [[ $a == "$expect" && $b == "$expect" ]]; then
        PASS=$((PASS + 1))
        printf '  ✓ %-34s %s\n' "$name" "$a"
    else
        FAIL=$((FAIL + 1))
        printf '  ✗ %-34s 期望 %s，实际 a=%s b=%s\n' "$name" "$expect" "${a:-无结果}" "${b:-无结果}"
        echo "    a: $(cat "$WORK/a.json")"
        echo "    b: $(cat "$WORK/b.json")"
        grep -h -E '出口探测完成|P2P|失败|error' "$WORK/a.log" "$WORK/b.log" | head -8 | sed 's/^/    /'
        # NATLAB_KEEP 指定目录时把完整日志留下来
        [[ -n ${NATLAB_KEEP:-} ]] && cp "$WORK"/*.log "$NATLAB_KEEP/" 2>/dev/null || true
    fi
}

echo "NAT 打洞实验室"
check "单出口 ↔ 家用路由器" single cone direct
check "两个出口（旧版本，对照）" dual cone relay legacy
check "两个出口 ↔ 家用路由器" dual cone direct
check "四个出口 ↔ 家用路由器" quad cone direct
check "两个出口，第二条没探测到，无历史" hidden cone relay
check "两个出口，第二条没探测到，有历史" hidden cone direct history
check "两个出口 ↔ 两个出口" dual dual direct
check "两个出口 ↔ 端口随机分配的 NAT" dual symmetric relay

echo
echo "通过 $PASS，失败 $FAIL"
[[ $FAIL -eq 0 ]]

#!/usr/bin/env bash
# 剪贴板隐私机制对照实验
#
# 同一套读取操作，分别在「隐私机制关闭」和「隐私机制开启」下各跑一轮，
# 直接对比哪些操作会被拦截。开启的那一轮预演的就是下一个 macOS 大版本的行为。
#
# 读方 = SpikeA (com.copysync.spike.privacy)
# 写方 = SpikeC (com.copysync.spike.writer)  ← 独立 bundle id 的真实 GUI App，
#        这点很重要：pbcopy/osascript 是命令行工具，TCC 可能不当作「其他 App」

set -uo pipefail
cd "$(dirname "$0")"

BUNDLE=com.copysync.spike.privacy
A=./build/SpikeA.app/Contents/MacOS/SpikeA
C=./build/SpikeC.app/Contents/MacOS/SpikeC
TESTFILE=/tmp/copysync-spike-test.txt
MODE=${1:-file}     # file | text | image

[[ -x "$A" && -x "$C" ]] || { echo "请先运行 ./build.sh"; exit 1; }
echo "CopySync spike 测试文件 $(date)" > "$TESTFILE"

# 等待进程结束，最多 $2 秒。返回 1 表示超时（通常意味着弹窗在等待用户点击）
wait_pid() {
    local pid=$1 limit=$2 i=0
    while kill -0 "$pid" 2>/dev/null; do
        sleep 1; i=$((i+1))
        [[ $i -ge $limit ]] && return 1
    done
    return 0
}

run_round() {
    local label="$1" log="$2" limit="$3"

    echo
    echo "┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "┃  $label"
    echo "┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    "$A" 1 > "$log" 2>&1 &
    local apid=$!
    sleep 2

    "$C" "$MODE" "$TESTFILE" > /dev/null 2>&1 &
    local cpid=$!

    if ! wait_pid "$apid" "$limit"; then
        echo
        echo "  ⚠️  超过 ${limit}s 未结束 —— 很可能有授权弹窗正等着你点击。"
        echo "      请看屏幕，点击弹窗后脚本会继续（再等 60s）。"
        wait_pid "$apid" 60 || { echo "  仍未结束，强制终止。"; kill -9 "$apid" 2>/dev/null; }
    fi

    kill "$cpid" 2>/dev/null
    wait "$cpid" 2>/dev/null
    cat "$log"
}

echo "模式：$MODE    测试文件：$TESTFILE"

# ── 第一轮：隐私机制关闭（当前 macOS 15.7 的默认行为）────────
defaults delete "$BUNDLE" EnablePasteboardPrivacyDeveloperPreview 2>/dev/null
run_round "第 1 轮 — 隐私机制【关闭】（今天的默认行为）" /tmp/spike-priv-off.log 30

# ── 第二轮：隐私机制开启（预演下一个 macOS 大版本）──────────
defaults write "$BUNDLE" EnablePasteboardPrivacyDeveloperPreview -bool yes
run_round "第 2 轮 — 隐私机制【开启】（预演未来 macOS 的强制行为）" /tmp/spike-priv-on.log 25

cat <<EOF

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
对比要点
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  · 阶段 1/2/3 在两轮中是否都无阻塞？
      → 决定 CopySync 能否免提示地「感知变化 + 判断类型」
  · 阶段 4（读实际内容）在第 2 轮是否阻塞/弹窗？
      → 决定 A 端拿源文件路径是否需要用户授权
  · 阶段 5 的 accessBehavior 是否从 Default 迁移到 Ask？
      → 验证「首次触发后出现在系统设置中」

日志：/tmp/spike-priv-off.log    /tmp/spike-priv-on.log

清理开关（重要，避免影响后续测试）：
  defaults delete $BUNDLE EnablePasteboardPrivacyDeveloperPreview
EOF

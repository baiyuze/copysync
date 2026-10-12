"""统计 GitHub Release 的下载次数：按平台汇总、按版本列出，并可记录快照看增长。

    python3 tools/stats/downloads.py            看当前的累计下载次数
    python3 tools/stats/downloads.py --record   同时记一份快照，并显示与上次、7 天前相比新增了多少

数据来自 GitHub 的 Release 接口（每个文件的 download_count），只有累计值，没有按天的明细，
所以想看趋势就要定期记快照。快照默认存在 ~/.copysync-stats/downloads.jsonl，不进仓库。
官网的下载按钮链接的就是这些文件，所以从官网下载的都算在里面；测试时自己下载的也算。

只用 Python 标准库。直连 api.github.com 不通时，设置 HTTPS_PROXY 即可。
"""

import argparse
import json
import re
import sys
import urllib.request
from datetime import datetime, timedelta
from pathlib import Path

REPO = "baiyuze/copysync"
DEFAULT_LOG = Path.home() / ".copysync-stats" / "downloads.jsonl"

# 文件名 → 类别。1.0.0 的文件名带版本号，后来的不带，两种都要认
# 1.6.0 起 Mac 分 Apple 芯片、Intel 与通用三个 DMG，合在一起算；Linux 的 .deb 与 tar.gz 也合在一起
CATEGORIES = [
    ("Mac 客户端", re.compile(r"^CopySync.*\.dmg$")),
    ("Windows 安装包", re.compile(r"^CopySync-Setup\.exe$")),
    ("Windows 便携版", re.compile(r"^CopySync-windows-x64\.zip$")),
    ("Linux 客户端", re.compile(r"^CopySync-linux-(amd64|arm64)\.(deb|tar\.gz)$")),
    ("服务器 Linux", re.compile(r"^copysync-server.*-linux-(amd64|arm64)\.tar\.gz$")),
    ("服务器 macOS", re.compile(r"^copysync-server.*-macos\.tar\.gz$")),
    ("服务器 Windows", re.compile(r"^copysync-server.*-windows-amd64\.zip$")),
]
NAMES = [name for name, _ in CATEGORIES]
CLIENTS = NAMES[:4]


def category(asset: str):
    for name, pattern in CATEGORIES:
        if pattern.match(asset):
            return name
    return None  # SHA256SUMS、RELEASE_NOTES.md 之类不算


def fetch() -> list[dict]:
    releases, page = [], 1
    while True:
        req = urllib.request.Request(
            f"https://api.github.com/repos/{REPO}/releases?per_page=100&page={page}",
            headers={"Accept": "application/vnd.github+json", "User-Agent": "copysync-stats"})
        with urllib.request.urlopen(req, timeout=30) as resp:
            batch = json.load(resp)
        releases += batch
        if len(batch) < 100:
            return releases
        page += 1


def snapshot(releases: list[dict]) -> dict:
    """{版本: {文件名: 下载次数}}"""
    return {r["tag_name"]: {a["name"]: a["download_count"] for a in r["assets"] if category(a["name"])}
            for r in releases if not r.get("draft")}


def by_category(counts: dict) -> dict:
    out = dict.fromkeys(NAMES, 0)
    for files in counts.values():
        for name, n in files.items():
            out[category(name)] += n
    return out


def width(s: str) -> int:
    """终端里的显示宽度：中文占两格。"""
    return sum(2 if ord(ch) > 0x2E7F else 1 for ch in s)


def pad(s: str, w: int, right: bool = False) -> str:
    gap = " " * max(0, w - width(s))
    return gap + s if right else s + gap


def table(header: list[str], rows: list[list[str]]) -> str:
    widths = [max(width(r[i]) for r in [header] + rows) for i in range(len(header))]
    line = lambda r: "  ".join(pad(c, widths[i], right=i > 0) for i, c in enumerate(r))
    return "\n".join([line(header), "  ".join("-" * w for w in widths)] + [line(r) for r in rows])


def load_log(path: Path) -> list[dict]:
    if not path.exists():
        return []
    return [json.loads(l) for l in path.read_text(encoding="utf-8").splitlines() if l.strip()]


def diff(now: dict, before: dict) -> dict:
    a, b = by_category(now), by_category(before)
    return {k: a[k] - b[k] for k in NAMES}


def main() -> None:
    ap = argparse.ArgumentParser(description="统计 CopySync 各版本的下载次数")
    ap.add_argument("--record", action="store_true", help="记录一份快照，并显示新增")
    ap.add_argument("--log", type=Path, default=DEFAULT_LOG, help=f"快照文件（默认 {DEFAULT_LOG}）")
    args = ap.parse_args()

    try:
        counts = snapshot(fetch())
    except OSError as e:
        sys.exit(f"连不上 GitHub：{e}\n直连不通时设置代理再试，例如 HTTPS_PROXY=http://127.0.0.1:1087")

    now = datetime.now().astimezone()
    totals = by_category(counts)
    clients = sum(totals[k] for k in CLIENTS)
    servers = sum(totals[k] for k in NAMES if k not in CLIENTS)
    print(f"CopySync 下载统计，{now:%Y-%m-%d %H:%M}（GitHub Release 累计）\n")
    print(table(["类别", "下载次数"], [[k, str(totals[k])] for k in NAMES]
                + [["客户端合计", str(clients)], ["服务器合计", str(servers)]]))

    print("\n按版本：\n")
    short = ["Mac", "Win 安装", "Win 便携", "Linux", "服务器 Linux", "服务器 Mac", "服务器 Win"]
    rows = []
    for tag, files in counts.items():  # 接口按发布时间从新到旧返回
        c = by_category({tag: files})
        rows.append([tag] + [str(c[k]) if c[k] else "·" for k in NAMES] + [str(sum(c.values()))])
    print(table(["版本"] + short + ["合计"], rows))

    if args.record:
        history = load_log(args.log)
        if history:
            last = history[-1]
            print(f"\n与上次快照（{last['time'][:16].replace('T', ' ')}）相比新增：")
            added = "，".join(f"{k} +{v}" for k, v in diff(counts, last["counts"]).items() if v)
            print("  " + (added or "没有新增"))
            week = [h for h in history if datetime.fromisoformat(h["time"]) <= now - timedelta(days=7)]
            if week:
                d = diff(counts, week[-1]["counts"])
                print(f"最近 7 天新增：客户端 +{sum(d[k] for k in CLIENTS)}，服务器 +{sum(d[k] for k in NAMES if k not in CLIENTS)}")
        args.log.parent.mkdir(parents=True, exist_ok=True)
        with args.log.open("a", encoding="utf-8") as f:
            f.write(json.dumps({"time": now.isoformat(timespec="seconds"), "counts": counts}, ensure_ascii=False) + "\n")
        print(f"\n快照已记录到 {args.log}")


if __name__ == "__main__":
    main()

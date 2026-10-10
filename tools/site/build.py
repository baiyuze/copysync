"""生成官网：docs/index.html（简体中文）、docs/en/、docs/ja/，以及 docs/design/ 下的技术方案页面。

    python3 tools/site/build.py

三种语言共用同一套页面结构，文案在 tools/site/content/ 下每种语言一份。
技术方案页面由 design/*.md 转换而来，放在网站上，不用跳到 GitHub 去看。
改了文案、结构或技术方案后运行一次，把生成的文件一起提交。只用 Python 标准库。
"""

import html
import importlib.util
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
DOCS = ROOT / "docs"
SITE = "https://baiyuze.github.io/copysync/"
REPO = "https://github.com/baiyuze/copysync"
DL = REPO + "/releases/latest/download/"
VERSION = "1.3.3"
DATE = "2026-10-11"

# 语言代码 → (页面目录, html 的 lang, og:locale, 语言名)
LANGS = {
    "zh": ("", "zh-CN", "zh_CN", "简体中文"),
    "en": ("en/", "en", "en_US", "English"),
    "ja": ("ja/", "ja", "ja_JP", "日本語"),
}

# 「在局域网里部署服务器」教程，每种语言目录下一份
GUIDE = "server/"

# 放到网站上的技术方案：文件名 → 网页名。其他 design/*.md 链接到 GitHub
DESIGN_DOCS = {"nat-traversal.md": "nat-traversal", "windows-client.md": "windows-client"}

DOWN = ('<svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.8" '
        'stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">'
        '<path d="M10 3v10m0 0-4-4m4 4 4-4M4 16h12"/></svg>')
LAPTOP = ('<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" '
          'stroke-linejoin="round" aria-hidden="true"><rect x="4" y="5" width="16" height="11" rx="1.5"/>'
          '<path d="M2 19h20"/></svg>')
DESKTOP = ('<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" '
           'stroke-linejoin="round" aria-hidden="true"><rect x="3" y="4" width="18" height="12" rx="1.5"/>'
           '<path d="M9 20h6M12 16v4"/></svg>')
SERVER = ('<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" '
          'stroke-linejoin="round" aria-hidden="true"><rect x="4" y="4" width="16" height="7" rx="1.5"/>'
          '<rect x="4" y="13" width="16" height="7" rx="1.5"/><path d="M8 7.5h.01M8 16.5h.01"/></svg>')
LINES = '<span class="lines"><i></i><i></i><i></i></span>'

# 首页按浏览器语言跳转：只在简体中文的首页（网站根目录）做，选过语言后以选的为准
REDIRECT = """  <script>
    // 第一次打开时按浏览器语言跳到对应页面；在页面上点过语言链接后，以点的为准
    (function () {
      try {
        var pick = localStorage.getItem("copysync-lang");
        if (!pick) {
          if (/bot|crawl|spider|slurp|lighthouse/i.test(navigator.userAgent)) return;
          var list = navigator.languages || [navigator.language || ""];
          for (var i = 0; i < list.length && !pick; i++) {
            var l = String(list[i]).toLowerCase();
            if (l.indexOf("zh") === 0) pick = "zh";
            else if (l.indexOf("ja") === 0) pick = "ja";
            else if (l.indexOf("en") === 0) pick = "en";
          }
          pick = pick || "en";
        }
        if (pick === "en" || pick === "ja") location.replace(pick + "/" + location.hash);
      } catch (e) {}
    })();
  </script>
"""


def esc(s: str) -> str:
    return html.escape(s, quote=True)


def text_of(fragment: str) -> str:
    """HTML 片段 → 纯文本，给 JSON-LD 与 meta 用。"""
    return html.unescape(re.sub(r"<[^>]+>", "", fragment)).strip()


def load(code: str) -> dict:
    path = Path(__file__).parent / "content" / f"{code}.py"
    spec = importlib.util.spec_from_file_location(f"content_{code}", path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module.C


# ─────────────────────────── 首页 ───────────────────────────


def lang_links(code: str, up: str, sub: str = "") -> str:
    """除当前语言外的其他语言链接。up 是当前页面到网站根目录的相对路径，sub 是子页面（如 server/）。"""
    out = []
    for other, (path, hl, _, name) in LANGS.items():
        if other == code:
            continue
        out.append(f'<a href="{up}{path + sub or "./"}" hreflang="{hl}" lang="{hl}" data-lang="{other}">{name}</a>')
    return "\n        ".join(out)


def footer(c: dict, code: str, up: str, sub: str = "") -> str:
    langs = []
    for other, (path, hl, _, name) in LANGS.items():
        cur = ' aria-current="page"' if other == code else ""
        langs.append(f'<a href="{up}{path + sub or "./"}" hreflang="{hl}" lang="{hl}" data-lang="{other}"{cur}>{name}</a>')
    links = "".join(f'\n        <a href="{href}">{label}</a>' for label, href in c["footer_links"])
    return f"""  <footer class="footer">
    <div class="wrap">
      <p>{c["footer_license"]}</p>
      <nav aria-label="{c["footer_label"]}">{links}
      </nav>
      <nav aria-label="{c["lang_label"]}">
        {(chr(10) + "        ").join(langs)}
      </nav>
    </div>
  </footer>"""


def cta(c: dict) -> str:
    return f"""
        <div class="cta">
          <a class="button" href="{DL}CopySync.dmg">{DOWN}{c["dl_mac"]}</a>
          <a class="button quiet" href="{DL}CopySync-Setup.exe">{DOWN}{c["dl_win"]}</a>
        </div>
        <p class="note">{c["dl_note"]}</p>"""


def jsonld(c: dict, code: str, path: str, shots: str) -> str:
    data = {
        "@context": "https://schema.org",
        "@graph": [
            {
                "@type": "SoftwareApplication",
                "@id": SITE + "#app",
                "name": "CopySync",
                "description": c["ld_description"],
                "applicationCategory": "UtilitiesApplication",
                "operatingSystem": c["ld_os"],
                "softwareVersion": VERSION,
                "datePublished": DATE,
                "inLanguage": LANGS[code][1],
                "license": "https://www.gnu.org/licenses/agpl-3.0.html",
                "isAccessibleForFree": True,
                "offers": {"@type": "Offer", "price": "0", "priceCurrency": c["ld_currency"]},
                "downloadUrl": REPO + "/releases/latest",
                "codeRepository": REPO,
                "image": SITE + "assets/icon.png",
                "screenshot": SITE + shots + "history-dark.png",
                "featureList": c["ld_features"],
                "author": {"@type": "Person", "name": "baiyuze", "url": "https://github.com/baiyuze"},
            },
            {
                "@type": "FAQPage",
                "mainEntity": [
                    {"@type": "Question", "name": q,
                     "acceptedAnswer": {"@type": "Answer", "text": " ".join(text_of(p) for p in answer)}}
                    for q, answer in c["faq"]
                ],
            },
        ],
    }
    body = json.dumps(data, ensure_ascii=False, indent=2)
    return "\n".join("  " + line for line in body.split("\n"))


def home(code: str) -> str:
    c = load(code)
    path, hl, locale, _ = LANGS[code]
    up = "../" * path.count("/")
    a = up + "assets/"
    shots = "assets/screenshots/" + ("" if code == "zh" else code + "/")
    alternates = "\n".join(
        f'  <link rel="alternate" hreflang="{h}" href="{SITE}{p}">' for p, h, _, _ in LANGS.values())
    others = [l for k, (_, _, l, _) in LANGS.items() if k != code]

    ours = ' class="ours"'
    compare_head = "".join(
        f'<th scope="col"{ours if i == 0 else ""}>{h}</th>' for i, h in enumerate(c["compare_cols"]))
    compare_rows = "\n".join(
        f'              <tr><th scope="row">{label}</th>' + "".join(
            f'<td class="{cls}{" ours" if i == 0 else ""}">{val}</td>' for i, (cls, val) in enumerate(cells)) + "</tr>"
        for label, cells in c["compare_rows"])
    flow = "\n".join(f"          <li><h3>{h}</h3><p>{p}</p></li>" for h, p in c["flow"])
    uplink = "\n".join(f"              <li><strong>{s}</strong>{t}</li>" for s, t in c["uplink_points"])
    lab = "\n".join(f"                <tr><td>{s}</td><td>{r}</td></tr>" for s, r in c["lab_rows"])
    security = "\n".join(f"            <li><strong>{s}</strong>{t}</li>" for s, t in c["security_points"])
    fallbacks = "\n".join(f"          <div><dt>{t}</dt><dd>{d}</dd></div>" for t, d in c["fallbacks"])
    mac_steps = "\n".join(f"              <li>{s}</li>" for s in c["mac_steps"])
    win_steps = "\n".join(f"              <li>{s}</li>" for s in c["win_steps"])
    ports = "\n".join(f"                <tr><td>{p}</td><td>{u}</td></tr>" for p, u in c["ports"])
    specs = "\n".join(f'            <tr><th scope="row">{k}</th><td>{v}</td></tr>' for k, v in c["specs"])
    faq = "\n".join(
        f"""          <details>
            <summary>{q}</summary>
            <div class="answer">{"".join(f"<p>{p}</p>" for p in answer)}</div>
          </details>""" for q, answer in c["faq"])
    d = c["diagram"]

    return f"""<!doctype html>
<html lang="{hl}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{c["title"]}</title>
  <meta name="description" content="{esc(c["description"])}">
  <link rel="canonical" href="{SITE}{path}">
{alternates}
  <link rel="alternate" hreflang="x-default" href="{SITE}">
  <meta name="theme-color" content="#141416">
  <meta property="og:type" content="website">
  <meta property="og:site_name" content="CopySync">
  <meta property="og:title" content="{esc(c["og_title"])}">
  <meta property="og:description" content="{esc(c["og_description"])}">
  <meta property="og:url" content="{SITE}{path}">
  <meta property="og:image" content="{SITE}assets/{c["og_image"]}">
  <meta property="og:image:width" content="1200">
  <meta property="og:image:height" content="630">
  <meta property="og:locale" content="{locale}">
{"".join(f'  <meta property="og:locale:alternate" content="{l}">{chr(10)}' for l in others)}  <meta name="twitter:card" content="summary_large_image">
  <link rel="icon" href="{a}favicon.svg" type="image/svg+xml">
  <link rel="icon" href="{a}favicon-32.png" sizes="32x32" type="image/png">
  <link rel="apple-touch-icon" href="{a}apple-touch-icon.png">
  <link rel="stylesheet" href="{a}site.css">
{REDIRECT if code == "zh" else ""}  <script type="application/ld+json">
{jsonld(c, code, path, shots)}
  </script>
</head>
<body>
  <a class="skip" href="#main">{c["skip"]}</a>

  <header class="bar">
    <div class="wrap">
      <a class="brand" href="./" aria-label="{c["home_label"]}"><img src="{a}icon.svg" alt="" width="26" height="26">CopySync</a>
      <nav aria-label="{c["nav_label"]}">
        <a class="wide" href="#how">{c["nav_how"]}</a>
        <a class="wide" href="#install">{c["nav_install"]}</a>
        <a class="wide" href="#faq">{c["nav_faq"]}</a>
        <a href="{GUIDE}">{c["nav_server"]}</a>
        <a class="wide" href="{REPO}">GitHub</a>
        {lang_links(code, up)}
      </nav>
    </div>
  </header>

  <main id="main">
    <section class="hero">
      <div class="wrap">
        <h1 class="display">{c["h1"]}</h1>
        <p class="lede">{c["lede"]}</p>{cta(c)}
        <div class="stage" data-play role="img" aria-label="{esc(c["stage_label"])}">
          <figure class="device mac">
            <div class="screen"><div class="desk">
              <div class="file"><span class="doc">{LINES}</span><span class="name">{c["stage_file"]}</span></div>
              <div class="keys"><kbd>⌘</kbd><kbd>C</kbd></div>
            </div></div>
            <div class="base"></div>
            <figcaption>MacBook</figcaption>
          </figure>
          <div></div>
          <figure class="device win">
            <div class="screen"><div class="desk">
              <div class="file"><span class="doc">{LINES}</span><span class="name">{c["stage_file"]}</span></div>
              <div class="keys"><kbd>Ctrl</kbd><kbd>V</kbd></div>
              <div class="taskbar"><i></i><i></i><i></i><i></i><i></i></div>
            </div></div>
            <div class="base"></div>
            <figcaption>{c["stage_win"]}</figcaption>
          </figure>
          <svg class="arc" aria-hidden="true"><path d=""/></svg>
          <span class="arc-label" aria-hidden="true">{c["stage_arc"]}</span>
          <span class="sheet-fly" aria-hidden="true">{LINES}</span>
        </div>
        <button class="replay" type="button">{c["replay"]}</button>
      </div>
    </section>

    <section class="section" aria-labelledby="history-title">
      <div class="wrap">
        <h2 id="history-title">{c["history_h2"]}</h2>
        <p class="intro">{c["history_intro"]}</p>
        <figure class="shot">
          <img src="{up}{shots}history-dark.png" width="2112" height="1472" alt="{esc(c["history_alt"])}">
        </figure>
      </div>
    </section>

    <section class="section" aria-labelledby="flow-title">
      <div class="wrap">
        <h2 id="flow-title">{c["flow_h2"]}</h2>
        <p class="intro">{c["flow_intro"]}</p>
        <ol class="flow">
{flow}
        </ol>
      </div>
    </section>

    <section class="section" aria-labelledby="why-title">
      <div class="wrap">
        <h2 id="why-title">{c["why_h2"]}</h2>
        <p class="intro">{c["why_intro"]}</p>
        <div class="block table-scroll">
          <table class="compare">
            <caption class="skip">{c["compare_caption"]}</caption>
            <thead><tr><td></td>{compare_head}</tr></thead>
            <tbody>
{compare_rows}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <section class="section" id="how" aria-labelledby="how-title">
      <div class="wrap">
        <h2 id="how-title">{c["how_h2"]}</h2>
        <p class="intro">{c["how_intro"]}</p>
        <figure class="diagram">
          <svg viewBox="0 0 960 400" role="img" aria-labelledby="diagram-title diagram-desc">
            <title id="diagram-title">{d["title"]}</title>
            <desc id="diagram-desc">{d["desc"]}</desc>
            <rect class="d-box" x="360" y="24" width="240" height="96" rx="14"/>
            <text class="d-title" x="480" y="64" text-anchor="middle">{d["server"]}</text>
            <text class="d-sub" x="480" y="92" text-anchor="middle">{d["server_sub"]}</text>
            <path class="d-sig" d="M140 210V58H360"/>
            <path class="d-sig" d="M820 210V58H600"/>
            <path class="d-relay" d="M196 210V100H360"/>
            <path class="d-relay" d="M764 210V100H600"/>
            <text class="d-note" x="210" y="50">{d["sig_note"]}</text>
            <text class="d-note" x="210" y="128">{d["relay_note"]}</text>
            <rect class="d-box" x="40" y="210" width="240" height="120" rx="14"/>
            <text class="d-title" x="160" y="262" text-anchor="middle">Mac</text>
            <text class="d-sub" x="160" y="290" text-anchor="middle">{d["copy_here"]}</text>
            <rect class="d-box" x="680" y="210" width="240" height="120" rx="14"/>
            <text class="d-title" x="800" y="262" text-anchor="middle">{d["other"]}</text>
            <text class="d-sub" x="800" y="290" text-anchor="middle">{d["paste_here"]}</text>
            <path class="d-direct" d="M280 270H680"/>
            <text class="d-label" x="480" y="256" text-anchor="middle">{d["direct"]}</text>
            <text class="d-note" x="480" y="298" text-anchor="middle">{d["direct_note"]}</text>
            <g class="d-legend" transform="translate(40 372)">
              <path class="d-direct" d="M0 0H28"/><text x="38" y="5">{d["legend_direct"]}</text>
              <path class="d-relay" d="M210 0H238"/><text x="248" y="5">{d["legend_relay"]}</text>
              <path class="d-sig" d="M470 0H498"/><text x="508" y="5">{d["legend_sig"]}</text>
            </g>
          </svg>
          <ul class="points diagram-text">
            <li><strong>{d["direct"]}</strong> {d["direct_note"]}</li>
            <li><strong>{d["server"]}</strong> {d["server_sub"]}</li>
            <li>{d["relay_note"]}</li>
          </ul>
        </figure>
      </div>
    </section>

    <section class="section" id="multi-uplink" aria-labelledby="multi-title">
      <div class="wrap">
        <h2 id="multi-title">{c["uplink_h2"]}</h2>
        <p class="intro">{c["uplink_intro"]}</p>
        <div class="cols">
          <div>
            <h3>{c["uplink_changed"]}</h3>
            <ul class="points">
{uplink}
            </ul>
          </div>
          <div>
            <h3>{c["lab_h3"]}</h3>
            <p class="small">{c["lab_intro"]}</p>
            <div class="table-scroll">
              <table>
                <thead><tr><th scope="col">{c["lab_cols"][0]}</th><th scope="col">{c["lab_cols"][1]}</th></tr></thead>
                <tbody>
{lab}
                </tbody>
              </table>
            </div>
            <p class="small">{c["uplink_measured"]}</p>
            <p class="small">{c["uplink_doc"].format(doc=up + "design/nat-traversal.html")}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="section" aria-labelledby="security-title">
      <div class="wrap">
        <h2 id="security-title">{c["security_h2"]}</h2>
        <div class="cols">
          <ul class="points">
{security}
          </ul>
          <figure class="figure-card">
            <img src="{up}{shots}card-verify-dark.png" width="1120" height="800" loading="lazy" alt="{esc(c["verify_alt"])}">
          </figure>
        </div>
      </div>
    </section>

    <section class="section" aria-labelledby="fallback-title">
      <div class="wrap">
        <h2 id="fallback-title">{c["fallback_h2"]}</h2>
        <p class="intro">{c["fallback_intro"]}</p>
        <dl class="fallbacks">
{fallbacks}
        </dl>
      </div>
    </section>

    <section class="section" id="install" aria-labelledby="install-title">
      <div class="wrap">
        <h2 id="install-title">{c["install_h2"]}</h2>
        <p class="intro">{c["install_intro"]}</p>
        <div class="install">
          <div>
            <h3>{LAPTOP}Mac</h3>
            <ol>
{mac_steps}
            </ol>
            <a class="button" href="{DL}CopySync.dmg">{DOWN}{c["dl_mac"]}</a>
          </div>
          <div>
            <h3>{DESKTOP}Windows</h3>
            <ol>
{win_steps}
            </ol>
            <a class="button" href="{DL}CopySync-Setup.exe">{DOWN}{c["dl_win"]}</a>
          </div>
          <div>
            <h3>{SERVER}{c["server_h3"]}</h3>
            <p class="small">{c["server_intro"]}</p>
            <div class="code">
              <button class="copy" type="button" data-done="{c["copied"]}">{c["copy"]}</button>
<pre><code>curl -LO {DL}copysync-server-linux-amd64.tar.gz
tar xzf copysync-server-linux-amd64.tar.gz
cd copysync-server-linux-amd64
sudo ./install.sh {c["public_ip"]}</code></pre>
            </div>
            <p class="small">{c["server_after"]}</p>
            <table>
              <thead><tr><th scope="col">{c["ports_cols"][0]}</th><th scope="col">{c["ports_cols"][1]}</th></tr></thead>
              <tbody>
{ports}
              </tbody>
            </table>
            <p class="small">{c["server_docker"]}</p>
            <a class="button quiet" href="{GUIDE}">{c["server_guide"]}</a>
          </div>
        </div>
        <p class="pairing">{c["pairing"]}</p>
      </div>
    </section>

    <section class="section" aria-labelledby="specs-title">
      <div class="wrap">
        <h2 id="specs-title">{c["specs_h2"]}</h2>
        <div class="table-scroll">
          <table class="specs">
            <tbody>
{specs}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <section class="section" id="faq" aria-labelledby="faq-title">
      <div class="wrap">
        <h2 id="faq-title">{c["faq_h2"]}</h2>
        <div class="faq">
{faq}
        </div>
      </div>
    </section>

    <section class="closing" aria-labelledby="closing-title">
      <div class="wrap">
        <h2 class="display" id="closing-title">{c["closing"]}</h2>{cta(c)}
        <a class="source" href="{REPO}">{c["source"]}</a>
      </div>
    </section>
  </main>

{footer(c, code, up)}

  <script src="{a}site.js" defer></script>
</body>
</html>
"""


# ─────────────────────────── 在局域网里部署服务器 ───────────────────────────

# 教程里的示例地址。页面上的终端输出照安装脚本的真实输出抄写，只把地址换成它
LAN_IP = "192.168.1.20"

# 安装脚本在中文系统上输出中文，其他系统输出英文；日文页面也是英文输出。
# 改了 server/deploy/ 下脚本的输出，这里跟着改。
TRANSCRIPTS = {
    "zh": {
        "win": f"""> 停止旧版本
> 安装到 C:\\Program Files\\CopySync Server
> 注册服务 CopySyncServer（中转地址 {LAN_IP}）
> 在 Windows 防火墙里放行
> 启动
  √ 已启动

安装完成：copysync-server {VERSION}，作为 Windows 服务开机自动运行
在每台电脑的 CopySync「设置 → 信令服务器地址」填：

    ws://{LAN_IP}:8787/signal

  程序   C:\\Program Files\\CopySync Server
  日志   C:\\ProgramData\\CopySync Server\\server.log
  重启   Restart-Service CopySyncServer（管理员 PowerShell）
  卸载   双击 uninstall.cmd""",
        "mac": f"""▶ 停止旧版本
▶ 安装到 /Users/me/Library/Application Support/CopySync Server
▶ 写入 LaunchAgent（中转地址 {LAN_IP}）
▶ 启动
  ✓ 已启动

安装完成：copysync-server {VERSION}，登录这台 Mac 后自动运行
在每台电脑的 CopySync「设置 → 信令服务器地址」填：

    ws://{LAN_IP}:8787/signal

  日志   /Users/me/Library/Logs/CopySync/server.log
  重启   launchctl kickstart -k gui/501/com.copysync.server
  卸载   ./install.sh uninstall""",
        "linux": f"""▶ 安装二进制到 /usr/local/bin
▶ 写入 systemd 单元
▶ 生成配置 /etc/copysync/server.env（中转地址 {LAN_IP}）
▶ 启动
  ✓ 已启动

安装完成：copysync-server {VERSION}
在每台电脑的 CopySync「设置 → 信令服务器地址」填：

    ws://{LAN_IP}:8787/signal

  配置   /etc/copysync/server.env
  日志   journalctl -u copysync-server -f
  重启   systemctl restart copysync-server
  卸载   sudo ./install.sh uninstall

防火墙需放行：TCP 8787（信令）、UDP 3478（STUN 与 TURN）、UDP 32768-60999（中转端口，系统随机分配）""",
    },
    "en": {
        "win": f"""> Stopping the running version
> Installing to C:\\Program Files\\CopySync Server
> Registering the service CopySyncServer (relay address {LAN_IP})
> Allowing it through Windows Firewall
> Starting
  √ Running

Installed: copysync-server {VERSION}. It runs as a Windows service and starts with Windows.
In CopySync on each computer, set Settings > Signaling server to:

    ws://{LAN_IP}:8787/signal

  Program    C:\\Program Files\\CopySync Server
  Logs       C:\\ProgramData\\CopySync Server\\server.log
  Restart    Restart-Service CopySyncServer (in an administrator PowerShell)
  Uninstall  double-click uninstall.cmd""",
        "mac": f"""▶ Stopping the running version
▶ Installing to /Users/me/Library/Application Support/CopySync Server
▶ Writing the LaunchAgent (relay address {LAN_IP})
▶ Starting
  ✓ Running

Installed: copysync-server {VERSION}. It starts whenever you log in to this Mac.
In CopySync on each computer, set Settings → Signaling server to:

    ws://{LAN_IP}:8787/signal

  Logs       /Users/me/Library/Logs/CopySync/server.log
  Restart    launchctl kickstart -k gui/501/com.copysync.server
  Uninstall  ./install.sh uninstall""",
        "linux": f"""▶ Installing the binary to /usr/local/bin
▶ Writing the systemd unit
▶ Writing /etc/copysync/server.env (relay address {LAN_IP})
▶ Starting
  ✓ Running

Installed: copysync-server {VERSION}
In CopySync on each computer, set Settings → Signaling server to:

    ws://{LAN_IP}:8787/signal

  Config     /etc/copysync/server.env
  Logs       journalctl -u copysync-server -f
  Restart    systemctl restart copysync-server
  Uninstall  sudo ./install.sh uninstall

Open in your firewall: TCP 8787 (signaling), UDP 3478 (STUN and TURN), UDP 32768-60999 (relay ports, assigned by the OS)""",
    },
}

WINDOWS_LOGO = ('<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">'
                '<path d="M3 5.5 10.5 4.5v7H3zM11.5 4.4 21 3v8.5h-9.5zM3 12.5h7.5v7L3 18.5zM11.5 12.5H21V21l-9.5-1.4z"/></svg>')


def transcript(text: str, prompt: str, command: str) -> str:
    """终端窗口里的内容：先是输入的命令，再是脚本的输出，客户端要填的地址用蓝色标出来。"""
    out = html.escape(text, quote=False)
    url = f"ws://{LAN_IP}:8787/signal"
    out = out.replace(url, f'<span class="t-url">{url}</span>')
    # 脚本每一步的提示行（Mac 与 Linux 是 ▶，Windows 是 >，转义后是 &gt;）
    out = re.sub(r"^((?:▶|&gt;) .*)$", r'<span class="t-step">\1</span>', out, flags=re.M)
    head = f'<span class="t-prompt">{html.escape(prompt)}</span> {html.escape(command)}\n' if command else ""
    return head + out


def term(kind: str, title: str, body: str, label: str) -> str:
    """一个终端窗口的示意图。kind 是 mac（红绿灯在左）或 win（窗口按钮在右）。"""
    if kind == "win":
        bar = (f'<div class="t-bar win"><span class="t-tab">{WINDOWS_LOGO}{title}</span>'
               '<span class="t-ctl" aria-hidden="true"><i></i><i></i><i></i></span></div>')
    else:
        bar = f'<div class="t-bar mac"><span class="t-lights" aria-hidden="true"><i></i><i></i><i></i></span><span class="t-title">{title}</span></div>'
    return f"""<figure class="term" role="img" aria-label="{esc(label)}">
                {bar}
<pre class="t-body">{body}</pre>
              </figure>"""


def code_block(c: dict, code: str) -> str:
    return f"""<div class="code">
                  <button class="copy" type="button" data-done="{c["copied"]}">{c["copy"]}</button>
<pre><code>{html.escape(code, quote=False)}</code></pre>
                </div>"""


def server_page(code: str) -> str:
    c = load(code)
    g = c["guide"]
    path, hl, locale, _ = LANGS[code]
    up = "../" * (path + GUIDE).count("/")
    home_href = up + (path or "./")
    a = up + "assets/"
    shots = up + "assets/screenshots/" + ("" if code == "zh" else code + "/")
    t = TRANSCRIPTS["zh" if code == "zh" else "en"]
    url = f"ws://{LAN_IP}:8787/signal"
    alternates = "\n".join(
        f'  <link rel="alternate" hreflang="{h}" href="{SITE}{p}{GUIDE}">' for p, h, _, _ in LANGS.values())

    def items(lst: list[str], indent: str = "                ") -> str:
        return "\n".join(f"{indent}<li>{x}</li>" for x in lst)

    def manage(rows: list[tuple[str, str]]) -> str:
        body = "\n".join(f'                  <tr><th scope="row">{k}</th><td>{v}</td></tr>' for k, v in rows)
        return f"""<div class="table-scroll">
                <table class="manage">
                  <caption>{g["manage_caption"]}</caption>
                  <tbody>
{body}
                  </tbody>
                </table>
              </div>"""

    linux_dl = f"""curl -LO {DL}copysync-server-linux-amd64.tar.gz
tar xzf copysync-server-linux-amd64.tar.gz
cd copysync-server-linux-amd64
sudo ./install.sh"""
    mac_dl = f"""curl -LO {DL}copysync-server-macos.tar.gz
tar xzf copysync-server-macos.tar.gz
cd copysync-server-macos
./install.sh"""

    w, m, l = g["win"], g["mac"], g["linux"]
    explorer = []
    for name, kind in [("copysync-server.exe", "exe"), ("install.cmd", "cmd"), ("install.ps1", "doc"),
                       ("README.md", "doc"), ("uninstall.cmd", "cmd")]:
        # 要双击的那个文件标出来
        mark = f'<em>{w["double_click"]}</em>' if name == "install.cmd" else ""
        cls = ' class="hl"' if mark else ""
        explorer.append(f'                  <li{cls}><span class="f-icon {kind}" aria-hidden="true"></span>'
                        f'<span class="f-name">{name}</span>{mark}</li>')
    explorer = "\n".join(explorer)
    prep = "\n".join(f"          <div><dt>{h}</dt><dd>{p}</dd></div>" for h, p in g["prep"])
    faq = "\n".join(
        f"""          <details>
            <summary>{q}</summary>
            <div class="answer">{"".join(f"<p>{p}</p>" for p in answer)}</div>
          </details>""" for q, answer in g["faq"])
    d = g["diagram"]
    toc = "\n".join(f'          <li><a href="#{anchor}">{label}</a></li>' for anchor, label in g["toc"])

    return f"""<!doctype html>
<html lang="{hl}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{g["title"]}</title>
  <meta name="description" content="{esc(g["description"])}">
  <link rel="canonical" href="{SITE}{path}{GUIDE}">
{alternates}
  <meta name="theme-color" content="#141416">
  <meta property="og:type" content="article">
  <meta property="og:site_name" content="CopySync">
  <meta property="og:title" content="{esc(g["title"])}">
  <meta property="og:description" content="{esc(g["description"])}">
  <meta property="og:url" content="{SITE}{path}{GUIDE}">
  <meta property="og:image" content="{SITE}assets/{c["og_image"]}">
  <meta property="og:locale" content="{locale}">
  <meta name="twitter:card" content="summary_large_image">
  <link rel="icon" href="{a}favicon.svg" type="image/svg+xml">
  <link rel="icon" href="{a}favicon-32.png" sizes="32x32" type="image/png">
  <link rel="apple-touch-icon" href="{a}apple-touch-icon.png">
  <link rel="stylesheet" href="{a}site.css">
</head>
<body>
  <a class="skip" href="#main">{c["skip"]}</a>

  <header class="bar">
    <div class="wrap">
      <a class="brand" href="{home_href}" aria-label="{c["home_label"]}"><img src="{a}icon.svg" alt="" width="26" height="26">CopySync</a>
      <nav aria-label="{c["nav_label"]}">
        <a class="wide" href="{home_href}#how">{c["nav_how"]}</a>
        <a class="wide" href="{home_href}#install">{c["nav_install"]}</a>
        <a href="./" aria-current="page">{c["nav_server"]}</a>
        <a class="wide" href="{REPO}">GitHub</a>
        {lang_links(code, up, GUIDE)}
      </nav>
    </div>
  </header>

  <main id="main">
    <section class="hero guide-hero">
      <div class="wrap">
        <p class="crumb"><a href="{home_href}">CopySync</a> / {c["nav_server"]}</p>
        <h1 class="display">{g["h1"]}</h1>
        <p class="lede">{g["lede"]}</p>
        <ol class="guide-toc">
{toc}
        </ol>
      </div>
    </section>

    <section class="section" id="overview" aria-labelledby="overview-title">
      <div class="wrap">
        <h2 id="overview-title">{g["overview_h2"]}</h2>
        <p class="intro">{g["overview_intro"]}</p>
        <figure class="diagram">
          <svg viewBox="0 0 960 440" role="img" aria-labelledby="lan-title lan-desc">
            <title id="lan-title">{d["title"]}</title>
            <desc id="lan-desc">{d["desc"]}</desc>
            <rect class="d-lan" x="8" y="8" width="944" height="384" rx="22"/>
            <text class="d-lan-label" x="32" y="42">{d["lan"]}</text>
            <rect class="d-box" x="340" y="64" width="280" height="104" rx="14"/>
            <text class="d-title" x="480" y="106" text-anchor="middle">{d["server"]}</text>
            <text class="d-ip" x="480" y="138" text-anchor="middle">{LAN_IP}</text>
            <path class="d-sig" d="M150 250V96H340"/>
            <path class="d-sig" d="M810 250V96H620"/>
            <path class="d-relay" d="M206 250V140H340"/>
            <path class="d-relay" d="M754 250V140H620"/>
            <text class="d-note" x="164" y="86">{d["sig_note"]}</text>
            <text class="d-note" x="220" y="200">{d["relay_note"]}</text>
            <rect class="d-box" x="40" y="250" width="240" height="112" rx="14"/>
            <text class="d-title" x="160" y="300" text-anchor="middle">Mac</text>
            <text class="d-sub" x="160" y="328" text-anchor="middle">{d["client"]}</text>
            <rect class="d-box" x="680" y="250" width="240" height="112" rx="14"/>
            <text class="d-title" x="800" y="300" text-anchor="middle">Windows</text>
            <text class="d-sub" x="800" y="328" text-anchor="middle">{d["client"]}</text>
            <path class="d-direct" d="M280 306H680"/>
            <text class="d-label" x="480" y="292" text-anchor="middle">{d["direct"]}</text>
            <text class="d-note" x="480" y="336" text-anchor="middle">{d["direct_note"]}</text>
            <g class="d-legend" transform="translate(40 424)">
              <path class="d-direct" d="M0 0H28"/><text x="38" y="5">{d["legend_direct"]}</text>
              <path class="d-relay" d="M250 0H278"/><text x="288" y="5">{d["legend_relay"]}</text>
              <path class="d-sig" d="M560 0H588"/><text x="598" y="5">{d["legend_sig"]}</text>
            </g>
          </svg>
        </figure>
        <ul class="points ports">
{items(g["overview_points"], "          ")}
        </ul>
      </div>
    </section>

    <section class="section" id="prepare" aria-labelledby="prepare-title">
      <div class="wrap">
        <h2 id="prepare-title">{g["prep_h2"]}</h2>
        <dl class="fallbacks">
{prep}
        </dl>
      </div>
    </section>

    <section class="section" id="install" aria-labelledby="install-title">
      <div class="wrap">
        <h2 id="install-title">{g["install_h2"]}</h2>
        <p class="intro">{g["install_intro"]}</p>
        <div class="os-tabs" data-tabs>
          <div class="tablist" role="tablist" aria-label="{g["os_label"]}" hidden>
            <button type="button" role="tab" id="tab-windows" aria-controls="windows">Windows</button>
            <button type="button" role="tab" id="tab-macos" aria-controls="macos">macOS</button>
            <button type="button" role="tab" id="tab-linux" aria-controls="linux">Linux</button>
          </div>

          <section class="os" id="windows" role="tabpanel" aria-labelledby="tab-windows">
            <h3 class="os-title">{DESKTOP}Windows</h3>
            <div class="os-grid">
              <div class="os-steps">
                <ol class="steps">
{items(w["steps"])}
                </ol>
                <a class="button" href="{DL}copysync-server-windows-amd64.zip">{DOWN}{w["download"]}</a>
                <ul class="points">
{items(w["notes"])}
                </ul>
              </div>
              <div class="os-figures">
                <figure class="explorer" role="img" aria-label="{esc(w["explorer_label"])}">
                  <div class="t-bar win"><span class="t-tab folder">copysync-server-windows-amd64</span><span class="t-ctl" aria-hidden="true"><i></i><i></i><i></i></span></div>
                  <ul class="files">
{explorer}
                  </ul>
                </figure>
                {term("win", w["term_title"], transcript(t["win"], "", ""), w["term_label"])}
              </div>
            </div>
            {manage(w["manage"])}
          </section>

          <section class="os" id="macos" role="tabpanel" aria-labelledby="tab-macos">
            <h3 class="os-title">{LAPTOP}macOS</h3>
            <div class="os-grid">
              <div class="os-steps">
                <ol class="steps">
                  <li>{m["steps"][0]}
                {code_block(c, mac_dl)}</li>
{items(m["steps"][1:])}
                </ol>
                <ul class="points">
{items(m["notes"])}
                </ul>
              </div>
              <div class="os-figures">
                {term("mac", "copysync-server-macos — zsh", transcript(t["mac"], "%", "./install.sh"), m["term_label"])}
              </div>
            </div>
            {manage(m["manage"])}
          </section>

          <section class="os" id="linux" role="tabpanel" aria-labelledby="tab-linux">
            <h3 class="os-title">{SERVER}Linux</h3>
            <div class="os-grid">
              <div class="os-steps">
                <ol class="steps">
                  <li>{l["steps"][0]}
                {code_block(c, linux_dl)}</li>
{items(l["steps"][1:])}
                </ol>
                <ul class="points">
{items(l["notes"])}
                </ul>
              </div>
              <div class="os-figures">
                {term("mac", f"ssh {LAN_IP}", transcript(t["linux"], "$", "sudo ./install.sh"), l["term_label"])}
              </div>
            </div>
            {manage(l["manage"])}
          </section>
        </div>
      </div>
    </section>

    <section class="section" id="check" aria-labelledby="check-title">
      <div class="wrap">
        <h2 id="check-title">{g["check_h2"]}</h2>
        <p class="intro">{g["check_intro"]}</p>
        <figure class="browser" role="img" aria-label="{esc(g["check_label"])}">
          <div class="b-bar"><span class="t-lights" aria-hidden="true"><i></i><i></i><i></i></span><span class="b-url">http://{LAN_IP}:8787/healthz</span></div>
          <pre class="b-body">{{"status":"ok","version":"{VERSION}"}}</pre>
        </figure>
        <p class="note">{g["check_after"]}</p>
      </div>
    </section>

    <section class="section" id="clients" aria-labelledby="clients-title">
      <div class="wrap">
        <h2 id="clients-title">{g["clients_h2"]}</h2>
        <ol class="flow">
{chr(10).join(f"          <li><h3>{h}</h3><p>{p}</p></li>" for h, p in g["clients_steps"])}
        </ol>
        <figure class="shot">
          <img src="{shots}settings-lan-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(g["clients_alt"])}">
        </figure>
        <p class="note">{g["clients_after"].format(home=home_href)}</p>
      </div>
    </section>

    <section class="section" id="faq" aria-labelledby="faq-title">
      <div class="wrap">
        <h2 id="faq-title">{g["faq_h2"]}</h2>
        <div class="faq">
{faq}
        </div>
      </div>
    </section>

    <section class="closing" aria-labelledby="closing-title">
      <div class="wrap">
        <h2 class="display" id="closing-title">{g["closing"]}</h2>{cta(c)}
        <a class="source" href="{REPO}/blob/main/server/deploy/README.md">{g["source"]}</a>
      </div>
    </section>
  </main>

{footer(c, code, up, GUIDE)}

  <script src="{a}site.js" defer></script>
</body>
</html>
"""


# ─────────────────────────── Markdown → HTML ───────────────────────────
# 只实现技术方案用到的那部分：标题、段落、列表（可嵌套）、引用、代码块、表格、分隔线，
# 行内的代码、粗体、链接。


def doc_link(target: str) -> str:
    if re.match(r"[a-z]+:", target) or target.startswith("#"):
        return target
    name, _, anchor = target.partition("#")
    if name in DESIGN_DOCS:
        return DESIGN_DOCS[name] + ".html" + ("#" + anchor if anchor else "")
    # 指向仓库里其他文件的相对链接：换成 GitHub 上的地址
    return f"{REPO}/blob/main/design/{target}"


def inline(s: str) -> str:
    codes = []

    def keep(m):
        codes.append("<code>" + html.escape(m.group(1), quote=False) + "</code>")
        return f"\x00{len(codes) - 1}\x00"

    s = re.sub(r"`([^`]+)`", keep, s)
    s = html.escape(s, quote=False)
    s = re.sub(r"\*\*(.+?)\*\*", r"<strong>\1</strong>", s)
    s = re.sub(r"\[([^\]]+)\]\(([^)\s]+)\)",
               lambda m: f'<a href="{esc(doc_link(html.unescape(m.group(2))))}">{m.group(1)}</a>', s)
    return re.sub(r"\x00(\d+)\x00", lambda m: codes[int(m.group(1))], s)


def join_lines(lines: list[str]) -> str:
    # 中文段落折行处直接相连，西文之间补一个空格
    out = ""
    for line in lines:
        line = line.strip()
        if out and not (re.search(r"[　-鿿＀-￯]$", out) or re.match(r"[　-鿿＀-￯]", line)):
            out += " "
        out += line
    return out


LIST_ITEM = re.compile(r"^(\s*)([-*]|\d+\.)\s+(.*)$")


def render_list(lines: list[str]) -> str:
    """lines 是一整块列表（含续行）。按缩进建嵌套列表。"""
    items = []  # (缩进, 是否有序, 起始编号, 文本行)
    for line in lines:
        m = LIST_ITEM.match(line)
        if m:
            items.append([len(m.group(1)), m.group(2)[-1] == ".", m.group(2).rstrip("."), [m.group(3)]])
        elif items:
            items[-1][3].append(line)

    def build(start: int, indent: int) -> tuple[str, int]:
        ordered = items[start][1]
        tag = "ol" if ordered else "ul"
        first = items[start][2]
        attrs = f' start="{first}"' if ordered and first != "1" else ""
        parts = [f"<{tag}{attrs}>"]
        i = start
        while i < len(items) and items[i][0] >= indent:
            if items[i][0] > indent:
                sub, i = build(i, items[i][0])
                parts[-1] = parts[-1][: -len("</li>")] + sub + "</li>"
                continue
            parts.append("<li>" + inline(join_lines(items[i][3])) + "</li>")
            i += 1
        parts.append(f"</{tag}>")
        return "".join(parts), i

    out, i = [], 0
    while i < len(items):
        block, i = build(i, items[i][0])
        out.append(block)
    return "\n".join(out)


def slug(text: str) -> str:
    return re.sub(r"[\s，。、：:（）()「」`*.]+", "-", text_of(text)).strip("-").lower()


def markdown(src: str) -> tuple[str, str, list[tuple[str, str]]]:
    """返回 (标题, 正文 HTML, 二级标题列表)。"""
    lines = src.split("\n")
    out, toc, title = [], [], ""
    i = 0
    while i < len(lines):
        line = lines[i]
        if not line.strip():
            i += 1
            continue
        if line.startswith("```"):
            j = i + 1
            while j < len(lines) and not lines[j].startswith("```"):
                j += 1
            code = "\n".join(lines[i + 1 : j])
            out.append("<pre><code>" + html.escape(code, quote=False) + "</code></pre>")
            i = j + 1
            continue
        m = re.match(r"^(#{1,4})\s+(.*)$", line)
        if m:
            level, text = len(m.group(1)), m.group(2).strip()
            if level == 1 and not title:
                title = text
                i += 1
                continue
            hid = slug(text)
            if level == 2:
                toc.append((hid, text))
            out.append(f'<h{level} id="{esc(hid)}">{inline(text)}</h{level}>')
            i += 1
            continue
        if re.match(r"^(-{3,}|\*{3,})\s*$", line):
            out.append("<hr>")
            i += 1
            continue
        if line.startswith(">"):
            quote = []
            while i < len(lines) and lines[i].startswith(">"):
                quote.append(lines[i].lstrip(">").strip())
                i += 1
            out.append("<blockquote>" + "".join(f"<p>{inline(q)}</p>" for q in quote if q) + "</blockquote>")
            continue
        if line.startswith("|") and i + 1 < len(lines) and re.match(r"^\|[\s:|-]+\|\s*$", lines[i + 1]):
            def cells(row):
                return [c.strip() for c in row.strip().strip("|").split("|")]
            head = cells(line)
            rows = []
            i += 2
            while i < len(lines) and lines[i].startswith("|"):
                rows.append(cells(lines[i]))
                i += 1
            t = ["<div class=\"table-scroll\"><table><thead><tr>"]
            t += [f'<th scope="col">{inline(h)}</th>' for h in head]
            t.append("</tr></thead><tbody>")
            for r in rows:
                t.append("<tr>" + "".join(f"<td>{inline(c)}</td>" for c in r) + "</tr>")
            t.append("</tbody></table></div>")
            out.append("".join(t))
            continue
        if LIST_ITEM.match(line):
            block = []
            while i < len(lines) and lines[i].strip() and (LIST_ITEM.match(lines[i]) or lines[i].startswith(" ")):
                block.append(lines[i])
                i += 1
            out.append(render_list(block))
            continue
        para = []
        while i < len(lines) and lines[i].strip() and not (
            lines[i].startswith(("#", ">", "```", "|")) or LIST_ITEM.match(lines[i])
        ):
            para.append(lines[i])
            i += 1
        out.append("<p>" + inline(join_lines(para)) + "</p>")
    return title, "\n".join(out), toc


def design_page(md_name: str, page: str) -> str:
    title, body, toc = markdown((ROOT / "design" / md_name).read_text(encoding="utf-8"))
    first_quote = re.search(r"<blockquote><p>(.*?)</p>", body)
    description = text_of(first_quote.group(1)) if first_quote else title
    toc_html = ""
    if len(toc) > 2:
        # 标题自带「一、二、」编号，目录不再加序号
        toc_html = ('<nav class="toc" aria-label="目录"><strong>目录</strong><ul>'
                    + "".join(f'<li><a href="#{esc(h)}">{inline(t)}</a></li>' for h, t in toc)
                    + "</ul></nav>")
    # 目录放在开头的状态说明之后
    if toc_html and body.startswith("<blockquote>"):
        cut = body.index("</blockquote>") + len("</blockquote>")
        body = body[:cut] + "\n" + toc_html + body[cut:]
    else:
        body = toc_html + body
    zh = load("zh")
    return f"""<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{esc(text_of(inline(title)))}：CopySync 技术方案</title>
  <meta name="description" content="{esc(description)}">
  <link rel="canonical" href="{SITE}design/{page}.html">
  <meta name="theme-color" content="#141416">
  <link rel="icon" href="../assets/favicon.svg" type="image/svg+xml">
  <link rel="stylesheet" href="../assets/site.css">
</head>
<body>
  <a class="skip" href="#main">跳到正文</a>
  <header class="bar">
    <div class="wrap">
      <a class="brand" href="../" aria-label="CopySync 首页"><img src="../assets/icon.svg" alt="" width="26" height="26">CopySync</a>
      <nav aria-label="页面导航">
        <a href="../{GUIDE}">{zh["nav_server"]}</a>
        <a class="wide" href="{REPO}/blob/main/design/{md_name}">在 GitHub 上查看</a>
        <a class="wide" href="{REPO}">GitHub</a>
      </nav>
    </div>
  </header>
  <main id="main">
    <article class="article">
      <p class="crumb"><a href="../">CopySync</a> 的技术方案</p>
      <h1>{inline(title)}</h1>
{body}
    </article>
  </main>

{footer(zh, "zh", "../")}
</body>
</html>
"""


# ─────────────────────────── 其他页面 ───────────────────────────


def not_found() -> str:
    return """<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="robots" content="noindex">
  <title>页面不存在 / Page not found：CopySync</title>
  <meta name="theme-color" content="#141416">
  <link rel="icon" href="/copysync/assets/favicon.svg" type="image/svg+xml">
  <link rel="stylesheet" href="/copysync/assets/site.css">
</head>
<body>
  <main class="lost">
    <div class="wrap">
      <h1 class="display">404</h1>
      <p class="lede">这个页面不存在。<br><span lang="en">This page doesn't exist.</span><br><span lang="ja">このページは存在しません。</span></p>
      <div class="cta">
        <a class="button" href="/copysync/">回到首页</a>
        <a class="button quiet" href="/copysync/en/" lang="en">English home</a>
        <a class="button quiet" href="/copysync/ja/" lang="ja">日本語のホーム</a>
      </div>
    </div>
  </main>
</body>
</html>
"""


def sitemap() -> str:
    alts = "".join(
        f'\n    <xhtml:link rel="alternate" hreflang="{h}" href="{SITE}{p}"/>' for p, h, _, _ in LANGS.values())
    alts += f'\n    <xhtml:link rel="alternate" hreflang="x-default" href="{SITE}"/>'
    urls = [f"  <url>\n    <loc>{SITE}{p}</loc>\n    <lastmod>{DATE}</lastmod>{alts}\n  </url>"
            for p, _, _, _ in LANGS.values()]
    server_alts = "".join(
        f'\n    <xhtml:link rel="alternate" hreflang="{h}" href="{SITE}{p}{GUIDE}"/>' for p, h, _, _ in LANGS.values())
    urls += [f"  <url>\n    <loc>{SITE}{p}{GUIDE}</loc>\n    <lastmod>{DATE}</lastmod>{server_alts}\n  </url>"
             for p, _, _, _ in LANGS.values()]
    urls += [f"  <url>\n    <loc>{SITE}design/{page}.html</loc>\n    <lastmod>{DATE}</lastmod>\n  </url>"
             for page in DESIGN_DOCS.values()]
    return ('<?xml version="1.0" encoding="UTF-8"?>\n'
            '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"\n'
            '        xmlns:xhtml="http://www.w3.org/1999/xhtml">\n' + "\n".join(urls) + "\n</urlset>\n")


def write(path: Path, text: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text, encoding="utf-8")
    print("✓", path.relative_to(ROOT))


def main() -> None:
    for code, (path, _, _, _) in LANGS.items():
        write(DOCS / path / "index.html", home(code))
        write(DOCS / path / GUIDE / "index.html", server_page(code))
    for md_name, page in DESIGN_DOCS.items():
        write(DOCS / "design" / f"{page}.html", design_page(md_name, page))
    write(DOCS / "404.html", not_found())
    write(DOCS / "sitemap.xml", sitemap())


if __name__ == "__main__":
    main()

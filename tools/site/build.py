"""生成官网：docs/index.html（简体中文）、docs/zh-tw/ 与 docs/zh-hk/（繁体）、docs/en/、docs/ja/，每种语言下还有快速开始（start/）与三种连接方式的教程——
同一网络（lan/，不需要服务器）、公网服务器（cloud/）、Tailscale（tailscale/），
在局域网里自建服务器的教程（server/），以及 docs/design/ 下的技术方案页面。

    python3 tools/site/build.py

各语言共用同一套页面结构，文案在 tools/site/content/ 下每种语言一份。
技术方案页面由 design/*.md 转换而来，放在网站上，不用跳到 GitHub 去看。
改了文案、结构或技术方案后运行一次，把生成的文件一起提交。只用 Python 标准库。
"""

import html
import importlib.util
import json
import re
from pathlib import Path

from content.tour import localized as tour_copy
from content.privacy import COPY as PRIVACY_COPY
from content.official import COPY as OFFICIAL_COPY

ROOT = Path(__file__).resolve().parents[2]
DOCS = ROOT / "docs"
SITE = "https://pasteferry.cc/"
REPO = "https://github.com/baiyuze/pasteferry"
# 网站部署在独立域名的根目录，404 页面用绝对路径引用资源。
BASE = "/"
DL = REPO + "/releases/latest/download/"
VERSION = "2.1.0"
DATE = "2026-10-16"

# 语言代码 → (页面目录, html 的 lang, og:locale, 语言名)。顶栏的语言菜单按这个顺序列出。
# 繁体的两份文案由 zh.py 转换后逐条校订，改了 zh.py 要一起改
LANGS = {
    "zh": ("", "zh-CN", "zh_CN", "简体中文"),
    "zh-tw": ("zh-tw/", "zh-TW", "zh_TW", "繁體中文（台灣）"),
    "zh-hk": ("zh-hk/", "zh-HK", "zh_HK", "繁體中文（香港）"),
    "en": ("en/", "en", "en_US", "English"),
    "ja": ("ja/", "ja", "ja_JP", "日本語"),
}


def shots_dir(code: str) -> str:
    """截图目录：简体中文的在 screenshots/ 下，其他语言在 screenshots/<语言>/ 下，由 tools/screenshots/render.sh 生成。"""
    return "assets/screenshots/" + ("" if code == "zh" else code + "/")

# Cloudflare Web Analytics：不用 Cookie，只统计浏览量、来源、国家与设备。
# token 会出现在每个页面的源码里，本来就是公开的，只能用来上报，看不到统计数据。
# 统计原先绑定 GitHub Pages，新域名的统计需在 Cloudflare 后台配置。
ANALYTICS = ("  <script type=\"module\" src=\"https://static.cloudflareinsights.com/beacon.min.js\" "
             "data-cf-beacon='{\"token\": \"7ee226baf4a949a397c4bfd69a9e7461\"}'></script>\n")

# 教程页面，每种语言目录下各一份：快速开始，以及三种连接方式各一篇——
# 同一网络（不需要服务器）、公网服务器、Tailscale。三种方式地位相同，在快速开始里由用户自己选。
# 在局域网里自建服务器（server/）不再是一种推荐的方式，页面保留：两个网段互相访问不到时还用得上
START = "start/"
LAN = "lan/"
CLOUD = "cloud/"
GUIDE = "server/"
TAILSCALE = "tailscale/"
# 下载页：所有「下载」按钮都指向它，按系统、芯片选具体的安装包
DOWNLOAD = "download/"
TOUR = "tour/"
PRIVACY = "privacy/"
OFFICIAL = "official/"
SUBPAGES = [START, OFFICIAL, LAN, CLOUD, GUIDE, TAILSCALE, DOWNLOAD, TOUR, PRIVACY]
DEPLOY_PAGES = [LAN, CLOUD, GUIDE, TAILSCALE]

# 放到网站上的技术方案：文件名 → 网页名。其他 design/*.md 链接到 GitHub
DESIGN_DOCS = {"nat-traversal.md": "nat-traversal", "windows-client.md": "windows-client",
               "lan-mode.md": "lan-mode", "linux-client.md": "linux-client"}

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
GLOBE = ('<svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.4" aria-hidden="true">'
         '<circle cx="10" cy="10" r="7.25"/><path d="M2.75 10h14.5M10 2.75c-2.1 2-3.1 4.4-3.1 7.25s1 5.25 3.1 7.25'
         'M10 2.75c2.1 2 3.1 4.4 3.1 7.25s-1 5.25-3.1 7.25"/></svg>')
CHEVRON = ('<svg class="chev" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5" '
           'stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m3 4.5 3 3 3-3"/></svg>')

# 同一份页面也由 GitHub Pages 发布；在那里先跳到官网，再加载样式和其他脚本。
GITHUB_PAGES_REDIRECT = """  <script>
    (function () {
      if (location.hostname !== "baiyuze.github.io") return;
      var prefix = /^\\/(?:pasteferry|copysync)(?=\\/|$)/.exec(location.pathname);
      if (!prefix) return;
      var path = location.pathname.slice(prefix[0].length) || "/";
      location.replace(SITE_ORIGIN + path + location.search + location.hash);
    })();
  </script>
""".replace("SITE_ORIGIN", json.dumps(SITE.rstrip("/")))

# 首页按浏览器语言跳转：只在简体中文的首页（网站根目录）做，选过语言后以选的为准。
# 香港、澳门的繁体去香港版，其他繁体（台湾、只写了 zh-Hant 的）去台湾版
REDIRECT = """  <script>
    // 第一次打开时按浏览器语言跳到对应页面；在页面上选过语言后，以选的为准
    (function () {
      if (location.hostname === "baiyuze.github.io") return;
      try {
        var paths = PATHS;
        var pick = localStorage.getItem("pasteferry-lang");
        if (!paths.hasOwnProperty(pick)) {
          pick = "";
          if (/bot|crawl|spider|slurp|lighthouse/i.test(navigator.userAgent)) return;
          var list = navigator.languages || [navigator.language || ""];
          for (var i = 0; i < list.length && !pick; i++) {
            var l = String(list[i]).toLowerCase();
            if (/^zh-(hk|mo|hant-hk|hant-mo)\\b/.test(l)) pick = "zh-hk";
            else if (/^zh-(tw|hant)\\b/.test(l)) pick = "zh-tw";
            else if (l.indexOf("zh") === 0) pick = "zh";
            else if (l.indexOf("ja") === 0) pick = "ja";
            else if (l.indexOf("en") === 0) pick = "en";
          }
          pick = pick || "en";
        }
        if (paths[pick]) location.replace(paths[pick] + location.search + location.hash);
      } catch (e) {}
    })();
  </script>
""".replace("PATHS", json.dumps({k: v[0] for k, v in LANGS.items()}))


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
    c = module.C
    o = OFFICIAL_COPY[code]
    c["official"] = o
    c["plans"]["official"] = {
        "name": o["name"], "fit": o["fit"], "summary": o["summary"], "go": o["go"],
        "facts": [(o["routing_title"], o["routing"]), (o["quota_title"], o["quota"])],
        "steps": o["steps"], "address": o["privacy"],
    }
    return c


# ─────────────────────────── 首页 ───────────────────────────


def lang_menu(c: dict, code: str, up: str, sub: str = "") -> str:
    """顶栏右侧的语言菜单：平时只显示当前语言，点开列出全部，链接到同一页面的其他语言版本。
    用 <details> 做，没有脚本也能展开。up 是当前页面到网站根目录的相对路径，sub 是子页面（如 server/）。"""
    items = []
    for other, (path, hl, _, name) in LANGS.items():
        cur = ' aria-current="page"' if other == code else ""
        items.append(f'<li><a href="{up}{path + sub or "./"}" hreflang="{hl}" lang="{hl}" data-lang="{other}"{cur}>{name}</a></li>')
    name = LANGS[code][3]
    return f"""<details class="lang">
          <summary aria-label="{c["lang_label"]}：{name}">{GLOBE}<span>{name}</span>{CHEVRON}</summary>
          <ul>
            {(chr(10) + "            ").join(items)}
          </ul>
        </details>"""


def footer(c: dict, code: str, up: str, sub: str = "") -> str:
    langs = []
    for other, (path, hl, _, name) in LANGS.items():
        cur = ' aria-current="page"' if other == code else ""
        langs.append(f'<a href="{up}{path + sub or "./"}" hreflang="{hl}" lang="{hl}" data-lang="{other}"{cur}>{name}</a>')
    links = "".join(f'\n        <a href="{href}">{label}</a>' for label, href in c["footer_links"])
    links += f'\n        <a href="{up}{LANGS[code][0]}{PRIVACY}">{PRIVACY_COPY[code][0]}</a>'
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


def cta(c: dict, base: str, start: str = "", tour_label: str = "") -> str:
    """下载按钮，指向下载页：在那里按系统和芯片选安装包。base 是语言根目录的相对路径，
    start 是快速开始页面的地址，给了就在下面加一行链接。"""
    link = f'\n        <a class="start-link" href="{start}">{c["start_link"]}</a>' if start else ""
    tour = (f'\n          <a class="button quiet tour-link" href="{base}{TOUR}" data-tour-open aria-haspopup="dialog">'
            f'<span aria-hidden="true">▶</span> {esc(tour_label)} <span class="tour-duration">1:25</span></a>') if tour_label else ""
    return f"""
        <div class="cta">
          <a class="button" href="{base}{DOWNLOAD}">{DOWN}{c["dl_button"]}</a>{tour}
        </div>
        <p class="note">{c["dl_note"]}</p>{link}"""


def header(c: dict, code: str, up: str, sub: str = "") -> str:
    """顶栏。sub 是当前的教程页面（START、GUIDE 等），首页为空。窄屏只留不带 wide 的几项，手机上再收起带 mid 的。
    「连接方式」指向快速开始里三种方式的选择，在各篇教程里显示为当前所在的栏目。"""
    path = LANGS[code][0]
    if sub:
        home = up + (path or "./")
        base = up + path  # 语言根目录，教程页面互相链接用
    else:
        home, base = "", ""
    if sub == START:
        start = '<a href="./" aria-current="page">'
        deploy = '<a class="mid" href="#server">'
    else:
        start = f'<a href="{base}{START}">'
        current = ' aria-current="true"' if sub in DEPLOY_PAGES else ""
        deploy = f'<a class="mid" href="{base}{START}#server"{current}>'
    download = '<a href="./" aria-current="page">' if sub == DOWNLOAD else f'<a href="{base}{DOWNLOAD}">'

    return f"""  <header class="bar">
    <div class="wrap">
      <a class="brand" href="{home or "./"}" aria-label="{c["home_label"]}"><img src="{up}assets/icon.svg" alt="" width="26" height="26">{c["brand"]}</a>
      <nav aria-label="{c["nav_label"]}">
        {start}{c["nav_start"]}</a>
        <a class="wide" href="{home}#how">{c["nav_how"]}</a>
        <a class="wide" href="{home}#faq">{c["nav_faq"]}</a>
        {deploy}{c["nav_server"]}</a>
        {download}{c["nav_download"]}</a>
        <a class="wide" href="{REPO}">GitHub</a>
        {lang_menu(c, code, up, sub)}
      </nav>
    </div>
  </header>"""


# 三种连接方式的小插图：设备是笔记本与服务器的轮廓，虚线框是一个网络，蓝线是复制的内容走的路
def _laptop(x: int, y: int) -> str:
    return (f'<rect class="p-dev" x="{x + 5}" y="{y}" width="38" height="24" rx="3"/>'
            f'<path class="p-dev" d="M{x} {y + 29}h48"/>')


def _server(x: int, y: int) -> str:
    return (f'<rect class="p-dev" x="{x}" y="{y}" width="44" height="12" rx="3"/>'
            f'<rect class="p-dev" x="{x}" y="{y + 15}" width="44" height="12" rx="3"/>')


# 三种连接方式，按这个顺序出现在首页、快速开始与安装一节：(文案里的键, 教程页面, 插图)
PLANS = [
    ("official", OFFICIAL, _server(98, 8)
     + '<rect class="p-net" x="6" y="50" width="100" height="56" rx="10"/>'
     + '<rect class="p-net" x="134" y="50" width="100" height="56" rx="10"/>'
     + '<path class="p-sig" d="M56 64V22h42M184 64V22h-42"/>'
     + '<path class="p-data" d="M80 80h80"/>'
     + _laptop(32, 64) + _laptop(160, 64)),
    # 同一网络：一个网络里只有两台电脑，没有服务器
    ("lan", LAN, '<rect class="p-net" x="6" y="6" width="228" height="100" rx="14"/>'
     + '<path class="p-data" d="M82 52h76"/>'
     + _laptop(34, 40) + _laptop(158, 40)),
    # 公网服务器：服务器在两个网络之外，两边都直接连它
    ("cloud", CLOUD, _server(98, 8)
     + '<rect class="p-net" x="6" y="50" width="100" height="56" rx="10"/>'
     + '<rect class="p-net" x="134" y="50" width="100" height="56" rx="10"/>'
     + '<path class="p-sig" d="M56 64V22h42M184 64V22h-42"/>'
     + '<path class="p-data" d="M80 80h80"/>'
     + _laptop(32, 64) + _laptop(160, 64)),
    # Tailscale：家与公司是两个网络，外面一圈点线把它们连成一个虚拟局域网，不需要服务器
    ("ts", TAILSCALE, '<rect class="p-ts" x="2" y="2" width="236" height="108" rx="16"/>'
     '<rect class="p-net" x="12" y="12" width="98" height="88" rx="10"/>'
     '<rect class="p-net" x="130" y="12" width="98" height="88" rx="10"/>'
     + '<path class="p-data" d="M85 52h71"/>'
     + _laptop(37, 40) + _laptop(156, 40)),
]


def paths(c: dict, base: str) -> str:
    """首页上三种连接方式的卡片，样式相同、不分主次。base 是语言根目录（教程页面在它下面）。"""
    cards = []
    for key, page, art in PLANS:
        p = c["plans"][key]
        cards.append(f"""          <a class="path" href="{base}{page}">
            <svg viewBox="0 0 240 112" aria-hidden="true">{art}</svg>
            <span class="path-when">{p["fit"]}</span>
            <h3>{p["name"]}</h3>
            <p>{p["summary"]}</p>
            <span class="path-go">{p["go"]}</span>
          </a>""")
    return '<div class="paths">\n' + "\n".join(cards) + "\n        </div>"


def plan_tabs(c: dict, base: str) -> str:
    """快速开始第 2 步：三种连接方式分页，由用户自己选。每页是适合谁、要准备什么、要注意什么，
    精简的步骤，最后一句要点，以及完整教程的链接。"""
    links = {"official": base + OFFICIAL, "lan": base + LAN, "cloud": base + CLOUD, "server": base + GUIDE, "ts": base + TAILSCALE}
    tabs, panels = [], []
    for key, page, art in PLANS:
        p = c["plans"][key]
        tabs.append(f'            <button type="button" role="tab" id="tab-plan-{key}" aria-controls="plan-{key}">{p["name"]}</button>')
        facts = "\n".join(f"                  <div><dt>{k}</dt><dd>{v}</dd></div>" for k, v in p["facts"])
        steps = "\n".join(f"                  <li>{fill(x, links)}</li>" for x in p["steps"])
        panels.append(f"""          <section class="os" id="plan-{key}" role="tabpanel" aria-labelledby="tab-plan-{key}">
            <h3 class="os-title">{p["name"]}</h3>
            <div class="plan-grid">
              <div>
                <ol class="steps">
{steps}
                </ol>
                <p class="plan-address">{fill(p["address"], links)}</p>
                <a class="button quiet" href="{base}{page}">{p["go"]}</a>
              </div>
              <div class="plan-card">
                <svg viewBox="0 0 240 112" aria-hidden="true">{art}</svg>
                <dl>
{facts}
                </dl>
              </div>
            </div>
          </section>""")
    return f"""<div class="os-tabs" data-tabs="plan">
          <div class="tablist" role="tablist" aria-label="{c["plans_label"]}" hidden>
{chr(10).join(tabs)}
          </div>

{(chr(10) + chr(10)).join(panels)}
        </div>"""


def jsonld(c: dict, code: str, path: str, shots: str) -> str:
    data = {
        "@context": "https://schema.org",
        "@graph": [
            {
                "@type": "SoftwareApplication",
                "@id": SITE + "#app",
                "name": c["brand"],
                "alternateName": c["brand_alt"],
                "description": c["ld_description"],
                "applicationCategory": "UtilitiesApplication",
                "operatingSystem": c["ld_os"],
                "softwareVersion": VERSION,
                "datePublished": DATE,
                "inLanguage": LANGS[code][1],
                "license": "https://www.gnu.org/licenses/agpl-3.0.html",
                "isAccessibleForFree": True,
                "offers": {"@type": "Offer", "price": "0", "priceCurrency": c["ld_currency"]},
                "downloadUrl": SITE + LANGS[code][0] + DOWNLOAD,
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
    _, tour_ui, _ = tour_copy(code, c)
    path, hl, locale, _ = LANGS[code]
    up = "../" * path.count("/")
    a = up + "assets/"
    shots = shots_dir(code)
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
    linux_steps = "\n".join(f"              <li>{s}</li>" for s in c["linux_steps"])
    opts = "\n".join(
        f'              <li><a href="{page}">{c["plans"][key]["name"]}</a><span>{c["plans"][key]["fit"]}</span></li>'
        for key, page, _ in PLANS)
    specs = "\n".join(f'            <tr><th scope="row">{k}</th><td>{v}</td></tr>' for k, v in c["specs"])
    links = {"official": OFFICIAL, "lan": LAN, "cloud": CLOUD, "server": GUIDE, "ts": TAILSCALE, "start": START}
    faq = faq_list(c["faq"], links)
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
  <meta property="og:site_name" content="{c["brand"]}">
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

{header(c, code, up)}

  <main id="main">
    <section class="hero">
      <div class="wrap">
        <a class="release-note" href="{OFFICIAL}">{esc(c["official"]["badge"])}</a>
        <h1 class="display">{c["h1"]}</h1>
        <p class="lede">{c["lede"]}</p>{cta(c, "", START, tour_ui["watch"])}
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

    <section class="section" id="deploy" aria-labelledby="deploy-title">
      <div class="wrap">
        <h2 id="deploy-title">{c["paths_h2"]}</h2>
        <p class="intro">{c["paths_intro"]}</p>
        {paths(c, "")}
        <p class="note">{fill(c["paths_after"], links)}</p>
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
        <p class="intro">{fill(c["install_intro"], links)}</p>
        <div class="install">
          <div>
            <h3>{LAPTOP}Mac</h3>
            <ol>
{mac_steps}
            </ol>
            <a class="button" href="{DOWNLOAD}#macos">{DOWN}{c["dl_mac"]}</a>
          </div>
          <div>
            <h3>{DESKTOP}Windows</h3>
            <ol>
{win_steps}
            </ol>
            <a class="button" href="{DOWNLOAD}#windows">{DOWN}{c["dl_win"]}</a>
          </div>
          <div>
            <h3>{DESKTOP}Linux</h3>
            <ol>
{linux_steps}
            </ol>
            <a class="button" href="{DOWNLOAD}#linux">{DOWN}{c["dl_linux"]}</a>
          </div>
        </div>
        <div class="install-more">
          <h3>{SERVER}{c["server_h3"]}</h3>
          <p class="small">{c["server_intro"]}</p>
          <ul class="opts">
{opts}
          </ul>
          <a class="button quiet" href="{START}#server">{c["server_compare"]}</a>
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
        <h2 class="display" id="closing-title">{c["closing"]}</h2>{cta(c, "", START)}
        <a class="source" href="{REPO}">{c["source"]}</a>
      </div>
    </section>
  </main>

{footer(c, code, up)}

  <dialog class="tour-dialog" aria-label="{esc(tour_ui["title"])}">
    <form method="dialog"><button class="tour-close" aria-label="{esc(tour_ui["close"])}">×</button></form>
    <iframe title="{esc(tour_ui["title"])}" data-src="{TOUR}?embed" allow="fullscreen" allowfullscreen></iframe>
  </dialog>
  <script src="{a}site.js" defer></script>
{ANALYTICS}</body>
</html>
"""


# ─────────────────────────── 在局域网里部署服务器 ───────────────────────────

# 教程里的示例地址。页面上的终端输出照安装脚本的真实输出抄写，只把地址换成它
LAN_IP = "192.168.1.20"

# 安装脚本在中文系统上输出简体中文（繁体系统也一样），其他系统输出英文；日文页面也是英文输出。
# 改了 server/deploy/ 下脚本的输出，这里跟着改。
TRANSCRIPTS = {
    "zh": {
        "win": f"""> 停止旧版本
> 安装到 C:\\Program Files\\PasteFerry Server
> 注册服务 PasteFerryServer（中转地址 {LAN_IP}）
> 在 Windows 防火墙里放行
> 启动
  √ 已启动

安装完成：pasteferry-server {VERSION}，作为 Windows 服务开机自动运行
在每台电脑的 PasteFerry「设置 → 信令服务器地址」填：

    ws://{LAN_IP}:8787/signal

  程序   C:\\Program Files\\PasteFerry Server
  日志   C:\\ProgramData\\PasteFerry Server\\server.log
  重启   Restart-Service PasteFerryServer（管理员 PowerShell）
  卸载   双击 uninstall.cmd""",
        "mac": f"""▶ 停止旧版本
▶ 安装到 /Users/me/Library/Application Support/PasteFerry Server
▶ 写入 LaunchAgent（中转地址 {LAN_IP}）
▶ 启动
  ✓ 已启动

安装完成：pasteferry-server {VERSION}，登录这台 Mac 后自动运行
在每台电脑的 PasteFerry「设置 → 信令服务器地址」填：

    ws://{LAN_IP}:8787/signal

  日志   /Users/me/Library/Logs/PasteFerry/server.log
  重启   launchctl kickstart -k gui/501/com.pasteferry.server
  卸载   ./install.sh uninstall""",
        "linux": f"""▶ 安装二进制到 /usr/local/bin
▶ 写入 systemd 单元
▶ 生成配置 /etc/pasteferry/server.env（中转地址 {LAN_IP}）
▶ 启动
  ✓ 已启动

安装完成：pasteferry-server {VERSION}
在每台电脑的 PasteFerry「设置 → 信令服务器地址」填：

    ws://{LAN_IP}:8787/signal

  配置   /etc/pasteferry/server.env
  日志   journalctl -u pasteferry-server -f
  重启   systemctl restart pasteferry-server
  卸载   sudo ./install.sh uninstall

防火墙需放行：TCP 8787（信令）、UDP 3478（STUN 与 TURN）、UDP 32768-60999（中转端口，系统随机分配）""",
    },
    "en": {
        "win": f"""> Stopping the running version
> Installing to C:\\Program Files\\PasteFerry Server
> Registering the service PasteFerryServer (relay address {LAN_IP})
> Allowing it through Windows Firewall
> Starting
  √ Running

Installed: pasteferry-server {VERSION}. It runs as a Windows service and starts with Windows.
In PasteFerry on each computer, set Settings > Signaling server to:

    ws://{LAN_IP}:8787/signal

  Program    C:\\Program Files\\PasteFerry Server
  Logs       C:\\ProgramData\\PasteFerry Server\\server.log
  Restart    Restart-Service PasteFerryServer (in an administrator PowerShell)
  Uninstall  double-click uninstall.cmd""",
        "mac": f"""▶ Stopping the running version
▶ Installing to /Users/me/Library/Application Support/PasteFerry Server
▶ Writing the LaunchAgent (relay address {LAN_IP})
▶ Starting
  ✓ Running

Installed: pasteferry-server {VERSION}. It starts whenever you log in to this Mac.
In PasteFerry on each computer, set Settings → Signaling server to:

    ws://{LAN_IP}:8787/signal

  Logs       /Users/me/Library/Logs/PasteFerry/server.log
  Restart    launchctl kickstart -k gui/501/com.pasteferry.server
  Uninstall  ./install.sh uninstall""",
        "linux": f"""▶ Installing the binary to /usr/local/bin
▶ Writing the systemd unit
▶ Writing /etc/pasteferry/server.env (relay address {LAN_IP})
▶ Starting
  ✓ Running

Installed: pasteferry-server {VERSION}
In PasteFerry on each computer, set Settings → Signaling server to:

    ws://{LAN_IP}:8787/signal

  Config     /etc/pasteferry/server.env
  Logs       journalctl -u pasteferry-server -f
  Restart    systemctl restart pasteferry-server
  Uninstall  sudo ./install.sh uninstall

Open in your firewall: TCP 8787 (signaling), UDP 3478 (STUN and TURN), UDP 32768-60999 (relay ports, assigned by the OS)""",
    },
}

WINDOWS_LOGO = ('<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">'
                '<path d="M3 5.5 10.5 4.5v7H3zM11.5 4.4 21 3v8.5h-9.5zM3 12.5h7.5v7L3 18.5zM11.5 12.5H21V21l-9.5-1.4z"/></svg>')


def transcript(text: str, prompt: str, command: str, ip: str = LAN_IP) -> str:
    """终端窗口里的内容：先是输入的命令，再是脚本的输出，客户端要填的地址用蓝色标出来。
    ip 是服务器的地址：脚本输出照局域网教程抄写，换成别的地址时一并替换。"""
    out = html.escape(text.replace(LAN_IP, ip), quote=False)
    url = f"ws://{ip}:8787/signal"
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


def items(lst: list[str], indent: str = "                ") -> str:
    return "\n".join(f"{indent}<li>{x}</li>" for x in lst)


def manage(caption: str, rows: list[tuple[str, str]]) -> str:
    """安装后会用到的命令表。"""
    body = "\n".join(f'                  <tr><th scope="row">{k}</th><td>{v}</td></tr>' for k, v in rows)
    return f"""<div class="table-scroll">
                <table class="manage">
                  <caption>{caption}</caption>
                  <tbody>
{body}
                  </tbody>
                </table>
              </div>"""


def fill(text: str, links: dict) -> str:
    """把文案里的 {home}、{ts} 这类占位换成链接地址。不用 str.format：文案里可能有别的花括号。"""
    for k, v in links.items():
        text = text.replace("{" + k + "}", v)
    return text


def faq_list(faq: list, links: dict) -> str:
    return "\n".join(
        f"""          <details>
            <summary>{q}</summary>
            <div class="answer">{"".join(f"<p>{fill(p, links)}</p>" for p in answer)}</div>
          </details>""" for q, answer in faq)


def toc_list(toc: list[tuple[str, str]]) -> str:
    return "\n".join(f'          <li><a href="#{anchor}">{label}</a></li>' for anchor, label in toc)


def crumb(c: dict, home: str, base: str) -> str:
    """教程顶部的面包屑：产品名 / 连接方式，后者回到快速开始里三种方式的选择。"""
    return f'<p class="crumb"><a href="{home}">{c["brand"]}</a> / <a href="{base}{START}#server">{c["nav_server"]}</a></p>'


def page_top(c: dict, code: str, sub: str, title: str, description: str) -> str:
    """教程页面从 <!doctype> 到顶栏的部分。sub 是页面目录（GUIDE、START 等）。"""
    path, hl, locale, _ = LANGS[code]
    up = "../" * (path + sub).count("/")
    a = up + "assets/"
    alternates = "\n".join(
        f'  <link rel="alternate" hreflang="{h}" href="{SITE}{p}{sub}">' for p, h, _, _ in LANGS.values())
    return f"""<!doctype html>
<html lang="{hl}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{title}</title>
  <meta name="description" content="{esc(description)}">
  <link rel="canonical" href="{SITE}{path}{sub}">
{alternates}
  <meta name="theme-color" content="#141416">
  <meta property="og:type" content="article">
  <meta property="og:site_name" content="{c["brand"]}">
  <meta property="og:title" content="{esc(title)}">
  <meta property="og:description" content="{esc(description)}">
  <meta property="og:url" content="{SITE}{path}{sub}">
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

{header(c, code, up, sub)}
"""


def page_bottom(c: dict, code: str, sub: str) -> str:
    up = "../" * (LANGS[code][0] + sub).count("/")
    return f"""
{footer(c, code, up, sub)}

  <script src="{up}assets/site.js" defer></script>
{ANALYTICS}</body>
</html>
"""


def server_page(code: str) -> str:
    c = load(code)
    g = c["guide"]
    path = LANGS[code][0]
    up = "../" * (path + GUIDE).count("/")
    home_href = up + (path or "./")
    base = up + path
    shots = up + shots_dir(code)
    t = TRANSCRIPTS["zh" if code.startswith("zh") else "en"]

    linux_dl = f"""curl -LO {DL}pasteferry-server-linux-amd64.tar.gz
tar xzf pasteferry-server-linux-amd64.tar.gz
cd pasteferry-server-linux-amd64
sudo ./install.sh"""
    mac_dl = f"""curl -LO {DL}pasteferry-server-macos.tar.gz
tar xzf pasteferry-server-macos.tar.gz
cd pasteferry-server-macos
./install.sh"""

    w, m, l = g["win"], g["mac"], g["linux"]
    explorer = []
    for name, kind in [("pasteferry-server.exe", "exe"), ("install.cmd", "cmd"), ("install.ps1", "doc"),
                       ("README.md", "doc"), ("uninstall.cmd", "cmd")]:
        # 要双击的那个文件标出来
        mark = f'<em>{w["double_click"]}</em>' if name == "install.cmd" else ""
        cls = ' class="hl"' if mark else ""
        explorer.append(f'                  <li{cls}><span class="f-icon {kind}" aria-hidden="true"></span>'
                        f'<span class="f-name">{name}</span>{mark}</li>')
    explorer = "\n".join(explorer)
    links = {"home": home_href, "lan": base + LAN, "cloud": base + CLOUD, "ts": base + TAILSCALE, "start": base + START}
    prep = "\n".join(f"          <div><dt>{h}</dt><dd>{fill(p, links)}</dd></div>" for h, p in g["prep"])
    d = g["diagram"]

    return page_top(c, code, GUIDE, g["title"], g["description"]) + f"""
  <main id="main">
    <section class="hero guide-hero">
      <div class="wrap">
        {crumb(c, home_href, base)}
        <h1 class="display">{g["h1"]}</h1>
        <p class="lede">{g["lede"]}</p>
        <p class="notice">{fill(g["notice"], links)}</p>
        <ol class="guide-toc">
{toc_list(g["toc"])}
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
        <div class="os-tabs" data-tabs="os">
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
                <a class="button" href="{DL}pasteferry-server-windows-amd64.zip">{DOWN}{w["download"]}</a>
                <ul class="points">
{items(w["notes"])}
                </ul>
              </div>
              <div class="os-figures">
                <figure class="explorer" role="img" aria-label="{esc(w["explorer_label"])}">
                  <div class="t-bar win"><span class="t-tab folder">pasteferry-server-windows-amd64</span><span class="t-ctl" aria-hidden="true"><i></i><i></i><i></i></span></div>
                  <ul class="files">
{explorer}
                  </ul>
                </figure>
                {term("win", w["term_title"], transcript(t["win"], "", ""), w["term_label"])}
              </div>
            </div>
            {manage(g["manage_caption"], w["manage"])}
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
                {term("mac", "pasteferry-server-macos — zsh", transcript(t["mac"], "%", "./install.sh"), m["term_label"])}
              </div>
            </div>
            {manage(g["manage_caption"], m["manage"])}
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
            {manage(g["manage_caption"], l["manage"])}
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
        <p class="note">{fill(g["clients_after"], links)}</p>
      </div>
    </section>

    <section class="section" id="faq" aria-labelledby="faq-title">
      <div class="wrap">
        <h2 id="faq-title">{g["faq_h2"]}</h2>
        <div class="faq">
{faq_list(g["faq"], links)}
        </div>
      </div>
    </section>

    <section class="closing" aria-labelledby="closing-title">
      <div class="wrap">
        <h2 class="display" id="closing-title">{g["closing"]}</h2>{cta(c, base)}
        <a class="source" href="{REPO}/blob/main/server/deploy/README.md">{g["source"]}</a>
      </div>
    </section>
  </main>
""" + page_bottom(c, code, GUIDE)


# ─────────────────────────── 部署到公网服务器 ───────────────────────────

# 教程里的公网 IP：RFC 5737 留给文档用的地址段，不会是任何人真实的服务器
CLOUD_IP = "203.0.113.10"


def cloud_page(code: str) -> str:
    c = load(code)
    g = c["guide"]
    cl = c["cloud"]
    path = LANGS[code][0]
    up = "../" * (path + CLOUD).count("/")
    home_href = up + (path or "./")
    base = up + path
    shots = up + shots_dir(code)
    t = TRANSCRIPTS["zh" if code.startswith("zh") else "en"]
    links = {"home": home_href, "lan": base + LAN, "server": base + GUIDE, "ts": base + TAILSCALE,
             "start": base + START, "deploy": f"{REPO}/blob/main/server/deploy/README.md"}
    d = cl["diagram"]
    install = f"""curl -LO {DL}pasteferry-server-linux-amd64.tar.gz
tar xzf pasteferry-server-linux-amd64.tar.gz
cd pasteferry-server-linux-amd64
sudo ./install.sh {CLOUD_IP}"""
    rules = "\n".join(
        "                <tr>" + "".join(f"<td>{x}</td>" for x in row) + "</tr>" for row in cl["rules"])
    cols = "".join(f'<th scope="col">{h}</th>' for h in cl["rules_cols"])

    def step(n: int) -> str:
        return f'<p class="eyebrow">{cl["step"].format(n=n)}</p>'

    return page_top(c, code, CLOUD, cl["title"], cl["description"]) + f"""
  <main id="main">
    <section class="hero guide-hero">
      <div class="wrap">
        {crumb(c, home_href, base)}
        <h1 class="display">{cl["h1"]}</h1>
        <p class="lede">{cl["lede"]}</p>
        <ol class="guide-toc">
{toc_list(cl["toc"])}
        </ol>
      </div>
    </section>

    <section class="section" id="overview" aria-labelledby="overview-title">
      <div class="wrap">
        <h2 id="overview-title">{cl["overview_h2"]}</h2>
        <p class="intro">{cl["overview_intro"]}</p>
        <figure class="diagram">
          <svg viewBox="0 0 960 460" role="img" aria-labelledby="cloud-title cloud-desc">
            <title id="cloud-title">{d["title"]}</title>
            <desc id="cloud-desc">{d["desc"]}</desc>
            <rect class="d-box" x="340" y="24" width="280" height="118" rx="14"/>
            <text class="d-title" x="480" y="62" text-anchor="middle">{d["server"]}</text>
            <text class="d-sub" x="480" y="90" text-anchor="middle">{d["server_sub"]}</text>
            <text class="d-ip" x="480" y="124" text-anchor="middle">{CLOUD_IP}</text>
            <rect class="d-lan" x="8" y="196" width="412" height="212" rx="22"/>
            <text class="d-lan-label" x="32" y="228">{d["office"]}</text>
            <rect class="d-lan" x="540" y="196" width="412" height="212" rx="22"/>
            <text class="d-lan-label" x="564" y="228">{d["home"]}</text>
            <path class="d-sig" d="M150 250V68H340"/>
            <path class="d-sig" d="M810 250V68H620"/>
            <path class="d-relay" d="M210 250V112H340"/>
            <path class="d-relay" d="M750 250V112H620"/>
            <text class="d-note" x="164" y="58">{d["sig_note"]}</text>
            <text class="d-note" x="224" y="170">{d["relay_note"]}</text>
            <rect class="d-box" x="60" y="250" width="240" height="120" rx="14"/>
            <text class="d-title" x="180" y="302" text-anchor="middle">Mac</text>
            <text class="d-sub" x="180" y="330" text-anchor="middle">{d["client"]}</text>
            <rect class="d-box" x="660" y="250" width="240" height="120" rx="14"/>
            <text class="d-title" x="780" y="302" text-anchor="middle">Windows</text>
            <text class="d-sub" x="780" y="330" text-anchor="middle">{d["client"]}</text>
            <path class="d-direct" d="M300 316H660"/>
            <text class="d-label" x="480" y="302" text-anchor="middle">{d["direct"]}</text>
            <g class="d-legend" transform="translate(40 444)">
              <path class="d-direct" d="M0 0H28"/><text x="38" y="5">{g["diagram"]["legend_direct"]}</text>
              <path class="d-relay" d="M250 0H278"/><text x="288" y="5">{g["diagram"]["legend_relay"]}</text>
              <path class="d-sig" d="M560 0H588"/><text x="598" y="5">{g["diagram"]["legend_sig"]}</text>
            </g>
          </svg>
        </figure>
        <ul class="points ports">
{items(cl["overview_points"], "          ")}
        </ul>
      </div>
    </section>

    <section class="section" id="prepare" aria-labelledby="prepare-title">
      <div class="wrap">
        <h2 id="prepare-title">{g["prep_h2"]}</h2>
        <dl class="fallbacks">
{chr(10).join(f"          <div><dt>{h}</dt><dd>{fill(p, links)}</dd></div>" for h, p in cl["prep"])}
        </dl>
      </div>
    </section>

    <section class="section" id="ports" aria-labelledby="ports-title">
      <div class="wrap">
        {step(1)}
        <h2 id="ports-title">{cl["ports_h2"]}</h2>
        <p class="intro">{cl["ports_intro"]}</p>
        <div class="os-grid qs-grid">
          <div class="os-figures">
            <figure class="browser listing" role="img" aria-label="{esc(cl["rules_label"])}">
              <div class="b-bar"><span class="t-lights" aria-hidden="true"><i></i><i></i><i></i></span><span class="b-url">{cl["rules_url"]}</span></div>
              <table>
                <thead><tr>{cols}</tr></thead>
                <tbody>
{rules}
                </tbody>
              </table>
            </figure>
          </div>
          <ul class="points">
{items(cl["ports_notes"], "            ")}
          </ul>
        </div>
      </div>
    </section>

    <section class="section" id="install" aria-labelledby="install-title">
      <div class="wrap">
        {step(2)}
        <h2 id="install-title">{cl["install_h2"]}</h2>
        <p class="intro">{cl["install_intro"]}</p>
        <div class="os-grid qs-grid">
          <div class="os-steps">
            <ol class="steps">
              <li>{cl["install_steps"][0]}</li>
              <li>{cl["install_steps"][1]}
                {code_block(c, install)}</li>
{items(cl["install_steps"][2:], "              ")}
            </ol>
            <ul class="points">
{items([fill(x, links) for x in cl["install_notes"]], "              ")}
            </ul>
          </div>
          <div class="os-figures">
            {term("mac", f"ssh root@{CLOUD_IP}", transcript(t["linux"], "$", f"sudo ./install.sh {CLOUD_IP}", CLOUD_IP), cl["term_label"])}
          </div>
        </div>
      </div>
    </section>

    <section class="section" id="check" aria-labelledby="check-title">
      <div class="wrap">
        {step(3)}
        <h2 id="check-title">{cl["check_h2"]}</h2>
        <p class="intro">{cl["check_intro"]}</p>
        <figure class="browser" role="img" aria-label="{esc(cl["check_label"])}">
          <div class="b-bar"><span class="t-lights" aria-hidden="true"><i></i><i></i><i></i></span><span class="b-url">http://{CLOUD_IP}:8787/healthz</span></div>
          <pre class="b-body">{{"status":"ok","version":"{VERSION}"}}</pre>
        </figure>
        <p class="note">{g["check_after"]}</p>
      </div>
    </section>

    <section class="section" id="clients" aria-labelledby="clients-title">
      <div class="wrap">
        {step(4)}
        <h2 id="clients-title">{g["clients_h2"]}</h2>
        <ol class="flow">
{chr(10).join(f"          <li><h3>{h}</h3><p>{p}</p></li>" for h, p in cl["clients_steps"])}
        </ol>
        <figure class="shot">
          <img src="{shots}settings-cloud-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(cl["clients_alt"])}">
        </figure>
        <p class="note">{fill(g["clients_after"], links)}</p>
        {manage(g["manage_caption"], g["linux"]["manage"])}
      </div>
    </section>

    <section class="section" id="faq" aria-labelledby="faq-title">
      <div class="wrap">
        <h2 id="faq-title">{g["faq_h2"]}</h2>
        <div class="faq">
{faq_list(cl["faq"], links)}
        </div>
      </div>
    </section>

    <section class="closing" aria-labelledby="closing-title">
      <div class="wrap">
        <h2 class="display" id="closing-title">{cl["closing"]}</h2>{cta(c, base)}
        <a class="source" href="{links["deploy"]}">{cl["source"]}</a>
      </div>
    </section>
  </main>
""" + page_bottom(c, code, CLOUD)


# ─────────────────────────── 快速开始 ───────────────────────────


def start_page(code: str) -> str:
    c = load(code)
    s = c["start"]
    path = LANGS[code][0]
    up = "../" * (path + START).count("/")
    home_href = up + (path or "./")
    base = up + path
    shots = up + shots_dir(code)
    links = {"home": home_href, "lan": base + LAN, "cloud": base + CLOUD, "server": base + GUIDE,
             "ts": base + TAILSCALE, "issues": REPO + "/issues"}

    def step(n: int) -> str:
        return f'<p class="eyebrow">{s["step"].format(n=n)}</p>'

    return page_top(c, code, START, s["title"], s["description"]) + f"""
  <main id="main">
    <section class="hero guide-hero">
      <div class="wrap">
        <p class="crumb"><a href="{home_href}">{c["brand"]}</a> / {c["nav_start"]}</p>
        <h1 class="display">{s["h1"]}</h1>
        <p class="lede">{s["lede"]}</p>
        <ol class="guide-toc">
{toc_list(s["toc"])}
        </ol>
      </div>
    </section>

    <section class="section" id="app" aria-labelledby="app-title">
      <div class="wrap">
        {step(1)}
        <h2 id="app-title">{s["app_h2"]}</h2>
        <p class="intro">{s["app_intro"]}</p>
        <div class="os-grid qs-grid">
          <div>
            <div class="qs-app">
              <h3>{LAPTOP}Mac</h3>
              <ol class="steps">
{items(s["mac"])}
              </ol>
              <a class="button" href="{base}{DOWNLOAD}#macos">{DOWN}{c["dl_mac"]}</a>
            </div>
            <div class="qs-app">
              <h3>{DESKTOP}Windows</h3>
              <ol class="steps">
{items(s["win"])}
              </ol>
              <a class="button" href="{base}{DOWNLOAD}#windows">{DOWN}{c["dl_win"]}</a>
            </div>
            <div class="qs-app">
              <h3>{DESKTOP}Linux</h3>
              <ol class="steps">
{items(s["linux"])}
              </ol>
              <a class="button" href="{base}{DOWNLOAD}#linux">{DOWN}{c["dl_linux"]}</a>
            </div>
          </div>
          <figure class="shot">
            <img src="{shots}welcome-dark.png" width="2112" height="1472" alt="{esc(s["app_alt"])}">
          </figure>
        </div>
      </div>
    </section>

    <section class="section" id="server" aria-labelledby="server-title">
      <div class="wrap">
        {step(2)}
        <h2 id="server-title">{s["server_h2"]}</h2>
        <p class="intro">{s["server_intro"]}</p>
        {plan_tabs(c, base)}
      </div>
    </section>

    <section class="section" id="pair" aria-labelledby="pair-title">
      <div class="wrap">
        {step(3)}
        <h2 id="pair-title">{s["pair_h2"]}</h2>
        <p class="intro">{s["pair_intro"]}</p>
        <div class="os-grid qs-grid">
          <div>
            <ol class="steps">
{items([fill(x, links) for x in s["pair_steps"]], "              ")}
            </ol>
            <figure class="figure-card qs-verify">
              <img src="{shots}card-verify-dark.png" width="1120" height="800" loading="lazy" alt="{esc(c["verify_alt"])}">
            </figure>
          </div>
          <figure class="shot">
            <img src="{shots}pairing-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(s["pair_alt"])}">
          </figure>
        </div>
      </div>
    </section>

    <section class="closing" aria-labelledby="closing-title">
      <div class="wrap">
        <h2 class="display" id="closing-title">{s["done_h2"]}</h2>
        <p class="lede">{s["done_text"]}</p>
        <p class="note">{fill(s["help"], links)}</p>
      </div>
    </section>
  </main>
""" + page_bottom(c, code, START)


# ─────────────────────────── 同一网络（不需要服务器） ───────────────────────────

# 教程里的两台电脑在同一个局域网里的地址
LAN_MAC = "192.168.1.20"
LAN_WIN = "192.168.1.35"


def lan_page(code: str) -> str:
    c = load(code)
    g = c["guide"]
    ln = c["lan"]
    path = LANGS[code][0]
    up = "../" * (path + LAN).count("/")
    home_href = up + (path or "./")
    base = up + path
    shots = up + shots_dir(code)
    links = {"home": home_href, "cloud": base + CLOUD, "server": base + GUIDE, "ts": base + TAILSCALE,
             "start": base + START, "design": up + "design/lan-mode.html"}
    d = ln["diagram"]

    def step(n: int) -> str:
        return f'<p class="eyebrow">{ln["step"].format(n=n)}</p>'

    return page_top(c, code, LAN, ln["title"], ln["description"]) + f"""
  <main id="main">
    <section class="hero guide-hero">
      <div class="wrap">
        {crumb(c, home_href, base)}
        <h1 class="display">{ln["h1"]}</h1>
        <p class="lede">{ln["lede"]}</p>
        <ol class="guide-toc">
{toc_list(ln["toc"])}
        </ol>
      </div>
    </section>

    <section class="section" id="overview" aria-labelledby="overview-title">
      <div class="wrap">
        <h2 id="overview-title">{ln["overview_h2"]}</h2>
        <p class="intro">{ln["overview_intro"]}</p>
        <figure class="diagram">
          <svg viewBox="0 0 960 330" role="img" aria-labelledby="lan-title lan-desc">
            <title id="lan-title">{d["title"]}</title>
            <desc id="lan-desc">{d["desc"]}</desc>
            <rect class="d-lan" x="8" y="8" width="944" height="314" rx="22"/>
            <text class="d-lan-label" x="32" y="42">{d["lan"]}</text>
            <rect class="d-box" x="40" y="120" width="240" height="140" rx="14"/>
            <text class="d-title" x="160" y="170" text-anchor="middle">Mac</text>
            <text class="d-sub" x="160" y="198" text-anchor="middle">{d["client"]}</text>
            <text class="d-ip" x="160" y="234" text-anchor="middle">{LAN_MAC}</text>
            <rect class="d-box" x="680" y="120" width="240" height="140" rx="14"/>
            <text class="d-title" x="800" y="170" text-anchor="middle">Windows</text>
            <text class="d-sub" x="800" y="198" text-anchor="middle">{d["client"]}</text>
            <text class="d-ip" x="800" y="234" text-anchor="middle">{LAN_WIN}</text>
            <path class="d-sig" d="M280 150H680"/>
            <text class="d-note" x="480" y="138" text-anchor="middle">{d["find_note"]}</text>
            <path class="d-direct" d="M280 214H680"/>
            <text class="d-label" x="480" y="200" text-anchor="middle">{d["direct"]}</text>
            <text class="d-note" x="480" y="242" text-anchor="middle">{d["direct_note"]}</text>
          </svg>
        </figure>
        <ul class="points ports">
{items(ln["overview_points"], "          ")}
        </ul>
      </div>
    </section>

    <section class="section" id="prepare" aria-labelledby="prepare-title">
      <div class="wrap">
        <h2 id="prepare-title">{g["prep_h2"]}</h2>
        <dl class="fallbacks">
{chr(10).join(f"          <div><dt>{h}</dt><dd>{fill(p, links)}</dd></div>" for h, p in ln["prep"])}
        </dl>
      </div>
    </section>

    <section class="section" id="pair" aria-labelledby="pair-title">
      <div class="wrap">
        {step(1)}
        <h2 id="pair-title">{ln["pair_h2"]}</h2>
        <p class="intro">{fill(ln["pair_intro"], links)}</p>
        <div class="os-grid qs-grid">
          <div>
            <ol class="steps">
{items(ln["pair_steps"], "              ")}
            </ol>
            <figure class="figure-card qs-verify">
              <img src="{shots}card-verify-dark.png" width="1120" height="800" loading="lazy" alt="{esc(c["verify_alt"])}">
            </figure>
          </div>
          <figure class="shot">
            <img src="{shots}pairing-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(ln["pair_alt"])}">
          </figure>
        </div>
        <div class="os-grid qs-grid">
          <div class="qs-app">
            <h3>{ln["nearby_h3"]}</h3>
            <ol class="steps">
{items(ln["nearby_steps"], "              ")}
            </ol>
          </div>
          <figure class="shot">
            <img src="{shots}nearby-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(ln["nearby_alt"])}">
          </figure>
        </div>
      </div>
    </section>

    <section class="section" id="subnet" aria-labelledby="subnet-title">
      <div class="wrap">
        {step(2)}
        <h2 id="subnet-title">{ln["subnet_h2"]}</h2>
        <p class="intro">{ln["subnet_intro"]}</p>
        <div class="os-grid qs-grid">
          <div>
            <ol class="steps">
{items(ln["subnet_steps"], "              ")}
            </ol>
          </div>
          <figure class="shot">
            <img src="{shots}pairing-join-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(ln["subnet_alt"])}">
          </figure>
        </div>
        <p class="note">{fill(ln["subnet_after"], links)}</p>
      </div>
    </section>

    <section class="section" id="faq" aria-labelledby="faq-title">
      <div class="wrap">
        <h2 id="faq-title">{g["faq_h2"]}</h2>
        <div class="faq">
{faq_list(ln["faq"], links)}
        </div>
      </div>
    </section>

    <section class="closing" aria-labelledby="closing-title">
      <div class="wrap">
        <h2 class="display" id="closing-title">{ln["closing"]}</h2>{cta(c, base)}
        <a class="source" href="{links["design"]}">{ln["source"]}</a>
      </div>
    </section>
  </main>
""" + page_bottom(c, code, LAN)


# ─────────────────────────── 下载 ───────────────────────────


def download_page(code: str) -> str:
    """下载页：按系统分页签（打开时选中访客自己的系统，见 site.js），Mac 再分芯片。"""
    c = load(code)
    d = c["download"]
    path = LANGS[code][0]
    up = "../" * (path + DOWNLOAD).count("/")
    home_href = up + (path or "./")
    base = up + path
    links = {"start": base + START, "server": base + START + "#server", "releases": REPO + "/releases",
             "sums": DL + "SHA256SUMS", "universal": DL + "PasteFerry.dmg", "portable": DL + "PasteFerry-windows-x64.zip",
             "notes": REPO + "/releases/latest", "lan": base + LAN}
    w, m, lx = d["win"], d["mac"], d["linux"]

    def button(name: str, label: str, quiet: bool = False) -> str:
        return f'<a class="button{" quiet" if quiet else ""}" href="{DL}{name}">{DOWN}{label}</a>'

    more = "\n".join(f"          <div><dt>{t}</dt><dd>{fill(v, links)}</dd></div>" for t, v in d["more"])

    return page_top(c, code, DOWNLOAD, d["title"], d["description"]) + f"""
  <main id="main">
    <section class="hero guide-hero">
      <div class="wrap">
        <p class="crumb"><a href="{home_href}">{c["brand"]}</a> / {c["nav_download"]}</p>
        <h1 class="display">{d["h1"]}</h1>
        <p class="lede">{fill(d["lede"], {"version": VERSION})}</p>
      </div>
    </section>

    <section class="section" id="get" aria-label="{d["os_label"]}">
      <div class="wrap">
        <div class="os-tabs" data-tabs="os">
          <div class="tablist" role="tablist" aria-label="{d["os_label"]}" hidden>
            <button type="button" role="tab" id="tab-windows" aria-controls="windows">Windows</button>
            <button type="button" role="tab" id="tab-macos" aria-controls="macos">macOS</button>
            <button type="button" role="tab" id="tab-linux" aria-controls="linux">Linux</button>
          </div>

          <section class="os" id="windows" role="tabpanel" aria-labelledby="tab-windows">
            <h3 class="os-title">{DESKTOP}Windows</h3>
            <div class="os-grid">
              <div class="os-steps">
                <div class="dl-row">
                  {button("PasteFerry-Setup.exe", w["setup"])}
                </div>
                <p class="small">{fill(w["portable"], links)}</p>
                <ol class="steps">
{items(w["steps"])}
                </ol>
              </div>
              <ul class="points">
{items(w["notes"])}
              </ul>
            </div>
          </section>

          <section class="os" id="macos" role="tabpanel" aria-labelledby="tab-macos">
            <h3 class="os-title">{LAPTOP}macOS</h3>
            <div class="os-grid">
              <div class="os-steps">
                <div class="dl-row" data-mac-chip>
                  {button("PasteFerry-mac-apple-silicon.dmg", m["apple"])}
                  {button("PasteFerry-mac-intel.dmg", m["intel"], True)}
                </div>
                <p class="small">{fill(m["universal"], links)}</p>
                <ol class="steps">
{items(m["steps"])}
                </ol>
              </div>
              <div>
                <h4>{m["chip_h"]}</h4>
                <ul class="points">
{items(m["chip"])}
                </ul>
              </div>
            </div>
          </section>

          <section class="os" id="linux" role="tabpanel" aria-labelledby="tab-linux">
            <h3 class="os-title">{DESKTOP}Linux</h3>
            <div class="os-grid">
              <div class="os-steps">
                <h4>{lx["deb_h"]}</h4>
                <div class="dl-row">
                  {button("PasteFerry-linux-amd64.deb", lx["x64"])}
                  {button("PasteFerry-linux-arm64.deb", lx["arm64"], True)}
                </div>
                {code_block(c, "sudo apt install ./PasteFerry-linux-amd64.deb")}
                <h4>{lx["tar_h"]}</h4>
                <div class="dl-row">
                  {button("PasteFerry-linux-amd64.tar.gz", lx["x64"], True)}
                  {button("PasteFerry-linux-arm64.tar.gz", lx["arm64"], True)}
                </div>
                {code_block(c, "tar xzf PasteFerry-linux-amd64.tar.gz" + chr(10) + "./PasteFerry/install.sh")}
              </div>
              <ul class="points">
{items(lx["notes"])}
              </ul>
            </div>
          </section>
        </div>
      </div>
    </section>

    <section class="section" aria-labelledby="more-title">
      <div class="wrap">
        <h2 id="more-title">{d["more_h2"]}</h2>
        <dl class="fallbacks">
{more}
        </dl>
      </div>
    </section>
  </main>
""" + page_bottom(c, code, DOWNLOAD)


# ─────────────────────────── 用 Tailscale 跨网络使用 ───────────────────────────

# 教程里的 Tailscale 地址：两台电脑各一个，都在 Tailscale 用的 100.64.0.0/10 里
TS_MAC = "100.88.20.5"
TS_WIN = "100.92.41.17"
# 管理后台设备列表里的示例：主机名、地址、系统
TS_MACHINES = [("macbook-pro", TS_MAC, "macOS"), ("home-pc", TS_WIN, "Windows")]


def tailscale_page(code: str) -> str:
    c = load(code)
    g = c["guide"]
    ts = c["ts"]
    path = LANGS[code][0]
    up = "../" * (path + TAILSCALE).count("/")
    home_href = up + (path or "./")
    base = up + path
    shots = up + shots_dir(code)
    links = {"home": home_href, "lan": base + LAN, "cloud": base + CLOUD, "server": base + GUIDE,
             "start": base + START}
    d = ts["diagram"]

    rows = "\n".join(f"                <tr><td>{name}</td><td>{ip}</td><td>{system}</td></tr>"
                     for name, ip, system in TS_MACHINES)
    cols = "".join(f'<th scope="col">{h}</th>' for h in ts["machines_cols"])

    return page_top(c, code, TAILSCALE, ts["title"], ts["description"]) + f"""
  <main id="main">
    <section class="hero guide-hero">
      <div class="wrap">
        {crumb(c, home_href, base)}
        <h1 class="display">{ts["h1"]}</h1>
        <p class="lede">{ts["lede"]}</p>
        <ol class="guide-toc">
{toc_list(ts["toc"])}
        </ol>
      </div>
    </section>

    <section class="section" id="overview" aria-labelledby="overview-title">
      <div class="wrap">
        <h2 id="overview-title">{ts["overview_h2"]}</h2>
        <p class="intro">{ts["overview_intro"]}</p>
        <figure class="diagram">
          <svg viewBox="0 0 960 330" role="img" aria-labelledby="ts-title ts-desc">
            <title id="ts-title">{d["title"]}</title>
            <desc id="ts-desc">{d["desc"]}</desc>
            <rect class="d-ts" x="8" y="8" width="944" height="314" rx="22"/>
            <text class="d-ts-label" x="32" y="42">{d["net"]}</text>
            <rect class="d-box" x="40" y="100" width="240" height="160" rx="14"/>
            <text class="d-title" x="160" y="146" text-anchor="middle">Mac</text>
            <text class="d-sub" x="160" y="174" text-anchor="middle">{d["mac_place"]}</text>
            <text class="d-ip" x="160" y="226" text-anchor="middle">{TS_MAC}</text>
            <rect class="d-box" x="680" y="100" width="240" height="160" rx="14"/>
            <text class="d-title" x="800" y="146" text-anchor="middle">Windows</text>
            <text class="d-sub" x="800" y="174" text-anchor="middle">{d["win_place"]}</text>
            <text class="d-ip" x="800" y="226" text-anchor="middle">{TS_WIN}</text>
            <path class="d-direct" d="M280 190H680"/>
            <text class="d-label" x="480" y="176" text-anchor="middle">{d["direct"]}</text>
            <text class="d-note" x="480" y="218" text-anchor="middle">{d["direct_note"]}</text>
          </svg>
        </figure>
        <ul class="points ports">
{items(ts["overview_points"], "          ")}
        </ul>
      </div>
    </section>

    <section class="section" id="prepare" aria-labelledby="prepare-title">
      <div class="wrap">
        <h2 id="prepare-title">{g["prep_h2"]}</h2>
        <dl class="fallbacks">
{chr(10).join(f"          <div><dt>{h}</dt><dd>{fill(p, links)}</dd></div>" for h, p in ts["prep"])}
        </dl>
      </div>
    </section>

    <section class="section" id="join" aria-labelledby="join-title">
      <div class="wrap">
        <p class="eyebrow">{ts["step"].format(n=1)}</p>
        <h2 id="join-title">{ts["join_h2"]}</h2>
        <p class="intro">{ts["join_intro"]}</p>
        <div class="os-grid qs-grid">
          <div class="os-steps">
            <ol class="steps">
{items(ts["join_steps"], "              ")}
            </ol>
            <a class="button" href="https://tailscale.com/download">{DOWN}{ts["join_download"]}</a>
          </div>
          <div class="os-figures">
            <figure class="browser listing" role="img" aria-label="{esc(ts["machines_label"])}">
              <div class="b-bar"><span class="t-lights" aria-hidden="true"><i></i><i></i><i></i></span><span class="b-url">login.tailscale.com/admin/machines</span></div>
              <table>
                <thead><tr>{cols}</tr></thead>
                <tbody>
{rows}
                </tbody>
              </table>
            </figure>
          </div>
        </div>
        <div class="server-notes">
          <h3>{ts["join_notes_h3"]}</h3>
          <ul class="points">
{items(ts["join_notes"], "            ")}
          </ul>
        </div>
      </div>
    </section>

    <section class="section" id="address" aria-labelledby="address-title">
      <div class="wrap">
        <p class="eyebrow">{ts["step"].format(n=2)}</p>
        <h2 id="address-title">{ts["address_h2"]}</h2>
        <p class="intro">{ts["address_intro"]}</p>
        <div class="address-term">
          {term("mac", "home-pc — PowerShell", f'<span class="t-prompt">&gt;</span> tailscale ip -4{chr(10)}<span class="t-url">{TS_WIN}</span>', ts["address_label"])}
        </div>
        <p class="note">{ts["address_after"]}</p>
      </div>
    </section>

    <section class="section" id="pair" aria-labelledby="pair-title">
      <div class="wrap">
        <p class="eyebrow">{ts["step"].format(n=3)}</p>
        <h2 id="pair-title">{ts["pair_h2"]}</h2>
        <p class="intro">{ts["pair_intro"]}</p>
        <div class="os-grid qs-grid">
          <div>
            <ol class="steps">
{items(ts["pair_steps"], "              ")}
            </ol>
          </div>
          <figure class="shot">
            <img src="{shots}pairing-join-ts-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(ts["pair_alt"])}">
          </figure>
        </div>
        <p class="note">{fill(ts["pair_after"], links)}</p>
      </div>
    </section>

    <section class="section" id="faq" aria-labelledby="faq-title">
      <div class="wrap">
        <h2 id="faq-title">{g["faq_h2"]}</h2>
        <div class="faq">
{faq_list(ts["faq"], links)}
        </div>
      </div>
    </section>

    <section class="closing" aria-labelledby="closing-title">
      <div class="wrap">
        <h2 class="display" id="closing-title">{ts["closing"]}</h2>{cta(c, base)}
        <a class="source" href="{links["lan"]}">{ts["source"]}</a>
      </div>
    </section>
  </main>
""" + page_bottom(c, code, TAILSCALE)


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
  <title>{esc(text_of(inline(title)))}：隔空贴技术方案</title>
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
      <a class="brand" href="../" aria-label="隔空贴首页"><img src="../assets/icon.svg" alt="" width="26" height="26">隔空贴</a>
      <nav aria-label="页面导航">
        <a href="../{START}">{zh["nav_start"]}</a>
        <a href="../{START}#server">{zh["nav_server"]}</a>
        <a class="wide" href="{REPO}/blob/main/design/{md_name}">在 GitHub 上查看</a>
        <a class="wide" href="{REPO}">GitHub</a>
      </nav>
    </div>
  </header>
  <main id="main">
    <article class="article">
      <p class="crumb"><a href="../">隔空贴</a>的技术方案</p>
      <h1>{inline(title)}</h1>
{body}
    </article>
  </main>

{footer(zh, "zh", "../")}
{ANALYTICS}</body>
</html>
"""


# ─────────────────────────── 其他页面 ───────────────────────────


def not_found() -> str:
    return f"""<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="robots" content="noindex">
  <title>页面不存在 / Page not found：隔空贴 PasteFerry</title>
  <meta name="theme-color" content="#141416">
  <link rel="icon" href="{BASE}assets/favicon.svg" type="image/svg+xml">
  <link rel="stylesheet" href="{BASE}assets/site.css">
</head>
<body>
  <main class="lost">
    <div class="wrap">
      <h1 class="display">404</h1>
      <p class="lede">这个页面不存在。<br><span lang="en">This page doesn't exist.</span><br><span lang="ja">このページは存在しません。</span></p>
      <div class="cta">
        <a class="button" href="{BASE}">回到首页</a>
        <a class="button quiet" href="{BASE}en/" lang="en">English home</a>
        <a class="button quiet" href="{BASE}ja/" lang="ja">日本語のホーム</a>
      </div>
    </div>
  </main>
{ANALYTICS}</body>
</html>
"""


def sitemap() -> str:
    alts = "".join(
        f'\n    <xhtml:link rel="alternate" hreflang="{h}" href="{SITE}{p}"/>' for p, h, _, _ in LANGS.values())
    alts += f'\n    <xhtml:link rel="alternate" hreflang="x-default" href="{SITE}"/>'
    urls = [f"  <url>\n    <loc>{SITE}{p}</loc>\n    <lastmod>{DATE}</lastmod>{alts}\n  </url>"
            for p, _, _, _ in LANGS.values()]
    for sub in SUBPAGES:
        sub_alts = "".join(
            f'\n    <xhtml:link rel="alternate" hreflang="{h}" href="{SITE}{p}{sub}"/>' for p, h, _, _ in LANGS.values())
        urls += [f"  <url>\n    <loc>{SITE}{p}{sub}</loc>\n    <lastmod>{DATE}</lastmod>{sub_alts}\n  </url>"
                 for p, _, _, _ in LANGS.values()]
    urls += [f"  <url>\n    <loc>{SITE}design/{page}.html</loc>\n    <lastmod>{DATE}</lastmod>\n  </url>"
             for page in DESIGN_DOCS.values()]
    return ('<?xml version="1.0" encoding="UTF-8"?>\n'
            '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"\n'
            '        xmlns:xhtml="http://www.w3.org/1999/xhtml">\n' + "\n".join(urls) + "\n</urlset>\n")


def tour_page(code: str) -> str:
    c = load(code)
    copy, ui, chapters = tour_copy(code, c)
    path, hl, _, _ = LANGS[code]
    up = "../" * (path.count("/") + 1)
    assets = up + "assets/"
    copy.update(shots=up + shots_dir(code), site_url=SITE.split("://", 1)[-1] + path)
    template = (Path(__file__).parent / "tour.html").read_text(encoding="utf-8")
    art = re.sub(r"\{\{(\w+)\}\}", lambda m: esc(copy[m[1]]), template)
    options = "\n".join(f'<option value="{i}">{i + 1:02d} · {esc(label)}</option>' for i, label in enumerate(chapters))
    alternates = "\n".join(f'<link rel="alternate" hreflang="{h}" href="{SITE}{p}{TOUR}">' for p, h, _, _ in LANGS.values())
    highlights = "\n".join(f'<li>{esc(copy[key])}</li>' for key in ("demo_file", "demo_text", "types", "retention", "pairing", "encrypted", "relay"))
    labels = esc(json.dumps(ui, ensure_ascii=False))
    return f"""<!doctype html>
<html lang="{hl}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{esc(ui["title"])} — {esc(c["brand"])}</title>
  <meta name="description" content="{esc(c["og_description"])}">
  <meta name="theme-color" content="#141416">
  <link rel="canonical" href="{SITE}{path}{TOUR}">
  {alternates}
  <link rel="icon" href="{assets}favicon.svg" type="image/svg+xml">
  <link rel="stylesheet" href="{assets}tour.css">
  <script>if (new URLSearchParams(location.search).has("embed")) document.documentElement.classList.add("embedded");</script>
</head>
<body>
  <header class="tour-header">
    <a class="tour-brand" href="../"><img src="{assets}icon.svg" width="26" height="26" alt="">{esc(c["brand"])}</a>
    <nav><a href="../">{esc(ui["back"])}</a>{lang_menu(c, code, up, TOUR)}</nav>
  </header>
  <main>
    <div class="tour-player" data-labels="{labels}">
      <div class="tour-viewport">{art}</div>
      <div class="tour-controls" hidden>
        <button class="tour-toggle" type="button">{esc(ui["play"])}</button>
        <output class="tour-time">0:00 / 1:25</output>
        <input class="tour-progress" type="range" min="0" max="84.5" step="0.1" value="0" aria-label="{esc(ui["progress"])}">
        <select class="tour-chapters" aria-label="{esc(ui["chapters"])}">{options}</select>
        <button class="tour-fullscreen" type="button" aria-label="{esc(ui["fullscreen"])}" title="{esc(ui["fullscreen"])}">⛶</button>
      </div>
      <p class="tour-status" role="status"></p>
    </div>
    <section class="tour-summary">
      <h1>{esc(ui["title"])}</h1>
      <p class="tour-silent">{esc(ui["silent"])} · macOS / Windows / Linux</p>
      <details><summary>{esc(ui["summary"])}</summary><ul>{highlights}</ul></details>
      <a class="tour-download" href="../{DOWNLOAD}">{esc(c["dl_button"])} →</a>
    </section>
  </main>
  <noscript><style>.tour-player {{ display:none; }}</style></noscript>
  <script src="{assets}tour-scenes.js" defer></script>
  <script src="{assets}tour.js" defer></script>
  <script src="{assets}site.js" defer></script>
{ANALYTICS}</body>
</html>
"""


def write(path: Path, text: str) -> None:
    if path.suffix == ".html":
        text = text.replace('  <meta charset="utf-8">\n',
                            '  <meta charset="utf-8">\n' + GITHUB_PAGES_REDIRECT, 1)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text, encoding="utf-8")
    print("✓", path.relative_to(ROOT))

def privacy_page(code: str) -> str:
    c = load(code)
    title, sections = PRIVACY_COPY[code]
    body = "\n".join(f'<section><h2>{esc(h)}</h2><p>{esc(p)}</p></section>' for h, p in sections)
    return (page_top(c, code, PRIVACY, f'{esc(title)} — PasteFerry', sections[0][1])
            + f'<main id="main" class="wrap" style="max-width:850px;padding-top:5rem;padding-bottom:5rem;line-height:1.9"><h1>{esc(title)}</h1>{body}'
            + f'<p><a href="{REPO}/blob/main/PRIVACY.md">PRIVACY.md</a> · <a href="{REPO}/blob/main/SECURITY.md">SECURITY.md</a> · <a href="{REPO}/blob/main/design/official-gateway.md">Gateway protocol v1</a></p></main>'
            + page_bottom(c, code, PRIVACY))


def official_page(code: str) -> str:
    c = load(code)
    o = c["official"]
    path = LANGS[code][0]
    up = "../" * (path + OFFICIAL).count("/")
    base = up + path
    shots = up + shots_dir(code)
    sections = "\n".join(
        f'<section class="section"><div class="wrap official-copy"><h2>{esc(o[key + "_title"])}</h2><p class="intro">{esc(o[key])}</p></div></section>'
        for key in ("routing", "quota", "privacy", "self")
    )
    return page_top(c, code, OFFICIAL, o["title"], o["summary"]) + f"""
  <main id="main">
    <section class="hero guide-hero"><div class="wrap">
      <p class="crumb"><a href="{base or './'}">{c['brand']}</a> / {o['name']}</p>
      <h1 class="display">{esc(o['fit'])}</h1>
      <p class="lede">{esc(o['lede'])}</p>
      <a class="button" href="{base}{DOWNLOAD}">{DOWN}{esc(o['download'])}</a>
    </div></section>
    <section class="section"><div class="wrap">
      <h2>{esc(o['steps_title'])}</h2>
      <div class="os-grid"><ol class="steps">{items(o['steps'])}</ol>
        <figure class="shot"><img src="{shots}official-settings-dark.png" width="2112" height="1472" alt="{esc(o['settings_alt'])}"></figure>
      </div>
    </div></section>
    {sections}
    <section class="section"><div class="wrap official-copy">
      <figure class="shot"><img src="{shots}official-usage-dark.png" width="2112" height="1472" loading="lazy" alt="{esc(o['usage_alt'])}"></figure>
      <p class="note"><a href="{base}{PRIVACY}">{esc(o['privacy_link'])}</a> · <a href="{base}{CLOUD}">{esc(o['self_link'])}</a> · <a href="https://pasteferry.cc/">pasteferry.cc</a></p>
    </div></section>
  </main>
""" + page_bottom(c, code, OFFICIAL)


def main() -> None:
    for code, (path, _, _, _) in LANGS.items():
        write(DOCS / path / "index.html", home(code))
        write(DOCS / path / TOUR / "index.html", tour_page(code))
        write(DOCS / path / GUIDE / "index.html", server_page(code))
        write(DOCS / path / START / "index.html", start_page(code))
        write(DOCS / path / LAN / "index.html", lan_page(code))
        write(DOCS / path / CLOUD / "index.html", cloud_page(code))
        write(DOCS / path / TAILSCALE / "index.html", tailscale_page(code))
        write(DOCS / path / DOWNLOAD / "index.html", download_page(code))
        write(DOCS / path / PRIVACY / "index.html", privacy_page(code))
        write(DOCS / path / OFFICIAL / "index.html", official_page(code))
    for md_name, page in DESIGN_DOCS.items():
        write(DOCS / "design" / f"{page}.html", design_page(md_name, page))
    write(DOCS / "404.html", not_found())
    write(DOCS / "sitemap.xml", sitemap())


if __name__ == "__main__":
    main()

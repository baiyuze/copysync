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
VERSION = "1.3.2"
DATE = "2026-10-10"

# 语言代码 → (页面目录, html 的 lang, og:locale, 语言名)
LANGS = {
    "zh": ("", "zh-CN", "zh_CN", "简体中文"),
    "en": ("en/", "en", "en_US", "English"),
    "ja": ("ja/", "ja", "ja_JP", "日本語"),
}

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


def lang_links(code: str, up: str) -> str:
    """除当前语言外的其他语言链接。up 是当前页面到网站根目录的相对路径。"""
    out = []
    for other, (path, hl, _, name) in LANGS.items():
        if other == code:
            continue
        out.append(f'<a href="{up}{path or "./"}" hreflang="{hl}" lang="{hl}" data-lang="{other}">{name}</a>')
    return "\n        ".join(out)


def footer(c: dict, code: str, up: str) -> str:
    langs = []
    for other, (path, hl, _, name) in LANGS.items():
        cur = ' aria-current="page"' if other == code else ""
        langs.append(f'<a href="{up}{path or "./"}" hreflang="{hl}" lang="{hl}" data-lang="{other}"{cur}>{name}</a>')
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
        <a href="{REPO}">GitHub</a>
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
        <a href="{REPO}/blob/main/design/{md_name}">在 GitHub 上查看</a>
        <a href="{REPO}">GitHub</a>
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
    for md_name, page in DESIGN_DOCS.items():
        write(DOCS / "design" / f"{page}.html", design_page(md_name, page))
    write(DOCS / "404.html", not_found())
    write(DOCS / "sitemap.xml", sitemap())


if __name__ == "__main__":
    main()

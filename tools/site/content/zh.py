"""官网文案：简体中文。结构见 tools/site/build.py，三种语言的键必须一致。"""

REPO = "https://github.com/baiyuze/copysync"
SITE = "https://baiyuze.github.io/copysync/"

C = {
    "title": "CopySync：在 Mac 与 Windows 之间同步剪贴板与文件，跨账号跨网络",
    "description": "CopySync 让你在一台电脑上复制，另一台上直接粘贴，Mac 与 Windows 互通。支持文本、图片、文件和文件夹，"
                   "不要求同一个账号，跨网络可用。设备间 P2P 直连、端到端加密，服务器自己部署且看不到内容。免费开源。",
    "og_title": "CopySync：在这台电脑复制，到那台电脑粘贴",
    "og_description": "Mac 与 Windows 互通。文本、图片、文件、文件夹，跨账号、跨网络，P2P 直连、端到端加密，服务器自己部署。免费开源。",
    "og_image": "og.png",
    "skip": "跳到正文",
    "home_label": "CopySync 首页",
    "nav_label": "页面导航",
    "nav_how": "工作原理",
    "nav_install": "安装",
    "nav_faq": "常见问题",
    "nav_server": "部署服务器",
    "lang_label": "语言",

    "h1": "在这台电脑复制，<br>到那台电脑粘贴。",
    "lede": "Mac 与 Windows 互通。文本、图片、文件、整个文件夹都可以，不要求同一个账号，也不要求在同一个网络。"
            "数据在你的设备之间直接传输，经过你自己的服务器时也是加密的。",
    "dl_mac": "下载 Mac 版",
    "dl_win": "下载 Windows 版",
    "dl_note": "版本 1.3.3，免费开源。macOS 13 及以上，Intel 与 Apple 芯片通用；Windows 10、11。",
    "stage_label": "演示：在 MacBook 上按下 Command C，文件化作一张纸，经直连飞到 Windows 笔记本，在那里按下 Ctrl V 粘贴出来",
    "stage_file": "季度汇报-终版.pptx",
    "stage_win": "Windows 笔记本",
    "stage_arc": "直连，端到端加密",
    "replay": "再演示一次",

    "history_h2": "所有设备上复制过的，都在这里",
    "history_intro": "复制记录保留最近 3 天的内容，随时可以放回剪贴板，图片点开就能预览。"
                     "50 MB 以内的文件在你复制的时候就已经传好了，更大的文件需要时点一下再拉取。",
    "history_alt": "CopySync 的复制记录：来自 MacBook Air 和 Mac mini 的文件、图片、文字，大文件显示「拉取到本机」按钮",

    "flow_h2": "按下复制之后",
    "flow_intro": "文字、图片、文件、整个文件夹，都走同一条路。",
    "flow": [
        ("察觉", "Mac 上，后台服务每隔一小段时间看一眼剪贴板的变更计数，只看计数和类型，不需要授权，也不会弹出系统提示；"
                 "Windows 上，复制的那一刻系统就会通知它。"),
        ("传输", "文本随消息一起发出；50 MB 以内的图片和文件打包成流立即推送；更大的文件只发一条记录，等你需要时再传。"),
        ("粘贴", "另一台电脑先把内容存到本地，再写入剪贴板。按下 ⌘V 或 Ctrl+V 粘贴出来的，就是本地的真实文件。"),
    ],

    "why_h2": "为什么不用现有的办法",
    "why_intro": "通用剪贴板要同一个 Apple ID、离得够近；隔空投送每次都要点接收；发给自己要手动下载，文件还留在别人的服务器上。"
                 "一台 Mac 一台 Windows 时，系统自带的办法一个都用不了。",
    "compare_caption": "CopySync 与其他方式的对比",
    "compare_cols": ["CopySync", "通用剪贴板", "隔空投送", "聊天软件"],
    "compare_rows": [
        ("不同 Apple ID 之间", [("yes", "可以"), ("no", "不行"), ("partial", "要对方接收"), ("yes", "可以")]),
        ("Mac 与 Windows 之间", [("yes", "可以"), ("no", "不行"), ("no", "不行"), ("yes", "可以")]),
        ("不在同一个网络（家里和公司）", [("yes", "可以"), ("no", "不行"), ("no", "不行"), ("yes", "可以")]),
        ("复制后直接粘贴，不用点接收", [("yes", "可以"), ("yes", "可以"), ("no", "不行"), ("no", "不行")]),
        ("文件与整个文件夹", [("yes", "可以"), ("partial", "不稳定"), ("yes", "可以"), ("partial", "要先打包")]),
        ("找回之前复制过的内容", [("yes", "可以"), ("no", "不行"), ("no", "不行"), ("partial", "翻聊天记录")]),
        ("内容不经过第三方服务器", [("yes", "不经过"), ("yes", "不经过"), ("yes", "不经过"), ("no", "经过")]),
    ],

    "how_h2": "工作原理",
    "how_intro": "两台电脑之间直接建立加密连接。你的服务器只负责让它们找到彼此，打不通直连的时候再帮忙转发。",
    "diagram": {
        "title": "CopySync 的连接方式",
        "desc": "电脑 A 与电脑 B（Mac 或 Windows）之间是 WebRTC 直连；两台电脑都通过信令连接到你的服务器；直连失败时，数据经服务器的 TURN 中转。",
        "server": "你的服务器",
        "server_sub": "信令、STUN、TURN 中转",
        "sig_note": "信令：交换连接信息",
        "relay_note": "中转：仅在打不通直连时",
        "copy_here": "在这里复制",
        "other": "Mac 或 Windows",
        "paste_here": "在这里粘贴",
        "direct": "WebRTC 直连，DTLS 加密",
        "direct_note": "数据只在两台电脑之间传输",
        "legend_direct": "数据（直连）",
        "legend_relay": "数据（中转，备用）",
        "legend_sig": "信令",
    },

    "uplink_h2": "公司网络也能直连",
    "uplink_intro": "公司双线、运营商地址池这类网络有多个出口，按目标地址选走哪一个。只探测一次只会知道其中一个出口，"
                    "只能经服务器中转。CopySync 会探测每个出口，把地址都告诉对端。",
    "uplink_changed": "做了什么",
    "uplink_points": [
        ("共用一个端口。", "所有连接从同一个本地 UDP 端口收发，向不同服务器问到的地址才对得上。"),
        ("探测每个出口。", "在这个端口上同时询问 7 台 IP 分散的服务器，网络有几个出口就拿到几个地址；被 DNS 污染的地址先过滤掉。"),
        ("全部告诉对端。", "对端逐个尝试，与实际出口吻合的那一个能打通。"),
        ("记住见过的出口。", "某次漏掉的出口如果不改端口，就用「出口 IP + 本地端口」推算补上。"),
        ("定好谁先发包。", "有的路由器收到对方先到的包会占住端口。一方先公布出口地址，另一方先往那边发包再公布自己的；重试时两种顺序轮流用。"),
        ("先中转，再换直连。", "第一次握手落到了中转时，空闲时再打一次洞，打通就换成直连。"),
    ],
    "lab_h3": "用真实的网络拓扑验证",
    "lab_intro": "NAT 打洞实验室用 Linux 网络命名空间与 iptables 搭出这些拓扑，用 CopySync 真实的连接代码各跑一遍，每次提交代码都自动运行。",
    "lab_cols": ["场景", "结果"],
    "lab_rows": [
        ("单出口 ↔ 家用路由器", "直连"),
        ("两个出口，旧版本（对照）", "中转"),
        ("两个出口 ↔ 家用路由器", "直连"),
        ("四个出口 ↔ 家用路由器", "直连"),
        ("两个出口，第二条没探测到，没有历史", "中转"),
        ("两个出口，第二条没探测到，有历史", "直连"),
        ("两个出口 ↔ 两个出口", "直连"),
        ("两个出口 ↔ 端口随机分配的 NAT", "中转"),
        ("一开始打不通，之后网络恢复", "先中转，随后直连"),
        ("公司路由器会为外来的包占住端口", "直连"),
        ("家用路由器会为外来的包占住端口", "先中转，重试一次后直连"),
    ],
    "uplink_measured": "实测：一台在 4 个出口的公司网络里的 Mac 与家里的 Mac，改造前只能中转，速度被服务器带宽限制在 0.6 MB/s；"
                       "改造后直连，连上服务器后约 3 秒建立。",
    "uplink_doc": '测量过程、设计取舍与踩过的坑见<a href="{doc}">技术方案</a>。',

    "security_h2": "服务器看不到你复制了什么",
    "security_points": [
        ("服务器不存账号，也不存内容。", "设备身份是一把 Ed25519 密钥，设备 ID 由公钥算出，服务器不需要用户表就能验证身份。"),
        ("配对时由你核对安全指纹。", "两块屏幕会显示同样的两行：双方各自公钥的指纹。服务器有能力替换转发中的公钥，但替换之后两块屏幕就对不上了，这一步就是防中间人的闸门。"),
        ("之后的每条消息都带签名，", "服务器无法冒充你的设备。直连通道的证书在握手后还会再核对一次，不一致立即断开。"),
        ("中转也只是转发密文。", "TURN 转发的是加密后的数据包，中转凭证按设备单独签发，12 小时后过期。"),
    ],
    "verify_alt": "配对时的「核对安全指纹」对话框，显示 R8NF-2WTC-QL5J，要求与另一台电脑上显示的完全一致",

    "fallback_h2": "出问题的时候会发生什么",
    "fallback_intro": "每个人的网络环境都不一样。下面这些情况都考虑过，也都有对应的测试。",
    "fallbacks": [
        ("两台电脑打不通直连", "常见于对称 NAT 和公司防火墙。自动改走你服务器上的 TURN 中转，内容依然加密，设备列表里显示「中转」。"),
        ("公共 STUN 服务器连不上", "国内访问公共 STUN 时通时断。服务器自带 STUN，客户端优先用它。"),
        ("先连上了中转", "中转握手快，有时会抢在直连之前。连上之后，发起连接的一方会在空闲时再打一次洞，打通就换成直连；正在传文件时不打扰。"),
        ("对方睡眠、断网或重启", "连接悄悄断了也能发现，立即关闭重连。这期间复制的内容会记下来，连上后补发到对方的复制记录里，但不会覆盖对方此刻剪贴板里的内容。"),
        ("服务器断开", "自动重连，间隔从 1 秒开始翻倍、最长 30 秒，并带随机抖动，避免一群设备同时重连。界面实时显示连接状态。"),
        ("文件太大", "超过上限的只同步一条记录，需要时再拉取。拉取时源文件已经被删，会明确告诉你，而不是不明不白地失败。"),
        ("粘贴到只认纯文本的地方", "从网页复制的带格式文字会同时写入 HTML 和纯文本两份，粘贴到哪里都不会冒出标签。"),
        ("文件名在 Windows 上不合法", "Mac 上可以用的冒号、问号等字符，收到时换成对应的全角字符，名字照样看得懂；只差大小写的两个文件也不会互相覆盖。"),
        ("复制的是密码", "Windows 上，密码管理器复制的内容会带一个「不要记录」的标记，CopySync 看到就跳过：不同步，也不进复制记录。"),
        ("还没授权读取剪贴板", "不去读内容，否则后台进程会卡在一个看不见的系统弹窗上。App 里会提示你去授权，并给出直达按钮。"),
        ("后台服务崩溃", "Mac 上由系统的 launchd、Windows 上由 CopySync 自带的守护进程自动拉起，10 秒节流，不会疯狂重启。"),
        ("更新了 App，或者把它挪了位置", "打开 App 时会自动校正登录项，并把后台服务重启到新版本。"),
    ],

    "install_h2": "安装",
    "install_intro": "每台电脑装一个 App，再准备一台两边都能访问到的服务器。两台电脑在同一个局域网的话，服务器放在其中一台 Windows、Mac 或 Linux 电脑上就行，见<a href=\"server/\">局域网部署教程</a>。",
    "mac_steps": [
        "下载 CopySync.dmg，打开后把 CopySync 拖进「应用程序」。",
        "从「应用程序」里打开。<span class=\"aside\">第一次会提示「无法验证开发者」：CopySync 还没有经过 Apple 公证。"
        "到「系统设置 → 隐私与安全性」，点页面下方的「仍要打开」，只需要这一次。</span>",
        "点「启用后台同步」，在「设置」里填入服务器地址 <code>ws://服务器地址:8787/signal</code>，按回车。",
        "系统询问是否允许读取剪贴板时，选择允许。",
    ],
    "win_steps": [
        "运行 CopySync-Setup.exe。装在当前用户目录下，不需要管理员权限。"
        "<span class=\"aside\">CopySync 还没有代码签名，Windows 会提示「Windows 已保护你的电脑」：点「更多信息」，再点「仍要运行」。"
        f"不想运行安装程序的话，可以下载<a href=\"{REPO}/releases/latest/download/CopySync-windows-x64.zip\">便携版</a>。</span>",
        "打开 CopySync，点「启用后台同步」。关掉窗口后它在右下角的托盘里继续运行。",
        "在「设置」里填入服务器地址，按回车。防火墙询问时选择允许。",
    ],
    "server_h3": "服务器",
    "server_intro": "跨公网使用时，一台有公网 IP 的 Linux 服务器，需要 root。ARM 服务器把 <code>amd64</code> 换成 <code>arm64</code>。",
    "copy": "复制",
    "copied": "已复制",
    "public_ip": "你的公网IP",
    "server_after": "脚本会注册开机自启的 systemd 服务，最后打印出 App 里要填的地址。防火墙或云安全组需要放行：",
    "ports_cols": ["端口", "用途"],
    "ports": [("8787/tcp", "信令"), ("3478/udp", "STUN 与 TURN"), ("32768–60999/udp", "中转端口")],
    "server_docker": f'也可以用 Docker 部署，见<a href="{REPO}/blob/main/server/deploy/README.md">服务器部署说明</a>。',
    "server_guide": "局域网部署教程",
    "pairing": "<strong>配对：</strong>一台点「设备 → 添加设备」生成 6 位配对码，另一台输入它；两块屏幕上的两行安全指纹完全一致后确认。"
               "Mac 与 Windows 之间、两台 Windows 之间都一样。",

    "specs_h2": "技术规格",
    "specs": [
        ("版本", '1.3.3，<time datetime="2026-10-11">2026 年 10 月 11 日</time>发布'),
        ("Mac 客户端", "macOS 13 Ventura 及以上；Intel 与 Apple 芯片通用，不需要 Rosetta；安装包约 47 MB"),
        ("Windows 客户端", "Windows 10（21H2 及以上）与 Windows 11，x64；按用户安装，不需要管理员权限；安装包约 23 MB"),
        ("服务器", "Linux x86_64 或 ARM64（systemd 或 Docker）、macOS（Intel 与 Apple 芯片）、Windows 10 与 11 x64（系统服务）；约 15 MB 内存"),
        ("界面语言", "简体中文、English、日本語，默认跟随系统，可在设置里切换"),
        ("同步的内容", "纯文本、带格式的文本、图片与截图、文件、文件夹，可以按类型关闭"),
        ("传输", "WebRTC 数据通道，直连优先，打不通时经 TURN 中转；多出口网络逐个探测每个出口；文件以 tar + zstd 流式打包"),
        ("加密与身份", "DTLS 端到端加密；设备身份为 Ed25519 密钥；信令消息逐条签名"),
        ("配对", "6 位配对码，5 分钟内有效，字符集去掉了 0/O、1/I 等易混字符；配对时人工核对安全指纹"),
        ("自动同步上限", "默认 50 MB，可在 1 MB 到 1 GB 之间调整"),
        ("记录保留", "默认 3 天，可在 1 小时到 30 天之间调整，到期自动清理"),
        ("实测速度", "局域网直连传输 20 MB 文件约 330 毫秒，校验和一致"),
        ("技术栈", "后台服务与服务器用 Go（Windows 上纯 Go 调用 Win32），界面用 Flutter（Mac 与 Windows 共用一套代码），消息定义用 Protocol Buffers"),
        ("许可", f'GNU AGPL-3.0；闭源产品等情况可购买<a href="{REPO}/blob/main/LICENSING.md">商业授权</a>'),
    ],

    "faq_h2": "常见问题",
    "faq": [
        ("UU 远程开启时，为什么不断出现 .uuremote_ 文件？", [
            "UU 远程会在系统剪贴板里写入临时占位文件，旧版 CopySync 把它们当成普通文件反复记录和转发。1.3.2 在发送和接收时过滤已确认的占位命名，启动时清理纯占位记录。正常空文件、隐藏文件和大文件不受此规则影响。50 MB 只决定是否自动传输，与 UU 的 UDP P2P 设置无关。建议两端都升级；其他软件如果也同步剪贴板，仍可能产生普通内容的重复记录。",
        ]),
        ("和 macOS 自带的「通用剪贴板」有什么区别？", [
            "通用剪贴板要求两台设备登录同一个 Apple ID、开着蓝牙和 Wi‑Fi，并且离得近，也只能在苹果设备之间用。"
            "CopySync 只要求两台电脑都能连上你部署的服务器：不同 Apple ID、一台在家一台在公司、一台 Mac 一台 Windows，都可以。"
            "另外它有复制记录，大文件也可以按需拉取。",
            "如果你的两台 Mac 是同一个 Apple ID、总在一起、主要复制文字，系统自带的已经够用了。",
        ]),
        ("服务器会看到我复制的内容吗？", [
            "不会。直连时数据根本不经过服务器；需要中转时，服务器转发的是 DTLS 加密后的数据包。服务器不存账号，也不存任何复制内容。代码是开源的，可以自己核对。",
        ]),
        ("公司网络（双线、多出口）能直连吗？", [
            "能。这类网络按目标地址选出口，CopySync 会在同一个端口上向多台服务器探测，拿到每个出口的地址一并发给对端，对端逐个尝试。"
            "实测公司 4 个出口与家用宽带之间可以直连。对端是端口随机分配的 NAT 时仍会走中转。",
        ]),
        ("必须自己准备服务器吗？", [
            "是的，CopySync 不提供公共服务器，这样就不会有任何人的数据经过别人的机器。两台电脑在同一个局域网时，在其中一台常开的 Windows、Mac 或 Linux 电脑上运行服务器即可，见<a href=\"server/\">局域网部署教程</a>；"
            "跨网络使用需要一台有公网 IP 的服务器，最低配置的云主机就够。",
        ]),
        ("支持 Windows 吗？iPhone 呢？", [
            "支持 Windows 10 与 Windows 11（x64），与 Mac 互通。复制的文件在 Windows 上粘贴出来就是真实文件，能直接粘进资源管理器、微信、Office。",
            "iPhone 和 iPad 暂不支持：iOS 不允许 App 在后台监听剪贴板，做不到复制后立即同步。",
        ]),
        ("为什么第一次打开时系统会拦一下？", [
            "CopySync 还没有经过 Apple 公证，Windows 版也还没有代码签名，两者都需要每年付费的证书。",
            "Mac 上到「系统设置 → 隐私与安全性」，点页面下方的「仍要打开」；Windows 上在「Windows 已保护你的电脑」提示里点「更多信息」，再点「仍要运行」。都只需要这一次。",
        ]),
        ("大文件怎么处理？", [
            "默认 50 MB 以内的内容在你复制时就推送到另一台电脑，那边可以直接粘贴；更大的文件只同步一条记录，需要时在复制记录里点「拉取到本机」再传。上限可以在设置里调整。",
        ]),
        ("免费吗？公司能用吗？", [
            "免费。CopySync 以 GNU AGPL-3.0 开源，个人和公司都可以免费使用和修改。分发它、或通过网络提供修改过的版本时，需要以同样的许可公开源码。",
            f'要放进闭源产品、修改后不想公开、或者作为托管服务运营，可以购买<a href="{REPO}/blob/main/LICENSING.md">商业授权</a>。',
        ]),
    ],

    "closing": "两台电脑，一个剪贴板。",
    "source": "在 GitHub 上查看源码",
    "footer_license": "CopySync 以 GNU AGPL-3.0 开源，另提供商业授权。",
    "footer_label": "页脚链接",
    "footer_links": [
        ("所有版本", f"{REPO}/releases"),
        ("更新记录", f"{REPO}/blob/main/CHANGELOG.md"),
        ("技术方案", f"{SITE}design/nat-traversal.html"),
        ("商业授权", f"{REPO}/blob/main/LICENSING.md"),
        ("反馈问题", f"{REPO}/issues"),
    ],

    # ── 「在局域网里部署服务器」教程（server/ 页面） ──
    "guide": {
        "title": "在局域网里部署 CopySync 服务器：Windows、Mac、Linux 教程",
        "description": "两台电脑在同一个网络里时，不需要云服务器：在一台常开的 Windows、Mac 或 Linux 电脑上运行 CopySync 的信令与 TURN 服务器，"
                       "下载后运行一个脚本即可，开机自动运行。附客户端设置与排错。",
        "h1": "在局域网里部署服务器",
        "lede": "家里或办公室的电脑在同一个网络里，就不需要云服务器。找一台常开的电脑，下载、运行一个脚本，"
                "再把它显示的地址填进每台电脑的 CopySync 里。",
        "toc": [("prepare", "准备"), ("install", "安装服务器"), ("check", "检查"), ("clients", "填地址"), ("faq", "常见问题")],

        "overview_h2": "服务器只负责牵线",
        "overview_intro": "每台电脑一直连着服务器，靠它找到对方、交换连接信息。连上之后，复制的内容在两台电脑之间直接传输；"
                          "只有直连打不通时，才经过它中转，内容依然加密。",
        "diagram": {
            "title": "局域网里的 CopySync",
            "desc": "同一个局域网里，Mac 与 Windows 电脑通过信令连接到地址为 192.168.1.20 的服务器，复制的内容在两台电脑之间直连传输；"
                    "直连不通时经服务器的 TURN 中转。",
            "lan": "局域网：同一个 Wi‑Fi 或有线网络",
            "server": "服务器",
            "sig_note": "信令 TCP 8787",
            "relay_note": "中转 UDP 3478，直连不通时",
            "client": "装着 CopySync",
            "direct": "直连，端到端加密",
            "direct_note": "复制的内容只在两台电脑之间传输",
            "legend_direct": "数据（直连）",
            "legend_relay": "数据（中转，备用）",
            "legend_sig": "信令",
        },
        "overview_points": [
            "<strong>信令，TCP 8787。</strong>每台电脑一直连着它，用来发现对方、配对、交换连接信息。",
            "<strong>直连。</strong>在同一个局域网里，内容几乎总是在两台电脑之间直接传输，速度就是局域网的速度。",
            "<strong>中转，UDP 3478。</strong>访客网络、开了 AP 隔离的 Wi‑Fi、不同网段之间互相访问不到时才会用到。",
        ],

        "prep_h2": "开始之前",
        "prep": [
            ("一台常开的电脑", "Windows 10、11，macOS 13 及以上，或者任意 Linux 都行，两台电脑中的一台也可以。"
                             "它关机或睡眠时同步会暂停，醒来后自动恢复。NAS、软路由、Mac mini 这类一直开着的设备最合适，内存只占约 15 MB。"),
            ("一个固定的地址", "在路由器的「DHCP 静态分配」（也叫「地址保留」「IP 与 MAC 绑定」）里给这台电脑固定一个地址。"
                             "不然路由器换了地址，每台电脑都要重新填。"),
            ("在同一个网络里", "设备连同一个 Wi‑Fi 或有线网络就行。访客网络、开了「AP 隔离」的 Wi‑Fi 不让设备互相访问，需要换个网络或关掉隔离。"),
            ("离开这个网络就不同步", "服务器只在局域网里访问得到，笔记本带出门后，回来才会继续同步。想在任何地方都能同步，"
                                   "需要一台有公网 IP 的服务器，见首页的「安装」。"),
        ],

        "install_h2": "安装服务器",
        "install_intro": "选服务器所在电脑的系统。安装脚本会自动找到本机的局域网地址，把服务器设成开机自动运行，最后显示客户端要填的地址。",
        "os_label": "服务器所在电脑的系统",
        "manage_caption": "以后会用到的",
        "win": {
            "steps": [
                f'下载 <a href="{REPO}/releases/latest/download/copysync-server-windows-amd64.zip">copysync-server-windows-amd64.zip</a>，'
                "右键选「全部解压缩」。",
                "打开解压出的文件夹，双击 <code>install.cmd</code>。"
                "<span class=\"aside\">提示「Windows 已保护你的电脑」时，点「更多信息」，再点「仍要运行」；弹出用户账户控制时点「是」。"
                "服务器还没有代码签名，所以会这样问一次。</span>",
                "安装在新打开的窗口里进行。最后显示的 <code>ws://…/signal</code> 就是客户端要填的地址。",
            ],
            "download": "下载 Windows 版服务器",
            "notes": [
                "<strong>作为 Windows 服务运行：</strong>开机就启动，不用登录，也没有窗口。在「服务」里叫 CopySync Server。",
                "<strong>防火墙已经放行：</strong>脚本添加了一条只针对 copysync-server.exe 的入站规则，不用自己开端口。",
                "<strong>地址认错了</strong>（比如装了虚拟机或 VPN，有好几块网卡）：在解压出的文件夹的地址栏里输入 <code>cmd</code> 回车，"
                "运行 <code>install.cmd 192.168.1.20</code>，换成你的地址。",
                "支持 Windows 10、11，x64。",
            ],
            "double_click": "双击",
            "explorer_label": "解压出的文件夹里有 copysync-server.exe、install.cmd、install.ps1、README.md、uninstall.cmd，双击其中的 install.cmd",
            "term_title": "管理员: Windows PowerShell",
            "term_label": "install.cmd 运行完的窗口：依次停止旧版本、安装、注册服务、放行防火墙、启动，最后显示客户端要填的地址 ws://192.168.1.20:8787/signal",
            "manage": [
                ("状态", "管理员 PowerShell 里运行 <code>Get-Service CopySyncServer</code>，或者在「服务」里找 CopySync Server"),
                ("日志", "<code>C:\\ProgramData\\CopySync Server\\server.log</code>"),
                ("重启", "管理员 PowerShell 里运行 <code>Restart-Service CopySyncServer</code>"),
                ("升级", "下载新版本，解压后再双击 <code>install.cmd</code>，原来的设置保留"),
                ("卸载", "双击 <code>uninstall.cmd</code>，服务、防火墙规则和文件一并删除"),
            ],
        },
        "mac": {
            "steps": [
                "打开「终端」，逐行运行：",
                "最后显示的 <code>ws://…/signal</code> 就是客户端要填的地址。",
            ],
            "notes": [
                "<strong>不需要管理员权限。</strong>登录这台 Mac 后自动在后台运行，意外退出会被系统重新拉起。",
                "<strong>这台 Mac 不能睡眠：</strong>在「系统设置 → 能耗」（笔记本在「电池 → 选项」）里打开防止自动进入睡眠的选项。",
                "<strong>开着系统防火墙的话，</strong>脚本最后会给出放行 copysync-server 的命令。",
                "Intel 与 Apple 芯片通用，macOS 13 及以上。",
            ],
            "term_label": "在终端里运行 ./install.sh 的输出：安装、写入 LaunchAgent、启动，最后显示客户端要填的地址 ws://192.168.1.20:8787/signal",
            "manage": [
                ("日志", "<code>~/Library/Logs/CopySync/server.log</code>"),
                ("重启", "<code>launchctl kickstart -k gui/$(id -u)/com.copysync.server</code>"),
                ("升级", "下载新版本，解压后再运行 <code>./install.sh</code>，原来的设置保留"),
                ("卸载", "<code>./install.sh uninstall</code>"),
            ],
        },
        "linux": {
            "steps": [
                "在这台机器上（或者 ssh 登录后）逐行运行：<span class=\"aside\">ARM 设备，比如装了 64 位系统的树莓派，把 <code>amd64</code> 换成 <code>arm64</code>。</span>",
                "最后显示的 <code>ws://…/signal</code> 就是客户端要填的地址。",
            ],
            "notes": [
                "<strong>注册为 systemd 服务，</strong>开机自动运行。需要 root。",
                "<strong>开着 ufw 或 firewalld 的话，</strong>脚本最后会给出放行端口的命令。",
                "<strong>也可以用 Docker：</strong>把 <code>env.example</code> 复制为 <code>.env</code>，<code>TURN_IP</code> 填这台机器的局域网地址，"
                "再运行 <code>docker compose up -d</code>。必须用 host 网络，所以只适用于 Linux。",
            ],
            "term_label": "运行 sudo ./install.sh 的输出：安装、写入 systemd 单元、生成配置、启动，最后显示客户端要填的地址 ws://192.168.1.20:8787/signal",
            "manage": [
                ("状态", "<code>systemctl status copysync-server</code>"),
                ("日志", "<code>journalctl -u copysync-server -f</code>"),
                ("重启", "<code>sudo systemctl restart copysync-server</code>"),
                ("配置", "<code>/etc/copysync/server.env</code>，改完重启服务"),
                ("升级", "下载新版本，解压后再运行 <code>sudo ./install.sh</code>，原来的配置保留"),
                ("卸载", "<code>sudo ./install.sh uninstall</code>"),
            ],
        },

        "check_h2": "检查其他电脑能不能访问到它",
        "check_intro": "在另一台电脑的浏览器里打开 <code>http://192.168.1.20:8787/healthz</code>（换成你的地址），看到下面这一行，网络就是通的。",
        "check_label": "浏览器打开 http://192.168.1.20:8787/healthz，页面显示服务器的状态 ok 与版本号",
        "check_after": "打不开的话，看下面的<a href=\"#faq\">常见问题</a>。",

        "clients_h2": "在每台电脑上填地址",
        "clients_steps": [
            ("打开设置", "在每台电脑上打开 CopySync，进入「设置」。服务器所在的那台电脑也一样。"),
            ("填入地址", "在「信令服务器地址」里填 <code>ws://192.168.1.20:8787/signal</code>（换成你的地址），按回车保存。"),
            ("看连接状态", "「连接状态」变成绿色的「已连接」就好了。接着在「设备」里配对，之后复制就会同步。"),
        ],
        "clients_alt": "CopySync 的设置页：信令服务器地址填着 ws://192.168.1.20:8787/signal，连接状态显示已连接",
        "clients_after": '还没装 App？见首页的<a href="{home}#install">安装</a>。',

        "faq_h2": "常见问题",
        "faq": [
            ("连接状态一直是「未连接」？", [
                "先核对地址：以 <code>ws://</code> 开头，端口是 8787，以 <code>/signal</code> 结尾，填完按过回车。",
                "再在这台电脑的浏览器里打开 <code>http://服务器地址:8787/healthz</code>。打不开，说明到服务器的网络不通："
                "服务器没在运行、那台电脑睡眠了、被防火墙拦住了，或者它的地址变了。",
                "Windows 的安装脚本已经放行了防火墙；Mac 开着系统防火墙时，运行安装脚本最后给出的放行命令；Linux 上的 ufw、firewalld 也一样。",
            ]),
            ("设备之间显示「中转」，不是「直连」？", [
                "同一个局域网里一般都是直连。显示中转，说明两台电脑之间的 UDP 不通，常见于访客网络、开了 AP 隔离的 Wi‑Fi，"
                "或者两台电脑在不同网段、中间有防火墙。中转经过服务器，内容同样加密，只是速度受那台电脑的网络限制。",
            ]),
            ("服务器的地址变了怎么办？", [
                "带上新地址再运行一次安装脚本：Windows 上是 <code>install.cmd 新地址</code>，Mac 上是 <code>./install.sh 新地址</code>；"
                "Linux 上编辑 <code>/etc/copysync/server.env</code> 里的 <code>-turn-ip</code>，再重启服务。然后在每台电脑上改掉服务器地址。",
                "在路由器里给它固定一个地址，就不会再遇到这件事。",
            ]),
            ("离开局域网以后还能同步吗？", [
                "不能，服务器只在局域网里访问得到。要在外面也能用，可以把服务器部署到有公网 IP 的机器上（首页的「安装」里有步骤），"
                "或者用 Tailscale、ZeroTier 这类组网工具把设备连进同一个虚拟局域网，服务器地址填它在虚拟网里的地址。",
            ]),
            ("用国内的云服务器，填域名连不上、填 IP 却能连？", [
                "域名没有备案时，国内的云厂商会拦截访问它的明文 HTTP 请求，返回一个跳转到拦截页的响应，实测 8787 这样的非常用端口也一样。"
                "WebSocket 握手也是一个 HTTP 请求，所以同样被拦，客户端就显示未连接。",
                "服务器地址里直接填 IP 就行，比如 <code>ws://服务器IP:8787/signal</code>；或者给域名备案。",
            ]),
            ("服务器所在的那台电脑也要用 CopySync，地址填什么？", [
                "填同样的局域网地址就行，填 <code>ws://127.0.0.1:8787/signal</code> 也可以。服务器和 App 互不影响。",
            ]),
        ],

        "closing": "一台常开的电脑，就够了。",
        "source": "服务器的完整部署说明（含公网与 Docker）",
    },

    "ld_description": "在 Mac 与 Windows 电脑之间同步剪贴板与文件的开源工具。复制文本、图片、文件或文件夹后，另一台电脑可直接粘贴。"
                      "设备间 WebRTC 直连、DTLS 端到端加密，打不通时经用户自建服务器的 TURN 中转。",
    "ld_os": "macOS 13 及以上（Intel 与 Apple 芯片）；Windows 10、Windows 11（x64）",
    "ld_currency": "CNY",
    "ld_features": [
        "同步纯文本、带格式的文本、图片、文件与文件夹",
        "Mac 与 Windows 互通，不要求同一个账号，跨网络可用",
        "WebRTC 直连，打不通时自动经 TURN 中转",
        "DTLS 端到端加密，配对时人工核对安全指纹",
        "50 MB 以内复制即推送，更大的文件按需拉取",
        "最近 3 天的复制记录，图片可以预览",
        "界面支持简体中文、英文、日文",
    ],
}

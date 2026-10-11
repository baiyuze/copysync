"""Site copy: English. Same keys as zh.py; see tools/site/build.py for the structure."""

REPO = "https://github.com/baiyuze/copysync"
SITE = "https://baiyuze.github.io/copysync/"

C = {
    "title": "CopySync: sync the clipboard and files between Mac and Windows, across accounts and networks",
    "description": "CopySync lets you copy on one computer and paste on another, Mac or Windows: text, images, files and folders. "
                   "No shared account or network required. Peer-to-peer, end-to-end encrypted, with a server you host that "
                   "can't see your content. Free and open source.",
    "og_title": "CopySync: copy on this computer, paste on that one",
    "og_description": "Mac and Windows. Text, images, files and folders across accounts and networks; peer-to-peer, end-to-end "
                      "encrypted, self-hosted server. Free and open source.",
    "og_image": "og-en.png",
    "skip": "Skip to content",
    "home_label": "CopySync home",
    "nav_label": "Page",
    "nav_start": "Quick start",
    "nav_how": "How it works",
    "nav_faq": "FAQ",
    "nav_server": "Self-host",
    "lang_label": "Language",

    "h1": "Copy on this computer.<br>Paste on that one.",
    "lede": "Mac and Windows. Text, images, files, whole folders. No shared account required, no shared network required. "
            "Data moves directly between your devices, and stays encrypted even when it passes through your own server.",
    "dl_mac": "Download for Mac",
    "dl_win": "Download for Windows",
    "dl_note": "Version 1.3.4, free and open source. macOS 13 or later, Intel and Apple silicon; Windows 10 and 11.",
    "start_link": "New here? Set it up in three steps →",
    "stage_label": "Demo: Command C on a MacBook turns a file into a sheet that flies over a direct connection to a Windows "
                   "laptop, where Ctrl V pastes it",
    "stage_file": "Q3 review - final.pptx",
    "stage_win": "Windows laptop",
    "stage_arc": "Direct, end-to-end encrypted",
    "replay": "Play again",

    # Three ways to run the server, on equal footing; the reader picks one. Home cards use name, fit, summary, go;
    # the quick start tabs use name, facts, steps, address, go; the install section uses name, fit.
    # {cloud}, {server} and {ts} in steps become links to the matching guide.
    "paths_h2": "Where to run the server",
    "paths_intro": "Your computers find each other through a server; nothing you copy is stored on it. There's no public server, "
                   "so pick one of three ways to run your own.",
    "paths_after": "Not sure? The <a href=\"{start}#server\">quick start</a> compares all three.",
    "plans_label": "Ways to run the server",
    "plans": {
        "cloud": {
            "name": "Public server",
            "fit": "Syncs from anywhere",
            "summary": "A Linux VM with a public IP; the smallest one is enough. Every computer connects to it directly, at home, "
                       "at the office or anywhere else, with nothing extra to install.",
            "go": "Public server guide →",
            "facts": [
                ("Best for", "Computers spread across networks that should sync from anywhere. The simplest way across networks."),
                ("You need", "A Linux VM with a public IP, as root; the smallest one is enough."),
                ("Keep in mind", "The VM costs money. When a direct connection fails, the relay is limited by the server's "
                                 "bandwidth. In mainland China, use the IP unless the domain has an ICP filing."),
            ],
            "steps": [
                "In your cloud provider's security group, open TCP 8787, UDP 3478 and UDP 32768–60999.",
                "Log in with ssh, download the install script and run it with the public IP: <code>sudo ./install.sh PUBLIC_IP</code>.",
                "Open <code>http://PUBLIC_IP:8787/healthz</code> in a browser; ok means it's running.",
            ],
            "address": "Address for every computer: <code>ws://PUBLIC_IP:8787/signal</code>",
        },
        "lan": {
            "name": "LAN",
            "fit": "All computers on one network",
            "summary": "At home or in the office, on the same Wi‑Fi or wired network. Run the server on one of your always-on "
                       "Windows, Mac or Linux computers with one script. Free.",
            "go": "LAN guide →",
            "facts": [
                ("Best for", "Computers on the same Wi‑Fi or wired network, such as at home or in the office."),
                ("You need", "An always-on Windows, Mac or Linux computer; one of the computers you sync is fine."),
                ("Keep in mind", "No sync away from that network; it resumes when you're back. Sync also pauses while the "
                                 "server computer is off or asleep."),
            ],
            "steps": [
                "Reserve an address for the server computer in your router.",
                "Download the server for its system and run the install script: double-click <code>install.cmd</code> on "
                "<a href=\"{server}#windows\">Windows</a>, run <code>./install.sh</code> on <a href=\"{server}#macos\">macOS</a> "
                "or <a href=\"{server}#linux\">Linux</a>.",
                "The script finds the computer's LAN address and prints the address for the apps.",
            ],
            "address": "Address for every computer: <code>ws://LAN_ADDRESS:8787/signal</code>, for example <code>ws://192.168.1.20:8787/signal</code>",
        },
        "ts": {
            "name": "LAN + Tailscale",
            "fit": "Across networks, no public server",
            "summary": "Run the server on an always-on computer at home, then install Tailscale on each computer under the same "
                       "account. They act as if they shared a LAN.",
            "go": "Tailscale guide →",
            "facts": [
                ("Best for", "Computers on different networks, without a public server."),
                ("You need", "An always-on computer for the server, and Tailscale on every computer under one account."),
                ("Keep in mind", "Every computer that runs CopySync needs Tailscale, and the one at home stays on. When "
                                 "Tailscale can't connect peer to peer it relays through DERP servers, which can be slow."),
            ],
            "steps": [
                "Install Tailscale on each computer and sign in to the same account.",
                "On the server computer, run <code>tailscale ip -4</code> and note the address starting with 100.",
                "Install the server with that address, as in <code>install.cmd 100.101.102.103</code> or <code>./install.sh 100.101.102.103</code>.",
            ],
            "address": "Address for every computer: <code>ws://TAILSCALE_ADDRESS:8787/signal</code>",
        },
    },

    "history_h2": "Everything you copied, on every computer",
    "history_intro": "History keeps the last 3 days, ready to go back on the clipboard; click an image to preview it. "
                     "Files under 50 MB have already arrived by the time you switch computers, and larger ones are a click away.",
    "history_alt": "CopySync history: files, images and text from a MacBook Air and a Mac mini; a large file shows a button "
                   "to pull it to this device",

    "flow_h2": "What happens when you copy",
    "flow_intro": "Text, images, files and whole folders all take the same path.",
    "flow": [
        ("Notice", "On Mac, the background service checks the clipboard's change count at short intervals; reading the count "
                   "and types needs no permission and never shows a prompt. On Windows, the system tells it the moment you copy."),
        ("Send", "Text travels inside the message. Images and files under 50 MB are streamed right away. Larger files are "
                 "announced as a record and sent when you ask for them."),
        ("Paste", "The other computer stores the content locally, then puts it on the clipboard. What you paste with ⌘V or "
                  "Ctrl+V is a real local file."),
    ],

    "why_h2": "Why not what's already there",
    "why_intro": "Universal Clipboard needs the same Apple ID and both Macs close together. AirDrop asks you to accept every "
                 "time. Messaging yourself means manual downloads and files left on someone else's server. With one Mac and "
                 "one Windows PC, none of the built-in options work at all.",
    "compare_caption": "CopySync compared with other options",
    "compare_cols": ["CopySync", "Universal Clipboard", "AirDrop", "Chat apps"],
    "compare_rows": [
        ("Different Apple IDs", [("yes", "Yes"), ("no", "No"), ("partial", "Needs accepting"), ("yes", "Yes")]),
        ("Between Mac and Windows", [("yes", "Yes"), ("no", "No"), ("no", "No"), ("yes", "Yes")]),
        ("Different networks (home and office)", [("yes", "Yes"), ("no", "No"), ("no", "No"), ("yes", "Yes")]),
        ("Paste right away, nothing to accept", [("yes", "Yes"), ("yes", "Yes"), ("no", "No"), ("no", "No")]),
        ("Files and whole folders", [("yes", "Yes"), ("partial", "Unreliable"), ("yes", "Yes"), ("partial", "Zip first")]),
        ("Find something you copied earlier", [("yes", "Yes"), ("no", "No"), ("no", "No"), ("partial", "Scroll back")]),
        ("Content stays off third-party servers", [("yes", "Yes"), ("yes", "Yes"), ("yes", "Yes"), ("no", "No")]),
    ],

    "how_h2": "How it works",
    "how_intro": "Your two computers open an encrypted connection to each other. Your server only helps them find one "
                 "another, and relays traffic when a direct path can't be made.",
    "diagram": {
        "title": "How CopySync connects",
        "desc": "Computer A and computer B (Mac or Windows) connect directly over WebRTC. Both connect to your server for "
                "signaling. If the direct connection fails, data is relayed through the server's TURN service.",
        "server": "Your server",
        "server_sub": "Signaling, STUN, TURN relay",
        "sig_note": "Signaling: connection details",
        "relay_note": "Relay: only if direct fails",
        "copy_here": "You copy here",
        "other": "Mac or Windows",
        "paste_here": "You paste here",
        "direct": "Direct WebRTC, DTLS encrypted",
        "direct_note": "Data moves only between your computers",
        "legend_direct": "Data (direct)",
        "legend_relay": "Data (relay, fallback)",
        "legend_sig": "Signaling",
    },

    "uplink_h2": "Direct on office networks too",
    "uplink_intro": "Office networks with two ISPs, or carriers with a pool of NAT addresses, have several uplinks and pick one "
                    "per destination. A single probe only reveals one of them, so connections end up on the relay. CopySync "
                    "probes every uplink and tells the peer about all of them.",
    "uplink_changed": "What it does",
    "uplink_points": [
        ("One shared port.", " All connections send and receive on one local UDP port, so the addresses learned from different servers line up."),
        ("Every uplink probed.", " Seven IP-diverse servers are queried from that port at once, yielding one address per uplink. Addresses poisoned by DNS are dropped first."),
        ("All of them sent to the peer.", " The peer tries each address; the one matching the real uplink gets through."),
        ("Uplinks remembered.", " If a round misses an uplink that keeps ports unchanged, its address is filled in as uplink IP plus local port."),
        ("A rule for who sends first.", " Some routers lock a port when the other side's packet arrives first. One side announces its addresses first; the other sends to them before announcing its own. Retries alternate the order."),
        ("Relay first, direct later.", " If the first handshake lands on the relay, CopySync punches again while idle and switches to direct when that works."),
    ],
    "lab_h3": "Checked against real network topologies",
    "lab_intro": "The NAT lab builds these topologies with Linux network namespaces and iptables and runs CopySync's actual "
                 "connection code through each. It runs on every commit.",
    "lab_cols": ["Scenario", "Result"],
    "lab_rows": [
        ("Single uplink ↔ home router", "Direct"),
        ("Two uplinks, previous version (control)", "Relay"),
        ("Two uplinks ↔ home router", "Direct"),
        ("Four uplinks ↔ home router", "Direct"),
        ("Two uplinks, second not probed, no history", "Relay"),
        ("Two uplinks, second not probed, with history", "Direct"),
        ("Two uplinks ↔ two uplinks", "Direct"),
        ("Two uplinks ↔ NAT that randomizes ports", "Relay"),
        ("Blocked at first, network recovers later", "Relay, then direct"),
        ("Office router locks a port for unsolicited packets", "Direct"),
        ("Home router locks a port for unsolicited packets", "Relay, then direct after one retry"),
    ],
    "uplink_measured": "Measured: a Mac on a 4-uplink office network and one at home could only relay before, capped at "
                       "0.6 MB/s by the server's bandwidth. Now they connect directly, about 3 seconds after reaching the server.",
    "uplink_doc": 'Measurements, design trade-offs and lessons learned are in the <a href="{doc}">design document</a> (in Chinese).',

    "security_h2": "Your server can't see what you copy",
    "security_points": [
        ("No accounts, no stored content.", " A device's identity is an Ed25519 key, and its ID is derived from the public key, so the server verifies identity without a user table."),
        ("You compare fingerprints when pairing.", " Both screens show the same two lines: the fingerprint of each device's key. The server could swap keys in transit, but then the two screens would no longer match. That check is what stops a man in the middle."),
        ("Every message after that is signed,", " so the server can't impersonate your device. The direct channel's certificate is checked again after the handshake; a mismatch drops the connection."),
        ("The relay only forwards ciphertext.", " TURN carries encrypted packets, and relay credentials are issued per device and expire after 12 hours."),
    ],
    "verify_alt": "The fingerprint check during pairing shows R8NF-2WTC-QL5J and asks you to confirm it matches the other computer exactly",

    "fallback_h2": "When things go wrong",
    "fallback_intro": "Every network is different. Each of these cases was considered, and each has a test.",
    "fallbacks": [
        ("No direct path between the computers", "Common behind symmetric NAT and corporate firewalls. Traffic switches to the TURN relay on your server, still encrypted, and the device list shows \"Relay\"."),
        ("Public STUN servers unreachable", "Public STUN is unreliable in some regions, including mainland China. Your server runs its own STUN, and clients prefer it."),
        ("Connected through the relay first", "The relay handshake is fast and sometimes beats the direct path. Once connected, the side that started the connection punches again while idle and switches to direct when that works, never in the middle of a transfer."),
        ("The other computer sleeps, drops off or restarts", "A connection that dies silently is noticed, closed and reconnected. Anything you copy meanwhile is kept and added to the other computer's history once it's back, without overwriting what's on its clipboard now."),
        ("Server connection drops", "Reconnects automatically: the delay starts at 1 second and doubles up to 30, with jitter so devices don't all reconnect at once. The app shows the status live."),
        ("File too large", "Above the limit only a record is synced, pulled when you need it. If the source file is gone by then, you get a clear message instead of a vague failure."),
        ("Pasting into plain-text fields", "Rich text copied from the web is written as both HTML and plain text, so no markup leaks into what you paste."),
        ("A file name Windows doesn't allow", "Colons, question marks and a few other characters that are fine on a Mac arrive as the matching full-width characters, so the name still reads the same. Two names that differ only in case never overwrite each other."),
        ("Copying a password", "On Windows, password managers mark what they copy as \"don't record\". CopySync skips it: not synced, not kept in the history."),
        ("Clipboard access not granted yet", "Content isn't read, which would block on a prompt a background process can't show. The app asks you to grant access and links straight to the setting."),
        ("Background service crashes", "launchd on Mac, and CopySync's own supervisor on Windows, restart it, throttled to once every 10 seconds."),
        ("App updated or moved", "On launch the app re-points its login item and restarts the service on the new version."),
    ],

    "install_h2": "Install",
    "install_intro": "Install the app on each computer, and set up a server both of them can reach. First time? Follow the "
                     "<a href=\"{start}\">quick start</a>.",
    "mac_steps": [
        "Download CopySync.dmg, open it and drag CopySync into Applications.",
        "Open it from Applications.<span class=\"aside\">The first time, macOS says it can't verify the developer: CopySync "
        "isn't notarized by Apple yet. Go to System Settings → Privacy &amp; Security and click Open Anyway near the bottom. "
        "You only do this once.</span>",
        "Click Turn on background sync. In Settings, enter your server as <code>ws://server-address:8787/signal</code> and press Return.",
        "When macOS asks whether CopySync may read the clipboard, allow it.",
    ],
    "win_steps": [
        "Run CopySync-Setup.exe. It installs for the current user, no administrator rights needed."
        "<span class=\"aside\">CopySync isn't code-signed yet, so Windows says \"Windows protected your PC\": click More info, "
        f"then Run anyway. To skip the installer, get the <a href=\"{REPO}/releases/latest/download/CopySync-windows-x64.zip\">portable zip</a>.</span>",
        "Open CopySync and click Turn on background sync. After you close the window it keeps running in the notification area.",
        "In Settings, enter your server address and press Enter. Allow it if Windows Firewall asks.",
    ],
    "server_h3": "Server",
    "server_intro": "Three ways to run it; pick the one that fits. Each has a step-by-step guide:",
    "server_compare": "Compare the three",
    "copy": "Copy",
    "copied": "Copied",
    "pairing": "<strong>Pairing:</strong> on one computer, choose Devices → Add device to show a 6-character code; enter it on "
               "the other, and confirm once both screens show the same two fingerprints. Mac with Windows, or two Windows PCs, "
               "works the same way.",

    "specs_h2": "Specifications",
    "specs": [
        ("Version", '1.3.4, released <time datetime="2026-10-11">October 11, 2026</time>'),
        ("Mac app", "macOS 13 Ventura or later; universal binary for Intel and Apple silicon, no Rosetta; about 47 MB download"),
        ("Windows app", "Windows 10 (21H2 or later) and Windows 11, x64; installs per user, no administrator rights; about 23 MB download"),
        ("Server", "Linux on x86_64 or ARM64 (systemd or Docker), macOS (Intel and Apple silicon), Windows 10 and 11 x64 (as a service); about 15 MB of RAM"),
        ("Languages", "English, Simplified Chinese and Japanese; follows the system by default, switchable in Settings"),
        ("Content", "Plain text, rich text, images and screenshots, files, folders; each type can be turned off"),
        ("Transport", "WebRTC data channels, direct first with TURN relay fallback; every uplink probed on multi-uplink networks; files streamed as tar + zstd"),
        ("Encryption and identity", "End-to-end DTLS; Ed25519 device keys; every signaling message signed"),
        ("Pairing", "6-character code valid for 5 minutes, using an alphabet without look-alikes such as 0/O and 1/I; fingerprints compared by you"),
        ("Auto-sync limit", "50 MB by default, adjustable from 1 MB to 1 GB"),
        ("History", "3 days by default, adjustable from 1 hour to 30 days; expired items removed automatically"),
        ("Measured speed", "A 20 MB file transfers in about 330 ms over a direct LAN connection, checksums matching"),
        ("Built with", "Go for the background service and server (pure Go calling Win32 on Windows), Flutter for the app (one codebase for Mac and Windows), Protocol Buffers for messages"),
        ("License", f'GNU AGPL-3.0; a <a href="{REPO}/blob/main/LICENSING.md">commercial license</a> is available for closed-source use'),
    ],

    "faq_h2": "Questions",
    "faq": [
        ("Why do .uuremote_ files keep appearing while UU Remote is running?", [
            "UU Remote writes temporary placeholder files to the system clipboard. Older CopySync clients recorded and forwarded them as ordinary files. Version 1.3.2 filters the confirmed placeholder naming pattern on send and receive, and removes placeholder-only history on startup. Ordinary empty, hidden and large files are unaffected. The 50 MB threshold only controls automatic transfer; UU’s UDP P2P setting is unrelated. Upgrade both clients. Other clipboard sync tools may still duplicate ordinary content.",
        ]),
        ("How is this different from Universal Clipboard?", [
            "Universal Clipboard needs both devices on the same Apple ID, with Bluetooth and Wi‑Fi on, close to each other, and "
            "works only between Apple devices. CopySync only needs both computers to reach your server, so different Apple IDs "
            "work, one Mac and one Windows PC work, and so does one at home and one at the office. It also keeps a history and "
            "lets you pull large files on demand.",
            "If your Macs share an Apple ID, sit side by side, and you mostly copy text, the built-in feature is already enough.",
        ]),
        ("Can the server see what I copy?", [
            "No. With a direct connection the data never touches the server. When it relays, it forwards DTLS-encrypted packets. "
            "It keeps no accounts and stores no clipboard content. The code is open, so you can check.",
        ]),
        ("Does it connect directly on office networks with several uplinks?", [
            "Yes. These networks pick an uplink per destination. CopySync probes several servers from one port to learn its "
            "address on every uplink and sends them all to the peer, which tries each one. A Mac behind a 4-uplink office "
            "network connects directly to a home connection in our tests. Peers behind NATs that randomize ports still fall "
            "back to the relay.",
        ]),
        ("Do I need my own server?", [
            "Yes. There is no public server, so nobody's data passes through anyone else's machine. There are three ways to run it: "
            "a <a href=\"{cloud}\">public server</a>, where the smallest VM is enough and sync works from anywhere; "
            "your <a href=\"{server}\">LAN</a>, on an always-on computer when all your computers share a network; "
            "or <a href=\"{ts}\">LAN + Tailscale</a>, across networks without a public server.",
        ]),
        ("No public IP at home. Can I sync between home and the office?", [
            "Yes, in two ways. With a public server, run CopySync's server there and connect both computers to it; see the "
            "<a href=\"{cloud}\">public server guide</a>. Without one, run the server on an always-on computer at home, install "
            "Tailscale on every computer under the same account, and use the server's Tailscale address; see "
            "<a href=\"{ts}\">LAN + Tailscale</a>.",
        ]),
        ("Windows? iPhone?", [
            "Windows 10 and 11 (x64) are supported and work together with Macs. Files you paste on Windows are real files you "
            "can drop into Explorer or Office.",
            "iPhone and iPad aren't supported: iOS doesn't let apps watch the clipboard in the background, so "
            "copy-and-it's-there can't work.",
        ]),
        ("Why does the system warn me the first time?", [
            "CopySync isn't notarized by Apple yet, and the Windows build isn't code-signed yet; both require paid certificates.",
            "On Mac, open System Settings → Privacy &amp; Security and click Open Anyway near the bottom. On Windows, click More "
            "info and then Run anyway in the \"Windows protected your PC\" prompt. You only do this once.",
        ]),
        ("What about large files?", [
            "By default, anything under 50 MB is pushed to the other computer as you copy it, so you can paste right away. "
            "Larger files sync as a record; pull them from the history when you need them. The limit is adjustable in Settings.",
        ]),
        ("Is it free? Can my company use it?", [
            "Yes. CopySync is open source under the GNU AGPL-3.0, so individuals and companies can use and modify it for free. "
            "If you distribute it, or offer a modified version over a network, you publish your source under the same license.",
            f'To ship it inside a closed-source product, keep modifications private or run it as a hosted service, get a '
            f'<a href="{REPO}/blob/main/LICENSING.md">commercial license</a>.',
        ]),
    ],

    "closing": "Two computers, one clipboard.",
    "source": "View the source on GitHub",
    "footer_license": "CopySync is open source under the GNU AGPL-3.0, with commercial licenses available.",
    "footer_label": "Footer",
    "footer_links": [
        ("All releases", f"{REPO}/releases"),
        ("Changelog", f"{REPO}/blob/main/CHANGELOG.md"),
        ("Design notes (Chinese)", f"{SITE}design/nat-traversal.html"),
        ("Commercial license", f"{REPO}/blob/main/LICENSING.md"),
        ("Report an issue", f"{REPO}/issues"),
    ],

    # ── Tutorial: run the server on your LAN (server/ page) ──
    "guide": {
        "title": "Run a CopySync server on your LAN: Windows, Mac and Linux",
        "description": "When your computers share a network you don't need a cloud server. Run CopySync's signaling and TURN server "
                       "on any always-on Windows, Mac or Linux machine with one script; it starts on boot. Includes app setup and troubleshooting.",
        "h1": "Run the server on your LAN",
        "lede": "If your computers are on the same home or office network, you don't need a cloud server. Pick a computer that stays on, "
                "download and run one script, then enter the address it prints into CopySync on each computer.",
        "toc": [("prepare", "Before you start"), ("install", "Install"), ("check", "Check"), ("clients", "Set up the apps"), ("faq", "Questions")],

        "overview_h2": "The server only introduces devices",
        "overview_intro": "Each computer stays connected to the server, which is how they find each other and exchange connection details. "
                          "Once connected, what you copy goes straight between the two computers. Only when that fails does it pass "
                          "through the server, still encrypted.",
        "diagram": {
            "title": "CopySync on a LAN",
            "desc": "On one LAN, a Mac and a Windows PC connect to the server at 192.168.1.20 for signaling. Copied content travels "
                    "directly between the two computers, and through the server's TURN relay only if a direct connection fails.",
            "lan": "Your LAN: the same Wi‑Fi or wired network",
            "server": "Server",
            "sig_note": "Signaling, TCP 8787",
            "relay_note": "Relay, UDP 3478, if direct fails",
            "client": "Running CopySync",
            "direct": "Direct, end-to-end encrypted",
            "direct_note": "Content stays between the two computers",
            "legend_direct": "Data (direct)",
            "legend_relay": "Data (relay, fallback)",
            "legend_sig": "Signaling",
        },
        "overview_points": [
            "<strong>Signaling, TCP 8787.</strong> Every computer stays connected to it to find the others, pair and exchange connection details.",
            "<strong>Direct.</strong> On a LAN, content almost always goes straight between the two computers, at the speed of your network.",
            "<strong>Relay, UDP 3478.</strong> Only used when devices can't reach each other: guest networks, Wi‑Fi with client isolation, separate subnets.",
        ],

        "prep_h2": "Before you start",
        "prep": [
            ("A computer that stays on", "Windows 10 or 11, macOS 13 or later, or any Linux, including one of the two computers you sync. "
                                        "Sync pauses while it's off or asleep and resumes when it wakes. A NAS, home router or Mac mini "
                                        "is ideal; the server uses about 15 MB of RAM."),
            ("A fixed address", "Reserve an address for it in your router's DHCP settings (often called DHCP reservation or static lease). "
                                "Otherwise, when the router hands out a new address, every computer has to be updated."),
            ("The same network", "Devices just need to be on the same Wi‑Fi or wired network. Guest networks and Wi‑Fi with client "
                                 "(AP) isolation block devices from reaching each other; use another network or turn isolation off."),
            ("No sync away from it", "The server is only reachable on your LAN, so a laptop taken elsewhere catches up when it's back."),
        ],

        "install_h2": "Install the server",
        "install_intro": "Choose the system of the computer that will run the server. The script finds the computer's LAN address, "
                         "sets the server to start automatically, and prints the address to enter in the apps.",
        "os_label": "System of the server computer",
        "manage_caption": "For later",
        "win": {
            "steps": [
                f'Download <a href="{REPO}/releases/latest/download/copysync-server-windows-amd64.zip">copysync-server-windows-amd64.zip</a> '
                "and choose Extract All.",
                "Open the extracted folder and double-click <code>install.cmd</code>."
                "<span class=\"aside\">If you see \"Windows protected your PC\", click More info, then Run anyway; when User Account "
                "Control asks, click Yes. The server isn't code-signed yet, which is why Windows asks once.</span>",
                "Installation continues in a new window. The <code>ws://…/signal</code> address at the end is what you enter in the apps.",
            ],
            "download": "Download the Windows server",
            "notes": [
                "<strong>Runs as a Windows service:</strong> starts with Windows, no sign-in needed, no window. It appears as CopySync Server in Services.",
                "<strong>Firewall already allowed:</strong> the script adds an inbound rule for copysync-server.exe only, so there are no ports to open.",
                "<strong>Wrong address picked</strong> (several network adapters, say from a VM or VPN): type <code>cmd</code> in the "
                "address bar of the extracted folder, press Enter and run <code>install.cmd 192.168.1.20</code> with your address.",
                "Windows 10 and 11, x64.",
            ],
            "double_click": "Double-click",
            "explorer_label": "The extracted folder holds copysync-server.exe, install.cmd, install.ps1, README.md and uninstall.cmd; double-click install.cmd",
            "term_title": "Administrator: Windows PowerShell",
            "term_label": "The window after install.cmd finishes: it stops the old version, installs, registers the service, allows it through the firewall and starts it, then prints the address ws://192.168.1.20:8787/signal",
            "manage": [
                ("Status", "<code>Get-Service CopySyncServer</code> in an administrator PowerShell, or find CopySync Server in Services"),
                ("Logs", "<code>C:\\ProgramData\\CopySync Server\\server.log</code>"),
                ("Restart", "<code>Restart-Service CopySyncServer</code> in an administrator PowerShell"),
                ("Upgrade", "Download the new version, extract it and double-click <code>install.cmd</code> again; your settings are kept"),
                ("Uninstall", "Double-click <code>uninstall.cmd</code>; it removes the service, the firewall rule and the files"),
            ],
        },
        "mac": {
            "steps": [
                "Open Terminal and run these lines:",
                "The <code>ws://…/signal</code> address at the end is what you enter in the apps.",
            ],
            "notes": [
                "<strong>No administrator rights needed.</strong> It runs in the background whenever you're logged in, and macOS restarts it if it exits.",
                "<strong>Keep this Mac awake:</strong> in System Settings → Energy (Battery → Options on a laptop), turn on the option that prevents automatic sleep.",
                "<strong>If the macOS firewall is on,</strong> the script ends with the command to allow copysync-server.",
                "Intel and Apple silicon, macOS 13 or later.",
            ],
            "term_label": "Output of ./install.sh in Terminal: it installs, writes the LaunchAgent and starts the server, then prints the address ws://192.168.1.20:8787/signal",
            "manage": [
                ("Logs", "<code>~/Library/Logs/CopySync/server.log</code>"),
                ("Restart", "<code>launchctl kickstart -k gui/$(id -u)/com.copysync.server</code>"),
                ("Upgrade", "Download the new version, extract it and run <code>./install.sh</code> again; your settings are kept"),
                ("Uninstall", "<code>./install.sh uninstall</code>"),
            ],
        },
        "linux": {
            "steps": [
                "On the machine itself, or over ssh, run these lines:<span class=\"aside\">On ARM, such as a Raspberry Pi with a 64-bit OS, replace <code>amd64</code> with <code>arm64</code>.</span>",
                "The <code>ws://…/signal</code> address at the end is what you enter in the apps.",
            ],
            "notes": [
                "<strong>Registers a systemd service</strong> that starts on boot. Needs root.",
                "<strong>If ufw or firewalld is active,</strong> the script ends with the commands to open the ports.",
                "<strong>Docker works too:</strong> copy <code>env.example</code> to <code>.env</code>, set <code>TURN_IP</code> to the "
                "machine's LAN address and run <code>docker compose up -d</code>. It needs host networking, so Linux only.",
            ],
            "term_label": "Output of sudo ./install.sh: it installs, writes the systemd unit and config and starts the server, then prints the address ws://192.168.1.20:8787/signal",
            "manage": [
                ("Status", "<code>systemctl status copysync-server</code>"),
                ("Logs", "<code>journalctl -u copysync-server -f</code>"),
                ("Restart", "<code>sudo systemctl restart copysync-server</code>"),
                ("Config", "<code>/etc/copysync/server.env</code>; restart the service after editing"),
                ("Upgrade", "Download the new version, extract it and run <code>sudo ./install.sh</code> again; your config is kept"),
                ("Uninstall", "<code>sudo ./install.sh uninstall</code>"),
            ],
        },

        "check_h2": "Check that other computers can reach it",
        "check_intro": "On another computer, open <code>http://192.168.1.20:8787/healthz</code> in a browser (with your address). "
                       "If you see this line, the network path works.",
        "check_label": "A browser showing http://192.168.1.20:8787/healthz, with the server's status ok and its version",
        "check_after": "If it doesn't load, see the <a href=\"#faq\">questions</a> below.",

        "clients_h2": "Enter the address on each computer",
        "clients_steps": [
            ("Open Settings", "Open CopySync on each computer, including the one running the server, and go to Settings."),
            ("Enter the address", "Under Signaling server, enter <code>ws://192.168.1.20:8787/signal</code> with your address and press Return or Enter."),
            ("Check the status", "When Connection shows a green Connected, you're set. Pair your devices under Devices and copying starts syncing."),
        ],
        "clients_alt": "CopySync Settings: the signaling server is ws://192.168.1.20:8787/signal and the connection status shows Connected",
        "clients_after": 'Haven\'t installed the app yet? See the <a href="{start}">quick start</a>.',

        "faq_h2": "Questions",
        "faq": [
            ("The status stays at Not connected", [
                "Check the address first: it starts with <code>ws://</code>, uses port 8787, ends in <code>/signal</code>, and you pressed Return after typing it.",
                "Then open <code>http://server-address:8787/healthz</code> in a browser on that computer. If it doesn't load, the network "
                "path is broken: the server isn't running, its computer is asleep, a firewall blocks it, or its address changed.",
                "The Windows script already allows the server through the firewall. On a Mac with the firewall on, run the command the "
                "script printed at the end; the same goes for ufw and firewalld on Linux.",
            ]),
            ("Devices show Relay instead of Direct", [
                "On a LAN connections are almost always direct. Relay means UDP can't get between the two computers, usually because of "
                "a guest network, Wi‑Fi client isolation, or separate subnets with a firewall in between. Relayed content is still "
                "encrypted; it's just limited by the server computer's network.",
            ]),
            ("The server's address changed", [
                "Run the script again with the new address: <code>install.cmd new-address</code> on Windows, <code>./install.sh new-address</code> "
                "on a Mac. On Linux, edit <code>-turn-ip</code> in <code>/etc/copysync/server.env</code> and restart the service. Then update "
                "the address on each computer.",
                "Reserving the address in your router avoids this.",
            ]),
            ("Does it sync away from the LAN?", [
                "No, the server is only reachable on your LAN. To sync from other networks, switch to one of the other two ways: "
                "a <a href=\"{cloud}\">public server</a>, or <a href=\"{ts}\">LAN + Tailscale</a>.",
            ]),
            ("The server computer also runs CopySync. Which address?", [
                "The same LAN address works, and so does <code>ws://127.0.0.1:8787/signal</code>. The server and the app don't interfere.",
            ]),
        ],

        "closing": "One computer that stays on is enough.",
        "source": "Full server guide, including public servers and Docker",
    },

    # ── Run it on a public server (cloud/ page) ──
    "cloud": {
        "title": "Run the CopySync server on a public server: cloud VM guide",
        "description": "One script installs CopySync's signaling and TURN server on a Linux VM with a public IP. Every computer "
                       "connects to it directly, so the clipboard syncs at home, at the office and anywhere else. Covers the "
                       "security group, checking, app setup and domains in mainland China.",
        "h1": "Run it on a public server",
        "lede": "A Linux VM with a public IP; the smallest one is enough. Every computer connects to it directly and syncs at "
                "home, at the office or anywhere else, with nothing extra to install.",
        "toc": [("prepare", "Before you start"), ("ports", "Open ports"), ("install", "Install"), ("check", "Check"),
                ("clients", "Set up the apps"), ("faq", "Questions")],
        "step": "Step {n}",

        "overview_h2": "Every computer connects to it",
        "overview_intro": "The server is on the public internet, so every computer can reach it from whatever network it's on. "
                          "What you copy goes directly between the two computers when possible; when that fails, it's relayed "
                          "through the server, still encrypted.",
        "diagram": {
            "title": "CopySync with a public server",
            "desc": "The public server is at 203.0.113.10. A Mac on an office network and a Windows PC on a home network both "
                    "connect to it for signaling; copied content goes directly between them, and through the server's relay if that fails.",
            "server": "Public server",
            "server_sub": "Cloud VM with a public IP",
            "office": "Office network",
            "home": "Home network",
            "client": "Running CopySync",
            "sig_note": "Signaling, TCP 8787",
            "relay_note": "Relay over UDP, if direct fails",
            "direct": "Direct",
        },
        "overview_points": [
            "<strong>Signaling, TCP 8787.</strong> Every computer stays connected to it to find the others, pair and exchange connection details.",
            "<strong>Direct.</strong> Computers on different networks connect directly through NAT hole punching, which works on most home and office networks.",
            "<strong>Relay, UDP 3478 and relay ports.</strong> When a direct path fails, the server forwards the traffic, still encrypted and limited by its bandwidth.",
        ],

        "prep": [
            ("A cloud VM", "Linux on x86_64 or ARM64; the smallest size is enough, as the server uses about 15 MB of RAM. You need root."),
            ("A public IP", "The VM's own public IP works; a static (elastic) one is best. If it changes, every computer has to be updated."),
            ("A way in", "ssh, or the web terminal in your provider's console."),
            ("VMs in mainland China", "Without an ICP filing for the domain, enter the IP in the apps, not the domain; see the <a href=\"#faq\">questions</a>."),
        ],

        "ports_h2": "Open the ports in the security group",
        "ports_intro": "In your cloud provider's console, add these three inbound rules to the VM's security group.",
        "rules_url": "Cloud console › Security group › Inbound rules",
        "rules_cols": ["Protocol", "Port range", "Source", "Purpose"],
        "rules": [
            ("TCP", "8787", "0.0.0.0/0", "Signaling"),
            ("UDP", "3478", "0.0.0.0/0", "STUN and TURN"),
            ("UDP", "32768-60999", "0.0.0.0/0", "Relay ports"),
        ],
        "rules_label": "Inbound rules in a cloud security group (illustration): TCP 8787 for signaling, UDP 3478 for STUN and TURN, UDP 32768-60999 for relay ports, all from 0.0.0.0/0",
        "ports_notes": [
            "<strong>Different names:</strong> some providers call it a firewall or a network ACL; it's the same setting.",
            "<strong>Relay ports</strong> are assigned by the OS; the range depends on the server. Run <code>cat /proc/sys/net/ipv4/ip_local_port_range</code> "
            "to see it; the install script prints it too.",
            "<strong>System firewall:</strong> if ufw or firewalld is active, open the same ports there; the install script prints the commands.",
        ],

        "install_h2": "Run the install script",
        "install_intro": "Give the script the public IP. It installs the server, registers a systemd service that starts on boot, "
                         "and prints the address to enter in the apps.",
        "install_steps": [
            "Log in with ssh, for example <code>ssh root@203.0.113.10</code> with your public IP.",
            "Run these lines, with your public IP on the last one:<span class=\"aside\">On ARM, replace <code>amd64</code> with <code>arm64</code>.</span>",
            "The <code>ws://203.0.113.10:8787/signal</code> address at the end is what you enter in the apps.",
        ],
        "install_notes": [
            "<strong>Registers a systemd service</strong> that starts on boot; the config is <code>/etc/copysync/server.env</code>.",
            "<strong>Upgrade:</strong> download the new version, extract it and run <code>sudo ./install.sh</code> again; your config is kept.",
            "<strong>Docker works too:</strong> see the <a href=\"{deploy}\">server deployment guide</a>.",
        ],
        "term_label": "Output of sudo ./install.sh 203.0.113.10: it installs, writes the systemd unit and config and starts the server, then prints the address ws://203.0.113.10:8787/signal",

        "check_h2": "Check that it's reachable",
        "check_intro": "Open <code>http://203.0.113.10:8787/healthz</code> (with your public IP) in a browser on any computer. If you see this line, it works.",
        "check_label": "A browser showing http://203.0.113.10:8787/healthz, with the server's status ok and its version",

        "clients_steps": [
            ("Open Settings", "Open CopySync on each computer and go to Settings."),
            ("Enter the public address", "Under Signaling server, enter <code>ws://203.0.113.10:8787/signal</code> with your public IP and press Return or Enter."),
            ("Check the status", "When Connection shows a green Connected, you're set. If you haven't paired yet, pair under Devices."),
        ],
        "clients_alt": "CopySync Settings: the signaling server is ws://203.0.113.10:8787/signal and the connection status shows Connected",

        "faq": [
            ("It works by IP but not by domain", [
                "Cloud providers in mainland China intercept plain HTTP requests for domains without an ICP filing and answer with a "
                "redirect to a block page; in our tests that included non-standard ports such as 8787. The WebSocket handshake is an "
                "HTTP request too, so it is intercepted and the app shows Not connected.",
                "Enter the server's IP instead, as in <code>ws://server-ip:8787/signal</code>, or file the domain.",
            ]),
            ("How big a server do I need?", [
                "The smallest one. The server uses about 15 MB of RAM and normally forwards only a few signaling messages. Content "
                "goes directly between computers when possible; only relayed transfers are limited by the server's bandwidth.",
            ]),
            ("Devices show Relay instead of Direct", [
                "Relay right after connecting, then Direct a little later, is normal: the relay gets you going and the direct path "
                "takes over while idle. Staying on Relay means the two computers can't reach each other directly, usually because "
                "of a NAT that randomizes ports or a strict firewall. It's still encrypted, just limited by the server's bandwidth.",
            ]),
            ("Can I use HTTPS (wss://)?", [
                "Yes, behind an HTTPS reverse proxy that supports WebSocket upgrades; enter an address starting with <code>wss://</code>. "
                "Domains in mainland China need an ICP filing. STUN and TURN use UDP and bypass the proxy, so keep those ports open.",
            ]),
            ("My home connection has a public IP. Can I use that?", [
                "Yes, but it's more work than a VM: forward TCP 8787, UDP 3478 and a large range of UDP relay ports to that computer, "
                "and a home IP may change. If all your computers are at home, the <a href=\"{server}\">LAN</a> setup is simpler; to sync "
                "across networks without port forwarding, see <a href=\"{ts}\">LAN + Tailscale</a>.",
            ]),
        ],

        "closing": "One cloud VM,<br>sync from anywhere.",
        "source": "Full server guide, including Docker and every option",
    },

    # ── Quick start (start/ page) ──
    "start": {
        "title": "Quick start: set up CopySync in three steps",
        "description": "New to CopySync? Install the app on each computer, pick a way to run the server (on your LAN, with "
                       "Tailscale, or in the cloud), enter its address and pair, then copy and paste between Mac and Windows. "
                       "No public IP needed.",
        "h1": "Quick start",
        "lede": "Install the app, set up a server, and pair your computers. The first time takes about ten minutes; "
                "after that, copying just syncs.",
        "toc": [("app", "Install the app"), ("server", "Set up a server"), ("connect", "Enter the address"), ("pair", "Pair")],
        "step": "Step {n}",

        "app_h2": "Install the app on each computer",
        "app_intro": "One on every computer you want to sync. Macs and Windows PCs mix freely.",
        "mac": [
            "Download CopySync.dmg, open it and drag CopySync into Applications.",
            "Open it from Applications.<span class=\"aside\">The first time, macOS says it can't verify the developer: go to "
            "System Settings → Privacy &amp; Security and click Open Anyway near the bottom.</span>",
            "Click Turn on background sync, and allow clipboard access when macOS asks.",
        ],
        "win": [
            "Run CopySync-Setup.exe; no administrator rights needed.<span class=\"aside\">If you see \"Windows protected your PC\", "
            "click More info, then Run anyway.</span>",
            "Open CopySync and click Turn on background sync. After you close the window it keeps running in the notification area.",
        ],
        "app_alt": "CopySync the first time it opens, asking you to turn on background sync",

        "server_h2": "Set up a server",
        "server_intro": "Your computers find each other through it; nothing you copy is stored there. There's no public server, "
                        "and there are three ways to run your own. Pick the one that fits.",

        "connect_h2": "Enter the server address on each computer",
        "connect_steps": [
            ("Open Settings", "Open CopySync on each computer, including the one running the server, and go to Settings."),
            ("Enter the address", "Under Signaling server, enter the address from the previous step and press Return or Enter."),
            ("Check the status", "When Connection shows a green Connected, that computer is ready."),
        ],
        "connect_alt": "CopySync Settings: the signaling server is ws://192.168.1.20:8787/signal and the connection status shows Connected",
        "connect_after": "Stuck at Not connected? See the <a href=\"{server}#faq\">questions</a>.",

        "pair_h2": "Pair your computers",
        "pair_intro": "Pair each two computers once, and you're done.",
        "pair_steps": [
            "On one computer, open Devices and click Add device. It shows a 6-character code, valid for 5 minutes.",
            "On the other, click Add device too, switch to Enter a code and type it in.",
            "Both screens show the same two fingerprint lines. Check them character by character, then confirm."
            "<span class=\"aside\">This stops anyone in the middle from posing as your device. Don't skip it.</span>",
        ],
        "pair_alt": "Add device: the pairing code K7MP2X, to be entered on the other device",

        "done_h2": "Now copy something.",
        "done_text": "Copy text, an image or a file on one computer and press ⌘V or Ctrl+V on the other. Everything copied on "
                     "either side is in the history.",
        "help": "Need help? <a href=\"{server}#faq\">LAN questions</a> · <a href=\"{ts}#faq\">Tailscale questions</a> · "
                "<a href=\"{issues}\">Report an issue</a>",
    },

    # ── Use it across networks with Tailscale (tailscale/ page) ──
    "ts": {
        "title": "Use CopySync across networks with Tailscale: sync home and office without a public IP",
        "description": "Run the server on an always-on computer at home and install Tailscale on each computer: they act as if "
                       "they shared a LAN, so the clipboard syncs from the office or anywhere else. No public IP, no router "
                       "settings. Step-by-step for Windows, Mac and Linux.",
        "h1": "Across networks with Tailscale",
        "lede": "Sync between home and the office without a public IP. The server stays at home; with Tailscale on each "
                "computer, they act as if they shared a LAN. No router settings, no cloud server.",
        "toc": [("join", "Install Tailscale"), ("address", "Find the address"), ("install", "Install the server"),
                ("check", "Check"), ("clients", "Set up the apps"), ("faq", "Questions")],
        "step": "Step {n}",

        "overview_h2": "As if on one LAN",
        "overview_intro": "Tailscale connects your devices into a private network: every device signed in to the same account "
                          "gets a fixed address starting with 100, reachable from the others wherever they are. Install the "
                          "CopySync server on that address and enter it in the apps; everything else works as on a LAN.",
        "diagram": {
            "title": "CopySync over Tailscale",
            "desc": "The server and a Windows PC are at home, a Mac is at the office. All three run Tailscale and each has an "
                    "address starting with 100 on the same private network. Both computers connect to the server at "
                    "100.101.102.103 for signaling; copied content goes directly between them, and through the server's relay "
                    "only if that fails.",
            "net": "Tailscale network: every device has a 100.x address",
            "server": "Server",
            "server_place": "At home, always on",
            "mac_place": "At the office",
            "win_place": "At home",
            "sig_note": "Signaling, TCP 8787",
            "relay_note": "Relay, UDP 3478, if direct fails",
            "direct": "Direct, encrypted twice",
            "direct_note": "Content stays between the two computers",
        },
        "overview_points": [
            "<strong>The server stays at home</strong> on any always-on Windows, Mac or Linux computer, just like on a LAN.",
            "<strong>Every computer runs Tailscale,</strong> signed in to the same account, including the server computer.",
            "<strong>Encrypted twice.</strong> Inside Tailscale's WireGuard tunnel runs CopySync's own end-to-end encrypted "
            "connection. Neither Tailscale nor your server can see the content.",
        ],

        "prep": [
            ("A computer that stays on", "It runs the server, with the same requirements as <a href=\"{server}#prepare\">on a LAN</a>; "
                                        "sync pauses while it's off or asleep. Its Tailscale address doesn't change, so there's "
                                        "nothing to reserve in the router."),
            ("A Tailscale account", "Free for personal use; sign in with Google, Microsoft, GitHub or Apple. Use the same account on every device."),
            ("Tailscale on every computer", "There are clients for Windows, macOS and Linux. Every computer that runs CopySync needs "
                                           "it, including one that never leaves home."),
            ("ZeroTier works too", "So does any tool that gives each device a fixed virtual address: install the server on that "
                                   "address and enter it in the apps. This guide uses Tailscale."),
        ],

        "join_h2": "Install Tailscale on each computer",
        "join_intro": "On the server computer and every computer you sync, install it and sign in to the same account.",
        "join_steps": [
            "Download the client for your system from the Tailscale website and install it."
            "<span class=\"aside\">On Linux, run <code>curl -fsSL https://tailscale.com/install.sh | sh</code>, then "
            "<code>sudo tailscale up</code>, and open the link it prints to sign in.</span>",
            "Open Tailscale and sign in, using the same account on every computer.",
            "The Tailscale admin console then lists all your devices, each with an address starting with 100.",
        ],
        "join_download": "Download Tailscale",
        "machines_cols": ["Machine", "Tailscale address", "OS"],
        "machines_server": "Server",
        "machines_label": "The machine list in the Tailscale admin console (illustration): mac-mini is the server at 100.101.102.103, "
                          "macbook-pro is 100.88.20.5 and home-pc is 100.92.41.17",
        "server_notes_h3": "One more thing on the server computer",
        "server_notes": [
            "<strong>Windows:</strong> right-click the Tailscale icon in the notification area and turn on Run unattended under "
            "Preferences. Otherwise Tailscale disconnects when you sign out, and computers elsewhere lose the server.",
            "<strong>Mac:</strong> stay logged in, keep it awake, and make sure Tailscale opens at login. The CopySync server "
            "also only runs while you're logged in.",
            "<strong>Turn off key expiry:</strong> by default Tailscale asks devices to sign in again every 180 days. On the "
            "Machines page of the admin console, open the menu at the end of the server's row and choose Disable key expiry.",
        ],

        "address_h2": "Find the server's Tailscale address",
        "address_intro": "On the server computer, click the Tailscale icon in the menu bar (the notification area on Windows); "
                         "the menu shows this device's address. Or run <code>tailscale ip -4</code> in a terminal.",
        "address_label": "tailscale ip -4 in a terminal prints this device's Tailscale address, 100.101.102.103",
        "address_after": "The steps below use <code>100.101.102.103</code>; use your own address.",

        "install_h2": "Install the server on that address",
        "install_intro": "The install scripts use the LAN address by default. Pass the Tailscale address instead so computers "
                         "elsewhere can use the server's relay. Already installed it with the <a href=\"{server}\">LAN guide</a>? "
                         "Run the script again with the address; your settings are kept.",
        "win": {
            "steps": [
                "Download copysync-server-windows-amd64.zip and choose Extract All. If it's already installed, use the folder you extracted before.",
                "Open the extracted folder, type <code>cmd</code> in the address bar at the top and press Enter to open a Command Prompt.",
                "Run <code>install.cmd 100.101.102.103</code> with your address, and click Yes when User Account Control asks.",
                "Installation continues in a new window. The address at the end should be <code>ws://100.101.102.103:8787/signal</code>.",
            ],
            "notes": [
                "<strong>Firewall already allowed:</strong> the script's inbound rule applies to every network, Tailscale's adapter included.",
                "If you see \"Windows protected your PC\", click More info, then Run anyway.",
            ],
            "cmd_title": "Command Prompt",
            "cmd_label": "A Command Prompt opened in the extracted folder runs install.cmd 100.101.102.103",
            "term_label": "The window after install.cmd finishes, ending with the address ws://100.101.102.103:8787/signal",
        },
        "mac": {
            "steps": [
                "Open Terminal and run these lines, with your address on the last one:",
                "The address at the end should be <code>ws://100.101.102.103:8787/signal</code>.",
            ],
            "notes": [
                "Already installed? Run <code>./install.sh 100.101.102.103</code> in the same folder; your settings are kept.",
                "<strong>If the macOS firewall is on,</strong> the script ends with the command to allow copysync-server.",
            ],
            "term_label": "Output of ./install.sh 100.101.102.103 in Terminal, ending with the address ws://100.101.102.103:8787/signal",
        },
        "linux": {
            "steps": [
                "Run these lines, with your address on the last one:<span class=\"aside\">On ARM, replace <code>amd64</code> with <code>arm64</code>.</span>",
                "The address at the end should be <code>ws://100.101.102.103:8787/signal</code>.",
            ],
            "switch": "<strong>Already installed?</strong> The Linux script leaves an existing config alone. Replace the old "
                      "address in it with the Tailscale address (put your old address in place of <code>192.168.1.20</code>) "
                      "and restart the service:",
            "term_label": "Output of sudo ./install.sh 100.101.102.103, ending with the address ws://100.101.102.103:8787/signal",
        },

        "check_h2": "Check from outside",
        "check_intro": "On a computer outside your home network, or a laptop on your phone's hotspot, make sure Tailscale is "
                       "connected and open <code>http://100.101.102.103:8787/healthz</code> in a browser. If you see this line, "
                       "the server is reachable from outside.",
        "check_label": "A browser showing http://100.101.102.103:8787/healthz, with the server's status ok and its version",
        "check_after": "If it doesn't load, see the <a href=\"#faq\">questions</a> below.",

        "clients_h2": "Enter the address on each computer",
        "clients_steps": [
            ("Open Settings", "Open CopySync on each computer, including the one running the server, and go to Settings."),
            ("Enter the Tailscale address", "Under Signaling server, enter <code>ws://100.101.102.103:8787/signal</code> with your address and press Return or Enter."),
            ("Check the status", "When Connection shows a green Connected, you're set. If you haven't paired yet, pair under Devices."),
        ],
        "clients_alt": "CopySync Settings: the signaling server is ws://100.101.102.103:8787/signal and the connection status shows Connected",
        "clients_after": "Use the Tailscale address on your home computers too, not the LAN address, so nothing changes when a "
                         "computer leaves the house. When they're all at home, the data still stays on your home network.",

        "faq": [
            ("Devices show Direct, but it's slow", [
                "Direct in CopySync means the two computers connect directly on the Tailscale network. When Tailscale can't "
                "connect them peer to peer, it relays through its DERP servers, which can be far away.",
                "Run <code>tailscale ping other-address</code> in a terminal: via DERP means Tailscale is relaying; via an IP "
                "and port means it's peer to peer. Office and campus networks relay more often; you can run your own DERP "
                "server or try another mesh VPN.",
            ]),
            ("The status stays at Not connected", [
                "First make sure Tailscale is connected on that computer, then open <code>http://server-tailscale-address:8787/healthz</code> in a browser.",
                "If it doesn't load, check whether the server computer is asleep, whether Tailscale is online there (Run unattended "
                "on Windows, staying logged in on a Mac), whether the server still uses its LAN address, and whether a firewall on "
                "a Mac or Linux server blocks it.",
            ]),
            ("Does every computer need Tailscale?", [
                "Yes. A computer without it can't reach the 100.x addresses, and computers elsewhere can't reach it. That includes "
                "one that never leaves home.",
            ]),
            ("I already followed the LAN guide. Start over?", [
                "No. Run the install script again with the Tailscale address (on Linux, edit the config), then change the server "
                "address on each computer to the Tailscale one. Paired devices stay paired.",
            ]),
            ("Can Tailscale see what I copy?", [
                "No. Tailscale connections are end-to-end encrypted with WireGuard, and inside runs CopySync's own DTLS-encrypted "
                "connection; only your two computers can open either layer. Tailscale's servers only exchange device keys and "
                "addresses, and its relays forward encrypted packets.",
            ]),
            ("Will the server's Tailscale address change?", [
                "Usually not. The address belongs to the device and survives restarts and network changes; it only changes if you "
                "remove the device from Tailscale and add it again.",
            ]),
            ("I have a public IP. Do I need Tailscale?", [
                "No. Run the server on the machine with the public IP and enter that address on every computer; see the "
                "<a href=\"{cloud}\">public server guide</a>.",
            ]),
        ],

        "closing": "The server stays home.<br>Your clipboard comes with you.",
        "source": "Back to the LAN guide",
    },

    "ld_description": "Open-source app that syncs the clipboard and files between Mac and Windows computers. Copy text, "
                      "images, files or folders on one computer and paste on the other. Devices connect directly over "
                      "WebRTC with end-to-end DTLS encryption, falling back to a TURN relay on a server the user hosts.",
    "ld_os": "macOS 13 or later (Intel and Apple silicon); Windows 10, Windows 11 (x64)",
    "ld_currency": "USD",
    "ld_features": [
        "Syncs plain text, rich text, images, files and folders",
        "Works between Mac and Windows, across accounts and networks",
        "Direct WebRTC connection with automatic TURN relay fallback",
        "End-to-end DTLS encryption with human-verified fingerprints at pairing",
        "Items under 50 MB pushed on copy, larger files pulled on demand",
        "3-day history with image preview",
        "Interface in English, Simplified Chinese and Japanese",
    ],
}

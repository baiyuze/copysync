"""Site copy: English. Same keys as zh.py; see tools/site/build.py for the structure."""

REPO = "https://github.com/baiyuze/pasteferry"
SITE = "https://pasteferry.cc/"

C = {
    "brand": "PasteFerry",
    "brand_alt": ["隔空贴", "CopySync"],
    "title": "PasteFerry: sync the clipboard and files between Mac, Windows and Linux, across accounts and networks",
    "description": 'PasteFerry syncs clipboard text, images, files and folders across Mac, Windows and Linux. Enable optional official access to connect across networks without hosting a server: direct first, encrypted relay when needed. Free and open source.',
    "og_title": "PasteFerry: copy on this computer, paste on that one",
    "og_description": 'PasteFerry syncs clipboard text, images, files and folders across Mac, Windows and Linux. Enable optional official access to connect across networks without hosting a server: direct first, encrypted relay when needed. Free and open source.',
    "og_image": "og-en.png",
    "skip": "Skip to content",
    "home_label": "PasteFerry home",
    "nav_label": "Page",
    "nav_start": "Quick start",
    "nav_how": "How it works",
    "nav_faq": "FAQ",
    "nav_server": "Ways to connect",
    "nav_download": "Download",
    "lang_label": "Language",

    "h1": "Copy on this computer.<br>Paste on that one.",
    "lede": 'Mac, Windows and Linux. Text, images, files, whole folders. No shared account required. Now connect across networks by enabling official access, with no server to host. Direct connections first, end-to-end encryption throughout.',
    "dl_button": "Download PasteFerry",
    "dl_mac": "Download for Mac",
    "dl_win": "Download for Windows",
    "dl_linux": "Download for Linux",
    "dl_note": "Version 2.1.0, free and open source. macOS 13 or later; Windows 10 and 11; Linux (x64 and ARM64).",
    "start_link": "New here? Set it up in three steps →",
    "stage_label": "Demo: Command C on a MacBook turns a file into a sheet that flies over a direct connection to a Windows "
                   "laptop, where Ctrl V pastes it",
    "stage_file": "Q3 review - final.pptx",
    "stage_win": "Windows laptop",
    "stage_arc": "Direct, end-to-end encrypted",
    "replay": "Play again",

    # Three ways to connect, on equal footing; the reader picks one. Home cards use name, fit, summary, go;
    # the quick start tabs use name, facts, steps, address, go; the install section uses name, fit.
    # {lan}, {cloud}, {server} and {ts} in steps and address become links to the matching guide.
    "paths_h2": "How your computers connect",
    "paths_intro": 'Same network? Install and pair. Different networks? Enable official access without hosting a server. Custom servers and Tailscale remain available when you prefer to manage the connection yourself.',
    "paths_after": 'New here? Follow the <a href="{start}#server">quick start</a>. Official access is optional; its current relay default is 100 MB/day, and direct traffic does not consume it.',
    "plans_label": "Ways to connect",
    "plans": {
        "lan": {
            "name": "Same network",
            "fit": "Nothing to set up",
            "summary": "At home or in the office, behind the same router. Install the app on both computers and enter a pairing code. "
                       "No server, and no computer that has to stay on.",
            "go": "Same network guide →",
            "facts": [
                ("Best for", "Computers on the same Wi‑Fi or wired network, such as at home or in the office."),
                ("You need", "Nothing. Just the app on each computer."),
                ("Keep in mind", "No sync away from that network; it resumes when you're back. If wired and Wi‑Fi are on "
                                 "different subnets, as in many offices, you also enter an address when pairing."),
            ],
            "steps": [
                "Install the app on each computer and click Turn on background sync.",
                "Allow PasteFerry if Windows Firewall asks, and allow it to find devices on your local network if macOS asks.",
                "Show a code on one, enter it on the other and compare fingerprints. On a different subnet, also enter the address "
                "shown under the code; see the <a href=\"{lan}#subnet\">same network guide</a>.",
            ],
            "address": "No server address to enter.",
        },
        "cloud": {
            "name": 'Self-hosted server',
            "fit": 'Control your own capacity',
            "summary": "A Linux VM with a public IP; the smallest one is enough. Every computer connects to it directly, at home, "
                       "at the office or anywhere else, with nothing extra to install.",
            "go": "Public server guide →",
            "facts": [
                ("Best for", 'Control your own capacity'),
                ("You need", "A Linux VM with a public IP, as root; the smallest one is enough."),
                ("Keep in mind", "The VM costs money. When a direct connection fails, the relay is limited by the server's "
                                 "bandwidth. In mainland China, use the IP unless the domain has an ICP filing."),
            ],
            "steps": [
                "In your cloud provider's security group, open TCP 8787, UDP 3478 and UDP 32768–60999.",
                "Log in with ssh, download the install script and run it with the public IP: <code>sudo ./install.sh PUBLIC_IP</code>.",
                "In Settings on each computer, enter <code>ws://PUBLIC_IP:8787/signal</code>, press Return, then pair.",
            ],
            "address": "Open <code>http://PUBLIC_IP:8787/healthz</code> in a browser; ok means it's running.",
        },
        "ts": {
            "name": "Tailscale",
            "fit": "Across networks, no server",
            "summary": "Install Tailscale on each computer under one account and enter the other computer's Tailscale address when "
                       "pairing. Syncs at home and at the office, with no server and no computer that has to stay on.",
            "go": "Tailscale guide →",
            "facts": [
                ("Best for", "Computers on different networks, without a public server."),
                ("You need", "Tailscale on every computer, signed in to one account."),
                ("Keep in mind", "Every computer that runs PasteFerry needs Tailscale, connected. When Tailscale can't connect "
                                 "peer to peer it relays through DERP servers, which can be slow."),
            ],
            "steps": [
                "Install Tailscale on each computer and sign in to the same account.",
                "On the computer you'll pair with, run <code>tailscale ip -4</code> and note the address starting with 100.",
                "When pairing, on the computer entering the code, click On a different subnet? and enter that address.",
            ],
            "address": "No server. A Tailscale address doesn't change, so you enter it once, when pairing.",
        },
    },

    "history_h2": "Everything you copied, on every computer",
    "history_intro": "History keeps the last 3 days, ready to go back on the clipboard; click an image to preview it. "
                     "Files under 50 MB have already arrived by the time you switch computers, and larger ones are a click away. "
                     "To find something you copied earlier, search by its text or file name.",
    "history_alt": "PasteFerry history: files, images and text from a MacBook Air and a Mac mini; a large file shows a button "
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
    "compare_caption": "PasteFerry compared with other options",
    "compare_cols": ["PasteFerry", "Universal Clipboard", "AirDrop", "Chat apps"],
    "compare_rows": [
        ("Different Apple IDs", [("yes", "Yes"), ("no", "No"), ("partial", "Needs accepting"), ("yes", "Yes")]),
        ("Between Mac and Windows", [("yes", "Yes"), ("no", "No"), ("no", "No"), ("yes", "Yes")]),
        ("Different networks (home and office)", [("yes", "Yes"), ("no", "No"), ("no", "No"), ("yes", "Yes")]),
        ("Paste right away, nothing to accept", [("yes", "Yes"), ("yes", "Yes"), ("no", "No"), ("no", "No")]),
        ("Files and whole folders", [("yes", "Yes"), ("partial", "Unreliable"), ("yes", "Yes"), ("partial", "Zip first")]),
        ("Find something you copied earlier", [("yes", "Yes"), ("no", "No"), ("no", "No"), ("partial", "Scroll back")]),
        ('Only your devices can decrypt content', [("yes", "Yes"), ("yes", "Yes"), ("yes", "Yes"), ("no", "No")]),
    ],

    "how_h2": "How it works",
    "how_intro": "Your two computers open an encrypted connection to each other. On the same network they find each other "
                 "on their own; across networks your server introduces them, and relays traffic when a direct path can't be "
                 "made. The diagram shows the cross-network case.",
    "diagram": {
        "title": "How PasteFerry connects",
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
                    "per destination. A single probe only reveals one of them, so connections end up on the relay. PasteFerry "
                    "probes every uplink and tells the peer about all of them.",
    "uplink_changed": "What it does",
    "uplink_points": [
        ("One shared port.", " All connections send and receive on one local UDP port, so the addresses learned from different servers line up."),
        ("Every uplink probed.", " Seven IP-diverse servers are queried from that port at once, yielding one address per uplink. Addresses poisoned by DNS are dropped first."),
        ("All of them sent to the peer.", " The peer tries each address; the one matching the real uplink gets through."),
        ("Uplinks remembered.", " If a round misses an uplink that keeps ports unchanged, its address is filled in as uplink IP plus local port."),
        ("A rule for who sends first.", " Some routers lock a port when the other side's packet arrives first. One side announces its addresses first; the other sends to them before announcing its own. Retries alternate the order."),
        ("Relay first, direct later.", " If the first handshake lands on the relay, PasteFerry punches again while idle and switches to direct when that works."),
    ],
    "lab_h3": "Checked against real network topologies",
    "lab_intro": "The NAT lab builds these topologies with Linux network namespaces and iptables and runs PasteFerry's actual "
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
        ("No server, same subnet, found over mDNS", "Direct"),
        ("No server, two subnets, one side remembers the address", "Direct"),
        ("No server, both block inbound, both remember", "Direct"),
        ("No server, both block inbound, only one remembers", "No connection (known limit)"),
        ("No server, subnets isolated from each other", "No connection"),
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
        ("Same on the local network.", " What a device broadcasts on the LAN is an anonymous token that changes daily, not its name; the name is added only while Search nearby is open. Pairing still means comparing fingerprints, and every message is still signed."),
    ],
    "verify_alt": "The fingerprint check during pairing shows R8NF-2WTC-QL5J and asks you to confirm it matches the other computer exactly",

    "fallback_h2": "When things go wrong",
    "fallback_intro": "Every network is different. Each of these cases was considered, and each has a test.",
    "fallbacks": [
        ("Can't find the other computer on the network", "Different subnets, guest networks and Wi‑Fi with client isolation block local discovery. Enter the address shown under the other computer's code when pairing; it's remembered afterwards, and if both computers change addresses you enter it once more in the device list."),
        ("The system blocks the local network", "From macOS 15, the background service needs your permission to reach devices on the local network; Windows asks about the firewall on first run. If either isn't allowed, the app says so and links straight to the fix."),
        ("No direct path between the computers", "Common behind symmetric NAT and corporate firewalls. Traffic switches to the TURN relay on your server, still encrypted, and the device list shows \"Relay\"."),
        ("Public STUN servers unreachable", "Public STUN is unreliable in some regions, including mainland China. Your server runs its own STUN, and clients prefer it."),
        ("Connected through the relay first", "The relay handshake is fast and sometimes beats the direct path. Once connected, the side that started the connection punches again while idle and switches to direct when that works, never in the middle of a transfer."),
        ("The other computer sleeps, drops off or restarts", "A connection that dies silently is noticed, closed and reconnected. Anything you copy meanwhile is kept and added to the other computer's history once it's back, without overwriting what's on its clipboard now."),
        ("Server connection drops", "Reconnects automatically: the delay starts at 1 second and doubles up to 30, with jitter so devices don't all reconnect at once. Devices on the same network keep syncing."),
        ("File too large", "Above the limit only a record is synced, pulled when you need it. If the source file is gone by then, you get a clear message instead of a vague failure."),
        ("Pasting into plain-text fields", "Rich text copied from the web is written as both HTML and plain text, so no markup leaks into what you paste."),
        ("A file name Windows doesn't allow", "Colons, question marks and a few other characters that are fine on a Mac arrive as the matching full-width characters, so the name still reads the same. Two names that differ only in case never overwrite each other."),
        ("Copying a password", "On Windows, password managers mark what they copy as \"don't record\". PasteFerry skips it: not synced, not kept in the history."),
        ("Clipboard access not granted yet", "Content isn't read, which would block on a prompt a background process can't show. The app asks you to grant access and links straight to the setting."),
        ("Background service crashes", "launchd on Mac, and PasteFerry's own supervisor on Windows, restart it, throttled to once every 10 seconds."),
        ("App updated or moved", "On launch the app re-points its login item and restarts the service on the new version."),
    ],

    "install_h2": "Install",
    "install_intro": "Install the app on each computer. On the same network you can pair right away; across networks, pick a way "
                     "for them to find each other. First time? Follow the <a href=\"{start}\">quick start</a>.",
    "mac_steps": [
        "Download the DMG, open it and drag PasteFerry into Applications.<span class=\"aside\">There's one for Apple silicon "
        "and one for Intel; the download page shows how to tell which your Mac has.</span>",
        "Open it from Applications.<span class=\"aside\">The first time, macOS says it can't verify the developer: PasteFerry "
        "isn't notarized by Apple yet. Go to System Settings → Privacy &amp; Security and click Open Anyway near the bottom. "
        "You only do this once.</span>",
        "Click Turn on background sync. When macOS asks whether PasteFerry may read the clipboard or find devices on your local "
        "network, allow both.",
        'For different networks, enable Settings → Use official service on each computer. No server address is needed. Custom hosting remains optional.',
    ],
    "win_steps": [
        "Run PasteFerry-Setup.exe. It installs for the current user, no administrator rights needed."
        "<span class=\"aside\">PasteFerry isn't code-signed yet, so Windows says \"Windows protected your PC\": click More info, "
        f"then Run anyway. To skip the installer, get the <a href=\"{REPO}/releases/latest/download/PasteFerry-windows-x64.zip\">portable zip</a>.</span>",
        "Open PasteFerry and click Turn on background sync. Allow it if Windows Firewall asks. After you close the window it "
        "keeps running in the notification area.",
        'For different networks, enable Settings → Use official service on each computer. No server address is needed. Custom hosting remains optional.',
    ],
    "linux_steps": [
        "On Ubuntu, Debian and similar, download the .deb and run <code>sudo apt install ./PasteFerry-linux-amd64.deb</code>; "
        "on other distributions, download the tar.gz, extract it and run <code>./PasteFerry/install.sh</code>.",
        "Open PasteFerry from the app menu and click Turn on background sync. From then on it starts when you log in.",
        'For different networks, enable Settings → Use official service on each computer. No server address is needed. Custom hosting remains optional.',
    ],
    "server_h3": "Ways to connect",
    "server_intro": 'Same network? Install and pair. Different networks? Enable official access without hosting a server. Custom servers and Tailscale remain available when you prefer to manage the connection yourself.',
    "server_compare": "Compare the three",
    "copy": "Copy",
    "copied": "Copied",
    "pairing": "<strong>Pairing:</strong> on one computer, choose Devices → Add device to show a 6-character code; enter it on "
               "the other, and confirm once both screens show the same two fingerprints. On a different subnet, also enter the "
               "address shown under the code. Any mix of Mac, Windows and Linux works the same way.",

    "specs_h2": "Specifications",
    "specs": [
        ("Version", '2.1.0, released <time datetime="2026-10-16">October 16, 2026</time>'),
        ("Mac app", "macOS 13 Ventura or later; separate downloads for Apple silicon and Intel, about 25 MB each, no Rosetta; "
                    "a universal version for both is about 47 MB"),
        ("Windows app", "Windows 10 and Windows 11, x64; installs per user, no administrator rights; about 23 MB download"),
        ("Linux app", "x64 and ARM64; .deb (Debian, Ubuntu and similar) or tar.gz (other distributions), about 22 MB; "
                      "needs an X11 desktop, or a Wayland desktop with XWayland (GNOME, KDE and others include it by default); GTK 3"),
        ("Server", "Optional, only for syncing across networks. Linux on x86_64 or ARM64 (systemd or Docker), macOS (Intel and Apple silicon), Windows 10 and 11 x64 (as a service); about 15 MB of RAM"),
        ("Languages", "English, Simplified Chinese, Traditional Chinese (Taiwan and Hong Kong) and Japanese; follows the system by default, switchable in Settings"),
        ("Content", "Plain text, rich text, images and screenshots, files, folders; each type can be turned off"),
        ("Transport", "WebRTC data channels, direct first with TURN relay fallback; on the same network, mDNS discovery and connection details over UDP 8788; every uplink probed on multi-uplink networks; files streamed as tar + zstd"),
        ("Encryption and identity", "End-to-end DTLS; Ed25519 device keys; every signaling message signed"),
        ("Pairing", "6-character code valid for 5 minutes, using an alphabet without look-alikes such as 0/O and 1/I; fingerprints compared by you"),
        ("Auto-sync limit", "50 MB by default, adjustable from 1 MB to 1 GB"),
        ("History", "3 days by default, adjustable from 1 hour to 30 days; expired items removed automatically"),
        ("Measured speed", "A 20 MB file transfers in about 330 ms over a direct LAN connection, checksums matching"),
        ("Built with", "Go for the background service and server (pure Go calling Win32 on Windows and speaking the X11 protocol on Linux), Flutter for the app (one codebase for all three systems), Protocol Buffers for messages"),
        ("License", f'<a href="{REPO}/blob/main/LICENSING.md">GNU AGPL-3.0</a>'),
    ],

    "faq_h2": "Questions",
    "faq": [('Why do .uuremote_ files keep appearing while UU Remote is running?',
      ['UU Remote writes temporary placeholder files to the system clipboard. Older PasteFerry clients recorded '
       'and forwarded them as ordinary files. Version 1.3.2 filters the confirmed placeholder naming pattern on '
       'send and receive, and removes placeholder-only history on startup. Ordinary empty, hidden and large '
       'files are unaffected. The 50 MB threshold only controls automatic transfer; UU’s UDP P2P setting is '
       'unrelated. Upgrade both clients. Other clipboard sync tools may still duplicate ordinary content.']),
     ('How is this different from Universal Clipboard?',
      ['Mac, Windows and Linux. Text, images, files, whole folders. No shared account required. Now connect '
       'across networks by enabling official access, with no server to host. Direct connections first, '
       'end-to-end encryption throughout.',
       'If your Macs share an Apple ID, sit side by side, and you mostly copy text, the built-in feature is '
       'already enough.']),
     ('Can the server see what I copy?',
      ['It cannot decrypt copied content: direct traffic stays between your devices, and relay forwards '
       'end-to-end encrypted packets. Official access retains essential device and quota statistics, described '
       'in the privacy notice. The public client and server source does not by itself prove how private '
       'production storage is operated.']),
     ('Does it connect directly on office networks with several uplinks?',
      ['Yes. These networks pick an uplink per destination. PasteFerry probes several servers from one port to '
       'learn its address on every uplink and sends them all to the peer, which tries each one. A Mac behind a '
       '4-uplink office network connects directly to a home connection in our tests. Peers behind NATs that '
       'randomize ports still fall back to the relay.']),
     ('Do I need my own server?',
      ['No. For different networks, enable <a href="{official}">official access</a> on both computers. No cloud '
       'VM or address to configure. It is off by default and currently includes 100 MB/day of relay traffic; '
       'direct connections do not use that allowance. LAN, <a href="{cloud}">self-hosting</a> and <a '
       'href="{ts}">Tailscale</a> remain available.']),
     ('No public IP at home. Can I sync between home and the office?',
      ['No. For different networks, enable <a href="{official}">official access</a> on both computers. No cloud '
       'VM or address to configure. It is off by default and currently includes 100 MB/day of relay traffic; '
       'direct connections do not use that allowance. LAN, <a href="{cloud}">self-hosting</a> and <a '
       'href="{ts}">Tailscale</a> remain available.']),
     ("Both computers are on the same network. Why can't they find each other?",
      ["Check that they're really on the same subnet: wired and Wi‑Fi in many offices, guest networks, and Wi‑Fi "
       "with client isolation don't let devices see each other. In that case, on the computer entering the code, "
       "click On a different subnet? and enter the address shown under the other computer's code.",
       'On Windows, allow PasteFerry when the firewall asks on first run; if you clicked Cancel, click Allow '
       'under Settings → Local network. On a Mac, turn on PasteFerry Daemon in System Settings → Privacy &amp; '
       'Security → Local Network. More in the <a href="{lan}#faq">same network guide</a>.']),
     ('Windows? Linux? iPhone?',
      ['Windows 10 and 11 (x64) are supported and work together with Macs. Files you paste on Windows are real '
       'files you can drop into Explorer or Office.',
       'Linux desktops on x64 and ARM64 are supported too: X11, or Wayland with XWayland (GNOME, KDE and others '
       'include it by default). Files you copy paste straight into the file manager.',
       "iPhone and iPad aren't supported: iOS doesn't let apps watch the clipboard in the background, so "
       "copy-and-it's-there can't work."]),
     ('Why does the system warn me the first time?',
      ["PasteFerry isn't notarized by Apple yet, and the Windows build isn't code-signed yet; both require paid "
       'certificates.',
       'On Mac, open System Settings → Privacy &amp; Security and click Open Anyway near the bottom. On Windows, '
       'click More info and then Run anyway in the "Windows protected your PC" prompt. You only do this once.']),
     ('What about large files?',
      ['By default, anything under 50 MB is pushed to the other computer as you copy it, so you can paste right '
       'away. Larger files sync as a record; pull them from the history when you need them. The limit is '
       'adjustable in Settings.']),
     ('Is it free? Can my company use it?',
      ['Yes. PasteFerry is open source under the GNU AGPL-3.0, so individuals and companies can use and modify '
       'it for free. If you distribute it, or offer a modified version over a network, you publish your source '
       'under the same license.']),
     ('Is official relay unlimited?',
      ['The current default is 100 MB per machine quota group per day, upload plus download, resetting at '
       'midnight Beijing time. Reinstalling does not guarantee a fresh allowance. The usage page shows the '
       'actual limit, remaining amount and update time. Direct connections still work when relay allowance runs '
       'out. Allowances and availability may change; check the app.'])],

    "closing": "Two computers, one clipboard.",
    "source": "View the source on GitHub",
    "footer_license": "PasteFerry is open source under the GNU AGPL-3.0.",
    "footer_label": "Footer",
    "footer_links": [
        ("All releases", f"{REPO}/releases"),
        ("Changelog", f"{REPO}/blob/main/CHANGELOG.md"),
        ("Design notes (Chinese)", f"{SITE}design/nat-traversal.html"),
        ("License", f"{REPO}/blob/main/LICENSING.md"),
        ("Report an issue", f"{REPO}/issues"),
    ],

    # ── Tutorial: run the server on your LAN (server/ page) ──
    "guide": {
        "title": "Run a PasteFerry server on your LAN: Windows, Mac and Linux",
        "description": "When your computers share a network you don't need a cloud server. Run PasteFerry's signaling and TURN server "
                       "on any always-on Windows, Mac or Linux machine with one script; it starts on boot. Includes app setup and troubleshooting.",
        "h1": "Run the server on your LAN",
        "lede": "If your computers are on the same home or office network, you don't need a cloud server. Pick a computer that stays on, "
                "download and run one script, then enter the address it prints into PasteFerry on each computer.",
        "notice": "<strong>On the same network you no longer need a server:</strong> install the app and pair; see the "
                  "<a href=\"{lan}\">same network guide</a>. This guide still helps when two subnets can't reach each other and "
                  "traffic has to go through a computer both can reach.",
        "toc": [("prepare", "Before you start"), ("install", "Install"), ("check", "Check"), ("clients", "Set up the apps"), ("faq", "Questions")],

        "overview_h2": "The server only introduces devices",
        "overview_intro": "Each computer stays connected to the server, which is how they find each other and exchange connection details. "
                          "Once connected, what you copy goes straight between the two computers. Only when that fails does it pass "
                          "through the server, still encrypted.",
        "diagram": {
            "title": "PasteFerry on a LAN",
            "desc": "On one LAN, a Mac and a Windows PC connect to the server at 192.168.1.20 for signaling. Copied content travels "
                    "directly between the two computers, and through the server's TURN relay only if a direct connection fails.",
            "lan": "Your LAN: the same Wi‑Fi or wired network",
            "server": "Server",
            "sig_note": "Signaling, TCP 8787",
            "relay_note": "Relay, UDP 3478, if direct fails",
            "client": "Running PasteFerry",
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
                f'Download <a href="{REPO}/releases/latest/download/pasteferry-server-windows-amd64.zip">pasteferry-server-windows-amd64.zip</a> '
                "and choose Extract All.",
                "Open the extracted folder and double-click <code>install.cmd</code>."
                "<span class=\"aside\">If you see \"Windows protected your PC\", click More info, then Run anyway; when User Account "
                "Control asks, click Yes. The server isn't code-signed yet, which is why Windows asks once.</span>",
                "Installation continues in a new window. The <code>ws://…/signal</code> address at the end is what you enter in the apps.",
            ],
            "download": "Download the Windows server",
            "notes": [
                "<strong>Runs as a Windows service:</strong> starts with Windows, no sign-in needed, no window. It appears as PasteFerry Server in Services.",
                "<strong>Firewall already allowed:</strong> the script adds an inbound rule for pasteferry-server.exe only, so there are no ports to open.",
                "<strong>Wrong address picked</strong> (several network adapters, say from a VM or VPN): type <code>cmd</code> in the "
                "address bar of the extracted folder, press Enter and run <code>install.cmd 192.168.1.20</code> with your address.",
                "Windows 10 and 11, x64.",
            ],
            "double_click": "Double-click",
            "explorer_label": "The extracted folder holds pasteferry-server.exe, install.cmd, install.ps1, README.md and uninstall.cmd; double-click install.cmd",
            "term_title": "Administrator: Windows PowerShell",
            "term_label": "The window after install.cmd finishes: it stops the old version, installs, registers the service, allows it through the firewall and starts it, then prints the address ws://192.168.1.20:8787/signal",
            "manage": [
                ("Status", "<code>Get-Service PasteFerryServer</code> in an administrator PowerShell, or find PasteFerry Server in Services"),
                ("Logs", "<code>C:\\ProgramData\\PasteFerry Server\\server.log</code>"),
                ("Restart", "<code>Restart-Service PasteFerryServer</code> in an administrator PowerShell"),
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
                "<strong>If the macOS firewall is on,</strong> the script ends with the command to allow pasteferry-server.",
                "Intel and Apple silicon, macOS 13 or later.",
            ],
            "term_label": "Output of ./install.sh in Terminal: it installs, writes the LaunchAgent and starts the server, then prints the address ws://192.168.1.20:8787/signal",
            "manage": [
                ("Logs", "<code>~/Library/Logs/PasteFerry/server.log</code>"),
                ("Restart", "<code>launchctl kickstart -k gui/$(id -u)/com.pasteferry.server</code>"),
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
                ("Status", "<code>systemctl status pasteferry-server</code>"),
                ("Logs", "<code>journalctl -u pasteferry-server -f</code>"),
                ("Restart", "<code>sudo systemctl restart pasteferry-server</code>"),
                ("Config", "<code>/etc/pasteferry/server.env</code>; restart the service after editing"),
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
            ("Open Settings", "Open PasteFerry on each computer, including the one running the server, and go to Settings."),
            ("Enter the address", "Under Signaling server, enter <code>ws://192.168.1.20:8787/signal</code> with your address and press Return or Enter."),
            ("Check the status", "When Connection shows a green Connected, you're set. Pair your devices under Devices and copying starts syncing."),
        ],
        "clients_alt": "PasteFerry Settings: the signaling server is ws://192.168.1.20:8787/signal and the connection status shows Connected",
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
                "on a Mac. On Linux, edit <code>-turn-ip</code> in <code>/etc/pasteferry/server.env</code> and restart the service. Then update "
                "the address on each computer.",
                "Reserving the address in your router avoids this.",
            ]),
            ("Does it sync away from the LAN?", [
                "No, the server is only reachable on your LAN. To sync from other networks, use a <a href=\"{cloud}\">public "
                "server</a> or <a href=\"{ts}\">Tailscale</a>.",
            ]),
            ("The server computer also runs PasteFerry. Which address?", [
                "The same LAN address works, and so does <code>ws://127.0.0.1:8787/signal</code>. The server and the app don't interfere.",
            ]),
        ],

        "closing": "One computer that stays on is enough.",
        "source": "Full server guide, including public servers and Docker",
    },

    # ── Run it on a public server (cloud/ page) ──
    "cloud": {
        "title": "Run the PasteFerry server on a public server: cloud VM guide",
        "description": "One script installs PasteFerry's signaling and TURN server on a Linux VM with a public IP. Every computer "
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
            "title": "PasteFerry with a public server",
            "desc": "The public server is at 203.0.113.10. A Mac on an office network and a Windows PC on a home network both "
                    "connect to it for signaling; copied content goes directly between them, and through the server's relay if that fails.",
            "server": "Public server",
            "server_sub": "Cloud VM with a public IP",
            "office": "Office network",
            "home": "Home network",
            "client": "Running PasteFerry",
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
            "<strong>Registers a systemd service</strong> that starts on boot; the config is <code>/etc/pasteferry/server.env</code>.",
            "<strong>Upgrade:</strong> download the new version, extract it and run <code>sudo ./install.sh</code> again; your config is kept.",
            "<strong>Docker works too:</strong> see the <a href=\"{deploy}\">server deployment guide</a>.",
        ],
        "term_label": "Output of sudo ./install.sh 203.0.113.10: it installs, writes the systemd unit and config and starts the server, then prints the address ws://203.0.113.10:8787/signal",

        "check_h2": "Check that it's reachable",
        "check_intro": "Open <code>http://203.0.113.10:8787/healthz</code> (with your public IP) in a browser on any computer. If you see this line, it works.",
        "check_label": "A browser showing http://203.0.113.10:8787/healthz, with the server's status ok and its version",

        "clients_steps": [
            ("Open Settings", "Open PasteFerry on each computer and go to Settings."),
            ("Enter the public address", "Under Signaling server, enter <code>ws://203.0.113.10:8787/signal</code> with your public IP and press Return or Enter."),
            ("Check the status", "When Connection shows a green Connected, you're set. If you haven't paired yet, pair under Devices."),
        ],
        "clients_alt": "PasteFerry Settings: the signaling server is ws://203.0.113.10:8787/signal and the connection status shows Connected",

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
                "and a home IP may change. If all your computers are at home, they need no server at all on the "
                "<a href=\"{lan}\">same network</a>; to sync across networks without port forwarding, see <a href=\"{ts}\">Tailscale</a>.",
            ]),
        ],

        "closing": "One cloud VM,<br>sync from anywhere.",
        "source": "Full server guide, including Docker and every option",
    },

    # ── Quick start (start/ page) ──
    "start": {
        "title": "Quick start: set up PasteFerry in three steps",
        "description": "New to PasteFerry? Install the app on each computer. On the same network, just pair, with no server; across "
                       "networks, use a public server or Tailscale. Then copy and paste between Mac and Windows.",
        "h1": "Quick start",
        "lede": "Install the app, pick a way to connect, and pair your computers. On the same network it takes a few minutes; "
                "after that, copying just syncs.",
        "toc": [("app", "Install the app"), ("server", "Pick a way to connect"), ("pair", "Pair")],
        "step": "Step {n}",

        "app_h2": "Install the app on each computer",
        "app_intro": "One on every computer you want to sync. Mac, Windows and Linux mix freely.",
        "mac": [
            "Download the DMG, open it and drag PasteFerry into Applications.<span class=\"aside\">There's one for Apple silicon "
            "and one for Intel; the download page helps you choose.</span>",
            "Open it from Applications.<span class=\"aside\">The first time, macOS says it can't verify the developer: go to "
            "System Settings → Privacy &amp; Security and click Open Anyway near the bottom.</span>",
            "Click Turn on background sync. When macOS asks about clipboard access and finding devices on your local network, "
            "allow both.",
        ],
        "win": [
            "Run PasteFerry-Setup.exe; no administrator rights needed.<span class=\"aside\">If you see \"Windows protected your PC\", "
            "click More info, then Run anyway.</span>",
            "Open PasteFerry and click Turn on background sync. Allow it if Windows Firewall asks. After you close the window it "
            "keeps running in the notification area.",
        ],
        "linux": [
            "On Ubuntu, Debian and similar, install the .deb; on other distributions, extract the tar.gz and run the "
            "<code>install.sh</code> inside.",
            "Open PasteFerry from the app menu and click Turn on background sync.",
        ],
        "app_alt": "PasteFerry the first time it opens, asking you to turn on background sync",

        "server_h2": "Pick a way to connect",
        "server_intro": 'Same network? Install and pair. Different networks? Enable official access without hosting a server. Custom servers and Tailscale remain available when you prefer to manage the connection yourself.',

        "pair_h2": "Pair your computers",
        "pair_intro": "Pair each two computers once, and you're done.",
        "pair_steps": [
            "On one computer, open Devices and click Add device. It shows a 6-character code, valid for 5 minutes, and this "
            "computer's address below it.",
            "On the other, click Add device too, switch to Enter a code and type it in."
            "<span class=\"aside\">On a different subnet, or with Tailscale, click On a different subnet? and enter the other "
            "computer's address as well.</span>",
            "Both screens show the same two fingerprint lines. Check them character by character, then confirm."
            "<span class=\"aside\">This stops anyone in the middle from posing as your device. Don't skip it.</span>",
        ],
        "pair_alt": "Add device: the pairing code K7MP2X, with this computer's address 192.168.1.20 below it",

        "done_h2": "Now copy something.",
        "done_text": "Copy text, an image or a file on one computer and press ⌘V or Ctrl+V on the other. Everything copied on "
                     "either side is in the history.",
        "help": "Need help? <a href=\"{lan}#faq\">Same network questions</a> · <a href=\"{ts}#faq\">Tailscale questions</a> · "
                "<a href=\"{cloud}#faq\">Public server questions</a> · <a href=\"{issues}\">Report an issue</a>",
    },

    # ── Same network (lan/ page), no server ──
    "lan": {
        "title": "Use PasteFerry on one network: no server, just pair",
        "description": "When your computers share a Wi‑Fi or wired network, PasteFerry needs no server: install the app, enter a "
                       "6-character pairing code, and sync the clipboard and files between Mac and Windows. On a different "
                       "subnet, you enter one address when pairing.",
        "h1": "On the same network",
        "lede": "Computers behind the same router don't need a server. Install the app, enter a pairing code, and copying syncs from then on.",
        "toc": [("prepare", "Before you start"), ("pair", "Pair"), ("subnet", "Different subnets"), ("faq", "Questions")],
        "step": "Step {n}",

        "overview_h2": "No server",
        "overview_intro": "Two computers on the same network find each other over mDNS, exchange connection details directly, "
                          "and open an encrypted direct connection. No server, and no computer that has to stay on.",
        "diagram": {
            "title": "PasteFerry on one network",
            "desc": "On one LAN, a Mac (192.168.1.20) and a Windows PC (192.168.1.35) find each other over mDNS and exchange "
                    "connection details over UDP 8788. Copied content goes directly between them; there is no server.",
            "lan": "LAN: the same Wi‑Fi or wired network",
            "client": "Running PasteFerry",
            "find_note": "Found over mDNS, connection details over UDP 8788",
            "direct": "Direct, end-to-end encrypted",
            "direct_note": "Content stays between the two computers",
        },
        "overview_points": [
            "<strong>Found automatically.</strong> Each computer announces itself on the LAN over mDNS, with only an anonymous "
            "token that changes daily, so nobody can tell whose computer it is. Its name is added only while Search nearby is open.",
            "<strong>Direct.</strong> Content goes straight between the two computers, at the speed of your network.",
            "<strong>Same security.</strong> You compare two fingerprint lines when pairing, and every message after that is "
            "signed, so others on the network can't pose as your device.",
        ],
        "prep": [
            ("The same network", "Both computers behind the same router; Wi‑Fi and wired can mix. Guest networks and Wi‑Fi with "
                                 "client isolation keep devices apart, and in many offices wired and Wi‑Fi are separate subnets: "
                                 "see <a href=\"#subnet\">different subnets</a> below."),
            ("Windows Firewall", "On first run, Windows asks whether PasteFerry may get through the firewall; allow it. If you "
                                 "clicked Cancel, click Allow under Settings → Local network in PasteFerry, which asks for "
                                 "administrator approval once."),
            ("Local Network on macOS", "From macOS 15, the first time it reaches a device on your network, macOS asks whether "
                                       "PasteFerry Daemon may find devices on your local network; allow it. If you declined, turn "
                                       "it on in System Settings → Privacy &amp; Security → Local Network."),
            ("No server address", "Leave the server field in Settings empty. To sync across networks later, add a server or use "
                                  "Tailscale; paired devices stay paired."),
        ],

        "pair_h2": "Pair your computers",
        "pair_intro": "Pair each two computers once. If the app isn't installed yet, see the <a href=\"{start}#app\">quick start</a>.",
        "pair_steps": [
            "On one computer, open Devices and click Add device. It shows a 6-character code, valid for 5 minutes.",
            "On the other, click Add device too, switch to Enter a code, type it in and click Continue.",
            "Both screens show the same two fingerprint lines. Check them character by character, then confirm. Copying syncs from now on.",
        ],
        "pair_alt": "Add device: the pairing code K7MP2X, with this computer's address 192.168.1.20 below it",
        "nearby_h3": "Without a code: Search nearby",
        "nearby_steps": [
            "On both computers, open Devices and click Search nearby. Only while this window is open can devices on the network see this computer's name.",
            "When the other one shows up on the radar, click it. Both screens show the same two fingerprint lines; check them, then confirm on each.",
            "This works within one subnet only. If the other computer doesn't show up, use a pairing code as above.",
        ],
        "nearby_alt": "Nearby devices: this computer in the middle, with Design PC and Meeting room iMac around it, plus the already paired MacBook Air",

        "subnet_h2": "Different subnets",
        "subnet_intro": "Local discovery only reaches one subnet. When your computers are on different subnets that can still "
                        "reach each other, as in many offices, you enter one address when pairing.",
        "subnet_steps": [
            "The computer showing the code lists its address below it, such as <code>192.168.1.20</code>.",
            "On the computer entering the code, click On a different subnet?, enter that address, and click Continue.",
            "After pairing, each computer remembers the other's address. If one changes address, the other finds it again on its "
            "own; if both change, open the device's … menu under Devices and choose Enter its address.",
        ],
        "subnet_alt": "Enter a code: the code K7MP2X, with the address field opened below it holding 192.168.1.20",
        "subnet_after": "If the subnets can't reach each other at all, such as a guest network kept apart from the office network, "
                        "traffic has to go through a computer both can reach: see <a href=\"{server}\">run a server on your LAN</a>.",

        "faq": [
            ("The other computer can't be found", [
                "Make sure both are on the network with background sync on. Then try entering the other computer's address on the "
                "one entering the code: if that works, they're on different subnets, or the network blocks mDNS.",
                "Still nothing? Check the other computer's firewall. On Windows, click Allow under Settings → Local network; on a "
                "Mac, make sure PasteFerry Daemon is on under System Settings → Privacy &amp; Security → Local Network and the "
                "macOS firewall isn't blocking it.",
            ]),
            ("Does it sync away from this network?", [
                "No. It resumes when you're back, and anything copied in between is added to the other computer's history. To sync "
                "across networks, add a <a href=\"{cloud}\">public server</a> or use <a href=\"{ts}\">Tailscale</a>; paired "
                "devices stay paired.",
            ]),
            ("What if an address changes?", [
                "Usually nothing to do. On one subnet, mDNS finds the computer again at any time; across subnets, it reconnects as "
                "long as one of the two keeps its address. If both change, enter the address once more under Devices.",
                "For a desktop that never moves, you can reserve its address in your router's DHCP settings.",
            ]),
            ("Can others on the network see my computer?", [
                "They can see that a PasteFerry device is there, but not whose: the broadcast carries an anonymous token that changes "
                "daily, with no device name or ID. The name is added only while you have Search nearby open, and removed when you close it. "
                "Messages from devices you haven't paired with are dropped, and pairing requires "
                "comparing fingerprints on both screens.",
            ]),
            ("Which ports does it use?", [
                "UDP 8788 for connection details (a random port if that's taken), UDP 5353 for mDNS, and another random UDP port "
                "for the data once connected. Nothing to configure on your router.",
            ]),
        ],

        "closing": "One network,<br>nothing to set up.",
        "source": "Design notes on local network mode (in Chinese)",
    },

    # ── Use it across networks with Tailscale (tailscale/ page), no server ──
    "ts": {
        "title": "Use PasteFerry across networks with Tailscale: no server, home and office in sync",
        "description": "Install Tailscale on each computer under one account and enter the other computer's Tailscale address "
                       "when pairing: no public IP, no server and no computer that has to stay on. The clipboard syncs from "
                       "the office or anywhere else.",
        "h1": "Across networks with Tailscale",
        "lede": "No public IP, and no wish to run a server? Install Tailscale on each computer and enter the other computer's "
                "Tailscale address when pairing. Syncs at home and at the office.",
        "toc": [("join", "Install Tailscale"), ("address", "Find the address"), ("pair", "Pair"), ("faq", "Questions")],
        "step": "Step {n}",

        "overview_h2": "As if on one LAN",
        "overview_intro": "Tailscale connects your devices into a private network: every device signed in to the same account "
                          "gets a fixed address starting with 100, reachable from the others wherever they are. PasteFerry "
                          "remembers the other computer's address when you pair, and connects as if both were on one network.",
        "diagram": {
            "title": "PasteFerry over Tailscale",
            "desc": "A Mac at the office and a Windows PC at home both run Tailscale and each has an address starting with 100 "
                    "on the same private network. Copied content goes directly between them; there is no server.",
            "net": "Tailscale network: every device has a 100.x address",
            "mac_place": "At the office",
            "win_place": "At home",
            "direct": "Direct, encrypted twice",
            "direct_note": "Content stays between the two computers",
        },
        "overview_points": [
            "<strong>No server,</strong> and no computer that has to stay on: two computers that are on and connected to "
            "Tailscale sync.",
            "<strong>Enter the address once.</strong> A Tailscale address doesn't change; you enter it when pairing and the two "
            "computers find each other from then on.",
            "<strong>Encrypted twice.</strong> Inside Tailscale's WireGuard tunnel runs PasteFerry's own end-to-end encrypted "
            "connection. Tailscale can't see the content.",
        ],

        "prep": [
            ("A Tailscale account", "Free for personal use; sign in with Google, Microsoft, GitHub or Apple. Use the same account on every device."),
            ("Tailscale on every computer", "There are clients for Windows, macOS and Linux. Every computer that runs PasteFerry "
                                           "needs it, connected."),
            ("PasteFerry installed", "The app on each computer with background sync on; see the <a href=\"{start}#app\">quick start</a>. "
                                   "Leave the server address empty."),
            ("ZeroTier works too", "So does any tool that gives each device a fixed virtual address: enter that address when "
                                   "pairing. This guide uses Tailscale."),
        ],

        "join_h2": "Install Tailscale on each computer",
        "join_intro": "On every computer you sync, install it and sign in to the same account.",
        "join_steps": [
            "Download the client for your system from the Tailscale website and install it."
            "<span class=\"aside\">On Linux, run <code>curl -fsSL https://tailscale.com/install.sh | sh</code>, then "
            "<code>sudo tailscale up</code>, and open the link it prints to sign in.</span>",
            "Open Tailscale and sign in, using the same account on every computer.",
            "The Tailscale admin console then lists all your devices, each with an address starting with 100.",
        ],
        "join_download": "Download Tailscale",
        "machines_cols": ["Machine", "Tailscale address", "OS"],
        "machines_label": "The Tailscale admin console's machine list (illustration): macbook-pro is 100.88.20.5 and home-pc is 100.92.41.17",
        "join_notes_h3": "Keep Tailscale connected",
        "join_notes": [
            "<strong>Windows:</strong> right-click the Tailscale icon in the notification area and check Run unattended under "
            "Preferences, so it stays connected after you sign out.",
            "<strong>Mac:</strong> make sure Tailscale opens at login.",
            "<strong>Disable key expiry:</strong> by default Tailscale asks devices to sign in again every 180 days. For a computer "
            "that rarely changes, open the menu at the end of its row on the Machines page and choose Disable key expiry.",
        ],

        "address_h2": "Find the other computer's Tailscale address",
        "address_intro": "On the computer you'll pair with, click the Tailscale icon (in the menu bar on a Mac, in the notification "
                         "area on Windows) to see its address, or run <code>tailscale ip -4</code> in a terminal.",
        "address_label": "Running tailscale ip -4 in PowerShell prints the computer's Tailscale address, 100.92.41.17",
        "address_after": "The steps below use <code>100.92.41.17</code>, the Windows PC at home. Use your own address.",

        "pair_h2": "Enter the address when pairing",
        "pair_intro": "Just like pairing on one network, with one extra address.",
        "pair_steps": [
            "On the Windows PC at home, open Devices, click Add device and note the 6-character code.",
            "On the Mac at the office, click Add device, switch to Enter a code and type it in; then click On a different "
            "subnet?, enter <code>100.92.41.17</code> and click Continue.",
            "Both screens show the same two fingerprint lines. Check them, then confirm. The two computers now sync wherever they are.",
        ],
        "pair_alt": "Enter a code: the code K7MP2X, with the Tailscale address 100.92.41.17 in the address field",
        "pair_after": "When both computers are at home, PasteFerry also finds the other one on the LAN directly, so data takes no detour.",

        "faq": [
            ("Devices show Direct, but it's slow", [
                "Direct in PasteFerry means the two computers connect directly on the Tailscale network. When Tailscale can't "
                "connect them peer to peer, it relays through its DERP servers, which can be far away.",
                "Run <code>tailscale ping other-address</code> in a terminal: via DERP means Tailscale is relaying; via an IP "
                "and port means it's peer to peer. Office and campus networks relay more often; you can run your own DERP "
                "server or try another mesh VPN.",
            ]),
            ("Pairing can't find the other computer", [
                "First make sure Tailscale is connected on both: run <code>tailscale ping other-address</code> and try again once "
                "it answers.",
                "If ping works but pairing doesn't, the other computer's firewall is likely blocking it. On Windows, click Allow "
                "under Settings → Local network in PasteFerry; on a Mac, turn on PasteFerry Daemon under System Settings → Privacy "
                "&amp; Security → Local Network.",
            ]),
            ("Does every computer need Tailscale?", [
                "Yes, a computer without it can't reach the 100.x addresses. Computers on the same network don't need it; see "
                "the <a href=\"{lan}\">same network guide</a>.",
            ]),
            ("I set up a server at home with the old guide. Do I still need it?", [
                "No. Clear the server address in Settings on each computer and uninstall the server at home; paired devices stay "
                "paired. Keeping it doesn't hurt either: a server and Tailscale work together.",
            ]),
            ("Can Tailscale see what I copy?", [
                "No. Tailscale connections are end-to-end encrypted with WireGuard, and inside runs PasteFerry's own DTLS-encrypted "
                "connection; only your two computers can open either layer. Tailscale's servers only exchange device keys and "
                "addresses, and its relays forward encrypted packets.",
            ]),
            ("Will a Tailscale address change?", [
                "Usually not. The address belongs to the device and survives restarts and network changes; it only changes if you "
                "remove the device from Tailscale and add it again. If it does, enter it once more under Devices.",
            ]),
            ("I have a public server. Do I need Tailscale?", [
                "No. Run the server on the machine with the public IP and enter that address on every computer; see the "
                "<a href=\"{cloud}\">public server guide</a>.",
            ]),
        ],

        "closing": "No server.<br>Your clipboard comes with you.",
        "source": "Back to the same network guide",
    },

    # ── Download (download/ page) ──
    "download": {
        "title": "Download PasteFerry: Mac (Apple silicon, Intel), Windows and Linux",
        "description": "PasteFerry for every system: Apple silicon and Intel versions for macOS 13 or later, an installer and a "
                       "portable zip for Windows 10 and 11, and .deb and tar.gz packages for Linux on x64 and ARM64. Free and open source.",
        "h1": "Download PasteFerry",
        "lede": "Free and open source, version {version}. Install it on each computer you want to sync; Mac, Windows and Linux mix freely.",
        "os_label": "Choose your system",
        "win": {
            "setup": "Download the installer",
            "portable": "Rather not run an installer? Get the <a href=\"{portable}\">portable zip</a>, extract it and run PasteFerry from there.",
            "steps": [
                "Run PasteFerry-Setup.exe. It installs for the current user, no administrator rights needed."
                "<span class=\"aside\">If you see \"Windows protected your PC\", click More info, then Run anyway.</span>",
                "Open PasteFerry and click Turn on background sync. Allow it if Windows Firewall asks.",
            ],
            "notes": [
                "Windows 10 and Windows 11, x64.",
                "The installer is about 23 MB, the portable zip about 29 MB.",
                "After you close the window, PasteFerry keeps running in the notification area.",
            ],
        },
        "mac": {
            "apple": "Apple silicon (M1 and later)",
            "intel": "Intel",
            "universal": "A <a href=\"{universal}\">universal version</a> that runs on both is also available, at twice the size.",
            "steps": [
                "Open the DMG you downloaded and drag PasteFerry into Applications.",
                "Open it from Applications.<span class=\"aside\">The first time, macOS says it can't verify the developer: go to "
                "System Settings → Privacy &amp; Security and click Open Anyway near the bottom.</span>",
                "Click Turn on background sync. When macOS asks about clipboard access and finding devices on your local network, "
                "allow both.",
            ],
            "chip_h": "Not sure which chip?",
            "chip": [
                "Click the Apple menu in the top-left corner of the screen and choose About This Mac.",
                "\"Chip: Apple M1\" (or M2, M3, M4…) means Apple silicon; \"Processor: … Intel Core …\" means Intel.",
                "Almost every Mac bought since late 2020 has Apple silicon. If you pick the wrong one, just download the other: "
                "on an Apple silicon Mac the Intel version needs Rosetta installed first, and the Apple silicon version won't "
                "open on an Intel Mac.",
                "Requires macOS 13 Ventura or later.",
            ],
        },
        "linux": {
            "deb_h": "Ubuntu, Debian and derivatives",
            "tar_h": "Other distributions",
            "x64": "x64",
            "arm64": "ARM64",
            "notes": [
                "Needs a graphical desktop: X11, or Wayland with XWayland, which GNOME, KDE and other major desktops include by default.",
                "Not sure whether you have x64 or ARM64? Run <code>uname -m</code> in a terminal: <code>x86_64</code> means x64, "
                "<code>aarch64</code> means ARM64.",
                "The .deb installs to /opt/PasteFerry and adds PasteFerry to the app menu. The tar.gz runs from wherever you extract "
                "it, and <code>install.sh</code> adds it to the app menu.",
                "Once it's open, click Turn on background sync; from then on it starts when you log in. GNOME shows no tray icon "
                "by default, so open the window from the app menu.",
                "Needs GTK 3. Tested on Ubuntu 24.04 with GNOME (Wayland).",
            ],
        },
        "more_h2": "More",
        "more": [
            ("Server", "Two computers on the same network need no server. The ways to connect across networks are in the "
                       "<a href=\"{server}\">quick start</a>; each guide has its own server download."),
            ("Checksums", "SHA-256 checksums of every file are in <a href=\"{sums}\">SHA256SUMS</a>."),
            ("What's new and older versions", "See the <a href=\"{notes}\">release notes</a> for what changed in this version; "
                                              "older versions are on <a href=\"{releases}\">GitHub Releases</a>."),
        ],
    },

    "ld_description": "Open-source app that syncs the clipboard and files between Mac, Windows and Linux computers. Copy text, "
                      "images, files or folders on one computer and paste on the other. No server is needed on the same "
                      "network. Devices connect directly over WebRTC with end-to-end DTLS encryption; across networks, a "
                      "server the user hosts introduces them and relays over TURN when needed.",
    "ld_os": "macOS 13 or later (Intel and Apple silicon); Windows 10, Windows 11 (x64); Linux (x64, ARM64)",
    "ld_currency": "USD",
    "ld_features": [
        "Syncs plain text, rich text, images, files and folders",
        "Works between Mac, Windows and Linux, across accounts and networks",
        "No server needed on the same network, with automatic mDNS discovery",
        "Direct WebRTC connection with automatic TURN relay fallback",
        "End-to-end DTLS encryption with human-verified fingerprints at pairing",
        "Items under 50 MB pushed on copy, larger files pulled on demand",
        "3-day history with image preview",
        "Interface in English, Simplified Chinese, Traditional Chinese and Japanese",
    ],
}

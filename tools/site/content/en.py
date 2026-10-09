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
    "nav_how": "How it works",
    "nav_install": "Install",
    "nav_faq": "FAQ",
    "lang_label": "Language",

    "h1": "Copy on this computer.<br>Paste on that one.",
    "lede": "Mac and Windows. Text, images, files, whole folders. No shared account required, no shared network required. "
            "Data moves directly between your devices, and stays encrypted even when it passes through your own server.",
    "dl_mac": "Download for Mac",
    "dl_win": "Download for Windows",
    "dl_note": "Version 1.3.0, free and open source. macOS 13 or later, Intel and Apple silicon; Windows 10 and 11.",
    "stage_label": "Demo: Command C on a MacBook turns a file into a sheet that flies over a direct connection to a Windows "
                   "laptop, where Ctrl V pastes it",
    "stage_file": "Q3 review - final.pptx",
    "stage_win": "Windows laptop",
    "stage_arc": "Direct, end-to-end encrypted",
    "replay": "Play again",

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
    "install_intro": "Install the app on each computer, and set up a server both of them can reach. If both share a LAN, the "
                     "server can run on one of your Macs.",
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
    "server_intro": "Any Linux, as root. On ARM, replace <code>amd64</code> with <code>arm64</code>.",
    "copy": "Copy",
    "copied": "Copied",
    "public_ip": "YOUR_PUBLIC_IP",
    "server_after": "The script registers a systemd service that starts on boot and prints the address to enter in the app. "
                    "Open these ports in your firewall or security group:",
    "ports_cols": ["Port", "Purpose"],
    "ports": [("8787/tcp", "Signaling"), ("3478/udp", "STUN and TURN"), ("32768–60999/udp", "Relay ports")],
    "server_docker": f'Docker works too; see the <a href="{REPO}/blob/main/server/deploy/README.md">server deployment guide</a>.',
    "pairing": "<strong>Pairing:</strong> on one computer, choose Devices → Add device to show a 6-character code; enter it on "
               "the other, and confirm once both screens show the same two fingerprints. Mac with Windows, or two Windows PCs, "
               "works the same way.",

    "specs_h2": "Specifications",
    "specs": [
        ("Version", '1.3.0, released <time datetime="2026-10-09">October 9, 2026</time>'),
        ("Mac app", "macOS 13 Ventura or later; universal binary for Intel and Apple silicon, no Rosetta; about 45 MB download"),
        ("Windows app", "Windows 10 (21H2 or later) and Windows 11, x64; installs per user, no administrator rights; about 21 MB download"),
        ("Server", "Linux on x86_64 or ARM64 (systemd or Docker), or macOS; about 15 MB of RAM"),
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
            "Yes. There is no public server, so nobody's data passes through anyone else's machine. On a LAN, run the server on "
            "one of your Macs. Across networks you need a machine with a public IP; the smallest cloud VM is enough.",
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

<p align="center">
  <img src="docs/assets/icon.png" width="112" height="112" alt="CopySync icon">
</p>

<h1 align="center">CopySync</h1>

<p align="center">
  Copy on one computer, paste on the other. Mac and Windows both.<br>
  Text, images, files and whole folders — across accounts and across networks.
</p>

<p align="center">
  <a href="https://github.com/baiyuze/copysync/releases/latest/download/CopySync.dmg"><strong>Download for Mac</strong></a>
  ·
  <a href="https://github.com/baiyuze/copysync/releases/latest/download/CopySync-Setup.exe"><strong>Download for Windows</strong></a>
  ·
  <a href="https://baiyuze.github.io/copysync/en/">Website</a>
  ·
  <a href="README.zh-CN.md">简体中文</a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/baiyuze/copysync?color=0A66D8" alt="Release">
  <img src="https://img.shields.io/badge/macOS-13%2B-1D1D1F" alt="macOS 13+">
  <img src="https://img.shields.io/badge/Windows-10%20%2F%2011-1D1D1F" alt="Windows 10 / 11">
  <img src="https://img.shields.io/badge/Intel%20%2B%20Apple%20silicon-universal-1D1D1F" alt="Universal binary">
  <img src="https://img.shields.io/badge/license-AGPL--3.0-6E6E73" alt="AGPL-3.0 License">
</p>

<p align="center">
  <img src="docs/assets/screenshots/en/history.png" alt="CopySync's history view: files, images and text arriving from other Macs">
</p>

## Why

If you work on two computers, you move things between them all day: a snippet of text, a screenshot, a file.

The built-in options each get in the way:

- **Universal Clipboard** needs both Macs on the same Apple ID, close to each other, with Bluetooth and Wi‑Fi on. A work Mac and a personal Mac rarely share an Apple ID, and large files often fail halfway.
- **AirDrop** makes you pick a device and accept on the other side every time, and it only moves files, not the clipboard.
- **Messaging yourself** through a chat app means manual downloads, recompressed images, and your files sitting on someone else's server.
- **With one Mac and one Windows PC**, none of the built-in options work at all.

With CopySync you pair two computers once. After that, copy on one and paste on the other — Mac to Mac, Mac to Windows, or Windows to Windows. Data goes directly between the devices; when a direct connection can't be made, it is relayed through a server you run yourself. Everything is end-to-end encrypted and the server cannot read it.

| | CopySync | Universal Clipboard | AirDrop | Chat apps |
|---|:-:|:-:|:-:|:-:|
| Different Apple IDs | ✓ | — | Needs accepting | ✓ |
| Different networks (home ↔ office) | ✓ | — | — | ✓ |
| Paste right away, nothing to accept | ✓ | ✓ | — | — |
| Between Mac and Windows | ✓ | — | — | ✓ |
| Files and whole folders | ✓ | Unreliable | ✓ | Zipped |
| History of what you copied | ✓ | — | — | Scroll back |
| Content stays off third-party servers | ✓ | ✓ | ✓ | — |

## Features

- **Text, rich text, images, files and folders.** What you copy is what comes out on the other side.
- **Size-aware transfer.** Anything under 50 MB (adjustable) is pushed as soon as you copy, so the other Mac can paste immediately. Larger files sync as a record you can pull when you need them.
- **History.** Everything copied on any of your computers in the last 3 days (adjustable), ready to put back on the clipboard. Click an image to preview and zoom it.
- **Direct first, relay as fallback.** Devices connect peer-to-peer; behind symmetric NAT or a strict firewall, traffic switches to the TURN relay on your server automatically.
- **Fingerprint check when pairing**, so neither the server nor anyone on the network can impersonate your device.
- **Runs in the background.** Keeps syncing with the window closed, starts at login, recovers from crashes. On Windows it lives in the notification area.
- **Mac and Windows.** Native on Intel and Apple silicon Macs; Windows 10 and 11. Pasted files are real files you can drop into Explorer, WeChat or Office. Light and dark appearance.
- **English, Chinese and Japanese.** The interface follows your system language, and you can switch it in Settings.
- **Passwords stay put** (Windows). Content that a password manager marks as "don't record" is skipped: not synced, not kept in the history.

<table>
  <tr>
    <td><img src="docs/assets/screenshots/en/devices.png" alt="Devices: paired Macs and whether each is connected directly or through the relay"></td>
    <td><img src="docs/assets/screenshots/en/verify.png" alt="Pairing: compare the fingerprint shown on both Macs before confirming"></td>
  </tr>
  <tr>
    <td><img src="docs/assets/screenshots/en/settings.png" alt="Settings: server address, auto-sync size limit, content types"></td>
    <td><img src="docs/assets/screenshots/en/history-dark.png" alt="History in dark appearance"></td>
  </tr>
</table>

## Supported systems

| | Requirement |
|---|---|
| **Mac app** | macOS 13 Ventura or later; universal binary for Intel and Apple silicon, no Rosetta needed |
| **Server** | Any Linux on x86_64 or ARM64 with systemd or Docker; it can also run on one of your Macs |
| **Windows app** | Windows 10 (21H2 or later) and Windows 11, x64; runs under emulation on Windows on ARM |
| **iPhone / iPad** | Not supported. iOS does not let apps watch the clipboard in the background, so copy-and-it's-there can't work the way it does on the Mac |

## Install

There are two parts: **the app on each computer**, and **a server both can reach**. If both computers are on the same LAN, the server can run on one of the Macs.

### 1. Deploy the server

On a Linux server, as root:

```bash
curl -LO https://github.com/baiyuze/copysync/releases/latest/download/copysync-server-linux-amd64.tar.gz
tar xzf copysync-server-linux-amd64.tar.gz
cd copysync-server-linux-amd64
sudo ./install.sh <public IP of the server>
```

The script registers a systemd service that starts on boot and prints the address to enter in the app. Use the `linux-arm64` package on ARM servers.

Open these ports in your firewall or cloud security group:

| Port | Protocol | Purpose |
|---|---|---|
| 8787 | TCP | Signaling (WebSocket) |
| 3478 | UDP | STUN / TURN |
| 32768–60999 | UDP | TURN relay ports, assigned by the OS |

For Docker, or for LAN-only use, see the [server deployment guide](server/deploy/README.md).

### 2. Install the app on each computer

**Mac**

1. Download [CopySync.dmg](https://github.com/baiyuze/copysync/releases/latest/download/CopySync.dmg), open it and drag CopySync into Applications.
2. Open CopySync from Applications.

   The first time, macOS says it can't verify the developer: CopySync isn't notarized by Apple yet (that requires a paid developer account). Go to **System Settings → Privacy & Security**, scroll down and click **Open Anyway**. You only need to do this once.

3. Click **Turn on background sync**.
4. In **Settings → Signaling server**, enter `ws://<server address>:8787/signal` and press Return.
5. Recent versions of macOS ask whether CopySync Daemon may read the clipboard. Allow it; you can switch it to always allow under **System Settings → Privacy & Security → Paste from Other Apps**.

**Windows**

1. Download [CopySync-Setup.exe](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-Setup.exe) and run it. It installs for the current user and needs no administrator rights.

   CopySync isn't code-signed yet, so Windows says "Windows protected your PC": click **More info → Run anyway**. To skip the installer, download the [portable zip](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-windows-x64.zip), extract it and open CopySync.

2. Open CopySync (it opens when setup finishes; there's also a desktop shortcut) and click **Turn on background sync**.
3. In **Settings → Signaling server**, enter `ws://<server address>:8787/signal` and press Enter.
4. If Windows Firewall asks about copysyncd, allow it.

After you close the window, CopySync keeps running in the notification area; open it again from the tray icon or the desktop shortcut.

### 3. Pair

1. On one computer: **Devices → Add device → Show a code**.
2. On the other: **Add device → Enter a code**, and type the 6 characters.
3. Both screens show the same two lines of fingerprints, one for each device. **Confirm only if the two screens match line for line.**

From then on, copy on either one and paste on the other.

## How it works

```mermaid
flowchart LR
    subgraph A["Computer A (Mac or Windows)"]
        UA["CopySync app"] <-- "gRPC" --> DA["Background service"]
    end
    subgraph B["Computer B (Mac or Windows)"]
        DB["Background service"] <-- "gRPC" --> UB["CopySync app"]
    end
    S[("Your server\nsignaling + TURN")]
    DA -- "Signaling (WebSocket)" --> S
    DB -- "Signaling (WebSocket)" --> S
    DA <== "WebRTC direct (DTLS)" ==> DB
    DA <-. "Relayed via TURN when direct fails" .-> S
    S <-.-> DB
```

| Component | Role | Built with |
|---|---|---|
| Background service `copysyncd` | Watches the clipboard, connects to peers, sends, receives and stores | Go. The clipboard uses cgo and AppKit on Mac, and calls Win32 directly on Windows (pure Go, no C compiler) |
| App `CopySync` | History, pairing, settings | Flutter, one codebase for Mac and Windows; talks to the service over local gRPC |
| Server `copysync-server` | Device discovery, pairing codes, signaling, TURN relay | Go, a single static binary using about 15 MB of RAM |
| Protocol `proto/` | Messages shared by all three | Protocol Buffers |

**What happens when you copy:**

1. The service notices the clipboard changed: on Mac it checks the change count at short intervals (reading the count and types needs no permission and never triggers a prompt); on Windows the system notifies it right away.
2. On a change it reads the content and picks a transport: text is inlined in the message; images and files under 50 MB are streamed right away (tar + zstd); anything larger is announced as a record only.
3. The receiving computer writes the data to a local cache and then onto its clipboard, so `⌘V` or `Ctrl+V` pastes real local files. File names that are fine on Mac but not on Windows (with a colon, say) get the matching full-width character on Windows.

**How the connection is made:**

1. Devices connect to the signaling server over WebSocket and receive STUN and TURN addresses.
2. All connections share one local UDP port. The background service probes several servers from that port to learn its public address, one per uplink on networks with several (see the next section).
3. During the handshake each side sends its LAN addresses, its public address on every uplink, and its TURN relay address; both sides try each other with WebRTC ICE at the same time. One side holds back its uplink addresses until it has received the other's and sent packets to them, because some routers lock a port when the other side's packet arrives first (see the next section).
4. A direct path wins if it works; relay paths must wait 5 seconds before they can be selected.
5. Once connected, each side tells the other whether it is relayed, and both show "relay" if either is.
6. If the connection still lands on the relay, the side that started it runs an ICE restart while idle to punch again, and moves to the direct path if that works (first after 10 seconds, then at intervals growing to 15 minutes).

Data flows over WebRTC data channels, encrypted with DTLS.

### Hole punching on networks with several uplinks

Office networks with two ISPs, or carriers with a pool of NAT addresses, have several uplinks and choose one per destination IP. One office network we measured had four:

| Uplink | Probe servers that landed on it | Public port |
|---|---|---|
| Uplink A | Our own server, Twilio | Unchanged |
| Uplink B | Bilibili, Google, Nextcloud | Unchanged |
| Uplink C | Cloudflare, FreeSWITCH | Rewritten, same for every destination |
| Uplink D | A service on Alibaba Cloud | Rewritten |

Asking one server reveals only one of these addresses. Packets to the peer may leave through another uplink, the peer sees an address it was never told about, and its router drops them. That is why CopySync 1.0 could only relay on such networks.

What 1.1 does:

- **One shared port.** pion used to open a new port for each STUN server, so each answer only applied to its own port. All connections now share one port, and every probe result applies to it.
- **Probe every uplink.** CopySync queries IP-diverse servers from that port in parallel (your own server plus public STUN servers such as Bilibili and Cloudflare), deduplicates by IP and drops addresses poisoned by DNS.
- **Tell the peer all of them.** Each uplink's address goes to the peer as a candidate; the peer checks each one, and the one matching the real uplink gets through. Uplinks that rewrite ports report their real mapping.
- **Remember uplinks.** If a round misses an uplink that preserves ports, its address is filled in as uplink IP plus local port.

A Mac behind that 4-uplink office network now connects directly to a home connection, about 3 seconds after reaching the server. The full design, measurements and validation are in [NAT traversal on multi-uplink networks](design/nat-traversal.md) (Chinese).

The **NAT lab** (`tools/natlab`) builds real topologies with Linux network namespaces and iptables and runs CopySync's actual connection code through eleven scenarios on every commit:

| Scenario | Result |
|---|---|
| Single uplink ↔ home router | Direct |
| Two uplinks, previous version (control) | Relay |
| Two uplinks ↔ home router | Direct |
| Four uplinks ↔ home router | Direct |
| Two uplinks, second one not probed: no history / with history | Relay / Direct |
| Two uplinks ↔ two uplinks | Direct |
| Two uplinks ↔ NAT that randomizes ports | Relay |
| Blocked at first, network recovers later | Relay first, then direct |
| Office router locks a port for unsolicited packets | Direct |
| Home router locks a port for unsolicited packets | Relay, then direct after the first retry |

### Security model

- **The server keeps no accounts.** A device's identity is an Ed25519 key, and its ID is derived from the public key, so the server can verify that an ID belongs to a key without a user table.
- **Fingerprints are compared by a person when pairing.** Both Macs show the same two lines: the fingerprint of each side's public key (the first 60 bits of its SHA-256, like `R8NF-2WTC-QL5J`), in a fixed order. The server could swap keys in transit, but then the two screens would no longer match — this step is what stops a man in the middle. The fingerprints are shown separately rather than combined into one short code: an attacker controlling both forged keys could find a matching combined code with a birthday attack, while separate fingerprints need a preimage attack per key, roughly a billion times harder.
- **Every signaling message after that is signed**, so the server cannot forge a device. The DTLS certificate fingerprint of the direct channel is also sent signed and checked against the actual certificate after the handshake; a mismatch drops the connection.
- **The relay can't read anything either.** TURN only forwards encrypted UDP packets. Relay credentials are issued per device and expire after 12 hours.
- **Public STUN servers see only your public IP.** To find every uplink on multi-uplink networks, CopySync probes a few public STUN servers by default. They see your public IP, as with any WebRTC app, never content. Turn off **Probe with public servers** in Settings to use only your own server.
- **Device key.** On Mac it's stored readable only by you; on Windows it's additionally encrypted with DPAPI, so another user or another computer can't decrypt a copied file.
- **No App Sandbox**, because CopySync needs to read files at whatever path you copy them from.

### Fallbacks

| Situation | What happens |
|---|---|
| No direct path (symmetric NAT, corporate firewall) | Switches to the TURN relay on your server, still encrypted; the app shows "relay" |
| Public STUN servers unreliable (common in mainland China) | The server runs its own STUN and clients prefer it |
| Multi-uplink networks that pick an uplink per destination (common in offices) | All connections share one local port; CopySync probes several servers from it to learn its address on every uplink and sends them all to the peer, filling gaps from history. See `copysync-cli nat` |
| Landed on the relay although a direct path works (the peer just started, its addresses arrived a moment late) | While idle, the side that started the connection runs an ICE restart to punch again and switches to direct if it works; never during a transfer. The connection stays up and only pauses for a few seconds |
| The other Mac sleeps, loses its network or restarts, and the connection dies silently | The dead connection is closed and reconnected right away; if the other side rebuilt its connection, this side follows |
| The other Mac can't be reached when you copy | The latest 20 items are kept and sent once the connection is back. They only go into the history and never overwrite what's on the other Mac's clipboard now |
| Signaling connection drops | Reconnects with exponential backoff (1 s up to 30 s, with jitter); status is shown live |
| Pairing confirmed on one side before the other | Early handshake messages are held and replayed once the other side confirms; the direct link is up within tens of milliseconds |
| Large files | Above the limit only a record is synced; if the source file is gone when you pull, you get a clear message |
| Apps that only accept plain text | Rich text is written as both HTML and plain text, so no markup leaks into the paste |
| Clipboard access not granted yet | Content isn't read (it would block on a prompt a background process can't show); the app asks you to grant access |
| Background service crashes | launchd restarts it, throttled to once every 10 seconds |
| App updated in place, or moved | On launch the app re-points the login item and restarts the service on the new version |
| History and cache grow | Expired items are removed automatically; 3 days by default |

### Lessons learned

Each of these was found by testing and has a regression test. The full investigation is in [spikes/results.md](spikes/results.md) (Chinese).

- **macOS clipboard calls must run on the main thread.** With the privacy features on, the clipboard API talks to a system UI service through the run loop; calling it from a plain thread crashes. The service locks the main thread at startup and dedicates it to the clipboard.
- **"Transfer only when the other side pastes" isn't possible on macOS.** The system resolves promised clipboard data about 0.1 s after it's written and caches it, and there is no notification that the clipboard was read. That's why small items are pushed ahead of time and large files are pulled on demand.
- **Three pitfalls in the pion WebRTC library:** `DetachDataChannels()` is a global switch that breaks regular send/receive on every channel; `OnMessage` doesn't buffer, so messages arriving before the handler is set are dropped (202 bytes sent, 0 received); `OnBufferedAmountLow` is edge-triggered, so a waiter that arrives late never wakes up (a 20 MB transfer stalled at 786 KB).
- **pion opens a separate port per STUN server**, so the answers don't apply to each other. Adding STUN servers can't fix multi-uplink networks; all connections now share one port and CopySync probes it itself.
- **One side can't tell on its own whether a connection is relayed.** When one side sends through TURN, the other may just see an ordinary address; once one Mac showed "direct" and the other "relay". Both sides now tell each other after connecting.
- **pion never switches paths once it has picked one.** The relay handshake is fast, and if it wins that race the whole connection stays relayed. We saw exactly that after the peer restarted: its addresses arrived after the relay wait had expired. The relay wait is now 5 seconds, probe hostnames resolve in parallel (a probe round went from about 1.6 s to under 0.2 s), and relayed connections retry the direct path with an ICE restart.
- **Who sends first matters.** Some routers (the office one we measured) keep a connection entry for a packet that arrives before you've sent anything to its source, locking the port; your own packets then leave from a different port and punching fails. That's why restarting the office Mac always connected directly while restarting the home Mac fell back to the relay: a freshly started Mac probes its uplinks first, so its addresses went out late and it ended up sending first. A rule now decides who sends first, and retries alternate the order.
- **A dead connection has to be cleaned up by hand.** When pion reports "failed" it doesn't close the connection; the control channel still looks open and whatever is written to it is lost, and the layer above thinks it's connected and never reconnects. In 1.1.0 this sent six files into a dead connection after the other Mac dropped off, and none arrived.
- **Public STUN domains can be poisoned by DNS**, resolving to `192.0.2.42` in our tests. Reserved ranges and proxy fake-IP ranges are filtered before probing.
- **Simulated routers need a firewall.** When both sides punch at once, a packet let through to the router itself leaves a Linux conntrack entry, the outgoing flow then gets a new port, and punching fails. Real routers drop such packets, and the lab does too.

Measured on a LAN with a direct connection: a 20 MB file transfers in about 330 ms with matching checksums.

## Limitations

- **Not notarized by Apple, and the Windows build isn't code-signed**, so the first launch needs one manual approval.
- **No "virtual files" on Windows.** Files copied straight out of an Outlook attachment or a zip have no real path on the clipboard and aren't synced. Save or extract them first.
- **The Mac app doesn't skip password-manager content yet**; the Windows app does.
- **Items aren't resent after a long absence.** The server stores nothing. While the other Mac is unreachable, this Mac keeps the latest 20 items (from the last 24 hours) in memory, and they're lost if its background service restarts.
- **Relayed transfers are limited by your server's bandwidth.** Direct connections aren't.

## Build from source

The Mac app needs Go 1.26, Flutter 3.44 and Xcode:

```bash
./scripts/build.sh       # build everything into dist/
./scripts/package.sh     # build the DMG and server packages into release/
```

The Windows app needs Go, Flutter, Visual Studio ("Desktop development with C++") and Inno Setup 6. On Windows:

```powershell
powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1   # installer and portable zip into dist\windows\
```

The background service is pure Go and also cross-compiles from a Mac: `GOOS=windows go build ./cmd/copysyncd`. CI builds the Windows installer on every commit.

During development:

```bash
cd client-core && go test ./...          # service tests
./tools/natlab/docker.sh                 # NAT lab: direct vs relay across 11 network topologies
cd ui && flutter test                    # app tests
./scripts/install-macos.sh               # install the service from dist/ as a login item
./dist/copysync-cli status               # service status
./dist/copysync-cli watch                # live event stream
```

After changing a `.proto` file, run `./scripts/gen-proto.sh` (uses buf; install with `go install`, no protoc needed). After changing the UI, run `./tools/screenshots/render.sh` to regenerate the screenshots from demo data.

```
client-core/    background service (Go)
  internal/clipboard    clipboard (Mac: cgo + Objective-C; Windows: Win32; main thread only)
  internal/transport    signaling (WebSocket) and P2P (WebRTC)
  internal/sync         sync engine: routing, packing, transfer, storage
ui/             app (Flutter)
server/         signaling + TURN relay server (Go); deployment files in server/deploy/
proto/          message definitions and device ID derivation, shared by all three
docs/           website (GitHub Pages)
spikes/         technical validation done before development
design/         design documents
tools/natlab/   NAT traversal lab
tools/installer/ Windows installer (Inno Setup)
```

## Troubleshooting

**`.uuremote_…` file records keep appearing while UU Remote is running**
Upgrade both clients to 1.3.2. UU Remote puts temporary placeholder files on the system clipboard; older CopySync clients recorded and forwarded them as ordinary files. The new version filters the confirmed `.uuremote_aeawv` followed by digits before scanning, recording or sending, keeps ordinary files in mixed selections, ignores placeholder-only incoming offers, and removes placeholder-only history on startup without deleting UU's source files. Ordinary empty, hidden and large files are unaffected. The 50 MB threshold controls automatic transfer; there is no need to disable UU's UDP P2P. Other clipboard sync tools can still duplicate ordinary content; use one clipboard sync provider if needed. See the [investigation and validation notes](design/uu-remote-clipboard.md).

**Nothing comes off the clipboard, and `pbpaste` is empty too**
The system pasteboard service may be stuck: an app declared clipboard content and then crashed, leaving types without data. Run `killall pboard` and copy again.

**A device stays offline**
Check **Settings → Connection** on both computers. If it isn't connected, verify the server address and that port 8787 is reachable.

**Always "relay", never "direct"**
Showing "relay" right after connecting and "direct" a little later is normal: the relay makes things work first, and the direct path takes over when the link is idle. If it stays "relay", NAT traversal didn't succeed, which is common with NATs that randomize ports or strict firewalls. Everything still works; speed is limited by the server's bandwidth. Run `copysync-cli nat` to see how many uplinks were found; zero means UDP is blocked.

**Logs**
The background service logs to `~/Library/Logs/CopySync/daemon.log` on Mac and `%LOCALAPPDATA%\CopySync\logs\daemon.log` on Windows. On the server, use `journalctl -u copysync-server -f`.

## License

CopySync is dual-licensed. It's free under the [GNU AGPL-3.0](LICENSE): anyone, companies included, can use and modify it, but distributing it or offering a modified version over a network means publishing your source under the same license. A [commercial license](LICENSING.md) is available for closed-source products, private modifications and hosted services; [open a "Commercial license" issue](https://github.com/baiyuze/copysync/issues/new?template=commercial-license.yml) to ask. Versions 1.2.0 and earlier remain under the MIT License. Details are in [LICENSING.md](LICENSING.md).

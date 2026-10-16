<p align="center">
  <img src="docs/assets/icon.png" width="112" height="112" alt="PasteFerry icon">
</p>

<h1 align="center">PasteFerry</h1>

<p align="center">
  Copy on one computer, paste on the other. Mac, Windows and Linux.<br>
  Text, images, files and whole folders — across accounts and across networks.
</p>

<p align="center">
  <a href="https://pasteferry.cc/en/download/"><strong>Download</strong></a>
  ·
  <a href="https://pasteferry.cc/en/">Website</a>
  ·
  <a href="README.zh-CN.md">简体中文</a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/baiyuze/pasteferry?color=0A66D8" alt="Release">
  <img src="https://img.shields.io/badge/macOS-13%2B-1D1D1F" alt="macOS 13+">
  <img src="https://img.shields.io/badge/Windows-10%20%2F%2011-1D1D1F" alt="Windows 10 / 11">
  <img src="https://img.shields.io/badge/Linux-x64%20%2F%20ARM64-1D1D1F" alt="Linux x64 / ARM64">
  <img src="https://img.shields.io/badge/license-AGPL--3.0-6E6E73" alt="AGPL-3.0 License">
</p>

<p align="center">
  <img src="docs/assets/screenshots/en/history.png" alt="PasteFerry's history view: files, images and text arriving from other Macs">
</p>

## Why

Universal Clipboard needs the same Apple ID and both Macs close by. AirDrop makes you accept every time and doesn't touch the clipboard. Messaging yourself means manual downloads and your files on someone else's server. With one Mac and one Windows PC, none of the built-in options work.

With PasteFerry you pair two computers once. After that, copy on one and paste on the other — Mac to Mac, Mac to Windows, or Windows to Windows. On the same network it works without a server. Across networks, enable **Settings → Use official service**: no server to host and no address to enter. Direct connections are tried first; if they fail, official relay forwards end-to-end encrypted packets. Custom servers and Tailscale remain available.

| | PasteFerry | Universal Clipboard | AirDrop | Chat apps |
|---|:-:|:-:|:-:|:-:|
| Different Apple IDs | ✓ | — | Needs accepting | ✓ |
| Different networks (home ↔ office) | ✓ | — | — | ✓ |
| Paste right away, nothing to accept | ✓ | ✓ | — | — |
| Between Mac and Windows | ✓ | — | — | ✓ |
| Files and whole folders | ✓ | Unreliable | ✓ | Zipped |
| History of what you copied | ✓ | — | — | Scroll back |
| Only your devices can decrypt content | ✓ | ✓ | ✓ | — |

## Features

- **Text, rich text, images, files and folders.** What you copy is what comes out on the other side.
- **Size-aware transfer.** Anything under 50 MB (adjustable) is pushed as soon as you copy; larger files sync as a record you pull when you need them.
- **History** of the last 3 days (adjustable) across all your computers, with image preview.
- **No server on the same network.** Computers find each other and pair with a code; on a different subnet you also enter an address.
- **Direct first, relay as fallback**, including office networks with several uplinks.
- **Fingerprint check when pairing**, so neither the server nor anyone on the network can impersonate your device.
- **Runs in the background**, starts at login, recovers from crashes. English, Simplified and Traditional Chinese, and Japanese.

<table>
  <tr>
    <td><img src="docs/assets/screenshots/en/devices.png" alt="Devices: paired Macs and whether each is connected directly or through the relay"></td>
    <td><img src="docs/assets/screenshots/en/verify.png" alt="Pairing: compare the fingerprint shown on both Macs before confirming"></td>
  </tr>
</table>

## Get started

New here? The **[quick start](https://pasteferry.cc/en/start/)** walks through it with screenshots. You need **the app on each computer**. On the same network that's all; across networks, enable official access or use your own server or Tailscale.

**1. Install the app** on each computer from the **[download page](https://pasteferry.cc/en/download/)**:

- **Mac** (macOS 13+): [Apple silicon](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-mac-apple-silicon.dmg) or [Intel](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-mac-intel.dmg); a [universal](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry.dmg) build runs on both.
- **Windows** 10/11: [installer](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-Setup.exe) or [portable zip](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-windows-x64.zip).
- **Linux** (X11, or Wayland with XWayland): `.deb` for Ubuntu and Debian ([x64](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-linux-amd64.deb), [ARM64](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-linux-arm64.deb)), or `tar.gz` for other distributions ([x64](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-linux-amd64.tar.gz), [ARM64](https://github.com/baiyuze/pasteferry/releases/latest/download/PasteFerry-linux-arm64.tar.gz)) — extract it and run `./PasteFerry/install.sh` to add it to the app menu.

The builds aren't notarized or code-signed yet: on Mac, allow it once under **System Settings → Privacy & Security → Open Anyway**; on Windows, click **More info → Run anyway**. Then click **Turn on background sync**, and allow PasteFerry when Windows Firewall asks, or when macOS asks to find devices on your local network.

**2. Pair.** On one computer, **Devices → Add device → Show a code**; on the other, **Enter a code**. On a different subnet, click **On a different subnet?** and also enter the address shown under the code. Confirm only if both screens show the same two fingerprint lines. From then on, copy on either one and paste on the other.

**3. Across networks (optional).** Enable **Use official service** on both computers before pairing; existing pairings are kept. New installs keep this off until you choose it. See the [official service guide](https://pasteferry.cc/en/official/). Other connection choices remain available:

| Way | Best for | You need |
|---|---|---|
| **[Official service](https://pasteferry.cc/en/official/)** | Different networks without hosting a server | Enable one switch on each computer; current default relay allowance is 100 MB/day per machine quota group, upload plus download; direct traffic does not consume it |
| **[Same network](https://pasteferry.cc/en/lan/)** | Computers behind one router, at home or in the office | Nothing |
| **[Public server](https://pasteferry.cc/en/cloud/)** | Computers on different networks that should sync from anywhere | A Linux VM with a public IP; the smallest is enough |
| **[Tailscale](https://pasteferry.cc/en/tailscale/)** | Across networks without a server | [Tailscale](https://tailscale.com) on every computer; enter the other's Tailscale address when pairing |

To self-host instead (optional):

```bash
curl -LO https://github.com/baiyuze/pasteferry/releases/latest/download/pasteferry-server-linux-amd64.tar.gz
tar xzf pasteferry-server-linux-amd64.tar.gz && cd pasteferry-server-linux-amd64
sudo ./install.sh <public IP of the server>
```

Open TCP 8787, UDP 3478 and UDP 32768–60999, then enter `ws://<server address>:8787/signal` under **Settings → Custom server**, then click **Save and connect** on each computer. Docker and all options: [server deployment guide](server/deploy/README.md).

| | Requirement |
|---|---|
| **Mac app** | macOS 13 Ventura or later; Apple silicon and Intel |
| **Windows app** | Windows 10 and 11, x64 |
| **Linux app** | x64 or ARM64; an X11 desktop, or Wayland with XWayland (GNOME and KDE have it by default); GTK 3 |
| **Server** | Optional for self-hosting; official access needs no deployment: Linux (x86_64, ARM64; systemd or Docker), macOS, or Windows 10/11 x64 |
| **iPhone / iPad** | Not supported: iOS doesn't let apps watch the clipboard in the background |

## Official access and relay allowance (2.1)

New installations keep official public service off. Enable **Settings → Use
official service** to try direct connections first and request official relay
only after direct failure. A separate sidebar shows the daily allowance and
usage, shown in MB with two decimal places. The current default is 100 MB per
machine quota group per day, upload plus download, resetting at midnight Beijing
time. Direct traffic does not consume this allowance; direct connections remain
available when relay allowance is exhausted. Actual limits and availability can
change and are shown in the app. Both peers need a compatible client for official relay; LAN and custom
servers continue to work without the gateway. See [privacy](PRIVACY.md),
[security](SECURITY.md) and the [public protocol](design/official-gateway.md).

The official website is **[pasteferry.cc](https://pasteferry.cc/)**. Old GitHub
Pages links redirect to it.

## Documentation

- [How it works](guide/how-it-works.md): architecture, NAT traversal, security model, fallbacks, lessons learned
- [Troubleshooting](guide/troubleshooting.md) and known limitations
- [Official service](https://pasteferry.cc/en/official/) · [Quick start](https://pasteferry.cc/en/start/) · [Same network](https://pasteferry.cc/en/lan/) · [Public server](https://pasteferry.cc/en/cloud/) · [Tailscale](https://pasteferry.cc/en/tailscale/) · [Run a server on your LAN](https://pasteferry.cc/en/server/), all with screenshots
- [Server deployment](server/deploy/README.md): Windows, macOS, Linux, Docker, Tailscale
- [Building from source](guide/building.md)
- Design notes (Chinese): [Local network mode](design/lan-mode.md), [NAT traversal on multi-uplink networks](design/nat-traversal.md), [Windows client](design/windows-client.md)
- [Changelog](CHANGELOG.md)

## License

PasteFerry is free and open source under the [GNU AGPL-3.0](LICENSE): anyone, companies included, can use and modify it, but distributing it or offering a modified version over a network means publishing your source under the same license. See [LICENSING.md](LICENSING.md) for details. Versions 1.2.0 and earlier remain under the MIT License.

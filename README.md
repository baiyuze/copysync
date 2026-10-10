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

Universal Clipboard needs the same Apple ID and both Macs close by. AirDrop makes you accept every time and doesn't touch the clipboard. Messaging yourself means manual downloads and your files on someone else's server. With one Mac and one Windows PC, none of the built-in options work.

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
- **Size-aware transfer.** Anything under 50 MB (adjustable) is pushed as soon as you copy; larger files sync as a record you pull when you need them.
- **History** of the last 3 days (adjustable) across all your computers, with image preview.
- **Direct first, relay as fallback**, including office networks with several uplinks.
- **Fingerprint check when pairing**, so neither the server nor anyone on the network can impersonate your device.
- **Runs in the background**, starts at login, recovers from crashes. English, Chinese and Japanese.

<table>
  <tr>
    <td><img src="docs/assets/screenshots/en/devices.png" alt="Devices: paired Macs and whether each is connected directly or through the relay"></td>
    <td><img src="docs/assets/screenshots/en/verify.png" alt="Pairing: compare the fingerprint shown on both Macs before confirming"></td>
  </tr>
</table>

## Get started

You need **the app on each computer** and **a server both can reach**.

**1. Run the server.**

- *Same LAN* (home or office): run it on any always-on Windows, Mac or Linux computer. Download the package, run one install script, done. Step-by-step guide: **[Run the server on your LAN](https://baiyuze.github.io/copysync/en/server/)**.
- *Across the internet*: a Linux machine with a public IP, as root:

  ```bash
  curl -LO https://github.com/baiyuze/copysync/releases/latest/download/copysync-server-linux-amd64.tar.gz
  tar xzf copysync-server-linux-amd64.tar.gz && cd copysync-server-linux-amd64
  sudo ./install.sh <public IP of the server>
  ```

  Open TCP 8787, UDP 3478 and UDP 32768–60999. Docker and all options: [server deployment guide](server/deploy/README.md).

**2. Install the app** on each computer: [CopySync.dmg](https://github.com/baiyuze/copysync/releases/latest/download/CopySync.dmg) for macOS 13+, [CopySync-Setup.exe](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-Setup.exe) for Windows 10/11 ([portable zip](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-windows-x64.zip)). The builds aren't notarized or code-signed yet: on Mac, allow it once under **System Settings → Privacy & Security → Open Anyway**; on Windows, click **More info → Run anyway**. Then click **Turn on background sync** and enter `ws://<server address>:8787/signal` under **Settings → Signaling server**.

**3. Pair.** On one computer, **Devices → Add device → Show a code**; on the other, **Enter a code**. Confirm only if both screens show the same two fingerprint lines. From then on, copy on either one and paste on the other.

| | Requirement |
|---|---|
| **Mac app** | macOS 13 Ventura or later; Intel and Apple silicon |
| **Windows app** | Windows 10 (21H2 or later) and 11, x64 |
| **Server** | Linux (x86_64, ARM64; systemd or Docker), macOS, or Windows 10/11 x64 |
| **iPhone / iPad** | Not supported: iOS doesn't let apps watch the clipboard in the background |

## Documentation

- [How it works](guide/how-it-works.md): architecture, NAT traversal, security model, fallbacks, lessons learned
- [Troubleshooting](guide/troubleshooting.md) and known limitations
- [Server deployment](server/deploy/README.md) · [LAN guide with screenshots](https://baiyuze.github.io/copysync/en/server/)
- [Building from source](guide/building.md)
- Design notes (Chinese): [NAT traversal on multi-uplink networks](design/nat-traversal.md), [Windows client](design/windows-client.md)
- [Changelog](CHANGELOG.md)

## License

CopySync is dual-licensed. It's free under the [GNU AGPL-3.0](LICENSE): anyone, companies included, can use and modify it, but distributing it or offering a modified version over a network means publishing your source under the same license. A [commercial license](LICENSING.md) is available for closed-source products, private modifications and hosted services; [open a "Commercial license" issue](https://github.com/baiyuze/copysync/issues/new?template=commercial-license.yml) to ask. Versions 1.2.0 and earlier remain under the MIT License.

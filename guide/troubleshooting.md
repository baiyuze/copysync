# Troubleshooting

[简体中文](troubleshooting.zh-CN.md) · [Back to README](../README.md)

**The other computer can't be found on the same network**
Make sure both have background sync on and are on the same network. Wired and Wi‑Fi in many offices, guest networks and Wi‑Fi with client isolation don't let devices see each other: on the computer entering the code, click **On a different subnet?** and enter the address shown under the other computer's code. Still nothing? Check the other computer's firewall: on Windows, click **Allow** under Settings → Local network (one administrator prompt); on a Mac, turn on CopySync Daemon in System Settings → Privacy & Security → Local Network. The Local network part of `copysync-cli nat` shows this computer's port, how many CopySync devices it found, whether each paired device is reachable and which addresses it remembers. See the [same network guide](https://baiyuze.github.io/copysync/en/lan/).

**Settings → Connection stays at Not connected**
This only applies when a server address is set; on the same network you can leave it empty. Check that the address starts with `ws://`, uses port 8787 and ends in `/signal`, then open `http://<server address>:8787/healthz` in a browser on that computer. If it doesn't load, the network path is broken: the server isn't running, its computer is asleep, a firewall blocks it, or its address changed. The Windows and Linux install scripts print the firewall commands you need; see [Server deployment](../server/deploy/README.md) and the [LAN guide](https://baiyuze.github.io/copysync/en/server/).

**A cloud server in mainland China works by IP but not by domain**
Providers in mainland China intercept plain HTTP requests for domains without an ICP filing, on any port, and answer with a redirect to a block page. The WebSocket handshake is a plain HTTP request, so it's intercepted too. Enter the IP, as in `ws://<server IP>:8787/signal`, or file the domain. See the [public server guide](https://baiyuze.github.io/copysync/en/cloud/).

**Using Tailscale: devices show Direct, but transfers are slow**
Direct means the two computers connect directly on the Tailscale network. When Tailscale can't connect them peer to peer, it relays through its DERP servers, which can be far away. Run `tailscale ping <other device>`: "via DERP" means Tailscale is relaying, "via <IP>:<port>" means peer to peer. See the [Tailscale guide](https://baiyuze.github.io/copysync/en/tailscale/).

**`.uuremote_…` file records keep appearing while UU Remote is running**
Upgrade both clients to 1.3.2 or later. UU Remote puts temporary placeholder files on the system clipboard; older CopySync clients recorded and forwarded them as ordinary files. The new version filters the confirmed `.uuremote_aeawv` followed by digits before scanning, recording or sending, keeps ordinary files in mixed selections, ignores placeholder-only incoming offers, and removes placeholder-only history on startup without deleting UU's source files. Ordinary empty, hidden and large files are unaffected. The 50 MB threshold controls automatic transfer; there is no need to disable UU's UDP P2P. Other clipboard sync tools can still duplicate ordinary content; use one clipboard sync provider if needed. See the [investigation and validation notes](../design/uu-remote-clipboard.md).

**Nothing comes off the clipboard, and `pbpaste` is empty too**
The system pasteboard service may be stuck: an app declared clipboard content and then crashed, leaving types without data. Run `killall pboard` and copy again.

**A device stays offline**
On the same network, see "The other computer can't be found" above; if both computers changed addresses, open the device's … menu and choose **Enter its address…**. Across networks, check **Settings → Connection** on both computers; if it isn't connected, verify the server address and that port 8787 is reachable.

**Always "relay", never "direct"**
Showing "relay" right after connecting and "direct" a little later is normal: the relay makes things work first, and the direct path takes over when the link is idle. If it stays "relay", NAT traversal didn't succeed, which is common with NATs that randomize ports or strict firewalls. Everything still works; speed is limited by the server's bandwidth. Run `copysync-cli nat` to see how many uplinks were found; zero means UDP is blocked.

**Logs**
The background service logs to `~/Library/Logs/CopySync/daemon.log` on Mac and `%LOCALAPPDATA%\CopySync\logs\daemon.log` on Windows. Server logs: `journalctl -u copysync-server -f` on Linux, `~/Library/Logs/CopySync/server.log` on Mac, `C:\ProgramData\CopySync Server\server.log` on Windows.

## Limitations

- **Not notarized by Apple, and the Windows build isn't code-signed**, so the first launch needs one manual approval.
- **No "virtual files" on Windows.** Files copied straight out of an Outlook attachment or a zip have no real path on the clipboard and aren't synced. Save or extract them first.
- **The Mac app doesn't skip password-manager content yet**; the Windows app does.
- **Items aren't resent after a long absence.** The server stores nothing. While the other Mac is unreachable, this Mac keeps the latest 20 items (from the last 24 hours) in memory, and they're lost if its background service restarts.
- **Relayed transfers are limited by your server's bandwidth.** Direct connections aren't.
- **On different subnets, if both computers change addresses**, you enter the address once more in the device menu; CopySync doesn't scan subnets. Two Windows PCs that both block inbound traffic, where only one knows the other's address, can't connect until the one being reached clicks **Allow**.
- **macOS 15 needs the Local Network permission**, or the background service can't reach devices on your network (it still connects when the other side reaches out first). Current builds are ad-hoc signed, so macOS may ask again after an update.

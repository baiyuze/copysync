# Server deployment / CopySync 服务器部署

The server only does two things: it lets your devices find each other (signaling), and relays encrypted traffic when a direct connection can't be established (TURN). It keeps no accounts and cannot read what you copy.

CopySync 的服务器只做两件事：让你的设备找到彼此（信令），以及在设备之间打不通直连时转发加密后的数据（TURN 中转）。它不存账号，也看不到你复制的内容。

## Requirements / 需要准备

- A machine both devices can reach. On a LAN, one of your Macs will do.
- Across the internet it needs a public IP, with these ports open:

一台两端设备都能访问到的机器，同一局域网内用其中一台 Mac 也行；跨公网使用时需要公网 IP，并放行以下端口：

| Port 端口 | Protocol 协议 | Purpose 用途 |
|---|---|---|
| 8787 | TCP | Signaling (WebSocket) 信令 |
| 3478 | UDP | STUN / TURN |
| 32768–60999 | UDP | TURN relay ports, assigned by the OS / 中转端口，系统随机分配 |

Uses about 15 MB of RAM; the smallest VPS is enough. 内存占用约 15 MB，1 核 512 MB 的机器足够。

## Option 1: systemd (recommended) / 方式一：systemd（推荐）

```bash
tar xzf copysync-server-linux-amd64.tar.gz
cd copysync-server-linux-amd64
sudo ./install.sh <public IP / 公网IP>
```

The script installs the binary to `/usr/local/bin`, writes `/etc/copysync/server.env` with a random TURN secret, and registers a systemd service that starts on boot. Running it again upgrades in place and keeps your config.

脚本会把程序装到 `/usr/local/bin`，生成配置 `/etc/copysync/server.env`（含随机的 TURN 密钥），注册为开机自启的 systemd 服务。重复执行即为升级，已有配置不会被覆盖。

```bash
journalctl -u copysync-server -f        # logs / 日志
systemctl restart copysync-server       # restart after editing config / 改完配置后重启
copysync-server -version                # version / 版本
```

## Option 2: Docker / 方式二：Docker

```bash
cp env.example .env       # set TURN_IP and TURN_SECRET / 填入 TURN_IP 与 TURN_SECRET
docker compose up -d
```

`network_mode: host` is required (already set in the compose file): TURN assigns a random UDP port per session, which port mapping can't expose.

必须使用 `network_mode: host`（compose 文件里已经写好）：TURN 为每个会话随机分配 UDP 端口，端口映射模式下这些端口不可达。

## Option 3: run it directly / 方式三：直接运行

```bash
./copysync-server -addr :8787 -turn-ip <public IP> -turn-secret <random string>
```

On a LAN you can omit `-turn-ip`; relaying won't be available. 只在局域网内使用、不需要中转时，`-turn-ip` 可以省略。

## Client setup / 客户端设置

On each computer (Mac or Windows), open CopySync → 设置 (Settings) → 信令服务器地址 (Signaling server) and enter:

在每台电脑（Mac 或 Windows）上打开 CopySync →「设置」→「信令服务器地址」，填入：

```
ws://<server address>:8787/signal
```

Behind an HTTPS reverse proxy, use `wss://` and make sure the proxy forwards WebSocket upgrades. 放在 HTTPS 反向代理后面时用 `wss://`，代理需要支持 WebSocket 升级。

## Uninstall / 卸载

```bash
sudo systemctl disable --now copysync-server
sudo rm /usr/local/bin/copysync-server /etc/systemd/system/copysync-server.service
sudo rm -r /etc/copysync
```

# CopySync 服务器部署 / Server deployment

CopySync 的服务器只做两件事：让你的设备找到彼此（信令），以及在设备之间打不通直连时转发加密后的数据（TURN 中转）。它不存账号，也看不到你复制的内容。

The server only does two things: it lets your devices find each other (signaling), and relays encrypted traffic when a direct connection can't be established (TURN). It keeps no accounts and cannot read what you copy.

## 需要准备 / Requirements

- 一台两端设备都能访问到的机器。同一局域网内的话，用其中一台 Mac 也行。
- 跨公网使用时需要公网 IP，并放行以下端口：

| 端口 Port | 协议 | 用途 Purpose |
|---|---|---|
| 8787 | TCP | 信令 Signaling (WebSocket) |
| 3478 | UDP | STUN / TURN |
| 32768–60999 | UDP | TURN 中转端口，系统随机分配 / relay ports, assigned by the OS |

内存占用约 15 MB，1 核 512 MB 的机器足够。 Uses about 15 MB of RAM; the smallest VPS is enough.

## 方式一：systemd（推荐） / Option 1: systemd (recommended)

```bash
tar xzf copysync-server-linux-amd64.tar.gz
cd copysync-server-linux-amd64
sudo ./install.sh <公网IP / public IP>
```

脚本会把程序装到 `/usr/local/bin`，生成配置 `/etc/copysync/server.env`（含随机的 TURN 密钥），注册为开机自启的 systemd 服务。重复执行即为升级，已有配置不会被覆盖。

The script installs the binary to `/usr/local/bin`, writes `/etc/copysync/server.env` with a random TURN secret, and registers a systemd service that starts on boot. Running it again upgrades in place and keeps your config.

```bash
journalctl -u copysync-server -f        # 日志 / logs
systemctl restart copysync-server       # 改完配置后重启 / restart after editing config
copysync-server -version                # 版本 / version
```

## 方式二：Docker / Option 2: Docker

```bash
cp env.example .env       # 填入 TURN_IP 与 TURN_SECRET / set TURN_IP and TURN_SECRET
docker compose up -d
```

必须使用 `network_mode: host`（compose 文件里已经写好）：TURN 为每个会话随机分配 UDP 端口，端口映射模式下这些端口不可达。

`network_mode: host` is required (already set in the compose file): TURN assigns a random UDP port per session, which port mapping can't expose.

## 方式三：直接运行 / Option 3: run it directly

```bash
./copysync-server -addr :8787 -turn-ip <公网IP> -turn-secret <随机字符串>
```

只在局域网内使用、不需要中转时，`-turn-ip` 可以省略。 On a LAN you can omit `-turn-ip`; relaying won't be available.

## 客户端设置 / Client setup

在每台 Mac 上打开 CopySync →「设置」→「信令服务器地址」，填入：

```
ws://<服务器地址>:8787/signal
```

放在 HTTPS 反向代理后面时用 `wss://`，代理需要支持 WebSocket 升级。 Behind an HTTPS reverse proxy, use `wss://` and make sure the proxy forwards WebSocket upgrades.

## 卸载 / Uninstall

```bash
sudo systemctl disable --now copysync-server
sudo rm /usr/local/bin/copysync-server /etc/systemd/system/copysync-server.service
sudo rm -r /etc/copysync
```

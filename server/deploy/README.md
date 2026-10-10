# Server deployment / CopySync 服务器部署

The server only does two things: it lets your devices find each other (signaling), and relays encrypted traffic when a direct connection can't be established (TURN). It keeps no accounts and cannot read what you copy.

CopySync 的服务器只做两件事：让你的设备找到彼此（信令），以及在设备之间打不通直连时转发加密后的数据（TURN 中转）。它不存账号，也看不到你复制的内容。

Step-by-step guide with screenshots / 带截图的图文教程：
<https://baiyuze.github.io/copysync/server/>

## Where to run it / 放在哪里

| | Requirement 要求 |
|---|---|
| **On your LAN** 局域网 | Any always-on Windows 10/11, macOS 13+ or Linux computer on the same network, one of your own computers included. Reserve its address in the router. 同一网络里任意一台常开的 Windows 10/11、macOS 13 以上或 Linux 电脑，自己的电脑也行；在路由器里给它固定地址。 |
| **Across the internet** 跨公网 | A Linux machine with a public IP; the smallest VPS is enough. 一台有公网 IP 的 Linux 机器，最低配的云主机就够。 |

It uses about 15 MB of RAM. 内存占用约 15 MB。

| Port 端口 | Protocol 协议 | Purpose 用途 |
|---|---|---|
| 8787 | TCP | Signaling (WebSocket) 信令 |
| 3478 | UDP | STUN / TURN |
| 32768–60999 (Linux), 49152–65535 (Windows, macOS) | UDP | TURN relay ports, assigned by the OS / 中转端口，系统随机分配 |

The install scripts below find the computer's LAN address, enable the TURN relay on it, register the server to start automatically, and print the address to enter in the apps. Running a script again upgrades in place and keeps your settings.

下面的安装脚本会找到本机的局域网地址并在它上面启用 TURN 中转、把服务器设为自动运行，最后打印出客户端要填的地址。重复运行即为升级，原来的设置保留。

## Windows

Download `copysync-server-windows-amd64.zip`, extract it, and double-click `install.cmd` (it asks for administrator rights).

下载 `copysync-server-windows-amd64.zip`，解压后双击 `install.cmd`（会请求管理员权限）。

- Runs as the Windows service `CopySyncServer` under the low-privilege LOCAL SERVICE account, starts with Windows and restarts after failures. 以低权限的 LOCAL SERVICE 账户注册为 Windows 服务 `CopySyncServer`，开机启动，出错后自动重启。
- Adds a Windows Firewall rule for `copysync-server.exe` only. 在 Windows 防火墙里只为 `copysync-server.exe` 添加一条入站规则。
- Wrong address picked? Run `install.cmd 192.168.1.20` with yours. 认错了地址就运行 `install.cmd 192.168.1.20`，换成你的地址。

| | |
|---|---|
| Program 程序 | `C:\Program Files\CopySync Server` |
| Settings 设置 | `C:\ProgramData\CopySync Server\server.conf` (administrators only 仅管理员可读) |
| Logs 日志 | `C:\ProgramData\CopySync Server\server.log` |
| Restart 重启 | `Restart-Service CopySyncServer` (administrator PowerShell 管理员 PowerShell) |
| Uninstall 卸载 | double-click `uninstall.cmd` 双击 `uninstall.cmd` |

## macOS

```bash
tar xzf copysync-server-macos.tar.gz
cd copysync-server-macos
./install.sh                 # or ./install.sh 192.168.1.20
```

No administrator rights needed: it registers a LaunchAgent that runs whenever you're logged in. Keep that Mac awake. If the macOS firewall is on, the script prints the command to allow the server through it.

不需要管理员权限：注册为 LaunchAgent，登录后自动运行。这台 Mac 不能睡眠。开着系统防火墙时，脚本最后会给出放行命令。

| | |
|---|---|
| Logs 日志 | `~/Library/Logs/CopySync/server.log` |
| Restart 重启 | `launchctl kickstart -k gui/$(id -u)/com.copysync.server` |
| Uninstall 卸载 | `./install.sh uninstall` |

## Linux: systemd (recommended) / Linux：systemd（推荐）

```bash
tar xzf copysync-server-linux-amd64.tar.gz
cd copysync-server-linux-amd64
sudo ./install.sh                  # on a LAN: uses this machine's LAN address / 局域网：用本机的局域网地址
sudo ./install.sh <public IP>      # on a cloud server / 云服务器：填公网 IP
```

The script installs the binary to `/usr/local/bin`, writes `/etc/copysync/server.env` with a random TURN secret, and registers a systemd service that starts on boot. If ufw or firewalld is active, it prints the commands to open the ports.

脚本会把程序装到 `/usr/local/bin`，生成配置 `/etc/copysync/server.env`（含随机的 TURN 密钥），注册为开机自启的 systemd 服务。开着 ufw 或 firewalld 时，最后会给出放行端口的命令。

```bash
journalctl -u copysync-server -f        # logs / 日志
systemctl restart copysync-server       # restart after editing the config / 改完配置后重启
copysync-server -version                # version / 版本
sudo ./install.sh uninstall             # uninstall / 卸载
```

On a cloud server, also open the ports above in the provider's security group. 云服务器还要在安全组里放行上面的端口。

## Linux: Docker / Linux：Docker

```bash
cp env.example .env       # set TURN_IP and TURN_SECRET / 填入 TURN_IP 与 TURN_SECRET
docker compose up -d
```

`network_mode: host` is required (already set in the compose file): TURN assigns a random UDP port per session, which port mapping can't expose. For the same reason Docker Desktop on Mac or Windows won't work; use the install scripts there.

必须使用 `network_mode: host`（compose 文件里已经写好）：TURN 为每个会话随机分配 UDP 端口，端口映射模式下这些端口不可达。同样的原因，Mac 与 Windows 上的 Docker Desktop 不适用，请用安装脚本。

## Run it directly / 直接运行

```bash
./copysync-server -addr :8787 -turn-ip <address> -turn-secret <random string>
```

`-turn-ip` is the address clients use to reach the relay: the LAN address on a LAN, the public IP on a cloud server. Without it there is no relay. `-log <file>` writes the log to a file instead of standard error. Run `copysync-server -h` for all options.

`-turn-ip` 是客户端访问中转用的地址：局域网里填局域网地址，云服务器填公网 IP；不填就没有中转。`-log <文件>` 把日志写到文件而不是标准错误。全部参数见 `copysync-server -h`。

## Client setup / 客户端设置

On each computer (Mac or Windows), open CopySync → Settings → Signaling server and enter the address the script printed:

在每台电脑（Mac 或 Windows）上打开 CopySync →「设置」→「信令服务器地址」，填入脚本打印的地址：

```
ws://<server address>:8787/signal
```

Then check `http://<server address>:8787/healthz` in a browser on another computer; it should show `{"status":"ok",...}`.

然后在另一台电脑的浏览器里打开 `http://<服务器地址>:8787/healthz`，应当显示 `{"status":"ok",...}`。

Behind an HTTPS reverse proxy, use `wss://` and make sure the proxy forwards WebSocket upgrades. 放在 HTTPS 反向代理后面时用 `wss://`，代理需要支持 WebSocket 升级。

On a cloud server in mainland China, enter the IP rather than a domain unless the domain has an ICP filing: providers intercept plain HTTP for unfiled domains on any port, and the WebSocket handshake is plain HTTP. 国内云服务器上，域名没有备案时请直接填 IP：云厂商会拦截未备案域名的明文 HTTP 请求，任何端口都一样，而 WebSocket 握手就是明文 HTTP。

<p align="center">
  <img src="docs/assets/icon.png" width="112" height="112" alt="CopySync 图标">
</p>

<h1 align="center">CopySync</h1>

<p align="center">
  在这台电脑复制，到那台电脑粘贴。Mac 与 Windows 互通。<br>
  文本、图片、文件、整个文件夹都可以，不限账号，不限是否在同一个网络。
</p>

<p align="center">
  <a href="https://github.com/baiyuze/copysync/releases/latest/download/CopySync.dmg"><strong>下载 Mac 版</strong></a>
  ·
  <a href="https://github.com/baiyuze/copysync/releases/latest/download/CopySync-Setup.exe"><strong>下载 Windows 版</strong></a>
  ·
  <a href="https://baiyuze.github.io/copysync/">项目主页</a>
  ·
  <a href="README.en.md">English</a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/baiyuze/copysync?label=%E7%89%88%E6%9C%AC&color=0A66D8" alt="版本">
  <img src="https://img.shields.io/badge/macOS-13%2B-1D1D1F" alt="macOS 13+">
  <img src="https://img.shields.io/badge/Windows-10%20%2F%2011-1D1D1F" alt="Windows 10 / 11">
  <img src="https://img.shields.io/badge/Intel%20%2B%20Apple%20%E8%8A%AF%E7%89%87-%E9%80%9A%E7%94%A8-1D1D1F" alt="Intel 与 Apple 芯片通用">
  <img src="https://img.shields.io/github/license/baiyuze/copysync?color=6E6E73" alt="MIT 许可">
</p>

<p align="center">
  <img src="docs/assets/screenshots/history.png" alt="CopySync 的复制记录界面：来自其他 Mac 的文件、图片和文字">
</p>

## 它解决什么问题

手边有两台电脑的人，每天都在做同一件事：这台上复制了一段文字、一张截图、一个文件，要拿到那台上用。

现有的办法各有各的不顺手：

- **系统自带的「通用剪贴板」**要求两台设备登录同一个 Apple ID、离得足够近、蓝牙和 Wi‑Fi 都开着。公司电脑和自己的电脑通常不是同一个 Apple ID；大文件经常复制到一半就没了。
- **隔空投送**每次都要选设备、在另一台上点接收，而且只传文件，不管剪贴板。
- **用聊天软件发给自己**要手动下载、图片会被压缩，文件还会留在第三方的服务器上。
- **一台 Mac、一台 Windows** 时，上面这些系统自带的办法都用不了。

CopySync 的做法是：两台电脑配对一次，之后在一台上复制，另一台上直接粘贴。Mac 与 Mac、Mac 与 Windows、Windows 与 Windows 都可以。数据在两台设备之间直接传输，打不通直连时经过你自己部署的服务器中转，全程端到端加密，服务器看不到内容。

| | CopySync | 通用剪贴板 | 隔空投送 | 聊天软件 |
|---|:-:|:-:|:-:|:-:|
| 不同 Apple ID 之间 | ✓ | — | 需对方接收 | ✓ |
| 不在同一个网络（家里 ↔ 公司） | ✓ | — | — | ✓ |
| 复制后直接粘贴，不用点接收 | ✓ | ✓ | — | — |
| Mac 与 Windows 之间 | ✓ | — | — | ✓ |
| 文件与整个文件夹 | ✓ | 不稳定 | ✓ | 打包后可以 |
| 复制历史，可以找回之前的内容 | ✓ | — | — | 翻聊天记录 |
| 内容不经过第三方服务器 | ✓ | ✓ | ✓ | — |

## 功能

- **文本、带格式的文本、图片、文件、文件夹**都能同步。复制什么，对面粘贴出来就是什么。
- **按大小分流**：50 MB 以内（可调）的内容复制时直接推送，对面立刻能粘贴；更大的文件只同步一条记录，需要时点「拉取到本机」，不会无谓地占用带宽和磁盘。
- **复制记录**：最近 3 天（可调）在所有设备上复制过的内容都在这里，随时可以放回剪贴板；图片点开就能预览，可以缩放。
- **直连优先，自动中转**：先尝试设备之间直连；在对称 NAT、公司防火墙后面打不通时，自动改走你服务器上的 TURN 中转。
- **配对时核对安全指纹**：防止服务器或网络上的第三方冒充你的设备。
- **后台常驻**：关掉窗口照常同步，开机自动启动，崩溃后自动恢复。Windows 上在任务栏右下角有托盘图标。
- **Mac 与 Windows 互通**：Mac 上 Intel 与 Apple 芯片原生运行；Windows 10、11 均可，粘贴文件时得到的是真实文件，能直接粘进资源管理器、微信、Office。支持浅色与深色外观。
- **不同步密码**（Windows）：密码管理器复制的内容带着「不要记录」的标记，CopySync 跳过它们，不同步也不进复制记录。

<table>
  <tr>
    <td><img src="docs/assets/screenshots/devices.png" alt="设备页：已配对的设备、直连或中转状态"></td>
    <td><img src="docs/assets/screenshots/verify.png" alt="配对时核对两台设备上显示的安全指纹"></td>
  </tr>
  <tr>
    <td><img src="docs/assets/screenshots/settings.png" alt="设置页：服务器地址、同步上限、同步的内容类型"></td>
    <td><img src="docs/assets/screenshots/history-dark.png" alt="深色外观下的复制记录"></td>
  </tr>
</table>

## 支持哪些电脑

| | 要求 |
|---|---|
| **Mac 客户端** | macOS 13 Ventura 及以上；Intel 与 Apple 芯片（M1–M4 等）通用，不需要 Rosetta |
| **Windows 客户端** | Windows 10（21H2 及以上）与 Windows 11，x64；ARM 版 Windows 可以模拟运行 |
| **服务器** | 任意 Linux（x86_64 或 ARM64），systemd 或 Docker 均可；也可以直接跑在其中一台 Mac 上 |
| **iPhone / iPad** | 暂不支持。iOS 不允许 App 在后台监听剪贴板，做不到 Mac 上这种复制即同步的体验 |

## 安装

整套东西由两部分组成：**每台电脑装一个 App**，再加**一台两端都能访问到的服务器**。两台电脑在同一个局域网时，服务器直接放在其中一台 Mac 上就行。

### 1. 部署服务器

在 Linux 服务器上（需要 root）：

```bash
curl -LO https://github.com/baiyuze/copysync/releases/latest/download/copysync-server-linux-amd64.tar.gz
tar xzf copysync-server-linux-amd64.tar.gz
cd copysync-server-linux-amd64
sudo ./install.sh <服务器公网IP>
```

脚本会注册一个开机自启的 systemd 服务，并在最后打印出客户端要填的地址。ARM 服务器下载 `linux-arm64` 版本。

云服务器的安全组或防火墙需要放行：

| 端口 | 协议 | 用途 |
|---|---|---|
| 8787 | TCP | 信令（WebSocket） |
| 3478 | UDP | STUN / TURN |
| 32768–60999 | UDP | TURN 中转端口，系统随机分配 |

用 Docker 部署、或者只在局域网里用，见[服务器部署说明](server/deploy/README.md)。

### 2. 在每台电脑上安装 App

**Mac**

1. 下载 [CopySync.dmg](https://github.com/baiyuze/copysync/releases/latest/download/CopySync.dmg)，打开后把 CopySync 拖进「应用程序」。
2. 从「应用程序」里打开 CopySync。

   第一次打开时 macOS 会提示「无法验证开发者」：CopySync 还没有经过 Apple 公证（需要付费开发者账号）。到 **系统设置 → 隐私与安全性**，在页面下方点 **仍要打开**。只需要这一次。

3. 点 **启用后台同步**。
4. 在 **设置 → 信令服务器地址** 填入 `ws://<服务器地址>:8787/signal`，按回车。
5. 较新的 macOS 会询问是否允许 CopySync Daemon 读取剪贴板，选择允许。之后可以在 **系统设置 → 隐私与安全性 → 从其他 App 粘贴** 里改为始终允许，App 里也有直达按钮。

**Windows**

1. 下载 [CopySync-Setup.exe](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-Setup.exe) 并运行。装在当前用户目录下，不需要管理员权限。

   CopySync 还没有代码签名，Windows 会提示「Windows 已保护你的电脑」：点 **更多信息 → 仍要运行**。不想运行安装程序的话，可以下载[便携版 zip](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-windows-x64.zip)，解压后直接打开。

2. 打开 CopySync（安装完会自动打开，桌面和开始菜单里也有），点 **启用后台同步**。
3. 在 **设置 → 信令服务器地址** 填入 `ws://<服务器地址>:8787/signal`，按回车。
4. Windows 防火墙询问是否允许 copysyncd 通信时，选择允许。

关掉窗口后，CopySync 在右下角的托盘里继续运行；从托盘图标或桌面快捷方式都能再打开它。

### 3. 配对

1. 在一台电脑上点 **设备 → 添加设备 → 生成配对码**。
2. 在另一台上点 **添加设备 → 输入配对码**，输入那 6 位字符。
3. 两台屏幕上会显示同样的两行安全指纹（这台的和对方的），**逐行核对两边完全一致后**再点确认。

之后在任意一台上复制，另一台就能直接粘贴。

## 工作原理

```mermaid
flowchart LR
    subgraph A["电脑 A（Mac 或 Windows）"]
        UA["CopySync 界面"] <-- "gRPC" --> DA["后台服务"]
    end
    subgraph B["电脑 B（Mac 或 Windows）"]
        DB["后台服务"] <-- "gRPC" --> UB["CopySync 界面"]
    end
    S[("你的服务器\n信令 + TURN")]
    DA -- "信令（WebSocket）" --> S
    DB -- "信令（WebSocket）" --> S
    DA <== "WebRTC 直连（DTLS 加密）" ==> DB
    DA <-. "打不通时经 TURN 中转" .-> S
    S <-.-> DB
```

| 组件 | 做什么 | 技术 |
|---|---|---|
| 后台服务 `copysyncd` | 监听剪贴板、与对端建立连接、收发与落盘 | Go。剪贴板部分在 Mac 上用 cgo 调用 AppKit，在 Windows 上直接调用 Win32（纯 Go，不需要 C 编译器） |
| 界面 `CopySync` | 复制记录、设备配对、设置 | Flutter，经本机 gRPC 与后台服务通信；Mac 与 Windows 共用一套代码 |
| 服务器 `copysync-server` | 设备发现、转达配对码、转发信令、TURN 中转 | Go，单个静态二进制，约 15 MB 内存 |
| 协议 `proto/` | 三方共用的消息定义 | Protocol Buffers |

**一次复制的经过：**

1. 后台服务发现剪贴板变了：Mac 上每隔一小段时间检查剪贴板的变更计数（只看计数和类型不需要授权，也不会触发系统提示）；Windows 上由系统即时通知。
2. 发现变化后读取内容，按类型和大小决定怎么传：文本直接内联在消息里；50 MB 以内的图片与文件打包成流（tar + zstd）立即推送；更大的只发一条记录。
3. 对端收到后先写入本地缓存，再写入剪贴板。在那台电脑上 `⌘V` 或 `Ctrl+V` 粘贴出来的，就是本地的真实文件。Mac 上合法、Windows 上不合法的文件名（比如带冒号的）在 Windows 上换成对应的全角字符。

**连接是怎么建立的：**

1. 设备通过 WebSocket 连上信令服务器，拿到 STUN 与 TURN 的地址。
2. 所有连接共用一个本地 UDP 端口。后台服务在这个端口上向多台服务器探测本机的公网地址，网络有几个出口就能拿到几个地址（见下一节）。
3. 握手时把局域网地址、每个出口的公网地址、TURN 中转地址都发给对端，两边用 WebRTC 的 ICE 同时互相尝试。其中一方先不公布出口地址，等收到对方的、先往那边发过包再公布：有的路由器会被对方先到的包占住端口（见下一节）。
4. 直连打通就走直连；中转路径至少等 5 秒才允许被选中，给直连留出时间。
5. 连上后两端互相告知自己是否经过中转，任一端经中转，两边都显示「中转」。
6. 若还是先落到了中转，发起连接的一方会在空闲时用 ICE 重启再打一次洞，打通就换成直连（10 秒后第一次，之后间隔逐步拉长到 15 分钟）。

数据走 WebRTC 数据通道，DTLS 加密。

### 多出口网络的打洞

公司双线、运营商 NAT 地址池这类网络有多个出口，按目标 IP 选择走哪一个。实测一个公司网络有 4 个出口：

| 出口 | 落在这个出口上的探测服务器 | 公网端口 |
|---|---|---|
| 出口 A | 自己的服务器、Twilio | 不变 |
| 出口 B | B 站、Google、Nextcloud | 不变 |
| 出口 C | Cloudflare、FreeSWITCH | 改写，但对不同目标相同 |
| 出口 D | 阿里云上的服务 | 改写 |

只问一台服务器，就只知道其中一个出口的地址；发往对端的包按目标 IP 很可能从另一个出口出去，对端收到的源地址对不上，被它的路由器丢弃。CopySync 1.0 就是这样只能走中转。

1.1 的做法：

- **共用一个端口**：之前 pion 每问一台 STUN 服务器就开一个新端口，问到的地址只对那个端口成立。现在所有连接共用一个端口，探测结果对它全部成立。
- **探测每个出口**：在这个端口上并发探测 IP 分散的服务器（自己的服务器，加上 B 站、Cloudflare 等公共 STUN），按 IP 去重，过滤被 DNS 污染的地址。
- **全部告诉对端**：每个出口的地址都作为候选地址发给对端，对端逐个检查，与实际出口吻合的那个能打通。会改写端口的出口拿到的是真实映射。
- **记住见过的出口**：某一轮漏掉的出口，若它不改端口，就按「出口 IP + 本地端口」推算补上。

实测公司 4 个出口与家用宽带之间直连，连上服务器后约 3 秒建立。设计细节、测量过程与验证见 [多出口网络的 NAT 打洞](design/nat-traversal.md)。

**NAT 打洞实验室**（`tools/natlab`）用 Linux 网络命名空间与 iptables 搭出真实拓扑，用 CopySync 真实的连接代码验证 11 种场景，每次提交都在 CI 里运行：

| 场景 | 结果 |
|---|---|
| 单出口 ↔ 家用路由器 | 直连 |
| 两个出口，旧版本（对照） | 中转 |
| 两个出口 ↔ 家用路由器 | 直连 |
| 四个出口 ↔ 家用路由器 | 直连 |
| 两个出口，第二条没探测到：无历史 / 有历史 | 中转 / 直连 |
| 两个出口 ↔ 两个出口 | 直连 |
| 两个出口 ↔ 端口随机分配的 NAT | 中转 |
| 一开始打不通，之后网络恢复 | 先中转，随后换成直连 |
| 公司路由器会为外来的包占住端口 | 直连 |
| 家用路由器会为外来的包占住端口 | 先中转，第一次重试后直连 |

### 安全模型

- **服务器不存账号**。设备身份是一把 Ed25519 密钥，设备 ID 由公钥派生，服务器据此就能验证"这个 ID 确实属于这把公钥"，无需用户表。
- **配对时人工核对指纹**。两台设备显示同样的两行：双方各自公钥的指纹（SHA-256 的前 60 位，形如 `R8NF-2WTC-QL5J`），按固定顺序排列。服务器可以在配对时替换转发中的公钥，但替换后两块屏幕上的内容就对不上了，这一步就是防中间人的闸门。两个指纹分开显示而不是合成一串短码，是因为攻击者同时控制两把伪造公钥时，凑出一串相同的短码只需生日攻击；分开显示则要为每把公钥各做一次原像攻击，难度高出约十亿倍。
- **之后的每条信令都带签名**，服务器无法伪造设备。直连通道的 DTLS 证书指纹也经签名传递，握手后比对实际证书，不一致立即断开。
- **中转也看不到内容**。TURN 只转发加密后的 UDP 包。中转凭证按设备单独签发，12 小时过期。
- **公共 STUN 服务器只看到公网 IP**。为了在多出口网络里找到每个出口，默认会向几台公共 STUN 服务器探测，它们能看到你的公网 IP（与任何 WebRTC 应用一样），看不到内容。介意的话在设置里关掉「用公共服务器探测网络出口」，只用自己的服务器。
- **设备私钥**：Mac 上以仅本人可读的权限存放；Windows 上另用系统的 DPAPI 加密，换个用户、或把文件拷到别的电脑上都解不开。
- **不启用 App 沙盒**：需要访问你复制的任意路径下的文件。

### 兜底：出问题时会发生什么

| 情况 | 处理方式 |
|---|---|
| 两端打不通直连（对称 NAT、公司防火墙） | 自动改走服务器的 TURN 中转，内容依然加密；界面上显示「中转」 |
| 公共 STUN 在国内时通时断 | 服务器自带 STUN，客户端优先用它 |
| 公司双线等多出口网络（按目标地址选出口） | 所有连接共用一个本地端口，在这个端口上向多台服务器探测每个出口的地址，全部发给对端由它逐个尝试；探测漏掉的出口按历史推算。`copysync-cli nat` 可以看到探测结果 |
| 先落到了中转，其实打得通（对端刚启动、出口地址晚到一步） | 发起方在空闲时用 ICE 重启再打一次洞，打通就换成直连；有传输在进行时不重启。重启期间连接不断，只停顿几秒 |
| 对方睡眠、断网、重启，连接悄悄断了 | 发现连接失效立即关闭、重连；对方重建了连接，本端也随之换新 |
| 复制时对方连不上 | 记下最近 20 条，连接恢复后补发。补发的只进复制记录，不写进对方的剪贴板，免得覆盖对方此刻正在用的内容 |
| 信令服务器断开 | 指数退避自动重连（1 秒起，最长 30 秒，带随机抖动），界面实时显示连接状态 |
| 配对时两端确认有先后 | 先确认一方发来的握手消息会被暂存，另一方确认后重放，配对完成几十毫秒内即可直连 |
| 大文件 | 超过阈值只同步记录，按需拉取；拉取时源文件已被删除会明确提示，而不是失败得不明不白 |
| 只认纯文本的输入框 | 带格式的文本同时写入 HTML 与纯文本两份，粘贴到哪里都不会出现标签 |
| 还没有授权读取剪贴板 | 不去读内容（否则会卡在一个后台进程看不见的系统弹窗上），界面上提示去授权 |
| 后台服务崩溃 | 由 launchd 自动拉起；10 秒节流，不会疯狂重启 |
| 覆盖安装了新版、或 App 换了位置 | 打开 App 时自动校正登录项，并重启到新版本 |
| 记录和缓存越来越多 | 到期自动清理，默认保留 3 天 |

### 踩过的坑

这些都是实测撞出来的，各自对应一个回归测试。完整的验证过程见 [spikes/results.md](spikes/results.md)。

- **macOS 的剪贴板调用必须在主线程。**开启隐私机制后，剪贴板 API 要经系统 UI 服务通信，依赖 run loop，在普通线程上调用会直接崩溃。后台服务启动时就把主线程锁住，专门留给剪贴板，业务逻辑投递过去执行。
- **"对方按下粘贴键才传输"在 macOS 上做不到。**系统在写入剪贴板后约 0.1 秒就会把延迟提供的数据兑现并缓存，真正按下 `⌘V` 时不再回调，也没有"剪贴板被读取"的通知。所以小文件是预先推送的，大文件改为手动拉取。
- **WebRTC 库（pion）的三个坑**：`DetachDataChannels()` 是全局开关，开了以后所有通道的普通收发都失效；`OnMessage` 不缓冲，回调设置之前到达的消息直接丢弃（实测发 202 字节，收到 0 字节）；`OnBufferedAmountLow` 是边沿触发，等待方晚一步进入等待就再也等不到唤醒（20 MB 传到 786 KB 卡死）。
- **pion 为每台 STUN 服务器单独开端口**，问到的地址互相用不上。多出口网络里多配几台 STUN 也打不通，所以改为所有连接共用一个端口、由自己探测。
- **单靠一端判断不了是否中转。**一端经 TURN 发送时，另一端看到的可能只是个普通地址，曾出现一台显示直连、一台显示中转。现在连接就绪后两端互相告知。
- **pion 选定线路后就不再更换。**中转握手快，一旦赢了这场赛跑，整个连接期间都走中转。实测对端重启后的第一次握手，出口地址比中转的等待时限晚到，就这样一直中转。现在中转等待放宽到 5 秒、探测时并发解析域名（整轮从约 1.6 秒降到 0.2 秒以内），并在中转时用 ICE 重启再试直连。
- **谁先发包很重要。**有的路由器（实测公司那台）收到对方先到的包，会留下一条连接记录、占住端口，本端随后往外发只好换端口，打洞就失败了。这解释了为什么本机重启时总能直连、家里那台一重启就中转：本机重启时出口要现探测，地址晚发出，本机反而先发包。现在由规则决定谁先发包，换直连的重试里两种顺序轮流用。
- **连接失效后要自己清理。**pion 报告「失败」后不会关闭连接，控制通道看上去还是打开的，写进去的消息石沉大海；上层以为连接还在，也不会重连。1.1.0 曾因此在对方断线后，把 6 个文件「发」进了一条死连接，对方一个也没收到。
- **公共 STUN 域名可能被 DNS 污染**，实测被解析到 `192.0.2.42`。这类保留地址和代理软件的 fake-ip 网段在探测前就过滤掉。
- **测试用的模拟路由器要有防火墙。**两端同时打洞时，对方的包若被放行到路由器本机，Linux 会为它留下一条连接记录，本端随后往外发时只好换端口，打洞就失败了。真实路由器直接丢弃这类包，实验室也照此配置。

局域网直连实测：20 MB 文件约 330 毫秒传完，校验和一致。

## 已知限制

- **没有经过 Apple 公证，Windows 版也没有代码签名**，第一次打开时需要手动放行一次。
- **Windows 上不支持「虚拟文件」**：从 Outlook 附件、压缩包里直接复制的文件，剪贴板里没有真实路径，不会同步。先保存或解压出来再复制即可。
- **Mac 版还不会跳过密码管理器复制的内容**，Windows 版已经会。
- **对方离线太久的内容不会补发**。服务器不存数据；对方暂时连不上时，本机只在内存里记下最近 20 条（24 小时内的），本机的后台服务重启后就清空了。
- **中转模式的速度取决于服务器带宽**。直连时不受影响。

## 从源码构建

Mac 版需要 Go 1.26、Flutter 3.44 和 Xcode：

```bash
./scripts/build.sh       # 构建全部到 dist/
./scripts/package.sh     # 打出 DMG 与服务器发布包到 release/
```

Windows 版需要 Go、Flutter、Visual Studio（「使用 C++ 的桌面开发」）与 Inno Setup 6，在 Windows 上运行：

```powershell
powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1   # 安装程序与便携版到 dist\windows\
```

后台服务是纯 Go，也可以在 Mac 上交叉编译：`GOOS=windows go build ./cmd/copysyncd`。CI 在每次提交时构建 Windows 安装包。

开发时可以单独运行各部分：

```bash
cd client-core && go test ./...          # 后台服务的测试
./tools/natlab/docker.sh                 # NAT 打洞实验室：11 种网络拓扑下的直连与中转
cd ui && flutter test                    # 界面的测试
./scripts/install-macos.sh               # 把 dist/ 里的后台服务装成登录项
./dist/copysync-cli status               # 查看后台服务状态
./dist/copysync-cli watch                # 实时事件流
```

改了 `.proto` 之后运行 `./scripts/gen-proto.sh` 重新生成代码（用 buf，`go install` 即可安装，无需 protoc）。改了界面之后运行 `./tools/screenshots/render.sh` 重新生成截图，截图用演示数据渲染，不会带出本机的真实内容。

```
client-core/    后台服务（Go）
  internal/clipboard    剪贴板（Mac：cgo + Objective-C；Windows：Win32；都必须跑在主线程）
  internal/transport    信令（WebSocket）与 P2P（WebRTC）
  internal/sync         同步引擎：分流、打包、传输、落盘
ui/             界面（Flutter）
server/         信令 + TURN 中转服务器（Go），部署文件在 server/deploy/
proto/          消息定义与设备 ID 派生规则，三方共用
docs/           项目主页（GitHub Pages）
spikes/         开工前的技术验证
design/         技术方案文档
tools/natlab/   NAT 打洞实验室
tools/installer/ Windows 安装程序（Inno Setup）
```

## 排障

**剪贴板读不到内容，`pbpaste` 也是空的**
系统的剪贴板服务可能卡住了：某个程序声明了剪贴板内容后异常退出，会留下一个"有类型声明、没有数据"的空壳。执行 `killall pboard`，再复制一次即可。

**设备一直显示「离线」**
检查两台电脑的 **设置 → 连接状态** 是否都是「已连接」。不是的话，确认服务器地址正确、8787 端口可以访问。

**一直显示「中转」而不是「直连」**
刚连上时显示中转、过一会儿变成直连是正常的：先用中转保证能用，空闲时再换直连。一直是中转，说明 NAT 打洞没成功，常见于端口随机分配的 NAT 或严格的防火墙。功能不受影响，只是速度受服务器带宽限制。运行 `copysync-cli nat` 看本机探测到了几个出口：出口为 0 说明 UDP 被拦截了。

**查看日志**
后台服务的日志：Mac 上在 `~/Library/Logs/CopySync/daemon.log`，Windows 上在 `%LOCALAPPDATA%\CopySync\logs\daemon.log`。服务器的日志用 `journalctl -u copysync-server -f` 查看。

## 许可

[MIT](LICENSE)

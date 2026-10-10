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
  <a href="README.md">English</a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/baiyuze/copysync?label=%E7%89%88%E6%9C%AC&color=0A66D8" alt="版本">
  <img src="https://img.shields.io/badge/macOS-13%2B-1D1D1F" alt="macOS 13+">
  <img src="https://img.shields.io/badge/Windows-10%20%2F%2011-1D1D1F" alt="Windows 10 / 11">
  <img src="https://img.shields.io/badge/Intel%20%2B%20Apple%20%E8%8A%AF%E7%89%87-%E9%80%9A%E7%94%A8-1D1D1F" alt="Intel 与 Apple 芯片通用">
  <img src="https://img.shields.io/badge/%E8%AE%B8%E5%8F%AF-AGPL--3.0-6E6E73" alt="AGPL-3.0 许可">
</p>

<p align="center">
  <img src="docs/assets/screenshots/history.png" alt="CopySync 的复制记录界面：来自其他 Mac 的文件、图片和文字">
</p>

## 它解决什么问题

通用剪贴板要同一个 Apple ID、两台离得近；隔空投送每次都要点接收，也不管剪贴板；发给自己要手动下载，文件还留在别人的服务器上。一台 Mac、一台 Windows 时，系统自带的办法一个都用不了。

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
- **按大小分流**：50 MB 以内（可调）复制时直接推送；更大的文件只同步一条记录，需要时再拉取。
- **复制记录**：所有设备上最近 3 天（可调）复制过的内容，图片可以预览。
- **直连优先，自动中转**，公司双线这类多出口网络也能直连。
- **配对时核对安全指纹**，防止服务器或网络上的第三方冒充你的设备。
- **后台常驻**，开机自启，崩溃自动恢复。中文、英文、日文界面。

<table>
  <tr>
    <td><img src="docs/assets/screenshots/devices.png" alt="设备页：已配对的设备、直连或中转状态"></td>
    <td><img src="docs/assets/screenshots/verify.png" alt="配对时核对两台设备上显示的安全指纹"></td>
  </tr>
</table>

## 开始使用

需要**每台电脑装一个 App**，再加**一台两端都能访问到的服务器**。

**1. 部署服务器。**

- *同一个局域网*（家里或办公室）：放在任意一台常开的 Windows、Mac 或 Linux 电脑上，下载后运行一个安装脚本即可。图文教程：**[在局域网里部署服务器](https://baiyuze.github.io/copysync/server/)**。
- *跨公网*：一台有公网 IP 的 Linux 服务器，需要 root：

  ```bash
  curl -LO https://github.com/baiyuze/copysync/releases/latest/download/copysync-server-linux-amd64.tar.gz
  tar xzf copysync-server-linux-amd64.tar.gz && cd copysync-server-linux-amd64
  sudo ./install.sh <服务器公网IP>
  ```

  放行 TCP 8787、UDP 3478 与 UDP 32768–60999。Docker 部署与全部参数见[服务器部署说明](server/deploy/README.md)。国内云服务器的域名没有备案时，客户端里请直接填 IP。

**2. 在每台电脑上安装 App**：macOS 13 以上用 [CopySync.dmg](https://github.com/baiyuze/copysync/releases/latest/download/CopySync.dmg)，Windows 10/11 用 [CopySync-Setup.exe](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-Setup.exe)（也有[便携版](https://github.com/baiyuze/copysync/releases/latest/download/CopySync-windows-x64.zip)）。还没有公证和代码签名：Mac 上第一次打开时到 **系统设置 → 隐私与安全性** 点 **仍要打开**；Windows 上点 **更多信息 → 仍要运行**。然后点 **启用后台同步**，在 **设置 → 信令服务器地址** 填入 `ws://<服务器地址>:8787/signal`。

**3. 配对。**一台点 **设备 → 添加设备 → 生成配对码**，另一台点 **输入配对码**。两边显示的两行安全指纹完全一致后再确认。之后在任意一台上复制，另一台就能直接粘贴。

| | 要求 |
|---|---|
| **Mac 客户端** | macOS 13 Ventura 及以上；Intel 与 Apple 芯片通用 |
| **Windows 客户端** | Windows 10（21H2 及以上）与 11，x64 |
| **服务器** | Linux（x86_64、ARM64；systemd 或 Docker）、macOS，或 Windows 10/11 x64 |
| **iPhone / iPad** | 暂不支持：iOS 不允许 App 在后台监听剪贴板 |

## 文档

- [工作原理](guide/how-it-works.zh-CN.md)：架构、NAT 打洞、安全模型、兜底、踩过的坑
- [排障](guide/troubleshooting.zh-CN.md)与已知限制
- [服务器部署说明](server/deploy/README.md) · [局域网部署图文教程](https://baiyuze.github.io/copysync/server/)
- [从源码构建](guide/building.zh-CN.md)
- 技术方案：[多出口网络的 NAT 打洞](design/nat-traversal.md)、[Windows 客户端](design/windows-client.md)
- [更新记录](CHANGELOG.md)

## 许可

CopySync 采用双授权。按 [GNU AGPL-3.0](LICENSE) 免费使用：任何人（包括公司）都可以使用和修改，但分发它、或通过网络提供修改过的版本时，必须以同样的许可公开你的源码。闭源产品、不公开的修改、托管服务等情况可以购买[商业授权](LICENSING.md)，请[提交「商业授权」issue](https://github.com/baiyuze/copysync/issues/new?template=commercial-license.yml)。1.2.0 及更早的版本仍适用 MIT 许可。

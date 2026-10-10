<#
在 Windows 上安装 copysync-server：注册为开机自启的系统服务，并在 Windows 防火墙里放行。

  双击 install.cmd               用本机的局域网地址启用 TURN 中转
  install.cmd 192.168.1.20       指定 TURN 对外公布的地址
  双击 uninstall.cmd             卸载

重复执行即为升级：替换程序并重启服务，已有设置（含 TURN 密钥）保留；指定了地址时改用新地址。
需要管理员权限，不是管理员时会弹出 UAC 确认，在新窗口里继续。

这个文件必须保存为带 BOM 的 UTF-8：Windows PowerShell 5.1 把不带 BOM 的脚本当作系统代码页读取，中文会乱码。
控制台按系统代码页输出，所以只用 GBK 与 CP437 里都有的符号（√、!），其他对勾、箭头会变成问号。
#>
param(
    [string]$Ip = "",
    [switch]$Uninstall
)

$ErrorActionPreference = 'Stop'

# 中文系统显示中文，其他显示英文
$zh = (Get-UICulture).Name -like 'zh*'
function Say([string]$zhText, [string]$enText) { if ($zh) { $zhText } else { $enText } }
function Step([string]$zhText, [string]$enText) { Write-Host ('> ' + (Say $zhText $enText)) }

$admin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole(
    [Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $admin) {
    # 以管理员身份在新窗口里重新运行；-NoExit 让窗口留着，看得到结果
    $argList = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-NoExit', '-File', "`"$PSCommandPath`"")
    if ($Ip) { $argList += @('-Ip', $Ip) }
    if ($Uninstall) { $argList += '-Uninstall' }
    try {
        Start-Process powershell.exe -Verb RunAs -ArgumentList $argList
    } catch {
        Write-Host (Say '需要管理员权限才能注册服务、设置防火墙。' 'Administrator rights are needed to register the service and configure the firewall.')
        exit 1
    }
    exit 0
}

$Name = 'CopySyncServer'
$Display = 'CopySync Server'
$InstallDir = Join-Path $env:ProgramFiles 'CopySync Server'
$DataDir = Join-Path $env:ProgramData 'CopySync Server'
$Exe = Join-Path $InstallDir 'copysync-server.exe'
$Conf = Join-Path $DataDir 'server.conf'
$Log = Join-Path $DataDir 'server.log'

function Stop-Server {
    if (Get-Service $Name -ErrorAction SilentlyContinue) {
        Stop-Service $Name -Force -ErrorAction SilentlyContinue
    }
    # 服务停了进程不一定立刻退出，退出前替换不了程序文件
    for ($i = 0; $i -lt 20 -and (Get-Process copysync-server -ErrorAction SilentlyContinue); $i++) {
        Start-Sleep -Milliseconds 250
    }
}

if ($Uninstall) {
    Stop-Server
    if (Get-Service $Name -ErrorAction SilentlyContinue) { sc.exe delete $Name | Out-Null }
    Get-NetFirewallRule -DisplayName $Display -ErrorAction SilentlyContinue | Remove-NetFirewallRule
    Remove-Item $InstallDir, $DataDir -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host (Say '√ 已卸载' '√ Uninstalled')
    exit 0
}

# 本机的局域网地址：默认路由所用网卡的地址，按跃点数取最优先的一个。
# 跳过代理软件 TUN 模式常用的 198.18.0.0/15 与链路本地地址。
function Get-LanIp {
    $routes = Get-NetRoute -DestinationPrefix '0.0.0.0/0' -ErrorAction SilentlyContinue |
        Sort-Object { $_.RouteMetric + (Get-NetIPInterface -InterfaceIndex $_.InterfaceIndex -AddressFamily IPv4).InterfaceMetric }
    foreach ($r in $routes) {
        foreach ($a in Get-NetIPAddress -InterfaceIndex $r.InterfaceIndex -AddressFamily IPv4 -ErrorAction SilentlyContinue) {
            if ($a.IPAddress -notmatch '^(169\.254\.|198\.1[89]\.|127\.)') { return $a.IPAddress }
        }
    }
}

$src = Join-Path $PSScriptRoot 'copysync-server.exe'
if (-not (Test-Path $src)) {
    throw (Say "找不到 $src，请把压缩包完整解压后再运行" "$src not found. Extract the whole zip first.")
}

# 读取已有设置
$cfg = @{}
if (Test-Path $Conf) {
    foreach ($line in Get-Content $Conf) {
        if ($line -match '^(\w+)=(.*)$') { $cfg[$Matches[1]] = $Matches[2] }
    }
}
if ($Ip) { $cfg.TURN_IP = $Ip } elseif (-not $cfg.TURN_IP) { $cfg.TURN_IP = Get-LanIp }
if (-not $cfg.TURN_IP) {
    throw (Say '找不到本机的局域网地址，请指定：install.cmd <IP>' "Couldn't find this PC's LAN address. Specify it: install.cmd <IP>")
}
if (-not $cfg.TURN_SECRET) {
    # 固定 secret：每次重启都换的话，客户端手里已签发的 TURN 凭证会失效
    $bytes = New-Object byte[] 32
    [Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
    $cfg.TURN_SECRET = -join ($bytes | ForEach-Object { $_.ToString('x2') })
}
$turnIp = $cfg.TURN_IP

Step '停止旧版本' 'Stopping the running version'
Stop-Server

Step "安装到 $InstallDir" "Installing to $InstallDir"
New-Item -ItemType Directory -Force $InstallDir, $DataDir | Out-Null
Copy-Item $src $Exe -Force
Unblock-File $Exe
# 服务以低权限的 LOCAL SERVICE 账户运行，只给它写日志的权限
icacls $DataDir /grant '*S-1-5-19:(OI)(CI)M' | Out-Null
Set-Content $Conf -Encoding ASCII -Value @("TURN_IP=$turnIp", "TURN_SECRET=$($cfg.TURN_SECRET)")
# 设置里有 TURN 密钥，只留给管理员与系统
icacls $Conf /inheritance:r /grant:r '*S-1-5-32-544:F' '*S-1-5-18:F' | Out-Null

Step "注册服务 $Name（中转地址 $turnIp）" "Registering the service $Name (relay address $turnIp)"
$bin = "`"$Exe`" -addr :8787 -turn-ip $turnIp -turn-secret $($cfg.TURN_SECRET) " +
       "-stun stun:${turnIp}:3478,stun:stun.l.google.com:19302 -log `"$Log`""
$desc = Say 'CopySync 的信令与 TURN 中转服务器' 'Signaling and TURN relay server for CopySync'
$existing = Get-CimInstance Win32_Service -Filter "Name='$Name'"
if ($existing) {
    $r = $existing | Invoke-CimMethod -MethodName Change -Arguments @{ PathName = $bin; StartMode = 'Automatic' }
    if ($r.ReturnValue -ne 0) { throw "Win32_Service.Change: $($r.ReturnValue)" }
} else {
    $cred = New-Object Management.Automation.PSCredential('NT AUTHORITY\LocalService', (New-Object Security.SecureString))
    New-Service -Name $Name -BinaryPathName $bin -DisplayName $Display -Description $desc `
        -StartupType Automatic -Credential $cred | Out-Null
}
# 意外退出（包括端口被占用）后自动重启：5 秒、5 秒、30 秒，一天后重新计数
sc.exe failure $Name reset= 86400 actions= restart/5000/restart/5000/restart/30000 | Out-Null
sc.exe failureflag $Name 1 | Out-Null

Step '在 Windows 防火墙里放行' 'Allowing it through Windows Firewall'
Get-NetFirewallRule -DisplayName $Display -ErrorAction SilentlyContinue | Remove-NetFirewallRule
# 按程序放行：信令、STUN/TURN 与系统随机分配的中转端口都覆盖到，不用逐个开端口
New-NetFirewallRule -DisplayName $Display -Description $desc -Direction Inbound -Program $Exe `
    -Action Allow -Profile Any | Out-Null

Step '启动' 'Starting'
Start-Service $Name
$ok = $false
$web = New-Object Net.WebClient
$web.Proxy = $null
for ($i = 0; $i -lt 20 -and -not $ok; $i++) {
    try { $null = $web.DownloadString('http://127.0.0.1:8787/healthz'); $ok = $true } catch { Start-Sleep -Milliseconds 250 }
}
if (-not $ok) {
    Write-Host (Say "  ! 启动失败，查看日志：$Log" "  ! Failed to start. See the log: $Log")
    exit 1
}
Write-Host (Say '  √ 已启动' '  √ Running')

$version = & $Exe -version
if ($zh) {
    Write-Host @"

安装完成：$version，作为 Windows 服务开机自动运行
在每台电脑的 CopySync「设置 → 信令服务器地址」填：

    ws://${turnIp}:8787/signal

  程序   $InstallDir
  日志   $Log
  重启   Restart-Service $Name（管理员 PowerShell）
  卸载   双击 uninstall.cmd
"@
} else {
    Write-Host @"

Installed: $version. It runs as a Windows service and starts with Windows.
In CopySync on each computer, set Settings > Signaling server to:

    ws://${turnIp}:8787/signal

  Program    $InstallDir
  Logs       $Log
  Restart    Restart-Service $Name (in an administrator PowerShell)
  Uninstall  double-click uninstall.cmd
"@
}

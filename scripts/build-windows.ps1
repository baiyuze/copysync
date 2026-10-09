# 在 Windows 上构建 CopySync，产物放在 dist\windows\：
#
#   CopySync\                    可以直接运行的目录：界面、后台服务、命令行工具
#   CopySync-Setup.exe           安装程序（需要 Inno Setup 6）
#   CopySync-windows-x64.zip     便携版
#
#   powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1 [-Version 1.3.0]
#
# 需要 Go、Flutter 与 Visual Studio（「使用 C++ 的桌面开发」）。
# 后台服务是纯 Go，也可以在 Mac 上交叉编译；界面只能在 Windows 上构建。

param([string]$Version)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Set-Location $root

# 版本号与 Mac 版共用 scripts/build.sh 里的那一个
if (-not $Version) {
    $m = Select-String -Path scripts\build.sh -Pattern 'VERSION=\$\{VERSION:-([0-9.]+)\}'
    $Version = $m.Matches[0].Groups[1].Value
}

$out = Join-Path $root 'dist\windows'
$app = Join-Path $out 'CopySync'
if (Test-Path $out) { Remove-Item -Recurse -Force $out }
New-Item -ItemType Directory -Force $app | Out-Null

function Invoke-Checked([scriptblock]$cmd, [string]$what) {
    & $cmd
    if ($LASTEXITCODE -ne 0) { throw "$what 失败（退出码 $LASTEXITCODE）" }
}

Write-Host "▶ 构建后台服务与命令行工具（$Version）"
$env:CGO_ENABLED = '0'; $env:GOOS = 'windows'; $env:GOARCH = 'amd64'
Push-Location client-core
try {
    # -H=windowsgui：后台服务没有控制台窗口，开机自启时不会闪一下黑框
    Invoke-Checked { go build -trimpath -ldflags "-s -w -H=windowsgui -X main.version=$Version" -o "$app\copysyncd.exe" ./cmd/copysyncd } 'go build copysyncd'
    Invoke-Checked { go build -trimpath -ldflags '-s -w' -o "$app\copysync-cli.exe" ./cmd/copysync-cli } 'go build copysync-cli'
} finally {
    Pop-Location
    Remove-Item Env:CGO_ENABLED, Env:GOOS, Env:GOARCH
}

Write-Host '▶ 构建图形界面'
Push-Location ui
try {
    Invoke-Checked { flutter build windows --release --build-name $Version "--dart-define=APP_VERSION=$Version" } 'flutter build windows'
} finally {
    Pop-Location
}
Copy-Item -Recurse -Force 'ui\build\windows\x64\runner\Release\*' $app

Write-Host '▶ 打包'
Compress-Archive -Path $app -DestinationPath (Join-Path $out 'CopySync-windows-x64.zip')

$iscc = Get-Command iscc.exe -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source
if (-not $iscc) {
    $iscc = @(
        "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
        "$env:ProgramFiles\Inno Setup 6\ISCC.exe",
        "$env:LOCALAPPDATA\Programs\Inno Setup 6\ISCC.exe"
    ) | Where-Object { Test-Path $_ } | Select-Object -First 1
}
if ($iscc) {
    Invoke-Checked { & $iscc /Qp "/DAppVersion=$Version" "/DSourceDir=$app" "/DOutputDir=$out" tools\installer\copysync.iss } 'Inno Setup'
} else {
    Write-Warning '没有找到 Inno Setup 6，跳过安装程序，只出便携版 zip'
}

Get-ChildItem $out -File | ForEach-Object { Write-Host "  ✓ $($_.Name)  $([math]::Round($_.Length / 1MB, 1)) MB" }

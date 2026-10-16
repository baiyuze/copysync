; PasteFerry 的 Windows 安装程序（Inno Setup 6）。由 scripts\build-windows.ps1 调用：
;
;   iscc /DAppVersion=1.3.2 /DSourceDir=<构建好的目录> /DOutputDir=<输出目录> tools\installer\pasteferry.iss
;
; 按用户安装：装到 %LOCALAPPDATA%\Programs\PasteFerry，不需要管理员权限、不弹 UAC。
; 数据目录 %LOCALAPPDATA%\PasteFerry 与安装目录分开，升级、重装都不受影响。

#ifndef AppVersion
  #define AppVersion "0.0.0"
#endif
#ifndef SourceDir
  #define SourceDir "..\..\dist\windows\PasteFerry"
#endif
#ifndef OutputDir
  #define OutputDir "..\..\dist\windows"
#endif

[Setup]
; AppId 一旦发布就不能再改：Windows 靠它认出这是同一个程序的升级
AppId={{13E09E3B-04CA-435D-A596-19D818B1BE24}
AppName=PasteFerry
AppVersion={#AppVersion}
AppVerName=PasteFerry {#AppVersion}
AppPublisher=baiyuze
AppPublisherURL=https://pasteferry.cc/
AppSupportURL=https://github.com/baiyuze/pasteferry/issues
DefaultDirName={localappdata}\Programs\PasteFerry
; 2.0.0 之前叫 CopySync，装在 Programs\CopySync：不沿用旧目录，旧目录在安装时删掉
UsePreviousAppDir=no
DisableDirPage=yes
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
OutputDir={#OutputDir}
OutputBaseFilename=PasteFerry-Setup
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
; 任意版本的 Windows 10 与 11：Flutter 引擎与 Go 都要求 Windows 10，Windows 7、8 跑不起来
MinVersion=10.0
SetupIconFile=..\..\ui\windows\runner\resources\app_icon.ico
UninstallDisplayIcon={app}\PasteFerry.exe
UninstallDisplayName=PasteFerry
; 用「重启管理器」关掉正在运行的界面，免得文件被占用
CloseApplications=yes
RestartApplications=no
ShowLanguageDialog=no

[Languages]
; 简体、繁体中文的翻译文件不是每个 Inno Setup 版本都自带，有就用
#if FileExists(AddBackslash(CompilerPath) + "Languages\ChineseSimplified.isl")
Name: "chinesesimplified"; MessagesFile: "compiler:Languages\ChineseSimplified.isl"
#endif
#if FileExists(AddBackslash(CompilerPath) + "Languages\ChineseTraditional.isl")
Name: "chinesetraditional"; MessagesFile: "compiler:Languages\ChineseTraditional.isl"
#endif
Name: "english"; MessagesFile: "compiler:Default.isl"

[InstallDelete]
Type: filesandordirs; Name: "{localappdata}\Programs\CopySync"

[Files]
Source: "{#SourceDir}\*"; DestDir: "{app}"; Flags: recursesubdirs ignoreversion

[Icons]
Name: "{autoprograms}\PasteFerry"; Filename: "{app}\PasteFerry.exe"
Name: "{autodesktop}\PasteFerry"; Filename: "{app}\PasteFerry.exe"

[Registry]
; 卸载时删掉开机自启项与窗口位置记录。ValueType: none 表示安装时不写，只在卸载时删
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; ValueName: "PasteFerry"; ValueType: none; Flags: uninsdeletevalue
Root: HKCU; Subkey: "Software\PasteFerry"; ValueType: none; Flags: uninsdeletekey
Root: HKCU; Subkey: "Software\Classes\pasteferry-notify"; ValueType: none; Flags: uninsdeletekey
Root: HKCU; Subkey: "Software\Classes\AppUserModelId\PasteFerry.Desktop"; ValueType: none; Flags: uninsdeletekey

[Run]
; 升级前停掉了后台服务；之前启用过后台同步的话，装完立刻重新拉起，不必等用户打开界面
Filename: "{app}\pasteferryd.exe"; Parameters: "-supervise"; Flags: nowait runhidden; Check: AutostartRegistered
Filename: "{app}\PasteFerry.exe"; Description: "{cm:LaunchProgram,PasteFerry}"; Flags: nowait postinstall skipifsilent

[UninstallRun]
Filename: "{sys}\taskkill.exe"; Parameters: "/F /T /IM pasteferryd.exe"; Flags: runhidden; RunOnceId: "StopDaemon"
Filename: "{sys}\taskkill.exe"; Parameters: "/F /IM PasteFerry.exe"; Flags: runhidden; RunOnceId: "StopUI"

[Code]
function AutostartRegistered: Boolean;
begin
  Result := RegValueExists(HKEY_CURRENT_USER, 'Software\Microsoft\Windows\CurrentVersion\Run', 'PasteFerry');
end;

// 后台服务（含守护进程）没有可见窗口，重启管理器关不掉它，安装前直接结束。
// 2.0.0 之前叫 CopySync：旧的界面与后台服务也结束掉，旧的开机自启项换成新的，
// 数据目录由新的后台服务第一次启动时挪过来
function PrepareToInstall(var NeedsRestart: Boolean): String;
var
  Code: Integer;
  RunKey: String;
begin
  Exec(ExpandConstant('{sys}\taskkill.exe'), '/F /T /IM pasteferryd.exe', '', SW_HIDE, ewWaitUntilTerminated, Code);
  Exec(ExpandConstant('{sys}\taskkill.exe'), '/F /T /IM copysyncd.exe', '', SW_HIDE, ewWaitUntilTerminated, Code);
  Exec(ExpandConstant('{sys}\taskkill.exe'), '/F /IM CopySync.exe', '', SW_HIDE, ewWaitUntilTerminated, Code);
  RunKey := 'Software\Microsoft\Windows\CurrentVersion\Run';
  if RegValueExists(HKEY_CURRENT_USER, RunKey, 'CopySync') then
  begin
    RegWriteStringValue(HKEY_CURRENT_USER, RunKey, 'PasteFerry',
      '"' + ExpandConstant('{app}\pasteferryd.exe') + '" -supervise');
    RegDeleteValue(HKEY_CURRENT_USER, RunKey, 'CopySync');
  end;
  RegDeleteKeyIncludingSubkeys(HKEY_CURRENT_USER, 'Software\CopySync');
  Result := '';
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if (CurUninstallStep = usPostUninstall) and not UninstallSilent then
    if MsgBox('是否同时删除复制记录、缓存与这台电脑的设备身份？' + #13#10 + #13#10 +
              '选「否」则保留，重新安装后配对关系和记录都还在。',
              mbConfirmation, MB_YESNO or MB_DEFBUTTON2) = IDYES then
      DelTree(ExpandConstant('{localappdata}\PasteFerry'), True, True, True);
end;

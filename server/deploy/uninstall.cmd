@echo off
rem Double-click to remove the copysync-server service, its firewall rule and its files.
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1" -Uninstall
if errorlevel 1 pause

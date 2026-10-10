@echo off
rem Double-click to install copysync-server as a Windows service (asks for administrator rights).
rem To set the relay address yourself: install.cmd 192.168.1.20
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1" %*
if errorlevel 1 pause

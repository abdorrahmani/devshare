@echo off
:: Batch script to run install.ps1 as Administrator

:: Check for Administrator rights
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Requesting Administrator privileges...
    powershell -Command "Start-Process '%~f0' -Verb RunAs"
    exit /b
)

:: Run the PowerShell installer
powershell -ExecutionPolicy Bypass -NoProfile -File "%~dp0install.ps1"
pause 
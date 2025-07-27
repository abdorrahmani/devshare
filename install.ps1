$IsAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $IsAdmin) {
    Write-Host "This installer must be run as Administrator to install to Program Files and update the system PATH."
    Write-Host "Right-click on PowerShell and select 'Run as administrator', then run this script again."
    Write-Host "Press Enter to exit."
    [void][System.Console]::ReadLine()
    exit 1
}

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$binaryName = "devshare.exe"
$binaryPath = Join-Path $scriptDir $binaryName
$installDir = "C:\Program Files\DevShare"
$installPath = "$installDir\$binaryName"

if (!(Test-Path $binaryPath)) {
    Write-Host "$binaryName not found in current directory. Please run this script from the extracted archive folder."
    Write-Host "Press Enter to exit."
    [void][System.Console]::ReadLine()
    exit 1
}

# Create install directory if it doesn't exist
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

Copy-Item -Path $binaryPath -Destination $installPath -Force
Write-Host "$binaryName installed to $installPath"

$envPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
if ($envPath -notlike "*${installDir}*") {
    try {
        [System.Environment]::SetEnvironmentVariable("Path", "$envPath;$installDir", [System.EnvironmentVariableTarget]::Machine)
        Write-Host "$installDir added to system PATH. You may need to restart your terminal or log out and back in."
    } catch {
        Write-Host "Could not update system PATH. Please add '$installDir' to your PATH manually."
    }
}

Write-Host "Installation complete!"
Write-Host "Built by Anophel"
Write-Host "Press Enter to exit."
[void][System.Console]::ReadLine()
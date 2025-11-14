# PowerShell install script for bt-phone-home on Windows

param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$Username,

    [Parameter(Mandatory=$true, Position=1)]
    [string]$RoomName
)

# Configuration
$BinaryHost = if ($env:KAINET_HOST) { $env:KAINET_HOST } else { "https://github.com/rickhallett/kainet/releases/download/v1.1.0" }
$BinaryName = "kainet.exe"

# Colors
function Write-Green { param($text) Write-Host $text -ForegroundColor Green }
function Write-Cyan { param($text) Write-Host $text -ForegroundColor Cyan }
function Write-Red { param($text) Write-Host $text -ForegroundColor Red }
function Write-Yellow { param($text) Write-Host $text -ForegroundColor Yellow }

Write-Green "╔══════════════════════════════════════════════════════════╗"
Write-Green "║           BT PHONE HOME - INSTALLER                      ║"
Write-Green "╚══════════════════════════════════════════════════════════╝"
Write-Host ""

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$Platform = "windows-$Arch"

Write-Cyan ">>> detected: Windows $Arch"

# Create temp directory
$TempDir = Join-Path $env:TEMP "bt-phone-home-$([Guid]::NewGuid())"
New-Item -ItemType Directory -Path $TempDir -Force | Out-Null
Set-Location $TempDir

Write-Cyan ">>> downloading binary for $Platform..."

# Download binary
$DownloadUrl = "$BinaryHost/$BinaryName"
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $BinaryName -UseBasicParsing
} catch {
    Write-Red "ERROR: download failed"
    Write-Host "Could not download from: $DownloadUrl"
    Write-Host ""
    Write-Host "Set custom host with: `$env:BT_PHONE_HOME_HOST='https://your-host.com'"
    exit 1
}

if (-not (Test-Path $BinaryName)) {
    Write-Red "ERROR: download failed"
    exit 1
}

Write-Green ">>> download complete"

# Unblock file
Unblock-File -Path $BinaryName
Write-Green ">>> unblocked file"

Write-Host ""
Write-Green "╔════════════════════════════════════════╗"
Write-Green "║  INSTALLATION COMPLETE                 ║"
Write-Green "║  LAUNCHING SECURE TERMINAL...          ║"
Write-Green "╚════════════════════════════════════════╝"
Write-Host ""

# Execute binary
& ".\$BinaryName" $Username $RoomName

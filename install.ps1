# kaiNET installer for Windows
# Usage: .\install.ps1 <username> <room-name>

param(
    [Parameter(Mandatory=$true)]
    [string]$Username,

    [Parameter(Mandatory=$true)]
    [string]$RoomName
)

# Configuration
$GithubReleaseUrl = if ($env:KAINET_HOST) { $env:KAINET_HOST } else { "https://github.com/rickhallett/kainet/releases/latest/download" }
$InstallDir = "$env:LOCALAPPDATA\kainet"
$BinaryName = "kainet.exe"
$BinaryUrl = "$GithubReleaseUrl/$BinaryName"

# Functions
function Write-GreenText {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Green
}

function Write-BoldGreenText {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Green -BackgroundColor Black
}

# Main installation
Write-BoldGreenText "=========================================="
Write-BoldGreenText "   kaiNET Installer"
Write-BoldGreenText "=========================================="
Write-Host ""

Write-GreenText "Download URL: $BinaryUrl"

# Create install directory if it doesn't exist
if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$InstallPath = Join-Path -Path $InstallDir -ChildPath $BinaryName

try {
    # Download binary
    Write-GreenText "Downloading kaiNET..."
    Invoke-WebRequest -Uri $BinaryUrl -OutFile $InstallPath -UseBasicParsing

    # Unblock file (remove Windows security warning)
    Write-GreenText "Unblocking file..."
    Unblock-File -Path $InstallPath

    Write-GreenText "Binary installed to: $InstallPath"
    Write-Host ""
    Write-BoldGreenText "Installation complete!"
    Write-GreenText "Starting kaiNET..."
    Write-Host ""

    # Execute kaiNET
    & $InstallPath $Username $RoomName
}
catch {
    Write-Host "Error during installation: $_" -ForegroundColor Red
    exit 1
}

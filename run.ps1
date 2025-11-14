# kaiNET runner for WhatsApp .mp4 downloads (Windows)
# This script renames, unblocks, and executes the kaiNET binary
# Usage: .\run.ps1 <username> <room-name>

param(
    [Parameter(Mandatory=$true)]
    [string]$Username,

    [Parameter(Mandatory=$true)]
    [string]$RoomName
)

# Functions
function Write-GreenText {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Green
}

function Write-BoldGreenText {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Green -BackgroundColor Black
}

# Main
Write-BoldGreenText "=========================================="
Write-BoldGreenText "   kaiNET Runner"
Write-BoldGreenText "=========================================="
Write-Host ""

# Find .mp4 file in current directory
$Mp4File = Get-ChildItem -Path . -Filter "*.mp4" -File | Select-Object -First 1

if (-not $Mp4File) {
    Write-Host "Error: No .mp4 file found in current directory" -ForegroundColor Red
    exit 1
}

Write-GreenText "Found binary: $($Mp4File.Name)"

$BinaryName = "kainet.exe"

try {
    # Rename file
    Write-GreenText "Renaming to: $BinaryName"
    Rename-Item -Path $Mp4File.FullName -NewName $BinaryName

    # Unblock file (remove Windows security warning)
    Write-GreenText "Unblocking file..."
    Unblock-File -Path $BinaryName

    Write-Host ""
    Write-BoldGreenText "Starting kaiNET..."
    Write-Host ""

    # Execute kaiNET
    & ".\$BinaryName" $Username $RoomName
}
catch {
    Write-Host "Error during execution: $_" -ForegroundColor Red
    exit 1
}

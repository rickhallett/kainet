# PowerShell script to run bt-phone-home on Windows

param(
    [string]$Username,
    [string]$RoomName
)

if (-not $Username -or -not $RoomName) {
    if ($args.Count -ge 2) {
        $Username = $args[0]
        $RoomName = $args[1]
    } else {
        Write-Host "usage: .\run.ps1 <username> <room-name>" -ForegroundColor Red
        exit 1
    }
}

# Change to Downloads folder
Set-Location "$env:USERPROFILE\Downloads"

# Find the .mp4 file
$mp4File = Get-ChildItem -Filter "bt-phone-home*.mp4" | Select-Object -First 1

if (-not $mp4File) {
    Write-Host "error: no bt-phone-home*.mp4 file found in Downloads folder" -ForegroundColor Red
    exit 1
}

# Remove .mp4 extension
$binaryName = $mp4File.Name -replace '\.mp4$', '.exe'
Rename-Item -Path $mp4File.Name -NewName $binaryName

Write-Host "renamed $($mp4File.Name) to $binaryName" -ForegroundColor Green

# Unblock file (Windows equivalent of removing quarantine)
Unblock-File -Path $binaryName
Write-Host "unblocked file" -ForegroundColor Green

# Run with username and room
Write-Host "starting chat as $Username in room $RoomName..." -ForegroundColor Cyan
& ".\$binaryName" $Username $RoomName

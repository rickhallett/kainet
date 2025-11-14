# kainet

ultra-lightweight private chat for 2 people using turso database.

## features

- military-grade terminal interface with color-coded output
- AES-256-GCM end-to-end encryption
- no server to maintain
- async messaging (works offline)
- private key authentication via turso credentials
- ~1 second message delivery
- single binary cli client
- classified communications aesthetic

## quick start (curl install)

### macOS / Linux

**one-line install and run:**

```bash
curl -fsSL https://raw.githubusercontent.com/rickhallett/kainet/main/install.sh | bash -s -- <username> <room-name>
```

**example:**
```bash
curl -fsSL https://raw.githubusercontent.com/rickhallett/kainet/main/install.sh | bash -s -- alice secret-room
```

### Windows

**PowerShell install:**

```powershell
iex "& { $(irm https://raw.githubusercontent.com/rickhallett/kainet/main/install.ps1) } alice secret-room"
```

**what happens:**
- detects your OS and architecture
- downloads the correct binary
- handles permissions and quarantine
- launches the terminal automatically

**room names:**
- acts as a shared secret
- only people with the room name can join
- different rooms are completely isolated
- use different rooms for different groups

see [HOSTING.md](HOSTING.md) for binary hosting setup.

## direct download

**latest release:** https://github.com/rickhallett/kainet/releases/latest

download pre-built binaries:
- [macOS Apple Silicon (M1/M2/M3)](https://github.com/rickhallett/kainet/releases/download/v1.1.0/kainet-darwin-arm64)
- [macOS Intel](https://github.com/rickhallett/kainet/releases/download/v1.1.0/kainet-darwin-amd64)
- [Linux (amd64)](https://github.com/rickhallett/kainet/releases/download/v1.1.0/kainet-linux-amd64)
- [Windows (64-bit)](https://github.com/rickhallett/kainet/releases/download/v1.1.0/kainet.exe)

after downloading, make executable and run:
```bash
# macOS/Linux
chmod +x kainet-*
./kainet-darwin-arm64 <username> <room-name>

# Windows
.\kainet.exe <username> <room-name>
```

## manual setup

### 1. create turso database

```bash
# install turso cli (first time only)
brew install tursodatabase/tap/turso

# login
turso auth login

# create database
turso db create bt-phone-home

# get credentials
turso db show bt-phone-home --url
turso db tokens create bt-phone-home
```

save the database url and auth token - these are your "private keys" that both users need.

### 2. embed credentials and build

edit `main.go` and set your turso credentials:

```go
const (
	embeddedDBURL     = "libsql://your-db-name.turso.io"
	embeddedAuthToken = "your-auth-token-here"
)
```

then build:

```bash
go mod download
go build -o bt-phone-home
```

**security note:** credentials are baked into the binary. share the binary securely with the other person (signal, keybase, etc). anyone with the binary can access the chat.

### 3. run

**macOS - if downloaded via whatsapp (gets renamed to .mp4):**

send both files to the other person:
- `bt-phone-home` (or `bt-phone-home.mp4`)
- `run.sh`

one-line command:
```bash
cd ~/Downloads && chmod +x run.sh && xattr -d com.apple.quarantine run.sh && ./run.sh $USER secret-room
```

replace `secret-room` with your room name.

the script automatically:
- removes .mp4 extension
- makes binary executable
- removes quarantine flag
- runs the chat

**macOS - if transferred normally:**

```bash
chmod +x bt-phone-home
xattr -d com.apple.quarantine bt-phone-home
./bt-phone-home alice secret-room
```

**windows - if downloaded via whatsapp (gets renamed to .mp4):**

send both files to the other person:
- `bt-phone-home.exe` (or `bt-phone-home.exe.mp4`)
- `run.ps1`

one-line command (PowerShell):
```powershell
cd $env:USERPROFILE\Downloads; .\run.ps1 $env:USERNAME
```

the script automatically:
- removes .mp4 extension
- unblocks file (removes Windows security warning)
- runs the chat

**windows - if transferred normally:**

```powershell
cd $env:USERPROFILE\Downloads
Unblock-File bt-phone-home.exe
.\bt-phone-home.exe alice
```

**alternative:** if you prefer environment variables instead of embedded credentials, leave the constants empty and use:

```bash
export DB_URL="libsql://your-db-name.turso.io"
export AUTH_TOKEN="your-auth-token-here"
./bt-phone-home alice
```

## usage

- type messages and press enter to send
- new messages appear automatically (polls every 1 second)
- shows last 20 messages on startup
- ctrl+c to exit

## architecture

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│  alice cli  │────────▶│  turso db    │◀────────│   bob cli   │
│             │         │  (sqlite)    │         │             │
└─────────────┘         └──────────────┘         └─────────────┘
```

- no websocket server needed
- clients poll database for new messages
- turso free tier: 500 dbs, 9gb transfer/month
- messages persist even when both offline

## deployment options

### option 1: both users run locally
just share the db_url and auth_token

### option 2: distribute binary
with embedded credentials, just send the binary:

```bash
# for intel mac
GOARCH=amd64 GOOS=darwin go build -o bt-phone-home-intel

# for m1/m2 mac
GOARCH=arm64 GOOS=darwin go build -o bt-phone-home-arm64

# for windows
GOARCH=amd64 GOOS=windows go build -o bt-phone-home.exe

# for linux
GOARCH=amd64 GOOS=linux go build -o bt-phone-home-linux
```

send them the binary. no separate credentials file needed.

## web version

**try it now:** https://web-byjukk5lm-rick-halletts-projects.vercel.app

a browser-based terminal version is available:

features:
- retro CRT terminal with scanlines
- same encryption as CLI
- works in any browser
- deployed on Vercel

**local development:**
```bash
cd web
bun install
bun run dev
```

see `web/README.md` for details.

## cost

$0/month for 2 people (well within turso free tier)

## troubleshooting

**connection failed:**
- verify db_url and auth_token are set correctly
- check internet connection
- ensure turso database exists: `turso db list`

**messages not appearing:**
- both users must use the same db_url and auth_token
- check if message was saved: `turso db shell bt-phone-home`, then `SELECT * FROM messages;`

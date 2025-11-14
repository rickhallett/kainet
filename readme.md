# kaiNET - Encrypted Terminal Chat

A secure, room-based encrypted chat CLI application built with Go.

## Features

- AES-256-GCM end-to-end encryption
- Room-based isolation for private conversations
- Real-time message polling (1 second intervals)
- Turso libSQL cloud database backend
- Military-style boot sequence with green terminal aesthetics
- `/burn` command to wipe room history
- Cross-platform support (macOS, Linux, Windows)

## Usage

```bash
kainet <username> <room-name>
```

### Examples

```bash
# Join a room as "alice" in room "secretops"
./bin/kainet-darwin-arm64 alice secretops

# Join the same room as "bob"
./bin/kainet-darwin-arm64 bob secretops
```

### Commands

- Type any message and press Enter to send
- `/burn` - Wipe all messages in the current room
- `Ctrl+C` - Exit the application

## Building

### Install Dependencies

```bash
go mod download
```

### Build All Platforms

```bash
chmod +x build-all.sh
./build-all.sh
```

This creates binaries for:
- macOS ARM64 (`bin/kainet-darwin-arm64`)
- macOS AMD64 (`bin/kainet-darwin-amd64`)
- Linux AMD64 (`bin/kainet-linux-amd64`)
- Windows AMD64 (`bin/kainet.exe`)

### Build for Current Platform Only

```bash
go build -o kainet main.go
./kainet <username> <room-name>
```

## Security

- Messages are encrypted using AES-256-GCM before being stored
- Encryption key derived from auth token using SHA-256
- Each message uses a unique nonce for security
- Room-based isolation ensures conversation privacy

## Database Schema

```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_name TEXT NOT NULL,
    username TEXT NOT NULL,
    message TEXT NOT NULL,  -- encrypted
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Requirements

- Go 1.21 or higher
- Internet connection for Turso database access

## Architecture

- **Encryption**: AES-256-GCM with SHA-256 key derivation
- **Database**: Turso libSQL (serverless SQLite)
- **Real-time**: 1-second polling for new messages
- **UI**: Fatih Color library for terminal colors

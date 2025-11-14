#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_HOST="${KAINET_HOST:-https://github.com/rickhallett/kainet/releases/download/v1.1.0}"
BINARY_NAME="kainet"

# Parse arguments
USERNAME="$1"
ROOM_NAME="$2"

if [ -z "$USERNAME" ] || [ -z "$ROOM_NAME" ]; then
    echo -e "${RED}ERROR: username and room name required${NC}"
    echo "usage: curl -fsSL <install-url> | bash -s -- <username> <room-name>"
    echo ""
    echo "example:"
    echo "  curl -fsSL https://example.com/install.sh | bash -s -- alice secret-room"
    exit 1
fi

echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║           BT PHONE HOME - INSTALLER                      ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Darwin*)
        if [ "$ARCH" = "arm64" ]; then
            PLATFORM="darwin-arm64"
            echo -e "${CYAN}>>> detected: macOS Apple Silicon${NC}"
        else
            PLATFORM="darwin-amd64"
            echo -e "${CYAN}>>> detected: macOS Intel${NC}"
        fi
        ;;
    Linux*)
        PLATFORM="linux-amd64"
        echo -e "${CYAN}>>> detected: Linux${NC}"
        ;;
    MINGW*|MSYS*|CYGWIN*)
        echo -e "${RED}Windows detected - please use PowerShell installer${NC}"
        echo "Run this command in PowerShell instead:"
        echo ""
        echo "  iex \"& { \$(irm https://raw.githubusercontent.com/rickhallett/kainet/main/install.ps1) } $USERNAME $ROOM_NAME\""
        echo ""
        exit 1
        ;;
    *)
        echo -e "${RED}unsupported OS: $OS${NC}"
        echo "This installer supports macOS and Linux."
        echo "For Windows, use the PowerShell installer."
        exit 1
        ;;
esac

# Create temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo -e "${CYAN}>>> downloading binary for $PLATFORM...${NC}"

# Download binary
DOWNLOAD_URL="$BINARY_HOST/$BINARY_NAME-$PLATFORM"
if command -v curl &> /dev/null; then
    curl -fsSL -o "$BINARY_NAME" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -q -O "$BINARY_NAME" "$DOWNLOAD_URL"
else
    echo -e "${RED}ERROR: neither curl nor wget found${NC}"
    exit 1
fi

if [ ! -f "$BINARY_NAME" ]; then
    echo -e "${RED}ERROR: download failed${NC}"
    echo "Could not download from: $DOWNLOAD_URL"
    echo ""
    echo "Set custom host with: export BT_PHONE_HOME_HOST=https://your-host.com"
    exit 1
fi

echo -e "${GREEN}>>> download complete${NC}"

# Make executable
chmod +x "$BINARY_NAME"
echo -e "${GREEN}>>> made executable${NC}"

# Remove quarantine on macOS
if [ "$OS" = "Darwin" ]; then
    xattr -d com.apple.quarantine "$BINARY_NAME" 2>/dev/null || true
    echo -e "${GREEN}>>> removed quarantine flag${NC}"
fi

echo ""
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  INSTALLATION COMPLETE                 ║${NC}"
echo -e "${GREEN}║  LAUNCHING SECURE TERMINAL...          ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
echo ""

# Execute binary with username and room
# Redirect stdin from terminal to handle pipe-to-bash installation
exec "./$BINARY_NAME" "$USERNAME" "$ROOM_NAME" < /dev/tty

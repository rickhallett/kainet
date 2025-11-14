#!/bin/bash
# kaiNET Bootstrap Installer
# Downloads and runs installer, then auto-launches the terminal

set -e

# Parse arguments
USERNAME="$1"
ROOM_NAME="$2"

if [ -z "$USERNAME" ] || [ -z "$ROOM_NAME" ]; then
    echo "ERROR: username and room name required"
    echo "usage: bash <(curl -fsSL ...) <username> <room-name>"
    echo ""
    echo "example:"
    echo "  bash <(curl -fsSL https://kainet.dev/bootstrap.sh) alice secret-room"
    exit 1
fi

# Download installer to temp file
TEMP_INSTALLER=$(mktemp)
curl -fsSL https://raw.githubusercontent.com/rickhallett/kainet/main/install.sh > "$TEMP_INSTALLER"
chmod +x "$TEMP_INSTALLER"

# Run installer (it will show instructions)
"$TEMP_INSTALLER" "$USERNAME" "$ROOM_NAME"

# Clean up
rm -f "$TEMP_INSTALLER"

# Now launch the binary in a new shell with proper stdin
INSTALL_DIR="$HOME/.local/bin"
BINARY_PATH="$INSTALL_DIR/kainet"

if [ -f "$BINARY_PATH" ]; then
    # Launch in new shell to restore stdin
    exec bash -c "exec '$BINARY_PATH' '$USERNAME' '$ROOM_NAME' < /dev/tty"
else
    echo "ERROR: Installation failed - binary not found at $BINARY_PATH"
    exit 1
fi

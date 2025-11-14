#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "usage: ./run.sh <username> <room-name>"
    exit 1
fi

USERNAME="$1"
ROOM_NAME="$2"

# find the .mp4 file in current directory
MP4_FILE=$(ls bt-phone-home*.mp4 2>/dev/null | head -1)

if [ -z "$MP4_FILE" ]; then
    echo "error: no bt-phone-home*.mp4 file found in current directory"
    exit 1
fi

# remove .mp4 extension
BINARY_NAME="${MP4_FILE%.mp4}"
mv "$MP4_FILE" "$BINARY_NAME"

echo "renamed $MP4_FILE to $BINARY_NAME"

# make executable
chmod +x "$BINARY_NAME"
echo "made executable"

# remove quarantine
xattr -d com.apple.quarantine "$BINARY_NAME" 2>/dev/null
echo "removed quarantine"

# run with username and room
echo "starting chat as $USERNAME in room $ROOM_NAME..."
./"$BINARY_NAME" "$USERNAME" "$ROOM_NAME"

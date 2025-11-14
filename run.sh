#!/bin/bash

# kaiNET runner for WhatsApp .mp4 downloads
# This script renames, sets permissions, and executes the kaiNET binary
# Usage: ./run.sh <username> <room-name>

set -e

# Colors (green branding)
GREEN='\033[0;32m'
BOLD_GREEN='\033[1;32m'
NC='\033[0m' # No Color

print_green() {
    echo -e "${GREEN}$1${NC}"
}

print_bold_green() {
    echo -e "${BOLD_GREEN}$1${NC}"
}

# Main
main() {
    if [ "$#" -ne 2 ]; then
        echo "Usage: $0 <username> <room-name>"
        echo "Example: $0 john chatroom"
        exit 1
    fi

    USERNAME="$1"
    ROOM_NAME="$2"

    print_bold_green "=========================================="
    print_bold_green "   kaiNET Runner"
    print_bold_green "=========================================="
    echo ""

    # Find .mp4 file in current directory
    MP4_FILE=$(find . -maxdepth 1 -name "*.mp4" -type f | head -n 1)

    if [ -z "$MP4_FILE" ]; then
        echo "Error: No .mp4 file found in current directory"
        exit 1
    fi

    print_green "Found binary: $MP4_FILE"

    # Detect platform
    case "$(uname -s)" in
        Darwin*)
            BINARY_NAME="kainet"
            ;;
        Linux*)
            BINARY_NAME="kainet"
            ;;
        *)
            echo "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac

    # Rename file
    print_green "Renaming to: $BINARY_NAME"
    mv "$MP4_FILE" "$BINARY_NAME"

    # Set executable permissions
    print_green "Setting executable permissions..."
    chmod +x "$BINARY_NAME"

    # Remove macOS quarantine attribute if on macOS
    if [[ "$(uname -s)" == "Darwin" ]]; then
        print_green "Removing macOS quarantine attribute..."
        xattr -d com.apple.quarantine "$BINARY_NAME" 2>/dev/null || true
    fi

    echo ""
    print_bold_green "Starting kaiNET..."
    echo ""

    # Execute kaiNET
    ./"$BINARY_NAME" "$USERNAME" "$ROOM_NAME"
}

main "$@"

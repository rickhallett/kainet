#!/bin/bash

# kaiNET installer for macOS/Linux
# Usage: ./install.sh <username> <room-name>

set -e

# Colors (green branding)
GREEN='\033[0;32m'
BOLD_GREEN='\033[1;32m'
NC='\033[0m' # No Color

# Configuration
GITHUB_RELEASE_URL="${KAINET_HOST:-https://github.com/rickhallett/kainet/releases/latest/download}"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="kainet"

print_green() {
    echo -e "${GREEN}$1${NC}"
}

print_bold_green() {
    echo -e "${BOLD_GREEN}$1${NC}"
}

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""

    case "$(uname -s)" in
        Darwin*)
            os="darwin"
            ;;
        Linux*)
            os="linux"
            ;;
        *)
            echo "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            echo "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

# Main installation
main() {
    print_bold_green "=========================================="
    print_bold_green "   kaiNET Installer"
    print_bold_green "=========================================="
    echo ""

    if [ "$#" -ne 2 ]; then
        echo "Usage: $0 <username> <room-name>"
        echo "Example: $0 john chatroom"
        exit 1
    fi

    USERNAME="$1"
    ROOM_NAME="$2"

    # Detect platform
    PLATFORM=$(detect_platform)
    print_green "Detected platform: ${PLATFORM}"

    # Construct download URL
    BINARY_URL="${GITHUB_RELEASE_URL}/kainet-${PLATFORM}"
    print_green "Download URL: ${BINARY_URL}"

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Download binary
    print_green "Downloading kaiNET..."
    TEMP_FILE=$(mktemp)

    if command -v curl &> /dev/null; then
        curl -L -o "${TEMP_FILE}" "${BINARY_URL}"
    elif command -v wget &> /dev/null; then
        wget -O "${TEMP_FILE}" "${BINARY_URL}"
    else
        echo "Error: Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    # Move to install directory
    INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"
    mv "${TEMP_FILE}" "${INSTALL_PATH}"

    # Set executable permissions
    chmod +x "${INSTALL_PATH}"
    print_green "Binary installed to: ${INSTALL_PATH}"

    # Remove macOS quarantine attribute if on macOS
    if [[ "$(uname -s)" == "Darwin" ]]; then
        print_green "Removing macOS quarantine attribute..."
        xattr -d com.apple.quarantine "${INSTALL_PATH}" 2>/dev/null || true
    fi

    echo ""
    print_bold_green "Installation complete!"
    print_green "Starting kaiNET..."
    echo ""

    # Execute kaiNET
    "${INSTALL_PATH}" "${USERNAME}" "${ROOM_NAME}"
}

main "$@"

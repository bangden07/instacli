#!/bin/bash
# InstaCli Installer Script
# Usage: curl -sSL https://raw.githubusercontent.com/bangden07/instacli/main/install.sh | bash

set -e

REPO="bangden07/instacli"
VERSION="latest"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║       InstaCli Installer v1.0            ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════╝${NC}"
echo ""

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    case "$OS" in
        linux)
            PLATFORM="linux-${ARCH}"
            ;;
        darwin)
            PLATFORM="darwin-${ARCH}"
            ;;
        mingw*|msys*|cygwin*)
            PLATFORM="windows-${ARCH}.exe"
            ;;
        *)
            echo -e "${RED}Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${YELLOW}Detected platform: ${PLATFORM}${NC}"
}

# Download and install
install_instacli() {
    BINARY_NAME="instacli-${PLATFORM}"
    
    if [ "$VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"
    else
        DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
    fi
    
    echo -e "${YELLOW}Downloading InstaCli...${NC}"
    echo "URL: $DOWNLOAD_URL"
    echo ""
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    TMP_FILE="${TMP_DIR}/instacli"
    
    # Download
    if command -v curl &> /dev/null; then
        curl -sSL "$DOWNLOAD_URL" -o "$TMP_FILE"
    elif command -v wget &> /dev/null; then
        wget -q "$DOWNLOAD_URL" -O "$TMP_FILE"
    else
        echo -e "${RED}Error: curl or wget is required${NC}"
        exit 1
    fi
    
    # Make executable
    chmod +x "$TMP_FILE"
    
    # Install to system directory
    echo -e "${YELLOW}Installing to ${INSTALL_DIR}...${NC}"
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "${INSTALL_DIR}/instacli"
    else
        sudo mv "$TMP_FILE" "${INSTALL_DIR}/instacli"
    fi
    
    # Cleanup
    rm -rf "$TMP_DIR"
    
    echo ""
    echo -e "${GREEN}✅ InstaCli installed successfully!${NC}"
    echo ""
    echo "Run 'instacli' to start the TUI installer."
    echo ""
}

# Main
detect_platform
install_instacli

# Verify installation
if command -v instacli &> /dev/null; then
    echo -e "${GREEN}Verification: instacli is ready to use!${NC}"
else
    echo -e "${YELLOW}Note: You may need to restart your terminal or add ${INSTALL_DIR} to your PATH.${NC}"
fi

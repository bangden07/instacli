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
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_banner() {
    echo ""
    echo -e "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║                                                          ║${NC}"
    echo -e "${CYAN}║   ██╗███╗   ██╗███████╗████████╗ █████╗  ██████╗██╗      ║${NC}"
    echo -e "${CYAN}║   ██║████╗  ██║██╔════╝╚══██╔══╝██╔══██╗██╔════╝██║      ║${NC}"
    echo -e "${CYAN}║   ██║██╔██╗ ██║███████╗   ██║   ███████║██║     ██║      ║${NC}"
    echo -e "${CYAN}║   ██║██║╚██╗██║╚════██║   ██║   ██╔══██║██║     ██║      ║${NC}"
    echo -e "${CYAN}║   ██║██║ ╚████║███████║   ██║   ██║  ██║╚██████╗███████╗ ║${NC}"
    echo -e "${CYAN}║   ╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝╚══════╝ ║${NC}"
    echo -e "${CYAN}║                                                          ║${NC}"
    echo -e "${CYAN}║          🚀 Universal Server Installer v1.2.0            ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

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
            echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
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
            echo -e "${RED}❌ Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${BLUE}📦 Detected platform: ${PLATFORM}${NC}"
}

# Check if Go is installed
check_go() {
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        echo -e "${GREEN}✅ Go found: ${GO_VERSION}${NC}"
        return 0
    else
        return 1
    fi
}

# Install from source
install_from_source() {
    echo -e "${YELLOW}📥 Installing from source...${NC}"
    echo ""
    
    # Check Go
    if ! check_go; then
        echo -e "${YELLOW}⚠️  Go not found. Installing Go first...${NC}"
        install_go
    fi
    
    # Clone and build
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    echo -e "${BLUE}📦 Cloning repository...${NC}"
    git clone --depth 1 https://github.com/${REPO}.git
    cd instacli
    
    echo -e "${BLUE}🔨 Building InstaCli...${NC}"
    go build -o instacli ./cmd/instacli
    
    # Install
    echo -e "${BLUE}📁 Installing to ${INSTALL_DIR}...${NC}"
    if [ -w "$INSTALL_DIR" ]; then
        mv instacli "${INSTALL_DIR}/instacli"
    else
        sudo mv instacli "${INSTALL_DIR}/instacli"
    fi
    
    # Cleanup
    cd /
    rm -rf "$TMP_DIR"
    
    echo ""
    echo -e "${GREEN}✅ InstaCli installed successfully from source!${NC}"
}

# Install Go
install_go() {
    echo -e "${YELLOW}📥 Installing Go...${NC}"
    
    GO_VERSION="1.21.5"
    
    case "$PLATFORM" in
        linux-amd64)
            GO_ARCHIVE="go${GO_VERSION}.linux-amd64.tar.gz"
            ;;
        linux-arm64)
            GO_ARCHIVE="go${GO_VERSION}.linux-arm64.tar.gz"
            ;;
        darwin-amd64)
            GO_ARCHIVE="go${GO_VERSION}.darwin-amd64.tar.gz"
            ;;
        darwin-arm64)
            GO_ARCHIVE="go${GO_VERSION}.darwin-arm64.tar.gz"
            ;;
        *)
            echo -e "${RED}❌ Cannot install Go automatically for ${PLATFORM}${NC}"
            echo "Please install Go manually: https://go.dev/dl/"
            exit 1
            ;;
    esac
    
    curl -sSL "https://go.dev/dl/${GO_ARCHIVE}" -o "/tmp/${GO_ARCHIVE}"
    sudo tar -C /usr/local -xzf "/tmp/${GO_ARCHIVE}"
    rm "/tmp/${GO_ARCHIVE}"
    
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    
    echo -e "${GREEN}✅ Go installed: $(go version)${NC}"
}

# Try to download from releases
install_from_release() {
    BINARY_NAME="instacli-${PLATFORM}"
    
    if [ "$VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"
    else
        DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
    fi
    
    echo -e "${YELLOW}📥 Downloading InstaCli...${NC}"
    echo -e "${BLUE}   URL: $DOWNLOAD_URL${NC}"
    echo ""
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    TMP_FILE="${TMP_DIR}/instacli"
    
    # Try to download
    if command -v curl &> /dev/null; then
        if curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE" 2>/dev/null; then
            # Download successful
            chmod +x "$TMP_FILE"
            
            echo -e "${BLUE}📁 Installing to ${INSTALL_DIR}...${NC}"
            
            if [ -w "$INSTALL_DIR" ]; then
                mv "$TMP_FILE" "${INSTALL_DIR}/instacli"
            else
                sudo mv "$TMP_FILE" "${INSTALL_DIR}/instacli"
            fi
            
            rm -rf "$TMP_DIR"
            
            echo ""
            echo -e "${GREEN}✅ InstaCli installed successfully!${NC}"
            return 0
        fi
    fi
    
    rm -rf "$TMP_DIR"
    return 1
}

# Print features
print_features() {
    echo ""
    echo -e "${CYAN}✨ Features:${NC}"
    echo -e "   ${GREEN}•${NC} 28+ One-click installers"
    echo -e "   ${GREEN}•${NC} Beautiful TUI interface"
    echo -e "   ${GREEN}•${NC} Clone & Setup repositories (NEW!)"
    echo -e "   ${GREEN}•${NC} SSH remote installation"
    echo -e "   ${GREEN}•${NC} Auto-detect project types"
    echo ""
    echo -e "${CYAN}📚 Categories:${NC}"
    echo -e "   Web Servers, Runtimes, Databases, Containers"
    echo -e "   Monitoring, Infrastructure, CI/CD, Security, CMS"
    echo ""
}

# Main
main() {
    print_banner
    detect_platform
    echo ""
    
    # Try release first, fallback to source
    if ! install_from_release; then
        echo -e "${YELLOW}⚠️  Binary release not available, building from source...${NC}"
        echo ""
        
        # Check for git
        if ! command -v git &> /dev/null; then
            echo -e "${RED}❌ Git is required to build from source${NC}"
            echo "Please install git: sudo apt install git"
            exit 1
        fi
        
        install_from_source
    fi
    
    print_features
    
    # Verify installation
    if command -v instacli &> /dev/null; then
        INSTALLED_VERSION=$(instacli --version 2>/dev/null || echo "unknown")
        echo -e "${GREEN}🎉 InstaCli is ready!${NC}"
        echo -e "   Version: ${INSTALLED_VERSION}"
        echo ""
        echo -e "   Run ${CYAN}instacli${NC} to start the TUI installer."
    else
        echo -e "${YELLOW}⚠️  Note: Restart your terminal or run:${NC}"
        echo -e "   ${CYAN}export PATH=\$PATH:${INSTALL_DIR}${NC}"
    fi
    
    echo ""
    echo -e "${BLUE}📖 Documentation: https://github.com/bangden07/instacli/wiki${NC}"
    echo ""
}

main

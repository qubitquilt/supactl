#!/bin/bash

# supactl installation script
# This script downloads and installs the latest version of supactl

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="yourusername/supactl"
BINARY_NAME="supactl"
INSTALL_DIR="/usr/local/bin"

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    print_info "Detected platform: $PLATFORM"
}

# Get latest release version
get_latest_version() {
    print_info "Fetching latest release version..."

    LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$LATEST_RELEASE" ]; then
        print_error "Failed to fetch latest release version"
        exit 1
    fi

    print_info "Latest version: $LATEST_RELEASE"
}

# Download binary
download_binary() {
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/${BINARY_NAME}-${PLATFORM}"

    print_info "Downloading from: $DOWNLOAD_URL"

    TMP_FILE=$(mktemp)

    if ! curl -L -f -o "$TMP_FILE" "$DOWNLOAD_URL"; then
        print_error "Failed to download binary"
        rm -f "$TMP_FILE"
        exit 1
    fi

    print_success "Download completed"
    echo "$TMP_FILE"
}

# Install binary
install_binary() {
    local tmp_file=$1

    print_info "Installing to $INSTALL_DIR/$BINARY_NAME..."

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_file" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        print_warning "Need sudo privileges to install to $INSTALL_DIR"
        sudo mv "$tmp_file" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi

    print_success "Installed $BINARY_NAME to $INSTALL_DIR"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        VERSION=$($BINARY_NAME --version)
        print_success "Installation verified: $VERSION"
        return 0
    else
        print_error "Installation verification failed"
        print_info "Make sure $INSTALL_DIR is in your PATH"
        return 1
    fi
}

# Main installation process
main() {
    echo ""
    echo "╔══════════════════════════════════════════════╗"
    echo "║     supactl Installation Script              ║"
    echo "╚══════════════════════════════════════════════╝"
    echo ""

    # Check dependencies
    if ! command -v curl &> /dev/null; then
        print_error "curl is required but not installed"
        exit 1
    fi

    detect_platform
    get_latest_version

    TMP_FILE=$(download_binary)
    install_binary "$TMP_FILE"

    echo ""
    if verify_installation; then
        echo ""
        print_success "Installation complete! 🎉"
        echo ""
        echo "Get started with:"
        echo "  supactl login https://your-supacontrol-server.com"
        echo "  supactl --help"
        echo ""
    else
        print_warning "Installation completed but verification failed"
        echo ""
        echo "Try running: export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
        exit 1
    fi
}

# Run main function
main

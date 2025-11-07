#!/bin/bash

# supactl uninstallation script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BINARY_NAME="supactl"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.supacontrol"

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

confirm() {
    read -p "$1 [y/N] " -n 1 -r
    echo
    [[ $REPLY =~ ^[Yy]$ ]]
}

main() {
    echo ""
    echo "╔══════════════════════════════════════════════╗"
    echo "║     supactl Uninstallation Script           ║"
    echo "╚══════════════════════════════════════════════╝"
    echo ""

    # Check if binary exists
    if [ ! -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        print_warning "$BINARY_NAME not found in $INSTALL_DIR"
        print_info "It may have been installed elsewhere or already removed"
    else
        print_info "Found $BINARY_NAME at $INSTALL_DIR/$BINARY_NAME"

        if confirm "Remove $BINARY_NAME binary?"; then
            if [ -w "$INSTALL_DIR" ]; then
                rm "$INSTALL_DIR/$BINARY_NAME"
            else
                sudo rm "$INSTALL_DIR/$BINARY_NAME"
            fi
            print_success "Binary removed"
        else
            print_info "Skipping binary removal"
        fi
    fi

    # Check for config directory
    if [ -d "$CONFIG_DIR" ]; then
        print_warning "Found configuration directory: $CONFIG_DIR"
        print_info "This contains your authentication credentials"

        if confirm "Remove configuration directory?"; then
            rm -rf "$CONFIG_DIR"
            print_success "Configuration removed"
        else
            print_info "Keeping configuration directory"
        fi
    fi

    echo ""
    print_success "Uninstallation complete"
    echo ""
}

main

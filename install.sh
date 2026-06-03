#!/usr/bin/env bash
set -euo pipefail

BLUE='\033[0;34m'
GREEN='\033[0;32m'
BOLD='\033[1m'
NC='\033[0m'

msg() { echo -e "${BLUE}${BOLD}::${NC} ${BOLD}$1${NC}"; }
error() { echo -e "\033[0;31m${BOLD}Error:${NC} $1"; exit 1; }

msg "Starting installation."
msg "Installing dependencies..."

if [ -f /etc/arch-release ]; then
    sudo pacman -Sy --noconfirm --needed go git imagemagick playerctl
elif [ -f /etc/debian_version ]; then
    sudo apt update && sudo apt install -y golang-go git imagemagick playerctl
elif [ -f /etc/fedora-release ]; then
    sudo dnf install -y golang git ImageMagick playerctl
elif command -v xbps-install &> /dev/null; then
    sudo xbps-install -Sy go git ImageMagick playerctl
else
    error "Unsupported distribution. Install dependencies manually."
fi

INSTALL_DIR="$HOME/Virga-player"

if [[ -d ".git" && $(basename "$(pwd)") == "Virga-player" ]]; then
    TARGET_DIR="$(pwd)"
    msg "Running from current directory: $TARGET_DIR"
else
    if [[ ! -d "$INSTALL_DIR" ]]; then
        msg "Cloning into $INSTALL_DIR..."
        git clone https://github.com/Glebsolopdf/Virga-player.git "$INSTALL_DIR"
    fi
    TARGET_DIR="$INSTALL_DIR"
    cd "$TARGET_DIR"
fi

if [[ -f "$TARGET_DIR/build.sh" ]]; then
    echo -n "Start compiling? (y/n): "
    read -r REPLY
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        msg "Building from $TARGET_DIR/build.sh..."
        chmod +x "$TARGET_DIR/build.sh"
        "$TARGET_DIR/build.sh"
    else
        msg "Compilation skipped by user."
    fi
else
    error "Security check failed: build.sh not found in $TARGET_DIR"
fi

echo -e "\n${GREEN}${BOLD}Done.${NC}"
echo -n "Run Virga-player now? (y/n): "
read -r REPLY

if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [[ -x "$TARGET_DIR/bin/virga-player" ]]; then
        "$TARGET_DIR/bin/virga-player"
    else
        error "Binary not found at $TARGET_DIR/bin/virga-player. Build it first."
    fi
fi

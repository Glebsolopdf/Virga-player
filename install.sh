#!/usr/bin/env bash
set -euo pipefail
[[ "$OSTYPE" == "linux-gnu"* ]] || { echo "Error: Only Linux is supported"; exit 1; }
command -v pacman >/dev/null 2>&1 || { echo "Error: This script requires pacman (Arch Linux)"; exit 1; }
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$ROOT_DIR/src"
BIN_NAME="virga"
INSTALL_DEST="/usr/local/bin/$BIN_NAME"
echo "==> Installing dependencies..."
pacman -Syu --noconfirm --needed go git imagemagick
if [[ ! -d "$SRC_DIR" ]]; then
    echo "Error: Source directory $SRC_DIR not found."
    exit 1
fi
echo "==> Building $BIN_NAME..."
go build -C "$SRC_DIR" -ldflags="-s -w" -o "$ROOT_DIR/$BIN_NAME" .
echo "==> Installing to $INSTALL_DEST..."
install -Dm755 "$ROOT_DIR/$BIN_NAME" "$INSTALL_DEST"
echo "==> Done. Run with: $BIN_NAME"
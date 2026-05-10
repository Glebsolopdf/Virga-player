#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if ! command -v pacman >/dev/null 2>&1; then
  echo "Error: pacman not found. This installer only supports Arch Linux / Arch based systems."
  exit 1
fi

if [[ $(id -u) -ne 0 ]]; then
  if ! command -v sudo >/dev/null 2>&1; then
    echo "Error: run this script as root or install sudo."
    exit 1
  fi
  SUDO="sudo"
else
  SUDO=""
fi

PACKAGES=(
  go
  git
  imagemagick
)

echo "Installing required packages: ${PACKAGES[*]}"
$SUDO pacman -Syu --noconfirm --needed "${PACKAGES[@]}"

echo "Dependencies installed."

if ! command -v go >/dev/null 2>&1; then
  echo "Error: Go installation failed or Go is not on PATH."
  exit 1
fi

echo "Go version: $(go version)"

echo "Installing Go module dependencies..."
cd "$ROOT_DIR/src"
go mod download

echo "Building virga-player..."
cd "$ROOT_DIR/src"
go build -ldflags="-s -w" -o "$ROOT_DIR/virga-player" .

echo "Installing /usr/bin/virga..."
$SUDO install -Dm755 "$ROOT_DIR/virga-player" /usr/bin/virga

echo "Installation complete. You can now run virga from the terminal."

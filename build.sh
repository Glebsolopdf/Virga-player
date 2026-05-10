#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$ROOT_DIR/src"
OUTPUT="$ROOT_DIR/virga-player"

if [[ ! -d "$SRC_DIR" ]]; then
    echo "Error: Directory $SRC_DIR not found."
    exit 1
fi

echo "Building virga-player..."
cd "$SRC_DIR"
go build -ldflags="-s -w" -o "$OUTPUT" .

echo "Success! Binary located at: $OUTPUT"
#!/bin/bash

APP_NAME="cli-app"
OUTPUT_DIR="./bin"


GOOS="$(uname -s | tr '[:upper:]' '[:lower:]')"
GOARCH="$(uname -m)"
case "$GOARCH" in
    x86_64) GOARCH="amd64" ;;
    aarch64) GOARCH="arm64" ;;
esac

BINARY="$OUTPUT_DIR/$APP_NAME-$GOOS-$GOARCH"

if [ ! -f "$BINARY" ]; then
    echo "Error: Binary not found. Run ./build.sh first."
    exit 1
fi

"$BINARY" "$@"
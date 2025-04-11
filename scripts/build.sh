#!/bin/bash

APP_NAME="cli-app"
OUTPUT_DIR="./bin"

GOOS="$(uname -s | tr '[:upper:]' '[:lower:]')"  # linux/darwin
GOARCH="$(uname -m)"                              # x86_64/arm64


case "$GOARCH" in
    x86_64) GOARCH="amd64" ;;
    aarch64) GOARCH="arm64" ;;
esac

mkdir -p "$OUTPUT_DIR"

echo "Building for $GOOS/$GOARCH..."
GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT_DIR/$APP_NAME-$GOOS-$GOARCH" ./cmd/cli

if [ -f "$OUTPUT_DIR/$APP_NAME-$GOOS-$GOARCH" ]; then
    echo "Build successful. Binary: $OUTPUT_DIR/$APP_NAME-$GOOS-$GOARCH"
else
    echo "Build failed!"
    exit 1
fi
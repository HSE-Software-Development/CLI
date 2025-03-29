#!/bin/bash

APP_NAME="cli-app"
GOOS="linux"
GOARCH="amd64"
OUTPUT_DIR="./bin"

mkdir -p $OUTPUT_DIR

echo "build for $GOOS/$GOARCH..."
GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT_DIR/$APP_NAME ./cmd/cli

if [ -f "$OUTPUT_DIR/$APP_NAME" ]; then
    echo "Build over. Path: $OUTPUT_DIR/$APP_NAME"
else
    echo "Building error"
    exit 1
fi
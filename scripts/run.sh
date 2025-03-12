#!/bin/bash

BIN_PATH="./bin/cli-app"

if [ ! -f "$BIN_PATH" ]; then
    echo "File not found. Build app (scripts/build.sh)."
    exit 1
fi

echo "Running..."
$BIN_PATH
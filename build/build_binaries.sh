#!/bin/sh

BUILD_DIR="${BUILD_DIR:-./dist_build}"
BIN_DIR="${BIN_DIR:-$BUILD_DIR/bin}"

if [ ! -d "$BIN_DIR" ]; then
	echo "create BIN_DIR=$BIN_DIR"
	mkdir -p "$BIN_DIR"
fi

EXAMPLE_BINARIES="issuelist grouplist"

for binary in $EXAMPLE_BINARIES; do
	go build -mod=vendor -o "$BIN_DIR/$binary" "example/$binary/main.go"
done

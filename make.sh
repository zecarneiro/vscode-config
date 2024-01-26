#!/bin/bash

declare DIRNAME=$(dirname "$0")
declare ROOT=$(realpath "$DIRNAME")
declare DOWNLOAD_DIR="$ROOT/download"
declare BINARY_DIR="$ROOT/bin"
declare BINARY="$BINARY_DIR/main.exe"

_check_go() {
    if [ ! "$(command -v go)" ]; then
        echo "[ERROR] Please install golang!"
        exit 1
    fi
    
}
_clean() {
    rm -rf "$BINARY_DIR"
    rm -rf "$DOWNLOAD_DIR"
}

_create_dir() {
    mkdir -p "$BINARY_DIR"
    mkdir -p "$DOWNLOAD_DIR"
}


run() {
    _clean
    _create_dir
    build
    eval "$BINARY $*"
}

build() {
    _clean
    _create_dir
    _check_go
    go build -o "$BINARY" main.go
}

pushd .
cd $ROOT
eval "$*"
#popd

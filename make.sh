#!/bin/bash

declare DIRNAME=$(dirname "$0")
declare ROOT=$(realpath "$DIRNAME")
declare DOWNLOAD_DIR="$ROOT/vscode-config-download"
declare DEPENDENCIES_DIR="$ROOT/dependencies"
declare BINARY_DIR="$ROOT/bin"
declare BINARY="$BINARY_DIR/vscode-config-linux"
declare SCRIPT_BASH_UTILS_DIR="$DEPENDENCIES_DIR/bash-utils-2.2.1"
declare SCRIPT_BASH_UTILS_MAIN="$SCRIPT_BASH_UTILS_DIR/main-utils.sh"

# ---------------------------------------------------------------------------- #
#                               OTHERS FUNCTIONS                               #
# ---------------------------------------------------------------------------- #
function _download() {
    local url="$1"
    local outFile="$2"
    if [ ! "$(command -v wget)" ]; then
        echo "[ERROR] Please install wget!"
        exit 1
    fi
    echo "Downloading File \"${outFile}\" from \"${url}\" ......."
    wget -O "$outFile" "$url" -q
    echo "Successfully Downloaded File \"${outFile}\" from \"${url}\""
}

# ---------------------------------------------------------------------------- #
#                                  OPERATIONS                                  #
# ---------------------------------------------------------------------------- #
_clean() {
    . "$SCRIPT_BASH_UTILS_MAIN"
    deletedirectory "$BINARY_DIR"
    deletedirectory "$DOWNLOAD_DIR"
}

_build() {
    local binaryWindows="$BINARY_DIR/vscode-config-win.exe"
	local binaryLinux="$BINARY_DIR/vscode-config-linux"
    local sourceDir="$ROOT/src"
    . "$SCRIPT_BASH_UTILS_MAIN"
	infolog "Build WINDOWS app..."
	export GOOS=windows
	export GOARCH=amd64
	go build -o "$binaryWindows" "$sourceDir"

	infolog "Build LINUX app..."
	export GOOS=linux
	export GOARCH=amd64
	go build -o "$binaryLinux" "$sourceDir"
}

_run() {
    . "$SCRIPT_BASH_UTILS_MAIN"
    if [ "$(confirm "Do you need reset vscode")" == "true" ]; then
        local configDir="${HOME}/.config/Code"
        local extensionsDir="${HOME}/.vscode"
        infolog "Delete: $configDir"
        deletedirectory "$configDir"
        infolog "Delete: $extensionsDir"
        deletedirectory "$extensionsDir"
    fi
    eval "$BINARY $*"
}

# ---------------------------------------------------------------------------- #
#                                     MAIN                                     #
# ---------------------------------------------------------------------------- #
_process_dependencies() {
    if [ ! -d "$SCRIPT_BASH_UTILS_DIR" ]; then
        mkdir -p "$SCRIPT_BASH_UTILS_DIR"
        _download "https://github.com/zecarneiro/bash-utils/archive/refs/tags/v2.2.1.zip" "$DEPENDENCIES_DIR/bash-utils.zip"
        unzip -q "$DEPENDENCIES_DIR/bash-utils.zip" -d "$DEPENDENCIES_DIR"
    fi
}

_check_vendor_dependencies() {
    if [ ! "$(command -v go)" ]; then
        echo "[ERROR] Please install golang!"
        exit 1
    fi
    if [ ! "$(command -v wget)" ]; then
        echo "[ERROR] Please install wget!"
        exit 1
    fi
    if [ ! "$(command -v unzip)" ]; then
        echo "[ERROR] Please install unzip!"
        exit 1
    fi
}

main() {
    _check_vendor_dependencies
    _process_dependencies
    case "${1}" in
        --run) _run "$2" ;;
        --build)
            _clean
            _build
        ;;
        --clean) _clean ;;
        --clean-with-dependencies)
            _clean
            if [ -d "$DEPENDENCIES_DIR" ]; then
                rm -rf "$DEPENDENCIES_DIR"
            fi
        ;;
        *) Write-Output "$SCRIPT_NAME --[run|build|clean|clean-with-dependencies] ARG_JSON_CONFIG" ;;
    esac
}

pushd .
cd $ROOT
main $@
popd

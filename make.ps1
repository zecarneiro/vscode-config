$ROOT=($PSScriptRoot)
$CURRENT_DIR="$pwd"
$DOWNLOAD_DIR="$ROOT\download"
$BINARY_DIR="$ROOT\bin"
$BINARY="$BINARY_DIR\main.exe"
$SCRIPT_NAME=([System.IO.Path]::GetFileName("$ROOT\make.ps1"))

function _check_go() {
    if ([string]::IsNullOrEmpty((which go.exe))) {
        Write-Output "[ERROR] Please install golang!"
        exit 1
    }    
}

function _run() {
    param([string] $jsonConfig)
    _build
    Invoke-Expression "$BINARY $jsonConfig"
}

function _build() {
    _check_go
    go build -o "$BINARY" main.go
}

function _clean {
    if ((Test-Path -Path "$BINARY_DIR")) {
        Remove-Item "$BINARY_DIR" -Force -Recurse
    }
    if ((Test-Path -Path "$DOWNLOAD_DIR")) {
        Remove-Item "$DOWNLOAD_DIR" -Force -Recurse
    }
}

function _create_dir {
    if (!(Test-Path -Path "$BINARY_DIR")) {
        mkdir "$BINARY_DIR"
    }
    if (!(Test-Path -Path "$DOWNLOAD_DIR")) {
        mkdir "$DOWNLOAD_DIR"
    }
}

function _start {
    param([string[]] $arguments)
    _create_dir
    switch ($arguments[0]) {
        run {
            _run $arguments[1]
        }
        build {
            _build
        }
        clean {
            _clean
        }
        Default {
            Write-Output "$SCRIPT_NAME run|build ARG_JSON_CONFIG"
        }
    }
}
Set-Location "$ROOT"
echo "$CURRENT_DIR"
_start $args
Set-Location "$CURRENT_DIR"


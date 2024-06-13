$ROOT=($PSScriptRoot)
$CURRENT_DIR="$pwd"
$DOWNLOAD_DIR="$ROOT\vscode-config-download"
$DEPENDENCIES_DIR="$ROOT\dependencies"
$SCRIPT_POWERSHELL_UTILS_DIR="$DEPENDENCIES_DIR\powershell-utils-1.3.0"
$SCRIPT_POWERSHELL_UTILS_MAIN="$SCRIPT_POWERSHELL_UTILS_DIR\MainUtils.ps1"
$BINARY_DIR="$ROOT\bin"
$BINARY="$BINARY_DIR\vscode-config-win.exe"
$SCRIPT_NAME=([System.IO.Path]::GetFileName("$ROOT\make.ps1"))

# ---------------------------------------------------------------------------- #
#                               OTHERS FUNCTIONS                               #
# ---------------------------------------------------------------------------- #
function _download($url, $outFile) {
    $WebClient = New-Object System.Net.WebClient
    Write-Output "Downloading File `"${outFile}`" from `"${url}`" ......."
    $WebClient.DownloadFile("$url","$outFile")
    Write-Output "Successfully Downloaded File `"${outFile}`" from `"${url}`""
}

# ---------------------------------------------------------------------------- #
#                                  OPERATIONS                                  #
# ---------------------------------------------------------------------------- #
function _clean($force) {
    . $SCRIPT_POWERSHELL_UTILS_MAIN
    deletedirectory "$BINARY_DIR"
    deletedirectory "$DOWNLOAD_DIR"
}

function _build() {
	$binaryWindows="$BINARY_DIR\vscode-config-win.exe"
	$binaryLinux="$BINARY_DIR\vscode-config-linux"
    $sourceDir = "$ROOT\src"
    . $SCRIPT_POWERSHELL_UTILS_MAIN
	infolog "Build WINDOWS app..."
	export GOOS=windows
	export GOARCH=amd64
	go build -o "$binaryWindows" "$sourceDir"

	infolog "Build LINUX app..."
	export GOOS=linux
	export GOARCH=amd64
	go build -o "$binaryLinux" "$sourceDir"
}

function _run() {
    param([string] $jsonConfig)
    . $SCRIPT_POWERSHELL_UTILS_MAIN
    if ((confirm "Do you need reset vscode")) {
        $configDir = "${home}\AppData\Roaming\Code"
        $extensionsDir = "${home}\.vscode"
        infolog "Delete: $configDir"
        deletedirectory "$configDir"
        infolog "Delete: $extensionsDir"
        deletedirectory "$extensionsDir"
    }
    Invoke-Expression "$BINARY $jsonConfig"
}

# ---------------------------------------------------------------------------- #
#                                     MAIN                                     #
# ---------------------------------------------------------------------------- #
function _process_dependencies {
    if (!(Test-Path -Path "$SCRIPT_POWERSHELL_UTILS_DIR")) {
        New-Item -ItemType Directory "$SCRIPT_POWERSHELL_UTILS_DIR" | Out-Null
        _download "https://github.com/zecarneiro/powershell-utils/archive/refs/tags/v1.3.0.zip" "$DEPENDENCIES_DIR\powershell-utils.zip"
        Expand-Archive -Path "$DEPENDENCIES_DIR\powershell-utils.zip" -DestinationPath "$DEPENDENCIES_DIR"
    }
}

function _check_vendor_dependencies() {
    if ([string]::IsNullOrEmpty((which go.exe))) {
        Write-Output "[ERROR] Please install golang!"
        exit 1
    }    
}

function main() {
    param([string[]] $arguments)
    _check_vendor_dependencies
    _process_dependencies
    switch ($arguments[0]) {
        --run {
            _run $arguments[1]
        }
        --build {
            _clean
            _build
        }
        --clean {
            _clean
        }
        --clean-with-dependencies {
            _clean
            if ((Test-Path -Path "$DEPENDENCIES_DIR")) {
                Remove-Item "$DEPENDENCIES_DIR" -Force -Recurse | Out-Null
            }
        }
        Default {
            Write-Output "$SCRIPT_NAME --[run|build|clean|clean-with-dependencies] ARG_JSON_CONFIG"
        }
    }
}
Set-Location "$ROOT"
Write-Output "$CURRENT_DIR"
main $args
Set-Location "$CURRENT_DIR"
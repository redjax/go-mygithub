#!/bin/bash

set -euo pipefail

if ! command -v go > /dev/null; then
    echo "Go is not installed"
    exit 1  
fi

## Default values
BIN_NAME="mygithub"
BUILD_OS="linux"
BUILD_ARCH="amd64"
BUILD_OUTPUT_DIR="dist/"
BUILD_TARGET="./main.go"

function usage {
    cat <<EOF
Usage: $0 [options]

Options:
    --bin-name NAME Name of the executable (default: $BIN_NAME)
    --build-os OS Build OS (default: $BUILD_OS)
    --build-arch ARCH Build architecture (default: $BUILD_ARCH)
    --build-output-dir DIR Build output directory (default: $BUILD_OUTPUT_DIR)
    --build-target TARGET Build target (default: $BUILD_TARGET)

    -h, --help Showw this help message
EOF
}

## Parse CLI args
while [[ $# -gt 0 ]]; do
    case "$1" in
        --bin-name)
            BIN_NAME="$2"
            shift 2
            ;;
        --build-os)
            BUILD_OS="$2"
            shift 2
            ;;
        --build-arch)
            BUILD_ARCH="$2"
            shift 2
            ;;
        --build-output-dir)
            BUILD_OUTPUT_DIR="$2"
            shift 2
            ;;
        --build-target)
            BUILD_TARGET="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

## Validate BIN_NAME
if [[ -z "$BIN_NAME" ]]; then
    echo "BIN_NAME is required"
    usage
    exit 1
fi

## Append .exe for Windows builds, if not present in build name
if [[ "$BUILD_OS" == "windows" && "$BIN_NAME" != *.exe ]]; then
    BIN_NAME="${BIN_NAME}.exe"
fi

## Create output directory if it doesn't exist
mkdir -pv "$BUILD_OUTPUT_DIR"

echo "Building $BIN_NAME for $BUILD_OS/$BUILD_ARCH"
echo ""
echo "--[ Build start"

export GOOS="$BUILD_OS"
export GOARCH="$BUILD_ARCH"

if go build -o "$BUILD_OUTPUT_DIR/$BIN_NAME" "$BUILD_TARGET"; then
    echo "Build successful"
else
    echo "Error building $BIN_NAME"

    echo ""
    echo "--[ Build complete"
    exit 1
fi

echo ""
echo "--[ Build complete"
exit 0

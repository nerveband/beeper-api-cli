#!/bin/bash
set -e

VERSION=${VERSION:-dev}
OUTPUT_DIR="dist"

echo "Building Beeper API CLI ${VERSION}..."

mkdir -p ${OUTPUT_DIR}

# Build for common platforms
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
    OS=$(echo $platform | cut -d'/' -f1)
    ARCH=$(echo $platform | cut -d'/' -f2)
    OUTPUT_NAME="beeper-api-${OS}-${ARCH}"
    
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "Building for ${OS}/${ARCH}..."
    GOOS=$OS GOARCH=$ARCH go build -o "${OUTPUT_DIR}/${OUTPUT_NAME}" \
        -ldflags "-X github.com/nerveband/beeper-api-cli/cmd.Version=${VERSION}" \
        .
done

echo "Build complete! Binaries in ${OUTPUT_DIR}/"
ls -lh ${OUTPUT_DIR}/

#!/bin/bash
# Build script for InstaCli - Cross-platform compilation
# This creates binaries for Linux, macOS, and Windows

set -e

VERSION="1.0.0"
BUILD_DIR="./dist"
APP_NAME="instacli"

echo "🔨 Building InstaCli v${VERSION}..."
echo ""

# Clean and create build directory
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Build for each platform
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    GOOS="${platform%/*}"
    GOARCH="${platform#*/}"
    
    output_name="${APP_NAME}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "📦 Building for ${GOOS}/${GOARCH}..."
    
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w -X main.Version=${VERSION}" \
        -o "${BUILD_DIR}/${output_name}" ./cmd/instacli
    
    echo "   ✓ ${output_name}"
done

echo ""
echo "✅ Build complete! Binaries are in ${BUILD_DIR}/"
echo ""
ls -lh $BUILD_DIR/

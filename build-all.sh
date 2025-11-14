#!/bin/bash

# Build script for all platforms

set -e

echo "Building kainet for all platforms..."
echo ""

# Create bin directory
mkdir -p bin

# macOS Apple Silicon
echo "Building macOS Apple Silicon..."
GOARCH=arm64 GOOS=darwin go build -o bin/kainet-darwin-arm64
echo "✓ bin/kainet-darwin-arm64"

# macOS Intel
echo "Building macOS Intel..."
GOARCH=amd64 GOOS=darwin go build -o bin/kainet-darwin-amd64
echo "✓ bin/kainet-darwin-amd64"

# Linux
echo "Building Linux..."
GOARCH=amd64 GOOS=linux go build -o bin/kainet-linux-amd64
echo "✓ bin/kainet-linux-amd64"

# Windows 64-bit
echo "Building Windows 64-bit..."
GOARCH=amd64 GOOS=windows go build -o bin/kainet.exe
echo "✓ bin/kainet.exe"

echo ""
echo "All binaries built successfully in bin/"
echo ""
ls -lh bin/

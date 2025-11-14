#!/bin/bash

set -e

echo "Building kaiNET for all platforms..."

# Create bin directory
mkdir -p bin

# Build for macOS ARM64
echo "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -o bin/kainet-darwin-arm64 main.go

# Build for macOS AMD64
echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -o bin/kainet-darwin-amd64 main.go

# Build for Linux AMD64
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o bin/kainet-linux-amd64 main.go

# Build for Windows AMD64
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -o bin/kainet.exe main.go

echo ""
echo "Build complete! Binaries created:"
ls -lh bin/

echo ""
echo "To run on your current platform:"
echo "  macOS ARM64: ./bin/kainet-darwin-arm64 <username> <room-name>"
echo "  macOS AMD64: ./bin/kainet-darwin-amd64 <username> <room-name>"
echo "  Linux AMD64: ./bin/kainet-linux-amd64 <username> <room-name>"
echo "  Windows:     bin\\kainet.exe <username> <room-name>"

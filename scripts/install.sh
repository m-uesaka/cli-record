#!/bin/bash

# Installation script for cli-record (Unix-like systems)

set -e

# Colors for output
GREEN='\033[0.32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="cli-record"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-0.2.0}"

echo -e "${GREEN}Installing ${APP_NAME}...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "Go version: $(go version)"

# Build the application
echo "Building ${APP_NAME}..."
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build \
    -ldflags "-X 'github.com/m-uesaka/cli-record/cmd.Version=${VERSION}' \
              -X 'github.com/m-uesaka/cli-record/cmd.Commit=${COMMIT}' \
              -X 'github.com/m-uesaka/cli-record/cmd.BuildTime=${BUILD_TIME}'" \
    -o "${APP_NAME}" .

# Check if we need sudo for installation
if [ -w "${INSTALL_DIR}" ]; then
    SUDO=""
else
    SUDO="sudo"
    echo -e "${YELLOW}Installation requires sudo access for ${INSTALL_DIR}${NC}"
fi

# Install the binary
echo "Installing to ${INSTALL_DIR}..."
${SUDO} mv "${APP_NAME}" "${INSTALL_DIR}/${APP_NAME}"
${SUDO} chmod +x "${INSTALL_DIR}/${APP_NAME}"

# Verify installation
if command -v ${APP_NAME} &> /dev/null; then
    echo -e "${GREEN}✓ ${APP_NAME} installed successfully!${NC}"
    echo ""
    echo "Version:"
    ${APP_NAME} --version
    echo ""
    echo "Get started with:"
    echo "  ${APP_NAME} start --task 'My first task'"
else
    echo -e "${RED}Installation failed. Please check if ${INSTALL_DIR} is in your PATH.${NC}"
    exit 1
fi

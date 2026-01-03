#!/bin/bash

# Uninstallation script for cli-record (Unix-like systems)

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="cli-record"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

echo -e "${YELLOW}Uninstalling ${APP_NAME}...${NC}"

# Check if the binary exists
if [ ! -f "${INSTALL_DIR}/${APP_NAME}" ]; then
    echo -e "${RED}${APP_NAME} is not installed in ${INSTALL_DIR}${NC}"
    exit 1
fi

# Check if we need sudo for uninstallation
if [ -w "${INSTALL_DIR}" ]; then
    SUDO=""
else
    SUDO="sudo"
    echo -e "${YELLOW}Uninstallation requires sudo access for ${INSTALL_DIR}${NC}"
fi

# Remove the binary
${SUDO} rm -f "${INSTALL_DIR}/${APP_NAME}"

echo -e "${GREEN}✓ ${APP_NAME} uninstalled successfully!${NC}"
echo ""
echo "Your data and configuration files are preserved at:"
echo "  Data: ~/.cli-record/data.json"
echo "  Config: ~/.config/cli-record/config.toml"
echo "  Archives: ~/.cli-record/archives/"
echo ""
echo "To remove all data, run:"
echo "  rm -rf ~/.cli-record ~/.config/cli-record"

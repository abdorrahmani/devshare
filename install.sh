#!/bin/bash
set -e

BINARY_NAME="devshare"
INSTALL_PATH="/usr/local/bin/$BINARY_NAME"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (e.g., with sudo)"
  exit 1
fi

if [ ! -f "$BINARY_NAME" ]; then
  echo "$BINARY_NAME not found in current directory. Please run this script from the extracted archive folder."
  exit 1
fi

cp "$BINARY_NAME" "$INSTALL_PATH"
chmod +x "$INSTALL_PATH"
echo "$BINARY_NAME installed to $INSTALL_PATH"

if ! echo "$PATH" | grep -q "/usr/local/bin"; then
  echo "export PATH=\"/usr/local/bin:$PATH\"" >> ~/.profile
  echo "/usr/local/bin added to PATH in ~/.profile. Please restart your terminal or run: source ~/.profile"
fi

echo "✅ Installation complete!" 
echo "Built with 💙 by Anophel" 
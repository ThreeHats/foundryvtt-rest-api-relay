#!/usr/bin/env bash
#
# Setup script for the Go relay server.
# Installs Go dependencies and Chromium for headless browser sessions.
#
# Usage: ./setup.sh

set -e

echo "=== Go Relay Server Setup ==="

# Check Go
if ! command -v go &>/dev/null; then
  if [ -x /usr/local/go/bin/go ]; then
    export PATH=$PATH:/usr/local/go/bin
  else
    echo "ERROR: Go is not installed. Install from https://go.dev/dl/"
    exit 1
  fi
fi
echo "Go: $(go version)"

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
echo "Done."

# Check/install Chromium
echo ""
echo "Checking for Chrome/Chromium..."
CHROME=""
for bin in chromium-browser chromium google-chrome google-chrome-stable; do
  if command -v "$bin" &>/dev/null; then
    CHROME=$(command -v "$bin")
    break
  fi
done

if [ -n "$CHROME" ]; then
  echo "Found: $CHROME"
else
  echo "Chrome/Chromium not found. Installing..."
  if command -v apt &>/dev/null; then
    sudo apt update && sudo apt install -y chromium-browser || sudo apt install -y chromium
  elif command -v apk &>/dev/null; then
    sudo apk add --no-cache chromium
  elif command -v dnf &>/dev/null; then
    sudo dnf install -y chromium
  elif command -v snap &>/dev/null; then
    sudo snap install chromium
  elif command -v brew &>/dev/null; then
    brew install --cask chromium
  else
    echo "ERROR: Could not detect package manager. Install Chromium manually."
    exit 1
  fi

  # Verify
  for bin in chromium-browser chromium google-chrome; do
    if command -v "$bin" &>/dev/null; then
      CHROME=$(command -v "$bin")
      break
    fi
  done

  if [ -n "$CHROME" ]; then
    echo "Installed: $CHROME"
  else
    echo "WARNING: Chromium installation may have failed. Headless sessions won't work."
  fi
fi

# Create .env if it doesn't exist
if [ ! -f .env ]; then
  echo ""
  echo "Creating .env from .env.example..."
  cp .env.example .env
  echo "Edit .env to configure your server."
fi

# Build
echo ""
echo "Building server..."
go build -o ./relay ./cmd/server/
echo "Done."

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Run the server:"
echo "  ./relay"
echo ""
echo "Or with hot-reload (install air first: go install github.com/air-verse/air@latest):"
echo "  air"

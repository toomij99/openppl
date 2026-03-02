#!/usr/bin/env bash
set -euo pipefail

REPO="${OPENPPL_REPO:-toomij99/openppl}"
VERSION="${OPENPPL_VERSION:-latest}"
INSTALL_DIR="${OPENPPL_INSTALL_DIR:-$HOME/bin}"
BIN_NAME="openppl"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64 | arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux | darwin)
    ;;
  *)
    echo "Unsupported operating system: $OS"
    exit 1
    ;;
esac

ASSET="${BIN_NAME}_${OS}_${ARCH}.tar.gz"
CHECKSUMS="checksums.txt"

if [ "$VERSION" = "latest" ]; then
  ASSET_URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
  CHECKSUMS_URL="https://github.com/${REPO}/releases/latest/download/${CHECKSUMS}"
else
  ASSET_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
  CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/${CHECKSUMS}"
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${ASSET_URL}"
ASSET_PATH="$TMP_DIR/$ASSET"
curl -fsSL "$ASSET_URL" -o "$ASSET_PATH"

if curl -fsSL "$CHECKSUMS_URL" -o "$TMP_DIR/checksums.txt"; then
  if command -v sha256sum >/dev/null 2>&1; then
    (cd "$TMP_DIR" && sha256sum -c --strict --ignore-missing checksums.txt)
  elif command -v shasum >/dev/null 2>&1; then
    EXPECTED="$(grep " ${ASSET}$" "$TMP_DIR/checksums.txt" | awk '{print $1}')"
    ACTUAL="$(shasum -a 256 "$ASSET_PATH" | awk '{print $1}')"
    if [ "$EXPECTED" != "$ACTUAL" ]; then
      echo "Checksum verification failed for ${ASSET}"
      exit 1
    fi
  else
    echo "Warning: no sha256 tool available, skipping checksum verification"
  fi
else
  echo "Warning: checksums.txt not found, skipping checksum verification"
fi

tar -xzf "$ASSET_PATH" -C "$TMP_DIR"

if [ -f "$TMP_DIR/$BIN_NAME" ]; then
  BIN_PATH="$TMP_DIR/$BIN_NAME"
elif [ -f "$TMP_DIR/$BIN_NAME/$BIN_NAME" ]; then
  BIN_PATH="$TMP_DIR/$BIN_NAME/$BIN_NAME"
else
  echo "No ${BIN_NAME} executable found in archive"
  echo "Archive contents:"
  tar -tzf "$ASSET_PATH" || true
  exit 1
fi

mkdir -p "$INSTALL_DIR"
install -m 0755 "$BIN_PATH" "$INSTALL_DIR/$BIN_NAME"

echo "Installed ${BIN_NAME} to ${INSTALL_DIR}/${BIN_NAME}"
echo "Run: ${BIN_NAME} help"

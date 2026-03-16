#!/usr/bin/env bash
set -euo pipefail

REPO="${OPENPPL_REPO:-toomij99/openppl}"
VERSION="${OPENPPL_VERSION:-latest}"
INSTALL_DIR="${OPENPPL_INSTALL_DIR:-$HOME/bin}"
BIN_NAME="openppl"
LATEST_API_URL="https://api.github.com/repos/${REPO}/releases/latest"

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

print_logo() {
  cat <<'EOF'
   ____  ____  _______  _   __ ____  ____  __
  / __ \/ __ \/ ____/ |/ / / __ \/ __ \/ /
 / / / / /_/ / __/  |   / / /_/ / /_/ / / 
/ /_/ / ____/ /___ /   | / ____/ ____/ /___
\____/_/   /_____//_/|_|/_/   /_/   /_____/
EOF
}

is_interactive() {
  [ -t 0 ] && [ -t 1 ] && [ "${CI:-}" != "true" ] && [ "${OPENPPL_INSTALL_NONINTERACTIVE:-0}" != "1" ]
}

fetch_latest_tag() {
  curl -fsSL "$LATEST_API_URL" 2>/dev/null | grep -oE '"tag_name"[[:space:]]*:[[:space:]]*"[^"]+"' | head -n1 | cut -d '"' -f4
}

resolve_current_version() {
  if [ -x "$INSTALL_DIR/$BIN_NAME" ]; then
    "$INSTALL_DIR/$BIN_NAME" version 2>/dev/null | awk '{print $3}'
  fi
}

print_command_summary() {
  cat <<'EOF'

Next commands:
  openppl help           Show all commands
  openppl version        Show installed version
  openppl                Launch the terminal app
  openppl onboard        Run setup wizard
  openppl web            Start the web dashboard
  openppl motd quiz      Run today's ACS quiz
  openppl motd progress  Show PPL readiness progress
EOF
}

TARGET_VERSION="$VERSION"
if [ "$TARGET_VERSION" = "latest" ]; then
  RESOLVED_LATEST="$(fetch_latest_tag || true)"
  if [ -n "$RESOLVED_LATEST" ]; then
    TARGET_VERSION="$RESOLVED_LATEST"
  fi
fi

CURRENT_VERSION="$(resolve_current_version || true)"

if is_interactive; then
  print_logo
  printf '\nWelcome to the openppl installer\n\n'
  printf 'Install directory: %s\n' "$INSTALL_DIR"
  printf 'Platform:          %s/%s\n' "$OS" "$ARCH"
  if [ -n "$CURRENT_VERSION" ]; then
    printf 'Current version:   %s\n' "$CURRENT_VERSION"
  else
    printf 'Current version:   not installed\n'
  fi
  printf 'New version:       %s\n\n' "$TARGET_VERSION"
  printf 'What you can do after install:\n'
  printf '  openppl help          Show all commands\n'
  printf '  openppl onboard       Configure training profile\n'
  printf '  openppl web           Launch browser dashboard\n'
  printf "  openppl motd quiz     Answer today's ACS question\n"
  printf '  openppl motd progress Check checkride readiness\n\n'
  printf 'Continue with installation? [Y/n] '
  read -r CONFIRM
  case "${CONFIRM:-Y}" in
    n|N|no|NO)
      echo "Installation cancelled"
      exit 0
      ;;
  esac
  echo
fi

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
if [ -x "$INSTALL_DIR/$BIN_NAME" ]; then
  INSTALLED_VERSION="$("$INSTALL_DIR/$BIN_NAME" version 2>/dev/null)"
  INSTALLED_VERSION="${INSTALLED_VERSION##* }"
  if [ -n "$INSTALLED_VERSION" ]; then
    echo "Installed version: ${INSTALLED_VERSION}"
  fi
fi
print_command_summary

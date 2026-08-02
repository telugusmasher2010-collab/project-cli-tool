#!/usr/bin/env bash
set -euo pipefail

REPO="telugusmasher2010-collab/project-cli-tool"
INSTALL_DIR="${PROJ_INIT_INSTALL_DIR:-/usr/local/bin}"
VERSION="${PROJ_INIT_VERSION:-latest}"

info()  { printf "\033[0;32m%s\033[0m\n" "$*"; }
error() { printf "\033[0;31m%s\033[0m\n" "$*" >&2; }

# --- detect OS / arch -------------------------------------------------------
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  aarch64|arm64) arch="arm64" ;;
  *) error "Unsupported architecture: $arch"; exit 1 ;;
esac

case "$os" in
  linux|darwin) ;;
  *) error "Unsupported OS: $os"; exit 1 ;;
esac

# --- resolve release asset --------------------------------------------------
if [ "$VERSION" = "latest" ]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep -m1 '"tag_name"' | sed 's/.*: "\(.*\)",*/\1/')"
fi

FILE="proj-init_${VERSION#v}_${os}_${arch}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILE"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

info "Downloading $URL"
curl -fsSL "$URL" -o "$TMP/proj-init.tar.gz"

tar -xzf "$TMP/proj-init.tar.gz" -C "$TMP"
chmod +x "$TMP/proj-init"

mkdir -p "$INSTALL_DIR"
mv "$TMP/proj-init" "$INSTALL_DIR/proj-init"

info "proj-init $VERSION installed to $INSTALL_DIR/proj-init"
info "Run 'proj-init --help' to get started."

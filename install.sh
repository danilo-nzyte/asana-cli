#!/usr/bin/env bash
set -euo pipefail

REPO="danilo-nzyte/asana-cli"
INSTALL_DIR="/usr/local/bin"
SKILL_DIR="$HOME/.claude/skills/asana"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    darwin|linux) ;;
    *)
        echo "Unsupported OS: $OS (use install.ps1 for Windows)"
        exit 1
        ;;
esac

ARCHIVE="asana-cli_${OS}_${ARCH}.tar.gz"

# Get latest release tag
echo "==> Fetching latest release..."
TAG=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
if [ -z "$TAG" ]; then
    echo "Error: could not determine latest release."
    exit 1
fi
echo "    Latest release: $TAG"

# Download and extract
URL="https://github.com/${REPO}/releases/download/${TAG}/${ARCHIVE}"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "==> Downloading ${ARCHIVE}..."
curl -sSfL "$URL" -o "$TMPDIR/$ARCHIVE"

echo "==> Extracting..."
tar -xzf "$TMPDIR/$ARCHIVE" -C "$TMPDIR"

echo "==> Installing to $INSTALL_DIR (may require sudo)..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMPDIR/asana-cli" "$INSTALL_DIR/asana-cli"
else
    sudo mv "$TMPDIR/asana-cli" "$INSTALL_DIR/asana-cli"
fi
chmod +x "$INSTALL_DIR/asana-cli"

# Install Claude Code skill
echo "==> Installing Claude Code skill..."
mkdir -p "$SKILL_DIR"
SKILL_URL="https://raw.githubusercontent.com/${REPO}/${TAG}/skill/SKILL.md"
curl -sSfL "$SKILL_URL" -o "$SKILL_DIR/SKILL.md"
echo "    Skill installed to $SKILL_DIR"

echo ""
echo "==> Installed asana-cli $TAG to $INSTALL_DIR/asana-cli"
echo ""
echo "Next steps:"
echo "  1. Run: asana-cli auth login --client-id <ID> --client-secret <SECRET>"
echo "     (credentials are saved — you only need to provide them once)"
echo "  2. Verify: asana-cli auth status"
echo "  3. Set ASANA_WORKSPACE_ID in your shell config (needed for search/list commands)"
echo ""
echo "See https://github.com/${REPO}#authentication-setup for details."

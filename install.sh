#!/usr/bin/env bash

set -e

REPO="sodikinnaa/youtube-autoposter"
BINARY_NAME="youtube-autoposter"
SKILL_FILE="SKILL.md"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "❌ Arsitektur $ARCH tidak didukung."
    exit 1
    ;;
esac

case "$OS" in
  linux)
    TARGET="youtube-autoposter-linux-${ARCH}"
    ;;
  darwin)
    TARGET="youtube-autoposter-darwin-${ARCH}"
    ;;
  *)
    echo "❌ Sistem Operasi $OS tidak didukung via bash script."
    exit 1
    ;;
esac

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${TARGET}"
RAW_SKILL_URL="https://raw.githubusercontent.com/${REPO}/main/SKILL.md"

if [ -f "$BINARY_NAME" ]; then
  echo "🔄 Meng-update binary '$BINARY_NAME' yang sudah ada ke versi terbaru..."
else
  echo "🚀 Mendownload binary YouTube Auto-Poster (${OS}/${ARCH})..."
fi

# Download binary ke file temporary untuk keamanan update
curl -sSL "$DOWNLOAD_URL" -o "${BINARY_NAME}.tmp" || {
  echo "⚠️ Download rilis spesifik gagal. Mengambil binary fallback..."
  curl -sSL "https://github.com/${REPO}/releases/latest/download/youtube-autoposter-linux-amd64" -o "${BINARY_NAME}.tmp"
}

mv "${BINARY_NAME}.tmp" "$BINARY_NAME"
chmod +x "$BINARY_NAME"

# Download file panduan AI Agent SKILL.md
echo "📄 Mendownload file panduan AI Agent (SKILL.md)..."
curl -sSL "$RAW_SKILL_URL" -o "$SKILL_FILE" || {
  echo "⚠️ Gagal mendownload SKILL.md, mengabaikan..."
}

echo ""
echo "======================================================="
echo "✅ SUKSES INSTALL / UPDATE YOUTUBE AUTO-POSTER!"
echo "======================================================="
echo "📁 Binary Executable : ./$BINARY_NAME"
echo "🤖 Panduan AI Agent : ./$SKILL_FILE"
echo "======================================================="
echo "Jalankan dengan perintah: ./$BINARY_NAME"

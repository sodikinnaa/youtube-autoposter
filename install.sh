#!/usr/bin/env bash

set -e

REPO="sodikinnaa/youtube-autoposter"
BINARY_NAME="youtube-autoposter"

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

echo "🚀 Mendownload binary YouTube Auto-Poster (${OS}/${ARCH})..."
curl -sSL "$DOWNLOAD_URL" -o "$BINARY_NAME" || {
  echo "⚠️ Download rilis spesifik gagal. Mengambil binary default..."
  curl -sSL "https://github.com/${REPO}/releases/latest/download/youtube-autoposter-linux-amd64" -o "$BINARY_NAME"
}

chmod +x "$BINARY_NAME"

echo "✅ Berhasil mendownload ./$BINARY_NAME!"
echo "Jalankan dengan perintah: ./$BINARY_NAME"

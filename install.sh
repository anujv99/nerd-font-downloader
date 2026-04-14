#!/bin/sh
set -eu

REPO="anujv99/nerd-font-downloader"
BINARY_NAME="nfdownloader"
INSTALL_DIR=${NFDOWNLOADER_INSTALL_DIR:-"$HOME/.local/bin"}

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

if [ "$os" != "linux" ]; then
    printf 'Unsupported operating system: %s\n' "$os" >&2
    exit 1
fi

case "$arch" in
    x86_64|amd64)
        archive_arch="x86_64"
        ;;
    aarch64|arm64)
        archive_arch="arm64"
        ;;
    *)
        printf 'Unsupported architecture: %s\n' "$arch" >&2
        exit 1
        ;;
esac

version=${NFDOWNLOADER_VERSION:-}
if [ -z "$version" ]; then
    version=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | awk -F '"' '/tag_name/ { print $4; exit }')
fi

if [ -z "$version" ]; then
    printf 'Failed to determine the latest release version.\n' >&2
    exit 1
fi

version_without_v=${version#v}
archive_name="${BINARY_NAME}_${version_without_v}_Linux_${archive_arch}"
archive_file="${archive_name}.tar.gz"
download_url="https://github.com/$REPO/releases/download/$version/$archive_file"

tmp_dir=$(mktemp -d)
cleanup() {
    rm -rf "$tmp_dir"
}
trap cleanup EXIT INT TERM

printf 'Installing %s %s for Linux %s\n' "$BINARY_NAME" "$version" "$archive_arch"
printf 'Downloading %s\n' "$download_url"

curl -fsSL "$download_url" -o "$tmp_dir/$archive_file"
tar -xzf "$tmp_dir/$archive_file" -C "$tmp_dir"

mkdir -p "$INSTALL_DIR"

if command -v install >/dev/null 2>&1; then
    install -m 0755 "$tmp_dir/$archive_name/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    cp "$tmp_dir/$archive_name/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod 0755 "$INSTALL_DIR/$BINARY_NAME"
fi

printf '\nInstalled to %s/%s\n' "$INSTALL_DIR" "$BINARY_NAME"

case ":$PATH:" in
    *":$INSTALL_DIR:"*)
        ;;
    *)
        printf 'Add %s to your PATH if it is not already available.\n' "$INSTALL_DIR"
        ;;
esac

printf 'Run `%s` to start the app.\n' "$BINARY_NAME"

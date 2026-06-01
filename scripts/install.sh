#!/bin/bash
set -e

PACKAGE_NAME="cmt"
BASE_URL="https://github.com/segersniels/${PACKAGE_NAME}/releases/latest/download"
DEFAULT_INSTALL_DIR="${HOME}/.local/bin"
INSTALL_DIR="${DEFAULT_INSTALL_DIR}"

if [ "$#" -eq 1 ]; then
  INSTALL_DIR="$1"
fi

if [ "$#" -gt 1 ]; then
  echo "Usage: install.sh [install-dir]" >&2
  exit 1
fi

# Identify OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "${OS}" in
Linux*) os=linux ;;
Darwin*) os=darwin ;;
*)
  echo "Unsupported OS. Exiting..."
  exit 1
  ;;
esac

case "${ARCH}" in
x86_64*) arch=amd64 ;;
arm64*) arch=arm64 ;;
*)
  echo "Unsupported architecture. Exiting..."
  exit 1
  ;;
esac

# Construct binary name and download URL
BIN_NAME="${PACKAGE_NAME}-${os}-${arch}"
DOWNLOAD_URL="${BASE_URL}/${BIN_NAME}"

# Full path to the target binary
FULL_PATH="${INSTALL_DIR}/${PACKAGE_NAME}"

echo "Downloading ${BIN_NAME} to ${FULL_PATH}..."

mkdir -p "${INSTALL_DIR}"
tmp_file="$(mktemp)"

# Download the binary
if command -v curl >/dev/null; then
  curl -L -o "${tmp_file}" "${DOWNLOAD_URL}"
elif command -v wget >/dev/null; then
  wget -O "${tmp_file}" "${DOWNLOAD_URL}"
else
  echo "Error: curl or wget is required to download the binary."
  exit 1
fi

# Make the binary executable
chmod +x "${tmp_file}"
mv "${tmp_file}" "${FULL_PATH}"
echo "Download completed. The binary is available at ${FULL_PATH}"

if [ "${INSTALL_DIR}" = "${DEFAULT_INSTALL_DIR}" ] && [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo ""
  echo "Note: ${INSTALL_DIR} is not in your PATH"
  echo ""
  echo "Add it to your shell config:"

  current_shell="$(basename "$SHELL")"
  show_bash=false
  show_zsh=false
  show_fish=false

  case "${current_shell}" in
  bash) show_bash=true ;;
  zsh) show_zsh=true ;;
  fish) show_fish=true ;;
  esac

  if [ -f "$HOME/.bashrc" ] || [ -f "$HOME/.bash_profile" ]; then show_bash=true; fi
  if [ -f "$HOME/.zshrc" ]; then show_zsh=true; fi
  if [ -f "$HOME/.config/fish/config.fish" ]; then show_fish=true; fi

  if [ "$show_bash" = true ]; then
    echo ""
    echo "  # bash"
    echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
    echo "  source ~/.bashrc"
  fi

  if [ "$show_zsh" = true ]; then
    echo ""
    echo "  # zsh"
    echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc"
    echo "  source ~/.zshrc"
  fi

  if [ "$show_fish" = true ]; then
    echo ""
    echo "  # fish"
    echo "  echo 'fish_add_path \$HOME/.local/bin' >> ~/.config/fish/config.fish"
    echo "  source ~/.config/fish/config.fish"
  fi

  if [ "$show_bash" != true ] && [ "$show_zsh" != true ] && [ "$show_fish" != true ]; then
    echo ""
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
  fi

  echo ""
fi

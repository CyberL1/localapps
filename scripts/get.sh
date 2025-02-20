#!/bin/sh

set -e

if ! command -v unzip > /dev/null; then
  echo "Error: Unzip is required to install localapps"
  exit 1
fi

dir="$HOME/localapps/bin"
zip="$dir/localapps.zip"
exe="$dir/localapps"

if [ "$OS" = "Windows_NT" ]; then
  target="windows"
else
  case $(uname -sm) in
  "Darwin x86_64") target="darwin-amd64" ;;
  "Darwin arm64") target="dawin-arm64" ;;
  "Linux aarch64") target="linux-arm64" ;;
  *) target="linux-amd64"
  esac
fi

download_url="https://github.com/CyberL1/localapps/releases/latest/download/localapps-${target}.zip"

if [ ! -d $dir ]; then
  mkdir -p $dir
fi

curl --fail --location --progress-bar --output $zip $download_url
unzip -d $dir -o $zip
chmod +x $exe
rm $zip

echo "Localapps CLI was installed to $exe"
if command -v localapps > /dev/null; then
  echo "Run 'localapps up' to get started"
else
  case $SHELL in
  /bin/zsh) shell_profile=".zshrc" ;;
  *) shell_profile=".bashrc" ;;
  esac
  echo "export PATH=\"$dir:\$PATH\"" >> $HOME/$shell_profile
fi

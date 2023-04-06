#!/bin/bash
set -o errexit -o nounset -o pipefail

check_cmd() {
    command -v "$1" > /dev/null 2>&1
}

echo "0-----------------------"
echo "* Install jq and curl"
if check_cmd jq; then
  echo "jq has been already installed"
elif check_cmd brew; then
      brew install jq
elif check_cmd apt; then
      apt install jq -y
elif check_cmd yum; then
      yum install jq
fi

if check_cmd curl; then
  echo "curl has been already installed"
elif check_cmd apt; then
  apt install curl
elif check_cmd yum; then
  yum install curl
fi

echo "1-----------------------"
echo "* Install Rust"
if check_cmd rustup; then
  echo "Rust has been already installed"
else
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
fi

source "$HOME/.cargo/env"

echo "2-----------------------"
echo "* Set Rust default toolchain"
rustup default stable

echo "3-----------------------"
echo "* Check cargo version"
cargo version
# If this is lower than 1.55.0+, update
echo "4-----------------------"
echo "* Update Rust default toolchain"
rustup update stable

echo "5-----------------------"
echo "* Add wasm32 target"
rustup target add wasm32-unknown-unknown

echo "------------------------"
echo "Build wasm environment successfully!"


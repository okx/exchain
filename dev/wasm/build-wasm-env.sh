#!/bin/bash
set -o errexit -o nounset -o pipefail

echo "0-----------------------"
echo "* Install Rust"
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

echo "1-----------------------"
echo "* Set Rust default toolchain"
rustup default stable

echo "2-----------------------"
echo "* Check cargo version"
cargo version
# If this is lower than 1.55.0+, update
echo "3-----------------------"
echo "* Update Rust default toolchain"
rustup update stable

echo "4-----------------------"
echo "* Add wasm32 target"
rustup target add wasm32-unknown-unknown

echo "------------------------"
echo "Build wasm environment successfully!"
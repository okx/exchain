#!/bin/sh
#set -e

echo "install wasm static lib with sudo"
installWasmLib() {
  if [ -r  /usr/local/lib/libwasmvm_muslc.a ]; then
      exit 0
  elif [ -r /lib/libwasmvm_muslc.a ]; then
      exit 0
  fi
  wget --no-check-certificate "https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0/libwasmvm_muslc.x86_64.a" -O /usr/local/lib/libwasmvm_muslc.x86_64.a
  cp /usr/local/lib/libwasmvm_muslc.x86_64.a /usr/local/lib/libwasmvm_muslc.a
  echo "install wasm static lib success"
}

installWasmLib



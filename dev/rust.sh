set -x

cd /Users/oker/workspace/github/wasm-native/wasm-native-code/wasmy-counter/
rm -rf depend/*
cp -rf /Users/oker/workspace/github/cosmwasm/packages/* depend/

cd /Users/oker/workspace/github/wasmvm/libwasmvm
cargo build --release
cp target/release/libwasmvm.dylib ../api/libwasmvm.dylib

cd /Users/oker/workspace/exchain/

go mod vendor


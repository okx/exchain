set -x

cd ../../wasm-test/contract/keccak-test/
rm -rf depend/*
cp -rf ../../../cosmwasm/packages/* depend/

cd -

cd ../../wasmvm/libwasmvm
cargo build --release
#cp target/release/libwasmvm.dylib ../api/libwasmvm.dylib
cp target/release/libwasmvm.dylib ../internal/api/libwasmvm.dylib

cd -

cd ../../exchain/

go mod tidy
go mod vendor


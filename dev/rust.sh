set -x

cd ../../wasm-test/contract/keccak-test/
rm -rf depend/*
rm -rf /Users/oker/workspace/wasm-native/vm-benchmark/contracts/wasm/counter/depend/*
cp -rf ../../../cosmwasm/packages/* depend/
cp -rf ../../../cosmwasm/packages/* /Users/oker/workspace/wasm-native/vm-benchmark/contracts/wasm/counter/depend/
cd -

cd ../../wasmvm/libwasmvm
cargo build --release
#cp target/release/libwasmvm.dylib ../api/libwasmvm.dylib
cp target/release/libwasmvm.dylib ../internal/api/libwasmvm.dylib
# cp target/release/libwasmvm.dylib ../../vm-benchmark/

cd -

cd ../../exchain/

go mod tidy
go mod vendor

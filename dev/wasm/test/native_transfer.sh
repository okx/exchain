#!/bin/bash
set -o errexit -o nounset -o pipefail

CHAIN_ID="exchain-67"
NODE="http://localhost:26657"
#CHAIN_ID="exchain-64"
#NODE="http://3.113.237.222:26657"
TX_EXTRA="--fees 0.01okt --gas 3000000 --chain-id=$CHAIN_ID --node $NODE -b block -y"

captain=$(exchaincli keys show captain -a)
admin18=$(exchaincli keys show admin18 -a)
admin17=$(exchaincli keys show admin17 -a)


# store wasm code
res=$(exchaincli tx wasm store ../native-transfer/artifacts/native_transfer.wasm --instantiate-everybody=true --from captain $TX_EXTRA)
echo $res | jq
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "store cw20 contract succeed, code id: $code_id"

# instantiate wasm contract
res=$(exchaincli tx wasm instantiate "$code_id" '{}' --label test1 --admin "$captain" --from captain $TX_EXTRA)
contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "instantiate cw20 succeed, contract address: $contractAddr"

# transfer okt
res=$(exchaincli tx wasm execute "$contractAddr" '{"transfer":{"recipient":"ex1qlruqeurjk9hcnfkfp90vzkh4z2vr0v46x5v9d"}}' --amount=10okt --from captain $TX_EXTRA)
echo $res | jq
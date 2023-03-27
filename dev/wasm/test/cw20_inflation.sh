#!/bin/bash
set -o errexit -o nounset -o pipefail

#CHAIN_ID="exchain-67"
#NODE="http://localhost:26657"
CHAIN_ID="exchain-64"
NODE="http://3.113.237.222:26657"
TX_EXTRA="--fees 0.01okt --gas 3000000 --chain-id=$CHAIN_ID --node $NODE -b block -y"

captain=$(exchaincli keys show captain -a)
admin18=$(exchaincli keys show admin18 -a)
admin17=$(exchaincli keys show admin17 -a)


# store wasm code
res=$(exchaincli tx wasm store ../cw20_inflation.wasm --instantiate-everybody=true --from captain $TX_EXTRA)
echo $res | jq
cw20_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "store cw20 contract succeed, code id: $cw20_code_id"

# instantiate wasm contract
res=$(exchaincli tx wasm instantiate "$cw20_code_id" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"10000000000"}],"name":"my test token", "symbol":"mtt"}' --label "cw20 inflation" --admin "$captain" --from captain $TX_EXTRA)
cw20contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "instantiate cw20 succeed, contract address: $cw20contractAddr"

res=$(exchaincli tx wasm instantiate "$cw20_code_id" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"10000000000"}],"name":"my test token", "symbol":"mtt"}' --label "cw20 inflation" --admin "$captain" --from captain $TX_EXTRA)
cw20contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "instantiate cw20 succeed, contract address: $cw20contractAddr"

# transfer cw20 token
#res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"100","recipient":"'$admin18'"}}' --from captain $TX_EXTRA)
#echo $res | jq

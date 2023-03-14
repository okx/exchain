#!/bin/bash
set -o errexit -o nounset -o pipefail

CHAIN_ID="exchain-67"
NODE="http://localhost:26657"
QUERY_EXTRA="--node=$NODE"
TX_EXTRA_UNBLOCKED="--fees 0.01okb --gas 3000000 --chain-id=$CHAIN_ID --node $NODE -b async -y"
TX_EXTRA="--fees 0.01okb --gas 3000000 --chain-id=$CHAIN_ID --node $NODE -b block -y"
captain=$(exchaincli keys show captain -a)


# claim cw20 from ce4-stake
totalAmount="100000000"
transferAmount="100"

res=$(exchaincli tx wasm store ../cw20-base/artifacts/cw20_base.wasm --from captain $TX_EXTRA)
cw20_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')

res=$(exchaincli tx wasm instantiate "$cw20_code_id" '{"decimals":10,"initial_balances":[{"address":"'"$captain"'","amount":"'$totalAmount'"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --from captain $TX_EXTRA)
cw20contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "cw20 contract address: $cw20contractAddr"

res=$(exchaincli tx wasm store ../cw4-stake/artifacts/cw4_stake.wasm --from $captain $TX_EXTRA)
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
res=$(exchaincli tx wasm instantiate "$code_id" '{"denom":{"cw20":"'$cw20contractAddr'"},"min_bond":"50","tokens_per_weight":"10","unbonding_period":{"height":0}}' --label test1 --admin $captain --from captain $TX_EXTRA)
contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "cw4-stake contract address: $contractAddr"

res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"send":{"amount":"'$transferAmount'","contract":"'$contractAddr'","msg":"eyJib25kIjp7fX0="}}' --from captain $TX_EXTRA)  # msg={"bond":{}}
echo $res | jq

res=$(exchaincli tx wasm execute "$contractAddr" '{"unbond":{"tokens":"'$transferAmount'"}}' --from captain $TX_EXTRA)
echo $res | jq

res=$(exchaincli tx wasm execute "$contractAddr" '{"claim":{}}' --from captain $TX_EXTRA)
echo $res | jq



# claim okb from cw4-stake
res=$(exchaincli tx wasm store ../cw4-stake/artifacts/cw4_stake.wasm --from $captain $TX_EXTRA)
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
# native token must be "okb", not "OKB" or tokens with other names
res=$(exchaincli tx wasm instantiate "$code_id" '{"denom":{"native":"okb"},"min_bond":"50","tokens_per_weight":"5","unbonding_period":{"height":0}}' --label test1 --admin $captain --from captain $TX_EXTRA)
contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "cw4-stake contract address: $contractAddr"

res=$(exchaincli query wasm contract-state smart "$contractAddr" '{"staked":{"address":"'$captain'"}}' $QUERY_EXTRA)
echo $res | jq

res=$(exchaincli tx wasm execute "$contractAddr" '{"bond":{}}' --amount=10okb --from captain $TX_EXTRA)
echo $res | jq

res=$(exchaincli query wasm contract-state smart "$contractAddr" '{"staked":{"address":"'$captain'"}}' $QUERY_EXTRA)
echo $res | jq

res=$(exchaincli query wasm contract-state smart "$contractAddr" '{"member":{"addr":"'$captain'"}}' $QUERY_EXTRA)
echo $res | jq

res=$(exchaincli tx wasm execute "$contractAddr" '{"unbond":{"tokens":"10000000000000000000"}}' --from captain $TX_EXTRA)
echo $res | jq

res=$(exchaincli tx wasm execute "$contractAddr" '{"claim":{}}' --from captain $TX_EXTRA)
echo $res | jq

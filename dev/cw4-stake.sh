captain=$(exchaincli keys show captain -a)
# cw20
# ex14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s6fqu27
# cw4-stake
# ex1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqku97gc

echo "## store cw20 contract..."
res=$(exchaincli tx wasm store ./wasm/cw20-base/artifacts/cw20_base.wasm --fees 0.01okt --from captain --gas=3000000 -b block -y)
echo "store cw20 contract succeed"
cw20_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "## instantiate cw20 contract..."
res=$(exchaincli tx wasm instantiate "$cw20_code_id" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin $captain --fees 0.001okt --from captain -b block -y)
echo "instantiate cw20 succeed"
cw20contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "cw20 contract address: $cw20contractAddr"
exchaincli query wasm contract-state smart "$cw20contractAddr" '{"balance":{"address":"'$captain'"}}'

echo "## store cw4-stake contract..."
res=$(exchaincli tx wasm store ./wasm/cw4-stake/artifacts/cw4_stake.wasm --fees 0.01okt --from captain --gas=2000000 -b block -y)
echo "store cw4-stake succeed"
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "## instantiate cw4-stake contract..."
res=$(exchaincli tx wasm instantiate "$code_id" '{"denom":{"cw20":"'$cw20contractAddr'"},"min_bond":"100","tokens_per_weight":"10","unbonding_period":{"height":100}}' --label test1 --admin $captain --fees 0.001okt --from captain -b block -y)
echo "instantiate cw4-stake succeed"
contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "$contractAddr"
exchaincli query wasm contract-state smart "$contractAddr" '{"staked":{"address":"'$captain'"}}'
exchaincli query wasm contract-state smart "$contractAddr" '{"member":{"addr":"'$captain'"}}'

echo "## send cw20 to cw4-stake and call Receive() method of cw4-stake"
res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"send":{"amount":"100","contract":"'$contractAddr'","msg":"eyJib25kIjp7fX0="}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
tx_hash=$(echo "$res" | jq '.txhash' | sed 's/\"//g')
echo "txhash: $tx_hash"
exchaincli query wasm contract-state smart "$cw20contractAddr" '{"balance":{"address":"'$captain'"}}'
exchaincli query wasm contract-state smart "$cw20contractAddr" '{"balance":{"address":"'$contractAddr'"}}'
exchaincli query wasm contract-state smart "$contractAddr" '{"staked":{"address":"'$captain'"}}'
exchaincli query wasm contract-state smart "$contractAddr" '{"member":{"addr":"'$captain'"}}'

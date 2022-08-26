captain=$(exchaincli keys show captain -a)
echo "captain addr: $captain"
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
res=$(exchaincli tx wasm store ./wasm/cw4-stake/artifacts/cw4_stake.wasm --fees 0.01okt --from admin18 --gas=2000000 -b block -y)
echo "store cw4-stake succeed"
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "## instantiate cw4-stake contract..."
res=$(exchaincli tx wasm instantiate "$code_id" '{"denom":{"cw20":"'$cw20contractAddr'"},"min_bond":"100","tokens_per_weight":"10","unbonding_period":{"height":100}}' --label test1 --admin $captain --fees 0.001okt --from captain -b block -y)
echo "instantiate cw4-stake succeed"
contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "cw4-stake contract address: $contractAddr"
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

res=$(exchaincli tx wasm store wasm/test/burner.wasm --from captain --fees 0.001okt --gas 1000000 -b block -y)
burner_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "burner_code_id: $burner_code_id"

# block contract to execute
echo "migrate cw20 contract to a new wasm code"
res=$(exchaincli tx gov submit-proposal migrate-contract "$cw20contractAddr" "$burner_code_id" '{"payout": "'$captain'"}' --deposit 10.1okt --title "test title" --description "test description" --fees 0.01okt --from captain --gas=3000000 -b block -y)
proposal_id=$(echo "$res" | jq '.logs[0].events[1].attributes[1].value' | sed 's/\"//g')
echo "proposal_id: $proposal_id"

res=$(exchaincli tx gov deposit "$proposal_id" 90000000okt --fees 0.001okt --from captain -b block -y)
res=$(exchaincli tx gov vote "$proposal_id" yes --fees 0.001okt --from captain -b block -y)

echo "## call send method of cw20 contract after migrating which is expected to fail"
res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"send":{"amount":"100","contract":"'$contractAddr'","msg":"eyJib25kIjp7fX0="}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
tx_hash=$(echo "$res" | jq '.txhash' | sed 's/\"//g')
echo "txhash: $tx_hash"
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
echo "expected to fail, raw_log: $raw_log"

# update whitelist
echo "## update deployment whitelist and store wasm code"
res=$(exchaincli tx gov submit-proposal update-wasm-deployment-whitelist "ex1h0j8x0v9hs4eq6ppgamemfyu4vuvp2sl0q9p3v,ex15nnhqdf9sds0s063kaaretxj3ftlnzrguhfdeq" --deposit 10.1okt --title "test title" --description "test description" --fees 0.001okt --from captain --gas=3000000 -b block -y)
proposal_id=$(echo "$res" | jq '.logs[0].events[1].attributes[1].value' | sed 's/\"//g')
echo "proposal_id: $proposal_id"
res=$(exchaincli tx gov deposit "$proposal_id" 90000000okt --fees 0.001okt --from captain -b block -y)
res=$(exchaincli tx gov vote "$proposal_id" yes --fees 0.001okt --from captain -b block -y)

res=$(exchaincli tx wasm store wasm/test/burner.wasm --from admin18 --fees 0.001okt --gas 1000000 -b block -y)
tx_hash=$(echo "$res" | jq '.txhash' | sed 's/\"//g')
echo "txhash: $tx_hash"
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
echo "expected to fail, raw_log: $raw_log"
res=$(exchaincli tx wasm store wasm/test/burner.wasm --from captain --fees 0.001okt --gas 1000000 -b block -y)
tx_hash=$(echo "$res" | jq '.txhash' | sed 's/\"//g')
echo "txhash: $tx_hash"
burner_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "burner_code_id: $burner_code_id"

# update whitelist
res=$(exchaincli tx gov submit-proposal update-wasm-deployment-whitelist all --deposit 10.1okt --title "test title" --description "test description" --fees 0.001okt --from captain --gas=3000000 -b block -y)
proposal_id=$(echo "$res" | jq '.logs[0].events[1].attributes[1].value' | sed 's/\"//g')
echo "proposal_id: $proposal_id"
res=$(exchaincli tx gov deposit "$proposal_id" 90000000okt --fees 0.001okt --from captain -b block -y)
res=$(exchaincli tx gov vote "$proposal_id" yes --fees 0.001okt --from captain -b block -y)
res=$(exchaincli tx wasm store wasm/test/burner.wasm --from admin18 --fees 0.001okt --gas 1000000 -b block -y)
tx_hash=$(echo "$res" | jq '.txhash' | sed 's/\"//g')
echo "txhash: $tx_hash"
burner_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "burner_code_id: $burner_code_id"
res=$(exchaincli tx wasm store wasm/test/burner.wasm --from captain --fees 0.001okt --gas 1000000 -b block -y)
tx_hash=$(echo "$res" | jq '.txhash' | sed 's/\"//g')
echo "txhash: $tx_hash"
burner_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
echo "burner_code_id: $burner_code_id"

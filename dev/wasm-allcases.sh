#!/bin/bash
set -o errexit -o nounset -o pipefail

captain=$(exchaincli keys show captain -a)
admin18=$(exchaincli keys show admin18 -a)
admin17=$(exchaincli keys show admin17 -a)

#####################################################
#############       store code       ################
#####################################################

echo "## store cw20 contract...everybody"
res=$(exchaincli tx wasm store ./wasm/cw20-base/artifacts/cw20_base.wasm --instantiate-everybody=true --fees 0.01okt --from captain --gas=3000000 -b block -y)
echo "store cw20 contract succeed"
cw20_code_id1=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')

echo "## store cw20 contract...nobody"
res=$(exchaincli tx wasm store ./wasm/cw20-base/artifacts/cw20_base.wasm --instantiate-nobody=true --fees 0.01okt --from captain --gas=3000000 -b block -y)
echo "store cw20 contract succeed"
cw20_code_id2=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')

echo "## store cw20 contract...only-address"
res=$(exchaincli tx wasm store ./wasm/cw20-base/artifacts/cw20_base.wasm --instantiate-only-address="${captain}" --fees 0.01okt --from captain --gas=3000000 -b block -y)
echo "store cw20 contract succeed"
cw20_code_id3=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')

echo "## store cw20 contract...null access"
res=$(exchaincli tx wasm store ./wasm/cw20-base/artifacts/cw20_base.wasm --fees 0.01okt --from captain --gas=3000000 -b block -y)
echo "store cw20 contract succeed"
cw20_code_id4=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')

echo "## store gzipped cw20 contract...null access"
res=$(exchaincli tx wasm store ./wasm/test/cw20_base_gzip.wasm --fees 0.01okt --from captain --gas=3000000 -b block -y)
echo "store cw20 contract succeed"
cw20_code_id5=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
data_hash4=$(exchaincli query wasm code-info "${cw20_code_id4}" | jq '.data_hash' | sed 's/\"//g')
data_hash5=$(exchaincli query wasm code-info "${cw20_code_id5}" | jq '.data_hash' | sed 's/\"//g')
if [[ "${data_hash4}" != "${data_hash5}" ]];
then
  echo "wrong data hash of gzipped cw20 contract"
  exit 1
fi;

echo "## store invalid cw20 contract...null access"
res=$(exchaincli tx wasm store ./wasm/test/invalid.wasm --fees 0.01okt --from captain --gas=3000000 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="create wasm contract failed: Error calling the VM: Error during static Wasm validation: Wasm bytecode could not be deserialized. Deserialization error: \I/O Error: UnexpectedEof\: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when store invalid wasm code"
  exit 1
fi;

#####################################################
#########    instantiate contract      ##############
#####################################################
echo "## instantiate everybody..."
res=$(exchaincli tx wasm instantiate "$cw20_code_id1" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from captain -b block -y)
echo "instantiate cw20 succeed"
if [[ $(echo "$res" | jq '.logs[0].events[0].attributes[0].key' | sed 's/\"//g') != "_contract_address" ]];
then
  echo "unexpected result of instantiate"
  exit 1
fi;
res=$(exchaincli tx wasm instantiate "$cw20_code_id1" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from admin18 -b block -y)
echo "instantiate cw20 succeed"
if [[ $(echo "$res" | jq '.logs[0].events[0].attributes[0].key' | sed 's/\"//g') != "_contract_address" ]];
then
  echo "unexpected result of instantiate"
  exit 1
fi;

echo "## instantiate nobody..."
res=$(exchaincli tx wasm instantiate "$cw20_code_id2" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from captain -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="unauthorized: can not instantiate: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  exit 1
fi;
res=$(exchaincli tx wasm instantiate "$cw20_code_id2" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from admin18 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="unauthorized: can not instantiate: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  exit 1
fi;

echo "## instantiate only address..."
res=$(exchaincli tx wasm instantiate "$cw20_code_id3" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from captain -b block -y)
echo "instantiate cw20 succeed"
if [[ $(echo "$res" | jq '.logs[0].events[0].attributes[0].key' | sed 's/\"//g') != "_contract_address" ]];
then
  echo "unexpected result of instantiate"
  exit 1
fi;
res=$(exchaincli tx wasm instantiate "$cw20_code_id3" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from admin18 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="unauthorized: can not instantiate: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  exit 1
fi;

echo "## instantiate nonexistent contract..."
res=$(exchaincli tx wasm instantiate 9999 '{"decimals":10,"initial_balances":[{"address":"","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from captain -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="not found: code: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo $res | jq
  echo "expect fail when instantiate nonexistent contract"
  exit 1
fi;

echo "## instantiate cw20 contract with invalid input..."
res=$(exchaincli tx wasm instantiate "$cw20_code_id5" '{"decimals":10,"initial_balances":[{"address":"","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from captain -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="instantiate wasm contract failed: Generic error: addr_validate errored: Input is empty: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when instantiate contract with invalid parameters"
  exit 1
fi;

echo "## instantiate cw20 contract with invalid amount..."
res=$(exchaincli tx wasm instantiate "$cw20_code_id5" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --amount=1000000000000okt --fees 0.001okt --from captain -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log_prefix="insufficient funds"
if [[ "${raw_log:0:18}" != "${failed_log_prefix}" ]];
then
  echo "expect fail when instantiate contract with invalid amount"
  exit 1
fi;

echo "## instantiate cw20 contract..."
res=$(exchaincli tx wasm instantiate "$cw20_code_id5" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"100000000"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --fees 0.001okt --from captain -b block -y)
echo "instantiate cw20 succeed"
if [[ $(echo "$res" | jq '.logs[0].events[0].attributes[0].key' | sed 's/\"//g') != "_contract_address" ]];
then
  echo "unexpected result of instantiate"
  exit 1
fi;

echo "## instantiate cw20 contract with deposit..."
totalAmount="100000000"
depositAmount="800"
depositDenom="okt"
res=$(exchaincli tx wasm instantiate "$cw20_code_id5" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"'${totalAmount}'"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --amount=${depositAmount}${depositDenom} --fees 0.001okt --from captain -b block -y)
echo "instantiate cw20 succeed"
if [[ $(echo "$res" | jq '.logs[0].events[0].attributes[0].key' | sed 's/\"//g') != "_contract_address" ]];
then
  echo "unexpected result of instantiate with deposit"
  exit 1
fi;
instantiate_gas_used=$(echo "$res" | jq '.gas_used' | sed 's/\"//g')
cw20contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "cw20 contract address: $cw20contractAddr"
res=$(exchaincli query account "$cw20contractAddr")
balanceAmount=$(echo "$res" | jq '.value.coins[0].amount' | sed 's/\"//g')
balanceAmount=${balanceAmount%.*}
# shellcheck disable=SC2053
if [[ ${balanceAmount} != ${depositAmount} ]];
then
  echo "invalid balance amount"
  exit 1
fi;
balanceDenom=$(echo "$res" | jq '.value.coins[0].denom' | sed 's/\"//g')
# shellcheck disable=SC2053
if [[ ${balanceDenom} != ${depositDenom} ]];
then
  echo "invalid balance denom"
  exit 1
fi;
cw20_balance=$(exchaincli query wasm contract-state smart "$cw20contractAddr" '{"balance":{"address":"'$captain'"}}' | jq '.data.balance' | sed 's/\"//g')
# shellcheck disable=SC2053
if [[ ${cw20_balance} != ${totalAmount} ]];
then
  echo "invalid cw20 balance"
fi;

#####################################################
#############    execute contract     ###############
#####################################################

transferAmount="100"
echo "## cw20 transfer to invalid recipient..."
res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"'$transferAmount'","recipient":""}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="execute wasm contract failed: Generic error: addr_validate errored: Input is empty: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when cw20 transfer to invalid recipient"
  exit 1
fi;

echo "## cw20 transfer..."
res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"100","recipient":"'$admin18'"}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
standard_gas_used=$(echo "$res" | jq '.gas_used' | sed 's/\"//g')
echo "standard_gas_used:$standard_gas_used"

echo "## cw20 transfer with okt transfer..."
res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"100","recipient":"'$admin18'"}}'  --amount=${depositAmount}${depositDenom} --fees 0.001okt --gas 2000000 --from captain -b block -y)
gas_used2=$(echo "$res" | jq '.gas_used' | sed 's/\"//g')
echo "gas_used2:$gas_used2"
if [[ "$standard_gas_used" -ge "$gas_used2" ]];
then
  echo "unexpected execute gas used2"
  exit 1
fi;

echo "## pin cw20 code..."
res=$(exchaincli tx gov submit-proposal pin-codes "$cw20_code_id5" --deposit 10.1okt --title "test title" --description "test description" --fees 0.001okt --from captain --gas=3000000 -b block -y)
proposal_id=$(echo "$res" | jq '.logs[0].events[1].attributes[1].value' | sed 's/\"//g')
res=$(exchaincli tx gov deposit "$proposal_id" 90000000okt --fees 0.001okt --from captain -b block -y)
res=$(exchaincli tx gov vote "$proposal_id" yes --fees 0.001okt --from captain -b block -y)

total_pinned=$(exchaincli query wasm pinned | jq '.code_ids|length')
if [[ $total_pinned -ne 1 ]];
then
  echo "unexpected total pinned: $total_pinned"
  exit 1
fi;

res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"100","recipient":"'$admin18'"}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
gas_used3=$(echo "$res" | jq '.gas_used' | sed 's/\"//g')
echo "gas_used3:$gas_used3"
if [[ "$standard_gas_used" -le "$gas_used3" ]];
then
  echo "unexpected execute gas used3"
  exit 1
fi;

res=$(exchaincli tx wasm instantiate "$cw20_code_id5" '{"decimals":10,"initial_balances":[{"address":"'$captain'","amount":"'${totalAmount}'"}],"name":"my test token", "symbol":"mtt"}' --label test1 --admin "$captain" --amount=${depositAmount}${depositDenom} --fees 0.001okt --from captain -b block -y)
instantiate_gas_used2=$(echo "$res" | jq '.gas_used' | sed 's/\"//g')
if [[ "$instantiate_gas_used" -le "$instantiate_gas_used2" ]];
then
  echo "unexpected instantiate gas_used2"
  exit 1
fi;

res=$(exchaincli tx wasm store ./wasm/cw4-stake/artifacts/cw4_stake.wasm --fees 0.01okt --from admin18 --gas=2000000 -b block -y)
echo "store cw4-stake succeed"
cw4_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
res=$(exchaincli tx wasm instantiate "$cw4_code_id" '{"denom":{"cw20":"'$cw20contractAddr'"},"min_bond":"100","tokens_per_weight":"10","unbonding_period":{"height":100}}' --label test1 --admin $captain --fees 0.001okt --from captain -b block -y)
echo "instantiate cw4-stake succeed"
cw4contractAddr=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
echo "cw4-stake contractAddr: $cw4contractAddr"
addr=$(exchaincli query wasm contract-state smart "$cw4contractAddr" '{"staked":{"address":"'$captain'"}}' | jq '.data.denom.cw20' | sed 's/\"//g')
# shellcheck disable=SC2053
if [[ $addr != $cw20contractAddr ]];
then
  echo "unexpected addr"
  exit 1
fi;

sendAmount="100"
res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"send":{"amount":"'$sendAmount'","contract":"'$cw4contractAddr'","msg":"eyJib25kIjp7fX0="}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
cw4balance=$(exchaincli query wasm contract-state smart "$cw20contractAddr" '{"balance":{"address":"'$cw4contractAddr'"}}' | jq '.data.balance' | sed 's/\"//g')
if [[ $cw4balance -ne $sendAmount ]];
then
  echo "unexpected cw4 contract balance"
  exit 1
fi;
cw4stake=$(exchaincli query wasm contract-state smart "$cw4contractAddr" '{"staked":{"address":"'$captain'"}}' | jq '.data.stake' | sed 's/\"//g')
if [[ $cw4stake -ne $sendAmount ]];
then
  echo "unexpected cw4 contract stake"
  exit 1
fi;

echo "## unpin cw20 code..."
res=$(exchaincli tx gov submit-proposal unpin-codes "$cw20_code_id5" --deposit 10.1okt --title "test title" --description "test description" --fees 0.001okt --from captain --gas=3000000 -b block -y)
proposal_id=$(echo "$res" | jq '.logs[0].events[1].attributes[1].value' | sed 's/\"//g')
res=$(exchaincli tx gov deposit "$proposal_id" 90000000okt --fees 0.001okt --from captain -b block -y)
res=$(exchaincli tx gov vote "$proposal_id" yes --fees 0.001okt --from captain -b block -y)

sleep 1
total_pinned=$(exchaincli query wasm pinned | jq '.code_ids|length')
if [[ $total_pinned -ne 0 ]];
then
  exit 1
fi;

res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"100","recipient":"'$admin18'"}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
gas_used4=$(echo "$res" | jq '.gas_used' | sed 's/\"//g')
echo "gas_used4:$gas_used4"
if [[ "$gas_used3" -ge "$gas_used4" ]];
then
  echo "unexpected execute gas used4"
  exit 1
fi;

#####################################################
#############    update&clear admin   ###############
#####################################################
echo "## update admin..."
res=$(exchaincli tx wasm set-contract-admin "$cw4contractAddr" "$admin18" --fees 0.001okt --from admin17 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="unauthorized: can not modify contract: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when update admin by other address"
  exit 1
fi;

res=$(exchaincli tx wasm set-contract-admin "$cw4contractAddr" "$admin17" --fees 0.001okt --from captain -b block -y)
actionName=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
if [[ "${actionName}" != "update-contract-admin" ]];
then
  echo "invalid action name"
  exit 1
fi;

res=$(exchaincli tx wasm set-contract-admin "$cw4contractAddr" "$admin18" --fees 0.001okt --from admin17 -b block -y)
actionName=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
if [[ "${actionName}" != "update-contract-admin" ]];
then
  echo "invalid action name"
  exit 1
fi;

echo "## clear admin..."
res=$(exchaincli tx wasm clear-contract-admin "$cw4contractAddr" --fees 0.001okt --from admin18 -b block -y)
actionName=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')
if [[ "${actionName}" != "clear-contract-admin" ]];
then
  echo "invalid action name: ${actionName}"
  exit 1
fi;

res=$(exchaincli tx wasm set-contract-admin "$cw4contractAddr" "$admin17" --fees 0.001okt --from admin18 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="unauthorized: can not modify contract: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when update admin after clear admin"
  exit 1
fi;

#####################################################
#############    migrate contract     ###############
#####################################################
res=$(exchaincli tx wasm store ./wasm/test/burner.wasm --fees 0.01okt --from admin18 --gas=2000000 -b block -y)
echo "store burner succeed"
burner_code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')

res=$(exchaincli tx wasm migrate "$cw20contractAddr" "$burner_code_id" "{}" --fees 0.01okt --from captain --gas=2000000 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="migrate wasm contract failed: Error parsing into type burner::msg::MigrateMsg: missing field \`payout\`: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when migrating with invalid parameters"
  exit 1
fi;

res=$(exchaincli tx wasm migrate "$cw20contractAddr" "$burner_code_id" '{"payout": "'$captain'"}' --fees 0.01okt --from admin18 --gas=2000000 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="unauthorized: can not migrate: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when migrating with address which is not admin"
  exit 1
fi;

res=$(exchaincli tx wasm migrate "$cw20contractAddr" "$burner_code_id" '{"payout": "'$captain'"}' --fees 0.01okt --from captain --gas=2000000 -b block -y)
new_code_id=$(exchaincli query wasm contract "$cw20contractAddr" | jq '.contract_info.code_id' | sed 's/\"//g')
if [[ $new_code_id -ne $burner_code_id ]];
then
  echo "migrate failed"
  exit 1
fi;

operation_name=$(exchaincli query wasm contract-history "$cw20contractAddr" | jq '.entries[1].operation' | sed 's/\"//g')
if [[ $operation_name != "CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE" ]];
then
  echo "migrate failed"
  exit 1
fi;

res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"100","recipient":"'$admin18'"}}'  --fees 0.001okt --gas 2000000 --from captain -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="execute wasm contract failed: Error calling the VM: Error resolving Wasm function: Could not get export: Missing export execute: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when execute after migrating contract"
  exit 1
fi;

res=$(exchaincli tx wasm migrate "$cw20contractAddr" "$burner_code_id" '{"payout": "'$captain'"}' --fees 0.01okt --from captain --gas=2000000 -b block -y)
new_code_id=$(exchaincli query wasm contract "$cw20contractAddr" | jq '.contract_info.code_id' | sed 's/\"//g')
if [[ $new_code_id -ne $burner_code_id ]];
then
  echo "migrate failed"
  exit 1
fi;

res=$(exchaincli tx wasm clear-contract-admin "$cw20contractAddr" --fees 0.001okt --from captain -b block -y)
res=$(exchaincli tx wasm migrate "$cw20contractAddr" "$burner_code_id" '{"payout": "'$captain'"}' --fees 0.01okt --from captain --gas=2000000 -b block -y)
raw_log=$(echo "$res" | jq '.raw_log' | sed 's/\"//g')
failed_log="unauthorized: can not migrate: failed to execute message; message index: 0"
if [[ "${raw_log}" != "${failed_log}" ]];
then
  echo "expect fail when migrating after clearing admin"
  exit 1
fi;

res=$(exchaincli tx gov submit-proposal migrate-contract "$cw20contractAddr" "$burner_code_id" '{"payout": "'$admin18'"}' --deposit 10.1okt --title "test title" --description "test description" --fees 0.01okt --from captain --gas=3000000 -b block -y)
proposal_id=$(echo "$res" | jq '.logs[0].events[1].attributes[1].value' | sed 's/\"//g')
echo "proposal_id: $proposal_id"
res=$(exchaincli tx gov deposit "$proposal_id" 90000000okt --fees 0.001okt --gas 3000000 --from captain -b block -y)
res=$(exchaincli tx gov vote "$proposal_id" yes --fees 0.001okt --gas 3000000 --from captain -b block -y)
result=$(echo "$res" | jq '.logs[0].events[2].attributes[2].value' | sed 's/\"//g')
if [[ $result != "passed" ]];
then
  echo "migrate by gov failed"
  exit 1
fi;

echo "migrate by gov succeed"



echo "$cw20_code_id1" "$cw20_code_id2" "$cw20_code_id3" "$cw20_code_id4" "$cw20_code_id5"


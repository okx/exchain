res=$(exchaincli tx wasm store ../cycle-counter/artifacts/counter.wasm --fees 0.01okt --from captain --gas=3000000 -b block -y -o json)
code_id=$(echo "$res" | jq -r '.logs[0].events[1].attributes[0].value')
#res=$(exchaincli tx gov submit-proposal pin-codes $code_id --deposit 100okt --title "test title" --description "test description" --fees 0.001okt --from captain --gas=3000000 -b block -y)
#proposal_id=$(echo "$res" | jq -r '.logs[0].events[1].attributes[1].value')
#exchaincli tx gov vote $proposal_id yes -y -b block --fees 0.004okt --gas 2000000 --from captain | jq
res=$(exchaincli tx wasm instantiate $code_id '{"count":"0"}' --label counter_Uint128 --admin 0xbbE4733d85bc2b90682147779DA49caB38C0aA1F --from captain --fees 0.01okt --gas 3000000 -y -b block)
contract_addr=$(echo "$res" | jq -r '.logs[0].events[0].attributes[0].value')
echo "contract_addr: $contract_addr"
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"1"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"10"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"100"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"1000"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"10000"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"100000"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"1000000"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"10000000"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli tx wasm execute ${contract_addr} '{"increment":{"count":"100000000"}}' --from captain --fees 0.01okt --gas 50000000 -y -b block | jq -r '.gas_used'
exchaincli query wasm contract-state smart ${contract_addr} '{"get_count":{}}' | jq


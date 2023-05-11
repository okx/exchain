#!/bin/bash

# config
privateKey=8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17

# deploy evm contract
evm_contract=$(go run main.go utils.go --action deploy --key $privateKey)
echo " ========================================================== "
echo "## deploy evm contract ##"
echo
echo "contarct address is $evm_contract"
echo

# deploy wasm contract
res=$(exchaincli tx wasm store ./wasmContract/counter.wasm --fees 0.01okt --from captain --gas=2000000 -b block -y)
code_id=$(echo "$res" | jq '.logs[0].events[1].attributes[0].value' | sed 's/\"//g')
res=$(exchaincli tx wasm instantiate "$code_id" '{}' --label test1 --admin ex1h0j8x0v9hs4eq6ppgamemfyu4vuvp2sl0q9p3v --fees 0.001okt --from captain -b block -y)
wasm_contract=$(echo "$res" | jq '.logs[0].events[0].attributes[0].value' | sed 's/\"//g')


echo " ========================================================== "
echo "## deploy wasm contract ##"
echo
echo "contarct address is $wasm_contract"
echo


# query evm state
res=$(go run main.go utils.go --action query --contract $evm_contract --key $privateKey)
echo " ========================================================== "
echo "## Query count in evm contarct ##"
echo
echo " The "count" is $res"
echo


# wasm call evm

res=$(exchaincli tx wasm execute "$wasm_contract" '{"add_counter_for_evm":{"evm_contract":"'$evm_contract'","delta":"1"}}' --fees 0.001okt --from captain -b block -y)
res=${res#*txhash}
res=${res:3:67}
echo " ========================================================== "
echo "## send a VM bridge tx to wasm contract ##"
echo
echo "tx hash $res "
echo





# check evm state
res=$(go run main.go utils.go --action query --contract $evm_contract --key $privateKey)
echo " ========================================================== "
echo "## Query count in evm contarct ##"
echo
echo "The "count" changed to $res"
echo



# check wasm state
res=$(exchaincli query wasm contract-state smart "$wasm_contract" '{"get_counter":{}}')
echo " ========================================================== "
echo "## Query count in wasm contarct ##"
echo
echo ""count" is $res "
echo

# evm call wasm

res=$(go run main.go utils.go --action execute --contract $evm_contract --wasmContract $wasm_contract --key $privateKey)
echo " ========================================================== "
echo "## send a VM bridge tx to evm contract ##"
echo
echo $res
echo


# check wasm state
res=$(exchaincli query wasm contract-state smart "$wasm_contract" '{"get_counter":{}}')
echo " ========================================================== "
echo "## Query count in wasm contarct ##"
echo
echo "The "count" changed to $res"
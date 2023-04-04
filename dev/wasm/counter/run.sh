#!/bin/bash

OPTIONS="--from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 6.1 -y"
OPTIONS="--from captain --fees 0.01okt --gas 6000000 -b block -y"

#exchaincli tx wasm store ./counter.wasm --from captain --gas-prices 0.00000001okt --gas auto -b block --gas-adjustment 1.1 -y
exchaincli tx wasm store ./counter.wasm --from captain --fees 0.01okt --gas 6000000 -b block -y
#
exchaincli tx wasm instantiate 1 '{}' --label test --from captain --fees 0.01okt --gas 6000000 --no-admin -b block -y
#
exchaincli tx wasm execute ex14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s6fqu27  \
'{"add":{"delta":"16"}}' \
--from captain --fees 0.01okt --gas 6000000 -b block -y

##
##exchaincli tx wasm execute  ${OPTIONS}
##
exchaincli query wasm contract-state smart ex14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s6fqu27 '{"get_counter":{}}'


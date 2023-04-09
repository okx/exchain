#!/bin/bash

OPTIONS="--from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.1 -y"

exchaincli tx wasm store ./counter.wasm ${OPTIONS}
exchaincli tx wasm instantiate 1 '{}' ${OPTIONS}

exchaincli tx wasm execute ex14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s6fqu27 '{"add":{"delta":"16"}}'  ${OPTIONS}

##
##exchaincli tx wasm execute  ${OPTIONS}
##
exchaincli query wasm contract-state smart ex14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s6fqu27 '{"get_counter":{}}'


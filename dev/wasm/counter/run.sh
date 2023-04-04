#!/bin/bash

OPTIONS="--from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 6.1 -y"
OPTIONS="--from captain --fees 0.01okt --gas 6000000 -b block -y"

#exchaincli tx wasm store ./counter.wasm --from captain --gas-prices 0.00000001okt --gas auto -b block --gas-adjustment 1.1 -y
exchaincli tx wasm store ./counter.wasm ${OPTIONS}
#
exchaincli tx wasm instantiate 1 --from captain --fees 0.01okt --gas 6000000 -b block -y
#
#exchaincli tx wasm execute  ${OPTIONS}
#
#exchaincli tx wasm execute  ${OPTIONS}
#
#exchaincli query wasm contract-state smart  '{"balance":{"address":"ex190227rqaps5nplhg2tg8hww7slvvquzy0qa0l0"}}'
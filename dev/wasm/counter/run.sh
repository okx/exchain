#!/bin/bash

OPTIONS="--from captain --gas-prices 0.00000001okt --gas auto -b block --gas-adjustment 1.1 -y"

exchaincli tx wasm store ./counter.wasm ${OPTIONS}

exchaincli tx wasm instantiate 1 ${OPTIONS}

exchaincli tx wasm execute  ${OPTIONS}

exchaincli tx wasm execute  ${OPTIONS}

exchaincli query wasm contract-state smart  '{"balance":{"address":"ex190227rqaps5nplhg2tg8hww7slvvquzy0qa0l0"}}'
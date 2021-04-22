#!/usr/bin/env sh

##
## Input parameters
##
ID=${ID:-0}
LOG=${LOG:-exchaind.log}

##
## Run binary with all parameters
##
export OKEXCHAINDHOME="/okexchaind/node${ID}/okexchaind"

if [ -d "$(dirname "${OKEXCHAINDHOME}"/"${LOG}")" ]; then
  exchaind --chain-id okexchain-1 --home "${OKEXCHAINDHOME}" "$@" | tee "${OKExCHAINDHOME}/${LOG}"
else
  exchaind --chain-id okexchain-1 --home "${OKEXCHAINDHOME}" "$@"
fi


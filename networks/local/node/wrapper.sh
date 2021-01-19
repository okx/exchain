#!/usr/bin/env sh

##
## Input parameters
##
ID=${ID:-0}
LOG=${LOG:-okexchaind.log}

##
## Run binary with all parameters
##
export OKEXCHAINDHOME="/okexchaind/node${ID}/okexchaind"

if [ -d "$(dirname "${OKEXCHAINDHOME}"/"${LOG}")" ]; then
  okexchaind --chain-id okexchain-1 --home "${OKEXCHAINDHOME}" "$@" | tee "${OKExCHAINDHOME}/${LOG}"
else
  okexchaind --chain-id okexchain-1 --home "${OKEXCHAINDHOME}" "$@"
fi


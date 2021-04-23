#!/usr/bin/env sh

##
## Input parameters
##
ID=${ID:-0}
LOG=${LOG:-exchaind.log}

##
## Run binary with all parameters
##
export OKEXCHAINDHOME="/exchaind/node${ID}/exchaind"

if [ -d "$(dirname "${OKEXCHAINDHOME}"/"${LOG}")" ]; then
  exchaind --chain-id exchain-1 --home "${OKEXCHAINDHOME}" "$@" | tee "${OKExCHAINDHOME}/${LOG}"
else
  exchaind --chain-id exchain-1 --home "${OKEXCHAINDHOME}" "$@"
fi


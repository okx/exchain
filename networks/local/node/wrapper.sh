#!/usr/bin/env sh

##
## Input parameters
##
ID=${ID:-0}
LOG=${LOG:-exchaind.log}

##
## Run binary with all parameters
##
export EXCHAINDHOME="/exchaind/node${ID}/exchaind"

if [ -d "$(dirname "${EXCHAINDHOME}"/"${LOG}")" ]; then
  exchaind --chain-id exchain-1 --home "${EXCHAINDHOME}" "$@" | tee "${EXCHAINDHOME}/${LOG}"
else
  exchaind --chain-id exchain-1 --home "${EXCHAINDHOME}" "$@"
fi


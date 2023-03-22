#!/usr/bin/env sh

##
## Input parameters
##
ID=${ID:-0}
LOG=${LOG:-okbchaind.log}

##
## Run binary with all parameters
##
export OKBCHAINDHOME="/okbchaind/node${ID}/okbchaind"

if [ -d "$(dirname "${OKBCHAINDHOME}"/"${LOG}")" ]; then
  okbchaind --chain-id okbchain-1 --home "${OKBCHAINDHOME}" "$@" | tee "${OKBCHAINDHOME}/${LOG}"
else
  okbchaind --chain-id okbchain-1 --home "${OKBCHAINDHOME}" "$@"
fi


#!/bin/bash

KEY="captain"
CHAINID="exchain-65"
MONIKER="oec"
CURDIR=`dirname $0`
HOME_SERVER=$CURDIR/"_cache_evm"

set -e
set -o errexit
set -a
set -m




run() {
    LOG_LEVEL=main:info,state:error,distr:error,auth:error,mint:error,farm:error
    #LOG_LEVEL=main:info,state:error,distr:error,auth:error,mint:error,farm:error,perf:info

    exchaind start --pruning=nothing --rpc.unsafe \
      --local-rpc-port 26657 \
      --log_level $LOG_LEVEL \
      --consensus.timeout_commit 3s \
      --trace --home $HOME_SERVER --chain-id $CHAINID \
      --rest.laddr "tcp://localhost:8545" > oec.log 2>&1 &

    exit
}


run
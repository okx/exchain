#!/bin/bash
set -o errexit -o nounset -o pipefail
CHAIN_ID="exchain-67"
NODE="http://localhost:26657"
TX_EXTRA="--fees 0.01okt --gas 3000000 --chain-id=$CHAIN_ID --node $NODE -b block -y"

cw20contractAddr="ex1zwv6feuzhy6a9wekh96cd57lsarmqlwxdypdsplw6zhfncqw6ftqwe39pr"


# Usage example:
#  ./transfer.sh ex1h0j8x0v9hs4eq6ppgamemfyu4vuvp2sl0q9p3v ex190227rqaps5nplhg2tg8hww7slvvquzy0qa0l0 1
fromAddr=$1
toAddr=$2
amount=$3

res=$(exchaincli tx wasm execute "$cw20contractAddr" '{"transfer":{"amount":"'$amount'","recipient":"'$toAddr'"}}' --from $fromAddr $TX_EXTRA)
echo $res | jq '.txhash' | sed 's/\"//g'
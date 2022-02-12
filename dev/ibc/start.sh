#!/bin/bash

KEY="captain"
CHAINID="exchain-67"
MONIKER="oec"
CURDIR=`dirname $0`

set -e
set -o errexit
set -a
set -m



CHAINID=$1
CHAINDIR=$2
RPCPORT=$3
P2PPORT=$4
RESTPORT=$5
#PROFPORT=$5
#GRPCPORT=$6



echorun() {
  echo "------------------------------------------------------------------------------------------------"
  echo "["$@"]"
  $@
  echo "------------------------------------------------------------------------------------------------"
}

rm -rf ~/.exchaind
rm -rf ~/.exchaincli


run() {
    LOG_LEVEL=main:debug,iavl:info,*:error,state:info,provider:info

     nohup exchaind start --pruning=nothing --rpc.unsafe \
      --home=$CHAINDIR/$CHAINID \
      --local-rpc-port ${RPCPORT} \
      --rpc.laddr="tcp://0.0.0.0:${RPCPORT}" \
      --rpc.external_laddr="0.0.0.0:${RPCPORT}" \
      --p2p.laddr="tcp://0.0.0.0:${P2PPORT}" \
      --log_level $LOG_LEVEL \
      --log_file json \
      --enable-dynamic-gp=false \
      --consensus.timeout_commit 2000ms \
      --enable-preruntx=false \
      --iavl-enable-async-commit \
      --enable-gid \
      --append-pid=true \
      --iavl-commit-interval-height 10 \
      --iavl-output-modules evm=0,acc=0 \
      --trace --home=$CHAINDIR/$CHAINID --chain-id $CHAINID \
      --rest.laddr "tcp://localhost:${RESTPORT}" \
      --elapsed Round=1,CommitRound=1,Produce=1 >$CHAINDIR/$CHAINID/oec.txt 2>&1 &
    exit
}




set -x # activate debugging

# run

# remove existing daemon and client
#rm -rf ~/.exchain*
#rm -rf $HOME_SERVER

(cd ../../ && make install VenusHeight=1)

# Set up config for CLI
#exchaincli config --home $CHAINDIR/$CHAINID chain-id $CHAINID
#exchaincli config --home $CHAINDIR/$CHAINID output json
#exchaincli config --home $CHAINDIR/$CHAINID indent true
#exchaincli config --home $CHAINDIR/$CHAINID trust-node true
#exchaincli config --home $CHAINDIR/$CHAINID keyring-backend test

HOME_SERVER=~/.exchaind

redirect() {
  if [ "$SILENT" -eq 1 ]; then
    "$@" > /dev/null 2>&1
  else
    "$@"
  fi
}


exchaincli config  chain-id $CHAINID
exchaincli config   output json
exchaincli config   indent true
exchaincli config   trust-node true
exchaincli config   keyring-backend test


redirect  exchaind init $MONIKER --chain-id $CHAINID

# if $KEY exists it should be deleted
#
#    "eth_address": "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F",
#     prikey: 8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17
exchaincli keys add --recover captain -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" --output json > ${HOME_SERVER}/validator_seed.json 2>&1

#    "eth_address": "0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0",
exchaincli keys add  --recover admin16 -m "palace cube bitter light woman side pave cereal donor bronze twice work" --output json > ${HOME_SERVER}/key_seed.json 2>&1

exchaincli keys add --recover admin17 -m "antique onion adult slot sad dizzy sure among cement demise submit scare" -y

exchaincli keys add  --recover admin18 -m "lazy cause kite fence gravity regret visa fuel tone clerk motor rent" -y

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
#

#exchaind init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to okt
cat $HOME_SERVER/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json

# Enable EVM

if [ "$(uname -s)" == "Darwin" ]; then
    sed -i "" 's/"enable_call": false/"enable_call": true/' $HOME_SERVER/config/genesis.json
    sed -i "" 's/"enable_create": false/"enable_create": true/' $HOME_SERVER/config/genesis.json
    sed -i "" 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' $HOME_SERVER/config/genesis.json
else
    sed -i 's/"enable_call": false/"enable_call": true/' $HOME_SERVER/config/genesis.json
    sed -i 's/"enable_create": false/"enable_create": true/' $HOME_SERVER/config/genesis.json
    sed -i 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' $HOME_SERVER/config/genesis.json
fi

# Allocate genesis accounts (cosmos formatted addresses)
exchaind add-genesis-account $(exchaincli keys show $KEY    -a ) 100000000okt
exchaind add-genesis-account $(exchaincli keys show admin16 -a ) 900000000okt
exchaind add-genesis-account $(exchaincli keys show admin17 -a ) 900000000okt
exchaind add-genesis-account $(exchaincli keys show admin18 -a ) 900000000okt

# Sign genesis transaction
#rm -rf ~/.exchaind
#mv $CHAINDIR/$CHAINID ~/.exchaind
#exchaind gentx --name $KEY --keyring-backend test --home $CHAINDIR/$CHAINID
exchaind gentx --name $KEY --keyring-backend test

# Collect genesis tx
exchaind collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
exchaind validate-genesis
exchaincli config keyring-backend test

mv ~/.exchaind $CHAINDIR/$CHAINID
run


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


killbyname() {
  NAME=$1
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2", "$8}'
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2}' | sh
  echo "All <$NAME> killed!"
}


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


killbyname exchaind
killbyname exchaincli

set -x # activate debugging

# run

# remove existing daemon and client
rm -rf ~/.exchain*
rm -rf $HOME_SERVER

(cd .. && make install)

# Set up config for CLI
exchaincli config chain-id $CHAINID
exchaincli config output json
exchaincli config indent true
exchaincli config trust-node true
exchaincli config keyring-backend test

# if $KEY exists it should be deleted
#
#    "eth_address": "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F",
exchaincli keys add --recover captain -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -y

#    "eth_address": "0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0",
exchaincli keys add --recover admin16 -m "palace cube bitter light woman side pave cereal donor bronze twice work" -y

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
exchaind init $MONIKER --chain-id $CHAINID --home $HOME_SERVER

# Change parameter token denominations to okt
cat $HOME_SERVER/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="okt"' > $HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json

# Enable EVM
sed -i "" 's/"enable_call": false/"enable_call": true/' $HOME_SERVER/config/genesis.json
sed -i "" 's/"enable_create": false/"enable_create": true/' $HOME_SERVER/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
exchaind add-genesis-account $(exchaincli keys show $KEY    -a) 100000000okt --home $HOME_SERVER
exchaind add-genesis-account $(exchaincli keys show admin16 -a) 900000000okt --home $HOME_SERVER

# Sign genesis transaction
exchaind gentx --name $KEY --keyring-backend test --home $HOME_SERVER

# Collect genesis tx
exchaind collect-gentxs --home $HOME_SERVER

# Run this to ensure everything worked and that the genesis file is setup correctly
exchaind validate-genesis --home $HOME_SERVER
exchaincli config keyring-backend test

run
#!/bin/bash

KEY1="alice"
KEY2="bob"
CHAINID="exchainevm-65"
MONIKER="okex"
CURDIR=$(dirname $0)
HOME_BASE=$CURDIR/"_cache_evm"
HOME_SERVER=$HOME_BASE/".exchaind"
HOME_CLI=$HOME_BASE/".exchaincli"

set -e

function killexchaind() {
  ps -ef | grep "exchaind" | grep -v grep | grep -v run.sh | awk '{print "kill -9 "$2", "$8}'
  ps -ef | grep "exchaind" | grep -v grep | grep -v run.sh | awk '{print "kill -9 "$2}' | sh
  echo "All <exchaind> killed!"
}

killexchaind

# remove existing daemon and client
rm -rf $HOME_BASE

cd ../../../
make install
cd ./app/rpc/tests

exchaincli config keyring-backend test --home $HOME_CLI

# Set up config for CLI
exchaincli config chain-id $CHAINID --home $HOME_CLI
exchaincli config output json --home $HOME_CLI
exchaincli config indent true --home $HOME_CLI
exchaincli config trust-node true --home $HOME_CLI
# if $KEY exists it should be deleted
exchaincli keys add $KEY1 --recover -m "plunge silk glide glass curve cycle snack garbage obscure express decade dirt" --home $HOME_CLI
exchaincli keys add $KEY2 --recover -m "lazy cupboard wealth canoe pumpkin gasp play dash antenna monitor material village" --home $HOME_CLI

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
exchaind init $MONIKER --chain-id $CHAINID --home $HOME_SERVER

# Change parameter token denominations to okt
cat $HOME_SERVER/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="okt"' >$HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="okt"' >$HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="okt"' >$HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json
cat $HOME_SERVER/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="okt"' >$HOME_SERVER/config/tmp_genesis.json && mv $HOME_SERVER/config/tmp_genesis.json $HOME_SERVER/config/genesis.json

# Enable EVM
sed -i "" 's/"enable_call": false/"enable_call": true/' $HOME_SERVER/config/genesis.json
sed -i "" 's/"enable_create": false/"enable_create": true/' $HOME_SERVER/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
exchaind add-genesis-account $(exchaincli keys show $KEY1 -a --home $HOME_CLI) 1000000000okt --home $HOME_SERVER
exchaind add-genesis-account $(exchaincli keys show $KEY2 -a --home $HOME_CLI) 1000000000okt --home $HOME_SERVER
## Sign genesis transaction
exchaind gentx --name $KEY1 --keyring-backend test --home $HOME_SERVER --home-client $HOME_CLI
# Collect genesis tx
exchaind collect-gentxs --home $HOME_SERVER
# Run this to ensure everything worked and that the genesis file is setup correctly
exchaind validate-genesis --home $HOME_SERVER

LOG_LEVEL=main:info,state:info,distr:debug,auth:info,mint:debug,farm:debug

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)

# start node with web3 rest
exchaind start \
  --pruning=nothing \
  --rpc.unsafe \
  --rest.laddr tcp://0.0.0.0:8545 \
  --chain-id $CHAINID \
  --log_level $LOG_LEVEL \
  --trace \
  --home $HOME_SERVER \
  --rest.unlock_key $KEY1,$KEY2 \
  --rest.unlock_key_home $HOME_CLI \
  --keyring-backend "test" \
  --minimum-gas-prices "0.000000001okt"

#go test ./

# cleanup
#killexchaind
#rm -rf $HOME_BASE

exit

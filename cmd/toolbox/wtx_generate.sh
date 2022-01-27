#!/usr/bin/env bash

BIN_NAME=exchaind
OKCHAIN_TOP=${GOPATH}/src/github.com/okex/exchain
OKCHAIN_BIN=${OKCHAIN_TOP}/build
OKCHAIN_BIN=${GOPATH}/bin
OKCHAIN_NET_TOP=`pwd`
OKCHAIN_NET_CACHE=${OKCHAIN_NET_TOP}/cache
CHAIN_ID="exchain-67"


BASE_PORT_PREFIX=26600
P2P_PORT_SUFFIX=56
RPC_PORT_SUFFIX=57
REST_PORT=8545
let BASE_PORT=${BASE_PORT_PREFIX}+${P2P_PORT_SUFFIX}
let seedp2pport=${BASE_PORT_PREFIX}+${P2P_PORT_SUFFIX}
let seedrpcport=${BASE_PORT_PREFIX}+${RPC_PORT_SUFFIX}
let seedrestport=${seedrpcport}+1

if [ -z ${IP} ]; then
  IP="127.0.0.1"
fi

exchaind testnet --v 5 -o workspace -l \
    --chain-id ${CHAIN_ID} \
    --node-dir-prefix n \
    --starting-ip-address ${IP} \
    --base-port ${BASE_PORT} \
    --keyring-backend test \
    --mnemonic=true
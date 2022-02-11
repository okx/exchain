#!/usr/bin/env bash

NUM_NODE=4

# tackle size chronic goose deny inquiry gesture fog front sea twin raise
# acid pulse trial pill stumble toilet annual upgrade gold zone void civil
# antique onion adult slot sad dizzy sure among cement demise submit scare
# lazy cause kite fence gravity regret visa fuel tone clerk motor rent
HARDCODED_MNEMONIC=true

set -e
set -o errexit
set -a
set -m

set -x # activate debugging

source oec.profile
WRAPPEDTX=false
PRERUN=false
NUM_RPC=0
WHITE_LIST=0b066ca0790f27a6595560b23bf1a1193f100797,\
3813c7011932b18f27f172f0de2347871d27e852,\
6ea83a21a43c30a280a3139f6f23d737104b6975,\
bab6c32fa95f3a54ecb7d32869e32e85a25d2e08,\
testnet-node-ids


while getopts "r:isn:b:p:c:Sxwk:" opt; do
  case $opt in
  i)
    echo "OKCHAIN_INIT"
    OKCHAIN_INIT=1
    ;;
  w)
    echo "WRAPPEDTX=$OPTARG"
    WRAPPEDTX=true
    ;;
  x)
    echo "PRERUN=$OPTARG"
    PRERUN=true
    ;;
  s)
    echo "OKCHAIN_START"
    OKCHAIN_START=1
    ;;
  k)
    echo "LOG_SERVER"
    LOG_SERVER="--log-server $OPTARG"
    ;;
  c)
    echo "Test_CASE"
    Test_CASE="--consensus-testcase $OPTARG"
    ;;
  n)
    echo "NUM_NODE=$OPTARG"
    NUM_NODE=$OPTARG
    ;;
  r)
    echo "NUM_RPC=$OPTARG"
    NUM_RPC=$OPTARG
    ;;
  b)
    echo "BIN_NAME=$OPTARG"
    BIN_NAME=$OPTARG
    ;;
  S)
    STREAM_ENGINE="analysis&mysql&localhost:3306,notify&redis&localhost:6379,kline&pulsar&localhost:6650"
    echo "$STREAM_ENGINE"
    ;;
  p)
    echo "IP=$OPTARG"
    IP=$OPTARG
    ;;
  \?)
    echo "Invalid option: -$OPTARG"
    ;;
  esac
done

echorun() {
  echo "------------------------------------------------------------------------------------------------"
  echo "["$@"]"
  $@
  echo "------------------------------------------------------------------------------------------------"
}

killbyname() {
  NAME=$1
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2", "$8}'
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2}' | sh
  echo "All <$NAME> killed!"
}

init() {
  killbyname ${BIN_NAME}

  (cd ${OKCHAIN_TOP} && make install VenusHeight=1)

  rm -rf cache

  echo "=================================================="
  echo "===== Generate testnet configurations files...===="
  echorun exchaind testnet --v $1 --r $2 -o cache -l \
    --chain-id ${CHAIN_ID} \
    --starting-ip-address ${IP} \
    --base-port ${BASE_PORT} \
    --keyring-backend test
}
recover() {
  killbyname ${BIN_NAME}
  (cd ${OKCHAIN_TOP} && make install VenusHeight=1)
  rm -rf cache
  cp -rf nodecache cache
}

run() {

  index=$1
  seed_mode=$2
  p2pport=$3
  rpcport=$4
  restport=$5
  p2p_seed_opt=$6
  p2p_seed_arg=$7


  if [ "$(uname -s)" == "Darwin" ]; then
      sed -i "" 's/"enable_call": false/"enable_call": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i "" 's/"enable_create": false/"enable_create": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i "" 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' cache/node${index}/exchaind/config/genesis.json
  else
      sed -i 's/"enable_call": false/"enable_call": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i 's/"enable_create": false/"enable_create": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' cache/node${index}/exchaind/config/genesis.json
  fi

  exchaind add-genesis-account 0xbbE4733d85bc2b90682147779DA49caB38C0aA1F 900000000okt --home cache/node${index}/exchaind
  exchaind add-genesis-account 0x4C12e733e58819A1d3520f1E7aDCc614Ca20De64 900000000okt --home cache/node${index}/exchaind
  exchaind add-genesis-account 0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0 900000000okt --home cache/node${index}/exchaind
  exchaind add-genesis-account 0x2Bd4AF0C1D0c2930fEE852D07bB9dE87D8C07044 900000000okt --home cache/node${index}/exchaind

  LOG_LEVEL=main:info,*:error,consensus:error,state:info

  echorun nohup exchaind start \
    --home cache/node${index}/exchaind \
    --p2p.seed_mode=$seed_mode \
    --p2p.allow_duplicate_ip \
    --enable-dynamic-gp=false \
    --enable-wtx=${WRAPPEDTX} \
    --mempool.node_key_whitelist ${WHITE_LIST} \
    --p2p.pex=false \
    --p2p.addr_book_strict=false \
    $p2p_seed_opt $p2p_seed_arg \
    --p2p.laddr tcp://${IP}:${p2pport} \
    --rpc.laddr tcp://${IP}:${rpcport} \
    --consensus.timeout_commit 600ms \
    --log_level ${LOG_LEVEL} \
    --chain-id ${CHAIN_ID} \
    --upload-delta=false \
    --enable-gid \
    --append-pid=true \
    ${LOG_SERVER} \
    --elapsed DeliverTxs=0,Round=1,CommitRound=1,Produce=1 \
    --rest.laddr tcp://localhost:$restport \
    --enable-preruntx=$PRERUN \
    --consensus-role=v$index \
    ${Test_CASE} \
    --keyring-backend test >cache/val${index}.log 2>&1 &

#     --iavl-enable-async-commit \    --consensus-testcase case12.json \
#     --upload-delta \
#     --enable-preruntx \
#     --mempool.node_key_whitelist="nodeKey1,nodeKey2" \
#    --mempool.node_key_whitelist ${WHITE_LIST} \
}

function start() {
  killbyname ${BIN_NAME}
  index=0

  echo "============================================"
  echo "=========== Startup seed node...============"
  ((restport = REST_PORT)) # for evm tx
  run $index true ${seedp2pport} ${seedrpcport} $restport
  seed=$(exchaind tendermint show-node-id --home cache/node${index}/exchaind)

  echo "============================================"
  echo "======== Startup validator nodes...========="
  for ((index = 1; index < ${1}; index++)); do
    ((p2pport = BASE_PORT_PREFIX + index * 100 + P2P_PORT_SUFFIX))
    ((rpcport = BASE_PORT_PREFIX + index * 100 + RPC_PORT_SUFFIX))  # for exchaincli
    ((restport = index * 100 + REST_PORT)) # for evm tx
    run $index false ${p2pport} ${rpcport} $restport --p2p.seeds ${seed}@${IP}:${seedp2pport}
  done
  echo "start node done"
}

if [ -z ${IP} ]; then
  IP="127.0.0.1"
fi

if [ ! -z "${OKCHAIN_INIT}" ]; then
  ((NUM_VAL=NUM_NODE-NUM_RPC))
  init ${NUM_VAL} ${NUM_RPC}
fi

#if [ ! -z "${OKCHAIN_RECOVER}" ]; then
#  recover ${NUM_NODE}
#fi

if [ ! -z "${OKCHAIN_START}" ]; then
  start ${NUM_NODE}
fi

#!/usr/bin/env bash

NUM_NODE=4

# tackle size chronic goose deny inquiry gesture fog front sea twin raise
# acid pulse trial pill stumble toilet annual upgrade gold zone void civil
# antique onion adult slot sad dizzy sure among cement demise submit scare
# lazy cause kite fence gravity regret visa fuel tone clerk motor rent
HARDCODED_MNEMONIC=false

set -e
set -o errexit
set -a
set -m

set -x # activate debugging

source exchain.profile

while getopts "isn:b:p:Sm" opt; do
  case $opt in
  i)
    echo "OKCHAIN_INIT"
    OKCHAIN_INIT=1
    ;;
  s)
    echo "OKCHAIN_START"
    OKCHAIN_START=1
    ;;
  n)
    echo "NUM_NODE=$OPTARG"
    NUM_NODE=$OPTARG
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
  m)
    echo "HARDCODED_MNEMONIC"
    HARDCODED_MNEMONIC=true
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

  (cd ${OKCHAIN_TOP} && make install)

  rm -rf cache

  echo "=================================================="
  echo "===== Generate testnet configurations files...===="
  echorun exchaind testnet --v $1 -o cache -l \
    --chain-id ${CHAIN_ID} \
    --starting-ip-address ${IP} \
    --base-port ${BASE_PORT} \
    --keyring-backend test \
    --mnemonic=${HARDCODED_MNEMONIC}
}

run() {

  index=$1
  seed_mode=$2
  p2pport=$3
  rpcport=$4
  p2p_seed_opt=$5
  p2p_seed_arg=$6
  parallel_run_tx=false

  if [ $index -eq 3 ];then
      parallel_run_tx=true
    else
      parallel_run_tx=false
    fi

  LOG_LEVEL=main:info,*:error,consensus:info,state:info

  if [ "$(uname -s)" == "Darwin" ]; then
      sed -i "" 's/"enable_call": false/"enable_call": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i "" 's/"enable_create": false/"enable_create": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i "" 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' cache/node${index}/exchaind/config/genesis.json
  else
      sed -i 's/"enable_call": false/"enable_call": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i 's/"enable_create": false/"enable_create": true/' cache/node${index}/exchaind/config/genesis.json
      sed -i 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' cache/node${index}/exchaind/config/genesis.json
  fi

  echorun nohup exchaind start \
    --home cache/node${index}/exchaind \
    --p2p.seed_mode=$seed_mode \
    --p2p.allow_duplicate_ip \
    --p2p.pex=false \
    --p2p.addr_book_strict=false \
    $p2p_seed_opt $p2p_seed_arg \
    --p2p.laddr tcp://${IP}:${p2pport} \
    --rpc.laddr tcp://${IP}:${rpcport} \
    --consensus.timeout_commit 200ms \
    --log_level ${LOG_LEVEL} \
    --chain-id ${CHAIN_ID} \
    --upload-delta \
    --elapsed DeliverTxs=0,Round=0,CommitRound=0,Produce=0 \
    --rest.laddr tcp://localhost:8545 \
    --enable-proactively-runtx=$parallel_run_tx \
    --prerun-testcase "./case.json" \
    --proactively-role=$index \
    --keyring-backend test >cache/val${index}.log 2>&1 &

#     --iavl-enable-async-commit \
#     --upload-delta \
#     --enable-proactively-runtx \
}

function start() {
  killbyname ${BIN_NAME}
  index=0

  echo "============================================"
  echo "=========== Startup seed node...============"
  run $index true ${seedp2pport} ${seedrpcport}
  seed=$(exchaind tendermint show-node-id --home cache/node${index}/exchaind)

  echo "============================================"
  echo "======== Startup validator nodes...========="
  for ((index = 1; index < ${1}; index++)); do

    ((p2pport = BASE_PORT_PREFIX + index * 100 + P2P_PORT_SUFFIX))
    ((rpcport = BASE_PORT_PREFIX + index * 100 + RPC_PORT_SUFFIX))
    run $index false ${p2pport} ${rpcport} --p2p.seeds ${seed}@${IP}:${seedp2pport}
  done
  echo "start node done"
}

if [ -z ${IP} ]; then
  IP="127.0.0.1"
fi

if [ ! -z "${OKCHAIN_INIT}" ]; then
  init ${NUM_NODE}
fi

if [ ! -z "${OKCHAIN_START}" ]; then
  start ${NUM_NODE}
fi

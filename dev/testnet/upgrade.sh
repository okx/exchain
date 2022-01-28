#!/usr/bin/env bash

NUM_NODE=4

set -e
set -o errexit
set -a
set -m

# set -x # activate debugging

source oec.profile
PRERUN=false

REST_PORT_MAP='{"val0":8545,"val1":8645,"val2":8745,"val3":8845,"rpc4":8945,"rpc5":9045}'
RPC_PORT_MAP='{"val0":26657,"val1":26757,"val2":26857,"val3":26957,"rpc4":27057,"rpc5":27157}'

function killbyname_gracefully() {
  NAME=$1
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill  "$2", "$8}'
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill  "$2}' | sh
  echo "All <$NAME> killed gracefully!"
}

function build_exchain() {
  version=$1
  (cd ../.. && git checkout dev && git pull && git checkout $version && make install)
  echo "exchaind version ////"
  exchaind version
}

function get_latest_height() {
  node=$1
  port=`echo $RPC_PORT_MAP | jq .$node`
  height=`exchaincli status --node http://${IP}:${port} | jq .sync_info.latest_block_height | awk '{ gsub(/"/,""); print $0 }'`
  echo $height
}

function get_tx_count_of_height() {
  node=$1
  height=$2
  port=`echo $REST_PORT_MAP   | jq .$node`
  hex_height=`printf "0x%x" $height`
  data_json='{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByNumber","params":["'${hex_height}'"],"id":1}'
  tx_count=`curl -X POST --data $data_json -H "Content-Type: application/json" http://${IP}:$port -s | jq .result | awk '{ gsub(/"/,""); print $0 }' `
  echo $tx_count
}

function check_block() {
  # node_name should like val1 rpc2 ....
  node_name=$1
  # extra_check is an optional parameter
  # This function also checks that new block contains tx when it be set to "tx"
  extra_check=$2

  echo "Check block: $node_name"
  i=0
  is_valid=0
  base_height=`get_latest_height $node_name`
  
  # To record the count of block that contains tx
  block_contains_tx=0
  while [ $i -le 600 ] 
  do
    latest_height=`get_latest_height $node_name`
    echo "latest_height,"$latest_height 
    latest_tx_count=`get_tx_count_of_height $node_name $latest_height`
    echo "tx count,"$latest_tx_count
    if [[ $latest_tx_count != "0x0" ]] ; then
      let block_contains_tx+=1
    fi

    if [[ $extra_check = "tx" ]] ; then
      # Need to get new blocks over 25 and make sure those blocks contains tx 
      if [ `expr $latest_height - $base_height` -gt 25 -a $block_contains_tx -gt 25 ] ;then
        echo "block_contains_tx ,"$block_contains_tx
        is_valid=1
        break
      fi
    else
      # Only check that the node generates blocks over 10
      if [ `expr $latest_height - $base_height` -gt 10 ] ;then
        is_valid=1
        break
      fi
    fi

    let i+=1
    sleep 1
    echo "Checking... $latest_height"
  done

  if [ $is_valid -eq 0 ] ;then
      echo "Check valid $node_name: Failed, not pass."
      exit 99
  else
      echo "Check valid $node_name: Successful, pass."
  fi
}

function check_block_all() {
  echo "Check all node block"
  check_block val0
  check_block val1
  check_block val2
  check_block val3
  check_block rpc4
  check_block rpc5
}

function send_tx() {
  echo "start sending tx ..."
  # (cd ../client/ && bash run.sh > /dev/null 2>&1 &)
  (cd ../client/ && bash run.sh > ./newrun.log 2>&1 &)
}

function start_node() {
  index=$1
  node_name=$2
  exchaind_opts=${@:3}

  if [[ $index == "0" ]] ; then
    p2pport=${seedp2pport}
    rpcport=${seedrpcport}
  else
    ((p2pport = BASE_PORT_PREFIX + index * 100 + P2P_PORT_SUFFIX))
    ((rpcport = BASE_PORT_PREFIX + index * 100 + RPC_PORT_SUFFIX))
  fi
  ((restport = index * 100 + REST_PORT))

  LOG_LEVEL=main:info,*:error,consensus:error,state:info,provider:info
  
  nohup ${BIN_NAME} start \
    --chain-id ${CHAIN_ID} \
    --home cache/node${index}/exchaind \
    --p2p.laddr tcp://${IP}:${p2pport} \
    --rpc.laddr tcp://${IP}:${rpcport} \
    --rest.laddr tcp://${IP}:${restport} \
    --log_level ${LOG_LEVEL} \
    --enable-gid \
    --append-pid=true \
    --p2p.addr_book_strict=false \
    --enable-preruntx=${PRERUN} \
    ${exchaind_opts} \
    > cache/${node_name}.log 2>&1 &
}

function add_val() {
  index=$1
  node_name=val${index}
  seed_addr=$(exchaind tendermint show-node-id --home cache/node0/exchaind)@${IP}:${seedp2pport}
  echo "add val >>> "$node_name

  exchaind_opts="--p2p.allow_duplicate_ip  --p2p.pex=false  --p2p.addr_book_strict=false  --consensus.timeout_commit 600ms    --upload-delta=false  --elapsed DeliverTxs=0,Round=1,CommitRound=1,Produce=1  --consensus-role=v${index}  --p2p.seeds ${seed_addr} "

  start_node $index $node_name $exchaind_opts
}

function add_seed() {
  index=0
  node_name=val${index}
  echo "add seed >>> "$node_name

  exchaind_opts="--p2p.seed_mode=true  --p2p.allow_duplicate_ip  --p2p.pex=false  --p2p.addr_book_strict=false  --consensus.timeout_commit 600ms  --upload-delta=false  --elapsed DeliverTxs=0,Round=1,CommitRound=1,Produce=1  --consensus-role=v$index "

  start_node $index $node_name $exchaind_opts
}

function add_rpc() {
  index=$1
  node_name=rpc${index}
  echo "add rpc >>> "$node_name
  
  seed_addr=$(exchaind tendermint show-node-id --home cache/node0/exchaind)@${IP}:${seedp2pport}
  echo $seed_addr

  exchaind_opts="--p2p.seeds ${seed_addr} "
  start_node $index $node_name $exchaind_opts
}

function case_prepare() {
  # Prepare 4 validators and 2 rpc nodes
  version1=$1

  killbyname_gracefully ${BIN_NAME}
  killbyname_gracefully "run.sh"
  killbyname_gracefully "./client"

  bash testnet.sh -i
  build_exchain $version1
  bash testnet.sh -s -n 4
  bash addnewnode.sh -n 4
  bash addnewnode.sh -n 5
}

function caseopt() {
  echo "caseopt()"
  version1=$1
  version2=$2

  case_1 $version1 $version2
  case_2 $version1 $version2
  case_3 $version1 $version2
  echo "All cases finished!"
}


function case_1() {
  # Upgrade 1 rpc node , then upgrade 1 validator node.
  echo "[][][][][][][][][][][][][][][][][][]"
  echo "[][][][][]    case_1      [][][][][]"
  echo "[][][][][][][][][][][][][][][][][][]"

  version1=$1
  version2=$2

  # pre
  case_prepare $version1
  # extend opts below....

  #STEP sleep
  sleep 20

  check_block_all

  #STEP send tx
  send_tx

  #STEP sleep
  sleep 30

  #STEP kill rpc
  killbyname_gracefully "cache/node4/exchaind"
  sleep 2

  #STEP BUILD version2
  build_exchain $version2

  #STEP add rpc
  add_rpc 4

  #STEP sleep
  sleep 30

  #STEP CHECK BLOCK ALL
  check_block rpc4 tx

  #STEP kill 25% v
  killbyname_gracefully "cache/node3/exchaind"
  sleep 3

  #STEP add v
  add_val 3
  sleep 30

  #STEP CHECK val
  check_block val3 tx
}

function case_2() {
  # Upgrade 25% validator, then upgrade rest of the validators ,and then upgrade all the rpc nodes.
  version1=$1
  version2=$2
  echo "[][][][][][][][][][][][][][][][][][]"
  echo "[][][][][]    case_2      [][][][][]"
  echo "[][][][][][][][][][][][][][][][][][]"

  # pre
  case_prepare $version1

  #STEP sleep
  sleep 20

  #STEP check block ,all
  check_block_all

  #STEP send tx
  send_tx
  sleep 30

  #STEP BUILD version2
  build_exchain $version2

  #STEP upgrade 25% v
  killbyname_gracefully "cache/node3/exchaind"
  sleep 3
  add_val 3
  sleep 20

  #STEP check block
  check_block val3 tx

  #STEP upgrade 100% v
  killbyname_gracefully "cache/node2/exchaind"
  sleep 3
  add_val 2

  killbyname_gracefully "cache/node1/exchaind"
  sleep 3
  add_val 1

  killbyname_gracefully "cache/node0/exchaind"
  sleep 3
  add_seed
  
  sleep 30

  #STEP check block
  check_block val2 tx
  check_block val1 tx
  check_block val0 tx

  #STEP upgrade 100% rpc
  #STEP kill rpc
  killbyname_gracefully "cache/node4/exchaind"
  killbyname_gracefully "cache/node5/exchaind"
  sleep 3

  #STEP add rpc
  add_rpc 4
  add_rpc 5

  #STEP check block
  sleep 30
  check_block rpc4 tx
  check_block rpc5 tx
}

function case_3() {
  # Upgrade all the validators,then upgrade 1 RPC
  version1=$1
  version2=$2
  echo "[][][][][][][][][][][][][][][][][][]"
  echo "[][][][][]    case_3      [][][][][]"
  echo "[][][][][][][][][][][][][][][][][][]"

  # pre
  case_prepare $version1

  #STEP sleep
  sleep 20

  #STEP check block val
  check_block_all

  #STEP send tx
  send_tx
  sleep 30

  #STEP BUILD version2
  build_exchain $version2

  #STEP upgrade 100% v
  killbyname_gracefully "cache/node3/exchaind"
  killbyname_gracefully "cache/node2/exchaind"
  killbyname_gracefully "cache/node1/exchaind"
  killbyname_gracefully "cache/node0/exchaind"
  sleep 3

  add_seed
  add_val 1
  add_val 2
  add_val 3
  sleep 20

  #STEP check block val
  check_block val0 tx 
  check_block val1 tx
  check_block val2 tx
  check_block val3 tx

  #STEP upgrade 1 rpc
  killbyname_gracefully "cache/node5/exchaind"
  add_rpc 5
  sleep 10

  #STEP check block rpc
  check_block rpc5 tx

}

if [ -z ${IP} ]; then
  IP="127.0.0.1"
fi

### send two params , the first is the old version of exchain, the second is the newer version.
exc_version1=$1
exc_version2=$2
caseopt $exc_version1 $exc_version2
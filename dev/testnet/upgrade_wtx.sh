#!/usr/bin/env bash

NUM_NODE=4

set -e
set -o errexit
set -a
set -m

# set -x # activate debugging

source oec.profile
PRERUN=false

# WTX enable case with number 
ENABLE_WTX_CASE=()
ENABLE_WTX_CASE+=('{"val0":"","val1":"","val2":"","val3":"","rpc4":""}')
ENABLE_WTX_CASE+=('{"val0":"","val1":"","val2":"","val3":"","rpc4":""}')
ENABLE_WTX_CASE+=('{"val0":"","val1":"--enable-wtx=true","val2":"--enable-wtx=true","val3":"","rpc4":"--enable-wtx=true"}')
ENABLE_WTX_CASE+=('{"val0":"","val1":"--enable-wtx=true","val2":"--enable-wtx=true","val3":"--enable-wtx=true","rpc4":"--enable-wtx=true"}')
ENABLE_WTX_CASE+=('{"val0":"","val1":"--enable-wtx=true","val2":"--enable-wtx=true","val3":"--enable-wtx=true","rpc4":"--enable-wtx=true"}')

# WTX version 
WTX_VERSIONS=()
WTX_VERSIONS+=('{"val0":"v1.1.4.3","val1":"yxq/wtx","val2":"yxq/wtx","val3":"v1.1.4.3","rpc4":"yxq/wtx"}')
WTX_VERSIONS+=('{"val0":"v1.1.4.3","val1":"yxq/wtx","val2":"v1.1.4.3","val3":"v1.1.4.3","rpc4":"v1.1.4.3"}')
WTX_VERSIONS+=('{"val0":"v1.1.4.3","val1":"yxq/wtx","val2":"yxq/wtx","val3":"v1.1.4.3","rpc4":"yxq/wtx"}')
WTX_VERSIONS+=('{"val0":"v1.1.4.3","val1":"yxq/wtx","val2":"yxq/wtx","val3":"yxq/wtx","rpc4":"yxq/wtx"}')
WTX_VERSIONS+=('{"val0":"v1.1.4.3","val1":"yxq/wtx","val2":"yxq/wtx","val3":"yxq/wtx","rpc4":"yxq/wtx"}')

# Extract case 
NODEKEYS=""

ADD_WTHIE_LIST_CASE=()
ADD_WTHIE_LIST_CASE+=('{"node0":"","node1":"","node2":"","node3":"","node4":""}')
ADD_WTHIE_LIST_CASE+=('{"node0":"","node1":"","node2":"","node3":"","node4":""}')
ADD_WTHIE_LIST_CASE+=('{"node0":"","node1":"--mempool.node_key_whitelist","node2":"--mempool.node_key_whitelist","node3":"","node4":"--mempool.node_key_whitelist"}')
ADD_WTHIE_LIST_CASE+=('{"node0":"","node1":"--mempool.node_key_whitelist","node2":"--mempool.node_key_whitelist","node3":"--mempool.node_key_whitelist","node4":"--mempool.node_key_whitelist"}')
ADD_WTHIE_LIST_CASE+=('{"node0":"","node1":"--mempool.node_key_whitelist","node2":"--mempool.node_key_whitelist","node3":"","node4":"--mempool.node_key_whitelist"}')


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
  port_map='{"val0":26657,"val1":26757,"val2":26857,"val3":26957,"rpc4":27057,"rpc5":27157}'
  port=`echo $port_map | jq .$node`
  height=`exchaincli status --node http://${IP}:${port} | jq .sync_info.latest_block_height | awk '{ gsub(/"/,""); print $0 }'`
  echo $height
}

function get_latest_tx_count() {
  node=$1
  height=$2
  port_map='{"val0":8545,"val1":8645,"val2":8745,"val3":8845,"rpc4":8945,"rpc5":9045}'
  port=`echo $port_map | jq .$node`
  hex_height=`printf "0x%x" $height`
  data_json='{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByNumber","params":["'${hex_height}'"],"id":1}'
  tx_count=`curl -X POST --data $data_json -H "Content-Type: application/json" http://${IP}:$port -s | jq .result | awk '{ gsub(/"/,""); print $0 }' `
  echo $tx_count
}

function get_white_list() {
  index=$1
  case_number=$2

  case_document=`echo ${ADD_WTHIE_LIST_CASE[$index]} | jq .node${index}`

  if [[ -z $case_document ]]; then 
    echo ""
  else 
    echo "--mempool.node_key_whitelist=${NODEKEYS}"
  fi 

  
}

function get_node_enable_wtx_opts() {
  case_number=$1
  node_name=$2
  enable_wtx=`echo ${ENABLE_WTX_CASE[$case_number]} | jq .$node_name`
  echo enable_wtx
}

function check_block() {
  echo "Check block: $1"
  _sec_i=0
  node_name=$1 # like val1 rpc2 ....
  extra_check=$2
  is_valid=0
  height0=`get_latest_height $node_name`
  
  # The count of block that the tx in it is great than 0
  block_contains_tx=0
  while [ $_sec_i -le 600 ] 
  do
    latest_height=`get_latest_height $node_name`
    echo "latest_height,"$latest_height 
    latest_tx_count=`get_latest_tx_count $node_name $latest_height`
    echo "tx count,"$latest_tx_count
    if [[ $latest_tx_count != "0x0" ]] ; then
      let block_contains_tx+=1
    fi

    if [[ $extra_check = "tx" ]] ; then
      # Get new blocks over 100 and those blocks contains tx 
      if [ `expr $latest_height - $height0` -gt 25 -a $block_contains_tx -gt 25 ] ;then
        echo "block_contains_tx ,"$block_contains_tx
        is_valid=1
        break
      fi
    else
      if [ `expr $latest_height - $height0` -gt 10 ] ;then
        is_valid=1
        break
      fi
    fi

    let _sec_i+=1
    sleep 1
    echo "Checking... $latest_height"
  done

  if [ $is_valid -eq 0 ] ;then
      echo "Check valid $node_name: Failed, not pass."
      exit 99
      # return 1
  else
      echo "Check valid $node_name: Successful, pass."
      # return 0
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
  parameter=$1
  echo "start sending tx ..."
  (cd ../client/ && bash run.sh $parameter > /dev/null 2>&1 &)
}

function add_val() {
  index=$1
  case_number=$2
  LOG_LEVEL=main:info,*:error,consensus:error,state:info,provider:info

  wtx_case=${ENABLE_WTX_CASE[$index]}
  enable_wtx=`echo $wtx_case | jq .val${index}`

  white_list_flag=`get_white_list $index $case_number`

  if [ $index == 0 ]
  then
    echo "add a seed node,val"$index
    seed_mode=true
    ((restport = REST_PORT)) # for evm tx

    nohup exchaind start \
      --home cache/node${index}/exchaind \
      --p2p.seed_mode=$seed_mode \
      --p2p.allow_duplicate_ip \
      --p2p.pex=false \
      --p2p.addr_book_strict=false \
      --p2p.laddr tcp://${IP}:${seedp2pport} \
      --rpc.laddr tcp://${IP}:${seedrpcport} \
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
      ${enable_wtx} \
      ${white_list_flag} \
      --keyring-backend test >cache/val${index}.log 2>&1 &

  else
    echo "add a val node,val"$index
    seed_mode=false

    seed=$(exchaind tendermint show-node-id --home cache/node0/exchaind)@${IP}:${seedp2pport}
    ((p2pport = BASE_PORT_PREFIX + index * 100 + P2P_PORT_SUFFIX))
    ((rpcport = BASE_PORT_PREFIX + index * 100 + RPC_PORT_SUFFIX))  # for exchaincli
    ((restport = index * 100 + REST_PORT)) # for evm tx

    nohup exchaind start \
      --home cache/node${index}/exchaind \
      --p2p.seed_mode=$seed_mode \
      --p2p.allow_duplicate_ip \
      --p2p.pex=false \
      --p2p.addr_book_strict=false \
      --p2p.seeds $seed \
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
      ${enable_wtx} \
      ${white_list_flag} \
      --keyring-backend test >cache/val${index}.log 2>&1 &
  fi
}

function add_rpc() {
  index=$1
  case_number=$2
  echo "add a rpc node,"$index

  wtx_case=${ENABLE_WTX_CASE[$index]}
  enable_wtx=`echo $wtx_case | jq .val${index} | awk '{ gsub(/"/,""); print $0 }'`

  white_list_flag=`get_white_list $index $case_number`

  NAME=node${index}
  let p2p_port=${BASE_PORT_PREFIX}+${index}*100+${P2P_PORT_SUFFIX}
  let rpc_port=${BASE_PORT_PREFIX}+${index}*100+${RPC_PORT_SUFFIX}
  seed_addr=$(exchaind tendermint show-node-id --home cache/node0/exchaind)@${IP}:${seedp2pport}
  ((restport = index * 100 + REST_PORT)) # for evm tx

  LOG_LEVEL=main:info,*:error,state:info
  # LOG_LEVEL=main:info,*:error,state:debug,consensus:debug

  nohup ${BIN_NAME} start \
  --chain-id ${CHAIN_ID} \
  --home cache/${NAME}/exchaind \
  --p2p.laddr tcp://${IP}:${p2p_port} \
  --rest.laddr tcp://${IP}:$restport \
  --p2p.seeds ${seed_addr} \
  --log_level ${LOG_LEVEL} \
  --enable-gid \
  --append-pid \
  --p2p.addr_book_strict=false \
  --enable-preruntx=${PRERUN} \
  ${enable_wtx} \
  ${white_list_flag} \
  --rpc.laddr tcp://${IP}:${rpc_port} > cache/rpc${index}.log 2>&1 &
}

function case_prepare() {
  # Prepare 4 validators and 2 rpc nodes
  v1=$1

  killbyname_gracefully ${BIN_NAME}
  killbyname_gracefully "run.sh"
  killbyname_gracefully "./client"

  bash testnet.sh -i
  build_exchain $v1
  bash testnet.sh -s -n 4
  bash addnewnode.sh -n 4
  bash addnewnode.sh -n 5
}

function caseopt() {
  echo "caseopt()"
  v1=$1
  v2=$2

  case_1 $v1 $v2
  case_2 $v1 $v2
  case_3 $v1 $v2
  echo "All cases finished!"
}

function wtx_case(){
  echo "[][][][][][][][][][][][][][][][][][][]"
  echo "[][][][][]    WTX-Case      [][][][][]"
  echo "[][][][][][][][][][][][][][][][][][][]"

  case_number=$1

  # kill all instance 
  killbyname_gracefully ${BIN_NAME}
  killbyname_gracefully "run.sh"
  killbyname_gracefully "./client"

  # direct build on stable branch 
  build_exchain "dev"

  # genesis
  bash testnet.sh -i -n 5

  # extract every node keys in global variables  
  node1_key=`exchaincli keys extract-node-key cache/node1/exchaind/config/node_key.json | awk '{print $2}'`
  node4_key=`exchaincli keys extract-node-key cache/node4/exchaind/config/node_key.json | awk '{print $2}'`

  NODEKEYS="${node1_key},${node4_key}"
  echo "--------NODEKEYS: ${NODEKEYS} ---------------------------------"
  
  build_exchain "v1.1.4.3"
  bash testnet.sh -s -n 5

  

  sleep 10 

  # restar all nodes with specific version 

  killbyname_gracefully "cache/node1/exchaind"
  sleep 2
  version=`echo ${WTX_VERSIONS[$case_number]} | jq .val1 | awk '{ gsub(/"/,""); print $0 }'`
  build_exchain $version
  add_val 1 $case_number

  killbyname_gracefully "cache/node2/exchaind"
  sleep 2
  version=`echo ${WTX_VERSIONS[$case_number]} | jq .val2 | awk '{ gsub(/"/,""); print $0 }'`
  build_exchain $version
  add_val 2 $case_number

  killbyname_gracefully "cache/node3/exchaind"
  sleep 2
  version=`echo ${WTX_VERSIONS[$case_number]} | jq .val3 | awk '{ gsub(/"/,""); print $0 }'`
  build_exchain $version
  add_val 3 $case_number

  killbyname_gracefully "cache/node4/exchaind"
  sleep 2
  version=`echo ${WTX_VERSIONS[$case_number]} | jq .rpc4 | awk '{ gsub(/"/,""); print $0 }'`
  build_exchain $version
  add_val 4 $case_number

  # bootstrap succress 
}

function case_1() {
  # Upgrade 1 rpc node , then upgrade 1 validator node.
  echo "[][][][][][][][][][][][][][][][][][]"
  echo "[][][][][]    case_1      [][][][][]"
  echo "[][][][][][][][][][][][][][][][][][]"

  v1=$1
  v2=$2

  # pre
  case_prepare $v1
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

  #STEP BUILD v2
  build_exchain $v2

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
  v1=$1
  v2=$2
  echo "[][][][][][][][][][][][][][][][][][]"
  echo "[][][][][]    case_2      [][][][][]"
  echo "[][][][][][][][][][][][][][][][][][]"

  # pre
  case_prepare $v1

  #STEP sleep
  sleep 20

  #STEP check block ,all
  check_block_all

  #STEP send tx
  send_tx
  sleep 30

  #STEP BUILD v2
  build_exchain $v2

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
  add_val 0
  
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
  v1=$1
  v2=$2
  echo "[][][][][][][][][][][][][][][][][][]"
  echo "[][][][][]    case_3      [][][][][]"
  echo "[][][][][][][][][][][][][][][][][][]"

  # pre
  case_prepare $v1

  #STEP sleep
  sleep 20

  #STEP check block val
  check_block_all

  #STEP send tx
  send_tx
  sleep 30

  #STEP BUILD v2
  build_exchain $v2

  #STEP upgrade 100% v
  killbyname_gracefully "cache/node3/exchaind"
  killbyname_gracefully "cache/node2/exchaind"
  killbyname_gracefully "cache/node1/exchaind"
  killbyname_gracefully "cache/node0/exchaind"
  sleep 3

  add_val 0
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
# exc_v1=$1
# exc_v2=$2
# caseopt $exc_v1 $exc_v2
wtx_case $1
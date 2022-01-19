#!/usr/bin/env bash



CLIENT_COUNT=1
while getopts "c:" opt; do
  case $opt in
  c)
    echo "CLIENT_COUNT=$OPTARG"
    CLIENT_COUNT=$OPTARG
    ;;
  \?)
    echo "Invalid option: -$OPTARG"
    ;;
  esac
done

testNetPath=${PWD}

killbyname() {
  NAME=$1
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2", "$8}'
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2}' | sh
  echo "All <$NAME> killed!"
}

clean(){
  killbyname counter
  killbyname redis-server
  rm dump.rdb
  rm redis.log
  rm -rf counterClient*.log
  rm -rf cache
  sleep 1
}

network(){
  nohup redis-server > redis.log &
  #
  ./testnet.sh -s -i -n 4 -u -x -f
  #
  sleep 1
  echo "add one consumer node "
  ./addnewnode.sh -n 4 -d -x -f
  echo "add one more consumer node"
  ./addnewnode.sh -n 5 -d -x -f

  echo "wait 5 seconds ,and we will try to catch up as fast as we can"
  sleep 5
  ./addnewnode.sh -n 6 -u -x -f
}



prepareClient() {
  cd ../client
  if [[ -f ${PWD}/client ]];then
    echo "client bin exists"
    cd ${testNetPath}
    return
  fi
  export GO111MODULE=on
  go build
  cd ${testNetPath}
}

runClinet(){
  cd ../client
  for i in $(seq 1 ${1})
  do
    echo "start client ${i}"
    nohup ./run.sh > ${testNetPath}/counterClient${i}.log 2>&1 &
  done
  cd ${testNetPath}
}

fail(){
   v=${1}
   failure=`grep -i "${v}" ./cache/*.log | grep -C 10 "${v}"`
   if [[ -n "${failure}" ]]; then
      echo "fail:${failure}"
      killbyname counter
      killbyname redis-server
      exit
   fi
}
periodGrep(){
   cd ${testNetPath}
   for (( i=0; ; ))
   do

     cd .

     grep -r 'delta is waitting prerun to be canceld or finished' ./
     sleep 2
     grep -r 'discard' ./
     sleep 2
     grep -r 'beginBlock execute twice' ./
     sleep 2
     fail 'CONSENSUS FAILURE'
     fail 'wrong Block.Header.AppHash'
     sleep 5
   done
}


clean

prepareClient

network

runClinet ${CLIENT_COUNT}

periodGrep


#cd ../ && ./addnewnode.sh -n 7 -u -x -f && cd ./cache/ && tail -f rpc7.log
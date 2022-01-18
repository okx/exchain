#!/usr/bin/env bash


killbyname() {
  NAME=$1
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2", "$8}'
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2}' | sh
  echo "All <$NAME> killed!"
}

killbyname redis-server
rm dump.rdb
rm redis.log

nohup redis-server > redis.log &
#
./testnet.sh -s -i -n 4 -u -x -f
#
sleep 1
echo "add one producer node "
./addnewnode.sh -n 4 -d -x -f
echo "add one more producer node"
./addnewnode.sh -n 5 -d -x -f

echo "wait 5 seconds ,and we will try to catch up as fast as we can"
sleep 5
./addnewnode.sh -n 6 -u -x -f


#killbyname redis-server




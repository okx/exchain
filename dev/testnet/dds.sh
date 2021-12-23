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
./testnet.sh -s -i -n 4
#
sleep 5

./addnewnode.sh -n 4

sleep 5

#killbyname redis-server




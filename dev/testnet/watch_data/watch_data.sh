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

./producer.sh -s -i -n 4

sleep 1

./consumer.sh -n 4

sleep 5

h=`redis-cli -h 127.0.0.1 -p 6379 get dds:2:LatestHeight`
echo "get latestHeight: $h"

p_res=`curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x2CF4ea7dF75b513509d95946B43062E26bD88035","0x0"],"id":1}' -H "Content-Type: application/json" 127.0.0.1:8545`
c_res=`curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x2CF4ea7dF75b513509d95946B43062E26bD88035","0x0"],"id":1}' -H "Content-Type: application/json" 127.0.0.1:8544`
if [ $p_res != $c_res ]; then
  echo "balance of consumer is not equal with producer"
  exit
else
  echo "init balance is equal:"
  echo $c_res
fi

echo "sendTransaction"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x2CF4ea7dF75b513509d95946B43062E26bD88035", "to":"0x0073F2E28ef8F117e53d858094086Defaf1837D5", "value":"0xde0b6b3a76400000"}],"id":1}' -H "Content-Type: application/json" 127.0.0.1:8544
sleep 10

p_res=`curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x2CF4ea7dF75b513509d95946B43062E26bD88035","0x0"],"id":1}' -H "Content-Type: application/json" 127.0.0.1:8545`
c_res=`curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x2CF4ea7dF75b513509d95946B43062E26bD88035","0x0"],"id":1}' -H "Content-Type: application/json" 127.0.0.1:8544`
if [ $p_res != $c_res ]; then
  echo "balance of consumer is not equal with producer"
  exit
else
  echo "get balance of consumer and producer is equal:"
  echo $c_res
fi

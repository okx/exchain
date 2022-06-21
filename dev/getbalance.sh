#!/bin/bash 
#set -x
balance=0x3DF95c73357f988F732c4c7a8Fa2f9beD7952862
while true
do
	block_num_json=`curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' -H "Content-Type: application/json" http://127.0.0.1:8545`
	block_num=`jq '.result'  <<< $block_num_json`
	c_cmd=`printf '{"jsonrpc":"2.0","method":"eth_getBalance","params":["%s",%s],"id":1}' $balance $block_num`
	c_res=`curl -X POST --data $c_cmd -H "Content-Type: application/json" http://127.0.0.1:8545`
	echo $c_res
done

#!/usr/bin/env bash

rm -rf cache
rm -rf nodecache

./testnet.sh -i -n 4

mv cache nodecache
./testnet.sh -r -s -n 4 -w


sleep 3

exchaincli status -n tcp://localhost:26657 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26757 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26857 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26957 |grep -v validator_info |grep id



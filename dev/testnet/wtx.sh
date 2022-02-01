#!/usr/bin/env bash

./testnet.sh -r -s -n 4 -w

sleep 3

exchaincli status -n tcp://localhost:26657 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26757 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26857 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26957 |grep -v validator_info |grep id


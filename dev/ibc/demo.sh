#!/bin/bash

./two-chainz
sleep 1
rly q bal ibc-1 testkey
sleep 1
rly q bal exchain-101 admin16
sleep 2
rly tx link oec101_ibc1  -d -o 3s

sleep 2
rly tx transfer exchain-101 ibc-1 10000okt $(rly chains address ibc-1) --path oec101_ibc1
sleep 1
rly tx relay-pkts oec101_ibc1 -d
sleep 1
echo "================================================"
rly q bal ibc-1 testkey
sleep 1
rly q bal exchain-101 admin16

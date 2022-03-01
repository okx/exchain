#!/bin/bash

./two-chainz
rly q bal ibc-1 testkey
rly q bal exchain-101 admin16
rly tx link oec101_ibc1  -d -o 3s

rly tx transfer exchain-101 ibc-1 10000okt $(rly chains address ibc-1) --path oec101_ibc1
rly tx relay-pkts oec101_ibc1 -d 
echo "================================================"
rly q bal ibc-1 testkey
rly q bal exchain-101 admin16

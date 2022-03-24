#!/bin/bash

echo "step1 ,link path firstly"
rly tx link demo -d -o 3s
rly q bal exchain-101 admin16
rly q bal ibc-1 testkey
echo "step2, transfer token"
rly tx transfer exchain-101 ibc-1 10000okt $(rly chains address ibc-1)
echo "step3, notify packet"
rly tx relay-pkts demo -d
echo "step4,ack packets "
rly tx relay-acks demo -d
rly q bal exchain-100
rly q bal exchain-101

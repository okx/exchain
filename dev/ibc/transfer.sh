#!/bin/bash

echo "step1 ,link path firstly"
rly tx link demo -d -o 3s
rly bal exchain-100
rly bal exchain-101
echo "step2, transfer token"
rly tx transfer exchain-100 exchain-101 1okt $(rly chains address exchain-101)
echo "step3, notify packet"
rly tx relay-pkts demo -d
echo "step4,ack packets "
rly tx relay-acks demo -d

rly bal exchain-100
rly bal exchain-101

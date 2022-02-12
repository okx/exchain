#!/bin/bash

echo "step1 ,link path firstly"
rly tx link demo -d -o 3s
rly q bal exchain-100 admin16
rly q bal exchain-101 admin16
echo "step2, transfer token"
rly tx transfer exchain-100 exchain-101 10000okt ex1s0vrf96rrsknl64jj65lhf89ltwj7lksr7m3r9
echo "step3, notify packet"
rly tx relay-pkts demo -d
echo "step4,ack packets "
rly tx relay-acks demo -d
rly q bal exchain-100
rly q bal exchain-101

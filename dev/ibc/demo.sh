#!/bin/bash

function failOnExit() {
    $@
   if [ $? -ne 0 ]; then
        echo "["$@"] failed"
        exit
    fi
}
function  qualBal() {
    echo "========================="
    echo "ibc-1 testKey account balance: $(rly q bal ibc-1 testkey)"
    echo "exchain-101 admin16 account balance:$(rly q bal exchain-101 admin16)"
    echo "========================="
}
./two-chainz
sleep 5
qualBal

failOnExit rly tx link oec101_ibc1  -d -o 3s --override
failOnExit rly tx link oec101_ibc0  -d -o 3s --override
failOnExit rly tx link oec101_oec100 -d -o 3s --override

failOnExit rly tx link oec100_ibc1  -d -o 3s --override
failOnExit rly tx link oec100_ibc0  -d -o 3s --override
failOnExit rly tx link oec100_oec101  -d -o 3s --override

failOnExit rly tx link ibc1_oec101  -d -o 3s --override
failOnExit rly tx link ibc1_oec100  -d -o 3s --override
failOnExit rly tx link ibc1_ibc0  -d -o 3s --override

failOnExit rly tx link ibc0_oec101  -d -o 3s --override
failOnExit rly tx link ibc0_oec100  -d -o 3s --override
failOnExit rly tx link ibc0_ibc1  -d -o 3s --override

#
rly chains list
rly paths list

sleep 2
rly tx transfer exchain-101 ibc-1 10000okt $(rly chains address ibc-1) --path oec101_ibc1
sleep 1
rly tx relay-pkts oec101_ibc1 -d
sleep 1
qualBal

rly tx transfer ibc-1 exchain-101 1000000samoleans raw:$(rly chains address exchain-101) --path ibc1_oec101
sleep 1
rly tx relay-pkts ibc1_oec101 -d
sleep 1
qualBal


rly tx transfer exchain-101 ibc-1 1000000ibc/49D321B40FCF56B0370E5673CF090389C8E9CD185209FBE1BEE5D94E58E69BDC $(rly chains addr  ibc-1) --path ibc1_oec101
sleep 1
rly tx relay-pkts ibc1_oec101 -d
sleep 1
qualBal
sleep 1
qualBal

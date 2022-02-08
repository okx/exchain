#!/usr/bin/env bash

./testnet.sh -s -i -n 4 -w

sleep 3

exchaincli status -n tcp://localhost:26657 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26757 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26857 |grep -v validator_info |grep id
exchaincli status -n tcp://localhost:26957 |grep -v validator_info |grep id

exit
    "id": "0b066ca0790f27a6595560b23bf1a1193f100797",
    "id": "3813c7011932b18f27f172f0de2347871d27e852",
    "id": "6ea83a21a43c30a280a3139f6f23d737104b6975",
    "id": "bab6c32fa95f3a54ecb7d32869e32e85a25d2e08",

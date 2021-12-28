#!/usr/bin/env bash

./testnet.sh -s -i -n 4 -c hang.json

sleep 5

./addnewnode.sh -n 4
./addnewnode.sh -n 5

#!/usr/bin/env bash

./testnet.sh -s -i -n 4 -c cases/allcases.json -x

sleep 5

./addnewnode.sh -n 4
./addnewnode.sh -n 5

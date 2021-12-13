#!/usr/bin/env bash

./testnet.sh -s -i -n 4

sleep 5

./addnewnode.sh -n 4

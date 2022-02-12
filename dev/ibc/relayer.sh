#!/bin/bash
rm -rf ~/.relayer
rly config init
rly config add-chains configs/demo/chains
SEED0=$(jq -r '.mnemonic' $GAIA_DATA/exchain-100/key_seed.json)
SEED1=$(jq -r '.mnemonic' $GAIA_DATA/exchain-101/key_seed.json)
echo "Key $(rly keys restore exchain-100 admin16 "$SEED0") imported from exchain-100 to relayer..."
echo "Key $(rly keys restore exchain-101 admin16 "$SEED1") imported from exchain-101 to relayer..."
rly config add-paths configs/demo/paths
rly chains list
rly paths list
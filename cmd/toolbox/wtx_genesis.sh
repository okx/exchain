#!/usr/bin/env bash

# TODO: only used in this situation to handle the WTX test for speed up time 

WORKSPACE="workspace"
NUMBER=4
PREFIX="n"

for ((index = 1; index < ${NUMBER}; index++)); do 

  if [ "$(uname -s)" == "Darwin" ]; then
      sed -i "" 's/"enable_call": false/"enable_call": true/' ${WORKSPACE}/${PREFIX}${index}/exchaind/config/genesis.json
      sed -i "" 's/"enable_create": false/"enable_create": true/' ${WORKSPACE}/${PREFIX}${index}/exchaind/config/genesis.json
      sed -i "" 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' ${WORKSPACE}/${PREFIX}${index}/exchaind/config/genesis.json
  else
      sed -i 's/"enable_call": false/"enable_call": true/' ${WORKSPACE}/${PREFIX}${index}/exchaind/config/genesis.json
      sed -i 's/"enable_create": false/"enable_create": true/' ${WORKSPACE}/${PREFIX}${index}/exchaind/config/genesis.json
      sed -i 's/"enable_contract_blocked_list": false/"enable_contract_blocked_list": true/' ${WORKSPACE}/${PREFIX}${index}/exchaind/config/genesis.json
  fi

 # has already install this 
  exchaind add-genesis-account 0xbbE4733d85bc2b90682147779DA49caB38C0aA1F 900000000okt --home ${WORKSPACE}/${PREFIX}${index}//exchaind
  exchaind add-genesis-account 0x4C12e733e58819A1d3520f1E7aDCc614Ca20De64 900000000okt --home ${WORKSPACE}/${PREFIX}${index}//exchaind
  exchaind add-genesis-account 0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0 900000000okt --home ${WORKSPACE}/${PREFIX}${index}//exchaind
  exchaind add-genesis-account 0x2Bd4AF0C1D0c2930fEE852D07bB9dE87D8C07044 900000000okt --home ${WORKSPACE}/${PREFIX}${index}//exchaind

done
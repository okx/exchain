#!/usr/bin/env bash

function startpeer {
    index=$1
    export OKCHAIN_PEER_ID=vp${index}
    export LOG_STDOUT_FILE=${OKCHAIN_PEER_ID}.json
    echo Run node ${OKCHAIN_PEER_ID} ...
    ./lock >> ${LOG_STDOUT_FILE} 2>${LOG_STDOUT_FILE} &
}

function main {
    go build
    /killbyname.sh lock
    sleep 3
    for ((index=0; index<$1; index++))
    do
        startpeer ${index}
    done
}

etcdctl --endpoints=localhost:2379 del distributed_lock_key
main $1

sleep 3

etcdctl --endpoints=localhost:2379 lease list
etcdctl --endpoints=localhost:2379 get distributed_lock_key

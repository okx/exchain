#!/bin/bash

HOME=$1
VERSION=$2
CMD=exchaind
NUM_EXECUTIONS=2

set -e
set -o errexit
set -a
set -m

killbyname() {
    NAME=$1
    ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2", "$8}'
    ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2}' | sh
    echo "All <$NAME> killed!"
}


killbyname exchaind
killbyname exchaincli

set -x # activate debugging


if [ "$1" == "-h" ];
then
    echo "Usage: exchaind [home] [s0|s1]"
    exit 0
fi

if [ -z "$HOME" ];
then
    echo specify home directory first please
    exit -1
fi

if [ -z "$VERSION" ];
then
    echo specify version first please
    exit -1
fi


echo using $VERSION mode

if [ "$VERSION" = "s0" ];
then
    for (( i=0; i<NUM_EXECUTIONS; i++ ))
    do
        echo runing $i-th times...
        $CMD data prune-compact all --home $HOME
    done
    rm -rf $HOME/data/cs.wal
    rm -rf $HOME/data/tx_index.db
    rm -rf $HOME/data/evidence.db
    rm -rf $HOME/data/watch.db
else
    for (( i=0; i<NUM_EXECUTIONS; i++ ))
    do
        echo runing $i-th times...
        $CMD data prune-compact state --home $HOME
    done
    rm -rf $HOME/data/watch.db
fi

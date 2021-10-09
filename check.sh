#!/bin/sh

CMD=pgrep
PPROC=exchaind_tmp
while :
do
    if [ -n "`$CMD $PPROC`" ]
    then  echo "test is ok"
    else
	killall exchaind_123
        echo "test is killed"
    fi
    sleep 10
done

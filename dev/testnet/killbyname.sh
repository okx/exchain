#!/bin/bash

NAME=$1
MYNAME="killbyname.sh"

if [ $# -eq 0 ]; then
    echo "$MYNAME <process name>"
    exit
fi

ps -ef|grep "$NAME"|grep -v grep |grep -v $MYNAME |awk '{print "kill -9 "$2", "$8}'
ps -ef|grep "$NAME"|grep -v grep |grep -v $MYNAME |awk '{print "kill -9 "$2}' | sh
echo "All <$NAME> killed!"



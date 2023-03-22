#!/bin/bash

if [ -z "$1" ];
then
  echo "Specify a directory to store logs from kafka."
  echo "For example: ./sub.sh your_logs_dir"
  exit
fi

rm -rf $1/*
okbchaind subscribe logs localhost:9092 $1
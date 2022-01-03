#!/bin/bash

rm -rf $1/*
exchaind subscribe logs localhost:9092 $1
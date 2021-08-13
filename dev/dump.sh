#!/bin/bash


grep ApplyBlock $1 | sed 's/>//g'| sed 's/</,/g'|sed 's/module=main//g'| sed 's/ms//g' |awk '{print $3 $2 $4 $5 $7 $10}'> $1.csv

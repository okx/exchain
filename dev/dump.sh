#!/bin/bash


# Height<5810707>, Tx<84>, GasUsed<13688568>, RunTx[Elapsed<208ms>, abci<177ms>, persist<29ms>, saveState<2ms>], Evm[read<74ms>, write<5ms>, execute<67ms>], DB[read<8329>, write<0>], Round[], CommitRound[], Produce[] module=main
echo "Height,Tx,GasUsed,RunTx-Elapsed,abci,persist,saveState,evm-r,evm-w,evm-e,kv-r,kv-w," > $1.csv
grep GasUsed $1 \
| sed 's/>//g' \
| sed 's/\[//g' \
| sed 's/\]//g' \
| sed 's/Height<//g' \
| sed 's/Tx</,/g' \
| sed 's/GasUsed</,/g' \
| sed 's/RunTxElapsed</,/g' \
| sed 's/abci</,/g' \
| sed 's/persist</,/g' \
| sed 's/saveState</,/g' \
| sed 's/Evmread</,/g' \
| sed 's/write</,/g' \
| sed 's/execute</,/g' \
| sed 's/DBread</,/g' \
| sed 's/CommitRound/,/g' \
| sed 's/Round/,/g' \
| sed 's/ProduceElapsed</,/g' \
| sed 's/, ,/,/g' \
| sed 's/ms//g' |awk '{print $2}' >> $1.csv


# grep ApplyBlock $1 | sed 's/>//g'| sed 's/</,/g'|sed 's/module=main//g'| sed 's/ms//g' |awk '{print $3 $2 $4 $5 $7 $10}'> $1.csv

#!/bin/bash


# I[2021-10-07|18:51:32.543][85554] Height<13>, Tx<0>, GasUsed<0>, RunTx[Elapsed<6026ms>, abci<6012ms>, persist<10ms>, saveState<3ms>], Round[0], CommitRound[0], Produce[Elapsed<6542ms>, NewRound-0<0ms>, Propose-0<487ms>, Precommit-0<6ms>, RunTx-0-0<6043ms>, Waiting<4ms>] module=main
echo ",Height,Tx,GasUsed,RunTx-Elapsed,abci,persist,saveState,Round,CommitRound,Produce," > $1.csv
grep GasUsed $1 \
| sed 's/>//g' \
| sed 's/\[//g' \
| sed 's/\]//g' \
| sed 's/Height</,/g' \
| sed 's/Tx</,/g' \
| sed 's/GasUsed</,/g' \
| sed 's/RunTxElapsed</,/g' \
| sed 's/abci</,/g' \
| sed 's/persist</,/g' \
| sed 's/saveState</,/g' \
| sed 's/CommitRound/,/g' \
| sed 's/Round/,/g' \
| sed 's/ProduceElapsed</,/g' \
| sed 's/, ,/,/g' \
| sed 's/module=main//g'| sed 's/ms//g' |awk '{print $2}' >> $1.csv


# grep ApplyBlock $1 | sed 's/>//g'| sed 's/</,/g'|sed 's/module=main//g'| sed 's/ms//g' |awk '{print $3 $2 $4 $5 $7 $10}'> $1.csv

#!/bin/bash

#wating a new block
echo "Waiting to deposit operate..."
for((i=1;i<=10;i++));
do
    curHeight=`exchaincli status | grep "latest_block_height" | sed -r 's/.*"(.*)".*/\1/'`
    if [ $curHeight -gt 0 ]
    then
      break
    fi
#    echo "Waiting for a new block, cur height:" $curHeight
    sleep 2
done

# deposit for different validators
beginHeight=`exchaincli status | grep "latest_block_height" | sed -r 's/.*"(.*)".*/\1/'`
exchaincli tx staking deposit 10okt --from admin16 --gas auto --gas-prices 0.0000000001okt --gas-adjustment 1.3 -y >/dev/null
exchaincli tx staking deposit 20okt --from admin18 --gas auto --gas-prices 0.0000000001okt --gas-adjustment 1.3 -y >/dev/null
exchaincli tx staking deposit 30okt --from admin17 --gas auto --gas-prices 0.0000000001okt --gas-adjustment 1.3 -y >/dev/null
echo "Waiting to add shares operate..."
for((i=1;i<=10;i++));
do
    sleep 2
    curHeight=`exchaincli status | grep "latest_block_height" | sed -r 's/.*"(.*)".*/\1/'`
    if [ $curHeight -gt $beginHeight ]
    then
#      echo $curHeight
#      echo "Mint a new block now."
      break
#    else
#      echo "Waiting a new block, cur height:" $curHeight
    fi
done

#add staking shares
exchaincli tx staking add-shares exvaloper1pt7xrmxul7sx54ml44lvv403r06clrdkehd8z7 --from admin17 --gas auto --gas-prices 0.0000000001okt --gas-adjustment 1.3 -y >/dev/null
exchaincli tx staking add-shares exvaloper1ve4mwgq9967gk338yptsg2fheur4ke322gzynt --from admin18 --gas auto --gas-prices 0.0000000001okt --gas-adjustment 1.3 -y >/dev/null
exchaincli tx staking add-shares exvaloper1q6ls3h64gkxq0r73u2eqwwr7d5mp583fm325zu --from admin16 --gas auto --gas-prices 0.0000000001okt --gas-adjustment 1.3 -y >/dev/null

echo "Init staking done!"


killall exchaind
rm -rf multi_run.log
make mainnet
export EXCHAIND_PATH=~/.exchaind

rm -rf ${EXCHAIND_PATH}
exchaind init multi_run --chain-id exchain-66 --home ${EXCHAIND_PATH}
cp /Users/oker/scf/genesis.json ${EXCHAIND_PATH}/config/genesis.json
rm -rf ${EXCHAIND_PATH}/data
cp -rf /Users/oker/scf/src/data ${EXCHAIND_PATH}
exchaind replay -d /Users/oker/scf/src/data-s1 --home ~/.exchaind > multi_run.log 2>&1 &

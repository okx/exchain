package state

import (
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	FlagParalleledTx = "paralleled-tx"
)

func execBlockOnProxyAppAsync(
	logger log.Logger,
	proxyAppConn proxy.AppConnConsensus,
	block *types.Block,
	stateDB dbm.DB,
) (*ABCIResponses, error) {
	var validTxs, invalidTxs = 0, 0

	txIndex := 0

	txReps := make([]abci.ExecuteRes, len(block.Txs))
	abciResponses := NewABCIResponses(block)

	commitInfo, byzVals := getBeginBlockValidatorInfo(block, stateDB)

	// Begin block
	var err error
	abciResponses.BeginBlock, err = proxyAppConn.BeginBlockSync(abci.RequestBeginBlock{
		Hash:                block.Hash(),
		Header:              types.TM2PB.Header(&block.Header),
		LastCommitInfo:      commitInfo,
		ByzantineValidators: byzVals,
	})
	if err != nil {
		logger.Error("Error in proxyAppConn.BeginBlock", "err", err)
		return nil, err
	}
	asCache := NewAsyncCache()
	signal := make(chan int, 1)
	rerunIdx := 0
	AsyncCb := func(execRes abci.ExecuteRes) {
		txReps[execRes.GetCounter()] = execRes
		for txReps[txIndex] != nil {
			res := txReps[txIndex]
			if res.Conflict(asCache) {
				rerunIdx++
				res = proxyAppConn.DeliverTxWithCache(abci.RequestDeliverTx{Tx: block.Txs[res.GetCounter()]})
				if proxyAppConn.Error() != nil {
					signal <- 0
					panic(proxyAppConn.Error())
				}
			}
			txRs := res.GetResponse()
			abciResponses.DeliverTxs[txIndex] = &txRs
			res.Collect(asCache)
			res.Commit()
			if abciResponses.DeliverTxs[txIndex].Code == abci.CodeTypeOK {
				validTxs++
			} else {
				invalidTxs++
			}

			txIndex++
			if txIndex == len(block.Txs) {
				logger.Info(fmt.Sprintf("BlockHeight %d With Tx %d : Paralle run %d, Conflected tx %d",
					block.Height, len(block.Txs), len(abciResponses.DeliverTxs)-rerunIdx, rerunIdx))
				signal <- 0
				return
			}
		}
	}

	// avoid panic when handle callback
	proxyCb := func(req *abci.Request, res *abci.Response) {
		return
	}
	proxyAppConn.SetResponseCallback(proxyCb)
	proxyAppConn.PrepareParallelTxs(AsyncCb, transTxsToBytes(block.Txs))

	// Run txs of block.
	for _, tx := range block.Txs {
		proxyAppConn.DeliverTxAsync(abci.RequestDeliverTx{Tx: tx})
		if err := proxyAppConn.Error(); err != nil {
			return nil, err
		}
	}

	if len(block.Txs) > 0 {
		//waiting for call back
		<-signal
		if err := proxyAppConn.Error(); err != nil {
			return nil, err
		}
		receiptsLogs := proxyAppConn.EndParallelTxs()
		for index, v := range receiptsLogs {
			if len(v) != 0 { // only update evm tx result
				abciResponses.DeliverTxs[index].Data = v
			}
		}
	}

	// End block.
	abciResponses.EndBlock, err = proxyAppConn.EndBlockSync(abci.RequestEndBlock{Height: block.Height})
	if err != nil {
		logger.Error("Error in proxyAppConn.EndBlock", "err", err)
		return nil, err
	}

	logger.Info("Executed block", "height", block.Height, "validTxs", validTxs, "invalidTxs", invalidTxs)

	return abciResponses, nil
}

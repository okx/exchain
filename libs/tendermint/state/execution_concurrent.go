package state

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/libs/tendermint/types"
	"time"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	dbm "github.com/okex/exchain/libs/tm-db"
)

// Executes block's transactions on proxyAppConn.
// Returns a list of transaction results and updates to the validator set
func execBlockOnProxyAppPartConcurrent(logger log.Logger,
	proxyAppConn proxy.AppConnConsensus,
	block *types.Block,
	stateDB dbm.DB,
	) (*ABCIResponses, error) {
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

	// Run txs of block.
	start := time.Now()
	abciResponses.DeliverTxs = proxyAppConn.DeliverTxsConcurrent(transTxsToBytes(block.Txs))
	elapsed := time.Since(start).Microseconds()
	deliverTxDuration += elapsed
	logger.Info("DeliverTxs duration", "cur", elapsed, "total", deliverTxDuration)

	var validTxs, invalidTxs = 0, 0
	for _, v := range abciResponses.DeliverTxs {
		if v.Code == abci.CodeTypeOK {
			validTxs++
		} else {
			invalidTxs++
		}
	}

	abciResponses.EndBlock, err = proxyAppConn.EndBlockSync(abci.RequestEndBlock{Height: block.Height})
	if err != nil {
		logger.Error("Error in proxyAppConn.EndBlock", "err", err)
		return nil, err
	}

	trace.GetElapsedInfo().AddInfo(trace.InvalidTxs, fmt.Sprintf("%d", invalidTxs))

	return abciResponses, nil
}
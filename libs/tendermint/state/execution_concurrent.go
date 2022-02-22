package state

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/libs/tendermint/types"

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
	var validTxs, invalidTxs = 0, 0
	txs := transTxsToBytes(block.Txs)
	logger.Error("transTxsToBytes", "origin", len(block.Txs), "after", len(txs))
	abciResponses.DeliverTxs = proxyAppConn.DeliverTxsConcurrent(txs)
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
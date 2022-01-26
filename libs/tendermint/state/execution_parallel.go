package state

import (
	"fmt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
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

	abciResponses.DeliverTxs = proxyAppConn.ParallelTxs(transTxsToBytes(block.Txs))
	for _, v := range abciResponses.DeliverTxs {
		if v.Code == abci.CodeTypeOK {
			validTxs++
		} else {
			invalidTxs++
		}
	}

	// End block.
	abciResponses.EndBlock, err = proxyAppConn.EndBlockSync(abci.RequestEndBlock{Height: block.Height})
	if err != nil {
		logger.Error("Error in proxyAppConn.EndBlock", "err", err)
		return nil, err
	}

	trace.GetElapsedInfo().AddInfo(trace.InvalidTxs, fmt.Sprintf("%d", invalidTxs))

	return abciResponses, nil
}

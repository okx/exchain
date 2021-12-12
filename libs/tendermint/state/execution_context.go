package state

import (
	"bytes"
	"github.com/okex/exchain/libs/iavl/trace"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)


type executionResult struct {
	res *ABCIResponses
	err error
}

type executionContext struct {
	height int64
	block *types.Block
	stopped bool
	result *executionResult

	prerunResultChan chan *executionContext
	proxyApp proxy.AppConnConsensus
	db dbm.DB
	logger log.Logger
	index int64
}


func (e *executionContext) dump(when string) {

	e.logger.Info(when,
		"gid", trace.GoRId,
		"stopped", e.stopped,
		"Height", e.block.Height,
		"index", e.index,
		"AppHash", e.block.AppHash,
	)
}

func (blockExec *BlockExecutor) prerunRoutine() {
	for context := range blockExec.prerunChan {
		prerun(context)
	}
}

func (blockExec *BlockExecutor) getPrerunResult(ctx *executionContext) (*ABCIResponses, error)  {

	for context := range blockExec.prerunResultChan {

		context.dump("Got prerun result")

		if context.stopped {
			continue
		}

		if context.height != ctx.block.Height{
			continue
		}

		if context.index != ctx.index {
			continue
		}

		if  bytes.Equal(context.block.AppHash, ctx.block.AppHash) {
			return context.result.res, context.result.err
		} else {
			// todo
			panic("wrong app hash")
		}
	}
	return nil, nil
}

func (blockExec *BlockExecutor) NotifyPrerun(height int64, block *types.Block) {

	context := blockExec.prerunContext
	// stop the existing prerun if any
	if blockExec.prerunContext != nil {
		if block.Height != blockExec.prerunContext.block.Height {
			context.dump("Prerun sanity check failed")

			// todo
			panic("Prerun sanity check failed")
		}
		context.dump("Stopping prerun")
		blockExec.prerunContext.stopped = true
	}

	blockExec.prerunIndex++
	blockExec.prerunContext = &executionContext{
		height:           height,
		block:            block,
		stopped:          false,
		db:               blockExec.db,
		proxyApp:         blockExec.proxyApp,
		logger:           blockExec.logger,
		prerunResultChan: blockExec.prerunResultChan,
		index:            blockExec.prerunIndex,
	}

	context = blockExec.prerunContext
	context.dump("Notify prerun")

	// start a new one
	blockExec.prerunChan <- blockExec.prerunContext
}

func prerun(context *executionContext) {
	context.dump("Start prerun")

	abciResponses, err := execBlockOnProxyApp(context)

	if !context.stopped {
		context.result = &executionResult{
			abciResponses, err,
		}
	}

	context.dump("Prerun completed")
	context.prerunResultChan <- context
}


func (blockExec *BlockExecutor) InitPrerun(open bool) {
	blockExec.proactivelyRunTx = open
	if blockExec.proactivelyRunTx {
		go blockExec.prerunRoutine()
	}
}
package state

import (
	"bytes"
	"fmt"
	gorid "github.com/okex/exchain/libs/goroutine"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/trace"

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
	height  int64
	block   *types.Block
	stopped bool
	result  *executionResult

	prerunResultChan chan *executionContext
	proxyApp         proxy.AppConnConsensus
	db               dbm.DB
	logger           log.Logger
	index            int64
}

func (e *executionContext) dump(when string) {

	e.logger.Info(when,
		"gid", gorid.GoRId,
		"stopped", e.stopped,
		"Height", e.block.Height,
		"index", e.index,
		"AppHash", e.block.AppHash,
	)
}

func (e *executionContext) stop() {
	e.stopped = true
	e.proxyApp.SetOptionSync(abci.RequestSetOption{
		Key: "ResetDeliverState",
	})
}

func (blockExec *BlockExecutor) flushPrerunResult() {
	for {
		select {
		case context := <-blockExec.prerunResultChan:
			context.dump("Flush prerun result")
		default:
			goto BREAK
		}
	}
BREAK:
}

func (blockExec *BlockExecutor) prerunRoutine() {
	for context := range blockExec.prerunChan {
		prerun(context)
	}
}

func (blockExec *BlockExecutor) getPrerunResult(ctx *executionContext) (*ABCIResponses, error) {

	for context := range blockExec.prerunResultChan {

		context.dump("Got prerun result")

		if context.stopped {
			continue
		}

		if context.height != ctx.block.Height {
			continue
		}

		if context.index != ctx.index {
			continue
		}

		if bytes.Equal(context.block.AppHash, ctx.block.AppHash) {
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
	if context != nil {
		if block.Height != context.block.Height {
			context.dump("Prerun sanity check failed")

			// todo
			panic("Prerun sanity check failed")
		}
		context.dump("Stopping prerun")
		if height != 1 {
			context.stop()
		}
	}
	blockExec.flushPrerunResult()
	blockExec.prerunIndex++
	context = &executionContext{
		height:           height,
		block:            block,
		stopped:          false,
		db:               blockExec.db,
		proxyApp:         blockExec.proxyApp,
		logger:           blockExec.logger,
		prerunResultChan: blockExec.prerunResultChan,
		index:            blockExec.prerunIndex,
	}

	context.dump("Notify prerun")
	blockExec.prerunContext = context

	// start a new one
	blockExec.prerunChan <- blockExec.prerunContext
}

func prerun(context *executionContext) {
	context.dump("Start prerun")

	trc := trace.NewTracer(fmt.Sprintf("prerun-%d-%d",
		context.block.Height, context.index))

	abciResponses, err := execBlockOnProxyApp(context)

	if !context.stopped {
		context.result = &executionResult{
			abciResponses, err,
		}
		trace.GetElapsedInfo().AddInfo(trace.Prerun, trc.Format())
	}
	PreTimeOut(context.block.Height, int(context.index))
	context.dump("Prerun completed")
	context.prerunResultChan <- context
}

func (blockExec *BlockExecutor) InitPrerun() {
	blockExec.proactivelyRunTx = true
	go blockExec.prerunRoutine()
}

//func FirstBlock(block *types.Block) bool {
//	if 	block.Height == 1{
//		return true
//	}
//	return false
//}

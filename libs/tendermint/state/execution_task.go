package state

import (
	"encoding/hex"
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	cfg "github.com/okex/exchain/libs/tendermint/config"

	"github.com/okex/exchain/libs/tendermint/libs/automation"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
)

type executionResult struct {
	res *ABCIResponses
	err error
}

type executionTask struct {
	height         int64
	index          int64
	block          *types.Block
	stopped        bool
	taskResultChan chan *executionTask
	result         *executionResult
	proxyApp       proxy.AppConnConsensus
	db             dbm.DB
	logger         log.Logger
	blockHash      string
}

func newExecutionTask(blockExec *BlockExecutor, block *types.Block, index int64) *executionTask {
	ret := &executionTask{
		height:         block.Height,
		block:          block,
		db:             blockExec.db,
		proxyApp:       blockExec.proxyApp,
		logger:         blockExec.logger,
		taskResultChan: blockExec.prerunCtx.taskResultChan,
		index:          index,
	}
	ret.blockHash = hex.EncodeToString(block.Hash())

	return ret
}

func (e *executionTask) dump(when string) {

	e.logger.Info(when,
		"stopped", e.stopped,
		"Height", e.block.Height,
		"index", e.index,
		"blockHash", e.blockHash,
		//"AppHash", e.block.AppHash,
	)
}

func (t *executionTask) stop() {
	if t.stopped {
		return
	}

	t.stopped = true
}

func (t *executionTask) run() {
	t.dump("Start prerun")
	trc := trace.NewTracer("lastRun")

	var abciResponses *ABCIResponses
	var err error

	mode := DeliverTxsExecMode(cfg.DynamicConfig.GetDeliverTxsExecuteMode())
	switch mode {
	case DeliverTxsExecModeSerial:
		abciResponses, err = execBlockOnProxyApp(t)
	case DeliverTxsExecModePartConcurrent:
		abciResponses, err = execBlockOnProxyAppPartConcurrent(t.logger, t.proxyApp, t.block, t.db)
	case DeliverTxsExecModeParallel:
		abciResponses, err = execBlockOnProxyAppAsync(t.logger, t.proxyApp, t.block, t.db)
	default:
		abciResponses, err = execBlockOnProxyApp(t)
	}

	if !t.stopped {
		t.result = &executionResult{
			abciResponses, err,
		}
		trace.GetElapsedInfo().AddInfo(trace.Prerun, fmt.Sprintf("num<%d>, lastRun<%s>", t.index, trc.Format()))
	}
	automation.PrerunTimeOut(t.block.Height, int(t.index)-1)
	t.dump("Prerun completed")
	t.taskResultChan <- t
}

//========================================================
func (blockExec *BlockExecutor) InitPrerun() {
	if blockExec.deltaContext.downloadDelta {
		panic("download delta is not allowed if prerun enabled")
	}
	go blockExec.prerunCtx.prerunRoutine()
}

func (blockExec *BlockExecutor) NotifyPrerun(block *types.Block) {
	if block.Height == 1+types.GetStartBlockHeight() {
		return
	}
	blockExec.prerunCtx.notifyPrerun(blockExec, block)
}

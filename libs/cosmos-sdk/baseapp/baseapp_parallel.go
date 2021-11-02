package baseapp

import (
	"encoding/hex"
	"fmt"
	"sync"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func (app *BaseApp) PrepareParallelTxs(txs [][]byte) []*abci.ResponseDeliverTx {
	app.parallelTxManage.isAsyncDeliverTx = true
	evmIndex := uint32(0)
	for k, v := range txs {
		tx, err := app.txDecoder(v)
		if err != nil {
			panic(err)
		}
		t := &txStatus{
			indexInBlock: uint32(k),
		}
		fee, isEvm := app.getTxFee(tx)
		if isEvm {
			t.evmIndex = evmIndex
			t.isEvmTx = true
			evmIndex++
		}

		vString := string(v)
		app.parallelTxManage.setFee(vString, fee)

		app.parallelTxManage.txStatus[vString] = t
		app.parallelTxManage.indexMapBytes = append(app.parallelTxManage.indexMapBytes, vString)
	}

	return app.runTxs(txs)

}

func (app *BaseApp) runTxs(txs [][]byte) []*abci.ResponseDeliverTx {
	maxGas := app.getMaximumBlockGas()

	sumGas := uint64(0)
	overFlow := func(sumGas uint64, currGas int64, maxGas uint64) bool {
		if maxGas <= 0 {
			return false
		}
		if sumGas+uint64(currGas) >= maxGas { // TODO : fix later
			return true
		}
		return false
	}

	asCache := newAsyncCache()
	signal := make(chan int, 1)
	rerunIdx := 0
	txIndex := 0
	txReps := make([]*executeResult, len(txs))
	deliverTxs := make([]*abci.ResponseDeliverTx, len(txs))
	asyncCb := func(execRes *executeResult) {
		txReps[execRes.GetCounter()] = execRes
		for txReps[txIndex] != nil {
			s := app.parallelTxManage.txStatus[app.parallelTxManage.indexMapBytes[txIndex]]

			res := txReps[txIndex]
			if res.Conflict(asCache) || overFlow(sumGas, txReps[txIndex].resp.GasUsed, maxGas) {
				rerunIdx++
				s.reRun = true
				res = app.deliverTxWithCache(abci.RequestDeliverTx{Tx: txs[res.GetCounter()]})

			}
			txRs := res.GetResponse()
			deliverTxs[txIndex] = &txRs
			res.Collect(asCache)
			res.Commit()
			if !app.parallelTxManage.isReRun(app.parallelTxManage.indexMapBytes[txIndex]) {
				app.deliverState.ctx.BlockGasMeter().ConsumeGas(sdk.Gas(res.resp.GasUsed), "unexpected error")
			}
			if s.anteErr != nil {
				app.logger.Error("paralleled-tx", "txIndex", txIndex, "anteErr", s.anteErr)
			}
			sumGas += uint64(res.resp.GasUsed)

			txIndex++
			if txIndex == len(txs) {
				app.logger.Info("Paralleled-tx", "len(txs)", len(txs), "Parallel run", len(txs)-rerunIdx, "ReRun", rerunIdx)
				signal <- 0
				return
			}
		}
	}

	app.parallelTxManage.workgroup.cb = asyncCb

	// Run txs of block.
	for _, tx := range txs {
		go app.DeliverTx(abci.RequestDeliverTx{Tx: tx})
	}

	if len(txs) > 0 {
		//waiting for call back
		<-signal
		receiptsLogs := app.endParallelTxs()
		for index, v := range receiptsLogs {
			if len(v) != 0 { // only update evm tx result
				deliverTxs[index].Data = v
			}
		}

	}
	return deliverTxs
}

func (app *BaseApp) endParallelTxs() [][]byte {
	txFeeInBlock := sdk.Coins{}
	feeMap := app.parallelTxManage.getFeeMap()
	refundMap := app.parallelTxManage.getRefundFeeMap()
	for tx, v := range feeMap {
		if app.parallelTxManage.txStatus[tx].anteErr != nil {
			continue
		}
		txFeeInBlock = txFeeInBlock.Add(v...)
		if refundFee, ok := refundMap[tx]; ok {
			txFeeInBlock = txFeeInBlock.Sub(refundFee)
		}
	}
	ctx, cache := app.cacheTxContext(app.getContextForTx(runTxModeDeliverInAsync, []byte{}), []byte{})
	if err := app.updateFeeCollectorAccHandler(ctx, txFeeInBlock); err != nil {
		panic(err)
	}
	cache.Write()

	txExecStats := make([][]string, 0)
	for _, v := range app.parallelTxManage.indexMapBytes {
		errMsg := ""
		if err := app.parallelTxManage.txStatus[v].anteErr; err != nil {
			errMsg = err.Error()
		}
		txExecStats = append(txExecStats, []string{v, errMsg})
	}
	app.parallelTxManage.clear()
	return app.logFix(txExecStats)
}

//we reuse the nonce that changed by the last async call
//if last ante handler has been failed, we need rerun it ? or not?
func (app *BaseApp) deliverTxWithCache(req abci.RequestDeliverTx) *executeResult {
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return nil
	}
	var (
		resp abci.ResponseDeliverTx
		mode runTxMode
	)
	mode = runTxModeDeliverInAsync
	g, r, m, e := app.runTx(mode, req.Tx, tx, LatestSimulateTxHeight)
	if e != nil {
		fmt.Println("err", err, g.GasWanted, g.GasUsed)
		resp = sdkerrors.ResponseDeliverTx(e, g.GasWanted, g.GasUsed, app.trace)
	} else {
		resp = abci.ResponseDeliverTx{
			GasWanted: int64(g.GasWanted), // TODO: Should type accept unsigned ints?
			GasUsed:   int64(g.GasUsed),   // TODO: Should type accept unsigned ints?
			Log:       r.Log,
			Data:      r.Data,
			Events:    r.Events.ToABCIEvents(),
		}
	}

	txStatus := app.parallelTxManage.txStatus[string(req.Tx)]
	asyncExe := newExecuteResult(resp, m, txStatus.indexInBlock, txStatus.evmIndex)
	asyncExe.err = e
	return asyncExe
}

type executeResult struct {
	resp       abci.ResponseDeliverTx
	ms         sdk.CacheMultiStore
	counter    uint32
	err        error
	evmCounter uint32
}

func (e executeResult) GetResponse() abci.ResponseDeliverTx {
	return e.resp
}

func (e executeResult) Conflict(cache *asyncCache) bool {
	rerun := false
	if e.ms == nil {
		return true //TODO fix later
	}

	e.ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		//the key we have read was wrote by pre txs
		if cache.Has(key) && !whiteAccountList[hex.EncodeToString(key)] {
			rerun = true
			return false // break
		}
		return true
	})
	return rerun
}

var (
	whiteAccountList = map[string]bool{
		"01f1829676db577682e944fc3493d451b67ff3e29f": true, //fee
	}
)

func (e executeResult) Collect(cache *asyncCache) {
	if e.ms == nil {
		return
	}
	e.ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		if isDirty {
			//push every data we have written in current tx
			cache.Push(key, value)
		}
		return true
	})
}

func (e executeResult) GetCounter() uint32 {
	return e.counter
}

func (e executeResult) Commit() {
	if e.ms == nil {
		return
	}
	e.ms.Write()
}

func newExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32, evmCounter uint32) *executeResult {
	return &executeResult{
		resp:       r,
		ms:         ms,
		counter:    counter,
		evmCounter: evmCounter,
	}
}

type asyncWorkGroup struct {
	workCh chan *executeResult
	cb     func(*executeResult)
}

func newAsyncWorkGroup() *asyncWorkGroup {
	return &asyncWorkGroup{
		workCh: make(chan *executeResult, 64),
		cb:     nil,
	}
}

func (a *asyncWorkGroup) Push(item *executeResult) {
	a.workCh <- item
}

func (a *asyncWorkGroup) Start() {
	go func() {
		for {
			select {
			case exec := <-a.workCh:
				if a.cb != nil {
					a.cb(exec)
				}
			}
		}
	}()
}

type parallelTxManager struct {
	mu               sync.RWMutex
	isAsyncDeliverTx bool
	workgroup        *asyncWorkGroup

	fee       map[string]sdk.Coins
	refundFee map[string]sdk.Coins

	txStatus      map[string]*txStatus
	indexMapBytes []string
}

type txStatus struct {
	reRun        bool
	isEvmTx      bool
	evmIndex     uint32
	indexInBlock uint32
	anteErr      error
}

func newParallelTxManager() *parallelTxManager {
	return &parallelTxManager{
		isAsyncDeliverTx: false,
		workgroup:        newAsyncWorkGroup(),
		fee:              make(map[string]sdk.Coins),
		refundFee:        make(map[string]sdk.Coins),

		txStatus:      make(map[string]*txStatus),
		indexMapBytes: make([]string, 0),
	}
}

func (f *parallelTxManager) clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee = make(map[string]sdk.Coins)
	f.refundFee = make(map[string]sdk.Coins)

	f.txStatus = make(map[string]*txStatus)
	f.indexMapBytes = make([]string, 0)

}
func (f *parallelTxManager) setFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee[key] = value
}

func (f *parallelTxManager) getFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.fee
}
func (f *parallelTxManager) setRefundFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.refundFee[key] = value
}

func (f *parallelTxManager) getRefundFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.refundFee
}

func (f *parallelTxManager) isReRun(tx string) bool {
	data, ok := f.txStatus[tx]
	if !ok {
		return false
	}
	return data.reRun
}

type asyncCache struct {
	mem map[string][]byte
}

func newAsyncCache() *asyncCache {
	return &asyncCache{mem: make(map[string][]byte)}
}

func (a *asyncCache) Push(key, value []byte) {
	a.mem[string(key)] = value
}

func (a *asyncCache) Has(key []byte) bool {
	_, ok := a.mem[string(key)]
	return ok
}

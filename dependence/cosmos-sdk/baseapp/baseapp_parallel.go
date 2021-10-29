package baseapp

import (
	"encoding/hex"
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (app *BaseApp) PrepareParallelTxs(cb abci.AsyncCallBack, txs [][]byte) {
	fmt.Println("start", len(app.parallelTxManage.workgroup.WorkCh))
	app.parallelTxManage.workgroup.Cb = cb
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
		app.parallelTxManage.SetFee(vString, fee)

		app.parallelTxManage.txStatus[vString] = t
		app.parallelTxManage.indexMapBytes = append(app.parallelTxManage.indexMapBytes, vString)
	}
}

func (app *BaseApp) EndParallelTxs() [][]byte {
	txFeeInBlock := sdk.Coins{}
	feeMap := app.parallelTxManage.GetFeeMap()
	refundMap := app.parallelTxManage.GetRefundFeeMap()
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
	app.parallelTxManage.Clear()
	fmt.Println("endpara", len(app.parallelTxManage.workgroup.WorkCh))
	return app.logFix(txExecStats)
}

//we reuse the nonce that changed by the last async call
//if last ante handler has been failed, we need rerun it ? or not?
func (app *BaseApp) DeliverTxWithCache(req abci.RequestDeliverTx) abci.ExecuteRes {
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return nil
	}
	var (
		gInfo sdk.GasInfo
		resp  abci.ResponseDeliverTx
		mode  runTxMode
	)
	mode = runTxModeDeliverInAsync
	g, r, m, e := app.runTx(mode, req.Tx, tx, LatestSimulateTxHeight)
	if e != nil {
		resp = sdkerrors.ResponseDeliverTx(e, gInfo.GasWanted, gInfo.GasUsed, app.trace)
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
	asyncExe := NewExecuteResult(resp, m, txStatus.indexInBlock, txStatus.evmIndex)
	asyncExe.err = e
	return asyncExe
}

type ExecuteResult struct {
	Resp       abci.ResponseDeliverTx
	Ms         sdk.CacheMultiStore
	Counter    uint32
	err        error
	evmCounter uint32
}

func (e ExecuteResult) GetResponse() abci.ResponseDeliverTx {
	return e.Resp
}

func (e ExecuteResult) Conflict(cache abci.AsyncCacheInterface) bool {
	rerun := false
	if e.Ms == nil {
		return true //TODO fix later
	}

	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		//the key we have read was wrote by pre txs
		if cache.Has(key) && !whiteAccountList[hex.EncodeToString(key)] {
			rerun = true
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

func (e ExecuteResult) Collect(cache abci.AsyncCacheInterface) {
	if e.Ms == nil {
		return
	}
	e.Ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		if isDirty {
			//push every data we have written in current tx
			cache.Push(key, value)
		}
		return true
	})
}

func (e ExecuteResult) GetCounter() uint32 {
	return e.Counter
}

func (e ExecuteResult) Commit() {
	if e.Ms == nil {
		return
	}
	e.Ms.Write()
}

func NewExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32, evmCounter uint32) ExecuteResult {
	return ExecuteResult{
		Resp:       r,
		Ms:         ms,
		Counter:    counter,
		evmCounter: evmCounter,
	}
}

type AsyncWorkGroup struct {
	WorkCh chan ExecuteResult
	Cb     abci.AsyncCallBack
}

func NewAsyncWorkGroup() *AsyncWorkGroup {
	return &AsyncWorkGroup{
		WorkCh: make(chan ExecuteResult, 64),
		Cb:     nil,
	}
}

func (a *AsyncWorkGroup) Push(item ExecuteResult) {
	fmt.Println("llll", len(a.WorkCh), item.Counter)
	a.WorkCh <- item
}

func (a *AsyncWorkGroup) Start() {
	go func() {
		for {
			select {
			case exec := <-a.WorkCh:
				if a.Cb != nil {
					a.Cb(exec)
				}
			}
		}
	}()
}

type parallelTxManager struct {
	mu               sync.RWMutex
	isAsyncDeliverTx bool
	workgroup        *AsyncWorkGroup

	fee       map[string]sdk.Coins
	refundFee map[string]sdk.Coins

	txStatus      map[string]*txStatus
	indexMapBytes []string
}

type txStatus struct {
	isEvmTx      bool
	evmIndex     uint32
	indexInBlock uint32
	anteErr      error
}

func newParallelTxManager() *parallelTxManager {
	return &parallelTxManager{
		isAsyncDeliverTx: false,
		workgroup:        NewAsyncWorkGroup(),
		fee:              make(map[string]sdk.Coins),
		refundFee:        make(map[string]sdk.Coins),

		txStatus:      make(map[string]*txStatus),
		indexMapBytes: make([]string, 0),
	}
}

func (f *parallelTxManager) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee = make(map[string]sdk.Coins)
	f.refundFee = make(map[string]sdk.Coins)

	f.txStatus = make(map[string]*txStatus)
	f.indexMapBytes = make([]string, 0)

}
func (f *parallelTxManager) SetFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee[key] = value
}

func (f *parallelTxManager) GetFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.fee
}
func (f *parallelTxManager) SetRefundFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.refundFee[key] = value
}

func (f *parallelTxManager) GetRefundFeeMap() map[string]sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.refundFee
}

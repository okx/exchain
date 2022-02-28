package baseapp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"sync"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

var (
	txIndexLen = 4
)

type extraDataForTx struct {
	fee       sdk.Coins
	isEvm     bool
	signCache sdk.SigCache
	to        *ethcommon.Address
}

// txByteWithIndex = txByte + index

func getTxByteWithIndex(txByte []byte, txIndex int) []byte {
	bs := make([]byte, txIndexLen)
	binary.LittleEndian.PutUint32(bs, uint32(txIndex))
	return append(txByte, bs...)
}

func getRealTxByte(txByteWithIndex []byte) []byte {
	return txByteWithIndex[:len(txByteWithIndex)-txIndexLen]

}

func (app *BaseApp) getExtraDataByTxs(txs [][]byte) []*extraDataForTx {
	res := make([]*extraDataForTx, len(txs), len(txs))
	var wg sync.WaitGroup
	for index, txBytes := range txs {
		wg.Add(1)
		index := index
		txBytes := txBytes
		go func() {
			defer wg.Done()
			tx, err := app.txDecoder(txBytes)
			if err != nil {
				res[index] = &extraDataForTx{}
				return
			}
			coin, isEvm, s, toAddr := app.getTxFee(app.getContextForTx(runTxModeDeliver, txBytes), tx)
			res[index] = &extraDataForTx{
				fee:       coin,
				isEvm:     isEvm,
				signCache: s,
				to:        toAddr,
			}
		}()
	}
	wg.Wait()
	return res
}

var (
	rootAddr = make(map[ethcommon.Address]ethcommon.Address, 0)
)

func Find(x ethcommon.Address) ethcommon.Address {
	if rootAddr[x] != x {
		rootAddr[x] = Find(rootAddr[x])
	}
	return rootAddr[x]
}

func Union(x ethcommon.Address, y *ethcommon.Address) {
	if _, ok := rootAddr[x]; !ok {
		rootAddr[x] = x
	}
	if y == nil {
		return
	}
	if _, ok := rootAddr[*y]; !ok {
		rootAddr[*y] = *y
	}
	fx := Find(x)
	fy := Find(*y)
	if fx != fy {
		rootAddr[fy] = fx
	}
}

func (app *BaseApp) calGroup(txsExtraData []*extraDataForTx) (map[int][]int, map[int]int) {
	rootAddr = make(map[ethcommon.Address]ethcommon.Address, 0)
	app.parallelTxManage.txReps = make([]*executeResult, len(txsExtraData))
	for index, tx := range txsExtraData {
		if tx.isEvm { //evmTx
			Union(tx.signCache.GetFrom(), tx.to)
		} else {
			fmt.Println("not evmTx", index)
			app.parallelTxManage.txReps[index] = &executeResult{}
		}
	}

	groupList := make(map[int][]int, 0)
	addrToID := make(map[ethcommon.Address]int, 0)

	for index, sender := range txsExtraData {
		if !sender.isEvm {
			continue
		}
		rootAddr := Find(sender.signCache.GetFrom())
		id, exist := addrToID[rootAddr]
		if !exist {
			id = len(groupList)
			addrToID[rootAddr] = id

		}
		groupList[id] = append(groupList[id], index)
	}

	nextTxIndexInGroup := make(map[int]int)
	preTxIndexInGroup := make(map[int]int)
	for _, list := range groupList {
		for index := 0; index < len(list); index++ {
			if index+1 <= len(list)-1 {
				nextTxIndexInGroup[list[index]] = list[index+1]
			}
			if index-1 >= 0 {
				preTxIndexInGroup[list[index]] = list[index-1]
			}
		}
	}
	app.parallelTxManage.nextTxInGroup = nextTxIndexInGroup
	app.parallelTxManage.preTxInGroup = preTxIndexInGroup
	return groupList, nextTxIndexInGroup

}

func (app *BaseApp) ParallelTxs(txs [][]byte) []*abci.ResponseDeliverTx {
	ts := time.Now()
	txWithIndex := make([][]byte, 0)
	for index, v := range txs {
		txWithIndex = append(txWithIndex, getTxByteWithIndex(v, index))
	}

	extraData := app.getExtraDataByTxs(txs)

	groupList, nextIndexInGroup := app.calGroup(extraData)

	app.parallelTxManage.isAsyncDeliverTx = true
	app.parallelTxManage.cms = app.deliverState.ms.CacheMultiStore()

	evmIndex := uint32(0)
	for k := range txs {
		t := &txStatus{
			indexInBlock: uint32(k),
			signCache:    extraData[k].signCache,
		}
		if extraData[k].isEvm {
			t.evmIndex = evmIndex
			t.isEvmTx = true
			evmIndex++
		}

		vString := string(txWithIndex[k])
		app.parallelTxManage.setFee(vString, extraData[k].fee)

		app.parallelTxManage.txStatus[vString] = t
		app.parallelTxManage.indexMapBytes = append(app.parallelTxManage.indexMapBytes, vString)
	}

	fmt.Println("calGroupTime", time.Now().Sub(ts).Milliseconds())
	return app.runTxs(txWithIndex, groupList, nextIndexInGroup)

}

//TODO: fuck
func (app *BaseApp) fixFeeCollector(txString string, ms sdk.CacheMultiStore) {
	if app.parallelTxManage.txStatus[txString].anteErr != nil {
		return
	}

	txFee := app.parallelTxManage.getFee(txString)
	refundFee := app.parallelTxManage.getRefundFee(txString)
	txFee = txFee.Sub(refundFee)
	app.parallelTxManage.currTxFee = app.parallelTxManage.currTxFee.Add(txFee...)

	ctx, cache := app.cacheTxContext(app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})

	ctx = ctx.WithMultiStore(ms)
	if err := app.updateFeeCollectorAccHandler(ctx, app.parallelTxManage.currTxFee); err != nil {
		panic(err)
	}
	cache.Write()
}

func (app *BaseApp) runTxs(txs [][]byte, groupList map[int][]int, nextTxInGroup map[int]int) []*abci.ResponseDeliverTx {
	//fmt.Println("detail", app.deliverState.ctx.BlockHeight(), "len(group)", len(groupList))
	//for index := 0; index < len(groupList); index++ {
	//	fmt.Println("groupIndex", index, "groupSize", len(groupList[index]), "list", groupList[index])
	//}
	maxGas := app.getMaximumBlockGas()
	currentGas := uint64(0)
	overFlow := func(sumGas uint64, currGas int64, maxGas uint64) bool {
		if maxGas <= 0 {
			return false
		}
		if sumGas+uint64(currGas) >= maxGas { // TODO : fix later
			return true
		}
		return false
	}

	asCache := newAsyncCache(len(txs))
	signal := make(chan int, 1)
	rerunIdx := 0
	txIndex := 0

	pm := app.parallelTxManage

	txReps := pm.txReps
	deliverTxs := make([]*abci.ResponseDeliverTx, len(txs))

	asyncCb := func(execRes *executeResult) {
		receiveTxIndex := int(execRes.GetCounter())

		pm.setTxStatus(int(receiveTxIndex), false)
		txReps[receiveTxIndex] = execRes

		if pm.isFailed(pm.runningStats(int(receiveTxIndex))) {
			txReps[receiveTxIndex] = nil
			pm.setTxStatus(receiveTxIndex, true)
			//fmt.Println("already failed:ReRnn", receiveTxIndex, "current", txIndex)
			go app.asyncDeliverTx(txs[receiveTxIndex], receiveTxIndex)

		}

		//fmt.Println("receiveTxIndex", receiveTxIndex, "currNeedTxIndex", txIndex)
		if receiveTxIndex >= txIndex || receiveTxIndex == 0 {
			if nextTx := nextTxInGroup[receiveTxIndex]; nextTx != 0 && pm.runningStats(nextTx) == -1 {

				needRun := true
				if psb, ok := pm.preTxInGroup[receiveTxIndex]; ok {
					if psb > txIndex {
						needRun = false
					}
				}

				if needRun {
					pm.setTxStatus(nextTx, true)
					//fmt.Println("scf:receive:next", receiveTxIndex, nextTx, "curr", txIndex)
					go app.asyncDeliverTx(txs[nextTx], nextTx)
				}
			}
		}

		if txIndex != receiveTxIndex {
			return
		}
		for txReps[txIndex] != nil {
			s := pm.txStatus[app.parallelTxManage.indexMapBytes[txIndex]]
			res := txReps[txIndex]

			if res.Conflict(pm.getRunBase(int(res.counter)), asCache) || overFlow(currentGas, res.resp.GasUsed, maxGas) {
				rerunIdx++
				s.reRun = true
				//fmt.Println("Need ReRun", txIndex)
				app.parallelTxManage.setTxStatus(txIndex, true)
				res = app.deliverTxWithCache(txs[txIndex])
				app.parallelTxManage.setTxStatus(txIndex, false)
				txReps[txIndex] = res

				nn, ok := app.parallelTxManage.nextTxInGroup[txIndex]

				if ok {
					pp := nn
					for true {
						txReps[pp] = nil
						pp, ok = app.parallelTxManage.nextTxInGroup[pp]
						if !ok {
							break
						}
					}

					runningTaskID := pm.runningStats(nn)

					if runningTaskID == -1 {
						//fmt.Println("ReRunTxEnd, StartNextInGroup", txIndex, nn)
						txReps[nn] = nil
						app.parallelTxManage.setTxStatus(nn, true)
						app.asyncDeliverTx(txs[nn], nn)
					} else {
						//fmt.Println("markFailed", nn, runningTaskID)
						pm.markFailed(runningTaskID)
					}

				}

			}
			if s.anteErr != nil {
				res.ms = nil
			}

			txRs := res.GetResponse()
			deliverTxs[txIndex] = &txRs

			res.Collect(txIndex, asCache)
			app.fixFeeCollector(app.parallelTxManage.indexMapBytes[txIndex], res.ms)

			if !s.reRun {
				app.deliverState.ctx.BlockGasMeter().ConsumeGas(sdk.Gas(res.resp.GasUsed), "unexpected error")
			}

			pm.SetCurrentIndex(txIndex, res) //Commit

			currentGas += uint64(res.resp.GasUsed)
			txIndex++
			if txIndex == len(txs) {
				ParaLog.Update(uint64(app.deliverState.ctx.BlockHeight()), len(txs), rerunIdx)
				app.logger.Info("Paralleled-tx", "blockHeight", app.deliverState.ctx.BlockHeight(), "len(txs)", len(txs), "Parallel run", len(txs)-rerunIdx, "ReRun", rerunIdx)
				signal <- 0
				return
			}
			if txReps[txIndex] == nil && !pm.isRunning(txIndex) {
				app.parallelTxManage.setTxStatus(txIndex, true)
				go app.asyncDeliverTx(txs[txIndex], txIndex)
			}

		}
	}

	pm.workgroup.cb = asyncCb

	for _, group := range groupList {
		txIndex := group[0]
		app.parallelTxManage.setTxStatus(txIndex, true)
		go app.asyncDeliverTx(txs[txIndex], txIndex)
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
	pm.cms.Write()
	return deliverTxs
}

func (app *BaseApp) endParallelTxs() [][]byte {

	txExecStats := make([][]string, 0)
	for _, v := range app.parallelTxManage.indexMapBytes {
		errMsg := ""
		if err := app.parallelTxManage.txStatus[v].anteErr; err != nil {
			errMsg = err.Error()
		}
		txExecStats = append(txExecStats, []string{string(getRealTxByte([]byte(v))), errMsg})
	}
	app.parallelTxManage.clear()
	return app.logFix(txExecStats)
}

//we reuse the nonce that changed by the last async call
//if last ante handler has been failed, we need rerun it ? or not?
func (app *BaseApp) deliverTxWithCache(txByte []byte) *executeResult {
	txStatus := app.parallelTxManage.txStatus[string(txByte)]

	tx, err := app.txDecoder(getRealTxByte(txByte))
	if err != nil {
		asyncExe := newExecuteResult(sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace), nil, txStatus.indexInBlock, txStatus.evmIndex)
		return asyncExe
	}
	var (
		resp abci.ResponseDeliverTx
		mode runTxMode
	)
	mode = runTxModeDeliverInAsync
	g, r, m, e := app.runTx(mode, txByte, tx, LatestSimulateTxHeight)
	if e != nil {
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

func (e executeResult) Conflict(from int, cache *asyncCache) bool {
	//fmt.Println("Ready CheckConflict", "base", from, "now", e.counter)
	rerun := false
	if e.ms == nil {
		return true //TODO fix later
	}

	e.ms.IteratorCache(func(key, value []byte, isDirty bool, isDelete bool, s sdk.StoreKey) bool {
		//the key we have read was written by pre txs
		if cache.Has(key, from, int(e.counter)) && !whiteAccountList[hex.EncodeToString(key)] {
			//fmt.Println("key", hex.EncodeToString(key), hex.EncodeToString(value))
			rerun = true
			return false // break
		}
		return true
	}, nil)
	return rerun
}

var (
	whiteAccountList = map[string]bool{
		"01f1829676db577682e944fc3493d451b67ff3e29f": true, //fee
	}
)

func (e executeResult) Collect(index int, cache *asyncCache) {
	if e.ms == nil {
		return
	}
	e.ms.IteratorCache(func(key, value []byte, isDirty bool, isDelete bool, s sdk.StoreKey) bool {
		if isDirty {
			//push every data we have written in current tx
			cache.Push(index, key)
		}
		return true
	}, nil)
}

func (e executeResult) GetCounter() uint32 {
	return e.counter
}

//func (e executeResult) Commit() {
//	if e.ms == nil {
//		return
//	}
//	e.ms.Write()
//}

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

	currTxFee     sdk.Coins
	runningStatus map[int]int
	isrunning     map[int]bool

	txReps        []*executeResult
	nextTxInGroup map[int]int
	preTxInGroup  map[int]int

	cms       sdk.CacheMultiStore
	currIndex int

	runBase         map[int]int
	indexInAll      int
	markFailedStats map[int]bool
}

type runStats struct {
	base       int
	indexInAll int
}

type txStatus struct {
	reRun        bool
	isEvmTx      bool
	evmIndex     uint32
	indexInBlock uint32
	anteErr      error
	signCache    sdk.SigCache
}

func newParallelTxManager() *parallelTxManager {
	return &parallelTxManager{
		isAsyncDeliverTx: false,
		workgroup:        newAsyncWorkGroup(),
		fee:              make(map[string]sdk.Coins),
		refundFee:        make(map[string]sdk.Coins),

		txStatus:      make(map[string]*txStatus),
		indexMapBytes: make([]string, 0),
		runningStatus: make(map[int]int),
		isrunning:     make(map[int]bool),

		nextTxInGroup: make(map[int]int),
		preTxInGroup:  make(map[int]int),

		currIndex:       -1,
		runBase:         make(map[int]int),
		markFailedStats: make(map[int]bool),
	}
}

func (f *parallelTxManager) clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee = make(map[string]sdk.Coins)
	f.refundFee = make(map[string]sdk.Coins)

	f.txStatus = make(map[string]*txStatus)
	f.indexMapBytes = make([]string, 0)
	f.currTxFee = sdk.Coins{}
	f.runningStatus = make(map[int]int)
	f.nextTxInGroup = make(map[int]int)
	f.preTxInGroup = make(map[int]int)
	f.runBase = make(map[int]int)
	f.currIndex = -1
	f.indexInAll = 0
	f.markFailedStats = make(map[int]bool)
}

func (f *parallelTxManager) markFailed(txIndexAll int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.markFailedStats[txIndexAll] = true
}

func (f *parallelTxManager) isFailed(txindexAll int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.markFailedStats[txindexAll]
}

func (f *parallelTxManager) setFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fee[key] = value
}

func (f *parallelTxManager) getFee(key string) sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.fee[key]
}
func (f *parallelTxManager) setRefundFee(key string, value sdk.Coins) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.refundFee[key] = value
}

func (f *parallelTxManager) getRefundFee(key string) sdk.Coins {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.refundFee[key]
}

func (f *parallelTxManager) isReRun(tx string) bool {
	data, ok := f.txStatus[tx]
	if !ok {
		return false
	}
	return data.reRun
}

func (f *parallelTxManager) setTxStatus(txIndex int, status bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if status == true {
		f.runningStatus[txIndex] = f.indexInAll
		f.indexInAll++
		f.isrunning[txIndex] = true
	} else {
		f.isrunning[txIndex] = false
	}
}

func (f *parallelTxManager) runningStats(txIndex int) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.runningStatus[txIndex]
}

func (f *parallelTxManager) isRunning(txIndex int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.isrunning[txIndex]
}

func (f *parallelTxManager) getTxResult(tx []byte) sdk.CacheMultiStore {
	index := f.txStatus[string(tx)].indexInBlock

	preIndexInGroup, ok := f.preTxInGroup[int(index)]
	f.mu.Lock()
	defer f.mu.Unlock()
	ms := f.cms.CacheMultiStore()
	base := f.currIndex
	if ok && preIndexInGroup > f.currIndex {
		if f.txStatus[f.indexMapBytes[preIndexInGroup]].anteErr == nil {
			ms = f.txReps[preIndexInGroup].ms.CacheMultiStore()
			base = preIndexInGroup
		} else {
			ms = f.cms.CacheMultiStore()
			base = f.currIndex
		}

	}

	f.runBase[int(index)] = base
	//fmt.Println("index", index, "base", base)
	return ms
}

func (f *parallelTxManager) getRunBase(now int) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.runBase[now]
}

func (f *parallelTxManager) SetCurrentIndex(d int, res *executeResult) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if res.ms == nil {
		return
	}

	res.ms.IteratorCache(func(key, value []byte, isDirty bool, isdelete bool, storeKey sdk.StoreKey) bool {
		if isDirty {
			if isdelete {
				f.cms.GetKVStore(storeKey).Delete(key)
			} else if value != nil {
				f.cms.GetKVStore(storeKey).Set(key, value)
			}
		}
		return true
	}, nil)
	f.cms.Write()

	f.currIndex = d
}

type asyncCache struct {
	mem []map[string]bool
}

func newAsyncCache(txSize int) *asyncCache {
	tt := make([]map[string]bool, 0)
	for index := 0; index < txSize; index++ {
		tt = append(tt, make(map[string]bool))
	}

	return &asyncCache{mem: tt}
}

func (a *asyncCache) Push(index int, key []byte) {
	a.mem[index][string(key)] = true
}

func (a *asyncCache) Has(key []byte, from int, to int) bool {
	for index := from + 1; index < to; index++ {
		if _, ok := a.mem[index][string(key)]; ok {
			return true
		}
	}
	return false
}

var (
	ParaLog *LogForParallel
)

func init() {
	ParaLog = NewLogForParallel()
}

type parallelBlockInfo struct {
	height   uint64
	txs      int
	reRunTxs int
}

func (p parallelBlockInfo) better(n parallelBlockInfo) bool {
	return 1-float64(p.reRunTxs)/float64(p.txs) > 1-float64(n.reRunTxs)/float64(n.txs)
}

func (p parallelBlockInfo) string() string {
	return fmt.Sprintf("Height:%d Txs %d ReRunTxs %d", p.height, p.txs, p.reRunTxs)
}

type LogForParallel struct {
	init         bool
	sumTx        int
	reRunTx      int
	blockNumbers int

	bestBlock     parallelBlockInfo
	terribleBlock parallelBlockInfo
}

func NewLogForParallel() *LogForParallel {
	return &LogForParallel{
		sumTx:        0,
		reRunTx:      0,
		blockNumbers: 0,
		bestBlock: parallelBlockInfo{
			height:   0,
			txs:      0,
			reRunTxs: 0,
		},
		terribleBlock: parallelBlockInfo{
			height:   0,
			txs:      0,
			reRunTxs: 0,
		},
	}
}

func (l *LogForParallel) Update(height uint64, txs int, reRunCnt int) {
	l.sumTx += txs
	l.reRunTx += reRunCnt
	l.blockNumbers++

	if txs < 20 {
		return
	}

	info := parallelBlockInfo{height: height, txs: txs, reRunTxs: reRunCnt}
	if !l.init {
		l.bestBlock = info
		l.terribleBlock = info
		l.init = true
		return
	}

	if info.better(l.bestBlock) {
		l.bestBlock = info
	}
	if l.terribleBlock.better(info) {
		l.terribleBlock = info
	}
}

func (l *LogForParallel) PrintLog() {
	fmt.Println("BlockNumbers", l.blockNumbers)
	fmt.Println("AllTxs", l.sumTx)
	fmt.Println("ReRunTxs", l.reRunTx)
	fmt.Println("All Concurrency Rate", float64(l.reRunTx)/float64(l.sumTx))
	fmt.Println("BestBlock", l.bestBlock.string(), "Concurrency Rate", 1-float64(l.bestBlock.reRunTxs)/float64(l.bestBlock.txs))
	fmt.Println("TerribleBlock", l.terribleBlock.string(), "Concurrency Rate", 1-float64(l.terribleBlock.reRunTxs)/float64(l.terribleBlock.txs))
}

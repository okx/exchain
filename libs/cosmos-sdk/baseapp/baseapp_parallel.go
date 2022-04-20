package baseapp

import (
	"bytes"
	"encoding/binary"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/spf13/viper"
	"runtime"
	"sync"
)

var (
	txIndexLen = 4
)

type extraDataForTx struct {
	fee          sdk.Coins
	isEvm        bool
	from         string
	to           *ethcommon.Address
	reRun        bool
	evmIndex     uint32
	indexInBlock uint32
	stdTx        sdk.Tx
	decodeErr    error
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

func (app *BaseApp) getExtraDataByTxs(txs [][]byte) {
	para := app.parallelTxManage
	if para.txSize > len(para.extraTxsInfo) {
		para.txReps = make([]*executeResult, para.txSize)
		para.extraTxsInfo = make([]*extraDataForTx, para.txSize)
		para.workgroup.runningStatus = make([]int, para.txSize)
		para.workgroup.isrunning = make([]bool, para.txSize)
	}

	var wg sync.WaitGroup
	for index, txBytes := range txs {
		wg.Add(1)
		index := index
		txBytes := txBytes
		go func() {
			defer wg.Done()
			tx, err := app.txDecoder(txBytes)
			if err != nil {
				para.extraTxsInfo[index] = &extraDataForTx{
					decodeErr: err,
				}
				return
			}
			coin, isEvm, s, toAddr := app.getTxFee(app.getContextForTx(runTxModeDeliver, txBytes), tx)
			para.extraTxsInfo[index] = &extraDataForTx{
				fee:   coin,
				isEvm: isEvm,
				from:  s,
				to:    toAddr,
				stdTx: tx,
			}
		}()
	}
	wg.Wait()
}

var (
	rootAddr = make(map[string]string, 0)
)

func Find(x string) string {
	if rootAddr[x] != x {
		rootAddr[x] = Find(rootAddr[x])
	}
	return rootAddr[x]
}

func Union(x string, y *ethcommon.Address) {
	if _, ok := rootAddr[x]; !ok {
		rootAddr[x] = x
	}
	if y == nil {
		return
	}
	yString := y.String()
	if _, ok := rootAddr[yString]; !ok {
		rootAddr[yString] = yString
	}
	fx := Find(x)
	fy := Find(yString)
	if fx != fy {
		rootAddr[fy] = fx
	}
}

func (app *BaseApp) calGroup() {

	para := app.parallelTxManage

	rootAddr = make(map[string]string, 0)
	for index, tx := range para.extraTxsInfo {
		if tx.isEvm { //evmTx
			Union(tx.from, tx.to)
		} else {
			app.parallelTxManage.txReps[index] = &executeResult{paraMsg: &sdk.ParaMsg{}}
		}
	}

	addrToID := make(map[string]int, 0)

	for index, sender := range para.extraTxsInfo {
		if !sender.isEvm {
			continue
		}
		rootAddr := Find(sender.from)
		id, exist := addrToID[rootAddr]
		if !exist {
			id = len(para.groupList)
			addrToID[rootAddr] = id

		}
		//if _, ok := para.groupList[id]; !ok {
		//	para.groupList[id] = make([]int, 0)
		//}
		para.groupList[id] = append(para.groupList[id], index)
	}

	groupSize := len(para.groupList)
	for groupIndex := 0; groupIndex < groupSize; groupIndex++ {
		list := para.groupList[groupIndex]
		for index := 0; index < len(list); index++ {
			if index+1 <= len(list)-1 {
				app.parallelTxManage.nextTxInGroup[list[index]] = list[index+1]
			}
			if index-1 >= 0 {
				app.parallelTxManage.preTxInGroup[list[index]] = list[index-1]
			}
			app.parallelTxManage.txIndexWithGroupID[list[index]] = groupIndex
		}
	}
}

func (app *BaseApp) paraLoadSender(txs [][]byte) {

	checkStateCtx := app.checkState.ctx
	checkStateCtx.SetBlockHeight(app.checkState.ctx.BlockHeight() + 1)

	maxNums := runtime.NumCPU()
	txSize := len(txs)
	if maxNums > txSize {
		maxNums = txSize
	}

	txJobChan := make(chan []byte)
	var wg sync.WaitGroup
	wg.Add(txSize)

	for index := 0; index < maxNums; index++ {
		go func(ch chan []byte, wg *sync.WaitGroup) {
			for txBytes := range ch {
				tx, err := app.txDecoder(txBytes)
				if err == nil {
					ctx := checkStateCtx
					ctx.SetTxBytes(txBytes)
					app.getTxFee(ctx, tx)
				}
				wg.Done()
			}
		}(txJobChan, &wg)
	}
	for _, v := range txs {
		txJobChan <- v
	}

	wg.Wait()
	close(txJobChan)
}

func (app *BaseApp) ParallelTxs(txs [][]byte, onlyCalSender bool) []*abci.ResponseDeliverTx {
	pm := app.parallelTxManage
	txSize := len(txs)
	pm.txSize = txSize
	pm.haveCosmosTxInBlock = false

	if txSize == 0 {
		return make([]*abci.ResponseDeliverTx, 0)
	}

	if onlyCalSender {
		app.paraLoadSender(txs)
		return nil
	}
	txWithIndex := make([][]byte, txSize)
	for index, v := range txs {
		txWithIndex[index] = getTxByteWithIndex(v, index)
	}

	app.getExtraDataByTxs(txs)

	app.calGroup()

	pm.isAsyncDeliverTx = true
	pm.cms = app.deliverState.ms.CacheMultiStore()
	pm.runBase = make([]int, txSize)

	evmIndex := uint32(0)
	for index := range txs {
		pm.extraTxsInfo[index].indexInBlock = uint32(index)
		if pm.extraTxsInfo[index].isEvm {
			pm.extraTxsInfo[index].evmIndex = evmIndex
			evmIndex++
		} else {
			pm.haveCosmosTxInBlock = true
			//fmt.Println("haveCosmosTxInBlock")
		}

		pm.indexMapBytes[string(txWithIndex[index])] = index
	}
	return app.runTxs(txWithIndex)
}

func (app *BaseApp) fixFeeCollector(index int, ms sdk.CacheMultiStore) {
	resp := app.parallelTxManage.txReps[index]
	//if resp == nil || resp.paraMsg == nil || resp.paraMsg.AnteErr != nil {
	//	continue
	//}

	if resp.paraMsg.AnteErr != nil {
		return
	}
	app.parallelTxManage.currTxFee = app.parallelTxManage.currTxFee.Add(app.parallelTxManage.extraTxsInfo[index].fee.Sub(resp.paraMsg.RefundFee)...)

	ctx, _ := app.cacheTxContext(app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})

	ctx.SetMultiStore(ms)
	if err := app.updateFeeCollectorAccHandler(ctx, app.parallelTxManage.currTxFee); err != nil {
		panic(err)
	}
}

func (app *BaseApp) runTxs(txs [][]byte) []*abci.ResponseDeliverTx {
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
	signal := make(chan int, 1)
	rerunIdx := 0
	txIndex := 0

	pm := app.parallelTxManage

	txReps := pm.txReps
	deliverTxs := make([]*abci.ResponseDeliverTx, pm.txSize)

	asyncCb := func(execRes *executeResult) {
		receiveTxIndex := int(execRes.counter)
		pm.workgroup.setTxStatus(receiveTxIndex, false)
		if receiveTxIndex < txIndex {
			return
		}
		txReps[receiveTxIndex] = execRes

		if pm.workgroup.isFailed(pm.workgroup.runningStats(receiveTxIndex)) {
			txReps[receiveTxIndex] = nil
			pm.workgroup.AddTask(txs[receiveTxIndex], receiveTxIndex)

		} else {
			if nextTx, ok := pm.nextTxInGroup[receiveTxIndex]; ok {
				if !pm.workgroup.isRunning(nextTx) {
					txReps[nextTx] = nil
					pm.workgroup.AddTask(txs[nextTx], nextTx)
				}
			}
		}

		if txIndex != receiveTxIndex {
			return
		}

		for txReps[txIndex] != nil {
			s := pm.extraTxsInfo[txIndex]
			res := txReps[txIndex]

			if pm.newIsConflict(res) || overFlow(currentGas, res.resp.GasUsed, maxGas) {
				rerunIdx++
				s.reRun = true
				res = app.deliverTxWithCache(txs[txIndex], txIndex)
				txReps[txIndex] = res

				nn, ok := app.parallelTxManage.nextTxInGroup[txIndex]

				if ok {
					if !pm.workgroup.isRunning(nn) {
						txReps[nn] = nil
						pm.workgroup.AddTask(txs[nn], nn)
					}
				}

			}
			if txReps[txIndex].paraMsg.AnteErr != nil {
				res.ms = nil
			}

			txRs := res.resp
			deliverTxs[txIndex] = &txRs

			if !s.reRun {
				app.deliverState.ctx.BlockGasMeter().ConsumeGas(sdk.Gas(res.resp.GasUsed), "unexpected error")
			}

			pm.SetCurrentIndex(txIndex, res) //Commit
			app.fixFeeCollector(txIndex, app.parallelTxManage.cms)
			currentGas += uint64(res.resp.GasUsed)
			txIndex++
			if txIndex == pm.txSize {
				app.logger.Info("Paralleled-tx", "blockHeight", app.deliverState.ctx.BlockHeight(), "len(txs)", pm.txSize,
					"Parallel run", pm.txSize-rerunIdx, "ReRun", rerunIdx, "len(group)", len(pm.groupList))
				signal <- 0
				return
			}
			if txReps[txIndex] == nil && !pm.workgroup.isRunning(txIndex) {
				pm.workgroup.AddTask(txs[txIndex], txIndex)
			}

		}
	}

	pm.workgroup.resultCb = asyncCb
	pm.workgroup.taskRun = app.asyncDeliverTx

	if len(pm.groupList) == 0 {
		pm.workgroup.AddTask(txs[0], 0)
	} else if pm.groupList[0][0] != 0 {
		pm.workgroup.AddTask(txs[0], 0)
	}

	for _, group := range pm.groupList {
		txIndex := group[0]
		pm.workgroup.AddTask(txs[txIndex], txIndex)
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

	logIndex := make([]int, app.parallelTxManage.txSize)
	errs := make([]error, app.parallelTxManage.txSize)
	for index := 0; index < app.parallelTxManage.txSize; index++ {
		paraM := app.parallelTxManage.txReps[index].paraMsg
		logIndex[index] = paraM.LogIndex
		errs[index] = paraM.AnteErr
	}
	app.parallelTxManage.clear()
	return app.logFix(logIndex, errs)
}

//we reuse the nonce that changed by the last async call
//if last ante handler has been failed, we need rerun it ? or not?
func (app *BaseApp) deliverTxWithCache(txByte []byte, txIndex int) *executeResult {
	app.parallelTxManage.workgroup.setTxStatus(txIndex, true)
	txStatus := app.parallelTxManage.extraTxsInfo[txIndex]

	if txStatus.stdTx == nil {
		asyncExe := newExecuteResult(sdkerrors.ResponseDeliverTx(txStatus.decodeErr, 0, 0, app.trace), nil, txStatus.indexInBlock, txStatus.evmIndex, nil)
		return asyncExe
	}
	var (
		resp abci.ResponseDeliverTx
		mode runTxMode
	)
	mode = runTxModeDeliverInAsync
	info, errM := app.runTx(mode, txByte, txStatus.stdTx, LatestSimulateTxHeight)
	g, r, m, e := info.gInfo, info.result, info.msCacheAnte, errM
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

	asyncExe := newExecuteResult(resp, m, txStatus.indexInBlock, txStatus.evmIndex, info.ctx.ParaMsg())
	asyncExe.err = e
	return asyncExe
}

type executeResult struct {
	resp       abci.ResponseDeliverTx
	ms         sdk.CacheMultiStore
	counter    uint32
	err        error
	evmCounter uint32
	//readList   map[string][]byte
	//writeList  map[string][]byte

	paraMsg *sdk.ParaMsg
}

func loadPreData(ms sdk.CacheMultiStore, rSet map[string][]byte, wSet map[string][]byte) {
	if ms == nil {
		return
	}

	ms.GetRWSet(rSet, wSet)
	return
}

func newExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32, evmCounter uint32, paraMsg *sdk.ParaMsg) *executeResult {
	ans := &executeResult{
		resp:       r,
		ms:         ms,
		counter:    counter,
		evmCounter: evmCounter,
		//readList:   make(map[string][]byte),
		//writeList:  make(map[string][]byte),
		paraMsg: paraMsg,
	}

	//loadPreData(ms, ans.readList, ans.writeList)
	//delete(ans.readList, whiteAcc)
	//delete(ans.writeList, whiteAcc)
	if paraMsg == nil {
		ans.paraMsg = &sdk.ParaMsg{}
	}
	return ans
}

type asyncWorkGroup struct {
	isAsync       bool
	runningStatus []int
	isrunning     []bool

	markFailedStats map[int]bool

	indexInAll int
	runningMu  sync.RWMutex

	resultCh chan *executeResult
	resultCb func(*executeResult)

	taskCh  chan *task
	taskRun func([]byte, int)
}

func newAsyncWorkGroup(isAsync bool) *asyncWorkGroup {
	return &asyncWorkGroup{
		isAsync:         isAsync,
		runningStatus:   make([]int, 0),
		isrunning:       make([]bool, 0),
		markFailedStats: make(map[int]bool),

		resultCh: make(chan *executeResult, 20000),
		resultCb: nil,

		taskCh:  make(chan *task, 20000),
		taskRun: nil,
	}
}

func (a *asyncWorkGroup) markFailed(txIndexAll int) {
	a.runningMu.Lock()
	defer a.runningMu.Unlock()
	a.markFailedStats[txIndexAll] = true
}

func (a *asyncWorkGroup) isFailed(txindexAll int) bool {
	a.runningMu.Lock()
	defer a.runningMu.Unlock()
	return a.markFailedStats[txindexAll]
}

func (a *asyncWorkGroup) setTxStatus(txIndex int, status bool) {
	a.runningMu.Lock()
	defer a.runningMu.Unlock()
	if status == true {
		a.runningStatus[txIndex] = a.indexInAll
		a.indexInAll++
	}
	a.isrunning[txIndex] = status
}

func (a *asyncWorkGroup) runningStats(txIndex int) int {
	a.runningMu.RLock()
	defer a.runningMu.RUnlock()
	return a.runningStatus[txIndex]
}

func (a *asyncWorkGroup) isRunning(txIndex int) bool {
	a.runningMu.RLock()
	defer a.runningMu.RUnlock()
	return a.isrunning[txIndex]
}

func (a *asyncWorkGroup) Push(item *executeResult) {
	a.resultCh <- item
}

func (a *asyncWorkGroup) AddTask(tx []byte, index int) {
	a.setTxStatus(index, true)
	a.taskCh <- &task{
		txBytes: tx,
		index:   index,
	}
}

func (a *asyncWorkGroup) Start() {
	if !a.isAsync {
		return
	}
	for index := 0; index < 16; index++ {
		go func() {
			for true {
				select {
				case task := <-a.taskCh:
					a.taskRun(task.txBytes, task.index)
				}
			}
		}()

	}

	go func() {
		for {
			select {
			case exec := <-a.resultCh:
				a.resultCb(exec)
			}
		}
	}()
}

type parallelTxManager struct {
	haveCosmosTxInBlock bool
	isAsyncDeliverTx    bool
	workgroup           *asyncWorkGroup

	indexMapBytes map[string]int

	extraTxsInfo []*extraDataForTx
	txReps       []*executeResult

	groupList          map[int][]int
	nextTxInGroup      map[int]int
	preTxInGroup       map[int]int
	txIndexWithGroupID map[int]int

	mu        sync.RWMutex
	cms       sdk.CacheMultiStore
	currTxFee sdk.Coins

	txSize    int
	cc        *conflictCheck
	currIndex int
	runBase   []int

	commitDone chan struct{}
}
type A struct {
	value   []byte
	txIndex int
}

type conflictCheck struct {
	items map[string]A
}

func newConflictCheck() *conflictCheck {
	return &conflictCheck{
		items: make(map[string]A),
	}
}

func (c *conflictCheck) update(key string, value []byte, txIndex int) {
	c.items[key] = A{
		value:   value,
		txIndex: txIndex,
	}
}

func (c *conflictCheck) deleteFee() {
	delete(c.items, whiteAcc)
}
func (c *conflictCheck) clear() {
	for key := range c.items {
		delete(c.items, key)
	}
}

var (
	whiteAcc = string(hexutil.MustDecode("0x01f1829676db577682e944fc3493d451b67ff3e29f")) //fee
)

func (pm *parallelTxManager) newIsConflict(e *executeResult) bool {
	//base := pm.runBase[e.counter]
	if e.ms == nil {
		return true //TODO fix later
	}
	conflict := false

	e.ms.IteratorCache(false, func(key string, value []byte, isDirty bool, isDelete bool, storeKey types.StoreKey) bool {
		if data, ok := pm.cc.items[key]; ok {
			if !bytes.Equal(data.value, value) {
				conflict = true
				return false
			}
		}
		return true
	}, nil)

	return conflict

}

func (p *parallelTxManager) isConflict(base int, key string, readValue []byte, txIndex int) bool {
	if dirtyTxIndex, ok := p.cc.items[key]; ok {
		if !bytes.Equal(dirtyTxIndex.value, readValue) {
			return true
		} else {
			if base < dirtyTxIndex.txIndex && p.txIndexWithGroupID[dirtyTxIndex.txIndex] != p.txIndexWithGroupID[txIndex] {
				return true
			}
		}
	}
	return false
}

type task struct {
	txBytes []byte
	index   int
}

func newParallelTxManager() *parallelTxManager {
	isAsync := viper.GetBool(sm.FlagParalleledTx)
	return &parallelTxManager{
		isAsyncDeliverTx: isAsync,
		workgroup:        newAsyncWorkGroup(isAsync),
		indexMapBytes:    make(map[string]int, 0),

		groupList:          make(map[int][]int),
		nextTxInGroup:      make(map[int]int),
		preTxInGroup:       make(map[int]int),
		txIndexWithGroupID: make(map[int]int),

		cc:         newConflictCheck(),
		currIndex:  -1,
		runBase:    make([]int, 0),
		commitDone: make(chan struct{}, 1),
	}
}

func (f *parallelTxManager) clear() {

	for key := range f.indexMapBytes {
		delete(f.indexMapBytes, key)
	}
	f.extraTxsInfo = make([]*extraDataForTx, 0)

	for key := range f.groupList {
		delete(f.groupList, key)
	}
	for key := range f.nextTxInGroup {
		delete(f.nextTxInGroup, key)
	}
	for key := range f.preTxInGroup {
		delete(f.preTxInGroup, key)
	}
	for key := range f.txIndexWithGroupID {
		delete(f.txIndexWithGroupID, key)
	}

	f.currIndex = -1
	f.currTxFee = sdk.Coins{}
	f.cc.clear()

	for key := range f.workgroup.markFailedStats {
		delete(f.workgroup.markFailedStats, key)
	}
	f.workgroup.indexInAll = 0
}

func (f *parallelTxManager) isReRun(tx string) bool {
	return f.extraTxsInfo[f.indexMapBytes[tx]].reRun

}

func (f *parallelTxManager) getTxResult(tx []byte) sdk.CacheMultiStore {
	index := f.indexMapBytes[string(tx)]
	preIndexInGroup, ok := f.preTxInGroup[index]
	f.mu.Lock()
	defer f.mu.Unlock()
	ms := f.cms.CacheMultiStore()
	base := f.currIndex
	if index <= base {
		return nil
	}
	if ok && preIndexInGroup > f.currIndex {
		if f.txReps[preIndexInGroup].paraMsg.AnteErr == nil {
			ms = f.txReps[preIndexInGroup].ms.CacheMultiStore()
		} else {
			ms = f.cms.CacheMultiStore()
		}

	}

	if next, ok := f.nextTxInGroup[index]; ok {
		if f.workgroup.isRunning(next) {
			f.workgroup.markFailed(f.workgroup.runningStats(next))
		} else {
			f.txReps[next] = nil
		}
	}
	f.runBase[index] = base

	return ms
}

func (f *parallelTxManager) getRunBase(now int) int {
	return f.runBase[now]
}

func (f *parallelTxManager) SetCurrentIndex(txIndex int, res *executeResult) {
	if res.ms == nil {
		return
	}

	f.mu.Lock()
	res.ms.IteratorCache(true, func(key string, value []byte, isDirty bool, isdelete bool, storeKey sdk.StoreKey) bool {
		if isDirty {
			f.cc.update(key, value, txIndex)
			if isdelete {
				f.cms.GetKVStore(storeKey).Delete([]byte(key))
			} else if value != nil {
				f.cms.GetKVStore(storeKey).Set([]byte(key), value)
			}
		}
		return true
	}, nil)
	f.currIndex = txIndex
	f.mu.Unlock()
	f.cc.deleteFee()

}

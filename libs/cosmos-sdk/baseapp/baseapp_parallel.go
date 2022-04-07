package baseapp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	fee   sdk.Coins
	isEvm bool
	from  string
	to    *ethcommon.Address
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
				fee:   coin,
				isEvm: isEvm,
				from:  s,
				to:    toAddr,
			}
		}()
	}
	wg.Wait()
	return res
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

func (app *BaseApp) calGroup(txsExtraData []*extraDataForTx) (map[int][]int, map[int]int) {
	rootAddr = make(map[string]string, 0)
	app.parallelTxManage.txReps = make([]*executeResult, len(txsExtraData))
	for index, tx := range txsExtraData {
		if tx.isEvm { //evmTx
			Union(tx.from, tx.to)
		} else {
			app.parallelTxManage.txReps[index] = &executeResult{}
		}
	}

	groupList := make(map[int][]int, 0)
	addrToID := make(map[string]int, 0)

	for index, sender := range txsExtraData {
		if !sender.isEvm {
			continue
		}
		rootAddr := Find(sender.from)
		id, exist := addrToID[rootAddr]
		if !exist {
			id = len(groupList)
			addrToID[rootAddr] = id

		}
		groupList[id] = append(groupList[id], index)
	}

	nextTxIndexInGroup := make(map[int]int)
	preTxIndexInGroup := make(map[int]int)
	txIndexWithGroupID := make(map[int]int)
	groupSize := len(groupList)
	for groupIndex := 0; groupIndex < groupSize; groupIndex++ {
		list := groupList[groupIndex]
		for index := 0; index < len(list); index++ {
			if index+1 <= len(list)-1 {
				nextTxIndexInGroup[list[index]] = list[index+1]
			}
			if index-1 >= 0 {
				preTxIndexInGroup[list[index]] = list[index-1]
			}
			txIndexWithGroupID[list[index]] = groupIndex
		}
	}
	app.parallelTxManage.nextTxInGroup = nextTxIndexInGroup
	app.parallelTxManage.preTxInGroup = preTxIndexInGroup
	app.parallelTxManage.txIndexWithGroupID = txIndexWithGroupID
	return groupList, nextTxIndexInGroup
}

func (app *BaseApp) paraLoadSender(txs [][]byte) {

	checkStateCtx := app.checkState.ctx.WithBlockHeight(app.checkState.ctx.BlockHeight() + 1)

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
					app.getTxFee(checkStateCtx.WithTxBytes(txBytes), tx)
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
	//sdk.DebugLogByScf.Clean()

	if len(txs) == 0 {
		return make([]*abci.ResponseDeliverTx, 0)
	}

	if onlyCalSender {
		app.paraLoadSender(txs)
		return nil
	}
	txWithIndex := make([][]byte, 0)
	for index, v := range txs {
		txWithIndex = append(txWithIndex, getTxByteWithIndex(v, index))
	}

	extraData := app.getExtraDataByTxs(txs)

	groupList, nextIndexInGroup := app.calGroup(extraData)

	app.parallelTxManage.isAsyncDeliverTx = true
	app.parallelTxManage.cms = app.deliverState.ms.CacheMultiStore()
	app.parallelTxManage.runBase = make([]int, len(txs))

	evmIndex := uint32(0)
	for k := range txs {
		t := &txStatus{
			indexInBlock: uint32(k),
		}
		if extraData[k].isEvm {
			t.evmIndex = evmIndex
			t.isEvmTx = true
			evmIndex++
		}

		vString := string(txWithIndex[k])
		app.parallelTxManage.fee[vString] = extraData[k].fee

		app.parallelTxManage.txStatus[vString] = t
		app.parallelTxManage.indexMapBytes = append(app.parallelTxManage.indexMapBytes, vString)
	}
	return app.runTxs(txWithIndex, groupList, nextIndexInGroup)
}

func (app *BaseApp) fixFeeCollector(txs [][]byte, ms sdk.CacheMultiStore) {
	currTxFee := sdk.Coins{}
	for index, v := range txs {
		txString := string(v)
		if app.parallelTxManage.txReps[index].paraMsg.AnteErr != nil {
			continue
		}
		txFee := app.parallelTxManage.fee[txString]
		refundFee := app.parallelTxManage.txReps[index].paraMsg.RefundFee
		txFee = txFee.Sub(refundFee)
		currTxFee = currTxFee.Add(txFee...)
	}

	ctx, _ := app.cacheTxContext(app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})

	ctx = ctx.WithMultiStore(ms)
	if err := app.updateFeeCollectorAccHandler(ctx, currTxFee); err != nil {
		panic(err)
	}
}

func (app *BaseApp) runTxs(txs [][]byte, groupList map[int][]int, nextTxInGroup map[int]int) []*abci.ResponseDeliverTx {
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
	deliverTxs := make([]*abci.ResponseDeliverTx, len(txs))

	asyncCb := func(execRes *executeResult) {
		receiveTxIndex := int(execRes.GetCounter())
		pm.workgroup.setTxStatus(receiveTxIndex, false)
		if receiveTxIndex < txIndex {
			return
		}
		txReps[receiveTxIndex] = execRes

		if pm.workgroup.isFailed(pm.workgroup.runningStats(receiveTxIndex)) {
			txReps[receiveTxIndex] = nil
			pm.workgroup.AddTask(txs[receiveTxIndex], receiveTxIndex)

		} else {
			if nextTx, ok := nextTxInGroup[receiveTxIndex]; ok {
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
			txBytes := app.parallelTxManage.indexMapBytes[txIndex]
			s := pm.txStatus[txBytes]
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

			txRs := res.GetResponse()
			deliverTxs[txIndex] = &txRs

			if !s.reRun {
				app.deliverState.ctx.BlockGasMeter().ConsumeGas(sdk.Gas(res.resp.GasUsed), "unexpected error")
			}

			pm.SetCurrentIndex(txIndex, res) //Commit
			currentGas += uint64(res.resp.GasUsed)
			txIndex++
			if txIndex == len(txs) {
				app.logger.Info("Paralleled-tx", "blockHeight", app.deliverState.ctx.BlockHeight(), "len(txs)", len(txs),
					"Parallel run", len(txs)-rerunIdx, "ReRun", rerunIdx, "len(group)", len(groupList))
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

	if len(groupList) == 0 {
		pm.workgroup.AddTask(txs[0], 0)
	} else if groupList[0][0] != 0 {
		pm.workgroup.AddTask(txs[0], 0)
	}

	for _, group := range groupList {
		txIndex := group[0]
		pm.workgroup.AddTask(txs[txIndex], txIndex)
	}

	if len(txs) > 0 {
		//waiting for call back
		<-signal
		app.fixFeeCollector(txs, pm.cms)
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

	logIndex := make([]int, 0)
	errs := make([]error, 0)
	for txIndex, _ := range app.parallelTxManage.indexMapBytes {
		paraM := app.parallelTxManage.txReps[txIndex].paraMsg
		logIndex = append(logIndex, paraM.LogIndex)
		errs = append(errs, paraM.AnteErr)
	}
	app.parallelTxManage.clear()
	return app.logFix(logIndex, errs)
}

//we reuse the nonce that changed by the last async call
//if last ante handler has been failed, we need rerun it ? or not?
func (app *BaseApp) deliverTxWithCache(txByte []byte, txIndex int) *executeResult {
	app.parallelTxManage.workgroup.setTxStatus(txIndex, true)
	txStatus := app.parallelTxManage.txStatus[string(txByte)]

	tx, err := app.txDecoder(getRealTxByte(txByte))
	if err != nil {
		asyncExe := newExecuteResult(sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace), nil, txStatus.indexInBlock, txStatus.evmIndex, nil)
		return asyncExe
	}
	var (
		resp abci.ResponseDeliverTx
		mode runTxMode
	)
	mode = runTxModeDeliverInAsync
	info, errM := app.runTx(mode, txByte, tx, LatestSimulateTxHeight)
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
	readList   map[string][]byte
	writeList  map[string][]byte

	paraMsg *sdk.ParaMsg
}

func (e executeResult) GetResponse() abci.ResponseDeliverTx {
	return e.resp
}

func (e executeResult) GetCounter() uint32 {
	return e.counter
}

func loadPreData(ms sdk.CacheMultiStore) (map[string][]byte, map[string][]byte) {
	if ms == nil {
		return nil, nil
	}

	rSet := make(map[string][]byte)
	wSet := make(map[string][]byte)

	ms.GetRWSet(rSet, wSet)
	return rSet, wSet
}

func newExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32, evmCounter uint32, paraMsg *sdk.ParaMsg) *executeResult {
	rSet, wSet := loadPreData(ms)
	delete(rSet, whiteAcc)
	delete(wSet, whiteAcc)
	if paraMsg == nil {
		paraMsg = &sdk.ParaMsg{}
	}

	return &executeResult{
		resp:       r,
		ms:         ms,
		counter:    counter,
		evmCounter: evmCounter,
		readList:   rSet,
		writeList:  wSet,
		paraMsg:    paraMsg,
	}
}

type asyncWorkGroup struct {
	isAsync       bool
	runningStatus map[int]int
	isrunning     map[int]bool

	markFailedStats map[int]bool

	indexInAll int
	runningMu  sync.RWMutex

	resultCh chan *executeResult
	resultCb func(*executeResult)

	taskCh  chan *task
	taskRun func([]byte)
}

func newAsyncWorkGroup(isAsync bool) *asyncWorkGroup {
	return &asyncWorkGroup{
		isAsync:         isAsync,
		runningStatus:   make(map[int]int),
		isrunning:       make(map[int]bool),
		markFailedStats: make(map[int]bool),

		resultCh: make(chan *executeResult, 100000),
		resultCb: nil,

		taskCh:  make(chan *task, 100000),
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
	for index := 0; index < 64; index++ {
		go func() {
			for true {
				select {
				case task := <-a.taskCh:
					a.taskRun(task.txBytes)
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
	isAsyncDeliverTx bool
	workgroup        *asyncWorkGroup

	fee map[string]sdk.Coins // not need mute

	txStatus      map[string]*txStatus
	indexMapBytes []string

	txReps             []*executeResult
	nextTxInGroup      map[int]int
	preTxInGroup       map[int]int
	txIndexWithGroupID map[int]int

	mu  sync.RWMutex
	cms sdk.CacheMultiStore

	cc         *conflictCheck
	currIndex  int
	runBase    []int
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
func (c *conflictCheck) clear() {
	c.items = make(map[string]A, 0)
}

var (
	whiteAcc = string(hexutil.MustDecode("0x01f1829676db577682e944fc3493d451b67ff3e29f")) //fee
)

func (pm *parallelTxManager) newIsConflict(e *executeResult) bool {
	base := pm.runBase[e.counter]
	if e.ms == nil {
		return true //TODO fix later
	}
	for k, readValue := range e.readList {

		if pm.isConflict(base, k, readValue, int(e.counter)) {
			return true
		}
	}
	return false

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

type txStatus struct {
	reRun        bool
	isEvmTx      bool
	evmIndex     uint32
	indexInBlock uint32
}

func newParallelTxManager() *parallelTxManager {
	isAsync := viper.GetBool(sm.FlagParalleledTx)
	return &parallelTxManager{
		isAsyncDeliverTx: isAsync,
		workgroup:        newAsyncWorkGroup(isAsync),
		fee:              make(map[string]sdk.Coins),

		txStatus:      make(map[string]*txStatus),
		indexMapBytes: make([]string, 0),

		nextTxInGroup:      make(map[int]int),
		preTxInGroup:       make(map[int]int),
		txIndexWithGroupID: make(map[int]int),

		cc:        newConflictCheck(),
		currIndex: -1,
		runBase:   make([]int, 0),

		commitDone: make(chan struct{}, 1),
	}
}

func (f *parallelTxManager) clear() {
	f.fee = make(map[string]sdk.Coins)

	f.txStatus = make(map[string]*txStatus)
	f.indexMapBytes = make([]string, 0)
	f.nextTxInGroup = make(map[int]int)
	f.preTxInGroup = make(map[int]int)
	f.txIndexWithGroupID = make(map[int]int)

	f.currIndex = -1
	f.cc.clear()
	f.workgroup.markFailedStats = make(map[int]bool)

	f.workgroup.runningStatus = make(map[int]int)
	f.workgroup.isrunning = make(map[int]bool)
	f.workgroup.indexInAll = 0
}

func (f *parallelTxManager) isReRun(tx string) bool {
	data, ok := f.txStatus[tx]
	if !ok {
		return false
	}
	return data.reRun
}

func (f *parallelTxManager) getTxResult(tx []byte) sdk.CacheMultiStore {
	index := int(f.txStatus[string(tx)].indexInBlock)
	preIndexInGroup, ok := f.preTxInGroup[index]
	f.mu.Lock()
	defer f.mu.Unlock()
	ms := f.cms.CacheMultiStore()
	base := f.currIndex
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

	chanStop := make(chan struct{}, 2)
	go func() {
		for k, v := range res.writeList {
			f.cc.update(k, v, txIndex)
		}
		chanStop <- struct{}{}
	}()

	go func() {
		tt := make([]string, 0)
		for k, v := range res.readList {
			if len(v) < 200 {
				tt = append(tt, fmt.Sprintf("read key:%s value:%s", hex.EncodeToString([]byte(k)), hex.EncodeToString(v)))
			}
		}

		for k, v := range res.writeList {
			tt = append(tt, fmt.Sprintf("write ket:%s value:%s", hex.EncodeToString([]byte(k)), hex.EncodeToString(v)))
		}
		tt = append(tt, fmt.Sprintf("txIndex %d base %d", txIndex, f.getRunBase(txIndex)))
		sdk.DebugLogByScf.AddRWSet(tt)
		chanStop <- struct{}{}
	}()

	f.mu.Lock()
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
	f.currIndex = txIndex
	f.mu.Unlock()
	<-chanStop
	<-chanStop
}

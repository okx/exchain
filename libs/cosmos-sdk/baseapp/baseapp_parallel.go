package baseapp

import (
	"bytes"
	"runtime"
	"sync"

	"github.com/spf13/viper"

	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	sm "github.com/okx/okbchain/libs/tendermint/state"
)

var (
	maxTxResultInChan           = 20000
	maxGoroutineNumberInParaTx  = runtime.NumCPU()
	multiCacheListClearInterval = int64(100)
)

type extraDataForTx struct {
	fee       sdk.Coins
	isEvm     bool
	from      string
	to        string
	stdTx     sdk.Tx
	decodeErr error
}

type txWithIndex struct {
	index   int
	txBytes []byte
}

// getExtraDataByTxs preprocessing tx : verify tx, get sender, get toAddress, get txFee
func (app *BaseApp) getExtraDataByTxs(txs [][]byte) {
	para := app.parallelTxManage

	var wg sync.WaitGroup
	wg.Add(len(txs))
	jobChan := make(chan txWithIndex, len(txs))
	for groupIndex := 0; groupIndex < maxGoroutineNumberInParaTx; groupIndex++ {
		go func(ch chan txWithIndex) {
			for j := range ch {
				index := j.index
				txBytes := j.txBytes
				var tx sdk.Tx
				var err error

				if mem := GetGlobalMempool(); mem != nil {
					tx, _ = mem.ReapEssentialTx(txBytes).(sdk.Tx)
				}
				if tx == nil {
					tx, err = app.txDecoder(txBytes)
					if err != nil {
						para.extraTxsInfo[index] = &extraDataForTx{
							decodeErr: err,
						}
						wg.Done()
						continue
					}
				}
				if tx != nil {
					app.blockDataCache.SetTx(txBytes, tx)
				}

				coin, isEvm, s, toAddr, _ := app.getTxFeeAndFromHandler(app.getContextForTx(runTxModeDeliver, txBytes), tx)
				para.extraTxsInfo[index] = &extraDataForTx{
					fee:   coin,
					isEvm: isEvm,
					from:  s,
					to:    toAddr,
					stdTx: tx,
				}
				wg.Done()
			}
		}(jobChan)
	}

	for index, v := range txs {
		jobChan <- txWithIndex{
			index:   index,
			txBytes: v,
		}
	}
	close(jobChan)
	wg.Wait()
}

var (
	rootAddr = make(map[string]string, 0)
)

// Find father node
func Find(x string) string {
	if rootAddr[x] != x {
		rootAddr[x] = Find(rootAddr[x])
	}
	return rootAddr[x]
}

// Union from and to
func Union(x string, yString string) {
	if _, ok := rootAddr[x]; !ok {
		rootAddr[x] = x
	}
	if yString == "" {
		return
	}
	if _, ok := rootAddr[yString]; !ok {
		rootAddr[yString] = yString
	}
	fx := Find(x)
	fy := Find(yString)
	if fx != fy {
		rootAddr[fy] = fx
	}
}

// calGroup cal group by txs
func (app *BaseApp) calGroup() {

	para := app.parallelTxManage

	rootAddr = make(map[string]string, 0)
	for index, tx := range para.extraTxsInfo {
		if tx.isEvm { //evmTx
			Union(tx.from, tx.to)
		} else {
			para.haveCosmosTxInBlock = true
			app.parallelTxManage.putResult(index, &executeResult{paraMsg: &sdk.ParaMsg{}, msIsNil: true})
		}
	}

	addrToID := make(map[string]int, 0)

	for index, txInfo := range para.extraTxsInfo {
		if !txInfo.isEvm {
			continue
		}
		rootAddr := Find(txInfo.from)
		id, exist := addrToID[rootAddr]
		if !exist {
			id = len(para.groupList)
			addrToID[rootAddr] = id

		}
		para.groupList[id] = append(para.groupList[id], index)
		para.txIndexWithGroup[index] = id
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
		}
	}
}

// ParallelTxs run txs
func (app *BaseApp) ParallelTxs(txs [][]byte, onlyCalSender bool) []*abci.ResponseDeliverTx {
	txSize := len(txs)

	if txSize == 0 {
		return make([]*abci.ResponseDeliverTx, 0)
	}

	pm := app.parallelTxManage
	pm.init(txs, app.deliverState.ctx.BlockHeight(), app.deliverState.ms)

	app.getExtraDataByTxs(txs)

	app.calGroup()

	return app.runTxs()
}

func (app *BaseApp) fixFeeCollector() {
	ctx, _ := app.cacheTxContext(app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})

	ctx.SetMultiStore(app.parallelTxManage.cms)
	// The feesplit is only processed at the endblock
	if err := app.updateFeeCollectorAccHandler(ctx, app.parallelTxManage.currTxFee, nil); err != nil {
		panic(err)
	}
}

func (app *BaseApp) runTxs() []*abci.ResponseDeliverTx {
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

	pm := app.parallelTxManage

	asyncCb := func(receiveTxIndex int) {
		if pm.alreadyEnd {
			return
		}
		//skip old txIndex
		if receiveTxIndex < pm.upComingTxIndex || receiveTxIndex >= pm.txSize {
			return
		}

		for true {
			res := pm.getTxResult(pm.upComingTxIndex)
			if res == nil {
				break
			}
			isReRun := false
			if pm.isConflict(res) || overFlow(currentGas, res.resp.GasUsed, maxGas) {
				rerunIdx++
				isReRun = true
				// conflict rerun tx
				if !pm.extraTxsInfo[pm.upComingTxIndex].isEvm {
					app.fixFeeCollector()
				}
				res = app.deliverTxWithCache(pm.upComingTxIndex)
			}
			if res.paraMsg.AnteErr != nil {
				res.msIsNil = true
			}

			pm.deliverTxs[pm.upComingTxIndex] = &res.resp
			pm.finalResult[pm.upComingTxIndex] = res

			pm.blockGasMeterMu.Lock()
			// Note : don't take care of the case of ErrorGasOverflow
			app.deliverState.ctx.BlockGasMeter().ConsumeGas(sdk.Gas(res.resp.GasUsed), "unexpected error")
			pm.blockGasMeterMu.Unlock()

			pm.SetCurrentIndex(pm.upComingTxIndex, res)
			currentGas += uint64(res.resp.GasUsed)

			if isReRun {
				if pm.nextTxInGroup[pm.upComingTxIndex] != 0 {
					pm.groupTasks[pm.txIndexWithGroup[pm.upComingTxIndex]].addRerun(pm.upComingTxIndex)
				}
			}
			pm.upComingTxIndex++

			if pm.upComingTxIndex == pm.txSize {
				app.logger.Info("Paralleled-tx", "blockHeight", app.deliverState.ctx.BlockHeight(), "len(txs)", pm.txSize,
					"Parallel run", pm.txSize-rerunIdx, "ReRun", rerunIdx, "len(group)", len(pm.groupList))
				signal <- 0
				return
			}
		}
	}

	pm.resultCb = asyncCb
	pm.StartResultHandle()
	for index := 0; index < len(pm.groupList); index++ {
		pm.groupTasks = append(pm.groupTasks, newGroupTask(len(pm.groupList[index]), pm.addMultiCache, pm.nextTxInThisGroup, app.asyncDeliverTx, pm.putResult))
		pm.groupTasks[index].addTask(pm.groupList[index][0])
	}
	if len(pm.groupList) == 0 {
		pm.resultCh <- 0
	}

	//waiting for call back
	<-signal

	for _, v := range pm.groupTasks {
		v.stopChan <- struct{}{}
	}
	pm.alreadyEnd = true
	pm.stop <- struct{}{}

	// fix logs
	app.feeChanged = true
	app.feeCollector = app.parallelTxManage.currTxFee
	receiptsLogs := app.endParallelTxs(pm.txSize)
	for index, v := range receiptsLogs {
		if len(v) != 0 { // only update evm tx result
			pm.deliverTxs[index].Data = v
		}
	}

	pm.cms.Write()
	return pm.deliverTxs
}

func (pm *parallelTxManager) nextTxInThisGroup(txindex int) (int, bool) {
	if pm.alreadyEnd {
		return 0, false
	}
	data, ok := pm.nextTxInGroup[txindex]
	return data, ok
}

func (app *BaseApp) endParallelTxs(txSize int) [][]byte {

	// handle receipt's logs
	logIndex := make([]int, txSize)
	errs := make([]error, txSize)
	hasEnterEvmTx := make([]bool, txSize)
	resp := make([]abci.ResponseDeliverTx, txSize)
	watchers := make([]sdk.IWatcher, txSize)
	txs := make([]sdk.Tx, txSize)
	app.FeeSplitCollector = make([]*sdk.FeeSplitInfo, 0)
	for index := 0; index < txSize; index++ {
		txRes := app.parallelTxManage.finalResult[index]
		logIndex[index] = txRes.paraMsg.LogIndex
		errs[index] = txRes.paraMsg.AnteErr
		hasEnterEvmTx[index] = txRes.paraMsg.HasRunEvmTx
		resp[index] = txRes.resp
		watchers[index] = txRes.watcher
		txs[index] = app.parallelTxManage.extraTxsInfo[index].stdTx
		if txRes.FeeSpiltInfo.HasFee {
			app.FeeSplitCollector = append(app.FeeSplitCollector, txRes.FeeSpiltInfo)
		}
	}
	app.watcherCollector(watchers...)
	app.parallelTxManage.clear()

	return app.logFix(txs, logIndex, hasEnterEvmTx, errs, resp)
}

//we reuse the nonce that changed by the last async call
//if last ante handler has been failed, we need rerun it ? or not?
func (app *BaseApp) deliverTxWithCache(txIndex int) *executeResult {
	app.parallelTxManage.currentRerunIndex = txIndex
	defer func() {
		app.parallelTxManage.currentRerunIndex = -1
	}()
	txStatus := app.parallelTxManage.extraTxsInfo[txIndex]

	if txStatus.stdTx == nil {
		asyncExe := newExecuteResult(sdkerrors.ResponseDeliverTx(txStatus.decodeErr,
			0, 0, app.trace), nil, uint32(txIndex), nil, 0, sdk.EmptyWatcher{}, nil, app.parallelTxManage, nil)
		return asyncExe
	}
	var (
		resp abci.ResponseDeliverTx
		mode runTxMode
	)
	mode = runTxModeDeliverInAsync
	info, errM := app.runTxWithIndex(txIndex, mode, app.parallelTxManage.txs[txIndex], txStatus.stdTx, LatestSimulateTxHeight)
	if errM != nil {
		resp = sdkerrors.ResponseDeliverTx(errM, info.gInfo.GasWanted, info.gInfo.GasUsed, app.trace)
	} else {
		resp = abci.ResponseDeliverTx{
			GasWanted: int64(info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
			GasUsed:   int64(info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
			Log:       info.result.Log,
			Data:      info.result.Data,
			Events:    info.result.Events.ToABCIEvents(),
		}
	}

	asyncExe := newExecuteResult(resp, info.msCacheAnte, uint32(txIndex), info.ctx.ParaMsg(),
		0, info.runMsgCtx.GetWatcher(), info.tx.GetMsgs(), app.parallelTxManage, info.ctx.GetFeeSplitInfo())
	app.parallelTxManage.addMultiCache(info.msCacheAnte, info.msCache)
	return asyncExe
}

type executeResult struct {
	resp         abci.ResponseDeliverTx
	ms           sdk.CacheMultiStore
	msIsNil      bool // TODO delete it
	counter      uint32
	paraMsg      *sdk.ParaMsg
	blockHeight  int64
	watcher      sdk.IWatcher
	msgs         []sdk.Msg
	FeeSpiltInfo *sdk.FeeSplitInfo

	rwSet types.MsRWSet
}

func newExecuteResult(r abci.ResponseDeliverTx, ms sdk.CacheMultiStore, counter uint32,
	paraMsg *sdk.ParaMsg, height int64, watcher sdk.IWatcher, msgs []sdk.Msg, para *parallelTxManager, feeSpiltInfo *sdk.FeeSplitInfo) *executeResult {

	rwSet := para.chainMpCache.GetRWSet()
	if ms != nil {
		ms.GetRWSet(rwSet)
	}
	para.blockMpCache.PutRwSet(rwSet)

	if feeSpiltInfo == nil {
		feeSpiltInfo = &sdk.FeeSplitInfo{}
	}
	ans := &executeResult{
		resp:         r,
		ms:           ms,
		msIsNil:      ms == nil,
		counter:      counter,
		paraMsg:      paraMsg,
		blockHeight:  height,
		watcher:      watcher,
		msgs:         msgs,
		rwSet:        rwSet,
		FeeSpiltInfo: feeSpiltInfo,
	}

	if paraMsg == nil {
		ans.paraMsg = &sdk.ParaMsg{}
	}
	return ans
}

type parallelTxManager struct {
	blockHeight         int64
	groupTasks          []*groupTask
	blockGasMeterMu     sync.Mutex
	haveCosmosTxInBlock bool
	isAsyncDeliverTx    bool
	txs                 [][]byte
	txSize              int
	alreadyEnd          bool

	resultCh chan int
	resultCb func(data int)
	stop     chan struct{}

	groupList        map[int][]int
	nextTxInGroup    map[int]int
	preTxInGroup     map[int]int
	txIndexWithGroup map[int]int

	currentRerunIndex int
	upComingTxIndex   int
	currTxFee         sdk.Coins
	cms               sdk.CacheMultiStore
	conflictCheck     types.MsRWSet

	blockMpCache     *cacheRWSetList
	chainMpCache     *cacheRWSetList
	blockMultiStores *cacheMultiStoreList
	chainMultiStores *cacheMultiStoreList

	extraTxsInfo []*extraDataForTx
	txReps       []*executeResult
	finalResult  []*executeResult
	deliverTxs   []*abci.ResponseDeliverTx
}

func (pm *parallelTxManager) putResult(txIndex int, res *executeResult) {
	if pm.alreadyEnd {
		return
	}

	pm.txReps[txIndex] = res
	if res != nil {
		pm.resultCh <- txIndex
	}
}

func (pm *parallelTxManager) getTxResult(txIndex int) *executeResult {
	if pm.alreadyEnd {
		return nil
	}
	return pm.txReps[txIndex]
}

func (pm *parallelTxManager) StartResultHandle() {
	go func() {
		for {
			select {
			case exec := <-pm.resultCh:
				pm.resultCb(exec)

			case <-pm.stop:
				return
			}
		}
	}()
}

type groupTask struct {
	addMultiCache func(msAnte types.CacheMultiStore, msCache types.CacheMultiStore)
	mu            sync.Mutex
	groupIndex    map[int]int

	nextTx    func(int) (int, bool)
	taskRun   func(int) *executeResult
	putResult func(index int, txResult *executeResult)

	txChan    chan int
	reRunChan chan int
	stopChan  chan struct{}
	ms        sdk.CacheMultiStore
}

func newGroupTask(txSizeInGroup int, addMultiCache func(msAnte types.CacheMultiStore, msCache types.CacheMultiStore), nextTx func(int2 int) (int, bool), task func(int2 int) *executeResult, putResult func(index int, txResult *executeResult)) *groupTask {
	g := &groupTask{
		addMultiCache: addMultiCache,
		mu:            sync.Mutex{},
		nextTx:        nextTx,
		taskRun:       task,
		txChan:        make(chan int, txSizeInGroup),
		reRunChan:     make(chan int, 1000),
		stopChan:      make(chan struct{}, 1),
		putResult:     putResult,
	}
	go g.run()
	return g
}

func (g *groupTask) addTask(txIndex int) {
	g.txChan <- txIndex
}

func (g *groupTask) addRerun(txIndex int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.clearResultChan(txIndex)
	g.reRunChan <- txIndex
}

func (g *groupTask) clearResultChan(rerunIndex int) {
	for true {
		next, ok := g.nextTx(rerunIndex) // TODO add currIndex
		if ok {
			g.putResult(next, nil)
		} else {
			return
		}
		rerunIndex = next
	}
}

func (g *groupTask) run() {

	for true {
		select {
		case txIndex := <-g.txChan:
			g.mu.Lock()
			res := g.taskRun(txIndex)
			if res.paraMsg.UseCurrentState {
				g.addMultiCache(g.ms, nil)
				g.ms = res.ms.CacheMultiStore()
			} else {
				if res.ms != nil {
					res.ms.Write()
				}
			}

			if len(g.reRunChan) == 0 {
				g.putResult(int(res.counter), res)
				if n, ok := g.nextTx(txIndex); ok {
					g.addTask(n)
				}
			}
			g.mu.Unlock()

		case rerunIndex := <-g.reRunChan:
			g.clearResultChan(rerunIndex)
			g.addMultiCache(g.ms, nil)
			g.ms = nil
			size := len(g.txChan)
			for index := 0; index < size; index++ {
				<-g.txChan
			}

			if n, ok := g.nextTx(rerunIndex); ok {
				g.addTask(n)
			}
		case <-g.stopChan:
			return
		}
	}
}

func newParallelTxManager() *parallelTxManager {
	isAsync := sm.DeliverTxsExecMode(viper.GetInt(sm.FlagDeliverTxsExecMode)) == sm.DeliverTxsExecModeParallel
	para := &parallelTxManager{
		blockGasMeterMu:  sync.Mutex{},
		isAsyncDeliverTx: isAsync,
		stop:             make(chan struct{}, 1),

		conflictCheck: make(types.MsRWSet),

		groupList:        make(map[int][]int),
		nextTxInGroup:    make(map[int]int),
		preTxInGroup:     make(map[int]int),
		txIndexWithGroup: make(map[int]int),
		resultCh:         make(chan int, maxTxResultInChan),

		blockMpCache:     newCacheRWSetList(),
		chainMpCache:     newCacheRWSetList(),
		blockMultiStores: newCacheMultiStoreList(),
		chainMultiStores: newCacheMultiStoreList(),
	}
	return para
}

func (pm *parallelTxManager) addMultiCache(ms1 types.CacheMultiStore, ms2 types.CacheMultiStore) {
	if ms1 != nil {
		pm.blockMultiStores.PushStore(ms1)
	}

	if ms2 != nil {
		pm.blockMultiStores.PushStore(ms2)
	}
}

func shouldCleanChainCache(height int64) bool {
	return height%multiCacheListClearInterval == 0
}

func (pm *parallelTxManager) addMpCacheToChainCache() {
	if shouldCleanChainCache(pm.blockHeight) {
		pm.chainMpCache.Clear()
	} else {
		jobChan := make(chan types.MsRWSet, pm.blockMpCache.Len())
		go func() {
			for index := 0; index < maxGoroutineNumberInParaTx; index++ {
				go func(ch chan types.MsRWSet) {
					for j := range ch {
						types.ClearMsRWSet(j)
						pm.chainMpCache.PutRwSet(j)
					}
				}(jobChan)
			}

		}()

		pm.blockMpCache.Range(func(c types.MsRWSet) {
			jobChan <- c
		})
		close(jobChan)
	}
	pm.blockMpCache.Clear()

}

func (pm *parallelTxManager) addBlockCacheToChainCache() {
	if shouldCleanChainCache(pm.blockHeight) {
		pm.chainMultiStores.Clear()
	} else {
		jobChan := make(chan types.CacheMultiStore, pm.blockMultiStores.Len())
		go func() {
			for index := 0; index < maxGoroutineNumberInParaTx; index++ {
				go func(ch chan types.CacheMultiStore) {
					for j := range ch {
						j.Clear()
						pm.chainMultiStores.PushStore(j)
					}
				}(jobChan)
			}
		}()

		pm.blockMultiStores.Range(func(c types.CacheMultiStore) {
			jobChan <- c
		})
		close(jobChan)
	}
	pm.blockMultiStores.Clear()
}

func (pm *parallelTxManager) isConflict(e *executeResult) bool {
	if e.msIsNil {
		return true //TODO fix later
	}
	for storeKey, rw := range e.rwSet {

		for key, value := range rw.Read {
			if data, ok := pm.conflictCheck[storeKey].Write[key]; ok {
				if !bytes.Equal(data.Value, value) {
					return true
				}
			}
		}
	}
	return false
}

func (pm *parallelTxManager) clear() {

	pm.addBlockCacheToChainCache()
	pm.addMpCacheToChainCache()

	for key := range pm.groupList {
		delete(pm.groupList, key)
	}
	for key := range pm.preTxInGroup {
		delete(pm.preTxInGroup, key)
	}
	for key := range pm.txIndexWithGroup {
		delete(pm.txIndexWithGroup, key)
	}

	for _, v := range pm.conflictCheck {
		for k := range v.Read {
			delete(v.Read, k)
		}
		for k := range v.Write {
			delete(v.Write, k)
		}
	}
}

func (pm *parallelTxManager) init(txs [][]byte, blockHeight int64, deliverStateMs sdk.CacheMultiStore) {

	txSize := len(txs)
	pm.blockHeight = blockHeight
	pm.groupTasks = make([]*groupTask, 0)
	pm.haveCosmosTxInBlock = false
	pm.isAsyncDeliverTx = true
	pm.txs = txs
	pm.txSize = txSize
	pm.alreadyEnd = false

	pm.currentRerunIndex = -1
	pm.upComingTxIndex = 0
	pm.currTxFee = sdk.Coins{}
	pm.cms = deliverStateMs.CacheMultiStore()
	pm.cms.DisableCacheReadList()
	deliverStateMs.DisableCacheReadList()

	if txSize > cap(pm.resultCh) {
		pm.resultCh = make(chan int, txSize)
	}

	pm.nextTxInGroup = make(map[int]int)

	pm.extraTxsInfo = make([]*extraDataForTx, txSize)
	pm.txReps = make([]*executeResult, txSize)
	pm.finalResult = make([]*executeResult, txSize)
	pm.deliverTxs = make([]*abci.ResponseDeliverTx, txSize)
}

func (pm *parallelTxManager) getParentMsByTxIndex(txIndex int) (sdk.CacheMultiStore, bool) {

	if txIndex <= pm.upComingTxIndex-1 {
		return nil, false
	}

	useCurrent := false
	var ms types.CacheMultiStore
	if pm.currentRerunIndex != txIndex && pm.preTxInGroup[txIndex] > pm.upComingTxIndex-1 {
		if groupMs := pm.groupTasks[pm.txIndexWithGroup[txIndex]].ms; groupMs != nil {
			ms = pm.chainMultiStores.GetStoreWithParent(groupMs)
		}
	}

	if ms == nil {
		useCurrent = true
		ms = pm.chainMultiStores.GetStoreWithParent(pm.cms)
	}
	return ms, useCurrent
}

func (pm *parallelTxManager) SetCurrentIndex(txIndex int, res *executeResult) {
	if res.msIsNil {
		return
	}

	for storeKey, rw := range res.rwSet {
		if _, ok := pm.conflictCheck[storeKey]; !ok {
			pm.conflictCheck[storeKey] = types.NewCacheKvRWSet()
		}

		ms := pm.cms.GetKVStore(storeKey)
		for key, value := range rw.Write {
			if value.Deleted {
				ms.Delete([]byte(key))
			} else {
				ms.Set([]byte(key), value.Value)
			}
			pm.conflictCheck[storeKey].Write[key] = value
		}
	}
	pm.currTxFee = pm.currTxFee.Add(pm.extraTxsInfo[txIndex].fee.Sub(pm.finalResult[txIndex].paraMsg.RefundFee)...)
}

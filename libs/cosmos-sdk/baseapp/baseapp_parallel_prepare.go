package baseapp

import (
	"sync"
	"sync/atomic"

	"github.com/okex/exchain/libs/utils/gopool"
)

var (
	gpTxPool             *parallelTxPool
	createCustomPoolOnce sync.Once
	preparedTxCount      uint64
)

type parallelTxPool struct {
	*gopool.CustomPool
}

type parallelTx struct {
	app     *BaseApp
	index   int
	txBytes []byte
}

func prepare(args interface{}) {
	defer atomic.AddUint64(&preparedTxCount, 1)
	ptx := args.(*parallelTx)
	app := ptx.app
	index := ptx.index
	txBytes := ptx.txBytes

	para := app.parallelTxManage

	tx, err := app.txDecoder(txBytes)
	if err != nil {
		para.extraTxsInfo[index] = &extraDataForTx{
			decodeErr: err,
		}
		return
	}
	coin, isEvm, s, toAddr, _ := app.getTxFeeAndFromHandler(app.getContextForTx(runTxModeDeliver, txBytes), tx)
	para.extraTxsInfo[index] = &extraDataForTx{
		fee:   coin,
		isEvm: isEvm,
		from:  s,
		to:    toAddr,
		stdTx: tx,
	}

}

func getParallelTxPool() *parallelTxPool {
	createCustomPoolOnce.Do(
		func() {
			pool, err := gopool.NewPool(gopool.CustomPoolConfig{Size: 10000}, prepare, gopool.WithNonblocking(true))
			if err != nil {
				panic(err)
			}
			gpTxPool = &parallelTxPool{pool}
		},
	)
	return gpTxPool
}

func (p *parallelTxPool) getExtraData(app *BaseApp, index int, txBytes []byte) error {
	return p.Invoke(&parallelTx{
		app:     app,
		index:   index,
		txBytes: txBytes,
	})
}

func waitTxPrepared(expectedTxCount uint64) {
	for !atomic.CompareAndSwapUint64(&preparedTxCount, expectedTxCount, 0) {
	}
}

func initParallelTxManage(txManager *parallelTxManager) {
	txSize := txManager.txSize
	txRepsCap := cap(txManager.txReps)
	if txManager.txReps == nil || txRepsCap < txSize {
		txManager.txReps = make([]*executeResult, txSize)
	} else if txRepsCap >= txSize {
		txManager.txReps = txManager.txReps[0:txSize:txRepsCap]
		// https://github.com/golang/go/issues/5373
		for i := range txManager.txReps {
			txManager.txReps[i] = nil
		}
	}

	txsInfoCap := cap(txManager.extraTxsInfo)
	if txManager.extraTxsInfo == nil || txsInfoCap < txSize {
		txManager.extraTxsInfo = make([]*extraDataForTx, txSize)
	} else if txsInfoCap >= txSize {
		txManager.extraTxsInfo = txManager.extraTxsInfo[0:txSize:txsInfoCap]
		for i := range txManager.extraTxsInfo {
			txManager.extraTxsInfo[i] = nil
		}
	}

	for key := range txManager.workgroup.runningStatus {
		delete(txManager.workgroup.runningStatus, key)
	}
	for key := range txManager.workgroup.isrunning {
		delete(txManager.workgroup.isrunning, key)
	}
}

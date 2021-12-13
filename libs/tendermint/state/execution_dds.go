package state

import (
	"github.com/okex/exchain/libs/tendermint/delta"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"time"

	"github.com/okex/exchain/libs/tendermint/types"
)

type DeltaContext struct {
	deltaBroker   delta.DeltaBroker
	deltas *types.Deltas

	applyDelta bool
	broadDelta bool
	downloadDelta bool
	uploadDelta bool
	useDeltas bool
	logger log.Logger
}

func newDeltaContext()  *DeltaContext {

	dp := &DeltaContext{}
	dp.applyDelta = types.EnableApplyP2PDelta()
	dp.broadDelta = types.EnableBroadcastP2PDelta()
	dp.downloadDelta = types.EnableDownloadDelta()
	dp.uploadDelta = types.EnableUploadDelta()

	return dp
}

func (dc *DeltaContext) postApplyBlock(block *types.Block)  {
	if dc.broadDelta || dc.uploadDelta {
		dc.deltas.Height = block.Height
	}
	if dc.uploadDelta {
		go dc.uploadData(block, dc.deltas)
	}

	dc.logger.Info("Begin abci",
		"len(deltas)", dc.deltas.Size(),
		"FlagUseDelta", dc.useDeltas)
}

//func (dc *DeltaContext) dump(caller string) {
//
//	dc.logger.Info(caller, "len(deltas)", dc.deltas.Size(),
//		"fastQuery", fastQuery,
//		"FlagUseDelta", useDeltas)
//}

func (dc *DeltaContext ) uploadData(block *types.Block, deltas *types.Deltas) {
	if err := dc.deltaBroker.SetDeltas(deltas); err != nil {
		dc.logger.Error("uploadData err:", err)
		return
	}
	dc.logger.Info("uploadData",
		"height", block.Height,
		"blockLen", block.Size(),
		"abciRspLen", len(deltas.ABCIRsp),
		"deltaLen", len(deltas.DeltasBytes),
		"watchLen", len(deltas.WatchBytes))
}

func (blockExec *BlockExecutor) prepareStateDelta(block *types.Block, deltas *types.Deltas) (bool, *types.Deltas) {
	fastQuery := types.IsFastQuery()
	applyDelta := types.EnableApplyP2PDelta()
	downloadDelta := types.EnableDownloadDelta()

	// not use delta, exe abci itself
	if !applyDelta && !downloadDelta {
		return false, deltas
	}

	// get watchData and Delta from p2p
	if applyDelta {
		if len(deltas.ABCIRsp) > 0 && len(deltas.DeltasBytes) > 0 {
			if !fastQuery || len(deltas.WatchBytes) > 0 {
				return true, deltas
			}
		}
	}

	if !downloadDelta {
		return false, deltas
	}

	var directDelta *types.Deltas
	var err error
	needDDS := true
	select {
	case directDelta = <-blockExec.deltaCh:
		if directDelta.Height == block.Height {
			needDDS = false
		}
		// already get delta of height, then request delta of height+1
		blockExec.deltaHeightCh <- block.Height + 1
	default:
		// can't get delta of height, request delta of height+1 and return
		blockExec.deltaHeightCh <- block.Height + 1
	}

	if needDDS {
		// request watchData and Delta from dds
		directDelta, err = blockExec.deltaContext.deltaBroker.GetDeltas(block.Height)
	}

	// can't get data from dds
	if directDelta == nil {
		if err != nil {
			blockExec.logger.Error("Download Delta err:", err)
		}
		return false, deltas
	}

	//// get watchData from dds
	if !fastQuery || len(directDelta.WatchBytes) > 0 {
		// get Delta from dds
		if len(directDelta.ABCIRsp) > 0 && len(directDelta.DeltasBytes) > 0 {
			return true, directDelta
		}
		// get Delta from p2p
		if len(deltas.ABCIRsp) > 0 && len(deltas.DeltasBytes) > 0 {
			deltas.WatchBytes = directDelta.WatchBytes
			return true, deltas
		}
		// can't get Delta
		return false, deltas
	}

	//// can't get watchData from dds
	{
		if len(deltas.WatchBytes) <= 0 {
			// can't get watchData
			return false, deltas
		}

		// get Delta from dds
		if len(directDelta.ABCIRsp) > 0 && len(directDelta.DeltasBytes) > 0 {
			directDelta.WatchBytes = deltas.WatchBytes
			return true, directDelta
		}
	}

	return false, deltas
}

func (blockExec *BlockExecutor) getDeltaFromDDS() {
	flag := false
	var height int64 = 0
	tryGetDDSTicker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case <-tryGetDDSTicker.C:
			if flag {
				directDelta, _ := blockExec.deltaContext.deltaBroker.GetDeltas(height)
				if directDelta != nil {
					flag = false
					blockExec.deltaCh <- directDelta
				}
			}

		case height = <-blockExec.deltaHeightCh:
			flag = true
		}
	}
}

package state

import (
	gorid "github.com/okex/exchain/libs/goroutine"
	"github.com/okex/exchain/libs/tendermint/delta"
	redis_cgi "github.com/okex/exchain/libs/tendermint/delta/redis-cgi"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"time"

	"github.com/okex/exchain/libs/tendermint/types"
)

type DeltaContext struct {
	deltaCh       chan *types.Deltas
	deltaHeightCh chan int64
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

	dp := &DeltaContext{
		deltaCh:        make(chan *types.Deltas, 1),
		deltaHeightCh:  make(chan int64, 1),
	}
	dp.applyDelta = types.EnableApplyP2PDelta()
	dp.broadDelta = types.EnableBroadcastP2PDelta()
	dp.downloadDelta = types.EnableDownloadDelta()
	dp.uploadDelta = types.EnableUploadDelta()

	return dp
}

func (dc *DeltaContext) init(l log.Logger) {
	dc.logger = l

	dc.logger.Info("DeltaContext init",
		"uploadDelta", dc.uploadDelta,
		"downloadDelta", dc.downloadDelta,
		"applyDelta", dc.applyDelta,
		"broadDelta", dc.broadDelta,
	)

	if dc.uploadDelta || dc.downloadDelta {
		dc.deltaBroker = redis_cgi.NewRedisClient(types.RedisUrl())
		dc.logger.Info("Init delta broker", "url", types.RedisUrl())
	}

	if dc.downloadDelta {
		go dc.getDeltaFromDDS()
	}

}

func (dc *DeltaContext) postApplyBlock(block *types.Block) {
	delta := dc.deltas

	if dc.broadDelta || dc.uploadDelta {
		delta.Height = block.Height
	}
	if dc.uploadDelta {
		go dc.uploadData(block, delta)
	}

	dc.logger.Info("Post apply block",
		"delta", delta,
		"useDeltas", dc.useDeltas)
}


func (dc *DeltaContext) uploadData(block *types.Block, deltas *types.Deltas) {
	//if !dc.uploadDelta {
	//	return
	//}

	if err := dc.deltaBroker.SetDeltas(deltas); err != nil {
		dc.logger.Error("Upload delta", "height", block.Height, "error", err)
		return
	}
	dc.logger.Info("Upload delta",
		"deltas", deltas,
		"blockLen", block.Size(),
		"gid", gorid.GoRId,
	)
}

func (dc *DeltaContext) prepareStateDelta(block *types.Block, deltas *types.Deltas) (bool, *types.Deltas) {
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
	case directDelta = <-dc.deltaCh:
		if directDelta.Height == block.Height {
			needDDS = false
		}
		// already get delta of height, then request delta of height+1
		dc.deltaHeightCh <- block.Height + 1
	default:
		// can't get delta of height, request delta of height+1 and return
		dc.deltaHeightCh <- block.Height + 1
	}

	if needDDS {
		// request watchData and Delta from dds
		directDelta, err = dc.deltaBroker.GetDeltas(block.Height)
	}

	// can't get data from dds
	if directDelta == nil {
		if err != nil {
			dc.logger.Error("Download Delta err:", err)
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

func (dc *DeltaContext) getDeltaFromDDS() {
	flag := false
	var height int64 = 0
	tryGetDDSTicker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case <-tryGetDDSTicker.C:
			if flag {
				directDelta, _ := dc.deltaBroker.GetDeltas(height)
				if directDelta != nil {
					dc.logger.Info("Download delta:",
						"delta", directDelta,
						"gid", gorid.GoRId)
					flag = false
					dc.deltaCh <- directDelta
				}
			}

		case height = <-dc.deltaHeightCh:
			flag = true
		}
	}
}

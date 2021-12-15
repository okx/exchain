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

	dp.deltas = &types.Deltas{}

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

func (dc *DeltaContext) setWatchData(wd []byte) {
	dc.deltas.WatchBytes = wd
}

func (dc *DeltaContext) setAbciRsp(ar []byte) {
	dc.deltas.ABCIRsp = ar
}

func (dc *DeltaContext) setStateDelta(sd []byte) {
	dc.deltas.DeltasBytes = sd
}

func (dc *DeltaContext) postApplyBlock(height int64, abciResponses *ABCIResponses, res []byte) {

	dc.deltas.Height = height
	dc.setStateDelta(res) // settt

	var abciResponsesBytes []byte
	var err error
	if dc.broadDelta || dc.uploadDelta {
		abciResponsesBytes, err = types.Json.Marshal(abciResponses)
		if err != nil {
			panic(err)
		}
		dc.setAbciRsp(abciResponsesBytes) // settt
	}

	if !dc.useDeltas {
		// get deliverTx WatchData and let wd = it
		dc.setWatchData(GetWatchData()) // settt
	} else {
		// commitBatch with wd in exchain
		UseWatchData(dc.deltas.WatchBytes)
	}

	delta4Upload :=  &types.Deltas {
		abciResponsesBytes,
		res,
		GetWatchData(),
		height,
	}


	dc.logger.Info("Post apply block",
		"delta", delta4Upload,
		"useDeltas", dc.useDeltas)

	dc.deltas = &types.Deltas{}

	if dc.uploadDelta {
		go dc.uploadData(delta4Upload)
	}

}


func (dc *DeltaContext) uploadData(deltas *types.Deltas) {
	if deltas == nil {
		return
	}


	if err := dc.deltaBroker.SetDeltas(deltas); err != nil {
		dc.logger.Error("Upload delta", "height", deltas.Height, "error", err)
		return
	}
	dc.logger.Info("Upload delta",
		"deltas", deltas,
		"gid", gorid.GoRId,
	)
}

func (dc *DeltaContext) prepareStateDelta(block *types.Block) {
	dc.useDeltas = false

	// not use delta, exe abci itself
	if !dc.applyDelta && !dc.downloadDelta {
		return
	}

	if !dc.downloadDelta {
		return
	}

	var dds *types.Deltas
	select {
	case dds = <-dc.deltaCh:
		dc.logger.Info("prepareStateDelta", "delta", dds)
		// already get delta of height
	default:
		dc.logger.Info("prepareStateDelta", "delta", dds)
		// can't get delta of height
	}
	// request delta of height+1 and return
	dc.deltaHeightCh <- block.Height + 1

	// can't get data from dds
	if dds == nil || dds.Height != block.Height ||
		len(dds.WatchBytes) < 0 || len(dds.ABCIRsp) < 0 || len(dds.DeltasBytes) < 0 {
		return
	}

	// get Delta from dds
	dc.useDeltas = true
	dc.deltas = dds
	return
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

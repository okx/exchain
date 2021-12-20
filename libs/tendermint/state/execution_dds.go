package state

import (
	gorid "github.com/okex/exchain/libs/goroutine"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/tendermint/delta"
	redis_cgi "github.com/okex/exchain/libs/tendermint/delta/redis-cgi"
	"github.com/okex/exchain/libs/tendermint/libs/compress"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"time"

	"github.com/okex/exchain/libs/tendermint/types"
)

type DeltaContext struct {
	deltaCh       chan *types.Deltas
	deltaHeightCh chan int64
	deltaBroker   delta.DeltaBroker
	deltas        *types.Deltas

	//applyDelta bool
	//broadDelta bool
	downloadDelta bool
	uploadDelta bool
	useDeltas bool
	logger log.Logger

	compressBroker compress.CompressBroker
}

func newDeltaContext() *DeltaContext {

	dp := &DeltaContext{
		deltaCh:       make(chan *types.Deltas, 1),
		deltaHeightCh: make(chan int64, 1),
	}
	//dp.applyDelta = types.EnableApplyP2PDelta()
	//dp.broadDelta = types.EnableBroadcastP2PDelta()
	dp.downloadDelta = types.EnableDownloadDelta()
	dp.uploadDelta = types.EnableUploadDelta()


	if dp.uploadDelta && dp.downloadDelta {
		panic("download delta is not allowed if upload delta enabled")
	}

	dp.deltas = &types.Deltas{}

	// todo can config different compress algorithm
	dp.compressBroker = &compress.Flate{}

	return dp
}

func (dc *DeltaContext) init(l log.Logger) {
	dc.logger = l

	dc.logger.Info("DeltaContext init",
		"uploadDelta", dc.uploadDelta,
		"downloadDelta", dc.downloadDelta,
		//"applyDelta", dc.applyDelta,
		//"broadDelta", dc.broadDelta,
	)

	if dc.uploadDelta || dc.downloadDelta {
		dc.deltaBroker = redis_cgi.NewRedisClient(types.RedisUrl(), l)
		dc.logger.Info("Init delta broker", "url", types.RedisUrl())
	}

	// control if iavl produce delta or not
	iavl.SetProduceDelta(dc.uploadDelta)

	if dc.downloadDelta {
		go dc.getDeltaFromDDS()
	}
}

func (dc *DeltaContext) reset() {
	dc.useDeltas = false
	dc.deltas = &types.Deltas{}
}

func (dc *DeltaContext) postApplyBlock(height int64, abciResponses *ABCIResponses, res []byte) {
	if dc.useDeltas {
		dc.logger.Info("Applied delta", "delta", dc.deltas, "gid", gorid.GoRId)
	}

	// rpc
	if dc.useDeltas {
		UseWatchData(dc.deltas.WatchBytes)
	}

	// validator
	if dc.uploadDelta {
		dc.upload(height, abciResponses, res)
	}
}

func (dc *DeltaContext) upload(height int64, abciResponses *ABCIResponses, res []byte) {

	var abciResponsesBytes []byte
	var err error
	abciResponsesBytes, err = types.Json.Marshal(abciResponses)
	if err != nil {
		dc.logger.Error("Failed to marshal abci Responses", "height", height, "error", err)
		return
	}

	// for outDelta log
	dc.deltas = &types.Deltas {
		ABCIRsp:     abciResponsesBytes,
		DeltasBytes: res,
		WatchBytes:  GetWatchData(),
		Height:      height,
	}

	delta4Upload := &types.Deltas {
		ABCIRsp:     abciResponsesBytes,
		DeltasBytes: res,
		WatchBytes:  GetWatchData(),
		Height:      height,
		Version:     types.DeltaVersion,
	}

	go dc.uploadData(delta4Upload)
}

func (dc *DeltaContext) uploadData(deltas *types.Deltas) {
	if deltas == nil {
		return
	}

	// todo get distributed lock, otherwise return

	// marshal deltas to bytes
	deltaBytes, err := deltas.Marshal()
	if err != nil {
		return
	}

	// compress
	compressBytes, err := dc.compressBroker.DefaultCompress(deltaBytes)
	if err != nil {
		return
	}

	// set into dds
	if err = dc.deltaBroker.SetDeltas(deltas.Height, compressBytes); err != nil {
		dc.logger.Error("Upload delta", "height", deltas.Height, "error", err)
		return
	}
	dc.logger.Info("Upload delta", "deltas", deltas, "gid", gorid.GoRId)
}

func (dc *DeltaContext) prepareStateDelta(block *types.Block) {
	if !dc.downloadDelta {
		return
	}

	var dds *types.Deltas
	select {
	case dds = <-dc.deltaCh:
		// already get delta of height
	default:
		// can't get delta of height
	}
	dc.logger.Info("prepareStateDelta", "delta", dds)

	// request delta of height+1 and return
	dc.deltaHeightCh <- block.Height + 1

	// can't get data from dds
	// or get invalid data
	if dds == nil || dds.Height != block.Height || types.DeltaVersion < dds.Version ||
		len(dds.WatchBytes) == 0 || len(dds.ABCIRsp) == 0 || len(dds.DeltasBytes) == 0 {
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
				deltaBytes, err := dc.deltaBroker.GetDeltas(height)
				if err != nil {
					continue
				}
				flag = false

				// uncompress
				compressBytes, err := dc.compressBroker.UnCompress(deltaBytes)
				if err != nil {
					continue
				}

				// unmarshal
				directDelta := &types.Deltas{}
				err = directDelta.Unmarshal(compressBytes)
				if err != nil {
					continue
				}
				
				dc.logger.Info("Download delta:", "delta", directDelta, "gid", gorid.GoRId)
				dc.deltaCh <- directDelta
			}

		case height = <-dc.deltaHeightCh:
			flag = true
		}
	}
}

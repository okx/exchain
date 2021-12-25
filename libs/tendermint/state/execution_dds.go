package state

import (
	"fmt"
	gorid "github.com/okex/exchain/libs/goroutine"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/tendermint/delta"
	redis_cgi "github.com/okex/exchain/libs/tendermint/delta/redis-cgi"
	"github.com/okex/exchain/libs/tendermint/libs/compress"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/trace"
	"sync/atomic"
	"time"

	"github.com/okex/exchain/libs/tendermint/types"
)

type DeltaContext struct {
	deltaBroker   delta.DeltaBroker
	lastCommitHeight int64
	dataMap *deltaMap

	downloadDelta bool
	uploadDelta bool
	applied float64
	missed float64
	logger log.Logger
	compressBroker compress.CompressBroker
}

func newDeltaContext() *DeltaContext {

	dp := &DeltaContext{
		dataMap: newDataMap(),
	}
	dp.downloadDelta = types.EnableDownloadDelta()
	dp.uploadDelta = types.EnableUploadDelta()

	if dp.uploadDelta && dp.downloadDelta {
		panic("download delta is not allowed if upload delta enabled")
	}

	// todo can config different compress algorithm
	dp.compressBroker = &compress.Flate{}
	dp.missed = 0.000001

	return dp
}

func (dc *DeltaContext) init(l log.Logger) {
	dc.logger = l

	dc.logger.Info("DeltaContext init",
		"uploadDelta", dc.uploadDelta,
		"downloadDelta", dc.downloadDelta,
	)

	if dc.uploadDelta || dc.downloadDelta {
		dc.deltaBroker = redis_cgi.NewRedisClient(types.RedisUrl(), types.RedisAuth(), types.RedisExpire(), l)
		dc.logger.Info("Init delta broker", "url", types.RedisUrl())
	}

	// control if iavl produce delta or not
	iavl.SetProduceDelta(dc.uploadDelta)

	if dc.downloadDelta {
		go dc.downloadRoutine()
	}
}

func (dc *DeltaContext) appliedRate() float64 {
	return dc.applied / (dc.applied + dc.missed)
}

func (dc *DeltaContext) postApplyBlock(height int64, delta *types.Deltas,
	abciResponses *ABCIResponses, res []byte) {

	// rpc
	if dc.downloadDelta {

		applied := false
		if delta != nil {
			applied = true
			dc.applied += float64(len(abciResponses.DeliverTxs))
		} else {
			dc.missed += float64(len(abciResponses.DeliverTxs))
		}

		trace.GetElapsedInfo().AddInfo(trace.Delta,
			fmt.Sprintf("applied<%t>, rate<%.2f>", applied, dc.appliedRate()))

		dc.logger.Info("Post apply block", "height", height, "delta-applied", applied,
			"applied-rate", dc.appliedRate(), "delta", delta)

		if applied && types.IsFastQuery() {
			UseWatchData(delta.WatchBytes)
		}
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

	wd := GetWatchData()

	delta4Upload := &types.Deltas {
		ABCIRsp:     abciResponsesBytes,
		DeltasBytes: res,
		WatchBytes:  wd,
		Height:      height,
		Version:     types.DeltaVersion,
	}

	go dc.uploadData(delta4Upload)
}

func (dc *DeltaContext) uploadData(deltas *types.Deltas) {
	if deltas == nil {
		return
	}

	dc.logger.Info("Upload delta started:", "target-height", deltas.Height, "gid", gorid.GoRId)
	locked := dc.deltaBroker.GetLocker()
	dc.logger.Info("Upload delta:", "locked", locked, "gid", gorid.GoRId)
	if !locked {
		return
	}

	defer dc.deltaBroker.ReleaseLocker()

	upload := func() bool {
		t0 := time.Now()
		// marshal deltas to bytes
		deltaBytes, err := deltas.Marshal()
		if err != nil {
			dc.logger.Error("Failed to upload delta", "target-height", deltas.Height, "error", err)
			return false
		}

		t1 := time.Now()
		// compress
		//compressBytes, err := dc.compressBroker.DefaultCompress(deltaBytes)
		//if err != nil {
		//	return
		//}
		t2 := time.Now()
		// set into dds
		if err = dc.deltaBroker.SetDeltas(deltas.Height, deltaBytes); err != nil {
			dc.logger.Error("Failed to upload delta", "target-height", deltas.Height, "error", err)
			return false

		}
		t3 := time.Now()
		dc.logger.Info("Upload delta finished",
			"target-height", deltas.Height,
			"marshal", t1.Sub(t0),
			"compress", t2.Sub(t1),
			"setRedis", t3.Sub(t2),
			"deltas", deltas,
			"gid", gorid.GoRId)
		return true
	}
	dc.deltaBroker.ResetLatestHeightAfterUpload(deltas.Height, upload)
}

// get delta from dds
func (dc *DeltaContext) prepareStateDelta(height int64) (dds *types.Deltas) {
	if !dc.downloadDelta {
		return
	}
	dds = dc.dataMap.fetch(height)
	var succeed bool
	if dds != nil {
		if !dds.Validate(height) {
			dc.logger.Error("Prepared an invalid delta!!!", "expected-height", height, "delta", dds)
			return nil
		}
		succeed = true
	}
	dc.logger.Info("Prepare delta", "expected-height", height, "succeed", succeed, "delta", dds)
	return
}

func (dc *DeltaContext) downloadRoutine() {
	var height int64
	var lastRemoved int64
	var buffer int64 = 5
	ticker := time.NewTicker(50 * time.Millisecond)

	for range ticker.C {
		lastCommitHeight := atomic.LoadInt64(&dc.lastCommitHeight)
		if height <= lastCommitHeight {
			// move to lastCommitHeight + 1
			height = lastCommitHeight + 1

			// git rid of all deltas before <height>
			removed, left := dc.dataMap.remove(lastCommitHeight)
			dc.logger.Info("Updated target delta height",
				"target-height", height,
				"lastCommitHeight", lastCommitHeight,
				"removed", removed,
				"left", left,
			)
		} else {
			if height % 10 == 0 && lastRemoved != lastCommitHeight {
				removed, left := dc.dataMap.remove(lastCommitHeight)
				dc.logger.Info("Remove stale delta",
					"target-height", height,
					"lastCommitHeight", lastCommitHeight,
					"removed", removed,
					"left", left,
				)
				lastRemoved = lastCommitHeight
			}
		}

		lastCommitHeight = atomic.LoadInt64(&dc.lastCommitHeight)
		if height > lastCommitHeight+buffer {
			continue
		}

		err, delta := dc.download(height)
		if err == nil {
			dc.dataMap.insert(height, delta)
			height++
		}
	}
}

func (dc *DeltaContext) download(height int64) (error, *types.Deltas){
	dc.logger.Debug("Download delta started:", "target-height", height, "gid", gorid.GoRId)

	t0 := time.Now()
	deltaBytes, err := dc.deltaBroker.GetDeltas(height)
	if err != nil {
		return err, nil
	}

	t1 := time.Now()
	// uncompress
	//compressBytes, err := dc.compressBroker.UnCompress(deltaBytes)
	//if err != nil {
	//	continue
	//}

	t2 := time.Now()
	// unmarshal
	delta := &types.Deltas{}
	err = delta.Unmarshal(deltaBytes)
	if err != nil {
		dc.logger.Error("Downloaded an invalid delta:", "target-height", height, "err", err,)
		return err, nil
	}

	cacheMap, cacheList := dc.dataMap.info()
	t3 := time.Now()
	dc.logger.Info("Download delta finished:",
		"target-height", height,
		"cacheMap", cacheMap,
		"cacheList", cacheList,
		"download", t1.Sub(t0),
		"uncompress", t2.Sub(t1),
		"unmarshal", t3.Sub(t2),
		"delta", delta,
		"gid", gorid.GoRId)

	return nil, delta
}

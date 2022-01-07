package state

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/okex/exchain/libs/tendermint/delta"
	redis_cgi "github.com/okex/exchain/libs/tendermint/delta/redis-cgi"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"os"
	"testing"
	"time"
)

func getRedisClient(t *testing.T) *redis_cgi.RedisClient {
	s := miniredis.RunT(t)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	ss := redis_cgi.NewRedisClient(s.Addr(), "", time.Minute, logger)
	return ss
}

func TestDeltaContext_upload(t *testing.T) {
	type fields struct {
		deltaBroker       delta.DeltaBroker
		lastFetchedHeight int64
		dataMap           *deltaMap
		downloadDelta     bool
		uploadDelta       bool
		hit               float64
		missed            float64
		logger            log.Logger
		compressType      int
		compressFlag      int
		bufferSize        int
		idMap             identityMapType
		identity          string
	}
	type args struct {
		deltas *tmtypes.Deltas
		txnum  float64
		mrh    int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DeltaContext{
				deltaBroker:       tt.fields.deltaBroker,
				lastFetchedHeight: tt.fields.lastFetchedHeight,
				dataMap:           tt.fields.dataMap,
				downloadDelta:     tt.fields.downloadDelta,
				uploadDelta:       tt.fields.uploadDelta,
				hit:               tt.fields.hit,
				missed:            tt.fields.missed,
				logger:            tt.fields.logger,
				compressType:      tt.fields.compressType,
				compressFlag:      tt.fields.compressFlag,
				bufferSize:        tt.fields.bufferSize,
				idMap:             tt.fields.idMap,
				identity:          tt.fields.identity,
			}
			if got := dc.upload(tt.args.deltas, tt.args.txnum, tt.args.mrh); got != tt.want {
				t.Errorf("upload() = %v, want %v", got, tt.want)
			}
		})
	}
}

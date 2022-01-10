package state

import (
	"encoding/hex"
	"github.com/alicebob/miniredis/v2"
	"github.com/okex/exchain/libs/tendermint/delta/redis-cgi"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"reflect"
	"testing"
	"time"
)

func getRedisClient(t *testing.T) *redis_cgi.RedisClient {
	s := miniredis.RunT(t)
	logger := log.TestingLogger()
	ss := redis_cgi.NewRedisClient(s.Addr(), "", time.Minute, logger)
	return ss
}

func failRedisClient() *redis_cgi.RedisClient {
	logger := log.TestingLogger()
	ss := redis_cgi.NewRedisClient("127.0.0.1:6378", "", time.Minute, logger)
	return ss
}

func setupTest(t *testing.T) *DeltaContext {
	dc := newDeltaContext(log.TestingLogger())
	dc.deltaBroker = getRedisClient(t)
	return dc
}

func bytesEqual(b1, b2 []byte) bool {
	return hex.EncodeToString(b1) == hex.EncodeToString(b2)
}
func deltaEqual(d1, d2 *types.Deltas) bool {
	if d1 == nil && d2 == nil { return true }
	if d1 == nil || d2 == nil { return false }
	return d1.Height == d2.Height &&
		d1.Version == d2.Version &&
		d1.From == d2.From &&
		d1.CompressType == d2.CompressType &&
		d1.CompressFlag == d2.CompressFlag &&
		bytesEqual(d1.ABCIRsp(), d2.ABCIRsp()) &&
		bytesEqual(d1.DeltasBytes(), d2.DeltasBytes()) &&
		bytesEqual(d1.WatchBytes(), d2.WatchBytes())
}

func TestDeltaContext_prepareStateDelta(t *testing.T) {
	dc := setupTest(t)

	type args struct {
		height int64
	}
	tests := []struct {
		name    string
		args    args
		wantDds *types.Deltas
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDds := dc.prepareStateDelta(tt.args.height); !reflect.DeepEqual(gotDds, tt.wantDds) {
				t.Errorf("prepareStateDelta() = %v, want %v", gotDds, tt.wantDds)
			}
		})
	}
}

func TestDeltaContext_download(t *testing.T) {
	dc := setupTest(t)
	deltas := &types.Deltas{Height: 1, Payload: types.DeltaPayload{ABCIRsp: []byte("ABCIRsp"), DeltasBytes: []byte("DeltasBytes"), WatchBytes: []byte("WatchBytes")}}
	dc.uploadRoutine(deltas, 0)

	tests := []struct {
		name   string
		height int64
		wants  *types.Deltas
	}{
		{"normal case", 1, deltas},
		{"wrong height", 11, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err, mrh := dc.deltaBroker.GetDeltas(tt.height)
			got, got1, got2 := dc.download(tt.height)
			if !reflect.DeepEqual(got, err) {
				t.Errorf("download() got = %v, want %v", got, err)
			}
			if !deltaEqual(got1, tt.wants) {
				t.Errorf("download() got = %v, want %v", got, deltas)
			}
			if got2 != mrh {
				t.Errorf("download() got2 = %v, want %v", got2, mrh)
			}
		})
	}
}

func TestDeltaContext_upload(t *testing.T) {
	dc := setupTest(t)
	deltas := &types.Deltas{Payload: types.DeltaPayload{ABCIRsp: []byte("ABCIRsp"), DeltasBytes: []byte("DeltasBytes"), WatchBytes: []byte("WatchBytes")}}
	okRedis := getRedisClient(t)
	failRedis := failRedisClient()

	tests := []struct {
		name   string
		r      *redis_cgi.RedisClient
		deltas *types.Deltas
		want   bool
	}{
		{"normal case", okRedis, deltas, true},
		{"nil delta", okRedis, nil, false},
		{"empty delta", okRedis, &types.Deltas{}, false},
		{"fail redis", failRedis, deltas, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc.deltaBroker = tt.r
			if got := dc.upload(tt.deltas, 0, 0); got != tt.want {
				t.Errorf("upload() = %v, want %v", got, tt.want)
			}
		})
	}
}

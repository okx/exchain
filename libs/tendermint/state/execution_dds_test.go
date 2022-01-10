package state

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/okex/exchain/libs/tendermint/delta/redis-cgi"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
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

func produceDelta() {

}

func TestDeltaContext_upload(t *testing.T) {
	dc := setupTest(t)
	type args struct {
		deltas *types.Deltas
		txnum  float64
		mrh    int64
	}
	tests := []struct {
		name   string
		args   args
		want   bool
	}{
		{"normal case", args{nil, 0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := dc.upload(tt.args.deltas, tt.args.txnum, tt.args.mrh); got != tt.want {
			//	t.Errorf("upload() = %v, want %v", got, tt.want)
			//}
			assert.Panics(t, func() {
				dc.upload(tt.args.deltas, tt.args.txnum, tt.args.mrh)
			})
		})
	}

	// connect to redis failed
	dc.deltaBroker = failRedisClient()
	got := dc.upload(&types.Deltas{}, 0, 0)
	assert.False(t, got)
}

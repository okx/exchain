package redis_cgi

import (
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	ConstRedisURL   = "127.0.0.1:6379"
	ConstDeltaBytes = "delta-bytes"
	ConstTestHeight = 1000000
)

func getRedisClient() *RedisClient {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	ss := NewRedisClient(ConstRedisURL, "", time.Minute, logger)
	return ss
}

func TestRedisClient_SetGetDeltas(t *testing.T) {
	r := getRedisClient()
	require.True(t, r != nil, r)

	height := int64(ConstTestHeight)
	// delta is empty
	re, err := r.GetDeltas(height)
	require.True(t, re == nil, re)
	require.True(t, err != nil, err)

	// set delta
	bytes := []byte(ConstDeltaBytes)
	err = r.SetDeltas(height, bytes)
	require.Nil(t, err)

	// get delta
	re, err = r.GetDeltas(height)
	require.True(t, re != nil, re)
	require.True(t, err == nil, err)

	// get wrong key
	fakeKey := height + 1
	noResult, err := r.GetDeltas(fakeKey)
	require.True(t, noResult == nil, noResult)
	require.True(t, err != nil, err)
}

func TestRedisClient_ResetLatestHeightAfterUpload1(t *testing.T) {
	r := getRedisClient()
	require.True(t, r != nil, r)
	uploadSuccess := func() bool {return true}
	uploadFailed := func() bool {return false}
	h := int64(ConstTestHeight)
	type args struct {
		height int64
		upload func() bool
	}
	tests := []struct {
		name   string
		args   args
		want   bool
	}{
		{"upload failed", args{h, uploadFailed}, false},
		{"first time set", args{h, uploadSuccess}, true},
		{"height < latestHeight", args{h-1, uploadSuccess}, false},
		{"height == latestHeight", args{h, uploadSuccess}, false},
		{"height > latestHeight", args{h+1, uploadSuccess}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.ResetLatestHeightAfterUpload(tt.args.height, tt.args.upload); got != tt.want {
				t.Errorf("ResetLatestHeightAfterUpload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_GetReleaseLocker(t *testing.T) {
	r := getRedisClient()
	require.True(t, r != nil, r)

	// first time lock
	locker := r.GetLocker()
	require.True(t, locker, locker)

	// already locked
	locker = r.GetLocker()
	require.True(t, !locker, locker)

	// release locker
	r.ReleaseLocker()
	locker = r.GetLocker()
	require.True(t, locker, locker)

	// when locker expire time, locker release itself
	time.Sleep(lockerExpire)
	locker = r.GetLocker()
	require.True(t, locker, locker)
}
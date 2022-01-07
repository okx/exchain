package redis_cgi

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	ConstDeltaBytes = "delta-bytes"
	ConstTestHeight = 1
)

func getRedisClient(t *testing.T) *RedisClient {
	s := miniredis.RunT(t)
	logger := log.TestingLogger()
	ss := NewRedisClient(s.Addr(), "", time.Minute, logger)
	return ss
}

func TestRedisClient_SetGetDeltas(t *testing.T) {
	r := getRedisClient(t)
	require.True(t, r != nil, r)

	height := int64(ConstTestHeight)
	// delta is empty
	re, err, _ := r.GetDeltas(height)
	require.True(t, re == nil, re)
	require.True(t, err != nil, err)

	// set delta
	bytes := []byte(ConstDeltaBytes)
	err = r.SetDeltas(height, bytes)
	require.Nil(t, err)

	// get delta
	re, err, _ = r.GetDeltas(height)
	require.True(t, re != nil, re)
	require.True(t, err == nil, err)

	// get wrong key
	fakeKey := height + 1
	noResult, err, _ := r.GetDeltas(fakeKey)
	require.True(t, noResult == nil, noResult)
	require.True(t, err != nil, err)
}

func TestRedisClient_ResetLatestHeightAfterUpload(t *testing.T) {
	r := getRedisClient(t)
	require.True(t, r != nil, r)
	uploadSuccess := func(int64) bool {return true}
	uploadFailed := func(int64) bool {return false}
	h := int64(ConstTestHeight)
	type args struct {
		height int64
		upload func(int64) bool
	}
	tests := []struct {
		name   string
		args   args
		want   bool
	}{
		{"upload failed", args{h, uploadFailed}, false},
		{"first time set", args{h, uploadSuccess}, true},
		{"height<latestHeight", args{h-1, uploadSuccess}, false},
		{"height==latestHeight", args{h, uploadSuccess}, false},
		{"height>latestHeight", args{h+1, uploadSuccess}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _ := r.ResetMostRecentHeightAfterUpload(tt.args.height, tt.args.upload)
			if got != tt.want {
				t.Errorf("ResetLatestHeightAfterUpload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_GetReleaseLocker(t *testing.T) {
	r := getRedisClient(t)
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
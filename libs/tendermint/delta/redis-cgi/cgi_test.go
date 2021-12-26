package redis_cgi

import (
	"github.com/google/uuid"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	ConstRedisURL   = "127.0.0.1:16379"
	ConstLocker     = "unittest"
	ConstStateKey   = "unittest-sk"
	ConstTestHeight = 1000000
)

func getRedisClient() (*RedisClient, error) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	ss, err := NewRedisClient(ConstRedisURL, "", ConstLocker, time.Minute, logger)
	return ss, err
}

func TestRedisClient_ResetLatestHeightAfterUpload(t *testing.T) {
	ss, err := getRedisClient()
	require.True(t, ss != nil, ss)
	require.True(t, err == nil, err)

	height := int64(ConstTestHeight)
	random, err := uuid.NewRandom()
	require.Nil(t, err)
	bytes := []byte(random.String())

	// 1. first time upload and set LatestHeight => success
	succeed := ss.ResetLatestHeightAfterUpload(height, bytes)
	latestHeight, err := ss.GetDistState(latestHeightKey)
	require.Nil(t, err)
	require.True(t, succeed, succeed)
	require.True(t, latestHeight == height, height)

	// 2. height < latestHeight => set failed
	height =
	succeed = ss.ResetLatestHeightAfterUpload(height, bytes)
	require.True(t, !succeed, succeed)
}

func TestRedisClient_SetGetDeltas(t *testing.T) {
	ss, err := getRedisClient()
	require.True(t, ss != nil, ss)
	require.True(t, err == nil, err)

	random, err := uuid.NewRandom()
	require.Nil(t, err)
	height := int64(ConstTestHeight)
	s, err := ss.GetDeltas(height)
	require.True(t, s == nil, s)
	require.True(t, err != nil, err)

	deltaBytes := []byte(random.String())
	err = ss.SetDeltas(height, deltaBytes)
	require.True(t, err == nil, err)

	result, err := ss.GetDeltas(height)
	require.True(t, result != nil, result)
	require.True(t, err == nil, err)

	fakeKey := int64(random.ID())
	noResult, err := ss.GetDeltas(fakeKey)
	require.True(t, noResult == nil, noResult)
	require.True(t, err == nil, err)
}

func TestRedisClient_SetGetBlock(t *testing.T) {
	ss, err := getRedisClient()
	require.True(t, ss != nil, ss)
	require.True(t, err == nil, err)

	random, err := uuid.NewRandom()
	require.Nil(t, err)
	height := int64(ConstTestHeight)
	s, err := ss.GetBlock(height)
	require.True(t, s == nil, s)
	require.True(t, err != nil, err)

	deltaBytes := []byte(random.String())
	err = ss.SetBlock(height, deltaBytes)
	require.True(t, err == nil, err)

	result, err := ss.GetBlock(height)
	require.True(t, result != nil, result)
	require.True(t, err == nil, err)

	fakeKey := int64(random.ID())
	noResult, err := ss.GetBlock(fakeKey)
	require.True(t, noResult == nil, noResult)
	require.True(t, err == nil, err)
}
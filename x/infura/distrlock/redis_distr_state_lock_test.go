package distrlock

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"
)

const ConstRedisURL = "127.0.0.1:6379"
const ConstLocker = "unittest"
const ConstStateKey = "unittest-sk"

func getRedisDistributeStateService() (*RedisDistributeStateService, error) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	ss, err := NewRedisDistributeStateService(ConstRedisURL, "", 0, logger, ConstLocker)
	ss.client.FlushDB(context.Background())
	return ss, err
}

func TestRedisDistributeStateService_SetGetDistState(t *testing.T) {
	ss, err := getRedisDistributeStateService()
	require.True(t, ss != nil, ss)
	require.True(t, err == nil, err)

	random, err := uuid.NewRandom()
	require.Nil(t, err)
	testKey := ConstStateKey
	s := ss.GetDistState(testKey)
	require.True(t, s == "", s)

	value := random.String()
	err = ss.SetDistState(testKey, value)
	require.True(t, err == nil, err)

	result := ss.GetDistState(testKey)
	require.True(t, result == value, result)

	fakeKey := random.String()
	noResult := ss.GetDistState(fakeKey)
	require.True(t, noResult == "", result)
}

func TestRedisDistributeStateService_FetchReleaseDistLock(t *testing.T) {
	ss, err := getRedisDistributeStateService()
	require.Nil(t, err)
	expiredInMS := 2000

	// 1. Origin locker success to get lock
	lockKey := "TestGlobalLock"
	success, err := ss.FetchDistLock(lockKey, ConstLocker, expiredInMS)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	// 2. Second locker try to get lock, but failed.
	fakeLocker := ConstLocker + "_fake"
	success, err = ss.FetchDistLock(lockKey, fakeLocker, expiredInMS)
	require.True(t, !success, success)
	require.True(t, err == nil, err)

	// 3. Wait utils lock expires. fakeLocker try to fetch dlock again.
	time.Sleep(2 * time.Duration(expiredInMS) * time.Millisecond)
	success, err = ss.FetchDistLock(lockKey, fakeLocker, expiredInMS)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	// 4. fakeLocker success to release dlock before expires
	success, err = ss.ReleaseDistLock(lockKey, fakeLocker)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	// 5. fakeLocker fail to release lock when the very lock has been released.
	success, err = ss.ReleaseDistLock(lockKey, fakeLocker)
	require.True(t, !success, success)
	require.True(t, err == nil, err)
}

func TestRedisDistributeStateService_UnlockDistLockWithState(t *testing.T) {
	ss, err := getRedisDistributeStateService()
	require.Nil(t, err)

	expiredInMS := 3000
	lockKey := "TestGlobalLock"
	success, err := ss.UnlockDistLockWithState(lockKey, ConstLocker, ConstStateKey, "_______")
	require.True(t, success == false, success)
	require.True(t, err == nil, err)

	success, err = ss.FetchDistLock(lockKey, ConstLocker, expiredInMS)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	r, err := uuid.NewRandom()
	require.Nil(t, err)
	nextStateValue := "_______" + r.String()
	success, err = ss.UnlockDistLockWithState(lockKey, ConstLocker, ConstStateKey, nextStateValue)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	s := ss.GetDistState(ConstStateKey)
	require.True(t, s == nextStateValue, s)

}

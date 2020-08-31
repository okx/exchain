package distrlock

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

const ConstRedisUrl = "redis://127.0.0.1:16379"
const ConstLocker = "unittest"
const ConstStateKey = "unittest-sk"

func getRedisDistributeStateService() (*RedisDistributeStateService, error) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	ss, err := NewRedisDistributeStateService(ConstRedisUrl, "", logger, ConstLocker)
	return ss, err
}

func TestRedisDistributeStateService_SetGetDistState(t *testing.T) {
	ss, err := getRedisDistributeStateService()
	require.True(t, ss != nil, ss)
	require.True(t, err == nil, err)

	random, _ := uuid.NewRandom()
	testKey := ConstStateKey
	s, err := ss.GetDistState(testKey)
	require.True(t, s == "", s)
	require.True(t, err == nil, err)

	err = ss.SetDistState(testKey, random.String())
	require.True(t, err == nil, err)

	result, err := ss.GetDistState(testKey)
	require.True(t, result != "", result)
	require.True(t, err == nil, err)

	fakeKey := random.String()
	noResult, err := ss.GetDistState(fakeKey)
	require.True(t, noResult == "", result)
	require.True(t, err == nil, err)

}

func TestRedisDistributeStateService_FetchReleaseDistLock(t *testing.T) {
	ss, _ := getRedisDistributeStateService()

	expiredInMS := 1000

	// 1. Origin locker success 2 get dlock
	lockKey := "TestGlobalLock"
	success, err := ss.FetchDistLock(lockKey, ConstLocker, expiredInMS)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	// 2. Second locker try 2 get dlock, but failed.
	fakeLocker := ConstLocker + "_fake"
	success, err = ss.FetchDistLock(lockKey, fakeLocker, expiredInMS)
	require.True(t, !success, success)
	require.True(t, err == nil, err)

	// 3. Wait utils dlock expires. fakeLocker try to fetch dlock again.
	time.Sleep(time.Second * 2)
	success, err = ss.FetchDistLock(lockKey, fakeLocker, expiredInMS)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	// 4. fakelocker success to release dlock before expires
	success, err = ss.ReleaseDistLock(lockKey, fakeLocker)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	//5. fakelocker fail to release dlock when the very lock has been released.
	success, err = ss.ReleaseDistLock(lockKey, fakeLocker)
	require.True(t, !success, success)
	require.True(t, err == nil, err)
}

func TestRedisDistributeStateService_UnlockDistLockWithState(t *testing.T) {
	ss, _ := getRedisDistributeStateService()

	expiredInMS := 5000
	lockKey := "TestGlobalLock"
	success, err := ss.UnlockDistLockWithState(lockKey, ConstLocker, ConstStateKey, "_______")
	require.True(t, success == false, success)
	require.True(t, err == nil, err)

	success, err = ss.FetchDistLock(lockKey, ConstLocker, expiredInMS)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	r, _ := uuid.NewRandom()
	nextStateValue := "_______" + r.String()
	success, err = ss.UnlockDistLockWithState(lockKey, ConstLocker, ConstStateKey, nextStateValue)
	require.True(t, success, success)
	require.True(t, err == nil, err)

	s, err := ss.GetDistState(ConstStateKey)
	require.True(t, s == nextStateValue, s)
	require.True(t, err == nil, err)

}

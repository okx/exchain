package distrlock

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

func TestLocalServiceSmoke(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	stateLockKey := "stream_lock_test"
	latestTaskKey := "stream_latest_task_test"
	stateValue := "{}"

	s, err := NewLocalStateService(logger, "TestLocalServiceSmoke", "/tmp/TestLocalServiceSmoke")
	require.Nil(t, err)
	err = s.RemoveStateFile(stateLockKey)
	require.Nil(t, err)

	require.Equal(t, "TestLocalServiceSmoke", s.GetLockerID())
	state, err := s.GetDistState(stateLockKey)
	require.Nil(t, err)
	require.Equal(t, "", state)

	err = s.SetDistState(stateLockKey, stateValue)
	require.Nil(t, err)

	state, err = s.GetDistState(stateLockKey)
	require.Nil(t, err)
	require.Equal(t, stateValue, state)

	got, err := s.FetchDistLock(stateLockKey, s.GetLockerID(), 1)
	require.True(t, got && err == nil)

	released, err := s.ReleaseDistLock(stateLockKey, s.GetLockerID())
	require.True(t, released && err == nil)

	ok, err := s.UnlockDistLockWithState(stateLockKey, s.GetLockerID(), latestTaskKey, stateValue)
	require.True(t, ok && err == nil)
}

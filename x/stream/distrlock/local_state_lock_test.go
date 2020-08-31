package distrlock

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func TestLocalServiceSmoke(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	stateLockKey := "stream_lock_test"
	latestTaskKey := "stream_latest_task_test"
	stateValue := "{}"

	s, err := NewLocalStateService(logger, "TestLocalServiceSmoke", "/tmp/TestLocalServiceSmoke")
	require.Nil(t, err)
	s.RemoveStateFile(stateLockKey)

	require.Equal(t, "TestLocalServiceSmoke", s.GetLockerId())
	state, err := s.GetDistState(stateLockKey)
	require.Nil(t, err)
	require.Equal(t, "", state)

	err = s.SetDistState(stateLockKey, stateValue)
	require.Nil(t, err)

	state, err = s.GetDistState(stateLockKey)
	require.Nil(t, err)
	require.Equal(t, stateValue, state)

	got, err := s.FetchDistLock(stateLockKey, s.GetLockerId(), 1)
	require.True(t, got && err == nil)

	released, err := s.ReleaseDistLock(stateLockKey, s.GetLockerId())
	require.True(t, released && err == nil)

	ok, err := s.UnlockDistLockWithState(stateLockKey, s.GetLockerId(), latestTaskKey, stateValue)
	require.True(t, ok && err == nil)
}

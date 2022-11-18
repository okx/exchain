package baseapp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecreateHguDB(t *testing.T) {
	CreateHguDB()
	t.Log("hguDB path:", hguPath)
	testKey, testValue := []byte("test"), []byte("test")
	err := guDB.Set(testKey, testValue)
	require.NoError(t, err)
	data, err := guDB.Get([]byte("test"))
	require.NoError(t, err)
	require.Equal(t, testValue, data)
	RecreateHguDB()
	data, err = guDB.Get([]byte("test"))
	require.NoError(t, err)
	require.True(t, data == nil)
	DeleteHguDB()
	require.True(t, guDB == nil)
}

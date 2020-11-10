package conn

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	TestRedisAddr     = "localhost:16379"
	TestRedisPassword = ""
	TestRedistDB      = 0
)

func TestClient_Get(t *testing.T) {
	client, err := NewClient(TestRedisAddr, TestRedisPassword, TestRedistDB, log.NewTMLogger(os.Stdout))
	require.NoError(t, err)

	key := time.Now().String()
	val := "test-value"

	err = client.Set(key, val)
	require.NoError(t, err)

	result, err := client.Get(key)
	require.NoError(t, err)
	require.Equal(t, val, result)
}

func TestClient_MGet(t *testing.T) {
	client, err := NewClient(TestRedisAddr, TestRedisPassword, TestRedistDB, log.NewTMLogger(os.Stdout))
	require.NoError(t, err)

	key1 := time.Now().String()
	val1 := "test-value1"

	err = client.Set(key1, val1)
	require.NoError(t, err)

	key2 := time.Now().Add(3 * time.Second).String()
	val2 := "test-value2"

	err = client.Set(key2, val2)
	require.NoError(t, err)

	result, err := client.MGet([]string{key1, key2})
	require.NoError(t, err)
	require.Equal(t, 2, len(result))
}

func TestClient_HGetAll(t *testing.T) {
	client, err := NewClient(TestRedisAddr, TestRedisPassword, TestRedistDB, log.NewTMLogger(os.Stdout))
	require.NoError(t, err)

	key := time.Now().String()

	field1 := "test-field1"
	val1 := "test-value1"

	field2 := "test-field2"
	val2 := "test-value2"

	err = client.redisCli.HMSet(key, map[string]interface{}{field1: val1, field2: val2}).Err()
	require.NoError(t, err)

	result, err := client.HGetAll(key)
	require.NoError(t, err)
	require.Equal(t, 2, len(result))
}

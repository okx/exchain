package watcher

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

const ConstRedisURL = "redis://10.0.240.21:6379"
const ConstRedisPW = ""

//const ConstRedisURL = "redis://127.0.0.1:16379"

func TestFun(t *testing.T) {
	key := []byte("aaaa")
	val := []byte("bbbb")

	redis := initRedisDB(ConstRedisURL, ConstRedisPW)
	redis.Set(key, val)
	result, err := redis.Get(key)
	require.Nil(t, err)
	require.True(t, bytes.Equal(val, result))
	require.True(t, redis.Has(key))
	redis.Delete(key)
	require.True(t, !redis.Has(key))
}

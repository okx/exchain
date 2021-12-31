/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2022/1/2 11:02 上午
# @File : memory_test.go.go
# @Description :
# @Attention :
*/
package memory

import (
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)


func TestExpire(t *testing.T) {
	broker := NewMemoryBroker(log.NewTMLogger(os.Stdout), MemoryBrokerWithTTL(time.Second * 3))
	broker.ValidateBasic()
	data := []byte{
		1, 2, 3,
	}
	broker.SetDeltas(1, data)

	time.Sleep(time.Second * 5)

	v, e, h := broker.GetDeltas(1)
	require.NoError(t, e)
	require.Equal(t, int64(0), h)
	require.Nil(t, v)
}

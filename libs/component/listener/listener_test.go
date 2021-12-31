/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2022/1/3 8:46 上午
# @File : listener_test.go.go
# @Description :
# @Attention :
*/
package listener

import (
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)



func TestRepeat(t *testing.T) {
	var l IListenerComponent = DefaultNewListenerComponent(log.NewTMLogger(os.Stdout), "listener")
	l.Start()
	ch1:=l.RegisterListener("1")
	ch2:=l.RegisterListener("1")
	l.NotifyListener(1,"1")
	require.Equal(t, 1,<-ch1)
	require.Equal(t, 1,<-ch2)
}

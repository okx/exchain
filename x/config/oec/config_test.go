package oec

import (
	"github.com/okex/exchain/x/config/cm"
	"github.com/okex/exchain/x/config/tm"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig(t *testing.T) {
	c := GetOecConfig()

	tm.SetConfig(c)
	cm.SetConfig(c)
	require.Equal(t, uint64(5000), cm.Config.GetMaxOpen())
	require.Equal(t, uint64(300), tm.Config.GetTpb())

	c.SetMaxOpen(3000)
	c.SetTpb(150)
	require.Equal(t, uint64(3000), cm.Config.GetMaxOpen())
	require.Equal(t, uint64(150), tm.Config.GetTpb())

}
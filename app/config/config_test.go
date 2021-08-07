package config

import (
	"github.com/stretchr/testify/require"
	"testing"

	cm "github.com/cosmos/cosmos-sdk/server/config"
	tm "github.com/tendermint/tendermint/config"
)

func TestConfig(t *testing.T) {
	c := GetOecConfig()

	tm.SetConfig(c)
	cm.SetConfig(c)
	require.Equal(t, int64(5000), cm.DynamicConfig.GetMaxOpen())
	require.Equal(t, int64(300), tm.DynamicConfig.GetTpb())

	c.SetMaxOpen(3000)
	c.SetTpb(150)
	require.Equal(t, int64(3000), cm.DynamicConfig.GetMaxOpen())
	require.Equal(t, int64(150), tm.DynamicConfig.GetTpb())
}

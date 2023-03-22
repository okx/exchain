package config

import (
	"testing"

	iavlconfig "github.com/okx/okbchain/libs/iavl/config"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"

	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	tm "github.com/okx/okbchain/libs/tendermint/config"
)

func TestConfig(t *testing.T) {
	c := GetOkbcConfig()

	tm.SetDynamicConfig(c)
	require.Equal(t, 0, tm.DynamicConfig.GetMempoolSize())

	c.SetMempoolSize(150)
	require.Equal(t, 150, tm.DynamicConfig.GetMempoolSize())

	iavlconfig.SetDynamicConfig(c)
	require.Equal(t, int64(100), iavlconfig.DynamicConfig.GetCommitGapHeight())

	c.SetCommitGapHeight(0)
	require.Equal(t, int64(100), iavlconfig.DynamicConfig.GetCommitGapHeight())

	c.SetCommitGapHeight(-1)
	require.Equal(t, int64(100), iavlconfig.DynamicConfig.GetCommitGapHeight())

	c.SetCommitGapHeight(10)
	require.Equal(t, int64(10), iavlconfig.DynamicConfig.GetCommitGapHeight())

	viper.SetDefault(server.FlagPruning, "nothing")
	c.SetCommitGapHeight(9)
	require.Equal(t, int64(1), iavlconfig.DynamicConfig.GetCommitGapHeight())
}

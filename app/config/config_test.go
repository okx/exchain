package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	cm "github.com/cosmos/cosmos-sdk/server/config"
	tm "github.com/tendermint/tendermint/config"
)

func TestConfig(t *testing.T) {
	c := GetOecConfig()

	tm.SetDynamicConfig(c)
	cm.SetDynamicConfig(c)
	require.Equal(t, false, cm.DynamicConfig.GetMempoolRecheck())
	require.Equal(t, 0, tm.DynamicConfig.GetMempoolSize())

	c.SetMempoolRecheck(true)
	c.SetMempoolSize(150)
	require.Equal(t, true, cm.DynamicConfig.GetMempoolRecheck())
	require.Equal(t, 150, tm.DynamicConfig.GetMempoolSize())
}

func TestApollo(t *testing.T) {
	oecConf := GetOecConfig()
	client := NewApolloClient(oecConf)
	client.LoadConfig()
}

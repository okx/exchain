package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	tm "github.com/tendermint/tendermint/config"
)

func TestConfig(t *testing.T) {
	c := GetOecConfig()

	tm.SetDynamicConfig(c)
	require.Equal(t, 0, tm.DynamicConfig.GetMempoolSize())

	c.SetMempoolSize(150)
	require.Equal(t, 150, tm.DynamicConfig.GetMempoolSize())
}

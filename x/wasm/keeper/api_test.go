package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddrValiate(t *testing.T) {
	bz, cost, err := canonicalAddress("0xbbe4733d85bc2b90682147779da49cab38c0aa1f")
	require.NoError(t, err)
	t.Log("canonicalAddress cost", cost)
	addr, cost, err := humanAddress(bz)
	t.Log("humanAddress cost", cost)
	t.Log("humanAddress addr", addr)
}

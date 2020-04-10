package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {
	ctx, _, k, _, _ := CreateTestInputDefault(t, false, 1000)
	require.NotNil(t, ctx)
	require.Equal(t, auth.FeeCollectorName, k.GetFeeCollectorName())
}

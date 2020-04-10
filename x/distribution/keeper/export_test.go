package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {
	_, _, k, _, _ := CreateTestInputDefault(t, false, 1000)

	require.NotNil(t, k.GetCdc())
	require.Equal(t, types.DefaultCodespace, k.GetCodespace())
	require.NotNil(t, k.GetParamSpace())
	require.NotNil(t, k.GetStoreKey())
	require.NotNil(t, k.GetStakingKeeper())
	require.NotNil(t, k.GetSupplyKeeper())
	require.Equal(t, auth.FeeCollectorName, k.GetFeeCollectorName())
	require.NotNil(t, k.GetBlackListedAddrs())
}

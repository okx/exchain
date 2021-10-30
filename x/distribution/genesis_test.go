package distribution

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

type testParam struct {
	commission sdk.SysCoins
}

func getTestParams() []testParam {
	return []testParam{
		{keeper.NewTestSysCoins(1000, 0)},
		{keeper.NewTestSysCoins(150, 2)},
		{keeper.NewTestSysCoins(50, 8)},
		{nil},
	}
}

// InitGenesis sets distribution information for genesis
func TestInitGenesis(t *testing.T) {
	tests := getTestParams()
	length := len(tests)
	ctx, _, k, _, supplyKeeper := keeper.CreateTestInputDefault(t, false, 1000)

	valOpAddrs, _, valConsAddrs := keeper.GetTestAddrs()
	dwis := make([]DelegatorWithdrawInfo, length)
	accs := make([]ValidatorAccumulatedCommissionRecord, length)
	for i, valAddr := range valOpAddrs {
		accs[i].ValidatorAddress = valAddr
		accs[i].Accumulated = tests[i].commission
		dwis[i].DelegatorAddress, dwis[i].WithdrawAddress = keeper.TestAddrs[i*2], keeper.TestAddrs[i*2+1]
	}

	genesisState := NewGenesisState(types.DefaultParams(), types.InitialFeePool(), dwis, valConsAddrs[0], accs)
	InitGenesis(ctx, k, supplyKeeper, genesisState)
	require.True(t, k.GetFeePoolCommunityCoins(ctx).IsZero())
	require.Equal(t, genesisState.Params.CommunityTax, k.GetCommunityTax(ctx))
	require.Equal(t, genesisState.Params.WithdrawAddrEnabled, k.GetWithdrawAddrEnabled(ctx))
	require.Equal(t, genesisState.PreviousProposer, k.GetPreviousProposerConsAddr(ctx))
	for i := range accs {
		require.Equal(t, genesisState.DelegatorWithdrawInfos[i].WithdrawAddress,
			k.GetDelegatorWithdrawAddr(ctx, dwis[i].DelegatorAddress))
		require.Equal(t, genesisState.ValidatorAccumulatedCommissions[i].Accumulated,
			k.GetValidatorAccumulatedCommission(ctx, accs[i].ValidatorAddress))
		require.Equal(t, tests[i].commission,
			k.GetValidatorAccumulatedCommission(ctx, accs[i].ValidatorAddress))
	}

	actualGenesis := ExportGenesis(ctx, k)
	require.Equal(t, genesisState.Params.CommunityTax, actualGenesis.Params.CommunityTax)
	require.Equal(t, genesisState.Params.WithdrawAddrEnabled, actualGenesis.Params.WithdrawAddrEnabled)
	require.ElementsMatch(t, genesisState.DelegatorWithdrawInfos, actualGenesis.DelegatorWithdrawInfos)
	require.Equal(t, genesisState.PreviousProposer, actualGenesis.PreviousProposer)
	require.ElementsMatch(t, genesisState.ValidatorAccumulatedCommissions, actualGenesis.ValidatorAccumulatedCommissions)
}

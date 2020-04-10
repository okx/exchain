package distribution

import (
	"testing"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// InitGenesis sets distribution information for genesis
func TestInitGenesis(t *testing.T) {
	ctx, _, keeper, _, supplyKeeper := CreateTestInputDefault(t, false, 1000)

	dwis := make([]DelegatorWithdrawInfo, 2)
	dwis[0].DelegatorAddress = TestAddrs[0]
	dwis[0].WithdrawAddress = TestAddrs[1]
	dwis[1].DelegatorAddress = TestAddrs[2]
	dwis[1].WithdrawAddress = TestAddrs[3]

	pp := sdk.ConsAddress(ed25519.GenPrivKey().PubKey().Address())

	acc := make([]ValidatorAccumulatedCommissionRecord, 2)
	acc[0].ValidatorAddress = sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address())
	acc[0].Accumulated = sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 123)}
	acc[1].ValidatorAddress = sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address())
	acc[1].Accumulated = sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 456)}

	genesisState := NewGenesisState(true, dwis, pp, acc)
	InitGenesis(ctx, keeper, supplyKeeper, genesisState)
	require.Equal(t, genesisState.WithdrawAddrEnabled,
		keeper.GetWithdrawAddrEnabled(ctx))
	require.Equal(t, genesisState.DelegatorWithdrawInfos[0].WithdrawAddress,
		keeper.GetDelegatorWithdrawAddr(ctx, dwis[0].DelegatorAddress))
	require.Equal(t, genesisState.DelegatorWithdrawInfos[1].WithdrawAddress,
		keeper.GetDelegatorWithdrawAddr(ctx, dwis[1].DelegatorAddress))
	require.Equal(t, genesisState.PreviousProposer,
		keeper.GetPreviousProposerConsAddr(ctx))
	require.Equal(t, genesisState.ValidatorAccumulatedCommissions[0].Accumulated,
		keeper.GetValidatorAccumulatedCommission(ctx, acc[0].ValidatorAddress))
	require.Equal(t, genesisState.ValidatorAccumulatedCommissions[1].Accumulated,
		keeper.GetValidatorAccumulatedCommission(ctx, acc[1].ValidatorAddress))

	actualGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, genesisState.WithdrawAddrEnabled, actualGenesis.WithdrawAddrEnabled)
	require.ElementsMatch(t, genesisState.DelegatorWithdrawInfos, actualGenesis.DelegatorWithdrawInfos)
	require.Equal(t, genesisState.PreviousProposer, actualGenesis.PreviousProposer)
	require.ElementsMatch(t, genesisState.ValidatorAccumulatedCommissions, actualGenesis.ValidatorAccumulatedCommissions)
}

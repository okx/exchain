package v0_36

import (
	"testing"

	"github.com/okex/okchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v034distr "github.com/okex/okchain/x/distribution/legacy/v0_34"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestMigrate(t *testing.T) {
	delPk1, delPk2 := ed25519.GenPrivKey().PubKey(), ed25519.GenPrivKey().PubKey()
	delAddr1, delAddr2 := sdk.AccAddress(delPk1.Address()), sdk.AccAddress(delPk2.Address())
	valConsPk := ed25519.GenPrivKey().PubKey()
	valConsAddr := sdk.ConsAddress(valConsPk.Address())
	valOpPk1 := ed25519.GenPrivKey().PubKey()
	valOpAddr := sdk.ValAddress(valOpPk1.Address())
	accmulatedComssion := sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(125, 4))

	oldGenesis := v034distr.GenesisState{
		WithdrawAddrEnabled: true,
		DelegatorWithdrawInfos: []v034distr.DelegatorWithdrawInfo{
			{DelegatorAddress: delAddr1, WithdrawAddress: delAddr2},
		},
		PreviousProposer: valConsAddr,
		ValidatorAccumulatedCommissions: []v034distr.ValidatorAccumulatedCommissionRecord{
			{ValidatorAddress: valOpAddr, Accumulated: accmulatedComssion},
		},
	}

	var genesisState GenesisState
	require.NotPanics(t, func() {
		genesisState = Migrate(oldGenesis)
	})
	require.Equal(t, true, genesisState.WithdrawAddrEnabled)
	require.Equal(t, delAddr1, genesisState.DelegatorWithdrawInfos[0].DelegatorAddress)
	require.Equal(t, delAddr2, genesisState.DelegatorWithdrawInfos[0].WithdrawAddress)
	require.Equal(t, valConsAddr, genesisState.PreviousProposer)
	require.Equal(t, valOpAddr, genesisState.ValidatorAccumulatedCommissions[0].ValidatorAddress)
	require.Equal(t, accmulatedComssion, genesisState.ValidatorAccumulatedCommissions[0].Accumulated)
	require.Equal(t, sdk.NewDecWithPrec(2,2), genesisState.CommunityTax)
	require.Equal(t, true, genesisState.FeePool.CommunityPool.IsZero())
}

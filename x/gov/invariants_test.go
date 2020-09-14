package gov

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/gov/keeper"
	"github.com/okex/okexchain/x/gov/types"
)

func TestModuleAccountInvariant(t *testing.T) {
	ctx, _, gk, _, crisisKeeper := keeper.CreateTestInput(t, false, 1000)
	govHandler := NewHandler(gk)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 50)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, initialDeposit, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newDepositMsg := NewMsgDeposit(keeper.Addrs[0], proposalID,
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)})
	res = govHandler(ctx, newDepositMsg)
	require.True(t, res.IsOK())

	invariant := ModuleAccountInvariant(gk)
	_, broken := invariant(ctx)
	require.False(t, broken)

	// todo: check diff after RegisterInvariants
	RegisterInvariants(&crisisKeeper, gk)
}

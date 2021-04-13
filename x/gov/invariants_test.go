package gov

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/x/gov/keeper"
	"github.com/okex/exchain/x/gov/types"
)

func TestModuleAccountInvariant(t *testing.T) {
	ctx, _, gk, _, crisisKeeper := keeper.CreateTestInput(t, false, 1000)
	govHandler := NewHandler(gk)

	initialDeposit := sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 50)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, initialDeposit, keeper.Addrs[0])
	res, err := govHandler(ctx, newProposalMsg)
	require.Nil(t, err)
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newDepositMsg := NewMsgDeposit(keeper.Addrs[0], proposalID,
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)})
	res, err = govHandler(ctx, newDepositMsg)
	require.Nil(t, err)

	invariant := ModuleAccountInvariant(gk)
	_, broken := invariant(ctx)
	require.False(t, broken)

	// todo: check diff after RegisterInvariants
	RegisterInvariants(&crisisKeeper, gk)
}

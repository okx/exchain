package ut

import (
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/gov"
	"github.com/okx/okbchain/x/gov/types"
	"github.com/stretchr/testify/require"
)

func TestModuleAccountInvariant(t *testing.T) {
	ctx, _, gk, _, crisisKeeper := CreateTestInput(t, false, 1000)
	govHandler := gov.NewHandler(gk)

	initialDeposit := sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 50)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := gov.NewMsgSubmitProposal(content, initialDeposit, Addrs[0])
	res, err := govHandler(ctx, newProposalMsg)
	require.Nil(t, err)
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newDepositMsg := gov.NewMsgDeposit(Addrs[0], proposalID,
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)})
	res, err = govHandler(ctx, newDepositMsg)
	require.Nil(t, err)

	invariant := gov.ModuleAccountInvariant(gk)
	_, broken := invariant(ctx)
	require.False(t, broken)

	// todo: check diff after RegisterInvariants
	gov.RegisterInvariants(&crisisKeeper, gk)
}

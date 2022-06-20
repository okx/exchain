package distribution

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	govtypes "github.com/okex/exchain/x/gov/types"

	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())

	amount = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)))
)

func makeCommunityPoolSendProposal(recipient sdk.AccAddress, amount sdk.Coins) govtypes.Proposal {
	return govtypes.Proposal{Content: types.NewCommunityPoolSpendProposal(
		"Test",
		"description",
		recipient,
		amount,
	)}
}

func TestCommunityPoolSendProposalHandlerPassed(t *testing.T) {
	ctx, accountKeeper, k, _, supplyKeeper := keeper.CreateTestInputDefault(t, false, 10)
	recipient := delAddr1

	// add coins to the module account
	macc := k.GetDistributionAccount(ctx)
	err := macc.SetCoins(macc.GetCoins().Add(amount...))
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, macc)

	account := accountKeeper.NewAccountWithAddress(ctx, recipient)
	require.True(t, account.GetCoins().IsZero())
	accountKeeper.SetAccount(ctx, account)

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = sdk.NewCoins(amount...)
	k.SetFeePool(ctx, feePool)

	tp := makeCommunityPoolSendProposal(recipient, amount)
	hdlr := NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &tp))
	require.Equal(t, accountKeeper.GetAccount(ctx, recipient).GetCoins(), amount)
}

func TestCommunityPoolSendProposalHandlerFailed(t *testing.T) {
	ctx, accountKeeper, k, _, _ := keeper.CreateTestInputDefault(t, false, 10)
	recipient := delAddr1

	account := accountKeeper.NewAccountWithAddress(ctx, recipient)
	require.True(t, account.GetCoins().IsZero())
	accountKeeper.SetAccount(ctx, account)

	tp := makeCommunityPoolSendProposal(recipient, amount)
	hdlr := NewDistributionProposalHandler(k)
	err := hdlr(ctx, &tp)
	fmt.Println(err)
	require.Error(t, err)
	require.True(t, accountKeeper.GetAccount(ctx, recipient).GetCoins().IsZero())
}

func makeChangeDistributionTypeProposal(distrType uint32) govtypes.Proposal {
	return govtypes.Proposal{Content: types.NewChangeDistributionTypeProposal(
		"Test",
		"description",
		distrType,
	)}
}

func TestChangeDistributionTypeProposalHandlerPassed(t *testing.T) {
	ctx, _, k, _, _ := keeper.CreateTestInputDefault(t, false, 10)

	//init status, distribution off chain
	queryDistrType := k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOffChain)

	//set same type
	proposal := makeChangeDistributionTypeProposal(types.DistributionTypeOffChain)
	hdlr := NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOffChain)

	//set diff type, first
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOnChain)

	//set diff type, second
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOffChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOffChain)

	//set diff type, third
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOnChain)

	//set same type
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOnChain)
}

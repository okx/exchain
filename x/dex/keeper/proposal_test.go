package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/dex/types"
	govTypes "github.com/okex/okexchain/x/gov/types"
	ordertypes "github.com/okex/okexchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetMinDeposit(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx

	p := types.Params{
		DelistMinDeposit: sdk.DecCoins{sdk.NewDecCoin(common.NativeToken, sdk.NewInt(12345))},
	}

	testInput.DexKeeper.SetParams(ctx, p)
	var contentImpl types.DelistProposal
	minDeposit := testInput.DexKeeper.GetMinDeposit(ctx, contentImpl)
	require.True(t, minDeposit.IsEqual(p.DelistMinDeposit))
}

func TestKeeper_GetMaxDepositPeriod(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx

	p := types.Params{
		DelistMaxDepositPeriod: time.Second * 123,
	}
	testInput.DexKeeper.SetParams(ctx, p)
	var contentImpl types.DelistProposal
	maxDepositPeriod := testInput.DexKeeper.GetMaxDepositPeriod(ctx, contentImpl)
	require.EqualValues(t, maxDepositPeriod, p.DelistMaxDepositPeriod)
}

func TestKeeper_GetVotingPeriod(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx

	p := types.Params{
		DelistVotingPeriod: time.Second * 123,
	}
	testInput.DexKeeper.SetParams(ctx, p)
	var contentImpl types.DelistProposal
	maxListVotingPeriod := testInput.DexKeeper.GetVotingPeriod(ctx, contentImpl)
	require.EqualValues(t, maxListVotingPeriod, p.DelistVotingPeriod)
}

func TestKeeper_CheckMsgSubmitProposal(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx

	testInput.DexKeeper.SetParams(ctx, *types.DefaultParams())
	tokenPair := GetBuiltInTokenPair()

	content := types.NewDelistProposal("delist xxb_okb", "delist asset from dex", tokenPair.Owner, tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	content.Proposer = tokenPair.Owner
	proposal := govTypes.NewMsgSubmitProposal(content, sdk.DecCoins{sdk.NewDecCoin(common.NativeToken, sdk.NewInt(150))}, tokenPair.Owner)

	// error case : fail to check proposal because product(token pair) not exist
	err := testInput.DexKeeper.CheckMsgSubmitProposal(ctx, proposal)
	require.Error(t, err)
	// SaveTokenPair
	saveErr := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, saveErr)

	// successful case : check proposal successfully
	err = testInput.DexKeeper.CheckMsgSubmitProposal(ctx, proposal)
	require.NoError(t, err)

	// error case:  fail to check proposal because the proposer can't afford the initial deposit
	proposal1 := govTypes.NewMsgSubmitProposal(content, sdk.DecCoins{sdk.NewDecCoin(common.NativeToken, sdk.NewInt(500000))}, tokenPair.Owner)
	err = testInput.DexKeeper.CheckMsgSubmitProposal(ctx, proposal1)
	require.Error(t, err)

	// error case: fail to check proposal because initial deposit must not be less than 100.00000000okb
	proposal2 := govTypes.NewMsgSubmitProposal(content, sdk.DecCoins{sdk.NewDecCoin(common.NativeToken, sdk.NewInt(1))}, tokenPair.Owner)
	err = testInput.DexKeeper.CheckMsgSubmitProposal(ctx, proposal2)
	require.Error(t, err)

	// error case: fail to check proposal because the proposal is nil
	proposal3 := govTypes.NewMsgSubmitProposal(nil, sdk.DecCoins{sdk.NewDecCoin(common.NativeToken, sdk.NewInt(1))}, tokenPair.Owner)
	err = testInput.DexKeeper.CheckMsgSubmitProposal(ctx, proposal3)
	require.Error(t, err)

}

func TestKeeper_RejectedHandler(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx

	testInput.DexKeeper.SetParams(ctx, *types.DefaultParams())
	tokenPair := GetBuiltInTokenPair()

	// SaveTokenPair
	saveErr := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, saveErr)

	content := types.NewDelistProposal("delist xxb_okb", "delist asset from dex", tokenPair.Owner, tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	testInput.DexKeeper.RejectedHandler(ctx, content)

}

func TestKeeper_AfterDepositPeriodPassed(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	testInput.DexKeeper.SetParams(ctx, *types.DefaultParams())
	tokenPair := GetBuiltInTokenPair()

	// SaveTokenPair
	saveErr := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, saveErr)

	content := types.NewDelistProposal("delist xxb_okb", "delist asset from dex", tokenPair.Owner, tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	content.Proposer = tokenPair.Owner
	proposal := govTypes.Proposal{Content: content}

	testInput.DexKeeper.AfterDepositPeriodPassed(ctx, proposal)

}

func TestKeeper_AfterSubmitProposalHandler(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	content := types.NewDelistProposal("delist xxb_okb", "delist asset from dex", nil, "", "")
	proposal := govTypes.Proposal{Content: content}

	testInput.DexKeeper.AfterSubmitProposalHandler(ctx, proposal)
}

func TestKeeper_VoteHandler(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	tokenPair := GetBuiltInTokenPair()

	content := types.NewDelistProposal("delist xxb_okb", "delist asset from dex", nil, tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	proposal := govTypes.Proposal{Content: content}
	_, err := testInput.DexKeeper.VoteHandler(ctx, proposal, govTypes.Vote{})
	require.Nil(t, err)

	// SaveTokenPair
	saveErr := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, saveErr)

	lock := ordertypes.ProductLock{}
	testInput.DexKeeper.LockTokenPair(ctx, tokenPair.Name(), &lock)
	_, err = testInput.DexKeeper.VoteHandler(ctx, proposal, govTypes.Vote{})
	require.NotNil(t, err)

}

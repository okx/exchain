package dex

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/exchain/x/dex/types"
	govTypes "github.com/okex/exchain/x/gov/types"
	ordertypes "github.com/okex/exchain/x/order/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestProposal_NewProposalHandler(t *testing.T) {

	fakeTokenKeeper := newMockTokenKeeper()
	fakeSupplyKeeper := newMockSupplyKeeper()

	mApp, mDexKeeper, err := newMockApp(fakeTokenKeeper, fakeSupplyKeeper, 10)
	require.True(t, err == nil)

	mApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mApp.BaseApp.NewContext(false, abci.Header{})

	proposalHandler := NewProposalHandler(mDexKeeper.Keeper)

	params := types.DefaultParams()
	require.NotNil(t, params.String())
	mDexKeeper.SetParams(ctx, *params)
	tokenPair := GetBuiltInTokenPair()

	content := types.NewDelistProposal("delist xxb_okb", "delist asset from dex",
		tokenPair.Owner, tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	content.Proposer = tokenPair.Owner
	proposal := govTypes.Proposal{Content: content}

	// error case : fail to handle proposal because product(token pair) not exist
	err = proposalHandler(ctx, &proposal)
	require.Error(t, err)

	// error case : fail to handle proposal because proposal not exist
	err = proposalHandler(ctx, &govTypes.Proposal{})
	require.Error(t, err)

	// save wrong tokenpair
	tokenPair.Deposits = sdk.NewDecCoin("xxb", sdk.NewInt(50))
	saveErr := mApp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, saveErr)

	// error case : fail to withdraw deposits because deposits is not okt
	err = proposalHandler(ctx, &proposal)
	require.Error(t, err)

	// save right tokenpair
	tokenPair.Deposits = sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50))
	saveErr = mApp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, saveErr)

	// successful case : withdraw successfully
	err = proposalHandler(ctx, &proposal)
	require.Nil(t, err)

	lock := ordertypes.ProductLock{}
	mDexKeeper.LockTokenPair(ctx, ordertypes.TestTokenPair, &lock)
	err = proposalHandler(ctx, &proposal)
	require.Error(t, err)

}

package dex

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/dex/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func getMockTestCaseEvn(t *testing.T) (mApp *mockApp,
	tkKeeper *mockTokenKeeper, spKeeper *mockSupplyKeeper, dexKeeper *mockDexKeeper, testContext sdk.Context) {
	fakeTokenKeeper := newMockTokenKeeper()
	fakeSupplyKeeper := newMockSupplyKeeper()

	mApp, mockDexKeeper, err := newMockApp(fakeTokenKeeper, fakeSupplyKeeper, 10)
	require.True(t, err == nil)

	mApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mApp.BaseApp.NewContext(false, abci.Header{})

	return mApp, fakeTokenKeeper, fakeSupplyKeeper, mockDexKeeper, ctx
}

func TestHandler_HandleMsgList(t *testing.T) {
	mApp, tkKeeper, spKeeper, mDexKeeper, ctx := getMockTestCaseEvn(t)

	address := mApp.GenesisAccounts[0].GetAddress()
	listMsg := NewMsgList(address, "btc", "okt", sdk.NewDec(10))

	handlerFunctor := NewHandler(mApp.dexKeeper)

	// fail case : failed to list because token is invalid
	tkKeeper.exist = false
	badResult := handlerFunctor(ctx, listMsg)
	require.True(t, badResult.Code != sdk.CodeOK)

	// fail case : failed to list because tokenpair has been exist
	tkKeeper.exist = true
	badResult = handlerFunctor(ctx, listMsg)
	require.True(t, badResult.Code != sdk.CodeOK)
	require.True(t, badResult.Events == nil)

	// fail case : failed to list because SendCoinsFromModuleToAccount return error
	tkKeeper.exist = true
	mDexKeeper.getFakeTokenPair = false
	spKeeper.behaveEvil = true
	badResult = handlerFunctor(ctx, listMsg)
	require.True(t, badResult.Code != sdk.CodeOK)

	// successful case
	tkKeeper.exist = true
	spKeeper.behaveEvil = false
	mDexKeeper.getFakeTokenPair = false
	goodResult := handlerFunctor(ctx, listMsg)
	require.True(t, goodResult.Code == sdk.CodeOK)
	require.True(t, goodResult.Events != nil)
}

func TestHandler_HandleMsgDeList(t *testing.T) {
	mApp, _, spKeeper, mDexKeeper, ctx := getMockTestCaseEvn(t)
	address := mApp.GenesisAccounts[0].GetAddress()
	builtInTP := GetBuiltInTokenPair()
	delistMsg := NewMsgDelist(address, builtInTP.Name())

	handlerFunctor := NewHandler(mApp.dexKeeper)

	// Case1: product(trading pair) not exist
	mDexKeeper.getFakeTokenPair = false
	badResult := handlerFunctor(ctx, delistMsg)
	require.True(t, badResult.Code != sdk.CodeOK)

	// Case1.1: product(trading pair) not exist
	mDexKeeper.getFakeTokenPair = true
	delistMsg.Owner = address
	badResult = handlerFunctor(ctx, delistMsg)
	require.True(t, badResult.Code != sdk.CodeOK)

	// Case2: product exists,
	//        product belong to owner in MsgDelist, but
	//        no sufficient money to delist
	mDexKeeper.getFakeTokenPair = true
	delistMsg.Owner = builtInTP.Owner

	mDexKeeper.getFakeTokenPair = true
	delistMsg.Owner = builtInTP.Owner
	spKeeper.behaveEvil = false
	mDexKeeper.failToDeleteTokenPair = true
	badResult = handlerFunctor(ctx, delistMsg)
	require.True(t, badResult.Code == sdk.CodeOK)
	require.True(t, badResult.Events != nil)

	// Case 4: failed to delist product which is under delisting
	builtInTP.Delisting = true
	err := mDexKeeper.SaveTokenPair(ctx, builtInTP)
	require.Nil(t, err)
	mDexKeeper.getFakeTokenPair = false
	badResult = handlerFunctor(ctx, delistMsg)
	require.True(t, badResult.Code != sdk.CodeOK)

	// Case4: product exists,
	//        product belong to owner in MsgDelist,
	//        sufficient money to delist,
	//        deleteTokenPairByName successfully
	mDexKeeper.getFakeTokenPair = true
	delistMsg.Owner = builtInTP.Owner
	spKeeper.behaveEvil = false
	mDexKeeper.failToDeleteTokenPair = false
	goodResult := handlerFunctor(ctx, delistMsg)
	require.True(t, goodResult.Code == sdk.CodeOK)
	require.True(t, goodResult.Events != nil)

}

func TestHandler_HandleMsgDeposit(t *testing.T) {
	mApp, _, _, mDexKeeper, ctx := getMockTestCaseEvn(t)
	builtInTP := GetBuiltInTokenPair()
	depositMsg := NewMsgDeposit(builtInTP.Name(),
		sdk.NewDecCoin(builtInTP.QuoteAssetSymbol, sdk.NewInt(100)), builtInTP.Owner)

	handlerFunctor := NewHandler(mApp.dexKeeper)

	// Case1: failed to deposit
	mDexKeeper.failToDeposit = true
	bad1 := handlerFunctor(ctx, depositMsg)
	require.True(t, bad1.Code != sdk.CodeOK)

	// Case2: success to deposit
	mDexKeeper.failToDeposit = false
	good1 := handlerFunctor(ctx, depositMsg)
	require.True(t, good1.Code == sdk.CodeOK)
	require.True(t, good1.Events != nil)
}

func TestHandler_HandleMsgWithdraw(t *testing.T) {
	mApp, _, _, mDexKeeper, ctx := getMockTestCaseEvn(t)
	builtInTP := GetBuiltInTokenPair()
	withdrawMsg := NewMsgWithdraw(builtInTP.Name(),
		sdk.NewDecCoin(builtInTP.QuoteAssetSymbol, sdk.NewInt(100)), builtInTP.Owner)

	handlerFunctor := NewHandler(mApp.dexKeeper)

	// Case1: failed to deposit
	mDexKeeper.failToWithdraw = true
	bad1 := handlerFunctor(ctx, withdrawMsg)
	require.True(t, bad1.Code != sdk.CodeOK)

	// Case2: success to deposit
	mDexKeeper.failToWithdraw = false
	good1 := handlerFunctor(ctx, withdrawMsg)
	require.True(t, good1.Code == sdk.CodeOK)
	require.True(t, good1.Events != nil)
}

func TestHandler_HandleMsgBad(t *testing.T) {
	mApp, _, _, _, ctx := getMockTestCaseEvn(t)
	handlerFunctor := NewHandler(mApp.dexKeeper)

	res := handlerFunctor(ctx, sdk.NewTestMsg())
	require.False(t, res.Code.IsOK())
}

func TestHandler_handleMsgTransferOwnership(t *testing.T) {
	mApp, _, spKeeper, mDexKeeper, ctx := getMockTestCaseEvn(t)

	tokenPair := GetBuiltInTokenPair()
	err := mDexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	handlerFunctor := NewHandler(mApp.dexKeeper)
	to := mApp.GenesisAccounts[0].GetAddress()

	// successful case
	msgTransferOwnership := types.NewMsgTransferOwnership(tokenPair.Owner, to, tokenPair.Name())
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgTransferOwnership)

	// fail case : failed to TransferOwnership because product is not exist
	msgFailedTransferOwnership := types.NewMsgTransferOwnership(tokenPair.Owner, to, "no-product")
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgFailedTransferOwnership)

	// fail case : failed to SendCoinsFromModuleToAccount return error
	msgFailedTransferOwnership = types.NewMsgTransferOwnership(tokenPair.Owner, to, tokenPair.Name())
	spKeeper.behaveEvil = true
	handlerFunctor(ctx, msgFailedTransferOwnership)
}

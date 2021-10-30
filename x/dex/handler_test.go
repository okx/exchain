package dex

import (
	"github.com/okex/exchain/x/common"
	"testing"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/dex/types"
	"github.com/stretchr/testify/require"
	abci "github.com/okex/exchain/dependence/tendermint/abci/types"
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
	listMsg := NewMsgList(address, "btc", common.NativeToken, sdk.NewDec(10))
	mDexKeeper.SetOperator(ctx, types.DEXOperator{Address: address, HandlingFeeAddress: address})

	handlerFunctor := NewHandler(mApp.dexKeeper)

	// fail case : failed to list because token is invalid
	tkKeeper.exist = false
	_, err := handlerFunctor(ctx, listMsg)
	require.NotNil(t, err)

	// fail case : failed to list because tokenpair has been exist
	tkKeeper.exist = true
	_, err = handlerFunctor(ctx, listMsg)
	require.NotNil(t, err)

	// fail case : failed to list because SendCoinsFromModuleToAccount return error
	tkKeeper.exist = true
	mDexKeeper.getFakeTokenPair = false
	spKeeper.behaveEvil = true
	_, err = handlerFunctor(ctx, listMsg)
	require.NotNil(t, err)

	// successful case
	tkKeeper.exist = true
	spKeeper.behaveEvil = false
	mDexKeeper.getFakeTokenPair = false
	goodResult, err := handlerFunctor(ctx, listMsg)
	require.Nil(t, err)
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
	_, err := handlerFunctor(ctx, depositMsg)
	require.NotNil(t, err)

	// Case2: success to deposit
	mDexKeeper.failToDeposit = false
	good1, err := handlerFunctor(ctx, depositMsg)
	require.Nil(t, err)
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
	_, err := handlerFunctor(ctx, withdrawMsg)
	require.NotNil(t, err)

	// Case2: success to deposit
	mDexKeeper.failToWithdraw = false
	good1, err := handlerFunctor(ctx, withdrawMsg)
	require.Nil(t, err)
	require.True(t, good1.Events != nil)
}

func TestHandler_HandleMsgBad(t *testing.T) {
	mApp, _, _, _, ctx := getMockTestCaseEvn(t)
	handlerFunctor := NewHandler(mApp.dexKeeper)

	_, err := handlerFunctor(ctx, sdk.NewTestMsg())
	require.NotNil(t, err)
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

	// fail case : failed to TransferOwnership because the address is not the owner of product
	msgFailedTransferOwnership = types.NewMsgTransferOwnership(to, to, tokenPair.Name())
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgFailedTransferOwnership)

	// confirm ownership successful case
	msgConfirmOwnership := types.NewMsgConfirmOwnership(to, tokenPair.Name())
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgConfirmOwnership)

	// fail case : failed to ConfirmOwnership because the address is not the new owner of product
	msgTransferOwnership = types.NewMsgTransferOwnership(to, tokenPair.Owner, tokenPair.Name())
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgTransferOwnership)
	msgFailedConfirmOwnership := types.NewMsgConfirmOwnership(to, tokenPair.Name())
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgFailedConfirmOwnership)

	// fail case : failed to ConfirmOwnership because there is not transfer-ownership list to confirm
	mDexKeeper.DeleteConfirmOwnership(ctx, tokenPair.Name())
	msgFailedConfirmOwnership = types.NewMsgConfirmOwnership(tokenPair.Owner, tokenPair.Name())
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgFailedConfirmOwnership)

	// fail case : failed to ConfirmOwnership because the product is not exist
	mDexKeeper.DeleteTokenPairByName(ctx, tokenPair.Owner, tokenPair.Name())
	msgFailedConfirmOwnership = types.NewMsgConfirmOwnership(tokenPair.Owner, tokenPair.Name())
	spKeeper.behaveEvil = false
	handlerFunctor(ctx, msgFailedConfirmOwnership)
}

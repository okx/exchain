package dex

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/dex/types"
	"github.com/stretchr/testify/require"
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

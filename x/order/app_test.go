package order

import (
	"fmt"
	"testing"

	"github.com/okex/exchain/x/common/monitor"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mock"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/okex/exchain/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/token"
)

type MockApp struct {
	*mock.App

	keyOrder     *sdk.KVStoreKey
	keyToken     *sdk.KVStoreKey
	keyLock      *sdk.KVStoreKey
	keyDex       *sdk.KVStoreKey
	keyTokenPair *sdk.KVStoreKey

	keySupply *sdk.KVStoreKey

	bankKeeper   bank.Keeper
	orderKeeper  Keeper
	tokenKeeper  token.Keeper
	supplyKeeper supply.Keeper
	dexKeeper    dex.Keeper
}

func registerCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	token.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
}

func getMockApp(t *testing.T, numGenAccs int) (mockApp *MockApp, addrKeysSlice mock.AddrKeysSlice) {
	return getMockAppWithBalance(t, numGenAccs, 100)
}

// initialize the mock application for this module
func getMockAppWithBalance(t *testing.T, numGenAccs int, balance int64) (mockApp *MockApp,
	addrKeysSlice mock.AddrKeysSlice) {
	mapp := mock.NewApp()
	registerCodec(mapp.Cdc)

	mockApp = &MockApp{
		App:      mapp,
		keyOrder: sdk.NewKVStoreKey(OrderStoreKey),

		keyToken:     sdk.NewKVStoreKey(token.StoreKey),
		keyLock:      sdk.NewKVStoreKey(token.KeyLock),
		keyDex:       sdk.NewKVStoreKey(dex.StoreKey),
		keyTokenPair: sdk.NewKVStoreKey(dex.TokenPairStoreKey),

		keySupply: sdk.NewKVStoreKey(supply.StoreKey),
	}

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.String()] = true

	mockApp.bankKeeper = bank.NewBaseKeeper(mockApp.AccountKeeper,
		mockApp.ParamsKeeper.Subspace(bank.DefaultParamspace),
		blacklistedAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		token.ModuleName:      {supply.Minter, supply.Burner},
	}
	mockApp.supplyKeeper = supply.NewKeeper(mockApp.Cdc, mockApp.keySupply, mockApp.AccountKeeper,
		mockApp.bankKeeper, maccPerms)

	mockApp.tokenKeeper = token.NewKeeper(
		mockApp.bankKeeper,
		mockApp.ParamsKeeper.Subspace(token.DefaultParamspace),
		auth.FeeCollectorName,
		mockApp.supplyKeeper,
		mockApp.keyToken,
		mockApp.keyLock,
		mockApp.Cdc,
		mockApp.AccountKeeper)

	mockApp.dexKeeper = dex.NewKeeper(
		auth.FeeCollectorName,
		mockApp.supplyKeeper,
		mockApp.ParamsKeeper.Subspace(dex.DefaultParamspace),
		mockApp.tokenKeeper,
		nil,
		mockApp.bankKeeper,
		mockApp.keyDex,
		mockApp.keyTokenPair,
		mockApp.Cdc)

	mockApp.orderKeeper = NewKeeper(
		mockApp.tokenKeeper,
		mockApp.supplyKeeper,
		mockApp.dexKeeper,
		mockApp.ParamsKeeper.Subspace(DefaultParamspace),
		auth.FeeCollectorName,
		mockApp.keyOrder,
		mockApp.Cdc,
		monitor.NopOrderMetrics())

	mockApp.Router().AddRoute(RouterKey, NewOrderHandler(mockApp.orderKeeper))
	mockApp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(mockApp.orderKeeper))

	mockApp.SetBeginBlocker(getBeginBlocker(mockApp.orderKeeper))
	mockApp.SetEndBlocker(getEndBlocker(mockApp.orderKeeper))
	mockApp.SetInitChainer(getInitChainer(mockApp.App, mockApp.supplyKeeper,
		[]exported.ModuleAccountI{feeCollector}))

	decCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		balance, common.NativeToken, balance, common.TestToken))
	require.Nil(t, err)
	coins := decCoins

	keysSlice, genAccs := CreateGenAccounts(numGenAccs, coins)
	addrKeysSlice = keysSlice

	// todo: checkTx in mock app
	mockApp.SetAnteHandler(nil)

	app := mockApp
	require.NoError(t, app.CompleteSetup(
		app.keyOrder,
		app.keyToken,
		app.keyDex,
		app.keyTokenPair,
		app.keyLock,
		app.keySupply,
	))
	mock.SetGenesis(mockApp.App, genAccs)

	for i := 0; i < numGenAccs; i++ {
		mock.CheckBalance(t, app.App, keysSlice[i].Address, coins)
		mockApp.TotalCoinsSupply = mockApp.TotalCoinsSupply.Add2(coins)
	}

	return mockApp, addrKeysSlice
}

func getBeginBlocker(keeper Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		BeginBlocker(ctx, keeper)
		return abci.ResponseBeginBlock{}
	}
}

func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{}
	}
}

func getInitChainer(mapp *mock.App, supplyKeeper types.SupplyKeeper,
	blacklistedAddrs []exported.ModuleAccountI) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}
		return abci.ResponseInitChain{}
	}
}

//func produceOrderTxs(app *MockApp, ctx sdk.Context, numToGenerate int, addrKeys mock.AddrKeys,
//	orderMsg *MsgNewOrder) []auth.StdTx {
//	txs := make([]auth.StdTx, numToGenerate)
//	orderMsg.Sender = addrKeys.Address
//	for i := 0; i < numToGenerate; i++ {
//		txs[i] = buildTx(app, ctx, addrKeys, *orderMsg)
//	}
//	return txs
//}

//func buildTx(app *MockApp, ctx sdk.Context, addrKeys mock.AddrKeys, msg sdk.Msg) auth.StdTx {
//	accs := app.AccountKeeper.GetAccount(ctx, addrKeys.Address)
//	accNum := accs.GetAccountNumber()
//	seqNum := accs.GetSequence()
//
//	tx := mock.GenTx(
//		[]sdk.Msg{msg}, []uint64{uint64(accNum)}, []uint64{uint64(seqNum)}, addrKeys.PrivKey)
//	res := app.Check(tx)
//	if !res.IsOK() {
//		panic(fmt.Sprintf("something wrong in checking transaction: %v", res))
//	}
//	return tx
//}
//
//func mockApplyBlock(app *MockApp, blockHeight int64, txs []auth.StdTx) {
//	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: blockHeight}})
//
//	newCtx := app.NewContext(false, abci.Header{})
//	param := DefaultTestParams()
//	app.orderKeeper.SetParams(newCtx, &param)
//	for _, tx := range txs {
//		app.Deliver(tx)
//	}
//
//	app.EndBlock(abci.RequestEndBlock{Height: blockHeight})
//	app.Commit()
//}

func CreateGenAccounts(numAccs int, genCoins sdk.Coins) (addrKeysSlice mock.AddrKeysSlice,
	genAccs []auth.Account) {
	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())

		addrKeys := mock.NewAddrKeys(addr, pubKey, privKey)
		account := &auth.BaseAccount{
			Address: addr,
			Coins:   genCoins,
		}
		genAccs = append(genAccs, account)
		addrKeysSlice = append(addrKeysSlice, addrKeys)
	}
	return
}

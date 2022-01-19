package ammswap

import (
	"fmt"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mock"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/okex/exchain/x/ammswap/types"
	staking "github.com/okex/exchain/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/x/token"
)

type MockApp struct {
	*mock.App

	keySwap   *sdk.KVStoreKey
	keyToken  *sdk.KVStoreKey
	keyLock   *sdk.KVStoreKey
	keySupply *sdk.KVStoreKey

	bankKeeper   bank.Keeper
	swapKeeper   Keeper
	tokenKeeper  token.Keeper
	supplyKeeper supply.Keeper
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
		App:       mapp,
		keySwap:   sdk.NewKVStoreKey(StoreKey),
		keyToken:  sdk.NewKVStoreKey(token.StoreKey),
		keyLock:   sdk.NewKVStoreKey(token.KeyLock),
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
		ModuleName:            {supply.Minter, supply.Burner},
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
		true, mockApp.AccountKeeper)

	mockApp.swapKeeper = NewKeeper(
		mockApp.supplyKeeper,
		mockApp.tokenKeeper,
		mockApp.Cdc,
		mockApp.keySwap,
		mockApp.ParamsKeeper.Subspace(DefaultParamspace),
	)

	mockApp.Router().AddRoute(RouterKey, NewHandler(mockApp.swapKeeper))
	mockApp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(mockApp.swapKeeper))

	mockApp.SetBeginBlocker(getBeginBlocker(mockApp.swapKeeper))
	mockApp.SetEndBlocker(getEndBlocker(mockApp.swapKeeper))
	mockApp.SetInitChainer(getInitChainer(mockApp.App, mockApp.supplyKeeper,
		[]exported.ModuleAccountI{feeCollector}))

	decCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s,%d%s,%d%s",
		balance, types.TestQuotePooledToken, balance, types.TestBasePooledToken, balance, types.TestBasePooledToken2, balance, types.TestBasePooledToken3))
	require.Nil(t, err)
	coins := decCoins

	keysSlice, genAccs := CreateGenAccounts(numGenAccs, coins)
	addrKeysSlice = keysSlice

	// todo: checkTx in mock app
	mockApp.SetAnteHandler(nil)

	app := mockApp
	require.NoError(t, app.CompleteSetup(
		app.keySwap,
		app.keyToken,
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

func getInitChainer(mapp *mock.App, supplyKeeper staking.SupplyKeeper,
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

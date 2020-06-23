package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okchain/x/poolswap/types"
	staking "github.com/okex/okchain/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/okex/okchain/x/token"
)

type TestInput struct {
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

func regCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
	token.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
}

func GetTestInput(t *testing.T, numGenAccs int) (mockApp *TestInput, addrKeysSlice mock.AddrKeysSlice) {
	return getTestInputWithBalance(t, numGenAccs, 100)
}

// initialize the mock application for this module
func getTestInputWithBalance(t *testing.T, numGenAccs int, balance int64) (mockApp *TestInput,
	addrKeysSlice mock.AddrKeysSlice) {
	mapp := mock.NewApp()
	regCodec(mapp.Cdc)

	mockApp = &TestInput{
		App:       mapp,
		keySwap:   sdk.NewKVStoreKey(types.StoreKey),
		keyToken:  sdk.NewKVStoreKey(token.StoreKey),
		keyLock:   sdk.NewKVStoreKey(token.KeyLock),
		keySupply: sdk.NewKVStoreKey(supply.StoreKey),
	}

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.String()] = true

	mockApp.bankKeeper = bank.NewBaseKeeper(mockApp.AccountKeeper,
		mockApp.ParamsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace, blacklistedAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		token.ModuleName:      {supply.Minter, supply.Burner},
		types.ModuleName:      {supply.Minter, supply.Burner},
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
		true)

	mockApp.swapKeeper = NewKeeper(
		mockApp.supplyKeeper,
		mockApp.tokenKeeper,
		mockApp.Cdc,
		mockApp.keySwap,
		mockApp.ParamsKeeper.Subspace(types.DefaultParamspace),
	)

	mockApp.QueryRouter().AddRoute(types.QuerierRoute, NewQuerier(mockApp.swapKeeper))

	mockApp.SetInitChainer(initChainer(mockApp.App, mockApp.supplyKeeper,
		[]exported.ModuleAccountI{feeCollector}))

	decCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s,%d%s,%d%s",
		balance, types.TestQuotePooledToken, balance, types.TestBasePooledToken, balance, types.TestBasePooledToken2, balance, types.TestBasePooledToken3))
	require.Nil(t, err)
	coins := decCoins

	keysSlice, genAccs := GenAccounts(numGenAccs, coins)
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
		mockApp.TotalCoinsSupply = mockApp.TotalCoinsSupply.Add(coins)
	}

	return mockApp, addrKeysSlice
}

func initChainer(mapp *mock.App, supplyKeeper staking.SupplyKeeper,
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

func GenAccounts(numAccs int, genCoins sdk.Coins) (addrKeysSlice mock.AddrKeysSlice,
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

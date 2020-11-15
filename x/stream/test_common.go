package stream

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/common/monitor"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/order"
	"github.com/okex/okexchain/x/order/keeper"
	ordertypes "github.com/okex/okexchain/x/order/types"
	stakingtypes "github.com/okex/okexchain/x/staking/types"
	"github.com/okex/okexchain/x/token"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type MockApp struct {
	*mock.App

	keyOrder *sdk.KVStoreKey

	keyToken     *sdk.KVStoreKey
	keyFreeze    *sdk.KVStoreKey
	keyLock      *sdk.KVStoreKey
	keyDex       *sdk.KVStoreKey
	keyTokenPair *sdk.KVStoreKey

	keySupply *sdk.KVStoreKey

	BankKeeper   bank.Keeper
	OrderKeeper  keeper.Keeper
	DexKeeper    dex.Keeper
	TokenKeeper  token.Keeper
	supplyKeeper supply.Keeper
	streamKeeper Keeper
}

func registerCodec(cdc *codec.Codec) {
	supply.RegisterCodec(cdc)
}
func GetMockApp(t *testing.T, numGenAccs int, cfg *appCfg.Config) (mockApp *MockApp, addrKeysSlice mock.AddrKeysSlice) {
	return getMockAppWithBalance(t, numGenAccs, 100, cfg)
}

// initialize the mock application for this module
func getMockAppWithBalance(t *testing.T, numGenAccs int, balance int64, cfg *appCfg.Config) (mockApp *MockApp, addrKeysSlice mock.AddrKeysSlice) {
	mapp := mock.NewApp()
	registerCodec(mapp.Cdc)

	mockApp = &MockApp{
		App:      mapp,
		keyOrder: sdk.NewKVStoreKey(ordertypes.OrderStoreKey),

		keyToken:     sdk.NewKVStoreKey("token"),
		keyFreeze:    sdk.NewKVStoreKey("freeze"),
		keyLock:      sdk.NewKVStoreKey("lock"),
		keyDex:       sdk.NewKVStoreKey(dex.StoreKey),
		keyTokenPair: sdk.NewKVStoreKey(dex.TokenPairStoreKey),

		keySupply: sdk.NewKVStoreKey(supply.StoreKey),
	}

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.String()] = true

	mockApp.BankKeeper = bank.NewBaseKeeper(mockApp.AccountKeeper, mockApp.ParamsKeeper.Subspace(bank.DefaultParamspace),
		blacklistedAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
	}
	mockApp.supplyKeeper = supply.NewKeeper(mockApp.Cdc, mockApp.keySupply, mockApp.AccountKeeper,
		mockApp.BankKeeper, maccPerms)

	mockApp.TokenKeeper = token.NewKeeper(
		mockApp.BankKeeper,
		mockApp.ParamsKeeper.Subspace(token.DefaultParamspace),
		auth.FeeCollectorName,
		mockApp.supplyKeeper,
		mockApp.keyToken,
		mockApp.keyLock,
		mockApp.Cdc,
		true)

	mockApp.DexKeeper = dex.NewKeeper(
		auth.FeeCollectorName,
		mockApp.supplyKeeper,
		mockApp.ParamsKeeper.Subspace(dex.DefaultParamspace),
		mockApp.TokenKeeper,
		nil,
		mockApp.BankKeeper,
		mockApp.keyDex,
		mockApp.keyTokenPair,
		mockApp.Cdc)

	mockApp.OrderKeeper = keeper.NewKeeper(
		mockApp.TokenKeeper,
		mockApp.supplyKeeper,
		mockApp.DexKeeper,
		mockApp.ParamsKeeper.Subspace(ordertypes.DefaultParamspace),
		auth.FeeCollectorName,
		mockApp.keyOrder,
		mockApp.Cdc,
		true,
		monitor.NopOrderMetrics())

	mockApp.streamKeeper = NewKeeper(mockApp.OrderKeeper, mockApp.TokenKeeper, &mockApp.DexKeeper, &mockApp.AccountKeeper, nil, mockApp.Cdc, mockApp.Logger(), cfg, monitor.NopStreamMetrics())

	mockApp.Router().AddRoute(ordertypes.RouterKey, order.NewOrderHandler(mockApp.OrderKeeper))
	mockApp.QueryRouter().AddRoute(order.QuerierRoute, keeper.NewQuerier(mockApp.OrderKeeper))

	mockApp.SetBeginBlocker(getBeginBlocker(mockApp))
	mockApp.SetEndBlocker(getEndBlocker(mockApp))
	mockApp.SetInitChainer(getInitChainer(mockApp.App, mockApp.supplyKeeper,
		[]exported.ModuleAccountI{feeCollector}))

	coins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		balance, common.NativeToken, balance, common.TestToken))
	if err != nil {
		panic(err)
	}

	keysSlice, genAccs := CreateGenAccounts(numGenAccs, coins)
	addrKeysSlice = keysSlice

	mockApp.SetAnteHandler(nil)

	app := mockApp
	mockApp.MountStores(
		app.keyOrder,
		app.keyToken,
		app.keyDex,
		app.keyTokenPair,
		app.keyFreeze,
		app.keyLock,
		app.keySupply,
	)

	require.NoError(t, mockApp.CompleteSetup(mockApp.keyOrder))
	mock.SetGenesis(mockApp.App, genAccs)
	return mockApp, addrKeysSlice
}

func getBeginBlocker(mapp *MockApp) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		order.BeginBlocker(ctx, mapp.OrderKeeper)
		return abci.ResponseBeginBlock{}
	}
}

func getEndBlocker(mapp *MockApp) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		order.EndBlocker(ctx, mapp.OrderKeeper)
		EndBlocker(ctx, mapp.streamKeeper)
		return abci.ResponseEndBlock{}
	}
}

func getInitChainer(mapp *mock.App, supplyKeeper stakingtypes.SupplyKeeper,
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

func ProduceOrderTxs(app *MockApp, ctx sdk.Context, numToGenerate int, addrKeys mock.AddrKeys,
	orderMsg *ordertypes.MsgNewOrders) []auth.StdTx {
	txs := make([]auth.StdTx, numToGenerate)
	orderMsg.Sender = addrKeys.Address
	for i := 0; i < numToGenerate; i++ {
		txs[i] = buildTx(app, ctx, addrKeys, *orderMsg)
	}
	return txs
}

func buildTx(app *MockApp, ctx sdk.Context, addrKeys mock.AddrKeys, msg sdk.Msg) auth.StdTx {
	accs := app.AccountKeeper.GetAccount(ctx, addrKeys.Address)
	accNum := accs.GetAccountNumber()
	seqNum := accs.GetSequence()

	tx := mock.GenTx([]sdk.Msg{msg}, []uint64{accNum}, []uint64{seqNum}, addrKeys.PrivKey)
	_, _, err := app.Check(tx)
	if err != nil {
		panic(fmt.Sprintf("something wrong in checking transaction: %v", err))
	}
	return tx
}

func MockApplyBlock(app *MockApp, blockHeight int64, txs []auth.StdTx) {
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: blockHeight}})

	newCtx := app.NewContext(false, abci.Header{})
	param := ordertypes.DefaultParams()
	app.OrderKeeper.SetParams(newCtx, &param)
	for _, tx := range txs {
		app.Deliver(tx)
	}

	app.EndBlock(abci.RequestEndBlock{Height: blockHeight})
	app.Commit()
}

func CreateGenAccounts(numAccs int, genCoins sdk.Coins) (addrKeysSlice mock.AddrKeysSlice, genAccs []auth.Account) {
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

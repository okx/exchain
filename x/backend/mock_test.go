package backend

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mock"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	types3 "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/ammswap"
	"github.com/okex/exchain/x/backend/client/cli"
	"github.com/okex/exchain/x/backend/config"
	"github.com/okex/exchain/x/backend/orm"
	"github.com/okex/exchain/x/backend/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/common/monitor"
	"github.com/okex/exchain/x/common/version"
	types2 "github.com/okex/exchain/x/dex/types"
	"github.com/okex/exchain/x/farm"
	"github.com/okex/exchain/x/order/keeper"
	ordertypes "github.com/okex/exchain/x/order/types"
	tokentypes "github.com/okex/exchain/x/token/types"

	//"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/order"

	//"github.com/okex/exchain/x/staking"
	"github.com/okex/exchain/x/token"
	"github.com/stretchr/testify/require"
)

type MockApp struct {
	*mock.App

	keyOrder *sdk.KVStoreKey

	keyToken     *sdk.KVStoreKey
	keyLock      *sdk.KVStoreKey
	keyDex       *sdk.KVStoreKey
	keyTokenPair *sdk.KVStoreKey

	keySupply *sdk.KVStoreKey

	bankKeeper    bank.Keeper
	orderKeeper   keeper.Keeper
	dexKeeper     dex.Keeper
	tokenKeeper   token.Keeper
	backendKeeper Keeper
	supplyKeeper  supply.Keeper
	swapKeeper    ammswap.Keeper
	keySwap       *sdk.KVStoreKey
	farmKeeper    farm.Keeper
	keyFarm       *sdk.KVStoreKey
}

func registerCdc(cdc *codec.Codec) {
	supply.RegisterCodec(cdc)
}

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int, enableBackend bool, dbDir string) (mockApp *MockApp, addrKeysSlice mock.AddrKeysSlice) {
	mapp := mock.NewApp()
	registerCdc(mapp.Cdc)

	mockApp = &MockApp{
		App:          mapp,
		keyOrder:     sdk.NewKVStoreKey(ordertypes.OrderStoreKey),
		keyToken:     sdk.NewKVStoreKey(tokentypes.ModuleName),
		keyLock:      sdk.NewKVStoreKey(tokentypes.KeyLock),
		keyDex:       sdk.NewKVStoreKey(dex.StoreKey),
		keyTokenPair: sdk.NewKVStoreKey(dex.TokenPairStoreKey),
		keySupply:    sdk.NewKVStoreKey(supply.StoreKey),
		keySwap:      sdk.NewKVStoreKey(ammswap.StoreKey),
		keyFarm:      sdk.NewKVStoreKey(farm.StoreKey),
	}

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.String()] = true

	mockApp.bankKeeper = bank.NewBaseKeeper(mockApp.AccountKeeper,
		mockApp.ParamsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)

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
		//mockApp.keyTokenPair,
		mockApp.Cdc,
		true, mockApp.AccountKeeper,
	)

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

	mockApp.orderKeeper = keeper.NewKeeper(
		mockApp.tokenKeeper,
		mockApp.supplyKeeper,
		mockApp.dexKeeper,
		mockApp.ParamsKeeper.Subspace(ordertypes.DefaultParamspace),
		auth.FeeCollectorName,
		mockApp.keyOrder,
		mockApp.Cdc,
		true,
		monitor.NopOrderMetrics())

	mockApp.swapKeeper = ammswap.NewKeeper(
		mockApp.supplyKeeper,
		mockApp.tokenKeeper,
		mockApp.Cdc,
		mockApp.keySwap,
		mockApp.ParamsKeeper.Subspace(ammswap.DefaultParamspace),
	)
	mockApp.farmKeeper = farm.NewKeeper(auth.FeeCollectorName, mockApp.supplyKeeper, mockApp.tokenKeeper,
		mockApp.swapKeeper, mockApp.ParamsKeeper.Subspace(farm.DefaultParamspace), mockApp.keyFarm, mockApp.Cdc,
	)
	// CleanUp data
	cfg, err := config.SafeLoadMaintainConfig(config.DefaultTestConfig)
	require.Nil(t, err)
	cfg.EnableBackend = enableBackend
	cfg.EnableMktCompute = enableBackend
	cfg.OrmEngine.EngineType = orm.EngineTypeSqlite
	cfg.OrmEngine.ConnectStr = config.DefaultTestDataHome + "/sqlite3/backend.db"
	if dbDir == "" {
		path := config.DefaultTestDataHome + "/sqlite3"
		if err := os.RemoveAll(path); err != nil {
			mockApp.Logger().Debug(err.Error())
		}
	} else {
		cfg.LogSQL = false
		cfg.OrmEngine.ConnectStr = dbDir + "/backend.db"
	}

	mockApp.backendKeeper = NewKeeper(
		mockApp.orderKeeper,
		mockApp.tokenKeeper,
		&mockApp.dexKeeper,
		&mockApp.swapKeeper,
		&mockApp.farmKeeper,
		nil,
		nil,
		mockApp.Cdc,
		mockApp.Logger(),
		cfg)

	mockApp.Router().AddRoute(ordertypes.RouterKey, order.NewOrderHandler(mockApp.orderKeeper))
	mockApp.QueryRouter().AddRoute(ordertypes.QuerierRoute, keeper.NewQuerier(mockApp.orderKeeper))
	//mockApp.Router().AddRoute(token.RouterKey, token.NewHandler(mockApp.tokenKeeper))
	mockApp.Router().AddRoute(token.RouterKey, token.NewTokenHandler(mockApp.tokenKeeper, version.ProtocolVersionV0))
	mockApp.QueryRouter().AddRoute(token.QuerierRoute, token.NewQuerier(mockApp.tokenKeeper))

	mockApp.SetEndBlocker(getEndBlocker(mockApp.orderKeeper, mockApp.backendKeeper))
	mockApp.SetInitChainer(getInitChainer(mockApp.App, mockApp.bankKeeper, mockApp.supplyKeeper,
		[]exported.ModuleAccountI{feeCollector}))

	intQuantity := 100000
	coins, _ := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		intQuantity, common.NativeToken, intQuantity, common.TestToken))

	keysSlice, genAccs := CreateGenAccounts(numGenAccs, coins)
	addrKeysSlice = keysSlice

	// todo: checkTx in mock app
	mockApp.SetAnteHandler(nil)

	app := mockApp
	mockApp.MountStores(
		//app.keyOrder,
		app.keyToken,
		app.keyTokenPair,
		app.keyLock,
		app.keySupply,
		app.keyDex,
		app.keySwap,
		app.keyFarm,
	)

	require.NoError(t, mockApp.CompleteSetup(mockApp.keyOrder))
	mock.SetGenesis(mockApp.App, genAccs)
	for i := 0; i < numGenAccs; i++ {
		mock.CheckBalance(t, app.App, keysSlice[i].Address, coins)
		mockApp.TotalCoinsSupply = mockApp.TotalCoinsSupply.Add(coins...)
	}
	return
}

func getEndBlocker(orderKeeper keeper.Keeper, backendKeeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		order.EndBlocker(ctx, orderKeeper)
		EndBlocker(ctx, backendKeeper)
		return abci.ResponseEndBlock{}
	}
}

func getInitChainer(mapp *mock.App, bankKeeper bank.Keeper, supplyKeeper supply.Keeper,
	blacklistedAddrs []exported.ModuleAccountI) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}
		bankKeeper.SetSendEnabled(ctx, true)
		supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
		return abci.ResponseInitChain{}
	}
}

func buildTx(app *MockApp, ctx sdk.Context, addrKeys mock.AddrKeys, msg []sdk.Msg) auth.StdTx {
	accs := app.AccountKeeper.GetAccount(ctx, addrKeys.Address)
	accNum := accs.GetAccountNumber()
	seqNum := accs.GetSequence()

	tx := mock.GenTx(msg, []uint64{uint64(accNum)}, []uint64{uint64(seqNum)}, addrKeys.PrivKey)
	_, _, err := app.Check(tx)
	if err != nil {
		panic("something wrong in checking transaction")
	}
	return tx
}

func mockApplyBlock(app *MockApp, ctx sdk.Context, txs []auth.StdTx) {
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: ctx.BlockHeight()}})

	orderParam := ordertypes.DefaultParams()
	app.orderKeeper.SetParams(ctx, &orderParam)
	tokenParam := tokentypes.DefaultParams()
	app.tokenKeeper.SetParams(ctx, tokenParam)
	for i, tx := range txs {
		_, _, err := app.Deliver(tx)
		if err == nil {
			txBytes, _ := auth.DefaultTxEncoder(app.Cdc)(tx)
			txHash := fmt.Sprintf("%X", types3.Tx(txBytes).Hash())
			app.Logger().Info(fmt.Sprintf("[Sync Tx(%s) to backend module]", txHash))
			app.backendKeeper.SyncTx(ctx, &txs[i], txHash, ctx.BlockHeader().Time.Unix()) // do not use tx
		} else {
			app.Logger().Error(fmt.Sprintf("DeliverTx failed: %v", err))
		}
	}

	app.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
	app.Commit(abci.RequestCommit{})
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

func mockOrder(orderID, product, side, price, quantity string) *ordertypes.Order {
	return &ordertypes.Order{
		OrderID:           orderID,
		Product:           product,
		Side:              side,
		Price:             sdk.MustNewDecFromStr(price),
		FilledAvgPrice:    sdk.ZeroDec(),
		Quantity:          sdk.MustNewDecFromStr(quantity),
		RemainQuantity:    sdk.MustNewDecFromStr(quantity),
		Status:            ordertypes.OrderStatusOpen,
		OrderExpireBlocks: ordertypes.DefaultOrderExpireBlocks,
		FeePerBlock:       ordertypes.DefaultFeePerBlock,
	}
}

func FireEndBlockerPeriodicMatch(t *testing.T, enableBackend bool) (mockDexApp *MockApp, orders []*ordertypes.Order) {
	mapp, addrKeysSlice := getMockApp(t, 2, enableBackend, "")
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{Time: time.Now()}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	feeParams := ordertypes.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)
	tokenPair := dex.GetBuiltInTokenPair()

	mapp.dexKeeper.SetOperator(ctx, types2.DEXOperator{Address: tokenPair.Owner, HandlingFeeAddress: tokenPair.Owner})
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	// mock orders
	orders = []*ordertypes.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "1.5"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	orders[1].Sender = addrKeysSlice[1].Address
	for i := 0; i < 2; i++ {
		err := mapp.orderKeeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}

	// call EndBlocker to execute periodic match

	order.EndBlocker(ctx, mapp.orderKeeper)
	EndBlocker(ctx, mapp.backendKeeper)
	return mapp, orders
}

func TestAppModule(t *testing.T) {
	mapp, _ := getMockApp(t, 2, false, "")
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{Time: time.Now()}).WithBlockHeight(10)
	app := NewAppModule(mapp.backendKeeper)

	require.Equal(t, true, app.InitGenesis(ctx, nil) == nil)
	require.Equal(t, nil, app.ValidateGenesis(nil))
	require.Equal(t, true, app.DefaultGenesis() == nil)
	require.Equal(t, true, app.ExportGenesis(ctx) == nil)
	require.Equal(t, true, app.NewHandler() == nil)
	require.Equal(t, true, app.GetTxCmd(mapp.Cdc) == nil)
	require.EqualValues(t, cli.GetQueryCmd(QuerierRoute, mapp.Cdc).Name(), app.GetQueryCmd(mapp.Cdc).Name())
	require.Equal(t, ModuleName, app.Name())
	require.Equal(t, ModuleName, app.AppModuleBasic.Name())
	require.Equal(t, true, app.NewQuerierHandler() != nil)
	require.Equal(t, RouterKey, app.Route())
	require.Equal(t, QuerierRoute, app.QuerierRoute())
	require.Equal(t, true, app.EndBlock(ctx, abci.RequestEndBlock{}) == nil)

	rs := lcd.NewRestServer(mapp.Cdc, nil)
	app.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

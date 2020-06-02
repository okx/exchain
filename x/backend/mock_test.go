package backend

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okchain/x/backend/client/cli"
	"github.com/okex/okchain/x/backend/config"
	"github.com/okex/okchain/x/backend/orm"
	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common/monitor"
	"github.com/okex/okchain/x/common/version"
	"github.com/okex/okchain/x/order/keeper"
	ordertypes "github.com/okex/okchain/x/order/types"
	tokentypes "github.com/okex/okchain/x/token/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/okex/okchain/x/common"

	//"github.com/okex/okchain/x/gov"
	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order"

	//"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/token"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
		true,
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
		nil,
		mockApp.ParamsKeeper.Subspace(ordertypes.DefaultParamspace),
		auth.FeeCollectorName,
		mockApp.keyOrder,
		mockApp.Cdc,
		true,
		monitor.NopOrderMetrics())

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
	mockApp.SetInitChainer(getInitChainer(mockApp.App, mockApp.supplyKeeper,
		[]exported.ModuleAccountI{feeCollector}))

	intQuantity := 100
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
	)

	require.NoError(t, mockApp.CompleteSetup(mockApp.keyOrder))
	mock.SetGenesis(mockApp.App, genAccs)
	for i := 0; i < numGenAccs; i++ {
		mock.CheckBalance(t, app.App, keysSlice[i].Address, coins)
		mockApp.TotalCoinsSupply = mockApp.TotalCoinsSupply.Add(coins)
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

func getInitChainer(mapp *mock.App, supplyKeeper supply.Keeper,
	blacklistedAddrs []exported.ModuleAccountI) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}
		supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
		return abci.ResponseInitChain{}
	}
}

func buildTx(app *MockApp, ctx sdk.Context, addrKeys mock.AddrKeys, msg sdk.Msg) auth.StdTx {
	accs := app.AccountKeeper.GetAccount(ctx, addrKeys.Address)
	accNum := accs.GetAccountNumber()
	seqNum := accs.GetSequence()

	tx := mock.GenTx([]sdk.Msg{msg}, []uint64{uint64(accNum)}, []uint64{uint64(seqNum)}, addrKeys.PrivKey)
	res := app.Check(tx)
	if !res.IsOK() {
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
		response := app.Deliver(tx)
		if response.IsOK() {
			txBytes, _ := auth.DefaultTxEncoder(app.Cdc)(tx)
			txHash := fmt.Sprintf("%X", tmhash.Sum(txBytes))
			app.Logger().Info(fmt.Sprintf("[Sync Tx(%s) to backend module]", txHash))
			app.backendKeeper.SyncTx(ctx, &txs[i], txHash, ctx.BlockHeader().Time.Unix()) // do not use tx
		} else {
			app.Logger().Error(fmt.Sprintf("DeliverTx failed: %v", response))
		}
	}

	app.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
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

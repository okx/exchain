package dex

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okchain/x/common"
	ordertypes "github.com/okex/okchain/x/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type mockTokenKeeper struct {
	exist bool
}

// TokenExist return true if token exist
func (k *mockTokenKeeper) TokenExist(ctx sdk.Context, symbol string) bool {
	return k.exist
}

// nolint
func newMockTokenKeeper() *mockTokenKeeper {
	return &mockTokenKeeper{
		exist: true,
	}
}

type mockSupplyKeeper struct {
	behaveEvil    bool
	moduleAccount exported.ModuleAccountI
}

func (k *mockSupplyKeeper) behave() sdk.Error {
	if k.behaveEvil {
		return sdk.ErrInternal("raise an mock exception here")
	}
	return nil
}

// SendCoinsFromAccountToModule mocks SendCoinsFromAccountToModule of supply.Keeper
func (k *mockSupplyKeeper) SendCoinsFromAccountToModule(
	ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error {
	return k.behave()
}

// SendCoinsFromModuleToAccount mocks SendCoinsFromModuleToAccount of supply.Keeper
func (k *mockSupplyKeeper) SendCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.behave()
}

// GetModuleAccount returns the ModuleAccount
func (k *mockSupplyKeeper) GetModuleAccount(
	ctx sdk.Context, moduleName string) exported.ModuleAccountI {
	return k.moduleAccount
}

// GetModuleAddress returns address of the ModuleAccount
func (k *mockSupplyKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return k.moduleAccount.GetAddress()
}

// MintCoins mocks MintCoins of supply.Keeper
func (k *mockSupplyKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error {
	return k.behave()
}

// nolint
func newMockSupplyKeeper() *mockSupplyKeeper {
	return &mockSupplyKeeper{
		behaveEvil:    true,
		moduleAccount: supply.NewEmptyModuleAccount(ModuleName),
	}
}

type mockDexKeeper struct {
	*Keeper
	getFakeTokenPair      bool
	failToDeleteTokenPair bool
	failToWithdraw        bool
	failToDeposit         bool
	failToMarshal         bool
}

// LockTokenPair mocks LockTokenPair of dex.Keeper
func (k *mockDexKeeper) LockTokenPair(ctx sdk.Context, product string, lock *ordertypes.ProductLock) {
	k.Keeper.LockTokenPair(ctx, product, lock)
}

// GetLockedProductsCopy returns map with product locked
func (k *mockDexKeeper) GetLockedProductsCopy(ctx sdk.Context) *ordertypes.ProductLockMap {
	return ordertypes.NewProductLockMap()
}

// LockTokenPair mocks LockTokenPair of dex.Keeper
func (k *mockDexKeeper) GetTokenPair(ctx sdk.Context, product string) *TokenPair {
	if k.getFakeTokenPair {
		return GetBuiltInTokenPair()
	}
	return k.Keeper.GetTokenPair(ctx, product)
}

// Deposit mocks Deposit of dex.Keeper
func (k *mockDexKeeper) Deposit(ctx sdk.Context, product string, from sdk.AccAddress, amount sdk.DecCoin) sdk.Error {
	if k.failToDeposit {
		return sdk.ErrInternal("raise an mock exception here")
	}
	return nil
}

// GetCDC mocks GetCDC of dex.Keeper
func (k *mockDexKeeper) GetCDC() *codec.Codec {
	if k.failToMarshal {
		return nil
	}
	return k.Keeper.GetCDC()
}

// Withdraw mocks Withdraw of dex.Keeper
func (k *mockDexKeeper) Withdraw(ctx sdk.Context, product string, from sdk.AccAddress, amount sdk.DecCoin) sdk.Error {
	if k.failToWithdraw {
		return sdk.ErrInternal("raise an mock exception here")
	}
	return nil
}

// DeleteTokenPairByName mocks DeleteTokenPairByName of dex.Keeper
func (k *mockDexKeeper) DeleteTokenPairByName(ctx sdk.Context, owner sdk.AccAddress, tokenPairName string) {
}

func newMockDexKeeper(baseDexKeeper *Keeper) *mockDexKeeper {
	m := mockDexKeeper{
		baseDexKeeper,
		true,
		false,
		false,
		false,
		false,
	}
	return &m
}

type mockApp struct {
	*mock.App

	// expected keeper
	tokenKeeper   TokenKeeper
	suppleyKeeper SupplyKeeper
	dexKeeper     IKeeper

	bankKeeper    BankKeeper
	stakingKeeper StakingKeeper

	// expected KVStoreKey to mount
	storeKey     *sdk.KVStoreKey
	keyTokenPair *sdk.KVStoreKey
	keySupply    *sdk.KVStoreKey
}

// nolint
func newMockApp(tokenKeeper TokenKeeper, supplyKeeper SupplyKeeper, accountsInGenisis int) (
	app *mockApp, mockDexKeeper *mockDexKeeper, err error) {

	mApp := mock.NewApp()
	RegisterCodec(mApp.Cdc)

	storeKey := sdk.NewKVStoreKey(StoreKey)
	keyTokenPair := sdk.NewKVStoreKey(TokenPairStoreKey)
	supplyKvStoreKey := sdk.NewKVStoreKey(supply.StoreKey)

	paramsKeeper := mApp.ParamsKeeper
	paramsSubspace := paramsKeeper.Subspace(DefaultParamspace)

	dexKeeper := NewKeeper(AuthFeeCollector, supplyKeeper, paramsSubspace, tokenKeeper, nil, nil,
		storeKey, keyTokenPair, mApp.Cdc)

	dexKeeper.SetGovKeeper(mockGovKeeper{})

	fakeDexKeeper := newMockDexKeeper(&dexKeeper)

	app = &mockApp{
		App:           mApp,
		storeKey:      storeKey,
		bankKeeper:    nil,
		keyTokenPair:  keyTokenPair,
		stakingKeeper: nil,
		keySupply:     supplyKvStoreKey,
		suppleyKeeper: supplyKeeper,
		tokenKeeper:   tokenKeeper,
		dexKeeper:     fakeDexKeeper,
	}

	dexHandler := NewHandler(fakeDexKeeper)
	dexQuerier := NewQuerier(fakeDexKeeper)
	app.Router().AddRoute(RouterKey, dexHandler)
	app.QueryRouter().AddRoute(QuerierRoute, dexQuerier)

	app.SetEndBlocker(getEndBlocker())
	app.SetInitChainer(getInitChainer(mApp, dexKeeper))

	initQuantity := 10000000
	var decCoins sdk.DecCoins
	decCoins, err = sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		initQuantity, common.NativeToken, initQuantity, common.TestToken))
	if err != nil {
		return nil, nil, err
	}
	genAccs, _, _, _ := mock.CreateGenAccounts(accountsInGenisis, decCoins)
	app.SetAnteHandler(nil)

	app.MountStores(app.storeKey, app.keyTokenPair, app.keySupply)

	err = app.CompleteSetup()
	if err == nil {
		mock.SetGenesis(app.App, genAccs)
	}

	return app, fakeDexKeeper, err
}

func getEndBlocker() sdk.EndBlocker {
	return func(_ sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{}
	}
}

func getInitChainer(mapp *mock.App, keeper IKeeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		InitGenesis(ctx, keeper, DefaultGenesisState())
		return abci.ResponseInitChain{}
	}
}

type mockGovKeeper struct{}

// RemoveFromActiveProposalQueue mocks RemoveFromActiveProposalQueue of gov.Keeper
func (k mockGovKeeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
}

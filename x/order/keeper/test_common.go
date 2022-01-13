package keeper

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/okex/exchain/x/common"

	"github.com/okex/exchain/x/common/monitor"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/params"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/order/types"
	"github.com/okex/exchain/x/token"
)

var mockOrder = types.MockOrder

// TestInput stores some variables for testing
type TestInput struct {
	Ctx       sdk.Context
	Cdc       *codec.Codec
	TestAddrs []sdk.AccAddress

	OrderKeeper   Keeper
	TokenKeeper   token.Keeper
	AccountKeeper auth.AccountKeeper
	SupplyKeeper  supply.Keeper
	DexKeeper     dex.Keeper
}

// MakeTestCodec creates a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	dex.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // order
	token.RegisterCodec(cdc) // token
	return cdc
}

// CreateTestInputWithBalance creates TestInput with the number of account and the quantity
func CreateTestInputWithBalance(t *testing.T, numAddrs, initQuantity int64) TestInput {

	db := dbm.NewMemDB()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// order module
	keyOrder := sdk.NewKVStoreKey(types.OrderStoreKey)

	// token module
	keyToken := sdk.NewKVStoreKey(token.StoreKey)
	keyLock := sdk.NewKVStoreKey(token.KeyLock)
	//keyTokenPair := sdk.NewKVStoreKey(token.KeyTokenPair)

	// dex module
	storeKey := sdk.NewKVStoreKey(dex.StoreKey)
	keyTokenPair := sdk.NewKVStoreKey(dex.TokenPairStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	ms.MountStoreWithDB(keyOrder, sdk.StoreTypeIAVL, db)

	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyLock, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTokenPair, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))
	cdc := MakeTestCodec()

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true

	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc,
		paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace),
		blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		token.ModuleName:      {supply.Minter, supply.Burner},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

	// set module accounts
	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)

	// token keeper
	tokenKeepr := token.NewKeeper(bankKeeper, paramsKeeper.Subspace(token.DefaultParamspace),
		auth.FeeCollectorName, supplyKeeper, keyToken, keyLock, cdc, accountKeeper)

	// dex keeper
	paramsSubspace := paramsKeeper.Subspace(dex.DefaultParamspace)
	dexKeeper := dex.NewKeeper(auth.FeeCollectorName, supplyKeeper, paramsSubspace, tokenKeepr, nil, bankKeeper, storeKey, keyTokenPair, cdc)

	// order keeper
	orderKeeper := NewKeeper(tokenKeepr, supplyKeeper, dexKeeper,
		paramsKeeper.Subspace(types.DefaultParamspace), auth.FeeCollectorName, keyOrder,
		cdc, monitor.NopOrderMetrics())

	defaultParams := types.DefaultTestParams()
	orderKeeper.SetParams(ctx, &defaultParams)

	// init account tokens
	decCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		initQuantity, common.NativeToken, initQuantity, common.TestToken))
	require.Nil(t, err)

	initCoins := decCoins

	var testAddrs []sdk.AccAddress
	for i := int64(0); i < numAddrs; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		addr := sdk.AccAddress(pk.Address())
		testAddrs = append(testAddrs, addr)
		//_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		err := orderKeeper.supplyKeeper.MintCoins(ctx, token.ModuleName, initCoins)
		require.Nil(t, err)
		err = orderKeeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, token.ModuleName, addr, initCoins)
		require.Nil(t, err)
	}

	return TestInput{ctx, cdc, testAddrs, orderKeeper, tokenKeepr, accountKeeper, supplyKeeper, dexKeeper}
}

// CreateTestInput creates TestInput with default params
func CreateTestInput(t *testing.T) TestInput {
	return CreateTestInputWithBalance(t, 2, 100)
}

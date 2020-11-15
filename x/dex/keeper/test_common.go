package keeper

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/okex/okexchain/x/staking"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/dex/types"
	"github.com/okex/okexchain/x/gov"
	"github.com/okex/okexchain/x/token"
)

type testInput struct {
	Ctx       sdk.Context
	Cdc       *codec.Codec
	TestAddrs []sdk.AccAddress

	DexKeeper Keeper
}

// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // dex
	return cdc
}

func createTestInputWithBalance(t *testing.T, numAddrs, initQuantity int64) testInput {
	db := dbm.NewMemDB()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// token module
	keyToken := sdk.NewKVStoreKey(token.StoreKey)
	keyLock := sdk.NewKVStoreKey(token.KeyLock)
	//keyTokenPair := sdk.NewKVStoreKey(token.KeyTokenPair)

	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)

	// dex module
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyTokenPair := sdk.NewKVStoreKey(types.TokenPairStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)

	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyLock, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTokenPair, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))
	cdc := makeTestCodec()

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true

	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc,
		paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace, blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		token.ModuleName:      {supply.Minter, supply.Burner},
		types.ModuleName:      nil,
		gov.ModuleName:        nil,
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.SysCoins{}))

	// set module accounts
	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)

	// token keeper
	tokenKeepr := token.NewKeeper(bankKeeper,
		paramsKeeper.Subspace(token.DefaultParamspace), auth.FeeCollectorName, supplyKeeper,
		keyToken, keyLock, cdc, true)

	paramsSubspace := paramsKeeper.Subspace(types.DefaultParamspace)

	mockStakingKeeper := &mockStakingKeeper{true}
	mockBankKeeper := mockBankKeeper{}

	// dex keeper
	dexKeeper := NewKeeper(auth.FeeCollectorName, supplyKeeper, paramsSubspace, tokenKeepr, mockStakingKeeper, mockBankKeeper, storeKey, keyTokenPair, cdc)

	// init account tokens
	decCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		initQuantity, common.NativeToken, initQuantity, common.TestToken))
	if err != nil {
		panic(err)
	}

	var testAddrs []sdk.AccAddress
	for i := int64(0); i < numAddrs; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		addr := sdk.AccAddress(pk.Address())
		testAddrs = append(testAddrs, addr)
		err := dexKeeper.supplyKeeper.MintCoins(ctx, token.ModuleName, decCoins)
		require.Nil(t, err)
		err = dexKeeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, token.ModuleName, addr, decCoins)
		require.Nil(t, err)
	}

	return testInput{ctx, cdc, testAddrs, dexKeeper}
}

// nolint
func createTestInput(t *testing.T) testInput {
	return createTestInputWithBalance(t, 2, 100)
}

type mockStakingKeeper struct {
	getFakeValidator bool
}

func (m *mockStakingKeeper) IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool {
	return m.getFakeValidator
}

func (m *mockStakingKeeper) SetFakeValidator(fakeValidator bool) {
	m.getFakeValidator = fakeValidator
}

type mockBankKeeper struct{}

// GetCoins returns coins for test
func (keeper mockBankKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.SysCoins {
	return sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDec(2500))
}

// GetBuiltInTokenPair returns built in token pair for test
func GetBuiltInTokenPair() *types.TokenPair {
	addr, err := sdk.AccAddressFromBech32(types.TestTokenPairOwner)
	if err != nil {
		panic(err)
	}
	return &types.TokenPair{
		ID:               1,
		BaseAssetSymbol:  common.TestToken,
		QuoteAssetSymbol: sdk.DefaultBondDenom,
		InitPrice:        sdk.MustNewDecFromStr("10.0"),
		MaxPriceDigit:    8,
		MaxQuantityDigit: 8,
		MinQuantity:      sdk.MustNewDecFromStr("0"),
		Owner:            addr,
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
	}
}

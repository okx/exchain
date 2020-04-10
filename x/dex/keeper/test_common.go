package keeper

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/okex/okchain/x/staking"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okchain/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/dex/types"
	"github.com/okex/okchain/x/gov"
	"github.com/okex/okchain/x/token"
)

type TestInput struct {
	Ctx       sdk.Context
	Cdc       *codec.Codec
	TestAddrs []sdk.AccAddress

	DexKeeper Keeper
}

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // dex
	return cdc
}

func CreateTestInputWithBalance(t *testing.T, numAddrs, initQuantity int64) TestInput {
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
	cdc := MakeTestCodec()

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
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

	// set module accounts
	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)

	// token keeper
	tokenKeepr := token.NewKeeper(bankKeeper, paramsKeeper,
		paramsKeeper.Subspace(token.DefaultParamspace), auth.FeeCollectorName, supplyKeeper,
		keyToken, keyLock, cdc, true)

	paramsSubspace := paramsKeeper.Subspace(types.DefaultParamspace)

	mockStakingKeeper := &MockStakingKeeper{true}
	mockBankKeeper := MockBankKeeper{}

	// dex keeper
	dexKeeper := NewKeeper(auth.FeeCollectorName, supplyKeeper, paramsSubspace, tokenKeepr, mockStakingKeeper, mockBankKeeper, storeKey, keyTokenPair, cdc)

	// init account tokens
	decCoins, _ := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		initQuantity, common.NativeToken, initQuantity, common.TestToken))
	initCoins := decCoins

	var testAddrs []sdk.AccAddress
	for i := int64(0); i < numAddrs; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		addr := sdk.AccAddress(pk.Address())
		testAddrs = append(testAddrs, addr)
		err := dexKeeper.supplyKeeper.MintCoins(ctx, token.ModuleName, initCoins)
		require.Nil(t, err)
		err = dexKeeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, token.ModuleName, addr, initCoins)
		require.Nil(t, err)
	}

	return TestInput{ctx, cdc, testAddrs, dexKeeper}
}

func CreateTestInput(t *testing.T) TestInput {
	return CreateTestInputWithBalance(t, 2, 100)
}

type MockStakingKeeper struct {
	getFakeValidator bool
}

func (m *MockStakingKeeper) IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool {
	return m.getFakeValidator
}

func (m *MockStakingKeeper) SetFakeValidator(fakeValidator bool) {
	m.getFakeValidator = fakeValidator
}

type MockBankKeeper struct{}

func (keeper MockBankKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDec(2500))
}

func GetBuiltInTokenPair() *types.TokenPair {
	addr, _ := sdk.AccAddressFromBech32(types.TestTokenPairOwner)
	return &types.TokenPair{
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

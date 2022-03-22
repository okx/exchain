package keeper

import (
	types2 "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/params"
	"github.com/okex/exchain/x/staking"
	"github.com/stretchr/testify/require"
)

//nolint: deadcode unused
var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delPk2   = ed25519.GenPrivKey().PubKey()
	delPk3   = ed25519.GenPrivKey().PubKey()
	delPk4   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
	delAddr2 = sdk.AccAddress(delPk2.Address())
	delAddr3 = sdk.AccAddress(delPk3.Address())
	delAddr4 = sdk.AccAddress(delPk4.Address())

	valOpPk1    = ed25519.GenPrivKey().PubKey()
	valOpPk2    = ed25519.GenPrivKey().PubKey()
	valOpPk3    = ed25519.GenPrivKey().PubKey()
	valOpPk4    = ed25519.GenPrivKey().PubKey()
	valOpAddr1  = sdk.ValAddress(valOpPk1.Address())
	valOpAddr2  = sdk.ValAddress(valOpPk2.Address())
	valOpAddr3  = sdk.ValAddress(valOpPk3.Address())
	valOpAddr4  = sdk.ValAddress(valOpPk4.Address())
	valAccAddr1 = sdk.AccAddress(valOpPk1.Address()) // generate acc addresses for these validator keys too
	valAccAddr2 = sdk.AccAddress(valOpPk2.Address())
	valAccAddr3 = sdk.AccAddress(valOpPk3.Address())
	valAccAddr4 = sdk.AccAddress(valOpPk4.Address())

	valConsPk1   = ed25519.GenPrivKey().PubKey()
	valConsPk2   = ed25519.GenPrivKey().PubKey()
	valConsPk3   = ed25519.GenPrivKey().PubKey()
	valConsPk4   = ed25519.GenPrivKey().PubKey()
	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())
	valConsAddr3 = sdk.ConsAddress(valConsPk3.Address())
	valConsAddr4 = sdk.ConsAddress(valConsPk4.Address())

	// TODO move to common testing package for all modules
	// test addresses
	TestAddrs = []sdk.AccAddress{
		delAddr1, delAddr2, delAddr3, delAddr4,
		valAccAddr1, valAccAddr2, valAccAddr3, valAccAddr4,
	}

	distrAcc = supply.NewEmptyModuleAccount(types.ModuleName)
)

func ReInit() {
	delPk1 = ed25519.GenPrivKey().PubKey()
	delPk2 = ed25519.GenPrivKey().PubKey()
	delPk3 = ed25519.GenPrivKey().PubKey()
	delPk4 = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
	delAddr2 = sdk.AccAddress(delPk2.Address())
	delAddr3 = sdk.AccAddress(delPk3.Address())
	delAddr4 = sdk.AccAddress(delPk4.Address())

	valOpPk1 = ed25519.GenPrivKey().PubKey()
	valOpPk2 = ed25519.GenPrivKey().PubKey()
	valOpPk3 = ed25519.GenPrivKey().PubKey()
	valOpPk4 = ed25519.GenPrivKey().PubKey()
	valOpAddr1 = sdk.ValAddress(valOpPk1.Address())
	valOpAddr2 = sdk.ValAddress(valOpPk2.Address())
	valOpAddr3 = sdk.ValAddress(valOpPk3.Address())
	valOpAddr4 = sdk.ValAddress(valOpPk4.Address())
	valAccAddr1 = sdk.AccAddress(valOpPk1.Address()) // generate acc addresses for these validator keys too
	valAccAddr2 = sdk.AccAddress(valOpPk2.Address())
	valAccAddr3 = sdk.AccAddress(valOpPk3.Address())
	valAccAddr4 = sdk.AccAddress(valOpPk4.Address())

	valConsPk1 = ed25519.GenPrivKey().PubKey()
	valConsPk2 = ed25519.GenPrivKey().PubKey()
	valConsPk3 = ed25519.GenPrivKey().PubKey()
	valConsPk4 = ed25519.GenPrivKey().PubKey()
	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())
	valConsAddr3 = sdk.ConsAddress(valConsPk3.Address())
	valConsAddr4 = sdk.ConsAddress(valConsPk4.Address())

	// TODO move to common testing package for all modules
	// test addresses
	TestAddrs = []sdk.AccAddress{
		delAddr1, delAddr2, delAddr3, delAddr4,
		valAccAddr1, valAccAddr2, valAccAddr3, valAccAddr4,
	}

	distrAcc = supply.NewEmptyModuleAccount(types.ModuleName)
}

// GetTestAddrs returns valOpAddrs, valConsPks, valConsAddrs for test
func GetTestAddrs() ([]sdk.ValAddress, []crypto.PubKey, []sdk.ConsAddress) {
	valOpAddrs := []sdk.ValAddress{valOpAddr1, valOpAddr2, valOpAddr3, valOpAddr4}
	valConsPks := []crypto.PubKey{valConsPk1, valConsPk2, valConsPk3, valConsPk4}
	valConsAddrs := []sdk.ConsAddress{valConsAddr1, valConsAddr2, valConsAddr3, valConsAddr4}
	return valOpAddrs, valConsPks, valConsAddrs
}

// NewTestSysCoins returns dec coins
func NewTestSysCoins(i int64, precison int64) sdk.SysCoins {
	return sdk.SysCoins{NewTestSysCoin(i, precison)}
}

// NewTestSysCoin returns one dec coin
func NewTestSysCoin(i int64, precison int64) sdk.SysCoin {
	return sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(i, precison))
}

// MakeTestCodec creates a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // distr
	return cdc
}

// CreateTestInputDefault test input with default values
func CreateTestInputDefault(t *testing.T, isCheckTx bool, initPower int64) (
	sdk.Context, auth.AccountKeeper, Keeper, staking.Keeper, types.SupplyKeeper) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, ak, _, dk, sk, _, supplyKeeper := CreateTestInputAdvanced(t, isCheckTx, initPower, communityTax)
	h := staking.NewHandler(sk)
	valOpAddrs, valConsPks, _ := GetTestAddrs()
	// create four validators
	for i := int64(0); i < 4; i++ {
		msg := staking.NewMsgCreateValidator(valOpAddrs[i], valConsPks[i],
			staking.Description{}, NewTestSysCoin(i+1, 0))
		// assert initial state: zero current rewards
		_, e := h(ctx, msg)
		require.Nil(t, e)
		require.True(t, dk.GetValidatorAccumulatedCommission(ctx, valOpAddrs[i]).IsZero())
	}
	return ctx, ak, dk, sk, supplyKeeper
}

// CreateTestInputAdvanced hogpodge of all sorts of input required for testing
func CreateTestInputAdvanced(t *testing.T, isCheckTx bool, initPower int64, communityTax sdk.Dec) (
	sdk.Context, auth.AccountKeeper, bank.Keeper, Keeper, staking.Keeper, params.Keeper, types.SupplyKeeper) {

	initTokens := sdk.TokensFromConsensusPower(initPower)

	keyDistr := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true
	blacklistedAddrs[distrAcc.GetAddress().String()] = true

	cdc := MakeTestCodec()
	reg := types2.NewInterfaceRegistry()
	cc := codec.NewProtoCodec(reg)
	pro := codec.NewCodecProxy(cc, cdc)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace),
		blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)

	sk := staking.NewKeeper(cdc, pro, keyStaking, supplyKeeper,
		pk.Subspace(staking.DefaultParamspace))
	sk.SetParams(ctx, staking.DefaultParams())

	keeper := NewKeeper(cdc, keyDistr, pk.Subspace(types.DefaultParamspace), sk, supplyKeeper,
		auth.FeeCollectorName, blacklistedAddrs)

	keeper.SetWithdrawAddrEnabled(ctx, true)
	initCoins := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range TestAddrs {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(t, err)
	}

	// set module accounts
	keeper.supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, distrAcc)

	// set the distribution hooks on staking
	sk.SetHooks(keeper.Hooks())

	// set genesis items required for distribution
	keeper.SetFeePool(ctx, types.InitialFeePool())
	keeper.SetCommunityTax(ctx, communityTax)

	return ctx, accountKeeper, bankKeeper, keeper, sk, pk, supplyKeeper
}

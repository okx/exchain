package keeper

import (
	"testing"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	types2 "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/supply"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	dbm "github.com/okx/okbchain/libs/tm-db"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/params"
	"github.com/okx/okbchain/x/staking"
	"github.com/stretchr/testify/require"
)

// CreateTestInputDefaultForBenchmark test input with default values
func CreateTestInputDefaultForBenchmark(b *testing.B, isCheckTx bool, initPower int64, newVersion bool) (
	sdk.Context, auth.AccountKeeper, Keeper, staking.Keeper, types.SupplyKeeper) {
	ctx, ak, _, dk, sk, _, supplyKeeper := CreateTestInputAdvancedForBenchmark(b, isCheckTx, initPower)
	h := staking.NewHandler(sk)
	valOpAddrs, valConsPks, _ := GetTestAddrs()

	if newVersion {
		dk.SetInitExistedValidatorFlag(ctx, true)
		dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	}

	// create four validators
	for i := int64(0); i < 4; i++ {
		msg := staking.NewMsgCreateValidator(valOpAddrs[i], valConsPks[i],
			staking.Description{}, NewTestSysCoin(i+1, 0))
		// assert initial state: zero current rewards
		_, e := h(ctx, msg)
		require.Nil(b, e)
		require.True(b, dk.GetValidatorAccumulatedCommission(ctx, valOpAddrs[i]).IsZero())
	}
	return ctx, ak, dk, sk, supplyKeeper
}

// CreateTestInputAdvancedForBenchmark hogpodge of all sorts of input required for testing
func CreateTestInputAdvancedForBenchmark(b *testing.B, isCheckTx bool, initPower int64) (
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
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeMPT, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(b, err)

	cdc := MakeTestCodec()
	reg := types2.NewInterfaceRegistry()
	cc := codec.NewProtoCodec(reg)
	pro := codec.NewCodecProxy(cc, cdc)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, log.NewNopLogger())
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), nil)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bank.NewBankKeeperAdapter(bankKeeper), maccPerms)
	sk := staking.NewKeeper(pro, keyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace))
	sk.SetParams(ctx, staking.DefaultDposParams())
	keeper := NewKeeper(cdc, keyDistr, pk.Subspace(types.DefaultParamspace), sk, supplyKeeper, auth.FeeCollectorName, nil)

	initCoins := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range TestAddrs {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(b, err)
	}

	// set the distribution hooks on staking
	sk.SetHooks(keeper.Hooks())
	return ctx, accountKeeper, bankKeeper, keeper, sk, pk, supplyKeeper
}

func BenchmarkAllocateTokensBefore(b *testing.B) {
	//start test
	ctx, _, k, sk, _ := CreateTestInputDefaultForBenchmark(b, false, 1000, false)
	val := sk.Validator(ctx, valOpAddr1)

	// allocate tokens
	tokens := NewTestSysCoins(123, 2)

	//reset benchmark timer
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		k.AllocateTokensToValidator(ctx, val, tokens)
	}

	require.Equal(b, tokens.MulDec(sdk.NewDec(int64(b.N))), k.GetValidatorAccumulatedCommission(ctx, val.GetOperator()))
}

func BenchmarkAllocateTokensAfter(b *testing.B) {
	//start test
	ctx, _, k, sk, _ := CreateTestInputDefaultForBenchmark(b, false, 1000, true)

	validator, found := sk.GetValidator(ctx, valOpAddr1)
	require.True(b, found)
	newRate, _ := sdk.NewDecFromStr("0.5")
	validator.Commission.Rate = newRate
	sk.SetValidator(ctx, validator)

	val := sk.Validator(ctx, valOpAddr1)
	// allocate tokens
	tokens := NewTestSysCoins(123, 2)

	//reset benchmark timer
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		k.AllocateTokensToValidator(ctx, val, tokens)
	}
	require.Equal(b, tokens.MulDec(sdk.NewDec(int64(b.N))).QuoDec(sdk.NewDec(int64(2))), k.GetValidatorAccumulatedCommission(ctx, val.GetOperator()))
}

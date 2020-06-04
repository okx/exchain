package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	wasmTypes "github.com/okex/okchain/x/wasm/internal/types"
)

const flagLRUCacheSize = "lru_size"
const flagQueryGasLimit = "query_gas_limit"

func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register AppAccount
	// cdc.RegisterInterface((*authexported.Account)(nil), nil)
	// cdc.RegisterConcrete(&auth.BaseAccount{}, "test/wasm/BaseAccount", nil)
	auth.AppModuleBasic{}.RegisterCodec(cdc)
	bank.AppModuleBasic{}.RegisterCodec(cdc)
	supply.AppModuleBasic{}.RegisterCodec(cdc)
	staking.AppModuleBasic{}.RegisterCodec(cdc)
	distribution.AppModuleBasic{}.RegisterCodec(cdc)
	wasmTypes.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

var TestingStakeParams = staking.Params{
	UnbondingTime:     100,
	MaxValidators:     10,
	MaxEntries:        10,
	//HistoricalEntries: 10,
	BondDenom:         "stake",
}

type TestKeepers struct {
	AccountKeeper auth.AccountKeeper
	StakingKeeper staking.Keeper
	WasmKeeper    Keeper
	DistKeeper    distribution.Keeper
	SupplyKeeper  supply.Keeper
}

// encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func CreateTestInput(t *testing.T, isCheckTx bool, tempDir string, supportedFeatures string, encoders *MessageEncoders, queriers *QueryPlugins) (sdk.Context, TestKeepers) {
	keyContract := sdk.NewKVStoreKey(wasmTypes.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyDistro := sdk.NewKVStoreKey(distribution.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyContract, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistro, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, isCheckTx, log.NewNopLogger())
	cdc := MakeTestCodec()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace, nil,
	)
	bankKeeper.SetSendEnabled(ctx, true)

	// this is also used to initialize module accounts (so nil is meaningful here)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:   nil,
		distribution.ModuleName: nil,
		//mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		//gov.ModuleName:            {supply.Burner},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(cdc, keyStaking, nil, supplyKeeper, pk.Subspace(staking.DefaultParamspace),staking.DefaultCodespace)
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	distKeeper := distribution.NewKeeper(cdc, keyDistro, pk.Subspace(distribution.DefaultParamspace), stakingKeeper, supplyKeeper, distribution.DefaultCodespace, auth.FeeCollectorName, nil)
	//distKeeper.SetParams(ctx, distribution.DefaultParams())
	stakingKeeper.SetHooks(distKeeper.Hooks())

	// set genesis items required for distribution
	distKeeper.SetFeePool(ctx, distribution.InitialFeePool())

	// total supply to track this
	totalSupply := sdk.NewCoins(sdk.NewInt64Coin("stake", 100000000))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	// set up initial accounts
	for name, perms := range maccPerms {
		mod := supply.NewEmptyModuleAccount(name, perms...)
		if name == staking.NotBondedPoolName {
			err = mod.SetCoins(totalSupply)
			require.NoError(t, err)
		} else if name == distribution.ModuleName {
			// some big pot to pay out
			err = mod.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("stake", 500000)))
			require.NoError(t, err)
		}
		supplyKeeper.SetModuleAccount(ctx, mod)
	}

	stakeAddr := supply.NewModuleAddress(staking.BondedPoolName)
	moduleAcct := accountKeeper.GetAccount(ctx, stakeAddr)
	require.NotNil(t, moduleAcct)

	router := baseapp.NewRouter()
	bh := bank.NewHandler(bankKeeper)
	router.AddRoute(bank.RouterKey, bh)
	sh := staking.NewHandler(stakingKeeper)
	router.AddRoute(staking.RouterKey, sh)
	dh := distribution.NewHandler(distKeeper)
	router.AddRoute(distribution.RouterKey, dh)

	// Load default wasm config
	wasmConfig := wasmTypes.DefaultWasmConfig()

	keeper := NewKeeper(cdc, keyContract, accountKeeper, bankKeeper, router, tempDir, wasmConfig, supportedFeatures, encoders, queriers)
	keeper.queryPlugins = keeper.queryPlugins.Merge(&QueryPlugins{
		Bank:    nil,
		Custom:  nil,
		Staking: StakingQuerier(stakingKeeper),
		Wasm:    nil,
	})
	// add wasm handler so we can loop-back (contracts calling contracts)
	router.AddRoute(wasmTypes.RouterKey, TestHandler(keeper))

	keepers := TestKeepers{
		AccountKeeper: accountKeeper,
		SupplyKeeper:  supplyKeeper,
		StakingKeeper: stakingKeeper,
		DistKeeper:    distKeeper,
		WasmKeeper:    keeper,
	}
	return ctx, keepers
}

// TestHandler returns a wasm handler for tests (to avoid circular imports)
func TestHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result{
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case wasmTypes.MsgInstantiateContract:
			return handleInstantiate(ctx, k, &msg)
		case *wasmTypes.MsgInstantiateContract:
			return handleInstantiate(ctx, k, msg)

		case wasmTypes.MsgExecuteContract:
			return handleExecute(ctx, k, &msg)
		case *wasmTypes.MsgExecuteContract:
			return handleExecute(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleInstantiate(ctx sdk.Context, k Keeper, msg *wasmTypes.MsgInstantiateContract) sdk.Result{
	contractAddr, err := k.Instantiate(ctx, msg.Code, msg.Sender, msg.InitMsg, msg.Label, msg.InitFunds)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	return sdk.Result{
		Data:   contractAddr,
		Events: ctx.EventManager().Events(),
	}
}

func handleExecute(ctx sdk.Context, k Keeper, msg *wasmTypes.MsgExecuteContract) sdk.Result{
	res, err := k.Execute(ctx, msg.Contract, msg.Sender, msg.Msg, msg.SentFunds)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	res.Events = ctx.EventManager().Events()
	return res
}

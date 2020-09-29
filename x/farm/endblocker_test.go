package farm

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	swap "github.com/okex/okexchain/x/ammswap/keeper"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/okex/okexchain/x/params"
	"github.com/okex/okexchain/x/token"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	TestChainID = "okexchain"
)

type MockFarmKeeper struct {
	Keeper
	StoreKey     sdk.StoreKey
	TkeyStoreKey sdk.StoreKey
	SupplyKeeper supply.Keeper
	MountedStore store.MultiStore
	AccKeeper    auth.AccountKeeper
}

type MockApp struct {
	*mock.App

	//keyOrder     *sdk.KVStoreKey
	//keyToken     *sdk.KVStoreKey
	//keyLock      *sdk.KVStoreKey
	//keyDex       *sdk.KVStoreKey
	//keyTokenPair *sdk.KVStoreKey
	//keySupply *sdk.KVStoreKey

	bankKeeper   bank.Keeper
	tokenKeeper  token.Keeper
	supplyKeeper supply.Keeper
	swapKeeper   swap.Keeper
}

func TestEndBlocker(t *testing.T) {

}

func getMockApp(t *testing.T) {
	mApp := mock.NewApp()

	// 0.1 init store key
	keyFarm := sdk.NewKVStoreKey(types.StoreKey)
	tkeyFarm := sdk.NewTransientStoreKey(types.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyToken := sdk.NewKVStoreKey(token.StoreKey)

	// 0.2 init db
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyFarm, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyFarm, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	// 0.3 init context
	ctx := sdk.NewContext(ms, abci.Header{ChainID: TestChainID}, false, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)

	// 0.4 init codec
	cdc := codec.New()
	types.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	token.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	swaptypes.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	// 1.1 init param keeper
	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	// 1.2 init account keeper
	ak := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	// 1.3 init bank keeper
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	farmAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Burner, supply.Minter)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[farmAcc.String()] = true

	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)

	// 1.4 init supply keeper
	maccPerms := map[string][]string{
		auth.FeeCollectorName:   nil,
		types.ModuleName: {supply.Burner, supply.Minter},
	}
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, sdk.NewDec(1000000000))))

	// 1.5 init token keeper
	tk := token.NewKeeper(bk, pk.Subspace(token.DefaultParamspace), auth.FeeCollectorName, sk, keyToken, keyToken, cdc, false)

	// 1.6 init farm keeper
	fk := NewKeeper(auth.FeeCollectorName, sk, tk, )
}
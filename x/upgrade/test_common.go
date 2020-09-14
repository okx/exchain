package upgrade

import (
	"encoding/hex"
	"os"
	"testing"

	//"github.com/okex/okexchain/x/staking/util"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/okex/okexchain/x/common/proto"

	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okexchain/x/params"

	//"github.com/okex/okexchain/x/staking"
	"github.com/okex/okexchain/x/staking"

	//"github.com/okex/okexchain/x/staking/types"
	"github.com/okex/okexchain/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	pubKeys = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB53"),
	}

	accAddrs = []sdk.AccAddress{
		sdk.AccAddress(pubKeys[0].Address()),
		sdk.AccAddress(pubKeys[1].Address()),
		sdk.AccAddress(pubKeys[2].Address()),
		sdk.AccAddress(pubKeys[3].Address()),
	}

	maccPerms = map[string][]string{
		staking.BondedPoolName:    {supply.Staking},
		staking.NotBondedPoolName: {supply.Staking},
	}
)

func testPrepare(t *testing.T) (ctx sdk.Context, keeper Keeper, stakingKeeper staking.Keeper, paramsKeeper params.Keeper) {
	skMap := sdk.NewKVStoreKeys(
		"main",
		auth.StoreKey,
		supply.StoreKey,

		// for staking/distr rollback to cosmos-sdk
		//staking.StoreKey, staking.DelegatorPoolKey, staking.RedelegationKeyM, staking.RedelegationActonKey, staking.UnbondingKey,
		staking.StoreKey,
		params.StoreKey,
		StoreKey,
	)
	tskMap := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	for _, v := range skMap {
		ms.MountStoreWithDB(v, sdk.StoreTypeIAVL, db)
	}

	for _, v := range tskMap {
		ms.MountStoreWithDB(v, sdk.StoreTypeTransient, db)
	}

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx = sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	cdc := getTestCodec()
	paramsKeeper = params.NewKeeper(cdc, skMap[params.StoreKey], tskMap[params.TStoreKey], params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, skMap[auth.StoreKey], paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)
	supplyKeeper := supply.NewKeeper(cdc, skMap[supply.StoreKey], accountKeeper, bankKeeper, maccPerms)

	// for staking/distr rollback to cosmos-sdk
	//stakingKeeper = staking.NewKeeper(
	//	cdc, skMap[staking.StoreKey], skMap[staking.DelegatorPoolKey], skMap[staking.RedelegationKeyM], skMap[staking.RedelegationActonKey], skMap[staking.UnbondingKey], tskMap[staking.TStoreKey],
	//	supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	stakingKeeper = staking.NewKeeper(
		cdc, skMap[staking.StoreKey], tskMap[staking.TStoreKey],
		supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)

	stakingKeeper.SetParams(ctx, types.DefaultParams())
	protocolKeeper := proto.NewProtocolKeeper(skMap["main"])
	keeper = NewKeeper(cdc, skMap[StoreKey], protocolKeeper, stakingKeeper, bankKeeper, paramsKeeper.Subspace(DefaultParamspace))
	return
}

func getTestCodec() *codec.Codec {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

func newPubKey(pubKey string) (res crypto.PubKey) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}
	var pubKeyEd25519 ed25519.PubKeyEd25519
	copy(pubKeyEd25519[:], pubKeyBytes[:])
	return pubKeyEd25519
}

package keeper

import (
	"encoding/hex"
	"os"
	"testing"

	//"github.com/okex/okchain/x/staking/util"
	"github.com/tendermint/tendermint/crypto/ed25519"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okchain/x/common/proto"
	"github.com/okex/okchain/x/params"

	//"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/upgrade/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	pubKeys = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
	}

	accAddrs = []sdk.AccAddress{
		sdk.AccAddress(pubKeys[0].Address()),
		sdk.AccAddress(pubKeys[1].Address()),
	}

	maccPerms = map[string][]string{
		staking.BondedPoolName:    {supply.Staking},
		staking.NotBondedPoolName: {supply.Staking},
	}
)

func testPrepare(t *testing.T) (ctx sdk.Context, keeper Keeper) {
	skMap := sdk.NewKVStoreKeys(
		"main",
		auth.StoreKey,
		supply.StoreKey,
		staking.StoreKey,
		params.StoreKey,
		types.StoreKey,
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

	require.NoError(t, ms.LoadLatestVersion())

	ctx = sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	cdc := getTestCodec()
	paramsKeeper := params.NewKeeper(cdc, skMap[params.StoreKey], tskMap[params.TStoreKey], params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, skMap[auth.StoreKey], paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)
	supplyKeeper := supply.NewKeeper(cdc, skMap[supply.StoreKey], accountKeeper, bankKeeper, maccPerms)

	stakingKeeper := staking.NewKeeper(
		cdc, skMap[staking.StoreKey], tskMap[staking.TStoreKey],
		supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	stakingKeeper.SetParams(ctx, staking.DefaultParams())
	protocolKeeper := proto.NewProtocolKeeper(skMap["main"])
	keeper = NewKeeper(cdc, skMap[types.StoreKey], protocolKeeper, stakingKeeper, bankKeeper, paramsKeeper.Subspace(types.DefaultParamspace))
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

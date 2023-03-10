package token

import (
	"testing"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/supply"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	dbm "github.com/okx/okbchain/libs/tm-db"
	"github.com/okx/okbchain/x/params"
	"github.com/okx/okbchain/x/token/types"
	"github.com/stretchr/testify/require"
)

// CreateParam create okbchain parm for test
func CreateParam(t *testing.T, isCheckTx bool) (sdk.Context, Keeper, *sdk.KVStoreKey, []byte) {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	keyToken := sdk.NewKVStoreKey("token")
	keyLock := sdk.NewKVStoreKey("lock")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeMPT, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyLock, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, nil)

	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, log.NewNopLogger())

	accountKeeper := auth.NewAccountKeeper(
		cdc, // amino codec
		keyAcc,
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)
	blacklistedAddrs := make(map[string]bool)

	bk := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)
	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		types.ModuleName:      nil,
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bank.NewBankKeeperAdapter(bk), maccPerms)
	tk := NewKeeper(bk,
		pk.Subspace(DefaultParamspace),
		auth.FeeCollectorName,
		supplyKeeper,
		keyToken,
		keyLock,
		cdc,
		true, accountKeeper)
	tk.SetParams(ctx, types.DefaultParams())

	return ctx, tk, keyParams, []byte("testToken")
}

func NewTestToken(t *testing.T, ctx sdk.Context, keeper Keeper, bankKeeper bank.Keeper, tokenName string, addrList []sdk.AccAddress) {
	require.NotEqual(t, 0, len(addrList))
	tokenObject := InitTestTokenWithOwner(tokenName, addrList[0])
	keeper.NewToken(ctx, tokenObject)

	initCoins := sdk.NewCoins(sdk.NewCoin(tokenName, sdk.NewInt(100000)))
	for _, addr := range addrList {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}
}

func InitTestToken(name string) types.Token {
	return InitTestTokenWithOwner(name, supply.NewModuleAddress(ModuleName))
}

func InitTestTokenWithOwner(name string, owner sdk.AccAddress) types.Token {
	return types.Token{
		Description:         name,
		Symbol:              name,
		OriginalSymbol:      name,
		WholeName:           name,
		OriginalTotalSupply: sdk.NewDec(0),
		Owner:               owner,
		Type:                1,
		Mintable:            true,
	}
}

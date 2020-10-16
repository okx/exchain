package farm

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
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

var (
	Addrs = createTestAddrs(10)
)

type MockFarmKeeper struct {
	Keeper
	StoreKey     sdk.StoreKey
	SupplyKeeper supply.Keeper
	MountedStore store.MultiStore
	AccKeeper    auth.AccountKeeper
}

func NewMockFarmKeeper(k Keeper, keyStoreKey sdk.StoreKey, sKeeper supply.Keeper,
	ms store.MultiStore, accKeeper auth.AccountKeeper) MockFarmKeeper {
	return MockFarmKeeper{
		k,
		keyStoreKey,
		sKeeper,
		ms,
		accKeeper,
	}
}

func TestBeginBlocker(t *testing.T) {
	ctx, mk := getKeeper(t)
	k := mk.Keeper

	// TODO issue token
	// TODO create swap
	// TODO test farm

	for i := int64(1); i < 10; i++ {
		ctx = ctx.WithBlockHeight(i)
		BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: i}}, k)
	}
}

func getKeeper(t *testing.T) (sdk.Context, MockFarmKeeper) {

	// 0.1 init store key
	keyFarm := sdk.NewKVStoreKey(types.StoreKey)
	tkeyFarm := sdk.NewTransientStoreKey(types.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyToken := sdk.NewKVStoreKey(token.StoreKey)
	keyLock := sdk.NewKVStoreKey(token.KeyLock)
	keySwap := sdk.NewKVStoreKey(swaptypes.StoreKey)

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
	ms.MountStoreWithDB(keySwap, sdk.StoreTypeIAVL, db)
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
	auth.RegisterCodec(cdc)
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
	// fill all the addresses with some coins
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000)))
	for _, addr := range Addrs {
		_, err := bk.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	// 1.4 init supply keeper
	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		types.ModuleName:      {supply.Burner, supply.Minter},
	}
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, sdk.NewDec(1000000000))))
	sk.SetModuleAccount(ctx, feeCollectorAcc)
	sk.SetModuleAccount(ctx, farmAcc)

	// 1.5 init token keeper
	tk := token.NewKeeper(bk, pk.Subspace(token.DefaultParamspace), auth.FeeCollectorName, sk, keyToken, keyLock, cdc, false)

	// 1.6 init swap keeper
	swapKeeper := swap.NewKeeper(sk, tk, cdc, keySwap, pk.Subspace(swaptypes.DefaultParamspace))

	// 1.7 init farm keeper
	fk := NewKeeper(auth.FeeCollectorName, sk, tk, swapKeeper, pk.Subspace(types.DefaultParamspace), keyFarm, cdc)
	fk.SetParams(ctx, types.DefaultParams())
	// 2. init mock keeper
	mk := NewMockFarmKeeper(fk, keyFarm, sk, ms, ak)

	//// 3. init mockApp
	//mApp := mock.NewApp()
	//mApp.Router().AddRoute(types.RouterKey, NewHandler(fk))
	//
	//require.NoError(t, mApp.CompleteSetup(mk.StoreKey))
	return ctx, mk
}

func createTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, err := sdk.AccAddressFromHex(buffer.String())
		if err != nil {
			fmt.Print("error")
		}
		bech := res.String()
		addresses = append(addresses, testAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

// TestAddr is designed for incode address generation
func testAddr(addr string, bech string) sdk.AccAddress {

	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res
}

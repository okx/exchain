package token

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/common/version"
	"github.com/okex/okchain/x/gov"
	govKeeper "github.com/okex/okchain/x/gov/keeper"
	"github.com/okex/okchain/x/params"
	"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var mockBlockHeight int64 = -1

type MockDexApp struct {
	*mock.App

	keyToken   *sdk.KVStoreKey
	keyLock    *sdk.KVStoreKey
	keySupply  *sdk.KVStoreKey
	keyGov     *sdk.KVStoreKey
	keyStaking *sdk.KVStoreKey

	bankKeeper    bank.Keeper
	tokenKeeper   Keeper
	supplyKeeper  supply.Keeper
	govKeeper     gov.Keeper
	stakingKeeper staking.Keeper
}

func registerCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
}

func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{}
	}
}

// initialize the mock application for this module
func getMockDexApp(t *testing.T, numGenAccs int) (mockDexApp *MockDexApp, keeper Keeper, addrs []sdk.AccAddress) {

	mapp := mock.NewApp()
	//mapp.Cdc = makeCodec()
	registerCodec(mapp.Cdc)

	mockDexApp = &MockDexApp{
		App: mapp,

		keyToken:   sdk.NewKVStoreKey("token"),
		keyLock:    sdk.NewKVStoreKey("lock"),
		keySupply:  sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking: sdk.NewKVStoreKey(staking.StoreKey),
	}

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true

	mockDexApp.bankKeeper = bank.NewBaseKeeper(
		mockDexApp.AccountKeeper,
		mockDexApp.ParamsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		types.ModuleName:      {supply.Minter, supply.Burner},
		gov.ModuleName:        nil,
	}
	mockDexApp.supplyKeeper = supply.NewKeeper(mockDexApp.Cdc, mockDexApp.keySupply, mockDexApp.AccountKeeper, mockDexApp.bankKeeper, maccPerms)
	mockDexApp.tokenKeeper = NewKeeper(
		mockDexApp.bankKeeper,
		mockDexApp.ParamsKeeper.Subspace(DefaultParamspace),
		auth.FeeCollectorName,
		mockDexApp.supplyKeeper,
		mockDexApp.keyToken,
		mockDexApp.keyLock,
		mockDexApp.Cdc,
		true)

	handler := NewTokenHandler(mockDexApp.tokenKeeper, version.CurrentProtocolVersion)

	mockDexApp.Router().AddRoute(RouterKey, handler)
	mockDexApp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(mockDexApp.tokenKeeper))

	mockDexApp.SetEndBlocker(getEndBlocker(mockDexApp.tokenKeeper))
	mockDexApp.SetInitChainer(getInitChainer(mockDexApp.App, mockDexApp.supplyKeeper, []exported.ModuleAccountI{feeCollectorAcc}))

	intQuantity := int64(100)
	valTokens := sdk.NewDec(intQuantity)
	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, valTokens),
		sdk.NewDecCoinFromDec(common.TestToken, valTokens),
	}

	genAccs, addrs, _, _ := mock.CreateGenAccounts(numGenAccs, coins)

	// todo: checkTx in mock app
	mockDexApp.SetAnteHandler(nil)

	app := mockDexApp
	require.NoError(t, app.CompleteSetup(
		app.keyToken,
		app.keyLock,
		app.keySupply,
	))
	// TODO: set genesis
	app.BaseApp.NewContext(true, abci.Header{})
	mock.SetGenesis(mockDexApp.App, genAccs)

	for i := 0; i < numGenAccs; i++ {
		mock.CheckBalance(t, app.App, addrs[i], coins)
		mockDexApp.TotalCoinsSupply = mockDexApp.TotalCoinsSupply.Add(coins)
	}

	return mockDexApp, mockDexApp.tokenKeeper, addrs
}

// initialize the mock application for this module
func getMockDexAppEx(t *testing.T, numGenAccs int) (mockDexApp *MockDexApp, keeper Keeper, h sdk.Handler) {

	mapp := mock.NewApp()
	//mapp.Cdc = makeCodec()
	registerCodec(mapp.Cdc)

	mockDexApp = &MockDexApp{
		App: mapp,

		keySupply:  sdk.NewKVStoreKey(supply.StoreKey),
		keyToken:   sdk.NewKVStoreKey("token"),
		keyLock:    sdk.NewKVStoreKey("lock"),
		keyGov:     sdk.NewKVStoreKey(gov.ModuleName),
		keyStaking: sdk.NewKVStoreKey(staking.StoreKey),
	}

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true

	mockDexApp.bankKeeper = bank.NewBaseKeeper(
		mockDexApp.AccountKeeper,
		mockDexApp.ParamsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          nil,
		gov.ModuleName:            nil,
		staking.BondedPoolName:    nil,
		staking.NotBondedPoolName: nil,
	}
	mockDexApp.supplyKeeper = supply.NewKeeper(
		mockDexApp.Cdc,
		mockDexApp.keySupply,
		mockDexApp.AccountKeeper,
		mockDexApp.bankKeeper,
		maccPerms)

	mockDexApp.tokenKeeper = NewKeeper(
		mockDexApp.bankKeeper,
		mockDexApp.ParamsKeeper.Subspace(DefaultParamspace),
		auth.FeeCollectorName,
		mockDexApp.supplyKeeper,
		mockDexApp.keyToken,
		mockDexApp.keyLock,
		mockDexApp.Cdc,
		true)

	stakingKeeper := staking.NewKeeper(
		mockDexApp.Cdc,
		mockDexApp.keyStaking,

		// for staking/distr rollback to cosmos-sdk
		//store.NewKVStoreKey(staking.DelegatorPoolKey),
		//store.NewKVStoreKey(staking.RedelegationKeyM),
		//store.NewKVStoreKey(staking.RedelegationActonKey),
		//store.NewKVStoreKey(staking.UnbondingKey),

		store.NewKVStoreKey(staking.TStoreKey),
		mockDexApp.supplyKeeper,
		mockDexApp.ParamsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	mockDexApp.stakingKeeper = stakingKeeper

	govKp := gov.NewKeeper(
		mockDexApp.Cdc, mockDexApp.keyGov,
		params.Keeper{Keeper: mockDexApp.ParamsKeeper},
		mockDexApp.ParamsKeeper.Subspace(gov.DefaultParamspace),
		mockDexApp.supplyKeeper, stakingKeeper, gov.DefaultCodespace,
		gov.NewRouter(), mockDexApp.bankKeeper, govKeeper.NewProposalHandlerRouter(),
		auth.FeeCollectorName,
	)
	//mockDexApp.tokenKeeper.SetGovKeeper(govKp)
	mockDexApp.govKeeper = govKp

	handler := NewTokenHandler(mockDexApp.tokenKeeper, version.CurrentProtocolVersion)

	mockDexApp.Router().AddRoute(RouterKey, handler)
	mockDexApp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(mockDexApp.tokenKeeper))

	mockDexApp.SetEndBlocker(getEndBlocker(mockDexApp.tokenKeeper))
	mockDexApp.SetInitChainer(getInitChainer(mockDexApp.App, mockDexApp.supplyKeeper, []exported.ModuleAccountI{feeCollectorAcc}))

	intQuantity := int64(10000000)
	valTokens := sdk.NewDec(intQuantity)
	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, valTokens),
		sdk.NewDecCoinFromDec(common.TestToken, valTokens),
	}

	genAccs, _, _, _ := mock.CreateGenAccounts(numGenAccs, coins)

	// todo: checkTx in mock app
	mockDexApp.SetAnteHandler(nil)

	app := mockDexApp
	mockDexApp.MountStores(
		app.keyToken,
		app.keyLock,
		app.keySupply,
		app.keyGov,
		app.keyStaking,
	)

	require.NoError(t, mockDexApp.CompleteSetup())
	mock.SetGenesis(mockDexApp.App, genAccs)
	//app.BaseApp.NewContext(true, abci.Header{})
	mockDexApp.stakingKeeper.SetParams(app.BaseApp.NewContext(true, abci.Header{}), staking.DefaultParams())
	return mockDexApp, mockDexApp.tokenKeeper, handler
}

func getInitChainer(mapp *mock.App, supplyKeeper supply.Keeper,
	blacklistedAddrs []exported.ModuleAccountI) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}
		supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
		return abci.ResponseInitChain{}
	}
}

func getTokenSymbol(ctx sdk.Context, keeper Keeper, prefix string) string {
	store := ctx.KVStore(keeper.tokenStoreKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenKey)
	defer iter.Close()
	for iter.Valid() {
		var token types.Token
		tokenBytes := iter.Value()
		keeper.cdc.MustUnmarshalBinaryBare(tokenBytes, &token)
		if strings.HasPrefix(token.Symbol, prefix) {
			return token.Symbol
		}
		iter.Next()
	}
	return ""
}

type testAccount struct {
	addrKeys    *mock.AddrKeys
	baseAccount types.DecAccount
}

func mockApplyBlock(t *testing.T, app *MockDexApp, txs []auth.StdTx, height int64) sdk.Context {
	mockBlockHeight++
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: height}})

	ctx := app.BaseApp.NewContext(false, abci.Header{})
	//ctx = ctx.WithTxBytes([]byte("90843555124EBF16EB13262400FB8CF639E6A772F437E37A0A141FE640A0B203"))
	param := types.DefaultParams()
	app.tokenKeeper.SetParams(ctx, param)
	for _, tx := range txs {
		app.Deliver(tx)
	}
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
	return ctx
}

func CreateGenAccounts(numAccs int, genCoins sdk.DecCoins) (genAccs []types.DecAccount, atList TestAccounts) {

	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())

		ak := mock.NewAddrKeys(addr, pubKey, privKey)
		testAccount := &testAccount{&ak,
			types.DecAccount{
				Address: addr,
				Coins:   genCoins},
		}
		atList = append(atList, testAccount)

		genAccs = append(genAccs, testAccount.baseAccount)
	}
	return
}

type TestAccounts []*testAccount

func createTokenMsg(t *testing.T, app *MockDexApp, ctx sdk.Context, addr sdk.AccAddress, priKey crypto.PrivKey, tokenMsg sdk.Msg) auth.StdTx {
	accs := app.AccountKeeper.GetAccount(ctx, addr)
	accNum := accs.GetAccountNumber()
	seqNum := accs.GetSequence()

	// todo:
	//tokenIssueMsg.Sender = account.addrKeys.Address
	tx := mock.GenTx([]sdk.Msg{tokenMsg}, []uint64{accNum}, []uint64{seqNum}, priKey)
	app.Check(tx)
	//if !res.IsOK() {
	//	panic("something wrong in checking transaction")
	//}
	return tx
}

type MsgFaked struct {
	Fakeid int
}

func (msg MsgFaked) Route() string { return "token" }

func (msg MsgFaked) Type() string { return "issue" }

// ValidateBasic Implements Msg.
func (msg MsgFaked) ValidateBasic() sdk.Error {
	// check owner
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgFaked) GetSignBytes() []byte {

	return sdk.MustSortJSON([]byte("1"))
}

// GetSigners Implements Msg.
func (msg MsgFaked) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func newFakeMsg() MsgFaked {
	return MsgFaked{
		Fakeid: 0,
	}
}

func TestMsgTokenChown(t *testing.T) {
	//change token owner
	intQuantity := int64(30000)
	// to
	toPriKey := secp256k1.GenPrivKey()
	toPubKey := toPriKey.PubKey()
	toAddr := sdk.AccAddress(toPubKey.Address())
	//init accounts
	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})
	//	fromPriKey := testAccounts[0].addrKeys.PrivKey
	//	fromPubKey := testAccounts[0].addrKeys.PubKey
	fromAddr := testAccounts[0].addrKeys.Address
	//gen app and keepper
	app, keeper, handler := getMockDexAppEx(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	//build context
	ctx := app.BaseApp.NewContext(true, abci.Header{})
	ctx = ctx.WithTxBytes([]byte("90843555124EBF16EB13262400FB8CF639E6A772F437E37A0A141FE640A0B203"))
	var TokenChown []auth.StdTx
	var TokenIssue []auth.StdTx

	//test fake message
	if handler != nil {
		handler(ctx, newFakeMsg())
	}

	//issue token to FromAddress
	tokenIssueMsg := types.NewMsgTokenIssue(common.NativeToken, common.NativeToken, common.NativeToken, "okcoin", "1000", testAccounts[0].baseAccount.Address, true)
	TokenIssue = append(TokenIssue, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))

	//test error supply coin issue(TotalSupply > (9*1e10))
	MsgErrorSupply := types.NewMsgTokenIssue("okc", "okc", "okc", "okccc", strconv.FormatInt(int64(10*1e10), 10), testAccounts[0].baseAccount.Address, true)
	TokenIssue = append(TokenIssue, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, MsgErrorSupply))

	//test error tokenDesc (length > 256)
	MsgErrorName := types.NewMsgTokenIssue(`ok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-b
ok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-b
ok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-b
ok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-bok-b`,
		common.NativeToken, common.NativeToken, "okcoin", "2100", testAccounts[0].baseAccount.Address, true)
	TokenIssue = append(TokenIssue, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, MsgErrorName))

	ctx = mockApplyBlock(t, app, TokenIssue, 3)

	require.NotNil(t, handleMsgTokenIssue(ctx, keeper, MsgErrorSupply, nil))
	//require.NotNil(t, handleMsgTokenIssue(ctx, keeper, MsgErrorName, nil))

	//test if zzb is not exist
	invalidmsg := types.NewMsgTransferOwnership(fromAddr, toAddr, "zzb")
	bSig, err := toPriKey.Sign(invalidmsg.GetSignBytes())
	require.NoError(t, err)
	invalidmsg.ToSignature.PubKey = toPubKey
	invalidmsg.ToSignature.Signature = bSig
	TokenChown = append(TokenChown, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, invalidmsg))

	//test if zzb is not exist
	tokenNotExist := types.NewMsgTransferOwnership(fromAddr, toAddr, "zzb")
	bSig, err = toPriKey.Sign(tokenNotExist.GetSignBytes())
	require.NoError(t, err)
	tokenNotExist.ToSignature.PubKey = toPubKey
	tokenNotExist.ToSignature.Signature = bSig
	TokenChown = append(TokenChown, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenNotExist))

	//test AddTokenSuffix->ValidSymbol
	AddTokenSuffix(ctx, keeper, "notexist")

	//normal test
	symbName := "okb-b85" //AddTokenSuffix(ctx,keeper,common.NativeToken)
	//change owner from F to T
	tokenChownMsg := types.NewMsgTransferOwnership(fromAddr, toAddr, symbName)
	bSig, err = toPriKey.Sign(tokenChownMsg.GetSignBytes())
	require.NoError(t, err)
	tokenChownMsg.ToSignature.PubKey = toPubKey
	tokenChownMsg.ToSignature.Signature = bSig
	TokenChown = append(TokenChown, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenChownMsg))

	ctx = mockApplyBlock(t, app, TokenChown, 4)
}

func TestUpdateUserTokenRelationship(t *testing.T) {
	toPriKey := secp256k1.GenPrivKey()
	toPubKey := toPriKey.PubKey()
	toAddr := sdk.AccAddress(toPubKey.Address())

	intQuantity := int64(30000)
	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.BaseApp.NewContext(true, abci.Header{})
	ctx = ctx.WithTxBytes([]byte("90843555124EBF16EB13262400FB8CF639E6A772F437E37A0A141FE640A0B203"))

	var tokenIssue []auth.StdTx

	totalSupplyStr := "500"
	tokenIssueMsg := types.NewMsgTokenIssue("bnb", "", "bnb", "binance coin", totalSupplyStr, testAccounts[0].baseAccount.Address, true)
	tokenIssue = append(tokenIssue, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))

	ctx = mockApplyBlock(t, app, tokenIssue, 3)

	tokens := keeper.GetUserTokensInfo(ctx, testAccounts[0].baseAccount.Address)
	require.EqualValues(t, 1, len(tokens))

	tokenName := getTokenSymbol(ctx, keeper, "bnb")
	// ===============

	var TokenChown []auth.StdTx

	//test if zzb is not exist
	chownMsg := types.NewMsgTransferOwnership(testAccounts[0].baseAccount.Address, toAddr, tokenName)
	bSig, err := toPriKey.Sign(chownMsg.GetSignBytes())
	require.NoError(t, err)
	chownMsg.ToSignature.PubKey = toPubKey
	chownMsg.ToSignature.Signature = bSig
	TokenChown = append(TokenChown, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, chownMsg))

	ctx = mockApplyBlock(t, app, TokenChown, 4)

	tokens = keeper.GetUserTokensInfo(ctx, testAccounts[0].baseAccount.Address)
	require.EqualValues(t, 0, len(tokens))
}

func TestCreateTokenIssue(t *testing.T) {
	intQuantity := int64(30000)
	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.BaseApp.NewContext(true, abci.Header{})
	ctx = ctx.WithTxBytes([]byte("90843555124EBF16EB13262400FB8CF639E6A772F437E37A0A141FE640A0B203"))

	var tokenIssue []auth.StdTx

	totalSupply := int64(500)
	totalSupplyStr := "500"
	tokenIssueMsg := types.NewMsgTokenIssue("bnb", "", "bnb", "binance coin", totalSupplyStr, testAccounts[0].baseAccount.Address, true)
	tokenIssue = append(tokenIssue, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))

	// not valid symbol
	//tokenIssueMsg = types.NewMsgTokenIssue("bnba123451fadfasdf", "bnba123451fadfasdf", "bnba123451fadfasdf", totalSupply, testAccounts[0].baseAccount.Address, true)
	//tokenIssue = append(tokenIssue, createTokenMsg(t, app, ctx, testAccounts[0], tokenIssueMsg))

	// Total exceeds the upper limit
	tokenIssueMsg = types.NewMsgTokenIssue("btc", "btc", "btc", "bitcoin", strconv.FormatInt(types.TotalSupplyUpperbound+1, 10), testAccounts[0].baseAccount.Address, true)
	tokenIssue = append(tokenIssue, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))

	// not enough okbs
	tokenIssueMsg = types.NewMsgTokenIssue("xmr", "xmr", "xmr", "Monero", totalSupplyStr, testAccounts[0].baseAccount.Address, true)
	tokenIssue = append(tokenIssue, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))

	ctx = mockApplyBlock(t, app, tokenIssue, 3)

	tokenName := getTokenSymbol(ctx, keeper, "bnb")
	//feeIssue, err := sdk.NewDecFromStr(DefaultFeeIssue)
	//require.EqualValues(t, nil, err)
	feeIssue := keeper.GetParams(ctx).FeeIssue.Amount
	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokenName, sdk.NewDec(totalSupply)),
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity).Sub(feeIssue)),
	}
	require.EqualValues(t, coins, app.AccountKeeper.GetAccount(ctx, testAccounts[0].addrKeys.Address).GetCoins())
	tokenStoreKeyNum, lockStoreKeyNum := keeper.GetNumKeys(ctx)
	require.Equal(t, int64(3), tokenStoreKeyNum)
	require.Equal(t, int64(0), lockStoreKeyNum)
	//require.Equal(t, int64(0), tokenPairStoreKeyNum)

	tokenInfo := keeper.GetTokenInfo(ctx, tokenName)
	require.EqualValues(t, sdk.MustNewDecFromStr("500"), tokenInfo.OriginalTotalSupply)
	require.EqualValues(t, sdk.MustNewDecFromStr("500"), tokenInfo.TotalSupply)
}

func TestCreateTokenBurn(t *testing.T) {
	intQuantity := int64(20011)

	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	_, testAccounts2 := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.NewContext(true, abci.Header{})
	var tokenMsgs []auth.StdTx

	tokenIssueMsg := types.NewMsgTokenIssue("btc", "btc", "btc", "bitcoin", "1000", testAccounts[0].baseAccount.Address, true)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 3)

	tokenMsgs = tokenMsgs[:0]

	burnNum := "100"
	//mockToken(app.tokenKeeper, ctx, testAccounts[0].baseAccount.Address, intQuantity)

	_, err := sdk.ParseDecCoin("-10000btc")
	require.Error(t, err)

	decCoin, err := sdk.ParseDecCoin("10000btc")
	require.Nil(t, err)
	// total exceeds the upper limit
	tokenBurnMsg := types.NewMsgTokenBurn(decCoin, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenBurnMsg))
	//mockApplyBlock(t, app, tokenMsgs)

	// not the token's owner
	tokenBurnMsg = types.NewMsgTokenBurn(decCoin, testAccounts2[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenBurnMsg))

	tokenSymbol := getTokenSymbol(ctx, keeper, "btc")

	decCoin, err = sdk.ParseDecCoin(burnNum + tokenSymbol)
	require.Nil(t, err)
	// normal case
	tokenBurnMsg = types.NewMsgTokenBurn(decCoin, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenBurnMsg))

	decCoin, err = sdk.ParseDecCoin(burnNum + "btc")
	require.Nil(t, err)
	// not enough fees
	tokenBurnMsg = types.NewMsgTokenBurn(decCoin, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenBurnMsg))

	ctx = mockApplyBlock(t, app, tokenMsgs, 4)

	//fee, err := sdk.NewDecFromStr("0.0125")
	fee, err := sdk.NewDecFromStr("0.0")
	require.Nil(t, err)
	validTxNum := sdk.NewDec(2)
	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokenSymbol, sdk.NewDec(900)),
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1).Add(fee.Mul(validTxNum))),
	}

	require.EqualValues(t, coins, app.AccountKeeper.GetAccount(ctx, testAccounts[0].addrKeys.Address).GetCoins())
}

func TestCreateTokenMint(t *testing.T) {
	intQuantity := int64(42001)

	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	_, testAccounts2 := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.NewContext(true, abci.Header{})
	var tokenMsgs []auth.StdTx

	tokenIssueMsg := types.NewMsgTokenIssue("btc", "btc", "btc", "bitcoin", "1000", testAccounts[0].baseAccount.Address, true)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 3)
	tokenMsgs = tokenMsgs[:0]

	tokenIssueMsg = types.NewMsgTokenIssue("xmr", "xmr", "xmr", "monero", "1000", testAccounts[0].baseAccount.Address, false)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 4)

	var mintNum int64 = 1000
	// normal case
	btcTokenSymbol := getTokenSymbol(ctx, keeper, "btc")
	decCoin := sdk.NewDecCoinFromDec(btcTokenSymbol, sdk.NewDec(mintNum))
	tokenMintMsg := types.NewMsgTokenMint(decCoin, testAccounts[0].baseAccount.Address)

	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenMintMsg))

	// Total exceeds the upper limit
	decCoin.Amount = sdk.NewDec(types.TotalSupplyUpperbound)
	tokenMintMsg = types.NewMsgTokenMint(decCoin, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenMintMsg))

	// not the token's owner
	decCoin.Amount = sdk.NewDec(mintNum)
	tokenMintMsg = types.NewMsgTokenMint(decCoin, testAccounts2[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenMintMsg))

	// token not mintable
	xmrTokenSymbol := getTokenSymbol(ctx, keeper, "xmr")
	decCoin.Denom = xmrTokenSymbol
	decCoin.Amount = sdk.NewDec(mintNum)
	tokenMintMsg = types.NewMsgTokenMint(decCoin, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenMintMsg))

	// not enough fees
	decCoin.Denom = btcTokenSymbol
	decCoin.Amount = sdk.NewDec(mintNum)
	tokenMintMsg = types.NewMsgTokenMint(decCoin, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenMintMsg))

	ctx = mockApplyBlock(t, app, tokenMsgs, 5)

	//validTxNum := sdk.NewInt(2)
	coins := sdk.MustParseCoins(btcTokenSymbol, "2000")
	//coins = append(coins, sdk.MustParseCoins(common.NativeToken, "1.0375")...)
	coins = append(coins, sdk.MustParseCoins(common.NativeToken, "1.0")...)
	coins = append(coins, sdk.MustParseCoins(xmrTokenSymbol, "1000")...)

	require.EqualValues(t, coins, app.AccountKeeper.GetAccount(ctx, testAccounts[0].addrKeys.Address).GetCoins())
}

func TestCreateMsgTokenSend(t *testing.T) {
	intQuantity := int64(100000)

	genAccs, testAccounts := CreateGenAccounts(2,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.NewContext(true, abci.Header{})
	var tokenMsgs []auth.StdTx

	tokenIssueMsg := types.NewMsgTokenIssue("btc", "btc", "btc", "bitcoin", "1000", testAccounts[0].baseAccount.Address, true)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 3)
	tokenMsgs = tokenMsgs[:0]

	tokenName := getTokenSymbol(ctx, keeper, "btc")
	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokenName, sdk.NewDec(100)),
	}
	tokenSendMsg := types.NewMsgTokenSend(testAccounts[0].baseAccount.Address, testAccounts[1].baseAccount.Address, coins)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenSendMsg))

	coins = sdk.DecCoins{
		sdk.NewDecCoinFromDec("btc", sdk.NewDec(10000)),
	}
	// not enough coins
	tokenSendMsg = types.NewMsgTokenSend(testAccounts[0].baseAccount.Address, testAccounts[1].baseAccount.Address, coins)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenSendMsg))

	ctx = mockApplyBlock(t, app, tokenMsgs, 4)

	accounts := app.AccountKeeper.GetAllAccounts(ctx)
	for _, acc := range accounts {
		if acc.GetAddress().Equals(testAccounts[0].baseAccount.Address) {
			senderCoins := sdk.MustParseCoins(tokenName, "900")
			//senderCoins = append(senderCoins, sdk.MustParseCoins(common.NativeToken, "80000.0125")...)
			senderCoins = append(senderCoins, sdk.MustParseCoins(common.NativeToken, "80000")...)
			require.EqualValues(t, senderCoins, acc.GetCoins())
		} else if acc.GetAddress().Equals(testAccounts[1].baseAccount.Address) {
			receiverCoins := sdk.MustParseCoins(tokenName, "100")
			receiverCoins = append(receiverCoins, sdk.MustParseCoins(common.NativeToken, "100000")...)
			require.EqualValues(t, receiverCoins, acc.GetCoins())
		}
	}

	// len(MsgTokenSend.Amount) > 1
	tokenMsgs = tokenMsgs[:0]
	coins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokenName, sdk.NewDec(100)),
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}
	tokenSendMsg = types.NewMsgTokenSend(testAccounts[0].baseAccount.Address, testAccounts[1].baseAccount.Address, coins)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenSendMsg))

	ctx = mockApplyBlock(t, app, tokenMsgs, 5)

	accounts = app.AccountKeeper.GetAllAccounts(ctx)
	for _, acc := range accounts {
		if acc.GetAddress().Equals(testAccounts[0].baseAccount.Address) {
			senderCoins := sdk.MustParseCoins(tokenName, "800")
			//senderCoins = append(senderCoins, sdk.MustParseCoins(common.NativeToken, "79900.005")...)
			senderCoins = append(senderCoins, sdk.MustParseCoins(common.NativeToken, "79899.9925")...)
			require.EqualValues(t, senderCoins.String(), acc.GetCoins().String())
		} else if acc.GetAddress().Equals(testAccounts[1].baseAccount.Address) {
			receiverCoins := sdk.MustParseCoins(tokenName, "200")
			receiverCoins = append(receiverCoins, sdk.MustParseCoins(common.NativeToken, "100100")...)
			require.EqualValues(t, receiverCoins, acc.GetCoins())
		}
	}
}

func TestChargeMultiCoinsFee(t *testing.T) {
	intQuantity := int64(1)
	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))
	ctx := app.NewContext(true, abci.Header{})
	keeper.SetParams(ctx, types.DefaultParams())

	goodRes := sdk.Result{}
	testCases := []struct {
		coinNum    int
		from       sdk.AccAddress
		feeCharged sdk.DecCoins
		res        sdk.Result
	}{
		{1, testAccounts[0].baseAccount.Address,
			sdk.ZeroFee().ToCoins(), goodRes},
		{2, testAccounts[0].baseAccount.Address,
			sdk.MustParseCoins(sdk.DefaultBondDenom, "0.0075"), goodRes},
		{3, testAccounts[0].baseAccount.Address,
			sdk.MustParseCoins(sdk.DefaultBondDenom, "0.0175"), goodRes},
		{4, testAccounts[0].baseAccount.Address,
			sdk.MustParseCoins(sdk.DefaultBondDenom, "0.0275"), goodRes},
		{100, testAccounts[0].baseAccount.Address,
			sdk.MustParseCoins(sdk.DefaultBondDenom, "0.9875"),
			sdk.ErrInsufficientCoins("insufficient fee coins(need 0.98750000okt)").Result()},
	}

	for _, tc := range testCases {
		feeCharged, res := chargeMultiCoinsFee(ctx, keeper, tc.from, tc.coinNum)
		require.Equal(t, tc.feeCharged, feeCharged)
		require.Equal(t, tc.res, res)
	}
}

func TestCreateMsgMultiSend(t *testing.T) {
	intQuantity := int64(100000)

	genAccs, testAccounts := CreateGenAccounts(2,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.NewContext(true, abci.Header{})
	var tokenMsgs []auth.StdTx

	tokenIssueMsg := types.NewMsgTokenIssue("btc", "btc", "btc", "bitcoin", "1000", testAccounts[0].baseAccount.Address, true)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 3)
	tokenMsgs = tokenMsgs[:0]

	btcSymbol := getTokenSymbol(ctx, keeper, "btc")

	multiSendStr := `[{"to":"` + testAccounts[1].baseAccount.Address.String() + `","amount":"1okt,2` + btcSymbol + `"}]`
	transfers, err := types.StrToTransfers(multiSendStr)
	require.Nil(t, err)
	multiSend := types.NewMsgMultiSend(testAccounts[0].baseAccount.Address, transfers)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, multiSend))

	// not enough coins
	multiSendStr = `[{"to":"` + testAccounts[1].baseAccount.Address.String() + `","amount":"1okt,2000` + btcSymbol + `"}]`
	transfers, err = types.StrToTransfers(multiSendStr)
	require.Nil(t, err)
	multiSend = types.NewMsgMultiSend(testAccounts[0].baseAccount.Address, transfers)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, multiSend))

	ctx = mockApplyBlock(t, app, tokenMsgs, 4)

	accounts := app.AccountKeeper.GetAllAccounts(ctx)
	for _, acc := range accounts {
		if acc.GetAddress().Equals(testAccounts[0].baseAccount.Address) {
			senderCoins := sdk.MustParseCoins(btcSymbol, "998")
			//senderCoins = append(senderCoins, sdk.MustParseCoins(common.NativeToken, "79999.005")...)
			senderCoins = append(senderCoins, sdk.MustParseCoins(common.NativeToken, "79998.99250000")...)
			require.EqualValues(t, senderCoins, acc.GetCoins())
		} else if acc.GetAddress().Equals(testAccounts[1].baseAccount.Address) {
			receiverCoins := sdk.MustParseCoins(btcSymbol, "2")
			receiverCoins = append(receiverCoins, sdk.MustParseCoins(common.NativeToken, "100001")...)
			require.EqualValues(t, receiverCoins, acc.GetCoins())
		}
	}
}

func TestCreateMsgTokenModify(t *testing.T) {
	intQuantity := int64(100000)

	genAccs, testAccounts := CreateGenAccounts(2,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.NewContext(true, abci.Header{})
	var tokenMsgs []auth.StdTx

	tokenIssueMsg := types.NewMsgTokenIssue("btc", "btc", "btc", "bitcoin", "1000", testAccounts[0].baseAccount.Address, true)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenIssueMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 3)

	tokenMsgs = tokenMsgs[:0]
	btcTokenSymbol := getTokenSymbol(ctx, keeper, "btc")

	// normal case
	tokenEditMsg := types.NewMsgTokenModify(btcTokenSymbol, "desc0", "whole name0", true, true, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 4)
	token := keeper.GetTokenInfo(ctx, btcTokenSymbol)
	require.EqualValues(t, "desc0", token.Description)
	require.EqualValues(t, "whole name0", token.WholeName)

	tokenMsgs = tokenMsgs[:0]
	tokenEditMsg = types.NewMsgTokenModify(btcTokenSymbol, "desc1", "whole name1", false, true, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 5)
	token = keeper.GetTokenInfo(ctx, btcTokenSymbol)
	require.EqualValues(t, "desc0", token.Description)
	require.EqualValues(t, "whole name1", token.WholeName)

	tokenMsgs = tokenMsgs[:0]
	tokenEditMsg = types.NewMsgTokenModify(btcTokenSymbol, "desc2", "whole name2", true, false, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 6)
	token = keeper.GetTokenInfo(ctx, btcTokenSymbol)
	require.EqualValues(t, "desc2", token.Description)
	require.EqualValues(t, "whole name1", token.WholeName)

	tokenMsgs = tokenMsgs[:0]
	tokenEditMsg = types.NewMsgTokenModify(btcTokenSymbol, "desc3", "whole name2", false, false, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 7)
	token = keeper.GetTokenInfo(ctx, btcTokenSymbol)
	require.EqualValues(t, "desc2", token.Description)
	require.EqualValues(t, "whole name1", token.WholeName)

	// error case
	tokenMsgs = tokenMsgs[:0]
	tokenEditMsg = types.NewMsgTokenModify("btcTokenSymbol", "desc4", "whole name4", true, true, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 8)

	tokenMsgs = tokenMsgs[:0]
	tokenEditMsg = types.NewMsgTokenModify(btcTokenSymbol, "desc5", "whole name5", true, true, testAccounts[1].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 9)

	tokenMsgs = tokenMsgs[:0]
	tokenEditMsg = types.NewMsgTokenModify(btcTokenSymbol, "desc6", "whole nasiangrueinvowfoij;oeasifnroeinagoirengodd   me6", true, true, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 10)

	tokenMsgs = tokenMsgs[:0]
	tokenEditMsg = types.NewMsgTokenModify(btcTokenSymbol, `bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234`, "whole name7", true, true, testAccounts[0].baseAccount.Address)
	tokenMsgs = append(tokenMsgs, createTokenMsg(t, app, ctx, testAccounts[0].baseAccount.Address, testAccounts[0].addrKeys.PrivKey, tokenEditMsg))
	ctx = mockApplyBlock(t, app, tokenMsgs, 11)

	token = keeper.GetTokenInfo(ctx, btcTokenSymbol)
	require.EqualValues(t, "desc2", token.Description)
	require.EqualValues(t, "whole name1", token.WholeName)
}

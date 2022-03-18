package token

import (
	"testing"

	"github.com/okex/exchain/x/token/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mock"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/common"
	"github.com/stretchr/testify/require"
)

func TestQueryOrder(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	_, testAccounts := CreateGenAccounts(2,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})

	token := types.Token{
		Description:         "okblockchain coin",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		OriginalTotalSupply: sdk.NewDec(1000000000),
		Owner:               testAccounts[0].baseAccount.Address,
		Mintable:            true,
	}

	coins := sdk.NewCoins(sdk.NewDecCoinFromDec(token.Symbol, token.OriginalTotalSupply))
	err := keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	keeper.NewToken(ctx, token)

	querier := NewQuerier(keeper)
	path := []string{types.QueryInfo, ""}
	res, err := querier(ctx, path, abci.RequestQuery{})
	require.NotNil(t, err)
	require.Equal(t, []byte(nil), res)
	path = []string{types.QueryInfo, common.NativeToken}
	res, err = querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)

	var token2 types.Token
	keeper.cdc.MustUnmarshalJSON(res, &token2)
	require.EqualValues(t, token, token2)
}

func TestQueryTokens(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	_, testAccounts := CreateGenAccounts(2,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})

	token := types.Token{
		Description:         "okblockchain coin",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		OriginalTotalSupply: sdk.NewDec(1000000000),
		Owner:               testAccounts[0].baseAccount.Address,
		Mintable:            true,
	}

	coins := sdk.NewCoins(sdk.NewDecCoinFromDec(token.Symbol, token.OriginalTotalSupply))
	err := keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	keeper.NewToken(ctx, token)

	var originTokens []types.Token
	originTokens = append(originTokens, token)

	querier := NewQuerier(keeper)

	path := []string{types.QueryTokens}
	res, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)

	var tokens []types.Token
	keeper.cdc.MustUnmarshalJSON(res, &tokens)
	require.EqualValues(t, originTokens, tokens)

	//token V2
	path = []string{types.QueryTokenV2, common.NativeToken}
	res, err = querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	require.Panics(t, func() { keeper.cdc.MustUnmarshalJSON(res, &tokens) })
	require.EqualValues(t, originTokens, tokens)
	//no symbol
	path = []string{types.QueryTokenV2, ""}
	res, err = querier(ctx, path, abci.RequestQuery{})
	//require.EqualValues(t, sdk.CodeInvalidCoins, err.Code())
	require.NotNil(t, err)
	require.Equal(t, []byte(nil), res)

	//tokens V2
	path = []string{types.QueryTokensV2}
	res, err = querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)

	keeper.cdc.MustUnmarshalJSON(res, &tokens)
	require.EqualValues(t, originTokens, tokens)

	//query with address
	path = []string{types.QueryTokens, testAccounts[0].baseAccount.Address.String()}
	res, err = querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)

	keeper.cdc.MustUnmarshalJSON(res, &tokens)
	require.EqualValues(t, originTokens, tokens)

	//query with invalid address
	token = types.Token{
		Description:         "okblockchain coin",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		OriginalTotalSupply: sdk.NewDec(1000000000),
		//TotalSupply:         sdk.NewDec(1000000000),
		Owner:    []byte("abc"),
		Mintable: true,
	}

	keeper.NewToken(ctx, token)
	path = []string{types.QueryTokens, "abc"}
	res, err = querier(ctx, path, abci.RequestQuery{})
	//require.EqualValues(t, sdk.CodeInvalidAddress, err.Code())
	require.NotNil(t, err)
	require.Equal(t, []byte(nil), res)
}

func TestQueryUserTokens(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	_, testAccounts := CreateGenAccounts(2,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})

	token := types.Token{
		Description:         "okblockchain coin",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		OriginalTotalSupply: sdk.NewDec(1000000000),
		Owner:               testAccounts[0].baseAccount.Address,
		Mintable:            true,
	}

	coins := sdk.NewCoins(sdk.NewDecCoinFromDec(token.Symbol, token.OriginalTotalSupply))
	err := keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	keeper.NewToken(ctx, token)

	var originTokens []types.Token
	originTokens = append(originTokens, token)

	querier := NewQuerier(keeper)

	path := []string{types.QueryTokens, testAccounts[0].baseAccount.Address.String()}
	res, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)

	var tokens []types.Token
	keeper.cdc.MustUnmarshalJSON(res, &tokens)
	require.EqualValues(t, originTokens, tokens)
}

func TestQueryCurrency(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	_, testAccounts := CreateGenAccounts(2,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})

	token := types.Token{
		Description:         "okblockchain coin",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		OriginalTotalSupply: sdk.NewDec(1000000000),
		Owner:               testAccounts[0].baseAccount.Address,
		Mintable:            true,
	}

	keeper.NewToken(ctx, token)
	keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDec(1000000000)))

	//var originTokens []types.Token
	originalTokens := []types.Currency{
		{
			Description: "okblockchain coin",
			Symbol:      common.NativeToken,
			TotalSupply: sdk.NewDec(1000000000),
		},
	}

	querier := NewQuerier(keeper)
	path := []string{types.QueryCurrency}
	res, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)

	var currency []types.Currency
	keeper.cdc.MustUnmarshalJSON(res, &currency)
	require.EqualValues(t, originalTokens, currency)
}

func TestQueryAccount(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	token := types.Token{
		Description:         "okblockchain coin",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		OriginalTotalSupply: sdk.NewDec(1000000000),
		Owner:               testAccounts[0].baseAccount.Address,
		Mintable:            true,
	}

	originalCoinsInfo := types.CoinsInfo{
		types.CoinInfo{
			Symbol:    common.NativeToken,
			Available: "1000000000.000000000000000000",
			Locked:    "0",
		},
	}

	keeper.NewToken(ctx, token)

	querier := NewQuerier(keeper)
	path := []string{types.QueryAccount, testAccounts[0].baseAccount.Address.String()}
	var accountParam types.AccountParam
	accountParam.Symbol = common.NativeToken
	accountParam.Show = "partial"

	bz, err := keeper.cdc.MarshalJSON(accountParam)
	require.Nil(t, err)
	res, err := querier(ctx, path, abci.RequestQuery{Data: nil})
	require.Error(t, err)
	require.Equal(t, []byte(nil), res)
	res, err = querier(ctx, path, abci.RequestQuery{Data: bz})
	require.Nil(t, err)

	//var coinsInfo CoinsInfo
	var accountResponse types.AccountResponse
	keeper.cdc.MustUnmarshalJSON(res, &accountResponse)

	require.EqualValues(t, originalCoinsInfo, accountResponse.Currencies)

	path = []string{types.QueryAccount, "notexist"}
	_, err = querier(ctx, path, abci.RequestQuery{})
	require.Error(t, err)

	path = []string{"not_exist", "notexist"}
	_, err = querier(ctx, path, abci.RequestQuery{})
	require.Error(t, err)

	//account V2
	var param types.AccountParamV2
	path = []string{types.QueryAccountV2, testAccounts[0].baseAccount.Address.String()}
	param.Currency = ""
	param.HideZero = "partial"

	bz = keeper.cdc.MustMarshalJSON(param)
	_, err = querier(ctx, path, abci.RequestQuery{Data: nil})
	require.Error(t, err)
	//Currency nil
	_, err = querier(ctx, path, abci.RequestQuery{Data: bz})
	require.Nil(t, err)
	//hide no
	param.Currency = ""
	param.HideZero = "no"
	bz = keeper.cdc.MustMarshalJSON(param)
	_, err = querier(ctx, path, abci.RequestQuery{Data: bz})
	require.Nil(t, err)

	param.Currency = common.NativeToken
	param.HideZero = "no"
	bz = keeper.cdc.MustMarshalJSON(param)
	_, err = querier(ctx, path, abci.RequestQuery{Data: bz})
	require.Nil(t, err)

	//err address
	path = []string{types.QueryAccountV2, "abc"}
	_, err = querier(ctx, path, abci.RequestQuery{Data: nil})
	require.Error(t, err)
}

func TestQueryAccount_ShowAll(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	tokens := []types.Token{
		{
			Description:         "okblockchain coin",
			Symbol:              common.NativeToken,
			OriginalSymbol:      common.NativeToken,
			OriginalTotalSupply: sdk.NewDec(1000000000),
			Owner:               testAccounts[0].baseAccount.Address,
			Mintable:            true,
		},
		{
			Description:         "not_exist",
			Symbol:              "xxb",
			OriginalSymbol:      "xxb",
			OriginalTotalSupply: sdk.ZeroDec(),
			Owner:               testAccounts[0].baseAccount.Address,
			Mintable:            false,
		},
	}
	for _, token := range tokens {
		keeper.NewToken(ctx, token)
	}

	originalCoinsInfo := types.CoinsInfo{
		types.CoinInfo{
			Symbol:    common.NativeToken,
			Available: "1000000000.000000000000000000",
			Locked:    "0",
		},
		types.CoinInfo{
			Symbol:    "xxb",
			Available: "0",
			Locked:    "0",
		},
	}

	querier := NewQuerier(keeper)

	accountParam := types.AccountParam{
		Symbol: "",
		Show:   "all",
	}
	bz := keeper.cdc.MustMarshalJSON(accountParam)
	path := []string{types.QueryAccount, testAccounts[0].baseAccount.Address.String()}
	res, err := querier(ctx, path, abci.RequestQuery{Data: bz})
	require.Nil(t, err)

	var accountResponse types.AccountResponse
	keeper.cdc.MustUnmarshalJSON(res, &accountResponse)
	require.EqualValues(t, originalCoinsInfo, accountResponse.Currencies)
}

func TestQueryParameters(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	genAccs, _ := CreateGenAccounts(1,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	querier := NewQuerier(keeper)
	path := []string{types.QueryParameters}

	res, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	//res, err := queryParameters(ctx, keeper)
	//require.Nil(t, err)

	var actualParams types.Params
	err1 := keeper.cdc.UnmarshalJSON(res, &actualParams)
	require.Nil(t, err1)
	require.EqualValues(t, actualParams, params)
}

func TestQueryKeysNum(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	genAccs, _ := CreateGenAccounts(1,
		sdk.SysCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	path := []string{types.QueryKeysNum}
	querier := NewQuerier(keeper)
	res, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)

	keyNums := map[string]int64{"token": 1,
		"lock":      0,
		"tokenPair": 1,
	}

	resMap := make(map[string]int64)
	err1 := keeper.cdc.UnmarshalJSON(res, &resMap)
	require.Nil(t, err1)
	require.EqualValues(t, keyNums, keyNums)
}

func TestCreateParam(t *testing.T) {
	ctx, kepper, kv, data := CreateParam(t, true)
	require.NotNil(t, ctx)
	require.NotNil(t, kepper)
	require.NotNil(t, kv)
	require.EqualValues(t, "testToken", string(data))
}

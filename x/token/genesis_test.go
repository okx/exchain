package token

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestDefault(t *testing.T) {
	genesisState := defaultGenesisState()
	err := validateGenesis(genesisState)
	require.NoError(t, err)
	defaultGenesisStateOKT()
}

func TestInitGenesis(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	//ctx, keeper, _, _ := CreateParam(t, false)
	keeper.SetParams(ctx, types.DefaultParams())
	params := keeper.GetParams(ctx)

	var tokens []types.Token
	tokens = append(tokens, defaultGenesisStateOKT())

	var lockedCoins []types.AccCoins
	decCoin := sdk.NewDecCoinFromDec(tokens[0].Symbol, sdk.NewDec(1234))
	lockedCoins = append(lockedCoins, types.AccCoins{
		Acc:   tokens[0].Owner,
		Coins: sdk.DecCoins{decCoin},
	})

	var lockedFees []types.AccCoins
	lockedFees = append(lockedFees, types.AccCoins{
		Acc:   tokens[0].Owner,
		Coins: sdk.DecCoins{decCoin},
	})

	initedGenesis := GenesisState{
		Params:       params,
		Tokens:       tokens,
		LockedAssets: lockedCoins,
		LockedFees:   lockedFees,
	}

	coins := sdk.NewDecCoinsFromDec(tokens[0].Symbol, tokens[0].OriginalTotalSupply)

	err := keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	initGenesis(ctx, keeper, initedGenesis)
	require.Equal(t, initedGenesis.Params, keeper.GetParams(ctx))
	require.Equal(t, initedGenesis.Tokens, keeper.GetTokensInfo(ctx))
	require.Equal(t, initedGenesis.LockedAssets, keeper.GetAllLockedCoins(ctx))
	require.Equal(t, uint64(len(initedGenesis.Tokens)), keeper.getTokenNum(ctx))
	require.Equal(t, initedGenesis.Tokens[0], keeper.GetUserTokensInfo(ctx, initedGenesis.Tokens[0].Owner)[0])
	var actualLockeedFees []types.AccCoins
	keeper.IterateLockedFees(ctx, func(acc sdk.AccAddress, coins sdk.DecCoins) bool {
		actualLockeedFees = append(actualLockeedFees,
			types.AccCoins{
				Acc:   acc,
				Coins: coins,
			})
		return false
	})
	require.Equal(t, initedGenesis.LockedFees, actualLockeedFees)

	exportGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, initedGenesis.Params, exportGenesis.Params)
	require.Equal(t, initedGenesis.Tokens, exportGenesis.Tokens)
	require.Equal(t, initedGenesis.LockedAssets, exportGenesis.LockedAssets)
	require.Equal(t, initedGenesis.LockedFees, exportGenesis.LockedFees)

	newMapp, newKeeper, _ := getMockDexApp(t, 0)
	newMapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	newCtx := newMapp.BaseApp.NewContext(false, abci.Header{})

	exportGenesis.Tokens[0].OriginalTotalSupply = sdk.NewDec(66666)
	//exportGenesis.Tokens[0].TotalSupply = sdk.NewDec(66666)
	decCoin.Denom = tokens[0].Symbol
	decCoin.Amount = sdk.NewDec(7777)
	exportGenesis.LockedAssets[0].Coins = sdk.DecCoins{decCoin}
	exportGenesis.LockedFees[0].Coins = sdk.DecCoins{decCoin}

	coins = sdk.NewCoins(sdk.NewDecCoinFromDec(exportGenesis.Tokens[0].Symbol, exportGenesis.Tokens[0].OriginalTotalSupply))
	err = newKeeper.supplyKeeper.MintCoins(newCtx, types.ModuleName, coins)
	require.NoError(t, err)

	initGenesis(newCtx, newKeeper, exportGenesis)
	require.Equal(t, exportGenesis.Params, newKeeper.GetParams(newCtx))
	require.Equal(t, exportGenesis.Tokens, newKeeper.GetTokensInfo(newCtx))
	require.Equal(t, exportGenesis.LockedAssets, newKeeper.GetAllLockedCoins(newCtx))
	require.Equal(t, uint64(len(exportGenesis.Tokens)), newKeeper.getTokenNum(newCtx))
	require.Equal(t, exportGenesis.Tokens[0], newKeeper.GetUserTokensInfo(newCtx, exportGenesis.Tokens[0].Owner)[0])
	actualLockeedFees = []types.AccCoins{}
	newKeeper.IterateLockedFees(newCtx, func(acc sdk.AccAddress, coins sdk.DecCoins) bool {
		actualLockeedFees = append(actualLockeedFees,
			types.AccCoins{
				Acc:   acc,
				Coins: coins,
			})
		return false
	})
	require.Equal(t, exportGenesis.LockedFees, actualLockeedFees)

	newExportGenesis := ExportGenesis(newCtx, newKeeper)
	require.Equal(t, newExportGenesis.Params, newKeeper.GetParams(newCtx))
	require.Equal(t, newExportGenesis.Tokens, newKeeper.GetTokensInfo(newCtx))
	require.Equal(t, newExportGenesis.LockedAssets, newKeeper.GetAllLockedCoins(newCtx))
	actualLockeedFees = []types.AccCoins{}
	newKeeper.IterateLockedFees(newCtx, func(acc sdk.AccAddress, coins sdk.DecCoins) bool {
		actualLockeedFees = append(actualLockeedFees,
			types.AccCoins{
				Acc:   acc,
				Coins: coins,
			})
		return false
	})
	require.Equal(t, newExportGenesis.LockedFees, actualLockeedFees)
}

func TestIssueToken(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	//ctx, keeper, _, _ := CreateParam(t, false)
	genesisState := defaultGenesisState()
	gs := types.ModuleCdc.MustMarshalJSON(genesisState)

	coins := sdk.NewCoins(sdk.NewDecCoinFromDec(genesisState.Tokens[0].Symbol,
		genesisState.Tokens[0].OriginalTotalSupply))
	err := keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	acc, _ := CreateGenAccounts(1, nil)
	err2 := IssueOKT(ctx, keeper, gs, acc[0].ToBaseAccount())
	require.NoError(t, err2)
	token := keeper.GetTokenInfo(ctx, genesisState.Tokens[0].Symbol)
	expectToken := types.Token{
		Description:         genesisState.Tokens[0].Description,
		Symbol:              genesisState.Tokens[0].Symbol,
		OriginalSymbol:      genesisState.Tokens[0].OriginalSymbol,
		WholeName:           genesisState.Tokens[0].WholeName,
		OriginalTotalSupply: genesisState.Tokens[0].OriginalTotalSupply,
		//TotalSupply:         genesisState.Tokens[0].TotalSupply,
		Owner:    genesisState.Tokens[0].Owner,
		Mintable: genesisState.Tokens[0].Mintable,
	}
	require.EqualValues(t, expectToken, token)

	//coin with no owner
	coin := []types.Token{{
		Description:         "OK Group Global Utility Token",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		WholeName:           common.NativeToken,
		OriginalTotalSupply: sdk.NewDec(1000000000),
		//TotalSupply:         sdk.NewDec(1000000000),
		Owner:    nil,
		Mintable: true,
	}}

	coins = sdk.NewCoins(sdk.NewDecCoinFromDec(coin[0].Symbol, coin[0].OriginalTotalSupply))
	err = keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	genesisState.Tokens = coin
	gs = types.ModuleCdc.MustMarshalJSON(genesisState)
	err2 = IssueOKT(ctx, keeper, gs, acc[0].ToBaseAccount())
	require.NoError(t, err2)

	err2 = validateGenesis(genesisState)
	require.Error(t, err2)
}

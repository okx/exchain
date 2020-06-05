package token

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestCreateCertifiedToken(t *testing.T) {
	intQuantity := int64(6000)
	genAccs, testAccounts := CreateGenAccounts(3,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity)),
		})

	app, keeper, _ := getMockDexApp(t, 0)
	mock.SetGenesis(app.App, types.DecAccountArrToBaseAccountArr(genAccs))

	ctx := app.NewContext(true, abci.Header{})
	keeper.SetParams(ctx, types.DefaultParams())

	tokenBtc := types.CertifiedToken{
		Description: "btc btc btc btc",
		Symbol:      "btc",
		WholeName:   "Bitcoin",
		TotalSupply: "21000000",
		Owner:       testAccounts[0].baseAccount.Address,
		Mintable:    false,
	}
	tokenXmr := types.CertifiedToken{
		Description: "xmr xmr xmr xmr xmr",
		Symbol:      "xmr",
		WholeName:   "xmr",
		TotalSupply: "21000000",
		Owner:       testAccounts[1].baseAccount.Address,
		Mintable:    false,
	}
	tokenBnb := types.CertifiedToken{
		Description: "bnb bnb",
		Symbol:      "bnb",
		WholeName:   "bnb",
		TotalSupply: "21000000",
		Owner:       testAccounts[1].baseAccount.Address,
		Mintable:    false,
	}

	//
	// set proposal
	//
	proposalHandler := NewCertifiedTokenProposalHandler(&keeper)

	content := types.NewCertifiedTokenProposal(
		"no suffix token proposal",
		"no suffix token proposal",
		tokenBtc,
	)
	proposalBtc := sdkGovTypes.NewProposal(ctx, sdk.ZeroDec(), content, 1, time.Now(),
		time.Now().Add(time.Hour))
	err := proposalHandler(ctx, &proposalBtc)
	require.Nil(t, err)

	proposalBtc = sdkGovTypes.NewProposal(ctx, sdk.ZeroDec(), content, 2, time.Now(),
		time.Now().Add(time.Hour))
	err = proposalHandler(ctx, &proposalBtc)
	require.Nil(t, err)

	content = types.NewCertifiedTokenProposal(
		"no suffix token proposal",
		"no suffix token proposal",
		tokenXmr,
	)
	proposalXmr := sdkGovTypes.NewProposal(ctx, sdk.ZeroDec(), content, 3, time.Now(),
		time.Now().Add(time.Hour))
	err = proposalHandler(ctx, &proposalXmr)
	require.Nil(t, err)

	content = types.NewCertifiedTokenProposal(
		"no suffix token proposal",
		"no suffix token proposal",
		tokenBnb,
	)
	proposalBnb := sdkGovTypes.NewProposal(ctx, sdk.ZeroDec(), content, 4, time.Now(),
		time.Now().Add(time.Hour))
	err = proposalHandler(ctx, &proposalBnb)
	require.Nil(t, err)

	//
	// active
	//
	tokenHandler := NewTokenHandler(keeper, 0)
	activeMsg := types.NewMsgTokenActive(1, testAccounts[0].baseAccount.Address)
	r := tokenHandler(ctx, activeMsg)
	require.Equal(t, sdk.CodeType(0), r.Code)

	tokenInfo := app.tokenKeeper.GetTokenInfo(ctx, "btc")
	require.Equal(t, parseToken(tokenBtc), tokenInfo)

	// not owner
	activeMsg = types.NewMsgTokenActive(3, testAccounts[0].baseAccount.Address)
	r = tokenHandler(ctx, activeMsg)
	require.NotEqual(t, sdk.CodeType(0), r.Code)

	// token exits
	activeMsg = types.NewMsgTokenActive(1, testAccounts[0].baseAccount.Address)
	r = tokenHandler(ctx, activeMsg)
	require.NotEqual(t, sdk.CodeType(0), r.Code)

	activeMsg = types.NewMsgTokenActive(3, testAccounts[1].baseAccount.Address)
	r = tokenHandler(ctx, activeMsg)
	require.Equal(t, sdk.CodeType(0), r.Code)

	tokenInfo = keeper.GetTokenInfo(ctx, "xmr")
	require.Equal(t, parseToken(tokenXmr), tokenInfo)

	// not enough okts
	activeMsg = types.NewMsgTokenActive(4, testAccounts[2].baseAccount.Address)
	r = tokenHandler(ctx, activeMsg)
	require.NotEqual(t, sdk.CodeType(0), r.Code)

	feeIssue := keeper.GetParams(ctx).FeeIssue.Amount
	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokenBtc.Symbol, sdk.MustNewDecFromStr(tokenBtc.TotalSupply)),
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity).Sub(feeIssue)),
	}
	require.EqualValues(t, coins, app.AccountKeeper.GetAccount(ctx, testAccounts[0].addrKeys.Address).GetCoins())

	coins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(intQuantity).Sub(feeIssue)),
		sdk.NewDecCoinFromDec(tokenXmr.Symbol, sdk.MustNewDecFromStr(tokenXmr.TotalSupply)),
	}
	require.EqualValues(t, coins, app.AccountKeeper.GetAccount(ctx, testAccounts[1].addrKeys.Address).GetCoins())
}

func parseToken(token types.CertifiedToken) types.Token {
	return types.Token{
		Description:         token.Description,
		Symbol:              token.Symbol,
		OriginalSymbol:      token.Symbol,
		WholeName:           token.WholeName,
		OriginalTotalSupply: sdk.MustNewDecFromStr(token.TotalSupply),
		TotalSupply:         sdk.MustNewDecFromStr(token.TotalSupply),
		Owner:               token.Owner,
		Mintable:            token.Mintable,
	}
}

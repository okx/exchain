package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okchain/x/poolswap/types"
	token "github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestKeeper_IsTokenExistTable(t *testing.T) {
	mapp, _ := GetTestInput(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	tests := []struct {
		testCase         string
		tokennames       []string
		tokentypes       []int
		tokenname        string
		exceptResultCode sdk.CodeType
	}{
		{"token is not exist", []string{"toa", "tob"}, []int{1, 1}, "nota", sdk.CodeInternal},
		{"token is not exist", nil, nil, "nota", sdk.CodeInternal},
		{"token is exist", []string{"boa", "bob"}, []int{1, 1}, "boa", sdk.CodeOK},
		{"token is pooltoken", []string{"tkoa", "tkob"}, []int{1, 2}, "tkob", sdk.CodeInvalidCoins},
	}

	for _, testCase := range tests {
		fmt.Println(testCase.testCase)
		genToken(mapp, ctx, testCase.tokennames, testCase.tokentypes)
		result := keeper.IsTokenExist(ctx, testCase.tokenname)
		if nil != result {
			require.Equal(t, testCase.exceptResultCode, result.(sdk.Error).Code())
		}
	}

}

func genToken(mapp *TestInput, ctx sdk.Context, tokennames []string, tokentypes []int) {
	for i, t := range tokennames {
		tok := token.Token{
			Description:         t,
			Symbol:              t,
			OriginalSymbol:      t,
			WholeName:           t,
			OriginalTotalSupply: sdk.NewDec(0),
			TotalSupply:         sdk.NewDec(0),
			Owner:               supply.NewModuleAddress(types.ModuleName),
			Mintable:            true,
			Type:                tokentypes[i],
		}
		mapp.tokenKeeper.NewToken(ctx, tok)
	}
}

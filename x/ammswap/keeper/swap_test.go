package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/supply"
	"github.com/okex/exchain/x/ammswap/types"
	token "github.com/okex/exchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
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
		exceptResultCode uint32
	}{
		{"token is not exist", []string{"toa", "tob"}, []int{1, 1}, "nota", sdk.CodeInternal},
		{"token is not exist", nil, nil, "nota", sdk.CodeInternal},
		{"token is exist", []string{"boa", "bob"}, []int{1, 1}, "boa", sdk.CodeOK},
		{"token is pool token", []string{"tkoa", "tkob"}, []int{1, 2}, "tkob", sdk.CodeInvalidCoins},
	}

	for _, testCase := range tests {
		fmt.Println(testCase.testCase)
		genToken(mapp, ctx, testCase.tokennames, testCase.tokentypes)
		result := keeper.IsTokenExist(ctx, testCase.tokenname)
		if nil != result {
			if testCase.exceptResultCode == 0 {
				require.Nil(t, result)
			}else {
				require.NotNil(t, result)
			}
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
			Owner:               supply.NewModuleAddress(types.ModuleName),
			Mintable:            true,
			Type:                tokentypes[i],
		}
		mapp.tokenKeeper.NewToken(ctx, tok)
	}
}

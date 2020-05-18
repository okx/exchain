package v0_9

import (
	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
	v08token "github.com/okex/okchain/x/token/legacy/v0_8"
	"github.com/okex/okchain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Migrate(oldGenState v08token.GenesisState, oldgovParams v08gov.GovParams) GenesisState {

	params := types.DefaultParams()

	tokens := make([]types.Token, len(oldGenState.Info))
	for k, token := range oldGenState.Info {
		tokens[k] = types.Token{
			Description:         token.Desc,
			Symbol:              token.Symbol,
			OriginalSymbol:      token.OriginalSymbol,
			WholeName:           token.WholeName,
			OriginalTotalSupply: sdk.NewDec(token.TotalSupply),
			TotalSupply:         sdk.NewDec(token.TotalSupply),
			Owner:               token.Owner,
			Mintable:            token.Mintable,
		}
	}
	return GenesisState{
		Params:       params,
		Tokens:       tokens,
		LockedAssets: oldGenState.LockCoins,
	}
}

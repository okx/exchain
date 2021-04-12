package v0_11

import "github.com/okex/exchain/x/token/legacy/v0_10"

func Migrate(oldGenState v0_10.GenesisState) GenesisState {
	tokens := make([]Token, len(oldGenState.Tokens))
	for i, token := range oldGenState.Tokens {
		tokens[i] = Token{
			Description:         token.Description,
			Symbol:              token.Symbol,
			OriginalSymbol:      token.OriginalSymbol,
			WholeName:           token.WholeName,
			OriginalTotalSupply: token.OriginalTotalSupply,
			Owner:               token.Owner,
			Mintable:            token.Mintable,
		}
	}

	return GenesisState{
		Params:       oldGenState.Params,
		Tokens:       tokens,
		LockedAssets: oldGenState.LockedAssets,
		LockedFees:   oldGenState.LockedFees,
	}
}

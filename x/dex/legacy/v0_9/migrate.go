// nolint
package v0_9

import (
	"github.com/okex/okchain/x/dex/types"
	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
	v08token "github.com/okex/okchain/x/token/legacy/v0_8"
)

func Migrate(oldGenState v08token.GenesisState, _ v08gov.GovParams) GenesisState {
	params := types.DefaultParams()
	tokenPairs := make([]*types.TokenPair, len(oldGenState.TokenPairs))
	for i, pair := range oldGenState.TokenPairs {
		tokenPairs[i] = &types.TokenPair{
			BaseAssetSymbol:  pair.BaseAssetSymbol,
			QuoteAssetSymbol: pair.QuoteAssetSymbol,
			InitPrice:        pair.InitPrice,
			MaxPriceDigit:    pair.MaxPriceDigit,
			MaxQuantityDigit: pair.MaxQuantityDigit,
			MinQuantity:      pair.MinQuantity,
			ID:               pair.ID,
			Delisting:        false,
			Owner:            nil,
			Deposits:         types.DefaultTokenPairDeposit,
			BlockHeight:      1,
		}
	}

	return GenesisState{
		Params:        *params,
		TokenPairs:    tokenPairs,
		WithdrawInfos: nil,
	}
}

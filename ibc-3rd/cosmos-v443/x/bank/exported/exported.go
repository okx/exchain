package exported

import (
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
)

// GenesisBalance defines a genesis balance interface that allows for account
// address and balance retrieval.
type GenesisBalance interface {
	GetAddress() sdk.AccAddress
	GetCoins() sdk.Coins
}

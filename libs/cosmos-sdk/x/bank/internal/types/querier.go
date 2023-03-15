package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address sdk.AccAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr sdk.AccAddress) QueryBalanceParams {
	return QueryBalanceParams{Address: addr}
}

type WrappedBalances struct {
	Balances sdk.Coins `json:"balances,omitempty"`
}

func NewWrappedBalances(coins sdk.Coins) WrappedBalances {
	return WrappedBalances{Balances: coins}
}

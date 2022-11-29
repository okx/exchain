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

type QueryBalanceWithDenomParams struct {
	Address sdk.AccAddress
	Denom   string
}

func NewQueryBalanceWithDenomParams(addr sdk.AccAddress, denom string) QueryBalanceWithDenomParams {
	return QueryBalanceWithDenomParams{Address: addr, Denom: denom}
}

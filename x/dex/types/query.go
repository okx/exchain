package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

const (
	// DefaultPage defines default number of page
	DefaultPage = 1
	// DefaultPerPage defines default number per page
	DefaultPerPage = 50
)

// QueryDexInfoParams defines query params of dex info
type QueryDexInfoParams struct {
	Owner   string
	Page    int
	PerPage int
}

// NewQueryDexInfoParams creates query params of dex info
func NewQueryDexInfoParams(owner string, page, perPage int) QueryDexInfoParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryDexInfoParams{
		Owner:   owner,
		Page:    page,
		PerPage: perPage,
	}
}

type QueryDexOperatorParams struct {
	Addr sdk.AccAddress
}

// creates a new instance of QueryDexOperatorParams
func NewQueryDexOperatorParams(addr sdk.AccAddress) QueryDexOperatorParams {
	return QueryDexOperatorParams{
		Addr: addr,
	}
}

// nolint
type QueryDepositParams struct {
	Address    string
	BaseAsset  string
	QuoteAsset string
	Page       int
	PerPage    int
}

// NewQueryDepositParams creates a new instance of QueryDepositParams
func NewQueryDepositParams(address, baseAsset, quoteAsset string, page, perPage int) QueryDepositParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryDepositParams{
		Address:    address,
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
		Page:       page,
		PerPage:    perPage,
	}
}

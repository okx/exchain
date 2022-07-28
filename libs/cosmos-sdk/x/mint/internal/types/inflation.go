package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

// WrappedInflation is used to wrap the Inflation, thus making the rest API response compatible with cosmos-sdk
type WrappedInflation struct {
	Inflation sdk.Dec `json:"inflation" yaml:"inflation"`
}

func NewWrappedInflation(inflation sdk.Dec) WrappedInflation {
	return WrappedInflation{
		Inflation: inflation,
	}
}

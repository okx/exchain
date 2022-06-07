package ibc_tx

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type DenomAdapterMsg interface {
	sdk.Msg
	DenomOpr
}

type DenomOpr interface {
	RulesFilter() (sdk.Msg, error)
}

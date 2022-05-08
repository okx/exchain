package ibc_tx

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

type DenomAdapterIbcTransferMsg interface {
	sdk.Msg
	DenomOpr
}

type DenomOpr interface {
	RulesFilter() (sdk.Msg, error)
}

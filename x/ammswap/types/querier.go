package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryBuyAmountParams struct {
	SoldToken    sdk.DecCoin
	TokenToBuy string
}

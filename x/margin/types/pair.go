package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TradePair struct {
	Owner                  sdk.AccAddress `json:"address"`
	Name                   string         `json:"name"`
	Deposit                sdk.DecCoin    `json:"deposit"`
	MaxLeverage            sdk.Dec        `json:"max-leverage"`
	BorrowRate             sdk.Dec        `json:"borrow-rate"`
	MaintenanceMarginRatio sdk.Dec        `json:"maintenance-margin-ratio"`
	BlockHeight            int64          `json:"block_height"`
}

func (tp TradePair) BaseSymbol() string {
	return strings.Split(tp.Name, "_")[0]
}

func (tp TradePair) QuoteSymbol() string {
	return strings.Split(tp.Name, "_")[1]
}

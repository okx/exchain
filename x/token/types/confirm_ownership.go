package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultConfirmPeriod defines default confirm  period
const DefaultConfirmPeriod = 24 * time.Hour

type ConfirmOwnership struct {
	Symbol  string         `json:"symbol"`
	Address sdk.AccAddress `json:"address"`
	Expire  time.Time      `json:expire`
}

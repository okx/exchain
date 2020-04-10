package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DefaultWithdrawPeriod = time.Hour * 24 * 3

type WithdrawInfo struct {
	Owner        sdk.AccAddress `json:"owner"`
	Deposits     sdk.DecCoin    `json:"deposits"`
	CompleteTime time.Time      `json:"complete_time"`
}

func (w WithdrawInfo) Equal(other WithdrawInfo) bool {
	return w.Owner.Equals(other.Owner) && w.Deposits.IsEqual(other.Deposits) && w.CompleteTime.Equal(other.CompleteTime)
}

type WithdrawInfos []WithdrawInfo

func (ws WithdrawInfos) Equal(other WithdrawInfos) bool {
	if len(ws) != len(other) {
		return false
	}

	for i := 0; i < len(ws); i++ {
		if !ws[i].Equal(other[i]) {
			return false
		}
	}
	return true
}

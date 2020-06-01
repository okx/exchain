package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DexWithdrawInfo represents infos for withdrawing
type DexWithdrawInfo struct {
	Owner        sdk.AccAddress `json:"owner"`
	Deposits     sdk.DecCoin    `json:"deposits"`
	CompleteTime time.Time      `json:"complete_time"`
}

// Equal returns boolean for whether two WithdrawInfo are Equal
func (w DexWithdrawInfo) Equal(other DexWithdrawInfo) bool {
	return w.Owner.Equals(other.Owner) && w.Deposits.IsEqual(other.Deposits) && w.CompleteTime.Equal(other.CompleteTime)
}

// DexWithdrawInfos defines list of WithdrawInfo
type DexWithdrawInfos []DexWithdrawInfo

// Equal returns boolean for whether two WithdrawInfos are Equal
func (ws DexWithdrawInfos) Equal(other DexWithdrawInfos) bool {
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

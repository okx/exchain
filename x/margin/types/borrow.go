package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BorrowInfo struct {
	Address      sdk.AccAddress `json:"address"`
	Product      string         `json:"product"`
	BorrowAmount sdk.DecCoins   `json:"amount"`
	BlockHeight  int64          `json:"block_height"`
	Rate         sdk.Dec        `json:"rate"`
	Leverage     sdk.Dec        `json:"leverage"`
}

// ShouldRepayEarlier returns true if the borrowInfo should be repaid earlier
// 1. compare rate
// 2. compare block height
func (borrowInfo *BorrowInfo) ShouldRepayEarlier(other *BorrowInfo) bool {
	if borrowInfo.Rate.GT(other.Rate) {
		return true
	}
	return borrowInfo.BlockHeight < other.BlockHeight
}

type BorrowInfoList []*BorrowInfo

// Len Implements Sort
func (list BorrowInfoList) Len() int { return len(list) }

// Less Implements Sort
func (list BorrowInfoList) Less(i, j int) bool { return list[i].ShouldRepayEarlier(list[j]) }

// Swap Implements Sort
func (list BorrowInfoList) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

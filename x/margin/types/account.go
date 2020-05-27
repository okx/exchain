package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RepayInterestAndPrincipal = 1
	RepayPartialPrincipal     = 2
	RepayInterest             = 3
)

type MarginAccountAssets struct {
	MarginProductAssets
}

type MarginProductAssets []AccountAssetOnProduct

type AccountAssetOnProduct struct {
	Product   string       `json:"product"`
	Available sdk.DecCoins `json:"available"`
	Locked    sdk.DecCoins `json:"locked"`
	Borrowed  sdk.DecCoins `json:"borrowed"`
	Deposits  sdk.DecCoins `json:"deposits"`
	//Interest  sdk.DecCoins `json:"interest"`
}

type BorrowInfo struct {
	BorrowAmount  sdk.DecCoin `json:"amount"`
	BorrowDeposit sdk.DecCoin `json:"deposit"`
	BlockHeight   int64       `json:"block_height"`
	Rate          sdk.Dec     `json:"rate"`
	Interest      sdk.DecCoin `json:"interest"`
}

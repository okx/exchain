package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	Deposit   sdk.DecCoins `json:"deposit"`
	//Interest  sdk.DecCoins `json:"interest"`
}

type BorrowInfo struct {
	Token       sdk.DecCoin `json:"amount"`
	BlockHeight int64       `json:"block_height"`
	Rate        sdk.Dec     `json:"rate"`
	Interest    sdk.DecCoin `json:"interest"`
}

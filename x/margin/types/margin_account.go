package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MarginAccountAssets struct {
	MarginProductAssets
}

type MarginProductAssets []MarginAssetOnProduct

type MarginAssetOnProduct struct {
	Product   string       `json:"product"`
	Available sdk.DecCoins `json:"available"`
	Locked    sdk.DecCoins `json:"locked"`
	Borrowed  sdk.DecCoins `json:"borrowed"`
	//Interest  sdk.DecCoins `json:"interest"`
}
type TradePair struct {
	Name        string `json:"name"`
	Leverage    int    `json:"leverage"`
	BlockHeight int64  `json:"block_height"`
}

type BorrowInfo struct {
	Amount      sdk.DecCoin `json:"amount"`
	BlockHeight int64       `json:"block_height"`
	Rate        sdk.Dec     `json:"rate"`
}

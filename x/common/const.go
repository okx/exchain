package common

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// const
const (
	//NativeToken = sdk.RawDefaultBondDenom
	TestToken = "xxb"

	blackHoleHex = "0000000000000000000000000000000000000000"
)

func NativeToken() string {
	return sdk.DefaultBondDenom()
}

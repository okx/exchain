package types

import (
	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"
	"github.com/okx/exchain/x/common"
)

// nolint
const (
	FeeTypeOrderNew     = "new"
	FeeTypeOrderCancel  = "cancel"
	FeeTypeOrderExpire  = "expire"
	FeeTypeOrderDeal    = "deal"
	FeeTypeOrderReceive = "receive"
	TestTokenPair       = common.TestToken + "_" + sdk.DefaultBondDenom
	BuyOrder            = "BUY"
	SellOrder           = "SELL"
)

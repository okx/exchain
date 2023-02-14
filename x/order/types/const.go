package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
)

// nolint
const (
	FeeTypeOrderNew     = "new"
	FeeTypeOrderCancel  = "cancel"
	FeeTypeOrderExpire  = "expire"
	FeeTypeOrderDeal    = "deal"
	FeeTypeOrderReceive = "receive"
	TestTokenPair       = common.TestToken + "_" + sdk.RawDefaultBondDenom
	BuyOrder            = "BUY"
	SellOrder           = "SELL"
)

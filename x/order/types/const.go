package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
)

// nolint
const (
	FeeTypeOrderNew         = "new"
	FeeTypeOrderCancel      = "cancel"
	FeeTypeOrderExpire      = "expire"
	FeeTypeOrderDeal        = "deal"
	FeeTypeOrderReceive     = "receive"
	TestTokenPair           = common.TestToken + "_" + sdk.DefaultBondDenom
	BuyOrder                = "BUY"
	SellOrder               = "SELL"
	DefaultNewOrderFeeRatio = "1"
)

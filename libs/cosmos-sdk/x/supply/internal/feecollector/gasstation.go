package feecollector

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"sync/atomic"
)

var feeCollectorAsReceiver sdk.Gas
var feeCollectorAsSender sdk.Gas

func SetFeeCollectorAsReceiver(accountGas sdk.Gas) {
	atomic.StoreUint64(&feeCollectorAsReceiver, accountGas)
}

func GetFeeCollectorAsReceiver() sdk.Gas {
	return atomic.LoadUint64(&feeCollectorAsReceiver)
}

func SetFeeCollectorAsSender(accountGas sdk.Gas) {
	atomic.StoreUint64(&feeCollectorAsSender, accountGas)
}

func GetFeeCollectorAsSender() sdk.Gas {
	return atomic.LoadUint64(&feeCollectorAsSender)
}

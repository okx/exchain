package backend

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	rpctypes "github.com/okex/exchain/app/rpc/types"
)

const (
	DefaultFeeHistoryCap int64 = 100
)

func (b *EthermintBackend) FeeHistory(
	userBlockCount rpc.DecimalOrHex, // number blocks to fetch, maximum is 100
	lastBlock rpc.BlockNumber, // the block to start search , to oldest
	rewardPercentiles []float64, // percentiles to fetch reward
) (*rpctypes.FeeHistoryResult, error) {

	return nil, fmt.Errorf("unsupported rpc function: eth_FeeHistory")
}

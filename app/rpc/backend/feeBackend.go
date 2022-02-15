package backend

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

const (
	DefaultFeeHistoryCap int64 = 100
)

func (b *EthermintBackend) FeeHistory(
	userBlockCount rpc.DecimalOrHex, // number blocks to fetch, maximum is 100
	lastBlock rpc.BlockNumber, // the block to start search , to oldest
	rewardPercentiles []float64, // percentiles to fetch reward
) (*rpctypes.FeeHistoryResult, error) {

	latestHeight, err := b.BlockNumber()
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get Height.")
	}
	if !tmtypes.IsLondon(int64(latestHeight)) {
		return nil, fmt.Errorf("unsupported rpc function: eth_FeeHistory")
	}
	blockEnd := int64(lastBlock)
	if blockEnd <= 0 {
		blockEnd = int64(latestHeight)
	}
	userBlockCountInt := int64(userBlockCount)
	if userBlockCountInt > DefaultFeeHistoryCap {
		return nil, fmt.Errorf("FeeHistory user block count %d higher than %d", userBlockCountInt, DefaultFeeHistoryCap)
	}
	blockStart := blockEnd - userBlockCountInt
	if blockStart < 0 {
		blockStart = 0
	}
	oldestBlock := (*hexutil.Big)(big.NewInt(blockStart))
	return &rpctypes.FeeHistoryResult{
		OldestBlock: oldestBlock,
	}, nil
}

package evmkeeper

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp/evmtx"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Keeper interface {
	SaveEvmTxAndSuccessReceipt(evmTx sdk.Tx, txIndexInBlock uint64, resultData evmtx.ResultData, gasUsed uint64) error
	SaveEvmTxAndFailedReceipt(evmTx sdk.Tx, txIndexInBlock uint64, txHash ethcmn.Hash, gasUsed uint64) error
	GetTxIndexInBlock() uint64
}

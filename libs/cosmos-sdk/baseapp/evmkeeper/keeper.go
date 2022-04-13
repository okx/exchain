package evmkeeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp/evmtx"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Keeper interface {
	SaveEvmTxAndSuccessReceipt(evmTx sdk.Tx, txIndexInBlock uint64, resultData evmtx.ResultData, gasUsed uint64) error
	SaveEvmTxAndFailedReceipt(evmTx sdk.Tx, txIndexInBlock uint64, resultData evmtx.ResultData, gasUsed uint64) error
	GetTxIndexInBlock() uint64
}

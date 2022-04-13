package evmkeeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp/evmtx"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Keeper interface {
	SaveTxAndSuccessReceipt(evmTx sdk.Tx, txIndexInBlock uint64, resultData evmtx.ResultData, gasUsed uint64)
	SaveTxAndFailedReceipt(evmTx sdk.Tx, txIndexInBlock uint64, resultData evmtx.ResultData, gasUsed uint64)
	GetTxIndexInBlock() uint64
}

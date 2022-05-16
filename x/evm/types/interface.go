package types

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to customize evm transaction processing.

// EvmHooks event hooks for evm tx processing
type EvmHooks interface {
	// PostTxProcessing Must be called after tx is processed successfully, if return an error, the whole transaction is reverted.
	PostTxProcessing(ctx sdk.Context, from common.Address, to *common.Address, receipt *ethtypes.Receipt) error
}

// EvmLogHandler defines the interface for evm log handler
type EvmLogHandler interface {
	// EventID Return the id of the log signature it handles
	EventID() common.Hash
	// Handle Process the log
	Handle(ctx sdk.Context, contract common.Address, data []byte) error
}

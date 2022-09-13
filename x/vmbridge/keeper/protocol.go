package keeper

import (
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/vmbridge/types"
)

type SendToWasmEventHandler struct {
	Keeper
}

func NewSendToWasmEventHandler(k Keeper) *SendToWasmEventHandler {
	return &SendToWasmEventHandler{k}
}

// EventID Return the id of the log signature it handles
func (h SendToWasmEventHandler) EventID() common.Hash {
	return types.SendToWasmEvent.ID
}

// Handle Process the log
func (h SendToWasmEventHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	logger := h.Keeper.Logger()
	unpacked, err := types.SendToWasmEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		logger.Error("log signature matches but failed to decode", "error", err)
		return nil
	}

}

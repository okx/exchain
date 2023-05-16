package vmbridge

import (
	"github.com/okex/exchain/x/vmbridge/keeper"
	"github.com/okex/exchain/x/vmbridge/types"
)

var (
	RegisterMsgServer         = types.RegisterMsgServer
	NewMsgServerImpl          = keeper.NewMsgServerImpl
	NewSendToWasmEventHandler = keeper.NewSendToWasmEventHandler
	NewCallToWasmEventHandler = keeper.NewCallToWasmEventHandler
	RegisterSendToEvmEncoder  = keeper.RegisterSendToEvmEncoder
	NewKeeper                 = keeper.NewKeeper
	RegisterInterface         = types.RegisterInterface
	PrecompileHooks           = keeper.PrecompileHooks
)

type (
	MsgSendToEvm = types.MsgSendToEvm
	Keeper       = keeper.Keeper
)

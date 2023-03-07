package vmbridge

import (
	"github.com/okx/okbchain/x/vmbridge/keeper"
	"github.com/okx/okbchain/x/vmbridge/types"
)

var (
	RegisterMsgServer         = types.RegisterMsgServer
	NewMsgServerImpl          = keeper.NewMsgServerImpl
	NewSendToWasmEventHandler = keeper.NewSendToWasmEventHandler
	RegisterSendToEvmEncoder  = keeper.RegisterSendToEvmEncoder
	NewKeeper                 = keeper.NewKeeper
	RegisterInterface         = types.RegisterInterface
)

type (
	MsgSendToEvm = types.MsgSendToEvm
	Keeper       = keeper.Keeper
)

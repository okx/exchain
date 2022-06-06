package simulator

import (
	cliContext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Simulator interface {
	Simulate(msg sdk.Msg) (*sdk.Result, error)
	Context() *sdk.Context
}

var SimulateCliCtx cliContext.CLIContext
var NewWasmSimulator func(cliCtx cliContext.CLIContext) Simulator

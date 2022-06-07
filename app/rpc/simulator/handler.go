package simulator

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Simulator interface {
	Simulate(msg sdk.Msg) (*sdk.Result, error)
	Context() *sdk.Context
}

var NewWasmSimulator func() Simulator

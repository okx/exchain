package dex

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)


type MockInvariantRegistry struct{}

func (ir MockInvariantRegistry) RegisterRoute(moduleName, route string, invar sdk.Invariant) {}

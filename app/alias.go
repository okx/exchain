package app

import (
	"github.com/okex/okchain/app/protocol"
)

var (
	// functions aliases
	MakeCodec = protocol.MakeCodec
	// variable aliases
	ModuleBasics    = protocol.ModuleBasics
	DefaultCLIHome  = protocol.DefaultCLIHome
	DefaultNodeHome = protocol.DefaultNodeHome
)

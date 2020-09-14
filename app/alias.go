package app

import (
	"github.com/okex/okexchain/app/protocol"
)

var (
	// MakeCodec is the function alias for codec maker
	MakeCodec = protocol.MakeCodec
	// ModuleBasics is the variable alias for NewBasicManager
	ModuleBasics = protocol.ModuleBasics
	// DefaultCLIHome is the directory for okexchaincli
	DefaultCLIHome = protocol.DefaultCLIHome
	// DefaultNodeHome is the directory for okexchaind
	DefaultNodeHome = protocol.DefaultNodeHome
)

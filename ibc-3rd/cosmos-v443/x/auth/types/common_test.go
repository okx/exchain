package types_test

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/simapp"
)

var (
	app                   = simapp.Setup(false)
	ecdc                  = simapp.MakeTestEncodingConfig()
	appCodec, legacyAmino = ecdc.Marshaler, ecdc.Amino
)

package ibc_tx

import (
	"fmt"
	signing2 "github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	signing "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibcsigning"
)

// DefaultSignModes are the default sign modes enabled for protobuf transactions.
var DefaultSignModes = []signing2.SignMode{
	signing2.SignMode_SIGN_MODE_DIRECT,
	//signing2.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
}

// makeSignModeHandler returns the default protobuf SignModeHandler supporting
// SIGN_MODE_DIRECT and SIGN_MODE_LEGACY_AMINO_JSON.
func makeSignModeHandler(modes []signing2.SignMode) signing.SignModeHandler {
	if len(modes) < 1 {
		panic(fmt.Errorf("no sign modes enabled"))
	}

	handlers := make([]signing.SignModeHandler, len(modes))

	for i, mode := range modes {
		switch mode {
		case signing2.SignMode_SIGN_MODE_DIRECT:
			handlers[i] = signModeDirectHandler{}
		//case signing2.SignMode_SIGN_MODE_LEGACY_AMINO_JSON:
		//	handlers[i] = signModeLegacyAminoJSONHandler{}
		default:
			panic(fmt.Errorf("unsupported sign mode %+v", mode))
		}
	}

	return signing.NewSignModeHandlerMap(
		modes[0],
		handlers,
	)
}

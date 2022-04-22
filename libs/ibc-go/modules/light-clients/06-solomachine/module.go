package solomachine

import "github.com/okex/exchain/libs/ibc-go/modules/light-clients/06-solomachine/types"

// Name returns the solo machine client name.
func Name() string {
	return types.SubModuleName
}

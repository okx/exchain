package types

import (
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)


//use MPT storage model to replace IAVL storage model
func HigherThanVenus(height int64) bool {
	return tmtypes.HigherThanVenus(height)
}

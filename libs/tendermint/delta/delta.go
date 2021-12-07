package delta

import tmtypes "github.com/okex/exchain/libs/tendermint/types"

type DeltaBroker interface {
	SetBlock(block *tmtypes.Block) error
	SetDelta(deltas *tmtypes.Deltas) error
	GetBlock(height int64) (*tmtypes.Block, error)
	GetDeltas(height int64) (*tmtypes.Deltas, error)
}

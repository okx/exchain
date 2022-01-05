package store

import (
	"github.com/pkg/errors"

	"github.com/okex/exchain/libs/tendermint/types"
)

// LoadToBlockPart set the Part at the given index
// from the block at the given height.
// If no part is found for the given height and index, it returns false.
func (bs *BlockStore) LoadToBlockPart(height int64, index int, part *types.Part) (found bool) {
	bz, err := bs.db.Get(calcBlockPartKey(height, index))
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		found = false
		return
	}
	*part = types.Part{}
	err = part.UnmarshalFromAmino(bz)
	if err != nil {
		*part = types.Part{}
		err = cdc.UnmarshalBinaryBare(bz, part)
		if err != nil {
			panic(errors.Wrap(err, "Error reading block part"))
		}
	}
	found = true
	return
}

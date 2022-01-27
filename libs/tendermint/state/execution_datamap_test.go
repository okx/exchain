package state

import (
	"fmt"
	"testing"

	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
)

func TestDataMap(t *testing.T) {
	mp := newDataMap()
	mp.insert(1, &types.Deltas{}, &DeltaInfo{}, 1)
	mp.insert(2, &types.Deltas{}, &DeltaInfo{}, 2)
	mp.insert(3, &types.Deltas{}, &DeltaInfo{}, 3)
	mp.insert(4, &types.Deltas{}, &DeltaInfo{}, 4)
	mp.insert(5, &types.Deltas{}, &DeltaInfo{}, 5)
	mp.insert(6, &types.Deltas{}, &DeltaInfo{}, 6)
	mp.insert(10, &types.Deltas{}, &DeltaInfo{}, 7)

	a, b := mp.remove(4)
	assert.EqualValues(t, a, 4)
	assert.EqualValues(t, b, 3)
	fmt.Printf("%d, %d\n", a, b)
}

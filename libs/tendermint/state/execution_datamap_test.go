package state

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataMap(t *testing.T) {
	mp := newDataMap()
	mp.insert(1, &types.Deltas{}, 1)
	mp.insert(2, &types.Deltas{}, 2)
	mp.insert(3, &types.Deltas{}, 3)
	mp.insert(4, &types.Deltas{}, 4)
	mp.insert(5, &types.Deltas{}, 5)
	mp.insert(6, &types.Deltas{}, 6)
	mp.insert(10, &types.Deltas{},7)

	a, b := mp.remove(4)
	assert.EqualValues(t, a, 4)
	assert.EqualValues(t, b, 3)
	fmt.Printf("%d, %d\n", a, b)
}
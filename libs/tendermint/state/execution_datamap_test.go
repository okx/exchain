package state

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataMap(t *testing.T) {
	mp := newDataMap()
	mp.insert(1, &types.Deltas{})
	mp.insert(2, &types.Deltas{})
	mp.insert(3, &types.Deltas{})
	mp.insert(4, &types.Deltas{})
	mp.insert(5, &types.Deltas{})
	mp.insert(6, &types.Deltas{})
	mp.insert(10, &types.Deltas{})

	a, b := mp.remove(4)
	assert.EqualValues(t, a, 4)
	assert.EqualValues(t, b, 3)
	fmt.Printf("%d, %d\n", a, b)
}
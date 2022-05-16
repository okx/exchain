package state

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataMap(t *testing.T) {
	mp := newDataMap()
	mp.insert(1, &DeltaInfo{}, 1)
	mp.insert(2, &DeltaInfo{}, 2)
	mp.insert(3, &DeltaInfo{}, 3)
	mp.insert(4, &DeltaInfo{}, 4)
	mp.insert(5, &DeltaInfo{}, 5)
	mp.insert(6, &DeltaInfo{}, 6)
	mp.insert(10, &DeltaInfo{}, 7)

	a, b := mp.remove(4)
	assert.EqualValues(t, a, 4)
	assert.EqualValues(t, b, 3)
	fmt.Printf("%d, %d\n", a, b)
}

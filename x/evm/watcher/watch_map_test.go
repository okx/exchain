package watcher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataMap(t *testing.T) {
	mp := newDataMap()
	mp.insert(1, &WatchData{}, 1)
	mp.insert(2, &WatchData{}, 2)
	mp.insert(3, &WatchData{}, 3)
	mp.insert(4, &WatchData{}, 4)
	mp.insert(5, &WatchData{}, 5)
	mp.insert(6, &WatchData{}, 6)
	mp.insert(10, &WatchData{},7)

	a, b := mp.remove(4)
	assert.EqualValues(t, a, 4)
	assert.EqualValues(t, b, 3)
	fmt.Printf("%d, %d\n", a, b)
}
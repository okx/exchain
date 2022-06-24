package types

import (
	tmrand "github.com/okex/exchain/libs/tendermint/libs/rand"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

const (
	testPartSize = 65536 // 64KB ...  4096 // 4KB
)

func TestBasicPartSetCache(t *testing.T) {
	// Construct random data of size partSize * 100
	data := tmrand.Bytes(testPartSize * 100)
	partSet := types.NewPartSetFromData(data, testPartSize)

	assert.NotEmpty(t, partSet.Hash())
	assert.Equal(t, 100, partSet.Total())
	assert.Equal(t, 100, partSet.BitArray().Size())
	assert.True(t, partSet.HashesTo(partSet.Hash()))
	assert.True(t, partSet.IsComplete())
	assert.Equal(t, 100, partSet.Count())

	// Test adding parts to a new partSet.
	partSet2 := types.NewPartSetFromHeader(partSet.Header())

	assert.True(t, partSet2.HasHeader(partSet.Header()))

	//
	var height int64 = 100
	hbc := NewBPCache(height)
	for i := 0; i < 100; i++ {
		hbc.AddBlockPart(height, partSet.GetPart(i))
	}

	hit := 0
	for _, part := range hbc.Cache() {
		if added, _ := partSet2.AddPart(part); added {
			hit++
		}
	}
	assert.Equal(t, 100, hit)

	assert.Equal(t, partSet.Hash(), partSet2.Hash())
	assert.Equal(t, 100, partSet2.Total())
	assert.True(t, partSet2.IsComplete())

	// Reconstruct data, assert that they are equal.
	data2Reader := partSet2.GetReader()
	data2, err := ioutil.ReadAll(data2Reader)
	require.NoError(t, err)

	assert.Equal(t, data, data2)
}

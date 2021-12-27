package types

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDelta(t *testing.T) {
	d := &Deltas{
		Payload: DeltaPayload{
			[]byte("abci"),
			[]byte("detal"),
			[]byte("watch"),
		},
		Height: 1,
		Version: 1,
		CompressType: 2,
	}

	marshaled, err := d.Marshal()
	require.NoError(t, err)

	unmarshaled := &Deltas{}
	err = unmarshaled.Unmarshal(marshaled)
	require.NoError(t, err)

	assert.True(t,	bytes.Compare(unmarshaled.ABCIRsp(), d.ABCIRsp()) == 0)
	assert.True(t,	bytes.Compare(unmarshaled.DeltasBytes(), d.DeltasBytes()) == 0)
	assert.True(t,	bytes.Compare(unmarshaled.WatchBytes(), d.WatchBytes()) == 0)
	assert.True(t, 	unmarshaled.Height == d.Height)
	assert.True(t, 	unmarshaled.Version == d.Version)
	assert.True(t, 	unmarshaled.CompressType == d.CompressType)
}


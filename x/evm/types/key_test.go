package types

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBloomKey(t *testing.T) {
	expectedHeight, expectedBloomKey := int64(1), []byte{0, 0, 0, 0, 0, 0, 0, 1}
	require.True(t, bytes.Equal(BloomKey(expectedHeight), expectedBloomKey))
}

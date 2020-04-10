package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQueryVersion(t *testing.T) {
	require.Equal(t, uint64(1), NewQueryVersion(1).Ver)
}

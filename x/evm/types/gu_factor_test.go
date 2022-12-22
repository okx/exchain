package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalGuFactor(t *testing.T) {
	str := "{\"gu_factor\":\"2\"}"
	_, err := UnmarshalGuFactor(str)
	require.NoError(t, err)

}

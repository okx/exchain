package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlackHoleAddress(t *testing.T) {
	addr := BlackHoleAddress()
	a := addr.String()
	fmt.Println(a)
	require.Equal(t, addr.String(), "okexchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqupa6dx")
}

package dydx

import (
	"math/big"
	"testing"
)

func TestBigFloat(t *testing.T) {
	num := big.NewFloat(1.001)
	t.Log(num.Text('f', 0))
}

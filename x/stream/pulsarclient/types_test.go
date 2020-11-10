package pulsarclient

import (
	"testing"
)

func TestNewPulsar(t *testing.T) {
	pd := NewPulsarData()
	pd.BlockHeight()
	pd.DataType()
}

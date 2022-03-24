package types

import (
	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
)

// DefaultTrustLevel is the tendermint light client default trust level
var DefaultTrustLevel = NewFractionFromTm(tmmath.Fraction{Numerator: 1, Denominator: 3})

// NewFractionFromTm returns a new Fraction instance from a tmmath.Fraction
func NewFractionFromTm(f tmmath.Fraction) Fraction {
	return Fraction{
		Numerator:   uint64(f.Numerator),
		Denominator: uint64(f.Denominator),
	}
}

// ToTendermint converts Fraction to tmmath.Fraction
func (f Fraction) ToTendermint() tmmath.Fraction {
	return tmmath.Fraction{
		Numerator:   int64(f.Numerator),
		Denominator: int64(f.Denominator),
	}
}

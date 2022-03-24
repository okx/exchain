package types

import "math/big"

func (d Dec) ToInt() Int {
	ret := Int{
		i: new(big.Int).Mul(d.BigInt(), precisionMultiplier(0)),
	}
	return ret
}

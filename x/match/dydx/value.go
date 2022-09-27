package dydx

import "math/big"

var (
	ShiftBase = new(big.Float)
)

func init() {
	ShiftBase.SetInt64(1000000000000000000)
}

type BaseValue struct {
	Value *big.Float
}

type Price = BaseValue

type Fee struct {
	BaseValue
}

func FeeFromBips(value *big.Float) Fee {
	num := new(big.Float).SetFloat64(1e-4)
	num.Mul(num, value)
	return Fee{BaseValue{num}}
}

func (v BaseValue) ToSolidity() string {
	if v.Value == nil {
		return ""
	}
	num := new(big.Float)
	num.Copy(v.Value)
	num.Mul(num, ShiftBase)
	num.Abs(num)
	intNum, _ := num.Int(nil)
	return intNum.String()
}

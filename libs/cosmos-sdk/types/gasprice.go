package types

import "math/big"

type DynamicGasInfo struct {
	// gas price
	gp *big.Int
	// gas used
	gu uint64
}

func NewDynamicGasInfo(gp *big.Int, gu uint64) DynamicGasInfo {
	return DynamicGasInfo{
		gp: new(big.Int).Set(gp),
		gu: gu,
	}
}

func NewEmptyDynamicGasInfo() DynamicGasInfo {
	return DynamicGasInfo{
		gp: big.NewInt(0),
		gu: 0,
	}
}

func (info DynamicGasInfo) GetGP() *big.Int {
	return info.gp
}

func (info *DynamicGasInfo) SetGP(gp *big.Int) {
	// deep copy
	gpCopy := new(big.Int).Set(gp)
	info.gp = gpCopy
}

func (info DynamicGasInfo) GetGU() uint64 {
	return info.gu
}

func (info *DynamicGasInfo) SetGU(gu uint64) {
	info.gu = gu
}

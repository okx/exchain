package types

import (
	"math"
	"sync"
)

// Gas consumption descriptors.
const (
	GasIterNextCostFlatDesc = "IterNextFlat"
	GasValuePerByteDesc     = "ValuePerByte"
	GasWritePerByteDesc     = "WritePerByte"
	GasReadPerByteDesc      = "ReadPerByte"
	GasWriteCostFlatDesc    = "WriteFlat"
	GasReadCostFlatDesc     = "ReadFlat"
	GasHasDesc              = "Has"
	GasDeleteDesc           = "Delete"

	defaultHasCost          = 1000
	defaultDeleteCost       = 1000
	defaultReadCostFlat     = 1000
	defaultReadCostPerByte  = 3
	defaultWriteCostFlat    = 2000
	defaultWriteCostPerByte = 30
	defaultIterNextCostFlat = 30
)

var (
	gGasConfig = &GasConfig{
		HasCost:          defaultHasCost,
		DeleteCost:       defaultDeleteCost,
		ReadCostFlat:     defaultReadCostFlat,
		ReadCostPerByte:  defaultReadCostPerByte,
		WriteCostFlat:    defaultWriteCostFlat,
		WriteCostPerByte: defaultWriteCostPerByte,
		IterNextCostFlat: defaultIterNextCostFlat,
	}
	mut = &sync.RWMutex{}
)

// Gas measured by the SDK
type Gas = uint64

// ErrorOutOfGas defines an error thrown when an action results in out of gas.
type ErrorOutOfGas struct {
	Descriptor string
}

// ErrorGasOverflow defines an error thrown when an action results gas consumption
// unsigned integer overflow.
type ErrorGasOverflow struct {
	Descriptor string
}

// GasMeter interface to track gas consumption
type GasMeter interface {
	GasConsumed() Gas
	GasConsumedToLimit() Gas
	Limit() Gas
	ConsumeGas(amount Gas, descriptor string)
	SetGas(val Gas)
	IsPastLimit() bool
	IsOutOfGas() bool
}

type ReusableGasMeter interface {
	GasMeter
	Reset()
}

type basicGasMeter struct {
	limit    Gas
	consumed Gas
}

// NewGasMeter returns a reference to a new basicGasMeter.
func NewGasMeter(limit Gas) GasMeter {
	return &basicGasMeter{
		limit:    limit,
		consumed: 0,
	}
}

func (g *basicGasMeter) GasConsumed() Gas {
	return g.consumed
}

func (g *basicGasMeter) Limit() Gas {
	return g.limit
}

func (g *basicGasMeter) GasConsumedToLimit() Gas {
	if g.IsPastLimit() {
		return g.limit
	}
	return g.consumed
}

// addUint64Overflow performs the addition operation on two uint64 integers and
// returns a boolean on whether or not the result overflows.
func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

func (g *basicGasMeter) ConsumeGas(amount Gas, descriptor string) {
	var overflow bool
	// TODO: Should we set the consumed field after overflow checking?
	g.consumed, overflow = addUint64Overflow(g.consumed, amount)
	if overflow {
		panic(ErrorGasOverflow{descriptor})
	}

	if g.consumed > g.limit {
		panic(ErrorOutOfGas{descriptor})
	}
}

func (g *basicGasMeter) SetGas(val Gas) {
	g.consumed = val
}

func (g *basicGasMeter) IsPastLimit() bool {
	return g.consumed > g.limit
}

func (g *basicGasMeter) IsOutOfGas() bool {
	return g.consumed >= g.limit
}

type infiniteGasMeter struct {
	consumed Gas
}

// NewInfiniteGasMeter returns a reference to a new infiniteGasMeter.
func NewInfiniteGasMeter() GasMeter {
	return &infiniteGasMeter{
		consumed: 0,
	}
}

func NewReusableInfiniteGasMeter() ReusableGasMeter {
	return &infiniteGasMeter{
		consumed: 0,
	}
}

func (g *infiniteGasMeter) Reset() {
	*g = infiniteGasMeter{
		consumed: 0,
	}
}

func (g *infiniteGasMeter) GasConsumed() Gas {
	return g.consumed
}

func (g *infiniteGasMeter) GasConsumedToLimit() Gas {
	return g.consumed
}

func (g *infiniteGasMeter) Limit() Gas {
	return 0
}

func (g *infiniteGasMeter) ConsumeGas(amount Gas, descriptor string) {
	var overflow bool
	// TODO: Should we set the consumed field after overflow checking?
	g.consumed, overflow = addUint64Overflow(g.consumed, amount)
	if overflow {
		panic(ErrorGasOverflow{descriptor})
	}
}

func (g *infiniteGasMeter) SetGas(val Gas) {
	g.consumed = val
}

func (g *infiniteGasMeter) IsPastLimit() bool {
	return false
}

func (g *infiniteGasMeter) IsOutOfGas() bool {
	return false
}

// GasConfig defines gas cost for each operation on KVStores
type GasConfig struct {
	HasCost          Gas `json:"hasCost"`
	DeleteCost       Gas `json:"deleteCost"`
	ReadCostFlat     Gas `json:"readCostFlat"`
	ReadCostPerByte  Gas `json:"readCostPerByte"`
	WriteCostFlat    Gas `json:"writeCostFlat"`
	WriteCostPerByte Gas `json:"writeCostPerByte"`
	IterNextCostFlat Gas `json:"iterNextCostFlat"`
}

// KVGasConfig returns a default gas config for KVStores.
func KVGasConfig() GasConfig {
	return GetGlobalGasConfig()
}

// TransientGasConfig returns a default gas config for TransientStores.
func TransientGasConfig() GasConfig {
	// TODO: define gasconfig for transient stores
	return KVGasConfig()
}

func UpdateGlobalGasConfig(gc *GasConfig) {
	mut.Lock()
	defer mut.Unlock()
	gGasConfig = gc
}

func AsDefaultGasConfig(gc *GasConfig) {
	if gc.HasCost == 0 {
		gc.HasCost = defaultHasCost
	}
	if gc.DeleteCost == 0 {
		gc.DeleteCost = defaultDeleteCost
	}
	if gc.ReadCostFlat == 0 {
		gc.ReadCostFlat = defaultReadCostFlat
	}
	if gc.ReadCostPerByte == 0 {
		gc.ReadCostPerByte = defaultReadCostPerByte
	}
	if gc.WriteCostFlat == 0 {
		gc.WriteCostFlat = defaultWriteCostFlat
	}
	if gc.WriteCostPerByte == 0 {
		gc.WriteCostPerByte = defaultWriteCostPerByte
	}
	if gc.IterNextCostFlat == 0 {
		gc.IterNextCostFlat = defaultIterNextCostFlat
	}
}

func GetGlobalGasConfig() GasConfig {
	mut.RLock()
	defer mut.RUnlock()
	return *gGasConfig
}

func GetDefaultGasConfig() *GasConfig {
	return &GasConfig{
		HasCost:          defaultHasCost,
		DeleteCost:       defaultDeleteCost,
		ReadCostFlat:     defaultReadCostFlat,
		ReadCostPerByte:  defaultReadCostPerByte,
		WriteCostFlat:    defaultWriteCostFlat,
		WriteCostPerByte: defaultWriteCostPerByte,
		IterNextCostFlat: defaultIterNextCostFlat,
	}
}

package types

import (
	"github.com/ethereum/go-ethereum/trie"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

type GasKVStore struct {
	parent   AccStore
	gsConfig GasConfig
	gasMeter GasMeter
}

func NewGasKvStore(parent AccStore, gasConfig GasConfig, gasMeter GasMeter) *GasKVStore {
	return &GasKVStore{
		parent:   parent,
		gsConfig: gasConfig,
		gasMeter: gasMeter,
	}
}

func (gs *GasKVStore) Get(key string) (data []byte) {
	defer func(){
		gs.gasMeter.ConsumeGas(gs.gsConfig.ReadCostFlat, types2.GasReadCostFlatDesc)
		gs.gasMeter.ConsumeGas(gs.gsConfig.ReadCostPerByte*types2.Gas(len(data)), types2.GasReadPerByteDesc)
	}()

	if data = gs.parent.Get(key); data != nil {
		return data
	}

	return nil
}

func (gs *GasKVStore) Set(key string, data []byte) {
	defer func(){
		gs.gasMeter.ConsumeGas(gs.gsConfig.WriteCostFlat, types2.GasWriteCostFlatDesc)
		gs.gasMeter.ConsumeGas(gs.gsConfig.WriteCostPerByte*types2.Gas(len(data)), types2.GasWritePerByteDesc)
	}()

	gs.parent.Set(key, data)
}

func (gs *GasKVStore) Delete(key string) {
	defer func(){
		gs.gasMeter.ConsumeGas(gs.gsConfig.DeleteCost, types2.GasDeleteDesc)
	}()

	gs.parent.Delete(key)
}

type GasIterator struct {
	gasMeter  GasMeter
	gasConfig GasConfig
	parent    *trie.Iterator
}

func (gs *GasKVStore) NewIterator(startKey []byte) *GasIterator {
	itr := gs.parent.(*AccCacheCommitStore).parent.NewIterator(startKey)

	return &GasIterator{
		gasMeter: gs.gasMeter,
		gasConfig: gs.gsConfig,
		parent: itr,
	}
}

func (gi *GasIterator) Next() bool {
	if valid:= gi.parent.Next() ; valid {
		gi.consumeSeekGas()

		return true
	}

	return false
}

// Key implements the Iterator interface. It returns the current key and it does
// not incur any gas cost.
func (gi *GasIterator) Key() (key []byte) {
	return gi.parent.Key
}

// Value implements the Iterator interface. It returns the current value and it
// does not incur any gas cost.
func (gi *GasIterator) Value() (value []byte) {
	return gi.parent.Value
}

func (gi *GasIterator) consumeSeekGas() {
	value := gi.Value()

	gi.gasMeter.ConsumeGas(gi.gasConfig.ReadCostPerByte*types2.Gas(len(value)), types2.GasValuePerByteDesc)
	gi.gasMeter.ConsumeGas(gi.gasConfig.IterNextCostFlat, types2.GasIterNextCostFlatDesc)
}

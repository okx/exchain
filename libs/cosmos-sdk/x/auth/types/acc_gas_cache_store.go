package types

import (
	"github.com/ethereum/go-ethereum/trie"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

type GasKVStore struct {
	parent   types.AccCacheStore
	gsConfig types.GasConfig
	gasMeter types.GasMeter
}

func NewGasKvStore(parent types.AccCacheStore, gasConfig types.GasConfig, gasMeter types.GasMeter) *GasKVStore {
	return &GasKVStore{
		parent:   parent,
		gsConfig: gasConfig,
		gasMeter: gasMeter,
	}
}

func (gs *GasKVStore) Get(addr types.AccAddress) (acc exported.Account) {
	defer func(){
		gs.gasMeter.ConsumeGas(gs.gsConfig.ReadCostFlat, types2.GasReadCostFlatDesc)
		gs.gasMeter.ConsumeGas(gs.gsConfig.ReadCostPerByte*types2.Gas(estimateAccByteLenForGasConsume(acc)), types2.GasReadPerByteDesc)
	}()

	if val := gs.parent.Get(addr); val != nil {
		return val.(exported.Account)
	}

	return nil
}

func (gs *GasKVStore) Set(acc exported.Account) {
	defer func(){
		gs.gasMeter.ConsumeGas(gs.gsConfig.WriteCostFlat, types2.GasWriteCostFlatDesc)
		gs.gasMeter.ConsumeGas(gs.gsConfig.WriteCostPerByte*types2.Gas(estimateAccByteLenForGasConsume(acc)), types2.GasWritePerByteDesc)
	}()

	gs.parent.Set(acc.GetAddress(), acc)
}

func (gs *GasKVStore) Delete(acc exported.Account) {
	defer func(){
		gs.gasMeter.ConsumeGas(gs.gsConfig.DeleteCost, types2.GasDeleteDesc)
	}()

	gs.parent.Delete(acc.GetAddress())
}

func (gs *GasKVStore) Write() {
	gs.parent.Write()
}

func estimateAccByteLenForGasConsume(acc exported.Account) int64{
	if acc == nil {
		return 0
	}

	if acc.IsEthAccount() {
		return 150
	}

	return 70
}

type GasIterator struct {
	gasMeter  types.GasMeter
	gasConfig types.GasConfig
	parent    *trie.Iterator
}

func (gs *GasKVStore) NewIterator(startKey []byte) *GasIterator {
	itr := gs.parent.(*CacheStore).parent.NewIterator(startKey)

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

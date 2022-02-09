package keeper

import (
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/feemarket/types"
)

// Keeper grants access to the Fee Market module state.
type Keeper struct {
	// Store key required for the Fee Market Prefix KVStore.
	storeKey sdk.StoreKey
	// module specific parameter space that can be configured through governance
	paramSpace Subspace
}

type Subspace interface {
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

// NewKeeper generates new fee market module keeper
func NewKeeper(storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// ----------------------------------------------------------------------------
// Parent Block Gas Used
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// GetBlockGasUsed returns the last block gas used value from the store.
func (k Keeper) GetBlockGasUsed(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixBlockGasUsed)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// SetBlockGasUsed gets the block gas consumed to the store.
// CONTRACT: this should be only called during EndBlock.
func (k Keeper) SetBlockGasUsed(ctx sdk.Context, gas uint64) {
	store := ctx.KVStore(k.storeKey)
	gasBz := sdk.Uint64ToBigEndian(gas)
	store.Set(types.KeyPrefixBlockGasUsed, gasBz)
}

// ----------------------------------------------------------------------------
// Parent Base Fee
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// GetBaseFee returns the last base fee value from the store.
// returns nil if base fee is not enabled.
func (k Keeper) GetBaseFee(ctx sdk.Context) *big.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixBaseFee)
	if len(bz) == 0 {
		return nil
	}

	return new(big.Int).SetBytes(bz)
}

// SetBaseFee set the last base fee value to the store.
// CONTRACT: this should be only called during EndBlock.
func (k Keeper) SetBaseFee(ctx sdk.Context, baseFee *big.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefixBaseFee, baseFee.Bytes())
}

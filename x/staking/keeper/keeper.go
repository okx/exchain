package keeper

import (
	"fmt"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	"strings"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/staking/exported"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/params"
	"github.com/okex/exchain/x/staking/types"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Keeper is the keeper struct of the staking store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	supplyKeeper types.SupplyKeeper
	hooks        types.StakingHooks
	paramstore   params.Subspace
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace) Keeper {

	// set KeyTable if it has not already been set
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(ParamKeyTable())
	}
	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		paramstore:  paramstore,
		hooks:        nil,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// SetHooks sets the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

// Codespace returns the codespace
func (k Keeper) Codespace() string {
	return types.ModuleName
}

// GetLastTotalPower loads the last total validator power
func (k Keeper) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// SetLastTotalPower sets the last total validator power
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.LastTotalPowerKey, b)
}

// IsValidator tells whether a validator is in bonded status
func (k Keeper) IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool {
	var curValidators []string
	// fetch all the bonded validators, insert them into currValidators
	k.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		curValidators = append(curValidators, validator.GetOperator().String())
		return false
	})

	valStr := sdk.ValAddress(addr).String()
	for _, val := range curValidators {
		if valStr == val {
			return true
		}
	}
	return false
}

// GetOperAddrFromValidatorAddr returns the validator address according to the consensus pubkey
// the validator has to exist
func (k Keeper) GetOperAddrFromValidatorAddr(ctx sdk.Context, va string) (sdk.ValAddress, bool) {
	validators := k.GetAllValidators(ctx)

	for _, validator := range validators {
		if strings.Compare(strings.ToUpper(va), validator.ConsPubKey.Address().String()) == 0 {
			return validator.OperatorAddress, true
		}
	}
	return nil, false
}

// GetOperAndValidatorAddr returns the operator addresses and consensus pubkeys of all the validators
func (k Keeper) GetOperAndValidatorAddr(ctx sdk.Context) types.OVPairs {
	validators := k.GetAllValidators(ctx)
	var ovPairs types.OVPairs

	for _, validator := range validators {
		ovPair := types.OVPair{OperAddr: validator.OperatorAddress, ValAddr: validator.ConsPubKey.Address().String()}
		ovPairs = append(ovPairs, ovPair)
	}
	return ovPairs
}


// GetHistoricalInfo gets the historical info at a given height
func (k Keeper) GetHistoricalInfo(ctx sdk.Context, height int64) (types2.HistoricalInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types2.GetHistoricalInfoKey(height)

	value := store.Get(key)
	if value == nil {
		return types2.HistoricalInfo{}, false
	}

	return types2.MustUnmarshalHistoricalInfo(k.cdc, value), true
}

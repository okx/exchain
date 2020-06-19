package keeper

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/okex/okchain/x/staking/exported"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/okex/okchain/x/staking/types"
)

const aminoCacheSize = 500

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Keeper is the keeper struct of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	storeTKey          sdk.StoreKey
	cdc                *codec.Codec
	supplyKeeper       types.SupplyKeeper
	hooks              types.StakingHooks
	paramstore         params.Subspace
	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {

	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:           key,
		storeTKey:          tkey,
		cdc:                cdc,
		supplyKeeper:       supplyKeeper,
		paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		hooks:              nil,
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
		codespace:          codespace,
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
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
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

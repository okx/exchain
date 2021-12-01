package keeper

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/slashing/internal/types"
)

// Keeper of the slashing store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	sk         types.StakingKeeper
	paramspace types.ParamSubspace
	cache      *paramCache
}
type paramCache struct {
	mpValSingInfo    map[ethcmn.Address]types.ValidatorSigningInfo
	mpValSingInfoGas map[ethcmn.Address]uint64
	pub              map[ethcmn.Address]crypto.PubKey
	pubGas           map[ethcmn.Address]uint64
}

func newParamcache() *paramCache {
	return &paramCache{
		mpValSingInfo:    make(map[ethcmn.Address]types.ValidatorSigningInfo, 0),
		mpValSingInfoGas: make(map[ethcmn.Address]uint64, 0),
		pub:              make(map[ethcmn.Address]crypto.PubKey),
		pubGas:           make(map[ethcmn.Address]uint64, 0),
	}
}
func (p *paramCache) getCacheValSignInfo(addr []byte) (types.ValidatorSigningInfo, uint64, bool) {
	ethAddr := ethcmn.BytesToAddress(addr)
	if data, ok := p.mpValSingInfo[ethAddr]; ok {
		return data, p.mpValSingInfoGas[ethAddr], true
	}
	return types.ValidatorSigningInfo{}, 0, false

}

func (p *paramCache) setCacheValSignInfo(addr []byte, value types.ValidatorSigningInfo, gas uint64) {
	ethAddr := ethcmn.BytesToAddress(addr)
	p.mpValSingInfo[ethAddr] = value
	p.mpValSingInfoGas[ethAddr] = gas
}

func (p *paramCache) getpub(addr []byte) (crypto.PubKey, uint64, bool) {
	ethAddr := ethcmn.BytesToAddress(addr)
	if data, ok := p.pub[ethAddr]; ok {
		return data, p.pubGas[ethAddr], true
	}
	return nil, 0, false
}
func (p *paramCache) setPub(addr []byte, pub crypto.PubKey, gas uint64) {
	ethAddr := ethcmn.BytesToAddress(addr)
	p.pub[ethAddr] = pub
	p.pubGas[ethAddr] = gas
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk types.StakingKeeper, paramspace types.ParamSubspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		sk:         sk,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
		cache:      newParamcache(),
	}
}

//Get Staking keeper object
func (k Keeper) GetStakingKeeper() types.StakingKeeper {
	return k.sk
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// AddPubkey sets a address-pubkey relation
func (k Keeper) AddPubkey(ctx sdk.Context, pubkey crypto.PubKey) {
	addr := pubkey.Address()
	k.setAddrPubkeyRelation(ctx, addr, pubkey)
}

// GetPubkey returns the pubkey from the adddress-pubkey relation
func (k Keeper) GetPubkey(ctx sdk.Context, address crypto.Address) (crypto.PubKey, error) {
	if data, gas, ok := k.cache.getpub(address); ok {
		ctx.GasMeter().ConsumeGas(gas, "exchain/x/slashing/internal/keeper/GetPubkey")
		return data, nil
	}
	beforeGas := ctx.GasMeter().GasConsumed()
	store := ctx.KVStore(k.storeKey)
	var pubkey crypto.PubKey
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.GetAddrPubkeyRelationKey(address)), &pubkey)
	if err != nil {
		return nil, fmt.Errorf("address %s not found", sdk.ConsAddress(address))
	}
	k.cache.setPub(address, pubkey, ctx.GasMeter().GasConsumed()-beforeGas)
	return pubkey, nil
}

// Slash attempts to slash a validator. The slash is delegated to the staking
// module to make the necessary validator changes.
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, fraction sdk.Dec, power, distributionHeight int64) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
			sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
			sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueDoubleSign),
		),
	)

	k.sk.Slash(ctx, consAddr, distributionHeight, power, fraction)
}

// Jail attempts to jail a validator. The slash is delegated to the staking module
// to make the necessary validator changes.
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute(types.AttributeKeyJailed, consAddr.String()),
		),
	)

	k.sk.Jail(ctx, consAddr)
}

func (k Keeper) setAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address, pubkey crypto.PubKey) {
	beforeGas := ctx.GasMeter().GasConsumed()
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pubkey)
	store.Set(types.GetAddrPubkeyRelationKey(addr), bz)
	k.cache.setPub(addr, pubkey, ctx.GasMeter().GasConsumed()-beforeGas)
}

func (k Keeper) deleteAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address) {
	beforeGas := ctx.GasMeter().GasConsumed()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAddrPubkeyRelationKey(addr))
	k.cache.setPub(addr, nil, ctx.GasMeter().GasConsumed()-beforeGas)
}

package params

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkparams "github.com/okex/exchain/libs/cosmos-sdk/x/params"

	"github.com/okex/exchain/x/params/types"
)

// Keeper is the struct of params keeper
type Keeper struct {
	cdc *codec.Codec
	sdkparams.Keeper
	// the reference to the Paramstore to get and set gov specific params
	paramSpace sdkparams.Subspace
	// the reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	sk StakingKeeper
	// the reference to the CoinKeeper to modify balances
	ck BankKeeper
	// the reference to the GovKeeper to insert waiting queue
	gk      GovKeeper
	signals []func()

	storeKey        *sdk.KVStoreKey
	upgradeReadyMap map[string]func(types.UpgradeInfo)
}

// NewKeeper creates a new instance of params keeper
func NewKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey) (
	k Keeper) {
	k = Keeper{
		Keeper:  sdkparams.NewKeeper(cdc, key, tkey),
		signals: make([]func(), 0),

		storeKey:        key,
		upgradeReadyMap: make(map[string]func(types.UpgradeInfo)),
	}
	k.cdc = cdc
	k.paramSpace = k.Subspace(DefaultParamspace).WithKeyTable(types.ParamKeyTable())
	return k
}

// SetStakingKeeper hooks the staking keeper into params keeper
func (keeper *Keeper) SetStakingKeeper(sk StakingKeeper) {
	keeper.sk = sk
}

// SetBankKeeper hooks the bank keeper into params keeper
func (keeper *Keeper) SetBankKeeper(ck BankKeeper) {
	keeper.ck = ck
}

// SetGovKeeper hooks the gov keeper into params keeper
func (keeper *Keeper) SetGovKeeper(gk GovKeeper) {
	keeper.gk = gk
}

// SetParams sets the params into the store
func (keeper *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	keeper.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the params info from the store
func (keeper Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	keeper.paramSpace.GetParamSet(ctx, &params)
	return params
}

// ClaimReadyForUpgrade tells Keeper that someone has get ready for the upgrade.
// cb could be nil if there's no code to be execute when the upgrade is take effective.
func (keeper *Keeper) ClaimReadyForUpgrade(ctx sdk.Context, name string, cb func(types.UpgradeInfo)) {
	if keeper.IsUpgradeEffective(ctx, name) {
		keeper.Logger(ctx).Info("upgrade has been effective, ready for it will do nothing", "upgrade name", name)
	}

	if _, ok := keeper.upgradeReadyMap[name]; ok {
		keeper.Logger(ctx).Error("more than one guys ready for the same upgrade, the front one will be cover", "upgrade name", name)
	}
	keeper.upgradeReadyMap[name] = cb
}

func (keeper *Keeper) QueryReadyForUpgrade(name string) (func(types.UpgradeInfo), bool) {
	cb, ok := keeper.upgradeReadyMap[name]
	return cb, ok
}

func (keeper *Keeper) IsUpgradeEffective(ctx sdk.Context, name string) bool {
	_, err := keeper.GetEffectiveUpgradeInfo(ctx, name)
	return err == nil
}

func (keeper *Keeper) GetEffectiveUpgradeInfo(ctx sdk.Context, name string) (types.UpgradeInfo, error) {
	exist, info, err := getUpgradeInfo(ctx, keeper, name)
	if err != nil {
		return types.UpgradeInfo{}, err
	}
	if !exist {
		return types.UpgradeInfo{}, fmt.Errorf("upgrade '%s' is not exist", name)
	}

	if !isUpgradeEffective(ctx, info) {
		keeper.Logger(ctx).Debug("upgrade is not effective", "name", name)
		return types.UpgradeInfo{}, fmt.Errorf("upgrade '%s' is not effective", name)
	}

	keeper.Logger(ctx).Debug("upgrade is effective", "name", name)
	return info, nil
}

func (keeper *Keeper) GetUpgradeInfo(ctx sdk.Context, name string) (bool, types.UpgradeInfo, error) {
	exist, info, err := getUpgradeInfo(ctx, keeper, name)
	if err != nil {
		return false, types.UpgradeInfo{}, err
	}
	if !exist {
		return false, types.UpgradeInfo{}, fmt.Errorf("upgrade '%s' is not exist", name)
	}

	return isUpgradeEffective(ctx, info), info, nil
}

func isUpgradeEffective(ctx sdk.Context, info types.UpgradeInfo) bool {
	return info.Status == types.UpgradeStatusEffective && uint64(ctx.BlockHeight()) >= info.EffectiveHeight
}

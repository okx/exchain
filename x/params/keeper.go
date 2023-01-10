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

	storeKey *sdk.KVStoreKey
}

// NewKeeper creates a new instance of params keeper
func NewKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey) (
	k Keeper) {
	k = Keeper{
		Keeper:   sdkparams.NewKeeper(cdc, key, tkey),
		signals:  make([]func(), 0),
		storeKey: key,
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
	return uint64(ctx.BlockHeight()) >= info.EffectiveHeight
}

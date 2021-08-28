package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/common/monitor"
	"github.com/okex/exchain/x/params"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/exchain/x/distribution/types"
)

// Keeper of the distribution store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSpace    params.Subspace
	stakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper

	blacklistedAddrs map[string]bool

	feeCollectorName string // name of the FeeCollector ModuleAccount

	metric              *monitor.DistrMetric
	monitoredValidators []string
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	sk types.StakingKeeper, supplyKeeper types.SupplyKeeper, feeCollectorName string,
	blacklistedAddrs map[string]bool, metrics *monitor.DistrMetric,
) Keeper {

	// ensure distribution module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	metrics.TotalFee.Add(viper.GetFloat64("test.init_totoal_fee"))
	metrics.FeeToControlledVals.Add(viper.GetFloat64("test.init_control_fee"))
	metrics.FeeToOtherVals.Add(viper.GetFloat64("test.init_other_fee"))
	metrics.FeeToCommunityPool.Add(viper.GetFloat64("test.init_community_fee"))
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace,
		stakingKeeper:    sk,
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
		blacklistedAddrs: blacklistedAddrs,
		metric:              metrics,
		monitoredValidators: viper.GetStringSlice("test.monitored_validators"),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ShortUseByCli)
}

// SetWithdrawAddr sets a new address that will receive the rewards upon withdrawal
func (k Keeper) SetWithdrawAddr(ctx sdk.Context, delegatorAddr sdk.AccAddress, withdrawAddr sdk.AccAddress) error {
	if k.blacklistedAddrs[withdrawAddr.String()] {
		return types.ErrWithdrawAddrInblacklist()
	}

	if !k.GetWithdrawAddrEnabled(ctx) {
		return types.ErrSetWithdrawAddrDisabled()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetWithdrawAddress,
			sdk.NewAttribute(types.AttributeKeyWithdrawAddress, withdrawAddr.String()),
		),
	)

	k.SetDelegatorWithdrawAddr(ctx, delegatorAddr, withdrawAddr)
	return nil
}

// WithdrawValidatorCommission withdraws validator commission
func (k Keeper) WithdrawValidatorCommission(ctx sdk.Context, valAddr sdk.ValAddress) (sdk.Coins, error) {
	// fetch validator accumulated commission
	accumCommission := k.GetValidatorAccumulatedCommission(ctx, valAddr)
	if accumCommission.IsZero() {
		return nil, types.ErrNoValidatorCommission()
	}

	commission, remainder := accumCommission.TruncateDecimal()
	k.SetValidatorAccumulatedCommission(ctx, valAddr, remainder) // leave remainder to withdraw later

	if !commission.IsZero() {
		accAddr := sdk.AccAddress(valAddr)
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, accAddr)
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, commission)
		if err != nil {
			return nil, types.ErrSendCoinsFromModuleToAccountFailed()
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		),
	)

	return commission, nil
}

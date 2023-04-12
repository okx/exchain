package keeper

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/common/monitor"
	"github.com/okex/exchain/x/params"
	"github.com/spf13/viper"

	"github.com/okex/exchain/x/distribution/types"
	stakingexported "github.com/okex/exchain/x/staking/exported"
	stakingtyps "github.com/okex/exchain/x/staking/types"
)

// Keeper of the distribution store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSpace    params.Subspace
	stakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper
	govKeeper     types.GovKeeper

	blacklistedAddrs map[string]bool

	feeCollectorName string // name of the FeeCollector ModuleAccount

	feeInfo             *FeeInfo
	metric              *monitor.DistrMetric
	monitoredValidators []string
	monitoredDelegators []string
}

type FeeInfo struct {
	TotalFee                       float64
	FeeToControlledVals            float64
	FeeToOtherVals                 float64
	FeeToCommunityPool             float64
	OfficialValidatorCommissionOKT float64
	OfficialDelegatorRewardsOKT    float64
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

	feeInfo := FeeInfo{
		viper.GetFloat64("test.init_totoal_fee"),
		viper.GetFloat64("test.init_control_fee"),
		viper.GetFloat64("test.init_other_fee"),
		viper.GetFloat64("test.init_community_fee"),
		0,
		0,
	}
	metrics.TotalFee.Set(feeInfo.TotalFee)
	metrics.FeeToControlledVals.Set(feeInfo.FeeToControlledVals)
	metrics.FeeToOtherVals.Set(feeInfo.FeeToOtherVals)
	metrics.FeeToCommunityPool.Set(feeInfo.FeeToCommunityPool)

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace,
		stakingKeeper:    sk,
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
		blacklistedAddrs: blacklistedAddrs,

		feeInfo:             &feeInfo,
		metric:              metrics,
		monitoredValidators: viper.GetStringSlice("test.monitored_validators"),
		monitoredDelegators: viper.GetStringSlice("test.monitored_delegators"),
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

	if k.CheckDistributionProposalValid(ctx) {
		// update outstanding
		outstanding := k.GetValidatorOutstandingRewards(ctx, valAddr)
		k.SetValidatorOutstandingRewards(ctx, valAddr, outstanding.Sub(sdk.NewDecCoinsFromCoins(commission...)))
	}

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

func (k Keeper) CheckStatistics(ctx sdk.Context) {
	logger := k.Logger(ctx)
	valFilter := getFilterFromAddress(k.monitoredValidators)
	delegatorFilter := getFilterFromAddress(k.monitoredDelegators)

	officialValidatorCommissionOKT, officialDelegatorRewardsOKT := sdk.ZeroDec(), sdk.ZeroDec()

	// iterate validators
	k.stakingKeeper.IterateValidators(ctx, func(index int64, val stakingexported.ValidatorI) (stop bool) {
		if _, ok := valFilter[val.GetOperator().String()]; ok {
			commissions := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
			if len(commissions) > 0 {
				officialValidatorCommissionOKT = officialValidatorCommissionOKT.Add(commissions[0].Amount)
			}
		}
		return false
	})

	// iterate delegators
	k.stakingKeeper.IterateDelegator(ctx, func(index int64, delegator stakingtyps.Delegator) bool {
		if _, ok := delegatorFilter[delegator.GetDelegatorAddress().String()]; ok {
			rewards := k.queryDelegatorRewards(ctx, delegator.GetDelegatorAddress())
			if len(rewards) > 0 {
				officialDelegatorRewardsOKT = officialDelegatorRewardsOKT.Add(rewards[0].Amount)
			}
		}
		return false
	})

	k.feeInfo.OfficialValidatorCommissionOKT = sdk.ConvertDecToFloat64(officialValidatorCommissionOKT)
	k.feeInfo.OfficialDelegatorRewardsOKT = sdk.ConvertDecToFloat64(officialDelegatorRewardsOKT)

	logger.Error("Distribution okt.", "official_validator_commission_okt", officialValidatorCommissionOKT,
		"official_delegator_rewards_okt", officialDelegatorRewardsOKT,
	)
}

func (k Keeper) queryDelegatorRewards(ctx sdk.Context, delegatorAddress sdk.AccAddress) sdk.DecCoins {
	ctx, _ = ctx.CacheContext()
	total := sdk.DecCoins{}
	del := k.stakingKeeper.Delegator(ctx, delegatorAddress)
	if del == nil {
		return total
	}

	for _, valAddr := range del.GetShareAddedValidatorAddresses() {
		val := k.stakingKeeper.Validator(ctx, valAddr)
		if val == nil {
			continue
		}

		if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), delegatorAddress) {
			if del.GetLastAddedShares().IsZero() {
				return total
			}
			k.initExistedDelegationStartInfo(ctx, val, del)
		}

		endingPeriod := k.incrementValidatorPeriod(ctx, val)
		delReward := k.calculateDelegationRewards(ctx, val, delegatorAddress, endingPeriod)
		if delReward == nil {
			delReward = sdk.DecCoins{}
		}
		total = total.Add(delReward...)
	}

	return total
}

// build a filter
func getFilterFromAddress(addrs []string) map[string]struct{} {
	valLen := len(addrs)
	valFilter := make(map[string]struct{}, valLen)
	for i := 0; i < valLen; i++ {
		valFilter[addrs[i]] = struct{}{}
	}

	return valFilter
}

func (k Keeper) GetOfficeRewards() float64 {
	return k.feeInfo.OfficialValidatorCommissionOKT + k.feeInfo.OfficialDelegatorRewardsOKT
}

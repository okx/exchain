package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okexchain/x/staking/exported"
	"github.com/okex/okexchain/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// InitGenesis sets the pool and parameters for the provided keeper
// For each validator in data, it sets that validator in the keeper along with manually setting the indexes
// In addition, it also sets any delegations found in data
// Finally, it updates the bonded validators
// Returns final validator set after applying all declaration and delegations
func InitGenesis(ctx sdk.Context, keeper Keeper, accountKeeper types.AccountKeeper,
	supplyKeeper types.SupplyKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	bondedTokens, notBondedTokens := sdk.ZeroDec(), sdk.ZeroDec()

	// We need to pretend to be "n blocks before genesis", where "n" is the validator update delay, so that e.g.
	// slashing periods are correctly initialized for the validator set e.g. with a one-block offset - the first
	// TM block is at height 1, so state updates applied from genesis.json are in block 0.
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)
	keeper.SetParams(ctx, data.Params)
	keeper.SetLastTotalPower(ctx, data.LastTotalPower)

	for _, validator := range data.Validators {
		initValidator(ctx, validator, keeper, &bondedTokens, data.Exported)
	}

	for _, delegator := range data.Delegators {
		initDelegator(ctx, delegator, keeper, &bondedTokens)
	}

	for _, ubd := range data.UnbondingDelegations {
		initUnbondingDelegation(ctx, ubd, keeper, &notBondedTokens)
	}
	for _, sharesExported := range data.AllShares {
		keeper.SetShares(ctx, sharesExported.DelAddress, sharesExported.ValidatorAddress, sharesExported.Shares)
	}
	for _, proxyDelegatorKeyExported := range data.ProxyDelegatorKeys {
		keeper.SetProxyBinding(ctx, proxyDelegatorKeyExported.ProxyAddr, proxyDelegatorKeyExported.DelAddr, false)
	}

	checkPools(ctx, keeper, sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, bondedTokens),
		sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, notBondedTokens), data.Exported)

	// don't need to run Tendermint updates if we exported
	if data.Exported {
		for _, lv := range data.LastValidatorPowers {
			keeper.SetLastValidatorPower(ctx, lv.Address, lv.Power)
			validator, found := keeper.GetValidator(ctx, lv.Address)
			if !found {
				panic(fmt.Sprintf("validator %s not found", lv.Address))
			}
			update := validator.ABCIValidatorUpdate()
			update.Power = lv.Power // keep the next-val-set offset, use the last power for the first block
			res = append(res, update)
		}
	} else {
		res = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	return res
}

// assume that there is only okt in pool, if not panics
func checkTokenSum(tokenSum sdk.SysCoin, pool supplyexported.ModuleAccountI) {
	poolCoins := pool.GetCoins()
	if !poolCoins.IsZero() {
		if len(poolCoins) != 1 {
			panic(fmt.Sprintf("only okt in %s, but there are %d kinds of coins", pool.GetName(), len(poolCoins)))
		}

		if !tokenSum.ToCoins().IsEqual(poolCoins) {
			panic(fmt.Sprintf("coins in %s don't match the token sum, tokenSum: %s, poolCoins: %s",
				pool.GetName(), tokenSum.String(), poolCoins.String()))
		}
	}
}

func checkPools(ctx sdk.Context, keeper Keeper, bondedDecCoin, notBondedDecCoin sdk.SysCoin, isExported bool) {
	bondedPool := keeper.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}
	notBondedPool := keeper.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}
	if isExported {
		checkTokenSum(bondedDecCoin, bondedPool)
		checkTokenSum(notBondedDecCoin, notBondedPool)
	}
}

func initUnbondingDelegation(ctx sdk.Context, ubd UndelegationInfo, keeper Keeper, notBondedTokens *sdk.Dec) {
	keeper.SetUndelegating(ctx, ubd)
	keeper.SetAddrByTimeKeyWithNilValue(ctx, ubd.CompletionTime, ubd.DelegatorAddress)
	*notBondedTokens = notBondedTokens.Add(ubd.Quantity)
}

func initDelegator(ctx sdk.Context, delegator Delegator, keeper Keeper, pBondedTokens *sdk.Dec) {
	keeper.SetDelegator(ctx, delegator)
	*pBondedTokens = pBondedTokens.Add(delegator.Tokens)
}

func initValidator(ctx sdk.Context, valExported ValidatorExport, keeper Keeper, pBondedTokens *sdk.Dec, exported bool) {
	validator := valExported.Import()
	keeper.SetValidator(ctx, validator)

	// manually set indices for the first time
	keeper.SetValidatorByConsAddr(ctx, validator)
	keeper.SetValidatorByPowerIndex(ctx, validator)

	// call the creation hook if not exported
	if !exported {
		keeper.AfterValidatorCreated(ctx, validator.OperatorAddress)
	}

	// update timeslice if necessary
	if validator.IsUnbonding() {
		keeper.InsertValidatorQueue(ctx, validator)
	}
	// all the msd on validator should be added into bonded pool
	*pBondedTokens = pBondedTokens.Add(validator.MinSelfDelegation)
}

// ExportGenesis returns a GenesisState for a given context and keeper
// The GenesisState will contain the pool, params, validators, and bonds found in the keeper
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	lastTotalPower := keeper.GetLastTotalPower(ctx)
	validators := keeper.GetAllValidators(ctx)
	var delegators []types.Delegator
	keeper.IterateDelegator(ctx, func(_ int64, delegator types.Delegator) (stop bool) {
		delegators = append(delegators, delegator)
		return false
	})
	var undelegationInfos []types.UndelegationInfo
	keeper.IterateUndelegationInfo(ctx, func(_ int64, ubd types.UndelegationInfo) (stop bool) {
		undelegationInfos = append(undelegationInfos, ubd)
		return false
	})
	var lastValidatorPowers []types.LastValidatorPower
	keeper.IterateLastValidatorPowers(ctx, func(addr sdk.ValAddress, power int64) (stop bool) {
		lastValidatorPowers = append(lastValidatorPowers, types.NewLastValidatorPower(addr, power))
		return false
	})
	var sharesExportedSlice []types.SharesExported
	keeper.IterateShares(ctx,
		func(_ int64, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares types.Shares) (stop bool) {
			sharesExportedSlice = append(sharesExportedSlice, types.NewSharesExported(delAddr, valAddr, shares))
			return false
		})

	var proxyDelegatorKeys []types.ProxyDelegatorKeyExported
	keeper.IterateProxy(ctx, []byte{}, false, func(_ int64, delAddr, proxyAddr sdk.AccAddress) (stop bool) {
		proxyDelegatorKeys = append(proxyDelegatorKeys, types.NewProxyDelegatorKeyExported(delAddr, proxyAddr))
		return false
	})

	return types.GenesisState{
		Params:               params,
		LastTotalPower:       lastTotalPower,
		LastValidatorPowers:  lastValidatorPowers,
		Validators:           validators.Export(),
		Delegators:           delegators,
		UnbondingDelegations: undelegationInfos,
		AllShares:            sharesExportedSlice,
		ProxyDelegatorKeys:   proxyDelegatorKeys,
		Exported:             true,
	}
}

// GetLatestGenesisValidator returns a slice of bonded genesis validators
func GetLatestGenesisValidator(ctx sdk.Context, keeper Keeper) (vals []tmtypes.GenesisValidator) {
	keeper.IterateLastValidators(ctx, func(_ int64, validator exported.ValidatorI) (stop bool) {
		vals = append(vals, tmtypes.GenesisValidator{
			PubKey: validator.GetConsPubKey(),
			Power:  validator.GetConsensusPower(),
			Name:   validator.GetMoniker(),
		})

		return false
	})

	return
}

// ValidateGenesis validates the provided staking genesis state to ensure the expected invariants holds
// (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data types.GenesisState) error {
	err := validateGenesisStateValidators(data.Validators)
	if err != nil {
		return err
	}
	return data.Params.Validate()
}

func validateGenesisStateValidators(valsExported []types.ValidatorExported) (err error) {
	valsLen := len(valsExported)
	addrMap := make(map[string]bool, valsLen)
	for i := 0; i < valsLen; i++ {
		valExported := valsExported[i]
		strKey := valExported.ConsPubKey
		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: moniker %v, address %v",
				valExported.Description.Moniker, valExported.ConsAddress())
		}
		if valExported.Jailed && valExported.IsBonded() {
			return fmt.Errorf("validator is bonded and jailed in genesis state: moniker %v, address %v",
				valExported.Description.Moniker, valExported.ConsAddress())
		}
		if valExported.DelegatorShares.IsZero() {
			return fmt.Errorf("it's impossible for a validator with zero delegator shares, validator: %v", valExported)
		}
		addrMap[strKey] = true
	}
	return
}

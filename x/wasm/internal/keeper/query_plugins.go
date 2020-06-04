package keeper

import (
	"encoding/json"
	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/okex/okchain/x/wasm/internal/types"
)

type QueryHandler struct {
	Ctx     sdk.Context
	Plugins QueryPlugins
}

var _ wasmTypes.Querier = QueryHandler{}

func (q QueryHandler) Query(request wasmTypes.QueryRequest) ([]byte, error) {
	if request.Bank != nil {
		return q.Plugins.Bank(q.Ctx, request.Bank)
	}
	if request.Custom != nil {
		return q.Plugins.Custom(q.Ctx, request.Custom)
	}
	if request.Staking != nil {
		return q.Plugins.Staking(q.Ctx, request.Staking)
	}
	if request.Wasm != nil {
		return q.Plugins.Wasm(q.Ctx, request.Wasm)
	}
	return nil, wasmTypes.Unknown{}
}

type CustomQuerier func(ctx sdk.Context, request json.RawMessage) ([]byte, error)

type QueryPlugins struct {
	Bank    func(ctx sdk.Context, request *wasmTypes.BankQuery) ([]byte, error)
	Custom  CustomQuerier
	Staking func(ctx sdk.Context, request *wasmTypes.StakingQuery) ([]byte, error)
	Wasm    func(ctx sdk.Context, request *wasmTypes.WasmQuery) ([]byte, error)
}

func DefaultQueryPlugins(bank bank.Keeper, wasm Keeper) QueryPlugins {
	return QueryPlugins{
		Bank:    BankQuerier(bank),
		Custom:  NoCustomQuerier,
		Wasm:    WasmQuerier(wasm),
	}
}

func (e QueryPlugins) Merge(o *QueryPlugins) QueryPlugins {
	// only update if this is non-nil and then only set values
	if o == nil {
		return e
	}
	if o.Bank != nil {
		e.Bank = o.Bank
	}
	if o.Custom != nil {
		e.Custom = o.Custom
	}
	if o.Staking != nil {
		e.Staking = o.Staking
	}
	if o.Wasm != nil {
		e.Wasm = o.Wasm
	}
	return e
}

func BankQuerier(bank bank.Keeper) func(ctx sdk.Context, request *wasmTypes.BankQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmTypes.BankQuery) ([]byte, error) {
		if request.AllBalances != nil {
			addr, err := sdk.AccAddressFromBech32(request.AllBalances.Address)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(request.AllBalances.Address)
			}
			coins := bank.GetCoins(ctx, addr)
			res := wasmTypes.AllBalancesResponse{
				Amount: convertSdkCoinsToWasmCoins(coins),
			}
			return json.Marshal(res)
		}
		if request.Balance != nil {
			addr, err := sdk.AccAddressFromBech32(request.Balance.Address)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(request.Balance.Address)
			}
			coins := bank.GetCoins(ctx, addr)
			amount := coins.AmountOf(request.Balance.Denom)
			res := wasmTypes.BalanceResponse{
				Amount: wasmTypes.Coin{
					Denom:  request.Balance.Denom,
					Amount: types.ConvertSdkDecToString(amount),
				},
			}
			return json.Marshal(res)
		}
		return nil, wasmTypes.UnsupportedRequest{"unknown BankQuery variant"}
	}
}

func NoCustomQuerier(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return nil, wasmTypes.UnsupportedRequest{"custom"}
}

func StakingQuerier(keeper staking.Keeper) func(ctx sdk.Context, request *wasmTypes.StakingQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmTypes.StakingQuery) ([]byte, error) {
		if request.BondedDenom != nil {
			denom := keeper.BondDenom(ctx)
			res := wasmTypes.BondedDenomResponse{
				Denom: denom,
			}
			return json.Marshal(res)
		}
		if request.Validators != nil {
			validators := keeper.GetBondedValidatorsByPower(ctx)
			//validators := keeper.GetAllValidators(ctx)
			wasmVals := make([]wasmTypes.Validator, len(validators))
			for i, v := range validators {
				wasmVals[i] = wasmTypes.Validator{
					Address:       v.OperatorAddress.String(),
					Commission:    v.Commission.Rate.String(),
					MaxCommission: v.Commission.MaxRate.String(),
					MaxChangeRate: v.Commission.MaxChangeRate.String(),
				}
			}
			res := wasmTypes.ValidatorsResponse{
				Validators: wasmVals,
			}
			return json.Marshal(res)
		}
		if request.AllDelegations != nil {
			delegator, err := sdk.AccAddressFromBech32(request.AllDelegations.Delegator)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(request.AllDelegations.Delegator)
			}
			sdkDels := keeper.GetAllDelegatorDelegations(ctx, delegator)
			delegations, err := sdkToDelegations(ctx, keeper, sdkDels)
			if err != nil {
				return nil, err
			}
			res := wasmTypes.AllDelegationsResponse{
				Delegations: delegations,
			}
			return json.Marshal(res)
		}
		if request.Delegation != nil {
			delegator, err := sdk.AccAddressFromBech32(request.Delegation.Delegator)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(request.Delegation.Delegator)
			}
			validator, err := sdk.ValAddressFromBech32(request.Delegation.Validator)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(request.Delegation.Validator)
			}

			var res wasmTypes.DelegationResponse
			d, found := keeper.GetDelegation(ctx, delegator, validator)
			if found {
				res.Delegation, err = sdkToFullDelegation(ctx, keeper, d)
				if err != nil {
					return nil, err
				}
			}
			return json.Marshal(res)
		}
		return nil, wasmTypes.UnsupportedRequest{"unknown Staking variant"}
	}
}

func sdkToDelegations(ctx sdk.Context, keeper staking.Keeper, delegations []staking.Delegation) (wasmTypes.Delegations, error) {
	result := make([]wasmTypes.Delegation, len(delegations))
	bondDenom := keeper.BondDenom(ctx)

	for i, d := range delegations {
		// shares to amount logic comes from here:
		// https://github.com/cosmos/cosmos-sdk/blob/v0.38.3/x/staking/keeper/querier.go#L404
		val, found := keeper.GetValidator(ctx, d.ValidatorAddress)
		if !found {
			return nil, staking.ErrNoValidatorFound(keeper.Codespace())
		}
		amount := sdk.NewCoin(bondDenom, val.TokensFromShares(d.Shares).TruncateInt())

		// Accumulated Rewards???

		// can relegate? other query for redelegations?
		// keeper.GetRedelegation

		result[i] = wasmTypes.Delegation{
			Delegator: d.DelegatorAddress.String(),
			Validator: d.ValidatorAddress.String(),
			Amount:    convertSdkCoinToWasmCoin(amount),
		}
	}
	return result, nil
}

func sdkToFullDelegation(ctx sdk.Context, keeper staking.Keeper, delegation staking.Delegation) (*wasmTypes.FullDelegation, error) {
	val, found := keeper.GetValidator(ctx, delegation.ValidatorAddress)
	if !found {
		return nil, staking.ErrNoValidatorFound(keeper.Codespace())
	}
	bondDenom := keeper.BondDenom(ctx)
	amount := sdk.NewCoin(bondDenom, val.TokensFromShares(delegation.Shares).TruncateInt())

	// can relegate? other query for redelegations?
	// keeper.GetRedelegation

	return &wasmTypes.FullDelegation{
		Delegator: delegation.DelegatorAddress.String(),
		Validator: delegation.ValidatorAddress.String(),
		Amount:    convertSdkCoinToWasmCoin(amount),
		// TODO: AccumulatedRewards
		AccumulatedRewards: wasmTypes.NewCoin(0, bondDenom),
		// TODO: Determine redelegate
		CanRedelegate: wasmTypes.NewCoin(0, bondDenom),
	}, nil
}

func WasmQuerier(wasm Keeper) func(ctx sdk.Context, request *wasmTypes.WasmQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmTypes.WasmQuery) ([]byte, error) {
		if request.Smart != nil {
			addr, err := sdk.AccAddressFromBech32(request.Smart.ContractAddr)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(request.Smart.ContractAddr)
			}
			return wasm.QuerySmart(ctx, addr, request.Smart.Msg)
		}
		if request.Raw != nil {
			addr, err := sdk.AccAddressFromBech32(request.Raw.ContractAddr)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(request.Raw.ContractAddr)
			}
			models := wasm.QueryRaw(ctx, addr, request.Raw.Key)
			// TODO: do we want to change the return value?
			return json.Marshal(models)
		}
		return nil, wasmTypes.UnsupportedRequest{"unknown WasmQuery variant"}
	}
}

func convertSdkCoinsToWasmCoins(coins []sdk.Coin) wasmTypes.Coins {
	converted := make(wasmTypes.Coins, len(coins))
	for i, c := range coins {
		converted[i] = convertSdkCoinToWasmCoin(c)
	}
	return converted
}

func convertSdkCoinToWasmCoin(coin sdk.Coin) wasmTypes.Coin {
	return wasmTypes.Coin{
		Denom:  coin.Denom,
		Amount:  types.ConvertSdkDecToString(coin.Amount),
	}
}

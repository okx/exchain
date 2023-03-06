package keeper

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/query"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/feesplit/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		if len(path) < 1 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
				"Insufficient parameters, at least 1 parameter is required")
		}

		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, keeper)
		case types.QueryFeeSplits:
			return queryFeeSplits(ctx, req, keeper)
		case types.QueryFeeSplit:
			return queryFeeSplit(ctx, req, keeper)
		case types.QueryDeployerFeeSplits:
			return queryDeployerFeeSplits(ctx, req, keeper)
		case types.QueryDeployerFeeSplitsDetail:
			return queryDeployerFeeSplitsDetail(ctx, req, keeper)
		case types.QueryWithdrawerFeeSplits:
			return queryWithdrawerFeeSplits(ctx, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}

// queryFeeSplits returns all FeeSplits that have been registered for fee distribution
func queryFeeSplits(
	ctx sdk.Context,
	req abci.RequestQuery,
	k Keeper,
) ([]byte, sdk.Error) {
	var params types.QueryWithdrawerFeeSplitsRequest
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var feeSplits []types.FeeSplitWithShare
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixFeeSplit)

	pageRes, err := query.Paginate(store, params.Pagination, func(_, value []byte) error {
		var fee types.FeeSplit
		if err := k.cdc.UnmarshalBinaryBare(value, &fee); err != nil {
			return err
		}
		share, found := k.GetContractShare(ctx, fee.ContractAddress)
		if !found {
			share = k.GetParams(ctx).DeveloperShares
		}
		feeSplits = append(feeSplits, types.FeeSplitWithShare{
			ContractAddress:   fee.ContractAddress.String(),
			DeployerAddress:   fee.DeployerAddress.String(),
			WithdrawerAddress: fee.WithdrawerAddress.String(),
			Share:             share,
		})
		return nil
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}

	resp := &types.QueryFeeSplitsResponse{
		FeeSplits:  feeSplits,
		Pagination: pageRes,
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryFeeSplit returns the FeeSplit that has been registered for fee distribution for a given
// contract
func queryFeeSplit(
	ctx sdk.Context,
	req abci.RequestQuery,
	k Keeper,
) ([]byte, sdk.Error) {
	var params types.QueryFeeSplitRequest
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if strings.TrimSpace(params.ContractAddress) == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "contract address is empty")
	}

	// check if the contract is a non-zero hex address
	if err := types.ValidateNonZeroAddress(params.ContractAddress); err != nil {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid format for contract %s, should be non-zero hex ('0x...')", params.ContractAddress),
		)
	}

	feeSplit, found := k.GetFeeSplit(ctx, common.HexToAddress(params.ContractAddress))
	if !found {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("not found fees registered contract '%s'", params.ContractAddress),
		)
	}
	share, found := k.GetContractShare(ctx, feeSplit.ContractAddress)
	if !found {
		share = k.GetParams(ctx).DeveloperShares
	}

	resp := &types.QueryFeeSplitResponse{FeeSplit: types.FeeSplitWithShare{
		ContractAddress:   feeSplit.ContractAddress.String(),
		DeployerAddress:   feeSplit.DeployerAddress.String(),
		WithdrawerAddress: feeSplit.WithdrawerAddress.String(),
		Share:             share,
	}}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryParams returns the fees module params
func queryParams(
	ctx sdk.Context,
	k Keeper,
) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(k.cdc, types.QueryParamsResponse{Params: params})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryDeployerFeeSplits returns all contracts that have been registered for fee
// distribution by a given deployer
func queryDeployerFeeSplits(
	ctx sdk.Context,
	req abci.RequestQuery,
	k Keeper,
) ([]byte, sdk.Error) {
	var params types.QueryDeployerFeeSplitsRequest
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if strings.TrimSpace(params.DeployerAddress) == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "deployer address is empty")
	}

	deployer, err := sdk.AccAddressFromBech32(params.DeployerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid format for deployer %s, should be bech32", params.DeployerAddress),
		)
	}

	var contracts []string
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.GetKeyPrefixDeployer(deployer),
	)

	pageRes, err := query.Paginate(store, params.Pagination, func(key, _ []byte) error {
		contracts = append(contracts, common.BytesToAddress(key).Hex())
		return nil
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}

	resp := &types.QueryDeployerFeeSplitsResponse{
		ContractAddresses: contracts,
		Pagination:        pageRes,
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryDeployerFeeSplitsDetail returns all contracts with feesplit info that have been registered for fee
// distribution by a given deployer
func queryDeployerFeeSplitsDetail(
	ctx sdk.Context,
	req abci.RequestQuery,
	k Keeper,
) ([]byte, sdk.Error) {
	var params types.QueryDeployerFeeSplitsRequest
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if strings.TrimSpace(params.DeployerAddress) == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "deployer address is empty")
	}

	deployer, err := sdk.AccAddressFromBech32(params.DeployerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid format for deployer %s, should be bech32", params.DeployerAddress),
		)
	}

	var contracts []common.Address
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.GetKeyPrefixDeployer(deployer),
	)

	pageRes, err := query.Paginate(store, params.Pagination, func(key, _ []byte) error {
		contracts = append(contracts, common.BytesToAddress(key))
		return nil
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}

	var feeSplits []types.FeeSplitWithShare
	for _, contract := range contracts {
		feeSplit, found := k.GetFeeSplit(ctx, contract)
		if !found {
			continue
		}
		share, found := k.GetContractShare(ctx, feeSplit.ContractAddress)
		if !found {
			share = k.GetParams(ctx).DeveloperShares
		}

		feeSplits = append(feeSplits, types.FeeSplitWithShare{
			ContractAddress:   feeSplit.ContractAddress.String(),
			DeployerAddress:   feeSplit.DeployerAddress.String(),
			WithdrawerAddress: feeSplit.WithdrawerAddress.String(),
			Share:             share,
		})
	}

	if len(feeSplits) == 0 {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("not found fees registered contract for deployer '%s'", params.DeployerAddress),
		)
	}
	resp := &types.QueryDeployerFeeSplitsResponseV2{
		FeeSplits:  feeSplits,
		Pagination: pageRes,
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryWithdrawerFeeSplits returns all fees for a given withdraw address
func queryWithdrawerFeeSplits(
	ctx sdk.Context,
	req abci.RequestQuery,
	k Keeper,
) ([]byte, sdk.Error) {
	var params types.QueryWithdrawerFeeSplitsRequest
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if strings.TrimSpace(params.WithdrawerAddress) == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "withdraw address is empty")
	}

	withdrawer, err := sdk.AccAddressFromBech32(params.WithdrawerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid format for withdraw addr %s, should be bech32", params.WithdrawerAddress),
		)
	}

	var contracts []string
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.GetKeyPrefixWithdrawer(withdrawer),
	)

	pageRes, err := query.Paginate(store, params.Pagination, func(key, _ []byte) error {
		contracts = append(contracts, common.BytesToAddress(key).Hex())

		return nil
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}

	resp := &types.QueryWithdrawerFeeSplitsResponse{
		ContractAddresses: contracts,
		Pagination:        pageRes,
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

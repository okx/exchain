package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/upgrade/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryUpgradeConfig        = "config"
	QueryUpgradeVersion       = "version"
	QueryUpgradeFailedVersion = "failed_version"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryUpgradeConfig:
			return queryUpgradeConfig(ctx, req, keeper)
		case QueryUpgradeVersion:
			return queryUpgradeVersion(ctx, req, keeper)
		case QueryUpgradeFailedVersion:
			return queryUpgradeLastFailedVersion(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown token query endpoint")
		}
	}
}

// nolint: unparam
func queryUpgradeConfig(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	appUpgradeConfig, found := keeper.GetAppUpgradeConfig(ctx)

	if !found {
		return nil, types.NewError(types.DefaultCodespace, types.CodeNoUpgradeConfig, "app upgrade config not found")
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, appUpgradeConfig)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

// nolint: unparam
func queryUpgradeVersion(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	version := keeper.protocolKeeper.GetCurrentVersion(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryVersion(version))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

// nolint: unparam
func queryUpgradeLastFailedVersion(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	version := keeper.protocolKeeper.GetLastFailedVersion(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryVersion(version))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

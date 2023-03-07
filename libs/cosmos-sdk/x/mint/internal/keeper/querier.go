package keeper

import (
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint/internal/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, k)
		case types.QueryTreasures:
			return queryTreasures(ctx, k)
		case types.QueryBlockRewards:
			return queryBlockRewards(ctx, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryTreasures(ctx sdk.Context, k Keeper) ([]byte, error) {
	treasures := k.GetTreasures(ctx)
	res, err := codec.MarshalJSONIndent(k.cdc, treasures)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBlockRewards(ctx sdk.Context, k Keeper) ([]byte, error) {
	minter := k.GetMinterCustom(ctx)

	blockAmount := minter.MintedPerBlock

	res, err := codec.MarshalJSONIndent(k.cdc, types.MinterCustom{
		MintedPerBlock:    blockAmount,
		NextBlockToUpdate: minter.NextBlockToUpdate,
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryInflation(ctx sdk.Context, k Keeper) ([]byte, error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, minter.Inflation)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryAnnualProvisions(ctx sdk.Context, k Keeper) ([]byte, error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, minter.AnnualProvisions)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

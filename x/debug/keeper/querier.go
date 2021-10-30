package keeper

import (
	"fmt"
	"github.com/okex/exchain/x/staking"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/debug/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

// NewDebugger returns query handler for module debug
func NewDebugger(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.DumpStore:
			return dumpStore(ctx, req, keeper)
		case types.SanityCheckShares:
			return sanityCheckShares(ctx, keeper)
		case types.InvariantCheck:
			return invariantCheck(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown common query endpoint")
		}
	}
}

func dumpStore(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	var params types.DumpInfoParams
	err := keeper.GetCDC().UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	keeper.DumpStore(ctx, params.Module)
	return nil, nil
}


func sanityCheckShares(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	stakingKeeper, ok := keeper.stakingKeeper.(staking.Keeper)
	if !ok {
		return nil, sdk.ErrInternal("staking keeper mismatch")
	}
	invariantFunc := staking.DelegatorAddSharesInvariant(stakingKeeper)
	msg, broken := invariantFunc(ctx)
	if broken {
		return nil, sdk.ErrInternal(msg)
	}
	return []byte("sanity check passed"), nil
}

func invariantCheck(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	defer func() {
		if e := recover(); e != nil {
			res, err = []byte(fmt.Sprintf("failed to check ivariant:\n\t%v", e)), nil
		}
	}()

	keeper.crisisKeeper.AssertInvariants(ctx)

	return []byte("invariant check passed"), nil
}

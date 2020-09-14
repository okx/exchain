package keeper

import (
	"strings"

	"github.com/okex/okexchain/x/staking"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/debug/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// NewDebugger returns query handler for module debug
func NewDebugger(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.DumpStore:
			return dumpStore(ctx, req, keeper)
		case types.SetLogLevel:
			return setLogLevel(path[1:])
		case types.SanityCheckShares:
			return sanityCheckShares(ctx, keeper)
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

func setLogLevel(paths []string) ([]byte, sdk.Error) {
	level := strings.Join(paths, "/")

	if err := tmlog.UpdateLogLevel(level); err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("log level set failed", err.Error()))
	}
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

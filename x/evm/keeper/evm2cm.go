package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

func (k *Keeper) IsMatchSysContractAddress(ctx sdk.Context, addr sdk.AccAddress) bool {
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	iaddr, err := csdb.GetSysContractAddress()
	if err != nil {
		return false
	}
	return iaddr.Equals(addr)
}

func querySysContractAddress(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res, err := types.CreateEmptyCommitStateDB(keeper.GeneratePureCSDBParams(), ctx).GetSysContractAddress()
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", err.Error()))
	}
	return res, nil
}

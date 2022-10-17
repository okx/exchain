package keeper

import (
	ethermint "github.com/okex/exchain/app/types"
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

func (k Keeper) IsContractAccount(ctx sdk.Context, addr sdk.AccAddress) bool {
	acct := k.accountKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return false
	}
	ethAcct, ok := acct.(*ethermint.EthAccount)
	if !ok {
		return false
	}
	return ethAcct.IsContract()
}

func querySysContractAddress(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res, err := types.CreateEmptyCommitStateDB(keeper.GeneratePureCSDBParams(), ctx).GetSysContractAddress()
	if err != nil {
		return nil, err
	}
	return res, nil
}

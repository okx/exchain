package farm

import (
	"strconv"

	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"
)

const (
	Limit = 5000 // about 6,000,000 gas
)

func init() {
	destroyPoolHandler = handleMsgRmKeys
}

func handleMsgRmKeys(ctx sdk.Context, k keeper.Keeper, msg types.MsgDestroyPool) (*sdk.Result, error) {
	total := 0
	contract := ethcmn.HexToAddress(msg.PoolName)
	evmKeeper := k.EvmKeeper()
	err := evmKeeper.ForEachStorage(ctx, contract, func(key, value ethcmn.Hash) bool {
		if total >= Limit {
			return true
		}
		evmKeeper.DeleteStateDirectly(ctx, contract, key)
		total++
		return false
	})
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeDestroyPool,
		sdk.NewAttribute("address", msg.PoolName),
		sdk.NewAttribute("count", strconv.Itoa(total)),
	))
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

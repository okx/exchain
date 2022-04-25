package token

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/token/types"
)

type WalletTokenTransfer interface {
	sdk.Msg
	GetFrom() sdk.AccAddress
	GetTo() sdk.AccAddress
	GetAmount() []sdk.DecCoin
}

func handleWalletMsgSend(ctx sdk.Context, keeper Keeper, msg types.MsgSend, logger log.Logger) (*sdk.Result, error) {
	if !keeper.bankKeeper.GetSendEnabled(ctx) {
		return types.ErrSendDisabled().Result()
	}
	err := keeper.SendCoinsFromAccountToAccount(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return types.ErrSendCoinsFromAccountToAccountFailed(err.Error()).Result()
	}

	var name = "handleMsgSend"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<From:%s,To:%s,Amount:%s>\n"+
			"                           result<Owner have enough okts to send a tx>\n",
			ctx.BlockHeight(), name,
			msg.FromAddress, msg.ToAddress, msg.Amount))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName)),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

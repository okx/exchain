package token

import (
	"fmt"

	"github.com/okex/okexchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common/perf"
	"github.com/okex/okexchain/x/common/version"
	"github.com/okex/okexchain/x/token/types"
	"github.com/tendermint/tendermint/libs/log"
)

// NewTokenHandler returns a handler for "token" type messages.
func NewTokenHandler(keeper Keeper, protocolVersion version.ProtocolVersionType) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		//logger := ctx.Logger().With("module", "token")
		// NOTE msg already has validate basic run
		var name string
		var handlerFun func() (*sdk.Result, error)
		logger := ctx.Logger().With("module", "token")
		switch msg := msg.(type) {
		case types.MsgTokenIssue:
			name = "handleMsgTokenIssue"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenIssue(ctx, keeper, msg, logger)
			}

		case types.MsgTokenBurn:
			name = "handleMsgTokenBurn"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenBurn(ctx, keeper, msg, logger)
			}

		case types.MsgTokenMint:
			name = "handleMsgTokenMint"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenMint(ctx, keeper, msg, logger)
			}

		case types.MsgMultiSend:
			name = "handleMsgMultiSend"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgMultiSend(ctx, keeper, msg, logger)
			}

		case types.MsgSend:
			name = "handleMsgSend"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgSend(ctx, keeper, msg, logger)
			}

		case types.MsgTransferOwnership:
			name = "handleMsgTransferOwnership"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTransferOwnership(ctx, keeper, msg, logger)
			}
		case types.MsgConfirmOwnership:
			name = "handleMsgConfirmOwnership"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgConfirmOwnership(ctx, keeper, msg, logger)
			}

		case types.MsgTokenModify:
			name = "handleMsgTokenModify"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenModify(ctx, keeper, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("Unrecognized token Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)
		res, err := handlerFun()
		common.SanityCheckHandler(res, err)
		return res, err
	}
}

func handleMsgTokenIssue(ctx sdk.Context, keeper Keeper, msg types.MsgTokenIssue, logger log.Logger) (*sdk.Result, error) {
	// check upper bound
	totalSupply, err := sdk.NewDecFromStr(msg.TotalSupply)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("invalid total supply(%s)", msg.TotalSupply)).Result()
	}
	if totalSupply.GT(sdk.NewDec(types.TotalSupplyUpperbound)) {
		return sdk.ErrInternal(fmt.Sprintf("total-supply(%s) exceeds the upper limit(%d)",
			msg.TotalSupply, types.TotalSupplyUpperbound)).Result()
	}

	token := types.Token{
		Description:         msg.Description,
		OriginalSymbol:      msg.OriginalSymbol,
		WholeName:           msg.WholeName,
		OriginalTotalSupply: totalSupply,
		Owner:               msg.Owner,
		Mintable:            msg.Mintable,
	}

	// generate a random symbol
	newName, valid := addTokenSuffix(ctx, keeper, msg.OriginalSymbol)
	if !valid {
		return sdk.ErrInvalidCoins(fmt.Sprintf(
			"temporarily failed to generate a unique symbol for %s. Try again.",
			msg.OriginalSymbol)).Result()
	}

	token.Symbol = newName

	coins := sdk.MustParseCoins(token.Symbol, msg.TotalSupply)
	// set supply
	err = keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply mint coins error:%s", err.Error())).Result()
	}

	// send coins to owner
	err = keeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, token.Owner, coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply send coins error:%s", err.Error())).Result()
	}

	// set token info
	keeper.NewToken(ctx, token)

	// deduction fee
	feeDecCoins := keeper.GetParams(ctx).FeeIssue.ToCoins()
	err = keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, token.Owner, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String())).Result()
	}

	var name = "handleMsgTokenIssue"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Description:%s,Symbol:%s,OriginalSymbol:%s,TotalSupply:%s,Owner:%v,Mintable:%v>\n"+
			"                           result<Owner have enough okts to issue %s>\n",
			ctx.BlockHeight(), name,
			msg.Description, msg.Symbol, msg.OriginalSymbol, msg.TotalSupply, msg.Owner, msg.Mintable,
			token.Symbol))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyFee, keeper.GetParams(ctx).FeeIssue.String()),
			sdk.NewAttribute("symbol", token.Symbol),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTokenBurn(ctx sdk.Context, keeper Keeper, msg types.MsgTokenBurn, logger log.Logger) (*sdk.Result, error) {

	token := keeper.GetTokenInfo(ctx, msg.Amount.Denom)

	// check owner
	if !token.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized("Not the token's owner").Result()
	}

	subCoins := msg.Amount.ToCoins()
	// send coins to moduleAcc
	err := keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, types.ModuleName, subCoins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply send coins error:%s", err.Error())).Result()
	}

	// set supply
	err = keeper.supplyKeeper.BurnCoins(ctx, types.ModuleName, subCoins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply burn coins error:%s", err.Error())).Result()
	}

	// deduction fee
	feeDecCoins := keeper.GetParams(ctx).FeeBurn.ToCoins()
	err = keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String())).Result()
	}

	var name = "handleMsgTokenBurn"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Owner:%s,Symbol:%s,Amount:%s>\n"+
			"                           result<Owner have enough okts to burn %s>\n",
			ctx.BlockHeight(), name,
			msg.Owner, msg.Amount.Denom, msg.Amount,
			msg.Amount.Denom))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyFee, feeDecCoins.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTokenMint(ctx sdk.Context, keeper Keeper, msg types.MsgTokenMint, logger log.Logger) (*sdk.Result, error) {
	token := keeper.GetTokenInfo(ctx, msg.Amount.Denom)
	// check owner
	if !token.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not the owner of token(%s)",
			msg.Owner.String(), msg.Amount.Denom)).Result()
	}

	// check whether token is mintable
	if !token.Mintable {
		return sdk.ErrUnauthorized(fmt.Sprintf("token(%s) is not mintable", token.Symbol)).Result()
	}

	// check upper bound
	totalSupplyAfterMint := keeper.supplyKeeper.GetSupplyByDenom(ctx, msg.Amount.Denom).Add(msg.Amount.Amount)
	if totalSupplyAfterMint.GT(sdk.NewDec(types.TotalSupplyUpperbound)) {
		return sdk.ErrInternal(fmt.Sprintf("total-supply(%s) exceeds the upper limit(%d)",
			totalSupplyAfterMint, types.TotalSupplyUpperbound)).Result()
	}

	mintCoins := msg.Amount.ToCoins()
	// set supply
	err := keeper.supplyKeeper.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply mint coins error:%s", err.Error())).Result()
	}

	// send coins to acc
	err = keeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.Owner, mintCoins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supply send coins error:%s", err.Error())).Result()
	}

	// deduction fee
	feeDecCoins := keeper.GetParams(ctx).FeeMint.ToCoins()
	err = keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String())).Result()
	}

	name := "handleMsgTokenMint"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Owner:%s,Symbol:%s,Amount:%s>\n"+
			"                           result<Owner have enough okts to Mint %s>\n",
			ctx.BlockHeight(), name,
			msg.Owner, msg.Amount.Denom, msg.Amount,
			msg.Amount.Denom))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyFee, feeDecCoins.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgMultiSend(ctx sdk.Context, keeper Keeper, msg types.MsgMultiSend, logger log.Logger) (*sdk.Result, error) {
	if !keeper.bankKeeper.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(DefaultCodespace).Result()
	}

	var transfers string
	var coinNum int
	for _, transferUnit := range msg.Transfers {
		coinNum += len(transferUnit.Coins)
		err := keeper.SendCoinsFromAccountToAccount(ctx, msg.From, transferUnit.To, transferUnit.Coins)
		if err != nil {
			return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient coins(need %s)",
				transferUnit.Coins.String())).Result()
		}
		transfers += fmt.Sprintf("                          msg<To:%s,Coin:%s>\n", transferUnit.To, transferUnit.Coins)
	}

	name := "handleMsgMultiSend"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<From:%s>\n"+
			transfers+
			"                           result<Owner have enough okts to send multi txs>\n",
			ctx.BlockHeight(), name,
			msg.From))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName)),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSend(ctx sdk.Context, keeper Keeper, msg types.MsgSend, logger log.Logger) (*sdk.Result, error) {
	if !keeper.bankKeeper.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(DefaultCodespace).Result()
	}

	err := keeper.SendCoinsFromAccountToAccount(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient coins(need %s)",
			msg.Amount.String())).Result()
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

func handleMsgTransferOwnership(ctx sdk.Context, keeper Keeper, msg types.MsgTransferOwnership, logger log.Logger) (*sdk.Result, error) {
	tokenInfo := keeper.GetTokenInfo(ctx, msg.Symbol)

	if !tokenInfo.Owner.Equals(msg.FromAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not the owner of token(%s)",
			msg.FromAddress.String(), msg.Symbol)).Result()
	}

	confirmOwnership, exist := keeper.GetConfirmOwnership(ctx, msg.Symbol)
	if exist && !ctx.BlockTime().After(confirmOwnership.Expire) {
		return sdk.ErrInternal(fmt.Sprintf("repeated transfer-ownership of token(%s) is not allowed", msg.Symbol)).Result()
	}

	if msg.ToAddress.Equals(common.BlackHoleAddress()) { // transfer ownership to black hole
		// first remove it from the raw owner
		keeper.DeleteUserToken(ctx, tokenInfo.Owner, tokenInfo.Symbol)
		tokenInfo.Owner = msg.ToAddress
		keeper.NewToken(ctx, tokenInfo)
	} else {
		// set confirm ownership info
		expireTime := ctx.BlockTime().Add(keeper.GetParams(ctx).OwnershipConfirmWindow)
		confirmOwnership = &types.ConfirmOwnership{
			Symbol:  msg.Symbol,
			Address: msg.ToAddress,
			Expire:  expireTime,
		}
		keeper.SetConfirmOwnership(ctx, confirmOwnership)
	}
	// deduction fee
	feeDecCoins := keeper.GetParams(ctx).FeeChown.ToCoins()
	err := keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FromAddress, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String())).Result()
	}

	var name = "handleMsgTransferOwnership"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<From:%s,To:%s,Symbol:%s>\n"+
			"                           result<Owner have enough okts to transfer the %s>\n",
			ctx.BlockHeight(), name,
			msg.FromAddress, msg.ToAddress, msg.Symbol,
			msg.Symbol))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyFee, keeper.GetParams(ctx).FeeChown.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgConfirmOwnership(ctx sdk.Context, keeper Keeper, msg types.MsgConfirmOwnership, logger log.Logger) (*sdk.Result, error) {
	confirmOwnership, exist := keeper.GetConfirmOwnership(ctx, msg.Symbol)
	if !exist {
		return sdk.ErrUnknownRequest(fmt.Sprintf("no transfer-ownership of token (%s) to confirm",
			msg.Address.String())).Result()
	}
	if ctx.BlockTime().After(confirmOwnership.Expire) {
		// delete ownership confirming information
		keeper.DeleteConfirmOwnership(ctx, confirmOwnership.Symbol)
		return sdk.ErrInternal(fmt.Sprintf("transfer-ownership is expired, expire time (%s)", confirmOwnership.Expire.String())).Result()
	}
	if !confirmOwnership.Address.Equals(msg.Address) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is expected as the new owner",
			confirmOwnership.Address.String())).Result()
	}

	tokenInfo := keeper.GetTokenInfo(ctx, msg.Symbol)
	// first remove it from the raw owner
	keeper.DeleteUserToken(ctx, tokenInfo.Owner, tokenInfo.Symbol)
	tokenInfo.Owner = msg.Address
	keeper.NewToken(ctx, tokenInfo)

	// delete ownership confirming information
	keeper.DeleteConfirmOwnership(ctx, confirmOwnership.Symbol)

	var name = "handleMsgConfirmOwnership"
	logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
		"                           msg<From:%s,Symbol:%s>\n"+
		"                           result<Owner have enough okts to transfer the %s>\n",
		ctx.BlockHeight(), name, msg.Address, msg.Symbol, msg.Symbol))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyFee, keeper.GetParams(ctx).FeeChown.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTokenModify(ctx sdk.Context, keeper Keeper, msg types.MsgTokenModify, logger log.Logger) (*sdk.Result, error) {
	token := keeper.GetTokenInfo(ctx, msg.Symbol)
	// check owner
	if !token.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not the owner of token(%s)",
			msg.Owner.String(), msg.Symbol)).Result()
	}
	if !msg.IsWholeNameModified && !msg.IsDescriptionModified {
		return sdk.ErrInternal("nothing modified").Result()
	}
	// modify
	if msg.IsWholeNameModified {
		token.WholeName = msg.WholeName
	}
	if msg.IsDescriptionModified {
		token.Description = msg.Description
	}

	keeper.UpdateToken(ctx, token)

	// deduction fee
	feeDecCoins := keeper.GetParams(ctx).FeeModify.ToCoins()
	err := keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String())).Result()
	}

	name := "handleMsgTokenModify"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Owner:%s,Symbol:%s,WholeName:%s,Description:%s>\n"+
			"                           result<Owner have enough okts to edit %s>\n",
			ctx.BlockHeight(), name,
			msg.Owner, msg.Symbol, msg.WholeName, msg.Description,
			msg.Symbol))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyFee, keeper.GetParams(ctx).FeeModify.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

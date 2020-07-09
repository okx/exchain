package dex

import (
	"fmt"
	"strconv"

	"github.com/okex/okchain/x/common/perf"
	"github.com/okex/okchain/x/dex/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
)

// NewHandler handles all "dex" type messages.
func NewHandler(k IKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", ModuleName)

		var handlerFun func() sdk.Result
		var name string
		switch msg := msg.(type) {
		case MsgList:
			name = "handleMsgList"
			handlerFun = func() sdk.Result {
				return handleMsgList(ctx, k, msg, logger)
			}
		case MsgDeposit:
			name = "handleMsgDeposit"
			handlerFun = func() sdk.Result {
				return handleMsgDeposit(ctx, k, msg, logger)
			}
		case MsgWithdraw:
			name = "handleMsgWithDraw"
			handlerFun = func() sdk.Result {
				return handleMsgWithDraw(ctx, k, msg, logger)
			}
		case MsgTransferOwnership:
			name = "handleMsgTransferOwnership"
			handlerFun = func() sdk.Result {
				return handleMsgTransferOwnership(ctx, k, msg, logger)
			}
		case MsgCreateOperator:
			name = "handleMsgCreateOperator"
			handlerFun = func() sdk.Result {
				return handleMsgCreateOperator(ctx, k, msg, logger)
			}
		case MsgUpdateOperator:
			name = "handleMsgUpdateOperator"
			handlerFun = func() sdk.Result {
				return handleMsgUpdateOperator(ctx, k, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("unrecognized dex message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, ModuleName, name, seq)
		return handlerFun()
	}
}

func handleMsgList(ctx sdk.Context, keeper IKeeper, msg MsgList, logger log.Logger) sdk.Result {

	if !keeper.GetTokenKeeper().TokenExist(ctx, msg.ListAsset) ||
		!keeper.GetTokenKeeper().TokenExist(ctx, msg.QuoteAsset) {
		return sdk.ErrInvalidCoins(
			fmt.Sprintf("%s or %s is not valid", msg.ListAsset, msg.QuoteAsset)).Result()
	}

	if _, exists := keeper.GetOperator(ctx, msg.Owner); !exists {
		return types.ErrUnknownOperator(msg.Owner).Result()
	}

	tokenPair := &TokenPair{
		BaseAssetSymbol:  msg.ListAsset,
		QuoteAssetSymbol: msg.QuoteAsset,
		InitPrice:        msg.InitPrice,
		MaxPriceDigit:    int64(DefaultMaxPriceDigitSize),
		MaxQuantityDigit: int64(DefaultMaxQuantityDigitSize),
		MinQuantity:      sdk.MustNewDecFromStr("0.00000001"),
		Owner:            msg.Owner,
		Delisting:        false,
		Deposits:         DefaultTokenPairDeposit,
		BlockHeight:      ctx.BlockHeight(),
	}

	// check whether a specific token pair exists with the symbols of base asset and quote asset
	// Note: aaa_bbb and bbb_aaa are actually one token pair
	if keeper.GetTokenPair(ctx, fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)) != nil ||
		keeper.GetTokenPair(ctx, fmt.Sprintf("%s_%s", tokenPair.QuoteAssetSymbol, tokenPair.BaseAssetSymbol)) != nil {
		return types.ErrTokenPairExisted(tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol).Result()
	}

	// deduction fee
	feeCoins := keeper.GetParams(ctx).ListFee.ToCoins()
	err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.GetFeeCollector(), feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeCoins.String())).Result()
	}

	err2 := keeper.SaveTokenPair(ctx, tokenPair)
	if err2 != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to SaveTokenPair: %s", err2.Error())).Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgList: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute("list-asset", tokenPair.BaseAssetSymbol),
			sdk.NewAttribute("quote-asset", tokenPair.QuoteAssetSymbol),
			sdk.NewAttribute("init-price", tokenPair.InitPrice.String()),
			sdk.NewAttribute("max-price-digit", strconv.FormatInt(tokenPair.MaxPriceDigit, 10)),
			sdk.NewAttribute("max-size-digit", strconv.FormatInt(tokenPair.MaxQuantityDigit, 10)),
			sdk.NewAttribute("min-trade-size", tokenPair.MinQuantity.String()),
			sdk.NewAttribute("delisting", fmt.Sprintf("%t", tokenPair.Delisting)),
			sdk.NewAttribute(sdk.AttributeKeyFee, feeCoins.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDeposit(ctx sdk.Context, keeper IKeeper, msg MsgDeposit, logger log.Logger) sdk.Result {
	if sdkErr := keeper.Deposit(ctx, msg.Product, msg.Depositor, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgDeposit: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}

}

func handleMsgWithDraw(ctx sdk.Context, keeper IKeeper, msg MsgWithdraw, logger log.Logger) sdk.Result {
	if sdkErr := keeper.Withdraw(ctx, msg.Product, msg.Depositor, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgWithDraw: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransferOwnership(ctx sdk.Context, keeper IKeeper, msg MsgTransferOwnership,
	logger log.Logger) sdk.Result {

	if _, exist := keeper.GetOperator(ctx, msg.FromAddress); !exist {
		return types.ErrUnknownOperator(msg.FromAddress).Result()
	}

	if _, exist := keeper.GetOperator(ctx, msg.ToAddress); !exist {
		return types.ErrUnknownOperator(msg.ToAddress).Result()
	}

	if sdkErr := keeper.TransferOwnership(ctx, msg.Product, msg.FromAddress, msg.ToAddress); sdkErr != nil {
		return sdkErr.Result()
	}

	// deduction fee
	feeCoins := keeper.GetParams(ctx).TransferOwnershipFee.ToCoins()
	err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.FromAddress, keeper.GetFeeCollector(), feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeCoins.String())).Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgTransferOwnership: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyFee, feeCoins.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgCreateOperator(ctx sdk.Context, keeper IKeeper, msg MsgCreateOperator, logger log.Logger) sdk.Result {

	logger.Debug(fmt.Sprintf("handleMsgCreateOperator msg: %+v", msg))

	if _, isExist := keeper.GetOperator(ctx, msg.Owner); isExist {
		return types.ErrExistOperator(msg.Owner).Result()
	}
	operator := types.DEXOperator{
		Address:            msg.Owner,
		HandlingFeeAddress: msg.HandlingFeeAddress,
		Website:            msg.Website,
		InitHeight:         ctx.BlockHeight(),
		TxHash:             fmt.Sprintf("%X", tmhash.Sum(ctx.TxBytes())),
	}
	keeper.SetOperator(ctx, operator)

	// deduction fee
	feeCoins := keeper.GetParams(ctx).RegisterOperatorFee.ToCoins()
	err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.GetFeeCollector(), feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeCoins.String())).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyFee, feeCoins.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgUpdateOperator(ctx sdk.Context, keeper IKeeper, msg MsgUpdateOperator, logger log.Logger) sdk.Result {

	logger.Debug(fmt.Sprintf("handleMsgUpdateOperator msg: %+v", msg))

	operator, isExist := keeper.GetOperator(ctx, msg.Owner)
	if !isExist {
		return types.ErrUnknownOperator(msg.Owner).Result()
	}
	if !operator.Address.Equals(msg.Owner) {
		return sdk.ErrUnauthorized("Not the operator's owner").Result()
	}

	operator.HandlingFeeAddress = msg.HandlingFeeAddress
	operator.Website = msg.Website

	keeper.SetOperator(ctx, operator)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

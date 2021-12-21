/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/21 5:08 上午
# @File : handler.go
# @Description :
# @Attention :
*/
package nameservice

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/nameservice/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var handlerFun func() (*sdk.Result, error)
		var name string

		switch msg := msg.(type) {
		case types.MsgSetName:
			name = "handleMsgSetName"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgSetName(ctx, k, msg)
			}
		case types.MsgBuyName:
			name = "handleMsgBuyName"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgBuyName(ctx, k, msg)
			}
		case types.MsgDeleteName:
			name = "handleMsgDeleteName"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgDeleteName(ctx, k, msg)
			}
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)

		res, err := handlerFun()
		common.SanityCheckHandler(res, err)
		return res, err
	}
}

// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg types.MsgSetName) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetCreator(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return &sdk.Result{}, nil                // return
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, k Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	// Checks if the the bid price is greater than the price paid by the current owner
	if k.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
	}
	if k.HasCreator(ctx, msg.Name) {
		err := k.CoinKeeper.SendCoins(ctx, msg.Buyer, k.GetCreator(ctx, msg.Name), msg.Bid)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO
		ctx.Logger().Info("bad")
		// _, err := k.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		// if err != nil {
		// 	return nil, err
		// }
	}
	k.SetCreator(ctx, msg.Name, msg.Buyer)
	k.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}

// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, k Keeper, msg types.MsgDeleteName) (*sdk.Result, error) {
	if !k.WhoisExists(ctx, msg.ID) {
		// replace with ErrKeyNotFound for 0.39+
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.ID)
	}
	if !msg.Creator.Equals(k.GetWhoisOwner(ctx, msg.ID)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	k.DeleteWhois(ctx, msg.ID)
	return &sdk.Result{}, nil
}

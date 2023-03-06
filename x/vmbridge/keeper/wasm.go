package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	ibcadapter "github.com/okx/okbchain/libs/cosmos-sdk/types/ibc-adapter"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
	"github.com/okx/okbchain/x/vmbridge/types"
	"github.com/okx/okbchain/x/wasm"
)

func (k Keeper) SendToWasm(ctx sdk.Context, caller sdk.AccAddress, wasmContractAddr, recipient string, amount sdk.Int) error {
	// must check recipient is ex address
	if !sdk.IsOKCAddress(recipient) {
		return types.ErrIsNotOKCAddr
	}
	to, err := sdk.AccAddressFromBech32(recipient)
	if err != nil {
		return err
	}

	if amount.IsNegative() {
		return types.ErrAmountNegative
	}
	input, err := types.GetMintCW20Input(amount.String(), to.String())
	if err != nil {
		return err
	}
	contractAddr, err := sdk.AccAddressFromBech32(wasmContractAddr)
	if err != nil {
		return err
	}
	if !sdk.IsWasmAddress(contractAddr) {
		return types.ErrIsNotWasmAddr
	}

	ret, err := k.wasmKeeper.Execute(ctx, contractAddr, caller, input, sdk.Coins{})
	if err != nil {
		k.Logger().Error("wasm return", string(ret))
	}
	return err
}

// RegisterSendToEvmEncoder needs to be registered in app setup to handle custom message callbacks
func RegisterSendToEvmEncoder(cdc *codec.ProtoCodec) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: sendToEvmEncoder(cdc),
	}
}

func sendToEvmEncoder(cdc *codec.ProtoCodec) wasm.CustomEncoder {
	return func(sender sdk.AccAddress, data json.RawMessage) ([]ibcadapter.Msg, error) {
		var msg types.MsgSendToEvm

		if err := cdc.UnmarshalJSON(data, &msg); err != nil {
			return nil, err
		}
		return []ibcadapter.Msg{&msg}, nil
	}
}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SendToEvmEvent(goCtx context.Context, msg *types.MsgSendToEvm) (*types.MsgSendToEvmResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !tmtypes.HigherThanEarth(ctx.BlockHeight()) {
		errMsg := fmt.Sprintf("vmbridger not supprt at height %d", ctx.BlockHeight())
		return &types.MsgSendToEvmResponse{Success: false}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
	params := k.wasmKeeper.GetParams(ctx)
	if !params.VmbridgeEnable {
		return &types.MsgSendToEvmResponse{Success: false}, types.ErrVMBridgeEnable
	}

	success, err := k.Keeper.SendToEvm(ctx, msg.Sender, msg.Contract, msg.Recipient, msg.Amount)
	if err != nil {
		return &types.MsgSendToEvmResponse{Success: false}, sdkerrors.Wrap(types.ErrEvmExecuteFailed, err.Error())
	}
	response := types.MsgSendToEvmResponse{Success: success}
	return &response, nil
}

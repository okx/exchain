package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	ibcadapter "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/vmbridge/types"
	"github.com/okex/exchain/x/wasm"
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
	input, err := types.GetMintCW20Input(amount.String(), sdk.AccToAWasmddress(to).String())
	if err != nil {
		return err
	}
	contractAddr, err := sdk.WasmAddressFromBech32(wasmContractAddr)
	if err != nil {
		return err
	}

	ret, err := k.wasmKeeper.Execute(ctx, contractAddr, sdk.AccToAWasmddress(caller), input, sdk.Coins{})
	if err != nil {
		k.Logger().Error("wasm return", string(ret))
	}
	return err
}

func (k Keeper) CallToWasm(ctx sdk.Context, caller sdk.AccAddress, wasmContractAddr string, value sdk.Int, calldata string) ([]byte, error) {
	if value.IsNegative() {
		return nil, types.ErrAmountNegative
	}

	contractAddr, err := sdk.WasmAddressFromBech32(wasmContractAddr)
	if err != nil {
		return nil, err
	}
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecFromBigIntWithPrec(value.BigInt(), sdk.Precision)))
	ret, err := k.wasmKeeper.Execute(ctx, contractAddr, sdk.AccToAWasmddress(caller), []byte(calldata), coins)
	if err != nil {
		k.Logger().Error("wasm return", string(ret))
	}
	return ret, err
}

func (k Keeper) QueryToWasm(ctx sdk.Context, wasmContractAddr string, calldata []byte) ([]byte, error) {
	contractAddr, err := sdk.WasmAddressFromBech32(wasmContractAddr)
	if err != nil {
		return nil, err
	}
	request, err := types.GetWasmVMQueryRequest(calldata)
	if err != nil {
		return nil, err
	}
	gaslimit := k.wasmKeeper.RuntimeGasForContract(ctx)
	queryHandler := k.wasmKeeper.NewQueryHandler(ctx, contractAddr)
	return queryHandler.Query(*request, gaslimit)
}

// RegisterSendToEvmEncoder needs to be registered in app setup to handle custom message callbacks
func RegisterSendToEvmEncoder(cdc *codec.ProtoCodec) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: sendToEvmEncoder(cdc),
	}
}

func sendToEvmEncoder(cdc *codec.ProtoCodec) wasm.CustomEncoder {
	return func(sender sdk.WasmAddress, data json.RawMessage) ([]ibcadapter.Msg, error) {
		var sendToEvmMsg types.MsgSendToEvm
		if err := cdc.UnmarshalJSON(data, &sendToEvmMsg); err != nil {
			var callToEvmMsg types.MsgCallToEvm
			if err := cdc.UnmarshalJSON(data, &callToEvmMsg); err != nil {
				return nil, err
			}
			return []ibcadapter.Msg{&callToEvmMsg}, nil
		}

		return []ibcadapter.Msg{&sendToEvmMsg}, nil
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

func (k msgServer) CallToEvmEvent(goCtx context.Context, msg *types.MsgCallToEvm) (*types.MsgCallToEvmResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !tmtypes.HigherThanEarth(ctx.BlockHeight()) {
		errMsg := fmt.Sprintf("vmbridger not supprt at height %d", ctx.BlockHeight())
		return &types.MsgCallToEvmResponse{Response: errMsg}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
	params := k.wasmKeeper.GetParams(ctx)
	if !params.VmbridgeEnable {
		return &types.MsgCallToEvmResponse{Response: types.ErrVMBridgeEnable.Error()}, types.ErrVMBridgeEnable
	}

	result, err := k.Keeper.CallToEvm(ctx, msg.Sender, msg.Evmaddr, msg.Calldata, msg.Value)
	if err != nil {
		return &types.MsgCallToEvmResponse{Response: sdkerrors.Wrap(types.ErrEvmExecuteFailed, err.Error()).Error()}, sdkerrors.Wrap(types.ErrEvmExecuteFailed, err.Error())
	}
	response := types.MsgCallToEvmResponse{Response: result}
	return &response, nil
}

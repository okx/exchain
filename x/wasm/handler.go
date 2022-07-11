package wasm

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/tendermint/libs/kv"
	types2 "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/wasm/keeper"
	"github.com/okex/exchain/x/wasm/types"
)

// NewHandler returns a handler for "wasm" type messages.
func NewHandler(k types.ContractOpsKeeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx.SetEventManager(sdk.NewEventManager())

		if !types2.HigherThanVenus2(ctx.BlockHeight()) {
			errMsg := fmt.Sprintf("wasm not supprt at height %d", ctx.BlockHeight())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		var (
			res proto.Message
			err error
		)
		switch msg := msg.(type) {
		case *MsgStoreCode: //nolint:typecheck
			res, err = msgServer.StoreCode(sdk.WrapSDKContext(ctx), msg)
		case *MsgInstantiateContract:
			res, err = msgServer.InstantiateContract(sdk.WrapSDKContext(ctx), msg)
		case *MsgExecuteContract:
			res, err = msgServer.ExecuteContract(sdk.WrapSDKContext(ctx), msg)
		case *MsgMigrateContract:
			res, err = msgServer.MigrateContract(sdk.WrapSDKContext(ctx), msg)
		case *MsgUpdateAdmin:
			res, err = msgServer.UpdateAdmin(sdk.WrapSDKContext(ctx), msg)
		case *MsgClearAdmin:
			res, err = msgServer.ClearAdmin(sdk.WrapSDKContext(ctx), msg)
		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		ctx.SetEventManager(filterMessageEvents(ctx))
		return sdk.WrapServiceResult(ctx, res, err)
	}
}

// filterMessageEvents returns the same events with all of type == EventTypeMessage removed except
// for wasm message types.
// this is so only our top-level message event comes through
func filterMessageEvents(ctx sdk.Context) *sdk.EventManager {
	m := sdk.NewEventManager()
	for _, e := range ctx.EventManager().Events() {
		if e.Type == sdk.EventTypeMessage &&
			!hasWasmModuleAttribute(e.Attributes) {
			continue
		}
		m.EmitEvent(e)
	}
	return m
}

func hasWasmModuleAttribute(attrs []kv.Pair) bool {
	for _, a := range attrs {
		if sdk.AttributeKeyModule == string(a.Key) &&
			types.ModuleName == string(a.Value) {
			return true
		}
	}
	return false
}

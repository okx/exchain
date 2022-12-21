package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
)

// EmitConnectionOpenTryEvent emits a connection open try event
func EmitConnectionOpenTryEvent(ctx sdk.Context, connectionID string, clientID string, counterparty types.Counterparty) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeConnectionOpenTry,
			sdk.NewAttribute(types.AttributeKeyConnectionID, connectionID),
			sdk.NewAttribute(types.AttributeKeyClientID, clientID),
			sdk.NewAttribute(types.AttributeKeyCounterpartyClientID, counterparty.ClientId),
			sdk.NewAttribute(types.AttributeKeyCounterpartyConnectionID, counterparty.ConnectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
}

// EmitConnectionOpenAckEvent emits a connection open acknowledge event
func EmitConnectionOpenAckEvent(ctx sdk.Context, connectionID string, connectionEnd types.ConnectionEnd) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeConnectionOpenAck,
			sdk.NewAttribute(types.AttributeKeyConnectionID, connectionID),
			sdk.NewAttribute(types.AttributeKeyClientID, connectionEnd.ClientId),
			sdk.NewAttribute(types.AttributeKeyCounterpartyClientID, connectionEnd.Counterparty.ClientId),
			sdk.NewAttribute(types.AttributeKeyCounterpartyConnectionID, connectionEnd.Counterparty.ConnectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
}

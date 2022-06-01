package wasmtesting

import (
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// MockMessageRouter mock for testing
type MockMessageRouter struct {
	HandlerFn func(msg sdk.Msg) baseapp.MsgServiceHandler
}

// Handler is the entry point
func (m MockMessageRouter) Handler(msg sdk.Msg) baseapp.MsgServiceHandler {
	if m.HandlerFn == nil {
		panic("not expected to be called")
	}
	return m.HandlerFn(msg)
}

// MessageRouterFunc convenient type to match the keeper.MessageRouter interface
type MessageRouterFunc func(methodName string) baseapp.MsgServiceHandler

// Handler is the entry point
func (m MessageRouterFunc) Handler(methodName string) baseapp.MsgServiceHandler {
	return m(methodName)
}

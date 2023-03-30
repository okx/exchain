package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgSendToEvm_GetSigners(t *testing.T) {
	testCases := []struct {
		name   string
		msg    MsgSendToEvm
		isErr  bool
		expect []sdk.AccAddress
	}{
		{
			name:   "normal",
			msg:    MsgSendToEvm{Sender: sdk.AccAddress{0x1}.String()},
			expect: []sdk.AccAddress{sdk.AccAddress{0x1}},
		},
		{
			name:  "sender is empty",
			msg:   MsgSendToEvm{},
			isErr: true,
		},
		{
			name:  "sender is error",
			msg:   MsgSendToEvm{Sender: "0x1111"},
			isErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {

			defer func() {
				r := recover()
				if tc.isErr {
					require.NotNil(t, r)
					require.Error(tt, r.(error))
				}
			}()
			result := tc.msg.GetSigners()
			require.Equal(tt, tc.expect, result)
		})
	}
}

func TestMsgSendToEvm_GetSignBytes(t *testing.T) {
	testCases := []struct {
		name  string
		msg   MsgSendToEvm
		isErr bool
	}{
		{
			name:  "normal",
			msg:   MsgSendToEvm{Sender: sdk.AccAddress{0x1}.String(), Contract: sdk.AccAddress{0x2}.String(), Recipient: sdk.AccAddress{0x3}.String(), Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "sender is empty",
			msg:   MsgSendToEvm{Contract: sdk.AccAddress{0x2}.String(), Recipient: sdk.AccAddress{0x3}.String(), Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "sender is error",
			msg:   MsgSendToEvm{Sender: "ex111", Contract: sdk.AccAddress{0x2}.String(), Recipient: sdk.AccAddress{0x3}.String(), Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "contract is error",
			msg:   MsgSendToEvm{Sender: sdk.AccAddress{0x1}.String(), Contract: "ex111", Recipient: sdk.AccAddress{0x3}.String(), Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "recipient is error",
			msg:   MsgSendToEvm{Sender: sdk.AccAddress{0x1}.String(), Contract: sdk.AccAddress{0x2}.String(), Recipient: "ex111", Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "amount is negative",
			msg:   MsgSendToEvm{Sender: sdk.AccAddress{0x1}.String(), Contract: sdk.AccAddress{0x2}.String(), Recipient: sdk.AccAddress{0x3}.String(), Amount: sdk.NewInt(-1)},
			isErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {

			defer func() {
				r := recover()
				if tc.isErr {
					require.NotNil(t, r)
					require.Error(tt, r.(error))
				}
			}()
			tc.msg.GetSignBytes()
		})
	}
}

func TestMsgSendToEvm_ValidateBasic(t *testing.T) {
	wasmaAddr := sdk.AccAddress(make([]byte, 64)).String()
	addr := sdk.AccAddress{0x1}.String()
	errAddr := "error addr"
	testCases := []struct {
		name  string
		msg   MsgSendToEvm
		isErr bool
	}{
		{
			name:  "normal",
			msg:   MsgSendToEvm{Sender: wasmaAddr, Contract: addr, Recipient: addr, Amount: sdk.NewInt(1)},
			isErr: false,
		},
		{
			name:  "sender is empty",
			msg:   MsgSendToEvm{Contract: addr, Recipient: addr, Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "sender is error",
			msg:   MsgSendToEvm{Sender: errAddr, Contract: addr, Recipient: addr, Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "sender is not wasm addr",
			msg:   MsgSendToEvm{Sender: addr, Contract: addr, Recipient: addr, Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "contract is error",
			msg:   MsgSendToEvm{Sender: wasmaAddr, Contract: errAddr, Recipient: addr, Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "recipient is error",
			msg:   MsgSendToEvm{Sender: wasmaAddr, Contract: addr, Recipient: errAddr, Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "recipient is wasm addr",
			msg:   MsgSendToEvm{Sender: wasmaAddr, Contract: addr, Recipient: wasmaAddr, Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "contract is wasm addr ",
			msg:   MsgSendToEvm{Sender: wasmaAddr, Contract: wasmaAddr, Recipient: addr, Amount: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "amount is negative",
			msg:   MsgSendToEvm{Sender: wasmaAddr, Contract: addr, Recipient: addr, Amount: sdk.NewInt(-1)},
			isErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			if err := tc.msg.ValidateBasic(); tc.isErr {
				require.Error(tt, err)
			}
		})
	}
}

func TestMsgCallToEvm_GetSigners(t *testing.T) {
	testCases := []struct {
		name   string
		msg    MsgCallToEvm
		isErr  bool
		expect []sdk.AccAddress
	}{
		{
			name:   "normal",
			msg:    MsgCallToEvm{Sender: sdk.AccAddress{0x1}.String()},
			expect: []sdk.AccAddress{sdk.AccAddress{0x1}},
		},
		{
			name:  "sender is empty",
			msg:   MsgCallToEvm{},
			isErr: true,
		},
		{
			name:  "sender is error",
			msg:   MsgCallToEvm{Sender: "0x1111"},
			isErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {

			defer func() {
				r := recover()
				if tc.isErr {
					require.NotNil(t, r)
					require.Error(tt, r.(error))
				}
			}()
			result := tc.msg.GetSigners()
			require.Equal(tt, tc.expect, result)
		})
	}
}

func TestMsgCallToEvm_GetSignBytes(t *testing.T) {
	testCases := []struct {
		name  string
		msg   MsgCallToEvm
		isErr bool
	}{
		{
			name:  "normal",
			msg:   MsgCallToEvm{Sender: sdk.AccAddress{0x1}.String(), Evmaddr: sdk.AccAddress{0x2}.String(), Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "sender is empty",
			msg:   MsgCallToEvm{Evmaddr: sdk.AccAddress{0x2}.String(), Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "sender is error",
			msg:   MsgCallToEvm{Sender: "ex111", Evmaddr: sdk.AccAddress{0x2}.String(), Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "contract is error",
			msg:   MsgCallToEvm{Sender: sdk.AccAddress{0x1}.String(), Evmaddr: "ex111", Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "Calldata is empty",
			msg:   MsgCallToEvm{Sender: sdk.AccAddress{0x1}.String(), Evmaddr: sdk.AccAddress{0x2}.String(), Value: sdk.NewInt(1)},
			isErr: true,
		},
		{
			name:  "amount is negative",
			msg:   MsgCallToEvm{Sender: sdk.AccAddress{0x1}.String(), Evmaddr: sdk.AccAddress{0x2}.String(), Value: sdk.NewInt(-1), Calldata: "CALL DATA"},
			isErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {

			defer func() {
				r := recover()
				if tc.isErr {
					require.NotNil(t, r)
					require.Error(tt, r.(error))
				}
			}()
			tc.msg.GetSignBytes()
		})
	}
}

func TestMsgCallToEvm_ValidateBasic(t *testing.T) {
	wasmaAddr := sdk.AccAddress(make([]byte, 64)).String()
	addr := sdk.AccAddress{0x1}.String()
	errAddr := "error addr"
	testCases := []struct {
		name  string
		msg   MsgCallToEvm
		isErr bool
	}{
		{
			name:  "normal",
			msg:   MsgCallToEvm{Sender: wasmaAddr, Evmaddr: addr, Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: false,
		},
		{
			name:  "sender is empty",
			msg:   MsgCallToEvm{Evmaddr: addr, Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "sender is error",
			msg:   MsgCallToEvm{Sender: errAddr, Evmaddr: addr, Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "sender is not wasm addr",
			msg:   MsgCallToEvm{Sender: addr, Evmaddr: addr, Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "contract is error",
			msg:   MsgCallToEvm{Sender: wasmaAddr, Evmaddr: errAddr, Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "contract is wasm addr ",
			msg:   MsgCallToEvm{Sender: wasmaAddr, Evmaddr: wasmaAddr, Value: sdk.NewInt(1), Calldata: "CALL DATA"},
			isErr: true,
		},
		{
			name:  "amount is negative",
			msg:   MsgCallToEvm{Sender: wasmaAddr, Evmaddr: addr, Value: sdk.NewInt(-1), Calldata: "CALL DATA"},
			isErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			if err := tc.msg.ValidateBasic(); tc.isErr {
				require.Error(tt, err)
			}
		})
	}
}

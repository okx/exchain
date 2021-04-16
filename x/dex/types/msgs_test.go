package types

import (
	"reflect"
	"testing"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/stretchr/testify/require"
)

func TestMsg(t *testing.T) {

	common.InitConfig()
	addr, err := sdk.AccAddressFromBech32(TestTokenPairOwner)
	require.Nil(t, err)
	product := common.TestToken + "_" + common.NativeToken

	msgList := NewMsgList(addr, common.TestToken, common.NativeToken, sdk.NewDec(10))
	msgDeposit := NewMsgDeposit(product, sdk.NewDecCoin(common.NativeToken, sdk.NewInt(100)), addr)
	msgWithdraw := NewMsgWithdraw(product, sdk.NewDecCoin(common.NativeToken, sdk.NewInt(100)), addr)
	msgTransferOwnership := NewMsgTransferOwnership(addr, addr, product)

	// test msg.Route()、msg.Type()、msg.GetSigners()、GetSignBytes()
	type Want struct {
		RouterKey string
		Type      string
		SignBytes []byte
		Signers   []sdk.AccAddress
	}
	tests := []struct {
		name string
		msg  sdk.Msg
		want Want
	}{
		{"msgList", msgList,
			Want{"dex", "list", sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgList)), []sdk.AccAddress{addr}}},
		{"msgDeposit", msgDeposit,
			Want{"dex", typeMsgDeposit, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgDeposit)), []sdk.AccAddress{addr}}},
		{"msgWithdraw", msgWithdraw,
			Want{"dex", typeMsgWithdraw, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgWithdraw)), []sdk.AccAddress{addr}}},
		{"msgTransferOwnership", msgTransferOwnership,
			Want{"dex", typeMsgTransferOwnership, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgTransferOwnership)), []sdk.AccAddress{addr}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.msg.Route(); !reflect.DeepEqual(got, tt.want.RouterKey) {
				t.Errorf("%s.Route() = %v, want %v", tt.name, got, tt.want.RouterKey)
			}
			if got := tt.msg.Type(); !reflect.DeepEqual(got, tt.want.Type) {
				t.Errorf("%s.Type() = %v, want %v", tt.name, got, tt.want.Type)
			}
			if got := tt.msg.GetSigners(); !reflect.DeepEqual(got, tt.want.Signers) {
				t.Errorf("%s.GetSigners() = %v, want %v", tt.name, got, tt.want.Signers)
			}
			if got := tt.msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want.SignBytes) {
				t.Errorf("%s.ValidateBasic() = %v, want %v", tt.name, got, tt.want.SignBytes)
			}
		})
	}

	// test msg.ValidateBasic()

	fromPriKey := secp256k1.GenPrivKey()
	fromPubKey := fromPriKey.PubKey()
	fromAddr := sdk.AccAddress(fromPubKey.Address())

	toPriKey := secp256k1.GenPrivKey()
	toPubKey := toPriKey.PubKey()
	toAddr := sdk.AccAddress(toPubKey.Address())
	testBasics := []struct {
		name   string
		msg    sdk.Msg
		result bool
	}{
		{"msgList", msgList, true},
		{"msgDeposit", msgDeposit, true},
		{"msgWithdraw", msgWithdraw, true},

		{"deposit-invalid-amount", NewMsgDeposit(product, sdk.SysCoin{"", sdk.NewDec(1)}, addr), false},
		{"deposit-no-depositor", NewMsgDeposit(product, sdk.NewDecCoin(common.NativeToken, sdk.NewInt(1)), nil), false},
		{"withdraw-invalid-amount", NewMsgWithdraw(product, sdk.SysCoin{"", sdk.NewDec(1)}, addr), false},
		{"withdraw-no-depositor", NewMsgWithdraw(product, sdk.NewDecCoin(common.NativeToken, sdk.NewInt(1)), nil), false},

		{"transfer-no-from", NewMsgTransferOwnership(nil, toAddr, product), false},
		{"transfer-no-to", NewMsgTransferOwnership(fromAddr, nil, product), false},
		{"transfer-no-product", NewMsgTransferOwnership(fromAddr, toAddr, ""), false},
	}
	for _, tb := range testBasics {
		t.Run(tb.name, func(t *testing.T) {

			if tb.result {
				require.Nil(t, tb.msg.ValidateBasic(), "test: %v", tb.name)
			} else {
				require.NotNil(t, tb.msg.ValidateBasic(), "test: %v", tb.name)
			}
		})
	}

}

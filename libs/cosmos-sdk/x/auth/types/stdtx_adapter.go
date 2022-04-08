package types

import (
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

type IbcTx struct {
	*StdTx
	AuthInfoBytes []byte
	BodyBytes     []byte
}

func (tx *IbcTx) GetSignBytes(ctx sdk.Context, acc exported.Account) []byte {
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	return IbcSignBytes(
		chainID, accNum, acc.GetSequence(), tx.Fee, tx.Msgs, tx.Memo, tx.AuthInfoBytes, tx.BodyBytes,
	)
}

// StdSignBytes returns the bytes to sign for a transaction.
func IbcSignBytes(chainID string, accnum uint64,
	sequence uint64, fee StdFee, msgs []sdk.Msg,
	memo string, authInfoBytes []byte, bodyBytes []byte) []byte {

	signDoc := SignDoc{
		BodyBytes:     bodyBytes,
		AuthInfoBytes: authInfoBytes,
		ChainId:       chainID,
		AccountNumber: accnum,
	}

	r, err := signDoc.Marshal()
	if err != nil {
		return nil
	}
	return r
}

//////

///////////
type IbcViewMsg struct {
	TypeStr string `json:"type"`
	Data    string `json:"data"`
}

func NewIbcViewMsg(typeStr string, data string) *IbcViewMsg {
	return &IbcViewMsg{TypeStr: typeStr, Data: data}
}

func (b IbcViewMsg) Route() string {
	return ""
}

func (b IbcViewMsg) Type() string {
	return b.TypeStr
}

func (b IbcViewMsg) ValidateBasic() error {
	return nil
}

func (b IbcViewMsg) GetSignBytes() []byte {
	return nil
}

func (b IbcViewMsg) GetSigners() []sdk.AccAddress {
	return nil
}

func FromRelayIBCTx(cdc *codec.CodecProxy, tx *IbcTx) (*StdTx, error) {
	msgs := make([]sdk.Msg, 0)
	for _, msg := range tx.GetMsgs() {
		m := (interface{})(msg).(sdk.MsgProtoAdapter)
		data, err := cdc.GetProtocMarshal().MarshalJSON(m)
		if nil != err {
			return nil, err
		}
		msgs = append(msgs, NewIbcViewMsg("/"+proto.MessageName(m), string(data)))
	}
	return NewStdTx(msgs, tx.Fee, tx.Signatures, tx.Memo), nil
}

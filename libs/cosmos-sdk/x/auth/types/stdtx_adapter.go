package types

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

type IbcTx struct {
	*StdTx
	AuthInfoBytes []byte
	BodyBytes     []byte
	SignMode      signing.SignMode
	SigFee        IbcFee
	SigMsgs       []sdk.Msg
}

type StdIBCSignDoc struct {
	AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	Sequence      uint64            `json:"sequence" yaml:"sequence"`
	TimeoutHeight uint64            `json:"timeout_height,omitempty" yaml:"timeout_height"`
	ChainID       string            `json:"chain_id" yaml:"chain_id"`
	Memo          string            `json:"memo" yaml:"memo"`
	Fee           json.RawMessage   `json:"fee" yaml:"fee"`
	Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
}

type IbcFee struct {
	Amount sdk.CoinAdapters `json:"amount" yaml:"amount"`
	Gas    uint64           `json:"gas" yaml:"gas"`
}

func (tx *IbcTx) GetSignBytes(ctx sdk.Context, acc exported.Account) []byte {
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	if tx.SignMode == signing.SignMode_SIGN_MODE_DIRECT {
		return IbcDirectSignBytes(
			chainID, accNum, acc.GetSequence(), tx.Fee, tx.Msgs, tx.Memo, tx.AuthInfoBytes, tx.BodyBytes,
		)
	}

	return IbcAminoSignBytes(
		chainID, accNum, acc.GetSequence(), tx.SigFee, tx.SigMsgs, tx.Memo, tx.TimeoutHeight,
	)
}

func IbcAminoSignBytes(chainID string, accNum uint64,
	sequence uint64, fee IbcFee, msgs []sdk.Msg,
	memo string, height uint64) []byte {

	msgsBytes := make([]json.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}
	bz, err := ModuleCdc.MarshalJSON(StdIBCSignDoc{
		AccountNumber: accNum,
		ChainID:       chainID,
		Fee:           ModuleCdc.MustMarshalJSON(fee),
		Memo:          memo,
		Msgs:          msgsBytes,
		Sequence:      sequence,
		TimeoutHeight: height,
	})
	if err != nil {
		return nil
	}
	return sdk.MustSortJSON(bz)
}

// IbcDirectSignBytes returns the bytes to sign for a transaction.
func IbcDirectSignBytes(chainID string, accnum uint64,
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
type ProtobufViewMsg struct {
	TypeStr string `json:"type"`
	Data    string `json:"data"`
}

func NewProtobufViewMsg(typeStr string, data string) *ProtobufViewMsg {
	return &ProtobufViewMsg{TypeStr: typeStr, Data: data}
}

func (b ProtobufViewMsg) Route() string {
	return ""
}

func (b ProtobufViewMsg) Type() string {
	return b.TypeStr
}

func (b ProtobufViewMsg) ValidateBasic() error {
	return nil
}

func (b ProtobufViewMsg) GetSignBytes() []byte {
	return nil
}

func (b ProtobufViewMsg) GetSigners() []sdk.AccAddress {
	return nil
}

func FromProtobufTx(cdc *codec.CodecProxy, tx *IbcTx) (*StdTx, error) {
	msgs := make([]sdk.Msg, 0)
	for _, msg := range tx.GetMsgs() {
		m := (interface{})(msg).(sdk.MsgProtoAdapter)
		data, err := cdc.GetProtocMarshal().MarshalJSON(m)
		if nil != err {
			return nil, err
		}
		msgs = append(msgs, NewProtobufViewMsg("/"+proto.MessageName(m), string(data)))
	}
	return NewStdTx(msgs, tx.Fee, tx.Signatures, tx.Memo), nil
}

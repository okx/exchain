package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

type IbcTx struct {
	*StdTx

	AuthInfoBytes                []byte
	BodyBytes                    []byte
	SignMode                     []signing.SignMode
	SigFee                       IbcFee
	SigMsgs                      []sdk.Msg
	Sequences                    []uint64
	TxBodyHasUnknownNonCriticals bool
	HasExtensionOpt              bool
	Payer                        string
	ValidateParams               func() error
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

const (
	aminoNonCriticalFieldsError = "protobuf transaction contains unknown non-critical fields. This is a transaction malleability issue and SIGN_MODE_LEGACY_AMINO_JSON cannot be used."
	aminoExtentionError         = "SIGN_MODE_LEGACY_AMINO_JSON does not support protobuf extension options."
)

func (tx *IbcTx) GetSignBytes(ctx sdk.Context, index int, acc exported.Account) []byte {
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}
	if index > len(tx.SignMode) {
		panic(fmt.Sprintf("GetSignBytes index %d is upper than tx.SignMode Length %d", index, len(tx.SignMode)))
	}
	switch tx.SignMode[index] {
	case signing.SignMode_SIGN_MODE_DIRECT:
		return IbcDirectSignBytes(
			chainID, accNum, acc.GetSequence(), tx.Fee, tx.Msgs, tx.Memo, tx.AuthInfoBytes, tx.BodyBytes,
		)
	case signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON:
		if tx.TxBodyHasUnknownNonCriticals {
			panic(aminoNonCriticalFieldsError)
		}
		if tx.HasExtensionOpt {
			panic(aminoExtentionError)
		}
		return IbcAminoSignBytes(
			chainID, accNum, acc.GetSequence(), tx.SigFee, tx.SigMsgs, tx.Memo, 0,
		)
	case signing.SignMode_SIGN_MODE_UNSPECIFIED:
		//Compatible with cosmojs for simulating tx
		return nil
	default:
		//does not other mode
		panic(fmt.Sprintf("ibctx not support sign mode: %s", tx.SignMode[index].String()))
	}
}

func (tx *IbcTx) ValidateBasic() error {
	err := tx.StdTx.ValidateBasic()
	if err != nil {
		return err
	}
	err = tx.ValidateParams()
	if err != nil {
		return err
	}

	return nil
}

func (tx *IbcTx) VerifySequence(index int, acc exported.Account) error {
	//check
	if index > len(tx.Sequences) {
		return errors.New("verify sequence error index not fit")
	}
	seq := tx.Sequences[index]
	if seq != acc.GetSequence() {
		return fmt.Errorf("verify sequence error expected:%d,got:%d", acc.GetSequence(), seq)
	}

	return nil
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
		m := (interface{})(msg).(sdk.MsgAdapter)
		data, err := cdc.GetProtocMarshal().MarshalJSON(m)
		if nil != err {
			return nil, err
		}
		msgs = append(msgs, NewProtobufViewMsg("/"+proto.MessageName(m), string(data)))
	}
	return NewStdTx(msgs, tx.Fee, tx.Signatures, tx.Memo), nil
}

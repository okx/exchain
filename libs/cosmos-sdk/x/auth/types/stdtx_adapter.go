package types

import (
	"encoding/json"
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

type StdIBCSignDoc struct {
	AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	Sequence      uint64            `json:"sequence" yaml:"sequence"`
	TimeoutHeight uint64            `json:"timeout_height,omitempty" yaml:"timeout_height"`
	ChainID       string            `json:"chain_id" yaml:"chain_id"`
	Memo          string            `json:"memo" yaml:"memo"`
	Fee           json.RawMessage   `json:"fee" yaml:"fee"`
	Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
}

type IBCCoin struct {
	Denom  string  `json:"denom"`
	Amount sdk.Int `json:"amount"`
}

type IBCFee struct {
	Amount []IBCCoin `json:"amount" yaml:"amount"`
	Gas    uint64    `json:"gas" yaml:"gas"`
}

func feeToIBCFeeBytes(fee StdFee) []byte {
	var ibcFee IBCFee
	ibcFee.Gas = fee.Gas
	for _, coin := range fee.Amount {
		ibcCoin := IBCCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.TruncateInt(),
		}
		ibcFee.Amount = append(ibcFee.Amount, ibcCoin)
	}
	bz, err := ModuleCdc.MarshalJSON(ibcFee)
	if err != nil {
		panic(err)
	}
	return bz
}

func (tx *IbcTx) GetSignBytes(ctx sdk.Context, acc exported.Account) []byte {
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	msgsBytes := make([]json.RawMessage, 0, len(tx.Msgs))
	for _, msg := range tx.Msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}
	bz, err := ModuleCdc.MarshalJSON(StdIBCSignDoc{
		AccountNumber: accNum,
		ChainID:       chainID,
		Fee:           json.RawMessage(feeToIBCFeeBytes(tx.Fee)),
		Memo:          tx.Memo,
		Msgs:          msgsBytes,
		Sequence:      acc.GetSequence(),
		TimeoutHeight: tx.TimeoutHeight,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
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

func FromRelayIBCTx(cdc *codec.CodecProxy, tx *IbcTx) (StdTx, error) {
	msgs := make([]sdk.Msg, 0)
	for _, msg := range tx.GetMsgs() {
		m := (interface{})(msg).(sdk.MsgProtoAdapter)
		data, err := cdc.GetProtocMarshal().MarshalJSON(m)
		if nil != err {
			return StdTx{}, err
		}
		msgs = append(msgs, NewIbcViewMsg("/"+proto.MessageName(m), string(data)))
	}
	return NewStdTx(msgs, tx.Fee, tx.Signatures, tx.Memo), nil
}

package types

import (
	"fmt"
	"github.com/tendermint/go-amino"

	"github.com/okex/exchain/app/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg    = MsgEthermint{}
)


// message type and route constants
const (
	// TypeMsgEthermint defines the type string of Ethermint message
	TypeMsgEthermint = "ethermint"
)

// MsgEthermint implements a cosmos equivalent structure for Ethereum transactions
type MsgEthermint struct {
	AccountNonce uint64          `json:"nonce"`
	Price        sdk.Int         `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *sdk.AccAddress `json:"to" rlp:"nil"` // nil means contract creation
	Amount       sdk.Int         `json:"value"`
	Payload      []byte          `json:"input"`

	// From address (formerly derived from signature)
	From sdk.AccAddress `json:"from"`
}

// NewMsgEthermint returns a reference to a new Ethermint transaction
func NewMsgEthermint(
	nonce uint64, to *sdk.AccAddress, amount sdk.Int,
	gasLimit uint64, gasPrice sdk.Int, payload []byte, from sdk.AccAddress,
) MsgEthermint {
	return MsgEthermint{
		AccountNonce: nonce,
		Price:        gasPrice,
		GasLimit:     gasLimit,
		Recipient:    to,
		Amount:       amount,
		Payload:      payload,
		From:         from,
	}
}

func (msg MsgEthermint) String() string {
	return fmt.Sprintf("nonce=%d gasPrice=%s gasLimit=%d recipient=%s amount=%s data=0x%x from=%s",
		msg.AccountNonce, msg.Price, msg.GasLimit, msg.Recipient, msg.Amount, msg.Payload, msg.From)
}

// Route should return the name of the module
func (msg MsgEthermint) Route() string { return RouterKey }

// Type returns the action of the message
func (msg MsgEthermint) Type() string { return TypeMsgEthermint }

// GetSignBytes encodes the message for signing
func (msg MsgEthermint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgEthermint) ValidateBasic() error {
	if msg.Price.IsZero() {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be 0")
	}

	if msg.Price.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be negative %s", msg.Price)
	}

	// Amount can be 0
	if msg.Amount.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "amount cannot be negative %s", msg.Amount)
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgEthermint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// To returns the recipient address of the transaction. It returns nil if the
// transaction is a contract creation.
func (msg MsgEthermint) To() *ethcmn.Address {
	if msg.Recipient == nil {
		return nil
	}

	addr := ethcmn.BytesToAddress(msg.Recipient.Bytes())
	return &addr
}

func (msg *MsgEthermint) UnmarshalFromAmino(data []byte) error {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]
		if len(data) == 0 {
			break
		}

		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		data = data[1:]

		if pbType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			data = data[n:]
			if len(data) < int(dataLen) {
				return fmt.Errorf("invalid tx data")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			var n int
			msg.AccountNonce, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
		case 2:
			msg.Price, err = sdk.NewIntFromAmino(subData)
			if err != nil {
				return err
			}
		case 3:
			var n int
			msg.GasLimit, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
		case 4:
			tmp := make(sdk.AccAddress, dataLen)
			msg.Recipient = &tmp
			copy(tmp[:], subData)
		case 5:
			msg.Amount, err = sdk.NewIntFromAmino(subData)
		case 6:
			msg.Payload = make([]byte, dataLen)
			copy(msg.Payload, subData)
		case 7:
			msg.From = make([]byte, dataLen)
			copy(msg.From, subData)
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

package types

import (
	"math/big"

	"github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/tendermint/libs/tlv"
	"github.com/okex/exchain/libs/tendermint/mempool"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/pkg/errors"

	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
)

var (
	_ sdk.Tx  = MsgEthereumCheckedTx{}
	_ sdk.Msg = MsgEthereumCheckedTx{}
)

// EthereumCheckedSignature carry both useful information to verify the confident node
type EthereumCheckedSignature struct {
	Signature []byte
	NodeKey   []byte
}

// WithPayload to construct the object
func (sig EthereumCheckedSignature) WithPayload(payload []byte) EthereumCheckedSignature {
	sig.decode(payload)
	return sig
}

// nolint
func (sig *EthereumCheckedSignature) encode() []byte {
	buf := tlv.NewBuffer()
	buf.Write(sig.Signature)
	buf.Write(sig.NodeKey)
	return buf.Bytes()
}

// nolint
func (sig *EthereumCheckedSignature) decode(pyload []byte) error {
	buf := tlv.With(pyload)
	if v, t := buf.Read(); t == tlv.Bytes {
		sig.Signature = v.([]byte)
	} else {
		return errors.New("TODO: should convert to sdkerrors")
	}
	if v, t := buf.Read(); t == tlv.Bytes {
		sig.NodeKey = v.([]byte)
	} else {

	}
	return nil
}

// Veify the signature of this Tx
// the pub keys will be set in the config file in tendermint
func (sig *EthereumCheckedSignature) Verify(msg []byte, pubs []ed25519.PubKeyEd25519) (confident bool, err error) {
	if len(sig.Signature) == 0 || len(sig.NodeKey) == 0 {
		return false, nil
	}
	if len(sig.NodeKey) != ed25519.PubKeyEd25519Size {
		return false, errors.New("node key size invalid")
	}
	pubKey := ed25519.PubKeyEd25519{}
	pubKey.UnmarshalFromAmino(sig.NodeKey)
	verified := pubKey.VerifyBytes(msg, sig.Signature)
	if !verified {
		return false, errors.New("signature invalid")
	}
	for _, v := range pubs {
		if v.Equals(pubKey) {
			confident = true
			return
		}
	}
	return
}

// MsgEthereumCheckedTx is the specific Tx
// which carries the signature of the publisher in the P2P network
// To figure out if this Tx could reduce some nouse check in CheckTx
type MsgEthereumCheckedTx struct {
	// convert for ethermint
	Data TxData
	// used to identify the message checked
	Payload []byte `json:"payload"`
}

func (msg MsgEthereumCheckedTx) String() string {
	return ""
}

//================================================
// specific methods

// Sign this struct with the p2p.NodeKey to specific this message is a Checked message
func (msg *MsgEthereumCheckedTx) Sign(_ *big.Int, nodeKey *p2p.NodeKey) error {
	wait, err := msg.Data.MarshalAmino()
	if err != nil {
		return err
	}
	signature, err := nodeKey.PrivKey.Sign(wait)
	if err != nil {
		return err
	}
	payload := EthereumCheckedSignature{
		Signature: signature,
		NodeKey:   nodeKey.PubKey().Bytes(),
	}
	msg.Payload = payload.encode()
	return nil
}

// ConvertToOriginTx for application to do ante
func (msg *MsgEthereumCheckedTx) ConvertToOriginTx() MsgEthereumTx {
	return MsgEthereumTx{Data: msg.Data}
}

//===============================================

//=================================================
//sdk.Msg impls

// nolint
func (msg MsgEthereumCheckedTx) Route() string {
	return RouterKey
}

// nolint
func (msg MsgEthereumCheckedTx) GetSignBytes() []byte {
	panic("must use the MsgEthereumTx to call ")
}

// nolint
func (msg MsgEthereumCheckedTx) Type() string {
	return TypeMsgEthereumCheckekTx
}

// nolint
func (msg MsgEthereumCheckedTx) ValidateBasic() error {
	if msg.Data.Price.Cmp(big.NewInt(0)) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be 0")
	}

	if msg.Data.Price.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be negative %s", msg.Data.Price)
	}

	// Amount can be 0
	if msg.Data.Amount.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "amount cannot be negative %s", msg.Data.Amount)
	}
	return nil
}

// nolint
func (msg MsgEthereumCheckedTx) GetSigners() []sdk.AccAddress {
	panic("unsupport method GetSigners in this object")
}

//================================================

//================================================
// sdk.Tx impls

// nolint
func (msg MsgEthereumCheckedTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// nolint
func (msg MsgEthereumCheckedTx) GetTxInfo(ctx sdk.Context) mempool.ExTxInfo {
	exTxInfo := mempool.ExTxInfo{
		Sender:   "",
		GasPrice: big.NewInt(0),
		Nonce:    msg.Data.AccountNonce,
	}
	chainIDEpoch, err := types.ParseChainID(ctx.ChainID())
	if err != nil {
		return exTxInfo
	}
	originTx := msg.ConvertToOriginTx()
	fromSigCache, err := originTx.VerifySig(chainIDEpoch, ctx.BlockHeight(), ctx.SigCache())
	if err != nil {
		return exTxInfo
	}
	from := fromSigCache.GetFrom()
	exTxInfo.Sender = from.String()
	exTxInfo.GasPrice = msg.Data.Price
	return exTxInfo
}

// nolint
func (msg MsgEthereumCheckedTx) GetGasPrice() *big.Int {
	return msg.Data.Price
}

// nolint
func (msg MsgEthereumCheckedTx) GetTxFnSignatureInfo() ([]byte, int) {
	if msg.Data.Recipient == nil {
		return DefaultDeployContractFnSignature, len(msg.Data.Payload)
	}

	// most case is transfer token
	if len(msg.Data.Payload) < 4 {
		return DefaultSendCoinFnSignature, 0
	}

	// call contract case (some times will together with transfer token case)
	recipient := msg.Data.Recipient.Bytes()
	methodId := msg.Data.Payload[0:4]
	return append(recipient, methodId...), 0
}

// nolint
func (msg MsgEthereumCheckedTx) GetTxCarriedData() []byte {
	return msg.Payload
}

//================================================

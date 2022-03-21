package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/tendermint/go-amino"
)

var (
	_ sdk.Msg    = (*MsgEthereumTx)(nil)
	_ sdk.Tx     = (*MsgEthereumTx)(nil)
	_ ante.FeeTx = (*MsgEthereumTx)(nil)
)

var big8 = big.NewInt(8)
var DefaultDeployContractFnSignature = ethcmn.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001")
var DefaultSendCoinFnSignature = ethcmn.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000010")

// message type and route constants
const (
	// TypeMsgEthereumTx defines the type string of an Ethereum tranasction
	TypeMsgEthereumTx = "ethereum"
)

// MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.
type MsgEthereumTx struct {
	Data TxData

	sdk.BaseTx `json:"-" rlp:"-"`
}

func (tx *MsgEthereumTx) GetType() sdk.TransactionType {
	return sdk.EvmTxType
}

func (tx *MsgEthereumTx) SetFrom(addr string) {
	tx.From = addr
}

func (msg *MsgEthereumTx) GetNonce() uint64 {
	return msg.Data.AccountNonce
}

func (msg *MsgEthereumTx) GetFee() sdk.Coins {
	fee := make(sdk.Coins, 1)
	fee[0] = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecFromBigIntWithPrec(msg.Fee(), sdk.Precision))
	return fee
}

func (msg MsgEthereumTx) FeePayer(ctx sdk.Context) sdk.AccAddress {
	if from := ctx.From(); len(from) > 0 {
		return sdk.AccAddress(ethcmn.HexToAddress(from).Bytes())
	}

	err := msg.VerifySig(msg.ChainID(), ctx.BlockHeight())
	if err != nil {
		return nil
	}

	return msg.AccountAddress()
}

// NewMsgEthereumTx returns a reference to a new Ethereum transaction message.
func NewMsgEthereumTx(
	nonce uint64, to *ethcmn.Address, amount *big.Int,
	gasLimit uint64, gasPrice *big.Int, payload []byte,
) *MsgEthereumTx {
	return newMsgEthereumTx(nonce, to, amount, gasLimit, gasPrice, payload)
}

// NewMsgEthereumTxContract returns a reference to a new Ethereum transaction
// message designated for contract creation.
func NewMsgEthereumTxContract(
	nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, payload []byte,
) *MsgEthereumTx {
	return newMsgEthereumTx(nonce, nil, amount, gasLimit, gasPrice, payload)
}

func newMsgEthereumTx(
	nonce uint64, to *ethcmn.Address, amount *big.Int,
	gasLimit uint64, gasPrice *big.Int, payload []byte,
) *MsgEthereumTx {
	if len(payload) > 0 {
		payload = ethcmn.CopyBytes(payload)
	}

	txData := TxData{
		AccountNonce: nonce,
		Recipient:    to,
		Payload:      payload,
		GasLimit:     gasLimit,
		Amount:       new(big.Int),
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}

	if amount != nil {
		txData.Amount.Set(amount)
	}
	if gasPrice != nil {
		txData.Price.Set(gasPrice)
	}

	return &MsgEthereumTx{Data: txData}
}

func (msg *MsgEthereumTx) String() string {
	return msg.Data.String()
}

// Route returns the route value of an MsgEthereumTx.
func (msg *MsgEthereumTx) Route() string { return RouterKey }

// Type returns the type value of an MsgEthereumTx.
func (msg *MsgEthereumTx) Type() string { return TypeMsgEthereumTx }

// ValidateBasic implements the sdk.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
func (msg *MsgEthereumTx) ValidateBasic() error {
	if msg.Data.Price.Sign() <= 0 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be non positive %s", msg.Data.Price)
	}

	// Amount can be 0
	if msg.Data.Amount.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "amount cannot be negative %s", msg.Data.Amount)
	}

	return nil
}

// To returns the recipient address of the transaction. It returns nil if the
// transaction is a contract creation.
func (msg *MsgEthereumTx) To() *ethcmn.Address {
	return msg.Data.Recipient
}

// GetSigners returns the expected signers for an Ethereum transaction message.
// For such a message, there should exist only a single 'signer'.
//
// NOTE: This method panics if 'VerifySig' hasn't been called first.
func (msg *MsgEthereumTx) GetSigners() []sdk.AccAddress {
	addr := msg.AccountAddress()
	if addr.Empty() {
		panic("must use 'VerifySig' with a chain ID to get the from addr")
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the Amino bytes of an Ethereum transaction message used
// for signing.
//
// NOTE: This method cannot be used as a chain ID is needed to create valid bytes
// to sign over. Use 'RLPSignBytes' instead.
func (msg *MsgEthereumTx) GetSignBytes() []byte {
	panic("must use 'RLPSignBytes' with a chain ID to get the valid bytes to sign")
}

// RLPSignBytes returns the RLP hash of an Ethereum transaction message with a
// given chainID used for signing.
func (msg *MsgEthereumTx) RLPSignBytes(chainID *big.Int) ethcmn.Hash {
	return rlpHash([]interface{}{
		msg.Data.AccountNonce,
		msg.Data.Price,
		msg.Data.GasLimit,
		msg.Data.Recipient,
		msg.Data.Amount,
		msg.Data.Payload,
		chainID, uint(0), uint(0),
	})
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (msg *MsgEthereumTx) HomesteadSignHash() ethcmn.Hash {
	return rlpHash([]interface{}{
		msg.Data.AccountNonce,
		msg.Data.Price,
		msg.Data.GasLimit,
		msg.Data.Recipient,
		msg.Data.Amount,
		msg.Data.Payload,
	})
}

// EncodeRLP implements the rlp.Encoder interface.
func (msg *MsgEthereumTx) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &msg.Data)
}

// DecodeRLP implements the rlp.Decoder interface.
func (msg *MsgEthereumTx) DecodeRLP(s *rlp.Stream) error {
	_, _, err := s.Kind()
	if err != nil {
		// return error if stream is too large
		return err
	}

	if err := s.Decode(&msg.Data); err != nil {
		return err
	}

	return nil
}

// Sign calculates a secp256k1 ECDSA signature and signs the transaction. It
// takes a private key and chainID to sign an Ethereum transaction according to
// EIP155 standard. It mutates the transaction as it populates the V, R, S
// fields of the Transaction's Signature.
func (msg *MsgEthereumTx) Sign(chainID *big.Int, priv *ecdsa.PrivateKey) error {
	txHash := msg.RLPSignBytes(chainID)

	sig, err := ethcrypto.Sign(txHash[:], priv)
	if err != nil {
		return err
	}

	if len(sig) != 65 {
		return fmt.Errorf("wrong size for signature: got %d, want 65", len(sig))
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])

	var v *big.Int

	if chainID.Sign() == 0 {
		v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	} else {
		v = big.NewInt(int64(sig[64] + 35))
		chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))

		v.Add(v, chainIDMul)
	}

	msg.Data.V = v
	msg.Data.R = r
	msg.Data.S = s
	return nil
}

func (msg *MsgEthereumTx) firstVerifySig(chainID *big.Int) error {
	if msg.GetFrom() != "" {
		return nil
	}

	from, ok := tmtypes.SignatureCache().Get(string(msg.TxHash()))
	if ok {
		msg.BaseTx.From = from
		return nil
	}

	var V *big.Int
	var sigHash ethcmn.Hash
	if isProtectedV(msg.Data.V) {
		// do not allow recovery for transactions with an unprotected chainID
		if chainID.Sign() == 0 {
			return errors.New("chainID cannot be zero")
		}

		chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))
		V = new(big.Int).Sub(msg.Data.V, chainIDMul)
		V.Sub(V, big8)

		sigHash = msg.RLPSignBytes(chainID)
	} else {
		V = msg.Data.V

		sigHash = msg.HomesteadSignHash()
	}

	sender, err := recoverEthSig(msg.Data.R, msg.Data.S, V, sigHash)
	if err != nil {
		return err
	}
	tmtypes.SignatureCache().Add(string(msg.TxHash()), sender.String())
	msg.BaseTx.From = sender.String()
	return nil
}

// VerifySig attempts to verify a Transaction's signature for a given chainID.
// A derived address is returned upon success or an error if recovery fails.
func (msg *MsgEthereumTx) VerifySig(chainID *big.Int, height int64) error {
	if !isProtectedV(msg.Data.V) && tmtypes.HigherThanMercury(height) {
		return errors.New("deprecated support for homestead Signer")
	}
	return msg.firstVerifySig(chainID)
}

// codes from go-ethereum/core/types/transaction.go:122
func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28
	}
	// anything not 27 or 28 is considered protected
	return true
}

// GetGas implements the GasTx interface. It returns the GasLimit of the transaction.
func (msg *MsgEthereumTx) GetGas() uint64 {
	return msg.Data.GasLimit
}

// Fee returns gasprice * gaslimit.
func (msg *MsgEthereumTx) Fee() *big.Int {
	return new(big.Int).Mul(msg.Data.Price, new(big.Int).SetUint64(msg.Data.GasLimit))
}

// ChainID returns which chain id this transaction was signed for (if at all)
func (msg *MsgEthereumTx) ChainID() *big.Int {
	return deriveChainID(msg.Data.V)
}

// Cost returns amount + gasprice * gaslimit.
func (msg *MsgEthereumTx) Cost() *big.Int {
	total := msg.Fee()
	total.Add(total, msg.Data.Amount)
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
func (msg *MsgEthereumTx) RawSignatureValues() (v, r, s *big.Int) {
	return msg.Data.V, msg.Data.R, msg.Data.S
}

// From loads the ethereum sender address from the sigcache and returns an
// sdk.AccAddress from its bytes
func (msg *MsgEthereumTx) AccountAddress() sdk.AccAddress {
	return ethcmn.FromHex(msg.GetFrom())
}

func (msg *MsgEthereumTx) EthereumAddress() ethcmn.Address {
	return ethcmn.HexToAddress(msg.GetFrom())
}

// deriveChainID derives the chain id from the given v parameter
func deriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}

func (msg *MsgEthereumTx) UnmarshalFromAmino(cdc *amino.Codec, data []byte) error {
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
			if err := msg.Data.UnmarshalFromAmino(cdc, subData); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

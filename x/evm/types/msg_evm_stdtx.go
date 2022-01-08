package types

import (
	"github.com/okex/exchain/libs/tendermint/mempool"
	"math/big"

	"github.com/okex/exchain/app/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

//type Tx interface {
//	GetMsgs() []Msg
//	// ValidateBasic does a simple and lightweight validation check that doesn't
//	// require access to any other information.
//	ValidateBasic() error
//	GetTxInfo(ctx Context) mempool.ExTxInfo
//	GetGasPrice() *big.Int
//	GetTxFnSignatureInfo() ([]byte, int)
//	GetTxCarriedData() []byte
//}


// ValidateBasic implements the sdk.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
//func (msg MsgEthereumTx) ValidateBasic() error {
//	if msg.Data.Price.Cmp(big.NewInt(0)) == 0 {
//		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be 0")
//	}
//
//	if msg.Data.Price.Sign() == -1 {
//		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be negative %s", msg.Data.Price)
//	}
//
//	// Amount can be 0
//	if msg.Data.Amount.Sign() == -1 {
//		return sdkerrors.Wrapf(types.ErrInvalidValue, "amount cannot be negative %s", msg.Data.Amount)
//	}
//
//	return nil
//}


//___________________std tx______________________

// GetMsgs returns a single MsgEthereumTx as an sdk.Msg.
func (msg MsgEthereumTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}


// Return tx sender and gas price
func (msg MsgEthereumTx) GetTxInfo(ctx sdk.Context) mempool.ExTxInfo {
	exTxInfo := mempool.ExTxInfo{
		Sender:   "",
		GasPrice: big.NewInt(0),
		Nonce:    msg.Data.AccountNonce,
	}

	chainIDEpoch, err := types.ParseChainID(ctx.ChainID())
	if err != nil {
		return exTxInfo
	}

	// Verify signature and retrieve sender address
	fromSigCache, err := msg.VerifySig(chainIDEpoch, ctx.BlockHeight(), ctx.SigCache())
	if err != nil {
		return exTxInfo
	}

	from := fromSigCache.GetFrom()
	exTxInfo.Sender = from.String()
	exTxInfo.GasPrice = msg.Data.Price

	return exTxInfo
}

// GetGasPrice return gas price
func (msg MsgEthereumTx) GetGasPrice() *big.Int {
	return msg.Data.Price
}

func (msg MsgEthereumTx) GetTxFnSignatureInfo() ([]byte, int) {
	// deploy contract case
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

// GetTxCarriedData implement the sdk.Tx interface
func (msg MsgEthereumTx) GetTxCarriedData() []byte {
	return nil
}

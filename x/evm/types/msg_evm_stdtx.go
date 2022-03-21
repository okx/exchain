package types

import (
	"github.com/okex/exchain/app/types"
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

//___________________std tx______________________

// GetMsgs returns a single MsgEthereumTx as an sdk.Msg.
func (msg *MsgEthereumTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetGasPrice return gas price
func (msg *MsgEthereumTx) GetGasPrice() *big.Int {
	return msg.Data.Price
}

func (msg *MsgEthereumTx) GetTxFnSignatureInfo() ([]byte, int) {
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

func (msg *MsgEthereumTx) GetPartnerInfo(ctx sdk.Context) (string, string) {
	chainIDEpoch, err := types.ParseChainID(ctx.ChainID())
	if err != nil {
		return "", ""
	}

	// Verify signature and retrieve sender address
	err = msg.VerifySig(chainIDEpoch, ctx.BlockHeight())
	if err != nil {
		return "", ""
	}

	to := ""
	if msg.Data.Recipient != nil {
		to = EthAddressStringer(*msg.Data.Recipient).String()
	}

	return msg.From, to
}

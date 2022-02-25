package types

import (
	"math/big"

	"github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/mempool"
)

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

	chainIDEpoch := types.GetChainIdEpoch()

	if ctx.From() == "" {
		// Verify signature and retrieve sender address
		fromSigCache, err := msg.VerifySig(chainIDEpoch, ctx.BlockHeight(), ctx.TxBytes())
		if err != nil {
			return exTxInfo
		}

		from := fromSigCache.GetFrom()
		exTxInfo.Sender = from.String()
	} else {
		exTxInfo.Sender = ctx.From()
	}

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

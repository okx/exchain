package types

import (
	"encoding/json"

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

	msgsBytes := make([]json.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}

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

package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)
var (
	_ sdk.Tx = (*WrappedTx)(nil)
)

type RawWrappedTx struct {
	Payload   []byte  `json:"payload"`   // std tx or evm tx
	Metadata  []byte  `json:"metadata"`  // customized message from the node who signs the tx
	Signature []byte  `json:"signature"` // signature for payload+metadata
	NodeKey   []byte  `json:"nodeKey"`   // pub key of the node who signs the tx
}

type WrappedTx struct {
	sdk.Tx
	Metadata  []byte
	Signature []byte
	NodeKey   []byte
}

func (msg WrappedTx) String() string {
	return fmt.Sprintf("StdTx=<%s>, Metadata=<%s>, Signature=<%s>, NodeKey=<%s>",
		msg.Tx,
		string(msg.Metadata),
		string(msg.Signature),
		string(msg.NodeKey),
		)
}

func (wtx WrappedTx) GetPayloadTx() sdk.Tx {
	return wtx.Tx
}

func EncodeWrappedTx(txbytes []byte, info *sdk.ExTxInfo, replace bool) ([]byte, error) {

	payload := txbytes
	if replace {
		// txbytes is a wrapped one
		raw := &RawWrappedTx{}
		err := json.Unmarshal(txbytes, raw)
		if err != nil {
			return nil, err
		}
		payload = raw.Payload
	}

	wrapped := &RawWrappedTx{
		Payload: payload,
		NodeKey: info.NodeKey,
		Signature: info.Signature,
		Metadata: info.Metadata,
	}

	return json.Marshal(wrapped)
}

func DecodeWrappedTx(txbytes []byte, payloadDecoder sdk.TxDecoder, heights ...int64) (sdk.Tx, error) {

	raw := &RawWrappedTx{}
	err := json.Unmarshal(txbytes, raw)
	if err != nil {
		return nil, err
	}

	payloadTx, err := payloadDecoder(raw.Payload, heights...)
	if err != nil {
		return nil, err
	}

	tx := WrappedTx{
		Tx: payloadTx,
		NodeKey: raw.NodeKey,
		Metadata: raw.Metadata,
		Signature: raw.Signature,
	}
	return tx, err
}
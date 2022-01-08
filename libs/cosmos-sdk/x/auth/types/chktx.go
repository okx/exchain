package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)
var (
	_ sdk.Tx = (*CheckedTx)(nil)

)

type RawCheckedTx struct {
	Payload   []byte  `json:"payload"`   // std tx or evm tx
	Metadata  []byte  `json:"metadata"`  // customized message from the node who signs the tx
	Signature []byte  `json:"signature"` // signature for payload+metadata
	NodeKey   []byte  `json:"nodeKey"`   // pub key of the node who signs the tx
}

type CheckedTx struct {
	sdk.Tx
	Metadata  []byte  `json:"metadata"`  // customized message from the node who signs the tx
	Signature []byte  `json:"signature"` // signature for payload+metadata
	NodeKey   []byte  `json:"nodeKey"`   // pub key of the node who signs the tx
}

func (msg CheckedTx) String() string {
	return fmt.Sprintf("StdTx=<%s>, Metadata=<%s>, Signature=<%s>, NodeKey=<%s>",
		msg.Tx,
		msg.Metadata,
		msg.Signature,
		msg.NodeKey,
		)
}


func EncodeCheckedTx(payload []byte, info *sdk.ExTxInfo) ([]byte, error) {

	raw := &RawCheckedTx{
		Payload: payload,
		NodeKey: info.NodeKey,
		Signature: info.Signature,
		Metadata: info.Metadata,
	}

	return json.Marshal(raw)
}

func DecodeCheckedTx(b []byte, payloadDecoder func(txBytes []byte) (sdk.Tx, error)) (sdk.Tx, error) {

	raw := &RawCheckedTx{}
	err := json.Unmarshal(b, raw)
	if err != nil {
		return nil, err
	}

	payloadTx, err := payloadDecoder(raw.Payload)
	if err != nil {
		return nil, err
	}

	tx := &CheckedTx{
		Tx: payloadTx,
		NodeKey: raw.NodeKey,
		Metadata: raw.Metadata,
		Signature: raw.Signature,
	}
	return tx, err
}
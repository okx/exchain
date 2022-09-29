package types

type EventAdapter struct {
	Type       string                  `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes []EventAttributeAdapter `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
}

// EventAttributeAdapter is a single key-value pair, associated with an event.
type EventAttributeAdapter struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Index bool   `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

// TxResult contains results of executing the transaction.
//
// One usage is indexing transaction results.
type TxResult struct {
	Height int64        `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Index  uint32       `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Tx     []byte       `protobuf:"bytes,3,opt,name=tx,proto3" json:"tx,omitempty"`
	Result ExecTxResult `protobuf:"bytes,4,opt,name=result,proto3" json:"result"`
}

// ExecTxResult contains results of executing one individual transaction.
//
// * Its structure is equivalent to #ResponseDeliverTx which will be deprecated/deleted
type ExecTxResult struct {
	Code      uint32         `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data      []byte         `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log       string         `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info      string         `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasWanted int64          `protobuf:"varint,5,opt,name=gas_wanted,json=gasWanted,proto3" json:"gas_wanted,omitempty"`
	GasUsed   int64          `protobuf:"varint,6,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Events    []EventAdapter `protobuf:"bytes,7,rep,name=events,proto3" json:"events,omitempty"`
	Codespace string         `protobuf:"bytes,8,opt,name=codespace,proto3" json:"codespace,omitempty"`
}

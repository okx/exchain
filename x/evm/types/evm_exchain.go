package types

import "time"

// from ethermint/x/evm/types/evm.pb.go

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	// custom javascript tracer
	Tracer string `protobuf:"bytes,1,opt,name=tracer,proto3" json:"tracer,omitempty"`
	// overrides the default timeout of 5 seconds for JavaScript-based tracing
	// calls
	Timeout string `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// number of blocks the tracer is willing to go back
	Reexec uint64 `protobuf:"varint,3,opt,name=reexec,proto3" json:"reexec,omitempty"`
	// disable stack capture
	DisableStack bool `protobuf:"varint,5,opt,name=disable_stack,json=disableStack,proto3" json:"disableStack"`
	// disable storage capture
	DisableStorage bool `protobuf:"varint,6,opt,name=disable_storage,json=disableStorage,proto3" json:"disableStorage"`
	// print output during capture end
	Debug bool `protobuf:"varint,8,opt,name=debug,proto3" json:"debug,omitempty"`
	// maximum length of output, but zero means unlimited
	Limit int32 `protobuf:"varint,9,opt,name=limit,proto3" json:"limit,omitempty"`
	// Chain overrides, can be used to execute a trace using future fork rules
	Overrides *ChainConfig `protobuf:"bytes,10,opt,name=overrides,proto3" json:"overrides,omitempty"`
	// enable memory capture
	EnableMemory bool `protobuf:"varint,11,opt,name=enable_memory,json=enableMemory,proto3" json:"enableMemory"`
	// enable return data capture
	EnableReturnData bool `protobuf:"varint,12,opt,name=enable_return_data,json=enableReturnData,proto3" json:"enableReturnData"`
}

// QueryTraceTxRequest defines TraceTx request
type QueryTraceTxRequest struct {
	// msgEthereumTx for the requested transaction
	Msg *MsgEthereumTx `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	// transaction index
	TxIndex uint64 `protobuf:"varint,2,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	// TraceConfig holds extra parameters to trace functions.
	TraceConfig *TraceConfig `protobuf:"bytes,3,opt,name=trace_config,json=traceConfig,proto3" json:"trace_config,omitempty"`
	// the predecessor transactions included in the same block
	// need to be replayed first to get correct context for tracing.
	Predecessors []*MsgEthereumTx `protobuf:"bytes,4,rep,name=predecessors,proto3" json:"predecessors,omitempty"`
	// block number of requested transaction
	BlockNumber int64 `protobuf:"varint,5,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	// block hex hash of requested transaction
	BlockHash string `protobuf:"bytes,6,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// block time of requested transaction
	BlockTime time.Time `protobuf:"bytes,7,opt,name=block_time,json=blockTime,proto3,stdtime" json:"block_time"`
}

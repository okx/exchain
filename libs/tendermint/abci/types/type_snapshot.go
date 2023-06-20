package types

type Snapshot struct {
	Height   uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Format   uint32 `protobuf:"varint,2,opt,name=format,proto3" json:"format,omitempty"`
	Chunks   uint32 `protobuf:"varint,3,opt,name=chunks,proto3" json:"chunks,omitempty"`
	Hash     []byte `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Metadata []byte `protobuf:"bytes,5,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

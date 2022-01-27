package types

type RequestCommit struct {
	DeltaMap             interface{}
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type ResponseCommit struct {
	// reserve 1
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	RetainHeight         int64    `protobuf:"varint,3,opt,name=retain_height,json=retainHeight,proto3" json:"retain_height,omitempty"`
	DeltaMap             interface{}
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

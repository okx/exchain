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

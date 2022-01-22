package types


type TreeDelta struct {
	NodesDelta         map[string]*NodeJson `json:"nodes_delta"`
	OrphansDelta       []*NodeJson   `json:"orphans_delta"`
	CommitOrphansDelta map[string]int64     `json:"commit_orphans_delta"`
}

// NodeJson provide json Marshal of Node.
type NodeJson struct {
	Key          []byte `json:"key"`
	Value        []byte `json:"value"`
	Hash         []byte `json:"hash"`
	LeftHash     []byte `json:"left_hash"`
	RightHash    []byte `json:"right_hash"`
	Version      int64  `json:"version"`
	Size         int64  `json:"size"`
	Height       int8   `json:"height"`
	Persisted    bool   `json:"persisted"`
	prePersisted bool   `json:"pre_persisted"`
}

type RequestCommit struct {
	Deltas               *Deltas
	TreeDelta2           TreeDelta
	DeltaMap             interface{}
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}


type ResponseCommit struct {
	// reserve 1
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	RetainHeight         int64    `protobuf:"varint,3,opt,name=retain_height,json=retainHeight,proto3" json:"retain_height,omitempty"`
	Deltas               *Deltas  `protobuf:"bytes,4,opt,name=deltas,proto3" json:"deltas,omitempty"`
	DeltaMap             interface{}
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

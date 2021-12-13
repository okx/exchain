package types

// DeltasMeta contains meta information.
type DeltasMeta struct {
	DeltasSize int `json:"Deltas_size"`
}

//-----------------------------------------------------------
// These methods are for Protobuf Compatibility

// Size returns the size of the amino encoding, in bytes.
func (bm *DeltasMeta) Size() int {
	bs, _ := bm.Marshal()
	return len(bs)
}

// Marshal returns the amino encoding.
func (bm *DeltasMeta) Marshal() ([]byte, error) {
	return cdc.MarshalBinaryBare(bm)
}

// MarshalTo calls Marshal and copies to the given buffer.
func (bm *DeltasMeta) MarshalTo(data []byte) (int, error) {
	bs, err := bm.Marshal()
	if err != nil {
		return -1, err
	}
	return copy(data, bs), nil
}

// Unmarshal deserializes from amino encoded form.
func (bm *DeltasMeta) Unmarshal(bs []byte) error {
	return cdc.UnmarshalBinaryBare(bs, bm)
}

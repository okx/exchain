package exported


// Prefix implements spec:CommitmentPrefix.
// Prefix represents the common "prefix" that a set of keys shares.
type Prefix interface {
	Bytes() []byte
	Empty() bool
}


type Root interface {
	GetHash() []byte
	Empty() bool
}


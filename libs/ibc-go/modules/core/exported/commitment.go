package exported

import ics23 "github.com/confio/ics23/go"

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


// Path implements spec:CommitmentPath.
// A path is the additional information provided to the verification function.
type Path interface {
	String() string
	Empty() bool
}


// Proof implements spec:CommitmentProof.
// Proof can prove whether the key-value pair is a part of the Root or not.
// Each proof has designated key-value pair it is able to prove.
// Proofs includes key but value is provided dynamically at the verification time.
type Proof interface {
	VerifyMembership([]*ics23.ProofSpec, Root, Path, []byte) error
	VerifyNonMembership([]*ics23.ProofSpec, Root, Path) error
	Empty() bool

	ValidateBasic() error
}

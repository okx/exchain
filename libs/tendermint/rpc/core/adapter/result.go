package adapter
//
//import (
//	"bytes"
//	"fmt"
//	gogotypes "github.com/gogo/protobuf/types"
//	"time"
//)
//
//// Commit and Header
//type ResultCommit struct {
//	SignedHeader `json:"signed_header"`
//	CanonicalCommit    bool `json:"canonical"`
//}
//
////-----------------------------------------------------------------------------
//
//// SignedHeader is a header along with the commits that prove it.
//type SignedHeader struct {
//	*Header `json:"header"`
//
//	Commit *Commit `json:"commit"`
//}
//
//
////-----------------------------------------------------------------------------
//
//// Header defines the structure of a Tendermint block header.
//// NOTE: changes to the Header should be duplicated in:
//// - header.Hash()
//// - abci.Header
//// - https://github.com/tendermint/spec/blob/master/spec/blockchain/blockchain.md
//type Header struct {
//	// basic block info
//	Version tmversion.Consensus `json:"version"`
//	ChainID string              `json:"chain_id"`
//	Height  int64               `json:"height"`
//	Time    time.Time           `json:"time"`
//
//	// prev block info
//	LastBlockID BlockID `json:"last_block_id"`
//
//	// hashes of block data
//	LastCommitHash tmbytes.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
//	DataHash       tmbytes.HexBytes `json:"data_hash"`        // transactions
//
//	// hashes from the app output from the prev block
//	ValidatorsHash     tmbytes.HexBytes `json:"validators_hash"`      // validators for the current block
//	NextValidatorsHash tmbytes.HexBytes `json:"next_validators_hash"` // validators for the next block
//	ConsensusHash      tmbytes.HexBytes `json:"consensus_hash"`       // consensus params for current block
//	AppHash            tmbytes.HexBytes `json:"app_hash"`             // state after txs from the previous block
//	// root hash of all results from the txs from the previous block
//	// see `deterministicResponseDeliverTx` to understand which parts of a tx is hashed into here
//	LastResultsHash tmbytes.HexBytes `json:"last_results_hash"`
//
//	// consensus info
//	EvidenceHash    tmbytes.HexBytes `json:"evidence_hash"`    // evidence included in the block
//	ProposerAddress Address          `json:"proposer_address"` // original proposer of the block
//}
//
//// Populate the Header with state-derived data.
//// Call this after MakeBlock to complete the Header.
//func (h *Header) Populate(
//	version tmversion.Consensus, chainID string,
//	timestamp time.Time, lastBlockID BlockID,
//	valHash, nextValHash []byte,
//	consensusHash, appHash, lastResultsHash []byte,
//	proposerAddress Address,
//) {
//	h.Version = version
//	h.ChainID = chainID
//	h.Time = timestamp
//	h.LastBlockID = lastBlockID
//	h.ValidatorsHash = valHash
//	h.NextValidatorsHash = nextValHash
//	h.ConsensusHash = consensusHash
//	h.AppHash = appHash
//	h.LastResultsHash = lastResultsHash
//	h.ProposerAddress = proposerAddress
//}
//
//// ValidateBasic performs stateless validation on a Header returning an error
//// if any validation fails.
////
//// NOTE: Timestamp validation is subtle and handled elsewhere.
//func (h Header) ValidateBasic() error {
//	if h.Version.Block != version.BlockProtocol {
//		return fmt.Errorf("block protocol is incorrect: got: %d, want: %d ", h.Version.Block, version.BlockProtocol)
//	}
//	if len(h.ChainID) > MaxChainIDLen {
//		return fmt.Errorf("chainID is too long; got: %d, max: %d", len(h.ChainID), MaxChainIDLen)
//	}
//
//	if h.Height < 0 {
//		return errors.New("negative Height")
//	} else if h.Height == 0 {
//		return errors.New("zero Height")
//	}
//
//	if err := h.LastBlockID.ValidateBasic(); err != nil {
//		return fmt.Errorf("wrong LastBlockID: %w", err)
//	}
//
//	if err := ValidateHash(h.LastCommitHash); err != nil {
//		return fmt.Errorf("wrong LastCommitHash: %v", err)
//	}
//
//	if err := ValidateHash(h.DataHash); err != nil {
//		return fmt.Errorf("wrong DataHash: %v", err)
//	}
//
//	if err := ValidateHash(h.EvidenceHash); err != nil {
//		return fmt.Errorf("wrong EvidenceHash: %v", err)
//	}
//
//	if len(h.ProposerAddress) != crypto.AddressSize {
//		return fmt.Errorf(
//			"invalid ProposerAddress length; got: %d, expected: %d",
//			len(h.ProposerAddress), crypto.AddressSize,
//		)
//	}
//
//	// Basic validation of hashes related to application data.
//	// Will validate fully against state in state#ValidateBlock.
//	if err := ValidateHash(h.ValidatorsHash); err != nil {
//		return fmt.Errorf("wrong ValidatorsHash: %v", err)
//	}
//	if err := ValidateHash(h.NextValidatorsHash); err != nil {
//		return fmt.Errorf("wrong NextValidatorsHash: %v", err)
//	}
//	if err := ValidateHash(h.ConsensusHash); err != nil {
//		return fmt.Errorf("wrong ConsensusHash: %v", err)
//	}
//	// NOTE: AppHash is arbitrary length
//	if err := ValidateHash(h.LastResultsHash); err != nil {
//		return fmt.Errorf("wrong LastResultsHash: %v", err)
//	}
//
//	return nil
//}
//
//// Hash returns the hash of the header.
//// It computes a Merkle tree from the header fields
//// ordered as they appear in the Header.
//// Returns nil if ValidatorHash is missing,
//// since a Header is not valid unless there is
//// a ValidatorsHash (corresponding to the validator set).
//func (h *Header) Hash() tmbytes.HexBytes {
//	if h == nil || len(h.ValidatorsHash) == 0 {
//		return nil
//	}
//	hbz, err := h.Version.Marshal()
//	if err != nil {
//		return nil
//	}
//
//	pbt, err := gogotypes.StdTimeMarshal(h.Time)
//	if err != nil {
//		return nil
//	}
//
//	pbbi := h.LastBlockID.ToProto()
//	bzbi, err := pbbi.Marshal()
//	if err != nil {
//		return nil
//	}
//	return merkle.HashFromByteSlices([][]byte{
//		hbz,
//		cdcEncode(h.ChainID),
//		cdcEncode(h.Height),
//		pbt,
//		bzbi,
//		cdcEncode(h.LastCommitHash),
//		cdcEncode(h.DataHash),
//		cdcEncode(h.ValidatorsHash),
//		cdcEncode(h.NextValidatorsHash),
//		cdcEncode(h.ConsensusHash),
//		cdcEncode(h.AppHash),
//		cdcEncode(h.LastResultsHash),
//		cdcEncode(h.EvidenceHash),
//		cdcEncode(h.ProposerAddress),
//	})
//}
//
//// ValidateBasic does basic consistency checks and makes sure the header
//// and commit are consistent.
////
//// NOTE: This does not actually check the cryptographic signatures.  Make sure
//// to use a Verifier to validate the signatures actually provide a
//// significantly strong proof for this header's validity.
//func (sh SignedHeader) ValidateBasic(chainID string) error {
//	if sh.Header == nil {
//		return errors.New("missing header")
//	}
//	if sh.Commit == nil {
//		return errors.New("missing commit")
//	}
//
//	if err := sh.Header.ValidateBasic(); err != nil {
//		return fmt.Errorf("invalid header: %w", err)
//	}
//	if err := sh.Commit.ValidateBasic(); err != nil {
//		return fmt.Errorf("invalid commit: %w", err)
//	}
//
//	if sh.ChainID != chainID {
//		return fmt.Errorf("header belongs to another chain %q, not %q", sh.ChainID, chainID)
//	}
//
//	// Make sure the header is consistent with the commit.
//	if sh.Commit.Height != sh.Height {
//		return fmt.Errorf("header and commit height mismatch: %d vs %d", sh.Height, sh.Commit.Height)
//	}
//	if hhash, chash := sh.Header.Hash(), sh.Commit.BlockID.Hash; !bytes.Equal(hhash, chash) {
//		return fmt.Errorf("commit signs block %X, header is block %X", chash, hhash)
//	}
//
//	return nil
//}
//
//// String returns a string representation of SignedHeader.
//func (sh SignedHeader) String() string {
//	return sh.StringIndented("")
//}
//
//// StringIndented returns an indented string representation of SignedHeader.
////
//// Header
//// Commit
//func (sh SignedHeader) StringIndented(indent string) string {
//	return fmt.Sprintf(`SignedHeader{
//%s  %v
//%s  %v
//%s}`,
//		indent, sh.Header.StringIndented(indent+"  "),
//		indent, sh.Commit.StringIndented(indent+"  "),
//		indent)
//}
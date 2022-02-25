package types

import (
	gobytes "bytes"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/bytes"

	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	"github.com/tendermint/go-amino"
)

//-----------------------------------------------------------------------------

// ABCIResult is the deterministic component of a ResponseDeliverTx.
// TODO: add tags and other fields
// https://github.com/tendermint/tendermint/issues/1007
type ABCIResult struct {
	Code uint32         `json:"code"`
	Data bytes.HexBytes `json:"data"`
}

func (a ABCIResult) AminoSize() int {
	var size int
	if a.Code != 0 {
		size += 1 + amino.UvarintSize(uint64(a.Code))
	}
	if len(a.Data) > 0 {
		size += 1 + amino.UvarintSize(uint64(len(a.Data))) + len(a.Data)
	}
	return size
}

func (a ABCIResult) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	buf := &gobytes.Buffer{}
	buf.Grow(a.AminoSize())

	if a.Code != 0 {
		const pbKey = 1<<3 | byte(amino.Typ3_Varint)
		err := amino.EncodeUvarintWithKeyToBuffer(buf, uint64(a.Code), pbKey)
		if err != nil {
			return nil, err
		}
	}

	if len(a.Data) != 0 {
		const pbKey = 2<<3 | byte(amino.Typ3_ByteLength)
		err := amino.EncodeByteSliceWithKeyToBuffer(buf, a.Data, pbKey)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Bytes returns the amino encoded ABCIResult
func (a ABCIResult) Bytes() []byte {
	bz, err := a.MarshalToAmino(cdc)
	if err != nil {
		return cdcEncode(a)
	} else {
		return bz
	}
}

// ABCIResults wraps the deliver tx results to return a proof
type ABCIResults []ABCIResult

// NewResults creates ABCIResults from the list of ResponseDeliverTx.
func NewResults(responses []*abci.ResponseDeliverTx) ABCIResults {
	res := make(ABCIResults, len(responses))
	for i, d := range responses {
		res[i] = NewResultFromResponse(d)
	}
	return res
}

// NewResultFromResponse creates ABCIResult from ResponseDeliverTx.
func NewResultFromResponse(response *abci.ResponseDeliverTx) ABCIResult {
	return ABCIResult{
		Code: response.Code,
		Data: response.Data,
	}
}

// Bytes serializes the ABCIResponse using amino
func (a ABCIResults) Bytes() []byte {
	bz, err := cdc.MarshalBinaryLengthPrefixed(a)
	if err != nil {
		panic(err)
	}
	return bz
}

// Hash returns a merkle hash of all results
func (a ABCIResults) Hash() []byte {
	// NOTE: we copy the impl of the merkle tree for txs -
	// we should be consistent and either do it for both or not.
	return merkle.SimpleHashFromByteSlices(a.toByteSlices())
}

// ProveResult returns a merkle proof of one result from the set
func (a ABCIResults) ProveResult(i int) merkle.SimpleProof {
	_, proofs := merkle.SimpleProofsFromByteSlices(a.toByteSlices())
	return *proofs[i]
}

func (a ABCIResults) toByteSlices() [][]byte {
	l := len(a)
	bzs := make([][]byte, l)
	for i := 0; i < l; i++ {
		bzs[i] = a[i].Bytes()
	}
	return bzs
}

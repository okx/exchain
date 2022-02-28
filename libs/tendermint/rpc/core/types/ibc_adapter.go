package coretypes

import (
	"github.com/okex/exchain/libs/tendermint/types"
)

// Commit and Header
type IBCResultCommit struct {
	types.IBCSignedHeader `json:"signed_header"`

	CanonicalCommit    bool `json:"canonical"`
}
func NewIBCResultCommit(header *types.IBCHeader, commit *types.IBCCommit,
	canonical bool) *IBCResultCommit {

	return &IBCResultCommit{
		IBCSignedHeader: types.IBCSignedHeader{
			IBCHeader: header,
			Commit: commit,
		},
		CanonicalCommit: canonical,
	}
}
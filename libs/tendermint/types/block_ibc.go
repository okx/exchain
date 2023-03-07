package types

import (
	tmbytes "github.com/okx/okbchain/libs/tendermint/libs/bytes"
)

func (sh SignedHeader) ValidateBasicForIBC(chainID string) error {
	return sh.commonValidateBasic(chainID, true)
}

func (h *Header) PureIBCHash() tmbytes.HexBytes {
	if h == nil || len(h.ValidatorsHash) == 0 {
		return nil
	}
	return h.IBCHash()
}

func (commit *Commit) IBCVoteSignBytes(chainID string, valIdx int) []byte {
	return commit.GetVote(valIdx).ibcSignBytes(chainID)
}

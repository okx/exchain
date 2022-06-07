package types

import (
	"bytes"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
)

func (vals *ValidatorSet) IBCVerifyCommitLightTrusting(chainID string, blockID BlockID,
	height int64, commit *Commit, trustLevel tmmath.Fraction) error {

	return vals.commonVerifyCommitLightTrusting(chainID, blockID, height, commit, trustLevel, true)
}

func (vals *ValidatorSet) IBCGetByAddress(address []byte) (index int, val *Validator) {
	for idx, val := range vals.Validators {
		if bytes.Equal(val.Address, address) {
			return idx, val.Copy()
		}
	}
	return -1, nil
}

func (vals *ValidatorSet) IBCHash() []byte {
	if len(vals.Validators) == 0 {
		return nil
	}
	bzs := make([][]byte, len(vals.Validators))
	for i, val := range vals.Validators {
		bzs[i] = val.IBCHeightBytes()
	}
	return merkle.SimpleHashFromByteSlices(bzs)
}

func (vals *ValidatorSet) IBCVerifyCommitLight(chainID string, blockID BlockID,
	height int64, commit *Commit) error {

	return vals.commonVerifyCommitLight(chainID, blockID, height, commit, true)
}

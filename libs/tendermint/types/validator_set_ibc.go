package types

import (
	"bytes"

	"github.com/okx/okbchain/libs/tendermint/crypto/merkle"
	tmmath "github.com/okx/okbchain/libs/tendermint/libs/math"
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

//-------------------------------------

// ValidatorsByVotingPower implements sort.Interface for []*Validator based on
// the VotingPower and Address fields.
type ValidatorsByVotingPower []*Validator

func (valz ValidatorsByVotingPower) Len() int { return len(valz) }

func (valz ValidatorsByVotingPower) Less(i, j int) bool {
	if valz[i].VotingPower == valz[j].VotingPower {
		return bytes.Compare(valz[i].Address, valz[j].Address) == -1
	}
	return valz[i].VotingPower > valz[j].VotingPower
}

func (valz ValidatorsByVotingPower) Swap(i, j int) {
	valz[i], valz[j] = valz[j], valz[i]
}

package types

import (
	"bytes"
	"fmt"
	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
	"github.com/pkg/errors"
)

func (vals *ValidatorSet) IBCVerifyCommitLightTrusting(chainID string, blockID BlockID,
	height int64, commit *Commit, trustLevel tmmath.Fraction) error {

	// sanity check
	if trustLevel.Numerator*3 < trustLevel.Denominator || // < 1/3
		trustLevel.Numerator > trustLevel.Denominator { // > 1
		panic(fmt.Sprintf("trustLevel must be within [1/3, 1], given %v", trustLevel))
	}

	if err := verifyCommitBasic(commit, height, blockID); err != nil {
		return err
	}

	var (
		talliedVotingPower int64
		seenVals           = make(map[int]int, len(commit.Signatures)) // validator index -> commit index
	)

	// Safely calculate voting power needed.
	totalVotingPowerMulByNumerator, overflow := safeMul(vals.TotalVotingPower(), trustLevel.Numerator)
	if overflow {
		return errors.New("int64 overflow while calculating voting power needed. please provide smaller trustLevel numerator")
	}
	votingPowerNeeded := totalVotingPowerMulByNumerator / trustLevel.Denominator

	for idx, commitSig := range commit.Signatures {
		// No need to verify absent or nil votes.
		if !commitSig.ForBlock() {
			continue
		}

		// We don't know the validators that committed this block, so we have to
		// check for each vote if its validator is already known.
		valIdx, val := vals.IBCGetByAddress(commitSig.ValidatorAddress)

		if val != nil {
			// check for double vote of validator on the same commit
			if firstIndex, ok := seenVals[valIdx]; ok {
				secondIndex := idx
				return errors.Errorf("double vote from %v (%d and %d)", val, firstIndex, secondIndex)
			}
			seenVals[valIdx] = idx

			// Validate signature.
			voteSignBytes := commit.VoteSignBytes(chainID, idx)
			if !val.PubKey.VerifyBytes(voteSignBytes, commitSig.Signature) {
				return errors.Errorf("wrong signature (#%d): %X", idx, commitSig.Signature)
			}

			talliedVotingPower += val.VotingPower

			if talliedVotingPower > votingPowerNeeded {
				return nil
			}
		}
	}

	return ErrNotEnoughVotingPowerSigned{Got: talliedVotingPower, Needed: votingPowerNeeded}
}

func (vals *ValidatorSet) IBCGetByAddress(address []byte) (index int, val *Validator) {
	for idx, val := range vals.Validators {
		if bytes.Equal(val.Address, address) {
			return idx, val.Copy()
		}
	}
	return -1, nil
}

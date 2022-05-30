package lite

import (
	"bytes"
	"errors"
	errpkg "github.com/pkg/errors"
	"time"

	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
	"github.com/okex/exchain/libs/tendermint/types"
)

func commonVerifyAdjacent(
	chainID string,
	trustedHeader *types.SignedHeader, // height=X
	untrustedHeader *types.SignedHeader, // height=X+1
	untrustedVals *types.ValidatorSet, // height=X+1
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration, isIbc bool) error {

	if untrustedHeader.Height != trustedHeader.Height+1 {
		return errors.New("headers must be adjacent in height")
	}

	if HeaderExpired(trustedHeader, trustingPeriod, now) {
		return ErrOldHeaderExpired{trustedHeader.Time.Add(trustingPeriod), now}
	}

	if err := commonVerifyNewHeaderAndVals(
		chainID,
		untrustedHeader, untrustedVals,
		trustedHeader,
		now, maxClockDrift, isIbc); err != nil {
		return ErrInvalidHeader{err}
	}

	// Check the validator hashes are the same
	if !bytes.Equal(untrustedHeader.ValidatorsHash, trustedHeader.NextValidatorsHash) {
		err := errpkg.Errorf("expected old header next validators (%X) to match those from new header (%X)",
			trustedHeader.NextValidatorsHash,
			untrustedHeader.ValidatorsHash,
		)
		return err
	}

	// Ensure that +2/3 of new validators signed correctly.
	if err := untrustedVals.CommonVerifyCommitLight(chainID, untrustedHeader.Commit.BlockID, untrustedHeader.Height,
		untrustedHeader.Commit, isIbc); err != nil {
		return ErrInvalidHeader{err}
	}

	return nil
}

func commonVerifyNonAdjacent(
	chainID string,
	trustedHeader *types.SignedHeader, // height=X
	trustedVals *types.ValidatorSet, // height=X or height=X+1
	untrustedHeader *types.SignedHeader, // height=Y
	untrustedVals *types.ValidatorSet, // height=Y
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration,
	trustLevel tmmath.Fraction, isIbc bool) error {

	if untrustedHeader.Height == trustedHeader.Height+1 {
		return errors.New("headers must be non adjacent in height")
	}

	if HeaderExpired(trustedHeader, trustingPeriod, now) {
		return ErrOldHeaderExpired{trustedHeader.Time.Add(trustingPeriod), now}
	}

	if err := commonVerifyNewHeaderAndVals(
		chainID,
		untrustedHeader, untrustedVals,
		trustedHeader,
		now, maxClockDrift, isIbc); err != nil {
		return ErrInvalidHeader{err}
	}

	// Ensure that +`trustLevel` (default 1/3) or more of last trusted validators signed correctly.
	err := trustedVals.CommonVerifyCommitLightTrusting(chainID, untrustedHeader.Commit.BlockID, untrustedHeader.Height,
		untrustedHeader.Commit, trustLevel, isIbc)
	if err != nil {
		switch e := err.(type) {
		case types.ErrNotEnoughVotingPowerSigned:
			return ErrNewValSetCantBeTrusted{e}
		default:
			return e
		}
	}

	// Ensure that +2/3 of new validators signed correctly.
	//
	// NOTE: this should always be the last check because untrustedVals can be
	// intentionally made very large to DOS the light client. not the case for
	// VerifyAdjacent, where validator set is known in advance.
	if err := untrustedVals.CommonVerifyCommitLight(chainID, untrustedHeader.Commit.BlockID, untrustedHeader.Height,
		untrustedHeader.Commit, isIbc); err != nil {
		return ErrInvalidHeader{err}
	}

	return nil
}

func commonVerifyNewHeaderAndVals(chainID string,
	untrustedHeader *types.SignedHeader,
	untrustedVals *types.ValidatorSet,
	trustedHeader *types.SignedHeader,
	now time.Time,
	maxClockDrift time.Duration, isIbc bool) error {

	if isIbc {
		if err := untrustedHeader.ValidateBasicForIBC(chainID); err != nil {
			return errpkg.Wrap(err, "untrustedHeader.ValidateBasic failed")
		}
	} else {
		if err := untrustedHeader.ValidateBasic(chainID); err != nil {
			return errpkg.Wrap(err, "untrustedHeader.ValidateBasic failed")
		}
	}

	if untrustedHeader.Height <= trustedHeader.Height {
		return errpkg.Errorf("expected new header height %d to be greater than one of old header %d",
			untrustedHeader.Height,
			trustedHeader.Height)
	}

	if !untrustedHeader.Time.After(trustedHeader.Time) {
		return errpkg.Errorf("expected new header time %v to be after old header time %v",
			untrustedHeader.Time,
			trustedHeader.Time)
	}

	if !untrustedHeader.Time.Before(now.Add(maxClockDrift)) {
		return errpkg.Errorf("new header has a time from the future %v (now: %v; max clock drift: %v)",
			untrustedHeader.Time,
			now,
			maxClockDrift)
	}
	if isIbc {
		if !bytes.Equal(untrustedHeader.ValidatorsHash, untrustedVals.IBCHash()) {
			return errpkg.Errorf("expected new header validators (%X) to match those that were supplied (%X) at height %d",
				untrustedHeader.ValidatorsHash,
				untrustedVals.Hash(untrustedHeader.Height),
				untrustedHeader.Height,
			)
		}
	} else {
		if !bytes.Equal(untrustedHeader.ValidatorsHash, untrustedVals.Hash(untrustedHeader.Height)) {
			return errpkg.Errorf("expected new header validators (%X) to match those that were supplied (%X) at height %d",
				untrustedHeader.ValidatorsHash,
				untrustedVals.Hash(untrustedHeader.Height),
				untrustedHeader.Height,
			)
		}
	}

	return nil
}

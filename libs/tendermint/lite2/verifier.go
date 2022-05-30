package lite

import (
	"bytes"
	"time"

	"github.com/pkg/errors"

	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
	"github.com/okex/exchain/libs/tendermint/types"
)

var (
	// DefaultTrustLevel - new header can be trusted if at least one correct
	// validator signed it.
	DefaultTrustLevel = tmmath.Fraction{Numerator: 1, Denominator: 3}
)

// VerifyNonAdjacent verifies non-adjacent untrustedHeader against
// trustedHeader. It ensures that:
//
//	a) trustedHeader can still be trusted (if not, ErrOldHeaderExpired is returned)
//	b) untrustedHeader is valid (if not, ErrInvalidHeader is returned)
//	c) trustLevel ([1/3, 1]) of trustedHeaderVals (or trustedHeaderNextVals)
//  signed correctly (if not, ErrNewValSetCantBeTrusted is returned)
//	d) more than 2/3 of untrustedVals have signed h2
//    (otherwise, ErrInvalidHeader is returned)
//  e) headers are non-adjacent.
//
// maxClockDrift defines how much untrustedHeader.Time can drift into the
// future.
func VerifyNonAdjacent(
	chainID string,
	trustedHeader *types.SignedHeader, // height=X
	trustedVals *types.ValidatorSet, // height=X or height=X+1
	untrustedHeader *types.SignedHeader, // height=Y
	untrustedVals *types.ValidatorSet, // height=Y
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration,
	trustLevel tmmath.Fraction) error {
	return commonVerifyNonAdjacent(chainID,
		trustedHeader,
		trustedVals,
		untrustedHeader,
		untrustedVals,
		trustingPeriod,
		now,
		maxClockDrift,
		trustLevel,
		false)
}

// VerifyAdjacent verifies directly adjacent untrustedHeader against
// trustedHeader. It ensures that:
//
//  a) trustedHeader can still be trusted (if not, ErrOldHeaderExpired is returned)
//  b) untrustedHeader is valid (if not, ErrInvalidHeader is returned)
//  c) untrustedHeader.ValidatorsHash equals trustedHeader.NextValidatorsHash
//  d) more than 2/3 of new validators (untrustedVals) have signed h2
//    (otherwise, ErrInvalidHeader is returned)
//  e) headers are adjacent.
//
// maxClockDrift defines how much untrustedHeader.Time can drift into the
// future.
func VerifyAdjacent(
	chainID string,
	trustedHeader *types.SignedHeader, // height=X
	untrustedHeader *types.SignedHeader, // height=X+1
	untrustedVals *types.ValidatorSet, // height=X+1
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration) error {
	return commonVerifyAdjacent(chainID,
		trustedHeader,
		untrustedHeader,
		untrustedVals,
		trustingPeriod,
		now,
		maxClockDrift,
		false)
}

// Verify combines both VerifyAdjacent and VerifyNonAdjacent functions.
func Verify(
	chainID string,
	trustedHeader *types.SignedHeader, // height=X
	trustedVals *types.ValidatorSet, // height=X or height=X+1
	untrustedHeader *types.SignedHeader, // height=Y
	untrustedVals *types.ValidatorSet, // height=Y
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration,
	trustLevel tmmath.Fraction) error {

	if untrustedHeader.Height != trustedHeader.Height+1 {
		return VerifyNonAdjacent(chainID, trustedHeader, trustedVals, untrustedHeader, untrustedVals,
			trustingPeriod, now, maxClockDrift, trustLevel)
	}

	return VerifyAdjacent(chainID, trustedHeader, untrustedHeader, untrustedVals, trustingPeriod, now, maxClockDrift)
}

func verifyNewHeaderAndVals(
	chainID string,
	untrustedHeader *types.SignedHeader,
	untrustedVals *types.ValidatorSet,
	trustedHeader *types.SignedHeader,
	now time.Time,
	maxClockDrift time.Duration) error {

	if err := untrustedHeader.ValidateBasic(chainID); err != nil {
		return errors.Wrap(err, "untrustedHeader.ValidateBasic failed")
	}

	if untrustedHeader.Height <= trustedHeader.Height {
		return errors.Errorf("expected new header height %d to be greater than one of old header %d",
			untrustedHeader.Height,
			trustedHeader.Height)
	}

	if !untrustedHeader.Time.After(trustedHeader.Time) {
		return errors.Errorf("expected new header time %v to be after old header time %v",
			untrustedHeader.Time,
			trustedHeader.Time)
	}

	if !untrustedHeader.Time.Before(now.Add(maxClockDrift)) {
		return errors.Errorf("new header has a time from the future %v (now: %v; max clock drift: %v)",
			untrustedHeader.Time,
			now,
			maxClockDrift)
	}

	if !bytes.Equal(untrustedHeader.ValidatorsHash, untrustedVals.Hash(untrustedHeader.Height)) {
		return errors.Errorf("expected new header validators (%X) to match those that were supplied (%X) at height %d",
			untrustedHeader.ValidatorsHash,
			untrustedVals.Hash(untrustedHeader.Height),
			untrustedHeader.Height,
		)
	}

	return nil
}

// ValidateTrustLevel checks that trustLevel is within the allowed range [1/3,
// 1]. If not, it returns an error. 1/3 is the minimum amount of trust needed
// which does not break the security model.
func ValidateTrustLevel(lvl tmmath.Fraction) error {
	if lvl.Numerator*3 < lvl.Denominator || // < 1/3
		lvl.Numerator > lvl.Denominator || // > 1
		lvl.Denominator == 0 {
		return errors.Errorf("trustLevel must be within [1/3, 1], given %v", lvl)
	}
	return nil
}

// HeaderExpired return true if the given header expired.
func HeaderExpired(h *types.SignedHeader, trustingPeriod time.Duration, now time.Time) bool {
	expirationTime := h.Time.Add(trustingPeriod)
	return !expirationTime.After(now)
}

// VerifyBackwards verifies an untrusted header with a height one less than
// that of an adjacent trusted header. It ensures that:
//
// 	a) untrusted header is valid
//  b) untrusted header has a time before the trusted header
//  c) that the LastBlockID hash of the trusted header is the same as the hash
//  of the trusted header
//
//  For any of these cases ErrInvalidHeader is returned.
func VerifyBackwards(chainID string, untrustedHeader, trustedHeader *types.SignedHeader) error {
	if err := untrustedHeader.ValidateBasic(chainID); err != nil {
		return ErrInvalidHeader{err}
	}

	if !untrustedHeader.Time.Before(trustedHeader.Time) {
		return ErrInvalidHeader{
			errors.Errorf("expected older header time %v to be before new header time %v",
				untrustedHeader.Time,
				trustedHeader.Time)}
	}

	if !bytes.Equal(untrustedHeader.Hash(), trustedHeader.LastBlockID.Hash) {
		return ErrInvalidHeader{
			errors.Errorf("older header hash %X does not match trusted header's last block %X",
				untrustedHeader.Hash(),
				trustedHeader.LastBlockID.Hash)}
	}

	return nil
}

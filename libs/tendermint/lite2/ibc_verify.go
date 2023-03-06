package lite

import (
	"time"

	tmmath "github.com/okx/okbchain/libs/tendermint/libs/math"
	"github.com/okx/okbchain/libs/tendermint/types"
)

func IBCVerify(
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
		return IBCVerifyNonAdjacent(chainID, trustedHeader, trustedVals, untrustedHeader, untrustedVals,
			trustingPeriod, now, maxClockDrift, trustLevel)
	}

	return IBCVerifyAdjacent(chainID, trustedHeader, untrustedHeader, untrustedVals, trustingPeriod, now, maxClockDrift)
}

func IBCVerifyNonAdjacent(
	chainID string,
	trustedHeader *types.SignedHeader, // height=X
	trustedVals *types.ValidatorSet, // height=X or height=X+1
	untrustedHeader *types.SignedHeader, // height=Y
	untrustedVals *types.ValidatorSet, // height=Y
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration,
	trustLevel tmmath.Fraction) error {

	return commonVerifyNonAdjacent(
		chainID, trustedHeader, trustedVals, untrustedHeader,
		untrustedVals, trustingPeriod, now, maxClockDrift, trustLevel, true)
}

func IBCVerifyAdjacent(
	chainID string,
	trustedHeader *types.SignedHeader, // height=X
	untrustedHeader *types.SignedHeader, // height=X+1
	untrustedVals *types.ValidatorSet, // height=X+1
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration) error {

	return commonVerifyAdjacent(
		chainID,
		trustedHeader,   // height=X
		untrustedHeader, // height=X+1
		untrustedVals,   // height=X+1
		trustingPeriod,
		now,
		maxClockDrift, true)
}

package types

import (
	"strconv"
	"sync"
)

// Disable followings after milestoneMercuryHeight
// 1. TransferToContractBlock
// 2. ChangeEvmDenomByProposal
// 3. BankTransferBlock
// 4. ibc

var (
	MILESTONE_GENESIS_HEIGHT string
	genesisHeight            int64

	MILESTONE_MERCURY_HEIGHT string
	milestoneMercuryHeight   int64

	MILESTONE_VENUS_HEIGHT string
	milestoneVenusHeight   int64

	MILESTONE_MARS_HEIGHT string
	milestoneMarsHeight   int64

	MILESTONE_VENUS1_HEIGHT string
	milestoneVenus1Height   int64

	MILESTONE_VENUS2_HEIGHT string
	milestoneVenus2Height   int64

	once sync.Once
)

func init() {
	once.Do(func() {
		genesisHeight = string2number(MILESTONE_GENESIS_HEIGHT)
		milestoneMercuryHeight = string2number(MILESTONE_MERCURY_HEIGHT)
		milestoneVenusHeight = string2number(MILESTONE_VENUS_HEIGHT)
		milestoneMarsHeight = string2number(MILESTONE_MARS_HEIGHT)
		milestoneVenus1Height = string2number(MILESTONE_VENUS1_HEIGHT)
		milestoneVenus2Height = string2number(MILESTONE_VENUS2_HEIGHT)
	})
}

func string2number(input string) int64 {
	if len(input) == 0 {
		input = "0"
	}
	res, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}

//depracate homstead signer support
func HigherThanMercury(height int64) bool {
	if milestoneMercuryHeight == 0 {
		// milestoneMercuryHeight not enabled
		return false
	}
	return height > milestoneMercuryHeight
}

func HigherThanVenus(height int64) bool {
	if milestoneVenusHeight == 0 {
		return false
	}
	return height >= milestoneVenusHeight
}

//use MPT storage model to replace IAVL storage model
func HigherThanMars(height int64) bool {
	if milestoneMarsHeight == 0 {
		return false
	}
	return height > milestoneMarsHeight
}

// GetMilestoneVenusHeight returns milestoneVenusHeight
func GetMilestoneVenusHeight() int64 {
	return milestoneVenusHeight
}

// 2322600 is mainnet GenesisHeight
func IsMainNet() bool {
	return MILESTONE_GENESIS_HEIGHT == "2322600"
}

// 1121818 is testnet GenesisHeight
func IsTestNet() bool {
	return MILESTONE_GENESIS_HEIGHT == "1121818"
}

func GetStartBlockHeight() int64 {
	return genesisHeight
}

func GetVenusHeight() int64 {
	return milestoneVenusHeight
}

func GetMercuryHeight() int64 {
	return milestoneMercuryHeight
}

func GetMarsHeight() int64 {
	return milestoneMarsHeight
}

// can be used in unit test only
func UnittestOnlySetMilestoneVenusHeight(height int64) {
	milestoneVenusHeight = height
}

// can be used in unit test only
func UnittestOnlySetMilestoneMarsHeight(height int64) {
	milestoneMarsHeight = height
}

// ==================================
// =========== Venus1 ===============
func HigherThanVenus1(h int64) bool {
	if milestoneVenus1Height == 0 {
		return false
	}
	return h >= milestoneVenus1Height
}

func UnittestOnlySetMilestoneVenus1Height(h int64) {
	milestoneVenus1Height = h
}

func GetVenus1Height() int64 {
	return milestoneVenus1Height
}

// =========== Venus1 ===============
// ==================================

// ==================================
// =========== Venus2 ===============
func HigherThanVenus2(h int64) bool {
	if milestoneVenus2Height == 0 {
		return false
	}
	return h >= milestoneVenus2Height
}

func GetVenus2Height() int64 {
	return milestoneVenus2Height
}

func UnittestOnlySetMilestoneVenus2Height(h int64) {
	milestoneVenus2Height = h
}

// =========== Venus2 ===============
// ==================================

package types

import (
	"math"
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

	MILESTONE_VENUS1_HEIGHT string
	milestoneVenus1Height   int64

	once sync.Once
)

func init() {
	once.Do(func() {
		genesisHeight = string2number(MILESTONE_GENESIS_HEIGHT)
		milestoneMercuryHeight = string2number(MILESTONE_MERCURY_HEIGHT)
		milestoneVenusHeight = string2number(MILESTONE_VENUS_HEIGHT)
		milestoneVenus1Height = string2number(MILESTONE_VENUS1_HEIGHT)
		if milestoneVenus1Height == 0 {
			milestoneVenus1Height = math.MaxInt64 - 2
		}
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

// HigherThanVenus1 will be true if milestoneVenus1Height is larger than current height
func HigherThanVenus1(h int64) bool {
	if milestoneVenus1Height == 0 {
		return false
	}
	return h > milestoneVenus1Height
}
func SetVenus1HeightForTest() {
	milestoneVenus1Height = math.MaxInt64 - 2
}
func SetVenus1HeightForIbcTest(h int64) {
	milestoneVenus1Height = h
}

func GetVenus1Height() int64 {
	return milestoneVenus1Height
}

func IsUpgradeVenus1FeatureInRuntime() bool {
	return milestoneVenus1Height >= 1
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

// can be used in unit test only
func UnittestOnlySetMilestoneVenusHeight(height int64) {
	milestoneVenusHeight = height
}

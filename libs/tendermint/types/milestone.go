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

	MILESTONE_IBC_HEIGHT string
	milestoreIbcHeight   int64

	once sync.Once
)

func init() {
	once.Do(func() {
		genesisHeight = string2number(MILESTONE_GENESIS_HEIGHT)
		milestoneMercuryHeight = string2number(MILESTONE_MERCURY_HEIGHT)
		milestoneVenusHeight = string2number(MILESTONE_VENUS_HEIGHT)
		milestoreIbcHeight = string2number(MILESTONE_IBC_HEIGHT)
		if milestoreIbcHeight == 0 {
			// as default: genesisHeight is zero
			milestoreIbcHeight = genesisHeight + 1
			if IsMainNet() || IsTestNet() {
				milestoreIbcHeight = math.MaxInt64 - 2
			}
		} else {
			if IsMainNet() || IsTestNet() {
				if milestoreIbcHeight < milestoneVenusHeight || milestoreIbcHeight < milestoneMercuryHeight {
					panic("invalid ibc height")
				}
			}
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

func HigherThanIBCHeight(h int64) bool {
	if milestoreIbcHeight == 0 {
		return false
	}
	return h > milestoreIbcHeight
}
func GetIBCHeight() int64 {
	return milestoreIbcHeight
}

func UpgradeIBCInRuntime() bool {
	return milestoreIbcHeight >= 1
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

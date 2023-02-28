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

	MILESTONE_VENUS_HEIGHT string
	milestoneVenusHeight   int64

	MILESTONE_MARS_HEIGHT string
	milestoneMarsHeight   int64

	MILESTONE_VENUS2_HEIGHT string
	milestoneVenus2Height   int64

	milestoneEarthHeight int64

	MILESTONE_VENUS4_HEIGHT string
	milestoneVenus4Height   int64

	// note: it stores the earlies height of the node,and it is used by cli
	nodePruneHeight int64

	once sync.Once
)

const (
	MainNet = "exchain-66"
	TestNet = "exchain-65"

	MILESTONE_EARTH = "earth"
)

const (
	MainNetVeneusHeight = 8200000
	TestNetVeneusHeight = 8510000

	MainNetGenesisHeight = 2322600
	TestNetGenesisHeight = 1121818

	TestNetChangeChainId = 2270901
	TestNetChainName1    = "okexchain-65"
)

func init() {
	once.Do(func() {
		genesisHeight = string2number(MILESTONE_GENESIS_HEIGHT)
		milestoneVenusHeight = string2number(MILESTONE_VENUS_HEIGHT)
		milestoneMarsHeight = string2number(MILESTONE_MARS_HEIGHT)
		milestoneVenus2Height = string2number(MILESTONE_VENUS2_HEIGHT)
		milestoneVenus4Height = string2number(MILESTONE_VENUS4_HEIGHT)
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

func SetupMainNetEnvironment(pruneH int64) {
	milestoneVenusHeight = MainNetVeneusHeight
	genesisHeight = MainNetGenesisHeight
	nodePruneHeight = pruneH
}

func SetupTestNetEnvironment(pruneH int64) {
	milestoneVenusHeight = TestNetVeneusHeight
	genesisHeight = TestNetGenesisHeight
	nodePruneHeight = pruneH
}

func HigherThanVenus(height int64) bool {
	if milestoneVenusHeight == 0 {
		return false
	}
	return height >= milestoneVenusHeight
}

// use MPT storage model to replace IAVL storage model
func HigherThanMars(height int64) bool {
	if milestoneMarsHeight == 0 {
		return false
	}
	return height >= milestoneMarsHeight
}

// GetMilestoneVenusHeight returns milestoneVenusHeight
func GetMilestoneVenusHeight() int64 {
	return milestoneVenusHeight
}

// 2322600 is mainnet GenesisHeight
func IsMainNet() bool {
	//return MILESTONE_GENESIS_HEIGHT == "2322600"
	return false
}

// 1121818 is testnet GenesisHeight
func IsTestNet() bool {
	//return MILESTONE_GENESIS_HEIGHT == "1121818"
	return false
}

func GetStartBlockHeight() int64 {
	return genesisHeight
}

func GetNodePruneHeight() int64 {
	return nodePruneHeight
}

func GetVenusHeight() int64 {
	return milestoneVenusHeight
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
// =========== Venus2 ===============
func HigherThanVenus2(h int64) bool {
	if milestoneVenus2Height == 0 {
		return false
	}
	return h >= milestoneVenus2Height
}

func UnittestOnlySetMilestoneVenus2Height(h int64) {
	milestoneVenus2Height = h
}

func GetVenus2Height() int64 {
	return milestoneVenus2Height
}

// =========== Venus2 ===============
// ==================================

// ==================================
// =========== Earth ===============
func UnittestOnlySetMilestoneEarthHeight(h int64) {
	milestoneEarthHeight = h
}

func SetMilestoneEarthHeight(h int64) {
	milestoneEarthHeight = h
}

func HigherThanEarth(h int64) bool {
	if milestoneEarthHeight == 0 {
		return false
	}
	return h >= milestoneEarthHeight
}

func GetEarthHeight() int64 {
	return milestoneEarthHeight
}

// =========== Earth ===============
// ==================================

// ==================================
// =========== Venus3 ===============
func HigherThanVenus4(h int64) bool {
	if milestoneVenus4Height == 0 {
		return false
	}
	return h > milestoneVenus4Height
}

func UnittestOnlySetMilestoneVenus4Height(h int64) {
	milestoneVenus4Height = h
}

func GetVenus4Height() int64 {
	return milestoneVenus4Height
}

// =========== Venus4 ===============
// ==================================

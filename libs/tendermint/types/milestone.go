package types

import (
	"strconv"
	"sync"
)

// Disable followings after milestoneMercuryHeight
// 1. TransferToContractBlock
// 2. ChangeEvmDenomByProposal
// 3. BankTransferBlock

var (
	MILESTONE_GENESIS_HEIGHT     string
	genesisHeight                int64

	MILESTONE_MERCURY_HEIGHT     string
	milestoneMercuryHeight       int64

	MILESTONE_VENUS_HEIGHT       string
	milestoneVenusHeight         int64

	once                         sync.Once
)

func init() {
	once.Do(func() {
		genesisHeight          = string2number(MILESTONE_GENESIS_HEIGHT)
		milestoneMercuryHeight = string2number(MILESTONE_MERCURY_HEIGHT)
		milestoneVenusHeight   = string2number(MILESTONE_VENUS_HEIGHT)
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
	return height > milestoneVenusHeight
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
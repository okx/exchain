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
	MILESTONE_MERCURY_HEIGHT string
	milestoneMercuryHeight   int64

	MILESTONE_VENUS_HEIGHT string
	milestoneVenusHeight   int64

	once sync.Once
)

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

func initVersionBlockHeight() {
	once.Do(func() {
		milestoneMercuryHeight = string2number(MILESTONE_MERCURY_HEIGHT)
		milestoneVenusHeight = string2number(MILESTONE_VENUS_HEIGHT)
	})
}

func init() {
	initVersionBlockHeight()
}

//depracate homstead signer support
func HigherThanMercury(height int64) bool {
	if milestoneMercuryHeight == 0 {
		// milestoneMercuryHeight not enabled
		return false
	}
	return height > milestoneMercuryHeight
}

//use MPT storage model to replace IAVL storage model
func HigherThanVenus(height int64) bool {
	if milestoneVenusHeight == 0 {
		// milestoneVenusHeight not enabled
		return false
	}
	return height > milestoneVenusHeight
}

////disable transfer tokens to contract address by cli
//func IsDisableTransferToContractBlock(height int64) bool {
//	return higherThanMercury(height)
//}
//
////disable change the param EvmDenom by proposal
//func IsDisableChangeEvmDenomByProposal(height int64) bool {
//	return higherThanMercury(height)
//}
//
////disable transfer tokens by module of cosmos-sdk/bank
//func IsDisableBankTransferBlock(height int64) bool {
//	return higherThanMercury(height)
//}

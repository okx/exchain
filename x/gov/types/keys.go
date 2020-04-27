package types

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// PrefixWaitingProposalQueue defines the prefix of waiting proposal queue
	PrefixWaitingProposalQueue = []byte{0x30}
)

// WaitingProposalByBlockHeightKey gets the waiting proposal queue key by block height
func WaitingProposalByBlockHeightKey(blockHeight uint64) []byte {
	return append(PrefixWaitingProposalQueue, sdk.Uint64ToBigEndian(blockHeight)...)
}

// WaitingProposalQueueKey returns the key for a proposalID in the WaitingProposalQueue
func WaitingProposalQueueKey(proposalID uint64, blockHeight uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)

	return append(WaitingProposalByBlockHeightKey(blockHeight), bz...)
}

// SplitWaitingProposalQueueKey split the active proposal key and returns the proposal id and endTime
func SplitWaitingProposalQueueKey(key []byte) (proposalID uint64, height uint64) {
	return splitKeyWithHeight(key)
}

func splitKeyWithHeight(key []byte) (proposalID uint64, height uint64) {
	// 16 is sum of proposalID length and height length
	if len(key[1:]) != 16 {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), 16))
	}

	height = binary.BigEndian.Uint64(key[1 : 1+8])
	proposalID = binary.LittleEndian.Uint64(key[1+8:])
	return
}

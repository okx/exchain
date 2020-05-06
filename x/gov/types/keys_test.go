package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitWaitingProposalQueueKey(t *testing.T) {
	var expectedProposalID uint64 = 1
	var expectedBlockHeight uint64 = 100
	keyBytes := WaitingProposalQueueKey(expectedProposalID, expectedBlockHeight)
	proposalID, height := SplitWaitingProposalQueueKey(keyBytes)
	require.Equal(t, expectedProposalID, proposalID)
	require.Equal(t, expectedBlockHeight, height)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	proposalID, height = SplitWaitingProposalQueueKey(keyBytes[1:])
}

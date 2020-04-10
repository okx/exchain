package distribution

import (
	"testing"

	"github.com/okex/okchain/x/distribution/keeper"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestBeginBlocker(t *testing.T) {
	_, valConsPks, valConsAddrs := keeper.GetTestAddrs()
	ctx, _, k, _, _ := keeper.CreateTestInputDefault(t, false, 1000)

	for i := int64(1); i < 10; i++ {
		ctx = ctx.WithBlockHeight(i)
		index := i % int64(len(valConsAddrs))
		votes := []abci.VoteInfo{
			{Validator: abci.Validator{Address: valConsPks[index].Address(), Power: 1}, SignedLastBlock: true},
		}
		req := abci.RequestBeginBlock{Header: abci.Header{Height: i, ProposerAddress: valConsAddrs[index].Bytes()},
			LastCommitInfo: abci.LastCommitInfo{Votes: votes}}
		BeginBlocker(ctx, req, k)
		require.Equal(t, k.GetPreviousProposerConsAddr(ctx), valConsAddrs[index])
	}
}

package distribution

import (
	"testing"

	"github.com/okex/okchain/x/distribution/keeper"
	"github.com/okex/okchain/x/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestBeginBlocker(t *testing.T) {
	valOpAddrs, valConsPks, valConsAddrs := keeper.GetAddrs()
	ctx, _, k, sk, _ := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)
	skMsg := staking.NewMsgCreateValidator(valOpAddrs[0], valConsPks[0], staking.Description{}, keeper.NewDecCoin(1))
	require.True(t, sh(ctx, skMsg).IsOK())

	votes := []abci.VoteInfo{
		{Validator: abci.Validator{Address: valConsPks[0].Address(), Power: 1}, SignedLastBlock: true},
	}

	ctx = ctx.WithBlockHeight(1)
	req := abci.RequestBeginBlock{Header: abci.Header{Height: 1, ProposerAddress: valConsAddrs[0].Bytes()},
		LastCommitInfo: abci.LastCommitInfo{Votes: votes}}
	BeginBlocker(ctx, req, k)
	require.Equal(t, k.GetPreviousProposerConsAddr(ctx), valConsAddrs[0])
	ctx = ctx.WithBlockHeight(2)
	req = abci.RequestBeginBlock{Header: abci.Header{Height: 2, ProposerAddress: valConsAddrs[0].Bytes()},
		LastCommitInfo: abci.LastCommitInfo{Votes: votes}}
	BeginBlocker(ctx, req, k)
	require.Equal(t, k.GetPreviousProposerConsAddr(ctx), valConsAddrs[0])
}

package temp

import (
	"encoding/json"
	"github.com/okex/exchain/x/staking/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgDeposit(t *testing.T) {
	//msgstr := `{"delegator_address":"cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity":{"denom":"okt","amount":"10"}}`

	msgstr := `{"delegator_address":"cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity":{"denom":"okt","amount":"10"}}`
	msg := &types.MsgDeposit{}
	err := json.Unmarshal([]byte(msgstr), msg)
	require.NoError(t, err)

	_, err = types.GenDepositMsg(nil, nil)
	require.NoError(t, err)
}

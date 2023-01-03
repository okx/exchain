package expired

import (
	"github.com/okex/exchain/x/evm/statistics/rediscli"
	"testing"
)

func Test_checkMintRewardIfNotEqualReturnTheLatestMint(t *testing.T) {
	rediscli.GetInstance().Init()
	t.Log(checkMintRewardIfNotEqualReturnTheLatestMint("c0x1826080876d1dfbb06aa4f722876fec7b243b59c"))
}

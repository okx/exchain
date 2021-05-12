package simulation

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

func TestEvmFactory(t *testing.T) {
	ef := EvmFactory{ChainId: "ok-1"}

	ef.BuildSimulator().DoCall(types.MsgEthermint{
		AccountNonce: 0,
		Price:        sdk.NewInt(100000),
		GasLimit:     30000000,
		Recipient:    nil,
		Amount:       sdk.NewInt(100),
		Payload:      nil,
		From:         nil,
	})
}

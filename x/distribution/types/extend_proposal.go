package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
import "encoding/json"

const (
	InvokeExtendProposalName = "InvokeExtendProposal"
)

const (
	MethodHello = "Hello"
)

type TestStatus struct {
	Passed bool
}

type TestExtend struct {
	Status TestStatus   `json:"status" yaml:"status"`
	Coins  sdk.SysCoins `json:"coins" yaml:"coins"`
}

func (extend TestExtend) ValidateBasic() error {
	if extend.Coins.Empty() || extend.Coins.AmountOf(sdk.DefaultBondDenom).IsZero() {
		return ErrExtendProposalParams("test extend proposal coins error")
	}

	return nil
}

func NewTestExtend(jsonData string) (TestExtend, error) {
	var extend TestExtend
	err := json.Unmarshal([]byte(jsonData), &extend)
	return extend, err
}

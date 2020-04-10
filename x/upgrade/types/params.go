package types

import (
	"fmt"
	"time"

	"github.com/okex/okchain/x/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DefaultParamspace = ModuleName

var (
	KeyAppUpgradeMaxDepositPeriod = []byte("AppUpgradeMaxDepositPeriod")
	KeyAppUpgradeMinDeposit       = []byte("AppUpgradeMinDeposit")
	KeyAppUpgradeVotingPeriod     = []byte("AppUpgradeVotingPeriod")
)

// UpgradeParams parameters
type UpgradeParams struct {
	AppUpgradeMaxDepositPeriod time.Duration `json:"app_upgrade_max_deposit_period"` //  Maximum period for okb holders to deposit on a AppUpgrade proposal. Initial value: 2 days
	AppUpgradeMinDeposit       sdk.DecCoins  `json:"app_upgrade_min_deposit"`        //  Minimum deposit for a critical AppUpgrade proposal to enter voting period.
	AppUpgradeVotingPeriod     time.Duration `json:"app_upgrade_voting_period"`      //  Length of the critical voting period for AppUpgrade proposal.
}

// ParamTable for upgrade module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&UpgradeParams{})
}

func (p *UpgradeParams) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyAppUpgradeMaxDepositPeriod, &p.AppUpgradeMaxDepositPeriod},
		{KeyAppUpgradeMinDeposit, &p.AppUpgradeMinDeposit},
		{KeyAppUpgradeVotingPeriod, &p.AppUpgradeVotingPeriod},
	}
}

func (p UpgradeParams) String() string {
	return fmt.Sprintf(`App Upgrade Params:
	App Upgrade Min Deposit:        %s
	App Upgrade Deposit Period:     %s
	App Upgrade Voting Period:      %s`, p.AppUpgradeMinDeposit, p.AppUpgradeMaxDepositPeriod, p.AppUpgradeVotingPeriod)
}

// default minting module parameters
func DefaultParams() UpgradeParams {
	var minDeposit = sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}

	return UpgradeParams{
		AppUpgradeMaxDepositPeriod: time.Hour * 24,
		AppUpgradeMinDeposit:       minDeposit,
		AppUpgradeVotingPeriod:     time.Hour * 72,
	}
}

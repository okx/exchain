package types

import (
	"fmt"
	"time"

	"github.com/okex/okexchain/x/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	keyAppUpgradeMaxDepositPeriod = []byte("AppUpgradeMaxDepositPeriod")
	keyAppUpgradeMinDeposit       = []byte("AppUpgradeMinDeposit")
	keyAppUpgradeVotingPeriod     = []byte("AppUpgradeVotingPeriod")
)

// UpgradeParams is the struct of upgrade module params
type UpgradeParams struct {
	// Maximum period for okb holders to deposit on a AppUpgrade proposal. Initial value: 2 days
	AppUpgradeMaxDepositPeriod time.Duration `json:"app_upgrade_max_deposit_period"`
	// Minimum deposit for a critical AppUpgrade proposal to enter voting period
	AppUpgradeMinDeposit sdk.DecCoins `json:"app_upgrade_min_deposit"`
	// Length of the critical voting period for AppUpgrade proposal
	AppUpgradeVotingPeriod time.Duration `json:"app_upgrade_voting_period"`
}

// ParamKeyTable gets KeyTable for upgrade module params
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&UpgradeParams{})
}

// ParamSetPairs sets upgrade module params
func (p *UpgradeParams) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: keyAppUpgradeMaxDepositPeriod, Value: &p.AppUpgradeMaxDepositPeriod},
		{Key: keyAppUpgradeMinDeposit, Value: &p.AppUpgradeMinDeposit},
		{Key: keyAppUpgradeVotingPeriod, Value: &p.AppUpgradeVotingPeriod},
	}
}

// String returns a human readable string representation of UpgradeParams
func (p UpgradeParams) String() string {
	return fmt.Sprintf(`App Upgrade Params:
	App Upgrade Min Deposit:        %s
	App Upgrade Deposit Period:     %s
	App Upgrade Voting Period:      %s`, p.AppUpgradeMinDeposit, p.AppUpgradeMaxDepositPeriod, p.AppUpgradeVotingPeriod)
}

// DefaultParams returns default module parameters
func DefaultParams() UpgradeParams {
	var minDeposit = sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}

	return UpgradeParams{
		AppUpgradeMaxDepositPeriod: time.Hour * 24,
		AppUpgradeMinDeposit:       minDeposit,
		AppUpgradeVotingPeriod:     time.Hour * 72,
	}
}

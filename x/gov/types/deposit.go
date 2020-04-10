package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Deposit defines the deposit of proposal
type Deposit struct {
	ProposalID uint64         `json:"proposal_id" yaml:"proposal_id"` //  proposalID of the proposal
	Depositor  sdk.AccAddress `json:"depositor" yaml:"depositor"`     //  Address of the depositor
	Amount     sdk.DecCoins   `json:"amount" yaml:"amount"`           //  Deposit amount
	DepositID  uint64         `json:"deposit_id" yaml:"deposit_id"`   //  Deposit ID
}

func (d Deposit) String() string {
	return fmt.Sprintf("deposit by %s on Proposal %d is for the amount %s, the deposit id is %d",
		d.Depositor, d.ProposalID, d.Amount, d.DepositID)
}

// Deposits is a collection of Deposit objects
type Deposits []Deposit

func (d Deposits) String() string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for Proposal %d:", d[0].ProposalID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %s", dep.Depositor, dep.Amount)
	}
	return out
}

// Equals returns whether two deposits are equal.
func (d Deposit) Equals(comp Deposit) bool {
	return d.Depositor.Equals(comp.Depositor) && d.ProposalID == comp.ProposalID &&
		d.Amount.IsEqual(comp.Amount) && d.DepositID == comp.DepositID
}

// Empty returns whether a deposit is empty.
func (d Deposit) Empty() bool {
	return d.Equals(Deposit{})
}

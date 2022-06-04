package types

import (
	"fmt"
)

const (
	ProposalTypeSoftwareUpgrade       string = "SoftwareUpgrade"
	ProposalTypeCancelSoftwareUpgrade string = "CancelSoftwareUpgrade"
)

func (sup SoftwareUpgradeProposal) String() string {
	return fmt.Sprintf(`Software Upgrade Proposal:
  Title:       %s
  Description: %s
`, sup.Title, sup.Description)
}

func (csup CancelSoftwareUpgradeProposal) String() string {
	return fmt.Sprintf(`Cancel Software Upgrade Proposal:
  Title:       %s
  Description: %s
`, csup.Title, csup.Description)
}

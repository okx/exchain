package types

import (
	"fmt"
)

const maxAddressListLength = 100

// ProposalRoute returns the routing key of a parameter change proposal.
func (p UpdateDeploymentWhitelistProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type
func (p UpdateDeploymentWhitelistProposal) ProposalType() string {
	return string(ProposalTypeUpdateDeploymentWhitelist)
}

// ValidateBasic validates the proposal
func (p UpdateDeploymentWhitelistProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	l := len(p.DistributorAddrs)
	if l == 0 || l > maxAddressListLength {
		return fmt.Errorf("invalid distributor addresses len: %d", l)
	}
	return nil
}

// MarshalYAML pretty prints the wasm byte code
func (p UpdateDeploymentWhitelistProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title            string   `yaml:"title"`
		Description      string   `yaml:"description"`
		DistributorAddrs []string `yaml:"distributorAddrs"`
		IsAdded          bool     `yaml:"isAdded"`
	}{
		Title:            p.Title,
		Description:      p.Description,
		DistributorAddrs: p.DistributorAddrs,
		IsAdded:          p.IsAdded,
	}, nil
}

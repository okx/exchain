package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
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
	return validateDistributorAddrs(p.DistributorAddrs)
}

// MarshalYAML pretty prints the wasm byte code
func (p UpdateDeploymentWhitelistProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title            string   `yaml:"title"`
		Description      string   `yaml:"description"`
		DistributorAddrs []string `yaml:"distributor_addresses"`
	}{
		Title:            p.Title,
		Description:      p.Description,
		DistributorAddrs: p.DistributorAddrs,
	}, nil
}

func validateDistributorAddrs(addrs []string) error {
	if IsAllAddress(addrs) {
		return nil
	}
	for _, addr := range addrs {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return err
		}
	}
	return nil
}

func IsAllAddress(addrs []string) bool {
	return len(addrs) == 1 && addrs[0] == "all"
}

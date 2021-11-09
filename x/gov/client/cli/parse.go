package cli

import (
	"encoding/json"
	"fmt"
	"os"

	govutils "github.com/okex/exchain/x/gov/client/utils"
	"github.com/spf13/viper"
)

func parseSubmitProposalFlags() (*proposal, error) {
	proposal := &proposal{}
	file := viper.GetString(flagProposal)

	if file == "" {
		proposal.Title = viper.GetString(flagTitle)
		proposal.Description = viper.GetString(flagDescription)
		proposal.Type = govutils.NormalizeProposalType(viper.GetString(flagProposalType))
		proposal.Deposit = viper.GetString(flagDeposit)
		return proposal, nil
	}

	for _, flag := range proposalFlags {
		if viper.GetString(flag) != "" {
			return nil, fmt.Errorf("--%s flag provided alongside --proposal, which is a noop", flag)
		}
	}

	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, proposal)
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

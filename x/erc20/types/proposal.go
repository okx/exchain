package types

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	govtypes "github.com/okex/exchain/x/gov/types"
)

const (
	// proposalTypeTokenMapping defines the type for a TokenMappingProposal
	proposalTypeTokenMapping = "TokenMapping"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeTokenMapping)
	govtypes.RegisterProposalTypeCodec(TokenMappingProposal{}, "okexchain/erc20/TokenMappingProposal")
}

var _ govtypes.Content = (*TokenMappingProposal)(nil)

type TokenMappingProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Denom       string `json:"denom" yaml:"denom"`
	Contract    string `json:"contract" yaml:"contract"`
}

func NewTokenMappingChangeProposal(title, description, denom string, contractAddr *common.Address) *TokenMappingProposal {
	contract := ""
	if contractAddr != nil {
		contract = contractAddr.Hex()
	}
	return &TokenMappingProposal{title, description, denom, contract}
}

func (tp TokenMappingProposal) GetTitle() string       { return tp.Title }
func (tp TokenMappingProposal) GetDescription() string { return tp.Description }
func (tp TokenMappingProposal) ProposalRoute() string  { return RouterKey }
func (tp TokenMappingProposal) ProposalType() string   { return proposalTypeTokenMapping }
func (tp TokenMappingProposal) ValidateBasic() sdk.Error {
	// TODO
	return nil
}
func (tp TokenMappingProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Token Mapping Proposal:
  Title:       %s
  Description: %s
  Denom:       %s
  Contract:    %s
`, tp.Title, tp.Description, tp.Denom, tp.Contract))

	return b.String()
}

package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/okex/okchain/x/gov/types"
)

const (
	ProposalTypeCertifiedToken = "CertifiedToken"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeCertifiedToken)
	govtypes.RegisterProposalTypeCodec(CertifiedTokenProposal{}, "okchain/token/CertifiedTokenProposal")

}

// Assert CertifiedTokenProposal implements govtypes.Content at compile-time
var _ govtypes.Content = (*CertifiedTokenProposal)(nil)

// CertifiedTokenProposal represents CertifiedToken proposal object
type CertifiedTokenProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Token       CertifiedToken `json:"token" yaml:"token"`
}

type CertifiedTokenExport struct {
	ProposalID uint64         `json:"id" yaml:"id"`
	Token      CertifiedToken `json:"token" yaml:"token"`
}

type CertifiedToken struct {
	Description string         `json:"description"`
	Symbol      string         `json:"symbol"`
	WholeName   string         `json:"whole_name"`
	TotalSupply string         `json:"total_supply"`
	Owner       sdk.AccAddress `json:"owner"`
	Mintable    bool           `json:"mintable"`
}

// NewCertifiedTokenProposal create a new CertifiedToken proposal object
func NewCertifiedTokenProposal(title, description string, token CertifiedToken) CertifiedTokenProposal {
	return CertifiedTokenProposal{
		Title:       title,
		Description: description,
		Token:       token,
	}
}

// GetTitle returns title of CertifiedToken proposal object
func (ctp CertifiedTokenProposal) GetTitle() string {
	return ctp.Title
}

// GetDescription returns description of CertifiedToken proposal object
func (ctp CertifiedTokenProposal) GetDescription() string {
	return ctp.Description
}

// ProposalRoute returns route key of CertifiedToken proposal object
func (CertifiedTokenProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of CertifiedToken proposal object
func (CertifiedTokenProposal) ProposalType() string {
	return ProposalTypeCertifiedToken
}

// ValidateBasic validates CertifiedToken proposal
func (ctp CertifiedTokenProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(ctp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace, "failed to submit CertifiedToken proposal because title is blank")
	}
	if len(ctp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace,
			fmt.Sprintf("failed to submit CertifiedToken proposal because title is longer than max length of %d", govtypes.MaxTitleLength))
	}

	if len(ctp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace,
			"failed to submit CertifiedToken proposal because description is blank")
	}

	if len(ctp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace,
			fmt.Sprintf("failed to submit CertifiedToken proposal because description is longer than max length of %d", govtypes.MaxDescriptionLength))
	}

	if ctp.ProposalType() != ProposalTypeCertifiedToken {
		return govtypes.ErrInvalidProposalType(DefaultCodespace,
			ctp.ProposalType())
	}

	// check owner
	if ctp.Token.Owner.Empty() {
		return sdk.ErrInvalidAddress(ctp.Token.Owner.String())
	}

	// check symbol
	if len(ctp.Token.Symbol) == 0 {
		return sdk.ErrUnknownRequest("failed to check CertifiedToken proposal because symbol cannot be empty")
	}
	if !ValidOriginalSymbol(ctp.Token.Symbol) {
		return sdk.ErrUnknownRequest("failed to check CertifiedToken proposal because invalid original symbol: " + ctp.Token.Symbol)
	}

	// check wholeName
	isValid := wholeNameValid(ctp.Token.WholeName)
	if !isValid {
		return sdk.ErrUnknownRequest("failed to check CertifiedToken proposal because invalid whole name")
	}
	// check desc
	if len(ctp.Token.Description) > DescLenLimit {
		return sdk.ErrUnknownRequest("failed to check issue msg because invalid desc")
	}
	// check totalSupply
	totalSupply, err := sdk.NewDecFromStr(ctp.Token.TotalSupply)
	if err != nil {
		return err
	}
	if totalSupply.GT(sdk.NewDec(TotalSupplyUpperbound)) || totalSupply.LTE(sdk.ZeroDec()) {
		return sdk.ErrUnknownRequest("failed to check issue msg because invalid total supply")
	}

	return nil
}

// String converts CertifiedToken proposal object to string
func (ctp CertifiedTokenProposal) String() string {
	return fmt.Sprintf(`CertifiedTokenProposal:
 Title:               %s
 Description:         %s
 Type:                %s
 Token Description:   %s
 Token Symbol:        %s
 Token WholeName:     %s
 Token Owner:         %s
 Token Mintable:      %v
 Token TotalSupply:   %s
`, ctp.Title, ctp.Description, ctp.ProposalType(),
		ctp.Token.Description, ctp.Token.Symbol,
		ctp.Token.WholeName, ctp.Token.Owner,
		ctp.Token.Mintable, ctp.Token.TotalSupply,
	)
}

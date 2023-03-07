package types

import (
	"encoding/json"
	"fmt"
	"github.com/okx/okbchain/libs/system"
	"strings"

	"github.com/okx/okbchain/x/common"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/tendermint/crypto"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

const (
	proposalTypeProposeValidator = "ProposeValidator"
	ProposeValidatorProposalName = system.Chain+"/staking/ProposeValidatorProposal"
)

var _ govtypes.Content = (*ProposeValidatorProposal)(nil)

func init() {
	govtypes.RegisterProposalType(proposalTypeProposeValidator)
	govtypes.RegisterProposalTypeCodec(ProposeValidatorProposal{}, ProposeValidatorProposalName)
}

type ProposeValidator struct {
	Description Description `json:"description" yaml:"description"`
	//Commission        CommissionRates `json:"commission" yaml:"commission"`
	MinSelfDelegation sdk.SysCoin    `json:"min_self_delegation" yaml:"min_self_delegation"`
	DelegatorAddress  sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress  sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	PubKey            crypto.PubKey  `json:"pubkey" yaml:"pubkey"`
}

type proposeValidatorJSON struct {
	Description Description `json:"description" yaml:"description"`
	//Commission        CommissionRates `json:"commission" yaml:"commission"`
	MinSelfDelegation sdk.SysCoin    `json:"min_self_delegation" yaml:"min_self_delegation"`
	DelegatorAddress  sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress  sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	PubKey            string         `json:"pubkey" yaml:"pubkey"`
}

// ProposeValidatorProposal - structure for the proposal of proposing validator
type ProposeValidatorProposal struct {
	Title       string           `json:"title" yaml:"title"`
	Description string           `json:"description" yaml:"description"`
	IsAdd       bool             `json:"is_add" yaml:"is_add"`
	Validator   ProposeValidator `json:"validator" yaml:"validator"`
}

// NewProposeValidatorProposal creates a new instance of ProposeValidatorProposal
func NewProposeValidatorProposal(title, description string, isAdd bool, validator ProposeValidator) ProposeValidatorProposal {
	return ProposeValidatorProposal{
		Title:       title,
		Description: description,
		IsAdd:       isAdd,
		Validator:   validator,
	}
}

// GetTitle returns title of the proposal object
func (pv ProposeValidatorProposal) GetTitle() string {
	return pv.Title
}

// GetDescription returns description of proposal object
func (pv ProposeValidatorProposal) GetDescription() string {
	return pv.Description
}

// ProposalRoute returns route key of the proposal object
func (pv ProposeValidatorProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of the proposal object
func (pv ProposeValidatorProposal) ProposalType() string {
	return proposalTypeProposeValidator
}

// ValidateBasic validates the proposal
func (pv ProposeValidatorProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(pv.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(pv.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is bigger than max title length")
	}

	if len(pv.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(pv.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is bigger than max description length")
	}

	if pv.ProposalType() != proposalTypeProposeValidator {
		return govtypes.ErrInvalidProposalType(pv.ProposalType())
	}

	if pv.Validator.ValidatorAddress.Empty() {
		return govtypes.ErrInvalidProposalContent("empty validator address")
	}
	if pv.IsAdd {
		if pv.Validator.DelegatorAddress.Empty() {
			return govtypes.ErrInvalidProposalContent("empty delegator address")
		}
		if !sdk.AccAddress(pv.Validator.ValidatorAddress).Equals(pv.Validator.DelegatorAddress) {
			return govtypes.ErrInvalidProposalContent("validator address is invalid")
		}
		if pv.Validator.MinSelfDelegation.Amount.LT(sdk.ZeroDec()) || !pv.Validator.MinSelfDelegation.IsValid() {
			return govtypes.ErrInvalidProposalContent("minimum self delegation is invalid")
		}
		if pv.Validator.Description == (Description{}) {
			return govtypes.ErrInvalidProposalContent("empty description")
		}
	}
	return nil
}

// String returns a human readable string representation of a ProposeValidatorProposal
func (pv ProposeValidatorProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ProposeValidatorProposal:
 Title:					%s
 Description:        	%s
 IsADD:                 %t
 Type:                	%s
`,
			pv.Title, pv.Description, pv.IsAdd, pv.ProposalType()),
	)

	return strings.TrimSpace(builder.String())
}

// MarshalJSON implements the json.Marshaler interface to provide custom JSON serialization
func (pv ProposeValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(proposeValidatorJSON{
		Description:       pv.Description,
		DelegatorAddress:  pv.DelegatorAddress,
		ValidatorAddress:  pv.ValidatorAddress,
		PubKey:            MustBech32ifyConsPub(pv.PubKey),
		MinSelfDelegation: pv.MinSelfDelegation,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface to provide custom JSON deserialization
func (pv *ProposeValidator) UnmarshalJSON(bz []byte) error {
	var pvJSON proposeValidatorJSON
	if err := json.Unmarshal(bz, &pvJSON); err != nil {
		return common.ErrUnMarshalJSONFailed(err.Error())
	}

	pv.Description = pvJSON.Description
	pv.DelegatorAddress = pvJSON.DelegatorAddress
	pv.ValidatorAddress = pvJSON.ValidatorAddress
	var err error
	pv.PubKey, err = GetConsPubKeyBech32(pvJSON.PubKey)
	if err != nil {
		return ErrGetConsPubKeyBech32()
	}
	pv.MinSelfDelegation = pvJSON.MinSelfDelegation

	return nil
}

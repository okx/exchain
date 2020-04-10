// nolint
package v0_8

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/common/proto"
)

var (
	_ Proposal = &BasicProposal{}
	_ Proposal = &TextProposal{}
	_ Proposal = &DexListProposal{}
	_ Proposal = &ParameterProposal{}
	_ Proposal = &AppUpgradeProposal{}
)

const (
	ModuleName = "gov"

	StatusNil           ProposalStatus = 0x00
	StatusDepositPeriod ProposalStatus = 0x01
	StatusVotingPeriod  ProposalStatus = 0x02
	StatusPassed        ProposalStatus = 0x03
	StatusRejected      ProposalStatus = 0x04

	ProposalTypeNil             ProposalKind = 0x00
	ProposalTypeText            ProposalKind = 0x01
	ProposalTypeParameterChange ProposalKind = 0x02
	ProposalTypeAppUpgrade      ProposalKind = 0x03
	ProposalTypeDexList         ProposalKind = 0x04

	OptionEmpty      VoteOption = 0x00
	OptionYes        VoteOption = 0x01
	OptionAbstain    VoteOption = 0x02
	OptionNo         VoteOption = 0x03
	OptionNoWithVeto VoteOption = 0x04
)

type (
	// Proposal interface
	Proposal interface {
		GetProposalID() uint64
		SetProposalID(uint64)

		GetTitle() string
		SetTitle(string)

		GetDescription() string
		SetDescription(string)

		GetProposalType() ProposalKind
		SetProposalType(ProposalKind)

		GetStatus() ProposalStatus
		SetStatus(ProposalStatus)

		GetFinalTallyResult() TallyResult
		SetFinalTallyResult(TallyResult)

		GetSubmitTime() time.Time
		SetSubmitTime(time.Time)

		GetDepositEndTime() time.Time
		SetDepositEndTime(time.Time)

		GetTotalDeposit() sdk.DecCoins
		SetTotalDeposit(sdk.DecCoins)

		GetVotingStartTime() time.Time
		SetVotingStartTime(time.Time)

		GetVotingEndTime() time.Time
		SetVotingEndTime(time.Time)

		String() string

		GetProtocolDefinition() proto.ProtocolDefinition
		SetProtocolDefinition(proto.ProtocolDefinition)
	}

	// Basic Proposal
	BasicProposal struct {
		ProposalID   uint64       `json:"proposal_id"`
		Title        string       `json:"title"`
		Description  string       `json:"description"`
		ProposalType ProposalKind `json:"proposal_type"`

		Status           ProposalStatus `json:"proposal_status"`
		FinalTallyResult TallyResult    `json:"tally_result"`

		SubmitTime     time.Time    `json:"submit_time"`
		DepositEndTime time.Time    `json:"deposit_end_time"`
		TotalDeposit   sdk.DecCoins `json:"total_deposit"`

		VotingStartTime time.Time `json:"voting_start_time"`
		VotingEndTime   time.Time `json:"voting_end_time"`
	}

	// Text Proposals
	TextProposal struct {
		BasicProposal
	}

	// DexList Proposals
	DexListProposal struct {
		BasicProposal
		Proposer         sdk.AccAddress `json:"proposer"`    //  Proposer of proposal
		ListAsset        string         `json:"list_asset"`  //  Symbol of asset listed on Dex.
		QuoteAsset       string         `json:"quote_asset"` //  Symbol of asset quoted by asset listed on Dex.
		InitPrice        sdk.Dec        `json:"init_price"`  //  Init price of asset listed on Dex.
		BlockHeight      uint64         `json:"block_height"`
		MaxPriceDigit    uint64         `json:"max_price_digit"` //  Decimal of price
		MaxSizeDigit     uint64         `json:"max_size_digit"`  //  Decimal of trade quantity
		MinTradeSize     sdk.Dec        `json:"min_trade_size"`
		DexListStartTime time.Time      `json:"dex_list_start_time"`
		DexListEndTime   time.Time      `json:"dex_list_end_time"`
	}

	ParameterProposal struct {
		BasicProposal
		Params Params `json:"params"`
		Height int64  `json:"height"`
	}

	AppUpgradeProposal struct {
		BasicProposal
		ProtocolDefinition proto.ProtocolDefinition
	}

	Param struct {
		Subspace string `json:"subspace"`
		Key      string `json:"key"`
		Value    string `json:"value"`
	}

	// Tally Results
	TallyResult struct {
		TotalBonded sdk.Dec `json:"total_bonded"`
		TotalVoting sdk.Dec `json:"total_voting"`
		Yes         sdk.Dec `json:"yes"`
		Abstain     sdk.Dec `json:"abstain"`
		No          sdk.Dec `json:"no"`
		NoWithVeto  sdk.Dec `json:"no_with_veto"`
	}

	// deposit
	Deposit struct {
		Depositor  sdk.AccAddress `json:"depositor"`   //  Address of the depositor
		ProposalID uint64         `json:"proposal_id"` //  proposalID of the proposal
		Amount     sdk.DecCoins   `json:"amount"`      //  deposit amount
		DepositID  uint64         `json:"deposit_id"`  //  id of deposit
	}

	// Vote
	Vote struct {
		Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
		ProposalID uint64         `json:"proposal_id"` //  proposalID of the proposal
		Option     VoteOption     `json:"option"`      //  option from OptionSet chosen by the voter
		VoteID     uint64         `json:"vote_id"`     //  id of vote
	}

	Deposits []Deposit

	Votes []Vote

	// Type that represents VoteOption as a byte
	VoteOption byte

	Params []Param

	ProposalStatus byte

	ProposalKind byte

	// mint parameters
	GovParams struct {
		// Text proposal params
		//  Maximum period for tokt holders to deposit on a Text proposal. Initial value: 2 days
		TextMaxDepositPeriod time.Duration `json:"text_max_deposit_period"`
		//  Minimum deposit for a critical Text proposal to enter voting period.
		TextMinDeposit sdk.DecCoins `json:"text_min_deposit"`
		//  Length of the critical voting period for Text proposal.
		TextVotingPeriod time.Duration `json:"text_voting_period"`

		// ParamChange proposal params
		//  Maximum period for tokt holders to deposit on a ParamChange proposal. Initial value: 2 days
		ParamChangeMaxDepositPeriod time.Duration `json:"param_change_max_deposit_period"`
		//  Minimum deposit for a critical ParamChange proposal to enter voting period.
		ParamChangeMinDeposit sdk.DecCoins `json:"param_change_min_deposit"`
		//  Length of the critical voting period for ParamChange proposal.
		ParamChangeVotingPeriod time.Duration `json:"param_change_voting_period"`
		//  block height for ParamChange can not be greater than current block height + MaxBlockHeightPeriod
		ParamChangeMaxBlockHeight int64 `json:"param_change_max_block_height_period"`

		// AppUpgrade proposal params
		//  Maximum period for tokt holders to deposit on a AppUpgrade proposal. Initial value: 2 days
		AppUpgradeMaxDepositPeriod time.Duration `json:"app_upgrade_max_deposit_period"`
		//  Minimum deposit for a critical AppUpgrade proposal to enter voting period.
		AppUpgradeMinDeposit sdk.DecCoins `json:"app_upgrade_min_deposit"`
		//  Length of the critical voting period for AppUpgrade proposal.
		AppUpgradeVotingPeriod time.Duration `json:"app_upgrade_voting_period"`

		// DexList proposal params
		//  Maximum period for tokt holders to deposit on a dex list proposal. Initial value: 2 days
		DexListMaxDepositPeriod time.Duration `json:"dex_list_max_deposit_period"`
		//  Minimum deposit for a critical dex list proposal to enter voting period.
		DexListMinDeposit sdk.DecCoins `json:"dex_list_min_deposit"`
		//  Length of the critical voting period for dex list proposal.
		DexListVotingPeriod time.Duration `json:"dex_list_voting_period"`
		//  Fee used for voting dex list proposal
		DexListVoteFee sdk.DecCoins `json:"dex_list_vote_fee"`
		//  block height for dex list can not be greater than DexListMaxBlockHeight
		DexListMaxBlockHeight uint64 `json:"dex_list_max_block_height"`
		//  fee for dex list
		DexListFee sdk.DecCoins `json:"dex_list_fee"`
		//  expire time for dex list
		DexListExpireTime time.Duration `json:"dex_list_expire_time"`

		// tally params
		Quorum sdk.Dec `json:"quorum"`
		//  Minimum proportion of Yes votes for proposal to pass. Initial value: 0.5
		Threshold sdk.Dec `json:"threshold"`
		//  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
		Veto sdk.Dec `json:"veto"`
		//  Minimum propotion of Yes votes before voting end time for proposal to pass. Initial value: 2/3
		YesInVotePeriod sdk.Dec `json:"yes_end_vote_period"`
		//  Max tx number per block
		MaxTxNumPerBlock int64 `json:"max_tx_num_per_block"`
	}

	// GenesisState - all staking state that must be provided at genesis
	GenesisState struct {
		StartingProposalID uint64     `json:"starting_proposal_id"`
		Deposits           Deposits   `json:"deposits" yaml:"deposits"`
		Votes              Votes      `json:"votes" yaml:"votes"`
		Proposals          []Proposal `json:"proposals"`
		Params             GovParams  `json:"params"` // params
	}
)

func (tp BasicProposal) GetProposalID() uint64                      { return tp.ProposalID }
func (tp *BasicProposal) SetProposalID(proposalID uint64)           { tp.ProposalID = proposalID }
func (tp BasicProposal) GetTitle() string                           { return tp.Title }
func (tp *BasicProposal) SetTitle(title string)                     { tp.Title = title }
func (tp BasicProposal) GetDescription() string                     { return tp.Description }
func (tp *BasicProposal) SetDescription(description string)         { tp.Description = description }
func (tp BasicProposal) GetProposalType() ProposalKind              { return tp.ProposalType }
func (tp *BasicProposal) SetProposalType(proposalType ProposalKind) { tp.ProposalType = proposalType }
func (tp BasicProposal) GetStatus() ProposalStatus                  { return tp.Status }
func (tp *BasicProposal) SetStatus(status ProposalStatus)           { tp.Status = status }
func (tp BasicProposal) GetFinalTallyResult() TallyResult           { return tp.FinalTallyResult }
func (tp *BasicProposal) SetFinalTallyResult(tallyResult TallyResult) {
	tp.FinalTallyResult = tallyResult
}
func (tp BasicProposal) GetSubmitTime() time.Time            { return tp.SubmitTime }
func (tp *BasicProposal) SetSubmitTime(submitTime time.Time) { tp.SubmitTime = submitTime }
func (tp BasicProposal) GetDepositEndTime() time.Time        { return tp.DepositEndTime }
func (tp *BasicProposal) SetDepositEndTime(depositEndTime time.Time) {
	tp.DepositEndTime = depositEndTime
}
func (tp BasicProposal) GetTotalDeposit() sdk.DecCoins              { return tp.TotalDeposit }
func (tp *BasicProposal) SetTotalDeposit(totalDeposit sdk.DecCoins) { tp.TotalDeposit = totalDeposit }
func (tp BasicProposal) GetVotingStartTime() time.Time              { return tp.VotingStartTime }
func (tp *BasicProposal) SetVotingStartTime(votingStartTime time.Time) {
	tp.VotingStartTime = votingStartTime
}
func (tp BasicProposal) GetVotingEndTime() time.Time { return tp.VotingEndTime }
func (tp *BasicProposal) SetVotingEndTime(votingEndTime time.Time) {
	tp.VotingEndTime = votingEndTime
}

func (tp BasicProposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  title:              %s
  proposalType:               %s
  Status:             %s
  Submit Time:        %s
  deposit End Time:   %s
  Total deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s`, tp.ProposalID, tp.Title, tp.ProposalType,
		tp.Status, tp.SubmitTime, tp.DepositEndTime,
		tp.TotalDeposit, tp.VotingStartTime, tp.VotingEndTime)
}

// software upgrade
func (tp BasicProposal) GetProtocolDefinition() proto.ProtocolDefinition {
	return proto.ProtocolDefinition{}
}
func (tp *BasicProposal) SetProtocolDefinition(proto.ProtocolDefinition) {}

// ProposalStatusToString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "DepositPeriod":
		return StatusDepositPeriod, nil
	case "VotingPeriod":
		return StatusVotingPeriod, nil
	case "Passed":
		return StatusPassed, nil
	case "Rejected":
		return StatusRejected, nil
	case "":
		return StatusNil, nil
	default:
		return ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", str)
	}
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}
	*status = bz2
	return nil
}

// Turns VoteStatus byte to String
func (status ProposalStatus) String() string {
	switch status {
	case StatusDepositPeriod:
		return "DepositPeriod"
	case StatusVotingPeriod:
		return "VotingPeriod"
	case StatusPassed:
		return "Passed"
	case StatusRejected:
		return "Rejected"
	default:
		return ""
	}
}

// String to proposalType byte. Returns 0xff if invalid.
func ProposalTypeFromString(str string) (ProposalKind, error) {
	switch str {
	case "Text":
		return ProposalTypeText, nil
	case "ParameterChange":
		return ProposalTypeParameterChange, nil
	case "AppUpgrade":
		return ProposalTypeAppUpgrade, nil
	case "DexList":
		return ProposalTypeDexList, nil
	default:
		return ProposalKind(0xff), fmt.Errorf("'%s' is not a valid proposal type", str)
	}
}

// Marshal needed for protobuf compatibility
func (pt ProposalKind) Marshal() ([]byte, error) {
	return []byte{byte(pt)}, nil
}

// Unmarshal needed for protobuf compatibility
func (pt *ProposalKind) Unmarshal(data []byte) error {
	*pt = ProposalKind(data[0])
	return nil
}

// Marshals to JSON using string
func (pt ProposalKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (pt *ProposalKind) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalTypeFromString(s)
	if err != nil {
		return err
	}
	*pt = bz2
	return nil
}

// Turns VoteOption byte to String
func (pt ProposalKind) String() string {
	switch pt {
	case ProposalTypeText:
		return "Text"
	case ProposalTypeParameterChange:
		return "ParameterChange"
	case ProposalTypeAppUpgrade:
		return "AppUpgrade"
	case ProposalTypeDexList:
		return "DexList"
	default:
		return ""
	}
}

// String to proposalType byte.  Returns ff if invalid.
func VoteOptionFromString(str string) (VoteOption, error) {
	switch str {
	case "Yes":
		return OptionYes, nil
	case "Abstain":
		return OptionAbstain, nil
	case "No":
		return OptionNo, nil
	case "NoWithVeto":
		return OptionNoWithVeto, nil
	default:
		return VoteOption(0xff), fmt.Errorf("'%s' is not a valid vote option", str)
	}
}

// Marshal needed for protobuf compatibility
func (vo VoteOption) Marshal() ([]byte, error) {
	return []byte{byte(vo)}, nil
}

// Unmarshal needed for protobuf compatibility
func (vo *VoteOption) Unmarshal(data []byte) error {
	*vo = VoteOption(data[0])
	return nil
}

// Marshals to JSON using string
func (vo VoteOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(vo.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (vo *VoteOption) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := VoteOptionFromString(s)
	if err != nil {
		return err
	}
	*vo = bz2
	return nil
}

// Turns VoteOption byte to String
func (vo VoteOption) String() string {
	switch vo {
	case OptionYes:
		return "Yes"
	case OptionAbstain:
		return "Abstain"
	case OptionNo:
		return "No"
	case OptionNoWithVeto:
		return "NoWithVeto"
	default:
		return ""
	}
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "okchain/gov/TextProposal", nil)
	cdc.RegisterConcrete(&DexListProposal{}, "okchain/gov/DexListProposal", nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "okchain/gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&AppUpgradeProposal{}, "okchain/gov/AppUpgradeProposal", nil)
}

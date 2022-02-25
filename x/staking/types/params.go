package types

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/params"
)

const (
	// Update the validator set every 252 blocks by default
	DefaultBlocksPerEpoch = 252

	// Default maximum number of validators to vote
	DefaultMaxValsToVote = 30

	// Default validate rate update interval by hours
	DefaultValidateRateUpdateInterval = 24
)

// Staking params default values
const (

	// Default unbonding duration, 14 days
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 7 * 2

	// Default maximum number of bonded validators
	DefaultMaxValidators uint16 = 21

	DefaultEpoch              uint16 = DefaultBlocksPerEpoch
	DefaultMaxValsToAddShares uint16 = DefaultMaxValsToVote
)

var (
	// DefaultMinDelegation is the limit value of delegation or undelegation
	DefaultMinDelegation = sdk.NewDecWithPrec(1, 4)
	// DefaultMinSelfDelegation is the default value of each validator's msd (hard code)
	DefaultMinSelfDelegation = sdk.NewDec(10000)
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime     = []byte("UnbondingTime")
	KeyMaxValidators     = []byte("MaxValidators")
	KeyEpoch             = []byte("BlocksPerEpoch")    // how many blocks each epoch has
	KeyTheEndOfLastEpoch = []byte("TheEndOfLastEpoch") // a block height that is the end of last epoch

	KeyMaxValsToAddShares = []byte("MaxValsToAddShares")
	KeyMinDelegation      = []byte("MinDelegation")
	KeyMinSelfDelegation  = []byte("MinSelfDelegation")

	KeyHistoricalEntries = []byte("HistoricalEntries")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	// time duration of unbonding
	UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
	// note: we need to be a bit careful about potential overflow here, since this is user-determined
	// maximum number of validators (max uint16 = 65535)
	MaxValidators uint16 `json:"max_bonded_validators" yaml:"max_bonded_validators"`
	// epoch for validator update
	Epoch              uint16 `json:"epoch" yaml:"epoch"`
	MaxValsToAddShares uint16 `json:"max_validators_to_add_shares" yaml:"max_validators_to_add_shares"`
	// limited amount of delegate
	MinDelegation sdk.Dec `json:"min_delegation" yaml:"min_delegation"`
	// validator's self declared minimum self delegation
	MinSelfDelegation sdk.Dec `json:"min_self_delegation" yaml:"min_self_delegation"`

	HistoricalEntries uint32 `protobuf:"varint,4,opt,name=historical_entries,json=historicalEntries,proto3" json:"historical_entries,omitempty" yaml:"historical_entries"`
}

// NewParams creates a new Params instance
func NewParams(unbondingTime time.Duration, maxValidators uint16, epoch uint16, maxValsToAddShares uint16, minDelegation sdk.Dec,
	minSelfDelegation sdk.Dec) Params {
	return Params{
		UnbondingTime:      unbondingTime,
		MaxValidators:      maxValidators,
		Epoch:              epoch,
		MaxValsToAddShares: maxValsToAddShares,
		MinDelegation:      minDelegation,
		MinSelfDelegation:  minSelfDelegation,
	}
}

// TODO: to supplement the validate function for every pair of param
func validateParams(value interface{}) error {
	return nil
}

// ParamSetPairs is the implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyUnbondingTime, Value: &p.UnbondingTime, ValidatorFn: common.ValidateDurationPositive("unbonding time")},
		{Key: KeyMaxValidators, Value: &p.MaxValidators, ValidatorFn: common.ValidateUint16Positive("max validators")},
		{Key: KeyEpoch, Value: &p.Epoch, ValidatorFn: common.ValidateUint16Positive("epoch")},
		{Key: KeyMaxValsToAddShares, Value: &p.MaxValsToAddShares, ValidatorFn: common.ValidateUint16Positive("max vals to add shares")},
		{Key: KeyMinDelegation, Value: &p.MinDelegation, ValidatorFn: common.ValidateDecPositive("min delegation")},
		{Key: KeyHistoricalEntries, Value: &p.HistoricalEntries, ValidatorFn: validateHistoricalEntries},
		{Key: KeyMinSelfDelegation, Value: &p.MinSelfDelegation, ValidatorFn: common.ValidateDecPositive("min self delegation")},
	}
}
func validateHistoricalEntries(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// Equal returns a boolean determining if two Param types are identical
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingTime,
		DefaultMaxValidators,
		DefaultEpoch,
		DefaultMaxValsToAddShares,
		DefaultMinDelegation,
		DefaultMinSelfDelegation,
	)
}

// String returns a human readable string representation of the Params
func (p *Params) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time:    		%s
  Max Validators:   	 	%d
  Epoch: 					%d
  MaxValsToAddShares:       %d
  MinDelegation				%d
  MinSelfDelegation         %d`,
		p.UnbondingTime, p.MaxValidators, p.Epoch, p.MaxValsToAddShares, p.MinDelegation, p.MinSelfDelegation)
}

// Validate gives a quick validity check for a set of params
func (p Params) Validate() error {
	if p.MaxValidators == 0 {
		return fmt.Errorf("staking parameter MaxValidators must be a positive integer")
	}
	if p.Epoch == 0 {
		return fmt.Errorf("staking parameter Epoch must be a positive integer")
	}
	if p.MaxValsToAddShares == 0 {
		return fmt.Errorf("staking parameter MaxValsToAddShares must be a positive integer")
	}

	return nil
}

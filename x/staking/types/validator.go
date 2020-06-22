package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
	"gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/exported"
)

// nolint
const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength  = 70
	MaxIdentityLength = 3000
	MaxWebsiteLength  = 140
	MaxDetailsLength  = 280
)

// Implements Validator interface
var _ exported.ValidatorI = Validator{}

// Validator defines the total amount of bond shares and their exchange rate to
// coins. Slashing results in a decrease in the exchange rate, allowing correct
// calculation of future undelegations without iterating over delegators.
// When coins are delegated to this validator, the validator is credited with a
// delegation whose number of bond shares is based on the amount of coins delegated
// divided by the current exchange rate. Voting power can be calculated as total
// bonded shares multiplied by exchange rate.
type Validator struct {
	// address of the validator's operator; bech encoded in JSON
	OperatorAddress sdk.ValAddress `json:"operator_address" yaml:"operator_address"`
	// the consensus public key of the validator; bech encoded in JSON
	ConsPubKey crypto.PubKey `json:"consensus_pubkey" yaml:"consensus_pubkey"`
	// has the validator been jailed from bonded status?
	Jailed bool `json:"jailed" yaml:"jailed"`
	// validator status (bonded/unbonding/unbonded)
	Status sdk.BondStatus `json:"status" yaml:"status"`
	// delegated tokens (incl. self-delegation)
	Tokens sdk.Int `json:"tokens" yaml:"tokens"`
	// total shares added to a validator
	DelegatorShares sdk.Dec `json:"delegator_shares" yaml:"delegator_shares"`
	// description terms for the validator
	Description Description `json:"description" yaml:"description"`
	// if unbonding, height at which this validator has begun unbonding
	UnbondingHeight int64 `json:"unbonding_height" yaml:"unbonding_height"`
	// if unbonding, min time for the validator to complete unbonding
	UnbondingCompletionTime time.Time `json:"unbonding_time" yaml:"unbonding_time"`
	// commission parameters
	Commission Commission `json:"commission" yaml:"commission"`
	// validator's self declared minimum self delegation
	MinSelfDelegation sdk.Dec `json:"min_self_delegation" yaml:"min_self_delegation"`
}

// MarshalYAML implements the text format for yaml marshaling due to consensus pubkey
func (v Validator) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		Status                  sdk.BondStatus
		Jailed                  bool
		UnbondingHeight         int64
		ConsPubKey              string
		OperatorAddress         sdk.ValAddress
		Tokens                  sdk.Int
		DelegatorShares         sdk.Dec
		Description             Description
		UnbondingCompletionTime time.Time
		Commission              Commission
		MinSelfDelegation       sdk.Dec
	}{
		OperatorAddress:         v.OperatorAddress,
		ConsPubKey:              sdk.MustBech32ifyConsPub(v.ConsPubKey),
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		Tokens:                  v.Tokens,
		DelegatorShares:         v.DelegatorShares,
		Description:             v.Description,
		UnbondingHeight:         v.UnbondingHeight,
		UnbondingCompletionTime: v.UnbondingCompletionTime,
		Commission:              v.Commission,
		MinSelfDelegation:       v.MinSelfDelegation,
	})
	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

// Validators is a collection of Validator
type Validators []Validator

// String returns a human readable string representation of Validators
func (v Validators) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ToSDKValidators converts []Validators to []sdk.Validators
func (v Validators) ToSDKValidators() (validators []exported.ValidatorI) {
	for _, val := range v {
		validators = append(validators, val)
	}
	return validators
}

// NewValidator initializes a new validator
func NewValidator(operator sdk.ValAddress, pubKey crypto.PubKey, description Description, minSelfDelegation sdk.Dec) Validator {
	return Validator{
		OperatorAddress:         operator,
		ConsPubKey:              pubKey,
		Jailed:                  false,
		Status:                  sdk.Unbonded,
		Tokens:                  sdk.ZeroInt(),
		DelegatorShares:         sdk.ZeroDec(),
		Description:             description,
		UnbondingHeight:         int64(0),
		UnbondingCompletionTime: time.Unix(0, 0).UTC(),
		Commission:              NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		MinSelfDelegation:       minSelfDelegation,
	}
}

// MustMarshalValidator must return the marshaling bytes of a validator
func MustMarshalValidator(cdc *codec.Codec, validator Validator) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(validator)
}

// MustUnmarshalValidator must return the validator entity by unmarshaling
func MustUnmarshalValidator(cdc *codec.Codec, value []byte) Validator {
	validator, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}
	return validator
}

// UnmarshalValidator unmarshals a validator from a store value
func UnmarshalValidator(cdc *codec.Codec, value []byte) (validator Validator, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &validator)
	return validator, err
}

// String returns a human readable string representation of a validator.
func (v Validator) String() string {
	bechConsPubKey, err := sdk.Bech32ifyConsPub(v.ConsPubKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`Validator
  Operator Address:           %s
  Validator Consensus Pubkey: %s
  Jailed:                     %v
  Status:                     %s
  Tokens:                     %s
  Delegator Shares:           %s
  Description:                %s
  Unbonding Height:           %d
  Unbonding Completion Time:  %v
  Minimum Self Delegation:    %v
  Commission:                 %s`,
		v.OperatorAddress, bechConsPubKey,
		v.Jailed, v.Status, v.Tokens,
		v.DelegatorShares, v.Description,
		v.UnbondingHeight, v.UnbondingCompletionTime, v.MinSelfDelegation,
		v.Commission)
}

// this is a helper struct used for JSON de- and encoding only
type bechValidator struct {
	// the bech32 address of the validator's operator
	OperatorAddress sdk.ValAddress `json:"operator_address" yaml:"operator_address"`
	// the bech32 consensus public key of the validator
	ConsPubKey string `json:"consensus_pubkey" yaml:"consensus_pubkey"`
	// has the validator been jailed from bonded status?
	Jailed bool `json:"jailed" yaml:"jailed"`
	// validator status (bonded/unbonding/unbonded)
	Status sdk.BondStatus `json:"status" yaml:"status"`
	// delegated tokens (incl. self-delegation)
	Tokens sdk.Int `json:"tokens" yaml:"tokens"`
	// total shares on a validator
	DelegatorShares sdk.Dec `json:"delegator_shares" yaml:"delegator_shares"`
	// description terms for the validator
	Description Description `json:"description" yaml:"description"`
	// if unbonding, height at which this validator has begun unbonding
	UnbondingHeight int64 `json:"unbonding_height" yaml:"unbonding_height"`
	// if unbonding, min time for the validator to complete unbonding
	UnbondingCompletionTime time.Time `json:"unbonding_time" yaml:"unbonding_time"`
	// commission parameters
	Commission Commission `json:"commission" yaml:"commission"`
	// minimum self delegation
	MinSelfDelegation sdk.Dec `json:"min_self_delegation" yaml:"min_self_delegation"`
}

// MarshalJSON marshals the validator to JSON using Bech32
func (v Validator) MarshalJSON() ([]byte, error) {
	bechConsPubKey, err := sdk.Bech32ifyConsPub(v.ConsPubKey)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSON(bechValidator{
		OperatorAddress:         v.OperatorAddress,
		ConsPubKey:              bechConsPubKey,
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		Tokens:                  v.Tokens,
		DelegatorShares:         v.DelegatorShares,
		Description:             v.Description,
		UnbondingHeight:         v.UnbondingHeight,
		UnbondingCompletionTime: v.UnbondingCompletionTime,
		MinSelfDelegation:       v.MinSelfDelegation,
		Commission:              v.Commission,
	})
}

// UnmarshalJSON unmarshals the validator from JSON using Bech32
func (v *Validator) UnmarshalJSON(data []byte) error {
	bv := &bechValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := sdk.GetConsPubKeyBech32(bv.ConsPubKey)
	if err != nil {
		return err
	}
	*v = Validator{
		OperatorAddress:         bv.OperatorAddress,
		ConsPubKey:              consPubKey,
		Jailed:                  bv.Jailed,
		Tokens:                  bv.Tokens,
		Status:                  bv.Status,
		DelegatorShares:         bv.DelegatorShares,
		Description:             bv.Description,
		UnbondingHeight:         bv.UnbondingHeight,
		UnbondingCompletionTime: bv.UnbondingCompletionTime,
		Commission:              bv.Commission,
		MinSelfDelegation:       bv.MinSelfDelegation,
	}
	return nil
}

// TestEquivalent is only for the ut
func (v Validator) TestEquivalent(v2 Validator) bool {
	return v.ConsPubKey.Equals(v2.ConsPubKey) &&
		bytes.Equal(v.OperatorAddress, v2.OperatorAddress) &&
		v.Status.Equal(v2.Status) &&
		v.Tokens.Equal(v2.Tokens) &&
		v.DelegatorShares.Equal(v2.DelegatorShares) &&
		v.Description == v2.Description &&
		v.Commission.Equal(v2.Commission)
}

// ConsAddress returns the TM validator address
func (v Validator) ConsAddress() sdk.ConsAddress {
	return sdk.ConsAddress(v.ConsPubKey.Address())
}

// IsBonded checks if the validator status equals Bonded
func (v Validator) IsBonded() bool {
	return v.GetStatus().Equal(sdk.Bonded)
}

// IsUnbonded checks if the validator status equals Unbonded
func (v Validator) IsUnbonded() bool {
	return v.GetStatus().Equal(sdk.Unbonded)
}

// IsUnbonding checks if the validator status equals Unbonding
func (v Validator) IsUnbonding() bool {
	return v.GetStatus().Equal(sdk.Unbonding)
}

// DoNotModifyDesc is the constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker" yaml:"moniker"`   // name
	Identity string `json:"identity" yaml:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website" yaml:"website"`   // optional website link
	Details  string `json:"details" yaml:"details"`   // optional details
}

// NewDescription returns a new Description with the provided values.
func NewDescription(moniker, identity, website, details string) Description {
	return Description{
		Moniker:  moniker,
		Identity: identity,
		Website:  website,
		Details:  details,
	}
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, sdk.Error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}
	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}
	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}
	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return Description{
		Moniker:  d2.Moniker,
		Identity: d2.Identity,
		Website:  d2.Website,
		Details:  d2.Details,
	}.EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, sdk.Error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, ErrDescriptionLength(DefaultCodespace, "moniker", len(d.Moniker), MaxMonikerLength)
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, ErrDescriptionLength(DefaultCodespace, "identity", len(d.Identity), MaxIdentityLength)
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, ErrDescriptionLength(DefaultCodespace, "website", len(d.Website), MaxWebsiteLength)
	}
	if len(d.Details) > MaxDetailsLength {
		return d, ErrDescriptionLength(DefaultCodespace, "details", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.ConsPubKey),
		Power:  v.ConsensusPower(),
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.ConsPubKey),
		Power:  0,
	}
}

// SetInitialCommission attempts to set a validator's initial commission. An
// error is returned if the commission is invalid.
func (v Validator) SetInitialCommission(commission Commission) (Validator, sdk.Error) {
	if err := commission.Validate(); err != nil {
		return v, err
	}

	v.Commission = commission
	return v, nil
}

// TODO : remove these functions that implements some origin interface later

// TokensFromShares calculates the token worth of provided shares
func (v Validator) TokensFromShares(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).Quo(v.DelegatorShares)
}

// TokensFromSharesTruncated calculates the token worth of provided shares, truncated
func (v Validator) TokensFromSharesTruncated(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoTruncate(v.DelegatorShares)
}

// TokensFromSharesRoundUp returns the token worth of provided shares, rounded up
// No usage found in All Places
func (v Validator) TokensFromSharesRoundUp(shares sdk.Dec) sdk.Dec {
	return sdk.ZeroDec()
	//return (shares.MulInt(v.Tokens)).QuoRoundUp(v.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount
// It returns an error if the validator has no tokens
// No usage found in All Places
func (v Validator) SharesFromTokens(amt sdk.Int) (sdk.Dec, sdk.Error) {
	return sdk.ZeroDec(), nil
	//if v.Tokens.IsZero() {
	//	return sdk.ZeroDec(), ErrInsufficientShares(DefaultCodespace)
	//}
	//
	//return v.GetDelegatorShares().MulInt(amt).QuoInt(v.GetTokens()), nil
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given a bond amount
// It returns an error if the validator has no tokens
// No usage found in All Places
func (v Validator) SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, sdk.Error) {
	return sdk.ZeroDec(), nil
	//if v.Tokens.IsZero() {
	//	return sdk.ZeroDec(), ErrInsufficientShares(DefaultCodespace)
	//}
	//
	//return v.GetDelegatorShares().MulInt(amt).QuoTruncate(v.GetTokens().ToDec()), nil
}

// BondedTokens gets the bonded tokens which the validator holds
func (v Validator) BondedTokens() sdk.Int {
	if v.IsBonded() {
		return v.Tokens
	}
	return sdk.ZeroInt()
}

// ConsensusPower gets the consensus-engine power
func (v Validator) ConsensusPower() int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPowerByShares()
	}
	return 0
}

// UpdateStatus updates the location of the shares within a validator
// to reflect the new status
func (v Validator) UpdateStatus(newStatus sdk.BondStatus) Validator {
	v.Status = newStatus
	return v
}

// nolint - for ValidatorI
func (v Validator) IsJailed() bool                { return v.Jailed }
func (v Validator) GetMoniker() string            { return v.Description.Moniker }
func (v Validator) GetStatus() sdk.BondStatus     { return v.Status }
func (v Validator) GetOperator() sdk.ValAddress   { return v.OperatorAddress }
func (v Validator) GetConsPubKey() crypto.PubKey  { return v.ConsPubKey }
func (v Validator) GetConsAddr() sdk.ConsAddress  { return sdk.ConsAddress(v.ConsPubKey.Address()) }
func (v Validator) GetTokens() sdk.Int            { return v.Tokens }
func (v Validator) GetBondedTokens() sdk.Int      { return sdk.ZeroInt() }
func (v Validator) GetConsensusPower() int64      { return v.ConsensusPower() }
func (v Validator) GetCommission() sdk.Dec        { return v.Commission.Rate }
func (v Validator) GetMinSelfDelegation() sdk.Dec { return v.MinSelfDelegation }
func (v Validator) GetDelegatorShares() sdk.Dec   { return v.DelegatorShares }

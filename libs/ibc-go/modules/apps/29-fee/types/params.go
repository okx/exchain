package types

import (
	"fmt"

	paramtypes "github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

const (
	DefaultAutoDispatch          = true
	DefaultFeePercent     uint32 = 100
	DefaultRecvPercent    uint32 = 50
	DefaultAckPercent     uint32 = 30
	DefaultTimeOutPercent uint32 = 20
)

var (
	KeyAllowAutoDispatch = []byte("AllowAutoDispatch")
	KeyFeePercent        = []byte("FeePercent")
	KeyRecvPercent       = []byte("RecvPercent")
	KeyAckPercent        = []byte("AckPercent")
	KeyTimeOutPercent    = []byte("TimeOutPercent")
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams is the default parameter configuration for the ibc-transfer module
func DefaultParams() *Params {
	return NewParams(DefaultAutoDispatch, DefaultFeePercent, DefaultRecvPercent, DefaultAckPercent, DefaultTimeOutPercent)
}

// NewParams creates a new parameter configuration for the ibc transfer module
func NewParams(allowAutoDispath bool,
	feePercent uint32,
	recvPercent uint32,
	ackPercent uint32,
	timeOutPercent uint32) *Params {
	return &Params{
		AllowAutoDispath: allowAutoDispath,
		FeePercent:       feePercent,
		RecvPercent:      recvPercent,
		AckPercent:       ackPercent,
		TimeOutPercent:   timeOutPercent}
}

// Validate all ibc-transfer module parameters
func (p Params) Validate() error {
	if err := validateEnabled(p.AllowAutoDispath); err != nil {
		return err
	}
	if err := validateUint(p.FeePercent); err != nil {
		return err
	}
	if err := validateUint(p.RecvPercent); err != nil {
		return err
	}
	if err := validateUint(p.AckPercent); err != nil {
		return err
	}
	return validateUint(p.TimeOutPercent)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllowAutoDispatch, p.AllowAutoDispath, validateEnabled),
		paramtypes.NewParamSetPair(KeyFeePercent, p.FeePercent, validateUint),
		paramtypes.NewParamSetPair(KeyRecvPercent, p.RecvPercent, validateUint),
		paramtypes.NewParamSetPair(KeyAckPercent, p.AckPercent, validateUint),
		paramtypes.NewParamSetPair(KeyTimeOutPercent, p.TimeOutPercent, validateUint),
	}
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateUint(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateFloat(i interface{}) error {
	_, ok := i.(float64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

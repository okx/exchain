package types

import (
	"time"

	sdktypes "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// PubkeyType is to be compatible with the response format of the standard cosmos REST API.
const PubkeyType = "/cosmos.crypto.ed25519.PubKey"

type CosmosAny struct {
	// nolint
	TypeUrl string `protobuf:"bytes,1,opt,name=type_url,json=typeUrl,proto3" json:"@type,omitempty"`
	// Must be a valid serialized protocol buffer of the above specified type.
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func WrapCosmosAny(v []byte) CosmosAny {
	return CosmosAny{
		TypeUrl: PubkeyType,
		Value:   v,
	}
}

// CM45Validator is constructed to be compatible with ATOMScan returning the latest cosmos REST API response
type CM45Validator struct {
	// address of the validator's operator; bech encoded in JSON
	OperatorAddress sdktypes.ValAddress `json:"operator_address" yaml:"operator_address"`
	// the consensus public key of the validator; bech encoded in JSON
	ConsPubKey *CosmosAny `json:"consensus_pubkey" yaml:"consensus_pubkey"`
	// has the validator been jailed from bonded status?
	Jailed bool `json:"jailed" yaml:"jailed"`
	// validator status (bonded/unbonding/unbonded)
	Status string `json:"status" yaml:"status"`
	// delegated tokens (incl. self-delegation)
	Tokens sdktypes.Int `json:"tokens" yaml:"tokens"`
	// total shares added to a validator
	DelegatorShares sdktypes.Dec `json:"delegator_shares" yaml:"delegator_shares"`
	// description terms for the validator
	Description Description `json:"description" yaml:"description"`
	// if unbonding, height at which this validator has begun unbonding
	UnbondingHeight int64 `json:"unbonding_height" yaml:"unbonding_height"`
	// if unbonding, min time for the validator to complete unbonding
	UnbondingCompletionTime time.Time `json:"unbonding_time" yaml:"unbonding_time"`
	// commission parameters
	Commission Commission `json:"commission" yaml:"commission"`
	// validator's self declared minimum self delegation
	MinSelfDelegation sdktypes.Dec `json:"min_self_delegation" yaml:"min_self_delegation"`
}

func WrapCM45Validator(v Validator, ca *CosmosAny) CM45Validator {
	return CM45Validator{
		OperatorAddress:         v.OperatorAddress,
		ConsPubKey:              ca,
		Jailed:                  v.Jailed,
		Status:                  v.Status.CM45String(),
		Tokens:                  v.Tokens,
		DelegatorShares:         v.DelegatorShares,
		Description:             v.Description,
		UnbondingHeight:         v.UnbondingHeight,
		UnbondingCompletionTime: v.UnbondingCompletionTime,
		Commission:              v.Commission,
		MinSelfDelegation:       v.MinSelfDelegation,
	}
}

type WrappedValidators struct {
	Vs []CM45Validator `json:"validators" yaml:"validator"`
}

func NewWrappedValidators(vs []CM45Validator) WrappedValidators {
	return WrappedValidators{
		Vs: vs,
	}
}

type WrappedValidator struct {
	V CM45Validator `json:"validator" yaml:"validator"`
}

func NewWrappedValidator(v CM45Validator) WrappedValidator {
	return WrappedValidator{
		V: v,
	}
}

type WrappedPool struct {
	P Pool `json:"pool" yaml:"pool"`
}

func NewWrappedPool(p Pool) WrappedPool {
	return WrappedPool{
		P: p,
	}
}

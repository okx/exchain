package simapp

import (
	"errors"
	"github.com/ethereum/go-ethereum/rlp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"io"
	"math/big"
)

var _ authexported.GenesisAccount = (*SimGenesisAccount)(nil)

func init() {
	authexported.RegisterConcreteAccountInfo(uint(authexported.SimGenesisAcc), &SimGenesisAccount{})
}

// SimGenesisAccount defines a type that implements the GenesisAccount interface
// to be used for simulation accounts in the genesis state.
type SimGenesisAccount struct {
	*authtypes.BaseAccount

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting" yaml:"original_vesting"`   // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free" yaml:"delegated_free"`       // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting" yaml:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64     `json:"start_time" yaml:"start_time"`               // vesting start time (UNIX Epoch time)
	EndTime          int64     `json:"end_time" yaml:"end_time"`                   // vesting end time (UNIX Epoch time)

	// module account fields
	ModuleName        string   `json:"module_name" yaml:"module_name"`               // name of the module account
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // permissions of module account
}

// Validate checks for errors on the vesting and module account parameters
func (sga SimGenesisAccount) Validate() error {
	if !sga.OriginalVesting.IsZero() {
		if sga.OriginalVesting.IsAnyGT(sga.Coins) {
			return errors.New("vesting amount cannot be greater than total amount")
		}
		if sga.StartTime >= sga.EndTime {
			return errors.New("vesting start-time cannot be before end-time")
		}
	}

	if sga.ModuleName != "" {
		ma := supply.ModuleAccount{
			BaseAccount: sga.BaseAccount, Name: sga.ModuleName, Permissions: sga.ModulePermissions,
		}
		if err := ma.Validate(); err != nil {
			return err
		}
	}

	return sga.BaseAccount.Validate()
}

type SimGenesisAccountPretty struct {
	authtypes.BaseAccountPretty `json:"base_account_pretty" yaml:"base_account_pretty"`

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting" yaml:"original_vesting"`   // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free" yaml:"delegated_free"`       // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting" yaml:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        *big.Int     `json:"start_time" yaml:"start_time"`               // vesting start time (UNIX Epoch time)
	EndTime          *big.Int     `json:"end_time" yaml:"end_time"`                   // vesting end time (UNIX Epoch time)

	// module account fields
	ModuleName        string   `json:"module_name" yaml:"module_name"`               // name of the module account
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // permissions of module account
}

func (alia SimGenesisAccountPretty) Pretty2Acc() (SimGenesisAccount, error) {
	bsAcc, err := alia.BaseAccountPretty.Pretty2Acc()
	if err != nil {
		return SimGenesisAccount{}, err
	}

	sga := SimGenesisAccount{
		BaseAccount: &bsAcc,
		OriginalVesting: alia.OriginalVesting,
		DelegatedFree: alia.DelegatedFree,
		DelegatedVesting: alia.DelegatedVesting,
		StartTime: alia.StartTime.Int64(),
		EndTime: alia.EndTime.Int64(),
		ModuleName: alia.ModuleName,
		ModulePermissions: alia.ModulePermissions,
	}

	return sga, nil
}

func (sga SimGenesisAccount) GetPrettyAccount() (SimGenesisAccountPretty, error) {
	if sga.BaseAccount == nil {
		return SimGenesisAccountPretty{}, errors.New("nil base account")
	}

	baseAccPretty, err := sga.BaseAccount.GetPrettyAccount()
	if err != nil {
		return SimGenesisAccountPretty{}, err
	}

	sgaPretty := SimGenesisAccountPretty{
		BaseAccountPretty: baseAccPretty,
		OriginalVesting: sga.OriginalVesting,
		DelegatedFree: sga.DelegatedFree,
		DelegatedVesting: sga.DelegatedVesting,
		StartTime: big.NewInt(sga.StartTime),
		EndTime: big.NewInt(sga.EndTime),
		ModuleName: sga.ModuleName,
		ModulePermissions: sga.ModulePermissions,
	}

	return sgaPretty, nil
}

// RLPEncodeToBytes returns the rlp representation of an account.
func (sga SimGenesisAccount) RLPEncodeToBytes() ([]byte, error) {
	alias, err := sga.GetPrettyAccount()
	if err != nil {
		return nil, err
	}

	return rlp.EncodeToBytes(alias)
}

// RLPDecodeBytes reduction account from rlp encode bytes
func (sga *SimGenesisAccount) RLPDecodeBytes(data []byte) error {
	var alia SimGenesisAccountPretty
	err := rlp.DecodeBytes(data, &alia)
	if err != nil {
		return err
	}
	*sga, err = alia.Pretty2Acc()
	return err
}

func (sga *SimGenesisAccount) EncodeRLP(w io.Writer) error {
	alias, err := sga.GetPrettyAccount()
	if err != nil {
		return err
	}

	if err = rlp.Encode(w, authexported.SimGenesisAcc); err != nil {
		return err
	}
	return rlp.Encode(w, alias)
}

func (sga *SimGenesisAccount) DecodeRLP(s *rlp.Stream) error {
	var alia SimGenesisAccountPretty
	err := s.Decode(&alia)
	if err != nil {
		return err
	}

	*sga, err = alia.Pretty2Acc()

	return err
}

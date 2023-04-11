package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/gov/types"
	"github.com/okex/exchain/libs/tendermint/global"
	govtypes "github.com/okex/exchain/x/gov/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

func TestValidateProposalCommons(t *testing.T) {
	type commonProposal struct {
		Title, Description string
	}

	specs := map[string]struct {
		src    commonProposal
		expErr bool
	}{
		"all good": {src: commonProposal{
			Title:       "Foo",
			Description: "Bar",
		}},
		"prevent empty title": {
			src: commonProposal{
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent white space only title": {
			src: commonProposal{
				Title:       " ",
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent leading white spaces in title": {
			src: commonProposal{
				Title:       " Foo",
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent title exceeds max length ": {
			src: commonProposal{
				Title:       strings.Repeat("a", govtypes.MaxTitleLength+1),
				Description: "Bar",
			},
			expErr: true,
		},
		"prevent empty description": {
			src: commonProposal{
				Title: "Foo",
			},
			expErr: true,
		},
		"prevent leading white spaces in description": {
			src: commonProposal{
				Title:       "Foo",
				Description: " Bar",
			},
			expErr: true,
		},
		"prevent white space only description": {
			src: commonProposal{
				Title:       "Foo",
				Description: " ",
			},
			expErr: true,
		},
		"prevent descr exceeds max length ": {
			src: commonProposal{
				Title:       "Foo",
				Description: strings.Repeat("a", govtypes.MaxDescriptionLength+1),
			},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := validateProposalCommons(spec.src.Title, spec.src.Description)
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateStoreCodeProposal(t *testing.T) {
	var (
		anyAddress     sdk.WasmAddress = bytes.Repeat([]byte{0x0}, SDKAddrLen)
		invalidAddress                 = "invalid address"
	)

	specs := map[string]struct {
		src    *StoreCodeProposal
		expErr bool
	}{
		"all good": {
			src: StoreCodeProposalFixture(),
		},
		"with instantiate permission": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				accessConfig := AccessTypeOnlyAddress.With(anyAddress)
				p.InstantiatePermission = &accessConfig
			}),
		},
		"base data missing": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"run_as missing": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.RunAs = ""
			}),
			expErr: true,
		},
		"run_as invalid": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.RunAs = invalidAddress
			}),
			expErr: true,
		},
		"wasm code missing": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = nil
			}),
			expErr: true,
		},
		"wasm code invalid": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = bytes.Repeat([]byte{0x0}, MaxWasmSize+1)
			}),
			expErr: true,
		},
		"with invalid instantiate permission": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.InstantiatePermission = &AccessConfig{}
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateInstantiateContractProposal(t *testing.T) {
	invalidAddress := "invalid address"

	specs := map[string]struct {
		src    *InstantiateContractProposal
		expErr bool
	}{
		"all good": {
			src: InstantiateContractProposalFixture(),
		},
		"without admin": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Admin = ""
			}),
		},
		"without init msg": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Msg = nil
			}),
			expErr: true,
		},
		"with invalid init msg": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Msg = []byte("not a json string")
			}),
			expErr: true,
		},
		"without init funds": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = nil
			}),
		},
		"base data missing": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"run_as missing": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.RunAs = ""
			}),
			expErr: true,
		},
		"run_as invalid": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.RunAs = invalidAddress
			}),
			expErr: true,
		},
		"admin invalid": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Admin = invalidAddress
			}),
			expErr: true,
		},
		"code id empty": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.CodeID = 0
			}),
			expErr: true,
		},
		"label empty": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Label = ""
			}),
			expErr: true,
		},
		"init funds negative": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.CoinAdapters{{Denom: "foo", Amount: sdk.NewInt(-1)}}
			}),
			expErr: true,
		},
		"init funds with duplicates": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.CoinAdapters{{Denom: "foo", Amount: sdk.NewInt(1)}, {Denom: "foo", Amount: sdk.NewInt(2)}}
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateMigrateContractProposal(t *testing.T) {
	invalidAddress := "invalid address2"

	specs := map[string]struct {
		src    *MigrateContractProposal
		expErr bool
	}{
		"all good": {
			src: MigrateContractProposalFixture(),
		},
		"without migrate msg": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Msg = nil
			}),
			expErr: true,
		},
		"migrate msg with invalid json": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Msg = []byte("not a json message")
			}),
			expErr: true,
		},
		"base data missing": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
		"code id empty": {
			src: MigrateContractProposalFixture(func(p *MigrateContractProposal) {
				p.CodeID = 0
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSudoContractProposal(t *testing.T) {
	invalidAddress := "invalid address"

	specs := map[string]struct {
		src    *SudoContractProposal
		expErr bool
	}{
		"all good": {
			src: SudoContractProposalFixture(),
		},
		"msg is nil": {
			src: SudoContractProposalFixture(func(p *SudoContractProposal) {
				p.Msg = nil
			}),
			expErr: true,
		},
		"msg with invalid json": {
			src: SudoContractProposalFixture(func(p *SudoContractProposal) {
				p.Msg = []byte("not a json message")
			}),
			expErr: true,
		},
		"base data missing": {
			src: SudoContractProposalFixture(func(p *SudoContractProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: SudoContractProposalFixture(func(p *SudoContractProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: SudoContractProposalFixture(func(p *SudoContractProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateExecuteContractProposal(t *testing.T) {
	invalidAddress := "invalid address"

	specs := map[string]struct {
		src    *ExecuteContractProposal
		expErr bool
	}{
		"all good": {
			src: ExecuteContractProposalFixture(),
		},
		"msg is nil": {
			src: ExecuteContractProposalFixture(func(p *ExecuteContractProposal) {
				p.Msg = nil
			}),
			expErr: true,
		},
		"msg with invalid json": {
			src: ExecuteContractProposalFixture(func(p *ExecuteContractProposal) {
				p.Msg = []byte("not a valid json message")
			}),
			expErr: true,
		},
		"base data missing": {
			src: ExecuteContractProposalFixture(func(p *ExecuteContractProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: ExecuteContractProposalFixture(func(p *ExecuteContractProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: ExecuteContractProposalFixture(func(p *ExecuteContractProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
		"run as is invalid": {
			src: ExecuteContractProposalFixture(func(p *ExecuteContractProposal) {
				p.RunAs = invalidAddress
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateAdminProposal(t *testing.T) {
	invalidAddress := "invalid address"

	specs := map[string]struct {
		src    *UpdateAdminProposal
		expErr bool
	}{
		"all good": {
			src: UpdateAdminProposalFixture(),
		},
		"base data missing": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
		"admin missing": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.NewAdmin = ""
			}),
			expErr: true,
		},
		"admin invalid": {
			src: UpdateAdminProposalFixture(func(p *UpdateAdminProposal) {
				p.NewAdmin = invalidAddress
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateClearAdminProposal(t *testing.T) {
	invalidAddress := "invalid address"

	specs := map[string]struct {
		src    *ClearAdminProposal
		expErr bool
	}{
		"all good": {
			src: ClearAdminProposalFixture(),
		},
		"base data missing": {
			src: ClearAdminProposalFixture(func(p *ClearAdminProposal) {
				p.Title = ""
			}),
			expErr: true,
		},
		"contract missing": {
			src: ClearAdminProposalFixture(func(p *ClearAdminProposal) {
				p.Contract = ""
			}),
			expErr: true,
		},
		"contract invalid": {
			src: ClearAdminProposalFixture(func(p *ClearAdminProposal) {
				p.Contract = invalidAddress
			}),
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProposalStrings(t *testing.T) {
	specs := map[string]struct {
		src govtypes.Content
		exp string
	}{
		"store code": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = []byte{0o1, 0o2, 0o3, 0o4, 0o5, 0o6, 0o7, 0x08, 0x09, 0x0a}
			}),
			exp: `Store Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      0x0101010101010101010101010101010101010101
  WasmCode:    0102030405060708090A
`,
		},
		"instantiate contract": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.CoinAdapters{{Denom: "foo", Amount: sdk.NewInt(1)}, {Denom: "bar", Amount: sdk.NewInt(2)}}
			}),
			exp: `Instantiate Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      0x0101010101010101010101010101010101010101
  Admin:       0x0101010101010101010101010101010101010101
  Code id:     1
  Label:       testing
  Msg:         "{\"verifier\":\"0x0101010101010101010101010101010101010101\",\"beneficiary\":\"0x0101010101010101010101010101010101010101\"}"
  Funds:       1foo,2bar
`,
		},
		"instantiate contract without funds": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Funds = nil }),
			exp: `Instantiate Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      0x0101010101010101010101010101010101010101
  Admin:       0x0101010101010101010101010101010101010101
  Code id:     1
  Label:       testing
  Msg:         "{\"verifier\":\"0x0101010101010101010101010101010101010101\",\"beneficiary\":\"0x0101010101010101010101010101010101010101\"}"
  Funds:       
`,
		},
		"instantiate contract without admin": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Admin = "" }),
			exp: `Instantiate Code Proposal:
  Title:       Foo
  Description: Bar
  Run as:      0x0101010101010101010101010101010101010101
  Admin:       
  Code id:     1
  Label:       testing
  Msg:         "{\"verifier\":\"0x0101010101010101010101010101010101010101\",\"beneficiary\":\"0x0101010101010101010101010101010101010101\"}"
  Funds:       
`,
		},
		"migrate contract": {
			src: MigrateContractProposalFixture(),
			exp: `Migrate Contract Proposal:
  Title:       Foo
  Description: Bar
  Contract:    0x5A8D648DEE57b2fc90D98DC17fa887159b69638b
  Code id:     1
  Msg:         "{\"verifier\":\"0x0101010101010101010101010101010101010101\"}"
`,
		},
		"update admin": {
			src: UpdateAdminProposalFixture(),
			exp: `Update Contract Admin Proposal:
  Title:       Foo
  Description: Bar
  Contract:    0x5A8D648DEE57b2fc90D98DC17fa887159b69638b
  New Admin:   0x0101010101010101010101010101010101010101
`,
		},
		"clear admin": {
			src: ClearAdminProposalFixture(),
			exp: `Clear Contract Admin Proposal:
  Title:       Foo
  Description: Bar
  Contract:    0x5A8D648DEE57b2fc90D98DC17fa887159b69638b
`,
		},
		"pin codes": {
			src: &PinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     []uint64{1, 2, 3},
			},
			exp: `Pin Wasm Codes Proposal:
  Title:       Foo
  Description: Bar
  Codes:       [1 2 3]
`,
		},
		"unpin codes": {
			src: &UnpinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     []uint64{3, 2, 1},
			},
			exp: `Unpin Wasm Codes Proposal:
  Title:       Foo
  Description: Bar
  Codes:       [3 2 1]
`,
		},
		"update deployment whitelist": {
			src: &UpdateDeploymentWhitelistProposal{
				Title:            "Foo",
				Description:      "Bar",
				DistributorAddrs: []string{"ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02", "ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf"},
			},
			exp: `title:"Foo" description:"Bar" distributorAddrs:"ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02" distributorAddrs:"ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf" `,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			assert.Equal(t, spec.exp, spec.src.String())
		})
	}
}

func TestProposalYaml(t *testing.T) {
	specs := map[string]struct {
		src govtypes.Content
		exp string
	}{
		"store code": {
			src: StoreCodeProposalFixture(func(p *StoreCodeProposal) {
				p.WASMByteCode = []byte{0o1, 0o2, 0o3, 0o4, 0o5, 0o6, 0o7, 0x08, 0x09, 0x0a}
			}),
			exp: `title: Foo
description: Bar
run_as: 0x0101010101010101010101010101010101010101
wasm_byte_code: AQIDBAUGBwgJCg==
instantiate_permission: null
`,
		},
		"instantiate contract": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) {
				p.Funds = sdk.CoinAdapters{{Denom: "foo", Amount: sdk.NewInt(1)}, {Denom: "bar", Amount: sdk.NewInt(2)}}
			}),
			exp: `title: Foo
description: Bar
run_as: 0x0101010101010101010101010101010101010101
admin: 0x0101010101010101010101010101010101010101
code_id: 1
label: testing
msg: '{"verifier":"0x0101010101010101010101010101010101010101","beneficiary":"0x0101010101010101010101010101010101010101"}'
funds:
- denom: foo
  amount: "0.000000000000000001"
- denom: bar
  amount: "0.000000000000000002"
`,
		},
		"instantiate contract without funds": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Funds = nil }),
			exp: `title: Foo
description: Bar
run_as: 0x0101010101010101010101010101010101010101
admin: 0x0101010101010101010101010101010101010101
code_id: 1
label: testing
msg: '{"verifier":"0x0101010101010101010101010101010101010101","beneficiary":"0x0101010101010101010101010101010101010101"}'
funds: []
`,
		},
		"instantiate contract without admin": {
			src: InstantiateContractProposalFixture(func(p *InstantiateContractProposal) { p.Admin = "" }),
			exp: `title: Foo
description: Bar
run_as: 0x0101010101010101010101010101010101010101
admin: ""
code_id: 1
label: testing
msg: '{"verifier":"0x0101010101010101010101010101010101010101","beneficiary":"0x0101010101010101010101010101010101010101"}'
funds: []
`,
		},
		"migrate contract": {
			src: MigrateContractProposalFixture(),
			exp: `title: Foo
description: Bar
contract: 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b
code_id: 1
msg: '{"verifier":"0x0101010101010101010101010101010101010101"}'
`,
		},
		"update admin": {
			src: UpdateAdminProposalFixture(),
			exp: `title: Foo
description: Bar
new_admin: 0x0101010101010101010101010101010101010101
contract: 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b
`,
		},
		"clear admin": {
			src: ClearAdminProposalFixture(),
			exp: `title: Foo
description: Bar
contract: 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b
`,
		},
		"pin codes": {
			src: &PinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     []uint64{1, 2, 3},
			},
			exp: `title: Foo
description: Bar
code_ids:
- 1
- 2
- 3
`,
		},
		"update deployment whitelist": {
			src: &UpdateDeploymentWhitelistProposal{
				Title:            "Foo",
				Description:      "Bar",
				DistributorAddrs: []string{"ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02", "ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf"},
			},
			exp: `title: Foo
description: Bar
distributor_addresses:
- ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02
- ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf
`,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			v, err := yaml.Marshal(&spec.src)
			require.NoError(t, err)
			assert.Equal(t, spec.exp, string(v))
		})
	}
}

func TestConvertToProposals(t *testing.T) {
	cases := map[string]struct {
		input     string
		isError   bool
		proposals []ProposalType
	}{
		"one proper item": {
			input:     "UpdateAdmin",
			proposals: []ProposalType{ProposalTypeUpdateAdmin},
		},
		"multiple proper items": {
			input:     "StoreCode,InstantiateContract,MigrateContract,UpdateDeploymentWhitelist",
			proposals: []ProposalType{ProposalTypeStoreCode, ProposalTypeInstantiateContract, ProposalTypeMigrateContract, ProposalTypeUpdateDeploymentWhitelist},
		},
		"empty trailing item": {
			input:   "StoreCode,",
			isError: true,
		},
		"invalid item": {
			input:   "StoreCode,InvalidProposalType",
			isError: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			chunks := strings.Split(tc.input, ",")
			proposals, err := ConvertToProposals(chunks)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, proposals, tc.proposals)
			}
		})
	}
}

func TestUnmarshalContentFromJson(t *testing.T) {
	specs := map[string]struct {
		src string
		got govtypes.Content
		exp govtypes.Content
	}{
		"instantiate ": {
			src: `
{
	"title": "foo",
	"description": "bar",
	"admin": "myAdminAddress",
	"code_id": 1,
	"funds": [{"denom": "ALX", "amount": "2"},{"denom": "BLX","amount": "3"}],
	"msg": {},
	"label": "testing",
	"run_as": "myRunAsAddress"
}`,
			got: &InstantiateContractProposal{},
			exp: &InstantiateContractProposal{
				Title:       "foo",
				Description: "bar",
				RunAs:       "myRunAsAddress",
				Admin:       "myAdminAddress",
				CodeID:      1,
				Label:       "testing",
				Msg:         []byte("{}"),
				Funds:       sdk.CoinAdapters{{"ALX", sdk.NewInt(2)}, {"BLX", sdk.NewInt(3)}},
			},
		},
		"migrate ": {
			src: `
{
	"title": "foo",
	"description": "bar",
	"code_id": 1,
	"contract": "myContractAddr",
	"msg": {},
	"run_as": "myRunAsAddress"
}`,
			got: &MigrateContractProposal{},
			exp: &MigrateContractProposal{
				Title:       "foo",
				Description: "bar",
				Contract:    "myContractAddr",
				CodeID:      1,
				Msg:         []byte("{}"),
			},
		},
		"update deployment whitelit": {
			src: `
{
	"title": "foo",
	"description": "bar",
	"distributorAddrs": ["ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02", "ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf"]
}`,
			got: &UpdateDeploymentWhitelistProposal{},
			exp: &UpdateDeploymentWhitelistProposal{
				Title:            "foo",
				Description:      "bar",
				DistributorAddrs: []string{"ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02", "ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf"},
			},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(spec.src), spec.got))
			assert.Equal(t, spec.exp, spec.got)
		})
	}
}

func TestProposalJsonSignBytes(t *testing.T) {
	const myInnerMsg = `{"foo":"bar"}`
	specs := map[string]struct {
		src govtypes.Content
		exp string
	}{
		"instantiate contract": {
			src: &InstantiateContractProposal{Msg: RawContractMessage(myInnerMsg)},
			exp: `
{
	"type":"okexchain/gov/MsgSubmitProposal",
	"value":{"content":{"type":"wasm/InstantiateContractProposal","value":{"funds":[],"msg":{"foo":"bar"}}},"initial_deposit":[],"proposer":""}
}`,
		},
		"migrate contract": {
			src: &MigrateContractProposal{Msg: RawContractMessage(myInnerMsg)},
			exp: `
{
	"type":"okexchain/gov/MsgSubmitProposal",
	"value":{"content":{"type":"wasm/MigrateContractProposal","value":{"msg":{"foo":"bar"}}},"initial_deposit":[],"proposer":""}
}`,
		},
		"update wasm deployment whitelist": {
			src: &UpdateDeploymentWhitelistProposal{DistributorAddrs: []string{"ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02", "ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf"}},
			exp: `
{	
	"type":"okexchain/gov/MsgSubmitProposal",
	"value":{"content":{"type":"wasm/UpdateDeploymentWhitelistProposal","value":{"distributorAddrs":["ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02","ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf"]}},"initial_deposit":[],"proposer":""}
}`,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			msg := govtypes.NewMsgSubmitProposal(spec.src, sdk.NewCoins(), []byte{})

			bz := msg.GetSignBytes()
			assert.JSONEq(t, spec.exp, string(bz), "raw: %s", string(bz))
		})
	}
}

type ProposalSuite struct {
	suite.Suite
}

func TestProposalSuite(t *testing.T) {
	suite.Run(t, new(ProposalSuite))
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func (suite *ProposalSuite) TestNewChangeDistributionTypeProposal() {
	testCases := []struct {
		name        string
		title       string
		description string
		action      string
		extra       string
		err         error
	}{
		{
			"no proposal title",
			"",
			"description",
			"",
			"",
			govtypes.ErrInvalidProposalContent("proposal title cannot be blank"),
		},
		{
			"gt max proposal title length",
			RandStr(types.MaxTitleLength + 1),
			"description",
			"",
			"",
			govtypes.ErrInvalidProposalContent(fmt.Sprintf("proposal title is longer than max length of %d", govtypes.MaxTitleLength)),
		},
		{
			"gt max proposal title length",
			RandStr(types.MaxTitleLength),
			"",
			"",
			"",
			govtypes.ErrInvalidProposalContent("proposal description cannot be blank"),
		},
		{
			"gt max proposal description length",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength + 1),
			"",
			"",
			govtypes.ErrInvalidProposalContent(fmt.Sprintf("proposal description is longer than max length of %d", govtypes.MaxDescriptionLength)),
		},
		{
			"no action",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			"",
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's action is required"),
		},
		{
			"action too large",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxExtraActionLength + 1),
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's action length is bigger than max length"),
		},
		{
			"no extra body",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxExtraActionLength),
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's extra is required"),
		},
		{
			"extra too large",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxTitleLength),
			RandStr(govtypes.MaxExtraBodyLength + 1),
			govtypes.ErrInvalidProposalContent("extra proposal's extra body length is bigger than max length"),
		},
		{
			"unknown action",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxTitleLength),
			RandStr(govtypes.MaxExtraBodyLength),
			ErrUnknownExtraProposalAction,
		},
		{
			"ActionModifyGasFactor, parse error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{dfafdasf}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionModifyGasFactor, action is nil",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			"",
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's action is required"),
		},
		{
			"ActionModifyGasFactor, extra is nil",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			"hello",
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's extra is required"),
		},
		{
			"ActionModifyGasFactor, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			"hello",
			"hello",
			ErrUnknownExtraProposalAction,
		},
		{
			"ActionModifyGasFactor, extra is nil",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's extra is required"),
		},
		{
			"ActionModifyGasFactor, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{}",
			ErrExtraProposalParams("parse factor error:"),
		},
		{
			"ActionModifyGasFactor, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"\"}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionModifyGasFactor, key error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"df\": \"\"}",
			ErrExtraProposalParams("parse factor error:"),
		},
		{
			"ActionModifyGasFactor, value error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\":19.7}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionModifyGasFactor, value error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\":19}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionModifyGasFactor, value error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\": \"adfasd\"}",
			ErrExtraProposalParams("parse factor error:adfasd"),
		},
		{
			"ActionModifyGasFactor, value -1",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\": \"-1\"}",
			ErrExtraProposalParams("parse factor error:-1"),
		},
		{
			"ActionModifyGasFactor, value 0",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\": \"0\"}",
			ErrExtraProposalParams("parse factor error:0"),
		},
		{
			"ActionModifyGasFactor, value > 18",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\": \"0.0000000000000000001\"}",
			ErrExtraProposalParams("parse factor error:0.0000000000000000001"),
		},
		{
			"ActionModifyGasFactor, value = 18",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\": \"0.000000000000000001\"}",
			nil,
		},
		{
			"ActionModifyGasFactor, value ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\": \"10000001\"}",
			ErrExtraProposalParams(fmt.Sprintf("max gas factor:%v", MaxGasFactor)),
		},
		{
			"ActionModifyGasFactor, value ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\": \"10000000\"}",
			nil,
		},
		{
			"ActionModifyGasFactor, value error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\":\"19.7a\"}",
			ErrExtraProposalParams("parse factor error:19.7a"),
		},
		{
			"ActionModifyGasFactor, value error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\":\"a19.7\"}",
			ErrExtraProposalParams("parse factor error:a19.7"),
		},
		{
			"ActionModifyGasFactor, value ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\":\"19.7\"}",
			nil,
		},
		{
			"ActionModifyGasFactor, value ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\":\"19.6757657657657567864554354354357\"}",
			ErrExtraProposalParams(fmt.Sprintf("parse factor error:19.6757657657657567864554354354357")),
		},
		{
			"ActionModifyGasFactor, value ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionModifyGasFactor,
			"{\"factor\":\"19.675765765767\"}",
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			global.SetGlobalHeight(100)
			proposal := ExtraProposal{
				Title:       tc.title,
				Description: tc.description,
				Action:      tc.action,
				Extra:       tc.extra,
			}

			require.Equal(suite.T(), tc.title, proposal.GetTitle())
			require.Equal(suite.T(), tc.description, proposal.GetDescription())
			require.Equal(suite.T(), RouterKey, proposal.ProposalRoute())
			require.Equal(suite.T(), string(ProposalTypeExtra), proposal.ProposalType())
			require.NotPanics(suite.T(), func() {
				_ = proposal.String()
			})

			err := proposal.ValidateBasic()
			require.Equal(suite.T(), tc.err, err)
		})
	}
}

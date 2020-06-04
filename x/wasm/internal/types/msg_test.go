package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilderRegexp(t *testing.T) {
	cases := map[string]struct {
		example string
		valid   bool
	}{
		"normal":                {"fedora/httpd:version1.0", true},
		"another valid org":     {"confio/js-builder-1:test", true},
		"no org name":           {"cosmwasm-opt:0.6.3", false},
		"invalid trailing char": {"someone/cosmwasm-opt-:0.6.3", false},
		"invalid leading char":  {"confio/.builder-1:manual", false},
		"multiple orgs":         {"confio/assembly-script/optimizer:v0.9.1", true},
		"too long":              {"over-128-character-limit/some-long-sub-path/and-yet-another-long-name/testtesttesttesttesttesttest/foobarfoobar/foobarfoobar:randomstringrandomstringrandomstringrandomstring", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateBuilder(tc.example)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})

	}
}

func TestStoreCodeValidation(t *testing.T) {
	badAddress, err := sdk.AccAddressFromHex("012345")
	require.NoError(t, err)
	// proper address size
	goodAddress := sdk.AccAddress(make([]byte, 20))

	cases := map[string]struct {
		msg   MsgStoreCode
		valid bool
	}{
		"empty": {
			msg:   MsgStoreCode{},
			valid: false,
		},
		"correct minimal": {
			msg: MsgStoreCode{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
			},
			valid: true,
		},
		"missing code": {
			msg: MsgStoreCode{
				Sender: goodAddress,
			},
			valid: false,
		},
		"bad sender minimal": {
			msg: MsgStoreCode{
				Sender:       badAddress,
				WASMByteCode: []byte("foo"),
			},
			valid: false,
		},
		"correct maximal": {
			msg: MsgStoreCode{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Builder:      "confio/cosmwasm-opt:0.6.2",
				Source:       "https://crates.io/api/v1/crates/cw-erc20/0.1.0/download",
			},
			valid: true,
		},
		"invalid builder": {
			msg: MsgStoreCode{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Builder:      "-bad-opt:0.6.2",
				Source:       "https://crates.io/api/v1/crates/cw-erc20/0.1.0/download",
			},
			valid: false,
		},
		"invalid source scheme": {
			msg: MsgStoreCode{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Builder:      "cosmwasm-opt:0.6.2",
				Source:       "ftp://crates.io/api/download.tar.gz",
			},
			valid: false,
		},
		"invalid source format": {
			msg: MsgStoreCode{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Builder:      "cosmwasm-opt:0.6.2",
				Source:       "/api/download-ss",
			},
			valid: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}

}

func TestInstantiateContractValidation(t *testing.T) {
	badAddress, err := sdk.AccAddressFromHex("012345")
	require.NoError(t, err)
	// proper address size
	goodAddress := sdk.AccAddress(make([]byte, 20))

	cases := map[string]struct {
		msg   MsgInstantiateContract
		valid bool
	}{
		"empty": {
			msg:   MsgInstantiateContract{},
			valid: false,
		},
		"correct minimal": {
			msg: MsgInstantiateContract{
				Sender:  goodAddress,
				Code:    1,
				Label:   "foo",
				InitMsg: []byte("{}"),
			},
			valid: true,
		},
		"missing code": {
			msg: MsgInstantiateContract{
				Sender:  goodAddress,
				Label:   "foo",
				InitMsg: []byte("{}"),
			},
			valid: false,
		},
		"missing label": {
			msg: MsgInstantiateContract{
				Sender:  goodAddress,
				InitMsg: []byte("{}"),
			},
			valid: false,
		},
		"label too long": {
			msg: MsgInstantiateContract{
				Sender: goodAddress,
				Label:  strings.Repeat("food", 33),
			},
			valid: false,
		},
		"bad sender minimal": {
			msg: MsgInstantiateContract{
				Sender:  badAddress,
				Code:    1,
				Label:   "foo",
				InitMsg: []byte("{}"),
			},
			valid: false,
		},
		"correct maximal": {
			msg: MsgInstantiateContract{
				Sender:    goodAddress,
				Code:      1,
				Label:     "foo",
				InitMsg:   []byte(`{"some": "data"}`),
				InitFunds: sdk.Coins{sdk.Coin{Denom: "foobar", Amount: sdk.MustNewDecFromStr("200")}},
			},
			valid: true,
		},
		"negative funds": {
			msg: MsgInstantiateContract{
				Sender:  goodAddress,
				Code:    1,
				Label:   "foo",
				InitMsg: []byte(`{"some": "data"}`),
				// we cannot use sdk.NewCoin() constructors as they panic on creating invalid data (before we can test)
				InitFunds: sdk.Coins{sdk.Coin{Denom: "foobar", Amount: sdk.MustNewDecFromStr("-200")}},
			},
			valid: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}

}

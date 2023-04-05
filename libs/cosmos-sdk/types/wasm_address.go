package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/tendermint/libs/bech32"
	"gopkg.in/yaml.v2"
	"strings"
)

var _ Address = AccAddress{}

// WasmAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type WasmAddress []byte

// WasmAddressFromHex creates an WasmAddress from a hex string.
func WasmAddressFromHex(address string) (addr WasmAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return WasmAddress(bz), nil
}

// MustWasmAddressFromBech32 calls WasmAddressFromBech32 and panics on error.
func MustWasmAddressFromBech32(address string) WasmAddress {
	addr, err := WasmAddressFromBech32(address)
	if err != nil {
		panic(err)
	}

	return addr
}

func WasmToAccAddress(addr WasmAddress) AccAddress {
	return AccAddress(addr)
}
func AccToAWasmddress(addr AccAddress) WasmAddress {
	return WasmAddress(addr)
}

// WasmAddressFromBech32 creates an WasmAddress from a Bech32 string.
func WasmAddressFromBech32(address string) (WasmAddress, error) {
	return WasmAddressFromBech32ByPrefix(address, GetConfig().GetBech32AccountAddrPrefix())
}

// WasmAddressFromBech32ByPrefix create an WasmAddress from a Bech32 string by address prefix
func WasmAddressFromBech32ByPrefix(address string, bech32PrefixAccAddr string) (addr WasmAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return nil, errors.New("empty address string is not allowed")
	}

	if !strings.HasPrefix(address, bech32PrefixAccAddr) {
		// strip 0x prefix if exists
		addrStr := strings.TrimPrefix(address, "0x")
		addr, err = WasmAddressFromHex(addrStr)
		if err != nil {
			return addr, err
		}
		return addr, VerifyAddressFormat(addr)
	}

	//decodes a bytestring from a Bech32 encoded string
	bz, err := GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}

	err = VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return WasmAddress(bz), nil
}

// Returns boolean for whether two WasmAddresses are Equal
func (aa WasmAddress) Equals(aa2 Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Returns boolean for whether an WasmAddress is empty
func (aa WasmAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := WasmAddress{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa WasmAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *WasmAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa WasmAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa WasmAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *WasmAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	if s == "" {
		*aa = WasmAddress{}
		return nil
	}

	aa2, err := WasmAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *WasmAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	if s == "" {
		*aa = WasmAddress{}
		return nil
	}

	aa2, err := WasmAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa WasmAddress) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa WasmAddress) String() string {
	if aa.Empty() {
		return ""
	}

	return common.BytesToAddress(aa).String()
}

// Bech32String convert account address to bech32 address.
func (aa WasmAddress) Bech32String(bech32PrefixAccAddr string) string {
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixAccAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa WasmAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

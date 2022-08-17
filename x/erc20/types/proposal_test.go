package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTokenMappingProposal_ValidateBasic(t *testing.T) {
	contractAddrStr := "0x7D4B7B8CA7E1a24928Bb96D59249c7a5bd1DfBe6"
	contractAddr := common.HexToAddress(contractAddrStr)

	proposal := NewTokenMappingProposal("proposal", "right delist proposal", "eth", &contractAddr)
	require.Equal(t, "proposal", proposal.GetTitle())
	require.Equal(t, "right delist proposal", proposal.GetDescription())
	require.Equal(t, RouterKey, proposal.ProposalRoute())
	require.Equal(t, proposalTypeTokenMapping, proposal.ProposalType())

	tests := []struct {
		name   string
		drp    TokenMappingProposal
		result bool
	}{
		{"valid-proposal", proposal, true},
		{"no-title", TokenMappingProposal{"", "delist proposal", "eth", contractAddrStr}, false},
		{"no-description", TokenMappingProposal{"proposal", "", "eth", contractAddrStr}, false},
		{"no-denom", TokenMappingProposal{"proposal", "delist proposal", "", contractAddrStr}, false},
		{"err-denom", TokenMappingProposal{"proposal", "delist proposal", ".@..", contractAddrStr}, false},
		{"no-contract", TokenMappingProposal{"proposal", "delist proposal", "btc", ""}, true},
		{"err-contract", TokenMappingProposal{"proposal", "delist proposal", "btc", "0xqwoifej923jd"}, false},
		{"long-title", TokenMappingProposal{getLongString(15),
			"right delist proposal", "eth", contractAddrStr}, false},
		{"long-description", TokenMappingProposal{"proposal",
			getLongString(501), "eth", contractAddrStr}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.result {
				require.Nil(t, tt.drp.ValidateBasic(), "test: %v", tt.name)
			} else {
				require.NotNil(t, tt.drp.ValidateBasic(), "test: %v", tt.name)
			}
		})
	}
}

func TestContractTemplateProposal_ValidateBasic(t *testing.T) {
	f := func(bin string) string {
		json := `[{	"inputs": [{"internalType": "uint256","name": "a","type": "uint256"	},{	"internalType": "uint256","name": "b","type": "uint256"}],"stateMutability": "nonpayable","type": "constructor"}]`
		str := fmt.Sprintf(`{"abi":%s,"bin":"%s"}`, json, bin)
		return str
	}
	tests := []struct {
		name   string
		drp    ContractTemplateProposal
		result bool
	}{
		{"invalid hex", ContractTemplateProposal{"title", "desc", ProposalTypeContextTemplateProxy, f("invalid hex")}, false},
		{"invalid type", ContractTemplateProposal{"title", "desc", "mis", f(hex.EncodeToString([]byte("valid hex")))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.result {
				require.Nil(t, tt.drp.ValidateBasic(), "test: %v", tt.name)
			} else {
				require.NotNil(t, tt.drp.ValidateBasic(), "test: %v", tt.name)
			}
		})
	}
}

func getLongString(n int) (s string) {
	str := "0123456789"
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%s", s, str)
	}
	fmt.Println(len(s))
	return s
}

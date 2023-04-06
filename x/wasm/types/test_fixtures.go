package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"math/rand"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func GenesisFixture(mutators ...func(*GenesisState)) GenesisState {
	const (
		numCodes     = 2
		numContracts = 2
		numSequences = 2
		numMsg       = 3
	)

	fixture := GenesisState{
		Params:    DefaultParams(),
		Codes:     make([]Code, numCodes),
		Contracts: make([]Contract, numContracts),
		Sequences: make([]Sequence, numSequences),
	}
	for i := 0; i < numCodes; i++ {
		fixture.Codes[i] = CodeFixture()
	}
	for i := 0; i < numContracts; i++ {
		fixture.Contracts[i] = ContractFixture()
	}
	for i := 0; i < numSequences; i++ {
		fixture.Sequences[i] = Sequence{
			IDKey: randBytes(5),
			Value: uint64(i),
		}
	}
	fixture.GenMsgs = []GenesisState_GenMsgs{
		{Sum: &GenesisState_GenMsgs_StoreCode{StoreCode: MsgStoreCodeFixture()}},
		{Sum: &GenesisState_GenMsgs_InstantiateContract{InstantiateContract: MsgInstantiateContractFixture()}},
		{Sum: &GenesisState_GenMsgs_ExecuteContract{ExecuteContract: MsgExecuteContractFixture()}},
	}
	for _, m := range mutators {
		m(&fixture)
	}
	return fixture
}

func randBytes(n int) []byte {
	r := make([]byte, n)
	rand.Read(r)
	return r
}

func CodeFixture(mutators ...func(*Code)) Code {
	wasmCode := randBytes(100)

	fixture := Code{
		CodeID:    1,
		CodeInfo:  CodeInfoFixture(WithSHA256CodeHash(wasmCode)),
		CodeBytes: wasmCode,
	}

	for _, m := range mutators {
		m(&fixture)
	}
	return fixture
}

func CodeInfoFixture(mutators ...func(*CodeInfo)) CodeInfo {
	wasmCode := bytes.Repeat([]byte{0x1}, 10)
	codeHash := sha256.Sum256(wasmCode)
	const anyAddress = "0x0101010101010101010101010101010101010101"
	fixture := CodeInfo{
		CodeHash:          codeHash[:],
		Creator:           anyAddress,
		InstantiateConfig: AllowEverybody,
	}
	for _, m := range mutators {
		m(&fixture)
	}
	return fixture
}

func ContractFixture(mutators ...func(*Contract)) Contract {
	const anyAddress = "0x0101010101010101010101010101010101010101"

	fixture := Contract{
		ContractAddress: anyAddress,
		ContractInfo:    ContractInfoFixture(OnlyGenesisFields),
		ContractState:   []Model{{Key: []byte("anyKey"), Value: []byte("anyValue")}},
	}

	for _, m := range mutators {
		m(&fixture)
	}
	return fixture
}

func OnlyGenesisFields(info *ContractInfo) {
	info.Created = nil
}

func ContractInfoFixture(mutators ...func(*ContractInfo)) ContractInfo {
	const anyAddress = "0x0101010101010101010101010101010101010101"

	fixture := ContractInfo{
		CodeID:  1,
		Creator: anyAddress,
		Label:   "any",
		Created: &AbsoluteTxPosition{BlockHeight: 1, TxIndex: 1},
	}

	for _, m := range mutators {
		m(&fixture)
	}
	return fixture
}

func WithSHA256CodeHash(wasmCode []byte) func(info *CodeInfo) {
	return func(info *CodeInfo) {
		codeHash := sha256.Sum256(wasmCode)
		info.CodeHash = codeHash[:]
	}
}

func MsgStoreCodeFixture(mutators ...func(*MsgStoreCode)) *MsgStoreCode {
	wasmIdent := []byte("\x00\x61\x73\x6D")
	const anyAddress = "0x0101010101010101010101010101010101010101"
	r := &MsgStoreCode{
		Sender:                anyAddress,
		WASMByteCode:          wasmIdent,
		InstantiatePermission: &AllowEverybody,
	}
	for _, m := range mutators {
		m(r)
	}
	return r
}

func MsgInstantiateContractFixture(mutators ...func(*MsgInstantiateContract)) *MsgInstantiateContract {
	const anyAddress = "0x0101010101010101010101010101010101010101"
	r := &MsgInstantiateContract{
		Sender: anyAddress,
		Admin:  anyAddress,
		CodeID: 1,
		Label:  "testing",
		Msg:    []byte(`{"foo":"bar"}`),
		Funds: sdk.CoinAdapters{{
			Denom:  "stake",
			Amount: sdk.NewInt(1),
		}},
	}
	for _, m := range mutators {
		m(r)
	}
	return r
}

func MsgExecuteContractFixture(mutators ...func(*MsgExecuteContract)) *MsgExecuteContract {
	const (
		anyAddress           = "0x0101010101010101010101010101010101010101"
		firstContractAddress = "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"
	)
	r := &MsgExecuteContract{
		Sender:   anyAddress,
		Contract: firstContractAddress,
		Msg:      []byte(`{"do":"something"}`),
		Funds: sdk.CoinAdapters{{
			Denom:  "stake",
			Amount: sdk.NewInt(1),
		}},
	}
	for _, m := range mutators {
		m(r)
	}
	return r
}

func StoreCodeProposalFixture(mutators ...func(*StoreCodeProposal)) *StoreCodeProposal {
	const anyAddress = "0x0101010101010101010101010101010101010101"
	p := &StoreCodeProposal{
		Title:        "Foo",
		Description:  "Bar",
		RunAs:        anyAddress,
		WASMByteCode: []byte{0x0},
	}
	for _, m := range mutators {
		m(p)
	}
	return p
}

func InstantiateContractProposalFixture(mutators ...func(p *InstantiateContractProposal)) *InstantiateContractProposal {
	var (
		anyValidAddress sdk.WasmAddress = bytes.Repeat([]byte{0x1}, SDKAddrLen)

		initMsg = struct {
			Verifier    sdk.WasmAddress `json:"verifier"`
			Beneficiary sdk.WasmAddress `json:"beneficiary"`
		}{
			Verifier:    anyValidAddress,
			Beneficiary: anyValidAddress,
		}
	)
	const anyAddress = "0x0101010101010101010101010101010101010101"

	initMsgBz, err := json.Marshal(initMsg)
	if err != nil {
		panic(err)
	}
	p := &InstantiateContractProposal{
		Title:       "Foo",
		Description: "Bar",
		RunAs:       anyAddress,
		Admin:       anyAddress,
		CodeID:      1,
		Label:       "testing",
		Msg:         initMsgBz,
		Funds:       nil,
	}

	for _, m := range mutators {
		m(p)
	}
	return p
}

func MigrateContractProposalFixture(mutators ...func(p *MigrateContractProposal)) *MigrateContractProposal {
	var (
		anyValidAddress sdk.WasmAddress = bytes.Repeat([]byte{0x1}, SDKAddrLen)

		migMsg = struct {
			Verifier sdk.WasmAddress `json:"verifier"`
		}{Verifier: anyValidAddress}
	)

	migMsgBz, err := json.Marshal(migMsg)
	if err != nil {
		panic(err)
	}
	const (
		contractAddr = "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"
		anyAddress   = "0x0101010101010101010101010101010101010101"
	)
	p := &MigrateContractProposal{
		Title:       "Foo",
		Description: "Bar",
		Contract:    contractAddr,
		CodeID:      1,
		Msg:         migMsgBz,
	}

	for _, m := range mutators {
		m(p)
	}
	return p
}

func SudoContractProposalFixture(mutators ...func(p *SudoContractProposal)) *SudoContractProposal {
	const (
		contractAddr = "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"
	)

	p := &SudoContractProposal{
		Title:       "Foo",
		Description: "Bar",
		Contract:    contractAddr,
		Msg:         []byte(`{"do":"something"}`),
	}

	for _, m := range mutators {
		m(p)
	}
	return p
}

func ExecuteContractProposalFixture(mutators ...func(p *ExecuteContractProposal)) *ExecuteContractProposal {
	const (
		contractAddr = "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"
		anyAddress   = "0x0101010101010101010101010101010101010101"
	)

	p := &ExecuteContractProposal{
		Title:       "Foo",
		Description: "Bar",
		Contract:    contractAddr,
		RunAs:       anyAddress,
		Msg:         []byte(`{"do":"something"}`),
		Funds: sdk.CoinsToCoinAdapters(sdk.Coins{{
			Denom:  "stake",
			Amount: sdk.NewDec(1),
		}}),
	}

	for _, m := range mutators {
		m(p)
	}
	return p
}

func UpdateAdminProposalFixture(mutators ...func(p *UpdateAdminProposal)) *UpdateAdminProposal {
	const (
		contractAddr = "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"
		anyAddress   = "0x0101010101010101010101010101010101010101"
	)

	p := &UpdateAdminProposal{
		Title:       "Foo",
		Description: "Bar",
		NewAdmin:    anyAddress,
		Contract:    contractAddr,
	}
	for _, m := range mutators {
		m(p)
	}
	return p
}

func ClearAdminProposalFixture(mutators ...func(p *ClearAdminProposal)) *ClearAdminProposal {
	const contractAddr = "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"
	p := &ClearAdminProposal{
		Title:       "Foo",
		Description: "Bar",
		Contract:    contractAddr,
	}
	for _, m := range mutators {
		m(p)
	}
	return p
}

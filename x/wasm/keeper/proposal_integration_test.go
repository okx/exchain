package keeper

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	wasmvm "github.com/CosmWasm/wasmvm"

	"github.com/CosmWasm/wasmd/x/wasm/keeper/wasmtesting"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CosmWasm/wasmd/x/wasm/types"
)

func TestStoreCodeProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper
	wasmKeeper.SetParams(ctx, types.Params{
		CodeUploadAccess:             types.AllowNobody,
		InstantiateDefaultPermission: types.AccessTypeNobody,
		MaxWasmCodeSize:              types.DefaultMaxWasmCodeSize,
	})
	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	myActorAddress := RandomBech32AccountAddress(t)

	src := types.StoreCodeProposalFixture(func(p *types.StoreCodeProposal) {
		p.RunAs = myActorAddress
		p.WASMByteCode = wasmCode
	})

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(ctx, src)
	require.NoError(t, err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx, storedProposal.GetContent())
	require.NoError(t, err)

	// then
	cInfo := wasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)
	assert.Equal(t, myActorAddress, cInfo.Creator)
	assert.True(t, wasmKeeper.IsPinnedCode(ctx, 1))

	storedCode, err := wasmKeeper.GetByteCode(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, wasmCode, storedCode)
}

func TestInstantiateProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper
	wasmKeeper.SetParams(ctx, types.Params{
		CodeUploadAccess:             types.AllowNobody,
		InstantiateDefaultPermission: types.AccessTypeNobody,
		MaxWasmCodeSize:              types.DefaultMaxWasmCodeSize,
	})

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	require.NoError(t, wasmKeeper.importCode(ctx, 1,
		types.CodeInfoFixture(types.WithSHA256CodeHash(wasmCode)),
		wasmCode),
	)

	var (
		oneAddress   sdk.AccAddress = bytes.Repeat([]byte{0x1}, types.ContractAddrLen)
		otherAddress sdk.AccAddress = bytes.Repeat([]byte{0x2}, types.ContractAddrLen)
	)
	src := types.InstantiateContractProposalFixture(func(p *types.InstantiateContractProposal) {
		p.CodeID = firstCodeID
		p.RunAs = oneAddress.String()
		p.Admin = otherAddress.String()
		p.Label = "testing"
	})
	em := sdk.NewEventManager()

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(ctx, src)
	require.NoError(t, err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx.WithEventManager(em), storedProposal.GetContent())
	require.NoError(t, err)

	// then
	contractAddr, err := sdk.AccAddressFromBech32("cosmos14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s4hmalr")
	require.NoError(t, err)

	cInfo := wasmKeeper.GetContractInfo(ctx, contractAddr)
	require.NotNil(t, cInfo)
	assert.Equal(t, uint64(1), cInfo.CodeID)
	assert.Equal(t, oneAddress.String(), cInfo.Creator)
	assert.Equal(t, otherAddress.String(), cInfo.Admin)
	assert.Equal(t, "testing", cInfo.Label)
	expHistory := []types.ContractCodeHistoryEntry{{
		Operation: types.ContractCodeHistoryOperationTypeInit,
		CodeID:    src.CodeID,
		Updated:   types.NewAbsoluteTxPosition(ctx),
		Msg:       src.Msg,
	}}
	assert.Equal(t, expHistory, wasmKeeper.GetContractHistory(ctx, contractAddr))
	// and event
	require.Len(t, em.Events(), 3, "%#v", em.Events())
	require.Equal(t, types.EventTypeInstantiate, em.Events()[0].Type)
	require.Equal(t, types.WasmModuleEventType, em.Events()[1].Type)
	require.Equal(t, types.EventTypeGovContractResult, em.Events()[2].Type)
	require.Len(t, em.Events()[2].Attributes, 1)
	require.NotEmpty(t, em.Events()[2].Attributes[0])
}

func TestMigrateProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper
	wasmKeeper.SetParams(ctx, types.Params{
		CodeUploadAccess:             types.AllowNobody,
		InstantiateDefaultPermission: types.AccessTypeNobody,
		MaxWasmCodeSize:              types.DefaultMaxWasmCodeSize,
	})

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeInfoFixture := types.CodeInfoFixture(types.WithSHA256CodeHash(wasmCode))
	require.NoError(t, wasmKeeper.importCode(ctx, 1, codeInfoFixture, wasmCode))
	require.NoError(t, wasmKeeper.importCode(ctx, 2, codeInfoFixture, wasmCode))

	var (
		anyAddress   sdk.AccAddress = bytes.Repeat([]byte{0x1}, types.ContractAddrLen)
		otherAddress sdk.AccAddress = bytes.Repeat([]byte{0x2}, types.ContractAddrLen)
		contractAddr                = BuildContractAddress(1, 1)
	)

	contractInfoFixture := types.ContractInfoFixture(func(c *types.ContractInfo) {
		c.Label = "testing"
		c.Admin = anyAddress.String()
	})
	key, err := hex.DecodeString("636F6E666967")
	require.NoError(t, err)
	m := types.Model{Key: key, Value: []byte(`{"verifier":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","beneficiary":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","funder":"AQEBAQEBAQEBAQEBAQEBAQEBAQE="}`)}
	require.NoError(t, wasmKeeper.importContract(ctx, contractAddr, &contractInfoFixture, []types.Model{m}))

	migMsg := struct {
		Verifier sdk.AccAddress `json:"verifier"`
	}{Verifier: otherAddress}
	migMsgBz, err := json.Marshal(migMsg)
	require.NoError(t, err)

	src := types.MigrateContractProposal{
		Title:       "Foo",
		Description: "Bar",
		CodeID:      2,
		Contract:    contractAddr.String(),
		Msg:         migMsgBz,
	}

	em := sdk.NewEventManager()

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(ctx, &src)
	require.NoError(t, err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx.WithEventManager(em), storedProposal.GetContent())
	require.NoError(t, err)

	// then
	require.NoError(t, err)
	cInfo := wasmKeeper.GetContractInfo(ctx, contractAddr)
	require.NotNil(t, cInfo)
	assert.Equal(t, uint64(2), cInfo.CodeID)
	assert.Equal(t, anyAddress.String(), cInfo.Admin)
	assert.Equal(t, "testing", cInfo.Label)
	expHistory := []types.ContractCodeHistoryEntry{{
		Operation: types.ContractCodeHistoryOperationTypeGenesis,
		CodeID:    firstCodeID,
		Updated:   types.NewAbsoluteTxPosition(ctx),
	}, {
		Operation: types.ContractCodeHistoryOperationTypeMigrate,
		CodeID:    src.CodeID,
		Updated:   types.NewAbsoluteTxPosition(ctx),
		Msg:       src.Msg,
	}}
	assert.Equal(t, expHistory, wasmKeeper.GetContractHistory(ctx, contractAddr))
	// and events emitted
	require.Len(t, em.Events(), 2)
	assert.Equal(t, types.EventTypeMigrate, em.Events()[0].Type)
	require.Equal(t, types.EventTypeGovContractResult, em.Events()[1].Type)
	require.Len(t, em.Events()[1].Attributes, 1)
	assert.Equal(t, types.AttributeKeyResultDataHex, string(em.Events()[1].Attributes[0].Key))
}

func TestExecuteProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, bankKeeper := keepers.GovKeeper, keepers.BankKeeper

	exampleContract := InstantiateHackatomExampleContract(t, ctx, keepers)
	contractAddr := exampleContract.Contract

	// check balance
	bal := bankKeeper.GetBalance(ctx, contractAddr, "denom")
	require.Equal(t, bal.Amount, sdk.NewInt(100))

	releaseMsg := struct {
		Release struct{} `json:"release"`
	}{}
	releaseMsgBz, err := json.Marshal(releaseMsg)
	require.NoError(t, err)

	// try with runAs that doesn't have pemission
	badSrc := types.ExecuteContractProposal{
		Title:       "First",
		Description: "Beneficiary has no permission to run",
		Contract:    contractAddr.String(),
		Msg:         releaseMsgBz,
		RunAs:       exampleContract.BeneficiaryAddr.String(),
	}

	em := sdk.NewEventManager()

	// fails on store - this doesn't have permission
	storedProposal, err := govKeeper.SubmitProposal(ctx, &badSrc)
	require.Error(t, err)
	// balance should not change
	bal = bankKeeper.GetBalance(ctx, contractAddr, "denom")
	require.Equal(t, bal.Amount, sdk.NewInt(100))

	// try again with the proper run-as
	src := types.ExecuteContractProposal{
		Title:       "Second",
		Description: "Verifier can execute",
		Contract:    contractAddr.String(),
		Msg:         releaseMsgBz,
		RunAs:       exampleContract.VerifierAddr.String(),
	}

	em = sdk.NewEventManager()

	// when stored
	storedProposal, err = govKeeper.SubmitProposal(ctx, &src)
	require.NoError(t, err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx.WithEventManager(em), storedProposal.GetContent())
	require.NoError(t, err)

	// balance should be empty (proper release)
	bal = bankKeeper.GetBalance(ctx, contractAddr, "denom")
	require.Equal(t, bal.Amount, sdk.NewInt(0))
}

func TestSudoProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, bankKeeper := keepers.GovKeeper, keepers.BankKeeper

	exampleContract := InstantiateHackatomExampleContract(t, ctx, keepers)
	contractAddr := exampleContract.Contract
	_, _, anyAddr := keyPubAddr()

	// check balance
	bal := bankKeeper.GetBalance(ctx, contractAddr, "denom")
	require.Equal(t, bal.Amount, sdk.NewInt(100))
	bal = bankKeeper.GetBalance(ctx, anyAddr, "denom")
	require.Equal(t, bal.Amount, sdk.NewInt(0))

	type StealMsg struct {
		Recipient string     `json:"recipient"`
		Amount    []sdk.Coin `json:"amount"`
	}
	stealMsg := struct {
		Steal StealMsg `json:"steal_funds"`
	}{Steal: StealMsg{
		Recipient: anyAddr.String(),
		Amount:    []sdk.Coin{sdk.NewInt64Coin("denom", 75)},
	}}
	stealMsgBz, err := json.Marshal(stealMsg)
	require.NoError(t, err)

	// sudo can do anything
	src := types.SudoContractProposal{
		Title:       "Sudo",
		Description: "Steal funds for the verifier",
		Contract:    contractAddr.String(),
		Msg:         stealMsgBz,
	}

	em := sdk.NewEventManager()

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(ctx, &src)
	require.NoError(t, err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx.WithEventManager(em), storedProposal.GetContent())
	require.NoError(t, err)

	// balance should be empty (and verifier richer)
	bal = bankKeeper.GetBalance(ctx, contractAddr, "denom")
	require.Equal(t, bal.Amount, sdk.NewInt(25))
	bal = bankKeeper.GetBalance(ctx, anyAddr, "denom")
	require.Equal(t, bal.Amount, sdk.NewInt(75))
}

func TestAdminProposals(t *testing.T) {
	var (
		otherAddress sdk.AccAddress = bytes.Repeat([]byte{0x2}, types.ContractAddrLen)
		contractAddr                = BuildContractAddress(1, 1)
	)
	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	specs := map[string]struct {
		state       types.ContractInfo
		srcProposal govtypes.Content
		expAdmin    sdk.AccAddress
	}{
		"update with different admin": {
			state: types.ContractInfoFixture(),
			srcProposal: &types.UpdateAdminProposal{
				Title:       "Foo",
				Description: "Bar",
				Contract:    contractAddr.String(),
				NewAdmin:    otherAddress.String(),
			},
			expAdmin: otherAddress,
		},
		"update with old admin empty": {
			state: types.ContractInfoFixture(func(info *types.ContractInfo) {
				info.Admin = ""
			}),
			srcProposal: &types.UpdateAdminProposal{
				Title:       "Foo",
				Description: "Bar",
				Contract:    contractAddr.String(),
				NewAdmin:    otherAddress.String(),
			},
			expAdmin: otherAddress,
		},
		"clear admin": {
			state: types.ContractInfoFixture(),
			srcProposal: &types.ClearAdminProposal{
				Title:       "Foo",
				Description: "Bar",
				Contract:    contractAddr.String(),
			},
			expAdmin: nil,
		},
		"clear with old admin empty": {
			state: types.ContractInfoFixture(func(info *types.ContractInfo) {
				info.Admin = ""
			}),
			srcProposal: &types.ClearAdminProposal{
				Title:       "Foo",
				Description: "Bar",
				Contract:    contractAddr.String(),
			},
			expAdmin: nil,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			ctx, keepers := CreateTestInput(t, false, "staking")
			govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper
			wasmKeeper.SetParams(ctx, types.Params{
				CodeUploadAccess:             types.AllowNobody,
				InstantiateDefaultPermission: types.AccessTypeNobody,
				MaxWasmCodeSize:              types.DefaultMaxWasmCodeSize,
			})

			codeInfoFixture := types.CodeInfoFixture(types.WithSHA256CodeHash(wasmCode))
			require.NoError(t, wasmKeeper.importCode(ctx, 1, codeInfoFixture, wasmCode))

			require.NoError(t, wasmKeeper.importContract(ctx, contractAddr, &spec.state, []types.Model{}))
			// when stored
			storedProposal, err := govKeeper.SubmitProposal(ctx, spec.srcProposal)
			require.NoError(t, err)

			// and execute proposal
			handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
			err = handler(ctx, storedProposal.GetContent())
			require.NoError(t, err)

			// then
			cInfo := wasmKeeper.GetContractInfo(ctx, contractAddr)
			require.NotNil(t, cInfo)
			assert.Equal(t, spec.expAdmin.String(), cInfo.Admin)
		})
	}
}

func TestUpdateParamsProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper

	var (
		cdc                                   = keepers.WasmKeeper.cdc
		myAddress              sdk.AccAddress = make([]byte, types.ContractAddrLen)
		oneAddressAccessConfig                = types.AccessTypeOnlyAddress.With(myAddress)
	)

	nobodyJson, err := json.Marshal(types.AccessTypeNobody)
	require.NoError(t, err)
	specs := map[string]struct {
		src                proposal.ParamChange
		expUploadConfig    types.AccessConfig
		expInstantiateType types.AccessType
	}{
		"update upload permission param": {
			src: proposal.ParamChange{
				Subspace: types.ModuleName,
				Key:      string(types.ParamStoreKeyUploadAccess),
				Value:    string(cdc.MustMarshalJSON(&types.AllowNobody)),
			},
			expUploadConfig:    types.AllowNobody,
			expInstantiateType: types.AccessTypeEverybody,
		},
		"update upload permission param with address": {
			src: proposal.ParamChange{
				Subspace: types.ModuleName,
				Key:      string(types.ParamStoreKeyUploadAccess),
				Value:    string(cdc.MustMarshalJSON(&oneAddressAccessConfig)),
			},
			expUploadConfig:    oneAddressAccessConfig,
			expInstantiateType: types.AccessTypeEverybody,
		},
		"update instantiate param": {
			src: proposal.ParamChange{
				Subspace: types.ModuleName,
				Key:      string(types.ParamStoreKeyInstantiateAccess),
				Value:    string(nobodyJson),
			},
			expUploadConfig:    types.AllowEverybody,
			expInstantiateType: types.AccessTypeNobody,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			wasmKeeper.SetParams(ctx, types.DefaultParams())

			proposal := proposal.ParameterChangeProposal{
				Title:       "Foo",
				Description: "Bar",
				Changes:     []proposal.ParamChange{spec.src},
			}

			// when stored
			storedProposal, err := govKeeper.SubmitProposal(ctx, &proposal)
			require.NoError(t, err)

			// and proposal execute
			handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
			err = handler(ctx, storedProposal.GetContent())
			require.NoError(t, err)

			// then
			assert.True(t, spec.expUploadConfig.Equals(wasmKeeper.getUploadAccessConfig(ctx)),
				"got %#v not %#v", wasmKeeper.getUploadAccessConfig(ctx), spec.expUploadConfig)
			assert.Equal(t, spec.expInstantiateType, wasmKeeper.getInstantiateAccessConfig(ctx))
		})
	}
}

func TestPinCodesProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper

	mock := wasmtesting.MockWasmer{
		CreateFn:      wasmtesting.NoOpCreateFn,
		AnalyzeCodeFn: wasmtesting.WithoutIBCAnalyzeFn,
	}
	var (
		hackatom           = StoreHackatomExampleContract(t, ctx, keepers)
		hackatomDuplicate  = StoreHackatomExampleContract(t, ctx, keepers)
		otherContract      = StoreRandomContract(t, ctx, keepers, &mock)
		gotPinnedChecksums []wasmvm.Checksum
	)
	checksumCollector := func(checksum wasmvm.Checksum) error {
		gotPinnedChecksums = append(gotPinnedChecksums, checksum)
		return nil
	}
	specs := map[string]struct {
		srcCodeIDs []uint64
		mockFn     func(checksum wasmvm.Checksum) error
		expPinned  []wasmvm.Checksum
		expErr     bool
	}{
		"pin one": {
			srcCodeIDs: []uint64{hackatom.CodeID},
			mockFn:     checksumCollector,
		},
		"pin multiple": {
			srcCodeIDs: []uint64{hackatom.CodeID, otherContract.CodeID},
			mockFn:     checksumCollector,
		},
		"pin same code id": {
			srcCodeIDs: []uint64{hackatom.CodeID, hackatomDuplicate.CodeID},
			mockFn:     checksumCollector,
		},
		"pin non existing code id": {
			srcCodeIDs: []uint64{999},
			mockFn:     checksumCollector,
			expErr:     true,
		},
		"pin empty code id list": {
			srcCodeIDs: []uint64{},
			mockFn:     checksumCollector,
			expErr:     true,
		},
		"wasmvm failed with error": {
			srcCodeIDs: []uint64{hackatom.CodeID},
			mockFn: func(_ wasmvm.Checksum) error {
				return errors.New("test, ignore")
			},
			expErr: true,
		},
	}
	parentCtx := ctx
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			gotPinnedChecksums = nil
			ctx, _ := parentCtx.CacheContext()
			mock.PinFn = spec.mockFn
			proposal := types.PinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     spec.srcCodeIDs,
			}

			// when stored
			storedProposal, gotErr := govKeeper.SubmitProposal(ctx, &proposal)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)

			// and proposal execute
			handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
			gotErr = handler(ctx, storedProposal.GetContent())
			require.NoError(t, gotErr)

			// then
			for i := range spec.srcCodeIDs {
				c := wasmKeeper.GetCodeInfo(ctx, spec.srcCodeIDs[i])
				require.Equal(t, wasmvm.Checksum(c.CodeHash), gotPinnedChecksums[i])
			}
		})
	}
}
func TestUnpinCodesProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking")
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper

	mock := wasmtesting.MockWasmer{
		CreateFn:      wasmtesting.NoOpCreateFn,
		AnalyzeCodeFn: wasmtesting.WithoutIBCAnalyzeFn,
	}
	var (
		hackatom             = StoreHackatomExampleContract(t, ctx, keepers)
		hackatomDuplicate    = StoreHackatomExampleContract(t, ctx, keepers)
		otherContract        = StoreRandomContract(t, ctx, keepers, &mock)
		gotUnpinnedChecksums []wasmvm.Checksum
	)
	checksumCollector := func(checksum wasmvm.Checksum) error {
		gotUnpinnedChecksums = append(gotUnpinnedChecksums, checksum)
		return nil
	}
	specs := map[string]struct {
		srcCodeIDs  []uint64
		mockFn      func(checksum wasmvm.Checksum) error
		expUnpinned []wasmvm.Checksum
		expErr      bool
	}{
		"unpin one": {
			srcCodeIDs: []uint64{hackatom.CodeID},
			mockFn:     checksumCollector,
		},
		"unpin multiple": {
			srcCodeIDs: []uint64{hackatom.CodeID, otherContract.CodeID},
			mockFn:     checksumCollector,
		},
		"unpin same code id": {
			srcCodeIDs: []uint64{hackatom.CodeID, hackatomDuplicate.CodeID},
			mockFn:     checksumCollector,
		},
		"unpin non existing code id": {
			srcCodeIDs: []uint64{999},
			mockFn:     checksumCollector,
			expErr:     true,
		},
		"unpin empty code id list": {
			srcCodeIDs: []uint64{},
			mockFn:     checksumCollector,
			expErr:     true,
		},
		"wasmvm failed with error": {
			srcCodeIDs: []uint64{hackatom.CodeID},
			mockFn: func(_ wasmvm.Checksum) error {
				return errors.New("test, ignore")
			},
			expErr: true,
		},
	}
	parentCtx := ctx
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			gotUnpinnedChecksums = nil
			ctx, _ := parentCtx.CacheContext()
			mock.UnpinFn = spec.mockFn
			proposal := types.UnpinCodesProposal{
				Title:       "Foo",
				Description: "Bar",
				CodeIDs:     spec.srcCodeIDs,
			}

			// when stored
			storedProposal, gotErr := govKeeper.SubmitProposal(ctx, &proposal)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)

			// and proposal execute
			handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
			gotErr = handler(ctx, storedProposal.GetContent())
			require.NoError(t, gotErr)

			// then
			for i := range spec.srcCodeIDs {
				c := wasmKeeper.GetCodeInfo(ctx, spec.srcCodeIDs[i])
				require.Equal(t, wasmvm.Checksum(c.CodeHash), gotUnpinnedChecksums[i])
			}
		})
	}
}

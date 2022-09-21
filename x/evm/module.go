package evm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"math/big"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"

	"github.com/okex/exchain/libs/temp"
	"github.com/okex/exchain/x/evm/client/cli"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
)

var _ module.AppModuleBasic = AppModuleBasic{}
var _ module.AppModule = AppModule{}

// AppModuleBasic struct
type AppModuleBasic struct{}

// Name for app module basic
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers types for module
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis is json default structure
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis is the validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var genesisState types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &genesisState)
	if err != nil {
		return err
	}

	return genesisState.Validate()
}

// RegisterRESTRoutes Registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
}

// GetQueryCmd Gets the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(types.ModuleName, cdc)
}

// GetTxCmd Gets the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

//____________________________________________________________________________

// AppModule implements an application module for the evm module.
type AppModule struct {
	AppModuleBasic
	keeper *Keeper
	ak     types.AccountKeeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k *Keeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		ak:             ak,
	}
}

// Name is module name
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants interface for registering invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	keeper.RegisterInvariants(ir, *am.keeper)

	// for temp test
	temp.RegisterEvm("EncodeResultData", EncodeResultData)
	temp.RegisterEvmParser(EvmMsgParser)
}

func EvmMsgParser(msg sdk.Msg) (string, string, []byte, error) {
	evmTx, ok := msg.(*types.MsgEthereumTx)
	if ok {
		if evmTx.Data.Recipient == nil {
			return "", "", nil, fmt.Errorf("deploy contract should not conver cosmos msg")
		}
		fmt.Println("EvmMsgParser", string(evmTx.Data.Payload))
		return ContractStringParamParse(evmTx.Data.Payload)
	}
	return "", "", nil, fmt.Errorf("not a MsgEthereumTx msg")
}

type CMTxParam struct {
	Module   string `json:"module"`
	Function string `json:"function"`
	Data     string `json:"data"`
}

func ContractStringParamParse(input []byte) (string, string, []byte, error) {
	const methodSite = 4
	const fixedSite = 32
	const padSite = methodSite + fixedSite                 // 36
	const dataLenSite = methodSite + fixedSite + fixedSite // 68
	if len(input) < dataLenSite {
		return "", "", nil, fmt.Errorf("the input data size is error")
	}

	size := new(big.Int).SetBytes(input[padSite:dataLenSite]) // 存放数据长度
	if len(input) < int(size.Int64())+dataLenSite {
		return "", "", nil, fmt.Errorf("the input data size is error")
	}
	data := input[dataLenSite : size.Int64()+dataLenSite] // 实际数据

	value, err := hex.DecodeString(string(data)) // this is json fmt
	if err != nil {
		return "", "", nil, err
	}
	cmtx := &CMTxParam{}
	err = json.Unmarshal(value, cmtx)
	if err != nil {
		return "", "", nil, err
	}
	return cmtx.Module, cmtx.Function, []byte(cmtx.Data), nil
	//value, err = hex.DecodeString(cmtx.Data) // this is json fmt
	//if err != nil {
	//	return "", "", nil, err
	//}
	//return cmtx.Module, cmtx.Function, value, nil
}

func EncodeResultData(data []byte) ([]byte, error) {
	ethHash := common.BytesToHash(data)
	return types.EncodeResultData(&types.ResultData{TxHash: ethHash})
}

// Route specifies path for transactions
func (am AppModule) Route() string {
	return types.RouterKey
}

// NewHandler sets up a new handler for module
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute sets up path for queries
func (am AppModule) QuerierRoute() string {
	return types.ModuleName
}

// NewQuerierHandler sets up new querier handler for module
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(*am.keeper)
}

// BeginBlock function for module at start of each block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	am.keeper.BeginBlock(ctx, req)
}

// EndBlock function for module at end of block
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.keeper.EndBlock(ctx, req)
}

// InitGenesis instantiates the genesis state
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, *am.keeper, am.ak, genesisState)
}

// ExportGenesis exports the genesis state to be used by daemon
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, *am.keeper, am.ak)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

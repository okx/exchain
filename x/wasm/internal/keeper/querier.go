package keeper

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/okex/okchain/x/wasm/internal/types/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/wasm/internal/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryListContractByCode = "list-contracts-by-code"
	QueryGetContract        = "contract-info"
	QueryGetContractState   = "contract-state"
	QueryGetCode            = "code"
	QueryListCode           = "list-code"
)

const (
	QueryMethodContractStateSmart = "smart"
	QueryMethodContractStateAll   = "all"
	QueryMethodContractStateRaw   = "raw"
)

// ContractInfoWithAddress adds the address (key) to the ContractInfo representation
type ContractInfoWithAddress struct {
	// embedded here, so all json items remain top level
	*types.ContractInfo
	Address sdk.AccAddress `json:"address"`
}

// controls error output on querier - set true when testing/debugging
const debug = false

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryGetContract:
			return queryContractInfo(ctx, path[1], req, keeper)
		case QueryListContractByCode:
			return queryContractListByCode(ctx, path[1], req, keeper)
		case QueryGetContractState:
			if len(path) < 3 {
				return nil, sdk.ErrUnknownRequest("unknown data query endpoint")
			}
			return queryContractState(ctx, path[1], path[2], req, keeper)
		case QueryGetCode:
			return queryCode(ctx, path[1], req, keeper)
		case QueryListCode:
			return queryCodeList(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown data query endpoint")
		}
	}
}

func queryContractInfo(ctx sdk.Context, bech string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, sdk.ErrInvalidAddress(err.Error())
	}
	info := keeper.GetContractInfo(ctx, addr)
	if info == nil {
		return []byte("null"), nil
	}
	// redact the Created field (just used for sorting, not part of public API)
	info.Created = nil

	infoWithAddress := ContractInfoWithAddress{
		Address:      addr,
		ContractInfo: info,
	}
	bz, err := json.MarshalIndent(infoWithAddress, "", "  ")
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}

func queryContractListByCode(ctx sdk.Context, codeIDstr string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	codeID, err := strconv.ParseUint(codeIDstr, 10, 64)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	var contracts []ContractInfoWithAddress
	keeper.ListContractInfo(ctx, func(addr sdk.AccAddress, info types.ContractInfo) bool {
		if info.CodeID == codeID {
			// remove init message on list
			info.InitMsg = nil
			// and add the address
			infoWithAddress := ContractInfoWithAddress{
				Address:      addr,
				ContractInfo: &info,
			}
			contracts = append(contracts, infoWithAddress)
		}
		return false
	})

	// now we sort them by CreatedAt
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].ContractInfo.Created.LessThan(contracts[j].ContractInfo.Created)
	})
	// and remove that info for the final json (yes, the json:"-" tag doesn't work)
	for i := range contracts {
		contracts[i].Created = nil
	}

	bz, err := json.MarshalIndent(contracts, "", "  ")
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}

func queryContractState(ctx sdk.Context, bech, queryMethod string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	contractAddr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, sdk.ErrInvalidAddress(bech)
	}

	var resultData []types.Model
	switch queryMethod {
	case QueryMethodContractStateAll:
		// this returns a serialized json object (which internally encoded binary fields properly)
		for iter := keeper.GetContractState(ctx, contractAddr); iter.Valid(); iter.Next() {
			resultData = append(resultData, types.Model{
				Key:   iter.Key(),
				Value: iter.Value(),
			})
		}
		if resultData == nil {
			resultData = make([]types.Model, 0)
		}
	case QueryMethodContractStateRaw:
		// this returns a serialized json object
		resultData = keeper.QueryRaw(ctx, contractAddr, req.Data)
	case QueryMethodContractStateSmart:
		// this returns raw bytes (must be base64-encoded)
		return keeper.QuerySmart(ctx, contractAddr, req.Data)
	default:
		return nil, sdk.ErrUnknownRequest(queryMethod)
	}
	bz, err := json.Marshal(resultData)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}

type GetCodeResponse struct {
	ListCodeResponse
	// Data is the entire wasm bytecode
	Data []byte `json:"data" yaml:"data"`
}

func queryCode(ctx sdk.Context, codeIDstr string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	codeID, err := strconv.ParseUint(codeIDstr, 10, 64)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("invalid codeID: " + err.Error())
	}

	res := keeper.GetCodeInfo(ctx, codeID)
	if res == nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}
	info := ListCodeResponse{
		ID:       codeID,
		Creator:  res.Creator,
		DataHash: res.CodeHash,
		Source:   res.Source,
		Builder:  res.Builder,
	}

	code, err := keeper.GetByteCode(ctx, codeID)
	if err != nil {
		return nil, sdk.ErrInternal("loading wasm code")
	}

	bz, err := json.MarshalIndent(GetCodeResponse{info, code}, "", "  ")
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}

type ListCodeResponse struct {
	ID       uint64         `json:"id"`
	Creator  sdk.AccAddress `json:"creator"`
	DataHash bytes.HexBytes `json:"data_hash"`
	Source   string         `json:"source"`
	Builder  string         `json:"builder"`
}

func queryCodeList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var info []ListCodeResponse

	var i uint64
	for true {
		i++
		res := keeper.GetCodeInfo(ctx, i)
		if res == nil {
			break
		}
		info = append(info, ListCodeResponse{
			ID:       i,
			Creator:  res.Creator,
			DataHash: res.CodeHash,
			Source:   res.Source,
			Builder:  res.Builder,
		})
	}

	bz, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}

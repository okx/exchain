package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/app/utils"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		if len(path) < 1 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
				"Insufficient parameters, at least 1 parameter is required")
		}

		switch path[0] {
		case types.QueryBalance:
			return queryBalance(ctx, path, keeper)
		case types.QueryBlockNumber:
			return queryBlockNumber(ctx, keeper)
		case types.QueryStorage:
			return queryStorage(ctx, path, keeper)
		case types.QueryStorageByKey:
			return queryStorageByKey(ctx, path, keeper)
		case types.QueryCode:
			return queryCode(ctx, path, keeper)
		case types.QueryCodeByHash:
			return queryCodeByHash(ctx, path, keeper)
		case types.QueryHashToHeight:
			return queryHashToHeight(ctx, path, keeper)
		case types.QueryBloom:
			return queryBlockBloom(ctx, path, keeper)
		case types.QueryAccount:
			return queryAccount(ctx, path, keeper)
		case types.QueryExportAccount:
			return queryExportAccount(ctx, path, keeper)
		case types.QueryParameters:
			return queryParams(ctx, keeper)
		case types.QueryHeightToHash:
			return queryHeightToHash(ctx, path, keeper)
		case types.QuerySection:
			return querySection(ctx, path, keeper)
		case types.QueryContractDeploymentWhitelist:
			return queryContractDeploymentWhitelist(ctx, keeper)
		case types.QueryContractBlockedList:
			return queryContractBlockedList(ctx, keeper)
		case types.QueryContractMethodBlockedList:
			return queryContractMethodBlockedList(ctx, keeper)
		case types.QueryAllMapping:
			return queryAllMapping(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}

func queryContractMethodBlockedList(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	blockedList := types.CreateEmptyCommitStateDB(keeper.GeneratePureCSDBParams(), ctx).GetContractMethodBlockedList()
	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, blockedList)
	if errUnmarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", errUnmarshal.Error()))
	}

	return res, nil
}

func queryContractBlockedList(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	blockedList := types.CreateEmptyCommitStateDB(keeper.GeneratePureCSDBParams(), ctx).GetContractBlockedList()
	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, blockedList)
	if errUnmarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", errUnmarshal.Error()))
	}

	return res, nil
}

func queryContractDeploymentWhitelist(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	whitelist := types.CreateEmptyCommitStateDB(keeper.GeneratePureCSDBParams(), ctx).GetContractDeploymentWhitelist()
	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, whitelist)
	if errUnmarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", errUnmarshal.Error()))
	}

	return res, nil
}

func queryBalance(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	addr := ethcmn.HexToAddress(path[1])
	balance := keeper.GetBalance(ctx, addr)
	balanceStr, err := utils.MarshalBigInt(balance)
	if err != nil {
		return nil, err
	}

	res := types.QueryResBalance{Balance: balanceStr}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBlockNumber(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	num := ctx.BlockHeight()
	bnRes := types.QueryResBlockNumber{Number: num}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, bnRes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryStorage(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 3 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 3 parameters is required")
	}

	addr := ethcmn.HexToAddress(path[1])
	key := ethcmn.HexToHash(path[2])
	val := keeper.GetState(ctx, addr, key)
	res := types.QueryResStorage{Value: val.Bytes()}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryStorageByKey(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 3 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 3 parameters is required")
	}

	addr := ethcmn.HexToAddress(path[1])
	key := ethcmn.HexToHash(path[2])
	val := keeper.GetStateByKey(ctx, addr, key)
	res := types.QueryResStorage{Value: val.Bytes()}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryCode(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	addr := ethcmn.HexToAddress(path[1])
	so := keeper.GetOrNewStateObject(ctx, addr)
	code := keeper.GetCodeByHash(ctx, ethcmn.BytesToHash(so.CodeHash()))
	res := types.QueryResCode{Code: code}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryCodeByHash(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	hash := ethcmn.HexToHash(path[1])
	code := keeper.GetCodeByHash(ctx, hash)
	res := types.QueryResCode{Code: code}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryHashToHeight(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	blockHash := ethcmn.FromHex(path[1])
	blockNumber, found := keeper.GetBlockHash(ctx, blockHash)
	if !found {
		return []byte{}, sdkerrors.Wrap(types.ErrKeyNotFound, fmt.Sprintf("block height not found for hash %s", path[1]))
	}

	res := types.QueryResBlockNumber{Number: blockNumber}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBlockBloom(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	num, err := strconv.ParseInt(path[1], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrStrConvertFailed, fmt.Sprintf("could not unmarshal block height: %s", err))
	}

	bloom := keeper.GetBlockBloom(ctx.WithBlockHeight(num), num)
	res := types.QueryBloomFilter{Bloom: bloom}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAccount(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	addr := ethcmn.HexToAddress(path[1])
	so := keeper.GetOrNewStateObject(ctx, addr)

	balance, err := utils.MarshalBigInt(so.Balance())
	if err != nil {
		return nil, err
	}

	res := types.QueryResAccount{
		Balance:  balance,
		CodeHash: so.CodeHash(),
		Nonce:    so.Nonce(),
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryExportAccount(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	hexAddress := path[1]
	addr := ethcmn.HexToAddress(hexAddress)

	var storage types.Storage
	err := keeper.ForEachStorage(ctx, addr, func(key, value ethcmn.Hash) bool {
		storage = append(storage, types.NewState(key, value))
		return false
	})
	if err != nil {
		return nil, err
	}

	res := types.GenesisAccount{
		Address: hexAddress,
		Code:    keeper.GetCode(ctx, addr),
		Storage: storage,
	}

	// TODO: codec.MarshalJSONIndent doesn't call the String() method of types properly
	bz, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	params := keeper.GetParams(ctx)
	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if errUnmarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", errUnmarshal.Error()))
	}
	return res, nil
}

func queryHeightToHash(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	height, err := strconv.Atoi(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, params[1] convert to int failed")
	}
	hash := keeper.GetHeightHash(ctx, uint64(height))

	return hash.Bytes(), nil
}

func querySection(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if !types.GetEnableBloomFilter() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"disable bloom filter")
	}

	if len(path) != 1 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"wrong parameters, need no parameters")
	}

	res, err := json.Marshal(types.GetIndexer().StoredSection())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func queryAllMapping(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	mapping := make(map[string]string)
	keeper.IterateMapping(ctx, func(denom, contract string) bool {
		mapping[denom] = contract
		return false
	})

	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, types.QueryResAllMapping{mapping})
	if errUnmarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", errUnmarshal.Error()))
	}
	return res, nil
}

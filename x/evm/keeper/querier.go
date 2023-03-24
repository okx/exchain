package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"

	ethcmn "github.com/ethereum/go-ethereum/common"
	apptypes "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/app/utils"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/evm/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
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
		case types.QueryStorageProof:
			return queryStorageProof(ctx, path, keeper)
		case types.QueryStorageRoot:
			return queryStorageRootHash(ctx, path, keeper, req.Height)
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
		case types.QuerySysContractAddress:
			return querySysContractAddress(ctx, keeper)
		case types.QueryEthBlockByHeight:
			return queryEthBlockByHeight(ctx, path, keeper)
		case types.QueryEthBlockByHash:
			return queryEthBlockByHash(ctx, path, keeper)
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

func queryStorageProof(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 3 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 3 parameters is required")
	}

	addr := ethcmn.HexToAddress(path[1])
	key := ethcmn.HexToHash(path[2])
	val := keeper.GetState(ctx, addr, key)
	proofList, err := keeper.GetStorageProof(ctx, addr, key)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}
	res := types.QueryResStorageProof{Value: val.Bytes(), Proof: proofList}
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

func queryStorageRootHash(ctx sdk.Context, path []string, keeper Keeper, height int64) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 1 parameters is required")
	}

	addr := ethcmn.HexToAddress(path[1])
	acc := keeper.accountKeeper.GetAccount(ctx, addr.Bytes())
	if acc == nil {
		return nil, fmt.Errorf("get %s storage root hash failed: acc is not exist", addr)
	}

	return acc.GetStateRoot().Bytes(), nil
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

	blockHash := ethcmn.HexToHash(path[1])
	blockNumber, found := keeper.GetBlockHeight(ctx, blockHash)
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

	copyCtx := ctx
	copyCtx.SetBlockHeight(num)
	bloom := keeper.GetBlockBloom(copyCtx, num)
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
	res, err := resolveEthAccount(ctx, keeper, addr)
	if err != nil {
		return nil, err
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func resolveEthAccount(ctx sdk.Context, k Keeper, addr ethcmn.Address) (*types.QueryResAccount, error) {
	codeHash := mpt.EmptyCodeHashBytes
	account := k.accountKeeper.GetAccount(ctx, addr.Bytes())
	if account == nil {
		return &types.QueryResAccount{Nonce: uint64(0), CodeHash: codeHash, Balance: "0x0"}, nil
	}
	ethAccount := account.(*apptypes.EthAccount)
	if ethAccount == nil {
		return &types.QueryResAccount{Nonce: uint64(0), CodeHash: codeHash, Balance: "0x0"}, nil
	}

	// get balance
	balance := ethAccount.Balance(sdk.DefaultBondDenom).BigInt()
	if balance == nil {
		balance = sdk.ZeroInt().BigInt()
	}
	balanceStr, err := utils.MarshalBigInt(balance)
	if err != nil {
		return nil, err
	}

	// get codeHash
	if ethAccount.CodeHash != nil {
		codeHash = ethAccount.CodeHash
	}

	//return
	return &types.QueryResAccount{
		Balance:  balanceStr,
		CodeHash: codeHash,
		Nonce:    ethAccount.Sequence,
	}, nil
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

func queryEthBlockByHeight(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	height, err := strconv.Atoi(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, params[1] convert to int failed")
	}

	res, found := keeper.GetEthBlockBytesByHeight(ctx, uint64(height))
	if !found {
		return nil, fmt.Errorf("not found block by heith(%d)", height)
	}

	return res, nil
}

func queryEthBlockByHash(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	blockHash := ethcmn.HexToHash(path[1])
	res, found := keeper.GetEthBlockBytesByHash(ctx, blockHash.Bytes())
	if !found {
		return nil, fmt.Errorf("not found block by hash(%s)", blockHash.Hex())
	}

	return res, nil
}

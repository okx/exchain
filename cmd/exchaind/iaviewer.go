package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/okex/exchain/libs/tendermint/crypto"

	"github.com/gogo/protobuf/proto"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	"github.com/okex/exchain/x/gov"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/tendermint/go-amino"

	"github.com/spf13/viper"

	"github.com/okex/exchain/app"
	minttypes "github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	supplytypes "github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	evmtypes "github.com/okex/exchain/x/evm/types"
	slashingtypes "github.com/okex/exchain/x/slashing"
	tokentypes "github.com/okex/exchain/x/token/types"
	"github.com/spf13/cobra"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	acctypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	distypes "github.com/okex/exchain/libs/cosmos-sdk/x/distribution/types"
	govtypes "github.com/okex/exchain/libs/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	"github.com/okex/exchain/x/distribution/types"
)

type (
	formatKeyValue func(cdc *codec.Codec, key []byte, value []byte) string
)

const (
	KeyDistribution = "s/k:distribution/"
	KeyGov          = "s/k:gov/"
	KeyMain         = "s/k:main/"
	KeyToken        = "s/k:token/"
	KeyMint         = "s/k:mint/"
	KeyAcc          = "s/k:acc/"
	KeySupply       = "s/k:supply/"
	KeyEvm          = "s/k:evm/"
	KeyParams       = "s/k:params/"
	KeyStaking      = "s/k:staking/"
	KeySlashing     = "s/k:slashing/"

	DefaultCacheSize int = 100000

	flagStart     = "start"
	flagLimit     = "limit"
	flagHex       = "hex"
	flagPrefix    = "prefix"
	flagKey       = "key"
	flagKeyPrefix = "keyprefix"
)

var printKeysDict = map[string]formatKeyValue{
	KeyEvm:          evmPrintKey,
	KeyAcc:          accPrintKey,
	KeyParams:       paramsPrintKey,
	KeyStaking:      stakingPrintKey,
	KeyGov:          govPrintKey,
	KeyDistribution: distributionPrintKey,
	KeySlashing:     slashingPrintKey,
	KeyMain:         mainPrintKey,
	KeyToken:        tokenPrintKey,
	KeyMint:         mintPrintKey,
	KeySupply:       supplyPrintKey,
}

type iaviewerFlags struct {
	Start     *int
	Limit     *int
	DbBackend *string
	Prefix    *string
}

type iaviewerContext struct {
	DataDir   string
	Prefix    string
	Module    string
	Version   int
	DbBackend dbm.BackendType
	Start     int
	Limit     int
	Codec     *codec.Codec

	flags iaviewerFlags
}

func iaviewerCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iaviewer",
		Short: "Read iavl tree data from db",
	}
	iavlCtx := &iaviewerContext{Codec: cdc, DbBackend: dbm.BackendType(ctx.Config.DBBackend)}

	cmd.AddCommand(
		iaviewerReadCmd(iavlCtx),
		iaviewerStatusCmd(iavlCtx),
		iaviewerDiffCmd(iavlCtx),
		iaviewerVersionsCmd(iavlCtx),
		iaviewerListModulesCmd(),
	)
	iavlCtx.flags.DbBackend = cmd.PersistentFlags().String(flagDBBackend, "", "Database backend: goleveldb | rocksdb")
	iavlCtx.flags.Start = cmd.PersistentFlags().Int(flagStart, 0, "index of result set start from")
	iavlCtx.flags.Limit = cmd.PersistentFlags().Int(flagLimit, 0, "limit of result set, 0 means no limit")
	iavlCtx.flags.Prefix = cmd.PersistentFlags().String(flagPrefix, "", "the prefix of iavl tree, module value must be \"\" if prefix is set")
	return cmd
}

func iaviewerCmdParseFlags(ctx *iaviewerContext) {
	if dbflag := ctx.flags.DbBackend; dbflag != nil && *dbflag != "" {
		ctx.DbBackend = dbm.BackendType(*dbflag)
	}

	if ctx.flags.Start != nil {
		ctx.Start = *ctx.flags.Start
	}
	if ctx.flags.Limit != nil {
		ctx.Limit = *ctx.flags.Limit
	}

	if ctx.flags.Prefix != nil && *ctx.flags.Prefix != "" {
		ctx.Prefix = *ctx.flags.Prefix
	}
}

func iaviewerCmdParseArgs(ctx *iaviewerContext, args []string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("must specify data_dir and module")
	}
	dataDir, module, version := args[0], args[1], 0
	if len(args) == 3 {
		version, err = strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid version: %s, error : %w\n", args[2], err)
		}
	}
	ctx.DataDir = dataDir
	ctx.Module = module
	ctx.Version = version
	if ctx.Module != "" {
		ctx.Prefix = fmt.Sprintf("s/k:%s/", ctx.Module)
	}
	return nil
}

func iaviewerListModulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-modules",
		Short: "List all module names",
		Run: func(cmd *cobra.Command, args []string) {
			moduleKeys := make([]string, 0, len(app.ModuleBasics))
			for moduleKey := range app.ModuleBasics {
				moduleKeys = append(moduleKeys, moduleKey)
			}
			sort.Strings(moduleKeys)
			fmt.Printf("there are %d modules:\n\n", len(moduleKeys))
			for _, key := range moduleKeys {
				fmt.Print("\t")
				fmt.Println(key)
			}
			fmt.Println()
		},
	}
	return cmd
}

func iaviewerReadCmd(ctx *iaviewerContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read <data_dir> <module> [version]",
		Short: "Read iavl tree key-value from db",
		Long:  "Read iavl tree key-value from db, you must specify data_dir and module, if version is 0 or not specified, read data from the latest version.\n",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			iaviewerCmdParseFlags(ctx)
			return iaviewerCmdParseArgs(ctx, args)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return iaviewerReadData(ctx)
		},
	}
	cmd.PersistentFlags().Bool(flagHex, false, "print key and value in hex format")
	cmd.PersistentFlags().String(flagKey, "", "print only the value for this key, key must be in hex format.\n"+
		"if specified, keyprefix, start and limit flags would be ignored")
	cmd.PersistentFlags().String(flagKeyPrefix, "", "print values for keys with specified prefix, prefix must be in hex format.")
	return cmd
}

func iaviewerStatusCmd(ctx *iaviewerContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <data_dir> <module> [version]",
		Short: "print iavl tree status",
		Long:  "print iavl tree status, you must specify data_dir and module, if version is 0 or not specified, read data from the latest version.\n",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			iaviewerCmdParseFlags(ctx)
			return iaviewerCmdParseArgs(ctx, args)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return iaviewerStatus(ctx)
		},
	}
	return cmd
}

func iaviewerVersionsCmd(ctx *iaviewerContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions <data_dir> <module> [version]",
		Short: "list iavl tree versions",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			iaviewerCmdParseFlags(ctx)
			return iaviewerCmdParseArgs(ctx, args)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return iaviewerVersions(ctx)
		},
	}
	return cmd
}

func iaviewerDiffCmd(ctx *iaviewerContext) *cobra.Command {
	var ver2 int
	cmd := &cobra.Command{
		Use:   "diff <data_dir> <module> <version1> <version2>",
		Short: "compare different key-value between two versions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			iaviewerCmdParseFlags(ctx)
			if len(args) != 4 {
				return fmt.Errorf("must specify data_dir, module, version1 and version2")
			}
			ctx.DataDir = args[0]
			ctx.Module = args[1]
			if ctx.Module != "" {
				ctx.Prefix = fmt.Sprintf("s/k:%s/", ctx.Module)
			}

			ver1, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid version1: %s, error : %w\n", args[2], err)
			}
			ctx.Version = ver1
			ver2, err = strconv.Atoi(args[3])
			if err != nil {
				return fmt.Errorf("invalid version2: %s, error : %w\n", args[3], err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return iaviewerPrintDiff(ctx, ver2)
		},
	}
	return cmd
}

// iaviewerPrintDiff reads different key-value from leveldb according two paths
func iaviewerPrintDiff(ctx *iaviewerContext, version2 int) error {
	db, err := OpenDB(ctx.DataDir, ctx.DbBackend)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}
	defer db.Close()

	tree, err := ReadTree(db, ctx.Version, []byte(ctx.Prefix), DefaultCacheSize)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	compareTree, err := ReadTree(db, version2, []byte(ctx.Prefix), DefaultCacheSize)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Printf("module: %s, prefix key: %s\n\n", ctx.Module, ctx.Prefix)

	if bytes.Equal(tree.Hash(), compareTree.Hash()) {
		fmt.Printf("tree version %d and %d are same, root hash: %X\n", ctx.Version, version2, tree.Hash())
		return nil
	}

	tree.Iterate(func(key, value []byte) bool {
		_, v2 := compareTree.Get(key)
		if v2 == nil {
			fmt.Printf("---only in ver1 %d, key: %X, value: %X\n", ctx.Version, key, value)
		} else {
			if !bytes.Equal(value, v2) {
				fmt.Printf("---diff ver1 %d, key: %X, value: %X\n", ctx.Version, key, value)
				fmt.Printf("+++diff ver2 %d, key: %X, value: %X\n", version2, key, v2)
			}
		}
		return false
	})

	compareTree.Iterate(func(key, value []byte) bool {
		_, v1 := tree.Get(key)
		if v1 == nil {
			fmt.Printf("+++only in ver2 %d, key: %X, value: %X\n", version2, key, value)
		}
		return false
	})
	return nil
}

// iaviewerReadData reads key-value from leveldb
func iaviewerReadData(ctx *iaviewerContext) error {
	db, err := OpenDB(ctx.DataDir, ctx.DbBackend)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}
	defer db.Close()

	tree, err := ReadTree(db, ctx.Version, []byte(ctx.Prefix), DefaultCacheSize)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Printf("module: %s, prefix key: %s\n\n", ctx.Module, ctx.Prefix)

	if key := viper.GetString(flagKey); key != "" {
		keyByte, err := hex.DecodeString(key)
		if err != nil {
			return fmt.Errorf("error decoding key: %w", err)
		}
		i, value := tree.Get(keyByte)
		fmt.Printf("key:\t%s\nvalue:\t%X\nindex:\t%d\n", key, value, i)
		return nil
	}

	printTree(ctx, tree)
	return nil
}

func iaviewerStatus(ctx *iaviewerContext) error {
	db, err := OpenDB(ctx.DataDir, ctx.DbBackend)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}
	defer db.Close()

	tree, err := ReadTree(db, ctx.Version, []byte(ctx.Prefix), DefaultCacheSize)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Printf("module: %s, prefix key: %s\n", ctx.Module, ctx.Prefix)
	printIaviewerStatus(tree)
	return nil
}

func printIaviewerStatus(tree *iavl.MutableTree) {
	fmt.Printf("iavl tree:\n"+
		"\troot hash: %X\n"+
		"\tsize: %d\n"+
		"\tcurrent version: %d\n", tree.Hash(), tree.Size(), tree.Version())
}

func iaviewerVersions(ctx *iaviewerContext) error {
	db, err := OpenDB(ctx.DataDir, ctx.DbBackend)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}
	defer db.Close()

	tree, err := ReadTree(db, ctx.Version, []byte(ctx.Prefix), DefaultCacheSize)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Printf("module: %s, prefix key: %s\n\n", ctx.Module, ctx.Prefix)
	iaviewerPrintVersions(ctx, tree)
	return nil
}

func iaviewerPrintVersions(ctx *iaviewerContext, tree *iavl.MutableTree) {
	versions := tree.AvailableVersions()
	fmt.Printf("total versions: %d\n", len(versions))

	if ctx.Start >= len(versions) {
		fmt.Printf("printed verions: 0\n")
		return
	}
	if ctx.Start+ctx.Limit > len(versions) {
		ctx.Limit = len(versions) - ctx.Start
	}
	if ctx.Limit == 0 {
		versions = versions[ctx.Start:]
	} else {
		versions = versions[ctx.Start : ctx.Start+ctx.Limit]
	}
	fmt.Printf("printed versions: %d\n\n", len(versions))

	for _, v := range versions {
		fmt.Printf("  %d\n", v)
	}
}

// getKVs, get all key-values by mutableTree
func getKVs(tree *iavl.MutableTree, dataMap map[string][32]byte, wg *sync.WaitGroup) {
	tree.Iterate(func(key []byte, value []byte) bool {
		dataMap[hex.EncodeToString(key)] = sha256.Sum256(value)
		return false
	})
	wg.Done()
}

func defaultKvFormatter(key []byte, value []byte) string {
	printKey := parseWeaveKey(key)
	return fmt.Sprintf("parsed key:\t%s\nhex key:\t%X\nhex value:\t%X", printKey, key, value)
}

func printKV(cdc *codec.Codec, modulePrefixKey string, key []byte, value []byte) {
	if impl, exit := printKeysDict[modulePrefixKey]; exit && !viper.GetBool(flagHex) {
		kvFormat := impl(cdc, key, value)
		if kvFormat != "" {
			fmt.Println(kvFormat)
			fmt.Println()
			return
		}
	}
	fmt.Println(defaultKvFormatter(key, value))
	fmt.Println()
}

func printTree(ctx *iaviewerContext, tree *iavl.MutableTree) {
	startKey := []byte(nil)
	endKey := []byte(nil)

	var keyPrefix string
	if keyPrefix = viper.GetString(flagKeyPrefix); keyPrefix != "" {
		index, _ := tree.Get(amino.StrToBytes(keyPrefix))
		ctx.Start += int(index)
	}

	if tree.Size() <= int64(ctx.Start) {
		return
	}
	printed := ctx.Limit
	if ctx.Start != 0 {
		startKey, _ = tree.GetByIndex(int64(ctx.Start))
	}
	if ctx.Limit != 0 && int64(ctx.Start+ctx.Limit) < tree.Size() {
		endKey, _ = tree.GetByIndex(int64(ctx.Start + ctx.Limit))
	} else {
		printed = int(tree.Size()) - ctx.Start
	}

	fmt.Printf("total: %d\n", tree.Size())
	fmt.Printf("printed: %d\n\n", printed)

	tree.IterateRange(startKey, endKey, true, func(key []byte, value []byte) bool {
		if keyPrefix != "" {
			if !bytes.HasPrefix(key, amino.StrToBytes(keyPrefix)) {
				return true
			}
		}
		printKV(ctx.Codec, ctx.Prefix, key, value)
		return false
	})

	//tree.Iterate(func(key []byte, value []byte) bool {
	//	printKV(ctx.Codec, ctx.Prefix, key, value)
	//	return false
	//})
}

func printByKey(cdc *codec.Codec, tree *iavl.MutableTree, module string, key []byte) {
	_, value := tree.Get(key)
	printKV(cdc, module, key, value)
}

func supplyPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case supplytypes.SupplyKey[0]:
		var supplyAmount sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &supplyAmount)
		return fmt.Sprintf("tokenSymbol:%s:info:%s", string(key[1:]), supplyAmount.String())
	default:
		return defaultKvFormatter(key, value)
	}
}

type MinterCustom struct {
	NextBlockToUpdate uint64       `json:"next_block_to_update" yaml:"next_block_to_update"` // record the block height for next year
	MintedPerBlock    sdk.DecCoins `json:"minted_per_block" yaml:"minted_per_block"`         // record the MintedPerBlock per block in this year
}

func mintPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case minttypes.MinterKey[0]:
		var minter MinterCustom
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &minter)
		return fmt.Sprintf("minter:%v", minter)
	default:
		return defaultKvFormatter(key, value)
	}
}

func tokenPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case tokentypes.TokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		return fmt.Sprintf("tokenName:%s:info:%s", string(key[1:]), token.String())
	case tokentypes.TokenNumberKey[0]:
		var tokenNumber uint64
		cdc.MustUnmarshalBinaryBare(value, &tokenNumber)
		return fmt.Sprintf("tokenNumber:%d", tokenNumber)
	case tokentypes.PrefixUserTokenKey[0]:
		return fmt.Sprintf("address:%s;symbol:%s", key[1:21], string(key[21:]))
	default:
		return defaultKvFormatter(key, value)
	}
}

func mainPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	if bytes.Equal(key, []byte("consensus_params")) {
		var cons abci.ConsensusParams
		err := proto.Unmarshal(value, &cons)
		if err != nil {
			return fmt.Sprintf("consensusParams:%X; unmarshal error, %s", value, err)
		}
		return fmt.Sprintf("consensusParams:%s", cons.String())
	}
	return defaultKvFormatter(key, value)
}

func slashingPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case slashingtypes.ValidatorSigningInfoKey[0]:
		var signingInfo slashingtypes.ValidatorSigningInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &signingInfo)
		return fmt.Sprintf("validatorAddr:%X;signingInfo:%s", key[1:], signingInfo.String())
	case slashingtypes.ValidatorMissedBlockBitArrayKey[0]:
		var index int64
		index = int64(binary.LittleEndian.Uint64(key[len(key)-8:]))
		var missed bool
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &missed)
		return fmt.Sprintf("validatorMissedBlockAddr:%X;index:%d;missed:%v", key[1:len(key)-8], index, missed)
	case slashingtypes.AddrPubkeyRelationKey[0]:
		var pubkey crypto.PubKey
		err := cdc.UnmarshalBinaryLengthPrefixed(value, &pubkey)
		if err != nil {
			return fmt.Sprintf("pubkeyAddr:%X;value %X unmarshal error, %s", key[1:], value, err)
		} else {
			return fmt.Sprintf("pubkeyAddr:%X;pubkey:%X", key[1:], pubkey.Bytes())
		}
	default:
		return defaultKvFormatter(key, value)
	}
}

func distributionPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case distypes.FeePoolKey[0]:
		var feePool distypes.FeePool
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &feePool)
		return fmt.Sprintf("feePool:%v", feePool)
	case distypes.ProposerKey[0]:
		var consAddr sdk.ConsAddress
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &consAddr)
		return fmt.Sprintf("proposerKey consAddress:%X", consAddr)
	case distypes.DelegatorWithdrawAddrPrefix[0]:
		return fmt.Sprintf("delegatorWithdrawAddr:%X;address:%X", key[1:], value)
	case distypes.ValidatorAccumulatedCommissionPrefix[0]:
		var commission types.ValidatorAccumulatedCommission
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &commission)
		return fmt.Sprintf("validatorAccumulatedAddr:%X;commission:%s", key[1:], commission.String())
	default:
		return defaultKvFormatter(key, value)
	}
}

func govPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case govtypes.ProposalsKeyPrefix[0]:
		var prop gov.Proposal
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &prop)
		return fmt.Sprintf("proposalId:%d;proposal:%s", binary.BigEndian.Uint64(key[1:]), prop.String())
	case govtypes.ActiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		return fmt.Sprintf("activeProposalEndTime:%s;proposalId:%d", time.String(), binary.BigEndian.Uint64(value))
	case govtypes.InactiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		return fmt.Sprintf("inactiveProposalEndTime:%s;proposalId:%d", time.String(), binary.BigEndian.Uint64(value))
	case govtypes.ProposalIDKey[0]:
		if len(value) != 0 {
			return fmt.Sprintf("proposalId:%d", binary.BigEndian.Uint64(value))
		} else {
			return fmt.Sprintf("proposalId:nil")
		}
	default:
		return defaultKvFormatter(key, value)
	}
}

func stakingPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case stakingtypes.LastValidatorPowerKey[0]:
		var power int64
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		return fmt.Sprintf("validatorAddress:%X;power:%d", key[1:], power)
	case stakingtypes.LastTotalPowerKey[0]:
		var power sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		return fmt.Sprintf("lastTotolValidatorPower:%s", power.String())
	case stakingtypes.ValidatorsKey[0]:
		var validator stakingtypes.Validator
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &validator)
		return fmt.Sprintf("validatorAddress:%X;validator:%s", key[1:], validator.String())
	case stakingtypes.ValidatorsByConsAddrKey[0]:
		return fmt.Sprintf("validatorConsAddr:%X;valAddress:%X", key[1:], value)
	case stakingtypes.ValidatorsByPowerIndexKey[0]:
		consensusPower := int64(binary.BigEndian.Uint64(key[1:9]))
		operAddr := key[9:]
		for i, b := range operAddr {
			operAddr[i] = ^b
		}
		return fmt.Sprintf("validatorPowerIndex consensusPower:%d;operAddr:%X;operatorAddress:%X", consensusPower, operAddr, value)
	default:
		return defaultKvFormatter(key, value)
	}
}

func paramsPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	return fmt.Sprintf("paramsKey:%s;value:%s", string(key), string(value))
}

func accPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	if key[0] == acctypes.AddressStoreKeyPrefix[0] {
		var acc exported.Account
		cdc.MustUnmarshalBinaryBare(value, &acc)
		return fmt.Sprintf("adress:%X;account:%s", key[1:], acc.String())
	} else if bytes.Equal(key, acctypes.GlobalAccountNumberKey) {
		var accNum uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &accNum)
		return fmt.Sprintf("%s:%d", string(key), accNum)
	} else {
		return defaultKvFormatter(key, value)
	}
}

func evmPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case evmtypes.KeyPrefixBlockHash[0]:
		blockHash := key[1:]
		height := int64(binary.BigEndian.Uint64(value))
		return fmt.Sprintf("blockHash:%X;height:%d", blockHash, height)
	case evmtypes.KeyPrefixBloom[0]:
		height := int64(binary.BigEndian.Uint64(key[1:]))
		bloom := ethtypes.BytesToBloom(value)
		return fmt.Sprintf("bloomHeight:%d;data:%X", height, bloom[:])
	case evmtypes.KeyPrefixCode[0]:
		return fmt.Sprintf("codeHash:%X;code:%X", key[1:], value)
	case evmtypes.KeyPrefixStorage[0]:
		return fmt.Sprintf("stroageAddr:%X;key:%X;data:%X", key[1:21], key[21:], value)
	case evmtypes.KeyPrefixChainConfig[0]:
		if len(value) != 0 {
			var config evmtypes.ChainConfig
			cdc.MustUnmarshalBinaryBare(value, &config)
			return fmt.Sprintf("chainConfig:%s", config.String())
		} else {
			return fmt.Sprintf("chainConfig:nil")
		}
	case evmtypes.KeyPrefixHeightHash[0]:
		height := binary.BigEndian.Uint64(key[1:])
		return fmt.Sprintf("height:%d;blockHash:%X", height, value)
	case evmtypes.KeyPrefixContractDeploymentWhitelist[0]:
		return fmt.Sprintf("contractWhiteAddress:%X", key[1:])
	case evmtypes.KeyPrefixContractBlockedList[0]:
		return fmt.Sprintf("contractBlockedAddres:%X;methods:%s", key[1:], value)
	default:
		return defaultKvFormatter(key, value)
	}
}

// ReadTree loads an iavl tree from the directory
// If version is 0, load latest, otherwise, load named version
// The prefix represents which iavl tree you want to read. The iaviwer will always set a prefix.
func ReadTree(db dbm.DB, version int, prefix []byte, cacheSize int) (*iavl.MutableTree, error) {
	if len(prefix) != 0 {
		db = dbm.NewPrefixDB(db, prefix)
	}

	tree, err := iavl.NewMutableTree(db, cacheSize)
	if err != nil {
		return nil, err
	}
	_, err = tree.LoadVersion(int64(version))
	return tree, err
}

func OpenDB(dir string, backend dbm.BackendType) (db dbm.DB, err error) {
	switch {
	case strings.HasSuffix(dir, ".db"):
		dir = dir[:len(dir)-3]
	case strings.HasSuffix(dir, ".db/"):
		dir = dir[:len(dir)-4]
	default:
		return nil, fmt.Errorf("database directory must end with .db")
	}
	//doesn't work on windows!
	cut := strings.LastIndex(dir, "/")
	if cut == -1 {
		return nil, fmt.Errorf("cannot cut paths on %s", dir)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("couldn't create db: %v", r)
		}
	}()
	name := dir[cut+1:]
	db = dbm.NewDB(name, backend, dir[:cut])
	return db, nil
}

// parseWeaveKey assumes a separating : where all in front should be ascii,
// and all afterwards may be ascii or binary
func parseWeaveKey(key []byte) string {
	cut := bytes.IndexRune(key, ':')
	if cut == -1 {
		return encodeID(key)
	}
	prefix := key[:cut]
	id := key[cut+1:]
	return fmt.Sprintf("%s:%s", encodeID(prefix), encodeID(id))
}

// casts to a string if it is printable ascii, hex-encodes otherwise
func encodeID(id []byte) string {
	for _, b := range id {
		if b < 0x20 || b >= 0x80 {
			return strings.ToUpper(hex.EncodeToString(id))
		}
	}
	return string(id)
}

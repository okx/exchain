package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

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

	flagStart  = "start"
	flagLimit  = "limit"
	flagHex    = "hex"
	flagPrefix = "prefix"
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
	iavlCtx.flags.Prefix = cmd.PersistentFlags().String(flagPrefix, "", "the prefix of keys, module value must be \"\" if prefix is set")
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
	cmd := &cobra.Command{
		Use:   "diff [data_dir] [compare_data_dir] [height] [module]",
		Short: "Read different key-value from leveldb according two paths",
		PreRun: func(cmd *cobra.Command, args []string) {
			iaviewerCmdParseFlags(ctx)
		},
		Run: func(cmd *cobra.Command, args []string) {
			var moduleList []string
			if len(args) == 4 {
				moduleList = []string{args[3]}
			} else {
				moduleList = make([]string, 0, len(app.ModuleBasics))
				for m := range app.ModuleBasics {
					moduleList = append(moduleList, fmt.Sprintf("s/k:%s/", m))
				}
			}
			height, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				panic("The input height is wrong")
			}
			iaviewerPrintDiff(ctx.Codec, args[0], ctx.DbBackend, args[1], moduleList, int(height))
		},
	}
	return cmd
}

// iaviewerPrintDiff reads different key-value from leveldb according two paths
func iaviewerPrintDiff(cdc *codec.Codec, dataDir string, backend dbm.BackendType, compareDir string, modules []string, height int) {
	if dataDir == compareDir {
		log.Fatal("data_dit and compare_data_dir should not be the same")
	}
	db, err := OpenDB(dataDir, backend)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}
	defer db.Close()

	compareDB, err := OpenDB(compareDir, backend)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}
	defer compareDB.Close()

	for _, module := range modules {
		//get all key-values
		tree, err := ReadTree(db, height, []byte(module), DefaultCacheSize)
		if err != nil {
			log.Println("Error reading data: ", err)
			os.Exit(1)
		}
		compareTree, err := ReadTree(compareDB, height, []byte(module), DefaultCacheSize)
		if err != nil {
			log.Println("Error reading compareTree data: ", err)
			os.Exit(1)
		}
		if bytes.Equal(tree.Hash(), compareTree.Hash()) {
			continue
		}

		var wg sync.WaitGroup
		wg.Add(2)
		dataMap := make(map[string][32]byte, tree.Size())
		compareDataMap := make(map[string][32]byte, compareTree.Size())
		go getKVs(tree, dataMap, &wg)
		go getKVs(compareTree, compareDataMap, &wg)
		wg.Wait()

		//get all keys
		keySize := tree.Size()
		if compareTree.Size() > keySize {
			keySize = compareTree.Size()
		}
		allKeys := make(map[string]bool, keySize)
		for k, _ := range dataMap {
			allKeys[k] = false
		}
		for k, _ := range compareDataMap {
			allKeys[k] = false
		}

		log.Println(fmt.Sprintf("==================================== %s begin ====================================", module))
		//find diff value by each key
		for key, _ := range allKeys {
			value, ok := dataMap[key]
			compareValue, compareOK := compareDataMap[key]
			keyByte, _ := hex.DecodeString(key)
			if ok && compareOK {
				if value == compareValue {
					continue
				}
				log.Println("\nvalue is different--------------------------------------------------------------------")
				log.Println("dir key-value :")
				printByKey(cdc, tree, module, keyByte)
				log.Println("compareDir key-value :")
				printByKey(cdc, compareTree, module, keyByte)
				log.Println("value is different--------------------------------------------------------------------")
				continue
			}
			if ok {
				log.Println("\nOnly be in dir--------------------------------------------------------------------")
				printByKey(cdc, tree, module, keyByte)
				continue
			}
			if compareOK {
				log.Println("\nOnly be in compare dir--------------------------------------------------------------------")
				printByKey(cdc, compareTree, module, keyByte)
				continue
			}

		}
		log.Println(fmt.Sprintf("==================================== %s end ====================================", module))
	}
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
	if tree.Size() <= int64(ctx.Start) {
		return
	}
	if ctx.Start != 0 {
		startKey, _ = tree.GetByIndex(int64(ctx.Start))
	}
	if ctx.Limit != 0 && int64(ctx.Start+ctx.Limit) < tree.Size() {
		endKey, _ = tree.GetByIndex(int64(ctx.Start + ctx.Limit))
	}

	tree.IterateRange(startKey, endKey, true, func(key []byte, value []byte) bool {
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
		return fmt.Sprintf("tokenNumber:%x", tokenNumber)
	case tokentypes.PrefixUserTokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		//address-token:tokenInfo
		return fmt.Sprintf("%s-%s:token:%s", hex.EncodeToString(key[1:21]), string(key[21:]), token.String())
	default:
		return defaultKvFormatter(key, value)
	}
}

func mainPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	if bytes.Equal(key, []byte("consensus_params")) {
		return fmt.Sprintf("consensusParams:%s", hex.EncodeToString(value))
	}
	return defaultKvFormatter(key, value)
}

func slashingPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case slashingtypes.ValidatorSigningInfoKey[0]:
		var signingInfo slashingtypes.ValidatorSigningInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &signingInfo)
		return fmt.Sprintf("validatorAddr:%s:signingInfo:%s", hex.EncodeToString(key[1:]), signingInfo.String())
	case slashingtypes.ValidatorMissedBlockBitArrayKey[0]:
		return fmt.Sprintf("validatorMissedBlockAddr:%s:index:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	case slashingtypes.AddrPubkeyRelationKey[0]:
		return fmt.Sprintf("pubkeyAddr:%s:pubkey:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
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
		return fmt.Sprintf("proposerKey:%s", hex.EncodeToString(value))
	case distypes.DelegatorWithdrawAddrPrefix[0]:
		return fmt.Sprintf("delegatorWithdrawAddr:%s:address:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	case distypes.ValidatorAccumulatedCommissionPrefix[0]:
		var commission types.ValidatorAccumulatedCommission
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &commission)
		return fmt.Sprintf("validatorAccumulatedAddr:%s:address:%s", hex.EncodeToString(key[1:]), commission.String())
	default:
		return defaultKvFormatter(key, value)
	}
}

func govPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case govtypes.ProposalsKeyPrefix[0]:
		return fmt.Sprintf("proposalId:%x;power:%x", binary.BigEndian.Uint64(key[1:]), hex.EncodeToString(value))
	case govtypes.ActiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		return fmt.Sprintf("activeProposalEndTime:%x;proposalId:%x", time.String(), binary.BigEndian.Uint64(value))
	case govtypes.InactiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		return fmt.Sprintf("inactiveProposalEndTime:%x;proposalId:%x", time.String(), binary.BigEndian.Uint64(value))
	case govtypes.ProposalIDKey[0]:
		return fmt.Sprintf("proposalId:%x", hex.EncodeToString(value))
	default:
		return defaultKvFormatter(key, value)
	}
}

func stakingPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case stakingtypes.LastValidatorPowerKey[0]:
		var power int64
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		return fmt.Sprintf("validatorAddress:%s;power:%x", hex.EncodeToString(key[1:]), power)
	case stakingtypes.LastTotalPowerKey[0]:
		var power sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		return fmt.Sprintf("lastTotolValidatorPower:%s", power.String())
	case stakingtypes.ValidatorsKey[0]:
		var validator stakingtypes.Validator
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &validator)
		return fmt.Sprintf("validator:%s;info:%s", hex.EncodeToString(key[1:]), validator)
	case stakingtypes.ValidatorsByConsAddrKey[0]:
		return fmt.Sprintf("validatorConsAddrKey:%s;address:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	case stakingtypes.ValidatorsByPowerIndexKey[0]:
		return fmt.Sprintf("validatorPowerIndexKey:%s;address:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	default:
		return defaultKvFormatter(key, value)
	}
}

func paramsPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	return fmt.Sprintf("%s:%s", string(key), string(value))
}

func accPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	if key[0] == acctypes.AddressStoreKeyPrefix[0] {
		var acc exported.Account
		bz := value
		cdc.MustUnmarshalBinaryBare(bz, &acc)
		return fmt.Sprintf("adress:%s;account:%s", hex.EncodeToString(key[1:]), acc.String())
	} else if bytes.Equal(key, acctypes.GlobalAccountNumberKey) {
		return fmt.Sprintf("%s:%s", string(key), hex.EncodeToString(value))
	} else {
		return defaultKvFormatter(key, value)
	}
}

func evmPrintKey(cdc *codec.Codec, key []byte, value []byte) string {
	switch key[0] {
	case evmtypes.KeyPrefixBlockHash[0]:
		return fmt.Sprintf("blockHash:%s;height:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	case evmtypes.KeyPrefixBloom[0]:
		return fmt.Sprintf("bloomHeight:%s;data:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	case evmtypes.KeyPrefixCode[0]:
		return fmt.Sprintf("code:%s;data:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	case evmtypes.KeyPrefixStorage[0]:
		return fmt.Sprintf("stroageHash:%s;keyHash:%s;data:%s", hex.EncodeToString(key[1:40]), hex.EncodeToString(key[41:]), hex.EncodeToString(value))
	case evmtypes.KeyPrefixChainConfig[0]:
		bz := value
		var config evmtypes.ChainConfig
		cdc.MustUnmarshalBinaryBare(bz, &config)
		return fmt.Sprintf("chainCofig:%s", config.String())
	case evmtypes.KeyPrefixHeightHash[0]:
		return fmt.Sprintf("height:%s;blockHash:%s", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
	case evmtypes.KeyPrefixContractDeploymentWhitelist[0]:
		return fmt.Sprintf("whiteAddress:%s", hex.EncodeToString(key[1:]))
	case evmtypes.KeyPrefixContractBlockedList[0]:
		return fmt.Sprintf("blockedAddres:%s", hex.EncodeToString(key[1:]))
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

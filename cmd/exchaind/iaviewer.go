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
)

var printKeysDict = map[string]printKey{
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

func iaviewerCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iaviewer",
		Short: "Iaviewer key-value from leveldb",
	}

	cmd.AddCommand(
		iaviewerReadCmd(ctx, cdc),
		readDiff(ctx, cdc),
		iaviewerListModulesCmd(),
	)
	cmd.PersistentFlags().String(flagDBBackend, "goleveldb", "Database backend: goleveldb | rocksdb")
	return cmd
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

func iaviewerReadCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read <data_dir> <module> [version]",
		Short: "Read iavl tree key-value from db",
		Long:  "Read iavl tree key-value from db, you must specify data_dir and module, if version is 0 or not specified, read data from the latest version.\n",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
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

			dbBackend := dbm.GoLevelDBBackend
			dbBackendStr := viper.GetString(flagDBBackend)
			if dbBackendStr != "" {
				dbBackend = dbm.BackendType(dbBackendStr)
			}

			return iaviewerReadData(cdc, dataDir, dbBackend, module, version)
		},
	}
	return cmd
}

func readDiff(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [data_dir] [compare_data_dir] [height] [module]",
		Short: "Read different key-value from leveldb according two paths",
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
			iaviewerPrintDiff(cdc, args[0], dbm.BackendType(ctx.Config.DBBackend), args[1], moduleList, int(height))
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
	defer db.Close()
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}
	compareDB, err := OpenDB(compareDir, backend)
	defer compareDB.Close()
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}

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
func iaviewerReadData(cdc *codec.Codec, dataDir string, backend dbm.BackendType, module string, version int) error {
	db, err := OpenDB(dataDir, backend)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}

	modulePrefix := fmt.Sprintf("s/k:%s/", module)

	fmt.Printf("==================================== %s begin ====================================\n", module)
	tree, err := ReadTree(db, version, []byte(modulePrefix), DefaultCacheSize)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	printTree(cdc, module, tree)
	fmt.Printf("Hash: %X\n", tree.Hash())
	fmt.Printf("Size: %d\n", tree.Size())
	fmt.Printf("==================================== %s end ====================================\n\n", module)
	return nil
}

// getKVs, get all key-values by mutableTree
func getKVs(tree *iavl.MutableTree, dataMap map[string][32]byte, wg *sync.WaitGroup) {
	tree.Iterate(func(key []byte, value []byte) bool {
		dataMap[hex.EncodeToString(key)] = sha256.Sum256(value)
		return false
	})
	wg.Done()
}

type (
	printKey func(cdc *codec.Codec, key []byte, value []byte)
)

func printTree(cdc *codec.Codec, module string, tree *iavl.MutableTree) {
	tree.Iterate(func(key []byte, value []byte) bool {
		if impl, exit := printKeysDict[module]; exit {
			impl(cdc, key, value)
		} else {
			printKey := parseWeaveKey(key)
			digest := hex.EncodeToString(value)
			log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
		}
		return false
	})
}

func printByKey(cdc *codec.Codec, tree *iavl.MutableTree, module string, key []byte) {
	_, value := tree.Get(key)
	if impl, exit := printKeysDict[module]; exit {
		impl(cdc, key, value)
	} else {
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func supplyPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case supplytypes.SupplyKey[0]:
		var supplyAmount sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &supplyAmount)
		log.Println(fmt.Sprintf("tokenSymbol:%s:info:%s\n", string(key[1:]), supplyAmount.String()))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

type MinterCustom struct {
	NextBlockToUpdate uint64       `json:"next_block_to_update" yaml:"next_block_to_update"` // record the block height for next year
	MintedPerBlock    sdk.DecCoins `json:"minted_per_block" yaml:"minted_per_block"`         // record the MintedPerBlock per block in this year
}

func mintPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case minttypes.MinterKey[0]:
		var minter MinterCustom
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &minter)
		log.Println(fmt.Sprintf("minter:%v\n", minter))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func tokenPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case tokentypes.TokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		log.Println(fmt.Sprintf("tokenName:%s:info:%s\n", string(key[1:]), token.String()))
		return
	case tokentypes.TokenNumberKey[0]:
		var tokenNumber uint64
		cdc.MustUnmarshalBinaryBare(value, &tokenNumber)
		log.Println(fmt.Sprintf("tokenNumber:%x\n", tokenNumber))
		return
	case tokentypes.PrefixUserTokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		//address-token:tokenInfo
		log.Println(fmt.Sprintf("%s-%s:token:%s\n", hex.EncodeToString(key[1:21]), string(key[21:]), token.String()))
		return

	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func mainPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	if bytes.Equal(key, []byte("consensus_params")) {
		log.Println(fmt.Sprintf("consensusParams:%s\n", hex.EncodeToString(value)))
		return
	}
	printKey := parseWeaveKey(key)
	digest := hex.EncodeToString(value)
	log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
}

func slashingPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case slashingtypes.ValidatorSigningInfoKey[0]:
		var signingInfo slashingtypes.ValidatorSigningInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &signingInfo)
		log.Println(fmt.Sprintf("validatorAddr:%s:signingInfo:%s\n", hex.EncodeToString(key[1:]), signingInfo.String()))
		return
	case slashingtypes.ValidatorMissedBlockBitArrayKey[0]:
		log.Println(fmt.Sprintf("validatorMissedBlockAddr:%s:index:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case slashingtypes.AddrPubkeyRelationKey[0]:
		log.Println(fmt.Sprintf("pubkeyAddr:%s:pubkey:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func distributionPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case distypes.FeePoolKey[0]:
		var feePool distypes.FeePool
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &feePool)
		log.Println(fmt.Sprintf("feePool:%v\n", feePool))
		return
	case distypes.ProposerKey[0]:
		log.Println(fmt.Sprintf("proposerKey:%s\n", hex.EncodeToString(value)))
		return
	case distypes.DelegatorWithdrawAddrPrefix[0]:
		log.Println(fmt.Sprintf("delegatorWithdrawAddr:%s:address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case distypes.ValidatorAccumulatedCommissionPrefix[0]:
		var commission types.ValidatorAccumulatedCommission
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &commission)
		log.Println(fmt.Sprintf("validatorAccumulatedAddr:%s:address:%s\n", hex.EncodeToString(key[1:]), commission.String()))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func govPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case govtypes.ProposalsKeyPrefix[0]:
		log.Println(fmt.Sprintf("proposalId:%x;power:%x\n", binary.BigEndian.Uint64(key[1:]), hex.EncodeToString(value)))
		return
	case govtypes.ActiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		log.Println(fmt.Sprintf("activeProposalEndTime:%x;proposalId:%x\n", time.String(), binary.BigEndian.Uint64(value)))
		return
	case govtypes.InactiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		log.Println(fmt.Sprintf("inactiveProposalEndTime:%x;proposalId:%x\n", time.String(), binary.BigEndian.Uint64(value)))
		return
	case govtypes.ProposalIDKey[0]:
		log.Println(fmt.Sprintf("proposalId:%x\n", hex.EncodeToString(value)))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func stakingPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case stakingtypes.LastValidatorPowerKey[0]:
		var power int64
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		log.Println(fmt.Sprintf("validatorAddress:%s;power:%x\n", hex.EncodeToString(key[1:]), power))
		return
	case stakingtypes.LastTotalPowerKey[0]:
		var power sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		log.Println(fmt.Sprintf("lastTotolValidatorPower:%s\n", power.String()))
		return
	case stakingtypes.ValidatorsKey[0]:
		var validator stakingtypes.Validator
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &validator)
		log.Println(fmt.Sprintf("validator:%s;info:%s\n", hex.EncodeToString(key[1:]), validator))
		return
	case stakingtypes.ValidatorsByConsAddrKey[0]:
		log.Println(fmt.Sprintf("validatorConsAddrKey:%s;address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case stakingtypes.ValidatorsByPowerIndexKey[0]:
		log.Println(fmt.Sprintf("validatorPowerIndexKey:%s;address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func paramsPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	log.Println(fmt.Sprintf("%s:%s\n", string(key), string(value)))
}

func accPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	if key[0] == acctypes.AddressStoreKeyPrefix[0] {
		var acc exported.Account
		bz := value
		cdc.MustUnmarshalBinaryBare(bz, &acc)
		log.Println(fmt.Sprintf("adress:%s;account:%s\n", hex.EncodeToString(key[1:]), acc.String()))
		return
	} else if bytes.Equal(key, acctypes.GlobalAccountNumberKey) {
		log.Println(fmt.Sprintf("%s:%s\n", string(key), hex.EncodeToString(value)))
		return
	} else {
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func evmPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case evmtypes.KeyPrefixBlockHash[0]:
		log.Println(fmt.Sprintf("blockHash:%s;height:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixBloom[0]:
		log.Println(fmt.Sprintf("bloomHeight:%s;data:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixCode[0]:
		log.Println(fmt.Sprintf("code:%s;data:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixStorage[0]:
		log.Println(fmt.Sprintf("stroageHash:%s;keyHash:%s;data:%s\n", hex.EncodeToString(key[1:40]), hex.EncodeToString(key[41:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixChainConfig[0]:
		bz := value
		var config evmtypes.ChainConfig
		cdc.MustUnmarshalBinaryBare(bz, &config)
		log.Println(fmt.Sprintf("chainCofig:%s\n", config.String()))
		return
	case evmtypes.KeyPrefixHeightHash[0]:
		log.Println(fmt.Sprintf("height:%s;blockHash:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixContractDeploymentWhitelist[0]:
		log.Println(fmt.Sprintf("whiteAddress:%s\n", hex.EncodeToString(key[1:])))
		return
	case evmtypes.KeyPrefixContractBlockedList[0]:
		log.Println(fmt.Sprintf("blockedAddres:%s\n", hex.EncodeToString(key[1:])))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
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
	ver, err := tree.LoadVersion(int64(version))
	log.Println(fmt.Sprintf("%s Got version: %d\n", string(prefix), ver))
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

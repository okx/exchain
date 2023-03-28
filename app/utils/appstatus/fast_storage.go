package appstatus

import (
	"fmt"
	"math"
	"path/filepath"

	bam "github.com/okx/okbchain/libs/cosmos-sdk/baseapp"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	capabilitytypes "github.com/okx/okbchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/params"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/supply"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/upgrade"
	"github.com/okx/okbchain/libs/iavl"
	ibctransfertypes "github.com/okx/okbchain/libs/ibc-go/modules/apps/transfer/types"
	ibchost "github.com/okx/okbchain/libs/ibc-go/modules/core/24-host"
	dbm "github.com/okx/okbchain/libs/tm-db"
	distr "github.com/okx/okbchain/x/distribution"
	"github.com/okx/okbchain/x/erc20"
	"github.com/okx/okbchain/x/evidence"
	"github.com/okx/okbchain/x/feesplit"
	"github.com/okx/okbchain/x/gov"
	"github.com/okx/okbchain/x/slashing"
	staking "github.com/okx/okbchain/x/staking/types"
	token "github.com/okx/okbchain/x/token/types"
	"github.com/spf13/viper"
)

const (
	applicationDB = "application"
	dbFolder      = "data"
)

func GetAllStoreKeys() []string {
	return []string{
		bam.MainStoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, upgrade.StoreKey, evidence.StoreKey,
		token.StoreKey, token.KeyLock,
		ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		ibchost.StoreKey,
		erc20.StoreKey,
		// mpt.StoreKey,
		// wasm.StoreKey,
		feesplit.StoreKey,
	}
}

func IsFastStorageStrategy() bool {
	return checkFastStorageStrategy(GetAllStoreKeys())
}

func checkFastStorageStrategy(storeKeys []string) bool {
	home := viper.GetString(flags.FlagHome)
	dataDir := filepath.Join(home, dbFolder)
	db, err := sdk.NewDB(applicationDB, dataDir)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, v := range storeKeys {
		if !isFss(db, v) {
			return false
		}
	}

	return true
}

func isFss(db dbm.DB, storeKey string) bool {
	prefix := fmt.Sprintf("s/k:%s/", storeKey)
	prefixDB := dbm.NewPrefixDB(db, []byte(prefix))

	return iavl.IsFastStorageStrategy(prefixDB)
}

func GetFastStorageVersion() int64 {
	home := viper.GetString(flags.FlagHome)
	dataDir := filepath.Join(home, dbFolder)
	db, err := sdk.NewDB(applicationDB, dataDir)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storeKeys := GetAllStoreKeys()
	var ret int64 = math.MaxInt64
	for _, v := range storeKeys {
		version := getVersion(db, v)
		if version < ret {
			ret = version
		}
	}

	return ret
}

func getVersion(db dbm.DB, storeKey string) int64 {
	prefix := fmt.Sprintf("s/k:%s/", storeKey)
	prefixDB := dbm.NewPrefixDB(db, []byte(prefix))

	version, _ := iavl.GetFastStorageVersion(prefixDB)

	return version
}

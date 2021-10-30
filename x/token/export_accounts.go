package token

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/spf13/viper"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
)

const (
	FlagOSSEnable          = "oss-enable"
	FlagOSSEndpoint        = "oss-endpoint"
	FlagOSSAccessKeyID     = "oss-access-key-id"
	FlagOSSAccessKeySecret = "oss-access-key-secret"
	FlagOSSBucketName      = "oss-bucket-name"
	FlagOSSObjectPath      = "oss-object-path"
)

var (
	logFileName = "export-upload-account.log"
)

type AccType int

const (
	userAccount AccType = iota
	contractAccount
)

func exportAccounts(ctx sdk.Context, keeper Keeper) (filePath string) {
	pt := time.Now().UTC().Format(time.RFC3339)
	rootDir := viper.GetString(cli.HomeFlag)

	accFileName := fmt.Sprintf("accounts-%d-%s.csv", ctx.BlockHeight(), pt)

	// 1. open log file
	logFile, logWr, err := openLogFile()
	if err != nil {
		return
	}
	defer logFile.Close()
	defer logWr.Flush()

	recodeLog(logWr, "===============")
	recodeLog(logWr, fmt.Sprintf("time: %s", pt))
	recodeLog(logWr, fmt.Sprintf("height: %d", ctx.BlockHeight()))
	recodeLog(logWr, fmt.Sprintf("file name: %s", accFileName))

	// 2. open account file
	accFile, err := os.OpenFile(path.Join(rootDir, accFileName), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("open account file error: %s", err))
		return
	}
	defer accFile.Close()
	accWr := bufio.NewWriter(accFile)
	defer accWr.Flush()
	defer func() {
		if err := recover(); err != nil {
			recodeLog(logWr, fmt.Sprintf("export accounts panic: %s", err))
		}
	}()

	count := 0
	startTime := time.Now()
	keeper.accountKeeper.IterateAccounts(ctx, func(account authexported.Account) bool {
		ethAcc, ok := account.(*ethermint.EthAccount)
		if !ok {
			return false
		}

		//account.SpendableCoins()
		oktBalance := account.GetCoins().AmountOf(sdk.DefaultBondDenom)
		if !oktBalance.GT(sdk.ZeroDec()) {
			return false
		}

		accType := userAccount
		if !bytes.Equal(ethAcc.CodeHash, ethcrypto.Keccak256(nil)) {
			accType = contractAccount
		}

		csvStr := fmt.Sprintf("%s,%d,%s,%d,%s",
			ethAcc.EthAddress().String(),
			accType,
			oktBalance.String(),
			ctx.BlockHeight(),
			pt,
		)
		fmt.Fprintln(accWr, csvStr)
		count++
		return false
	})
	recodeLog(logWr, fmt.Sprintf("count: %d", count))
	recodeLog(logWr, fmt.Sprintf("export duration: %s", time.Since(startTime).String()))
	return path.Join(rootDir, accFileName)
}

func uploadOSS(filePath string) {
	// 1. open log file
	logFile, logWr, err := openLogFile()
	if err != nil {
		return
	}
	defer logFile.Close()
	defer logWr.Flush()
	defer func() {
		if err := recover(); err != nil {
			recodeLog(logWr, fmt.Sprintf("upload OSS panic: %s", err))
		}
	}()

	startTime := time.Now()
	// create OSSClient
	ossClient, err := oss.New(viper.GetString(FlagOSSEndpoint), viper.GetString(FlagOSSAccessKeyID), viper.GetString(FlagOSSAccessKeySecret))
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("creates oss lcient error: %s", err))
		return
	}

	// gets the bucket instance
	bucket, err := ossClient.Bucket(viper.GetString(FlagOSSBucketName))
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("gets the bucket instance error: %s", err))
		return
	}

	_, fileName := path.Split(filePath)
	objectName := viper.GetString(FlagOSSObjectPath) + fmt.Sprintf("accounts-%s/", time.Now().Format("20060102")) + fileName
	// multipart file upload
	err = bucket.UploadFile(objectName, filePath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("multipart file upload error: %s", err))
		return
	}
	recodeLog(logWr, fmt.Sprintf("oss file: %s", objectName))
	recodeLog(logWr, fmt.Sprintf("upload duration: %s", time.Since(startTime).String()))
}

func openLogFile() (*os.File, *bufio.Writer, error) {
	rootDir := viper.GetString(cli.HomeFlag)

	file, err := os.OpenFile(path.Join(rootDir, logFileName), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	logWr := bufio.NewWriter(file)
	return file, logWr, nil
}
func recodeLog(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}

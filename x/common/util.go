package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/okex/exchain/x/params/subspace"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	apptypes "github.com/okex/exchain/app/types"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func InitConfig() {
	config := sdk.GetConfig()
	if config.GetBech32ConsensusAddrPrefix() == apptypes.Bech32PrefixConsAddr {
		return
	}
	apptypes.SetBech32Prefixes(config)
	apptypes.SetBip44CoinType(config)
	config.Seal()
}

// Int64ToBytes converts int64 to bytes
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BytesToInt64 converts bytes to int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// Paginate converts page params for a paginated query,
func Paginate(pageStr, perPageStr string) (page int, perPage int, err error) {
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return
		}
	}
	if perPageStr != "" {
		perPage, err = strconv.Atoi(perPageStr)
		if err != nil {
			return
		}
	}
	if page < 0 || perPage < 0 {
		err = fmt.Errorf("negative page %d or per_page %d is invalid", page, perPage)
		return
	}
	return
}

// GetPage returns the offset and limit for data query
func GetPage(page, perPage int) (offset, limit int) {
	if page <= 0 || perPage <= 0 {
		return
	}
	offset = (page - 1) * perPage
	limit = perPage
	return
}

// HandleErrorMsg handles the error msg
func HandleErrorMsg(w http.ResponseWriter, cliCtx context.CLIContext, code uint32, msg string) {
	response := GetErrorResponseJSON(code, msg, msg)
	rest.PostProcessResponse(w, cliCtx, response)
}

// HasSufficientCoins checks whether the account has sufficient coins
func HasSufficientCoins(addr sdk.AccAddress, availableCoins, amt sdk.Coins) (err error) {
	//availableCoins := availCoins[:]
	if !amt.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}
	if !availableCoins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}

	_, hasNeg := availableCoins.SafeSub(amt)
	if hasNeg {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds,
			fmt.Sprintf("insufficient account funds;address: %s, availableCoin: %s, needCoin: %s",
				addr.String(), availableCoins, amt),
		)
	}
	return nil
}

// SkipSysTestChecker is supported to used in System Unit Test
// (described in http://gitlab.okcoin-inc.com/dex/exchain/issues/472)
// if System environment variables "SYS_TEST_ALL" is set to 1, all of the system test will be enable. \n
// if System environment variables "ORM_MYSQL_SYS_TEST" is set to 1,
// 				all of the system test in orm_mysql_sys_test.go will be enble.
func SkipSysTestChecker(t *testing.T) {
	_, fname, _, ok := runtime.Caller(0)
	enable := ok
	if enable {
		enableAllEnv := "SYS_TEST_ALL"

		sysTestName := strings.Split(fname, ".go")[0]
		enableCurrent := strings.ToUpper(sysTestName)

		enable = os.Getenv(enableAllEnv) == "1" ||
			(strings.HasSuffix(sysTestName, "sys_test") && os.Getenv(enableCurrent) == "1")
	}

	if !enable {
		t.SkipNow()
	}
}

// mulAndQuo returns a * b / c
func MulAndQuo(a, b, c sdk.Dec) sdk.Dec {
	// 10^8
	auxiliaryDec := sdk.NewDecFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil))
	a = a.MulTruncate(auxiliaryDec)
	return a.MulTruncate(b).QuoTruncate(c).QuoTruncate(auxiliaryDec)
}

// BlackHoleAddress returns the black hole address
func BlackHoleAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHex(blackHoleHex)
	return addr
}

func GetFixedLengthRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func PanicTrace(kb int) {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	fmt.Print(string(stack))
}

func SanityCheckHandler(res *sdk.Result, err error) {
	if res == nil && err == nil {
		panic("Invalid handler")
	}
	if res != nil && err != nil {
		panic("Invalid handler")
	}
}

func ValidateSysCoin(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(sdk.SysCoin)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if !v.IsValid() {
			return fmt.Errorf("invalid %s: %s", param, v)
		}

		return nil
	}
}

func ValidateSysCoins(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(sdk.SysCoins)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if !v.IsValid() {
			return fmt.Errorf("invalid %s: %s", param, v)
		}

		return nil
	}
}

func ValidateDurationPositive(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(time.Duration)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if v <= 0 {
			return fmt.Errorf("%s must be positive: %d", param, v)
		}

		return nil
	}
}

func ValidateBool(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		_, ok := i.(bool)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		return nil
	}
}

func ValidateInt64Positive(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(int64)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if v <= 0 {
			return fmt.Errorf("%s must be positive: %d", param, v)
		}

		return nil
	}
}

func ValidateUint64Positive(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(uint64)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if v == 0 {
			return fmt.Errorf("%s must be positive: %d", param, v)
		}

		return nil
	}
}

func ValidateRateNotNeg(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(sdk.Dec)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}
		if v.IsNegative() {
			return fmt.Errorf("%s cannot be negative: %s", param, v)
		}
		if v.GT(sdk.OneDec()) {
			return fmt.Errorf("%s is too large: %s", param, v)
		}
		return nil
	}
}

func ValidateDecPositive(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(sdk.Dec)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}
		if !v.IsPositive() {
			return fmt.Errorf("%s must be positive: %s", param, v)
		}
		return nil
	}
}

func ValidateDenom(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(string)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if sdk.ValidateDenom(v) != nil {
			return fmt.Errorf("invalid %s", param)
		}

		return nil
	}
}

func ValidateUint16Positive(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(uint16)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if v == 0 {
			return fmt.Errorf("%s must be positive: %d", param, v)
		}

		return nil
	}
}

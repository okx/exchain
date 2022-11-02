package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	tm "github.com/buger/goterm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/inancgumus/screen"
	dydxlib "github.com/okex/exchain/libs/dydx"
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/mempool/dydx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	Address string
	Name    string
}

type Config struct {
	ServerAddr string
	Users      []User
}

func DefaultConfig() Config {
	return Config{
		ServerAddr: "localhost:7070",
		Users: []User{
			{
				Name:    "Alice",
				Address: "0x2CF4ea7dF75b513509d95946B43062E26bD88035",
			},
			{
				Name:    "Bob",
				Address: "0x0073F2E28ef8F117e53d858094086Defaf1837D5",
			},
			{
				Name:    "Captain",
				Address: "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F",
			},
		},
	}
}

//var dydxConfig = dydx.DydxConfig{
//	// PrivKeyHex:                 "fefac29bfa769d8a6c17b685816dadbd30e3f395e997ed955a5461914be75ed5",
//	PrivKeyHex:                 "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3",
//	ChainID:                    "65",
//	EthWsRpcUrl:                "wss://exchaintestws.okex.org:8443",
//	EthHttpRpcUrl:              "https://exchaintestrpc.okex.org",
//	PerpetualV1ContractAddress: "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
//	P1OrdersContractAddress:    "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
//	P1MakerOracleAddress:       "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
//	P1MarginAddress:            "0xC87EF36830A0D94E42bB2D82a0b2bB939368b10B",
//}

var dydxConfig = dydx.DydxConfig{
	PrivKeyHex:                 "2438019d3fccd8ffdff4d526c0f7fae4136866130affb3aa375d95835fa8f60f",
	ChainID:                    "8",
	EthWsRpcUrl:                "ws://localhost:8546",
	EthHttpRpcUrl:              "http://localhost:8545",
	PerpetualV1ContractAddress: "0xbc0Bf2Bf737344570c02d8D8335ceDc02cECee71",
	P1OrdersContractAddress:    "0x632D131CCCE01206F08390cB66D1AdEf9b264C61",
	P1MakerOracleAddress:       "0xF306F8B7531561d0f92BA965a163B6C6d422ade1",
	P1MarginAddress:            "0xeb95A3D1f7Ca2B8Ba61F326fC4dA9124b6C057b9",
}

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "dydx-viewer")
	var options []log.Option
	options = append(options, log.AllowDebug())
	logger = log.NewFilter(logger, options...)

	configPath := "./config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	var config Config
	configBz, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error("read config file failed", "err", err)
		config = DefaultConfig()
		logger.Info("use default config", "config", config)
	} else {
		err = json.Unmarshal(configBz, &config)
		if err != nil {
			logger.Error("unmarshal config file failed", "err", err)
		}
	}

	ethCli, err := ethclient.Dial(dydxConfig.EthHttpRpcUrl)
	if err != nil {
		logger.Error("failed to connect to dydx", "err", err)
		return
	}

	dydxContracts, err := dydxlib.NewContracts(
		common.HexToAddress(dydxConfig.PerpetualV1ContractAddress),
		common.HexToAddress(dydxConfig.P1OrdersContractAddress),
		common.HexToAddress(dydxConfig.P1MakerOracleAddress),
		common.HexToAddress(dydxConfig.P1MarginAddress),
		nil,
		ethCli,
	)
	if err != nil {
		logger.Error("get dydx contracts error", "err", err)
		return
	}

	conn, err := grpc.Dial(config.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Error("failed to connect to dydx server", "err", err)
		return
	}
	client := dydx.NewOrderBookUpdaterClient(conn)
	stream, err := client.WatchOrderBookLevel(context.Background(), new(dydx.Empty))
	if err != nil {
		logger.Error("failed to watch order book", "err", err)
		return
	}

	for {
		ob, err := stream.Recv()
		if err != nil {
			logger.Error("failed to receive order book", "err", err)
			return
		}
		var usersBalance []UserBalance
		for _, user := range config.Users {
			b, err := dydxContracts.PerpetualV1.GetAccountBalance(nil, common.HexToAddress(user.Address))
			if err != nil {
				logger.Error("failed to get user balance", "err", err)
				continue
			}
			usersBalance = append(usersBalance, UserBalance{
				User:    user,
				Balance: b,
			})
		}

		Print(ob, usersBalance)

		go func() {
			time.Sleep(3 * time.Second)
			var usersBalance []UserBalance
			for _, user := range config.Users {
				b, err := dydxContracts.PerpetualV1.GetAccountBalance(nil, common.HexToAddress(user.Address))
				if err != nil {
					logger.Error("failed to get user balance", "err", err)
					continue
				}
				usersBalance = append(usersBalance, UserBalance{
					User:    user,
					Balance: b,
				})
			}

			Print(ob, usersBalance)
		}()

		//Print(nil, usersBalance)
		//time.Sleep(3 * time.Second)
	}
}

func clear() {
	tm.Clear()
	screen.Clear()
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Print(orderBook *dydx.OrderBookLevel, usersBalance []UserBalance) {
	clear()
	if orderBook != nil {
		fmt.Println("===========================================")
		fmt.Println("OrderBook:")
		fmt.Println("===========================================")
		fmt.Println("Sell:")
		for i := 0; i < len(orderBook.SellLevels); i++ {
			order := orderBook.SellLevels[len(orderBook.SellLevels)-1-i]
			fmt.Printf("price: %s, amount: %d\n", order.Price, order.Amount)
		}
		fmt.Println()
		fmt.Println("Buy:")
		for _, order := range orderBook.BuyLevels {
			fmt.Printf("price: %s, amount: %d\n", order.Price, order.Amount)
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("Users Balance:")
	fmt.Println("===========================================")
	fmt.Println()
	for _, userBalance := range usersBalance {
		_, _ = fmt.Println(userBalance.String())
		fmt.Println()
	}
	fmt.Println("===========================================")
}

type UserBalance struct {
	User    User
	Balance contracts.P1TypesBalance
}

func (ub UserBalance) String() string {
	margin := ub.Balance.Margin.String()
	position := ub.Balance.Position.String()
	if !ub.Balance.PositionIsPositive {
		position = "-" + position
	}
	if !ub.Balance.MarginIsPositive {
		margin = "-" + margin
	}
	return fmt.Sprintf("%s: margin: %s, position: %s", ub.User.Name, margin, position)
}

package dydx

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	grpc "google.golang.org/grpc"
)

type OrderBookServer struct {
	UnimplementedOrderBookUpdaterServer
	logger      log.Logger
	s           *grpc.Server
	book        *DepthBook
	clientIndex int

	mtx      sync.Mutex
	clientCh map[int]*int64
}

func NewOrderBookServer(book *DepthBook, logger log.Logger) *OrderBookServer {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &OrderBookServer{
		logger:   logger,
		book:     book,
		clientCh: map[int]*int64{},
	}
}

func (s *OrderBookServer) SetLogger(logger log.Logger) {
	s.logger = logger
}

func (s *OrderBookServer) UpdateClient() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for k := range s.clientCh {
		atomic.StoreInt64(s.clientCh[k], 1)
	}
}

func bookToOrderBook(book *DepthBook) *OrderBook {
	return nil
	//book.buyOrders.List()
	//return &OrderBook{
	//	BuyOrders:  bookToOrders(book.buyOrders),
	//	SellOrders: bookToOrders(book.sellOrders),
	//}
}

func bookToLevel(book *DepthBook) *OrderBookLevel {
	buyLevels := []*OrderLevel{{"0", 0}}
	sellLevels := []*OrderLevel{{"0", 0}}
	for _, order := range book.buyOrders.List() {
		if order.LeftAndFrozen().Sign() == 0 {
			continue
		}
		if order.GetLimitPrice().String() == buyLevels[len(buyLevels)-1].Price {
			buyLevels[len(buyLevels)-1].Amount += int64(order.LeftAndFrozen().Uint64())
		} else {
			buyLevels = append(buyLevels, &OrderLevel{
				Price:  new(big.Int).Div(order.GetLimitPrice(), exp18).String(),
				Amount: int64(order.LeftAndFrozen().Uint64()),
			})
		}
	}
	for _, order := range book.sellOrders.List() {
		if order.LeftAndFrozen().Sign() == 0 {
			continue
		}
		if order.GetLimitPrice().String() == sellLevels[len(sellLevels)-1].Price {
			sellLevels[len(sellLevels)-1].Amount += int64(order.LeftAndFrozen().Uint64())
		} else {
			sellLevels = append(sellLevels, &OrderLevel{
				Price:  new(big.Int).Div(order.GetLimitPrice(), exp18).String(),
				Amount: int64(order.LeftAndFrozen().Uint64()),
			})
		}
	}
	return &OrderBookLevel{
		BuyLevels:  buyLevels[1:],
		SellLevels: sellLevels[1:],
	}
	//book.buyOrders.List()
	//return &OrderBook{
	//	BuyOrders:  bookToOrders(book.buyOrders),
	//	SellOrders: bookToOrders(book.sellOrders),
	//}
}

//func (s *OrderBookServer) WatchOrderBook(_ *Empty, stream OrderBookUpdater_WatchOrderBookServer) error {
//	s.mtx.Lock()
//	chIndex := s.clientIndex
//	s.clientIndex++
//	ch := new(int64)
//	s.clientCh[chIndex] = ch
//	s.mtx.Unlock()
//
//	defer func() {
//		s.mtx.Lock()
//		delete(s.clientCh, chIndex)
//		s.mtx.Unlock()
//	}()
//
//	for {
//		select {
//		case <-time.After(1 * time.Second):
//			if atomic.LoadInt64(ch) == 1 {
//				atomic.StoreInt64(ch, 0)
//				b := bookToOrderBook(s.book)
//				if b == nil {
//					continue
//				}
//				if err := stream.Send(b); err != nil {
//					return err
//				}
//			}
//		}
//	}
//}

func (s *OrderBookServer) WatchOrderBookLevel(_ *Empty, stream OrderBookUpdater_WatchOrderBookLevelServer) error {
	s.mtx.Lock()
	chIndex := s.clientIndex
	s.clientIndex++
	ch := new(int64)
	s.clientCh[chIndex] = ch
	s.mtx.Unlock()

	defer func() {
		s.mtx.Lock()
		delete(s.clientCh, chIndex)
		s.mtx.Unlock()
	}()

	init := true

	for {
		select {
		case <-time.After(1 * time.Second):
			if atomic.LoadInt64(ch) == 1 || init {
				atomic.StoreInt64(ch, 0)
				init = false
				b := bookToLevel(s.book)
				if b == nil {
					continue
				}
				if err := stream.Send(b); err != nil {
					return err
				}
			}
		}
	}
}

func (s *OrderBookServer) Start(configPort string) error {
	configPort = strings.ToLower(configPort)
	if configPort == "off" {
		return nil
	} else if configPort == "auto" {
		configPort = "0"
	}

	if port, err := strconv.Atoi(configPort); err == nil {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		} else {
			var options = []grpc.ServerOption{
				//grpc.KeepaliveParams(keepalive.ServerParameters{
				//	Time:    30 * time.Second,
				//	Timeout: 10 * time.Second,
				//}),
				//grpc.ForceServerCodec(newGogoCodec()),
			}
			s.s = grpc.NewServer(options...)
			RegisterOrderBookUpdaterServer(s.s, s)
			s.logger.Info("orderbook grpc server listening on", "addr", lis.Addr().String())
			go func() {
				if err := s.s.Serve(lis); err != nil {
					s.logger.Error("Failed to start orderbook grpc server", "err", err)
				}
			}()
		}
	}
	return nil
}

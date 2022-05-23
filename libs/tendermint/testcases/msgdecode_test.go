package testcases

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	mempl "github.com/okex/exchain/libs/tendermint/mempool"
	"github.com/okex/exchain/libs/tendermint/types"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/tendermint/consensus"
	"github.com/okex/exchain/libs/tendermint/libs/autofile"
	"github.com/okex/exchain/libs/tendermint/p2p/conn"
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()
var mtx sync.Mutex

type WriteWalFn func(msg consensus.Message)
type RoutineFinishedFn func()

type MsgProcessRoutine struct {
	msgCh chan []byte
	//done       chan int8
	writeWalFn WriteWalFn
	finishedFn RoutineFinishedFn
	stopSend   bool
	msgNum     int
	index      int
	decode     bool

	reactor *mempl.Reactor
}

func newMsgProcessRoutine(index int, decode bool, wwFn WriteWalFn, fFn RoutineFinishedFn) *MsgProcessRoutine {
	m := &MsgProcessRoutine{
		msgCh:      make(chan []byte, 1000000),
		writeWalFn: wwFn,
		finishedFn: fFn,
		index:      index,
		decode:     decode,
	}

	//config := cfg.TestConfig()
	//const N = 4
	//m.reactor = mempl.NewReactor(nil, nil)

	return m
}

func (mr *MsgProcessRoutine) start(count int) {
	//mr.done = make(chan int8)
	mr.msgNum = count
	go mr.receiveRoutine()
}

func (mr *MsgProcessRoutine) receiveRoutine() {
	count := 0
	start := time.Now()
	for {
		select {
		//case <-mr.done:
		//	mr.finishedFn()
		//	return
		case msgBytes := <-mr.msgCh:
			msgpkt := conn.PacketMsg{}
			//err := msgpkt.UnmarshalFromAmino(cdc, msgBytes[4:])
			//if err != nil {
			//	fmt.Println(err)
			//}

			if mr.decode {
				switch msgpkt.ChannelID {
				case consensus.StateChannel:
					msg, err := consensus.DecodeMsg(msgpkt.Bytes)
					if err != nil {
						fmt.Println(err)
					}

					if err = msg.ValidateBasic(); err != nil {
						fmt.Println(err)
					}

					//wal.Write(msg)
					mr.writeWalFn(msg)

				//case mempl.MempoolChannel:
				default:
					msg, err := mempl.DecodeMsg(msgBytes)
					if err != nil {
						fmt.Println(err)
					}
					switch msg := msg.(type) {
					case *mempl.TxMessage:
						tx := msg.Tx
						var retHash [sha256.Size]byte
						//mtx.Lock()
						copy(retHash[:], tx.Hash(types.GetVenusHeight())[:sha256.Size])
						//mtx.Unlock()
					}
				}
			} else {
				// compute hash for msg bytes
				ethcrypto.Keccak256Hash(msgpkt.Bytes)
			}

			count++
			if count == mr.msgNum {
				mr.finishedFn()
				fmt.Print(time.Since(start).Milliseconds(), ", ")
				return
			}
		}
	}
}

type MsgProcessMgr struct {
	mtx            sync.Mutex
	wal            *consensus.BaseWAL
	done           chan int8
	end            chan int8
	walCh          chan consensus.Message
	routineList    []*MsgProcessRoutine
	prevoteBytes   []byte
	precommitBytes []byte
	hasVoteBytes   []byte
	newStepBytes   []byte
	blockpartBytes []byte
	txBytes []byte
	finished       int
	decode         bool
}

func newMsgProcessMgr(count int, decode bool) *MsgProcessMgr {
	m := &MsgProcessMgr{
		walCh:  make(chan consensus.Message, 1000000),
		end:    make(chan int8),
		decode: decode,
	}
	m.routineList = make([]*MsgProcessRoutine, 0, count)
	for i := 0; i < count; i++ {
		r := newMsgProcessRoutine(i, decode, m.newMsgToWal, m.routineFinished)
		m.routineList = append(m.routineList, r)
	}

	walDir, _ := os.Getwd()
	walFile := filepath.Join(walDir, "wal")
	m.wal, _ = consensus.NewWAL(walFile,
		autofile.GroupHeadSizeLimit(4096),
		autofile.GroupCheckDuration(1*time.Millisecond),
	)
	m.wal.Start()

	prevote := "b05b4f2c082210011abd01c1da2d940ab6010802100222480a207a81b9db16c3cb82bae88a058bb795620e4f2caca193a1141006bc82fc37efaf1224080c1220b0e26207ce3829dd522d5c1b6001975ca2ce17f1d3c3ab9e494407e1f5841b432a0c08fdcb8d940610f8f8ffad023214821ba59ec2a5e1152f25ded9691b0b7d3b0af677380142406d06f52a03913cfa0db4d43f868b0ec7d0625a274c29efef1fe5a1afb357f42ac9e63849b1ede9e3540c4853e01c1b3d3590adc3356db290e53f482193735504"
	m.prevoteBytes, _ = hex.DecodeString(prevote)

	precommit := "b05b4f2c082210011abd01c1da2d940ab6010801100222480a207a81b9db16c3cb82bae88a058bb795620e4f2caca193a1141006bc82fc37efaf1224080c1220b0e26207ce3829dd522d5c1b6001975ca2ce17f1d3c3ab9e494407e1f5841b432a0c08fdcb8d940610b084fadd013214a881e35fac8a1ddcaa861e7d17b0122d956cb651380242409014eeafbdd9beef8afc73e35d00a690ebc592a36a4e1910a3fc245ae858fb0c2c404b039a14dcba0d4cc45a94d3eb845b847358cc21904e21e6e04e03d90f02"
	m.precommitBytes, _ = hex.DecodeString(precommit)

	hasVote := "b05b4f2c082010011a0a1919b3d5080318012002"
	m.hasVoteBytes, _ = hex.DecodeString(hasVote)

	newstep := "b05b4f2c082010011a0a1919b3d5080318012001"
	m.newStepBytes, _ = hex.DecodeString(newstep)

	blockpart := "b05b4f2c082110011afe0129e9ae8208031af501080112409aa51051d22a35521433c9c3373cef86e76dbaf3feb9fbeffe7fb6cfbecff30c20d0c3001a40a171e8d7de63d554b23994279243601006756afca80b8195aca91aae01080c10011a2028fcd60d68180f43dbf4fd7e4fb1b5005a500d15f598b1a05af3170adae9cf942220b7cf8d5a83aea7cc2d410b5358182e053010d232620c2269186714fc42613b012220c82dca7f5d7d21f14faef6faaa10e529ad09cd79a1edc5eb03a83f4dd9c3a08f222030ad77703cb8cd0616479f1b15a59a5a1870793de495b5a7fd5524733561c3d92220f41c34b3692a91b403273f4a29d54b68ea0da37c849c571cbee90d0d049a56c6"
	m.blockpartBytes, _ = hex.DecodeString(blockpart)

	tx := "2b06579d0aac01f8aa038405f5e100832dc6c0941033796b018b2bf0fc9cb88c0793b2f275edb62480b844a9059cbb00000000000000000000000083d83497431c2d3feab296a9fba4e5fadd2f7ed000000000000000000000000000000000000000000000152d02c7e14af680000081a9a0c5d31169294a8c1298fa4a178072b479666dc38fb994df02710ca165f91b6d14a047752e35109c52ac61496bf44ae545a2d5896a2e4bc59b10d9271748cc15f730"
	m.txBytes, _ = hex.DecodeString(tx)

	return m
}

func (mgr *MsgProcessMgr) start() {
	mgr.done = make(chan int8)
	mgr.finished = 0

	//mgr.sendMsgBytes(mgr.newStepBytes, 20*100)
	//mgr.sendMsgBytes(mgr.blockpartBytes, 13*100)
	//mgr.sendMsgBytes(mgr.prevoteBytes, 20*100)
	//mgr.sendMsgBytes(mgr.hasVoteBytes, 20*100)
	//mgr.sendMsgBytes(mgr.precommitBytes, 20*100)
	//mgr.sendMsgBytes(mgr.hasVoteBytes, 20*100)
	mgr.sendMsgBytes(mgr.txBytes, 200*100)

	//if mgr.decode {
		go mgr.writeWalRoutine()
	//}

	for i := 0; i < len(mgr.routineList); i++ {
		r := mgr.routineList[i]
		//r.start(113 * 100)
		r.start(200 * 100)
	}
}

func (mgr *MsgProcessMgr) endWork() {
	mgr.wal.Stop()
	mgr.wal.Wait()
}

func (mgr *MsgProcessMgr) writeWalRoutine() {
	for {
		select {
		case <-mgr.done:
			mgr.end <- 0
			return
		case msg := <-mgr.walCh:
			mgr.wal.Write(msg)
		}
	}
}

func (mgr *MsgProcessMgr) sendMsgBytes(msgBytes []byte, repeat int) {
	for i := 0; i < len(mgr.routineList); i++ {
		for j := 0; j < repeat; j++ {
			r := mgr.routineList[i]
			r.msgCh <- msgBytes
		}
	}
}

func (mgr *MsgProcessMgr) newMsgToWal(msg consensus.Message) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()
	mgr.walCh <- msg
}

func (mgr *MsgProcessMgr) routineFinished() {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()
	mgr.finished++
	//fmt.Println("finished: ", mgr.finished)
	if mgr.finished == len(mgr.routineList) {
		mgr.done <- 0
	}
}

func TestConsensusMsgProcessOnce(t *testing.T) {
	goroutinesList := [5]int{1, 5, 10, 15, 20}
	for i := 0; i < 5; i++ {
		fmt.Print(goroutinesList[i], " : ")
		mgr := newMsgProcessMgr(i, false)
		mgr.start()
		<-mgr.end
		mgr.endWork()
		fmt.Println()
	}
}

func TestConsensusMsgProcess5(t *testing.T) {
	boolList := [2]bool{true, false}
	goroutinesList := [5]int{1, 5, 10, 15, 20}
	for j := 0; j < len(boolList); j++ {
		for i := 0; i < 5; i++ {
			fmt.Print(goroutinesList[i], "-", j, " : ")
			mgr := newMsgProcessMgr(goroutinesList[i], boolList[j])
			mgr.start()
			<-mgr.end
			mgr.endWork()
			fmt.Println()
		}
	}
}

func BenchmarkConsensusMsgProcessOnce(b *testing.B) {
	mgr := newMsgProcessMgr(1, false)
	b.ResetTimer()
	start := time.Now()
	//for i := 0; i < b.N; i++ {
	mgr.start()
	<-mgr.end
	mgr.endWork()
	//}
	//avgDur := time.Since(start).Nanoseconds()/int64(b.N)
	//fmt.Println(avgDur)
	fmt.Println(time.Since(start).Nanoseconds())
}

func BenchmarkConsensusMsgProcess5(b *testing.B) {
	mgr := newMsgProcessMgr(5, true)
	b.ResetTimer()
	start := time.Now()
	//for i := 0; i < b.N; i++ {
	mgr.start()
	<-mgr.end
	mgr.endWork()
	//}
	//avgDur := time.Since(start).Nanoseconds()/int64(b.N)
	//fmt.Println(avgDur)
	fmt.Println(time.Since(start).Nanoseconds())
}

//func BenchmarkP2PMsgProcessOnce(b *testing.B) {
//	// wal
//	walDir, _ := os.Getwd()
//	walFile := filepath.Join(walDir, "wal")
//	wal, err := consensus.NewWAL(walFile,
//		autofile.GroupHeadSizeLimit(4096),
//		autofile.GroupCheckDuration(1*time.Millisecond),
//	)
//	if err != nil {
//		b.Fatal(err)
//	}
//	err = wal.Start()
//	defer func() {
//		wal.Stop()
//		wal.Wait()
//	}()
//
//	// prepare msg bytes
//	msginfo := "b05b4f2c082210011abb01c1da2d940ab4010802100222480a20f6c1c76201852e755fb4ba22f07ed027c335c48e6d8f7269d4e7f028a53f96d71224080d12205350e84037be21d7ffeb6eb9e541d10a2c5e906a34da005e88fa07989e5b89fe2a0c0896d488940610b08ef1d1013214006e126d9b97ccd450361317d8caa4c5d24e84cf4240209ab24a1d1d6dc6fd81e306b36b365c050e6fb9a1c583613976e8ad4e931d7dbe6a1ad4c31e7611ca5896f69125df10503e9f6fdeaf52ca50e086e4b64d3401"
//	msgBytes, err := hex.DecodeString(msginfo)
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	cdc := amino.NewCodec()
//	start := time.Now()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		//for k := 0; k < 28; k++ {
//		//	err = parseMsg(msgBytes, wal)
//		//	if err != nil {
//		//		b.Fatal(err)
//		//	}
//		peersCount := 1
//		var wg sync.WaitGroup
//		wg.Add(peersCount)
//		for index := 0; index < peersCount; index++ {
//			go func(i int, wg *sync.WaitGroup) {
//				msgpkt := conn.PacketMsg{}
//				err = msgpkt.UnmarshalFromAmino(cdc, msgBytes[4:])
//				if err != nil {
//					b.Fatal(err)
//				}
//
//				msg, err := consensus.decodeMsg(msgpkt.Bytes)
//				if err != nil {
//					b.Fatal(err)
//				}
//
//				if err = msg.ValidateBasic(); err != nil {
//					b.Fatal(err)
//				}
//
//				wal.Write(msg)
//				//err = parseMsg(msgBytes, wal)
//				//fmt.Println(i)
//				//if err != nil {
//				//	b.Fatal(err)
//				//}
//				wg.Done()
//			}(index, &wg)
//		}
//		wg.Wait()
//		//}
//	}
//	fmt.Printf("totalTime: %d, count: %d\n", time.Since(start).Nanoseconds(), b.N)
//
//	//for k := 1; k < 16; k++ {
//	//	peersCount := 15
//	//	var wg sync.WaitGroup
//	//	wg.Add(peersCount)
//	//	for index := 0; index < peersCount; index++ {
//	//		go func(i int, wg *sync.WaitGroup) {
//	//			err = parseMsg(msgBytes, wal)
//	//			//fmt.Println(i)
//	//			if err != nil {
//	//				b.Fatal(err)
//	//			}
//	//			wg.Done()
//	//		}(index, &wg)
//	//	}
//	//	wg.Wait()
//	//}
//}
//
//func BenchmarkP2pMsgProcessRepeatedly5(b *testing.B) {
//	// wal
//	walDir, _ := os.Getwd()
//	walFile := filepath.Join(walDir, "wal")
//	wal, err := consensus.NewWAL(walFile,
//		autofile.GroupHeadSizeLimit(4096),
//		autofile.GroupCheckDuration(1*time.Millisecond),
//	)
//	if err != nil {
//		b.Fatal(err)
//	}
//	err = wal.Start()
//	defer func() {
//		wal.Stop()
//		wal.Wait()
//	}()
//
//	// prepare msg bytes
//	msginfo := "b05b4f2c082210011abb01c1da2d940ab4010802100222480a20f6c1c76201852e755fb4ba22f07ed027c335c48e6d8f7269d4e7f028a53f96d71224080d12205350e84037be21d7ffeb6eb9e541d10a2c5e906a34da005e88fa07989e5b89fe2a0c0896d488940610b08ef1d1013214006e126d9b97ccd450361317d8caa4c5d24e84cf4240209ab24a1d1d6dc6fd81e306b36b365c050e6fb9a1c583613976e8ad4e931d7dbe6a1ad4c31e7611ca5896f69125df10503e9f6fdeaf52ca50e086e4b64d3401"
//	msgBytes, err := hex.DecodeString(msginfo)
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	cdc := amino.NewCodec()
//
//	b.ResetTimer()
//	start := time.Now()
//	for i := 0; i < b.N; i++ {
//		//for k := 0; k < 28; k++ {
//		peersCount := 5
//		var wg sync.WaitGroup
//		wg.Add(peersCount)
//		for index := 0; index < peersCount; index++ {
//			go func(i int, wg *sync.WaitGroup) {
//				msgpkt := conn.PacketMsg{}
//				err = msgpkt.UnmarshalFromAmino(cdc, msgBytes[4:])
//				if err != nil {
//					b.Fatal(err)
//				}
//
//				msg, err := consensus.decodeMsg(msgpkt.Bytes)
//				if err != nil {
//					b.Fatal(err)
//				}
//
//				if err = msg.ValidateBasic(); err != nil {
//					b.Fatal(err)
//				}
//
//				wal.Write(msg)
//				//err = parseMsg(msgBytes, wal)
//				//if err != nil {
//				//	b.Fatal(err)
//				//}
//				wg.Done()
//			}(index, &wg)
//		}
//		wg.Wait()
//		//}
//	}
//	fmt.Printf("totalTime: %d, count: %d\n", time.Since(start).Nanoseconds(), b.N)
//}
//
//func BenchmarkP2pMsgProcessRepeatedly10(b *testing.B) {
//	// wal
//	walDir, _ := os.Getwd()
//	walFile := filepath.Join(walDir, "wal")
//	wal, err := consensus.NewWAL(walFile,
//		autofile.GroupHeadSizeLimit(4096),
//		autofile.GroupCheckDuration(1*time.Millisecond),
//	)
//	if err != nil {
//		b.Fatal(err)
//	}
//	err = wal.Start()
//	defer func() {
//		wal.Stop()
//		wal.Wait()
//	}()
//
//	// prepare msg bytes
//	msginfo := "b05b4f2c082210011abb01c1da2d940ab4010802100222480a20f6c1c76201852e755fb4ba22f07ed027c335c48e6d8f7269d4e7f028a53f96d71224080d12205350e84037be21d7ffeb6eb9e541d10a2c5e906a34da005e88fa07989e5b89fe2a0c0896d488940610b08ef1d1013214006e126d9b97ccd450361317d8caa4c5d24e84cf4240209ab24a1d1d6dc6fd81e306b36b365c050e6fb9a1c583613976e8ad4e931d7dbe6a1ad4c31e7611ca5896f69125df10503e9f6fdeaf52ca50e086e4b64d3401"
//	msgBytes, err := hex.DecodeString(msginfo)
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	cdc := amino.NewCodec()
//
//	b.ResetTimer()
//	start := time.Now()
//	for i := 0; i < b.N; i++ {
//		//for k := 0; k < 28; k++ {
//		peersCount := 5
//		var wg sync.WaitGroup
//		wg.Add(peersCount)
//		for index := 0; index < peersCount; index++ {
//			go func(i int, wg *sync.WaitGroup) {
//				msgpkt := conn.PacketMsg{}
//				err = msgpkt.UnmarshalFromAmino(cdc, msgBytes[4:])
//				if err != nil {
//					b.Fatal(err)
//				}
//
//				msg, err := consensus.decodeMsg(msgpkt.Bytes)
//				if err != nil {
//					b.Fatal(err)
//				}
//
//				if err = msg.ValidateBasic(); err != nil {
//					b.Fatal(err)
//				}
//
//				wal.Write(msg)
//				//err = parseMsg(msgBytes, wal)
//				//if err != nil {
//				//	b.Fatal(err)
//				//}
//				wg.Done()
//			}(index, &wg)
//		}
//		wg.Wait()
//		//}
//	}
//	fmt.Printf("totalTime: %d, count: %d\n", time.Since(start).Nanoseconds(), b.N)
//}
//
//func parseMsg(bz []byte, wal *consensus.BaseWAL) error {
//	cdc := amino.NewCodec()
//	msgpkt := conn.PacketMsg{}
//	err := msgpkt.UnmarshalFromAmino(cdc, bz[4:])
//	if err != nil {
//		return err
//	}
//
//	msg, err := consensus.decodeMsg(msgpkt.Bytes)
//	if err != nil {
//		return err
//	}
//
//	if err = msg.ValidateBasic(); err != nil {
//		return err
//	}
//
//	wal.Write(msg)
//	return nil
//}

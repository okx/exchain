package distrlock

/*
import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/log"
)

func task1(task *string, exit chan struct{}) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	lock, err := ParseEtcdLock(logger)
	if err != nil {
		fmt.Println("ParseEtcdLock failed!", err.Error())
		return
	}

	fmt.Println("worker0 try get lock mysql:", 0, "...")
	lock.TryLockBlock(fmt.Sprintf("mysql:%d", 0), "worker0")
	fmt.Println("worker0 lock mysql:", 0, "succeed!")
	*task = "task1"
	lock.Client.Close()
	<-exit
	//lock.UnLock(fmt.Sprintf("mysql:%d", i), "worker0")
	//fmt.Println("worker0  unlock mysql:", i, "succeed!")
}

func task2(task *string, exit chan struct{}) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	lock, err := ParseEtcdLock(logger)
	if err != nil {
		fmt.Println("ParseEtcdLock failed!", err.Error())
		return
	}

	fmt.Println("worker1 try get lock mysql:", 0, "...")
	lock.TryLockBlock(fmt.Sprintf("mysql:%d", 0), "worker1")
	fmt.Println("worker1 lock mysql:", 0, "succeed!")
	*task = "task2"
	lock.UnLock(fmt.Sprintf("mysql:%d", 0), "worker1")
	fmt.Println("worker1  unlock mysql:", 0, "succeed!")
}

// test a process can lock when another is killed
func TestEtcdLock_TryLockBlock(t *testing.T) {
	return
	viper.Reset()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	flagSet.String(FlagEtcdLock, "127.0.0.1:2379", "")
	viper.BindPFlags(flagSet)
	exitTask1 := make(chan struct{}, 1)
	exitTask2 := make(chan struct{}, 1)

	var task = ""
	go task1(&task, exitTask1)
	time.Sleep(1 * time.Second)
	require.Equal(t, "task1", task)
	time.Sleep(1 * time.Second)
	go task2(&task, exitTask2)
	exitTask1 <- struct{}{}
	time.Sleep(3 * time.Second)
	require.Equal(t, "task2", task)
}

func task3(task *string, info chan struct{}) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	lock, err := ParseEtcdLock(logger)
	if err != nil {
		fmt.Println("ParseEtcdLock failed!", err.Error())
		return
	}

	fmt.Println("worker0 try get lock mysql:", 0, "...")
	lock.TryLockBlock(fmt.Sprintf("mysql:%d", 0), "worker0")
	fmt.Println("worker0 lock mysql:", 0, "succeed!")
	*task = "task3"
	time.Sleep(2 * time.Second)
	lock.UnLock(fmt.Sprintf("mysql:%d", 0), "worker0")
	time.Sleep(100 * time.Millisecond)
	info <- struct{}{}
	fmt.Println("worker0  unlock mysql:", 0, "succeed!")
}

func task4(task *string) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	lock, err := ParseEtcdLock(logger)
	if err != nil {
		fmt.Println("ParseEtcdLock failed!", err.Error())
		return
	}

	fmt.Println("worker1 try get lock mysql:", 0, "...")
	lock.TryLockBlock(fmt.Sprintf("mysql:%d", 0), "worker1")
	fmt.Println("worker1 lock mysql:", 0, "succeed!")
	*task = "task4"
	lock.UnLock(fmt.Sprintf("mysql:%d", 0), "worker1")
	fmt.Println("worker1  unlock mysql:", 0, "succeed!")
}

// test a process can lock when another unlock
func TestEtcdLock_TryLockBlock1(t *testing.T) {
	return

	viper.Reset()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	flagSet.String(FlagEtcdLock, "127.0.0.1:2379", "")
	viper.BindPFlags(flagSet)

	var task = ""
	info := make(chan struct{}, 1)
	go task3(&task, info)
	time.Sleep(1 * time.Second)
	require.Equal(t, "task3", task)
	go task4(&task)
	<-info
	require.Equal(t, "task4", task)
}
*/
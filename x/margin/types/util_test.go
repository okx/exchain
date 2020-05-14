package types

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetMarginAccountByHash(t *testing.T) {

	//var initString = "okchain10q0rk5qnyag7wfvvt7rtphlw589m7frsmyq4ya"

	readBT := time.Now()
	addressList := readLine("/Users/fpchen/Desktop/TestProgram/Python/address")
	eT := time.Since(readBT)
	fmt.Println("read file time: ", eT)

	fmt.Println("addressList size : ", len(addressList))

	f, _ := os.OpenFile("/Users/fpchen/Desktop/TestProgram/address", os.O_WRONLY|os.O_APPEND, 0666)
	writeHash := time.Now()
	for _, addr := range addressList {
		addrHash := GetMarginAccount(addr)
		f.WriteString(addrHash.String() + "\n")
	}

	hashET := time.Since(writeHash)
	fmt.Println("write hash file time: ", hashET)

	f.WriteString("====================================================================================================\n")

	writeSwap := time.Now()
	for _, addr := range addressList {
		addrSwap, _ := GetMarginAccountBySwap(addr)
		f.WriteString(addrSwap.String() + "\n")
	}
	swapET := time.Since(writeSwap)

	fmt.Println("write swap file time: ", swapET)

	//
	//addrSwap, err := GetMarginAccountBySwap(initString)
	//fmt.Println("hash : ", addrHash.String())
	//fmt.Println(err)
	//fmt.Println("swap : ", addrSwap.String())

}

func readLine(filename string) (addressList []string) {

	r, _ := os.Open(filename)
	defer r.Close()
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		line = strings.Replace(line, "\t", "", -1)
		addressList = append(addressList, line)
	}
	return
}

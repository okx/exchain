package backend

import (
	"fmt"
	"testing"
)

func TestChan(t *testing.T)  {
	c := make(chan interface{}, 10)
	c <- "1"
	c <- "2"

	for true {
		if len(c) > 0 {
			v, ok := <- c
			if ok {
				fmt.Println(v)
			} else {
				break
			}
		} else {
			break
		}
	}
}

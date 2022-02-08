package logevents

import (
	"fmt"
	"os"
	"testing"
)

func TestProvider(t *testing.T) {

}

func TestSubscriber(t *testing.T) {
	s := &subscriber{
		fileMap: make(map[string]*os.File),
	}

	for i := 1; i < 100; i++ {
		s.onEvent(fmt.Sprintf("192.168.0.%d.log", i%6), "test\n")
	}
}

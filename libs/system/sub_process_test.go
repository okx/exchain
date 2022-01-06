package system_test

import (
	"bytes"
	"fmt"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func testTmp1(t *testing.T) {
	viper.Set(tmtypes.FlagDownloadDDS, true)
	v := tmtypes.EnableDownloadDelta()
	assert.True(t, v, "tmp1")
}
func testTmp2(t *testing.T) {
	viper.Set(tmtypes.FlagDownloadDDS, false)
	v := tmtypes.EnableDownloadDelta()
	assert.False(t, v, "tmp2")
}
func TestCommitDelta(t *testing.T) {
	var funcs = []func(t *testing.T) {
		testTmp1,
		testTmp2,
	}
	for i, f := range funcs {
		if os.Getenv("SUB_PROCESS") == fmt.Sprintf("%d", i) {
			f(t)
			return
		}
	}

	for i, _ := range funcs {
		var outb, errb bytes.Buffer
		cmd := exec.Command(os.Args[0], "-test.run=TestCommitDelta")
		cmd.Env = append(os.Environ(), fmt.Sprintf("SUB_PROCESS=%d", i))
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			isFailed := false
			if strings.Contains(outb.String(), "FAIL:") ||
				strings.Contains(errb.String(), "FAIL:") {
				fmt.Print(cmd.Stderr)
				fmt.Print(cmd.Stdout)
				isFailed = true
			}
			assert.Equal(t, isFailed, false)
		}
	}

}

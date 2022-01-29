package sanity

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type item interface {
	// label: get item's name
	label() string
	// check: whether the userSetting value is equal to the conflicts value
	check() bool
	// verbose: show the readable flag
	verbose() string
}

type boolItem struct {
	name  string
	value bool
}

func (b boolItem) label() string {
	return b.name
}

func (b boolItem) check() bool {
	return viper.GetBool(b.label()) == b.value
}

func (b boolItem) verbose() string {
	return fmt.Sprintf("--%v=%v", b.name, b.value)
}

type stringItem struct {
	name  string
	value string
}

func (s stringItem) label() string {
	return s.name
}

func (s stringItem) check() bool {
	return strings.ToLower(viper.GetString(s.label())) == s.value
}

func (s stringItem) verbose() string {
	return fmt.Sprintf("--%v=%v", s.name, s.value)
}

package sanity

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	FlagDisableSanity = "disable-sanity"
)

// item: app's flags
type item interface {
	// label: get item's name
	label() string
	// check: whether the userSetting value is equal to the value
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

type intItem struct {
	name  string
	value int
}

func (i intItem) label() string {
	return i.name
}

func (i intItem) check() bool {
	return viper.GetInt(i.label()) == i.value
}

func (i intItem) verbose() string {
	return fmt.Sprintf("--%v=%v", i.name, i.value)
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

type dependentPair struct {
	config       item
	reliedConfig item
}

func (cp *dependentPair) Check() error {
	//if config is true,  then the reliedConfig must be checked as true
	if cp.config.check() &&
		!cp.reliedConfig.check() {
		return fmt.Errorf(" %v must be set explicitly, as %v", cp.reliedConfig.verbose(), cp.config.verbose())
	}
	return nil
}

// conflictPair: configA and configB are conflict pair
type conflictPair struct {
	configA item
	configB item
}

// checkConflict: check configA vs configB
// and the value is equal to the conflicts value then complain it
func (cp *conflictPair) Check() error {
	if cp.configA.check() &&
		cp.configB.check() {
		return fmt.Errorf(" %v conflict with %v", cp.configA.verbose(), cp.configB.verbose())
	}

	return nil
}

type rangeItem struct {
	enumRange []int
	value     int
	name      string
}

func (i rangeItem) label() string {
	return i.name
}

func (i rangeItem) checkRange() error {
	i.value = viper.GetInt(i.label())

	for _, v := range i.enumRange {
		if v == i.value {
			return nil
		}
	}

	return fmt.Errorf(" %v", i.verbose())
}

func (b rangeItem) verbose() string {
	return fmt.Sprintf("--%v=%v not in %v", b.name, b.value, b.enumRange)
}

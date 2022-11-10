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
	// check: whether the userSetting value is equal to expect value
	check() bool
	// verbose: show the readable flag
	verbose() string
}

type boolItem struct {
	name   string
	expect bool
}

func (b boolItem) label() string {
	return b.name
}

func (b boolItem) check() bool {
	return viper.GetBool(b.label()) == b.expect
}

func (b boolItem) verbose() string {
	return fmt.Sprintf("--%v=%v", b.name, b.expect)
}

type stringItem struct {
	name   string
	expect string
}

func (s stringItem) label() string {
	return s.name
}

func (s stringItem) check() bool {
	return strings.ToLower(viper.GetString(s.label())) == s.expect
}

func (s stringItem) verbose() string {
	return fmt.Sprintf("--%v=%v", s.name, s.expect)
}

type funcItem struct {
	name   string
	expect bool
	actual bool
	f      func() bool
}

func (f funcItem) label() string {
	return f.name
}

func (f funcItem) check() bool {
	f.actual = f.f()
	return f.actual == f.expect
}

func (f funcItem) verbose() string {
	return fmt.Sprintf("%v=%v", f.name, f.actual)
}

type dependentPair struct {
	config       item
	reliedConfig item
}

func (cp *dependentPair) check() error {
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
	tips    string
}

// checkConflict: check configA vs configB
// if both configA and configB are got expect values
// then complain it. if there is a custom tips use it.
func (cp *conflictPair) check() error {
	if cp.configA.check() &&
		cp.configB.check() {
		if cp.tips == "" {
			return fmt.Errorf(" %v conflict with %v", cp.configA.verbose(), cp.configB.verbose())
		}
		return fmt.Errorf(cp.tips)
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

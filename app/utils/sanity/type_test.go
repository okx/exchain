package sanity

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"testing"
)

// universeFlag used to build command
type universeFlag interface {
	// add flag to cmd
	add(cmd *cobra.Command)
	// args get formatted flags
	args() string
	// changed If the user set the value (or if left to default)
	changed() bool
}

// boolFlag bool type flag
type boolFlag struct {
	Name    string
	Default bool
	Changed bool
	Value   bool
}

func (bf *boolFlag) add(cmd *cobra.Command) {
	cmd.Flags().Bool(bf.Name, bf.Default, "")
	viper.BindPFlag(bf.Name, cmd.Flags().Lookup(bf.Name))
}

func (bf *boolFlag) args() string {
	return fmt.Sprintf("--%v=%v", bf.Name, bf.Value)
}

func (bf *boolFlag) changed() bool {
	return bf.Changed
}

// stringFlag string type flag
type stringFlag struct {
	Name    string
	Default string
	Changed bool
	Value   string
}

func (sf *stringFlag) add(cmd *cobra.Command) {
	cmd.Flags().String(sf.Name, sf.Default, "")
	viper.BindPFlag(sf.Name, cmd.Flags().Lookup(sf.Name))
}

func (sf *stringFlag) args() string {
	return fmt.Sprintf("--%v=%v", sf.Name, sf.Value)
}

func (sf *stringFlag) changed() bool {
	return sf.Changed
}

// getCommand build command by flags
func getCommand(flags []universeFlag) *cobra.Command {
	cmd := &cobra.Command{}
	var args []string
	for _, v := range flags {
		v.add(cmd)
		if v.changed() {
			args = append(args, v.args())
		}
	}
	cmd.ParseFlags(args)

	cmd.Execute()
	return cmd
}

func getCommandBool() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    "b1",
			Default: false,
			Changed: true,
			Value:   true,
		},
		&boolFlag{
			Name:    "b2",
			Default: false,
			Changed: true,
			Value:   true,
		},
	})
}

func getCommandBoolDiff() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    "b1",
			Default: false,
			Changed: true,
			Value:   true,
		},
		&boolFlag{
			Name:    "b3",
			Default: false,
			Changed: true,
			Value:   false,
		},
	})
}

func getCommandBoolString() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    "b1",
			Default: false,
			Changed: true,
			Value:   true,
		},
		&stringFlag{
			Name:    "s1",
			Default: "none",
			Changed: true,
			Value:   "conflict",
		},
	})
}

func Test_conflictPair_checkConflict(t *testing.T) {
	type fields struct {
		configA item
		configB item
	}
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "1. bool item and bool item both true",
			fields: fields{configA: boolItem{name: "b1", value: true}, configB: boolItem{name: "b2", value: true}},
			args:   args{cmd: getCommandBool()}, wantErr: true},
		{name: "2. bool item and bool item true vs false",
			fields: fields{configA: boolItem{name: "b1", value: true}, configB: boolItem{name: "b3", value: false}},
			args:   args{cmd: getCommandBoolDiff()}, wantErr: true},
		{name: "3. bool item and string item",
			fields: fields{configA: boolItem{name: "b1", value: true}, configB: stringItem{name: "s1", value: "conflict"}},
			args:   args{cmd: getCommandBoolString()}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &conflictPair{
				configA: tt.fields.configA,
				configB: tt.fields.configB,
			}
			var err error
			if err = cp.checkConflict(); (err != nil) != tt.wantErr {
				t.Errorf("checkConflict() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(err)
		})
	}
}

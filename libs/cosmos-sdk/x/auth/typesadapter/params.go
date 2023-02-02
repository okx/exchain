package types

import (
	"sigs.k8s.io/yaml"
)

// String implements the stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

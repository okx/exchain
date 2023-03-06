package localhost

import "github.com/okx/okbchain/libs/ibc-go/modules/light-clients/09-localhost/types"

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}

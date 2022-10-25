package rootmulti

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type storeSmbInfos []storeInfo

func (s storeSmbInfos) String() string {
	sb := strings.Builder{}
	for _, v := range s {
		sb.WriteString(fmt.Sprintf("%s:%v\n", v.Name, hex.EncodeToString(v.Core.CommitID.Hash)))
	}

	return sb.String()
}

func (s storeSmbInfos) Len() int {
	return len(s)
}

func (s storeSmbInfos) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s storeSmbInfos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
